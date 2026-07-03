---
description: Implement tasks from docs/aics/<slug>/plan.md following the plan status contract — one task at a time, status TODO→TAKEN→DONE, tests/lint, one commit per task. Use when the user runs /arcdlc:execute (all pending tasks) or /arcdlc:execute <TASK-ID> (single task), or invokes arcdlc-execute.
argument-hint: "[--aic <slug>] [TASK-ID]"
---

# ArcDLC Execute (/arcdlc:execute)

Implement tasks from `docs/aics/<slug>/plan.md`. You are the executor: the plan format and status lifecycle defined in
`../plan/references/plan-format.md` (flat installs: `../arcdlc-plan/references/plan-format.md`) are the contract you
enforce — read that file before starting.

## Tooling: prefer `arctool`

Probe once at the start: `command -v arctool`. If present, drive the queue with the `arctool` commands below — they read
one block at a time and mutate status through guarded, byte-preserving, atomic writes, so you never re-read the whole
plan or hand-edit a status line. If `arctool` is absent, say so once (`arctool not found — operating on plan.md by hand`)
and use the manual fallback noted in each step. Either way `plan.md` stays the single source of truth.

Pass the resolved initiative to every `arctool` call as `--aic <slug>` (or `--plan <path>`); when you omit it,
`arctool` auto-detects the single initiative under `docs/aics/`.

## Initiative selection

Each initiative has its own `docs/aics/<slug>/plan.md`. Resolve which one to work:

- `--aic <slug>` selects it explicitly (put it before the `TASK-ID`, e.g. `/arcdlc:execute --aic payments-v2 AIC-1`).
- With no `--aic`, auto-detect the single initiative under `docs/aics/`. If several exist, **stop and ask which**
  (`arctool` exits `2` and lists the slugs) — never guess.

**One initiative per run.** A single `/arcdlc:execute` works exactly one `plan.md` to keep the "one plan is the source
of truth, one commit per task" discipline. To work another initiative, run again with its `--aic <slug>`.

## Argument: task selection

- `/arcdlc:execute` — execute every pending task in the resolved plan, top-to-bottom, one at a time.
- `/arcdlc:execute <TASK-ID>` — execute only that task (e.g. `/arcdlc:execute AIC-1`); read it with
  `arctool show <TASK-ID> --json` (add `--aic <slug>` when needed). Task IDs are unique **within** a plan, so pair an
  ID with `--aic <slug>` when more than one initiative exists.
  - If its status is not `TODO`, stop and report the status; only proceed on `DONE`/`BLOCKED`/`TAKEN` if the user
    explicitly confirms a redo or takeover (`arctool take` refuses a non-`TODO` task unless you pass `--force`).

## Per-task contract

For each task, in order:

1. Get the task: `arctool next --json` (whole queue) or `arctool show <TASK-ID> --json` (single task). Read only the
   files named in its `references` and `where`/`whereLayers` — not the whole plan. Note its `acceptance` criteria:
   they are the definition of done you must satisfy in step 5. *Fallback: read `plan.md` and take the first `### `
   block whose `- Status:` is `TODO`, including its `- Acceptance:` section.*
2. Claim it before touching code: `arctool take <id>` (flips `TODO`→`TAKEN`; refuses a non-`TODO` task). A `TAKEN` block
   with no commit marks a crashed session. *Fallback: edit the block's `- Status: TODO.` to `- Status: TAKEN.`*
3. Implement ONLY this task, exactly as written — including intentional breaking changes when the task says so.
   The whole repository is context; changes go in the files/modules named in `WHERE` (extend within the same
   subproject when strictly needed to complete the task).
4. Run the relevant tests and lint for the touched areas. If the subproject `Makefile` has `test`/`lint` targets, use
   `make test` and `make lint`; otherwise use the project's documented commands.
5. Verify acceptance, then mark done. Walk **every** `acceptance` criterion from step 1 and confirm each is
   demonstrably met — by the test that exercises it or by the observable behavior it describes. Only when tests/lint
   are green **and** all criteria hold: `arctool done <id>` (flips `TAKEN`→`DONE`; touches no other task's status).
   If a criterion is not covered by an existing test, add one (in the `Tests` files named in `WHERE`) so "met" is
   evidenced, not asserted. *Fallback: edit the status line to `- Status: DONE.`*
6. Commit ONLY this task's changes plus the plan status update. Do not include unrelated pre-existing worktree
   changes. The commit message must include `#AI-assisted` and a concise summary. When more than one initiative
   exists, prefix the subject with the slug (e.g. `[payments-v2] AIC-1: add health endpoint`). Do not push.
7. If the task cannot be completed — including any acceptance criterion you cannot satisfy: `arctool block <id> -m
   "<one-line reason>"` naming the failing criterion (or `arctool todo <id>` to release it back to the queue), report
   why, and stop — do not continue to the next task on failure. Never `arctool done` a task whose acceptance criteria
   are unmet. *Fallback: set the status line to `- Status: BLOCKED — <reason>.` or back to `TODO`.*
8. Repeat from step 1. When running the whole queue, stop when `arctool next` exits non-zero (code `3` = no `TODO`
   left). *Fallback: stop when no `TODO` block remains.*

## Verification phase (after the queue is empty)

When running the full queue (no task-ID argument), finish with a whole-project check:

- Run `make test` and `make lint` in the subproject (skip targets that don't exist).
- Fix any failures, commit fixes separately with `#AI-assisted`, and re-run until clean or the same failure repeats
  without progress — then stop and report.

## Report

Summarize per task: what changed, validation results, and the commit. Suggest `/arcdlc:archive` when several `DONE`
blocks have accumulated in the plan.
