# ADR-0005: Mute by zeroing volume with a persisted prior level

- Status: accepted
- Date: 2026-06-28

## Context

Music.app exposes no mute property through AppleScript — only `sound volume`
(0–100). Muting therefore means setting the volume to 0, and unmuting means
restoring the level that was in effect before. Because the CLI is one-shot, that
prior level cannot live in memory between `mute` and `unmute`; it must be
persisted somewhere.

## Decision

- `mute` reads the current volume and, if it is above 0, stores it and sets the
  volume to 0. If already 0, it does nothing (so a double mute does not lose the
  remembered level).
- `unmute` restores the stored level; if nothing is stored it falls back to a
  default level.
- Persistence is modelled as a new driven port, `VolumeStore`, implemented by a
  small file adapter writing under the user config directory. The application
  depends only on the port, consistent with ADR-0002.

## Consequences

- Mute/unmute behave intuitively across separate CLI invocations.
- A second driven port demonstrates the architecture scaling to more than one
  outbound dependency; the file adapter is unit-tested against a temp directory.
- A tiny state file is written to disk. Accepted; it is namespaced and
  self-healing (a missing or unreadable file falls back to the default).

## Alternatives considered

- **Unmute to a fixed default every time.** Simpler, no persistence, but loses
  the user's prior level — surprising after raising the volume well above the
  default.
- **Store the level inside a Music playlist or comment field.** Abuses the
  library and is slower; rejected.
