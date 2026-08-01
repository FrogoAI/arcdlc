# ADR-0008 — Domain imports are judged by purity, not by folder

- Status: Accepted
- Date: 2026-08-01
- Initiative: [source-library-cleanup](../aics/source-library-cleanup/aic.md)

## Context

MDCA's dependency rule is that the domain layer must not depend on infrastructure. The standard
expresses this through folder names: §7.2 (line 201) defines `pkg/` as "shared infrastructure", so
§7.1 states the domain may import "standard library only", and Appendix A marks
`internal/domain/... → pkg/...` as "discouraged; allowed only for pure utilities (e.g. `pkg/clock`)".

The folder is a poor proxy for the rule. In Go, `pkg/` conventionally means *importable by other
modules* and `internal/` means *private to this module* — a visibility distinction, not an
architectural-layer one. A package like `pkg/clock` is not infrastructure at all; it merely lives in
the folder the standard reserved for infrastructure. The result is a rule that says "stdlib only" in
one section and "pure utilities are fine" in another, with no test for what counts as pure.

## Decision

State the rule as what it actually means: **the domain layer must not import infrastructure.**
Whether an import is allowed is decided by the *properties of the imported package*, not by the
folder it sits in.

A package is importable by the domain layer when both hold:

- its own imports are standard library only (transitively, no third-party SDK, no other project
  package that fails this test); and
- it performs no I/O and holds no mutable global state — no network, no filesystem, no
  `database/sql`, no process or environment access.

Consequently:

- §7.1's "standard library only" is amended to the purity rule above.
- Appendix A's `internal/domain/... → pkg/...` row becomes "allowed when the target package is pure
  (per §7.1); forbidden otherwise" — replacing "discouraged", which was unenforceable.
- The prohibitions on `internal/repository`, `internal/handler`, `database/sql`, `net/http`, vendor
  SDKs, and cross-context `domain/<a> → domain/<b>` are unchanged: those fail the purity test anyway,
  and the cross-context ban is a separate boundary rule.

## Justification

- **The rule becomes checkable.** "Is this package pure?" is answered by reading its import list;
  "is this discouraged?" is not answerable at all, so an auditor had to guess.
- **It matches the intent.** The dependency rule was always about keeping infrastructure concerns out
  of the domain; encoding it as a folder ban was shorthand that broke the moment a non-infrastructure
  package lived in the infrastructure folder.
- **No layout churn.** The alternative fix — splitting `pkg/` into pure and infrastructure trees —
  changes the canonical layout in §7.2 and every project built on it, to enforce by folder name a
  rule the purity test already enforces directly.

## Trade-offs

- **Purity is transitive and must be re-checked.** A pure package that later grows an HTTP client
  silently becomes an illegal domain import. Mitigation: the compliance checklist gains an explicit
  item, and the rule is phrased so a linter can walk the import graph.
- **A clock is a boundary case.** `pkg/clock` passes the "no I/O" bar in the ordinary reading but
  reads wall-clock time, which is exactly the dependency time-injection exists to control. The
  merged document must say plainly that a domain needing controllable time defines a port, and that
  importing a concrete clock is a testability decision, not a layering violation.
- **Folder names no longer tell the whole story.** A reader must inspect a package to know whether
  the domain may import it; the guidance carries the purity test prominently to compensate.
