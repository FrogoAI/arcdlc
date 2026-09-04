# ADR-0014 — Exit code 5 means "self-validation failed", not "archive self-validation failed"

- Status: Accepted
- Date: 2026-09-04
- Initiative: [ordering](../aics/ordering/aic.md)

## Context

`arctool` exit codes are part of its interface. Skills branch on them, so `AGENTS.md` forbids
renumbering: `0` ok, `1` contract failure, `2` usage, `3` not found or empty, `4` I/O, `5` archive
self-validation.

Code `5` exists because `arctool archive` rewrites the whole plan file. Before writing, it re-parses
the bytes it is about to write and checks that no `DONE` block survived, that the pending counts did
not change, and that every archived ID appears once in the ledger and once in the archive section.
If that check fails, it writes nothing and exits `5`.

`arctool order` is the second whole-file rewrite. A permutation bug there silently destroys a plan
the engineer cannot reconstruct, so it needs the same guard. That guard needs an exit code.

## Decision

Widen code `5` from "archive self-validation" to "self-validation failed (nothing was written)".

`arctool order` re-parses its own output before writing and asserts four things:

- the number of task blocks is unchanged,
- the set of task IDs is unchanged,
- the status counts per status are unchanged,
- every block's content is byte-identical to the original, ignoring trailing newlines.

Any failure means the file is left untouched and the command exits `5`. No new exit code is added.

## Justification

- **The meaning a skill acts on does not change.** A skill that already treats `5` as "the tool
  refused to write a corrupted file, your file is untouched" keeps working with no edit. That is the
  behaviour every caller branches on, and it is identical for both commands.
- **Widening wording is not renumbering.** `AGENTS.md` bans renumbering because a skill keyed to a
  number would break. Nothing here changes a number.
- **A sixth code buys nothing.** The caller's recovery is the same either way: stop, tell the
  engineer, do not retry. A code that never leads to a different action is a code nobody reads.

## Trade-offs

- **Every place that spells out the exit table has to change in the same change set.** That is
  `AGENTS.md`, the `usage` string in `cmd/arctool/main.go`, `README.md`, and any `SKILL.md` that
  lists the codes. Missing one leaves a stale table, which is a documentation bug rather than a
  behaviour bug, but it is still wrong.
- **An operator reading `5` in a log now has to look at the command to know which check failed.**
  Small cost: the message on stderr names the command and the failed invariant.
