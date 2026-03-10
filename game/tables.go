package game

// EventID identifies a specific event
type EventID string

// TravelTableEntry defines lost/event thresholds and event reference tables for a terrain type
type TravelTableEntry struct {
	LostThreshold  int       // 2d6 >= this = lost (0 = can't get lost)
	EventThreshold int       // 2d6 >= this = event occurs
	EventRefs      [6]string // r231-r281 reference codes, one per 1d6 roll
	RoadEventRef   string    // event table for road travel
}

// r207TravelTable is the travel event table indexed by terrain
// Values from Barbarian Prince rule tables
var r207TravelTable = map[TerrainType]TravelTableEntry{
	Farmland: {
		LostThreshold:  0, // can't get lost in farmland
		EventThreshold: 7,
		EventRefs:      [6]string{"r231", "r232", "r233", "r234", "r235", "r236"},
		RoadEventRef:   "r230",
	},
	Countryside: {
		LostThreshold:  0,
		EventThreshold: 8,
		EventRefs:      [6]string{"r241", "r242", "r243", "r244", "r245", "r246"},
		RoadEventRef:   "r230",
	},
	Forest: {
		LostThreshold:  9,
		EventThreshold: 7,
		EventRefs:      [6]string{"r251", "r252", "r253", "r254", "r255", "r256"},
		RoadEventRef:   "r230",
	},
	Hills: {
		LostThreshold:  10,
		EventThreshold: 7,
		EventRefs:      [6]string{"r261", "r262", "r263", "r264", "r265", "r266"},
		RoadEventRef:   "r230",
	},
	Mountains: {
		LostThreshold:  8,
		EventThreshold: 7,
		EventRefs:      [6]string{"r271", "r272", "r273", "r274", "r275", "r276"},
		RoadEventRef:   "r230",
	},
	Swamp: {
		LostThreshold:  8,
		EventThreshold: 7,
		EventRefs:      [6]string{"sw1", "sw2", "sw3", "sw4", "sw5", "sw6"},
		RoadEventRef:   "r230",
	},
	Desert: {
		LostThreshold:  8,
		EventThreshold: 7,
		EventRefs:      [6]string{"r281", "r282", "r281", "r282", "r281", "r282"},
		RoadEventRef:   "r230",
	},
}

// EventRefTable maps a reference code to a 1d6 event lookup
type EventRefTable [6]EventID

// r231-r281 event reference tables
// Each entry maps roll 1-6 to an event ID
var eventRefTables = map[string]EventRefTable{
	// Farmland event tables
	"r231": {"e009", "e009", "e010", "e011", "e012", "e013"}, // Farm/settlement
	"r232": {"e014", "e015", "e016", "e017", "e017", "e018"}, // Peasants/priest
	"r233": {"e019", "e020", "e021", "e022", "e023", "e024"}, // Travellers/bandits
	"r234": {"e025", "e026", "e027", "e028", "e029", "e030"}, // Treasure/merchants
	"r235": {"e031", "e032", "e033", "e034", "e166", "e004"}, // tournament, misc
	"r236": {"e178", "e006", "e007", "e008", "e169", "e050"}, // pretender bounty, priest

	// Countryside event tables
	"r241": {"e009", "e009", "e051", "e052", "e053", "e054"},
	"r242": {"e055", "e056", "e057", "e058", "e059", "e060"},
	"r243": {"e061", "e062", "e063", "e064", "e065", "e066"},
	"r244": {"e003", "e004", "e005", "e006", "e007", "e008"},
	"r245": {"e025", "e027", "e028", "e067", "e068", "e069"},
	"r246": {"e167", "e168", "e174", "e179", "e050", "e128"}, // wandering noble, wagon, bandit lord, sage

	// Forest event tables
	"r251": {"e051", "e057", "e070", "e071", "e072", "e073"},
	"r252": {"e074", "e075", "e076", "e077", "e078", "e079"},
	"r253": {"e080", "e081", "e006", "e007", "e008", "e082"},
	"r254": {"e064", "e065", "e066", "e068", "e083", "e084"},
	"r255": {"e172", "e176", "e028", "e003", "e004", "e005"}, // forest witch, elven ambassador
	"r256": {"e085", "e086", "e087", "e088", "e089", "e128"},

	// Hills event tables
	"r261": {"e051", "e057", "e070", "e075", "e085", "e086"},
	"r262": {"e087", "e088", "e089", "e090", "e091", "e092"},
	"r263": {"e093", "e094", "e095", "e096", "e097", "e098"},
	"r264": {"e003", "e006", "e008", "e064", "e066", "e068"},
	"r265": {"e025", "e027", "e028", "e175", "e129", "e128"}, // dwarven forge
	"r266": {"e078", "e079", "e022", "e023", "e051", "e057"},

	// Mountains event tables
	"r271": {"e085", "e086", "e087", "e088", "e089", "e090"},
	"r272": {"e091", "e092", "e093", "e095", "e096", "e097"},
	"r273": {"e098", "e099", "e100", "e070", "e071", "e051"},
	"r274": {"e006", "e008", "e064", "e066", "e068", "e025"},
	"r275": {"e027", "e028", "e040", "e171", "e128", "e128"}, // mountain hermit king
	"r276": {"e078", "e079", "e120", "e121", "e085", "e086"},

	// Swamp event tables (thematic: creatures, hazards, wilderness, NPCs)
	"sw1": {"e094", "e056", "e074", "e091", "e090", "e092"}, // Swamp creatures
	"sw2": {"e089", "e079", "e078", "e091", "e122", "e125"}, // Hazards & terrain
	"sw3": {"e063", "e066", "e083", "e080", "e082", "e128"}, // NPCs & magic
	"sw4": {"e025", "e058", "e027", "e028", "e124", "e127"}, // Treasure & desert crossover
	"sw5": {"e051", "e057", "e003", "e006", "e008", "e128"}, // Encounters
	"sw6": {"e126", "e092", "e054", "e090", "e129", "e128"}, // Misc

	// Desert event tables
	"r281": {"e120", "e122", "e121", "e123", "e092", "e091"}, // Heat, nomads, oasis, flood, snakes
	"r282": {"e124", "e125", "e126", "e127", "e170", "e128"}, // Ruins, storms, merchants, oracle

	// Road event table
	"r230": {"e009", "e022", "e051", "e003", "e128", "e050"},
}

// LookupEventRef looks up an event from a reference table
func LookupEventRef(tableCode string, roll1d6 int) EventID {
	table, ok := eventRefTables[tableCode]
	if !ok {
		return "e128" // nothing happens
	}
	idx := roll1d6 - 1
	if idx < 0 {
		idx = 0
	}
	if idx > 5 {
		idx = 5
	}
	return table[idx]
}

// r220c combat result table: (attacker CS - defender CS + 2d6) => wounds inflicted
// Only non-zero entries listed; anything else = 0 (miss)
var combatWoundsTable = map[int]int{
	2:  0,
	3:  0,
	4:  0,
	5:  1,
	6:  1,
	7:  1,
	8:  1,
	9:  1,
	10: 2,
	11: 1,
	12: 2,
	13: 2,
	14: 3,
	15: 2,
	16: 5,
	17: 2,
	18: 5,
	19: 5,
	20: 6,
}

// CombatWounds returns the number of wounds inflicted based on net combat roll
func CombatWounds(netRoll int) int {
	if netRoll <= 0 {
		return 0
	}
	if netRoll > 20 {
		netRoll = 20
	}
	w, ok := combatWoundsTable[netRoll]
	if !ok {
		return 0
	}
	return w
}

// r226 treasure table: wealthCode -> 1d6 -> gold amount
var treasureGoldTable = map[int][6]int{
	1: {5, 10, 15, 20, 25, 30},
	2: {10, 20, 30, 40, 50, 60},
	3: {20, 40, 60, 80, 100, 150},
	4: {30, 60, 90, 120, 150, 200},
	5: {50, 100, 150, 200, 250, 300},
	6: {100, 150, 200, 250, 300, 400},
	7: {150, 200, 300, 400, 500, 600},
}

// TreasureRoll returns gold for a given wealth code and 1d6 roll
func TreasureRoll(wealthCode, roll1d6 int) int {
	table, ok := treasureGoldTable[wealthCode]
	if !ok {
		return 0
	}
	idx := roll1d6 - 1
	if idx < 0 {
		idx = 0
	}
	if idx > 5 {
		idx = 5
	}
	return table[idx]
}

// SpecialPossessionRoll determines if a treasure includes a special possession
// Returns PossNone if no special item, otherwise returns the item
func SpecialPossessionRoll(wealthCode int) PossessionType {
	roll := Roll1d6()
	// Higher wealth codes have better chances
	threshold := 7 - wealthCode
	if roll >= threshold {
		// On a second roll of 6, yield a rare magic item instead of a common one
		if Roll1d6() == 6 {
			rareItems := []PossessionType{
				PossAlcoveOfSending,
				PossArchOfTravel,
				PossGatewayToDarkness,
				PossMirrorOfReversal,
			}
			return rareItems[Roll1d6()%len(rareItems)]
		}
		// Roll for which common item
		itemRoll := Roll1d6()
		items := []PossessionType{
			PossRingOfCommand,
			PossAmuletOfPower,
			PossElvenBoots,
			PossHealingPotion,
			PossMagicSword,
			PossHolySymbol,
		}
		return items[itemRoll-1]
	}
	return PossNone
}

// r208 ruins search table — 2d6 (2-12) for 11 outcomes
var ruinsSearchTable = [11]EventID{
	"e136", // 2  — empty
	"e134", // 3  — ancient trap
	"e131", // 4  — gold cache
	"e139", // 5  — trapped corridor (new)
	"e132", // 6  — undead guardians
	"e131", // 7  — gold cache (most common)
	"e140", // 8  — ancient library (new)
	"e137", // 9  — rival treasure hunters
	"e133", // 10 — rich chamber + item
	"e141", // 11 — undead swarm (new)
	"e138", // 12 — Royal Helm vault!
}

// LookupRuinsEvent returns the ruins search event for a 2d6 roll (2-12)
func LookupRuinsEvent(roll2d6 int) EventID {
	if roll2d6 < 2 {
		roll2d6 = 2
	}
	if roll2d6 > 12 {
		roll2d6 = 12
	}
	return ruinsSearchTable[roll2d6-2]
}

// r209 seek news table — expanded with new NPC rumours
var seekNewsTable = [6]EventID{
	"e144", "e165", "e064", "e065", "e066", "e143",
}

func LookupNewsEvent(roll1d6 int) EventID {
	idx := roll1d6 - 1
	if idx < 0 {
		idx = 0
	}
	if idx > 5 {
		idx = 5
	}
	return seekNewsTable[idx]
}

// r210 seek followers table
var seekFollowersTable = [6]EventID{
	"e003", "e004", "e005", "e006", "e007", "e008",
}

func LookupFollowerEvent(roll1d6 int) EventID {
	idx := roll1d6 - 1
	if idx < 0 {
		idx = 0
	}
	if idx > 5 {
		idx = 5
	}
	return seekFollowersTable[idx]
}

// r211 seek audience table (for castle) — e162/e163 add new NPC outcomes
var seekAudienceTable = [6]EventID{
	"e148", "e162", "e150", "e151", "e163", "e152",
}

func LookupAudienceEvent(roll1d6 int) EventID {
	idx := roll1d6 - 1
	if idx < 0 {
		idx = 0
	}
	if idx > 5 {
		idx = 5
	}
	return seekAudienceTable[idx]
}

// r212 offering table (for temple)
var seekOfferingTable = [6]EventID{
	"e154", "e155", "e156", "e157", "e158", "e159",
}

func LookupOfferingEvent(roll1d6 int) EventID {
	idx := roll1d6 - 1
	if idx < 0 {
		idx = 0
	}
	if idx > 5 {
		idx = 5
	}
	return seekOfferingTable[idx]
}
