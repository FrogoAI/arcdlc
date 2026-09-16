---
name: arcdlc-execute
description: Implement the planned tasks, one at a time, each with tests, lint and its own commit, flipping status TODO to TAKEN to DONE. Use when someone says start implementing, build it, work the plan, do the next task, pick up the backlog, ship this, or get it done, or runs /arcdlc:execute <slug> [TASK-ID], or invokes arcdlc-execute.
argument-hint: "<slug> [TASK-ID]"
---

# ArcDLC Execute (/arcdlc:execute)

Implement tasks from `docs/aics/<slug>/plan.md`. You are the executor: the plan format and status lifecycle defined in
`../plan/references/plan-format.md` (flat installs: `../arcdlc-plan/references/plan-format.md`) are the contract you
enforce — read that file before starting.

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

## Before the first task: is the plan still current?

Run `arctool validate --aic <slug>` once at the start. A `source-changed` warning means the
architecture document moved after this plan was written, so some tasks may build a design that no
longer exists. You cannot tell which: the blocks are mechanical and carry no trace of the reasoning.

Stop and tell the engineer. They either re-run `/arcdlc:plan <slug>`, or read the diff and re-stamp
with `arctool stamp --aic <slug>` to say the tasks still hold. Never make that call yourself, and
never execute past the warning. *Fallback without `arctool`: check the plan's
`<!-- arcdlc:source … -->` line against `sha256sum` of the file it names.*

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

Give every task a fresh context. The plan carries all the state: statuses in `plan.md`, work in per-task
commits, so nothing needs to survive in conversation between tasks.

This is the point of the whole shape, and it matters more the longer the queue. One agent working 150
tasks in one session accumulates every file it read and every diff it produced, compacts somewhere in
the middle, and a squeezed executor drops acceptance criteria first. So quality falls exactly where the
risk is highest: late in the queue, where the most work already rests on what came before. One hundred
and fifty fresh subagents do not degrade. Task 150 is executed as well as task 1, which is a property no
choice of model can buy you.

The commits are the same bargain. One hundred and fifty small commits are reviewable, revertable and
bisectable, and each is a resume point. One session producing one enormous commit is none of those.

Probe once: can this harness spawn subagents with their own clean context (e.g. the Agent/Task tool in Claude Code)?

### Executor tier: ask, never assume

A task block is written to be mechanical: executable by a model with no context but the block and the
files it names. That is a property of the plan, not a bet on which model is cheap this year. What it
buys you is freedom to run the queue on a cheaper tier, and a tier is two dials, not one: the model,
and the effort level it runs at. A strong model at its lowest effort often fits a block that already
carries every decision, and usually beats a weaker model at the same price.

You pick neither dial on your own. Both belong to the engineer, and you ask once per run.

The floor is capability, not price. Whatever tier the run lands on, the subagent must be able to:

- run shell commands,
- read and edit files in the repository,
- run the project's test and lint commands (`make test` and `make lint`, or the ones the project documents),
- write a git commit.

Resolve the tier in this order, and stop at the first line that answers:

1. **The task names it.** A task's `HOW` names a model or an effort level. The planner decided it, so it binds that
   task, exactly like every other `HOW` decision.
2. **The project pin.** A line in `CONTEXT.md` that starts with `Executor tier:`, for example
   `Executor tier: haiku` or `Executor tier: opus, low effort`. Use it exactly as written and do not ask. The
   engineer wrote it for the harness they actually run. (`arctool sync` rewrites `AGENTS.md` and `README.md` only,
   so a pin in `CONTEXT.md` survives every sync.)
3. **Ask the engineer.** One question, one turn, before the first spawn: which model and which effort level should
   run this queue. Name the options this harness actually exposes, and carry a recommended answer, which is the
   cheapest combination that clears the floor. Then use that answer for every task in the run and for the
   verification subagent, and offer to write it into `CONTEXT.md` as the pin above so the next run does not ask
   again.
4. **Nobody to ask, or nothing to choose.** A non-interactive run (`claude -p ...`), or a harness that cannot vary
   the model or the effort of a subagent. Then use in-session mode below and say so in the report. Never spawn a
   subagent at your own tier and report it as a cheaper run: a fresh same-tier subagent buys clean context, not a
   cheaper executor.

Ask only when you are about to spawn. Single-task mode (`/arcdlc:execute <slug> <TASK-ID>`) runs in the current
session, so there is no tier to choose and no question to ask.

The tier is not a quality dial. When the executor meets a decision its block does not carry, it grills the
engineer or blocks the task (per-task contract, step 7). Never retry at a higher tier. A block that is not
mechanical is a plan defect, and raising the tier hides it: the task gets done by judgement nobody
recorded, and the next run of that plan behaves differently. `/arcdlc:plan <slug>` sharpens the block
instead. Name the tier the run used, and where it came from, in the report.

**Orchestrator mode (subagents available).** Run the queue as a thin dispatcher and implement nothing yourself:

1. Get the next task ID: `arctool next --json` (fallback: the first `TODO` block in `plan.md`). None left → go to the
   Verification phase.
2. Spawn ONE fresh subagent, never several in parallel. This is a decision, not a limitation: the plan is
   an ordered list, not a dependency graph, so it runs top to bottom one task at a time. That keeps the
   history linear and the run resumable from a single pointer, and it is why no task declares what it
   depends on. Where throughput matters, run separate initiatives, which are independent by construction. Spawn it at the tier resolved above, never at one you picked yourself. Its prompt must name the
   initiative slug, the task ID, the per-task contract to follow (point it at this skill file and
   `../plan/references/plan-format.md`; flat installs: `../arcdlc-plan/references/plan-format.md`), and the
   accumulated notes from earlier task reports.
3. **Tell the subagent it owns its own verification, and has nobody to ask.** Its prompt must say both
   plainly:
   - It verifies its own work. It reaches `DONE` only when every `Acceptance` criterion demonstrably
     holds and tests and lint are green (step 5 of the contract). Nobody downstream re-checks a `DONE`
     task, so `DONE` is its assertion and it must mean it.
   - There is no engineer in its session. It must never grill and never guess in place of grilling. Its
     task is mechanical; if it turns out not to be, that is not its problem to solve. On any ambiguity,
     contradiction, missing decision, or criterion it cannot satisfy, it runs
     `arctool block <id> -m "<the question or the defect>"` and returns it to you, unanswered and
     unguessed. **Returning a blocked task is a successful outcome for a subagent, not a failure**, and
     saying so in the prompt is what stops a cheaper model guessing to look helpful.

   It then executes the per-task contract below (take → implement → verify acceptance → done → commit) and
   reports back its final status, the question or defect if it blocked, and at most one line of notes
   useful to later tasks. Not the diff, not the test output: you do not need them.
4. **A `DONE` task is done. Move on.** The subagent verified its own acceptance criteria in step 5 of the
   per-task contract and could not reach `DONE` without them holding. Re-running the same criteria here
   would catch only a subagent that lied about running them, and would cost you that output on every task
   in the queue. You are not the per-task check. Take the one line of notes it returned and spawn the next.
5. **A `BLOCKED` task is the only thing that needs you, and it needs the thing only you have.** The subagent
   had no engineer in its session; you do. Read the block reason, then:
   - **It returned a question**: grill for it, per "When a task is unclear" below, one question per turn.
     Then re-spawn the same task with the answer in the prompt, and note in the report that the block was
     not mechanical, so `/arcdlc:plan` can harden it.
   - **The answer changes the plan**: leave it `BLOCKED`, stop the run, hand back to `/arcdlc:plan <slug>`.
     You never edit `plan.md`.
   - **It hit a defect rather than a question** (a criterion that cannot be satisfied, a contradiction in the
     code): stop the run and report. Never re-spawn at a higher tier to get past it: that is the
     tier-as-quality-dial mistake, and it buries a defect the plan should carry.
   - **Nobody to ask** (a non-interactive run): leave it `BLOCKED`, stop, and report the question verbatim so
     it is the first thing the engineer sees. Never answer it yourself to keep the queue moving.

   Contradictions between tasks are not visible here and you should not go looking: a task that is fine on
   its own can still conflict with one from forty tasks ago. That is a whole-queue property, and the
   Verification phase is where it is checked.

Keep your own context small. In this loop you read plan state, the one line of notes each subagent
returns, and the block reason when one comes back `BLOCKED`. Nothing else: not reports of successful
work, not diffs, not test output. A dispatcher that stays this thin can drive a long queue in one
invocation, which is the point of the mode.

If it does fill up anyway, on a very long queue or after several blocked tasks pulled you into detail,
stop at a task boundary and tell the engineer to clear the session and re-run `/arcdlc:execute <slug>`,
saying how many tasks are done and how many remain. Nothing is lost: the plan carries every status and
the commits carry every change, so the next invocation resumes at `arctool next` exactly where this one
ended. A 150-task queue is any number of sessions over an unchanged plan, not one heroic run.

**In-session mode (no subagents, no tier to choose, or nobody to ask: flat installs and other harnesses).**
Execute tasks yourself, one at a time, with a hard boundary discipline. You cannot measure your own context size,
so use proxies:

- Task boundaries (after `done` + commit) are the only legitimate stopping points.
- After each non-trivial task — or roughly every third small one, or immediately when the harness signals compaction
  or low context — finish the current task, commit, then stop and tell the user: clear the session and re-run
  `/arcdlc:execute <slug>`. Nothing is lost; the run resumes from `plan.md` exactly where it stopped.
- Never start a new task in a nearly-exhausted context, and never let compaction land mid-task.

A queue can also be driven externally — one non-interactive run per task (e.g. `claude -p "/arcdlc:execute <slug>"`
in a loop) until `arctool next` exits `3` (no `TODO` left).

Single-task mode (`/arcdlc:execute <slug> <TASK-ID>`) needs none of this: execute it directly in the current session.

## Per-task contract

Work trunk-based. Branch off trunk for the task, keep the branch short-lived, commit once per task, and
merge back the same day. Ship incomplete work behind a flag rather than on a branch that lives for
days. Never open a long-lived release or feature branch, and never batch several tasks into one branch:
the plan's one-task-one-commit rule is what keeps trunk releasable.

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
   If the task is unclear, contradicts something, or would need a decision it does not carry, do not guess:
   grill the engineer first, per "When a task is unclear" below.
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
7. If the task cannot be completed — including any acceptance criterion you cannot satisfy: when the blocker is
   a question rather than a defect, grill the engineer if you are the session they started, and if you are a
   spawned subagent put the question in the block reason and return it unanswered (see below). Then `arctool block <id> -m
   "<one-line reason>"` naming the failing criterion (or `arctool todo <id>` to release it back to the queue), report
   why, and stop — do not continue to the next task on failure. Never `arctool done` a task whose acceptance criteria
   are unmet. The same rule covers an order inversion: while implementing, you find the task needs something a task
   **below** it will build. Block it the same way, with a reason that names the other task: `arctool block <id> -m
   "needs <OTHER-ID>, which is below it"`. Then add one line to your report with the exact command that would fix the
   order, for example `arctool order <OTHER-ID> <id> --aic <slug>`, which swaps those two positions and leaves every
   other block where it is. You never run that command: the engineer decides whether the plan changes, and the run
   stops here either way. *Fallback: set the status line to `- Status: BLOCKED — <reason>.` or back to `TODO`.*
8. Repeat from step 1. When running the whole queue, stop when `arctool next` exits non-zero (code `3` = no `TODO`
   left). *Fallback: stop when no `TODO` block remains.*

## When a task is unclear, never guess (mandatory)

**Who answers depends on who is in the room, and a spawned subagent is alone in it.**

- **You are a spawned subagent** (orchestrator mode): there is no engineer in your session. Do not grill,
  do not wait, and above all do not guess in its place. Run `arctool block <id> -m "<the question>"`, state
  the question in your report, and stop. Returning an unanswered question is a correct outcome for you.
  Someone with an engineer will answer it.
- **You are the session the engineer started** (in-session mode, or the dispatcher handling a subagent's
  returned question): you have someone to ask, so ask. Grill, one question per turn, and record the answer
  where the next session reads it.

Getting this backwards is how a queue goes quietly wrong: a subagent that guesses produces a task that
looks `DONE`, and a dispatcher that does not ask leaves the same gap for the next run.

A task the plan left ambiguous is not yours to guess. A mechanical plan should make this rare, so
every time it happens one of two things is true, and they are handled differently:

- **The block was defective**: it was not mechanical, and the planner should have settled this. Grill
  for the answer so the run continues, then **name it in your report as a plan defect** with the task
  ID and the decision that was missing. That is how `/arcdlc:plan <slug>` hardens the block instead of
  the same gap biting the next run. Do not skip the report because you got an answer.
- **Reality contradicts the plan**: the code, a dependency, or a test says something the planner could
  not have known. Grill for the decision the same way. This is not a plan defect, and saying so in the
  report keeps the signal honest.

Either way, stop coding and ask the engineer the moment you hit one of these:

- The task contradicts itself, the code, an ADR, `CONTEXT.md`, or the architecture document.
- `HOW` is missing a decision you would otherwise invent: a signature, a name, a data shape, a
  threshold, an error path, an edge case.
- An `Acceptance` criterion cannot be satisfied as written, or two criteria pull opposite ways.
- The change is risky or hard to reverse and the task does not say it is wanted: a data migration, a
  deleted public interface, a dropped column, a rewritten configuration format, a dependency swap.
- Two tasks in the queue want the opposite thing.
- The task tells you to change code you can find no caller for, and the answer decides whether it is
  dead or called from outside this repository.

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

Where the answer goes, so nobody has to answer it twice:

- A decision that outlives this task: write an ADR in `docs/adr/`, name it in the commit body, and
  cite it in the code comment that needs it.
- A decision local to this task: put it in the commit body, under the task line.
- A new term the project will reuse: add it to `CONTEXT.md`.
- Do not edit another task, and do not rewrite this task's text in `plan.md`. Text edits race with the
  status writes of a parallel run, and the plan belongs to `/arcdlc:plan`.
- No answer, or an answer that changes the plan: `arctool block <id> -m "<one-line reason>"`, report
  what you asked, and stop. The engineer decides whether the plan changes.
- Every grilled decision, answered or not, goes in the run report: the task ID, the question, the
  answer, and whether it was a plan defect or something reality forced. A run that grilled three times
  and reported none looks like a clean run, and the plan never gets fixed.

## Commit message: Conventional Commits

Every commit this skill makes, per-task commits and verification-phase fixes alike, follows the
Conventional Commits 1.0.0 specification. The rules below are the whole of it that ArcDLC uses.

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
- **scope** — always the initiative slug, including in single-initiative repos: `feat(payments-v2): …`.
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

**This is the real gate, and it is yours.** Each task verified itself in isolation, which is the only
thing a task can verify: its own criteria, against the block it was given. Nothing so far has asked
whether the tasks agree with each other. Two of them can each be correct and still contradict, and the
queue is exactly where that hides, because tasks that never touch are never compared.

Run it whenever the whole queue was worked (no task-ID argument), in three parts.

**1. The queue really is finished.** Cheap, and it catches a crashed session that a passing test suite
never would:

- `arctool list --status TAKEN` and `--status BLOCKED` are both empty. A leftover `TAKEN` is a subagent
  that died mid-task; its work is half-applied and no commit exists. Reset it with `arctool todo <id>` and
  report, do not quietly finish it yourself.
- Every `DONE` task has a commit. `git log --oneline` should show one per task, and the count should match.
- `arctool validate --strict --aic <slug>` exits 0. This also re-checks the source stamp, so a design that
  moved during a long run surfaces here rather than never.

**2. The project is whole.** Mechanical, so delegate it to one final subagent at the task tier, and take
back only the verdict:

- Build, test and lint at the repository level, not only the subproject the tasks touched:
  the project's documented commands, or `make test` and `make lint`, skipping targets that do not exist.
- Fix failures, commit each separately per the Conventional Commits rules above, and re-run until clean, or
  until the same failure repeats without progress, then stop and report.

**3. The tasks agree with each other.** Judgement, so this part stays with you. Do not delegate it, and do
not skip it because part 2 was green: a passing test suite is exactly what a pair of contradicting tasks
looks like when neither has a test for the other's assumption.

- List the tasks whose `WHERE` overlapped. Those are where one task can have undone or redefined another's
  work. Read the final state of the shared files, not each commit.
- Check the decisions. Two tasks that each recorded a decision in a commit body, about the same interface,
  the same name, the same data shape, may have recorded opposite ones.
- Check the seams between independent tasks: a producer and a consumer planned separately, a writer and a
  reader of the same file or table, two tasks that each added a case to the same switch.
- A contradiction found here is a plan defect, not an execution one. Report it with both task IDs and hand
  back to `/arcdlc:plan <slug>`. Do not paper over it with a fix-up commit: the plan produced two tasks
  that disagree, and the next run of it will do so again.

## Report

Summarize per task: what changed, validation results, and the commit. List every decision you grilled
for, marked as a plan defect or as something reality forced, so the planner can harden the blocks that
were not mechanical. Name the executor tier the run used and where
it came from: a task's `HOW`, the `CONTEXT.md` pin, the engineer's answer this run, or in-session because there was
no tier to choose or nobody to ask. Suggest `/arcdlc:archive <slug>` when several `DONE` blocks have accumulated in
the plan.
