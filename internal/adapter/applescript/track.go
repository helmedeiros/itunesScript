package applescript

import "github.com/helmedeiros/amp/internal/music"

// trackDTO mirrors the JSON shape Music.app track queries emit.
type trackDTO struct {
	Name     string  `json:"name"`
	Artist   string  `json:"artist"`
	Album    string  `json:"album"`
	Duration float64 `json:"duration"`
}

// toTrack maps the DTO onto the domain Track.
func (d trackDTO) toTrack() music.Track {
	return music.Track{
		Name:     d.Name,
		Artist:   d.Artist,
		Album:    d.Album,
		Duration: seconds(d.Duration),
	}
}
