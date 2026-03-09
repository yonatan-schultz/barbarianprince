package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Base colors
	colorFarmland    = lipgloss.Color("#90EE90") // light green
	colorCountryside = lipgloss.Color("#228B22") // forest green
	colorForest      = lipgloss.Color("#006400") // dark green
	colorHills       = lipgloss.Color("#8B7355") // tan/brown
	colorMountains   = lipgloss.Color("#808080") // gray
	colorSwamp       = lipgloss.Color("#2E8B57") // sea green
	colorDesert      = lipgloss.Color("#DAA520") // goldenrod

	colorRiver     = lipgloss.Color("#4488FF") // blue for rivers
	colorRoad      = lipgloss.Color("#CC8844") // brown for roads
	colorStructure = lipgloss.Color("#FFD700") // gold for structures
	colorPlayer    = lipgloss.Color("#FF4500") // orange-red for player
	colorVisited   = lipgloss.Color("#4169E1") // royal blue for visited
	colorBorder    = lipgloss.Color("#444444")
	colorText      = lipgloss.Color("#FFFFFF")
	colorMuted     = lipgloss.Color("#888888")
	colorGold      = lipgloss.Color("#FFD700")
	colorDanger    = lipgloss.Color("#FF4444")
	colorSuccess   = lipgloss.Color("#44FF44")
	colorWarning   = lipgloss.Color("#FFA500")

	// Styles
	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	StyleTitle = lipgloss.NewStyle().
			Foreground(colorGold).
			Bold(true)

	StyleLabel = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleValue = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)

	StylePlayer = lipgloss.NewStyle().
			Foreground(colorPlayer).
			Bold(true).
			Reverse(true)

	StyleVisited = lipgloss.NewStyle().
			Foreground(colorVisited)

	StyleStructure = lipgloss.NewStyle().
			Foreground(colorStructure)

	StyleDanger = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess)

	StyleWarning = lipgloss.NewStyle().
			Foreground(colorWarning)

	StyleMenuKey = lipgloss.NewStyle().
			Foreground(colorGold).
			Bold(true)

	StyleMenuText = lipgloss.NewStyle().
			Foreground(colorText)

	StyleLog = lipgloss.NewStyle().
			Foreground(colorText)

	StyleLogOld = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleGameOver = lipgloss.NewStyle().
			Foreground(colorGold).
			Bold(true).
			Align(lipgloss.Center)

	StyleTutorial = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00CFFF")).
			Bold(true)
)

// TerrainColor returns the foreground color for a terrain type
func TerrainColor(t int) lipgloss.Color {
	switch t {
	case 0: // Farmland
		return colorFarmland
	case 1: // Countryside
		return colorCountryside
	case 2: // Forest
		return colorForest
	case 3: // Hills
		return colorHills
	case 4: // Mountains
		return colorMountains
	case 5: // Swamp
		return colorSwamp
	case 6: // Desert
		return colorDesert
	}
	return colorText
}