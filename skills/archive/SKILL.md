---
description: Archive DONE task blocks from docs/aics/<slug>/plan.md into docs/aics/<slug>/plan-archive.md and compact the plan, keeping it short for agent context while preserving history. Use when the user runs /arcdlc:archive, invokes arcdlc-archive, or asks to archive/compact the plan.
argument-hint: "[--aic <slug>]"
---

# ArcDLC Archive (/arcdlc:archive)

Shrink `docs/aics/<slug>/plan.md` after tasks complete, without breaking the format contract that `/arcdlc:execute`
depends on (defined in `../plan/references/plan-format.md`; flat installs:
`../arcdlc-plan/references/plan-format.md`).

## Initiative selection

Resolve which initiative to compact with `--aic <slug>`, or auto-detect the single one under `docs/aics/`; if several
exist and no `--aic` is given, list them and ask. `plan-archive.md` is always written **beside** that folder's
`plan.md`.

## Prefer `arctool`

Probe once: `command -v arctool`. If present, run `arctool archive --aic <slug>` (or add `--dry-run` first to preview;
omit `--aic` to auto-detect). It performs every step below deterministically — moves the `DONE` blocks into
`docs/aics/<slug>/plan-archive.md` under a dated `## Archived <YYYY-MM-DD>` section, extends the single compact ledger,
leaves `TODO`/`TAKEN`/`BLOCKED` blocks untouched, and self-validates (writing nothing, exit 5, if an invariant fails).
It writes the archive before the plan so a crash never loses a `DONE` block. Report its `archived N, pending M` line.
If `arctool` is absent, say so once and follow the manual steps below.

## What to do

1. Read `docs/aics/<slug>/plan.md` (the resolved initiative's plan).
2. Move every `###` block whose `- Status:` is `DONE` into `docs/aics/<slug>/plan-archive.md`:
   - Create the archive file if missing, with a header linking back to the plan.
   - Append the full, unmodified task blocks under a dated section (`## Archived <YYYY-MM-DD>`), preserving order.
3. In `plan.md`, replace the archived blocks with a single compact ledger near the top:

   ```md
   Completed (archived to docs/aics/<slug>/plan-archive.md):
   - <TASK-ID>: <Short Title>
   - <TASK-ID>: <Short Title>
   ```

   Use a plain bullet list — never `###` headings — so the runner does not count ledger lines as task blocks.
4. Never modify, reorder, or reword blocks with status `TODO`, `TAKEN`, or `BLOCKED`. A `TAKEN` block means an agent
   is mid-task — leave it alone.
5. If `docs/aics/<slug>/gap.md` mirrors the plan, leave it untouched; it is the evidence register, not the queue.

## Validate before finishing

- The number of pending `TODO` blocks in `plan.md` is identical before and after compaction.
- Every archived task ID appears exactly once in the ledger and once in the archive.
- No `DONE` blocks remain in `plan.md`.

Report to the user: how many tasks were archived, how many remain pending, and the archive path.
