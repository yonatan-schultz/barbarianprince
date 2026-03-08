package game

import "fmt"

// EventContext provides context for event resolution
type EventContext struct {
	TriggerTerrain TerrainType
	SurpriseCode   string
}

// EventResult is returned by event handlers
type EventResult struct {
	Messages       []string
	CombatTriggered bool
	Enemy          *Character
	PlayerAttFirst  bool
	Choices        []string    // if non-empty, pause for player choice
	ChoiceHandler  func(*GameState, int) EventResult
	GoldChange     int
	FoodChange     int
	NewFollower    *Character
	GameOver       bool
	Note              string // if non-empty, auto-added to FieldNotes
	BlocksTravel      bool   // if true, travel to this hex is cancelled
	AdvanceDayOnChoice bool  // if true, advance day after player resolves a choice
}

// EventHandler is the function signature for all events
type EventHandler func(s *GameState, ctx EventContext) EventResult

// eventRegistry maps EventID to handler functions
var eventRegistry map[EventID]EventHandler

func init() {
	eventRegistry = make(map[EventID]EventHandler)
	registerEventsE001()
	registerEventsE051()
	registerEventsE100()
	registerEventsE110()
	registerEventsE180()
}

// RegisterEvent adds an event to the registry
func RegisterEvent(id EventID, handler EventHandler) {
	eventRegistry[id] = handler
}

// DispatchEvent runs an event by ID
func DispatchEvent(s *GameState, id EventID, ctx EventContext) EventResult {
	handler, ok := eventRegistry[id]
	if !ok {
		// Stub: unknown event
		return EventResult{
			Messages: []string{fmt.Sprintf("Nothing of note occurs. (event %s)", id)},
		}
	}
	return handler(s, ctx)
}

// ApplyEventResultToState is the exported version of applyEventResult.
// Returns log messages for the caller to add to the game log.
func ApplyEventResultToState(s *GameState, result EventResult) []string {
	return applyEventResult(s, result)
}

// applyEventResult applies gold/food changes and follower additions from an event.
// Returns log messages for the caller to add to the game log.
func applyEventResult(s *GameState, result EventResult) []string {
	var msgs []string
	if result.GoldChange != 0 {
		s.Gold += result.GoldChange
		if result.GoldChange > 0 {
			msgs = append(msgs, fmt.Sprintf("Gained %d gold. Total: %d", result.GoldChange, s.Gold))
		} else {
			msgs = append(msgs, fmt.Sprintf("Lost %d gold. Total: %d", -result.GoldChange, s.Gold))
		}
	}
	if result.FoodChange != 0 {
		s.FoodUnits += result.FoodChange
		if result.FoodChange > 0 {
			msgs = append(msgs, fmt.Sprintf("Gained %d food units.", result.FoodChange))
		} else {
			msgs = append(msgs, fmt.Sprintf("Lost %d food units.", -result.FoodChange))
		}
	}
	if result.NewFollower != nil {
		f := result.NewFollower
		alreadyHired := false
		for _, existing := range s.Party {
			if existing.Name == f.Name {
				alreadyHired = true
				break
			}
		}
		if alreadyHired {
			msgs = append(msgs, fmt.Sprintf("%s is already in your party.", f.Name))
		} else {
			s.AddFollower(*f)
			msgs = append(msgs, fmt.Sprintf("%s joins your party!", f.Name))
			msgs = append(msgs, followerEffectDesc(f))
			if f.DailyWage > 0 {
				msgs = append(msgs, fmt.Sprintf("Daily wage: %d gold.", f.DailyWage))
			}
		}
	}
	if result.Note != "" {
		s.AddNote(result.Note)
	}
	return msgs
}

// followerEffectDesc returns a short description of the mechanical benefit of a follower.
func followerEffectDesc(f *Character) string {
	if f.IsGuide {
		return "As a guide they know these lands well, reducing your chance of getting lost while traveling."
	}
	switch f.Type {
	case TypeSwordsman, TypeMercenary, TypeBandit:
		return fmt.Sprintf("A trained fighter (CS %d, End %d) who will fight alongside you in combat.", f.CombatSkill, f.MaxEndurance)
	case TypeAmazon:
		return fmt.Sprintf("A skilled warrior (CS %d, End %d) who fights independently and shares in any treasure found.", f.CombatSkill, f.MaxEndurance)
	case TypeDwarf:
		return fmt.Sprintf("A sturdy dwarf (CS %d, End %d) who excels in close combat and shares in treasure.", f.CombatSkill, f.MaxEndurance)
	case TypeElf:
		return fmt.Sprintf("An elven companion (CS %d, End %d) with keen senses; adds their Combat Skill to your party.", f.CombatSkill, f.MaxEndurance)
	case TypeHalfling:
		return fmt.Sprintf("A nimble halfling (CS %d, End %d) who can scout ahead and assist in combat.", f.CombatSkill, f.MaxEndurance)
	case TypePriest:
		return fmt.Sprintf("A holy priest (CS %d) who may intercede for you at temples and tends wounds after battle.", f.CombatSkill)
	case TypeWizard:
		return fmt.Sprintf("A powerful wizard (CS %d, W/W %d) whose magic can turn the tide of difficult encounters.", f.CombatSkill, f.WitWiles)
	default:
		if f.CombatSkill > 0 {
			return fmt.Sprintf("Fights alongside you in combat (CS %d, End %d).", f.CombatSkill, f.MaxEndurance)
		}
		return "Accompanies you on your quest."
	}
}
