package applescript

import (
	"encoding/json"
	"fmt"
)

// namesScript builds a JXA program that returns the sorted, de-duplicated,
// non-empty values of a track field ("artist" or "album") across the library.
// field is a fixed internal literal, never user input.
func namesScript(field string) string {
	return fmt.Sprintf(`
const Music = Application('Music');
let out = [];
if (Music.running()) {
  const vals = Music.libraryPlaylists[0].tracks.%s();
  out = [...new Set(vals)].filter(v => v && v.length).sort((a, b) => a.localeCompare(b));
}
JSON.stringify(out);
`, field)
}

// parseNames decodes a JSON array of strings.
func parseNames(raw []byte) ([]string, error) {
	var names []string
	if err := json.Unmarshal(raw, &names); err != nil {
		return nil, fmt.Errorf("decode names: %w", err)
	}
	return names, nil
}
