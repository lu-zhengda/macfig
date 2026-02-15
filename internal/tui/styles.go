package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette.
var (
	colorCyan   = lipgloss.Color("6")
	colorGreen  = lipgloss.Color("2")
	colorYellow = lipgloss.Color("3")
	colorRed    = lipgloss.Color("1")
	colorGray   = lipgloss.Color("8")
	colorWhite  = lipgloss.Color("15")
	colorPink   = lipgloss.Color("212")
)

// Layout styles.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan).
			PaddingBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPink)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			PaddingTop(1)

	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Underline(true)

	onStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorGreen)

	offStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	recommendedStyle = lipgloss.NewStyle().
				Foreground(colorYellow)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRed)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGreen)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGray)
)

// renderBar creates a simple horizontal bar chart element.
func renderBar(value, maxValue int64, width int) string {
	if maxValue == 0 {
		return ""
	}
	ratio := float64(value) / float64(maxValue)
	filled := int(ratio * float64(width))
	if filled == 0 && value > 0 {
		filled = 1
	}
	return "|" + strings.Repeat("#", filled) + strings.Repeat(".", width-filled) + "|"
}

// pad right-pads a string to display width w with spaces.
// Uses lipgloss.Width for correct handling of multi-byte/emoji characters.
func pad(s string, w int) string {
	visible := lipgloss.Width(s)
	if visible >= w {
		return s
	}
	return s + strings.Repeat(" ", w-visible)
}

// truncate truncates a string to max display length, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if maxLen < 4 {
		maxLen = 4
	}
	if lipgloss.Width(s) <= maxLen {
		return s
	}
	// Trim runes until we fit.
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		candidate := string(runes[:i]) + "..."
		if lipgloss.Width(candidate) <= maxLen {
			return candidate
		}
	}
	return "..."
}
