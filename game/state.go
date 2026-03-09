package game

import "fmt"

// TurnPhase represents the current phase of gameplay
type TurnPhase int

const (
	PhaseActionSelect TurnPhase = iota // Player choosing what to do
	PhaseEventResolve                  // Event text shown, awaiting acknowledgment
	PhaseCombat                        // Combat round in progress
	PhaseTravel                        // Selecting hex to travel to
	PhaseGameOver                      // Win or lose
)

// Note records a significant discovery for the field notes log.
type Note struct {
	Day  int    `json:"day"`
	Hex  HexID  `json:"hex"`
	Text string `json:"text"`
}

// GameState holds all mutable game data
type GameState struct {
	Day        int // 1-70
	Week       int
	DayOfWeek  int // 1-7

	CurrentHex HexID
	Prince     Character
	Party      []Character // followers

	Gold      int
	FoodUnits int // each unit feeds 1 character for 1 day

	Caches      []Cache
	HexFlags    map[HexID]*HexFlags
	Flags       GlobalFlags
	AudienceBarred map[HexID]int // day until unbarred (0 = not barred)

	Log           []string // narrative log
	Phase         TurnPhase
	WinReason     string
	LoseReason    string

	// For multi-step event resolution
	PendingEventID EventID
	PendingChoices []string
	PendingChoice  int // player's selection index

	// For combat
	ActiveEnemy    *Character
	CombatLog      []string
	PlayerAttacks  bool

	// For travel
	SelectedHex HexID
	VisitedHexes map[HexID]bool

	// Pending duel prize — awarded only when the prince wins the duel (e151)
	PendingDuelGold int

	// Travel hops remaining this day (1 = on foot, 2 = mounted)
	RemainingTravelHops int

	// Lodging
	InLodging bool
	LodgingHex HexID

	FieldNotes []Note `json:"fieldNotes,omitempty"`

	// Tutorial tracks in-game guided tutorial progress. Nil when inactive.
	Tutorial *TutorialState `json:"tutorial,omitempty"`
}

// NewGameState creates a fresh game state
func NewGameState() *GameState {
	s := &GameState{
		Day:          1,
		Week:         1,
		DayOfWeek:    1,
		CurrentHex:   NewHexID(1, 1), // Start at Ogon
		Prince:       NewPrince(),
		Gold:         10,
		FoodUnits:    14,
		HexFlags:     make(map[HexID]*HexFlags),
		Flags:        NewGlobalFlags(),
		AudienceBarred: make(map[HexID]int),
		VisitedHexes: make(map[HexID]bool),
		Phase:        PhaseActionSelect,
	}
	s.VisitedHexes[s.CurrentHex] = true
	s.AddLog("Cal Arath, deposed prince of the realm, begins his quest.")
	s.AddLog(fmt.Sprintf("Day 1 of 70. Gold: %d. Food: %d units.", s.Gold, s.FoodUnits))
	s.AddLog("You must accumulate 500 gold within 70 days to reclaim your throne.")
	s.AddLog("You stand in Ogon, the northern city of your homeland.")
	return s
}

// AddLog appends a message to the narrative log
func (s *GameState) AddLog(msg string) {
	s.Log = append(s.Log, msg)
}

// AddNote appends a discovery note to the field notes log.
func (s *GameState) AddNote(text string) {
	s.FieldNotes = append(s.FieldNotes, Note{Day: s.Day, Hex: s.CurrentHex, Text: text})
}

// GetHexFlags returns (or creates) flags for a hex
func (s *GameState) GetHexFlags(id HexID) *HexFlags {
	if s.HexFlags[id] == nil {
		s.HexFlags[id] = &HexFlags{
			EventUsed: make(map[string]bool),
		}
	}
	return s.HexFlags[id]
}

// PartySize returns the total number of people (prince + followers)
func (s *GameState) PartySize() int {
	return 1 + len(s.Party)
}

// TotalMounts returns the number of mounts in the party
func (s *GameState) TotalMounts() int {
	count := 0
	if s.Prince.HasMount {
		count++
	}
	for _, f := range s.Party {
		if f.HasMount {
			count++
		}
	}
	return count
}

// DailyFoodNeeded returns food units needed per day.
// Mounts forage free in farmland/countryside/forest/hills (r215f).
// In a desert hex, cost doubles for men and mounts (r215a).
func (s *GameState) DailyFoodNeeded() int {
	people := s.PartySize()
	hex := GetHex(s.CurrentHex)

	// Mounts forage for free in open terrain (r215f)
	mountFood := 0
	if hex == nil || !mountCanForage(hex.Terrain) {
		mountFood = s.TotalMounts() * 2
	}

	base := people + mountFood
	if hex != nil && hex.Terrain == Desert {
		base *= 2 // carrying water doubles food requirement
	}
	return base
}

// mountCanForage returns true for terrain where mounts forage for themselves (r215f).
func mountCanForage(t TerrainType) bool {
	return t == Farmland || t == Countryside || t == Forest || t == Hills
}

// AllMounted returns true when every member of the party has a mount.
func (s *GameState) AllMounted() bool {
	if !s.Prince.HasMount {
		return false
	}
	for _, f := range s.Party {
		if !f.HasMount {
			return false
		}
	}
	return true
}

// EffectiveWitWiles returns W&W plus +1 if a true-love companion is present.
func (s *GameState) EffectiveWitWiles() int {
	ww := s.Prince.WitWiles
	for _, f := range s.Party {
		if f.IsTrueLove {
			ww++
			break
		}
	}
	return ww
}

// AvailableActions returns which actions are valid given current hex and state
func (s *GameState) AvailableActions() []Action {
	hex := GetHex(s.CurrentHex)
	if hex == nil {
		return nil
	}

	var actions []Action
	actions = append(actions, ActionTravel)
	actions = append(actions, ActionRest)

	if hex.IsSettlement() {
		actions = append(actions, ActionSeekNews)
		actions = append(actions, ActionSeekFollowers)
		actions = append(actions, ActionBuyFood)
		if hex.Structure == StructCastle || hex.Structure == StructKeep {
			if s.AudienceBarred[s.CurrentHex] <= s.Day {
				actions = append(actions, ActionSeekAudience)
			}
		}
		if hex.Structure == StructTemple {
			actions = append(actions, ActionSubmitOffering)
		}
	}

	if s.GetHexFlags(s.CurrentHex).HiddenTemple {
		actions = append(actions, ActionSubmitOffering)
	}

	if hex.IsRuins() {
		if !s.GetHexFlags(s.CurrentHex).Searched {
			actions = append(actions, ActionSearchRuins)
		}
	}

	hexFlags := s.GetHexFlags(s.CurrentHex)
	if hexFlags.CacheHidden && !hexFlags.CacheFound {
		actions = append(actions, ActionSearchCache)
	}

	// Use Item: available when prince carries any usable consumable
	usable := []PossessionType{PossHealingPotion, PossPoisonAntidote}
	for _, p := range usable {
		if s.Prince.HasPossession(p) {
			actions = append(actions, ActionUseItem)
			break
		}
	}

	// Hunt: available in non-farmland wilderness (no game in towns)
	if hex.Terrain != Farmland && hex.Structure != StructTown {
		actions = append(actions, ActionHunt)
	}

	// Buy Raft: available at settlements when prince doesn't already have one
	if hex.IsSettlement() && !s.Prince.HasPossession(PossRaft) {
		actions = append(actions, ActionBuyRaft)
	}

	return actions
}

// Action represents a player action choice
type Action int

const (
	ActionTravel Action = iota
	ActionRest
	ActionSeekNews
	ActionSeekFollowers
	ActionBuyFood
	ActionSeekAudience
	ActionSubmitOffering
	ActionSearchRuins
	ActionSearchCache
	ActionUseItem
	ActionHunt
	ActionBuyRaft
)

func (a Action) String() string {
	switch a {
	case ActionTravel:
		return "[T]ravel"
	case ActionRest:
		return "[R]est"
	case ActionSeekNews:
		return "[N]ews"
	case ActionSeekFollowers:
		return "[H]ire Followers"
	case ActionBuyFood:
		return "[B]uy Food"
	case ActionSeekAudience:
		return "[A]udience"
	case ActionSubmitOffering:
		return "[O]ffering"
	case ActionSearchRuins:
		return "[S]earch Ruins"
	case ActionSearchCache:
		return "[C]ache"
	case ActionUseItem:
		return "[U]se Item"
	case ActionHunt:
		return "[G]o Hunting"
	case ActionBuyRaft:
		return "[P]urchase Raft"
	}
	return "Unknown"
}

// ActionKey returns the keyboard shortcut character
func (a Action) ActionKey() string {
	switch a {
	case ActionTravel:
		return "t"
	case ActionRest:
		return "r"
	case ActionSeekNews:
		return "n"
	case ActionSeekFollowers:
		return "h"
	case ActionBuyFood:
		return "b"
	case ActionSeekAudience:
		return "a"
	case ActionSubmitOffering:
		return "o"
	case ActionSearchRuins:
		return "s"
	case ActionSearchCache:
		return "c"
	case ActionUseItem:
		return "u"
	case ActionHunt:
		return "g"
	case ActionBuyRaft:
		return "p"
	}
	return ""
}

// CheckWinConditions returns true if the player has won
func CheckWinConditions(s *GameState) (bool, string) {
	// Win: 500+ gold AND north of Tragoth
	if s.Gold >= 500 && IsNorthOfTragoth(s.CurrentHex) {
		return true, "You have amassed enough wealth to reclaim your throne! Cal Arath returns triumphant to the northern lands!"
	}
	// Win: Royal Helm (or Golden Crown substitute) returned to Ogon or Weshor
	if s.Prince.HasPossession(PossRoyalHelm) || s.Prince.HasPossession(PossGoldenCrown) {
		if s.CurrentHex == NewHexID(1, 1) || s.CurrentHex == NewHexID(15, 1) {
			return true, "You have recovered the crown of your dynasty and returned it to the north! Your throne is restored!"
		}
	}
	// Win: Noble Ally secured
	if s.Flags.NobleAllySecured {
		if IsNorthOfTragoth(s.CurrentHex) {
			return true, "With the support of your noble ally, you march north to reclaim your birthright!"
		}
	}
	// Win: Staff of Command + reach north in time
	if s.Prince.HasPossession(PossStaffOfCommand) && IsNorthOfTragoth(s.CurrentHex) {
		return true, "The Staff of Command bends armies to your will. Your throne awaits!"
	}
	return false, ""
}

// CheckLoseConditions returns true if the player has lost
func CheckLoseConditions(s *GameState) (bool, string) {
	if s.Prince.IsDead() {
		return true, "Cal Arath has fallen in battle. His quest ends here."
	}
	if s.Day > 70 {
		return true, "70 days have passed. Without the resources to reclaim your throne, you fade into obscurity."
	}
	return false, ""
}

// AdvanceDay processes end-of-day bookkeeping
func AdvanceDay(s *GameState) {
	// Feed the party
	needed := s.DailyFoodNeeded()
	if s.FoodUnits >= needed {
		s.FoodUnits -= needed
		// Reset starvation for everyone
		s.Prince.StarvationDays = 0
		for i := range s.Party {
			s.Party[i].StarvationDays = 0
		}
	} else {
		// Not enough food — starvation (r216b): CS -1/day, carry halved; no death
		shortage := needed - s.FoodUnits
		s.FoodUnits = 0
		s.AddLog(fmt.Sprintf("You are %d food short! The party goes hungry.", shortage))
		s.Prince.StarvationDays++
		s.AddLog(fmt.Sprintf("Cal Arath starves (day %d): -1 CS, carry capacity halved.", s.Prince.StarvationDays))
		// Follower starvation desertion (r216a): roll 2d6 - W&W - (Morale-3); >= 4 = deserts
		// True-love companions never desert.
		var starvDeserters []int
		for i, f := range s.Party {
			s.Party[i].StarvationDays++
			if f.IsTrueLove {
				continue
			}
			roll := Roll2d6() - s.EffectiveWitWiles() - (f.Morale - 3)
			if roll >= 4 {
				starvDeserters = append(starvDeserters, i)
				s.AddLog(fmt.Sprintf("%s abandons you, unable to endure the hunger!", f.Name))
			}
		}
		for i := len(starvDeserters) - 1; i >= 0; i-- {
			idx := starvDeserters[i]
			s.Party = append(s.Party[:idx], s.Party[idx+1:]...)
		}
	}

	// Pay followers
	totalWages := 0
	var deserters []int
	for i, f := range s.Party {
		totalWages += f.DailyWage
		// Desertion check for unpaid wages — true-love never deserts
		// Roll 2d6 - W&W - (Morale-3) >= 4 means desert
		if !f.IsTrueLove && s.Gold < f.DailyWage {
			roll := Roll2d6() - s.EffectiveWitWiles() - (f.Morale - 3)
			if roll >= 4 {
				deserters = append(deserters, i)
				s.AddLog(fmt.Sprintf("%s deserts due to unpaid wages!", f.Name))
			}
		}
	}
	if totalWages > 0 && s.Gold >= totalWages {
		s.Gold -= totalWages
	}

	// Remove deserters (reverse order to keep indices valid)
	for i := len(deserters) - 1; i >= 0; i-- {
		idx := deserters[i]
		s.Party = append(s.Party[:idx], s.Party[idx+1:]...)
	}

	// Plague dust (r227): deal wounds daily until recovery roll (1d6 >= 4)
	if s.Prince.PlagueDustActive {
		wounds := (Roll1d6() + 1) / 2 // 1d6/2 round up
		s.Prince.Wounds += wounds
		s.AddLog(fmt.Sprintf("The plague dust festers — you suffer %d wound(s).", wounds))
		if Roll1d6() >= 4 {
			s.Prince.PlagueDustActive = false
			s.AddLog("Your body finally purges the plague dust. You begin to recover.")
		}
	}

	// Check for escapees in town
	if GetHex(s.CurrentHex) != nil && GetHex(s.CurrentHex).IsSettlement() {
		var remaining []Character
		for _, f := range s.Party {
			if f.IsEscapee {
				s.AddLog(fmt.Sprintf("%s slips away into the town, free at last.", f.Name))
			} else {
				remaining = append(remaining, f)
			}
		}
		s.Party = remaining
	}

	// Advance time
	s.Day++
	s.DayOfWeek++
	if s.DayOfWeek > 7 {
		s.DayOfWeek = 1
		s.Week++
	}

	s.AddLog(fmt.Sprintf("--- Day %d ---", s.Day))

	// If prince has the Ancient Map, reveal neighbours each day
	revealMapNeighbours(s)

	// Check win/lose
	if won, reason := CheckWinConditions(s); won {
		s.Phase = PhaseGameOver
		s.WinReason = reason
		return
	}
	if lost, reason := CheckLoseConditions(s); lost {
		s.Phase = PhaseGameOver
		s.LoseReason = reason
		return
	}
}

// HasGuide returns true if any follower is a guide
func (s *GameState) HasGuide() bool {
	for _, f := range s.Party {
		if f.IsGuide {
			return true
		}
	}
	return false
}

// TotalCombatSkill returns combined party combat skill (for group combats)
func (s *GameState) TotalCombatSkill() int {
	cs := s.Prince.EffectiveCombatSkill()
	for _, f := range s.Party {
		cs += f.EffectiveCombatSkill()
	}
	return cs
}

// AddFollower adds a follower to the party
func (s *GameState) AddFollower(f Character) {
	s.Party = append(s.Party, f)
}

// RemoveFollower removes a follower by name
func (s *GameState) RemoveFollower(name string) bool {
	for i, f := range s.Party {
		if f.Name == name {
			s.Party = append(s.Party[:i], s.Party[i+1:]...)
			return true
		}
	}
	return false
}

// StatusLine returns a brief status string
func (s *GameState) StatusLine() string {
	line := fmt.Sprintf("Day %d/70 | Gold: %d | Food: %d | Endurance: %d/%d | Followers: %d",
		s.Day, s.Gold, s.FoodUnits,
		s.Prince.CurrentEndurance(), s.Prince.MaxEndurance,
		len(s.Party))
	if s.Prince.PoisonWounds > 0 {
		line += fmt.Sprintf(" | Poison: %d", s.Prince.PoisonWounds)
	}
	return line
}
