---
description: Decompose an approved architecture document (AIC by default; arc42, TOGAF) into the executable docs/aics/<slug>/plan.md task queue consumed by /arcdlc:execute. Use when the user runs /arcdlc:plan, invokes arcdlc-plan, or asks to turn an architecture document into an implementation plan.
argument-hint: "[aic|arc42|togaf|path] [--aic <slug>]"
---

# ArcDLC Plan (/arcdlc:plan)

Turn an approved architecture document into `docs/aics/<slug>/plan.md` — the executable task queue of the ArcDLC
delivery pipeline:

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`

The plan format is defined in `references/plan-format.md` next to this file. Read it before writing the plan — it is
the contract `/arcdlc:execute` parses mechanically.

## Initiative selection

Each initiative lives in `docs/aics/<slug>/`. Resolve the slug the same way the whole pipeline does:
if the user passes `--aic <slug>`, use it; otherwise auto-detect the single initiative folder under
`docs/aics/`. If several exist and no `--aic` is given, list them and ask which one. All inputs and
outputs below are inside that folder.

## Step 1 — Locate the inputs

- Architecture document: `docs/aics/<slug>/aic.md` by default; accept an explicit path or format argument
  (e.g. `/arcdlc:plan arc42` reads `docs/aics/<slug>/arc42.md`).
- If no architecture document exists, stop and tell the user to run `/arcdlc:aic` first. Do not plan from a verbal
  description — the pipeline requires the grilled, written document as the source of truth.
- Also read `docs/aics/<slug>/gap.md` if present (evidence register, possibly produced by `/arcdlc:examinate`),
  `CONTEXT.md`, and `docs/adr/` for constraints.

## Step 2 — Decompose into tasks

- One task per `###` block, exactly in the format from `references/plan-format.md`
  (keys `WHAT`, `WHERE`, `WHY`, `References`, `Status` — exact casing).
- Task IDs: unique, prefixed by the initiative (e.g. `AIC-1`, `AIC-2`, or a project code like `WA240-VER-03`).
  These IDs are what `/arcdlc:execute <TASK-ID>` targets — keep them short and stable.
- Size each task so a single agent session can implement, test, and commit it. Split anything larger.
- Order blocks by dependency: the runner executes top-to-bottom, so a task may only depend on tasks above it.
- `WHERE` lists the exact files/modules expected to change, per layer of the target project.
- `Acceptance` gives at least one testable success criterion — the task's definition of done that
  `/arcdlc:execute` must demonstrate before marking it `DONE`. Prefer `GIVEN … WHEN … THEN …`
  scenarios; each criterion must be confirmable by a test or observable behavior, not a paraphrase of
  `WHAT`. This is the contract's teeth: a task with no acceptance criteria is not plannable.
- `References` must include the architecture document and any ADRs the task relies on.
- Every block ends with `- Status: TODO.`
- If `docs/aics/<slug>/gap.md` exists, keep it in sync per the Gap Register Sync rules in the format guide.

## Step 3 — Write and validate

- Write `docs/aics/<slug>/plan.md`, starting with a one-line link back to the format guide, then the task blocks. No
  runner instructions inside the plan.
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
- Report the task count and order to the user, and confirm the decomposition before handing off.
- Next step: `/arcdlc:execute` to implement the queue.
