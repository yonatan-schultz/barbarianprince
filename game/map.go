package game

// WorldMap holds all hex definitions
var WorldMap map[HexID]*Hex

func init() {
	WorldMap = buildWorldMap()
}

// buildWorldMap constructs the Barbarian Prince hex map (20 cols × 23 rows).
// Col 01 is the western edge, row 01 is the northern edge.
// The Tragoth River runs east-west along the row 2/3 boundary.
// The Nesser River runs north-south between cols 4-5 (rows 9-16),
// then as the Great Nesser between cols 13-14 (rows 17-23).
func buildWorldMap() map[HexID]*Hex {
	m := make(map[HexID]*Hex)

	add := func(col, row int, terrain TerrainType, structure StructureType, name string) {
		id := NewHexID(col, row)
		m[id] = &Hex{
			ID:        id,
			Terrain:   terrain,
			Structure: structure,
			Name:      name,
		}
	}

	// ── Row 1 — Northern strip, north of Tragoth ──────────────────────────
	add(1, 1, Hills, StructTown, "Ogon")
	add(2, 1, Hills, StructNone, "")
	add(3, 1, Forest, StructNone, "")
	add(4, 1, Forest, StructNone, "")
	add(5, 1, Countryside, StructNone, "")
	add(6, 1, Countryside, StructNone, "")
	add(7, 1, Countryside, StructNone, "")
	add(8, 1, Forest, StructNone, "")
	add(9, 1, Hills, StructRuins, "Jakor's Keep")
	add(10, 1, Forest, StructNone, "")
	add(11, 1, Forest, StructNone, "")
	add(12, 1, Forest, StructNone, "")
	add(13, 1, Forest, StructNone, "")
	add(14, 1, Forest, StructNone, "")
	add(15, 1, Hills, StructTown, "Weshor")
	add(16, 1, Countryside, StructNone, "")
	add(17, 1, Countryside, StructNone, "")
	add(18, 1, Hills, StructNone, "")
	add(19, 1, Hills, StructNone, "")
	add(20, 1, Mountains, StructNone, "")

	// ── Row 2 — North of Tragoth (Tragoth River on south boundary) ─────────
	add(1, 2, Hills, StructNone, "")
	add(2, 2, Mountains, StructNone, "")
	add(3, 2, Forest, StructNone, "")
	add(4, 2, Forest, StructNone, "")
	add(5, 2, Forest, StructNone, "")
	add(6, 2, Countryside, StructNone, "")
	add(7, 2, Countryside, StructNone, "")
	add(8, 2, Forest, StructNone, "")
	add(9, 2, Forest, StructNone, "")
	add(10, 2, Forest, StructNone, "")
	add(11, 2, Forest, StructNone, "")
	add(12, 2, Forest, StructNone, "")
	add(13, 2, Forest, StructNone, "")
	add(14, 2, Forest, StructNone, "")
	add(15, 2, Countryside, StructNone, "")
	add(16, 2, Countryside, StructNone, "")
	add(17, 2, Countryside, StructNone, "")
	add(18, 2, Hills, StructNone, "")
	add(19, 2, Hills, StructNone, "")
	add(20, 2, Mountains, StructNone, "")

	// ── Row 3 — South of Tragoth; Barrier Peaks begin ─────────────────────
	add(1, 3, Hills, StructNone, "")
	add(2, 3, Mountains, StructNone, "")
	add(3, 3, Mountains, StructNone, "")
	add(4, 3, Mountains, StructNone, "")
	add(5, 3, Mountains, StructNone, "")
	add(6, 3, Mountains, StructNone, "")
	add(7, 3, Mountains, StructNone, "")
	add(8, 3, Forest, StructNone, "")
	add(9, 3, Forest, StructNone, "")
	add(10, 3, Countryside, StructNone, "")
	add(11, 3, Forest, StructNone, "")
	add(12, 3, Forest, StructNone, "")
	add(13, 3, Forest, StructNone, "")
	add(14, 3, Countryside, StructNone, "")
	add(15, 3, Mountains, StructNone, "")
	add(16, 3, Mountains, StructNone, "")
	add(17, 3, Mountains, StructNone, "")
	add(18, 3, Mountains, StructNone, "")
	add(19, 3, Mountains, StructNone, "")
	add(20, 3, Mountains, StructNone, "")

	// ── Row 4 — Barrier Peaks; Mountains of Zhor; North Pass ──────────────
	add(1, 4, Hills, StructNone, "")
	add(2, 4, Mountains, StructNone, "")
	add(3, 4, Mountains, StructNone, "")
	add(4, 4, Mountains, StructNone, "")
	add(5, 4, Mountains, StructNone, "")
	add(6, 4, Mountains, StructNone, "")
	add(7, 4, Mountains, StructNone, "")
	add(8, 4, Forest, StructNone, "")
	add(9, 4, Countryside, StructNone, "")
	add(10, 4, Countryside, StructNone, "")
	add(11, 4, Countryside, StructNone, "")
	add(12, 4, Countryside, StructNone, "")
	add(13, 4, Forest, StructNone, "")
	add(14, 4, Forest, StructNone, "")
	add(15, 4, Mountains, StructNone, "")
	add(16, 4, Mountains, StructNone, "")
	add(17, 4, Mountains, StructNone, "")
	add(18, 4, Mountains, StructNone, "")
	add(19, 4, Mountains, StructNone, "")
	add(20, 4, Mountains, StructNone, "")

	// ── Row 5 — Dead Plains; Cumry ────────────────────────────────────────
	add(1, 5, Desert, StructNone, "")
	add(2, 5, Hills, StructNone, "")
	add(3, 5, Mountains, StructNone, "")
	add(4, 5, Mountains, StructNone, "")
	add(5, 5, Mountains, StructNone, "")
	add(6, 5, Mountains, StructNone, "")
	add(7, 5, Forest, StructNone, "")
	add(8, 5, Forest, StructNone, "")
	add(9, 5, Countryside, StructNone, "")
	add(10, 5, Hills, StructKeep, "Cumry")
	add(11, 5, Forest, StructNone, "")
	add(12, 5, Forest, StructNone, "")
	add(13, 5, Forest, StructNone, "")
	add(14, 5, Forest, StructNone, "")
	add(15, 5, Mountains, StructNone, "")
	add(16, 5, Mountains, StructNone, "")
	add(17, 5, Mountains, StructNone, "")
	add(18, 5, Mountains, StructNone, "")
	add(19, 5, Mountains, StructNone, "")
	add(20, 5, Mountains, StructNone, "")

	// ── Row 6 — Temple of Zhor ────────────────────────────────────────────
	add(1, 6, Desert, StructNone, "")
	add(2, 6, Hills, StructNone, "")
	add(3, 6, Hills, StructNone, "")
	add(4, 6, Hills, StructNone, "")
	add(5, 6, Forest, StructNone, "")
	add(6, 6, Forest, StructNone, "")
	add(7, 6, Countryside, StructNone, "")
	add(8, 6, Countryside, StructNone, "")
	add(9, 6, Countryside, StructNone, "")
	add(10, 6, Countryside, StructNone, "")
	add(11, 6, Forest, StructNone, "")
	add(12, 6, Forest, StructNone, "")
	add(13, 6, Forest, StructNone, "")
	add(14, 6, Forest, StructNone, "")
	add(15, 6, Mountains, StructNone, "")
	add(16, 6, Mountains, StructNone, "")
	add(17, 6, Mountains, StructTemple, "Temple of Zhor")
	add(18, 6, Mountains, StructNone, "")
	add(19, 6, Mountains, StructNone, "")
	add(20, 6, Mountains, StructNone, "")

	// ── Row 7 ─────────────────────────────────────────────────────────────
	add(1, 7, Desert, StructNone, "")
	add(2, 7, Countryside, StructNone, "")
	add(3, 7, Hills, StructNone, "")
	add(4, 7, Hills, StructNone, "")
	add(5, 7, Hills, StructNone, "")
	add(6, 7, Countryside, StructNone, "")
	add(7, 7, Countryside, StructNone, "")
	add(8, 7, Countryside, StructNone, "")
	add(9, 7, Countryside, StructNone, "")
	add(10, 7, Countryside, StructNone, "")
	add(11, 7, Forest, StructNone, "")
	add(12, 7, Forest, StructNone, "")
	add(13, 7, Forest, StructNone, "")
	add(14, 7, Forest, StructNone, "")
	add(15, 7, Forest, StructNone, "")
	add(16, 7, Forest, StructNone, "")
	add(17, 7, Forest, StructNone, "")
	add(18, 7, Hills, StructNone, "")
	add(19, 7, Hills, StructNone, "")
	add(20, 7, Mountains, StructNone, "")

	// ── Row 8 — Dead Plains ruins ─────────────────────────────────────────
	add(1, 8, Desert, StructRuins, "Dead Plains Ruins")
	add(2, 8, Forest, StructNone, "")
	add(3, 8, Forest, StructNone, "")
	add(4, 8, Hills, StructNone, "")
	add(5, 8, Forest, StructNone, "")
	add(6, 8, Countryside, StructNone, "")
	add(7, 8, Countryside, StructNone, "")
	add(8, 8, Countryside, StructNone, "")
	add(9, 8, Forest, StructNone, "")
	add(10, 8, Forest, StructNone, "")
	add(11, 8, Countryside, StructNone, "")
	add(12, 8, Countryside, StructNone, "")
	add(13, 8, Countryside, StructNone, "")
	add(14, 8, Forest, StructNone, "")
	add(15, 8, Forest, StructNone, "")
	add(16, 8, Forest, StructNone, "")
	add(17, 8, Forest, StructNone, "")
	add(18, 8, Hills, StructNone, "")
	add(19, 8, Hills, StructNone, "")
	add(20, 8, Hills, StructNone, "")

	// ── Row 9 — Cawther; Kabir Desert; Ruins of Pelgar ────────────────────
	add(1, 9, Countryside, StructNone, "")
	add(2, 9, Forest, StructNone, "")
	add(3, 9, Hills, StructNone, "")
	add(4, 9, Hills, StructNone, "")
	add(5, 9, Forest, StructNone, "")
	add(6, 9, Forest, StructNone, "")
	add(7, 9, Forest, StructNone, "")
	add(8, 9, Forest, StructNone, "")
	add(9, 9, Countryside, StructTown, "Cawther")
	add(10, 9, Countryside, StructNone, "")
	add(11, 9, Hills, StructNone, "")
	add(12, 9, Hills, StructNone, "")
	add(13, 9, Hills, StructNone, "")
	add(14, 9, Hills, StructNone, "")
	add(15, 9, Desert, StructNone, "")
	add(16, 9, Desert, StructNone, "")
	add(17, 9, Desert, StructNone, "")
	add(18, 9, Desert, StructNone, "")
	add(19, 9, Hills, StructRuins, "Ruins of Pelgar")
	add(20, 9, Mountains, StructNone, "")

	// ── Row 10 — Llewylla Moor; Angleae; Donry's Temple; Kabir Desert ──────
	add(1, 10, Swamp, StructNone, "")
	add(2, 10, Countryside, StructVillage, "Angleae")
	add(3, 10, Swamp, StructNone, "")
	add(4, 10, Swamp, StructNone, "")
	add(5, 10, Forest, StructNone, "")
	add(6, 10, Forest, StructNone, "")
	add(7, 10, Forest, StructNone, "")
	add(8, 10, Countryside, StructNone, "")
	add(9, 10, Countryside, StructNone, "")
	add(10, 10, Countryside, StructNone, "")
	add(11, 10, Hills, StructTemple, "Donry's Temple")
	add(12, 10, Hills, StructNone, "")
	add(13, 10, Hills, StructNone, "")
	add(14, 10, Desert, StructNone, "")
	add(15, 10, Desert, StructNone, "")
	add(16, 10, Desert, StructNone, "")
	add(17, 10, Desert, StructNone, "")
	add(18, 10, Desert, StructNone, "")
	add(19, 10, Hills, StructNone, "")
	add(20, 10, Mountains, StructNone, "")

	// ── Row 11 — Llewylla Moor; Kabir Desert ──────────────────────────────
	add(1, 11, Swamp, StructNone, "")
	add(2, 11, Swamp, StructNone, "")
	add(3, 11, Swamp, StructNone, "")
	add(4, 11, Swamp, StructNone, "")
	add(5, 11, Countryside, StructNone, "")
	add(6, 11, Forest, StructNone, "")
	add(7, 11, Forest, StructNone, "")
	add(8, 11, Countryside, StructNone, "")
	add(9, 11, Countryside, StructNone, "")
	add(10, 11, Countryside, StructNone, "")
	add(11, 11, Hills, StructNone, "")
	add(12, 11, Hills, StructNone, "")
	add(13, 11, Hills, StructNone, "")
	add(14, 11, Hills, StructNone, "")
	add(15, 11, Desert, StructNone, "")
	add(16, 11, Desert, StructNone, "")
	add(17, 11, Hills, StructNone, "")
	add(18, 11, Hills, StructNone, "")
	add(19, 11, Hills, StructNone, "")
	add(20, 11, Mountains, StructNone, "")

	// ── Row 12 — Llewylla Moor; Branwyn's Castle ──────────────────────────
	add(1, 12, Swamp, StructNone, "")
	add(2, 12, Swamp, StructNone, "")
	add(3, 12, Swamp, StructNone, "")
	add(4, 12, Countryside, StructNone, "")
	add(5, 12, Countryside, StructCastle, "Branwyn's Castle")
	add(6, 12, Forest, StructNone, "")
	add(7, 12, Forest, StructNone, "")
	add(8, 12, Forest, StructNone, "")
	add(9, 12, Hills, StructNone, "")
	add(10, 12, Countryside, StructNone, "")
	add(11, 12, Hills, StructNone, "")
	add(12, 12, Hills, StructNone, "")
	add(13, 12, Hills, StructNone, "")
	add(14, 12, Hills, StructNone, "")
	add(15, 12, Mountains, StructNone, "")
	add(16, 12, Mountains, StructNone, "")
	add(17, 12, Hills, StructNone, "")
	add(18, 12, Hills, StructNone, "")
	add(19, 12, Hills, StructNone, "")
	add(20, 12, Hills, StructNone, "")

	// ── Row 13 — Hulora Castle ────────────────────────────────────────────
	add(1, 13, Swamp, StructNone, "")
	add(2, 13, Swamp, StructNone, "")
	add(3, 13, Countryside, StructNone, "")
	add(4, 13, Forest, StructNone, "")
	add(5, 13, Countryside, StructNone, "")
	add(6, 13, Forest, StructNone, "")
	add(7, 13, Forest, StructNone, "")
	add(8, 13, Forest, StructNone, "")
	add(9, 13, Hills, StructNone, "")
	add(10, 13, Countryside, StructNone, "")
	add(11, 13, Hills, StructNone, "")
	add(12, 13, Hills, StructCastle, "Hulora Castle")
	add(13, 13, Hills, StructNone, "")
	add(14, 13, Hills, StructNone, "")
	add(15, 13, Mountains, StructNone, "")
	add(16, 13, Mountains, StructNone, "")
	add(17, 13, Hills, StructNone, "")
	add(18, 13, Countryside, StructNone, "")
	add(19, 13, Hills, StructNone, "")
	add(20, 13, Hills, StructNone, "")

	// ── Row 14 ─────────────────────────────────────────────────────────────
	add(1, 14, Countryside, StructNone, "")
	add(2, 14, Countryside, StructNone, "")
	add(3, 14, Countryside, StructNone, "")
	add(4, 14, Forest, StructNone, "")
	add(5, 14, Countryside, StructNone, "")
	add(6, 14, Forest, StructNone, "")
	add(7, 14, Forest, StructNone, "")
	add(8, 14, Forest, StructNone, "")
	add(9, 14, Forest, StructNone, "")
	add(10, 14, Forest, StructNone, "")
	add(11, 14, Hills, StructNone, "")
	add(12, 14, Hills, StructNone, "")
	add(13, 14, Forest, StructNone, "")
	add(14, 14, Hills, StructNone, "")
	add(15, 14, Countryside, StructNone, "")
	add(16, 14, Forest, StructNone, "")
	add(17, 14, Countryside, StructNone, "")
	add(18, 14, Countryside, StructNone, "")
	add(19, 14, Countryside, StructNone, "")
	add(20, 14, Countryside, StructNone, "")

	// ── Row 15 ─────────────────────────────────────────────────────────────
	add(1, 15, Countryside, StructNone, "")
	add(2, 15, Countryside, StructNone, "")
	add(3, 15, Countryside, StructNone, "")
	add(4, 15, Forest, StructNone, "")
	add(5, 15, Countryside, StructNone, "")
	add(6, 15, Forest, StructNone, "")
	add(7, 15, Forest, StructNone, "")
	add(8, 15, Forest, StructNone, "")
	add(9, 15, Forest, StructNone, "")
	add(10, 15, Forest, StructNone, "")
	add(11, 15, Hills, StructNone, "")
	add(12, 15, Hills, StructNone, "")
	add(13, 15, Forest, StructNone, "")
	add(14, 15, Countryside, StructNone, "")
	add(15, 15, Countryside, StructNone, "")
	add(16, 15, Countryside, StructNone, "")
	add(17, 15, Countryside, StructNone, "")
	add(18, 15, Countryside, StructNone, "")
	add(19, 15, Countryside, StructNone, "")
	add(20, 15, Countryside, StructNone, "")

	// ── Row 16 — Galden; Lilith ───────────────────────────────────────────
	add(1, 16, Hills, StructTown, "Galden")
	add(2, 16, Countryside, StructNone, "")
	add(3, 16, Countryside, StructNone, "")
	add(4, 16, Farmland, StructNone, "")
	add(5, 16, Farmland, StructNone, "")
	add(6, 16, Countryside, StructNone, "")
	add(7, 16, Countryside, StructNone, "")
	add(8, 16, Countryside, StructNone, "")
	add(9, 16, Countryside, StructNone, "")
	add(10, 16, Countryside, StructNone, "")
	add(11, 16, Countryside, StructNone, "")
	add(12, 16, Countryside, StructVillage, "Lilith")
	add(13, 16, Countryside, StructNone, "")
	add(14, 16, Forest, StructNone, "")
	add(15, 16, Countryside, StructNone, "")
	add(16, 16, Countryside, StructNone, "")
	add(17, 16, Countryside, StructNone, "")
	add(18, 16, Forest, StructNone, "")
	add(19, 16, Countryside, StructNone, "")
	add(20, 16, Hills, StructNone, "")

	// ── Row 17 — Plains of Datha; Erwyna ──────────────────────────────────
	add(1, 17, Hills, StructNone, "")
	add(2, 17, Countryside, StructNone, "")
	add(3, 17, Farmland, StructNone, "")
	add(4, 17, Farmland, StructNone, "")
	add(5, 17, Farmland, StructNone, "")
	add(6, 17, Farmland, StructNone, "")
	add(7, 17, Farmland, StructNone, "")
	add(8, 17, Farmland, StructNone, "")
	add(9, 17, Farmland, StructNone, "")
	add(10, 17, Countryside, StructVillage, "Erwyna")
	add(11, 17, Forest, StructNone, "")
	add(12, 17, Forest, StructNone, "")
	add(13, 17, Forest, StructNone, "")
	add(14, 17, Forest, StructNone, "")
	add(15, 17, Forest, StructNone, "")
	add(16, 17, Forest, StructNone, "")
	add(17, 17, Forest, StructNone, "")
	add(18, 17, Countryside, StructNone, "")
	add(19, 17, Countryside, StructNone, "")
	add(20, 17, Hills, StructNone, "")

	// ── Row 18 — Temple of Duffyd; Disental Branch ────────────────────────
	add(1, 18, Countryside, StructNone, "")
	add(2, 18, Countryside, StructNone, "")
	add(3, 18, Countryside, StructNone, "")
	add(4, 18, Farmland, StructNone, "")
	add(5, 18, Farmland, StructNone, "")
	add(6, 18, Farmland, StructNone, "")
	add(7, 18, Farmland, StructNone, "")
	add(8, 18, Countryside, StructNone, "")
	add(9, 18, Countryside, StructNone, "")
	add(10, 18, Countryside, StructNone, "")
	add(11, 18, Forest, StructNone, "")
	add(12, 18, Forest, StructNone, "")
	add(13, 18, Forest, StructNone, "")
	add(14, 18, Swamp, StructNone, "")
	add(15, 18, Countryside, StructNone, "")
	add(16, 18, Forest, StructNone, "")
	add(17, 18, Forest, StructNone, "")
	add(18, 18, Countryside, StructNone, "")
	add(19, 18, Countryside, StructTemple, "Temple of Duffyd")
	add(20, 18, Mountains, StructNone, "")

	// ── Row 19 — Brigud ───────────────────────────────────────────────────
	add(1, 19, Countryside, StructNone, "")
	add(2, 19, Countryside, StructNone, "")
	add(3, 19, Farmland, StructNone, "")
	add(4, 19, Farmland, StructNone, "")
	add(5, 19, Countryside, StructNone, "")
	add(6, 19, Farmland, StructNone, "")
	add(7, 19, Farmland, StructNone, "")
	add(8, 19, Farmland, StructTown, "Brigud")
	add(9, 19, Countryside, StructNone, "")
	add(10, 19, Forest, StructNone, "")
	add(11, 19, Countryside, StructNone, "")
	add(12, 19, Forest, StructNone, "")
	add(13, 19, Forest, StructNone, "")
	add(14, 19, Swamp, StructNone, "")
	add(15, 19, Countryside, StructNone, "")
	add(16, 19, Forest, StructNone, "")
	add(17, 19, Countryside, StructNone, "")
	add(18, 19, Countryside, StructNone, "")
	add(19, 19, Hills, StructNone, "")
	add(20, 19, Mountains, StructNone, "")

	// ── Row 20 — Halowich; Lullwyn ────────────────────────────────────────
	add(1, 20, Countryside, StructNone, "")
	add(2, 20, Countryside, StructNone, "")
	add(3, 20, Countryside, StructTown, "Halowich")
	add(4, 20, Countryside, StructNone, "")
	add(5, 20, Countryside, StructNone, "")
	add(6, 20, Farmland, StructNone, "")
	add(7, 20, Farmland, StructNone, "")
	add(8, 20, Countryside, StructNone, "")
	add(9, 20, Forest, StructNone, "")
	add(10, 20, Forest, StructNone, "")
	add(11, 20, Countryside, StructNone, "")
	add(12, 20, Countryside, StructNone, "")
	add(13, 20, Countryside, StructNone, "")
	add(14, 20, Swamp, StructNone, "")
	add(15, 20, Countryside, StructNone, "")
	add(16, 20, Countryside, StructTown, "Lullwyn")
	add(17, 20, Countryside, StructNone, "")
	add(18, 20, Forest, StructNone, "")
	add(19, 20, Countryside, StructNone, "")
	add(20, 20, Mountains, StructNone, "")

	// ── Row 21 — Sulwyth Temple ───────────────────────────────────────────
	add(1, 21, Forest, StructNone, "")
	add(2, 21, Countryside, StructNone, "")
	add(3, 21, Countryside, StructNone, "")
	add(4, 21, Countryside, StructNone, "")
	add(5, 21, Forest, StructNone, "")
	add(6, 21, Forest, StructNone, "")
	add(7, 21, Countryside, StructNone, "")
	add(8, 21, Forest, StructNone, "")
	add(9, 21, Countryside, StructTemple, "Sulwyth Temple")
	add(10, 21, Countryside, StructNone, "")
	add(11, 21, Countryside, StructNone, "")
	add(12, 21, Countryside, StructNone, "")
	add(13, 21, Countryside, StructNone, "")
	add(14, 21, Swamp, StructNone, "")
	add(15, 21, Countryside, StructNone, "")
	add(16, 21, Countryside, StructNone, "")
	add(17, 21, Countryside, StructNone, "")
	add(18, 21, Hills, StructNone, "")
	add(19, 21, Hills, StructNone, "")
	add(20, 21, Mountains, StructNone, "")

	// ── Row 22 — Saman Marshes; Lower Drogat ─────────────────────────────
	add(1, 22, Hills, StructNone, "")
	add(2, 22, Hills, StructNone, "")
	add(3, 22, Hills, StructNone, "")
	add(4, 22, Countryside, StructNone, "")
	add(5, 22, Countryside, StructVillage, "Lower Drogat")
	add(6, 22, Forest, StructNone, "")
	add(7, 22, Forest, StructNone, "")
	add(8, 22, Countryside, StructNone, "")
	add(9, 22, Hills, StructNone, "")
	add(10, 22, Hills, StructNone, "")
	add(11, 22, Swamp, StructNone, "")
	add(12, 22, Swamp, StructNone, "")
	add(13, 22, Swamp, StructNone, "")
	add(14, 22, Swamp, StructNone, "")
	add(15, 22, Countryside, StructNone, "")
	add(16, 22, Countryside, StructNone, "")
	add(17, 22, Countryside, StructNone, "")
	add(18, 22, Hills, StructNone, "")
	add(19, 22, Hills, StructNone, "")
	add(20, 22, Mountains, StructNone, "")

	// ── Row 23 — Adrogat Castle; Saman Marshes; Aeravir Castle ───────────
	add(1, 23, Mountains, StructNone, "")
	add(2, 23, Mountains, StructNone, "")
	add(3, 23, Mountains, StructNone, "")
	add(4, 23, Hills, StructNone, "")
	add(5, 23, Hills, StructCastle, "Adrogat Castle")
	add(6, 23, Countryside, StructNone, "")
	add(7, 23, Countryside, StructNone, "")
	add(8, 23, Countryside, StructNone, "")
	add(9, 23, Hills, StructNone, "")
	add(10, 23, Hills, StructNone, "")
	add(11, 23, Swamp, StructNone, "")
	add(12, 23, Swamp, StructNone, "")
	add(13, 23, Swamp, StructNone, "")
	add(14, 23, Swamp, StructNone, "")
	add(15, 23, Countryside, StructNone, "")
	add(16, 23, Countryside, StructNone, "")
	add(17, 23, Hills, StructNone, "")
	add(18, 23, Hills, StructNone, "")
	add(19, 23, Hills, StructNone, "")
	add(20, 23, Hills, StructCastle, "Aeravir Castle")

	initRiversAndRoads(m)
	return m
}

// initRiversAndRoads populates RiverSides on world map hexes.
//
// Tragoth River: east-west boundary between rows 2 and 3 (full map width).
//
// Nesser River: runs north-south between cols 12 and 13, rows 3-16.
//   Col 12 is even; its east neighbours are DirNE=(13,row-1) and DirSE=(13,row).
//   Col 13 is odd;  its west neighbours are DirNW=(12,row)   and DirSW=(12,row+1).
//   Both edge directions are marked so all crossings are blocked.
//
// Largos River: runs north-south between cols 13 and 14, rows 17-23.
//   Col 13 is odd;  its east neighbours are DirNE=(14,row) and DirSE=(14,row+1).
//   Col 14 is even; its west neighbours are DirSW=(13,row) and DirNW=(13,row-1).
func initRiversAndRoads(m map[HexID]*Hex) {
	setRiver := func(col, row int, dir Direction) {
		if h := m[NewHexID(col, row)]; h != nil {
			h.RiverSides[dir] = true
		}
	}

	// setRoad marks both sides of a road edge between two adjacent hexes.
	setRoad := func(col1, row1, col2, row2 int) {
		h1 := m[NewHexID(col1, row1)]
		h2 := m[NewHexID(col2, row2)]
		if h1 == nil || h2 == nil {
			return
		}
		id1 := NewHexID(col1, row1)
		id2 := NewHexID(col2, row2)
		if d := id1.DirectionTo(id2); d >= 0 {
			h1.RoadSides[d] = true
		}
		if d := id2.DirectionTo(id1); d >= 0 {
			h2.RoadSides[d] = true
		}
	}

	// ── Tragoth River: east-west boundary between rows 2 and 3 ───────────────
	for col := 1; col <= 20; col++ {
		setRiver(col, 2, DirS)
		setRiver(col, 3, DirN)
	}

	// ── Nesser River: north-south between cols 12 and 13, rows 3–16 ──────────
	for row := 3; row <= 16; row++ {
		setRiver(12, row, DirSE) // (12,row) → (13,row)
		setRiver(12, row, DirNE) // (12,row) → (13,row-1)
		setRiver(13, row, DirNW) // (13,row) → (12,row)
		setRiver(13, row, DirSW) // (13,row) → (12,row+1)
	}

	// ── Largos River: north-south between cols 13 and 14, rows 17–23 ─────────
	for row := 17; row <= 23; row++ {
		setRiver(13, row, DirNE) // (13,row) → (14,row)
		setRiver(13, row, DirSE) // (13,row) → (14,row+1)
		setRiver(14, row, DirSW) // (14,row) → (13,row)
		setRiver(14, row, DirNW) // (14,row) → (13,row-1)
	}

	// ── Road network ──────────────────────────────────────────────────────────
	//
	// Main Road: Ogon → row-1 east → Tragoth bridge at (8,2)↔(8,3) → Cawther →
	//            south through col 9 → Brigud → west → Halowich → NW → Galden →
	//            south → Lower Drogat → Adrogat Castle.
	//
	// Roads have bridges where they cross rivers (no raft required at those hexes).

	// Row 1: Ogon (1,1) east to (8,1), then east again to Weshor (15,1)
	// — zigzag NE/SE along the northern edge
	setRoad(1, 1, 2, 1)
	setRoad(2, 1, 3, 1)
	setRoad(3, 1, 4, 1)
	setRoad(4, 1, 5, 1)
	setRoad(5, 1, 6, 1)
	setRoad(6, 1, 7, 1)
	setRoad(7, 1, 8, 1)
	// Continue east to Weshor
	setRoad(8, 1, 9, 1)
	setRoad(9, 1, 10, 1)
	setRoad(10, 1, 11, 1)
	setRoad(11, 1, 12, 1)
	setRoad(12, 1, 13, 1)
	setRoad(13, 1, 14, 1)
	setRoad(14, 1, 15, 1) // arrive Weshor

	// Approach and bridge over Tragoth River at col 8
	setRoad(8, 1, 8, 2)
	setRoad(8, 2, 8, 3) // ← bridge: road crosses Tragoth here

	// South through central hills to Cawther (9,9)
	setRoad(8, 3, 8, 4)
	setRoad(8, 4, 8, 5)
	setRoad(8, 5, 8, 6)
	setRoad(8, 6, 8, 7)
	setRoad(8, 7, 8, 8)
	setRoad(8, 8, 9, 8) // jog east
	setRoad(9, 8, 9, 9) // arrive Cawther

	// South from Cawther through col 9 to Brigud (8,19)
	setRoad(9, 9, 9, 10)
	setRoad(9, 10, 9, 11)
	setRoad(9, 11, 9, 12)
	setRoad(9, 12, 9, 13)
	setRoad(9, 13, 9, 14)
	setRoad(9, 14, 9, 15)
	setRoad(9, 15, 9, 16)
	setRoad(9, 16, 9, 17)
	setRoad(9, 17, 9, 18)
	setRoad(9, 18, 9, 19)
	setRoad(9, 19, 8, 19) // arrive Brigud

	// West through farmland: Brigud → Halowich (3,20)
	setRoad(8, 19, 7, 19)
	setRoad(7, 19, 6, 19)
	setRoad(6, 19, 5, 19)
	setRoad(5, 19, 4, 19)
	setRoad(4, 19, 3, 19)
	setRoad(3, 19, 3, 20) // arrive Halowich

	// NW from Halowich to Galden (1,16)
	setRoad(3, 20, 2, 20)
	setRoad(2, 20, 1, 20)
	setRoad(1, 20, 1, 19)
	setRoad(1, 19, 1, 18)
	setRoad(1, 18, 1, 17)
	setRoad(1, 17, 1, 16) // arrive Galden

	// South from Galden (1,16) through col 1, then SE to Lower Drogat (5,22) → Adrogat (5,23)
	// (segments (1,16)→(1,20) already set above by Halowich→Galden route)
	setRoad(1, 20, 1, 21)
	setRoad(1, 21, 2, 21)
	setRoad(2, 21, 3, 21)
	setRoad(3, 21, 4, 21)
	setRoad(4, 21, 5, 21)
	setRoad(5, 21, 5, 22) // arrive Lower Drogat
	setRoad(5, 22, 5, 23) // arrive Adrogat Castle
}

// TragotheRow is the first row south of the Tragoth River.
// Rows 1-2 are the prince's starting territory (north of Tragoth).
const TragotheRow = 3

// IsNorthOfTragoth returns true if the hex is in the northern home territory
func IsNorthOfTragoth(id HexID) bool {
	return id.Row() < TragotheRow
}

// GetHex returns the hex at the given ID, or nil if not on map
func GetHex(id HexID) *Hex {
	return WorldMap[id]
}

// IsOnMap returns true if the hex exists in the world map
func IsOnMap(id HexID) bool {
	_, ok := WorldMap[id]
	return ok
}

// AdjacentHexes returns only valid (on-map) adjacent hexes
func AdjacentHexes(id HexID) []HexID {
	ns := id.Neighbors()
	var result []HexID
	for _, n := range ns {
		if IsOnMap(n) {
			result = append(result, n)
		}
	}
	return result
}

// HasStructure returns true if the hex has the given structure type
func (h *Hex) HasStructure(s StructureType) bool {
	return h.Structure == s
}

// IsSettlement returns true if the hex has a town, village, castle, keep, or temple
func (h *Hex) IsSettlement() bool {
	return h.Structure == StructTown ||
		h.Structure == StructVillage ||
		h.Structure == StructCastle ||
		h.Structure == StructKeep ||
		h.Structure == StructTemple
}

// IsRuins returns true if the hex has ruins
func (h *Hex) IsRuins() bool {
	return h.Structure == StructRuins
}
