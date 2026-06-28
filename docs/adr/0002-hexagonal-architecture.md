# ADR-0002: Hexagonal architecture

- Status: accepted
- Date: 2026-06-28

## Context

The application must control Apple Music from the terminal and will grow several
delivery mechanisms over time — a one-shot CLI, a background daemon, and an
interactive TUI — all driving the same underlying behaviour. The mechanism we
use to control the music engine (AppleScript via `osascript` today, possibly a
daemon connection later) must be swappable without rewriting the clients, and
the core logic must be testable without a running Music.app.

## Decision

We adopt a hexagonal (ports & adapters) architecture:

- **Domain** (`internal/music`): pure entities and value objects (player state,
  track, volume). No I/O, no dependencies.
- **Ports** (`internal/port`): interfaces describing what the application needs.
  The primary **driven port** is the music engine (play, pause, status, …).
- **Application** (`internal/app`): use-case services that orchestrate the
  domain and depend only on ports.
- **Adapters** (`internal/adapter/...`): the outside edges.
  - Driven adapter `applescript` implements the engine port via `osascript`.
  - Driving adapter `cli` translates terminal commands into application calls.

Composition (wiring concrete adapters to ports) happens only in `cmd/*/main.go`.

## Consequences

- The domain and application layers are unit-testable with fakes, giving a wide,
  fast base for the test pyramid; `osascript` is exercised by a thin layer of
  integration tests behind a build tag.
- A later daemon transport is just another adapter implementing the same engine
  port — clients do not change. (See the roadmap; tracked separately.)
- Slightly more indirection up front than calling `osascript` directly from
  command handlers; accepted for the testability and swappability it buys.

## Alternatives considered

- **Direct `osascript` calls from command handlers** (the legacy shell design).
  Fast to write but couples every feature to the shell-out, untestable without
  Music.app, and offers no seam for the future daemon/TUI.
- **Layered/MVC.** Does not express the "swappable outside" requirement as
  cleanly as ports & adapters for a tool whose entire job is driving an external
  engine through interchangeable transports.
