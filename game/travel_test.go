package game

import "testing"

func TestTravelEventOccurs(t *testing.T) {
	// Below threshold: no event
	if TravelEventOccurs(Farmland, 6) {
		t.Error("Farmland threshold is 7, roll 6 should not trigger event")
	}
	// At threshold: event
	if !TravelEventOccurs(Farmland, 7) {
		t.Error("Farmland threshold is 7, roll 7 should trigger event")
	}
	// Above threshold: event
	if !TravelEventOccurs(Farmland, 12) {
		t.Error("roll 12 should always trigger event")
	}
}

func TestDoTravel_ValidMove(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon

	neighbors := AdjacentHexes(s.CurrentHex)
	if len(neighbors) == 0 {
		t.Fatal("Ogon should have adjacent hexes")
	}

	target := neighbors[0]
	result := DoTravel(s, target)

	if len(result.Messages) == 0 {
		t.Error("DoTravel should return messages")
	}
	// Either arrived or got lost - either way, hex changed
	if s.CurrentHex == NewHexID(1, 1) && !result.Lost {
		t.Error("hex should have changed after successful travel")
	}
}

func TestDoTravel_InvalidMove(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)

	// Try to travel to a non-adjacent hex
	farHex := NewHexID(10, 10)
	result := DoTravel(s, farHex)

	if result.Success {
		t.Error("travel to non-adjacent hex should not succeed")
	}
	if s.CurrentHex != NewHexID(1, 1) {
		t.Error("current hex should not change on invalid move")
	}
}

func TestDoTravel_VisitedTracked(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)
	neighbors := AdjacentHexes(s.CurrentHex)
	if len(neighbors) == 0 {
		t.Fatal("need adjacent hexes")
	}

	target := neighbors[0]
	if s.VisitedHexes[target] {
		t.Skip("target already marked visited")
	}

	DoTravel(s, target)

	// The hex we ended up in (might be different if lost) should be visited
	if !s.VisitedHexes[s.CurrentHex] {
		t.Error("current hex should be marked visited after travel")
	}
}

func TestAdjacentHexes(t *testing.T) {
	// Central hex should have 6 neighbors
	center := NewHexID(12, 12)
	ns := AdjacentHexes(center)
	if len(ns) != 6 {
		t.Errorf("central hex should have 6 neighbors, got %d", len(ns))
	}

	// Corner hex should have fewer
	corner := NewHexID(1, 1)
	nsCorner := AdjacentHexes(corner)
	if len(nsCorner) >= 6 {
		t.Errorf("corner hex should have fewer than 6 valid neighbors, got %d", len(nsCorner))
	}
}

func TestLookupTablesBoundsCheck(t *testing.T) {
	// All lookup tables should handle out-of-range rolls gracefully
	tables := []struct {
		name string
		fn   func(int) EventID
	}{
		{"ruins", LookupRuinsEvent},
		{"news", LookupNewsEvent},
		{"followers", LookupFollowerEvent},
		{"audience", LookupAudienceEvent},
		{"offering", LookupOfferingEvent},
	}

	for _, tc := range tables {
		for _, roll := range []int{0, 1, 3, 6, 7, 100} {
			id := tc.fn(roll)
			if id == "" {
				t.Errorf("%s lookup returned empty EventID for roll %d", tc.name, roll)
			}
		}
	}
}

func TestDoTravel_RiverBlockedWithoutRaft(t *testing.T) {
	s := NewGameState()
	// Tragoth River runs between row 2 (DirS) and row 3 (DirN) for all cols.
	// Travel from (5,2) south to (5,3) should be blocked without a raft.
	s.CurrentHex = NewHexID(5, 2)
	target := NewHexID(5, 3)

	result := DoTravel(s, target)

	if result.Success {
		t.Error("crossing Tragoth without raft should not succeed")
	}
	if s.CurrentHex != NewHexID(5, 2) {
		t.Error("hex should not change when river crossing is blocked")
	}
	found := false
	for _, msg := range result.Messages {
		if len(msg) > 5 {
			found = true
			break
		}
	}
	if !found {
		t.Error("blocked crossing should return an explanatory message")
	}
}

func TestDoTravel_RiverCrossWithRaft(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 2)
	s.Prince.AddPossession(PossRaft)
	target := NewHexID(5, 3)

	result := DoTravel(s, target)

	// Should succeed (possibly with event)
	if !result.Success && !result.Lost {
		t.Error("crossing Tragoth with raft should succeed or trigger an event")
	}
}

func TestDoTravel_RoadBridgeCrossesWithoutRaft(t *testing.T) {
	s := NewGameState()
	// The road bridge is at (8,2)↔(8,3) crossing the Tragoth River.
	// Without a raft but ON the road, the crossing should succeed.
	s.CurrentHex = NewHexID(8, 2)
	target := NewHexID(8, 3)

	// Verify the road side is actually set
	h := GetHex(NewHexID(8, 2))
	if h == nil {
		t.Fatal("hex (8,2) should exist")
	}
	dir := NewHexID(8, 2).DirectionTo(NewHexID(8, 3))
	if !h.RoadSides[dir] {
		t.Fatal("(8,2) should have a road side toward (8,3)")
	}
	if !h.RiverSides[dir] {
		t.Fatal("(8,2) should have a river side toward (8,3) — Tragoth")
	}

	result := DoTravel(s, target)

	if !result.Success && !result.Lost {
		t.Error("road bridge should allow Tragoth crossing without a raft")
	}
}

func TestDoTravel_RoadBridgeBlockedOffRoad(t *testing.T) {
	s := NewGameState()
	// Off-road Tragoth crossing at col 5 (no bridge) should still be blocked.
	s.CurrentHex = NewHexID(5, 2)
	target := NewHexID(5, 3)

	result := DoTravel(s, target)
	if result.Success {
		t.Error("non-bridge Tragoth crossing should be blocked without raft")
	}
}

func TestRoadSidesSet(t *testing.T) {
	// Verify road data was populated at several expected road hexes.
	pairs := [][4]int{
		{1, 1, 2, 1},   // Ogon east
		{8, 2, 8, 3},   // Tragoth bridge
		{9, 8, 9, 9},   // approach to Cawther
		{9, 19, 8, 19}, // to Brigud
		{3, 20, 2, 20}, // Halowich NW
		{1, 17, 1, 16}, // to Galden
		{5, 22, 5, 23}, // to Adrogat
	}
	for _, p := range pairs {
		h1 := GetHex(NewHexID(p[0], p[1]))
		h2 := GetHex(NewHexID(p[2], p[3]))
		if h1 == nil || h2 == nil {
			t.Errorf("hex (%d,%d) or (%d,%d) missing", p[0], p[1], p[2], p[3])
			continue
		}
		d1 := NewHexID(p[0], p[1]).DirectionTo(NewHexID(p[2], p[3]))
		d2 := NewHexID(p[2], p[3]).DirectionTo(NewHexID(p[0], p[1]))
		if d1 < 0 || !h1.RoadSides[d1] {
			t.Errorf("road not set on (%d,%d)→(%d,%d)", p[0], p[1], p[2], p[3])
		}
		if d2 < 0 || !h2.RoadSides[d2] {
			t.Errorf("road not set on (%d,%d)→(%d,%d)", p[2], p[3], p[0], p[1])
		}
	}
}

func TestDoTravel_PossMapRevealsNeighbours(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(10, 10)
	s.Prince.AddPossession(PossMap)

	target := NewHexID(10, 11)
	// Clear visited state for target's neighbours
	for _, adj := range AdjacentHexes(target) {
		delete(s.VisitedHexes, adj)
	}

	DoTravel(s, target)

	// After arriving, all neighbours of the new position should be visible
	revealed := 0
	for _, adj := range AdjacentHexes(s.CurrentHex) {
		if s.VisitedHexes[adj] {
			revealed++
		}
	}
	if revealed == 0 {
		t.Error("PossMap should reveal neighbours after moving")
	}
}
