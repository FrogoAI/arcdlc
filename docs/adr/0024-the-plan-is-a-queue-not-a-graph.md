# ADR-0024 — The plan is a queue, not a graph, and it runs in order

- Status: Accepted
- Date: 2026-09-16
- Initiative: none (bundle-wide, settles a question deferred twice)
- Relates to: [ADR-0013](0013-order-is-a-slot-permutation.md) (`arctool order` permutes slots) and
  [ADR-0023](0023-the-dispatcher-verifies-and-the-subagent-never-guesses.md) (one subagent at a time)

## Context

The `ordering` initiative deferred a question in 2026-08: should a task block gain a `- DEPENDS:` key,
with `arctool order --auto` deriving the run order from it? It was raised again on 2026-09-16 with a
stronger reason: orchestrator mode spawns strictly one subagent at a time, so a 150-task queue is 150
serial spawns, and most of those tasks are probably independent of each other.

Both times the pull was toward turning the plan into a dependency graph, which is the shape that would
let independent tasks run in parallel.

## Decision

**No `- DEPENDS:` key, no derived ordering, no parallel execution.** A plan is an ordered list. It runs
top to bottom, one task at a time, and the order is whatever the file says. `arctool order` remains the
only way to change it, and it is a deliberate act by a person.

## Consequences

**A queue is resumable with one pointer.** `arctool next` returns the first `TODO` block and that is the
entire scheduling state. A graph needs per-node state, a readiness computation, and an answer for what
happens when a node fails with others in flight. The queue answers all of that by not having the
question: a failed task stops the run, and the next invocation resumes exactly where it stopped. That
property is what makes a 150-task plan survivable across any number of sessions, and it is worth more
than wall-clock speed.

**Order stays visible.** A dependency graph is read by a tool; a list is read by a person. The engineer
who opens `plan.md` sees the run order without computing it, and fixes a wrong order by moving a block.
An implicit order derived from declared dependencies would be correct more often and understood less
often.

**Commits stay linear.** One task, one commit, never interleaved. The history is bisectable, each commit
is reviewable against exactly one task block, and a revert has an obvious boundary. Parallel execution
buys throughput and spends all three.

**Concurrency failures cannot be expressed in a mechanical block.** [ADR-0021](0021-a-planned-task-must-be-mechanical.md)
requires a task to be executable by a model with no context but the block. Two tasks running at once can
touch the same file, race on the same status line, or each assume the other has not landed yet. None of
that fits in a block, and a mechanical executor has no way to notice it.

**The cost is wall-clock time, and it is accepted.** A long queue takes as long as its tasks take. Where
throughput matters, run separate initiatives, which are already independent by construction: different
folders, different plans, different queues.

**The `- DEPENDS:` question is closed, not deferred.** Re-open it only with a concrete failure that
ordering cannot express, not with a wish for speed.

## Alternatives considered

- **`- DEPENDS:` with `arctool order --auto`.** Rejected above. It solves an ordering problem the file
  already solves, and its real payoff is parallelism, which costs linear history and adds failure modes
  a mechanical block cannot carry.
- **Parallel spawning with a file-level lock.** Rejected: it makes commits interleave, which breaks the
  one-task-one-commit rule that makes the history reviewable.
- **Parallel spawning for tasks with disjoint `WHERE`.** Rejected: `WHERE` names the files a task is
  expected to change, not every file it reads. Disjoint `WHERE` does not prove independence.
