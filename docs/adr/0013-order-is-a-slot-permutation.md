# ADR-0013 — `arctool order` permutes slots, it does not move tasks

- Status: Accepted
- Date: 2026-09-04
- Initiative: [ordering](../aics/ordering/aic.md)

## Context

`/arcdlc:execute` takes the first `TODO` block in `docs/aics/<slug>/plan.md`, reading top to bottom.
The order of blocks in the file is the run order. Today the only way to change that order is to cut
and paste Markdown blocks by hand. That is exactly the kind of whole-file edit `arctool` exists to
make safe.

Engineers asked for `arctool order <ID> <ID> <ID>`. That command line can mean three different
things, and they give three different results for the same input:

1. **Full permutation.** The list must name every task. The list is the new order.
2. **Move to front.** The named tasks go to the top in the given order. Everything else follows.
3. **Slot permutation.** The named tasks keep the positions they already occupy. Only their contents
   are permuted among those positions.

## Decision

`arctool order` is a slot permutation.

- Collect the positions of the named tasks in the current file. Sort those positions ascending.
- Put the named tasks into those positions, in the order given on the command line.
- Every task that is not named keeps its position.

Given `T1 T2 T3 T4 T5 T6`, the command `arctool order T5 T2 T3` fills positions 2, 3 and 5 with T5,
T2 and T3. The result is `T1 T5 T2 T4 T3 T6`.

The command is mechanical. It applies the order it is given and never judges it. The plan format
gains no dependency field, `arctool validate` gains no ordering check, and there is no automatic
sort. There is one verb: no `arctool move --before/--after/--top/--bottom`.

## Justification

- **An unnamed task can never drift.** The plan is a dependency-ordered queue, and no dependency is
  written down anywhere a machine can read. Under "move to front", naming one late task silently
  hoists it above tasks it may need. Under slot permutation, a task you did not name sits in the
  same place afterwards, always. That property is what makes the command safe to hand to an agent.
- **Full permutation is correct and unusable.** Naming all twenty tasks to swap two is not a command
  anyone runs twice.
- **A dependency field is a different initiative.** `plan-format.md` is a contract. Changing it
  forces matching changes in the parser, validator, mutator, archiver, their tests, and every skill
  that references the format, all in one change set. It is also new state that a weaker executor
  would have to keep correct. A checkable ordering can be added on top of this decision later
  without changing this command line, so nothing here is a dead end.
- **One verb, one mental model.** The caller is almost always an agent that has just run
  `arctool list` and already holds every ID, so naming a span costs it nothing. A second mutation
  verb would be a second thing for every skill to document and probe for.

## Trade-offs

- **"Move to first" is a swap, not a shift.** `arctool order T3 T1` on `T1 T2 T3` gives `T3 T2 T1`.
  T2 stays in the middle. To hoist T3 and keep the rest in relative order you name the span:
  `arctool order T3 T1 T2` gives `T3 T1 T2`. Every move is expressible, and the cost is extra IDs on
  a command line that an agent generates anyway. This is written out in these words because "move to
  first" is the phrase that misleads.
- **Naming a subset that is already in order does nothing.** `arctool order T1 T3` on `T1 T2 T3` is a
  no-op, because T1 already holds the first named slot. This surprises people. It follows directly
  from the property in the first justification, so we accept it. The command prints the resulting
  order, so a no-op is visible rather than silent.
- **A wrong order is applied without complaint.** `arctool order` will put a task above the task it
  depends on. Correctness of the order stays the author's job, per authoring rule 5 in
  `plan-format.md`. `/arcdlc:execute` catches the inversion when it runs into it.
