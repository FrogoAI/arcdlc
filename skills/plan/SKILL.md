---
name: arcdlc-plan
description: Decompose an approved architecture document (AIC by default; arc42, TOGAF) into the executable docs/aics/<slug>/plan.md task queue consumed by /arcdlc:execute. The initiative slug is the required first argument (e.g. /arcdlc:plan payments). Use when the user runs /arcdlc:plan, invokes arcdlc-plan, or asks to turn an architecture document into an implementation plan.
argument-hint: "<slug> [aic|arc42|togaf|path]"
---

# ArcDLC Plan (/arcdlc:plan)

Turn an approved architecture document into `docs/aics/<slug>/plan.md` — the executable task queue of the ArcDLC
delivery pipeline:

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`

The plan format is defined in `references/plan-format.md` next to this file. Read it before writing the plan — it is
the contract `/arcdlc:execute` parses mechanically.

## Initiative selection

The initiative slug is the **required first positional argument**: `/arcdlc:plan <slug> [format]`.
If it is missing, stop and report the error, listing the existing initiatives under `docs/aics/` —
never guess. All inputs and outputs below live inside `docs/aics/<slug>/`, and the slug is passed to
`arctool` as `--aic <slug>`. A legacy flat `docs/aics/plan.md` has no slug; tell the user to migrate
it into a `docs/aics/<slug>/` folder.

## Step 1 — Locate the inputs

- Architecture document: `docs/aics/<slug>/aic.md` by default; accept an explicit path or format argument
  (e.g. `/arcdlc:plan <slug> arc42` reads `docs/aics/<slug>/arc42.md`).
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
- Size each task so a single agent session can implement, test, and commit it: one coherent slice,
  roughly ≤5–6 files in `WHERE`. If it spans unrelated modules, split it.
- Order blocks by dependency: the runner executes top-to-bottom, so a task may only depend on tasks above it.
- `HOW` records the design decisions the executor must not re-derive or guess: signatures, naming,
  data shapes, algorithm choice, edge cases, error handling — resolved by you from the architecture
  document. Name the relevant section of a long reference here (e.g. `see aic.md §"Data model"`).
  When adjacent work must not be touched, fence it: `Out of scope: <thing> (covered by <TASK-ID>).`
- `WHERE` lists the exact files/modules expected to change, per layer of the target project.
- `Acceptance` gives at least one testable success criterion — the task's definition of done that
  `/arcdlc:execute` must demonstrate before marking it `DONE`. Prefer `GIVEN … WHEN … THEN …`
  scenarios; each criterion must be confirmable by a test or observable behavior, not a paraphrase of
  `WHAT`. Where a criterion is test-verifiable, name the runnable check (test file or command, e.g.
  ``go test ./internal/x/...``). This is the contract's teeth: a task with no acceptance criteria is
  not plannable.
- `References` must include the architecture document and any ADRs the task relies on — clean file
  paths only (section pointers go in `HOW`).
- Every block ends with `- Status: TODO.`
- If `docs/aics/<slug>/gap.md` exists, keep it in sync per the Gap Register Sync rules in the format guide.

## Step 2.5 — Risk coverage gate (mandatory)

Risks named in the architecture document must not evaporate during decomposition. After decomposing,
before handing off, reconcile the plan against the document's **Technical Challenges & Risks** and
**Open questions** sections:

- For each risk (and open question), decide whether it is **covered** — addressed by at least one plan
  task, or by an explicit process mitigation you record — or consciously **accepted/deferred** with a
  short rationale. Nothing may be silently dropped.
- If any risk is neither covered nor accepted, **run a grilling session with the engineer** focused on
  the uncovered risks: invoke the `grilling` skill if available, otherwise interview inline (one
  question at a time, each with your recommended answer). Turn each outcome into the plan — a new task
  block when it needs implementation, or an accepted-risk note when it does not.
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
  - Probe once with `command -v arctool` (or install it from the arcdlc repo root via
    `make build`/`make release`). If found, run `arctool validate --strict --aic <slug>` (or `--plan <path>`) and fix
    every finding before handoff — exit `0` means clean.
  - If `arctool` is unavailable, say so once and hand-check the same rules from `plan-format.md`:
    - Every `###` block contains a `- Status:` line (missing status = silently skipped by the runner).
    - Status values are uppercase with optional trailing period.
    - Task IDs are unique.
    - Every block has a non-empty `- Acceptance:` section (`--strict` fails otherwise — it implies
      `--require-acceptance`).
- Self-sufficiency check (the litmus test): reread each block as if you were a weaker model that has
  read **only** the block and its `References`. If implementing it would require asking a question,
  guessing a design decision, or hunting for an unnamed file, fix the block now — put the decision in
  `HOW`, the file in `WHERE` — do not defer it to the executor.
- Report the task count and order to the user, and confirm the decomposition before handing off. Do not hand off until
  the Step 2.5 risk-coverage gate has passed (every risk covered by a task or explicitly accepted).
- Next step: `/arcdlc:execute <slug>` to implement the queue.
