# ADR-0001: Record architecture decisions

- Status: accepted
- Date: 2026-06-28

## Context

This project is being rebuilt from a small shell script into a real Go
application that will grow incrementally across a CLI, a daemon, and a TUI. We
want the reasoning behind significant structural choices to be durable and
reviewable rather than living only in commit messages or someone's memory.

## Decision

We record every architecturally significant decision as a numbered Architecture
Decision Record (ADR) in `docs/adr/`, using the lightweight template in
`0000-template.md`. ADRs are immutable once accepted; a reversal is a new ADR
that supersedes the old one.

## Consequences

- New contributors can read `docs/adr/` to understand why the system looks the
  way it does.
- A small ongoing cost: significant decisions must be written down before or
  alongside the code that implements them.

## Alternatives considered

- **A single design doc.** Drifts out of date and loses the history of *why*
  choices changed.
- **Commit messages only.** Hard to discover and read as a coherent narrative.
