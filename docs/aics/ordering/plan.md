# Task Ordering — Plan

Task format: `skills/plan/references/plan-format.md`. Architecture: [aic.md](aic.md).

All paths are relative to the repo root. Run `go build ./... && go test ./... && gofmt -l . && go vet ./...`
before every commit: the repo must stay green.

## Risk Coverage

- **R1, a permutation bug destroys a plan** covered by ORD-1: the self-check runs inside `Reorder`, so
  no caller can skip it, and the acceptance criteria name the index-maths cases (one ID, all IDs,
  already-in-order, first and last swapped, single-block plan).
- **R2, the exit table lives in several places** covered by ORD-2 (the `usage` string in
  `cmd/arctool/main.go`) and ORD-3 (`AGENTS.md`). `README.md` has no exit table, so ORD-3 covers only
  its command list and the sentence that stops being true.
- **R3, inter-block content moves with the block above it** covered by ORD-1: a test pins the
  behaviour so it cannot change by accident.
- **R4, line endings** covered by ORD-1: a test pins that a CRLF plan keeps CRLF inside each block and
  gets LF between blocks, matching `Archive`.
- **R5, duplicate task IDs in a plan** covered by ORD-1 (the ambiguous-ID rejection) and ORD-2 (exit 3).
- **R6, "move to first" is the phrase people will use** covered by ORD-2 (worked example in `--help`)
  and ORD-4 (worked example in `/arcdlc:plan`).
- **R7, reordering while a task is TAKEN** covered by ORD-4: the `/arcdlc:plan` step says to reorder
  between runs, not during one.
- **R8, an uncommitted change set sits in the same files** accepted as a process mitigation, not a
  task: land the in-flight writing-style change set first, then execute this plan on top. Never run
  two agents across `AGENTS.md`, `README.md` and the `SKILL.md` files at once.
- **Q1, should `/arcdlc:examinate` place gap tasks by order instead of appending?** Accepted and
  deferred. It is a change to a different skill's behaviour and needs its own interview.
- **Q2, should a later initiative add `- DEPENDS:` plus `arctool order --auto`?** Accepted and
  deferred by [ADR-0013](../../adr/0013-order-is-a-slot-permutation.md). The command line designed
  here does not change if it is added.
- **Q3, should `arctool validate` warn on irregular block spacing?** Accepted and deferred. Both
  `order` and `archive` normalise spacing silently today, which is consistent, so nothing is broken
  while this stays open.

Completed (archived to docs/aics/ordering/plan-archive.md):
- ORD-1: Add slot-permutation reordering to `internal/plan`
- ORD-2: Add the `arctool order` command
- ORD-3: Update AGENTS.md and README.md for the new command
- ORD-4: Teach `/arcdlc:plan` and `/arcdlc:execute` about ordering
- ORD-5: Bump arctool to 0.11.0 and the plugin bundle to 0.13.0
