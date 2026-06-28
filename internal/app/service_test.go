package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helmedeiros/itunesScript/internal/app"
	"github.com/helmedeiros/itunesScript/internal/music"
)

// fakePlayer is a hand-rolled in-memory Player for testing the service in
// isolation from any real engine.
type fakePlayer struct {
	status    music.Status
	statusErr error
	volumeSet *music.Volume
	calls     []string
	setVolErr error
}

func (f *fakePlayer) Status(context.Context) (music.Status, error) {
	f.calls = append(f.calls, "Status")
	return f.status, f.statusErr
}

func (f *fakePlayer) Play(context.Context) error  { f.calls = append(f.calls, "Play"); return nil }
func (f *fakePlayer) Pause(context.Context) error { f.calls = append(f.calls, "Pause"); return nil }
func (f *fakePlayer) Stop(context.Context) error  { f.calls = append(f.calls, "Stop"); return nil }
func (f *fakePlayer) Next(context.Context) error  { f.calls = append(f.calls, "Next"); return nil }
func (f *fakePlayer) Previous(context.Context) error {
	f.calls = append(f.calls, "Previous")
	return nil
}

func (f *fakePlayer) TogglePlayPause(context.Context) error {
	f.calls = append(f.calls, "TogglePlayPause")
	return nil
}

func (f *fakePlayer) SetVolume(_ context.Context, v music.Volume) error {
	f.calls = append(f.calls, "SetVolume")
	if f.setVolErr != nil {
		return f.setVolErr
	}
	f.volumeSet = &v
	return nil
}

func TestServiceStatusPassesThrough(t *testing.T) {
	t.Parallel()

	want := music.Status{State: music.Playing, Track: music.Track{Name: "Gorgon"}}
	fake := &fakePlayer{status: want}
	svc := app.NewService(fake)

	got, err := svc.Status(context.Background())

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestServiceTransportDelegates(t *testing.T) {
	t.Parallel()

	fake := &fakePlayer{}
	svc := app.NewService(fake)
	ctx := context.Background()

	require.NoError(t, svc.Play(ctx))
	require.NoError(t, svc.Pause(ctx))
	require.NoError(t, svc.Toggle(ctx))
	require.NoError(t, svc.Stop(ctx))
	require.NoError(t, svc.Next(ctx))
	require.NoError(t, svc.Previous(ctx))

	assert.Equal(t,
		[]string{"Play", "Pause", "TogglePlayPause", "Stop", "Next", "Previous"},
		fake.calls,
	)
}

func TestServiceSetVolumeClamps(t *testing.T) {
	t.Parallel()

	fake := &fakePlayer{}
	svc := app.NewService(fake)

	got, err := svc.SetVolume(context.Background(), 250)

	require.NoError(t, err)
	assert.Equal(t, 100, got.Int())
	require.NotNil(t, fake.volumeSet)
	assert.Equal(t, 100, fake.volumeSet.Int())
}

func TestServiceAdjustVolumeReadsThenSets(t *testing.T) {
	t.Parallel()

	fake := &fakePlayer{status: music.Status{Volume: music.NewVolume(95)}}
	svc := app.NewService(fake)

	got, err := svc.AdjustVolume(context.Background(), 10)

	require.NoError(t, err)
	assert.Equal(t, 100, got.Int(), "95 + 10 clamps to 100")
	require.NotNil(t, fake.volumeSet)
	assert.Equal(t, 100, fake.volumeSet.Int())
	assert.Equal(t, []string{"Status", "SetVolume"}, fake.calls)
}

func TestServiceAdjustVolumeSurfacesStatusError(t *testing.T) {
	t.Parallel()

	boom := errors.New("osascript failed")
	fake := &fakePlayer{statusErr: boom}
	svc := app.NewService(fake)

	_, err := svc.AdjustVolume(context.Background(), 10)

	require.ErrorIs(t, err, boom)
	assert.NotContains(t, fake.calls, "SetVolume", "must not set volume if read failed")
}
