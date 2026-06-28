package port

// VolumeStore is the driven port for persisting a volume level between runs.
// It backs mute/unmute, which must remember the level in effect before muting.
type VolumeStore interface {
	// Save persists a volume level.
	Save(level int) error
	// Load returns the persisted level and whether one was present.
	Load() (level int, ok bool, err error)
}
