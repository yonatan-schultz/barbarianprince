package game

import "fmt"

func registerEventsE001() {
	// e001 - stub (game start handled in NewGameState)
	RegisterEvent("e001", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"Your adventure begins."}}
	})

	// e002 - stub
	RegisterEvent("e002", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"The road is quiet today."}}
	})

	// ── Follower events (e003–e008) ───────────────────────────────────────────
	// All have a hire/decline choice so the player actually gets to decide.

	// e003 - Swordsman
	RegisterEvent("e003", func(s *GameState, ctx EventContext) EventResult {
		roll := Roll1d6()
		if roll <= 2 {
			enemy := MakeEnemy("Mercenary Swordsman", 4, 8, 3)
			return EventResult{
				Messages:        []string{"A mercenary swordsman challenges you to a duel for your gold!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  SurpriseCheck(""),
			}
		}
		if roll <= 4 {
			wage := 3
			follower := Character{
				Name: "Hired Swordsman", Type: TypeSwordsman,
				CombatSkill: 4, MaxEndurance: 8, DailyWage: wage, Morale: 4,
			}
			return EventResult{
				Messages: []string{fmt.Sprintf("A swordsman offers his services for %d gold/day.", wage)},
				Choices:  []string{"Hire him", "Decline"},
				ChoiceHandler: func(s *GameState, choice int) EventResult {
					if choice == 0 {
						return EventResult{NewFollower: &follower}
					}
					return EventResult{Messages: []string{"You move on without hiring him."}}
				},
			}
		}
		gold := Roll1d6() * 5
		return EventResult{
			Messages:   []string{fmt.Sprintf("A friendly swordsman shares news and tips you %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e004 - Mercenary Band
	RegisterEvent("e004", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d6() + 2
		if Roll1d6() <= 3 {
			enemy := MakeEnemy(fmt.Sprintf("Mercenary Band (%d)", count), count+2, count*2, 4)
			return EventResult{
				Messages:        []string{fmt.Sprintf("A band of %d mercenaries blocks your path!", count)},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  false,
			}
		}
		follower := Character{
			Name: "Mercenary Captain", Type: TypeMercenary,
			CombatSkill: 5, MaxEndurance: 10, DailyWage: 5, Morale: 4,
		}
		return EventResult{
			Messages: []string{fmt.Sprintf("A band of %d mercenaries offers their captain's services for 5 gold/day.", count)},
			Choices:  []string{"Hire the captain", "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					return EventResult{NewFollower: &follower}
				}
				return EventResult{Messages: []string{"The band moves on."}}
			},
		}
	})

	// e005 - Amazons
	RegisterEvent("e005", func(s *GameState, ctx EventContext) EventResult {
		roll := Roll1d6()
		if roll <= 2 {
			enemy := MakeEnemy("Amazon Warriors", 5, 9, 4)
			return EventResult{
				Messages:        []string{"A company of Amazons regards you with hostile intent!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  SurpriseCheck("enemy_first"),
			}
		}
		if roll <= 4 {
			follower := Character{
				Name: "Amazon Scout", Type: TypeAmazon,
				CombatSkill: 5, MaxEndurance: 9, DailyWage: 4, Morale: 5, IsGuide: true,
			}
			return EventResult{
				Messages: []string{"An Amazon warrior, impressed by your bearing, offers to guide you for 4 gold/day."},
				Choices:  []string{"Accept her service", "Decline"},
				ChoiceHandler: func(s *GameState, choice int) EventResult {
					if choice == 0 {
						return EventResult{NewFollower: &follower}
					}
					return EventResult{Messages: []string{"She returns to her companions."}}
				},
			}
		}
		return EventResult{Messages: []string{"A company of Amazons passes without incident."}}
	})

	// e006 - Dwarf
	RegisterEvent("e006", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() <= 3 {
			follower := Character{
				Name: "Gromli the Dwarf", Type: TypeDwarf,
				CombatSkill: 6, MaxEndurance: 12, DailyWage: 4, Morale: 5, WitWiles: 2,
			}
			return EventResult{
				Messages: []string{"A grizzled dwarf with a battle-axe eyes you appraisingly. He wants 4 gold/day."},
				Choices:  []string{"Hire Gromli", "Move on"},
				ChoiceHandler: func(s *GameState, choice int) EventResult {
					if choice == 0 {
						return EventResult{NewFollower: &follower}
					}
					return EventResult{Messages: []string{"The dwarf shrugs and returns to his ale."}}
				},
			}
		}
		gold := Roll1d6() * 10
		return EventResult{
			Messages:   []string{fmt.Sprintf("A friendly dwarf trader pays you %d gold for news from afar.", gold)},
			GoldChange: gold,
		}
	})

	// e007 - Elf
	RegisterEvent("e007", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() <= 3 {
			follower := Character{
				Name: "Sylara the Elf", Type: TypeElf,
				CombatSkill: 5, MaxEndurance: 8, DailyWage: 5, Morale: 5, IsGuide: true, WitWiles: 3,
			}
			return EventResult{
				Messages: []string{"An elven ranger emerges from the shadows and offers her skills for 5 gold/day."},
				Choices:  []string{"Accept (5 gold/day)", "Decline"},
				ChoiceHandler: func(s *GameState, choice int) EventResult {
					if choice == 0 {
						return EventResult{NewFollower: &follower}
					}
					return EventResult{Messages: []string{"She melts back into the forest."}}
				},
			}
		}
		return EventResult{Messages: []string{"An elf warns you of dangers to the south and vanishes."}}
	})

	// e008 - Halfling
	RegisterEvent("e008", func(s *GameState, ctx EventContext) EventResult {
		follower := Character{
			Name: "Pip the Halfling", Type: TypeHalfling,
			CombatSkill: 3, MaxEndurance: 6, DailyWage: 2, Morale: 4, WitWiles: 5, IsGuide: true,
		}
		return EventResult{
			Messages: []string{"A cheerful halfling offers to serve as scout and guide for 2 gold/day."},
			Choices:  []string{"Hire Pip", "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					return EventResult{NewFollower: &follower}
				}
				return EventResult{Messages: []string{"Pip tips his cap and wanders off humming."}}
			},
		}
	})

	// ── Farmland / settlement events (e009–e034, e040, e050) ─────────────────

	// e009 - Farm (buy food at market price, afford what you can)
	RegisterEvent("e009", func(s *GameState, ctx EventContext) EventResult {
		available := Roll1d6() + 2
		cost := available * FoodCostPerUnit
		if s.Gold < cost {
			affordable := s.Gold / FoodCostPerUnit
			if affordable == 0 {
				return EventResult{Messages: []string{"You come upon a farmstead, but cannot afford to buy food."}}
			}
			available = affordable
			cost = available * FoodCostPerUnit
		}
		return EventResult{
			Messages:   []string{fmt.Sprintf("A farmstead sells you %d food units for %d gold.", available, cost)},
			FoodChange: available,
			GoldChange: -cost,
		}
	})

	// e010 - Village inn
	RegisterEvent("e010", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You find a comfortable inn. For 2 gold you can rest and recover a wound."},
			Choices:  []string{"Pay 2 gold and rest", "Move on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold < 2 {
						return EventResult{Messages: []string{"You can't afford the room."}}
					}
					s.Gold -= 2
					if s.Prince.Wounds > 0 {
						s.Prince.Wounds--
					}
					return EventResult{Messages: []string{"You rest well and recover your strength."}}
				}
				return EventResult{Messages: []string{"You move on without stopping."}}
			},
		}
	})

	// e011 - Merchant caravan under attack
	RegisterEvent("e011", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() <= 2 {
			gold := Roll1d6()*10 + 30
			enemy := MakeEnemy("Bandit Raiders", 4, 14, 3)
			return EventResult{
				Messages: []string{
					"Bandits are attacking a merchant caravan! The merchants beg for help.",
					fmt.Sprintf("They promise %d gold if you drive off the attackers.", gold),
				},
				Choices: []string{"Defend the caravan", "Slip away"},
				ChoiceHandler: func(s *GameState, choice int) EventResult {
					if choice == 0 {
						return EventResult{
							Messages:        []string{"You charge into the fray!"},
							CombatTriggered: true,
							Enemy:           &enemy,
							PlayerAttFirst:  true,
							GoldChange:      gold,
						}
					}
					return EventResult{Messages: []string{"You slip away as the merchants cry for help."}}
				},
			}
		}
		gold := Roll1d6()*10 + 20
		return EventResult{
			Messages:   []string{fmt.Sprintf("A merchant caravan pays you %d gold for armed escort through the area.", gold)},
			GoldChange: gold,
		}
	})

	// e012 - Pilgrims
	RegisterEvent("e012", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6()*5 + 10
		return EventResult{
			Messages: []string{"A group of pilgrims asks for your protection to the next settlement."},
			Choices:  []string{fmt.Sprintf("Escort them (reward: ~%d gold)", gold), "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					return EventResult{
						Messages:   []string{fmt.Sprintf("You escort the pilgrims safely. They reward you with %d gold.", gold)},
						GoldChange: gold,
					}
				}
				return EventResult{Messages: []string{"You leave the pilgrims to their journey."}}
			},
		}
	})

	// e013 - Beggar
	RegisterEvent("e013", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A ragged beggar holds out a trembling hand."},
			Choices:  []string{"Give 5 gold", "Ignore him"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold < 5 {
						return EventResult{Messages: []string{"You have nothing to spare."}}
					}
					s.Gold -= 5
					if Roll1d6() == 6 {
						gold := 30
						s.Gold += gold
						return EventResult{Messages: []string{
							"The beggar is a disguised spy!",
							fmt.Sprintf("He rewards your charity with %d gold and a whispered secret.", gold),
						}}
					}
					return EventResult{Messages: []string{"The beggar thanks you with tears in his eyes."}}
				}
				return EventResult{Messages: []string{"You walk past without meeting his eyes."}}
			},
		}
	})

	// e014 - Hunting party
	RegisterEvent("e014", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d6() + 1
		return EventResult{
			Messages:   []string{fmt.Sprintf("A hunting party shares %d food from their catch.", food)},
			FoodChange: food,
		}
	})

	// e015 - Lost child
	RegisterEvent("e015", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A lost child cries beside the road."},
			Choices:  []string{"Return the child to the nearest town", "Move on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					gold := Roll1d6()*5 + 10
					return EventResult{
						Messages:   []string{fmt.Sprintf("Grateful parents reward you with %d gold!", gold)},
						GoldChange: gold,
					}
				}
				return EventResult{Messages: []string{"The child's cries fade behind you."}}
			},
		}
	})

	// e016 - Soldiers patrol
	RegisterEvent("e016", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A patrol of soldiers stops you, looking for fugitives from the realm."},
			Choices:  []string{"Bribe them (10 gold)", "Brazen it out", "Run!"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if s.Gold < 10 {
						return EventResult{Messages: []string{"You lack the gold to bribe them. They grow suspicious."}}
					}
					s.Gold -= 10
					return EventResult{Messages: []string{"The soldiers pocket your coin and wave you through."}}
				case 1:
					if Roll1d6()+s.Prince.WitWiles >= 8 {
						return EventResult{Messages: []string{"You bluff your way past the patrol."}}
					}
					enemy := MakeEnemy("Soldiers", 4, 12, 2)
					return EventResult{
						Messages:        []string{"The soldiers don't buy it and attack!"},
						CombatTriggered: true,
						Enemy:           &enemy,
					}
				default:
					if Roll1d6() >= 4 {
						return EventResult{Messages: []string{"You vanish into the wilderness before they can react."}}
					}
					enemy := MakeEnemy("Pursuing Soldiers", 4, 12, 2)
					return EventResult{
						Messages:        []string{"They give chase and corner you!"},
						CombatTriggered: true,
						Enemy:           &enemy,
					}
				}
			},
		}
	})

	// e017 - Peasant mob
	RegisterEvent("e017", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"An angry peasant mob blocks the road, cursing all nobles!"},
			Choices:  []string{"Fight through", "Try to talk them down", "Detour around them"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					mob := MakeEnemy("Peasant Mob", 2, 20, 1)
					return EventResult{
						Messages:        []string{"The mob surges forward!"},
						CombatTriggered: true,
						Enemy:           &mob,
						PlayerAttFirst:  false,
					}
				case 1:
					if Roll1d6()+s.Prince.WitWiles >= 7 {
						return EventResult{Messages: []string{"You convince them you are no enemy of the common folk. They let you pass."}}
					}
					mob := MakeEnemy("Peasant Mob", 2, 20, 1)
					return EventResult{
						Messages:        []string{"Your words fall on deaf ears. They attack!"},
						CombatTriggered: true,
						Enemy:           &mob,
					}
				default:
					return EventResult{
						Messages:   []string{"You take a long detour. It costs you an extra food unit."},
						FoodChange: -1,
					}
				}
			},
		}
	})

	// e018 - Wandering priest
	RegisterEvent("e018", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() <= 3 {
			healed := Roll1d3()
			old := s.Prince.Wounds
			s.Prince.Wounds -= healed
			if s.Prince.Wounds < 0 {
				s.Prince.Wounds = 0
			}
			actual := old - s.Prince.Wounds
			return EventResult{
				Messages: []string{fmt.Sprintf("A wandering priest lays hands on you, healing %d wound(s).", actual)},
			}
		}
		follower := Character{
			Name: "Brother Aldric", Type: TypePriest,
			CombatSkill: 3, MaxEndurance: 7, DailyWage: 3, WitWiles: 4, Morale: 5,
		}
		return EventResult{
			Messages: []string{"A traveling priest offers to join your quest for the glory of the gods (3 gold/day)."},
			Choices:  []string{"Accept Brother Aldric", "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					return EventResult{NewFollower: &follower}
				}
				return EventResult{Messages: []string{"The priest blesses you and continues his journey."}}
			},
		}
	})

	// e019-e024: Farmland encounters
	RegisterEvent("e019", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6() * 5
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find %d gold coins scattered along the road.", gold)},
			GoldChange: gold,
		}
	})
	RegisterEvent("e020", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d6() + 3
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find an abandoned cart full of provisions — %d food units.", food)},
			FoodChange: food,
		}
	})
	RegisterEvent("e021", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages:   []string{"You find a waystation with fresh water and shelter."},
			FoodChange: 2,
		}
	})
	RegisterEvent("e022", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d6() + 1
		enemy := MakeEnemy(fmt.Sprintf("Light Bandits (%d)", count), 3, count*3, 2)
		return EventResult{
			Messages:        []string{fmt.Sprintf("%d bandits leap from hiding and demand your gold!", count)},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 4,
		}
	})
	RegisterEvent("e023", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() <= 3 {
			enemy := MakeEnemy("Wizard's Construct", 6, 10, 6)
			return EventResult{
				Messages:        []string{"A robed mage sends his magical construct against you!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  false,
			}
		}
		gold := Roll1d6() * 20
		return EventResult{
			Messages:   []string{fmt.Sprintf("A wizard pays %d gold for your silence about his location.", gold)},
			GoldChange: gold,
		}
	})
	RegisterEvent("e024", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A cloaked figure offers to sell you information."},
			Choices:  []string{"Pay 15 gold", "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if s.Gold < 15 {
						return EventResult{Messages: []string{"You can't afford it."}}
					}
					s.Gold -= 15
					return EventResult{Messages: []string{
						"The spy reveals that a cache of gold is buried in a nearby ruin.",
						"Search ruins in adjacent hexes to find it.",
					}}
				}
				return EventResult{Messages: []string{"You decline and the figure disappears into the crowd."}}
			},
		}
	})

	// e025-e034, e040: Treasure and misc
	RegisterEvent("e025", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(3, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("You discover hidden treasure — %d gold!", gold)},
			GoldChange: gold,
		}
	})
	RegisterEvent("e026", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"A rockslide blocks the road. You lose time finding a way around."}}
	})
	RegisterEvent("e027", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(4, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find the remains of a wealthy traveler — %d gold among their belongings.", gold)},
			GoldChange: gold,
		}
	})
	RegisterEvent("e028", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(5, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("A buried cache of treasure! %d gold and valuables!", gold)},
			GoldChange: gold,
		}
	})
	RegisterEvent("e029", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d6()*2 + 4
		gold := Roll1d6() * 5
		return EventResult{
			Messages:   []string{fmt.Sprintf("An abandoned campsite: %d food and %d gold.", food, gold)},
			FoodChange: food,
			GoldChange: gold,
		}
	})
	RegisterEvent("e030", func(s *GameState, ctx EventContext) EventResult {
		gold := Roll1d6() * 15
		return EventResult{
			Messages:   []string{fmt.Sprintf("Merchants pay %d gold for your armed company during the road stretch.", gold)},
			GoldChange: gold,
		}
	})
	RegisterEvent("e031", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"The road is quiet and you make good progress."}}
	})
	RegisterEvent("e032", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages:   []string{"Heavy rain turns the road to mud and spoils a day's rations."},
			FoodChange: -1,
		}
	})
	RegisterEvent("e033", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d3()
		return EventResult{
			Messages:   []string{fmt.Sprintf("You forage wild berries and edible plants — %d food units.", food)},
			FoodChange: food,
		}
	})
	RegisterEvent("e034", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{Messages: []string{"Passing travellers whisper of great treasure somewhere to the south."}}
	})

	// e040 - Rich find
	RegisterEvent("e040", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(6, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("Exceptional luck — a hidden cache of %d gold!", gold)},
			GoldChange: gold,
		}
	})

	// e050 - Constabulary
	RegisterEvent("e050", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"The local constabulary recognises you as a wanted man!"},
			Choices:  []string{"Fight your way free", "Surrender and bribe (20 gold)"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					enemy := MakeEnemy("Town Guard", 4, 10, 2)
					return EventResult{
						Messages:        []string{"You draw your blade!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				}
				if s.Gold >= 20 {
					s.Gold -= 20
					return EventResult{Messages: []string{"You bribe the guard captain. He looks the other way."}}
				}
				enemy := MakeEnemy("Town Guard", 4, 10, 2)
				return EventResult{
					Messages:        []string{"You can't afford the bribe — they move to arrest you!"},
					CombatTriggered: true,
					Enemy:           &enemy,
				}
			},
		}
	})
}
