package applescript

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlaySearchScript(t *testing.T) {
	t.Parallel()

	script := playSearchScript(`the "one"`, 25, 3)

	assert.Contains(t, script, `Music.search(lib, {for: "the \"one\"", only: 'all'})`)
	assert.Contains(t, script, "const limit = 25;") // limit
	assert.Contains(t, script, "res.slice(0, limit)")
	assert.Contains(t, script, "(((3) % res.length)") // rotation by start
	assert.Contains(t, script, `Music.userPlaylists.byName("amp queue")`)
	assert.Contains(t, script, "pl.play();")
}

func TestPlayerPlaySearchRunsJavaScript(t *testing.T) {
	t.Parallel()

	fake := &fakeRunner{out: []byte(`{"queued":5}`)}
	p := newPlayer(fake)

	err := p.PlaySearch(context.Background(), "daft", 50, 2)

	require.NoError(t, err)
	require.Len(t, fake.calls, 1)
	assert.Equal(t, javaScript, fake.calls[0].lang)
	assert.Contains(t, fake.calls[0].script, `"daft"`)
}
