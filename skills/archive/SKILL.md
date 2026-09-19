---
name: arcdlc-archive
description: Move finished tasks out of the plan into a dated archive and compact what remains, so the plan stays small enough to stay in an agent's context without losing history. Use when someone says the plan is getting long, clean up the finished tasks, compact the plan, or archive what is done, or runs /arcdlc:archive <slug>, or invokes arcdlc-archive.
argument-hint: "<slug>"
---

# ArcDLC Archive (/arcdlc:archive)

Shrink `docs/aics/<slug>/plan.md` after tasks complete, without breaking the format contract that `/arcdlc:execute`
depends on (defined in `../plan/references/plan-format.md`; flat installs:
`../arcdlc-plan/references/plan-format.md`).

## Talk simple, write like a human

Every rule here governs everything you emit, chat messages exactly as much as the files you write.
Plain English: short sentences, common words, one idea each, active voice, a named actor.

- **No AI filler, anywhere.** Not "Furthermore", "In conclusion", "It is important to note", "delve",
  "leverage", "robust", "seamless". No warm-up opener, no praise, no restating the request back, no
  invented summary. Concrete names, numbers, and paths, never "significantly improves performance".
- **Be direct, name the thing.** The fewest words that carry the fact. Cut empty intensifiers:
  `honestly`, `genuinely`, `truly`, `clearly`, `obviously`. Never a metaphor, a narrative line, or a
  question. Every title that names work is a verb plus its object: "Implement the alerter service".
- **No long dashes.** A full stop, a comma, a colon, or brackets instead of `—` and `–`. Hyphens,
  flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.
- **Shape.** In chat, bullets rather than paragraphs: what you did, what you found, what comes next.
  In files, vary sentence length so the prose does not read as a list.

Nothing is dropped to save space, in a reply as much as in a file. Never leave out a rule, path,
decision, trade-off, open question, or acceptance criterion because the answer is getting long: if it
bears on what the reader does next, it goes in. Short means no padding, never less content. Cut filler
words, repetition, and throat-clearing; never cut a fact. When completeness makes a reply long, let it
be long. The full standard, with examples and a pre-save check, is
`../grilling/references/Writing Style.md` (flat installs:
`../arcdlc-grilling/references/Writing Style.md`).

## Judge by the four virtues

- **Wisdom.** Unclear is a question, not a guess. A contradiction, a missing decision, code that looks
  dead, a change that is risky or hard to reverse: grill it, never pick for the engineer.
- **Courage.** Say the hard thing. A task that is not mechanical is a plan defect, not a reason to
  raise the tier. Never stop to report a helper skill as missing.
- **Justice.** Write every answer where the next session reads it, never only in the chat. Name the
  tier you actually used. Never report a same-tier spawn as a cheaper run.
- **Temperance.** Touch only what you were pointed at. Cutting a comment out of the code is asked for,
  not assumed. No invented summary, no scope you were not given.

## Initiative selection

The initiative slug is the **required first positional argument**: `/arcdlc:archive <slug>`. If it is
missing, stop and report the error, listing the existing initiatives under `docs/aics/` — never guess.
Compact that folder's `docs/aics/<slug>/plan.md` (passing `--aic <slug>` to `arctool`);
`plan-archive.md` is always written **beside** it.

## In a workspace

The working directory is a workspace when it is not a git work tree and `git -C docs rev-parse
--show-toplevel` resolves to `<cwd>/docs` (see the `Workspace` and `Hub` terms in `CONTEXT.md`). In a
workspace, `docs/` is a repository of its own. Every write this skill makes there, `plan-archive.md`,
is committed in the hub with `git -C docs commit -- <hub-relative paths>` under the message
`docs(<slug>): <what changed>` plus `#AI-assisted`, and pushed to the hub's default branch at once. A
rejected push is retried once after `git -C docs pull --ff-only`. Nothing is written into a product
repository by these skills. `arctool sync` runs from the workspace root and writes through the root
links.

## Prefer `arctool`

Probe once: `command -v arctool`. If present, run `arctool archive --aic <slug>` (`--dry-run` previews).
It does everything below deterministically, writes the archive before the plan so a crash never loses a
`DONE` block, and self-validates (writing nothing, exit 5, on a failed invariant). Report its
`archived N, pending M` line. If `arctool` is absent, say so once and do it by hand.

## Manual fallback

1. Move every `###` block in `docs/aics/<slug>/plan.md` whose `- Status:` is `DONE` into
   `docs/aics/<slug>/plan-archive.md`, full and unmodified, in order, under a dated
   `## Archived <YYYY-MM-DD>` section. Create the archive with a header linking back to the plan if it
   does not exist.
2. Replace them in `plan.md` with one compact ledger near the top. Plain bullets, never `###`, or the
   runner counts ledger lines as task blocks:

   ```md
   Completed (archived to docs/aics/<slug>/plan-archive.md):
   - <TASK-ID>: <Short Title>
   ```

3. Never modify, reorder, or reword a `TODO`, `TAKEN`, or `BLOCKED` block. `TAKEN` means an agent is
   mid-task. Leave `gap.md` alone: it is the evidence register, not the queue.

Before finishing, check all three: the `TODO` count is identical before and after, every archived task
ID appears exactly once in the ledger and once in the archive, and no `DONE` block remains in
`plan.md`. Report how many were archived, how many remain, and the archive path.
