package game

import "fmt"

func registerEventsE100() {
	// e102-e108: Airborne weather events (if pegasus/griffon travel)
	RegisterEvent("e102", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"Fierce winds buffet your flying mount! -1 food lost to the effort."}, FoodChange: -1}
	})
	RegisterEvent("e103", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"A lightning storm forces you to land early. Travel interrupted."}}
	})
	RegisterEvent("e104", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"Clear skies and favorable winds. You make excellent progress!"}}
	})
	RegisterEvent("e105", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"Thick clouds obscure your vision. You become uncertain of your position."}}
	})
	RegisterEvent("e106", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"Ice crystals form on your mount's wings. Travel is slowed."}}
	})
	RegisterEvent("e107", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6() * 15
		return EventResult{
			Messages:   []string{fmt.Sprintf("From above you spot a glinting cache below! You land to recover %d gold.", gold)},
			GoldChange: gold,
		}
	})
	RegisterEvent("e108", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"Your mount spots a thermal and rides it gracefully. A swift journey."}}
	})

	// e109 - Wild Pegasus
	RegisterEvent("e109", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A beautiful winged horse lands nearby, eyeing you curiously!"},
			Choices:  []string{"Attempt to tame it", "Leave it be"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if Roll2d6()+s.Prince.WitWiles >= 10 {
						s.Prince.HasMount = true
						s.Prince.MountType = MountPegasus
						s.Flags.PegasusCaptured = true
						return EventResult{Messages: []string{"You gentle the pegasus with patience and skill. A mighty mount is yours!"}}
					}
					wounds := Roll1d3()
					s.Prince.Wounds += wounds
					return EventResult{Messages: []string{fmt.Sprintf("The pegasus bucks and kicks you for %d wounds before fleeing!", wounds)}}
				}
				return EventResult{Messages: []string{"You watch the magnificent creature take flight."}}
			},
		}
	})

	// e120 - Desert Exhaustion
	RegisterEvent("e120", func(s *GameState, ctx EventContext) EventResult {
		s.Prince.Wounds++
		return EventResult{
			Messages: []string{"The merciless desert heat saps your strength.", "You suffer from exhaustion. +1 wound."},
		}
	})

	// e121 - Sunstroke
	RegisterEvent("e121", func(s *GameState, ctx EventContext) EventResult {
		wounds := Roll1d3() + 1
		s.Prince.Wounds += wounds
		return EventResult{
			Messages: []string{fmt.Sprintf("Terrible sunstroke! The sun beats down mercilessly. +%d wounds!", wounds),
				"Without shelter or rest, you may not survive the desert."},
		}
	})

	// e128 - Nothing (placeholder)
	RegisterEvent("e128", func(s *GameState, ctx EventContext) EventResult {
		messages := []string{
			"The road is quiet today.",
			"Nothing of note occurs.",
			"The day passes uneventfully.",
			"All is calm.",
			"You travel without incident.",
			"The countryside is peaceful.",
		}
		return EventResult{Messages: []string{messages[Roll1d6()-1]}}
	})

	// e129 - Small group
	RegisterEvent("e129", func(s *GameState, ctx EventContext) EventResult {
		roll := Roll1d6()
		if roll <= 3 {
			return EventResult{
				Messages: []string{"You encounter a small band of travelers. They share their campfire and news."},
				FoodChange: 1,
			}
		}
		return EventResult{
			Messages: []string{"A small group of refugees passes. They have little to share."},
		}
	})

	// Ruins events e131-e138
	RegisterEvent("e131", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(4, Roll1d6())
		s.GetHexFlags(s.CurrentHex).Searched = true
		return EventResult{
			Messages:   []string{fmt.Sprintf("In the ruins you find a cache of ancient coins! %d gold!", gold)},
			GoldChange: gold,
		}
	})

	RegisterEvent("e132", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		enemy := MakeEnemy("Ruins Undead", 4, 10, 4)
		enemy.IsUndead = true
		return EventResult{
			Messages:        []string{"Undead guardians rise from the rubble!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 4,
		}
	})

	RegisterEvent("e133", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(6, Roll1d6())
		item := SpecialPossessionRoll(5)
		s.GetHexFlags(s.CurrentHex).Searched = true
		result := EventResult{
			Messages:   []string{fmt.Sprintf("You find a richly appointed chamber with %d gold!", gold)},
			GoldChange: gold,
		}
		if item != PossNone {
			s.Prince.AddPossession(item)
			result.Messages = append(result.Messages, fmt.Sprintf("You also find: %s!", PossessionName(item)))
		}
		return result
	})

	RegisterEvent("e134", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		if Roll1d6() >= 4 {
			s.Prince.PoisonWounds += 2
			return EventResult{
				Messages: []string{"You trigger an ancient trap! Poison darts strike you! +2 poison wounds."},
			}
		}
		return EventResult{
			Messages: []string{"You find and disarm an ancient trap. The chamber beyond is empty."},
		}
	})

	RegisterEvent("e135", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		// Big monster guard big treasure
		enemy := MakeEnemy("Ancient Golem", 7, 16, 7)
		enemy.IsUndead = true
		gold := TreasureRoll(7, Roll1d6())
		return EventResult{
			Messages:        []string{fmt.Sprintf("A massive stone golem guards a treasure of %d gold!", gold)},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  false,
			GoldChange:      gold,
		}
	})

	RegisterEvent("e136", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		return EventResult{
			Messages: []string{"The ruins are thoroughly picked over. You find nothing of value."},
		}
	})

	RegisterEvent("e137", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		gold := TreasureRoll(5, Roll1d6())
		enemy := MakeEnemy("Ruins Bandits", 4, 12, 4)
		return EventResult{
			Messages:        []string{fmt.Sprintf("Other treasure hunters are here! They have found %d gold and won't share!", gold)},
			CombatTriggered: true,
			Enemy:           &enemy,
			GoldChange:      gold,
		}
	})

	RegisterEvent("e138", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		// Royal Helm location - but only if not already found
		if !s.Flags.HasRoyalHelm {
			s.Prince.AddPossession(PossRoyalHelm)
			s.Flags.HasRoyalHelm = true
			return EventResult{
				Messages: []string{"Deep in the ruins, on a stone pedestal, gleams the ROYAL HELM!",
					"This is the lost crown of your dynasty! Your heart surges with hope!",
					"Recover this to Ogon or Weshor and your throne is restored!"},
			}
		}
		gold := TreasureRoll(6, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find ancient treasury records and %d gold in coin.", gold)},
			GoldChange: gold,
		}
	})

	// Secrets e143-e147
	RegisterEvent("e143", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You hear whispers of a secret cache in the old fortress ruins.",
				"\"Go to the tower where three roads meet, and dig beneath the oak.\""},
			Note: fmt.Sprintf("Secret cache rumoured near fortress ruins (day %d, hex %s).", s.Day, s.CurrentHex),
		}
	})

	RegisterEvent("e144", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A dying noble whispers that the Baron's heir lives in the south!",
				"\"Find him and he will support your claim to the throne.\""},
			Note: fmt.Sprintf("Baron's heir said to live in the south (day %d, hex %s).", s.Day, s.CurrentHex),
		}
	})

	RegisterEvent("e145", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"Ancient texts reveal: the Staff of Command lies in the deepest desert ruin.",
				"\"He who wields the Staff shall command the armies of the east.\""},
			Note: fmt.Sprintf("Staff of Command: deepest desert ruin (day %d, hex %s).", s.Day, s.CurrentHex),
		}
	})

	RegisterEvent("e146", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(7, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("A hidden treasure room! Ancient gold worth %d pieces!", gold)},
			GoldChange: gold,
			Note:       fmt.Sprintf("Hidden treasure room at hex %s yielded %d gold (day %d).", s.CurrentHex, gold, s.Day),
		}
	})

	RegisterEvent("e147", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You discover the location of a secret passage beneath Huldra Castle.",
				"It leads directly to the throne room."},
			Note: fmt.Sprintf("Secret passage beneath Huldra Castle (day %d, hex %s).", s.Day, s.CurrentHex),
		}
	})

	// Audience events e148-e161
	RegisterEvent("e148", func(s *GameState, ctx EventContext) EventResult {
		s.AudienceBarred[s.CurrentHex] = s.Day + 7
		return EventResult{
			Messages: []string{"The lord is in a foul mood and bars you from his court for a week!"},
		}
	})

	RegisterEvent("e149", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6() * 20
		return EventResult{
			Messages:   []string{fmt.Sprintf("The lord is sympathetic to your cause and grants you %d gold for your quest.", gold)},
			GoldChange: gold,
		}
	})

	RegisterEvent("e150", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"The lord offers you lodging and news in exchange for news from afar.",
				"You spend a day gathering intelligence about the region."},
			FoodChange: Roll1d6(),
		}
	})

	RegisterEvent("e151", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6()*30 + 20
		enemy := MakeEnemy("Rival Champion", 6, 12, 0) // WealthCode=0: no random loot; prize awarded on win
		return EventResult{
			Messages: []string{fmt.Sprintf("The lord offers %d gold if you defeat his champion in single combat!", gold)},
			Choices:  []string{"Accept the challenge", "Decline"},
			ChoiceHandler: func(gs *GameState, choice int) EventResult {
				if choice != 0 {
					return EventResult{Messages: []string{"You decline the challenge. The lord dismisses you."}}
				}
				// Store the prize; it is awarded in handleCombatKey only if the player wins
				gs.PendingDuelGold = gold
				e := enemy
				return EventResult{
					Messages:        []string{"You step onto the field of honor!"},
					CombatTriggered: true,
					Enemy:           &e,
					PlayerAttFirst:  true,
				}
			},
		}
	})

	RegisterEvent("e152", func(s *GameState, ctx EventContext) EventResult {
		if s.Flags.NobleAllySecured {
			gold := Roll1d6() * 25
			return EventResult{
				Messages:   []string{fmt.Sprintf("Your noble ally grants you additional support: %d gold.", gold)},
				GoldChange: gold,
			}
		}
		s.Flags.NobleAllySecured = true
		s.Prince.AddPossession(PossNobleParchment)
		return EventResult{
			Messages: []string{
				"The lord is impressed by your bearing and deeds!",
				"He pledges his house's support to your cause — you have a noble ally!",
				"The Noble Parchment is sealed. Return north and your throne awaits!",
			},
		}
	})

	RegisterEvent("e153", func(s *GameState, ctx EventContext) EventResult {
		follower := Character{
			Name:         "Knight Errant",
			Type:         TypeMercenary,
			CombatSkill:  6,
			MaxEndurance: 12,
			DailyWage:    8,
			Morale:       6,
			HasMount:     true,
			MountType:    MountHorse,
		}
		return EventResult{
			Messages: []string{"The lord offers you one of his knights as escort for your mission."},
			Choices:  []string{"Accept the knight's service (8 gold/day)", "Respectfully decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					return EventResult{NewFollower: &follower}
				}
				return EventResult{Messages: []string{"You thank the lord and decline the escort."}}
			},
		}
	})

	// Temple offering events e154-e159
	RegisterEvent("e154", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"The priests reject your offering as insufficient. They bar you from the temple."},
		}
	})

	RegisterEvent("e155", func(s *GameState, ctx EventContext) EventResult {
		s.Prince.Wounds = 0
		return EventResult{
			Messages: []string{
				"The gods accept your offering.",
				"The priests bless your journey — all wounds are healed by divine grace!",
			},
		}
	})

	RegisterEvent("e156", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d6() + 4
		return EventResult{
			Messages:   []string{fmt.Sprintf("The temple community shares their surplus with you: %d food units.", food)},
			FoodChange: food,
		}
	})

	RegisterEvent("e157", func(s *GameState, ctx EventContext) EventResult {
		if !s.Flags.HasRoyalHelm {
			return EventResult{
				Messages: []string{"The oracle speaks: \"The Royal Helm lies in the ruins of the old kingdom.\"",
					"\"It waits in the deepest chamber, where stone guardians sleep.\""},
			}
		}
		gold := Roll1d6() * 30
		return EventResult{
			Messages:   []string{fmt.Sprintf("The temple treasury contributes %d gold to your noble cause!", gold)},
			GoldChange: gold,
		}
	})

	RegisterEvent("e158", func(s *GameState, ctx EventContext) EventResult {
		s.Prince.PoisonWounds = 0
		return EventResult{
			Messages: []string{
				"The priests perform a purification ritual, cleansing poison from your blood.",
				"All poison is purged from your body!",
			},
		}
	})

	RegisterEvent("e159", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossStaffOfCommand) && Roll1d6() == 6 {
			s.Prince.AddPossession(PossStaffOfCommand)
			s.Flags.HasStaffOfCommand = true
			return EventResult{
				Messages: []string{"The High Priest emerges bearing an ancient staff.",
					"\"This is the Staff of Command,\" he says. \"Use it to reclaim what is yours.\"",
					"You have received the STAFF OF COMMAND! Return north to win the game!"},
			}
		}
		gold := Roll1d6() * 20
		return EventResult{
			Messages:   []string{fmt.Sprintf("The gods reward your faith with %d gold.", gold)},
			GoldChange: gold,
		}
	})

	RegisterEvent("e160", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"The temple offers sanctuary. You may rest here safely for one day.",
				"All wounds are healed at a cost of 20 gold."},
			Choices: []string{"Pay 20 gold to heal", "Rest without healing (free)"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold < 20 {
						return EventResult{Messages: []string{"You cannot afford the 20 gold healing. You rest in the temple's guest quarters instead."}}
					}
					s.Gold -= 20
					s.Prince.Wounds = 0
					s.Prince.PoisonWounds = 0
					return EventResult{Messages: []string{"Temple healers restore your health completely."}}
				}
				return EventResult{Messages: []string{"You rest in the temple's guest quarters."}}
			},
		}
	})

	RegisterEvent("e161", func(s *GameState, ctx EventContext) EventResult {
		// Count Drogat audience
		if s.CurrentHex == NewHexID(10, 18) || s.CurrentHex == NewHexID(3, 23) {
			if !s.Flags.NobleAllySecured {
				s.Flags.NobleAllySecured = true
				s.Prince.AddPossession(PossNobleParchment)
				return EventResult{
					Messages: []string{
						"You stand before Count Drogat himself!",
						"The old count listens to your tale and strokes his grey beard.",
						"\"Your father was a good man,\" he says at last. \"You have my support.\"",
						"Count Drogat pledges his forces and grants you 200 gold!",
						"Return to the north with this alliance and your throne is restored!",
					},
					GoldChange: 200,
				}
			}
		}
		gold := Roll1d6() * 30
		return EventResult{
			Messages:   []string{fmt.Sprintf("A generous lord grants you %d gold in support of your quest.", gold)},
			GoldChange: gold,
		}
	})
}
