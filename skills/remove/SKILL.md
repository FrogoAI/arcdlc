---
name: arcdlc-remove
description: Delete an initiative's folder for good, design included, after showing what is lost and an explicit confirmation. Use when someone says delete the design for good, purge the initiative, or throw away that initiative entirely, or runs /arcdlc:remove <slug>, or invokes arcdlc-remove.
argument-hint: "<slug>"
---

# ArcDLC Remove (/arcdlc:remove)

Delete an initiative's folder for good. This skill deletes the whole `docs/aics/<slug>/` folder, open
or closed, and refreshes the initiative registry in `AGENTS.md` and `README.md`. Git history is the
only copy afterwards. Finishing an initiative, with the design kept, is `/arcdlc:close`.

Removal is **destructive** and always requires an explicit human confirmation. `arctool` itself has no
delete command (it stays non-destructive); this skill does the deletion.

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

## Argument: initiative slug (required, first positional)

The initiative slug is the **first positional argument** and is **required**: `/arcdlc:remove <slug>`.

- If no slug is given, stop and report the error, listing the existing initiatives under `docs/aics/`
  so the user can pick one. Never guess.
- Resolve the folder `docs/aics/<slug>/`. If it does not exist, stop and say so, listing what does.

## In a workspace

The working directory is a workspace when it is not a git work tree and `git -C docs rev-parse
--show-toplevel` resolves to `<cwd>/docs` (see the `Workspace` and `Hub` terms in `CONTEXT.md`). In a
workspace, `docs/` is a repository of its own. Every write this skill makes there, the deleted folder
and the two registry files, `AGENTS.md` and `README.md`, is committed in the hub with `git -C docs
commit -- <hub-relative paths>` under the message `docs(<slug>): <what changed>` plus `#AI-assisted`,
and pushed to the hub's default branch at once. A rejected push is retried once after `git -C docs
pull --ff-only`. Nothing is written into a product repository by these skills. `arctool sync` runs
from the workspace root and writes through the root links. This skill deletes the folder with
`git -C docs rm -r aics/<slug>` (hub-relative) and refreshes the registry with `arctool sync` run from
the root; the deletion and the two registry files are one hub commit, `docs(<slug>): remove the
initiative`, pushed at once.

## Step 1 — Show what will be removed

Before deleting anything, report the blast radius so the engineer decides with full information:

- The initiative **title** (the first `# ` heading of its architecture document).
- **Task counts by status.** Prefer `arctool` (probe once with `command -v arctool`):
  `arctool list --aic <slug>` prints the `TODO/TAKEN/DONE/BLOCKED` tallies. *Fallback (no `arctool`):
  read `docs/aics/<slug>/plan.md` and count `- Status:` lines by value.*
- The **file list** under `docs/aics/<slug>/` (e.g. `aic.md`, `plan.md`, `gap.md`, `plan-archive.md`,
  `CLOSED.md`).
- When `CLOSED.md` exists, its `- Closed:` and `- Outcome:` lines, so the engineer sees this
  initiative was already finished.

If any task is **`TAKEN`**, stop and change nothing: name every `TAKEN` task by ID, because an agent is
working on it right now.

If any task is `TODO` or `BLOCKED`, warn loudly and name the counts: removal discards not-yet-completed
work. Removal is still allowed (abandoned initiatives also need cleanup), but the engineer must
acknowledge it.

## Step 2 — Require explicit confirmation (always)

Ask the engineer to confirm the deletion, every time. There is no flag that skips this. State, every
time:

- The folder path.
- That the design leaves the working tree and survives afterwards only in git history.
- That other files may still reference this design, and remove does not rewrite links: those
  references will be left to point at nothing.
- That the engineer can cancel now, to handle the references first.

Proceed only on an explicit yes. On anything else, stop and change nothing.

## Step 3 — Delete and clean the registry

1. Delete the folder. If it is tracked by git, use `git rm -r docs/aics/<slug>/` (leaves the deletion
   staged for the user to commit). Otherwise remove it from the working tree (`rm -rf
   docs/aics/<slug>/`).
2. Refresh the registry so the removed initiative disappears from `AGENTS.md` and `README.md`. Prefer
   `arctool sync` (writes only the `<!-- arcdlc:initiatives -->` marker blocks); it also refreshes the
   closed-count line at the end of the block. *Fallback (no `arctool`): delete the initiative's bullet
   from those marker blocks by hand; if no initiatives remain, leave the block reading `_none_`. Then
   write or update the closed-count line the same way `arctool sync` does: a blank line, then ``N closed
   initiatives: run `arctool status`, or see `CLOSED.md` in each folder.`` (`1 closed initiative:
   ...` for one), counting the folders that still hold `CLOSED.md`.*

Do not commit on the user's behalf unless they ask; leave the deletion and registry edit staged so they
can review.

## Step 4 — Report

Confirm what was removed: the slug, the folder, and that the registry was refreshed. If the initiative
had unfinished tasks, restate that they were discarded and note that `git` history still holds them.

The plan format and status contract this skill reads task counts from live in
`../plan/references/plan-format.md` (flat installs: `../arcdlc-plan/references/plan-format.md`).
