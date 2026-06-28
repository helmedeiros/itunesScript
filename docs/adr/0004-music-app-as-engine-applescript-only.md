# ADR-0004: Music.app is the engine, driven via AppleScript only

- Status: accepted
- Date: 2026-06-28

## Context

Unlike MPD, we are not building a playback engine. macOS already ships one:
**Music.app** owns playback, the library, playlists, search, ratings, and
volume. We need to decide how to observe and control it. The options span
AppleScript/JXA via `osascript`, Apple's private **MediaRemote** framework, and
system audio capture for a spectrum visualizer.

## Decision

Treat **Music.app as the engine** and drive it **exclusively through
AppleScript/JXA via `osascript`**. No private frameworks, no audio capture.

Three constraints follow and are designed for, not fought:

1. **No real-time push.** Now-playing state is obtained by polling.
2. **The play queue is emulated.** Music's live "Up Next" is not reliably
   scriptable, so our queue is modelled with a managed playlist.
3. **No spectrum visualizer.** Music.app exposes no PCM stream; the spectrum
   pane from the reference UI is dropped.

Performance rule: read a full status record in a **single batched `osascript`
call** returning structured data, never one call per field.

## Consequences

- Maximum stability and zero dependency on Apple-private/entitled APIs; works
  via a tool present on every Mac (`/usr/bin/osascript`).
- The engine port has a single, well-understood adapter; all `osascript`
  knowledge is isolated there and covered by integration tests.
- We accept polling latency and a feature gap (no live push, no spectrum, an
  emulated queue) in exchange for that stability.

## Alternatives considered

- **MediaRemote adapter** for real-time now-playing and fast artwork. Richer and
  more responsive, but relies on a private framework that Apple has been locking
  down; added compatibility/fragility risk we chose not to take on now.
- **System audio tap** (BlackHole / CoreAudio) to power a real spectrum
  analyzer. Authentic to the reference UI but invasive to set up and fragile;
  out of scope.
