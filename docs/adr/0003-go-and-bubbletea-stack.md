# ADR-0003: Go with a Bubble Tea TUI

- Status: accepted
- Date: 2026-06-28

## Context

The project was a collection of POSIX-ish shell scripts. To build a testable,
maintainable application with a rich terminal UI we need a real language and a
TUI toolkit. The reference experience is a modern MPD TUI client (tabs, live
header, album art, vim-style keys).

## Decision

Implement the application in **Go**, distributed as a single static binary.

- CLI: a conventional command framework (one-shot, scriptable, JSON output).
- TUI: **Bubble Tea** with **Lipgloss** and **Bubbles** widgets.
- Engine access: shell out to `osascript` (see ADR-0004).

Quality gates: `golangci-lint` v2 (lint + gofumpt/goimports formatting), `go
vet`, and race-enabled tests run via `gotestsum`, all wired through the Makefile
and CI.

## Consequences

- Single self-contained binary, easy to install and run; strong standard
  library for process execution and JSON.
- Fast compile/test loop supports TDD.
- Bubble Tea's Elm-style model fits an event-driven TUI that consumes engine
  state updates.
- Team must follow Go idioms and keep the gates green on every commit.

## Alternatives considered

- **Rust + ratatui** (what the reference client uses). Best-in-class TUI and
  single binary, but a steeper build/iteration loop and a larger jump from the
  existing scripts.
- **Stay in shell.** No path to a testable core or a real TUI.
- **Python + Textual.** Quick to prototype but heavier, slower runtime and more
  awkward distribution than a Go binary.
