// Package tui holds interactive terminal components built with Bubble Tea.
// The picker here is a generic single-choice list selector; it will also seed
// the Phase 3 full TUI.
package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true)
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	hintStyle     = lipgloss.NewStyle().Faint(true)
)

// model is the picker's Bubble Tea state.
type model struct {
	title  string
	items  []string
	cursor int
	chosen int // index chosen with enter, or -1
	height int // number of visible rows
	offset int // index of the first visible row
}

func newModel(title string, items []string) model {
	return model{title: title, items: items, chosen: -1, height: 15}
}

func (m model) Init() tea.Cmd { return nil }

// Update handles key and window-size messages. It is pure and unit-testable.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if h := msg.Height - 2; h > 0 { // leave room for title and hint
			m.height = h
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "home", "g":
			m.cursor = 0
		case "end", "G":
			m.cursor = len(m.items) - 1
		case "enter":
			m.chosen = m.cursor
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
		m.clampOffset()
	}
	return m, nil
}

// clampOffset scrolls the window so the cursor stays visible.
func (m *model) clampOffset() {
	switch {
	case m.cursor < m.offset:
		m.offset = m.cursor
	case m.cursor >= m.offset+m.height:
		m.offset = m.cursor - m.height + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m model) View() string {
	var b strings.Builder
	if m.title != "" {
		b.WriteString(titleStyle.Render(m.title) + "\n")
	}

	end := m.offset + m.height
	if end > len(m.items) {
		end = len(m.items)
	}
	for i := m.offset; i < end; i++ {
		if i == m.cursor {
			b.WriteString("▶ " + selectedStyle.Render(m.items[i]) + "\n")
			continue
		}
		b.WriteString("  " + m.items[i] + "\n")
	}

	b.WriteString(hintStyle.Render("↑/↓ move · enter play · q cancel"))
	return b.String()
}

// Pick runs an interactive single-choice list and returns the chosen index and
// whether a choice was confirmed (false if cancelled or the list is empty).
func Pick(title string, items []string) (int, bool, error) {
	if len(items) == 0 {
		return -1, false, nil
	}

	res, err := tea.NewProgram(newModel(title, items)).Run()
	if err != nil {
		return -1, false, err
	}

	final := res.(model)
	return final.chosen, final.chosen >= 0, nil
}
