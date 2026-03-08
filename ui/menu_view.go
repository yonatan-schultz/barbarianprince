package ui

import (
	"fmt"
	"strings"

	"barbarianprince/game"
)

// RenderMenu renders the action menu for the current game state
func RenderMenu(state *game.GameState) string {
	if state.Phase == game.PhaseGameOver {
		return RenderGameOver(state)
	}

	if state.Phase == game.PhaseTravel {
		return RenderAdjacentHexes(state, 0)
	}

	if state.Phase == game.PhaseEventResolve || len(state.PendingChoices) > 0 {
		return RenderChoices(state)
	}

	if state.Phase == game.PhaseCombat {
		return RenderCombatMenu(state)
	}

	return RenderActionMenu(state)
}

// RenderActionMenu renders the main action selection menu
func RenderActionMenu(state *game.GameState) string {
	var lines []string

	hex := game.GetHex(state.CurrentHex)
	hexName := ""
	if hex != nil {
		if hex.Name != "" {
			hexName = hex.Name
		} else {
			hexName = hex.Terrain.String()
		}
	}

	lines = append(lines, StyleTitle.Render(fmt.Sprintf("[ %s ] Day %d/70", hexName, state.Day)))
	lines = append(lines, "")

	actions := state.AvailableActions()
	for _, a := range actions {
		key := strings.ToUpper(a.ActionKey())
		lines = append(lines, "  "+StyleMenuKey.Render("["+key+"]")+" "+StyleMenuText.Render(a.String()))
	}

	lines = append(lines, "")
	lines = append(lines, StyleLabel.Render("[?] Help  [Q] Quit"))

	return strings.Join(lines, "\n")
}

// RenderChoices renders event choices for the player
func RenderChoices(state *game.GameState) string {
	var lines []string
	lines = append(lines, StyleTitle.Render("Choose your action:"))
	lines = append(lines, "")

	for i, choice := range state.PendingChoices {
		key := fmt.Sprintf("%d", i+1)
		lines = append(lines, StyleMenuKey.Render("["+key+"] ")+StyleMenuText.Render(choice))
	}

	lines = append(lines, "")
	lines = append(lines, StyleLabel.Render("[1-9] select"))

	return strings.Join(lines, "\n")
}

// RenderCombatMenu renders combat options
func RenderCombatMenu(state *game.GameState) string {
	var lines []string

	if state.ActiveEnemy != nil {
		e := state.ActiveEnemy
		lines = append(lines, StyleDanger.Render(fmt.Sprintf("COMBAT: %s", e.Name)))
		lines = append(lines, fmt.Sprintf("  Enemy CS: %d  Endurance: %d/%d",
			e.EffectiveCombatSkill(), e.CurrentEndurance(), e.MaxEndurance))
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("  Your CS: %d  Endurance: %d/%d",
			state.Prince.EffectiveCombatSkill(),
			state.Prince.CurrentEndurance(), state.Prince.MaxEndurance))
		lines = append(lines, "")
	}

	// Show combat log
	for _, msg := range state.CombatLog {
		lines = append(lines, StyleWarning.Render(msg))
	}
	lines = append(lines, "")

	lines = append(lines, StyleMenuKey.Render("[F] ")+ StyleMenuText.Render("Fight (attack)"))
	lines = append(lines, StyleMenuKey.Render("[R] ")+ StyleMenuText.Render("Retreat (attempt to flee)"))

	return strings.Join(lines, "\n")
}

// RenderGameOver renders the win or lose screen
func RenderGameOver(state *game.GameState) string {
	var lines []string
	lines = append(lines, "")
	lines = append(lines, "")

	followers := len(state.Party)
	notes := len(state.FieldNotes)

	if state.WinReason != "" {
		lines = append(lines, StyleSuccess.Render("╔══════════════════════╗"))
		lines = append(lines, StyleSuccess.Render("║     VICTORY!         ║"))
		lines = append(lines, StyleSuccess.Render("╚══════════════════════╝"))
		lines = append(lines, "")
		lines = append(lines, StyleSuccess.Render(state.WinReason))
		lines = append(lines, "")
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Days survived:    %d / 70", state.Day)))
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Gold:             %d", state.Gold)))
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Followers:        %d", followers)))
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Field notes:      %d", notes)))
	} else {
		lines = append(lines, StyleDanger.Render("╔══════════════════════╗"))
		lines = append(lines, StyleDanger.Render("║     GAME OVER        ║"))
		lines = append(lines, StyleDanger.Render("╚══════════════════════╝"))
		lines = append(lines, "")
		lines = append(lines, StyleDanger.Render(state.LoseReason))
		lines = append(lines, "")
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Days survived:    %d / 70", state.Day-1)))
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Gold:             %d", state.Gold)))
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Followers:        %d", followers)))
		lines = append(lines, StyleValue.Render(fmt.Sprintf("Field notes:      %d", notes)))
	}

	lines = append(lines, "")
	lines = append(lines, StyleLabel.Render("[Q] Quit  [Enter] New Game"))

	return strings.Join(lines, "\n")
}

// RenderStatus renders the status bar. maxWidth is the available character width
// (including ANSI sequences the border will consume); items are dropped from the
// right when the bar would overflow a narrow terminal.
func RenderStatus(state *game.GameState, maxWidth int) string {
	prince := &state.Prince

	// Day progress bar (omitted on very narrow terminals, < 60 cols)
	dayFrac := float64(state.Day) / 70.0
	barWidth := 20
	filled := int(dayFrac * float64(barWidth))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	// Endurance bar
	endFrac := float64(prince.CurrentEndurance()) / float64(prince.MaxEndurance)
	endFilled := int(endFrac * float64(10))
	if endFilled < 0 {
		endFilled = 0
	}
	endBar := strings.Repeat("♥", endFilled) + strings.Repeat("♡", 10-endFilled)

	var parts []string

	// Day counter: highlight in red when ≤ 10 days remain
	daysLeft := 70 - state.Day + 1
	dayStr := fmt.Sprintf("%d/70", state.Day)
	if daysLeft <= 10 {
		parts = append(parts, StyleLabel.Render("Day: ")+StyleDanger.Render(dayStr+" !!"))
	} else {
		parts = append(parts, StyleLabel.Render("Day: ")+StyleValue.Render(dayStr))
	}
	if maxWidth >= 60 {
		parts = append(parts, StyleLabel.Render(" [")+StyleWarning.Render(bar)+StyleLabel.Render("]"))
	}
	parts = append(parts, StyleLabel.Render("  Gold: ")+StyleGold.Render(fmt.Sprintf("%d", state.Gold)))

	// Win progress hints
	switch {
	case prince.HasPossession(game.PossRoyalHelm) || prince.HasPossession(game.PossGoldenCrown):
		parts = append(parts, StyleSuccess.Render(" → Return to Ogon/Weshor!"))
	case state.Flags.NobleAllySecured && !game.IsNorthOfTragoth(state.CurrentHex):
		parts = append(parts, StyleSuccess.Render(" → Go north! (noble ally)"))
	case prince.HasPossession(game.PossStaffOfCommand) && !game.IsNorthOfTragoth(state.CurrentHex):
		parts = append(parts, StyleSuccess.Render(" → Go north! (staff)"))
	case state.Gold >= 500 && !game.IsNorthOfTragoth(state.CurrentHex):
		parts = append(parts, StyleSuccess.Render(" GO NORTH!"))
	case state.Gold < 500:
		parts = append(parts, StyleMuted.Render(fmt.Sprintf("(need %d)", 500-state.Gold)))
	}

	parts = append(parts, StyleLabel.Render("  Food: ")+StyleValue.Render(fmt.Sprintf("%d", state.FoodUnits)))
	parts = append(parts, StyleLabel.Render("  CS: ")+StyleValue.Render(fmt.Sprintf("%d", prince.EffectiveCombatSkill())))
	parts = append(parts, StyleLabel.Render("  HP: ")+StyleSuccess.Render(endBar))

	if prince.StarvationDays > 0 {
		parts = append(parts, StyleDanger.Render(fmt.Sprintf("  STARVING:%d", prince.StarvationDays)))
	}
	if prince.PoisonWounds > 0 {
		parts = append(parts, StyleDanger.Render(fmt.Sprintf("  POISON:%d", prince.PoisonWounds)))
	}
	if prince.HasPossession(game.PossRaft) {
		parts = append(parts, StyleSuccess.Render("  Raft"))
	}

	// Followers and items: compact counts on medium terminals, full labels on wide ones
	if maxWidth >= 80 {
		if len(state.Party) > 0 {
			if maxWidth >= 100 {
				parts = append(parts, StyleLabel.Render(fmt.Sprintf("  Followers: %d", len(state.Party))))
			} else {
				parts = append(parts, StyleLabel.Render(fmt.Sprintf("  Flw:%d", len(state.Party))))
			}
		}
		// Non-raft items
		var nonRaftPoss []game.PossessionType
		for _, p := range prince.Possessions {
			if p != game.PossRaft {
				nonRaftPoss = append(nonRaftPoss, p)
			}
		}
		if len(nonRaftPoss) > 0 {
			if maxWidth >= 100 {
				poss := "  Items:"
				for _, p := range nonRaftPoss {
					name := game.PossessionName(p)
					if len(name) > 3 {
						name = name[:3]
					}
					poss += " " + name
				}
				parts = append(parts, StyleLabel.Render(poss))
			} else {
				parts = append(parts, StyleLabel.Render(fmt.Sprintf("  Items:%d", len(nonRaftPoss))))
			}
		}
	}

	return strings.Join(parts, "")
}

// RenderNotes renders the field notes panel with scroll support.
func RenderNotes(notes []game.Note, height int, scroll int) string {
	var noteLines []string
	if len(notes) == 0 {
		noteLines = append(noteLines, StyleMuted.Render("No discoveries recorded yet."))
	} else {
		for _, n := range notes {
			noteLines = append(noteLines, StyleLabel.Render(fmt.Sprintf("Day %d [%s]: ", n.Day, n.Hex))+StyleMenuText.Render(n.Text))
		}
	}

	maxScroll := len(noteLines) - 1
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}
	if scroll < 0 {
		scroll = 0
	}

	var lines []string
	lines = append(lines, StyleTitle.Render("── Field Notes ──"))
	lines = append(lines, "")
	lines = append(lines, noteLines[scroll:]...)
	lines = append(lines, "")

	footer := "[any key] close"
	if scroll > 0 {
		footer = "[↑/k] up  " + footer
	}
	if scroll < maxScroll {
		footer = footer + "  [↓/j] down"
	}
	lines = append(lines, StyleLabel.Render(footer))

	return clipLines(strings.Join(lines, "\n"), height)
}

// StyleGold is defined here to avoid circular dependency
var StyleGold = StyleWarning
