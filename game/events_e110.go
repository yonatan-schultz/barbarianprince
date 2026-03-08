package game

import "fmt"

func registerEventsE110() {
	// ── Airborne encounters e110-e119 (pegasus/griffon travel) ───────────────

	// e110 - Sky Pirates
	RegisterEvent("e110", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d3() + 1
		enemy := MakeEnemy(fmt.Sprintf("Sky Pirates (%d)", count), count+3, count*4+4, 4)
		return EventResult{
			Messages:        []string{fmt.Sprintf("%d sky pirates on winged steeds attack from out of the sun!", count)},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 5,
		}
	})

	// e111 - Wyvern
	RegisterEvent("e111", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("Wyvern", 8, 20, 6)
		return EventResult{
			Messages:        []string{"A wyvern rises from a mountain peak, drawn by your flying mount!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  false,
		}
	})

	// e112 - Thunderstorm (airborne)
	RegisterEvent("e112", func(s *GameState, ctx EventContext) EventResult {
		wounds := Roll1d3()
		s.Prince.Wounds += wounds
		return EventResult{
			Messages: []string{
				"A violent thunderstorm strikes without warning!",
				fmt.Sprintf("Lightning and hail batter you and your mount. %d wounds from the storm!", wounds),
				"You are forced to land early and take shelter.",
			},
			FoodChange: -1,
		}
	})

	// e113 - Enemy camp spotted from above
	RegisterEvent("e113", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"From above you spot an enemy war camp spread across the valley below."},
			Choices:  []string{"Swoop down and raid the supply wagon", "Continue past unseen"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if Roll1d6() >= 4 {
						gold := TreasureRoll(5, Roll1d6())
						return EventResult{
							Messages:   []string{fmt.Sprintf("A swift raid nets you %d gold from the supply wagon!", gold)},
							GoldChange: gold,
						}
					}
					enemy := MakeEnemy("Camp Guards", 4, 14, 3)
					return EventResult{
						Messages:        []string{"The guards spot you! They loose arrows and close to fight!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				}
				return EventResult{Messages: []string{"You fly past the camp, unseen in the clouds."}}
			},
		}
	})

	// e114 - Favorable winds
	RegisterEvent("e114", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d3()
		return EventResult{
			Messages:   []string{"Powerful air currents carry you and your mount swiftly forward.", fmt.Sprintf("The easy flight means your mount needs less effort — you conserve %d food.", food)},
			FoodChange: food,
		}
	})

	// e115 - Cloud Giant
	RegisterEvent("e115", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A massive cloud giant steps out of the mist, hurling boulders at you!"},
			Choices:  []string{"Fight the giant!", "Outmaneuver it (Wit/Wiles check)"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					enemy := MakeEnemy("Cloud Giant", 8, 24, 7)
					return EventResult{
						Messages:        []string{"You close with the cloud giant sword in hand!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				}
				if Roll1d6()+s.Prince.WitWiles >= 9 {
					return EventResult{Messages: []string{"You dart through the clouds and outmaneuver the lumbering giant!"}}
				}
				wounds := Roll1d3() + 1
				s.Prince.Wounds += wounds
				return EventResult{
					Messages: []string{fmt.Sprintf("A thrown boulder clips you! %d wounds.", wounds)},
				}
			},
		}
	})

	// e116 - Flying Drake
	RegisterEvent("e116", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("Flying Drake", 7, 16, 5)
		return EventResult{
			Messages:        []string{"A flying drake rivals your mount's territory and attacks with tooth and claw!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 4,
		}
	})

	// e117 - Giant bird flock
	RegisterEvent("e117", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 4 {
			return EventResult{
				Messages: []string{"A vast flock of giant birds engulfs you!", "Your mount is disoriented and you lose your bearings."},
				FoodChange: -1,
			}
		}
		return EventResult{
			Messages: []string{"A flock of giant birds wheels past you in perfect formation — a breathtaking sight."},
		}
	})

	// e118 - Ruins spotted from air
	RegisterEvent("e118", func(s *GameState, ctx EventContext) EventResult {
		for _, adj := range AdjacentHexes(s.CurrentHex) {
			h := GetHex(adj)
			if h != nil && h.IsRuins() && !s.VisitedHexes[adj] {
				s.VisitedHexes[adj] = true
				return EventResult{
					Messages: []string{fmt.Sprintf("From the air you spot the unmistakable outline of ancient ruins at hex %s!", adj)},
					Note:     fmt.Sprintf("Ruins spotted from the air at hex %s (day %d).", adj, s.Day),
				}
			}
		}
		gold := Roll1d6() * 10
		return EventResult{
			Messages:   []string{fmt.Sprintf("From above you spot a glinting cache of %d gold half-buried in the ground!", gold)},
			GoldChange: gold,
		}
	})

	// e119 - Settlement spotted from air
	RegisterEvent("e119", func(s *GameState, ctx EventContext) EventResult {
		for _, adj := range AdjacentHexes(s.CurrentHex) {
			h := GetHex(adj)
			if h != nil && h.IsSettlement() && !s.VisitedHexes[adj] {
				s.VisitedHexes[adj] = true
				return EventResult{
					Messages: []string{fmt.Sprintf("From the air you spot a settlement at hex %s — it wasn't on your maps!", adj)},
					Note:     fmt.Sprintf("Settlement spotted from the air at hex %s (day %d).", adj, s.Day),
				}
			}
		}
		return EventResult{
			Messages: []string{"The view from above gives you a clear picture of the land ahead.", "You orient yourself precisely — no chance of getting lost today."},
		}
	})

	// ── Desert and swamp terrain events e122-e127 ─────────────────────────────

	// e122 - Desert Nomads
	RegisterEvent("e122", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A band of desert nomads approaches, their camels laden with trade goods."},
			Choices:  []string{"Trade with them (buy 6 food for 12 gold)", "Ignore them"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold >= 12 {
						s.Gold -= 12
						s.FoodUnits += 6
						return EventResult{Messages: []string{"You trade with the nomads. Gained 6 food for 12 gold."}}
					}
					if s.Gold >= 4 {
						units := s.Gold / 2
						cost := units * 2
						s.Gold -= cost
						s.FoodUnits += units
						return EventResult{Messages: []string{fmt.Sprintf("You can only afford %d food for %d gold.", units, cost)}}
					}
					return EventResult{Messages: []string{"You cannot afford their prices. They ride on."}}
				}
				return EventResult{Messages: []string{"The nomads pass by without incident."}}
			},
		}
	})

	// e123 - Oasis
	RegisterEvent("e123", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d3() + 2
		healed := 0
		if s.Prince.Wounds > 0 {
			healed = 1
			s.Prince.Wounds--
		}
		msg := fmt.Sprintf("You discover a lush oasis! Fresh water and shade restore your strength. Gained %d food.", food)
		if healed > 0 {
			msg += " The rest heals 1 wound."
		}
		return EventResult{
			Messages:   []string{msg},
			FoodChange: food,
			Note:       fmt.Sprintf("Oasis at hex %s (day %d) — a welcome sight in this wasteland.", s.CurrentHex, s.Day),
		}
	})

	// e124 - Desert Ruins
	RegisterEvent("e124", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"Half-buried in sand you find the ruins of a desert city."},
			Choices:  []string{"Search the ruins (takes time)", "Press on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					roll := Roll1d6()
					if roll >= 5 {
						gold := TreasureRoll(5, Roll1d6())
						return EventResult{
							Messages:   []string{fmt.Sprintf("The ruins yield ancient treasure! %d gold!", gold)},
							GoldChange: gold,
						}
					}
					if roll >= 3 {
						enemy := MakeEnemy("Desert Wraith", 5, 10, 4)
						enemy.IsUndead = true
						return EventResult{
							Messages:        []string{"Something stirs in the sand — a desert wraith rises to protect its tomb!"},
							CombatTriggered: true,
							Enemy:           &enemy,
						}
					}
					return EventResult{
						Messages:   []string{"The ruins are empty save for stinging sand. The search costs you food."},
						FoodChange: -1,
					}
				}
				return EventResult{Messages: []string{"You leave the half-buried city behind."}}
			},
		}
	})

	// e125 - Sandstorm
	RegisterEvent("e125", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d3()
		wounds := 1
		s.Prince.Wounds += wounds
		return EventResult{
			Messages: []string{
				"A massive sandstorm descends without warning!",
				fmt.Sprintf("You shelter as best you can, but the storm costs you %d food and inflicts 1 wound.", food),
				"After hours of howling grit, the storm passes.",
			},
			FoodChange: -food,
		}
	})

	// e126 - Desert Merchant
	RegisterEvent("e126", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A lone merchant on a camel halts you. He sells desert supplies at a steep premium."},
			Choices:  []string{"Buy 4 food for 10 gold (premium price)", "Bargain (Wit/Wiles check)", "Ignore him"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if s.Gold >= 10 {
						s.Gold -= 10
						s.FoodUnits += 4
						return EventResult{Messages: []string{"You pay the premium. 4 food acquired."}}
					}
					return EventResult{Messages: []string{"You can't afford his prices."}}
				case 1:
					if Roll1d6()+s.Prince.WitWiles >= 8 {
						if s.Gold >= 6 {
							s.Gold -= 6
							s.FoodUnits += 4
							return EventResult{Messages: []string{"A bit of haggling wins you 4 food for just 6 gold."}}
						}
					}
					return EventResult{Messages: []string{"The merchant won't budge on his price."}}
				default:
					return EventResult{Messages: []string{"You wave the merchant away and carry on."}}
				}
			},
		}
	})

	// e127 - Mirage / Lost Day
	RegisterEvent("e127", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{
				"Shimmering mirages of water dance on the horizon, leading you in circles.",
				"You lose a day chasing phantom oases. Precious food is consumed.",
			},
			FoodChange: -Roll1d3(),
		}
	})

	// ── Additional ruins events e139-e142 ─────────────────────────────────────

	// e139 - Trapped Corridor
	RegisterEvent("e139", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		return EventResult{
			Messages: []string{"You find a long corridor lined with ancient mechanisms."},
			Choices:  []string{"Proceed carefully", "Rush through"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if Roll2d6()+s.Prince.WitWiles >= 9 {
						gold := TreasureRoll(4, Roll1d6())
						return EventResult{
							Messages:   []string{fmt.Sprintf("Your careful approach disarms the traps. The far chamber holds %d gold.", gold)},
							GoldChange: gold,
						}
					}
					wounds := Roll1d3()
					s.Prince.Wounds += wounds
					return EventResult{
						Messages: []string{fmt.Sprintf("A pressure plate triggers a volley of darts! %d wounds.", wounds)},
					}
				}
				wounds := Roll1d3() + 1
				s.Prince.Wounds += wounds
				return EventResult{
					Messages: []string{fmt.Sprintf("You trigger multiple traps rushing through! %d wounds from blades and darts!", wounds)},
				}
			},
		}
	})

	// e140 - Ancient Library
	RegisterEvent("e140", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		roll := Roll1d6()
		if roll >= 5 {
			gold := Roll1d6() * 30
			return EventResult{
				Messages:   []string{"You find an intact ancient library! A traveling scholar pays handsomely for copies of rare texts.", fmt.Sprintf("You receive %d gold.", gold)},
				GoldChange: gold,
				Note:       fmt.Sprintf("Ancient library at hex %s (day %d) — texts copied and sold.", s.CurrentHex, s.Day),
			}
		}
		if roll >= 3 {
			// Knowledge of Win condition
			return EventResult{
				Messages: []string{
					"The library contains fragmentary records of the old kingdom.",
					"One scroll describes the Royal Helm: 'Hidden in the deepest vault, where the crown of kings was laid to rest.'",
				},
				Note: fmt.Sprintf("Library at hex %s (day %d): Royal Helm lies in the deep vault of ancient ruins.", s.CurrentHex, s.Day),
			}
		}
		return EventResult{
			Messages: []string{"The library's texts have crumbled to dust. You find nothing of value."},
		}
	})

	// e141 - Undead Swarm
	RegisterEvent("e141", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		enemy := MakeEnemy("Undead Swarm", 6, 18, 5)
		enemy.IsUndead = true
		return EventResult{
			Messages: []string{
				"The ruins erupt with the dead! Dozens of skeletal figures claw their way from the rubble!",
				"This is no ordinary haunt — the entire ruin is cursed!",
			},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  false,
		}
	})

	// e142 - Crystal Guardian
	RegisterEvent("e142", func(s *GameState, ctx EventContext) EventResult {
		s.GetHexFlags(s.CurrentHex).Searched = true
		gold := TreasureRoll(6, Roll1d6())
		enemy := MakeEnemy("Crystal Guardian", 6, 14, 6)
		return EventResult{
			Messages: []string{
				"At the heart of the ruins stands a towering construct of living crystal.",
				fmt.Sprintf("Beyond it gleams a treasury worth %d gold!", gold),
			},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 4,
			GoldChange:      gold,
		}
	})

	// ── Special NPC / audience outcomes e162-e179 ─────────────────────────────

	// e162 - Lord Demands Tribute
	RegisterEvent("e162", func(s *GameState, ctx EventContext) EventResult {
		tribute := Roll1d6()*10 + 20
		return EventResult{
			Messages: []string{fmt.Sprintf("The lord eyes you coldly. \"Strangers in my lands pay a toll of %d gold.\"", tribute)},
			Choices:  []string{fmt.Sprintf("Pay %d gold", tribute), "Refuse and fight your way out"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold >= tribute {
						s.Gold -= tribute
						return EventResult{Messages: []string{"You pay the toll. The lord waves you through with obvious satisfaction."}}
					}
					return EventResult{Messages: []string{"You lack the gold. The lord has you thrown out."}}
				}
				enemy := MakeEnemy("Lord's Guard Captain", 5, 12, 3)
				return EventResult{
					Messages:        []string{"You draw steel! The lord's guards close in!"},
					CombatTriggered: true,
					Enemy:           &enemy,
					PlayerAttFirst:  false,
				}
			},
		}
	})

	// e163 - Lady's Quest
	RegisterEvent("e163", func(s *GameState, ctx EventContext) EventResult {
		reward := Roll1d6()*20 + 50
		return EventResult{
			Messages: []string{
				"A noblewoman takes you aside after the audience.",
				fmt.Sprintf("\"My husband is held captive two days' journey east. Rescue him and I will pay you %d gold.\"", reward),
			},
			Choices: []string{"Accept the quest", "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					// Simulate the rescue — roll to see outcome
					if Roll2d6()+s.Prince.EffectiveCombatSkill() >= 10 {
						s.Gold += reward
						follower := Character{
							Name:         "Rescued Noble",
							Type:         TypeSwordsman,
							CombatSkill:  4,
							MaxEndurance: 8,
							DailyWage:    0,
							Morale:       5,
						}
						return EventResult{
							Messages:    []string{fmt.Sprintf("You rescue the nobleman! The grateful lady pays you %d gold.", reward), "The nobleman insists on accompanying you as recompense."},
							NewFollower: &follower,
						}
					}
					wounds := Roll1d3()
					s.Prince.Wounds += wounds
					return EventResult{
						Messages: []string{fmt.Sprintf("The rescue attempt goes wrong. %d wounds, and the captive is gone.", wounds)},
					}
				}
				return EventResult{Messages: []string{"You decline politely. The lady nods, her expression unreadable."}}
			},
		}
	})

	// e164 - Merchant Spy
	RegisterEvent("e164", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A well-dressed merchant corners you in the market. \"I have information that could help you — for a price.\""},
			Choices:  []string{"Buy the information (15 gold)", "Ignore him"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold < 15 {
						return EventResult{Messages: []string{"You can't afford his price."}}
					}
					s.Gold -= 15
					roll := Roll1d6()
					switch {
					case roll <= 2:
						gold := Roll1d6() * 20
						s.Gold += gold
						return EventResult{
							Messages: []string{fmt.Sprintf("He reveals a nearby hidden stash! You find %d gold.", gold)},
							Note:     fmt.Sprintf("Merchant's tip: stash near hex %s (day %d).", s.CurrentHex, s.Day),
						}
					case roll <= 4:
						return EventResult{
							Messages: []string{"He tells you the location of an unsearched ruin two days north.", "\"The guardians sleep during daylight — go then.\""},
							Note:     fmt.Sprintf("Merchant's tip: unsearched ruin north of hex %s (day %d).", s.CurrentHex, s.Day),
						}
					default:
						return EventResult{Messages: []string{"His \"information\" turns out to be common knowledge. You've been swindled."}}
					}
				}
				return EventResult{Messages: []string{"You brush past the merchant."}}
			},
		}
	})

	// e165 - Rebel Faction
	RegisterEvent("e165", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{
				"You are approached by hooded figures — rebels against your usurper!",
				"They know who you are and offer what aid they can.",
			},
			Choices: []string{"Accept their aid (free food)", "Decline — too risky"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					food := Roll1d6() + 3
					return EventResult{
						Messages:   []string{fmt.Sprintf("The rebels share their stores. You gain %d food and their blessing.", food)},
						FoodChange: food,
						Note:       fmt.Sprintf("Rebel contact at hex %s (day %d) — they know your name.", s.CurrentHex, s.Day),
					}
				}
				return EventResult{Messages: []string{"You decline and slip away before anyone notices the meeting."}}
			},
		}
	})

	// e166 - Tournament at the Castle
	RegisterEvent("e166", func(s *GameState, ctx EventContext) EventResult {
		prize := Roll1d6()*25 + 50
		entry := 10
		return EventResult{
			Messages: []string{fmt.Sprintf("A tournament is underway at the castle! The prize is %d gold. Entry costs %d gold.", prize, entry)},
			Choices:  []string{fmt.Sprintf("Enter the tournament (%d gold entry)", entry), "Watch from the sidelines"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold < entry {
						return EventResult{Messages: []string{"You cannot afford the entry fee."}}
					}
					s.Gold -= entry
					if Roll1d6()+s.Prince.EffectiveCombatSkill() >= 9 {
						s.Gold += prize
						return EventResult{Messages: []string{fmt.Sprintf("You fight your way through the brackets and WIN the tournament! %d gold is yours!", prize)}}
					}
					return EventResult{Messages: []string{"You fight well but are eliminated in the semifinals. No prize money."}}
				}
				return EventResult{Messages: []string{"You watch the jousting and wrestling from the crowd — a welcome distraction."}}
			},
		}
	})

	// e167 - Wandering Noble
	RegisterEvent("e167", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A young noble on horseback hails you. He is travelling to the south and is dangerously naive about the road."},
			Choices:  []string{"Escort him safely (gain 30 gold)", "Send him on his way"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if Roll1d6() >= 3 {
						return EventResult{
							Messages:   []string{"You safely escort the noble to his destination. He presses 30 gold into your hand."},
							GoldChange: 30,
						}
					}
					enemy := MakeEnemy("Road Bandits", 4, 10, 3)
					return EventResult{
						Messages:        []string{"Bandits ambush you on the road! The noble hides behind a tree."},
						CombatTriggered: true,
						Enemy:           &enemy,
						GoldChange:      30, // paid after combat
					}
				}
				return EventResult{Messages: []string{"You wish the noble well and go your separate ways."}}
			},
		}
	})

	// e168 - Wagon Robbery
	RegisterEvent("e168", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You come upon an overturned merchant wagon surrounded by bandits!"},
			Choices:  []string{"Intervene and fight", "Slip past unnoticed"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					enemy := MakeEnemy("Road Bandits", 4, 12, 3)
					return EventResult{
						Messages:        []string{"You charge the bandits! The merchant cowers in his wagon."},
						CombatTriggered: true,
						Enemy:           &enemy,
						GoldChange:      Roll1d6() * 15, // reward from grateful merchant
					}
				}
				if Roll1d6() >= 4 {
					return EventResult{Messages: []string{"You slip past in the commotion, unseen."}}
				}
				food := Roll1d3()
				return EventResult{
					Messages:   []string{"A stray bandit spots you! You grab some fallen supplies and run."},
					FoodChange: food,
				}
			},
		}
	})

	// e169 - High Priest Travelling
	RegisterEvent("e169", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You encounter a high priest and his retinue travelling the road."},
			Choices:  []string{"Request a blessing", "Request healing (cost 10 gold)", "Let them pass"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if Roll1d6() >= 4 {
						return EventResult{Messages: []string{"The priest blesses your quest. You feel renewed purpose — and the gods are watching."}}
					}
					return EventResult{Messages: []string{"The priest offers a brief prayer and continues on his way."}}
				case 1:
					if s.Gold >= 10 {
						s.Gold -= 10
						healed := Roll1d3()
						if healed > s.Prince.Wounds {
							healed = s.Prince.Wounds
						}
						s.Prince.Wounds -= healed
						if s.Prince.PoisonWounds > 0 {
							s.Prince.PoisonWounds--
						}
						return EventResult{Messages: []string{fmt.Sprintf("The priest's healing arts restore %d wounds and ease the poison.", healed)}}
					}
					return EventResult{Messages: []string{"You cannot afford the priest's services."}}
				default:
					return EventResult{Messages: []string{"The priest nods to you as he passes."}}
				}
			},
		}
	})

	// e170 - Desert Oracle
	RegisterEvent("e170", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"Buried to the neck in sand with only their head showing, an ancient oracle speaks to you."},
			Choices:  []string{"Ask about the Royal Helm (5 gold)", "Ask about gold (5 gold)", "Leave quickly"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if s.Gold < 5 {
						return EventResult{Messages: []string{"The oracle demands 5 gold. You cannot pay."}}
					}
					s.Gold -= 5
					if !s.Flags.HasRoyalHelm {
						return EventResult{
							Messages: []string{
								"\"The Helm of Kings sleeps beneath shattered stone,\"",
								"\"in the ruin where three roads once crossed,\"",
								"\"guarded by the memory of the old dynasty.\"",
							},
							Note: fmt.Sprintf("Oracle at hex %s (day %d): Helm in ruins where three roads crossed.", s.CurrentHex, s.Day),
						}
					}
					return EventResult{Messages: []string{"\"You already carry the crown's twin. The path north opens before you.\""}}
				case 1:
					if s.Gold < 5 {
						return EventResult{Messages: []string{"The oracle demands 5 gold. You cannot pay."}}
					}
					s.Gold -= 5
					gold := TreasureRoll(4, Roll1d6())
					s.Gold += gold
					return EventResult{
						Messages: []string{fmt.Sprintf("\"Dig three paces east of the cracked stone.\" You find %d gold!", gold)},
					}
				default:
					return EventResult{Messages: []string{"You back away from the disturbing sight."}}
				}
			},
		}
	})

	// e171 - Mountain Hermit King
	RegisterEvent("e171", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{
				"A greybearded old man sits on a throne of stones, wrapped in animal furs.",
				"\"I was once a king,\" he says. \"Sit. I have wisdom to share.\"",
			},
			Choices: []string{"Sit and listen", "Move on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					roll := Roll1d6()
					if roll >= 5 {
						s.Prince.WitWiles++
						return EventResult{Messages: []string{"The old king's words of statecraft are invaluable. +1 Wit/Wiles."}}
					}
					if roll >= 3 {
						food := Roll1d3() + 2
						return EventResult{
							Messages:   []string{fmt.Sprintf("The hermit king shares a hearty meal. Gained %d food.", food)},
							FoodChange: food,
						}
					}
					return EventResult{Messages: []string{"His tales are long and mostly rambling, but not unpleasant."}}
				}
				return EventResult{Messages: []string{"The old man shrugs and returns to his solitude."}}
			},
		}
	})

	// e172 - Forest Witch
	RegisterEvent("e172", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A forest witch steps from between the trees, green eyes sharp and unreadable."},
			Choices:  []string{"Ask for a potion", "Ask for knowledge", "Back away carefully"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if Roll1d6() >= 3 {
						s.Prince.AddPossession(PossHealingPotion)
						return EventResult{Messages: []string{"The witch hands you a glowing vial. \"Drink when you must.\" You receive a Healing Potion!"}}
					}
					s.Prince.PoisonWounds += 1
					return EventResult{Messages: []string{"The witch cackles and throws a cursed powder at you! +1 poison wound."}}
				case 1:
					if Roll1d6() >= 4 {
						return EventResult{
							Messages: []string{
								"\"The thing you seek lies where the old magic is strongest.\"",
								"She points south. You note the direction.",
							},
							Note: fmt.Sprintf("Forest witch at hex %s (day %d): seek where old magic is strongest, to the south.", s.CurrentHex, s.Day),
						}
					}
					return EventResult{Messages: []string{"\"Ask the stones,\" she says, then vanishes into the wood."}}
				default:
					if Roll1d6() >= 5 {
						wounds := Roll1d3()
						s.Prince.Wounds += wounds
						return EventResult{
							Messages: []string{fmt.Sprintf("The witch takes offense at your rudeness! %d wounds from a hurled curse!", wounds)},
						}
					}
					return EventResult{Messages: []string{"You back away slowly. The witch watches you go without moving."}}
				}
			},
		}
	})

	// e173 - Coastal Chart
	RegisterEvent("e173", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6() * 15
		return EventResult{
			Messages: []string{
				"You find a sea captain's log and coastal navigation charts washed ashore.",
				fmt.Sprintf("A local trader pays %d gold for the valuable charts.", gold),
			},
			GoldChange: gold,
			Note:       fmt.Sprintf("Coastal charts found near hex %s (day %d) — sold to trader.", s.CurrentHex, s.Day),
		}
	})

	// e174 - Bandit Lord
	RegisterEvent("e174", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A richly dressed bandit lord bars your path with a dozen armed men."},
			Choices:  []string{"Pay his toll (30 gold)", "Challenge him to single combat", "Try to bluff past"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if s.Gold >= 30 {
						s.Gold -= 30
						return EventResult{Messages: []string{"You pay the toll. The bandit lord bows mockingly and stands aside."}}
					}
					enemy := MakeEnemy("Bandit Lord", 5, 12, 5)
					return EventResult{
						Messages:        []string{"You can't pay! The bandit lord draws his sword."},
						CombatTriggered: true,
						Enemy:           &enemy,
					}
				case 1:
					enemy := MakeEnemy("Bandit Lord", 5, 12, 5)
					return EventResult{
						Messages:        []string{"The bandit lord grins. \"A man of courage! En garde!\""},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  true,
						GoldChange:      Roll1d6() * 30, // his loot if you win
					}
				default:
					if Roll1d6()+s.Prince.WitWiles >= 9 {
						return EventResult{Messages: []string{"Your confident bearing convinces the bandits you are a dangerous target. They let you pass."}}
					}
					enemy := MakeEnemy("Bandit Lord", 5, 12, 5)
					return EventResult{
						Messages:        []string{"The bluff fails! The bandit lord laughs and attacks!"},
						CombatTriggered: true,
						Enemy:           &enemy,
					}
				}
			},
		}
	})

	// e175 - Dwarven Forge
	RegisterEvent("e175", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"Smoke rises from a hidden cleft in the rock — a dwarven forge!"},
			Choices:  []string{"Trade with the dwarves", "Request weapon upgrade (50 gold)", "Continue on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					food := Roll1d6() + 2
					return EventResult{
						Messages:   []string{fmt.Sprintf("The dwarves trade food and provisions. Gained %d food.", food)},
						FoodChange: food,
					}
				case 1:
					if s.Gold >= 50 {
						s.Gold -= 50
						s.Prince.CombatSkill++
						return EventResult{Messages: []string{"The master smith reforges your blade with dwarven steel. +1 permanent Combat Skill!"}}
					}
					return EventResult{Messages: []string{"You cannot afford their craft prices."}}
				default:
					return EventResult{Messages: []string{"You leave the dwarves to their work."}}
				}
			},
		}
	})

	// e176 - Elven Ambassador
	RegisterEvent("e176", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"An elven ambassador hails you from the treetops in accented but perfect speech."},
			Choices:  []string{"Speak with her", "Continue on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					roll := Roll1d6()
					if roll >= 5 {
						follower := Character{
							Name:         "Elven Scout",
							Type:         TypeElf,
							CombatSkill:  5,
							MaxEndurance: 9,
							DailyWage:    4,
							Morale:       6,
							IsGuide:      true,
						}
						return EventResult{
							Messages: []string{"The ambassador is impressed by your bearing.", "She offers the service of her scout — a skilled elven warrior and guide."},
							Choices:  []string{"Accept the elven scout (4 gold/day)", "Decline"},
							ChoiceHandler: func(s *GameState, c int) EventResult {
								if c == 0 {
									return EventResult{NewFollower: &follower}
								}
								return EventResult{Messages: []string{"You decline. The ambassador nods with quiet respect."}}
							},
						}
					}
					gold := Roll1d6() * 15
					return EventResult{
						Messages:   []string{fmt.Sprintf("The ambassador shares elven provisions and %d gold — a gesture of goodwill.", gold)},
						FoodChange: Roll1d3(),
						GoldChange: gold,
					}
				}
				return EventResult{Messages: []string{"You walk on. From above, elven eyes follow your path in silence."}}
			},
		}
	})

	// e177 - Ancient Knight
	RegisterEvent("e177", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{
				"In the deepest chamber of the ruins, a figure in ancient armour stands motionless.",
				"As you approach, it turns. Its eyes gleam with cold blue light — but it does not attack.",
				"\"Are you of the old blood?\" it rasps.",
			},
			Choices: []string{"\"I am Cal Arath, rightful prince of this realm.\"", "Attack it", "Flee"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					gold := TreasureRoll(6, Roll1d6())
					s.Gold += gold
					return EventResult{
						Messages: []string{
							"The ancient knight bows slowly.",
							"\"Then take what is yours, my prince. I have guarded it long enough.\"",
							fmt.Sprintf("It gestures to an alcove filled with %d gold — the treasury of your ancestors!", gold),
						},
					}
				case 1:
					enemy := MakeEnemy("Ancient Knight", 7, 16, 6)
					enemy.IsUndead = true
					return EventResult{
						Messages:        []string{"The ancient knight draws a blade of pale fire!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				default:
					return EventResult{Messages: []string{"You flee the chamber. Behind you, the blue-lit eyes watch you go."}}
				}
			},
		}
	})

	// e178 - Pretender Prince
	RegisterEvent("e178", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6()*20 + 40
		enemy := MakeEnemy("Pretender's Champion", 5, 12, 4)
		return EventResult{
			Messages: []string{
				"A proclamation nailed to a tree: \"REWARD — for the capture of the rebel Cal Arath.\"",
				fmt.Sprintf("A hunter has read the poster and steps out with levelled crossbow. The bounty is %d gold.", gold),
			},
			Choices: []string{"Fight your way past", "Try to persuade him (Wit/Wiles)", "Run!"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					return EventResult{
						Messages:        []string{"You draw your blade!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				case 1:
					if Roll1d6()+s.Prince.WitWiles >= 9 {
						return EventResult{Messages: []string{"You convince the hunter that you are merely a travelling merchant. He lowers his weapon and lets you pass."}}
					}
					return EventResult{
						Messages:        []string{"The persuasion fails. The hunter attacks!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				default:
					if Roll1d6() >= 3 {
						return EventResult{Messages: []string{"You sprint away and lose the hunter in the undergrowth!"}}
					}
					wounds := Roll1d3()
					s.Prince.Wounds += wounds
					return EventResult{
						Messages: []string{fmt.Sprintf("A crossbow bolt grazes you as you flee. %d wounds.", wounds)},
					}
				}
			},
		}
	})

	// e179 - Sage's Library (countryside)
	RegisterEvent("e179", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A travelling sage has set up camp, his pack mule laden with books and scrolls."},
			Choices:  []string{"Buy information (20 gold)", "Share news in exchange for food", "Move on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if s.Gold < 20 {
						return EventResult{Messages: []string{"You cannot afford his consultation fee."}}
					}
					s.Gold -= 20
					roll := Roll1d6()
					switch {
					case roll >= 5:
						return EventResult{
							Messages: []string{
								"The sage's maps are extraordinary — he reveals two hidden hexes on your route.",
								"\"The eastern path is faster but far more dangerous,\" he warns.",
							},
							Note: fmt.Sprintf("Sage near hex %s (day %d): hidden paths east — faster but dangerous.", s.CurrentHex, s.Day),
						}
					case roll >= 3:
						return EventResult{
							Messages: []string{
								"The sage knows of the Royal Helm.",
								"\"In my travels I passed ruins with that exact description — marked with the dynasty's crest.\"",
							},
							Note: fmt.Sprintf("Sage near hex %s (day %d): ruins with dynasty's crest possibly contain the Royal Helm.", s.CurrentHex, s.Day),
						}
					default:
						return EventResult{Messages: []string{"The sage's information is outdated. A waste of 20 gold."}}
					}
				case 1:
					food := Roll1d3() + 1
					return EventResult{
						Messages:   []string{fmt.Sprintf("You trade news for %d food and a warm conversation.", food)},
						FoodChange: food,
					}
				default:
					return EventResult{Messages: []string{"You pass the sage with a nod."}}
				}
			},
		}
	})

	// ── Improved stub events ──────────────────────────────────────────────────

	// e034 - Treasure rumour (now with actual mechanic)
	RegisterEvent("e034", func(s *GameState, ctx EventContext) EventResult {
		roll := Roll1d6()
		if roll >= 4 {
			// Reveal a nearby cache hint
			gold := Roll1d6() * 5
			s.Gold += gold
			return EventResult{
				Messages: []string{
					"Passing travellers whisper of great treasure somewhere to the south.",
					fmt.Sprintf("One presses a map fragment into your hand for %d gold.", gold),
				},
				Note: fmt.Sprintf("Treasure rumour near hex %s (day %d) — said to be south.", s.CurrentHex, s.Day),
			}
		}
		return EventResult{
			Messages: []string{
				"Passing travellers whisper of great treasure somewhere to the south.",
				"The rumours are vague, but you note the direction.",
			},
			Note: fmt.Sprintf("Treasure rumour heard near hex %s (day %d).", s.CurrentHex, s.Day),
		}
	})

	// e067 - Ancient map (now with actual reveal mechanic)
	RegisterEvent("e067", func(s *GameState, ctx EventContext) EventResult {
		for _, adj := range AdjacentHexes(s.CurrentHex) {
			if !s.VisitedHexes[adj] {
				h := GetHex(adj)
				if h != nil && (h.IsRuins() || h.IsSettlement()) {
					s.VisitedHexes[adj] = true
					label := "ruins"
					if h.IsSettlement() {
						label = "settlement"
					}
					return EventResult{
						Messages: []string{
							fmt.Sprintf("You find an ancient map showing the layout of nearby %s!", label),
							fmt.Sprintf("It reveals a %s at hex %s.", label, adj),
						},
						Note: fmt.Sprintf("Ancient map revealed %s at hex %s (day %d).", label, adj, s.Day),
					}
				}
			}
		}
		gold := Roll1d6() * 10
		return EventResult{
			Messages: []string{
				"You find an ancient map, but the locations it marks are too far away to use directly.",
				fmt.Sprintf("A local scholar pays you %d gold for the cartographic curiosity.", gold),
			},
			GoldChange: gold,
		}
	})
}
