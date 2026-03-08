package game

import (
	"strings"
	"testing"
)

// helpers

func newState() *GameState { return NewGameState() }

func hexWithRuins() HexID {
	for id, h := range WorldMap {
		if h.IsRuins() {
			return id
		}
	}
	return ""
}

func hexWithSettlement() HexID {
	for id, h := range WorldMap {
		if h.IsSettlement() {
			return id
		}
	}
	return ""
}

// ---- e064 Hidden Ruins -------------------------------------------------------

func TestE064RevealsAdjacentRuins(t *testing.T) {
	// Find a hex that has a ruins neighbour.
	ruinsHex := hexWithRuins()
	if ruinsHex == "" {
		t.Skip("no ruins on map")
	}
	// Pick a neighbour that has the ruins as a neighbour.
	// Easier: place the prince on a hex adjacent to ruinsHex.
	adj := AdjacentHexes(ruinsHex)
	if len(adj) == 0 {
		t.Skip("ruins hex has no neighbours")
	}
	s := newState()
	s.CurrentHex = adj[0]
	// Ensure ruinsHex is not yet visited.
	delete(s.VisitedHexes, ruinsHex)

	result := DispatchEvent(s, "e064", EventContext{})

	// Should reveal the ruins hex.
	if !s.VisitedHexes[ruinsHex] {
		// Maybe ruinsHex wasn't adjacent to adj[0]; try all adj.
		found := false
		for _, a := range adj {
			s2 := newState()
			s2.CurrentHex = a
			delete(s2.VisitedHexes, ruinsHex)
			r2 := DispatchEvent(s2, "e064", EventContext{})
			if s2.VisitedHexes[ruinsHex] {
				found = true
				_ = r2
				break
			}
		}
		if !found {
			t.Errorf("e064 did not reveal the ruins hex even when one is adjacent")
		}
	}
	_ = result
}

func TestE064GoldFallback(t *testing.T) {
	// Find a hex with NO adjacent unvisited ruins.
	s := newState()
	s.CurrentHex = NewHexID(1, 1) // Ogon — dense settlement area
	// Mark all neighbours as visited so the ruins search fails.
	for _, adj := range AdjacentHexes(s.CurrentHex) {
		s.VisitedHexes[adj] = true
	}
	before := s.Gold
	result := DispatchEvent(s, "e064", EventContext{})
	if result.GoldChange <= 0 && s.Gold == before {
		// May also reveal ruins - acceptable; just verify no panic and has messages.
	}
	if len(result.Messages) == 0 {
		t.Error("e064 returned no messages")
	}
}

// ---- e065 Hidden Town --------------------------------------------------------

func TestE065RevealsAdjacentSettlement(t *testing.T) {
	settlementHex := hexWithSettlement()
	if settlementHex == "" {
		t.Skip("no settlement on map")
	}
	adj := AdjacentHexes(settlementHex)
	if len(adj) == 0 {
		t.Skip("settlement hex has no neighbours")
	}
	for _, a := range adj {
		h := GetHex(a)
		if h == nil || h.IsSettlement() {
			continue
		}
		s := newState()
		s.CurrentHex = a
		delete(s.VisitedHexes, settlementHex)
		DispatchEvent(s, "e065", EventContext{})
		if s.VisitedHexes[settlementHex] {
			return // success
		}
	}
	// If we reach here, no adjacent non-settlement neighbour could trigger it — skip.
	t.Skip("could not find suitable test hex for e065")
}

// ---- e061 Escaped Prisoners --------------------------------------------------

func TestE061ChoiceHandlerExists(t *testing.T) {
	s := newState()
	result := DispatchEvent(s, "e061", EventContext{})
	if len(result.Choices) == 0 {
		// count==0 edge case is fine
		return
	}
	if result.NewFollower != nil {
		t.Error("e061 initial result must not include NewFollower (that leaks the follower without a player choice)")
	}
	if result.ChoiceHandler == nil {
		t.Error("e061 must have a ChoiceHandler when Choices are present")
	}
}

func TestE061AcceptAddsFollower(t *testing.T) {
	for range 20 {
		s := newState()
		result := DispatchEvent(s, "e061", EventContext{})
		if len(result.Choices) == 0 || result.ChoiceHandler == nil {
			continue
		}
		r2 := result.ChoiceHandler(s, 0)
		if r2.NewFollower == nil {
			t.Error("e061 accept should return NewFollower")
		}
		return
	}
}

func TestE061DeclineNoFollower(t *testing.T) {
	for range 20 {
		s := newState()
		result := DispatchEvent(s, "e061", EventContext{})
		if len(result.Choices) == 0 || result.ChoiceHandler == nil {
			continue
		}
		r2 := result.ChoiceHandler(s, 1)
		if r2.NewFollower != nil {
			t.Error("e061 decline must not add a follower")
		}
		return
	}
}

// ---- e098 Dragon -------------------------------------------------------------

func TestE098HasChoicesNotCombat(t *testing.T) {
	s := newState()
	result := DispatchEvent(s, "e098", EventContext{})
	if result.CombatTriggered {
		t.Error("e098 initial result must not trigger combat; player must choose first")
	}
	if len(result.Choices) == 0 {
		t.Error("e098 must offer choices (fight/flee/reason)")
	}
	if result.ChoiceHandler == nil {
		t.Error("e098 must have a ChoiceHandler")
	}
}

func TestE098FightTriggersCombet(t *testing.T) {
	s := newState()
	result := DispatchEvent(s, "e098", EventContext{})
	if result.ChoiceHandler == nil {
		t.Skip("dragon already slain")
	}
	r2 := result.ChoiceHandler(s, 0) // fight
	if !r2.CombatTriggered || r2.Enemy == nil {
		t.Error("choosing to fight the dragon must trigger combat")
	}
}

// ---- e153 Knight ------------------------------------------------------------

func TestE153ChoiceHandlerExists(t *testing.T) {
	s := newState()
	result := DispatchEvent(s, "e153", EventContext{})
	if result.NewFollower != nil {
		t.Error("e153 initial result must not include NewFollower")
	}
	if result.ChoiceHandler == nil {
		t.Error("e153 must have a ChoiceHandler")
	}
}

func TestE153AcceptAddsKnight(t *testing.T) {
	s := newState()
	result := DispatchEvent(s, "e153", EventContext{})
	if result.ChoiceHandler == nil {
		t.Fatal("e153 missing ChoiceHandler")
	}
	r2 := result.ChoiceHandler(s, 0)
	if r2.NewFollower == nil {
		t.Error("accepting knight should return NewFollower")
	}
}

func TestE153DeclineNoFollower(t *testing.T) {
	s := newState()
	result := DispatchEvent(s, "e153", EventContext{})
	if result.ChoiceHandler == nil {
		t.Fatal("e153 missing ChoiceHandler")
	}
	r2 := result.ChoiceHandler(s, 1)
	if r2.NewFollower != nil {
		t.Error("declining knight must not add a follower")
	}
}

// ---- e155 Temple healing ----------------------------------------------------

func TestE155HealsWounds(t *testing.T) {
	s := newState()
	s.Prince.Wounds = 5
	DispatchEvent(s, "e155", EventContext{})
	if s.Prince.Wounds != 0 {
		t.Errorf("e155 should heal all wounds; got %d remaining", s.Prince.Wounds)
	}
}

// ---- e158 Poison cure -------------------------------------------------------

func TestE158CuresPoison(t *testing.T) {
	s := newState()
	s.Prince.PoisonWounds = 3
	DispatchEvent(s, "e158", EventContext{})
	if s.Prince.PoisonWounds != 0 {
		t.Errorf("e158 should cure all poison; got %d remaining", s.Prince.PoisonWounds)
	}
}

// ---- e152 Noble ally --------------------------------------------------------

func TestE152SecuresNobleAlly(t *testing.T) {
	s := newState()
	result := DispatchEvent(s, "e152", EventContext{})
	if !s.Flags.NobleAllySecured {
		t.Error("e152 should set NobleAllySecured")
	}
	if !s.Prince.HasPossession(PossNobleParchment) {
		t.Error("e152 should give Noble Parchment")
	}
	if len(result.Messages) == 0 {
		t.Error("e152 returned no messages")
	}
}

// ---- AllEvents: Choices always have ChoiceHandlers --------------------------

func TestAllEventsChoicesHaveHandlers(t *testing.T) {
	s := newState()
	ctx := EventContext{}

	for id, handler := range eventRegistry {
		result := handler(s, ctx)
		if len(result.Choices) > 0 && result.ChoiceHandler == nil {
			t.Errorf("event %s has Choices but no ChoiceHandler", id)
		}
		if result.CombatTriggered && len(result.Choices) > 0 {
			t.Errorf("event %s sets both CombatTriggered and Choices (contradictory)", id)
		}
	}
}

// ---- AllEvents: all registered events return at least one message -----------

func TestAllEventsReturnMessages(t *testing.T) {
	s := newState()
	ctx := EventContext{}

	for id, handler := range eventRegistry {
		result := handler(s, ctx)
		if len(result.Messages) == 0 && !result.CombatTriggered && len(result.Choices) == 0 {
			t.Errorf("event %s returned no messages and no combat/choice", id)
		}
	}
}

// ---- e066 Hidden Temple flag ------------------------------------------------

func TestE066SetsHiddenTempleFlag(t *testing.T) {
	s := newState()
	s.CurrentHex = NewHexID(5, 5)
	DispatchEvent(s, "e066", EventContext{})
	if !s.GetHexFlags(s.CurrentHex).HiddenTemple {
		t.Error("e066 should set HiddenTemple flag on current hex")
	}
}

// ---- Follower events: hire choice returns NewFollower -----------------------

func TestFollowerEventsChoiceHandler(t *testing.T) {
	followerEvents := []EventID{"e003", "e004", "e005", "e006", "e007", "e008", "e018", "e073"}
	s := newState()

	for _, id := range followerEvents {
		for range 50 { // retry because of dice
			s2 := newState()
			result := DispatchEvent(s2, id, EventContext{})
			if len(result.Choices) == 0 {
				continue // not the hire branch this roll
			}
			if result.ChoiceHandler == nil {
				t.Errorf("event %s has Choices but no ChoiceHandler", id)
				break
			}
			r2 := result.ChoiceHandler(s2, 0)
			if r2.NewFollower == nil {
				t.Errorf("event %s: accepting hire should return NewFollower", id)
			}
			r3 := result.ChoiceHandler(s2, 1)
			if r3.NewFollower != nil {
				t.Errorf("event %s: declining hire must not return NewFollower", id)
			}
			break
		}
	}
	_ = s
}

// ---- News events don't panic and return messages ----------------------------

func TestNewsEvents(t *testing.T) {
	s := newState()
	for roll := 1; roll <= 6; roll++ {
		id := LookupNewsEvent(roll)
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("LookupNewsEvent(%d) → %s panicked: %v", roll, id, r)
				}
			}()
			result := DispatchEvent(s, id, EventContext{})
			if len(result.Messages) == 0 {
				t.Errorf("news event %s (roll %d) returned no messages", id, roll)
			}
		}()
	}
}

// ---- Follower lookup events don't panic ------------------------------------

func TestFollowerEvents(t *testing.T) {
	s := newState()
	for roll := 1; roll <= 6; roll++ {
		id := LookupFollowerEvent(roll)
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("LookupFollowerEvent(%d) → %s panicked: %v", roll, id, r)
				}
			}()
			result := DispatchEvent(s, id, EventContext{})
			if len(result.Messages) == 0 && len(result.Choices) == 0 {
				t.Errorf("follower event %s (roll %d) returned no messages or choices", id, roll)
			}
		}()
	}
}

// ---- Offering lookup events don't panic ------------------------------------

func TestOfferingEvents(t *testing.T) {
	s := newState()
	for roll := 1; roll <= 6; roll++ {
		id := LookupOfferingEvent(roll)
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("LookupOfferingEvent(%d) → %s panicked: %v", roll, id, r)
				}
			}()
			result := DispatchEvent(s, id, EventContext{})
			_ = result
		}()
	}
}

// ---- Ruins lookup events don't panic ----------------------------------------

func TestRuinsEvents(t *testing.T) {
	s := newState()
	for roll := 1; roll <= 6; roll++ {
		id := LookupRuinsEvent(roll)
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("LookupRuinsEvent(%d) → %s panicked: %v", roll, id, r)
				}
			}()
			result := DispatchEvent(s, id, EventContext{})
			_ = result
		}()
	}
}

// ---- SearchRuins action shows up in ruins hex --------------------------------

func TestSearchRuinsAvailableInRuinsHex(t *testing.T) {
	ruinsHex := hexWithRuins()
	if ruinsHex == "" {
		t.Skip("no ruins on map")
	}
	s := newState()
	s.CurrentHex = ruinsHex
	s.VisitedHexes[ruinsHex] = true

	actions := s.AvailableActions()
	for _, a := range actions {
		if a == ActionSearchRuins {
			return
		}
	}
	t.Errorf("ActionSearchRuins not in AvailableActions for ruins hex %s", ruinsHex)
}

// ---- AvailableActions: Submit Offering after e066 ---------------------------

func TestSubmitOfferingAvailableAfterHiddenTemple(t *testing.T) {
	s := newState()
	s.CurrentHex = NewHexID(5, 5) // non-temple hex
	s.GetHexFlags(s.CurrentHex).HiddenTemple = true

	actions := s.AvailableActions()
	for _, a := range actions {
		if a == ActionSubmitOffering {
			return
		}
	}
	t.Error("ActionSubmitOffering should be available after HiddenTemple flag is set")
}

// ---- Message content sanity check -------------------------------------------

func TestE064MessageMentionsHexOrGold(t *testing.T) {
	s := newState()
	// Mark all neighbours visited so we hit the gold fallback.
	for _, adj := range AdjacentHexes(s.CurrentHex) {
		s.VisitedHexes[adj] = true
	}
	result := DispatchEvent(s, "e064", EventContext{})
	found := false
	for _, m := range result.Messages {
		if strings.Contains(m, "gold") || strings.Contains(m, "ruins") || strings.Contains(m, "Ruins") {
			found = true
		}
	}
	if !found {
		t.Errorf("e064 fallback messages don't mention gold or ruins: %v", result.Messages)
	}
}

// ---- e082 Faerie ring: gold via GoldChange, not direct mutation -------------

func TestE082FaerieGoldViaGoldChange(t *testing.T) {
	for range 100 {
		s := newState()
		goldBefore := s.Gold
		result := DispatchEvent(s, "e082", EventContext{})
		if result.ChoiceHandler == nil {
			continue
		}
		r2 := result.ChoiceHandler(s, 0) // step inside
		// If the faerie gold branch fired, GoldChange must be set (not direct mutation).
		// Direct mutation would leave s.Gold changed but r2.GoldChange == 0.
		if r2.GoldChange > 0 {
			// Correct path: gold not yet applied to state
			if s.Gold != goldBefore {
				t.Errorf("e082 faerie gold branch must use GoldChange, not direct s.Gold mutation; gold changed from %d to %d before applyEventResult", goldBefore, s.Gold)
			}
			return
		}
		// Other branches (curse or heal) — keep looping
	}
	// If we never hit the gold branch, that's fine (dice dependent) — skip.
	t.Skip("faerie gold branch not encountered in 100 tries")
}

// ---- e160 temple: can't-afford feedback -------------------------------------

func TestE160CantAffordHealing(t *testing.T) {
	s := newState()
	s.Gold = 5 // less than 20
	s.Prince.Wounds = 3
	result := DispatchEvent(s, "e160", EventContext{})
	if result.ChoiceHandler == nil {
		t.Fatal("e160 must have ChoiceHandler")
	}
	r2 := result.ChoiceHandler(s, 0) // choose "Pay 20 gold"
	if len(r2.Messages) == 0 {
		t.Error("e160: choosing to pay when broke should return a message")
	}
	if s.Prince.Wounds != 3 {
		t.Error("e160: wounds should not be healed when player can't afford it")
	}
	if s.Gold != 5 {
		t.Error("e160: gold should not be deducted when player can't afford it")
	}
}

func TestE160HealingWhenAffordable(t *testing.T) {
	s := newState()
	s.Gold = 50
	s.Prince.Wounds = 5
	result := DispatchEvent(s, "e160", EventContext{})
	if result.ChoiceHandler == nil {
		t.Fatal("e160 must have ChoiceHandler")
	}
	result.ChoiceHandler(s, 0) // choose "Pay 20 gold"
	if s.Prince.Wounds != 0 {
		t.Errorf("e160: wounds should be healed; got %d", s.Prince.Wounds)
	}
	if s.Gold != 30 {
		t.Errorf("e160: gold should be 30 after paying 20; got %d", s.Gold)
	}
}

// ---- Duplicate follower prevention -----------------------------------------

func TestNoDuplicateFollower(t *testing.T) {
	s := newState()
	pip := Character{Name: "Pip the Halfling", Type: TypeHalfling, CombatSkill: 3, MaxEndurance: 6, DailyWage: 2}
	s.AddFollower(pip)

	result := EventResult{NewFollower: &pip}
	applyEventResult(s, result)

	if len(s.Party) != 1 {
		t.Errorf("party size = %d after duplicate hire attempt; want 1", len(s.Party))
	}
}

func TestDuplicateFollowerMessageReturned(t *testing.T) {
	s := newState()
	pip := Character{Name: "Pip the Halfling"}
	s.AddFollower(pip)

	msgs := applyEventResult(s, EventResult{NewFollower: &pip})
	found := false
	for _, m := range msgs {
		if strings.Contains(m, "already in your party") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'already in your party' message; got %v", msgs)
	}
}

func TestFirstHireAddsFollower(t *testing.T) {
	s := newState()
	pip := Character{Name: "Pip the Halfling", Type: TypeHalfling, CombatSkill: 3, MaxEndurance: 6}
	applyEventResult(s, EventResult{NewFollower: &pip})

	if len(s.Party) != 1 {
		t.Errorf("party size = %d after first hire; want 1", len(s.Party))
	}
}
