package game

import "testing"

func TestBuyFood_Normal(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon - a town
	s.Gold = 50
	foodBefore := s.FoodUnits

	msgs := BuyFood(s, 5)

	if len(msgs) == 0 {
		t.Error("BuyFood should return messages")
	}
	if s.FoodUnits != foodBefore+5 {
		t.Errorf("food = %d, want %d", s.FoodUnits, foodBefore+5)
	}
	if s.Gold != 50-5*FoodCostPerUnit {
		t.Errorf("gold = %d, want %d", s.Gold, 50-5*FoodCostPerUnit)
	}
}

func TestBuyFood_NotEnoughGold(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)
	s.Gold = 3 // can afford 1 unit at 2 gold each
	foodBefore := s.FoodUnits

	msgs := BuyFood(s, 10)

	if len(msgs) == 0 {
		t.Error("BuyFood should return messages")
	}
	// Should have bought only what was affordable
	if s.FoodUnits <= foodBefore {
		t.Error("should have bought at least 1 food unit")
	}
	if s.Gold < 0 {
		t.Error("gold should not go negative")
	}
}

func TestBuyFood_NoGold(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)
	s.Gold = 0

	msgs := BuyFood(s, 5)

	if len(msgs) == 0 {
		t.Error("BuyFood should return a message")
	}
	// Should not have bought anything
	if s.FoodUnits != 14 { // starting food
		t.Errorf("food changed unexpectedly: %d", s.FoodUnits)
	}
}

func TestBuyFood_NoSettlement(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(5, 5) // wilderness
	s.Gold = 100

	msgs := BuyFood(s, 5)

	if len(msgs) == 0 {
		t.Error("should return an error message")
	}
}

func TestHuntForFood_Success(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(8, 5) // forest hex - huntable
	s.Prince.CombatSkill = 20     // CS 20 + endurance/2 always exceeds 2d6 max of 12

	foodBefore := s.FoodUnits
	gained, msgs := HuntForFood(s)

	if gained == 0 {
		t.Error("hunt should succeed with CS 20 in forest")
	}
	if s.FoodUnits != foodBefore+gained {
		t.Errorf("food not updated correctly: got %d, want %d", s.FoodUnits, foodBefore+gained)
	}
	if len(msgs) == 0 {
		t.Error("should return messages")
	}
}

func TestHuntForFood_Farmland(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // Ogon - farmland/town, can't hunt

	gained, _ := HuntForFood(s)
	if gained != 0 {
		t.Error("should not be able to hunt in farmland/town")
	}
}

func TestLodgingCost(t *testing.T) {
	if LodgingCost(1) != 1 {
		t.Errorf("LodgingCost(1) = %d, want 1", LodgingCost(1))
	}
	if LodgingCost(4) != 4 {
		t.Errorf("LodgingCost(4) = %d, want 4", LodgingCost(4))
	}
}

func TestHealRest_Settlement(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1) // settlement
	s.Prince.Wounds = 3

	msgs := HealRest(s)

	if s.Prince.Wounds != 2 {
		t.Errorf("wounds = %d, want 2 after resting in settlement", s.Prince.Wounds)
	}
	if len(msgs) == 0 {
		t.Error("should return a message when healing occurs")
	}
}

func TestHealRest_NoWounds(t *testing.T) {
	s := NewGameState()
	s.CurrentHex = NewHexID(1, 1)
	s.Prince.Wounds = 0

	msgs := HealRest(s)

	if s.Prince.Wounds != 0 {
		t.Error("wounds should stay 0 when already healthy")
	}
	if len(msgs) != 0 {
		t.Error("should return no message when not healing")
	}
}
