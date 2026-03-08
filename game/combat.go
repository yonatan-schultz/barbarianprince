package game

import "fmt"

// CombatResult summarizes the outcome of a combat
type CombatResult struct {
	PlayerWon    bool
	PlayerFled   bool
	EnemyDead    bool
	EnemyFled    bool
	Messages     []string
	LootGold     int
	LootItem     PossessionType
}

// ResolveCombatRound executes one round of combat between prince (+party) vs enemy
// Returns messages describing what happened
func ResolveCombatRound(s *GameState, enemy *Character, playerAttacksFirst bool) (msgs []string, combatOver bool, result CombatResult) {
	// Compute party combat skill: prince + followers
	playerCS := s.TotalCombatSkill()
	// Holy Symbol gives +2 CS against undead enemies
	if s.Prince.HasPossession(PossHolySymbol) && enemy.IsUndead {
		playerCS += 2
	}

	// applyStrike resolves one attack, using explicit attacker CS to support party bonuses
	applyStrike := func(atkCS int, atk, def *Character) (int, int, string) {
		roll := Roll2d6()
		net := atkCS - def.EffectiveCombatSkill() + roll
		wounds := CombatWounds(net)
		def.Wounds += wounds
		msg := fmt.Sprintf("%s strikes %s (roll %d+%d=%d): %d wound(s)",
			atk.Name, def.Name,
			atkCS-def.EffectiveCombatSkill(), roll, net,
			wounds)
		return wounds, net, msg
	}

	playerStrike := func(def *Character) (int, int, string) {
		return applyStrike(playerCS, &s.Prince, def)
	}
	enemyStrike := func(def *Character) (int, int, string) {
		return applyStrike(enemy.EffectiveCombatSkill(), enemy, def)
	}

	resolveFirstStrike := func() (int, int, string) {
		if playerAttacksFirst {
			return playerStrike(enemy)
		}
		return enemyStrike(&s.Prince)
	}

	firstDef := enemy
	if !playerAttacksFirst {
		firstDef = &s.Prince
	}

	w1, net1, m1 := resolveFirstStrike()
	msgs = append(msgs, m1)

	// Magic Sword: extra wound on a high net roll when player attacks first
	if playerAttacksFirst && net1 >= 9 && s.Prince.HasPossession(PossMagicSword) {
		enemy.Wounds++
		msgs = append(msgs, "The Magic Sword bites deep! +1 extra wound.")
	}

	// Check if first defender is dead
	if firstDef.IsDead() {
		if firstDef == enemy {
			result.EnemyDead = true
			result.PlayerWon = true
			loot := TreasureRoll(enemy.WealthCode, Roll1d6())
			result.LootGold = loot
			s.Gold += loot
			msgs = append(msgs, fmt.Sprintf("%s is slain! You gain %d gold.", enemy.Name, loot))
		} else {
			msgs = append(msgs, "Cal Arath has been slain!")
		}
		combatOver = true
		return
	}

	_ = w1

	// Second strike (if first defender survives)
	var w2, net2 int
	var m2 string
	if playerAttacksFirst {
		w2, net2, m2 = enemyStrike(&s.Prince)
	} else {
		w2, net2, m2 = playerStrike(enemy)
		// Magic Sword extra wound on player's counter-attack
		if net2 >= 9 && s.Prince.HasPossession(PossMagicSword) {
			enemy.Wounds++
			m2 += " (Magic Sword bites deep! +1 extra wound)"
		}
	}
	msgs = append(msgs, m2)
	_ = w2

	secondDef := &s.Prince
	if playerAttacksFirst {
		secondDef = &s.Prince
	} else {
		secondDef = enemy
	}

	if secondDef.IsDead() {
		if secondDef == enemy {
			result.EnemyDead = true
			result.PlayerWon = true
			loot := TreasureRoll(enemy.WealthCode, Roll1d6())
			result.LootGold = loot
			s.Gold += loot
			msgs = append(msgs, fmt.Sprintf("%s is slain! You gain %d gold.", enemy.Name, loot))
		} else {
			msgs = append(msgs, "Cal Arath has been slain!")
		}
		combatOver = true
		return
	}

	// Rout check (r220f): after any round, enemy may flee — no loot when they rout
	if Roll1d6() == 6 {
		result.EnemyFled = true
		result.PlayerWon = true
		result.LootGold = 0 // cannot loot a routing enemy
		msgs = append(msgs, fmt.Sprintf("%s turns and flees! You cannot loot a fleeing foe.", enemy.Name))
		combatOver = true
		return
	}

	return msgs, combatOver, result
}

// AttemptFlee tries to escape from combat
// Returns true if escape is successful
func AttemptFlee(s *GameState, enemy *Character) (bool, string) {
	roll := Roll1d6()
	if roll >= 4 {
		return true, fmt.Sprintf("You manage to escape from %s!", enemy.Name)
	}
	// Enemy gets a free strike
	roll2 := Roll2d6()
	net := enemy.EffectiveCombatSkill() - s.Prince.EffectiveCombatSkill() + roll2
	wounds := CombatWounds(net)
	s.Prince.Wounds += wounds
	return false, fmt.Sprintf("Escape failed! %s strikes you for %d wound(s).", enemy.Name, wounds)
}

// SurpriseCheck determines who attacks first based on reference code
// Returns true if player attacks first
func SurpriseCheck(surpriseCode string) bool {
	roll := Roll1d6()
	switch surpriseCode {
	case "player_first":
		return true
	case "enemy_first":
		return false
	case "both_surprise":
		// Both sides surprised - defender goes first (enemy)
		return false
	default:
		// Normal: player goes first on 4+
		return roll >= 4
	}
}

// CheckUnconsciousFollowers handles r221b: when the prince falls unconscious,
// each follower rolls 1d6: 4+ = carries him; 1-3 = all desert and steal everything.
// Returns messages. Call this after any combat round where the prince may be unconscious.
func CheckUnconsciousFollowers(s *GameState) []string {
	if !s.Prince.IsUnconscious() || len(s.Party) == 0 {
		return nil
	}
	var msgs []string
	msgs = append(msgs, "Cal Arath is unconscious!")
	roll := Roll1d6()
	if roll >= 4 {
		msgs = append(msgs, "Your followers rally and carry you to safety.")
		// Prince's CS becomes 0 while unconscious (handled via IsUnconscious in EffectiveCombatSkill)
	} else {
		msgs = append(msgs, "Your followers abandon you, taking everything!")
		s.Gold = 0
		s.Prince.Possessions = nil
		s.Party = nil
	}
	return msgs
}

// MakeEnemy creates an enemy character for combat
func MakeEnemy(name string, cs, endurance, wealthCode int) Character {
	return Character{
		Name:         name,
		Type:         TypeGeneric,
		CombatSkill:  cs,
		MaxEndurance: endurance,
		WealthCode:   wealthCode,
	}
}
