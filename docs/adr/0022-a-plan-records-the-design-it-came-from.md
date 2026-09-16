# ADR-0022 — A plan records the design it came from

- Status: Accepted
- Date: 2026-09-16
- Initiative: none (bundle-wide, extends the plan format and `arctool validate`)
- Builds on: [ADR-0021](0021-a-planned-task-must-be-mechanical.md), which made execution mechanical
  and, in doing so, created the problem this record solves.

## Context

An architecture document matures by re-running `/arcdlc:aic`. Each round adds detail, tests an idea,
or reverses a call. That is the intended workflow, and the skills now support it: a document is
amended rather than regenerated, and a reversed decision supersedes its ADR in the same turn.

Once a plan exists, that workflow has a hole. `/arcdlc:plan` decomposes the document as it stood, and
nothing afterwards notices when the document moves. The plan keeps its tasks, `/arcdlc:execute` keeps
running them, and the work builds a design that was abandoned two rounds ago.

ADR-0021 makes this worse rather than better. A mechanical task carries every decision and no
reasoning, which is exactly what lets a cheap executor run it safely. It also means the executor has
nothing to compare against: it cannot notice that the design changed, because it never saw the design.
It implements the block faithfully and reports success. The failure is silent, and it surfaces after
commits exist.

Prose cannot close this. Telling `/arcdlc:aic` to say "re-run the plan" helps only the engineer who
reads the report and acts on it.

## Decision

A plan records which architecture document it was decomposed from, and what that document said at the
time, as a hash in an HTML comment in the preamble:

```
<!-- arcdlc:source docs/aics/checkout/aic.md sha256:3f78685… -->
```

An HTML comment because it never renders and never parses as a `###` task block, so the format gains
a field without gaining a way to break.

- `/arcdlc:plan` writes it, once per document it used, with `arctool stamp <path> --aic <slug>`.
- `arctool validate` checks it. A changed document is a **warning**, so `--strict` fails while an
  ordinary run reports and continues. A missing document is an **error**: the plan points at
  something that is gone.
- `/arcdlc:execute` runs that check before the first task and stops on the warning.
- `arctool stamp --aic <slug>` refreshes the hash. That is how an engineer says "I read the diff and
  these tasks still hold". Nothing decides it for them, because only a person can tell whether a
  wording change matters.

A plan with no stamp is not checked at all.

## Consequences

- Plans written before this keep validating. The feature is opt-in by presence of the stamp, which is
  what makes it safe to add to a contract other tools parse.
- `arctool` gains one command and one behaviour, and stays pure standard library: `crypto/sha256` and
  `encoding/hex`.
- `Validate` stays a pure function of the plan text. The stamp needs the filesystem, so it is checked
  by `plan.CheckSources`, called from `cmd/arctool`, rather than folded into the validator.
- Re-stamping is a deliberate, recorded act rather than a silent one. `arctool stamp --dry-run` names
  what would change first.
- The warning is advisory by design. Making it an error would block a run over a typo fix in the
  architecture document, and people who cannot proceed learn to skip the check.

## Alternatives considered

- **Compare modification times.** Rejected: a fresh checkout rewrites every mtime, so the check would
  fire constantly and be ignored within a week.
- **Make the mismatch an error.** Rejected for the reason above. `--strict` already promotes it for
  anyone who wants the hard gate, and `/arcdlc:execute` stops on it regardless.
- **Hash each task's `References` instead of the plan's source.** Rejected: it multiplies stamps by
  task for the same signal, and the question being answered is about the design as a whole.
- **Leave it as prose in `/arcdlc:aic`'s report.** Rejected: that is the current state, and it only
  works when the engineer reads the report and acts on it. The pattern this repository follows is to
  turn a prose promise into a mechanical check.
