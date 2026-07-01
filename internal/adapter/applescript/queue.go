package applescript

import (
	"encoding/json"
	"fmt"
)

// queuePlaylistName is the amp-owned playlist used as the play queue.
const queuePlaylistName = "amp queue"

// playSearchScript builds a JXA program that re-runs the library search, loads
// the results into the managed queue playlist rotated so the chosen track is
// first, and plays the playlist from the top. Everything after the pick plays
// next; earlier results sit behind it. Music.app's live "Up Next" is not
// scriptable, so rotating and playing from the top is how we honour the pick
// (see ADR-0004).
func playSearchScript(query string, limit, start int) string {
	q, _ := json.Marshal(query)
	name, _ := json.Marshal(queuePlaylistName)

	return fmt.Sprintf(`
const Music = Application('Music');
const lib = Music.libraryPlaylists[0];
let res = Music.search(lib, {for: %s, only: 'all'});
const limit = %d;
if (limit > 0) res = res.slice(0, limit);
if (res.length > 0) {
  const s = (((%d) %% res.length) + res.length) %% res.length;
  res = res.slice(s).concat(res.slice(0, s));
  let pl;
  try { pl = Music.userPlaylists.byName(%s); pl.name(); Music.delete(pl.tracks); }
  catch (e) { pl = Music.make({new: 'playlist', withProperties: {name: %s}}); }
  for (const t of res) Music.duplicate(t, {to: pl});
  pl.play();
}
JSON.stringify({queued: res.length});
`, q, limit, start, name, name)
}
