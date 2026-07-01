package applescript

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTracks(t *testing.T) {
	t.Parallel()

	raw := []byte(`[
		{"name": "Gorgon", "artist": "Utsu-P", "album": "X", "duration": 255.0},
		{"name": "Vulgar", "artist": "Utsu-P", "album": "X", "duration": 215.0}
	]`)

	got, err := parseTracks(raw)

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Gorgon", got[0].Name)
	assert.Equal(t, 255*time.Second, got[0].Duration)
	assert.Equal(t, "Vulgar", got[1].Name)
}

func TestParseTracksEmpty(t *testing.T) {
	t.Parallel()

	got, err := parseTracks([]byte(`[]`))

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestSearchScriptEscapesQuery(t *testing.T) {
	t.Parallel()

	// A query with quotes must be embedded safely as a JSON literal.
	script := searchScript(`the "quoted" one`, 25)

	assert.Contains(t, script, `{for: "the \"quoted\" one", only: 'all'}`)
	assert.Contains(t, script, "const limit = 25;")
}

func TestPlayerSearchRunsJavaScriptAndParses(t *testing.T) {
	t.Parallel()

	fake := &fakeRunner{out: []byte(`[{"name":"Gorgon","artist":"Utsu-P"}]`)}
	p := newPlayer(fake)

	got, err := p.Search(context.Background(), "utsu", 10)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Gorgon", got[0].Name)
	require.Len(t, fake.calls, 1)
	assert.Equal(t, javaScript, fake.calls[0].lang)
	assert.True(t, strings.Contains(fake.calls[0].script, `"utsu"`))
}
