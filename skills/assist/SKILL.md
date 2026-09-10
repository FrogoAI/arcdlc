---
name: arcdlc-assist
description: Turn the code comment markers a team already wrote (TODO by default, also FIXME, HACK, XXX, BUG) into planned work. Sweeps the source with `arctool scan`, records every marker in docs/aics/<slug>/comments.md, grills the engineer about the unclear ones, then adds a TODO task to docs/aics/<slug>/plan.md for each marker that names real work. The initiative slug is the required first argument; the marker is an optional second one (e.g. /arcdlc:assist payments FIXME). Use when the user runs /arcdlc:assist, invokes arcdlc-assist, or asks to turn code TODOs into tasks.
argument-hint: "<slug> [TODO|FIXME|HACK|XXX|BUG]"
---

# ArcDLC Assist (/arcdlc:assist)

Collect the markers the team left in the code, judge each one with the engineer, and feed the real
work into the executable plan. The sweep moves each marker out of the code and into the register, so
the comment is recorded once and the code stops carrying it.

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

The initiative slug is the **required first positional argument**: `/arcdlc:assist <slug> [MARKER]`.
The register and its mirrored tasks are filed into `docs/aics/<slug>/` (as `comments.md` and
`plan.md`), and the slug is passed to `arctool` as `--aic <slug>`. If the slug is missing, stop and
report the error, listing the existing initiatives under `docs/aics/` — never guess. If the named
initiative folder does not exist yet (a fresh sweep, e.g. `comment-debt`), confirm the slug with the
engineer before you sweep: `arctool scan` creates the folder itself and says so, so a mistyped slug
would invent an initiative nobody asked for. If `docs/aics/<slug>/plan.md` does not exist yet, create it
per `../plan/references/plan-format.md` (flat installs: `../arcdlc-plan/references/plan-format.md`).

`MARKER` is the optional second argument and defaults to `TODO`. One marker per run keeps the
register and the report about one kind of debt. Pass a comma-separated list only when the engineer
asks for several at once.

## Step 1 — Sweep the code

Prefer `arctool`, which does the sweep for free and keeps you out of the files:

- Probe once with `command -v arctool` (or install it from the arcdlc repo root: `make install`).
- Found: `arctool scan --marker <MARKER> --aic <slug>`. It writes `docs/aics/<slug>/comments.md`,
  one block per marker, **and deletes each marker comment from the code**, so the register becomes the
  one place that marker text lives. It prints a summary: markers found, blocks appended, comments
  removed per file, and any marker it refused to touch. Add `--json` for the findings as data,
  `--path <dir>` to sweep one subtree, and `--dry-run` to look first: a dry run writes nothing and
  edits nothing.
- The sweep changes the working tree. Say so in your report and point the engineer at `git diff`. Run
  `--dry-run` first when the engineer has not seen the list yet, or when the tree already has
  uncommitted changes.
- A marker inside a multi-line block comment is left in the code and listed as skipped, because
  deleting one line of such a comment can leave a dangling opener. Its block is still in the register.
  Delete the comment yourself when you finish the task it produced.
- Exit `3` means no marker in the tree and nothing to resolve. Report that and stop.
- Exit `1` means the existing register broke its own format (a block with no `- Marker:` line). Fix
  that block by hand, then scan again.
- If `arctool` is unavailable, say so once and sweep by hand with the same rules: the marker must be
  the first word of a comment (`// TODO …`, `# TODO …`), matched with word boundaries, so `TODOLIST`
  and a sentence that merely mentions TODO are not findings; skip `.git`, `vendor`, `node_modules`,
  `dist`, `bin`, `docs`, and every document (`.md`, `.txt`, `.rst`). Then write the register by hand
  in the block shape below.

**What `arctool scan` owns, and what you must not touch:** the task ID and the `- Marker:` line (the
finding's identity, file plus marker text). Never edit either: the next scan matches on them, and a
changed one is re-registered as a new finding. The register is append-only, so a block is never
rewritten or removed, only filled in.

## Step 2 — Judge every finding

Read the register (or `arctool scan --json` for the list alone). Every block arrives as
`- Verdict: NEW.`; replace that with one of four:

- `ACTIONABLE`: the comment names what changes and where, so you can write `WHAT`, `HOW` and
  `WHERE` without asking anybody. This is the only verdict that becomes a task.
- `UNCLEAR`: the comment names a wish, not work ("TODO maybe rethink this"). Goes to Step 3.
- `STALE`: the code already does what the comment asks. No task, and nothing left to do: the sweep
  already took the comment out.
- `DEFERRED`: real work the engineer decided not to plan now. Write the reason in `WHY`. The register
  is the only record of it now, so the reason has to be readable a year later.

Rewrite each heading title into an instruction: an imperative verb plus its object, with the real
component named (`Move the index rebuild into its own package`), never the marker text as it stands.

## Step 3 — Grill the unclear markers (mandatory)

Do not write a task from a guess, and do not drop an unclear marker in silence. Run a grilled
interview on the `UNCLEAR` findings only:

- Prefer the bundle's `arcdlc-grilling` skill. If it cannot be invoked here, read its `SKILL.md`
  (`../grilling/SKILL.md`, flat installs `../arcdlc-grilling/SKILL.md`) and run the same protocol
  inline. Never stop to report a helper skill as missing.
- **One question per turn.** Ask, wait for the answer, then ask the next. Every question carries a
  recommended answer. Never a numbered round of questions.
- Each answer settles the finding: it becomes `ACTIONABLE` with the decision written into `HOW`, or
  `STALE` or `DEFERRED` with the reason written down.

## Step 4 — Finish the register

Fill the keys `arctool scan` left empty, writing for a **less capable executor**: the model running
`/arcdlc:execute` sees the task block and its references, nothing else.

```md
### <MARKER>-CMT-NN: <Short Title, an imperative verb plus its object>

- Marker: `<file>:<line>` `<the marker text, written by arctool scan>`
- Verdict: ACTIONABLE.
- WHAT: <The change, one line. Name the change, not the complaint.>
- HOW:
  <The decision the comment already carries, plus whatever the grilling settled: target package,
  new name, exported signature, data shape, edge cases, error handling.
  Out of scope: <adjacent code the executor must leave alone>.>
- WHERE: <the marker's file and line, plus every file the change touches, one `Layer` line per layer>
- WHY: <What the marker costs while it stands, one line.>
- Acceptance:
  - GIVEN <precondition> WHEN <the runnable check: a named test, a command> THEN <observable result>.
  - GIVEN <precondition> WHEN <another check> THEN <observable result>.
```

Rules that keep the register usable:

- Keep `Layer <name>: <files>` lines in `WHERE`, and never write a bare `file:line` there: the runner
  reads any `name: targets` line in `WHERE` as a layer, so `internal.go:3` would arrive as a layer
  called `internal.go`. `arctool scan` writes `internal.go (marker at line 3)` for that reason.

- Acceptance criteria are about the work, never about the comment: the sweep already removed it.
- Prefer a runnable check (a named test, a lint rule, a command) over "look and see". This is also
  what lets the mirrored task pass `arctool validate --strict`, which requires an `Acceptance`
  section.
- Number blocks sequentially per marker, continuing from what is there: `TODO-CMT-01`,
  `FIXME-CMT-02`. Never renumber, never delete a block, never reuse an ID. `arctool scan` does this
  for you.
- Never delete a block, not even for a finding you judged `STALE`. The comment is gone from the code,
  so the register is the only history left.

## Step 5 — Mirror the register into the plan

Per the Register Sync rules in `../plan/references/plan-format.md` (flat installs:
`../arcdlc-plan/references/plan-format.md`), append a task block to `docs/aics/<slug>/plan.md` for
every `ACTIONABLE` finding:

- Same task ID and heading; same `WHAT`, `HOW`, `WHERE`, `WHY` and `Acceptance` content.
- Drop the `- Marker:` and `- Verdict:` lines: they belong to the register, not to the runner.
- Add runner metadata: `- References:` must include `docs/aics/<slug>/comments.md`, plus the
  architecture document or any ADR the task relies on. Close the block with `- Status: TODO.`
- **No parenthetical tag in the heading.** The only valid tags are `MISSING`, `PARTIAL` and `DRIFT`,
  which belong to gap-derived tasks; `arctool validate` warns on anything else.
- Append after the existing blocks. Never modify an existing task, never reuse an ID already in the
  plan.
- Order matters: the runner works top to bottom, so a task may only depend on tasks above it. When
  two markers touch the same file, put the one the other needs first (`arctool order` fixes an
  inversion later).
- Validate before handing off. Prefer `arctool validate --strict --aic <slug>` and fix every finding;
  exit `0` means clean. If `arctool` is unavailable, say so once and hand-check unique IDs, present
  and uppercase `Status`, and the required keys per the format guide.

## Step 6 — Report

Say it in four numbers: markers found, findings per verdict, tasks added to the plan, comments removed
from the code (and from how many files). Name every marker you grilled and what the engineer decided,
and every marker the sweep left in place. Tell the engineer to review the code change with `git diff`.
Then name the next step: `/arcdlc:execute <slug>` implements the queue.
