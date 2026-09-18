# ADR-0025 — A task pins its executor tier with an `Executor` key

- Status: Accepted
- Date: 2026-09-18
- Initiative: none (bundle-wide, changes the plan contract and `/arcdlc:execute`)
- Amends: [ADR-0019](0019-the-executor-tier-is-asked-for-not-guessed.md). The tier is still
  asked for, never guessed, and the order still runs task, then pin, then question, then
  in-session. What changes is how a task names its tier: an `Executor` key, not a sentence in `HOW`.
- Re-opens and closes: the backlog question "Where should the executor tier pin live?", decided on
  2026-09-16 in favour of `CONTEXT.md` alone. See [docs/backlog.md](../backlog.md).

## Context

ADR-0019 put the tier resolution in a fixed order and said a task could bind its own tier when its
`HOW` "names a model or an effort level". That was the only per-task route, and it had three defects.

It was prose. `HOW` records decisions about the code: signatures, data shapes, edge cases. A sentence
about which model runs the task sits among them and looks like one, and a dispatcher had to read the
whole section to find it. `arctool next --json` returns `how` as one string, so nothing mechanical
could tell a task with a tier from a task without one.

It was a trap. "The task names it" invites the planner to write `HOW: run this on opus` whenever a
block feels hard, which is the tier-as-quality-dial mistake ADR-0021 forbids. The route existed, and
nothing said when it was legitimate.

And it could not do what the engineer actually wanted. The backlog closed the question of a
per-initiative pin on 2026-09-16 on the grounds that the `HOW` route covered the residual. In
practice an engineer who has pinned `sonnet` for a repository still wants one migration task on
`opus, high effort`, or one rename on `haiku`, and wants to say so in the plan, in a line a reader
finds at a glance and a tool can print. A format change was judged too much machinery then. It is
one optional single-line key.

## Decision

**The plan format gains an optional single-line key, `Executor`.** Its value is the tier, written
the way the `CONTEXT.md` pin is written: a model, optionally followed by an effort level, for example
`- Executor: opus, high effort.` or `- Executor: sonnet`. The parser trims surrounding whitespace and
one trailing period and interprets nothing else. It is the eighth defined key; the other seven and
the custom-key rule are unchanged.

**The `Executor` key is step one of the resolution order.** `/arcdlc:execute` resolves the tier for
each task in this order, first answer wins:

1. The task's `Executor` key. It binds that one task and nothing else.
2. The `Executor tier:` pin in `CONTEXT.md`, exactly as written, for every task that has no key.
3. One question to the engineer before the first spawn, for every task that has no key, and for the
   verification subagent. If every `TODO` task carries an `Executor` key there is nothing to ask.
4. In-session mode when there is nobody to ask or nothing to choose.

A `HOW` section no longer names a tier. A plan that still does so is not read for it; `/arcdlc:plan`
moves the tier onto an `Executor` line when it next touches the task.

**A task pin is asked for, like every other tier.** `/arcdlc:plan` writes an `Executor` line only when
the engineer names the task and the tier, in the interview or in the request. It never adds one
because a block looks hard: a block that needs a stronger model to be executed correctly is not
mechanical, and the fix is a sharper block, per ADR-0021. Legitimate reasons are the engineer's:
a task whose cost of a wrong answer is high, or one so small the run tier is a waste on it.

**The capability floor still applies to a task pin.** Whatever the line says, the subagent must be
able to run shell commands, edit files, run the project's test and lint commands, and commit.

**`arctool` carries the pin, and checks only that it is not empty.** `next --json` and `show --json`
return it as `executor`, an empty string when the task has none. `validate --strict` reports
`empty-executor` for a line with nothing on it, because an empty pin is a typo waiting to be read as
"no pin". Nothing in `arctool` judges the value: model names are the engineer's and the harness's.

**Single-task mode notes the pin and cannot apply it.** `/arcdlc:execute <slug> <TASK-ID>` runs in the
current session, where no model can be changed, so the report names the pin and says it was not used.

**The run reports the pin where it was used.** The report names the tier per task where an
`Executor` line overrode the run tier, so a reader can tell a followed pin from an ignored one.

## Justification

- **A tool can print a key; it cannot print a sentence.** `executor` in the JSON gives the dispatcher
  the pin without reading `HOW`, and gives a human `arctool show` to check what a task will run on.
- **One place per decision.** The tier of a task lives on its `Executor` line, the tier of a run in
  `CONTEXT.md`. `HOW` goes back to being about the code only, which is what ADR-0021 asks of it.
- **The override the backlog accepted as residual is now explicit.** An engineer who wants one task
  on a different tier writes one line in the plan instead of editing the repository pin around a
  run.
- **The change is small.** One optional key, one JSON field, one strict rule. Every existing plan
  parses as before, and a plan written for an older `arctool` is unaffected because the key was never
  required.
- **It keeps the tier out of the planner's hands.** Making the key explicit made it possible to say
  plainly when it may be written, which the `HOW` route never did.

## Trade-offs

- **Eight keys, not seven.** Every document that said "seven" now says "eight", and a reader of the
  format has one more thing to know. It is optional, so a plan that never needs it looks the same.
- **A pin in `HOW` written before this ADR is no longer honoured.** The dispatcher reads one place.
  The next `/arcdlc:plan` pass over the task moves it, and until then the run tier applies. This was
  accepted over reading both places, which would have kept the prose route alive.
- **The value is still free text, still honour-based.** A model name the harness does not have fails
  the spawn loudly, as before. `arctool` prints the pin and does not verify the spawn used it.
- **A per-task pin can be misused as a quality dial.** The rule against it is written in the plan
  skill and here. Nothing mechanical can tell a legitimate pin from a lazy one, so the reviewer of a
  plan reads the `Executor` lines with that question in mind.
