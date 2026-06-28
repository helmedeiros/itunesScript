package music

// VolumeMin and VolumeMax bound the sound volume Music.app accepts.
const (
	VolumeMin = 0
	VolumeMax = 100
)

// Volume is a sound-volume level constrained to the range [VolumeMin, VolumeMax].
type Volume int

// NewVolume returns a Volume, clamping n into the valid range.
func NewVolume(n int) Volume {
	switch {
	case n < VolumeMin:
		return VolumeMin
	case n > VolumeMax:
		return VolumeMax
	default:
		return Volume(n)
	}
}

// Adjust returns the volume shifted by delta, clamped to the valid range.
func (v Volume) Adjust(delta int) Volume {
	return NewVolume(int(v) + delta)
}

// Int returns the volume as a plain integer.
func (v Volume) Int() int {
	return int(v)
}

// IsMuted reports whether the volume is at its minimum.
func (v Volume) IsMuted() bool {
	return v == VolumeMin
}
