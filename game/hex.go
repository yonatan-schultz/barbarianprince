package game

import "fmt"

// TerrainType represents hex terrain
type TerrainType int

const (
	Farmland    TerrainType = iota
	Countryside             // gentle hills, open fields
	Forest
	Hills
	Mountains
	Swamp
	Desert
)

func (t TerrainType) String() string {
	switch t {
	case Farmland:
		return "Farmland"
	case Countryside:
		return "Countryside"
	case Forest:
		return "Forest"
	case Hills:
		return "Hills"
	case Mountains:
		return "Mountains"
	case Swamp:
		return "Swamp"
	case Desert:
		return "Desert"
	}
	return "Unknown"
}

// TerrainSymbol returns the ASCII symbol for a terrain type
func (t TerrainType) Symbol() string {
	switch t {
	case Farmland:
		return "."
	case Countryside:
		return "~"
	case Forest:
		return "f"
	case Hills:
		return "^"
	case Mountains:
		return "M"
	case Swamp:
		return "s"
	case Desert:
		return "o"
	}
	return "?"
}

// StructureType represents what structure (if any) is on the hex
type StructureType int

const (
	StructNone    StructureType = iota
	StructTown                  // e.g. Ogon, Weshor
	StructCastle                // noble castle
	StructTemple                // temple/shrine
	StructRuins                 // ancient ruins
	StructVillage               // small village
	StructKeep                  // minor fortress
)

// Direction constants for hex sides (flat-top hex)
type Direction int

const (
	DirN  Direction = iota // 0 North
	DirNE                  // 1 Northeast
	DirSE                  // 2 Southeast
	DirS                   // 3 South
	DirSW                  // 4 Southwest
	DirNW                  // 5 Northwest
)

// HexID is the 4-character identifier "CCRR" (column/row, 1-based, zero-padded)
type HexID string

// NewHexID creates a HexID from column and row (1-based)
func NewHexID(col, row int) HexID {
	return HexID(fmt.Sprintf("%02d%02d", col, row))
}

// Col returns the column number (1-based)
func (h HexID) Col() int {
	if len(h) != 4 {
		return 0
	}
	var c int
	fmt.Sscanf(string(h[:2]), "%d", &c)
	return c
}

// Row returns the row number (1-based)
func (h HexID) Row() int {
	if len(h) != 4 {
		return 0
	}
	var r int
	fmt.Sscanf(string(h[2:]), "%d", &r)
	return r
}

// Hex represents a single hex on the map
type Hex struct {
	ID         HexID
	Terrain    TerrainType
	Structure  StructureType
	Name       string        // named location (empty if unnamed)
	RiverSides [6]bool       // river on hex side indexed by Direction
	RoadSides  [6]bool       // road on hex side indexed by Direction
	OffLimits  bool          // e.g. sea hexes
}

// Neighbors returns the adjacent HexIDs (some may not exist in the map)
// Uses offset-grid rules: even columns shift differently than odd
func (h HexID) Neighbors() [6]HexID {
	col := h.Col()
	row := h.Row()
	var ns [6]HexID
	if col%2 == 0 {
		// even column
		ns[DirN]  = NewHexID(col, row-1)
		ns[DirNE] = NewHexID(col+1, row-1)
		ns[DirSE] = NewHexID(col+1, row)
		ns[DirS]  = NewHexID(col, row+1)
		ns[DirSW] = NewHexID(col-1, row)
		ns[DirNW] = NewHexID(col-1, row-1)
	} else {
		// odd column
		ns[DirN]  = NewHexID(col, row-1)
		ns[DirNE] = NewHexID(col+1, row)
		ns[DirSE] = NewHexID(col+1, row+1)
		ns[DirS]  = NewHexID(col, row+1)
		ns[DirSW] = NewHexID(col-1, row+1)
		ns[DirNW] = NewHexID(col-1, row)
	}
	return ns
}

// DirectionTo returns the direction from h to neighbor n, or -1 if not adjacent
func (h HexID) DirectionTo(n HexID) Direction {
	ns := h.Neighbors()
	for d, nb := range ns {
		if nb == n {
			return Direction(d)
		}
	}
	return -1
}