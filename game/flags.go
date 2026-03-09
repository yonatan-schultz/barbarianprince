package game

// HexFlags stores per-hex persistent state
type HexFlags struct {
	Searched      bool            // ruins have been searched
	CacheHidden   bool            // a cache is hidden here
	CacheFound    bool            // cache has been found/retrieved
	EventUsed     map[string]bool // one-shot events used
	BaronMet      bool            // baron/noble has been met here
	TempleBarred  bool            // prince barred from this temple
	HiddenTemple  bool            // secret temple discovered here (e066)
}

// GlobalFlags stores game-wide boolean flags
type GlobalFlags struct {
	HasRoyalHelm        bool
	HasStaffOfCommand   bool
	HasNobleParchment   bool // noble ally secured
	NobleAllySecured    bool
	StaffQuestActive    bool
	StaffQuestDays      int  // days remaining in staff quest
	PegasusCaptured     bool
	TrueLoveMet         bool // true love companion has joined the party
	DragonSlain         map[HexID]bool
	WizardsTowerVisited map[HexID]bool
	SecretFound         map[string]bool
}

func NewGlobalFlags() GlobalFlags {
	return GlobalFlags{
		DragonSlain:         make(map[HexID]bool),
		WizardsTowerVisited: make(map[HexID]bool),
		SecretFound:         make(map[string]bool),
	}
}
