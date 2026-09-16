# ADR-0021 — A planned task must be mechanical, not merely cheap to run

- Status: Accepted
- Date: 2026-09-16
- Initiative: none (bundle-wide, changes the premise of the plan contract)
- Amends: [ADR-0019](0019-the-executor-tier-is-asked-for-not-guessed.md). The tier is still asked for
  and never guessed. What changes is why a cheap tier is safe: the plan being mechanical, not the
  planner having aimed at a weaker model.

## Context

`plan-format.md` opened by saying a plan is "written to be executed by a less capable model than the
one that planned it". That sentence carried two different claims at once, and only one of them holds.

The engineering claim is that a task block must contain every decision, so the executor never has to
judge. That is right and it is what makes the format work.

The economic claim is that the executor should therefore be a cheaper model. That is a bet on pricing,
and pricing moved. `/arcdlc:execute` already said the truer thing in its own tier section: a tier is
two dials, model and effort, and a strong model at its lowest effort often beats a weaker model at the
same price. The two documents disagreed.

Measuring the 55 tasks this repository actually executed showed where the framing leads when it is
followed literally:

| Initiative | Tasks | Median `HOW` | Max |
|---|---|---|---|
| ordering | 5 | 33 lines | 52 |
| antigravity-cli | 5 | 11 | 23 |
| source-library-cleanup | 30 | 10 | 25 |
| cursor-support | 4 | 6 | 7 |
| initiative-lifecycle | 11 | 0 | 0 |

All 55 pass `arctool validate --strict`, and 80% carry a runnable `GIVEN … WHEN … THEN`, so the
contract itself is sound. But ORD-2's 52-line `HOW` names the exact function signature, every flag and
its help string, every exit-code mapping, the `fmt.Printf` format string, and the literal usage text.
That is the implementation written in prose. Aiming at a weaker model pushed the planner into writing
the code twice, once as prose and once as code, which costs more than the cheaper tier saves.

At the other end, `initiative-lifecycle` ran eleven tasks with no `HOW` at all, because the key did not
exist yet. There the executor was guessing, which is what `HOW` was added to stop.

## Decision

The requirement on a task is that it is **mechanical**: executable by a model that has no context but
the block and the files it names. No judgement call, no open decision, nothing to infer from a
conversation the executor never saw.

Mechanical does not mean small. It means nothing is left to judgement.

Two corollaries follow, and they are the point of the pipeline having three stages:

- **A task that cannot be made mechanical is not a task that needs a stronger executor.** It is a task
  that is too big and must be split, or a design that is not finished. An unfinished design goes back
  to `/arcdlc:aic`, which is where the hard thinking belongs. Raising the tier hides the defect: the
  work gets done by judgement nobody recorded, and the next run of that plan behaves differently.
- **`HOW` records decisions, never code.** Signatures, naming, data shapes, algorithm choice, edge
  cases, error handling belong there. The implementation line by line does not. Writing the code in
  prose is the signal that the task is too big or a decision is still open.

Which tier then runs the queue becomes an ordinary economics decision, made per initiative rather than
baked into the format. Repetitive tasks with a small `HOW` suit a cheap tier; design-heavy ones suit
the same model at low effort.

## Consequences

- `plan-format.md` states the mechanical bar and adds authoring rule 8, which forbids code in `HOW`.
- `/arcdlc:plan` Step 2 names the two things an unmechanical task is telling you, and its litmus test
  is now the mechanical check rather than a self-sufficiency check.
- `/arcdlc:execute` explains that a cheap tier is safe because the plan is mechanical, and that a block
  which is not mechanical is a defect the tier must not paper over.
- The courage virtue, verbatim in all ten skills, now reads "a task that is not mechanical is a plan
  defect", which is the rule it was always pointing at.
- `/arcdlc:assist` and `/arcdlc:examinate` write mirrored blocks to the same bar.
- Nothing in `arctool` changes. `unverifiable-acceptance` already enforces the half of this that can be
  checked mechanically; the rest is a planning judgement no parser can make.

## Alternatives considered

- **Keep the weaker-model framing and accept the prose-code duplication.** Rejected: the duplication is
  a real cost, it is invisible in the plan, and it grows with task difficulty.
- **Drop the tier question and always run the planner's model.** Rejected: it throws away a genuine
  saving on repetitive work, and ADR-0019's reason for asking still stands.
- **Cap `HOW` length mechanically.** Rejected: length is a symptom, not the rule. A long `HOW` full of
  real decisions is fine; a short one that hides a decision is not.
