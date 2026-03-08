package game

import "fmt"

// IsLost determines if the party becomes lost when traveling.
// onRoad is true when the party is moving along a road hex side.
func IsLost(terrain TerrainType, roll2d6 int, hasGuide bool, hasElvenBoots bool, onRoad bool) bool {
	// Road travel: you can't get lost following a road
	if onRoad {
		return false
	}
	// Elven Boots suppress lost checks in forest terrain
	if hasElvenBoots && terrain == Forest {
		return false
	}
	entry, ok := r207TravelTable[terrain]
	if !ok {
		return false
	}
	threshold := entry.LostThreshold
	if threshold == 0 {
		return false // can't get lost in farmland or countryside
	}
	if hasGuide {
		threshold += 2 // guides reduce lost chance
	}
	return roll2d6 >= threshold
}

// isLostRiverCrossing checks if the party gets lost trying to cross a river (r204e).
// Lost threshold for river crossing = 8 (per r207 "Cross River" row).
func isLostRiverCrossing(roll2d6 int, hasGuide bool) bool {
	threshold := 8
	if hasGuide {
		threshold += 2
	}
	return roll2d6 >= threshold
}

// TravelEventOccurs determines if a travel event is triggered
func TravelEventOccurs(terrain TerrainType, roll2d6 int) bool {
	entry, ok := r207TravelTable[terrain]
	if !ok {
		return false
	}
	return roll2d6 >= entry.EventThreshold
}

// LookupTravelEvent looks up an event from the travel table
func LookupTravelEvent(terrain TerrainType, roll1d6First, roll1d6Second int) EventID {
	entry, ok := r207TravelTable[terrain]
	if !ok {
		return "e128"
	}
	idx := roll1d6First - 1
	if idx < 0 {
		idx = 0
	}
	if idx > 5 {
		idx = 5
	}
	tableCode := entry.EventRefs[idx]
	return LookupEventRef(tableCode, roll1d6Second)
}

// lookupRoadEvent returns a road travel event (r204c / r230 table).
func lookupRoadEvent(roll1d6 int) EventID {
	return LookupEventRef("r230", roll1d6)
}

// TravelResult holds the result of a travel action
type TravelResult struct {
	Success  bool
	Lost     bool
	LostHex  HexID // where they ended up if lost
	EventID  EventID
	HasEvent bool
	Messages []string
}

// DoTravel executes the travel action to an adjacent hex.
// Implements r204c (road events + terrain fallback), r204e (river crossing two-step).
func DoTravel(s *GameState, targetHex HexID) TravelResult {
	result := TravelResult{}

	current := GetHex(s.CurrentHex)
	if current == nil {
		result.Messages = append(result.Messages, "Error: current hex not found!")
		return result
	}

	// Check target is adjacent
	neighbors := AdjacentHexes(s.CurrentHex)
	isAdjacent := false
	for _, n := range neighbors {
		if n == targetHex {
			isAdjacent = true
			break
		}
	}
	if !isAdjacent {
		result.Messages = append(result.Messages, "That hex is not adjacent to your current position.")
		return result
	}

	target := GetHex(targetHex)
	if target == nil {
		result.Messages = append(result.Messages, "Cannot travel there.")
		return result
	}

	dir := s.CurrentHex.DirectionTo(targetHex)
	onRoad := dir >= 0 && current.RoadSides[dir]
	crossingRiver := dir >= 0 && current.RiverSides[dir]

	// ── River crossing (r204c / r204e) ───────────────────────────────────────
	if crossingRiver {
		if onRoad {
			// Road has a bridge — cross freely without a raft (r204c)
			result.Messages = append(result.Messages, "You cross the bridge over the river.")
			// Fall through to normal road travel / event checks below
		} else {
			// No bridge — raft required (r204e two-step)
			if !s.Prince.HasPossession(PossRaft) {
				result.Messages = append(result.Messages, "A river blocks your path. You need a raft to cross.")
				return result
			}

			// Step 1: lost check for the crossing itself
			riverLostRoll := Roll2d6()
			if isLostRiverCrossing(riverLostRoll, s.HasGuide()) {
				result.Lost = true
				result.Messages = append(result.Messages, "You cannot find a suitable ford — the crossing defeats you today.")
				// Still check for a river event (r204e: event occurs even when lost crossing)
				if Roll2d6() >= 10 {
					result.EventID = lookupRoadEvent(Roll1d6())
					result.HasEvent = true
				}
				return result
			}

			// Step 2: river crossing event check (threshold 10, r230 table)
			if Roll2d6() >= 10 {
				result.EventID = lookupRoadEvent(Roll1d6())
				result.HasEvent = true
			}

			// Raft survival roll
			if Roll1d6() == 1 {
				s.Prince.RemovePossession(PossRaft)
				result.Messages = append(result.Messages, "Your raft is wrecked in the crossing! You swim to the far bank.")
			} else {
				result.Messages = append(result.Messages, "You pole your raft across the river.")
			}

			// If a river event was triggered, return now — don't also check terrain event
			if result.HasEvent {
				result.Success = true
				s.CurrentHex = targetHex
				s.VisitedHexes[targetHex] = true
				revealMapNeighbours(s)
				result.Messages = append(result.Messages, fmt.Sprintf("You cross the river and arrive at %s.", hexDisplayName(targetHex)))
				return result
			}
		}
	}

	// ── Movement announcement ─────────────────────────────────────────────────
	if onRoad {
		result.Messages = append(result.Messages, fmt.Sprintf("Following the road to %s...", hexDisplayName(targetHex)))
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("Traveling to %s (%s)...",
			hexDisplayName(targetHex), target.Terrain.String()))
	}

	// ── Lost check for destination terrain ───────────────────────────────────
	lostRoll := Roll2d6()
	if IsLost(target.Terrain, lostRoll, s.HasGuide(), s.Prince.HasPossession(PossElvenBoots), onRoad) {
		result.Lost = true
		neighbors := AdjacentHexes(targetHex)
		if len(neighbors) > 0 {
			lostDest := neighbors[Roll1d6()%len(neighbors)]
			result.LostHex = lostDest
			s.CurrentHex = lostDest
			s.VisitedHexes[lostDest] = true
			revealMapNeighbours(s)
			result.Messages = append(result.Messages, fmt.Sprintf("You become lost in the %s and end up in hex %s!",
				target.Terrain.String(), hexDisplayName(lostDest)))
		}
	} else {
		result.Success = true
		s.CurrentHex = targetHex
		s.VisitedHexes[targetHex] = true
		revealMapNeighbours(s)
		result.Messages = append(result.Messages, fmt.Sprintf("You arrive at %s.", hexDisplayName(targetHex)))
	}

	// ── Event check ──────────────────────────────────────────────────────────
	// Road travel (r204c): check road event table first; if none, also check terrain.
	// Off-road: check terrain event table only.
	actualHex := GetHex(s.CurrentHex)
	if actualHex != nil && !result.HasEvent {
		if onRoad {
			roadRoll := Roll2d6()
			if roadRoll >= 9 { // road event threshold
				result.EventID = lookupRoadEvent(Roll1d6())
				result.HasEvent = true
			}
			// r204c: if no road event, also check terrain event
			if !result.HasEvent {
				eventRoll := Roll2d6()
				if TravelEventOccurs(actualHex.Terrain, eventRoll) {
					result.EventID = LookupTravelEvent(actualHex.Terrain, Roll1d6(), Roll1d6())
					result.HasEvent = true
				}
			}
		} else {
			eventRoll := Roll2d6()
			if TravelEventOccurs(actualHex.Terrain, eventRoll) {
				result.EventID = LookupTravelEvent(actualHex.Terrain, Roll1d6(), Roll1d6())
				result.HasEvent = true
			}
		}
	}

	return result
}

// revealMapNeighbours marks all hexes adjacent to the current position as visited.
// Used when the prince carries PossMap (Ancient Map).
func revealMapNeighbours(s *GameState) {
	if !s.Prince.HasPossession(PossMap) {
		return
	}
	for _, adj := range AdjacentHexes(s.CurrentHex) {
		s.VisitedHexes[adj] = true
	}
}

// hexDisplayName returns a short display string for a hex
func hexDisplayName(id HexID) string {
	hex := GetHex(id)
	if hex == nil {
		return string(id)
	}
	if hex.Name != "" {
		return fmt.Sprintf("%s (%s)", hex.Name, id)
	}
	return fmt.Sprintf("%s %s", hex.Terrain.String(), id)
}

// DoRest processes a rest day in the current hex
func DoRest(s *GameState) []string {
	var msgs []string
	msgs = append(msgs, "You rest for the day.")

	// Heal wounds
	healMsgs := HealRest(s)
	msgs = append(msgs, healMsgs...)

	// Pay lodging if in settlement
	hex := GetHex(s.CurrentHex)
	if hex != nil && hex.IsSettlement() {
		lodgeMsgs := PayLodging(s)
		msgs = append(msgs, lodgeMsgs...)
	}

	return msgs
}
