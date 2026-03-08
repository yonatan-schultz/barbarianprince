package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"barbarianprince/game"
)

// uiPhase is the top-level UI phase, separate from game.TurnPhase
type uiPhase int

const (
	uiPhaseStartMenu uiPhase = iota // title / new-game / load screen
	uiPhasePlaying                  // normal gameplay
	uiPhaseNotes                    // field notes overlay
)

// Model is the bubbletea model
type Model struct {
	ui           uiPhase
	state        *game.GameState
	width        int
	height       int
	travelIndex  int               // selected hex when traveling
	neighbors    []game.HexID      // current travel options
	eventResult  *game.EventResult // pending event with choices
	combatDone   bool
	saveMsg      string            // transient feedback message (save ok/err)
	startOpts    []string          // start menu options
	startIndex   int               // cursor on start menu
	notesScroll  int               // scroll offset for field notes panel
	tutorialStep int               // -1 = no tutorial; 0..N = active slide
	quitConfirm  bool             // true = waiting for second Q to confirm quit
}

// NewModel creates a new UI model starting at the title screen
func NewModel() Model {
	opts := []string{"New Game", "Tutorial"}
	if game.SaveExists() {
		opts = append(opts, "Continue")
	}
	return Model{
		ui:           uiPhaseStartMenu,
		startOpts:    opts,
		tutorialStep: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global quit
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	// Start menu
	if m.ui == uiPhaseStartMenu {
		return m.handleStartMenuKey(key)
	}

	// Field notes overlay
	if m.ui == uiPhaseNotes {
		switch key {
		case "up", "k":
			if m.notesScroll > 0 {
				m.notesScroll--
			}
		case "down", "j":
			m.notesScroll++
		default:
			m.ui = uiPhasePlaying
			m.notesScroll = 0
		}
		return m, nil
	}

	// Tutorial overlay — block all game input until slides are dismissed
	if m.tutorialStep >= 0 && m.tutorialStep < len(tutorialSteps) {
		switch key {
		case " ", "enter":
			m.tutorialStep++
		case "esc":
			m.tutorialStep = len(tutorialSteps) // skip to end
		}
		return m, nil
	}

	if key == "q" && m.state.Phase != game.PhaseTravel {
		if m.quitConfirm {
			return m, tea.Quit
		}
		m.quitConfirm = true
		m.state.AddLog("Press [Q] again to quit, or any other key to cancel.")
		return m, nil
	}
	// Any non-Q key cancels the quit confirmation
	if m.quitConfirm {
		m.quitConfirm = false
	}

	// Game over
	if m.state.Phase == game.PhaseGameOver {
		if key == "enter" {
			// Return to start menu
			opts := []string{"New Game", "Tutorial"}
			if game.SaveExists() {
				opts = append(opts, "Continue")
			}
			m.ui = uiPhaseStartMenu
			m.startOpts = opts
			m.startIndex = 0
			m.state = nil
			m.eventResult = nil
			m.tutorialStep = -1
		}
		return m, nil
	}

	// Clear transient save message on any keypress
	m.saveMsg = ""

	// Combat phase
	if m.state.Phase == game.PhaseCombat {
		return m.handleCombatKey(key)
	}

	// Event resolve / choices
	if m.eventResult != nil && len(m.state.PendingChoices) > 0 {
		return m.handleChoiceKey(key)
	}

	// Travel phase
	if m.state.Phase == game.PhaseTravel {
		return m.handleTravelKey(key)
	}

	// Normal action phase
	return m.handleActionKey(key)
}

func (m Model) handleActionKey(key string) (tea.Model, tea.Cmd) {
	s := m.state

	switch key {
	case "t", "T":
		// Mounted party gets 2 hops (r204a); rope & grapnel also grants 2 hops in mountains
		m.neighbors = game.AdjacentHexes(s.CurrentHex)
		m.travelIndex = 0
		hops := 1
		if s.AllMounted() {
			hops = 2
		} else if s.Prince.HasPossession(game.PossRopeAndGrapnel) {
			hex := game.GetHex(s.CurrentHex)
			if hex != nil && hex.Terrain == game.Mountains {
				hops = 2
				s.AddLog("Your rope & grapnel lets you scale the peaks twice today.")
			}
		}
		s.RemainingTravelHops = hops
		s.Phase = game.PhaseTravel
		return m, nil

	case "r", "R":
		game.ExecuteAction(s, game.ActionRest)

	case "n", "N":
		if er := game.ExecuteAction(s, game.ActionSeekNews); er != nil {
			return m.storeActionResult(er), nil
		}

	case "h", "H":
		if er := game.ExecuteAction(s, game.ActionSeekFollowers); er != nil {
			return m.storeActionResult(er), nil
		}

	case "b", "B":
		game.ExecuteAction(s, game.ActionBuyFood)

	case "a", "A":
		if er := game.ExecuteAction(s, game.ActionSeekAudience); er != nil {
			return m.storeActionResult(er), nil
		}

	case "o", "O":
		if er := game.ExecuteAction(s, game.ActionSubmitOffering); er != nil {
			return m.storeActionResult(er), nil
		}

	case "s", "S":
		if er := game.ExecuteAction(s, game.ActionSearchRuins); er != nil {
			return m.storeActionResult(er), nil
		}

	case "c", "C":
		game.ExecuteAction(s, game.ActionSearchCache)

	case "u", "U":
		if er := game.ExecuteAction(s, game.ActionUseItem); er != nil {
			return m.storeActionResult(er), nil
		}

	case "g", "G":
		game.ExecuteAction(s, game.ActionHunt)

	case "p", "P":
		game.ExecuteAction(s, game.ActionBuyRaft)

	case "f", "F":
		m.ui = uiPhaseNotes
		return m, nil

	case "ctrl+s":
		if err := game.Save(s); err != nil {
			m.saveMsg = "Save failed: " + err.Error()
		} else {
			m.saveMsg = "Game saved."
		}

	case "w", "W":
		s.AddLog("=== PARTY ===")
		s.AddLog(fmt.Sprintf("Cal Arath — CS:%d  HP:%d/%d  W&W:%d",
			s.Prince.EffectiveCombatSkill(),
			s.Prince.CurrentEndurance(), s.Prince.MaxEndurance,
			s.EffectiveWitWiles()))
		if len(s.Party) == 0 {
			s.AddLog("  (travelling alone)")
		}
		for _, f := range s.Party {
			mount := ""
			if f.HasMount {
				mount = " [mounted]"
			}
			wage := ""
			if f.DailyWage > 0 {
				wage = fmt.Sprintf(" wage:%dg/day", f.DailyWage)
			}
			s.AddLog(fmt.Sprintf("  %s — CS:%d  HP:%d/%d  Morale:%d%s%s",
				f.Name, f.EffectiveCombatSkill(),
				f.CurrentEndurance(), f.MaxEndurance,
				f.Morale, mount, wage))
		}

	case "?":
		s.AddLog("=== HELP ===")
		s.AddLog("[T]ravel [R]est [N]ews [H]ire [B]uy food")
		s.AddLog("[A]udience [O]ffering [S]earch ruins [C]ache")
		s.AddLog("[U]se Item  [G]o Hunting  [P]urchase Raft")
		s.AddLog("[W]ho's in party  [F]ield Notes")
		s.AddLog("[Ctrl+S] Save  [Q] Quit (confirm)")
		s.AddLog("Goal: 500 gold north of Tragoth, Royal Helm to")
		s.AddLog("  Ogon/Weshor, noble ally north, or Staff of Command")
		s.AddLog("Tragoth crossing: road bridge at col 8 (free) or raft (15g)")

	}

	return m, nil
}

func (m Model) handleStartMenuKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		m.startIndex--
		if m.startIndex < 0 {
			m.startIndex = len(m.startOpts) - 1
		}
	case "down", "j":
		m.startIndex++
		if m.startIndex >= len(m.startOpts) {
			m.startIndex = 0
		}
	case "1":
		m.startIndex = 0
		return m.startMenuConfirm()
	case "2":
		if len(m.startOpts) > 1 {
			m.startIndex = 1
			return m.startMenuConfirm()
		}
	case "enter", " ":
		return m.startMenuConfirm()
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) startMenuConfirm() (tea.Model, tea.Cmd) {
	choice := m.startOpts[m.startIndex]
	switch choice {
	case "New Game":
		m.state = game.NewGameState()
		m.ui = uiPhasePlaying
		m.tutorialStep = -1
		m.eventResult = nil
		m.saveMsg = ""
	case "Tutorial":
		m.state = game.NewGameState()
		m.ui = uiPhasePlaying
		m.tutorialStep = 0
		m.eventResult = nil
		m.saveMsg = ""
	case "Continue":
		s, err := game.Load()
		if err != nil {
			m.saveMsg = "Load failed: " + err.Error()
			return m, nil
		}
		m.state = s
		m.ui = uiPhasePlaying
		m.eventResult = nil
		m.saveMsg = ""
		m.state.AddLog("Game loaded. Welcome back, Cal Arath.")
	}
	return m, nil
}

func (m Model) handleTravelKey(key string) (tea.Model, tea.Cmd) {
	s := m.state

	switch key {
	case "1", "2", "3", "4", "5", "6":
		idx := int(key[0]-'0') - 1
		if idx >= 0 && idx < len(m.neighbors) {
			m.travelIndex = idx
		}

	case "up", "k":
		m.travelIndex--
		if m.travelIndex < 0 {
			m.travelIndex = len(m.neighbors) - 1
		}

	case "down", "j":
		m.travelIndex++
		if m.travelIndex >= len(m.neighbors) {
			m.travelIndex = 0
		}

	case "esc", "q":
		// Esc during second hop still ends the day
		if s.RemainingTravelHops < 2 { // already used first hop
			game.AdvanceDay(s)
		}
		s.Phase = game.PhaseActionSelect
		return m, nil

	case "enter", " ":
		if m.travelIndex >= 0 && m.travelIndex < len(m.neighbors) {
			target := m.neighbors[m.travelIndex]
			previousHex := s.CurrentHex
			result := game.DoTravel(s, target)
			for _, msg := range result.Messages {
				s.AddLog(msg)
			}

			s.RemainingTravelHops--

			// If successful and hops remain (mounted second hop), offer another move
			if result.Success && s.RemainingTravelHops > 0 && !result.HasEvent {
				s.AddLog("You may move again. [Enter] to continue or [Esc] to make camp.")
				m.neighbors = game.AdjacentHexes(s.CurrentHex)
				m.travelIndex = 0
				return m, nil
			}

			s.Phase = game.PhaseActionSelect

			if result.HasEvent {
				// Dispatch the event
				ctx := game.EventContext{}
				if t := game.GetHex(s.CurrentHex); t != nil {
					ctx.TriggerTerrain = t.Terrain
				}
				evResult := game.DispatchEvent(s, result.EventID, ctx)
				for _, msg := range evResult.Messages {
					s.AddLog(msg)
				}

				if evResult.BlocksTravel {
					// Undo movement — player stays in previous hex
					s.CurrentHex = previousHex
					game.AdvanceDay(s)
					return m, nil
				}

				if evResult.CombatTriggered && evResult.Enemy != nil {
					s.ActiveEnemy = evResult.Enemy
					s.PlayerAttacks = evResult.PlayerAttFirst
					s.CombatLog = nil
					s.Phase = game.PhaseCombat
					return m, nil
				}

				if len(evResult.Choices) > 0 {
					// Apply non-choice effects (Notes, immediate gold/food) before pausing for choice
					for _, msg := range game.ApplyEventResultToState(s, evResult) {
						s.AddLog(msg)
					}
					m.eventResult = &evResult
					s.PendingChoices = evResult.Choices
					return m, nil
				}

				// Apply changes
				for _, msg := range game.ApplyEventResultToState(s, evResult) {
					s.AddLog(msg)
				}
			}

			game.AdvanceDay(s)
		}
	}

	return m, nil
}

func (m Model) handleChoiceKey(key string) (tea.Model, tea.Cmd) {
	s := m.state

	if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
		idx := int(key[0]-'0') - 1
		if idx >= 0 && idx < len(s.PendingChoices) && m.eventResult != nil {
			advanceDay := m.eventResult.AdvanceDayOnChoice
			if m.eventResult.ChoiceHandler != nil {
				result := m.eventResult.ChoiceHandler(s, idx)
				for _, msg := range result.Messages {
					s.AddLog(msg)
				}
				if result.CombatTriggered && result.Enemy != nil {
					s.ActiveEnemy = result.Enemy
					s.PlayerAttacks = result.PlayerAttFirst
					s.CombatLog = nil
					s.Phase = game.PhaseCombat
					s.PendingChoices = nil
					m.eventResult = nil
					return m, nil
				}
				// Nested choices: ChoiceHandler itself returned choices (e.g. wizard enchant)
				if len(result.Choices) > 0 && result.ChoiceHandler != nil {
					result.AdvanceDayOnChoice = advanceDay // propagate day-advance flag
					m.eventResult = &result
					s.PendingChoices = result.Choices
					return m, nil
				}
				for _, msg := range game.ApplyEventResultToState(s, result) {
					s.AddLog(msg)
				}
			}
			s.PendingChoices = nil
			m.eventResult = nil
			if advanceDay {
				game.AdvanceDay(s)
			}
		}
	}

	return m, nil
}

func (m Model) handleCombatKey(key string) (tea.Model, tea.Cmd) {
	s := m.state

	switch key {
	case "f", "F", "enter":
		if s.ActiveEnemy == nil {
			s.Phase = game.PhaseActionSelect
			return m, nil
		}

		// Execute one combat round
		msgs, over, result := game.ResolveCombatRound(s, s.ActiveEnemy, s.PlayerAttacks)
		// Append live HP status after each round
		if !over {
			msgs = append(msgs, fmt.Sprintf("  You: %d/%d HP  |  %s: %d/%d HP",
				s.Prince.CurrentEndurance(), s.Prince.MaxEndurance,
				s.ActiveEnemy.Name, s.ActiveEnemy.CurrentEndurance(), s.ActiveEnemy.MaxEndurance))
		}
		s.CombatLog = msgs
		for _, msg := range msgs {
			s.AddLog(msg)
		}

		if over {
			s.ActiveEnemy = nil
			s.Phase = game.PhaseActionSelect
			s.CombatLog = nil
			m.combatDone = true

			if result.PlayerWon {
				if result.EnemyFled {
					s.AddLog("The enemy flees! You find nothing of value.")
				} else {
					s.AddLog(fmt.Sprintf("Victory! Gained %d gold.", result.LootGold))
				}
				// Award pending duel prize (e151) only on victory
				if s.PendingDuelGold > 0 {
					s.Gold += s.PendingDuelGold
					s.AddLog(fmt.Sprintf("The lord awards you the promised prize: %d gold!", s.PendingDuelGold))
					s.PendingDuelGold = 0
				}
			} else {
				// Clear any pending duel prize if player lost
				s.PendingDuelGold = 0
			}

			// Check if prince died
			if s.Prince.IsDead() {
				s.Phase = game.PhaseGameOver
				s.LoseReason = "Cal Arath has fallen in battle. His quest ends here."
				return m, nil // don't advance day after death
			}

			// Unconsciousness follower check (r221b)
			for _, msg := range game.CheckUnconsciousFollowers(s) {
				s.AddLog(msg)
			}

			// Advance day after combat
			game.AdvanceDay(s)
		}

	case "r", "R":
		// Flee
		success, msg := game.AttemptFlee(s, s.ActiveEnemy)
		s.AddLog(msg)
		// Check death first — flee attack may have killed the prince
		if s.Prince.IsDead() {
			s.ActiveEnemy = nil
			s.Phase = game.PhaseGameOver
			s.LoseReason = "Cal Arath has fallen in battle. His quest ends here."
			return m, nil
		}
		if success {
			s.ActiveEnemy = nil
			s.Phase = game.PhaseActionSelect
			s.CombatLog = nil
			game.AdvanceDay(s)
		} else {
			s.CombatLog = []string{msg}
		}
	}

	return m, nil
}

func (m Model) checkPendingEvent() {}

// storeActionResult sets up combat or choice state from a pending action EventResult.
func (m Model) storeActionResult(er *game.EventResult) Model {
	s := m.state
	if er.CombatTriggered && er.Enemy != nil {
		s.ActiveEnemy = er.Enemy
		s.PlayerAttacks = er.PlayerAttFirst
		s.CombatLog = nil
		s.Phase = game.PhaseCombat
		return m
	}
	if len(er.Choices) > 0 {
		m.eventResult = er
		s.PendingChoices = er.Choices
	}
	return m
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}
	if m.width < 60 || m.height < 16 {
		return fmt.Sprintf("Terminal too small (%dx%d). Please resize to at least 60×16.", m.width, m.height)
	}

	if m.ui == uiPhaseStartMenu {
		return m.renderStartMenu()
	}

	s := m.state

	// Layout: map on left (left 60%), log+menu on right (right 40%)
	rightWidth := 38
	if m.width >= 120 {
		rightWidth = 44
	}
	mapWidth := m.width - rightWidth - 2 // -2 for the gap between panels

	// statusBarHeight = 1 content line + 2 border rows = 3 total rendered lines
	// mapPanel border adds 2 rows, so inner height = total - statusBar - 2 borders
	statusBarHeight := 3
	mapHeight := m.height - statusBarHeight - 2

	// Render map — show travel target info in travel mode, current hex info otherwise
	var infoHex game.HexID
	if s.Phase == game.PhaseTravel && m.travelIndex >= 0 && m.travelIndex < len(m.neighbors) {
		infoHex = m.neighbors[m.travelIndex]
	} else {
		infoHex = s.CurrentHex
	}
	mapContent := RenderMap(s, mapWidth, mapHeight, infoHex)
	mapPanel := StyleBorder.
		Width(mapWidth).
		Height(mapHeight).
		Render(mapContent)

	// Right panel: log (top half) + menu (bottom half).
	// Outer height of each panel = inner + 2 (borders).
	// For right panel to match map panel: logH+2 + menuH+2 = mapH+2
	// → logH + menuH = mapH - 2
	logHeight := (mapHeight - 2) / 2
	if logHeight < 4 {
		logHeight = 4
	}
	menuHeight := mapHeight - 2 - logHeight
	if menuHeight < 3 {
		menuHeight = 3
	}

	logInnerWidth := rightWidth - 4 // subtract 2 border chars on each side
	logContent := RenderLog(s.Log, logHeight, logInnerWidth)
	logPanel := StyleBorder.
		Width(rightWidth - 2).
		Height(logHeight).
		Render(logContent)

	// Menu content depends on phase
	var menuContent string
	if m.tutorialStep >= 0 && m.tutorialStep < len(tutorialSteps) {
		// rightWidth-2 for border, -2 again for inner padding = rightWidth-4
		menuContent = RenderTutorial(m.tutorialStep, rightWidth-4)
	} else if m.ui == uiPhaseNotes {
		menuContent = RenderNotes(s.FieldNotes, menuHeight, m.notesScroll)
	} else if s.Phase == game.PhaseTravel {
		menuContent = RenderAdjacentHexes(s, m.travelIndex)
	} else if s.Phase == game.PhaseCombat {
		menuContent = RenderCombatMenu(s)
	} else if len(s.PendingChoices) > 0 {
		menuContent = RenderChoices(s)
	} else if s.Phase == game.PhaseGameOver {
		menuContent = RenderGameOver(s)
	} else {
		menuContent = RenderActionMenuFull(s)
	}
	// Clip menu content — lipgloss Height() pads but does NOT clip, so we must
	// truncate the content ourselves before handing it to the border renderer.
	menuContent = clipLines(menuContent, menuHeight)

	menuPanel := StyleBorder.
		Width(rightWidth - 2).
		Height(menuHeight).
		Render(menuContent)

	rightPanel := lipgloss.JoinVertical(lipgloss.Left, logPanel, menuPanel)

	// Status bar at bottom (includes transient save message)
	statusContent := RenderStatus(s, m.width-4)
	if m.saveMsg != "" {
		statusContent += "  " + StyleSuccess.Render(m.saveMsg)
	}
	statusBar := StyleBorder.
		Width(m.width - 4).
		Render(statusContent)

	// Main content: map + right
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, mapPanel, rightPanel)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, statusBar)
}

func (m Model) renderStartMenu() string {
	title := []string{
		"",
		StyleTitle.Render("  ██████╗  █████╗ ██████╗ ██████╗  █████╗ ██████╗ ██╗ █████╗ ███╗  ██╗"),
		StyleTitle.Render("  ██╔══██╗██╔══██╗██╔══██╗██╔══██╗██╔══██╗██╔══██╗██║██╔══██╗████╗ ██║"),
		StyleTitle.Render("  ██████╔╝███████║██████╔╝██████╔╝███████║██████╔╝██║███████║██╔██╗██║"),
		StyleTitle.Render("  ██╔══██╗██╔══██║██╔══██╗██╔══██╗██╔══██║██╔══██╗██║██╔══██║██║╚████║"),
		StyleTitle.Render("  ██████╔╝██║  ██║██║  ██║██████╔╝██║  ██║██║  ██║██║██║  ██║██║ ╚███║"),
		StyleTitle.Render("  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚═╝  ╚═╝"),
		"",
		StyleLabel.Render("  ██████╗ ██████╗ ██╗███╗  ██╗ ██████╗███████╗"),
		StyleLabel.Render("  ██╔══██╗██╔══██╗██║████╗ ██║██╔════╝██╔════╝"),
		StyleLabel.Render("  ██████╔╝██████╔╝██║██╔██╗██║██║     █████╗  "),
		StyleLabel.Render("  ██╔═══╝ ██╔══██╗██║██║╚████║██║     ██╔══╝  "),
		StyleLabel.Render("  ██║     ██║  ██║██║██║ ╚███║╚██████╗███████╗"),
		StyleLabel.Render("  ╚═╝     ╚═╝  ╚═╝╚═╝╚═╝  ╚══╝ ╚═════╝╚══════╝"),
		"",
		StyleMuted.Render("  Cal Arath, deposed prince. 70 days. 500 gold. Reclaim your throne."),
		"",
	}

	for i, opt := range m.startOpts {
		if i == m.startIndex {
			title = append(title, StyleMenuKey.Render("  > ")+StyleValue.Render(opt))
		} else {
			title = append(title, StyleLabel.Render("    ")+StyleMenuText.Render(opt))
		}
	}

	title = append(title, "")
	title = append(title, StyleLabel.Render("  [↑↓] navigate  [Enter] select  [Q] quit"))

	if m.saveMsg != "" {
		title = append(title, "")
		title = append(title, StyleDanger.Render("  "+m.saveMsg))
	}

	return strings.Join(title, "\n")
}

// RenderActionMenuFull renders the full action menu
func RenderActionMenuFull(state *game.GameState) string {
	var lines []string

	hex := game.GetHex(state.CurrentHex)
	hexDesc := ""
	if hex != nil {
		hexDesc = hex.Terrain.String()
		if hex.Name != "" {
			hexDesc = hex.Name
		}
		if hex.Structure != game.StructNone {
			hexDesc += " (" + structureDesc(hex.Structure) + ")"
		}
	}

	lines = append(lines, StyleTitle.Render("BARBARIAN PRINCE"))
	lines = append(lines, StyleLabel.Render("Location: ")+StyleValue.Render(hexDesc))
	lines = append(lines, "")

	actions := state.AvailableActions()
	for _, a := range actions {
		key := a.ActionKey()
		str := a.String()
		lines = append(lines,
			"  "+StyleMenuKey.Render("["+strings.ToUpper(key)+"]")+
				StyleMenuText.Render(" "+str[len("["+strings.ToUpper(key)+"]"):]))
	}

	lines = append(lines, "")
	lines = append(lines, StyleLabel.Render("[?] Help  [W]ho  [F]ield Notes  [Ctrl+S] Save  [Q] Quit"))
	lines = append(lines, "")
	lines = append(lines, StyleLabel.Render("── Map Legend ──"))
	lines = append(lines, StyleMenuText.Render("[*] You  >>> Target  ~~~~ River  = Road"))
	lines = append(lines, StyleMenuText.Render("Named places show first 3 letters"))
	lines = append(lines, StyleMenuText.Render("[T] Town  [C] Castle  [K] Keep"))
	lines = append(lines, StyleMenuText.Render("[R] Ruins  (R) Searched  [+] Temple"))
	lines = append(lines, StyleMenuText.Render("[v] Village   $  Cache hidden here"))
	lines = append(lines, StyleMuted.Render(". Farm  ~ Country  f Forest  ^ Hills"))
	lines = append(lines, StyleMuted.Render("M Mtn  s Swamp  o Desert"))

	return strings.Join(lines, "\n")
}

// clipLines truncates s to at most maxLines lines, joining with newline.
func clipLines(s string, maxLines int) string {
	if maxLines <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	return strings.Join(lines[:maxLines], "\n")
}

func structureDesc(s game.StructureType) string {
	switch s {
	case game.StructTown:
		return "Town"
	case game.StructCastle:
		return "Castle"
	case game.StructTemple:
		return "Temple"
	case game.StructRuins:
		return "Ruins"
	case game.StructVillage:
		return "Village"
	case game.StructKeep:
		return "Keep"
	}
	return ""
}
