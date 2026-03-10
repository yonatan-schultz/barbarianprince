package game

import "fmt"

func registerEventsE180() {
	// e180 - Ring of Command
	RegisterEvent("e180", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossRingOfCommand) {
			s.Prince.AddPossession(PossRingOfCommand)
			return EventResult{
				Messages: []string{"You find a golden ring engraved with a command rune!",
					"The Ring of Command grants +2 Combat Skill in all battles!"},
			}
		}
		gold := 80
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find valuable jewelry worth %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e181 - Amulet of Power
	RegisterEvent("e181", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossAmuletOfPower) {
			s.Prince.AddPossession(PossAmuletOfPower)
			return EventResult{
				Messages: []string{"You find a glowing amulet of ancient power!",
					"The Amulet of Power grants +1 Combat Skill!"},
			}
		}
		gold := 50
		return EventResult{
			Messages:   []string{fmt.Sprintf("The amulet is a replica worth %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e182 - Elven Boots
	RegisterEvent("e182", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossElvenBoots) {
			s.Prince.AddPossession(PossElvenBoots)
			return EventResult{
				Messages: []string{"A pair of supple green boots with intricate leaf designs!",
					"Elven Boots: You never get lost in forest terrain!"},
			}
		}
		return EventResult{Messages: []string{"You find ordinary boots. Not worth taking."}}
	})

	// e183 - Rope and Grapnel
	RegisterEvent("e183", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossRopeAndGrapnel) {
			s.Prince.AddPossession(PossRopeAndGrapnel)
			return EventResult{
				Messages: []string{"You find a strong rope with a well-forged grapnel.",
					"Rope & Grapnel: Can scale cliffs and cross one extra mountain hex per day."},
			}
		}
		return EventResult{Messages: []string{"An extra rope. You already have one."}}
	})

	// e184 - Lantern
	RegisterEvent("e184", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossLantern) {
			s.Prince.AddPossession(PossLantern)
			return EventResult{
				Messages: []string{"You find a brass lantern that never needs oil.",
					"Lantern: Bonus when searching ruins."},
			}
		}
		return EventResult{Messages: []string{"An ordinary lantern. You already have a better one."}}
	})

	// e185 - Poison Antidote
	RegisterEvent("e185", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossPoisonAntidote) {
			s.Prince.AddPossession(PossPoisonAntidote)
			return EventResult{
				Messages: []string{"You find a vial of universal antidote.",
					"Poison Antidote: Cures all poison wounds when used."},
			}
		}
		// Use the antidote if already poisoned
		if s.Prince.PoisonWounds > 0 {
			s.Prince.PoisonWounds = 0
			return EventResult{Messages: []string{"You use the antidote to cure your poison wounds!"}}
		}
		gold := 30
		return EventResult{
			Messages:   []string{fmt.Sprintf("You find a spare antidote and sell it for %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e186 - Healing Potion
	RegisterEvent("e186", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossHealingPotion) {
			s.Prince.AddPossession(PossHealingPotion)
			return EventResult{
				Messages: []string{"You find a gleaming red potion in a crystal vial.",
					"Healing Potion: Restores 1d6 endurance when used."},
			}
		}
		// Use it immediately
		healed := Roll1d6()
		s.Prince.Wounds -= healed
		if s.Prince.Wounds < 0 {
			s.Prince.Wounds = 0
		}
		return EventResult{Messages: []string{fmt.Sprintf("You find and drink a healing potion! +%d endurance restored.", healed)}}
	})

	// e187 - Magic Sword
	RegisterEvent("e187", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossMagicSword) {
			s.Prince.AddPossession(PossMagicSword)
			return EventResult{
				Messages: []string{"A sword glowing with blue runes rests on the stone!",
					"Magic Sword: +2 Combat Skill, and any roll of 9+ causes an extra wound."},
			}
		}
		gold := 100
		return EventResult{
			Messages:   []string{fmt.Sprintf("A fine enchanted blade. You already have one, so you sell it for %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e188 - Holy Symbol
	RegisterEvent("e188", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossHolySymbol) {
			s.Prince.AddPossession(PossHolySymbol)
			return EventResult{
				Messages: []string{"An ancient holy symbol carved from pure white stone.",
					"Holy Symbol: Repels undead in ruins; +1 CS against undead enemies."},
			}
		}
		return EventResult{Messages: []string{"Another holy symbol. The gods must favor you."}}
	})

	// e189 - Noble Parchment
	RegisterEvent("e189", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossNobleParchment) {
			s.Prince.AddPossession(PossNobleParchment)
			s.Flags.NobleAllySecured = true
			return EventResult{
				Messages: []string{"You find a sealed letter bearing a noble crest!",
					"Noble Parchment: Evidence of noble support for your claim.",
					"Return north with this document and your throne is within reach!"},
			}
		}
		gold := 40
		return EventResult{
			Messages:   []string{fmt.Sprintf("An old document. Worth %d gold as curiosity.", gold)},
			GoldChange: gold,
		}
	})

	// e190 - Ancient Map
	RegisterEvent("e190", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossMap) {
			s.Prince.AddPossession(PossMap)
			return EventResult{
				Messages: []string{"An oilskin map of unknown lands!",
					"Ancient Map: Reveals adjacent hex terrain when you enter a new area."},
			}
		}
		gold := 20
		return EventResult{
			Messages:   []string{fmt.Sprintf("A duplicate map. Worth %d gold to a collector.", gold)},
			GoldChange: gold,
		}
	})

	// e191 - Raft
	RegisterEvent("e191", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossRaft) {
			s.Prince.AddPossession(PossRaft)
			return EventResult{
				Messages: []string{"A serviceable raft, lashed to the riverbank.",
					"Raft: Allows river travel downstream at 3 hexsides per day."},
			}
		}
		return EventResult{Messages: []string{"Another raft. You already have one."}}
	})

	// e192 - Golden Crown
	RegisterEvent("e192", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossGoldenCrown) && !s.Flags.HasRoyalHelm {
			s.Prince.AddPossession(PossGoldenCrown)
			return EventResult{
				Messages: []string{"A golden crown set with ancient gems!",
					"Golden Crown: May substitute for the Royal Helm to win the game."},
			}
		}
		gold := TreasureRoll(6, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("A beautiful golden crown worth %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e193 - Royal Helm (also reachable via ruins search e138)
	RegisterEvent("e193", func(s *GameState, ctx EventContext) EventResult {
		if !s.Flags.HasRoyalHelm {
			s.Prince.AddPossession(PossRoyalHelm)
			s.Flags.HasRoyalHelm = true
			return EventResult{
				Messages: []string{"THE ROYAL HELM OF CAL ARATH!",
					"The lost crown of your dynasty, gleaming in the darkness!",
					"Return this to Ogon or Weshor and your throne is restored!"},
			}
		}
		gold := TreasureRoll(7, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("Ancient treasury: %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e194 - Staff of Command (also reachable via temple e159)
	RegisterEvent("e194", func(s *GameState, ctx EventContext) EventResult {
		if !s.Flags.HasStaffOfCommand {
			s.Prince.AddPossession(PossStaffOfCommand)
			s.Flags.HasStaffOfCommand = true
			return EventResult{
				Messages: []string{"THE STAFF OF COMMAND!",
					"This legendary artifact bends armies to its wielder's will.",
					"With the Staff, you can march north and reclaim your throne!"},
			}
		}
		gold := TreasureRoll(7, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("Powerful enchanted items worth %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e195 - True Love (r228): a noble companion joins the party permanently.
	// Never deserts, +1 W&W, and guarantees rescue if the prince falls unconscious.
	// One-shot: if already met, falls back to a minor gold reward.
	RegisterEvent("e195", func(s *GameState, ctx EventContext) EventResult {
		if s.Flags.TrueLoveMet {
			// Already met — a generous patron instead
			gold := Roll2d6() * 3
			return EventResult{
				Messages:   []string{fmt.Sprintf("A generous admirer rewards your growing reputation: %d gold.", gold)},
				GoldChange: gold,
			}
		}
		s.Flags.TrueLoveMet = true
		names := []string{"Lady Mira", "Lady Sora", "Lady Elara", "Lady Vessa"}
		name := names[Roll1d6()%len(names)]
		tl := Character{
			Name:         name,
			Type:         TypeGeneric,
			CombatSkill:  2,
			MaxEndurance: 6,
			Morale:       6,
			IsTrueLove:   true,
			DailyWage:    0,
		}
		s.AddFollower(tl)
		return EventResult{
			Messages: []string{
				fmt.Sprintf("%s — a noble driven from her home by the usurper — approaches you.", name),
				"Her eyes hold both grief and steel. She asks to stand at your side.",
				fmt.Sprintf("%s joins your party. (CS 2, E 6, no wage) Never deserts. +1 W&W.", name),
			},
			Note: fmt.Sprintf("%s joined at %s (day %d). True Love: never deserts, +1 W&W, guarantees rescue.", name, s.CurrentHex, s.Day),
		}
	})

	// e196 - Alcove of Sending
	RegisterEvent("e196", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossAlcoveOfSending) {
			s.Prince.AddPossession(PossAlcoveOfSending)
			return EventResult{
				Messages: []string{
					"Built into the wall is a shimmering alcove carved with sending-runes.",
					"Alcove of Sending: Deposit gold for safekeeping; withdraw it any time.",
				},
			}
		}
		gold := Roll2d6() * 10
		return EventResult{
			Messages:   []string{fmt.Sprintf("The magical alcove holds %d gold abandoned by a previous traveller.", gold)},
			GoldChange: gold,
		}
	})

	// e197 - Arch of Travel
	RegisterEvent("e197", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossArchOfTravel) {
			s.Prince.AddPossession(PossArchOfTravel)
			return EventResult{
				Messages: []string{
					"A stone archway crackles with teleportation magic!",
					"Arch of Travel: One-use item — teleport instantly to any major settlement.",
				},
			}
		}
		gold := TreasureRoll(5, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("Magical residue in the arch — worth %d gold to an enchanter.", gold)},
			GoldChange: gold,
		}
	})

	// e198 - Gateway to Darkness
	RegisterEvent("e198", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossGatewayToDarkness) {
			s.Prince.AddPossession(PossGatewayToDarkness)
			return EventResult{
				Messages: []string{
					"A dark portal hums with baleful energy. A daemon waits on the other side.",
					"Gateway to Darkness: Step through to face the Daemon and claim its hoard.",
					"Warning: The Daemon is powerful. The gateway is consumed when used.",
				},
			}
		}
		gold := TreasureRoll(4, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("Dark crystal fragments from a sealed gateway — worth %d gold.", gold)},
			GoldChange: gold,
		}
	})

	// e199 - Mirror of Reversal
	RegisterEvent("e199", func(s *GameState, ctx EventContext) EventResult {
		if !s.Prince.HasPossession(PossMirrorOfReversal) {
			s.Prince.AddPossession(PossMirrorOfReversal)
			return EventResult{
				Messages: []string{
					"A silver mirror inscribed with reversal runes leans against the wall.",
					"Mirror of Reversal: One-use item — instantly heals all your wounds.",
				},
			}
		}
		// Already have one — use it immediately
		healed := s.Prince.Wounds + s.Prince.PoisonWounds
		if healed > 0 {
			s.Prince.Wounds = 0
			s.Prince.PoisonWounds = 0
			return EventResult{
				Messages: []string{fmt.Sprintf("A second mirror shatters as you touch it, reversing your %d wound(s)!", healed)},
			}
		}
		gold := TreasureRoll(3, Roll1d6())
		return EventResult{
			Messages:   []string{fmt.Sprintf("A cracked mirror. Worth %d gold as scrap silver.", gold)},
			GoldChange: gold,
		}
	})
}
