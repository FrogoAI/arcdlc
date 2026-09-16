---
name: arcdlc-plan
description: Decompose an approved architecture document (AIC by default; arc42, TOGAF) into the executable docs/aics/<slug>/plan.md task queue consumed by /arcdlc:execute. The initiative slug is the required first argument (e.g. /arcdlc:plan payments). Use when the user runs /arcdlc:plan, invokes arcdlc-plan, or asks to turn an architecture document into an implementation plan.
argument-hint: "<slug> [aic|arc42|tsc|togaf|path]"
---

# ArcDLC Plan (/arcdlc:plan)

Turn an approved architecture document into `docs/aics/<slug>/plan.md` — the executable task queue of the ArcDLC
delivery pipeline:

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`

The plan format is defined in `references/plan-format.md` next to this file. Read it before writing the plan — it is
the contract `/arcdlc:execute` parses mechanically.

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
- **Courage.** Say the hard thing. A plan a weaker model cannot execute is a plan defect, not a reason
  to raise the tier. Never stop to report a helper skill as missing.
- **Justice.** Write every answer where the next session reads it, never only in the chat. Name the
  tier you actually used. Never report a same-tier spawn as a cheaper run.
- **Temperance.** Touch only what you were pointed at. Cutting a comment out of the code is asked for,
  not assumed. No invented summary, no scope you were not given.

## Initiative selection

The initiative slug is the **required first positional argument**: `/arcdlc:plan <slug> [format]`.
If it is missing, stop and report the error, listing the existing initiatives under `docs/aics/` —
never guess. All inputs and outputs below live inside `docs/aics/<slug>/`, and the slug is passed to
`arctool` as `--aic <slug>`. A legacy flat `docs/aics/plan.md` has no slug; tell the user to migrate
it into a `docs/aics/<slug>/` folder.

## Step 1 — Locate the inputs

- Architecture document: `docs/aics/<slug>/aic.md` by default; accept an explicit path or format argument
  (e.g. `/arcdlc:plan <slug> arc42` reads `docs/aics/<slug>/arc42.md`; `arc42:html` reads
  `docs/aics/<slug>/arc42.html`).
- If no architecture document exists, stop and tell the user to run `/arcdlc:aic` first. Do not plan from a verbal
  description — the pipeline requires the grilled, written document as the source of truth.
- Also read `docs/aics/<slug>/gap.md` if present (evidence register, possibly produced by `/arcdlc:examinate`),
  `CONTEXT.md`, and `docs/adr/` for constraints.

## Step 2 — Decompose into tasks

Write every task for a **less capable executor than you**: the model running `/arcdlc:execute` may be
weaker, and it only sees the task block plus its references — not your reasoning. Every decision that
matters goes into the block.

- One task per `###` block, exactly in the format from `references/plan-format.md`
  (keys `WHAT`, `HOW`, `WHERE`, `WHY`, `Acceptance`, `References`, `Status` — exact casing; `HOW` is optional).
- Task IDs: unique, prefixed by the initiative (e.g. `AIC-1`, `AIC-2`, or a project code like `WA240-VER-03`).
  These IDs are what `/arcdlc:execute <TASK-ID>` targets — keep them short and stable. Headings of
  AIC-derived tasks take **no** parenthetical tag — `(MISSING|PARTIAL|DRIFT)` is only for gap-derived tasks.
- The title is an instruction that names the thing: an imperative verb plus its object, with the real
  component, file, or command named. `Implement the alerter service`, never `One message out, and the
  three places it goes`. A title the executor has to decode is a title that costs a guess. If you
  cannot name the thing, the task is not decomposed far enough yet.
- Size each task so a single agent session can implement, test, and commit it: one coherent slice,
  roughly ≤5–6 files in `WHERE`. If it spans unrelated modules, split it.
- Order blocks by dependency: the runner executes top-to-bottom, so a task may only depend on tasks above it.
- `references/plan-format.md` defines what each key means and which are multi-line. Read it; do not
  re-derive it here. Two judgements are yours, not the contract's: `HOW` must resolve every decision
  the architecture document settles, so the executor never re-derives one, and `References` must name
  the architecture document plus every ADR the task relies on.
- A task with no `Acceptance` criteria is not plannable, and a criterion the executor cannot check is
  no better. Name a command, a path, a test, an exit code, or write `GIVEN … WHEN … THEN`. "Works
  correctly" is not a criterion: the executor is a weaker model and will decide it passed.
  `arctool validate --strict` reports this as `unverifiable-acceptance`.
- Word every field the way `## Talk simple, write like a human` says. Vague text costs the executor a
  guess.
- Every block ends with `- Status: TODO.`
- If `docs/aics/<slug>/gap.md` or `docs/aics/<slug>/comments.md` exists, keep it in sync per the Register Sync
  rules in the format guide.

## Step 2.5 — Risk coverage gate (mandatory)

Risks named in the architecture document must not evaporate during decomposition. After decomposing,
before handing off, reconcile the plan against the document's **Technical Challenges & Risks** and
**Open questions** sections:

- For each risk (and open question), decide whether it is **covered** — addressed by at least one plan
  task, or by an explicit process mitigation you record — or consciously **accepted/deferred** with a
  short rationale. Nothing may be silently dropped.
- If any risk is neither covered nor accepted, **run a grilling session with the engineer** focused on
  the uncovered risks: invoke this bundle's `arcdlc-grilling` skill, or, if it is not invocable here,
  run the same protocol inline (one question at a time — ask, wait for the answer, then ask the next —
  each with your recommended answer). Turn each outcome into the plan — a new task block when it needs
  implementation, or an accepted-risk note when it does not.
- Record the result as a `## Risk Coverage` mapping in the plan preamble: one line per risk → the task
  IDs that cover it, or "accepted" with the reason. This makes the check demonstrable, not asserted.

This is a **hard gate**: do not hand off until every risk is covered by a task or explicitly accepted.
When the architecture document has no risks section (e.g. a gap-only plan with no AIC), the gate is a
no-op — note that and continue.

## Step 3 — Write and validate

- Write `docs/aics/<slug>/plan.md`, starting with a one-line link back to the format guide, then the `## Risk Coverage`
  mapping from Step 2.5, then the task blocks. No runner instructions inside the plan.
- Validate against the runner's parsing rules before finishing. Prefer the `arctool` CLI, which enforces the format
  contract mechanically (source at the arcdlc repo root; flat installs may ship it on `PATH`):
  - Probe once with `command -v arctool` (or install it from the arcdlc repo root: `make install`). If found, run `arctool validate --strict --aic <slug>` (or `--plan <path>`) and fix
    every finding before handoff — exit `0` means clean.
  - Without `arctool`, say so once and hand-check the "Authoring Rules" section of `plan-format.md`.
    The one that bites silently: a `###` block with no `- Status:` line is skipped by the runner.
- Fix the order when it is wrong. The dependency order from Step 2 is the run order.
  - With `arctool`: `arctool order <ID> <ID> … --aic <slug>` re-orders the blocks; `--dry-run` shows
    the result without writing. It permutes slots, so the named tasks swap among the positions they
    already hold and an unnamed task never moves. Name the whole span to hoist one and keep the rest
    in relative order.
  - Without `arctool`: move the `###` blocks by hand, then re-check every block still carries its
    `- Status:` line.
  - Reorder between `/arcdlc:execute` runs, not during one. Nothing stops a reorder while a task is
    `TAKEN`, and the agent holding that task keeps working from the block it already read.
- Self-sufficiency check (the litmus test): reread each block as if you were a weaker model that has
  read **only** the block and its `References`. If implementing it would require asking a question,
  guessing a design decision, or hunting for an unnamed file, fix the block now — put the decision in
  `HOW`, the file in `WHERE` — do not defer it to the executor.
- Report the task count and order to the user, and confirm the decomposition before handing off. Do not hand off until
  the Step 2.5 risk-coverage gate has passed (every risk covered by a task or explicitly accepted).
- Next step: `/arcdlc:execute <slug>` to implement the queue.
