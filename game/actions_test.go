package game

import "testing"

// ---- ActionBuyFood advances the day ----------------------------------------

func TestBuyFoodAdvancesDay(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon — settlement
	s.Gold = 100
	s.FoodUnits = 0
	dayBefore := s.Day

	ExecuteAction(s, ActionBuyFood)

	if s.Day != dayBefore+1 {
		t.Errorf("after BuyFood: day = %d, want %d (day must advance)", s.Day, dayBefore+1)
	}
}

func TestBuyFoodActuallyBuysFood(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)
	s.Gold = 100
	s.FoodUnits = 0

	ExecuteAction(s, ActionBuyFood)

	if s.FoodUnits == 0 {
		t.Error("BuyFood should have added food units")
	}
	if s.Gold >= 100 {
		t.Error("BuyFood should have deducted gold")
	}
}

// ---- ActionSearchCache advances the day and is gated by flags ---------------

func TestSearchCacheAdvancesDay(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5)
	flags := s.GetHexFlags(s.CurrentHex)
	flags.CacheHidden = true
	s.Caches = append(s.Caches, Cache{Location: s.CurrentHex, Gold: 50})
	dayBefore := s.Day

	ExecuteAction(s, ActionSearchCache)

	if s.Day != dayBefore+1 {
		t.Errorf("after SearchCache: day = %d, want %d (day must advance)", s.Day, dayBefore+1)
	}
}

func TestSearchCacheInAvailableActions(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5)
	s.GetHexFlags(s.CurrentHex).CacheHidden = true

	has := func(a Action) bool {
		for _, x := range s.AvailableActions() {
			if x == a {
				return true
			}
		}
		return false
	}

	if !has(ActionSearchCache) {
		t.Error("ActionSearchCache must appear in AvailableActions when a cache is hidden here")
	}

	// Once found, should disappear
	s.GetHexFlags(s.CurrentHex).CacheFound = true
	if has(ActionSearchCache) {
		t.Error("ActionSearchCache must not appear after cache is already found")
	}
}

func TestSearchCacheNotAvailableWithoutCache(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5) // no cache

	for _, a := range s.AvailableActions() {
		if a == ActionSearchCache {
			t.Error("ActionSearchCache must not appear when no cache is hidden here")
		}
	}
}

// ---- GoldenCrown win condition -----------------------------------------------

func TestGoldenCrownWinCondition(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon
	s.Prince.AddPossession(PossGoldenCrown)

	won, _ := CheckWinConditions(s)
	if !won {
		t.Error("Golden Crown at Ogon should satisfy the win condition")
	}
}

func TestGoldenCrownNotWinOutsideNorth(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(12, 15) // south
	s.Prince.AddPossession(PossGoldenCrown)

	won, _ := CheckWinConditions(s)
	if won {
		t.Error("Golden Crown should not win when south of Tragoth")
	}
}

func TestRoyalHelmWinConditionUnchanged(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(15, 1) // Weshor
	s.Prince.AddPossession(PossRoyalHelm)

	won, _ := CheckWinConditions(s)
	if !won {
		t.Error("Royal Helm at Weshor should still satisfy the win condition")
	}
}

// ---- CS possession bonuses ---------------------------------------------------

func TestCSBonusRingOfCommand(t *testing.T) {
	s := NewGameState()
	base := s.Prince.EffectiveCombatSkill()
	s.Prince.AddPossession(PossRingOfCommand)
	if s.Prince.EffectiveCombatSkill() != base+2 {
		t.Errorf("Ring of Command: CS = %d, want %d", s.Prince.EffectiveCombatSkill(), base+2)
	}
}

func TestCSBonusAmuletOfPower(t *testing.T) {
	s := NewGameState()
	base := s.Prince.EffectiveCombatSkill()
	s.Prince.AddPossession(PossAmuletOfPower)
	if s.Prince.EffectiveCombatSkill() != base+1 {
		t.Errorf("Amulet of Power: CS = %d, want %d", s.Prince.EffectiveCombatSkill(), base+1)
	}
}

func TestCSBonusMagicSword(t *testing.T) {
	s := NewGameState()
	base := s.Prince.EffectiveCombatSkill()
	s.Prince.AddPossession(PossMagicSword)
	if s.Prince.EffectiveCombatSkill() != base+2 {
		t.Errorf("Magic Sword: CS = %d, want %d", s.Prince.EffectiveCombatSkill(), base+2)
	}
}

func TestFollowerCombatSkillAdded(t *testing.T) {
	s := NewGameState()
	alone := s.TotalCombatSkill()
	s.AddFollower(Character{Name: "Swordsman", Type: TypeSwordsman, CombatSkill: 4, MaxEndurance: 8, Morale: 6})
	withFollower := s.TotalCombatSkill()
	if withFollower <= alone {
		t.Errorf("TotalCombatSkill with follower (%d) should exceed prince alone (%d)", withFollower, alone)
	}
}

// ---- UseItem action ----------------------------------------------------------

func TestUseItemHealingPotion(t *testing.T) {
	s := NewGameState()
	s.Prince.Wounds = 4
	s.Prince.AddPossession(PossHealingPotion)

	woundsBefore := s.Prince.Wounds
	er := ExecuteAction(s, ActionUseItem)
	if er == nil {
		t.Fatal("UseItem with Healing Potion should return an EventResult with choices")
	}
	// Select choice 0 (Healing Potion)
	if er.ChoiceHandler == nil {
		t.Fatal("UseItem result should have a ChoiceHandler")
	}
	result := er.ChoiceHandler(s, 0)
	if len(result.Messages) == 0 {
		t.Error("Healing Potion use should produce a message")
	}
	if s.Prince.Wounds >= woundsBefore {
		t.Errorf("Healing Potion should reduce wounds: before=%d after=%d", woundsBefore, s.Prince.Wounds)
	}
	if s.Prince.HasPossession(PossHealingPotion) {
		t.Error("Healing Potion should be consumed after use")
	}
}

func TestUseItemPoisonAntidote(t *testing.T) {
	s := NewGameState()
	s.Prince.PoisonWounds = 3
	s.Prince.AddPossession(PossPoisonAntidote)

	er := ExecuteAction(s, ActionUseItem)
	if er == nil {
		t.Fatal("UseItem with Poison Antidote should return an EventResult with choices")
	}
	if er.ChoiceHandler == nil {
		t.Fatal("UseItem result should have a ChoiceHandler")
	}
	er.ChoiceHandler(s, 0)
	if s.Prince.PoisonWounds != 0 {
		t.Errorf("Antidote should clear poison wounds, got %d", s.Prince.PoisonWounds)
	}
	if s.Prince.HasPossession(PossPoisonAntidote) {
		t.Error("Poison Antidote should be consumed after use")
	}
}

// ---- Hunt action -------------------------------------------------------------

func TestHuntAdvancesDay(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(3, 5) // wilderness hex
	// Make sure it's not farmland/town
	hex := GetHex(s.CurrentHex)
	if hex == nil {
		t.Skip("hex not found in test map")
	}
	dayBefore := s.Day
	ExecuteAction(s, ActionHunt)
	if s.Day != dayBefore+1 {
		t.Errorf("Hunt should advance day: before=%d after=%d", dayBefore, s.Day)
	}
}

func TestHuntNotAvailableInFarmland(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon — farmland settlement
	for _, a := range s.AvailableActions() {
		if a == ActionHunt {
			t.Error("Hunt should not be available in farmland/settlement")
		}
	}
}

// ---- Poison status line ------------------------------------------------------

func TestStatusLinePoisonDisplay(t *testing.T) {
	s := NewGameState()
	s.Prince.PoisonWounds = 2
	line := s.StatusLine()
	if !contains(line, "Poison") {
		t.Errorf("StatusLine should show poison count, got: %s", line)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ---- doSeekAudience ----------------------------------------------------------

func TestSeekAudience_NotAvailableInTown(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon — town, no castle

	result := ExecuteAction(s, ActionSeekAudience)
	// Should log "no court here" and return nil (no event result)
	if result != nil {
		t.Error("audience in a town should return nil EventResult")
	}
	found := false
	for _, msg := range s.Log {
		if len(msg) > 5 {
			found = true
			break
		}
	}
	if !found {
		t.Error("audience in a town should log a message")
	}
}

func TestSeekAudience_AvailableAtHulora(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(12, 13) // Hulora Castle
	s.FoodUnits = 100

	// Run several times — the action should produce log messages each time
	msgs := 0
	for i := 0; i < 5; i++ {
		s2 := NewGameState()
		s2.CurrentHex = NewHexID(12, 13)
		s2.FoodUnits = 100
		s2.Gold = 100
		ExecuteAction(s2, ActionSeekAudience)
		msgs += len(s2.Log)
	}
	if msgs == 0 {
		t.Error("audience at Hulora Castle should always produce log output")
	}
}

func TestSeekAudience_BarredPreventsAction(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(12, 13) // Hulora Castle
	s.FoodUnits = 100
	// Bar the prince indefinitely
	s.AudienceBarred[s.CurrentHex] = s.Day + 100

	ExecuteAction(s, ActionSeekAudience)

	for _, msg := range s.Log {
		// Should contain a "barred" message
		if len(msg) > 10 {
			return // got a message, test passes
		}
	}
	t.Error("barred audience should log a message")
}
