# Task Ordering

> Add an arctool command that re-orders task blocks in plan.md, because /arcdlc:execute runs them top to bottom.

Format: AIC (template: `source-map/source/AIC Template.md`). Status: Draft. Decompose with
`/arcdlc:plan ordering`. Interview date: 2026-09-04.

## Goals

### 🟢 Business Case

`/arcdlc:execute` reads `docs/aics/<slug>/plan.md` from the top and runs the first block whose
status is `TODO`. The order of blocks in the file is the run order. Nothing else decides it.

Today the only way to change that order is to cut and paste Markdown blocks by hand. That is the one
plan mutation `arctool` does not cover, and it is the riskiest one. Every other write to `plan.md`
goes through the tool: `take`, `done`, `block` and `todo` rewrite a single `- Status:` line and leave
every other byte alone, and `archive` re-parses its own output before it writes. A hand edit has none
of that. It can drop a block, split a multi-line `Acceptance` section, or leave a half-pasted task in
a queue that an agent is part-way through.

This initiative adds `arctool order`. It is the missing mutation, built to the same rules as the
ones that already exist.

### 🟢 Functional Overview

1. **New command.** `arctool order <ID> <ID> [<ID>…]`, with the standard selection flags
   `--aic SLUG` or `--plan PATH`, plus `--dry-run`.
2. **Slot permutation.** The command collects the positions the named tasks occupy today, sorts
   those positions ascending, and fills them with the named tasks in the order given on the command
   line. A task that is not named keeps its position. For `T1 T2 T3 T4 T5 T6`, the command
   `arctool order T5 T2 T3` fills positions 2, 3 and 5, giving `T1 T5 T2 T4 T3 T6`.
3. **Mechanical.** The command applies the order it is given. It does not check the order against
   anything, because nothing machine-readable records which task needs which.
4. **Self-check before writing.** The new bytes are re-parsed and checked: same number of blocks,
   same set of task IDs, same count per status, and every block byte-identical to the original once
   trailing newlines are ignored. Any mismatch means nothing is written and the command exits `5`.
5. **All-or-nothing failure.** Fewer than two IDs, or the same ID twice on the command line, exits
   `2`. An ID that is not in the plan, an ID that matches two blocks, or a plan with no tasks at all
   exits `3`. A plan already in the requested order prints that and exits `0` without writing. Every
   bad ID is reported in one message, not one per run.
6. **Output.** On success the command prints the resulting document order using the same columns as
   `arctool list`, so the caller does not need a second call to see the result. `--dry-run` prints
   exactly that and writes nothing. There is no `--json`.
7. **Layout.** Blocks are re-joined the way `archive` already joins them: trailing newlines trimmed
   per block, one blank line between blocks, the preamble and the "Completed (archived to …)" ledger
   byte-identical, the file written by temp-file plus rename.
8. **`/arcdlc:plan` gains an ordering step.** Authoring rule 5 in `plan-format.md` already makes
   order the author's job. The skill now names the command that discharges it, behind the mandatory
   `command -v arctool` probe, with the manual fallback spelled out (move the blocks by hand).
9. **`/arcdlc:execute` may propose, never apply.** When the executor is implementing a task and finds
   it needs something a task below it will build, it does what step 7 of its per-task contract
   already says: `arctool block <id> -m "needs <OTHER-ID>, which is below it"`, then stop the run. It
   adds one line to its report, the exact `arctool order …` command that would fix the plan. It never
   runs that command.

Code lands in `internal/plan/order.go` (with `order_test.go`) and a new `cmdOrder` in
`cmd/arctool/main.go`.

### 🟢 Quality Goals

1. **Safety.** A whole-file rewrite must never produce a plan the engineer cannot reconstruct. The
   self-check plus the atomic write mean the file is either the old one or a correct new one.
2. **Predictability.** A task nobody named is in the same place afterwards. Always. An agent can run
   the command without holding the whole file in its head.
3. **Fit.** No new concepts. The command reuses the selection flags, the exit codes, the atomic
   write, and the block layout that `archive` already established.

### 🟢 Organizational Constraints

- Skills stay install-agnostic. Any `SKILL.md` edit works both as `/arcdlc:<name>` in the plugin and
  as `arcdlc-<name>` flat, with both path forms intact.
- `arctool` stays optional in skills. `/arcdlc:plan` probes `command -v arctool` and documents the
  manual fallback.
- No new sub-skill, so `SUBSKILLS` in `install.sh` and the CI skill-layout checks are untouched.
- Everything written here follows `skills/source-map/source/Writing Style.md`: plain English, no long
  dashes in prose.
- `arctool` goes to `0.11.0`. Both plugin manifests go `0.12.0` to `0.13.0`, in lockstep.

### 🟢 Technical Constraints

- `arctool` stays pure standard library, `CGO_ENABLED=0`, static release binaries.
- Exit codes are interface: `0` ok, `1` contract failure, `2` usage, `3` not found or empty, `4` I/O,
  `5` self-validation. Nothing is renumbered. Code `5` widens in wording only, from "archive
  self-validation" to "self-validation".
- Writes stay byte-preserving where they can be and atomic always: temp-file plus rename, permissions
  copied from the original.
- `skills/plan/references/plan-format.md` does not change. The file states its own boundary, "keep
  format rules in this file, no runner instructions", and a tool command is a runner instruction.
- `go build ./...`, `go test ./...`, `gofmt -l .` and `go vet ./...` all stay clean.

### 🟢 Business Context

The system under construction is the ArcDLC bundle, meaning the skills plus `arctool`. The new
command sits beside the existing plan mutations and touches one file.

```
 engineer ──arctool order T5 T2 T3 --aic ordering──▶ arctool
 agent    ──/arcdlc:plan <slug>────────────────────▶ skill ──▶ arctool

                                arctool order
                                     │ reads
                                     ▼
                        docs/aics/<slug>/plan.md
                                     │ permute blocks in place
                                     ▼
                        self-check (re-parse, compare)
                             │ pass              │ fail
                             ▼                   ▼
                 temp file + rename         exit 5, file untouched
                             │
                             ▼
                   new order printed to stdout

 /arcdlc:execute ──reads first TODO top-to-bottom──▶ docs/aics/<slug>/plan.md
 /arcdlc:execute ──proposes an order command, never runs it──▶ engineer
```

## Architectural Hypotheses

### 🔵 H1 — `arctool order` is a slot permutation ([ADR-0013](../../adr/0013-order-is-a-slot-permutation.md))

- **Context:** `arctool order A B C` could mean a full permutation of the whole plan, a move of the
  named tasks to the front, or a permutation of the slots those tasks already occupy. The three give
  different results for the same command line.
- **Decision:** Slot permutation. Named tasks keep their positions and swap contents among them.
  Unnamed tasks never move.
- **Justification:** No dependency is recorded anywhere a machine can read, so a rule that moves
  tasks the caller did not name can silently break the queue. Slot permutation is the only reading
  where an unnamed task is guaranteed to stay put.
- **Trade-offs:** "Move to first" becomes a swap. `arctool order T3 T1` on `T1 T2 T3` gives
  `T3 T2 T1`, leaving T2 in the middle. Hoisting T3 while keeping the rest in relative order needs
  the whole span: `arctool order T3 T1 T2`. Naming a subset that already sits in the requested order
  does nothing at all.

### 🔵 H2 — Ordering stays mechanical, with no dependency field

- **Context:** The obvious next step is a `- DEPENDS: <ID>, <ID>` key, so `order` could refuse a bad
  permutation and `validate` could flag a plan that is already wrong.
- **Decision:** Out of scope. No format change, no ordering check in `validate`, no `--auto` sort.
- **Justification:** `plan-format.md` is a contract. Touching it forces matching changes in the
  parser, validator, mutator, archiver, their tests and every skill that references the format, all
  in one change set. It is also new state a weaker executor would have to keep correct. The command
  line designed here does not change if dependencies are added later, so this is not a dead end.
- **Trade-offs:** `arctool order` will happily put a task above the task it depends on and say
  nothing. Correctness of the order stays the author's job.

### 🔵 H3 — One verb, no `move`

- **Context:** A one-task hoist under H1 costs naming every task it jumps over. A second verb
  (`arctool move <ID> --before <ID>`, with shift semantics) would make that one short command.
- **Decision:** Ship `order` alone.
- **Justification:** The caller is nearly always an agent that has just run `arctool list` and holds
  every ID, so a long argument list costs it nothing. A second mutation verb is a second set of edge
  cases (move before itself, move where it already is) and a second thing every skill has to
  document and probe for.
- **Trade-offs:** A human typing this by hand writes more than they would like. `move` stays addable
  later without breaking `order`.

### 🔵 H4 — No status guardrail

- **Context:** `take`, `done`, `block` and `todo` refuse transitions that do not fit the lifecycle,
  unless given `--force`. Ordering could copy that and refuse to move a `DONE` or `TAKEN` task.
- **Decision:** No guardrail. Any task may be named, whatever its status. There is no `--force`.
- **Justification:** Ordering is not a status transition, so there is no lifecycle to protect.
  `/arcdlc:execute` reads the first `TODO`, so a misplaced `DONE` block cannot change what runs next,
  and `/arcdlc:archive` removes `DONE` blocks from the plan anyway.
- **Trade-offs:** A typo that names a completed task quietly shuffles a `DONE` block down the file.
  The printed result is the only place it shows up.

### 🔵 H5 — Self-check before writing, on exit 5 ([ADR-0014](../../adr/0014-exit-5-means-self-validation-failed.md))

- **Context:** `order` is the second whole-file rewrite after `archive`. `archive` already re-parses
  its output and exits `5` when the result is wrong.
- **Decision:** `order` runs the same kind of check and reuses exit `5`, whose documented meaning
  widens from "archive self-validation" to "self-validation".
- **Justification:** The parser is already there, so the check is nearly free, and a permutation bug
  is the one failure that destroys work the engineer cannot get back. Reusing `5` rather than adding
  `6` keeps every skill that already reads `5` as "nothing was written, your file is intact" working
  with no edit. Widening wording is not renumbering.
- **Trade-offs:** Every place that spells out the exit table changes in the same change set:
  `AGENTS.md`, the `usage` string in `cmd/arctool/main.go`, `README.md`, and any `SKILL.md` that
  lists codes. An operator reading `5` in a log now checks which command produced it.

### 🔵 H6 — Strict, all-or-nothing failure contract

- **Context:** Bad input comes in several shapes: too few IDs, a repeated ID, an unknown ID, an ID
  that matches two blocks, an empty plan.
- **Decision:** A malformed command line exits `2` (fewer than two IDs, repeated ID). A lookup miss
  exits `3` (unknown ID, ambiguous ID, no tasks). Nothing is written in any failing case, and one
  message names every bad ID.
- **Justification:** The split already exists in `arctool`: `2` is what you can detect without
  reading the plan, `3` needs the plan, and `cmdMutate` already answers an ambiguous ID with `3`.
  Reporting every bad ID at once matters because the caller is an agent: one message listing three
  typos is one retry instead of three.
- **Trade-offs:** A plan that is silently reordered except for one fat-fingered ID would be worse
  than no reorder, so leniency is rejected outright. The cost is that one bad ID in a long list
  cancels the whole command.

### 🔵 H7 — Blocks are re-joined the way `archive` joins them

- **Context:** The parser gives each block a raw span that runs to the next `###` heading. A block's
  span swallows the blank lines after it, and the last block's span swallows everything to end of
  file. Moving spans verbatim would carry that spacing around with them.
- **Decision:** Trim trailing newlines per block, re-join with exactly one blank line between blocks,
  leave the preamble and ledger byte-identical, end the file with a single newline. In the H5
  self-check, "byte-identical" means the block content with trailing newlines trimmed.
- **Justification:** `archive` already does this, so the repo has one answer for how a whole-file
  rewrite lays blocks out. `plan.md` is a generated queue, so uniform spacing is correct.
- **Trade-offs:** Anything sitting between two blocks, such as a stray note line, belongs to the
  preceding block's span and travels with that block to its new position. `archive` has the same
  hazard today. Fixing it means teaching the parser about inter-block regions, which is the format
  change H2 rules out.

### 🔵 H8 — `/arcdlc:plan` owns ordering, `/arcdlc:execute` only proposes

- **Context:** Order is decided when the plan is authored, but the inversion is felt when the plan
  runs.
- **Decision:** `/arcdlc:plan` gains the `arctool order` step. `/arcdlc:execute` may propose a
  reorder, but only at discovery time, meaning while implementing a task it finds it needs something
  a lower task will build. It reports the command and stops. It never runs it.
- **Justification:** A runner that rewrites its own queue mid-run is very hard to debug, and the
  skill already has the right answer for a task it cannot finish: `block` with a reason, then stop.
  Discovery time is a fact rather than an inference. A pre-flight scan would have to read every
  block's prose, which contradicts the orchestrator rule "never read source files or diffs, only
  plan state, subagent reports, and commit subjects", the rule that lets a long queue finish in one
  invocation.
- **Trade-offs:** An inversion is only found when the executor walks into it, after the task above it
  has already been claimed and blocked. Nothing warns earlier.

### 🔵 H9 — `plan-format.md` does not change

- **Context:** Authoring rule 5 there says "order blocks by dependency", which is exactly what this
  command serves.
- **Decision:** Leave the file alone. Document the command in `arctool --help`, `README.md` and
  `/arcdlc:plan`.
- **Justification:** The file states its own boundary: keep format rules here, no runner
  instructions. `AGENTS.md` also makes edits there expensive on purpose, and rule 5 stays true either
  way. The obligation is the author's; the tool is only how they discharge it.
- **Trade-offs:** Someone reading only `plan-format.md` will not learn that the command exists.

## Assessment

### 🔴 Technical Challenges & Risks

- **A permutation bug destroys a plan.** This is the reason H5 exists. The tests must cover the
  cases that break naive index maths: one ID named, every ID named, named IDs already in the
  requested order, the first and last blocks swapped, and a plan whose only block is named.
- **The exit table lives in four places.** `AGENTS.md`, the `usage` string in `cmd/arctool/main.go`,
  `README.md`, and skill files that list codes. H5 changes the wording of code `5` in all of them, in
  one change set. Missing one leaves a stale table.
- **Inter-block content moves with the block above it.** Named in H7. It is inherited from `archive`
  and accepted here, but it needs a test that pins the behaviour so it does not change by accident.
- **Line endings.** `archive` joins blocks with `\n` even when `Plan.CRLF` is true, so a CRLF plan
  keeps CRLF inside each block and gets LF between them. `order` copies that behaviour for
  consistency. If that is wrong, it is wrong in both places and should be fixed in both.
- **Duplicate task IDs in a plan.** `Plan.ByID` can return more than one block. `cmdMutate` refuses
  to guess and exits `3`; `order` must do the same, and a test must cover it, because a plan built by
  merging gap tasks from `/arcdlc:examinate` is where duplicates come from.
- **"Move to first" is the phrase people will use.** Under H1 it is a swap, not a shift. Both the
  `--help` text and the `/arcdlc:plan` step need the worked example, or engineers will file the
  no-op case as a bug.
- **An uncommitted change set sits in the same files.** At interview time the working tree holds
  unreleased writing-style work touching `AGENTS.md`, `README.md`, all nine `SKILL.md` files, the
  plugin manifests and `.github/workflows/ci.yml`. This initiative edits `AGENTS.md`,
  `skills/plan/SKILL.md`, `skills/execute/SKILL.md` and both manifests too. Land that change set
  first, then execute this plan on top. Never run two agents across these files at once.
- **Reordering while a task is `TAKEN`.** Nothing stops it, by H4. An agent holding a `TAKEN` task
  keeps working on the block it already read, so the risk is confusion rather than corruption, but
  the `/arcdlc:plan` step should say to reorder between runs.

### 🔴 Open questions

- Should `/arcdlc:examinate` place its gap-derived tasks by order instead of appending them to the
  end of `plan.md`? It appends today. Now that a reorder command exists, appending may be the wrong
  default, but this was not decided in the interview.
- Should a later initiative add a `- DEPENDS:` key plus `arctool order --auto`, making the order
  checkable rather than merely applicable? H2 leaves the door open and decides nothing.
- Should `arctool validate` warn when a plan's block spacing is irregular, given that `order` and
  `archive` both normalise it silently?

## References

- [ADR-0013 — arctool order permutes slots](../../adr/0013-order-is-a-slot-permutation.md)
- [ADR-0014 — Exit 5 means self-validation failed](../../adr/0014-exit-5-means-self-validation-failed.md)
- `skills/plan/references/plan-format.md`, the plan contract this command mutates but does not change.
- `internal/plan/archive.go`, the existing whole-file rewrite this command copies.
- `internal/plan/mutate.go`, the existing byte-preserving mutation and its guardrails.
- `skills/source-map/source/Writing Style.md`, the writing standard for this document.
- `CONTEXT.md`, glossary for the terms used above.
