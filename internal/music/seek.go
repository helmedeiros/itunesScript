package music

// SeekMode describes how a seek value is interpreted.
type SeekMode int

const (
	// SeekAbsolute sets the player position to an absolute number of seconds.
	SeekAbsolute SeekMode = iota
	// SeekRelative shifts the current position by a number of seconds.
	SeekRelative
	// SeekPercent sets the position to a percentage of the track duration.
	SeekPercent
)
