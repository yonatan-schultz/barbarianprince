package game

import "testing"

func TestResolveCombatRound_EnemyDies(t *testing.T) {
	s := NewGameState()
	// Give prince overwhelming CS so enemy dies in one hit
	s.Prince.CombatSkill = 20
	enemy := MakeEnemy("Weakling", 1, 1, 3)

	_, over, result := ResolveCombatRound(s, &enemy, true)

	if !over {
		t.Error("combat should be over when enemy endurance reaches 0")
	}
	if !result.EnemyDead {
		t.Error("EnemyDead should be true")
	}
	if !result.PlayerWon {
		t.Error("PlayerWon should be true")
	}
	if result.LootGold <= 0 {
		t.Error("should have received loot gold")
	}
	if s.Gold <= 10 { // started with 10
		t.Error("gold should have increased from loot")
	}
}

func TestResolveCombatRound_WoundsApplied(t *testing.T) {
	s := NewGameState()
	s.Prince.CombatSkill = 1
	enemy := MakeEnemy("Brute", 20, 100, 2) // enemy will almost certainly wound prince

	woundsBefore := s.Prince.Wounds
	ResolveCombatRound(s, &enemy, false) // enemy attacks first

	// With CS 20 vs CS 1, enemy will almost certainly land wounds
	// Run 10 rounds to be statistically sure
	for i := 0; i < 10 && s.Prince.Wounds == woundsBefore; i++ {
		ResolveCombatRound(s, &enemy, false)
	}
	if s.Prince.Wounds == woundsBefore {
		t.Error("prince should have taken wounds against overwhelming enemy")
	}
}

func TestResolveCombatRound_PrinceDies(t *testing.T) {
	s := NewGameState()
	s.Prince.CombatSkill = 1
	s.Prince.MaxEndurance = 1
	s.Prince.Wounds = 0
	enemy := MakeEnemy("Executioner", 20, 100, 2)

	// Fight until prince dies (max 20 rounds to avoid infinite loop)
	for i := 0; i < 20; i++ {
		_, over, _ := ResolveCombatRound(s, &enemy, false)
		if over || s.Prince.IsDead() {
			break
		}
	}

	if !s.Prince.IsDead() {
		t.Error("prince with 1 endurance should die against overwhelming enemy")
	}
}

func TestAttemptFlee(t *testing.T) {
	s := NewGameState()
	enemy := MakeEnemy("Pursuer", 5, 10, 2)

	// Run many flee attempts; should sometimes succeed, sometimes fail
	successes := 0
	failures := 0
	for i := 0; i < 100; i++ {
		s2 := NewGameState()
		ok, msg := AttemptFlee(s2, &enemy)
		if msg == "" {
			t.Error("AttemptFlee should always return a message")
		}
		if ok {
			successes++
		} else {
			failures++
		}
	}
	_ = s

	// With threshold roll >=4 on 1d6, ~50% success rate
	if successes == 0 {
		t.Error("should sometimes succeed in fleeing")
	}
	if failures == 0 {
		t.Error("should sometimes fail to flee")
	}
}

func TestAttemptFlee_FailureCausesDamage(t *testing.T) {
	s := NewGameState()
	// Use overwhelming enemy so flee failure nearly always deals wounds
	enemy := MakeEnemy("Giant", 20, 100, 5)

	woundsBefore := s.Prince.Wounds
	damageOccurred := false
	for i := 0; i < 30; i++ {
		s2 := NewGameState()
		ok, _ := AttemptFlee(s2, &enemy)
		if !ok && s2.Prince.Wounds > woundsBefore {
			damageOccurred = true
			break
		}
	}
	if !damageOccurred {
		t.Error("failed flee attempt against overwhelming enemy should cause damage")
	}
}

func TestMakeEnemy(t *testing.T) {
	e := MakeEnemy("Dragon", 9, 24, 9)
	if e.Name != "Dragon" {
		t.Errorf("Name = %q, want Dragon", e.Name)
	}
	if e.CombatSkill != 9 {
		t.Errorf("CombatSkill = %d, want 9", e.CombatSkill)
	}
	if e.MaxEndurance != 24 {
		t.Errorf("MaxEndurance = %d, want 24", e.MaxEndurance)
	}
	if !e.IsAlive() {
		t.Error("fresh enemy should be alive")
	}
}

func TestCheckUnconsciousFollowers_CarriedByFollower(t *testing.T) {
	s := NewGameState()
	// Make prince unconscious: wounds == MaxEndurance-1
	s.Prince.MaxEndurance = 5
	s.Prince.Wounds = 4
	s.Gold = 50
	s.AddFollower(Character{Name: "Guard", CombatSkill: 4, MaxEndurance: 8, Morale: 5})

	// Run enough times that the carry case fires (Roll1d6 >= 4 = ~50%)
	carried := false
	for i := 0; i < 50; i++ {
		s2 := NewGameState()
		s2.Prince.MaxEndurance = 5
		s2.Prince.Wounds = 4
		s2.Gold = 50
		s2.AddFollower(Character{Name: "Guard", CombatSkill: 4, MaxEndurance: 8, Morale: 5})
		msgs := CheckUnconsciousFollowers(s2)
		if len(msgs) > 0 {
			carried = true
			break
		}
	}
	if !carried {
		t.Error("unconscious prince should sometimes be carried by followers")
	}
}

func TestCheckUnconsciousFollowers_NotTriggeredWhenConscious(t *testing.T) {
	s := NewGameState()
	s.Prince.MaxEndurance = 9
	s.Prince.Wounds = 0
	s.AddFollower(Character{Name: "Guard", CombatSkill: 4, MaxEndurance: 8})

	msgs := CheckUnconsciousFollowers(s)
	if len(msgs) != 0 {
		t.Errorf("conscious prince should not trigger unconscious check, got: %v", msgs)
	}
}
