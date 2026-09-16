# ADR-0023 — The dispatcher verifies, and a spawned subagent never guesses

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

**The dispatcher verifies every finished task independently**, in four checks, cheapest first:

1. Status `DONE` and a commit exists. Necessary, nowhere near sufficient.
2. **Run the task's `Acceptance` criteria.** Every criterion names a command, a path, a test, or an exit
   code, because `arctool validate --strict` refuses a plan where one does not. That requirement was
   added for the executor's benefit; it turns out to be what makes independent verification cheap.
3. `git show --stat HEAD` against `WHERE`. Files touched outside it, or a `WHERE` file left untouched,
   means a different problem was solved.
4. A report naming no grilled decisions on a task that had an ambiguity is the signature of a guess.

The context rule is relaxed from "never read source files or diffs" to: read plan state, reports, commit
subjects, `--stat`, and the output of the acceptance commands. Read a diff or a source file only when a
check fails and you need to name why.

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

## Alternatives considered

- **Have the subagent wait for an answer.** Rejected: it has no channel to a human, so waiting is
  hanging.
- **Let the subagent grill through the dispatcher as a relay.** Rejected as complexity for no gain. The
  block-and-return path already carries the question, and it leaves the task in a state the next run can
  resume from.
- **Have the dispatcher re-read the diff for every task.** Rejected: that is the context cost the
  original rule was protecting against, and the acceptance criteria answer the same question for less.
- **Trust the subagent and verify once at the end.** Rejected: the verification phase already runs at
  the end, and by then several commits rest on the bad one.
