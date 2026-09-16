# ADR-0023 — A task verifies itself, the queue is verified once, and a subagent never guesses

- Status: Accepted
- Date: 2026-09-16
- Initiative: none (bundle-wide, changes `/arcdlc:execute` orchestrator mode)
- Builds on: [ADR-0021](0021-a-planned-task-must-be-mechanical.md) (tasks are mechanical) and
  [ADR-0019](0019-the-executor-tier-is-asked-for-not-guessed.md) (the tier is asked for)

## Context

Orchestrator mode spawns one subagent per task, at a tier the engineer chose, usually cheaper than the
dispatcher's. Two things were wrong with it, and they compound.

**The executor was its own judge.** Step 5 of the per-task contract has the executor verify its own
acceptance criteria and then mark the task `DONE`. The dispatcher's only check was that the status said
`DONE` and a commit existed. A wrong implementation passes both. A cheaper model self-certifies
optimistically: it meets a contradiction, implements something plausible, decides it passed, and
commits. ADR-0021 makes that failure silent by construction, because a mechanical block carries no
reasoning, so the executor has nothing to notice it contradicted. The dispatcher was also told
"never read source files or diffs", so it could not have caught it even in principle.

**The skill told a subagent to do something impossible.** "When a task is unclear, grill before you
code" instructed the executor to interview the engineer. A spawned subagent has no engineer in its
session. It cannot grill, it cannot wait, and the instruction it could actually follow was the one
thing nobody wants: guess and carry on.

The dispatcher, meanwhile, is the session the engineer started. It has the human, and it is the more
capable model. It was the only party able to answer, and it was not asked to.

## Decision

**A spawned subagent never grills and never guesses.** On any ambiguity, contradiction, or missing
decision it runs `arctool block <id> -m "<the question>"` and returns the question unanswered.
Returning a question is a correct outcome for a subagent, not a failure, and its spawn prompt says so.

**The dispatcher answers.** It grills the engineer, one question per turn, then either re-spawns the
task with the answer or, when the answer changes the plan, leaves the task blocked and hands back to
`/arcdlc:plan`. In a non-interactive run it reports the question verbatim rather than answering it
itself to keep the queue moving.

**A task verifies itself, and the dispatcher does not re-check it.** The subagent reaches `DONE` only
when every `Acceptance` criterion demonstrably holds, so `DONE` is its assertion and its spawn prompt
says nobody downstream will re-check it. The dispatcher re-running those same criteria would catch only
a subagent that lied about running them, and would pay that output on every task in the queue. It is not
the per-task check.

**`BLOCKED` is the only thing the dispatcher engages with**, because a block is the one outcome that
needs what only the dispatcher has: an engineer to ask.

**The queue is verified once, at the end, and that is the real gate.** A task can only verify itself
against the block it was given. Nothing in a per-task check, by the executor or the dispatcher, can ask
whether two tasks agree. Two can each be correct and still contradict, and that hides precisely where
tasks never touch. So the Verification phase has three parts: the queue really is finished (no leftover
`TAKEN` from a crashed session, a commit per `DONE` task, `validate --strict` clean including the source
stamp); the project is whole (repository-level build, test and lint, delegated to one subagent because it
is mechanical); and the tasks agree (overlapping `WHERE`, opposing decisions in commit bodies, seams
between independently planned producers and consumers). The third part stays with the dispatcher because
it is judgement, and a green test suite is exactly what two contradicting tasks look like when neither
has a test for the other's assumption.

A contradiction found there is a plan defect. It goes back to `/arcdlc:plan`, never into a fix-up commit,
because the plan produced two tasks that disagree and will do so again.

## Consequences

- The dispatcher's context grows by a `--stat` and some command output per task, not by diffs. A long
  queue still finishes in one invocation, which was the point of the original restriction.
- `unverifiable-acceptance` becomes load-bearing twice over. It was added so a weak executor had
  something real to check; it is now also what lets the dispatcher check the executor.
- A failed verification releases the task (`arctool todo`) or blocks it, and stops the run. It never
  re-spawns the same task at a higher tier: that is the tier-as-quality-dial mistake, and it hides a
  defect the plan should carry.
- In-session mode is unchanged. There the executor *is* the session the engineer started, so it grills
  directly, and there is no second party to verify it. That asymmetry is real and worth knowing: a
  whole-queue run with subagents is checked twice, an in-session run once.
- The dispatcher stays thin by construction rather than by discipline: in the loop it reads plan state,
  one line of notes per task, and a block reason. On a 150-task queue that is the difference between a
  dispatcher that finishes and one that compacts halfway. It still stops at a task boundary and hands back
  if it fills up, because nothing makes a session immortal.
- Detection of a cross-task contradiction is deferred to the end of the queue. That is the cost of this
  shape, and it is accepted: per-task checking cannot find these at all, so the choice is late detection
  or none. Dependency ordering limits the blast radius, because a task that broke something downstream
  usually fails the next task that touches it.

## Alternatives considered

- **Have the subagent wait for an answer.** Rejected: it has no channel to a human, so waiting is
  hanging.
- **Let the subagent grill through the dispatcher as a relay.** Rejected as complexity for no gain. The
  block-and-return path already carries the question, and it leaves the task in a state the next run can
  resume from.
- **Have the dispatcher re-run every task's acceptance criteria.** Considered and rejected after being
  written: it re-runs exactly what the executor just ran, so it catches only a subagent that lied about
  running them, and it pays that output on every task. The check that sounded strongest added the least.
- **Spawn a separate verifier subagent per task.** Rejected: it keeps the dispatcher thin, which is real,
  but doubles the spawn count to check something the executor already checked. The same reasoning that
  rejects dispatcher re-verification rejects this.
- **Have the dispatcher read the diff for every task.** Rejected: that is the context cost the original
  rule was protecting against, and it invites rejecting work on taste rather than on the contract.
