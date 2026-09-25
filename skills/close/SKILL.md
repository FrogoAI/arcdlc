---
name: arcdlc-close
description: Close a finished initiative in place: write a CLOSED.md note with the outcome, delete its plan, board stories and registers, and stop it taking new work, after showing what is unfinished and asking for a yes. Use when someone says we are done with X, wrap up X, close the initiative, mark it finished, or retire that initiative, or runs /arcdlc:close <slug>, or invokes arcdlc-close.
argument-hint: "<slug>"
---

# ArcDLC Close (/arcdlc:close)

Close keeps the design in `docs/aics/<slug>/` where every link points, marks it with `CLOSED.md`,
deletes `plan.md`, `plan-archive.md`, `plan-human.md`, `gap.md` and `comments.md`, and leaves the
registry through `arctool sync`. A closed initiative is final: no skill changes it or deletes
`CLOSED.md`; follow-up work is a new initiative.

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive` → `/arcdlc:close`, with
`/arcdlc:remove` beside it for deleting a design for good.

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

The initiative slug is the **first positional argument** and is **required**: `/arcdlc:close <slug>`.

- If no slug is given, stop and report the error, listing the existing initiatives under `docs/aics/`
  so the user can pick one. Never guess.
- Resolve the folder `docs/aics/<slug>/`. If it does not exist, stop and say so, listing what does.

## In a workspace

The working directory is a workspace when it is not a git work tree and `git -C docs rev-parse
--show-toplevel` resolves to `<cwd>/docs` (see the `Workspace` and `Hub` terms in `CONTEXT.md`). In a
workspace, `docs/` is a repository of its own. Every write this skill makes there, `CLOSED.md`, the
deleted files, and the two registry files, `AGENTS.md` and `README.md`, is committed in the hub with
`git -C docs commit -- <hub-relative paths>` under the message `docs(<slug>): <what changed>` plus
`#AI-assisted`, and pushed to the hub's default branch at once. A rejected push is retried once after
`git -C docs pull --ff-only`. Nothing is written into a product repository by these skills. `arctool
sync` runs from the workspace root and writes through the root links. This skill deletes the files
with `git -C docs rm aics/<slug>/<file>` (hub-relative); everything lands in one hub commit,
`docs(<slug>): close the initiative` plus `#AI-assisted`, pushed at once, retried once after `git -C
docs pull --ff-only`.

## Ask through the grilling protocol

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

## Step 1 — Read the state

When `CLOSED.md` exists and none of the five files (`plan.md`, `plan-archive.md`, `plan-human.md`,
`gap.md`, `comments.md`) remain, stop: the initiative was already closed, on the date in its
`- Closed:` line.

When `CLOSED.md` exists and some of the five files remain, an earlier close was interrupted. Skip
straight to Step 4, then Step 6: the note is already written, so only the confirmation and the
deletion are left.

Otherwise, read the task counts. Probe once with `command -v arctool`: if present, run `arctool list
--aic <slug>` for the `TODO`/`TAKEN`/`DONE`/`BLOCKED` tallies. *Fallback (no `arctool`): read
`docs/aics/<slug>/plan.md` and count `- Status:` lines by value.* No `plan.md` means zero tasks. The
archived count is the number of `### ` task blocks in `plan-archive.md`; a missing `plan-archive.md`
means zero archived.

If any task is `TAKEN`, stop and change nothing: name each `TAKEN` task by ID, because an agent is
working on it right now.

## Step 2 — Show what closing does

Before asking anything, report the blast radius:

- The task counts, by status, including the archived count.
- Each `TODO` and `BLOCKED` task, by ID, title and status, called out so they are not missed.
- The files that stay: every architecture document, image and other file in the folder.
- The files that will be deleted: whichever of `plan.md`, `plan-archive.md`, `plan-human.md`,
  `gap.md` and `comments.md` exist.

State plainly that a link pointing at one of the deleted files is not rewritten and may break; git
history keeps the target.

## Step 3 — Settle the outcome

With every task `DONE`, or no plan at all, the outcome is `delivered` and nobody is asked. Never write
`delivered` while a task is `TODO` or `BLOCKED`.

Otherwise ask one question, through the grilling protocol above: `partly delivered` or `stopped`.
Recommend `partly delivered` when at least one task is `DONE`, else recommend `stopped`. Then ask for
one line of reason. Both answers go into `CLOSED.md`.

## Step 4 — Confirm

Ask one question naming the folder, the files that will be deleted, and that git history keeps them.
Proceed only on an explicit yes. On anything else, stop and change nothing.

## Step 5 — Write CLOSED.md

Write `docs/aics/<slug>/CLOSED.md` in this exact shape.

Line 1 is `# <title>`, the first `# ` heading of the architecture document: try `aic.md`, then
`arc42.md`, `togaf.md`, `c4.md`, `tsc.md`, then the first other `.md` file left in the folder,
alphabetically, that is not one of the five deleted files. With no architecture document at all, the
title is the slug.

Then a blank line and three list lines:

```
- Closed: YYYY-MM-DD
- Outcome: <delivered | partly delivered | stopped>
- Tasks: <total> (<done> done, <todo> todo, <blocked> blocked)
```

`Closed` is today's date. `Tasks` counts include the archived tasks in `total` and `done`.

When the outcome is not `delivered`, add:

- `## Reason`, with the engineer's one line from Step 3.
- `## Not done`, one bullet per `TODO` or `BLOCKED` task: `` - `<ID>`: <title> (<STATUS>) ``. When
  `comments.md` holds a `### ` block with the same task ID, add indented sub-bullets copying that
  block's `- Marker:` lines and its `- WHERE:` value as they stand, so a marker's text is not lost
  once `comments.md` is deleted.

Then always:

- `## Left open`, with the body under the architecture document's `Open questions` heading, copied as
  it stands, or `None recorded.` when there is none.
- `## Documents`, one bullet `- [<file>](<file>)` per `.md` or `.html` file left in the folder other
  than `CLOSED.md`, or `Documents: none` when none are left.

Close never writes a `## Superseded` section; that heading is added later, by another initiative's
design round.

## Step 6 — Delete and refresh the registry

Delete each of the five files that exists: `git rm <file>` when it is tracked, plain `rm <file>` when
it is not. Then refresh the registry. Probe once with `command -v arctool`: if present, run `arctool
sync`. *Fallback (no `arctool`): delete the initiative's bullet from the `<!-- arcdlc:initiatives -->`
blocks in `AGENTS.md` and `README.md`, and write or update the closed-count line the same way `arctool
sync` does: a blank line, then ``N closed initiatives: run `arctool status`, or see `CLOSED.md` in
each folder.`` (`1 closed initiative: ...` for one).*

Leave the changes staged; do not commit unless asked. In a workspace, commit and push as the workspace
section above says.

## Step 7 — Report

State: the `CLOSED.md` path, the outcome, which files were deleted, that the registry was refreshed,
and that follow-up work is a new initiative with its own slug.
