package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

// send feeds keys to the model and returns the resulting model.
func send(m model, keys ...string) model {
	for _, k := range keys {
		next, _ := m.Update(key(k))
		m = next.(model)
	}
	return m
}

func TestPickerMovesAndChooses(t *testing.T) {
	t.Parallel()

	m := newModel("Search", []string{"a", "b", "c"})

	m = send(m, "down", "down", "enter")

	assert.Equal(t, 2, m.chosen)
}

func TestPickerClampsAtEnds(t *testing.T) {
	t.Parallel()

	m := newModel("Search", []string{"a", "b"})

	m = send(m, "up") // already at top
	assert.Equal(t, 0, m.cursor)

	m = send(m, "down", "down", "down") // past the end
	assert.Equal(t, 1, m.cursor)
}

func TestPickerVimKeysAndJumps(t *testing.T) {
	t.Parallel()

	m := newModel("Search", []string{"a", "b", "c", "d"})

	m = send(m, "j", "j", "G")
	assert.Equal(t, 3, m.cursor)

	m = send(m, "g")
	assert.Equal(t, 0, m.cursor)
}

func TestPickerCancelDoesNotChoose(t *testing.T) {
	t.Parallel()

	m := send(newModel("Search", []string{"a", "b"}), "down", "q")
	assert.Equal(t, -1, m.chosen)

	m = send(newModel("Search", []string{"a", "b"}), "esc")
	assert.Equal(t, -1, m.chosen)
}

func TestPickerScrollsToKeepCursorVisible(t *testing.T) {
	t.Parallel()

	m := newModel("Search", []string{"0", "1", "2", "3", "4", "5"})
	next, _ := m.Update(tea.WindowSizeMsg{Height: 5}) // height -> 3 rows
	m = next.(model)

	m = send(m, "down", "down", "down", "down") // cursor 4, must scroll
	assert.GreaterOrEqual(t, m.cursor, m.offset)
	assert.Less(t, m.cursor, m.offset+m.height)
}

func TestPickEmptyReturnsNoChoice(t *testing.T) {
	t.Parallel()

	idx, ok, err := Pick("Search", nil)

	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, -1, idx)
}

func TestPickerViewShowsCursorAndItems(t *testing.T) {
	t.Parallel()

	view := newModel("Search: daft", []string{"Robot Rock", "Around the World"}).View()

	assert.Contains(t, view, "Search: daft")
	assert.Contains(t, view, "▶ ")
	assert.Contains(t, view, "Robot Rock")
	assert.Contains(t, view, "Around the World")
}
