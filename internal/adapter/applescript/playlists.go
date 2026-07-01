package applescript

import (
	"encoding/json"
	"fmt"

	"github.com/helmedeiros/amp/internal/music"
)

// playlistsScript lists the user's playlists as JSON, skipping the built-in
// library playlists (which report a non-"none" special kind).
const playlistsScript = `
const Music = Application('Music');
const out = [];
if (Music.running()) {
  for (const p of Music.userPlaylists()) {
    if (p.specialKind() !== 'none') continue;
    out.push({name: p.name(), count: p.tracks().length});
  }
}
JSON.stringify(out);
`

type playlistDTO struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// parsePlaylists decodes a JSON array of playlists into the domain type.
func parsePlaylists(raw []byte) ([]music.Playlist, error) {
	var dtos []playlistDTO
	if err := json.Unmarshal(raw, &dtos); err != nil {
		return nil, fmt.Errorf("decode playlists: %w", err)
	}

	playlists := make([]music.Playlist, len(dtos))
	for i, d := range dtos {
		playlists[i] = music.Playlist{Name: d.Name, Count: d.Count}
	}
	return playlists, nil
}
