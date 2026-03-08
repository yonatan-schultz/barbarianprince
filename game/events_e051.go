package game

import "fmt"

func registerEventsE051() {
	// e051 - Bandits (wilderness)
	RegisterEvent("e051", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d6() + 2
		cs := 3 + Roll1d3()
		enemy := MakeEnemy(fmt.Sprintf("Brigands (%d)", count), cs, count*3+4, 3)
		return EventResult{
			Messages:        []string{fmt.Sprintf("%d armed brigands ambush you from the tree line!", count)},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 5,
		}
	})

	// e052 - Wolf
	RegisterEvent("e052", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 5 {
			enemy := MakeEnemy("Lone Wolf", 4, 8, 1)
			return EventResult{
				Messages:        []string{"A large wolf snarls and attacks!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  Roll1d6() >= 4,
			}
		}
		return EventResult{Messages: []string{"A wolf watches you from a distance, then slinks away."}}
	})

	// e053 - Deserters
	RegisterEvent("e053", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d6() + 1
		enemy := MakeEnemy(fmt.Sprintf("Deserters (%d)", count), 3, count*4, 2)
		return EventResult{
			Messages:        []string{fmt.Sprintf("%d military deserters demand your equipment!", count)},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 4,
		}
	})

	// e054 - Strange fog
	RegisterEvent("e054", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"An unnatural fog envelops the land. You lose your bearings briefly.", "After hours of wandering, you find the path again."},
		}
	})

	// e055 - Wild boar
	RegisterEvent("e055", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 4 {
			enemy := MakeEnemy("Wild Boar", 3, 10, 1)
			return EventResult{
				Messages:        []string{"A large wild boar charges from the undergrowth!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  false,
			}
		}
		food := Roll1d3() + 1
		return EventResult{
			Messages:   []string{fmt.Sprintf("You successfully hunt a boar. Gained %d food units.", food)},
			FoodChange: food,
		}
	})

	// e056 - Giant ant colony
	RegisterEvent("e056", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("Giant Ants", 3, 15, 1)
		return EventResult{
			Messages:        []string{"You stumble into a colony of giant ants! They attack en masse!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  false,
		}
	})

	// e057 - Troll
	RegisterEvent("e057", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("Mountain Troll", 6, 16, 5)
		return EventResult{
			Messages:        []string{"A massive troll blocks the path! It raises its club with a roar!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 5,
		}
	})

	// e058 - Hidden cache
	RegisterEvent("e058", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(3, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find a hidden cache of %d gold buried under a tree.", gold)},
			GoldChange: gold,
		}
	})

	// e059 - Mysterious ruins (leads to search)
	RegisterEvent("e059", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You discover the crumbling walls of an ancient structure nearby."},
			Choices:  []string{"Investigate the ruins", "Continue on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					gold := TreasureRoll(4, Roll1d6())
					if Roll1d6() >= 5 {
						enemy := MakeEnemy("Ruins Undead", 4, 8, 3)
						return EventResult{
							Messages:        []string{"Ancient guardians rise to defend the ruins!"},
							CombatTriggered: true,
							Enemy:           &enemy,
							GoldChange:      gold,
						}
					}
					return EventResult{
						Messages:   []string{fmt.Sprintf("The ruins yield %d gold in ancient coins.", gold)},
						GoldChange: gold,
					}
				}
				return EventResult{Messages: []string{"You press on without investigating."}}
			},
		}
	})

	// e060 - Arrested
	RegisterEvent("e060", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"Local authorities arrest you on charges of trespassing!"},
			Choices:  []string{"Bribe the guard (20 gold)", "Fight your way out", "Submit to arrest"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if s.Gold >= 20 {
						s.Gold -= 20
						return EventResult{Messages: []string{"The guard takes your coin and releases you."}}
					}
					return EventResult{Messages: []string{"You don't have enough gold! You are imprisoned."},
						GameOver: false}
				case 1:
					enemy := MakeEnemy("Town Guard Captain", 5, 10, 2)
					return EventResult{
						Messages:        []string{"You fight your way free of the guards!"},
						CombatTriggered: true,
						Enemy:           &enemy,
					}
				default:
					// Imprisoned - lose 1d6 days
					days := Roll1d6()
					s.Day += days
					return EventResult{Messages: []string{fmt.Sprintf("You spend %d days in prison before escaping.", days)}}
				}
			},
		}
	})

	// e061 - Escaped prisoners
	RegisterEvent("e061", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d3()
		var followers []Character
		for i := 0; i < count; i++ {
			followers = append(followers, Character{
				Name:         fmt.Sprintf("Escapee %d", i+1),
				Type:         TypeEscapee,
				CombatSkill:  2 + Roll1d3(),
				MaxEndurance: 5 + Roll1d3(),
				DailyWage:    0,
				Morale:       3,
				IsEscapee:    true,
			})
		}
		if count == 0 {
			return EventResult{Messages: []string{"A group of escaped prisoners pass you by without stopping."}}
		}
		return EventResult{
			Messages: []string{fmt.Sprintf("%d escaped prisoners beg to join your group for protection.", count)},
			Choices:  []string{"Accept them", "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					for i := 1; i < len(followers); i++ {
						s.AddFollower(followers[i])
					}
					return EventResult{
						Messages:    []string{"The prisoners join your group, grateful for the protection."},
						NewFollower: &followers[0],
					}
				}
				return EventResult{Messages: []string{"You send them on their way."}}
			},
		}
	})

	// e062 - Abandoned farmstead
	RegisterEvent("e062", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d6() + 2
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find an abandoned farmstead with food stores. Gained %d food.", food)},
			FoodChange: food,
		}
	})

	// e063 - Wounded traveler
	RegisterEvent("e063", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A wounded traveler lies on the road, begging for help."},
			Choices:  []string{"Aid him (lose 2 food)", "Ignore him"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 && s.FoodUnits >= 2 {
					s.FoodUnits -= 2
					gold := Roll1d6() * 10
					s.Gold += gold
					return EventResult{Messages: []string{fmt.Sprintf("The traveler is grateful and shares his hidden gold: %d pieces.", gold)}}
				}
				return EventResult{Messages: []string{"You leave the traveler to his fate."}}
			},
		}
	})

	// e064 - Hidden Ruins: reveal an adjacent unvisited ruins hex, or pay gold if none
	RegisterEvent("e064", func(s *GameState, ctx EventContext) EventResult {
		hexID := s.CurrentHex
		for _, adj := range AdjacentHexes(hexID) {
			if s.VisitedHexes[adj] {
				continue
			}
			h := GetHex(adj)
			if h != nil && h.IsRuins() {
				s.VisitedHexes[adj] = true
				return EventResult{
					Messages: []string{
						fmt.Sprintf("You discover a map fragment revealing hidden ruins at hex %s!", adj),
						"Travel there and use [Search Ruins] from the action menu.",
					},
					Note: fmt.Sprintf("Hidden ruins at hex %s (day %d).", adj, s.Day),
				}
			}
		}
		// No adjacent unvisited ruins — fragment still has value
		gold := Roll1d6() * 5
		return EventResult{
			Messages: []string{
				"You find an old map fragment hinting at nearby ruins.",
				fmt.Sprintf("A scholar pays you %d gold for the information.", gold),
			},
			GoldChange: gold,
			Note:       fmt.Sprintf("Map fragment near hex %s (day %d).", hexID, s.Day),
		}
	})

	// e065 - Hidden Town: reveal a random adjacent unvisited settlement, or give gold if none
	RegisterEvent("e065", func(s *GameState, ctx EventContext) EventResult {
		hexID := s.CurrentHex
		// Search adjacent hexes for an unvisited settlement to reveal
		for _, adj := range AdjacentHexes(hexID) {
			if s.VisitedHexes[adj] {
				continue
			}
			h := GetHex(adj)
			if h != nil && h.IsSettlement() {
				s.VisitedHexes[adj] = true
				return EventResult{
					Messages: []string{
						fmt.Sprintf("Locals tell you of a hidden settlement at %s!", adj),
						"The inhabitants are wary of strangers but have gold to trade.",
					},
					Note: fmt.Sprintf("Hidden settlement at hex %s (day %d).", adj, s.Day),
				}
			}
		}
		// No adjacent hidden settlement — rumour still worth something
		gold := Roll1d6() * 10
		return EventResult{
			Messages: []string{
				fmt.Sprintf("Locals near hex %s speak of a hidden village somewhere to the south.", hexID),
				fmt.Sprintf("They press %d gold into your hand to carry word to distant traders.", gold),
			},
			GoldChange: gold,
			Note:       fmt.Sprintf("Rumour of hidden settlement near hex %s (day %d).", hexID, s.Day),
		}
	})

	// e066 - Secret Temple
	RegisterEvent("e066", func(s *GameState, ctx EventContext) EventResult {
		hexID := s.CurrentHex
		s.GetHexFlags(hexID).HiddenTemple = true
		return EventResult{
			Messages: []string{
				fmt.Sprintf("You discover the entrance to a hidden temple at hex %s, covered in vines.", hexID),
				"Ancient priests once performed powerful rituals here.",
				"You may enter now, or return to this hex later and use [Submit Offering] from the action menu.",
			},
			Note:    fmt.Sprintf("Hidden temple at hex %s (day %d).", hexID, s.Day),
			Choices: []string{"Enter the temple now", "Mark location and move on"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					gold := TreasureRoll(5, Roll1d6())
					if Roll1d6() >= 4 {
						enemy := MakeEnemy("Temple Guardian", 6, 12, 5)
						return EventResult{
							Messages:        []string{"Ancient guardians defend the temple!"},
							CombatTriggered: true,
							Enemy:           &enemy,
							GoldChange:      gold,
						}
					}
					return EventResult{
						Messages:   []string{fmt.Sprintf("The temple contains offerings worth %d gold!", gold)},
						GoldChange: gold,
					}
				}
				return EventResult{Messages: []string{
					fmt.Sprintf("You note the temple location (hex %s) for later.", hexID),
					"Return to this hex and choose [Submit Offering] to interact with it.",
				}}
			},
		}
	})

	// e067 - Ancient map
	RegisterEvent("e067", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You find an ancient map showing the layout of nearby ruins!"},
		}
	})

	// e068 - Wizard's Tower
	RegisterEvent("e068", func(s *GameState, ctx EventContext) EventResult {
		if s.Flags.WizardsTowerVisited[s.CurrentHex] {
			return EventResult{Messages: []string{"You pass by the wizard's tower. It appears deserted."}}
		}
		return EventResult{
			Messages: []string{"A tall tower stands alone on the hilltop.",
				"Magical light flickers in its windows."},
			Note:    fmt.Sprintf("Wizard's tower at hex %s (day %d).", s.CurrentHex, s.Day),
			Choices: []string{"Approach the tower", "Give it a wide berth"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				s.Flags.WizardsTowerVisited[s.CurrentHex] = true
				if choice == 0 {
					roll := Roll1d6()
					if roll <= 2 {
						enemy := MakeEnemy("Wizard", 7, 10, 7)
						return EventResult{
							Messages:        []string{"The wizard is hostile and attacks!"},
							CombatTriggered: true,
							Enemy:           &enemy,
						}
					} else if roll <= 4 {
						gold := TreasureRoll(6, Roll1d6())
						return EventResult{
							Messages:   []string{fmt.Sprintf("The wizard is friendly and pays you %d gold for news of the outside world.", gold)},
							GoldChange: gold,
						}
					}
					return EventResult{
						Messages: []string{"The wizard offers to enchant one of your weapons for 30 gold."},
						Choices:  []string{"Pay 30 gold", "Decline"},
						ChoiceHandler: func(s *GameState, c int) EventResult {
							if c == 0 && s.Gold >= 30 {
								s.Gold -= 30
								s.Prince.CombatSkill++
								return EventResult{Messages: []string{"Your weapon is imbued with magical power! +1 Combat Skill."}}
							}
							return EventResult{Messages: []string{"You decline the wizard's offer."}}
						},
					}
				}
				return EventResult{Messages: []string{"You give the tower a wide berth."}}
			},
		}
	})

	// e069 - Treasure hoard
	RegisterEvent("e069", func(s *GameState, ctx EventContext) EventResult {
		gold := TreasureRoll(5, Roll1d6())
		if Roll1d6() >= 5 {
			enemy := MakeEnemy("Treasure Guardian", 6, 14, 6)
			return EventResult{
				Messages:        []string{fmt.Sprintf("You find a treasure hoard worth %d gold, guarded by a monster!", gold)},
				CombatTriggered: true,
				Enemy:           &enemy,
				GoldChange:      gold,
			}
		}
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find an unguarded treasure cache worth %d gold!", gold)},
			GoldChange: gold,
		}
	})

	// e070 - Bear
	RegisterEvent("e070", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 4 {
			enemy := MakeEnemy("Forest Bear", 5, 14, 2)
			return EventResult{
				Messages:        []string{"A huge bear rears up and attacks!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  false,
			}
		}
		return EventResult{Messages: []string{"You spot a bear in the distance. It ignores you."}}
	})

	// e071 - Deer
	RegisterEvent("e071", func(s *GameState, ctx EventContext) EventResult {
		food := Roll1d6() + 2
		return EventResult{
			Messages:   []string{fmt.Sprintf("You successfully hunt a deer. Gained %d food.", food)},
			FoodChange: food,
		}
	})

	// e072 - Boar hunt
	RegisterEvent("e072", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6()+s.Prince.EffectiveCombatSkill() >= 8 {
			food := Roll1d6() + 3
			return EventResult{
				Messages:   []string{fmt.Sprintf("Successful hunt! Gained %d food units.", food)},
				FoodChange: food,
			}
		}
		return EventResult{Messages: []string{"The game is elusive today."}}
	})

	// e073 - Woodsman
	RegisterEvent("e073", func(s *GameState, ctx EventContext) EventResult {
		follower := Character{
			Name:         "Woodsman Guide",
			Type:         TypeGuide,
			CombatSkill:  3,
			MaxEndurance: 8,
			DailyWage:    2,
			Morale:       4,
			IsGuide:      true,
		}
		return EventResult{
			Messages: []string{"A skilled woodsman offers to guide you through the wilderness for 2 gold/day."},
			Choices:  []string{"Hire the woodsman", "Decline"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					return EventResult{NewFollower: &follower}
				}
				return EventResult{Messages: []string{"The woodsman tips his hat and disappears into the trees."}}
			},
		}
	})

	// e074 - Giant Spiders
	RegisterEvent("e074", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d3() + 1
		enemy := MakeEnemy(fmt.Sprintf("Giant Spiders (%d)", count), count+2, count*4, 2)
		return EventResult{
			Messages:        []string{fmt.Sprintf("%d giant spiders drop from the trees above you!", count)},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 5,
		}
	})

	// e075 - Wolf Pack
	RegisterEvent("e075", func(s *GameState, ctx EventContext) EventResult {
		count := Roll1d6() + 3
		enemy := MakeEnemy(fmt.Sprintf("Wolf Pack (%d)", count), count+1, count*3, 1)
		return EventResult{
			Messages:        []string{fmt.Sprintf("A pack of %d wolves surrounds you, growling hungrily!", count)},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 5,
		}
	})

	// e076 - Giant eagle
	RegisterEvent("e076", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 5 {
			enemy := MakeEnemy("Giant Eagle", 6, 12, 4)
			return EventResult{
				Messages:        []string{"A giant eagle dives at you, mistaking you for prey!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  false,
			}
		}
		return EventResult{Messages: []string{"A magnificent giant eagle soars overhead, a good omen."}}
	})

	// e077 - Treant
	RegisterEvent("e077", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 4 {
			enemy := MakeEnemy("Forest Treant", 7, 20, 6)
			return EventResult{
				Messages:        []string{"An ancient tree-being blocks your path! Its roots crack the earth around you!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  false,
			}
		}
		return EventResult{Messages: []string{"A forest spirit watches you pass. It seems benign."}}
	})

	// e078 - Bad Going
	RegisterEvent("e078", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"The terrain becomes treacherous. You make very slow progress.",
				"The extra effort costs you 1 food unit."},
			FoodChange: -1,
		}
	})

	// e079 - Heavy Rains
	RegisterEvent("e079", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"Heavy rains turn the ground to mud. Travel is extremely difficult.",
				"Your food stores are partly ruined by the damp."},
			FoodChange: -2,
		}
	})

	// e080 - Old hermit
	RegisterEvent("e080", func(s *GameState, ctx EventContext) EventResult {
		roll := Roll1d6()
		if roll <= 3 {
			return EventResult{
				Messages: []string{"An old hermit gives you cryptic advice about your quest.",
					"\"Seek the ancient crown where rivers meet\", he says."},
			}
		}
		if roll <= 5 {
			food := Roll1d6() + 2
			return EventResult{
				Messages:   []string{fmt.Sprintf("A hermit shares his simple provisions. Gained %d food.", food)},
				FoodChange: food,
			}
		}
		// Hermit is wise - heals wounds
		s.Prince.Wounds -= 2
		if s.Prince.Wounds < 0 {
			s.Prince.Wounds = 0
		}
		return EventResult{Messages: []string{"The hermit treats your wounds with herbal medicine. -2 wounds."}}
	})

	// e081 - Strange mushrooms
	RegisterEvent("e081", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You find unusual mushrooms. Are they edible?"},
			Choices:  []string{"Eat them", "Leave them"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					if Roll1d6() >= 4 {
						food := Roll1d6()
						return EventResult{
							Messages:   []string{fmt.Sprintf("The mushrooms are delicious and filling! Gained %d food.", food)},
							FoodChange: food,
						}
					}
					s.Prince.PoisonWounds++
					return EventResult{Messages: []string{"The mushrooms are poisonous! +1 poison wound."}}
				}
				return EventResult{Messages: []string{"You leave the mushrooms alone."}}
			},
		}
	})

	// e082 - Faerie circle
	RegisterEvent("e082", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You find a circle of glowing mushrooms. A faerie ring!"},
			Choices:  []string{"Step inside the ring", "Avoid it"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				if choice == 0 {
					roll := Roll1d6()
					switch {
					case roll <= 2:
						s.Prince.Wounds += 3
						return EventResult{Messages: []string{"The faeries curse you! +3 wounds."}}
					case roll <= 4:
						gold := Roll1d6() * 20
						return EventResult{
							Messages:   []string{fmt.Sprintf("The faeries reward your boldness with %d gold!", gold)},
							GoldChange: gold,
						}
					default:
						s.Prince.Wounds = 0
						return EventResult{Messages: []string{"The faerie magic heals all your wounds!"}}
					}
				}
				return EventResult{Messages: []string{"Wisely, you avoid the faerie ring."}}
			},
		}
	})

	// e083 - Forest shrine
	RegisterEvent("e083", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You find an ancient forest shrine draped in moss."},
			Choices:  []string{"Leave an offering (5 gold)", "Pray without offering", "Ignore it"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					if s.Gold >= 5 {
						s.Gold -= 5
						if Roll1d6() >= 3 {
							s.Prince.Wounds = 0
							return EventResult{Messages: []string{"The forest spirits heal your wounds completely!"}}
						}
						return EventResult{Messages: []string{"The forest seems to grow lighter around you. You feel refreshed."}}
					}
					return EventResult{Messages: []string{"You don't have enough gold to make an offering."}}
				case 1:
					if Roll1d6() == 6 {
						return EventResult{Messages: []string{"The spirits bless you with a vision of the path ahead."}}
					}
					return EventResult{Messages: []string{"Your prayer is met with silence."}}
				default:
					return EventResult{Messages: []string{"You ignore the shrine and press on."}}
				}
			},
		}
	})

	// e084 - Elven outpost
	RegisterEvent("e084", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You find an elven outpost. The elves watch you from the shadows.",
				"After a tense moment, they lower their bows."},
			GoldChange: 0,
			FoodChange: Roll1d3() + 1,
		}
	})

	// e085 - Narrow Ledge
	RegisterEvent("e085", func(s *GameState, ctx EventContext) EventResult {
		if Roll2d6() >= 8 {
			wounds := Roll1d3()
			s.Prince.Wounds += wounds
			return EventResult{
				Messages: []string{fmt.Sprintf("You lose your footing on a narrow ledge! %d wounds from the fall.", wounds)},
			}
		}
		return EventResult{
			Messages: []string{"You carefully navigate a treacherous narrow ledge and make it safely."},
		}
	})

	// e086 - High Pass
	RegisterEvent("e086", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"You find a high mountain pass. The thin air saps your strength."},
			FoodChange: -1,
		}
	})

	// e087 - Impassable
	RegisterEvent("e087", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages:     []string{"The way ahead is completely impassable! A sheer cliff blocks all progress.", "You must turn back and find another route."},
			BlocksTravel: true,
		}
	})

	// e088 - Rockslide
	RegisterEvent("e088", func(s *GameState, ctx EventContext) EventResult {
		if Roll2d6() >= 9 {
			wounds := Roll1d3() + 1
			s.Prince.Wounds += wounds
			return EventResult{
				Messages: []string{fmt.Sprintf("A rockslide crashes down around you! %d wounds from flying debris!", wounds)},
			}
		}
		return EventResult{
			Messages: []string{"You hear a rockslide in the distance. A near miss, but you are unharmed."},
		}
	})

	// e089 - Morass
	RegisterEvent("e089", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"The ground turns into a sucking morass. Extra effort is needed to push through.",
				"This cost you an additional food unit."},
			FoodChange: -1,
		}
	})

	// e090 - Quicksand
	RegisterEvent("e090", func(s *GameState, ctx EventContext) EventResult {
		if Roll2d6() >= 8 {
			return EventResult{
				Messages: []string{"You step into quicksand and sink up to your waist!", "By supreme effort you drag yourself free, exhausted."},
				FoodChange: -2,
			}
		}
		return EventResult{
			Messages: []string{"You spot quicksand patches and carefully navigate around them."},
		}
	})

	// e091 - Poison Snake
	RegisterEvent("e091", func(s *GameState, ctx EventContext) EventResult {
		if Roll2d6() >= 7 {
			s.Prince.PoisonWounds += 2
			return EventResult{
				Messages: []string{"A venomous snake strikes before you can react! +2 poison wounds.",
					"You must find an antidote soon or the poison will drain your strength."},
			}
		}
		return EventResult{
			Messages: []string{"A large snake rears up but you avoid it in time."},
		}
	})

	// e092 - Flood
	RegisterEvent("e092", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"Sudden flooding! The river bursts its banks, sweeping away supplies!"},
			FoodChange: -Roll1d3(),
			GoldChange: -Roll1d6() * 5,
		}
	})

	// e093 - Mountain lion
	RegisterEvent("e093", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("Mountain Lion", 5, 10, 2)
		return EventResult{
			Messages:        []string{"A mountain lion springs at you from above!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  false,
		}
	})

	// e094 - Crocodiles
	RegisterEvent("e094", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 4 {
			enemy := MakeEnemy("River Crocodile", 5, 16, 1)
			return EventResult{
				Messages:        []string{"Enormous crocodiles lunge from the river crossing!"},
				CombatTriggered: true,
				Enemy:           &enemy,
				PlayerAttFirst:  false,
			}
		}
		return EventResult{
			Messages: []string{"You spot crocodiles in the water ahead and find a safer crossing."},
		}
	})

	// e095 - Landslide
	RegisterEvent("e095", func(s *GameState, ctx EventContext) EventResult {
		return EventResult{
			Messages: []string{"A landslide blocks the path. You spend precious time finding a way around."},
		}
	})

	// e096 - Eagle attack
	RegisterEvent("e096", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("War Eagle", 4, 8, 2)
		return EventResult{
			Messages:        []string{"A massive eagle with razored talons dives at you!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 4,
		}
	})

	// e097 - Giant scorpion
	RegisterEvent("e097", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("Giant Scorpion", 5, 12, 2)
		result := EventResult{
			Messages:        []string{"A giant scorpion emerges from beneath a rock! Its stinger drips venom!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  Roll1d6() >= 5,
		}
		if Roll1d6() >= 5 {
			s.Prince.PoisonWounds++
			result.Messages = append(result.Messages, "Its tail scores a glancing blow! +1 poison wound.")
		}
		return result
	})

	// e098 - Dragon
	RegisterEvent("e098", func(s *GameState, ctx EventContext) EventResult {
		if s.Flags.DragonSlain[s.CurrentHex] {
			return EventResult{Messages: []string{"The dragon's cave is empty. You already slew it."}}
		}
		return EventResult{
			Messages: []string{
				"A DRAGON circles overhead then descends with a deafening roar!",
				"Fire scorches the earth around you. This will be the fight of your life!",
			},
			Choices: []string{"Fight the dragon!", "Flee!", "Try to reason with it"},
			ChoiceHandler: func(s *GameState, choice int) EventResult {
				switch choice {
				case 0:
					enemy := MakeEnemy("Fire Dragon", 9, 24, 9)
					return EventResult{
						Messages:        []string{"You draw your blade against the mighty dragon!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				case 1:
					if Roll1d6() >= 4 {
						return EventResult{
							Messages:   []string{"You flee! The dragon scorches the ground behind you but you escape!"},
							FoodChange: -2,
						}
					}
					enemy := MakeEnemy("Fire Dragon", 9, 24, 9)
					return EventResult{
						Messages:        []string{"The dragon catches you as you flee!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				default:
					if Roll1d6()+s.Prince.WitWiles >= 10 {
						gold := TreasureRoll(7, Roll1d6())
						return EventResult{
							Messages:   []string{fmt.Sprintf("The ancient dragon is amused by your boldness. It gifts you %d gold from its hoard!", gold)},
							GoldChange: gold,
						}
					}
					enemy := MakeEnemy("Fire Dragon", 9, 24, 9)
					return EventResult{
						Messages:        []string{"The dragon is not amused. It attacks!"},
						CombatTriggered: true,
						Enemy:           &enemy,
						PlayerAttFirst:  false,
					}
				}
			},
		}
	})

	// e099 - Roc
	RegisterEvent("e099", func(s *GameState, ctx EventContext) EventResult {
		enemy := MakeEnemy("Roc", 8, 20, 7)
		return EventResult{
			Messages:        []string{"The shadow of an enormous bird passes over you! A ROC has spotted you as prey!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  false,
		}
	})

	// e100 - Griffon
	RegisterEvent("e100", func(s *GameState, ctx EventContext) EventResult {
		if Roll1d6() >= 5 {
			// Could capture/tame griffon
			return EventResult{
				Messages: []string{"A magnificent griffon lands near you and eyes you curiously."},
				Choices:  []string{"Try to tame it (dangerous)", "Back away slowly"},
				ChoiceHandler: func(s *GameState, choice int) EventResult {
					if choice == 0 {
						if Roll1d6()+s.Prince.WitWiles >= 9 {
							s.Prince.HasMount = true
							s.Prince.MountType = MountPegasus // use pegasus type for griffon
							return EventResult{Messages: []string{"You tame the griffon! You now have a fearsome mount!"}}
						}
						enemy := MakeEnemy("Angry Griffon", 7, 18, 6)
						return EventResult{
							Messages:        []string{"The griffon attacks!"},
							CombatTriggered: true,
							Enemy:           &enemy,
						}
					}
					return EventResult{Messages: []string{"The griffon watches you leave, then takes flight."}}
				},
			}
		}
		enemy := MakeEnemy("Wild Griffon", 7, 18, 6)
		return EventResult{
			Messages:        []string{"A griffon dives from the sky, talons extended!"},
			CombatTriggered: true,
			Enemy:           &enemy,
			PlayerAttFirst:  false,
		}
	})
}
