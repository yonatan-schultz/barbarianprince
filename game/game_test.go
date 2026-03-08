package game

import (
	"testing"
)

// TestDice verifies dice rolls are within valid ranges
func TestDice(t *testing.T) {
	for i := 0; i < 1000; i++ {
		r := Roll1d6()
		if r < 1 || r > 6 {
			t.Errorf("Roll1d6() = %d, want 1-6", r)
		}
	}
	for i := 0; i < 1000; i++ {
		r := Roll2d6()
		if r < 2 || r > 12 {
			t.Errorf("Roll2d6() = %d, want 2-12", r)
		}
	}
	for i := 0; i < 1000; i++ {
		r := Roll1d3()
		if r < 1 || r > 3 {
			t.Errorf("Roll1d3() = %d, want 1-3", r)
		}
	}
}

// TestCombatTable verifies key entries in the combat wounds table
func TestCombatTable(t *testing.T) {
	cases := []struct {
		netRoll int
		want    int
	}{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{5, 1},
		{10, 2},
		{14, 3},
		{16, 5},
		{20, 6},
		{25, 6}, // clamped
	}
	for _, tc := range cases {
		got := CombatWounds(tc.netRoll)
		if got != tc.want {
			t.Errorf("CombatWounds(%d) = %d, want %d", tc.netRoll, got, tc.want)
		}
	}
}

// TestTreasureTable verifies treasure rolls return positive values
func TestTreasureTable(t *testing.T) {
	for wealthCode := 1; wealthCode <= 7; wealthCode++ {
		for roll := 1; roll <= 6; roll++ {
			gold := TreasureRoll(wealthCode, roll)
			if gold <= 0 {
				t.Errorf("TreasureRoll(%d, %d) = %d, want > 0", wealthCode, roll, gold)
			}
		}
	}

	// Higher wealth codes should give more gold
	low := TreasureRoll(1, 6)
	high := TreasureRoll(7, 6)
	if high <= low {
		t.Errorf("wealth code 7 roll 6 (%d) should exceed wealth code 1 roll 6 (%d)", high, low)
	}
}

// TestHexAdjacency verifies adjacency computation for even and odd columns
func TestHexAdjacency(t *testing.T) {
	// Even column (2): neighbors should be correctly offset
	h := NewHexID(2, 5)
	ns := h.Neighbors()

	if ns[DirN] != NewHexID(2, 4) {
		t.Errorf("Even col N neighbor: got %s, want %s", ns[DirN], NewHexID(2, 4))
	}
	if ns[DirS] != NewHexID(2, 6) {
		t.Errorf("Even col S neighbor: got %s, want %s", ns[DirS], NewHexID(2, 6))
	}
	if ns[DirNE] != NewHexID(3, 4) {
		t.Errorf("Even col NE neighbor: got %s, want %s", ns[DirNE], NewHexID(3, 4))
	}
	if ns[DirSE] != NewHexID(3, 5) {
		t.Errorf("Even col SE neighbor: got %s, want %s", ns[DirSE], NewHexID(3, 5))
	}

	// Odd column (3): neighbors are offset differently
	h2 := NewHexID(3, 5)
	ns2 := h2.Neighbors()
	if ns2[DirN] != NewHexID(3, 4) {
		t.Errorf("Odd col N neighbor: got %s, want %s", ns2[DirN], NewHexID(3, 4))
	}
	if ns2[DirNE] != NewHexID(4, 5) {
		t.Errorf("Odd col NE neighbor: got %s, want %s", ns2[DirNE], NewHexID(4, 5))
	}
	if ns2[DirSE] != NewHexID(4, 6) {
		t.Errorf("Odd col SE neighbor: got %s, want %s", ns2[DirSE], NewHexID(4, 6))
	}
	if ns2[DirSW] != NewHexID(2, 6) {
		t.Errorf("Odd col SW neighbor: got %s, want %s", ns2[DirSW], NewHexID(2, 6))
	}
}

// TestHexID verifies col/row parsing
func TestHexID(t *testing.T) {
	id := NewHexID(12, 7)
	if id != "1207" {
		t.Errorf("NewHexID(12,7) = %q, want %q", id, "1207")
	}
	if id.Col() != 12 {
		t.Errorf("Col() = %d, want 12", id.Col())
	}
	if id.Row() != 7 {
		t.Errorf("Row() = %d, want 7", id.Row())
	}
}

// TestTravelEventLookup verifies travel event lookup is deterministic
func TestTravelEventLookup(t *testing.T) {
	// All terrain types should return a valid event ID
	terrains := []TerrainType{Farmland, Countryside, Forest, Hills, Mountains, Swamp, Desert}
	for _, terrain := range terrains {
		for r1 := 1; r1 <= 6; r1++ {
			for r2 := 1; r2 <= 6; r2++ {
				id := LookupTravelEvent(terrain, r1, r2)
				if id == "" {
					t.Errorf("LookupTravelEvent(%v, %d, %d) returned empty EventID", terrain, r1, r2)
				}
			}
		}
	}
}

// TestIsLost verifies lost checks respect guide bonus and Elven Boots
func TestIsLost(t *testing.T) {
	// Farmland: can't get lost
	if IsLost(Farmland, 12, false, false, false) {
		t.Error("Should not get lost in Farmland")
	}

	// Mountains: threshold 8, so roll of 9 with no guide = lost
	if !IsLost(Mountains, 9, false, false, false) {
		t.Error("Should get lost in Mountains on roll 9 without guide")
	}

	// With guide: threshold raises by 2, so 9+2=11, roll of 9 = not lost
	if IsLost(Mountains, 9, true, false, false) {
		t.Error("Should not get lost in Mountains on roll 9 with guide")
	}

	// Elven Boots suppress lost in Forest regardless of roll
	if IsLost(Forest, 12, false, true, false) {
		t.Error("Elven Boots should prevent getting lost in Forest")
	}

	// Elven Boots do NOT suppress lost in other terrain
	if !IsLost(Mountains, 9, false, true, false) {
		t.Error("Elven Boots should not prevent getting lost in Mountains")
	}

	// Road travel: can never get lost
	if IsLost(Mountains, 12, false, false, true) {
		t.Error("Road travel should prevent getting lost")
	}
}

// TestWinConditions verifies win condition checks
func TestWinConditions(t *testing.T) {
	s := NewGameState()
	s.Gold = 500
	s.CurrentHex = NewHexID(1, 1) // Ogon (north of Tragoth)

	won, _ := CheckWinConditions(s)
	if !won {
		t.Error("Should win with 500 gold north of Tragoth")
	}

	// Not enough gold
	s.Gold = 499
	won, _ = CheckWinConditions(s)
	if won {
		t.Error("Should not win with only 499 gold")
	}

	// South of Tragoth
	s.Gold = 500
	s.CurrentHex = NewHexID(1, 15) // south of Tragoth
	won, _ = CheckWinConditions(s)
	if won {
		t.Error("Should not win when south of Tragoth")
	}
}

// TestLoseConditions verifies lose conditions
func TestLoseConditions(t *testing.T) {
	s := NewGameState()

	// Day > 70
	s.Day = 71
	lost, _ := CheckLoseConditions(s)
	if !lost {
		t.Error("Should lose when Day > 70")
	}

	// Prince dead
	s.Day = 10
	s.Prince.Wounds = s.Prince.MaxEndurance + 1
	lost, _ = CheckLoseConditions(s)
	if !lost {
		t.Error("Should lose when prince is dead")
	}
}

// TestCharacterCombatSkill verifies wound penalties apply
func TestCharacterCombatSkill(t *testing.T) {
	c := Character{
		CombatSkill:  5,
		MaxEndurance: 10,
		Wounds:       3,
	}
	// 3 wounds = -1 CS
	if c.EffectiveCombatSkill() != 4 {
		t.Errorf("EffectiveCombatSkill() = %d, want 4", c.EffectiveCombatSkill())
	}
}

// TestMapIntegrity verifies the world map has the required named hexes
func TestMapIntegrity(t *testing.T) {
	requiredHexes := map[HexID]string{
		NewHexID(1, 1):  "Ogon",
		NewHexID(15, 1): "Weshor",
	}

	for id, name := range requiredHexes {
		hex := GetHex(id)
		if hex == nil {
			t.Errorf("Required hex %s not found in map", id)
			continue
		}
		if hex.Name != name {
			t.Errorf("Hex %s name = %q, want %q", id, hex.Name, name)
		}
	}
}

// TestEventDispatch verifies events can be dispatched without panic
func TestEventDispatch(t *testing.T) {
	s := NewGameState()
	ctx := EventContext{}

	// Test a sample of events
	testEvents := []EventID{
		"e003", "e005", "e009", "e017", "e022",
		"e051", "e057", "e074", "e075", "e078",
		"e098", "e128", "e131", "e133", "e148",
		"e155", "e180", "e193",
	}

	for _, id := range testEvents {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("DispatchEvent(%s) panicked: %v", id, r)
				}
			}()
			result := DispatchEvent(s, id, ctx)
			if len(result.Messages) == 0 {
				t.Errorf("DispatchEvent(%s) returned no messages", id)
			}
		}()
	}
}
