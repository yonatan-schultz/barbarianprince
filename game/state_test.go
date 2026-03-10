package game

import (
	"fmt"
	"testing"
)

func TestAdvanceDay_FoodConsumed(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 20
	s.Day = 1

	AdvanceDay(s)

	// Prince alone = 1 food/day
	if s.FoodUnits != 19 {
		t.Errorf("food = %d, want 19 after one day alone", s.FoodUnits)
	}
	if s.Day != 2 {
		t.Errorf("day = %d, want 2", s.Day)
	}
}

func TestAdvanceDay_WithMount(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 20
	s.Prince.HasMount = true
	s.Prince.MountType = MountHorse
	// Use mountains hex — mounts cannot forage there (r215f)
	s.CurrentHex = NewHexID(4, 3) // Mountains

	AdvanceDay(s)

	// Prince (1) + horse (2) = 3 food/day in mountains
	if s.FoodUnits != 17 {
		t.Errorf("food = %d, want 17 (prince + horse in mountains)", s.FoodUnits)
	}
}

func TestAdvanceDay_WithFollower(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 20
	s.Gold = 100
	follower := Character{
		Name:         "Guard",
		CombatSkill:  4,
		MaxEndurance: 8,
		DailyWage:    3,
		Morale:       5,
	}
	s.AddFollower(follower)

	AdvanceDay(s)

	// Prince (1) + follower (1) = 2 food/day
	if s.FoodUnits != 18 {
		t.Errorf("food = %d, want 18 (prince + follower)", s.FoodUnits)
	}
	// Wages paid
	if s.Gold != 97 {
		t.Errorf("gold = %d, want 97 after paying 3 gold wage", s.Gold)
	}
}

func TestAdvanceDay_Starvation(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 0

	// r216b: starvation reduces CS (via StarvationDays) — no wounds
	AdvanceDay(s)
	if s.Prince.StarvationDays != 1 {
		t.Errorf("StarvationDays = %d, want 1", s.Prince.StarvationDays)
	}

	AdvanceDay(s)
	if s.Prince.StarvationDays != 2 {
		t.Errorf("StarvationDays = %d, want 2", s.Prince.StarvationDays)
	}

	AdvanceDay(s)
	if s.Prince.StarvationDays != 3 {
		t.Errorf("StarvationDays = %d, want 3", s.Prince.StarvationDays)
	}
	// CS penalty is StarvationDays; verify via EffectiveCombatSkill
	baseCS := s.Prince.CombatSkill
	effectiveCS := s.Prince.EffectiveCombatSkill()
	if effectiveCS >= baseCS {
		t.Errorf("EffectiveCombatSkill %d should be < base %d after starvation", effectiveCS, baseCS)
	}
}

func TestAdvanceDay_StarvationResets(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 0
	AdvanceDay(s)
	if s.Prince.StarvationDays != 1 {
		t.Fatalf("StarvationDays = %d, want 1", s.Prince.StarvationDays)
	}

	// Feed the party
	s.FoodUnits = 10
	AdvanceDay(s)
	if s.Prince.StarvationDays != 0 {
		t.Errorf("StarvationDays = %d, want 0 after eating", s.Prince.StarvationDays)
	}
}

func TestAdvanceDay_DayOfWeek(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 100

	for i := 1; i <= 7; i++ {
		AdvanceDay(s)
	}

	if s.Week != 2 {
		t.Errorf("Week = %d, want 2 after 7 days", s.Week)
	}
	if s.DayOfWeek != 1 {
		t.Errorf("DayOfWeek = %d, want 1 at start of week 2", s.DayOfWeek)
	}
}

func TestAvailableActions_Wilderness(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5) // forest, no structure

	actions := s.AvailableActions()

	has := func(a Action) bool {
		for _, x := range actions {
			if x == a {
				return true
			}
		}
		return false
	}

	if !has(ActionTravel) {
		t.Error("Travel should always be available")
	}
	if !has(ActionRest) {
		t.Error("Rest should always be available")
	}
	if has(ActionSeekAudience) {
		t.Error("Audience should not be available in wilderness")
	}
	if has(ActionSubmitOffering) {
		t.Error("Offering should not be available in wilderness")
	}
	if has(ActionBuyFood) {
		t.Error("Buy food should not be available in wilderness")
	}
}

func TestAvailableActions_Town(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon - town

	actions := s.AvailableActions()

	has := func(a Action) bool {
		for _, x := range actions {
			if x == a {
				return true
			}
		}
		return false
	}

	if !has(ActionTravel) {
		t.Error("Travel should be available in town")
	}
	if !has(ActionBuyFood) {
		t.Error("Buy food should be available in town")
	}
	if !has(ActionSeekNews) {
		t.Error("Seek news should be available in town")
	}
	if !has(ActionSeekFollowers) {
		t.Error("Hire followers should be available in town")
	}
}

func TestAvailableActions_Castle(t *testing.T) {
	s := NewGameState()
	// Find a castle hex
	s.CurrentHex = NewHexID(12, 13) // Hulora Castle

	actions := s.AvailableActions()

	has := func(a Action) bool {
		for _, x := range actions {
			if x == a {
				return true
			}
		}
		return false
	}

	if !has(ActionSeekAudience) {
		t.Error("Seek audience should be available at castle")
	}
}

func TestAvailableActions_Ruins(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(9, 1) // Jakor's Keep ruins

	actions := s.AvailableActions()

	has := func(a Action) bool {
		for _, x := range actions {
			if x == a {
				return true
			}
		}
		return false
	}

	if !has(ActionSearchRuins) {
		t.Error("Search ruins should be available at ruins")
	}
}

func TestAvailableActions_RuinsAlreadySearched(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(11, 9)
	s.GetHexFlags(s.CurrentHex).Searched = true

	actions := s.AvailableActions()

	for _, a := range actions {
		if a == ActionSearchRuins {
			t.Error("Search ruins should not be available after already searched")
		}
	}
}

func TestPartySize(t *testing.T) {
	s := NewGameState()
	if s.PartySize() != 1 {
		t.Errorf("PartySize = %d, want 1 (prince only)", s.PartySize())
	}

	s.AddFollower(Character{Name: "A"})
	s.AddFollower(Character{Name: "B"})
	if s.PartySize() != 3 {
		t.Errorf("PartySize = %d, want 3", s.PartySize())
	}
}

func TestDailyFoodNeeded(t *testing.T) {
	s := NewGameState()
	// Use mountains — mounts don't forage there (r215f)
	s.CurrentHex = NewHexID(4, 3) // Mountains
	if s.DailyFoodNeeded() != 1 {
		t.Errorf("DailyFoodNeeded = %d, want 1 (prince alone)", s.DailyFoodNeeded())
	}

	s.Prince.HasMount = true
	if s.DailyFoodNeeded() != 3 {
		t.Errorf("DailyFoodNeeded = %d, want 3 (prince + horse in mountains)", s.DailyFoodNeeded())
	}

	s.AddFollower(Character{Name: "Guard"})
	if s.DailyFoodNeeded() != 4 {
		t.Errorf("DailyFoodNeeded = %d, want 4 (prince + horse + follower in mountains)", s.DailyFoodNeeded())
	}

	// In forest, mount forages for free (r215f)
	s.CurrentHex = NewHexID(8, 5) // Forest
	if s.DailyFoodNeeded() != 2 { // prince + follower, horse free
		t.Errorf("DailyFoodNeeded = %d, want 2 (prince + follower, horse forages in forest)", s.DailyFoodNeeded())
	}
}

func TestHasGuide(t *testing.T) {
	s := NewGameState()
	if s.HasGuide() {
		t.Error("should have no guide initially")
	}

	s.AddFollower(Character{Name: "Guide", IsGuide: true})
	if !s.HasGuide() {
		t.Error("should have guide after adding one")
	}
}

func TestRemoveFollower(t *testing.T) {
	s := NewGameState()
	s.AddFollower(Character{Name: "Alice"})
	s.AddFollower(Character{Name: "Bob"})

	removed := s.RemoveFollower("Alice")
	if !removed {
		t.Error("RemoveFollower should return true when found")
	}
	if s.PartySize() != 2 {
		t.Errorf("PartySize = %d, want 2 after removal", s.PartySize())
	}
	if s.RemoveFollower("Charlie") {
		t.Error("RemoveFollower should return false for unknown name")
	}
}

func TestAdvanceDay_PlagueDust(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 100
	s.Prince.PlagueDustActive = true
	woundsBefore := s.Prince.Wounds

	AdvanceDay(s)

	if s.Prince.Wounds <= woundsBefore {
		t.Error("plague dust should deal wounds each day")
	}
}

func TestAdvanceDay_PlagueDustEventuallyClears(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 1000
	s.Prince.PlagueDustActive = true
	s.Prince.MaxEndurance = 999 // prevent death

	cleared := false
	for i := 0; i < 100; i++ {
		AdvanceDay(s)
		if !s.Prince.PlagueDustActive {
			cleared = true
			break
		}
	}
	if !cleared {
		t.Error("plague dust should eventually clear via recovery roll")
	}
}

func TestAdvanceDay_StarvationFollowerDesertion(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 0
	s.Prince.WitWiles = 0 // no wit/wiles bonus means desertion is more likely
	// Add many followers so at least one deserts across multiple trials
	for i := 0; i < 10; i++ {
		s.AddFollower(Character{Name: fmt.Sprintf("Guard%d", i), CombatSkill: 4, MaxEndurance: 8, Morale: 3})
	}

	deserted := false
	for attempt := 0; attempt < 20; attempt++ {
		s2 := NewGameState()
		s2.FoodUnits = 0
		s2.Prince.WitWiles = 0
		for i := 0; i < 5; i++ {
			s2.AddFollower(Character{Name: fmt.Sprintf("Guard%d", i), CombatSkill: 4, MaxEndurance: 8, Morale: 3})
		}
		before := len(s2.Party)
		AdvanceDay(s2)
		if len(s2.Party) < before {
			deserted = true
			break
		}
	}
	if !deserted {
		t.Error("followers should sometimes desert when starving")
	}
}

func TestAdvanceDay_TrueLoveNeverDeserts(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 0
	s.Prince.WitWiles = 0
	tl := Character{Name: "True Love", IsTrueLove: true, CombatSkill: 3, MaxEndurance: 6, DailyWage: 0}
	s.AddFollower(tl)

	for i := 0; i < 20; i++ {
		s.FoodUnits = 0
		AdvanceDay(s)
	}

	found := false
	for _, f := range s.Party {
		if f.IsTrueLove {
			found = true
			break
		}
	}
	if !found {
		t.Error("true love follower should never desert from starvation")
	}
}

func TestEffectiveWitWiles_TrueLoveBonus(t *testing.T) {
	s := NewGameState()
	base := s.EffectiveWitWiles()

	s.AddFollower(Character{Name: "True Love", IsTrueLove: true})
	withTL := s.EffectiveWitWiles()

	if withTL != base+1 {
		t.Errorf("EffectiveWitWiles = %d, want %d (base %d + 1 for true love)", withTL, base+1, base)
	}
}

func TestEffectiveWitWiles_NoDouble(t *testing.T) {
	s := NewGameState()
	s.AddFollower(Character{Name: "True Love 1", IsTrueLove: true})
	s.AddFollower(Character{Name: "True Love 2", IsTrueLove: true})
	base := s.Prince.WitWiles

	ww := s.EffectiveWitWiles()
	if ww != base+1 {
		t.Errorf("two true loves should only give +1 W&W, got %d (base %d)", ww, base)
	}
}

func TestAllMounted_AllHaveMounts(t *testing.T) {
	s := NewGameState()
	s.Prince.HasMount = true
	s.AddFollower(Character{Name: "Lancer", HasMount: true})

	if !s.AllMounted() {
		t.Error("AllMounted should be true when prince and all followers have mounts")
	}
}

func TestAllMounted_FollowerWithoutMount(t *testing.T) {
	s := NewGameState()
	s.Prince.HasMount = true
	s.AddFollower(Character{Name: "Foot soldier"})

	if s.AllMounted() {
		t.Error("AllMounted should be false when any follower lacks a mount")
	}
}

func TestAllMounted_PrinceWithoutMount(t *testing.T) {
	s := NewGameState()
	s.AddFollower(Character{Name: "Lancer", HasMount: true})

	if s.AllMounted() {
		t.Error("AllMounted should be false when prince lacks a mount")
	}
}

// ---- Difficulty and DayLimit ------------------------------------------------

func TestNewGameStateWithDifficulty_Easy(t *testing.T) {
	s := NewGameStateWithDifficulty(DifficultyEasy)
	if s.Gold != 20 {
		t.Errorf("Easy mode gold = %d, want 20", s.Gold)
	}
	if s.FoodUnits != 20 {
		t.Errorf("Easy mode food = %d, want 20", s.FoodUnits)
	}
	if s.Difficulty != DifficultyEasy {
		t.Errorf("Difficulty = %v, want Easy", s.Difficulty)
	}
}

func TestNewGameStateWithDifficulty_Hard(t *testing.T) {
	s := NewGameStateWithDifficulty(DifficultyHard)
	if s.Gold != 5 {
		t.Errorf("Hard mode gold = %d, want 5", s.Gold)
	}
	if s.FoodUnits != 7 {
		t.Errorf("Hard mode food = %d, want 7", s.FoodUnits)
	}
	if s.Difficulty != DifficultyHard {
		t.Errorf("Difficulty = %v, want Hard", s.Difficulty)
	}
}

func TestNewGameStateWithDifficulty_Normal(t *testing.T) {
	s := NewGameStateWithDifficulty(DifficultyNormal)
	if s.Difficulty != DifficultyNormal {
		t.Errorf("Difficulty = %v, want Normal", s.Difficulty)
	}
	// Normal starts with default gold/food
	if s.Gold != 10 {
		t.Errorf("Normal mode gold = %d, want 10", s.Gold)
	}
}

func TestDayLimit_Easy(t *testing.T) {
	s := NewGameStateWithDifficulty(DifficultyEasy)
	if s.DayLimit() != 80 {
		t.Errorf("Easy DayLimit = %d, want 80", s.DayLimit())
	}
}

func TestDayLimit_Normal(t *testing.T) {
	s := NewGameState()
	if s.DayLimit() != 70 {
		t.Errorf("Normal DayLimit = %d, want 70", s.DayLimit())
	}
}

func TestDifficultyString(t *testing.T) {
	if DifficultyEasy.String() != "Easy" {
		t.Errorf("DifficultyEasy.String() = %q, want Easy", DifficultyEasy.String())
	}
	if DifficultyHard.String() != "Hard" {
		t.Errorf("DifficultyHard.String() = %q, want Hard", DifficultyHard.String())
	}
	if DifficultyNormal.String() != "Normal" {
		t.Errorf("DifficultyNormal.String() = %q, want Normal", DifficultyNormal.String())
	}
}

// ---- TotalMounts ------------------------------------------------------------

func TestTotalMounts(t *testing.T) {
	s := NewGameState()
	if s.TotalMounts() != 0 {
		t.Errorf("TotalMounts = %d, want 0 with no mounts", s.TotalMounts())
	}

	s.Prince.HasMount = true
	if s.TotalMounts() != 1 {
		t.Errorf("TotalMounts = %d, want 1 (prince only)", s.TotalMounts())
	}

	s.AddFollower(Character{Name: "Lancer", HasMount: true})
	if s.TotalMounts() != 2 {
		t.Errorf("TotalMounts = %d, want 2 (prince + follower)", s.TotalMounts())
	}
}

// ---- CheckLoseConditions (hard difficulty) ----------------------------------

func TestCheckLoseConditions_HardDayLimit(t *testing.T) {
	s := NewGameStateWithDifficulty(DifficultyHard)
	s.Day = 71 // past the 70-day limit
	lost, _ := CheckLoseConditions(s)
	if !lost {
		t.Error("should lose when day > 70 on hard mode")
	}
}

// ---- Win condition: Noble Ally + north of Tragoth ---------------------------

func TestWinCondition_NobleAllyNorth(t *testing.T) {
	s := NewGameState()
	s.Flags.NobleAllySecured = true
	s.CurrentHex = NewHexID(1, 1) // Ogon, north of Tragoth

	won, _ := CheckWinConditions(s)
	if !won {
		t.Error("noble ally + north of Tragoth should win")
	}
}

func TestWinCondition_NobleAllySouth(t *testing.T) {
	s := NewGameState()
	s.Flags.NobleAllySecured = true
	s.CurrentHex = NewHexID(5, 15) // south of Tragoth

	won, _ := CheckWinConditions(s)
	if won {
		t.Error("noble ally south of Tragoth should not win")
	}
}

// ---- Win condition: Staff of Command + north --------------------------------

func TestWinCondition_StaffOfCommandNorth(t *testing.T) {
	s := NewGameState()
	s.Prince.AddPossession(PossStaffOfCommand)
	s.CurrentHex = NewHexID(1, 1)

	won, _ := CheckWinConditions(s)
	if !won {
		t.Error("Staff of Command + north should win")
	}
}

// ---- StatusLine content checks ----------------------------------------------

func TestStatusLine_ContainsDay(t *testing.T) {
	s := NewGameState()
	line := s.StatusLine()
	if line == "" {
		t.Error("StatusLine should not be empty")
	}
	// Should contain "Day"
	found := false
	for i := 0; i <= len(line)-3; i++ {
		if line[i:i+3] == "Day" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("StatusLine should contain 'Day', got: %s", line)
	}
}

// ---- AddNote ----------------------------------------------------------------

func TestAddNote(t *testing.T) {
	s := NewGameState()
	s.AddNote("Found a hidden temple")

	if len(s.FieldNotes) != 1 {
		t.Errorf("FieldNotes len = %d, want 1", len(s.FieldNotes))
	}
	if s.FieldNotes[0].Text != "Found a hidden temple" {
		t.Errorf("FieldNotes[0].Text = %q, want 'Found a hidden temple'", s.FieldNotes[0].Text)
	}
	if s.FieldNotes[0].Day != s.Day {
		t.Errorf("FieldNotes[0].Day = %d, want %d", s.FieldNotes[0].Day, s.Day)
	}
}

// ---- AdvanceDay: day advances -----------------------------------------------

func TestAdvanceDay_DayIncrements(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 100
	dayBefore := s.Day
	AdvanceDay(s)
	if s.Day != dayBefore+1 {
		t.Errorf("Day = %d, want %d after AdvanceDay", s.Day, dayBefore+1)
	}
}

func TestAdvanceDay_WageDeduction(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 100
	s.Gold = 50
	s.AddFollower(Character{Name: "Hire", CombatSkill: 3, MaxEndurance: 6, DailyWage: 5, Morale: 6})
	AdvanceDay(s)
	if s.Gold != 45 {
		t.Errorf("Gold = %d, want 45 after paying 5 gold wage", s.Gold)
	}
}

// ---- Escapee deserts in settlement ------------------------------------------

func TestAdvanceDay_EscapeeDesertsInSettlement(t *testing.T) {
	s := NewGameState()
	s.FoodUnits = 100
	s.Gold = 100
	s.CurrentHex = NewHexID(1, 1) // Ogon - settlement
	s.AddFollower(Character{Name: "Freed Prisoner", IsEscapee: true, CombatSkill: 1, MaxEndurance: 3})

	AdvanceDay(s)

	for _, f := range s.Party {
		if f.IsEscapee {
			t.Error("Escapee should have deserted when entering a settlement")
		}
	}
}

// ---- TutorialState ----------------------------------------------------------

func TestTutorial_IsActive(t *testing.T) {
	ts := NewTutorialState()
	if !ts.IsActive() {
		t.Error("fresh tutorial should be active")
	}

	ts.Skip()
	if ts.IsActive() {
		t.Error("skipped tutorial should not be active")
	}
}

func TestTutorial_NilIsNotActive(t *testing.T) {
	var ts *TutorialState
	if ts.IsActive() {
		t.Error("nil tutorial should not be active")
	}
}

func TestTutorial_Hint(t *testing.T) {
	ts := NewTutorialState()

	h := ts.Hint()
	if h == "" {
		t.Error("fresh tutorial should return a non-empty hint")
	}

	ts.Skip()
	if ts.Hint() != "" {
		t.Error("done tutorial should return empty hint")
	}
}

func TestTutorial_NilHint(t *testing.T) {
	var ts *TutorialState
	if ts.Hint() != "" {
		t.Error("nil tutorial Hint should return empty string")
	}
}

func TestTutorial_OnAction_BuyFoodAdvances(t *testing.T) {
	ts := NewTutorialState()
	if ts.Step != TutStepBuyFood {
		t.Fatalf("initial step should be TutStepBuyFood, got %d", ts.Step)
	}

	ts.OnAction(ActionBuyFood)
	if ts.Step != TutStepTravel {
		t.Errorf("after BuyFood action, step = %d, want TutStepTravel", ts.Step)
	}
}

func TestTutorial_OnTravel_Advances(t *testing.T) {
	ts := NewTutorialState()
	ts.Step = TutStepTravel

	ts.OnTravel()
	if ts.Step != TutStepRest {
		t.Errorf("after OnTravel, step = %d, want TutStepRest", ts.Step)
	}
}

func TestTutorial_OnAction_RestAdvancesToDone(t *testing.T) {
	ts := NewTutorialState()
	ts.Step = TutStepRest

	ts.OnAction(ActionRest)
	if ts.Step != TutStepDone {
		t.Errorf("after Rest action, step = %d, want TutStepDone", ts.Step)
	}
}

func TestTutorial_AllHintSteps(t *testing.T) {
	ts := NewTutorialState()
	steps := []TutorialStep{TutStepBuyFood, TutStepTravel, TutStepRest}
	for _, step := range steps {
		ts.Step = step
		if ts.Hint() == "" {
			t.Errorf("step %d should have a non-empty hint", step)
		}
	}
}

func TestTutorial_OnAction_DoneDoesNothing(t *testing.T) {
	ts := NewTutorialState()
	ts.Step = TutStepDone
	ts.OnAction(ActionBuyFood) // should be a no-op
	if ts.Step != TutStepDone {
		t.Error("OnAction when done should not change step")
	}
}

func TestTutorial_OnTravel_WrongStepIgnored(t *testing.T) {
	ts := NewTutorialState()
	ts.Step = TutStepBuyFood
	ts.OnTravel() // should be a no-op
	if ts.Step != TutStepBuyFood {
		t.Error("OnTravel at wrong step should not advance")
	}
}

func TestTutorial_SkipNil(t *testing.T) {
	var ts *TutorialState
	ts.Skip() // should not panic
}

// ---- Character: CurrentEndurance with PoisonWounds --------------------------

func TestCurrentEndurance_WithPoison(t *testing.T) {
	c := Character{
		CombatSkill:  5,
		MaxEndurance: 10,
		Wounds:       2,
		PoisonWounds: 3,
	}
	// 10 - 2 - 3 = 5
	if c.CurrentEndurance() != 5 {
		t.Errorf("CurrentEndurance = %d, want 5", c.CurrentEndurance())
	}
}

func TestIsUnconscious(t *testing.T) {
	c := Character{MaxEndurance: 9, Wounds: 8}
	if !c.IsUnconscious() {
		t.Error("wounds == MaxEndurance-1 should be unconscious")
	}

	c.Wounds = 7
	if c.IsUnconscious() {
		t.Error("wounds < MaxEndurance-1 should not be unconscious")
	}
}

func TestIsUnconscious_ZeroWoundsNotUnconscious(t *testing.T) {
	c := Character{MaxEndurance: 9, Wounds: 0}
	if c.IsUnconscious() {
		t.Error("0 wounds should not be unconscious (requires wounds > 0)")
	}
}

// ---- EffectiveCombatSkill: seriously wounded --------------------------------

func TestEffectiveCombatSkill_SeriouslyWounded(t *testing.T) {
	c := Character{
		CombatSkill:  6,
		MaxEndurance: 8,
		Wounds:       4, // 4/8 = 50% = seriously wounded (wounds*2 >= MaxEndurance)
	}
	// -1 for any wound, -1 for seriously wounded = CS 4
	if c.EffectiveCombatSkill() != 4 {
		t.Errorf("EffectiveCombatSkill seriously wounded = %d, want 4", c.EffectiveCombatSkill())
	}
}

func TestEffectiveCombatSkill_MinOne(t *testing.T) {
	c := Character{
		CombatSkill:    1,
		MaxEndurance:   10,
		Wounds:         5, // seriously wounded
		StarvationDays: 5,
	}
	// Would go negative but should clamp to 1
	if c.EffectiveCombatSkill() < 1 {
		t.Errorf("EffectiveCombatSkill should never be < 1, got %d", c.EffectiveCombatSkill())
	}
}

func TestEffectiveCombatSkill_PoisonPenalty(t *testing.T) {
	c := Character{
		CombatSkill:  6,
		MaxEndurance: 10,
		PoisonWounds: 3, // -1 CS
	}
	// -1 for poison (3/3 = 1)
	if c.EffectiveCombatSkill() != 5 {
		t.Errorf("EffectiveCombatSkill with 3 poison wounds = %d, want 5", c.EffectiveCombatSkill())
	}
}

// ---- PossessionName coverage ------------------------------------------------

func TestPossessionName_AllItems(t *testing.T) {
	items := []PossessionType{
		PossRoyalHelm, PossStaffOfCommand, PossRingOfCommand, PossAmuletOfPower,
		PossElvenBoots, PossRopeAndGrapnel, PossLantern, PossPoisonAntidote,
		PossHealingPotion, PossMagicSword, PossHolySymbol, PossNobleParchment,
		PossMap, PossRaft, PossGoldenCrown,
	}
	for _, item := range items {
		name := PossessionName(item)
		if name == "" || name == "Unknown Item" {
			t.Errorf("PossessionName(%d) = %q, want a real name", item, name)
		}
	}
	// PossNone should return "Unknown Item"
	if PossessionName(PossNone) != "Unknown Item" {
		t.Errorf("PossessionName(PossNone) = %q, want 'Unknown Item'", PossessionName(PossNone))
	}
}
