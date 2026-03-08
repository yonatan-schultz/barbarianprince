package game

import "fmt"

// FoodCostPerUnit is the gold cost to buy one food unit
const FoodCostPerUnit = 2

// BuyFood attempts to buy food at a settlement
// Returns messages
func BuyFood(s *GameState, units int) []string {
	var msgs []string
	hex := GetHex(s.CurrentHex)
	if hex == nil || !hex.IsSettlement() {
		return []string{"No market available here."}
	}
	cost := units * FoodCostPerUnit
	if s.Gold < cost {
		affordable := s.Gold / FoodCostPerUnit
		if affordable == 0 {
			return []string{"You cannot afford any food."}
		}
		units = affordable
		cost = units * FoodCostPerUnit
		msgs = append(msgs, fmt.Sprintf("You can only afford %d food units.", units))
	}
	s.Gold -= cost
	s.FoodUnits += units
	msgs = append(msgs, fmt.Sprintf("Bought %d food units for %d gold. Food: %d, Gold: %d", units, cost, s.FoodUnits, s.Gold))
	return msgs
}

// HuntForFood attempts to hunt in the current hex (r215b).
// Formula: CS + (current endurance / 2) + guide bonus - 2d6 = food units gained.
// Returns food gained and messages.
func HuntForFood(s *GameState) (int, []string) {
	var msgs []string
	hex := GetHex(s.CurrentHex)
	if hex == nil {
		return 0, []string{"Cannot hunt here."}
	}

	// Can't hunt in towns or settlements with a town/castle/temple structure (r215c)
	if hex.Structure == StructTown || hex.Structure == StructCastle || hex.Structure == StructTemple {
		return 0, []string{"Hunting is prohibited here."}
	}
	// No game in mountains or desert (r207)
	if hex.Terrain == Mountains || hex.Terrain == Desert {
		return 0, []string{"There is no game to hunt in this terrain."}
	}

	roll := Roll2d6()

	// Roll of exactly 12: hunter takes 1d6 wounds regardless of success (r215b)
	if roll == 12 {
		woundRoll := Roll1d6()
		s.Prince.Wounds += woundRoll
		msgs = append(msgs, fmt.Sprintf("An accident during the hunt! You take %d wound(s).", woundRoll))
	}

	// Guide adds +1 to hunt (r215b)
	guideBonus := 0
	if s.HasGuide() {
		guideBonus = 1
	}

	result := s.Prince.EffectiveCombatSkill() + s.Prince.CurrentEndurance()/2 + guideBonus - roll
	if result <= 0 {
		msgs = append(msgs, "You find no game today.")
		return 0, msgs
	}

	s.FoodUnits += result
	msgs = append(msgs, fmt.Sprintf("Successful hunt! Gained %d food unit(s).", result))
	return result, msgs
}

// LodgingCost returns the gold cost for people to stay at a settlement for one night.
// Animal stabling is handled separately in PayLodging (r215f: 1 gold/mount).
func LodgingCost(partySize int) int {
	return partySize // 1 gold per person
}

// PayLodging deducts lodging and stabling costs; handles desertion when sleeping rough (r217).
func PayLodging(s *GameState) []string {
	hex := GetHex(s.CurrentHex)
	if hex == nil || !hex.IsSettlement() {
		return nil
	}
	mounts := s.TotalMounts()
	cost := LodgingCost(s.PartySize()) + mounts // +1 gold per mount for stabling (r215f)
	if s.Gold >= cost {
		s.Gold -= cost
		if mounts > 0 {
			return []string{fmt.Sprintf("Lodging costs %d gold (incl. stabling for %d mount(s)).", cost, mounts)}
		}
		return []string{fmt.Sprintf("Lodging costs %d gold.", cost)}
	}
	// Can't afford lodging — followers may desert (r217): 2d6 - W&W - (Morale-3) >= 4
	msgs := []string{"You cannot afford lodging. The party sleeps rough."}
	var deserters []int
	for i, f := range s.Party {
		if f.IsTrueLove {
			continue // true love never deserts
		}
		roll := Roll2d6() - s.EffectiveWitWiles() - (f.Morale - 3)
		if roll >= 4 {
			deserters = append(deserters, i)
			msgs = append(msgs, fmt.Sprintf("%s deserts after a night sleeping rough!", f.Name))
		}
	}
	for i := len(deserters) - 1; i >= 0; i-- {
		idx := deserters[i]
		s.Party = append(s.Party[:idx], s.Party[idx+1:]...)
	}
	return msgs
}

// HealRest heals wounds during a rest day (staying at settlement)
func HealRest(s *GameState) []string {
	var msgs []string
	hex := GetHex(s.CurrentHex)
	inSettlement := hex != nil && hex.IsSettlement()

	// Heal 1 wound if resting in settlement
	if inSettlement && s.Prince.Wounds > 0 {
		s.Prince.Wounds--
		msgs = append(msgs, "Cal Arath rests and recovers 1 wound.")
	}
	// Heal 1 wound if just resting anywhere (slower)
	if !inSettlement && s.Prince.Wounds > 0 && Roll1d6() >= 5 {
		s.Prince.Wounds--
		msgs = append(msgs, "Cal Arath rests and recovers 1 wound.")
	}

	// Poison slowly clears with rest
	if s.Prince.PoisonWounds > 0 {
		if inSettlement {
			s.Prince.PoisonWounds--
			msgs = append(msgs, "The healers' care reduces the poison's effect. (-1 poison wound)")
		} else if Roll1d6() >= 5 {
			s.Prince.PoisonWounds--
			msgs = append(msgs, "Your body fights off some of the poison. (-1 poison wound)")
		}
	}

	// Followers heal too
	for i := range s.Party {
		if inSettlement && s.Party[i].Wounds > 0 {
			s.Party[i].Wounds--
		}
	}

	return msgs
}
