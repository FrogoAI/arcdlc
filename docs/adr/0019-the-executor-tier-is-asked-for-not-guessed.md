# ADR-0019 — The executor tier is asked for, not guessed

- Status: Accepted
- Date: 2026-09-11
- Initiative: none (bundle-wide, tightens `/arcdlc:execute`)
- Relates to: `skills/plan/references/plan-format.md` (opening line) and Step 2 of
  `skills/plan/SKILL.md`, which both already write each task for a weaker executor

## Context

The plan format has said from the start that a task block is executed by a less capable model than the
one that planned it. `/arcdlc:plan` repeats it in Step 2 and checks for it in the self-sufficiency
pass. Two files carried the promise, and no file collected on it.

`/arcdlc:execute` probes whether the harness can spawn subagents, then spawns one per task for the
clean context. It never said at which tier, so the dispatcher spawned at its own tier: the most
expensive model in the run implemented every task, and an engineer who asked for a cheaper executor
got the same tier back. The report named no tier, so the gap stayed invisible.

Two things then became clear about picking the tier automatically.

A tier is not one dial. A model has an effort or reasoning level, and a strong model at its lowest
effort is a real option, often the best one for a block that already carries every decision. "The
lowest model that clears the floor" cannot express that.

And the choice is not ours. Which models a harness exposes, what they cost on this account, and how
much effort this queue deserves are all facts the engineer has and the dispatcher does not. A rule
that guesses is a rule that quietly overrides them.

## Decision

**The floor is capability.** A task subagent must be able to run shell commands, read and edit files
in the repository, run the project's test and lint commands, and write a commit. That is what the
per-task contract needs, nothing more, and it is the sanity check on any answer.

**The tier is the model plus its effort level.** Both dials count, and the pin and the question both
carry both, for example `opus, low effort`.

**`/arcdlc:execute` never picks the tier by itself.** It resolves in a fixed order, first answer wins:

1. A task's `HOW` names a model or an effort level. The planner decided it, so it binds that task.
2. The pin: a line that starts with `Executor tier:` in `CONTEXT.md`, used exactly as written.
3. The question: one question, one turn, before the first spawn, with a recommended answer, naming the
   options this harness exposes. The answer covers every task in the run and the verification
   subagent, and the skill offers to write it into `CONTEXT.md` as the pin.
4. Nobody to ask (a non-interactive run) or nothing to choose (a harness with one tier): in-session
   mode, and the report says so.

**Never spawn at the dispatcher's own tier and call it cheaper.** A same-tier subagent buys clean
context, not a cheaper executor, and the report must not claim otherwise.

**The pin lives in `CONTEXT.md`.** Every ArcDLC skill already reads that file, and `arctool sync`
rewrites only `AGENTS.md` and `README.md`, so nothing in `CONTEXT.md` gets clobbered.

**A thin task block is a plan defect, not a reason to raise the tier.** When the executor meets a
decision its block does not carry, it grills the engineer or blocks the task, as the per-task contract
already says. Nobody retries that task at a higher tier. `/arcdlc:plan <slug>` sharpens the block at
planner tier instead.

**The run names the tier it used.** The report says which tier ran and where the choice came from: a
task's `HOW`, the pin, the answer given this run, or in-session because there was nothing to choose.

## Justification

- **The promise was already made twice.** The plan format is written for a weaker executor. A
  dispatcher that spawns at its own tier makes that a decorative sentence.
- **Asking is cheaper than guessing wrong.** One question costs a turn. A wrong guess costs a full
  queue at the wrong price, or a queue of blocked tasks, and the engineer finds out afterwards.
- **Only the engineer knows the lineup.** Model names, availability, effort settings and account
  pricing change every few months and differ per harness. This bundle should name none of them.
- **A pin makes the question one-time.** Nothing here forces a team to answer the same question twice:
  the answer becomes a line in `CONTEXT.md`, which is where the next session reads it.
- **A functional floor is checkable, a price is not.** "Shell, file edits, test and lint, a commit" is
  a list an agent can verify against the harness it is standing in. "Cheaper" is not.
- **Reporting the tier makes the rule enforceable.** An engineer who cannot see which tier ran cannot
  tell a followed rule from an ignored one. That is how this gap survived.

## Trade-offs

- **A whole-queue run stops for a question.** The run is no longer a single unattended command on a
  project with no pin. The pin removes the stop, and a non-interactive run skips the question by
  falling to in-session mode instead of hanging on it.
- **A weaker executor blocks more often.** More tasks come back as `BLOCKED` with a question, and the
  fix is a `/arcdlc:plan` pass. That is the intended pressure: it pushes decisions into the plan,
  where they are written once instead of re-derived per run.
- **Losing orchestrator mode when there is nothing to choose.** A harness with one subagent tier falls
  back to in-session mode, so the run gives up per-task fresh context and stops at task boundaries.
  The report says why.
- **The answer is free text.** A model name the harness does not have makes the spawn fail. The
  failure is loud and names the line or the answer to fix, and no fallback guessing hides it.
- **The tier is honour-based.** Nothing in `arctool` reads the pin or checks the spawn today, so the
  rule holds as far as the dispatcher obeys it. A later version can have `arctool` print the resolved
  tier, which turns the report line into something a human can check.
