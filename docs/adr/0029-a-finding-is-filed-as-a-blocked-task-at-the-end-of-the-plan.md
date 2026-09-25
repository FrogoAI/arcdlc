# ADR-0029 — A finding is filed as a blocked task at the end of the plan

- Status: Accepted
- Date: 2026-09-25
- Initiative: [findings](../aics/findings/aic.md)
- Amends: [ADR-0023](0023-the-dispatcher-verifies-and-the-subagent-never-guesses.md). Verification
  part 1 no longer requires zero `BLOCKED` tasks: it ignores findings, and a review of them becomes
  the run's last step. Also narrows [ADR-0021](0021-a-planned-task-must-be-mechanical.md): a task an
  executor files is not yet mechanical, so it enters the plan `BLOCKED` and runs only after a person
  releases it.

## Context

An `/arcdlc:execute` executor reads code for its task and sometimes sees a problem its task does not
cover: a bug in a neighbouring module, in code a closed initiative delivered, or in something `HOW`
fences off. Step 3 of the per-task contract tells it to leave that alone, which is right, and the only
place left to mention it is the one line of notes it returns to the dispatcher. That line lives in the
dispatcher's chat and is gone when the session ends, so the problem is lost until it breaks something.

## Decision

1. **A finding** is a defect, risk or inconsistency the executor sees in a file it read, that its own
   task does not cover and that does not stop its task reaching `DONE`. A problem that stops the task
   keeps the existing path: a `BLOCKED` question or defect. Ideas and improvements are not findings.
2. **The executor files it as an ordinary task block at the end of the same `plan.md`**, with
   `- Status: BLOCKED — found during <TASK-ID>: <reason>.` and the ID `<TASK-ID>-F<n>`. The block
   holds only facts: `WHAT`, `WHERE` (`file:line`), `WHY`, one runnable `Acceptance` re-check,
   `References`. No `HOW`. At most three per task. It lands in the task's own commit.
3. **`arctool add` is the writer.** It appends one block atomically and refuses (exit 1) a block that
   is malformed, reuses an ID from `plan.md` or `plan-archive.md`, is `TAKEN` or `DONE`, fails the
   strict task checks, or, when `BLOCKED`, lacks the `found during <ID>:` reason or the `<ID>-F<n>`
   ID. It self-validates with exit 5. Without `arctool` the executor appends by hand.
4. **The dispatcher reviews findings at the end of the run**, after Verification, one question per
   turn: release (`arctool todo`), sharpen (`/arcdlc:plan`), redesign (`/arcdlc:aic`), or dismiss
   (`/arcdlc:plan` deletes the block). A released finding runs in the next run, never this one. A
   non-interactive run lists them first in its report.

## Consequences

- The plan format contract does not change. `next` skips a `BLOCKED` block, so the running queue
  never picks up a finding by itself.
- `/arcdlc:close` keeps refusing `delivered` while a finding is `BLOCKED`, and lists it under
  `Not done`. That is intended: an undecided finding is work left open.
- `execute` gains its one text write to `plan.md`, through `arctool add`, and `plan` gains the job of
  sharpening or deleting findings.

## Alternatives rejected

- A `## Findings` section at the end of `plan.md` with its own block shape: the parser folds trailing
  text into the last task, so it needed a parser and archiver change.
- A task appended as `TODO`: unreviewed work written by a cheap executor goes straight into the
  running queue, and a finding task can file findings of its own.
- A gap block in `gap.md` with a register-only question line: `gap.md` would gain a second writer and
  a second meaning.
- A new sibling register `findings.md`: one more file nobody opens.
- The dispatcher writes findings from the subagent's report: its context grows on every task.
