package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"barbarianprince/game"
)

const (
	mapCols   = 20
	mapRows   = 23
	cellWidth = 4 // characters per cell (3 content + 1 separator)
	labelW    = 3 // "NN " row-label prefix
)

// RenderMap renders the map as a simple square grid centered on the player.
// width, height are the available terminal character dimensions.
// highlightHex, if non-empty, is drawn with a travel-target style.
func RenderMap(state *game.GameState, width, height int, highlightHex ...game.HexID) string {
	var highlight game.HexID
	if len(highlightHex) > 0 {
		highlight = highlightHex[0]
	}

	// 1 header line + content rows; reserve 1 line for hex info when highlighting
	innerHeight := height - 1
	if innerHeight < 1 {
		innerHeight = 1
	}
	hexRows := innerHeight
	if highlight != "" {
		hexRows-- // reserve bottom line for hex info
	}
	if hexRows < 1 {
		hexRows = 1
	}
	if hexRows > mapRows {
		hexRows = mapRows
	}

	// Available width after label prefix
	hexCols := (width - labelW) / cellWidth
	if hexCols < 1 {
		hexCols = 1
	}
	if hexCols > mapCols {
		hexCols = mapCols
	}

	// Center viewport on player
	playerCol := state.CurrentHex.Col()
	playerRow := state.CurrentHex.Row()
	startCol := playerCol - hexCols/2
	startRow := playerRow - hexRows/2
	if startCol < 1 {
		startCol = 1
	}
	if startRow < 1 {
		startRow = 1
	}
	if startCol+hexCols-1 > mapCols {
		startCol = mapCols - hexCols + 1
	}
	if startRow+hexRows-1 > mapRows {
		startRow = mapRows - hexRows + 1
	}
	if startCol < 1 {
		startCol = 1
	}
	if startRow < 1 {
		startRow = 1
	}

	var rows []string

	// Column header row
	var hdr strings.Builder
	hdr.WriteString(StyleLabel.Render("   "))
	for col := startCol; col < startCol+hexCols && col <= mapCols; col++ {
		hdr.WriteString(StyleLabel.Render(fmt.Sprintf("%3d ", col)))
	}
	rows = append(rows, hdr.String())

	for row := startRow; row < startRow+hexRows && row <= mapRows; row++ {
		var line strings.Builder
		line.WriteString(StyleLabel.Render(fmt.Sprintf("%2d ", row)))
		for col := startCol; col < startCol+hexCols && col <= mapCols; col++ {
			id := game.NewHexID(col, row)
			cell := renderCell(state, id, highlight)
			line.WriteString(cell)
		}
		rows = append(rows, line.String())
	}

	// Hex info line at the bottom when a travel target is highlighted.
	// Truncate to map width so lipgloss's word-wrap (triggered by Width())
	// never splits it into two lines, which would make the panel 1 row taller.
	if highlight != "" {
		info := hexInfoLine(highlight)
		info = ansi.Truncate(info, width, "")
		rows = append(rows, info)
	}

	// Enforce exactly height lines so the panel never overflows or underflows.
	// The Tragoth separator can add one extra line; trim it from the hex rows
	// (never the header at index 0 or the info line at the end).
	for len(rows) > height {
		if highlight != "" && len(rows) >= 2 {
			// Remove the last hex row, keeping the info line at the bottom.
			infoLine := rows[len(rows)-1]
			rows = rows[:len(rows)-2]
			rows = append(rows, infoLine)
		} else {
			rows = rows[:len(rows)-1]
		}
	}
	for len(rows) < height {
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

// hexInfoLine returns a one-line description of a hex for the map info bar.
func hexInfoLine(id game.HexID) string {
	hex := game.GetHex(id)
	if hex == nil {
		return StyleMuted.Render(fmt.Sprintf("  %s — unknown", id))
	}
	parts := []string{fmt.Sprintf("  %s", id)}
	if hex.Name != "" {
		parts = append(parts, hex.Name)
	}
	parts = append(parts, hex.Terrain.String())
	if hex.Structure != game.StructNone {
		parts = append(parts, structureDesc(hex.Structure))
	}
	hasRiver := false
	for _, r := range hex.RiverSides {
		if r {
			hasRiver = true
			break
		}
	}
	if hasRiver {
		parts = append(parts, "river")
	}
	hasRoad := false
	for _, r := range hex.RoadSides {
		if r {
			hasRoad = true
			break
		}
	}
	if hasRoad {
		parts = append(parts, "road")
	}
	return StyleLabel.Render(strings.Join(parts, " · "))
}

// renderCell returns a cellWidth-wide styled string for one hex.
func renderCell(state *game.GameState, id game.HexID, highlight game.HexID) string {
	hex := game.GetHex(id)
	if hex == nil {
		return "    "
	}

	isPlayer := id == state.CurrentHex
	isTarget := highlight != "" && id == highlight
	isVisited := state.VisitedHexes[id]

	flags := state.GetHexFlags(id)

	var content string
	if isPlayer {
		content = "[*]"
	} else if isTarget {
		content = ">>>"
	} else if hex.Name != "" {
		runes := []rune(hex.Name)
		if len(runes) >= 3 {
			content = string(runes[:3])
		} else {
			content = fmt.Sprintf("%-3s", hex.Name)
		}
	} else if hex.Structure != game.StructNone {
		sym := structSymbol(hex.Structure)
		// Show searched ruins as (R) and known cache locations as $xx
		if hex.Structure == game.StructRuins && flags.Searched {
			sym = "(R)"
		} else if flags.CacheHidden && !flags.CacheFound {
			sym = " $ "
		}
		content = sym
	} else {
		// Non-structure hex: show cache marker if present
		if flags.CacheHidden && !flags.CacheFound {
			content = " $ "
		} else {
			content = fmt.Sprintf(" %s ", hex.Terrain.Symbol())
		}
	}

	var styled string
	switch {
	case isPlayer:
		styled = StylePlayer.Render(content)
	case isTarget:
		styled = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8C00")).Bold(true).Blink(true).Render(content)
	case hex.Structure != game.StructNone:
		c := lipgloss.NewStyle().Foreground(colorStructure)
		if !isVisited {
			c = c.Faint(true)
		}
		styled = c.Render(content)
	case isVisited:
		styled = lipgloss.NewStyle().Foreground(TerrainColor(int(hex.Terrain))).Render(content)
	default:
		styled = lipgloss.NewStyle().Foreground(TerrainColor(int(hex.Terrain))).Faint(true).Render(content)
	}

	return styled + " "
}

func structSymbol(s game.StructureType) string {
	switch s {
	case game.StructTown:
		return "[T]"
	case game.StructCastle:
		return "[C]"
	case game.StructTemple:
		return "[+]"
	case game.StructRuins:
		return "[R]"
	case game.StructVillage:
		return "[v]"
	case game.StructKeep:
		return "[K]"
	}
	return " . "
}

// dirName returns the short compass name for a direction index
func dirName(dir int) string {
	switch dir {
	case 0:
		return "N "
	case 1:
		return "NE"
	case 2:
		return "SE"
	case 3:
		return "S "
	case 4:
		return "SW"
	case 5:
		return "NW"
	}
	return "? "
}

// RenderAdjacentHexes renders the list of adjacent hexes for travel selection
func RenderAdjacentHexes(state *game.GameState, selected int) string {
	neighbors := game.AdjacentHexes(state.CurrentHex)
	currentHex := game.GetHex(state.CurrentHex)
	var lines []string
	lines = append(lines, StyleTitle.Render("Choose destination:"))

	for i, n := range neighbors {
		hex := game.GetHex(n)
		if hex == nil {
			continue
		}
		marker := "  "
		if i == selected {
			marker = "> "
		}
		name := hex.Name
		if name == "" {
			name = hex.Terrain.String()
		}
		dir := state.CurrentHex.DirectionTo(n)

		// Road/bridge indicator
		roadTag := ""
		if currentHex != nil && dir >= 0 && currentHex.RoadSides[dir] {
			if currentHex.RiverSides[dir] {
				roadTag = " [bridge]"
			} else {
				roadTag = " [road]"
			}
		}

		line := fmt.Sprintf("%s[%d] %s %s (%s)%s", marker, i+1, dirName(int(dir)), n, name, roadTag)
		if i == selected {
			lines = append(lines, StyleMenuKey.Render(line))
		} else {
			lines = append(lines, StyleMenuText.Render(line))
		}
	}
	lines = append(lines, "")
	lines = append(lines, StyleLabel.Render("[1-6] select  [Enter] confirm  [Esc] cancel"))
	return strings.Join(lines, "\n")
}
