---
name: arcdlc-execute
description: Implement tasks from docs/aics/<slug>/plan.md following the plan status contract — one task at a time, status TODO→TAKEN→DONE, tests/lint, one commit per task. The initiative slug is the required first argument. Use when the user runs /arcdlc:execute <slug> (all pending tasks) or /arcdlc:execute <slug> <TASK-ID> (single task), or invokes arcdlc-execute.
argument-hint: "<slug> [TASK-ID]"
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

Pass the resolved initiative to every `arctool` call as `--aic <slug>` (or `--plan <path>`); `arctool`
always requires an explicit selection.

## Initiative selection

The initiative slug is the **required first positional argument**: `/arcdlc:execute <slug> [TASK-ID]`
(e.g. `/arcdlc:execute payments-v2 AIC-1`). If it is missing, stop and report the error, listing the
existing initiatives under `docs/aics/` — never guess. Work that one `docs/aics/<slug>/plan.md`, passing
`--aic <slug>` to every `arctool` call. A legacy flat `docs/aics/plan.md` has no slug; tell the user to
migrate it into a `docs/aics/<slug>/` folder.

**One initiative per run.** A single `/arcdlc:execute` works exactly one `plan.md` to keep the "one plan
is the source of truth, one commit per task" discipline. To work another initiative, run again with its
slug.

## Argument: task selection

- `/arcdlc:execute <slug>` — execute every pending task in that initiative's plan, top-to-bottom, one at a time.
- `/arcdlc:execute <slug> <TASK-ID>` — execute only that task (e.g. `/arcdlc:execute payments-v2 AIC-1`); read it
  with `arctool show <TASK-ID> --aic <slug> --json`. Task IDs are unique **within** a plan.
  - If its status is not `TODO`, stop and report the status; only proceed on `DONE`/`BLOCKED`/`TAKEN` if the user
    explicitly confirms a redo or takeover (`arctool take` refuses a non-`TODO` task unless you pass `--force`).

## Session strategy: fresh context per task (whole-queue mode)

The plan is the only state carrier — statuses live in `plan.md`, work lives in per-task commits — so nothing needs to
survive in conversation context between tasks, and quality drops when it does: a long session accumulates earlier
tasks' file reads and diffs until the harness compacts mid-work, and a squeezed executor loses acceptance criteria
first. Give every task a fresh context.

Probe once: can this harness spawn subagents with their own clean context (e.g. the Agent/Task tool in Claude Code)?

**Orchestrator mode (subagents available).** Run the queue as a thin dispatcher and implement nothing yourself:

1. Get the next task ID: `arctool next --json` (fallback: the first `TODO` block in `plan.md`). None left → go to the
   Verification phase.
2. Spawn ONE fresh subagent — never several in parallel: the queue is dependency-ordered and commits must not
   interleave. Its prompt must name the initiative slug, the task ID, the per-task contract to follow (point it at
   this skill file and `../plan/references/plan-format.md`; flat installs: `../arcdlc-plan/references/plan-format.md`),
   and the accumulated notes from earlier task reports.
3. The subagent executes the full per-task contract below (take → implement → verify acceptance → done → commit) and
   reports back: files changed, test results, commit subject, final status — plus at most one line of notes useful to
   later tasks (e.g. a project convention it discovered).
4. Verify the outcome yourself before looping: the task's status is `DONE` (`arctool show <id>`) and the commit exists
   (`git log -1`). A subagent that reports success without both counts as failed — reset per step 7 of the contract.
5. On `BLOCKED` or failure, stop the whole run and report — same rule as step 7.

Keep your own context small: in orchestrator mode never read source files or diffs — only plan state, subagent
reports, and commit subjects. That is what lets a long queue finish in a single `/arcdlc:execute <slug>` invocation.

**In-session mode (no subagents — flat installs and other harnesses).** Execute tasks yourself, one at a time, with a
hard boundary discipline. You cannot measure your own context size, so use proxies:

- Task boundaries (after `done` + commit) are the only legitimate stopping points.
- After each non-trivial task — or roughly every third small one, or immediately when the harness signals compaction
  or low context — finish the current task, commit, then stop and tell the user: clear the session and re-run
  `/arcdlc:execute <slug>`. Nothing is lost; the run resumes from `plan.md` exactly where it stopped.
- Never start a new task in a nearly-exhausted context, and never let compaction land mid-task.

A queue can also be driven externally — one non-interactive run per task (e.g. `claude -p "/arcdlc:execute <slug>"`
in a loop) until `arctool next` exits `3` (no `TODO` left).

Single-task mode (`/arcdlc:execute <slug> <TASK-ID>`) needs none of this: execute it directly in the current session.

## Per-task contract

For each task, in order (in orchestrator mode, the spawned subagent performs these steps for its one task):

1. Get the task: `arctool next --json` (whole queue) or `arctool show <TASK-ID> --json` (single task). Read only the
   files named in its `references` and `where`/`whereLayers` — not the whole plan. Note its `how` field (when
   present): those are the planner's design decisions — signatures, naming, edge cases, out-of-scope fences — and
   they are binding, not suggestions. Note its `acceptance` criteria: they are the definition of done you must
   satisfy in step 5. *Fallback: read `plan.md` and take the first `### ` block whose `- Status:` is `TODO`,
   including its `- HOW:` and `- Acceptance:` sections.*
2. Claim it before touching code: `arctool take <id>` (flips `TODO`→`TAKEN`; refuses a non-`TODO` task). A `TAKEN` block
   with no commit marks a crashed session. *Fallback: edit the block's `- Status: TODO.` to `- Status: TAKEN.`*
3. Implement ONLY this task, exactly as written — including intentional breaking changes when the task says so.
   Follow the `HOW` decisions when present and leave anything it marks `Out of scope:` untouched, even if you see
   an adjacent improvement. The whole repository is context; changes go in the files/modules named in `WHERE`
   (extend within the same subproject when strictly needed to complete the task).
4. Run the relevant tests and lint for the touched areas. If the subproject `Makefile` has `test`/`lint` targets, use
   `make test` and `make lint`; otherwise use the project's documented commands.
5. Verify acceptance, then mark done. Walk **every** `acceptance` criterion from step 1 and confirm each is
   demonstrably met — by the test that exercises it or by the observable behavior it describes. Only when tests/lint
   are green **and** all criteria hold: `arctool done <id>` (flips `TAKEN`→`DONE`; touches no other task's status).
   If a criterion is not covered by an existing test, add one (in the `Tests` files named in `WHERE`) so "met" is
   evidenced, not asserted. *Fallback: edit the status line to `- Status: DONE.`*
6. Commit ONLY this task's changes plus the plan status update. Do not include unrelated pre-existing worktree
   changes. Write the message exactly as specified in "Commit message: Conventional Commits" below. Do not push.
7. If the task cannot be completed — including any acceptance criterion you cannot satisfy: `arctool block <id> -m
   "<one-line reason>"` naming the failing criterion (or `arctool todo <id>` to release it back to the queue), report
   why, and stop — do not continue to the next task on failure. Never `arctool done` a task whose acceptance criteria
   are unmet. *Fallback: set the status line to `- Status: BLOCKED — <reason>.` or back to `TODO`.*
8. Repeat from step 1. When running the whole queue, stop when `arctool next` exits non-zero (code `3` = no `TODO`
   left). *Fallback: stop when no `TODO` block remains.*

## Commit message: Conventional Commits

Every commit this skill makes — per-task commits and verification-phase fixes alike — follows the Conventional
Commits 1.0.0 specification, bundled at `../source-map/source/Conventional Commits.md` (flat installs:
`../arcdlc-source-map/source/Conventional Commits.md`). Read that reference when a case below is not covered.

The shape:

```
<type>(<slug>): <description>

<body: why the change, plus anything a reviewer needs>

Refs: <TASK-ID>
#AI-assisted
```

Mechanical rules:

- **type** — derived from what the diff actually does, not from the task's wording: `feat` (new behavior), `fix`
  (bug fix), `docs`, `test`, `refactor` (no behavior change), `perf`, `build`, `ci`, `chore`. If a task spans two
  types, the dominant one wins — never split one task across two commits.
- **scope** — always the initiative slug, including in single-initiative repos: `feat(payments-v2): …`. This
  replaces the former `[payments-v2]` subject prefix.
- **description** — imperative mood, lower case, no trailing period; keep the whole subject line ≤ 72 characters.
- **breaking change** — when the task intentionally breaks an interface (the plan says so), append `!` after the
  scope (`feat(payments-v2)!: …`) **and** add a `BREAKING CHANGE: <what breaks and what callers must do>` footer.
- **footers** — `Refs: <TASK-ID>` (the plan task ID, e.g. `AIC-1`), then the literal `#AI-assisted` marker as the
  last line. `#AI-assisted` is required on every commit, `Refs:` on every task commit; the task ID no longer
  belongs in the subject.
- **fix-up commits** (verification phase, or repairing an earlier task) — same shape with the type that fits
  (`fix`, `test`, `ci`, …) and `Refs:` the task whose work they repair; omit `Refs:` when they belong to no single
  task.
- Pass the message with `git commit -m` (or `-F -`) — an editor-based commit strips the `#AI-assisted` line as a
  comment.

Example, per task:

```
feat(payments-v2): add health endpoint

Expose /healthz returning 200 with build info so the load balancer can
drop unhealthy instances.

Refs: AIC-1
#AI-assisted
```

## Verification phase (after the queue is empty)

When running the full queue (no task-ID argument), finish with a whole-project check:

- Run `make test` and `make lint` in the subproject (skip targets that don't exist).
- Fix any failures, commit fixes separately (Conventional Commits, as above), and re-run until clean or the same failure repeats
  without progress — then stop and report.
- In orchestrator mode, delegate this phase to one final subagent (running tests and fixing failures is
  implementation work, and its output does not belong in the dispatcher's context).

## Report

Summarize per task: what changed, validation results, and the commit. Suggest `/arcdlc:archive` when several `DONE`
blocks have accumulated in the plan.
