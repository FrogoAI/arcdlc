---
name: arcdlc-archive
description: Archive DONE task blocks from docs/aics/<slug>/plan.md into docs/aics/<slug>/plan-archive.md and compact the plan, keeping it short for agent context while preserving history. The initiative slug is the required first argument (e.g. /arcdlc:archive payments). Use when the user runs /arcdlc:archive, invokes arcdlc-archive, or asks to archive/compact the plan.
argument-hint: "<slug>"
---

# ArcDLC Archive (/arcdlc:archive)

Shrink `docs/aics/<slug>/plan.md` after tasks complete, without breaking the format contract that `/arcdlc:execute`
depends on (defined in `../plan/references/plan-format.md`; flat installs:
`../arcdlc-plan/references/plan-format.md`).

## Talk simple, write like a human

Plain B2 English, everywhere: short sentences, one idea each, active voice, a named actor. Cut
empty intensifiers: `honestly`, `genuinely`, `truly`, `clearly`, `obviously` add nothing.

- **Replies to the user:** bullets, not paragraphs. No filler, no praise, no restating the request.
  Say what you did, what you found, what comes next.
- **Files you write:** no AI filler ("Furthermore", "In conclusion", "It is important to note",
  "delve", "leverage", "robust", "seamless", "In today's fast-paced world"), no warm-up opener, no
  invented summary. Vary sentence length. Concrete names, numbers, and paths, never "significantly
  improves performance".
- **Be direct, name the thing.** The fewest words that carry the fact; a word that only adds tone
  comes out. Never a metaphor, a narrative line, or a question. Say what must happen: "The alerter
  retries three times, then dead-letters." Every title that names work is an instruction, a verb
  plus its object: "Implement the alerter service", never "One message out, and the three places it
  goes".
- **No long dashes.** Use a full stop, a comma, a colon, or brackets instead of `—` and `–`.
  Hyphens, flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.

Short talk, full content. Brevity is for your replies, never for the files: never drop a rule, path,
decision, trade-off, or acceptance criterion to save space. The full standard, with examples and a
pre-save check, is `source/Writing Style.md` in the bundle's `source-map` skill.

## Initiative selection

The initiative slug is the **required first positional argument**: `/arcdlc:archive <slug>`. If it is
missing, stop and report the error, listing the existing initiatives under `docs/aics/` — never guess.
Compact that folder's `docs/aics/<slug>/plan.md` (passing `--aic <slug>` to `arctool`);
`plan-archive.md` is always written **beside** it.

## Prefer `arctool`

Probe once: `command -v arctool`. If present, run `arctool archive --aic <slug>` (or add `--dry-run` first to preview).
It performs every step below deterministically — moves the `DONE` blocks into
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
