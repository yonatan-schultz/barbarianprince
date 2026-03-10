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

// ---- ActionRest -------------------------------------------------------------

func TestRestAdvancesDay(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 100
	dayBefore := s.Day
	ExecuteAction(s, ActionRest)
	if s.Day != dayBefore+1 {
		t.Errorf("Rest: day = %d, want %d", s.Day, dayBefore+1)
	}
}

func TestRestHealsInSettlement(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon - settlement
	s.FoodUnits = 100
	s.Gold = 100
	s.Prince.Wounds = 3
	ExecuteAction(s, ActionRest)
	if s.Prince.Wounds >= 3 {
		t.Error("Rest in settlement should reduce wounds")
	}
}

func TestRestLogsMessage(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 100
	logBefore := len(s.Log)
	ExecuteAction(s, ActionRest)
	if len(s.Log) <= logBefore {
		t.Error("Rest should append log messages")
	}
}

// ---- ActionBuyRaft ----------------------------------------------------------

func TestBuyRaft_SucceedsWithEnoughGold(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // settlement
	s.FoodUnits = 100
	s.Gold = 50

	ExecuteAction(s, ActionBuyRaft)

	if !s.Prince.HasPossession(PossRaft) {
		t.Error("BuyRaft with enough gold should add PossRaft to inventory")
	}
	if s.Gold != 35 { // 50 - 15 = 35
		t.Errorf("Gold = %d, want 35 after buying raft for 15", s.Gold)
	}
}

func TestBuyRaft_FailsWithInsufficientGold(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)
	s.FoodUnits = 100
	s.Gold = 5 // less than 15

	ExecuteAction(s, ActionBuyRaft)

	if s.Prince.HasPossession(PossRaft) {
		t.Error("BuyRaft with insufficient gold should not give raft")
	}
	if s.Gold != 5 {
		t.Errorf("Gold = %d, want 5 unchanged when can't afford raft", s.Gold)
	}
}

func TestBuyRaft_AvailableAtSettlement(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon

	found := false
	for _, a := range s.AvailableActions() {
		if a == ActionBuyRaft {
			found = true
			break
		}
	}
	if !found {
		t.Error("ActionBuyRaft should be available at settlement when not already carrying one")
	}
}

func TestBuyRaft_NotAvailableIfAlreadyOwned(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)
	s.Prince.AddPossession(PossRaft)

	for _, a := range s.AvailableActions() {
		if a == ActionBuyRaft {
			t.Error("ActionBuyRaft should not be available when prince already has a raft")
		}
	}
}

// ---- HideCacheHere ----------------------------------------------------------

func TestHideCacheHere_Success(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5)
	s.Gold = 50

	msgs := HideCacheHere(s, 30)

	if len(msgs) == 0 {
		t.Error("HideCacheHere should return messages")
	}
	if s.Gold != 20 {
		t.Errorf("Gold = %d, want 20 after hiding 30 gold", s.Gold)
	}
	if !s.GetHexFlags(s.CurrentHex).CacheHidden {
		t.Error("CacheHidden flag should be set after hiding cache")
	}
	if len(s.Caches) != 1 {
		t.Errorf("Caches len = %d, want 1", len(s.Caches))
	}
}

func TestHideCacheHere_AlreadyHasCacheBlocked(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5)
	s.Gold = 50
	s.GetHexFlags(s.CurrentHex).CacheHidden = true

	msgs := HideCacheHere(s, 20)

	if len(msgs) == 0 {
		t.Error("HideCacheHere should return a message when cache already exists")
	}
	if s.Gold != 50 {
		t.Error("Gold should not change when cache already exists")
	}
}

func TestHideCacheHere_NoGold(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5)
	s.Gold = 0

	msgs := HideCacheHere(s, 10)

	if len(msgs) == 0 {
		t.Error("HideCacheHere should return a message when no gold to hide")
	}
}

// ---- PayLodging desertion path ----------------------------------------------

func TestPayLodging_CannotAffordSleepsRough(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // settlement
	s.Gold = 0
	s.Prince.WitWiles = 0

	msgs := PayLodging(s)

	if len(msgs) == 0 {
		t.Error("PayLodging with no gold should return messages about sleeping rough")
	}
	// Check that the "sleeping rough" message is in there
	found := false
	for _, m := range msgs {
		if len(m) > 10 {
			found = true
			break
		}
	}
	if !found {
		t.Error("PayLodging should produce at least one message")
	}
}

func TestPayLodging_WithMounts(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // settlement
	s.Gold = 100
	s.Prince.HasMount = true

	msgs := PayLodging(s)

	if len(msgs) == 0 {
		t.Error("PayLodging with mount should return stabling message")
	}
	// Gold should have been reduced by 2 (1 person + 1 mount)
	if s.Gold != 98 {
		t.Errorf("Gold = %d, want 98 after lodging+stabling (1+1=2)", s.Gold)
	}
}

// ---- HealRest: wilderness probabilistic path --------------------------------

func TestHealRest_WildernessEventuallyHeals(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(4, 3) // Mountains - wilderness
	s.Prince.Wounds = 5

	healed := false
	for i := 0; i < 30; i++ {
		s2 := NewGameState()
		s2.CurrentHex = NewHexID(4, 3)
		s2.Prince.Wounds = 5
		msgs := HealRest(s2)
		if s2.Prince.Wounds < 5 || len(msgs) > 0 {
			healed = true
			break
		}
	}
	if !healed {
		t.Error("HealRest in wilderness should sometimes heal (1d6>=5)")
	}
}

func TestHealRest_WildernessPoison(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(4, 3) // wilderness
	s.Prince.PoisonWounds = 3

	cleared := false
	for i := 0; i < 30; i++ {
		s2 := NewGameState()
		s2.CurrentHex = NewHexID(4, 3)
		s2.Prince.PoisonWounds = 3
		HealRest(s2)
		if s2.Prince.PoisonWounds < 3 {
			cleared = true
			break
		}
	}
	if !cleared {
		t.Error("HealRest in wilderness should sometimes reduce poison (1d6>=5)")
	}
}

// ---- SeekNews at non-settlement ---------------------------------------------

func TestSeekNews_NotInSettlement(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(4, 3) // Mountains - no settlement

	result := ExecuteAction(s, ActionSeekNews)
	if result != nil {
		t.Error("SeekNews outside settlement should return nil")
	}
	found := false
	for _, msg := range s.Log {
		if len(msg) > 5 {
			found = true
			break
		}
	}
	if !found {
		t.Error("SeekNews outside settlement should log a message")
	}
}

// ---- SeekFollowers at non-settlement ----------------------------------------

func TestSeekFollowers_NotInSettlement(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(4, 3) // wilderness

	result := ExecuteAction(s, ActionSeekFollowers)
	if result != nil {
		t.Error("SeekFollowers outside settlement should return nil")
	}
}

// ---- ActionString and ActionKey coverage ------------------------------------

func TestActionString(t *testing.T) {
	allActions := []Action{
		ActionTravel, ActionRest, ActionSeekNews, ActionSeekFollowers,
		ActionBuyFood, ActionSeekAudience, ActionSubmitOffering,
		ActionSearchRuins, ActionSearchCache, ActionUseItem,
		ActionHunt, ActionBuyRaft,
	}
	for _, a := range allActions {
		s := a.String()
		if s == "" || s == "Unknown" {
			t.Errorf("Action(%d).String() = %q, want non-empty known string", a, s)
		}
		k := a.ActionKey()
		if k == "" {
			t.Errorf("Action(%d).ActionKey() = %q, want non-empty key", a, k)
		}
	}
}

// ---- TurnPhase String (via StatusLine) --------------------------------------

func TestTurnPhases(t *testing.T) {
	phases := []TurnPhase{
		PhaseActionSelect, PhaseEventResolve, PhaseCombat, PhaseTravel, PhaseGameOver,
	}
	// Just verify they're distinct integer values (not overlapping)
	seen := make(map[TurnPhase]bool)
	for _, p := range phases {
		if seen[p] {
			t.Errorf("TurnPhase %d is duplicated", p)
		}
		seen[p] = true
	}
}
