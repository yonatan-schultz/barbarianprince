package game

// TutorialStep represents the current in-game guided tutorial step.
type TutorialStep int

const (
	TutStepBuyFood TutorialStep = iota // hint: buy food in starting town
	TutStepTravel                       // hint: travel to an adjacent hex
	TutStepRest                         // hint: rest to recover wounds
	TutStepDone                         // all steps complete
)

// TutorialState tracks in-game guided tutorial progress.
// Nil means no tutorial is active.
type TutorialState struct {
	Step TutorialStep `json:"step"`
}

// NewTutorialState creates a fresh tutorial starting at step 1.
func NewTutorialState() *TutorialState {
	return &TutorialState{Step: TutStepBuyFood}
}

// IsActive reports whether the in-game hints are still showing.
func (t *TutorialState) IsActive() bool {
	return t != nil && t.Step < TutStepDone
}

// Hint returns the current hint line to display above the action menu.
func (t *TutorialState) Hint() string {
	if t == nil || t.Step >= TutStepDone {
		return ""
	}
	const skip = " [X] skip"
	switch t.Step {
	case TutStepBuyFood:
		return "▶ GOAL: Buy food to feed your party — press [B] Buy Food." + skip
	case TutStepTravel:
		return "▶ GOAL: Travel to an adjacent hex — press [T] then choose a direction." + skip
	case TutStepRest:
		return "▶ GOAL: Rest to recover wounds — press [R] Rest." + skip
	}
	return ""
}

// OnAction advances the tutorial when the player executes an action.
func (t *TutorialState) OnAction(a Action) {
	if t == nil || t.Step >= TutStepDone {
		return
	}
	switch t.Step {
	case TutStepBuyFood:
		if a == ActionBuyFood {
			t.Step = TutStepTravel
		}
	case TutStepRest:
		if a == ActionRest {
			t.Step = TutStepDone
		}
	}
}

// OnTravel advances the tutorial after a successful travel action.
func (t *TutorialState) OnTravel() {
	if t == nil || t.Step != TutStepTravel {
		return
	}
	t.Step = TutStepRest
}

// Skip dismisses all remaining tutorial hints.
func (t *TutorialState) Skip() {
	if t != nil {
		t.Step = TutStepDone
	}
}
