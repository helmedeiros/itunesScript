# amp

A terminal player for **Apple Music** on macOS, modeled on the `mpd` + `mpc`
trio: control playback and read what's playing straight from the console —
handy for hotkeys, remote shells, or scripting your music without leaving the
keyboard.

`amp` treats **Music.app as the engine** and drives it through AppleScript, so
there's no playback daemon to run and nothing to configure. The CLI binary is
`am`.

## Install

Requires Go 1.26+ and macOS with Music.app.

```sh
go install github.com/helmedeiros/amp/cmd/am@latest
```

Or build from a clone:

```sh
git clone https://github.com/helmedeiros/amp.git
cd amp
make build      # produces ./bin/am
```

## Usage

```sh
am status            # show what's playing
am status --json     # machine-readable, for scripts

am play              # transport
am pause
am toggle            # play/pause
am stop
am next
am prev

am vol 60            # set volume to 60%
am vol +10           # raise by 10
am vol -10           # lower by 10
am vol up | down     # ±10
am mute              # silence, remembering the level
am unmute            # restore it

am shuffle           # toggle shuffle
am shuffle on | off
am repeat off | one | all
```

Run `am --help` for the full command list.

## Architecture

`amp` is built with a hexagonal (ports & adapters) architecture in Go:

- **`internal/music`** — pure domain: player state, volume, repeat mode, track,
  status.
- **`internal/port`** — the ports: `Player` and `VolumeStore` (driven),
  `Controller` (driving).
- **`internal/app`** — the application service (use cases).
- **`internal/adapter`** — the edges: `applescript` (the Music.app engine via
  `osascript`), `store` (file-backed volume state), and `cli` (the command
  tree).
- **`cmd/am`** — wiring.

Architectural decisions are recorded in [`docs/adr/`](docs/adr/). The design
keeps `osascript` behind a single seam, so the logic is covered by a wide base
of fast unit tests with a thin layer of integration tests against the real
binary (`make integration`, requires Music.app).

## Development

```sh
make ci            # format check, vet, lint, race-enabled tests
make test          # tests only
make integration   # osascript integration tests (needs Music.app)
make build         # build ./bin/am
```

## License

[MIT](LICENSE)
