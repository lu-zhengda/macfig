package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengda-lu/macfig/internal/backup"
	"github.com/zhengda-lu/macfig/internal/defaults"
	"github.com/zhengda-lu/macfig/internal/preset"
)

// viewState represents the current TUI screen.
type viewState int

const (
	viewCategories viewState = iota
	viewSettings
	viewConfirmApply
	viewApplyResult
	viewBackupResult
)

// settingRow holds a setting and its current live value.
type settingRow struct {
	setting preset.Setting
	current string
	err     error
}

// readDoneMsg is sent when background reading of current values completes.
type readDoneMsg struct {
	rows []settingRow
}

// applyDoneMsg is sent when a setting has been applied.
type applyDoneMsg struct {
	setting preset.Setting
	newVal  interface{}
	err     error
}

// backupDoneMsg is sent when a backup has been saved.
type backupDoneMsg struct {
	path string
	err  error
}

// Model is the main Bubbletea model for the macfig TUI.
type Model struct {
	exec   defaults.Executor
	runner defaults.CmdRunner

	categories []preset.Category
	catCursor  int

	rows         []settingRow
	setCursor    int
	setOffset    int
	loadingRows  bool

	currentView viewState

	// Apply state
	applyTarget  *settingRow
	applyValue   interface{} // the value to apply (recommended or toggled)
	lastApplyErr error

	// Backup state
	lastBackupPath string
	lastBackupErr  error

	spinner spinner.Model

	width  int
	height int
}

// New creates a new TUI model.
func New(exec defaults.Executor, runner defaults.CmdRunner) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colorCyan)

	return Model{
		exec:       exec,
		runner:     runner,
		categories: preset.AllCategories(),
		spinner:    sp,
	}
}

// Init starts the spinner and loads settings for the first category.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadSettings())
}

// loadSettings reads current values for the selected category.
func (m Model) loadSettings() tea.Cmd {
	cat := m.categories[m.catCursor]
	exec := m.exec
	return func() tea.Msg {
		rows := make([]settingRow, len(cat.Settings))
		for i, s := range cat.Settings {
			val, err := exec.Read(s.Domain, s.Key)
			rows[i] = settingRow{
				setting: s,
				current: val,
				err:     err,
			}
		}
		return readDoneMsg{rows: rows}
	}
}

// doApply applies a specific value for a setting.
func (m Model) doApply(row settingRow, value interface{}) tea.Cmd {
	exec := m.exec
	runner := m.runner
	s := row.setting
	return func() tea.Msg {
		err := defaults.WriteAndRestart(exec, runner, s.Domain, s.Key, value, s.Type, s.RequiresRestart)
		return applyDoneMsg{setting: s, newVal: value, err: err}
	}
}

// doBackup saves a backup of all tracked settings.
func (m Model) doBackup() tea.Cmd {
	exec := m.exec
	return func() tea.Msg {
		path := backup.DefaultPath()
		err := backup.Save(exec, path)
		return backupDoneMsg{path: path, err: err}
	}
}

// Update handles all messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		if m.loadingRows {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case readDoneMsg:
		m.loadingRows = false
		m.rows = msg.rows
		m.setCursor = 0
		m.setOffset = 0
		return m, nil

	case applyDoneMsg:
		m.lastApplyErr = msg.err
		m.currentView = viewApplyResult
		// Reload settings to show updated values.
		return m, m.loadSettings()

	case backupDoneMsg:
		m.lastBackupPath = msg.path
		m.lastBackupErr = msg.err
		m.currentView = viewBackupResult
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.currentView {
		case viewCategories:
			return m.updateCategories(msg)
		case viewSettings:
			return m.updateSettings(msg)
		case viewConfirmApply:
			return m.updateConfirmApply(msg)
		case viewApplyResult:
			return m.updateApplyResult(msg)
		case viewBackupResult:
			return m.updateBackupResult(msg)
		}
	}

	return m, nil
}

func (m Model) updateCategories(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.catCursor > 0 {
			m.catCursor--
			m.loadingRows = true
			return m, tea.Batch(m.loadSettings(), m.spinner.Tick)
		}
	case "down", "j":
		if m.catCursor < len(m.categories)-1 {
			m.catCursor++
			m.loadingRows = true
			return m, tea.Batch(m.loadSettings(), m.spinner.Tick)
		}
	case "enter", "right", "l":
		m.currentView = viewSettings
		m.setCursor = 0
		m.setOffset = 0
	case "b":
		m.currentView = viewBackupResult
		return m, m.doBackup()
	}
	return m, nil
}

func (m Model) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.setCursor > 0 {
			m.setCursor--
			m.ensureSettingVisible()
		}
	case "down", "j":
		if m.setCursor < len(m.rows)-1 {
			m.setCursor++
			m.ensureSettingVisible()
		}
	case " ":
		// Space toggles booleans directly.
		if m.setCursor < len(m.rows) {
			row := m.rows[m.setCursor]
			if row.setting.Type == defaults.TypeBool {
				cur, err := parseBoolValue(row.current)
				if err != nil {
					// Value not set yet — default to setting it ON.
					cur = false
				}
				toggled := !cur
				m.applyTarget = &row
				m.applyValue = toggled
				m.currentView = viewConfirmApply
			}
		}
	case "enter":
		// Enter applies the recommended value.
		if m.setCursor < len(m.rows) {
			row := m.rows[m.setCursor]
			if row.setting.Recommended != nil {
				m.applyTarget = &row
				m.applyValue = row.setting.Recommended
				m.currentView = viewConfirmApply
			}
		}
	case "esc", "left", "h", "backspace":
		m.currentView = viewCategories
	case "b":
		return m, m.doBackup()
	case "r":
		m.loadingRows = true
		return m, tea.Batch(m.loadSettings(), m.spinner.Tick)
	}
	return m, nil
}

func (m Model) updateConfirmApply(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		if m.applyTarget != nil && m.applyValue != nil {
			return m, m.doApply(*m.applyTarget, m.applyValue)
		}
		m.currentView = viewSettings
	case "n", "esc", "backspace":
		m.applyTarget = nil
		m.applyValue = nil
		m.currentView = viewSettings
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateApplyResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc", "backspace":
		m.currentView = viewSettings
		m.applyTarget = nil
		m.applyValue = nil
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateBackupResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc", "backspace":
		m.currentView = viewSettings
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) ensureSettingVisible() {
	visible := m.settingsVisibleCount()
	if m.setCursor < m.setOffset {
		m.setOffset = m.setCursor
	}
	if m.setCursor >= m.setOffset+visible {
		m.setOffset = m.setCursor - visible + 1
	}
}

func (m Model) settingsVisibleCount() int {
	// Reserve: title(2) + category header(2) + help(2) + padding(2) = 8
	available := m.height - 8
	if available < 5 {
		available = 5
	}
	return available
}

// --- Views ---

// View renders the TUI.
func (m Model) View() string {
	switch m.currentView {
	case viewConfirmApply:
		return m.viewConfirmApply()
	case viewApplyResult:
		return m.viewApplyResult()
	case viewBackupResult:
		return m.viewBackupResult()
	default:
		return m.viewTwoPane()
	}
}

// viewTwoPane renders the main two-pane layout.
func (m Model) viewTwoPane() string {
	leftWidth := 20
	rightWidth := m.width - leftWidth - 5 // separators and padding
	if rightWidth < 40 {
		rightWidth = 40
	}

	// Left pane: categories
	left := m.renderCategories(leftWidth)

	// Right pane: settings
	right := m.renderSettings(rightWidth)

	// Join panes side by side.
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		"  ",
		right,
	)

	// Title bar.
	title := titleStyle.Render("macfig")
	help := m.renderHelp()

	return title + "\n" + content + "\n" + help
}

func (m Model) renderCategories(width int) string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("Categories"))
	b.WriteString("\n")

	for i, cat := range m.categories {
		cursor := "  "
		if i == m.catCursor {
			cursor = "> "
		}

		name := cat.Icon + " " + cat.Name
		if lipgloss.Width(name) > width-2 {
			name = truncate(name, width-2)
		}

		line := cursor + pad(name, width-2)
		if i == m.catCursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	// Pad remaining lines so both panes are the same height.
	visible := m.settingsVisibleCount()
	rendered := len(m.categories) + 1 // +1 for header
	for i := rendered; i < visible+1; i++ {
		b.WriteString(strings.Repeat(" ", width))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderSettings(width int) string {
	var b strings.Builder

	cat := m.categories[m.catCursor]
	b.WriteString(headerStyle.Render(cat.Name + " Settings"))
	b.WriteString("\n")

	if m.loadingRows {
		b.WriteString(m.spinner.View() + " Loading...")
		return b.String()
	}

	if len(m.rows) == 0 {
		b.WriteString(dimStyle.Render("No settings in this category."))
		return b.String()
	}

	nameWidth := 30
	if width > 80 {
		nameWidth = 35
	}

	visible := m.settingsVisibleCount()
	total := len(m.rows)
	end := m.setOffset + visible
	if end > total {
		end = total
	}

	// Fixed column width for the current-value field so arrows align.
	valColWidth := 12

	for i := m.setOffset; i < end; i++ {
		row := m.rows[i]
		s := row.setting

		cursor := "  "
		isCursor := false
		if m.currentView == viewSettings && i == m.setCursor {
			cursor = "> "
			isCursor = true
		}

		name := truncate(s.Name, nameWidth)
		currentVal := m.formatValue(row)
		recVal := m.formatRecommended(s)

		// Pad current value to fixed width using visible length
		// so the "->" arrows align regardless of ANSI styling.
		visibleLen := lipgloss.Width(currentVal)
		valPad := ""
		if visibleLen < valColWidth {
			valPad = strings.Repeat(" ", valColWidth-visibleLen)
		}

		var line string
		if recVal != "" {
			line = fmt.Sprintf("%s%-*s  %s%s -> %s",
				cursor, nameWidth, name, currentVal, valPad, recVal)
		} else {
			line = fmt.Sprintf("%s%-*s  %s",
				cursor, nameWidth, name, currentVal)
		}

		if isCursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	// Scroll indicator.
	if total > visible {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  [%d-%d of %d]", m.setOffset+1, end, total)))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) formatValue(row settingRow) string {
	if row.err != nil {
		return dimStyle.Render("(not set)")
	}

	switch row.setting.Type {
	case defaults.TypeBool:
		v, err := parseBoolValue(row.current)
		if err != nil {
			return dimStyle.Render(row.current)
		}
		if v {
			return onStyle.Render("[ON]")
		}
		return offStyle.Render("[OFF]")
	default:
		return row.current
	}
}

func (m Model) formatRecommended(s preset.Setting) string {
	if s.Recommended == nil {
		return ""
	}

	switch s.Type {
	case defaults.TypeBool:
		v, ok := s.Recommended.(bool)
		if !ok {
			return ""
		}
		if v {
			return recommendedStyle.Render("ON")
		}
		return recommendedStyle.Render("OFF")
	default:
		return recommendedStyle.Render(fmt.Sprintf("%v", s.Recommended))
	}
}

func (m Model) renderHelp() string {
	var parts []string
	switch m.currentView {
	case viewCategories:
		parts = []string{"j/k:navigate", "enter:view settings", "b:backup", "q:quit"}
	case viewSettings:
		parts = []string{"j/k:navigate", "space:toggle", "enter:apply recommended", "h/esc:back", "b:backup", "r:refresh", "q:quit"}
	}
	return helpStyle.Render(strings.Join(parts, "  "))
}

func (m Model) viewConfirmApply() string {
	if m.applyTarget == nil {
		return "Nothing to confirm"
	}

	s := m.applyTarget.setting
	var b strings.Builder

	b.WriteString(titleStyle.Render("Confirm Change"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  Setting:     %s\n", s.Name))
	b.WriteString(fmt.Sprintf("  Description: %s\n", s.Description))
	b.WriteString(fmt.Sprintf("  Domain:      %s\n", s.Domain))
	b.WriteString(fmt.Sprintf("  Key:         %s\n", s.Key))
	b.WriteString(fmt.Sprintf("  Current:     %s\n", m.applyTarget.current))
	b.WriteString(fmt.Sprintf("  New value:   %v\n", m.formatNewValue(s, m.applyValue)))

	if s.RequiresRestart != "" {
		b.WriteString(fmt.Sprintf("\n  Note: This will restart %s\n", s.RequiresRestart))
	}

	b.WriteString(helpStyle.Render("\n  y:apply  n/esc:cancel  q:quit"))
	return b.String()
}

// formatNewValue renders the value being applied in a human-readable way.
func (m Model) formatNewValue(s preset.Setting, val interface{}) string {
	if s.Type == defaults.TypeBool {
		if v, ok := val.(bool); ok {
			if v {
				return "ON (true)"
			}
			return "OFF (false)"
		}
	}
	return fmt.Sprintf("%v", val)
}

func (m Model) viewApplyResult() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Result"))
	b.WriteString("\n\n")

	if m.lastApplyErr != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("  Failed: %v", m.lastApplyErr)))
	} else {
		b.WriteString(successStyle.Render("  Setting applied successfully!"))
	}

	b.WriteString(helpStyle.Render("\n\n  enter/esc:back  q:quit"))
	return b.String()
}

func (m Model) viewBackupResult() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Backup"))
	b.WriteString("\n\n")

	if m.lastBackupErr != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("  Backup failed: %v", m.lastBackupErr)))
	} else {
		b.WriteString(successStyle.Render("  Backup saved successfully!"))
		b.WriteString("\n")
		b.WriteString(dimStyle.Render(fmt.Sprintf("  %s", m.lastBackupPath)))
	}

	b.WriteString(helpStyle.Render("\n\n  enter/esc:back  q:quit"))
	return b.String()
}

// parseBoolValue handles macOS defaults boolean representations.
func parseBoolValue(s string) (bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "1", "true", "yes":
		return true, nil
	case "0", "false", "no":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse %q as bool", s)
	}
}
