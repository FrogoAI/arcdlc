# Task Ordering — Plan

Task format: `skills/plan/references/plan-format.md`. Architecture: [aic.md](aic.md).

All paths are relative to the repo root. Run `go build ./... && go test ./... && gofmt -l . && go vet ./...`
before every commit: the repo must stay green.

## Risk Coverage

- **R1, a permutation bug destroys a plan** covered by ORD-1: the self-check runs inside `Reorder`, so
  no caller can skip it, and the acceptance criteria name the index-maths cases (one ID, all IDs,
  already-in-order, first and last swapped, single-block plan).
- **R2, the exit table lives in several places** covered by ORD-2 (the `usage` string in
  `cmd/arctool/main.go`) and ORD-3 (`AGENTS.md`). `README.md` has no exit table, so ORD-3 covers only
  its command list and the sentence that stops being true.
- **R3, inter-block content moves with the block above it** covered by ORD-1: a test pins the
  behaviour so it cannot change by accident.
- **R4, line endings** covered by ORD-1: a test pins that a CRLF plan keeps CRLF inside each block and
  gets LF between blocks, matching `Archive`.
- **R5, duplicate task IDs in a plan** covered by ORD-1 (the ambiguous-ID rejection) and ORD-2 (exit 3).
- **R6, "move to first" is the phrase people will use** covered by ORD-2 (worked example in `--help`)
  and ORD-4 (worked example in `/arcdlc:plan`).
- **R7, reordering while a task is TAKEN** covered by ORD-4: the `/arcdlc:plan` step says to reorder
  between runs, not during one.
- **R8, an uncommitted change set sits in the same files** accepted as a process mitigation, not a
  task: land the in-flight writing-style change set first, then execute this plan on top. Never run
  two agents across `AGENTS.md`, `README.md` and the `SKILL.md` files at once.
- **Q1, should `/arcdlc:examinate` place gap tasks by order instead of appending?** Accepted and
  deferred. It is a change to a different skill's behaviour and needs its own interview.
- **Q2, should a later initiative add `- DEPENDS:` plus `arctool order --auto`?** Accepted and
  deferred by [ADR-0013](../../adr/0013-order-is-a-slot-permutation.md). The command line designed
  here does not change if it is added.
- **Q3, should `arctool validate` warn on irregular block spacing?** Accepted and deferred. Both
  `order` and `archive` normalise spacing silently today, which is consistent, so nothing is broken
  while this stays open.

### ORD-1: Add slot-permutation reordering to `internal/plan`

- WHAT: Add `(*Plan).Reorder(ids []string)`, which returns plan bytes with the named tasks permuted among the positions they already occupy, self-checked before it returns.
- HOW:
  New file `internal/plan/order.go`, pure and standard-library only, in the style of `archive.go`.

  Error type, so the CLI can map a rejection to an exit code without string matching:

  `type ReorderErrKind uint8` with `const ( ReorderBadArgs ReorderErrKind = iota + 1; ReorderNotFound; ReorderSelfCheck )`.
  `ReorderBadArgs` maps to exit 2, `ReorderNotFound` to exit 3, `ReorderSelfCheck` to exit 5 (ORD-2 does the mapping).
  `type ReorderError struct { Kind ReorderErrKind; Msg string; IDs []string }` with a pointer receiver
  `func (e *ReorderError) Error() string` that renders `Msg`, then `": "` and `strings.Join(e.IDs, ", ")` when `IDs` is non-empty.

  Signature: `func (p *Plan) Reorder(ids []string) (out []byte, changed bool, err error)`.

  Validate in exactly this order, so the message an agent gets is deterministic:
  1. `len(ids) < 2` gives `ReorderBadArgs`, Msg `"order needs at least two task IDs"`, no IDs.
  2. An ID repeated in `ids` gives `ReorderBadArgs`, Msg `"task ID given more than once"`, IDs = the repeated IDs in first-seen order.
  3. `len(p.Tasks) == 0` gives `ReorderNotFound`, Msg `"plan has no task blocks"`, no IDs.
  4. Walk `ids` once, calling `p.ByID`. Collect IDs with zero matches into `unknown` and IDs with more
     than one match into `ambiguous`. If `unknown` is non-empty, return `ReorderNotFound` with Msg
     `"no such task"` and IDs = all of `unknown`. Otherwise, if `ambiguous` is non-empty, return
     `ReorderNotFound` with Msg `"ambiguous task ID (duplicate in the plan)"` and IDs = all of `ambiguous`.
     Report the whole set, never just the first: the caller is an agent and one message is one retry.

  Permute: build `slots`, the index in `p.Tasks` of each named task, then `sort.Ints(slots)`. Start
  `newOrder` as `[0,1,…,len(p.Tasks)-1]`, then for each `i`, set `newOrder[slots[i]]` to the index of
  `ids[i]`. Set `changed` false when `newOrder` is the identity, and in that case return `(nil, false, nil)`.

  Render, following `Archive`'s layout rules and nothing else:
  - Preamble is `p.Bytes[:p.Tasks[0].RawStart]` copied verbatim, with trailing newlines trimmed, then one `"\n"`.
    Do **not** call `splitPreamble` and do not regenerate the ledger: `order` must leave the
    "Completed (archived to …)" block exactly as it found it.
  - Then, for each index in `newOrder`: write `"\n"` first only when the buffer is already non-empty
    (a plan whose first byte starts a `###` block must not gain a leading blank line), then
    `bytes.TrimRight(p.Raw(&p.Tasks[idx]), "\n")`, then `"\n"`. That yields exactly one blank line
    between blocks and a single trailing newline, the same shape `Archive` produces.
  - Join with `"\n"` even when `p.CRLF` is true. `Archive` already does this, and the two must agree.

  Self-check before returning, in an unexported `verifyReorder(orig *Plan, out []byte, want []string) error`
  that `Reorder` always calls. Any failure becomes `ReorderError{Kind: ReorderSelfCheck}` and `out` is
  discarded (return `nil` bytes). Assert four things against a re-parse of `out`:
  1. the block count is unchanged;
  2. the sequence of task IDs equals `want`, the ID sequence `newOrder` asked for;
  3. `StatusCounts()` is equal for every status;
  4. the multiset of block bodies is unchanged. Build a sorted `[]string` of
     `string(bytes.TrimRight(raw, "\n"))` for the original and for the re-parse, and compare. Use a
     multiset, not a map keyed by ID, because a plan may hold duplicate IDs among the tasks nobody named.

  Verification lives inside `Reorder` rather than beside it, unlike `VerifyArchive`. `archive` writes two
  files in a fixed order so its caller needs the intermediate result; `order` writes one file, so nothing
  needs the bytes before the check, and folding it in makes the check impossible to skip.

  Out of scope: any CLI wiring (ORD-2), any change to `plan-format.md`, and any dependency field or
  ordering validation (ADR-0013 rules both out).
- WHERE:
  Layer `internal/plan`: `internal/plan/order.go`.
  Tests: `internal/plan/order_test.go`.
- WHY: This is the whole feature. Without a self-checked permutation the command cannot exist, and a permutation bug destroys a plan the engineer cannot reconstruct.
- Acceptance:
  - GIVEN a plan with blocks `T1 T2 T3 T4 T5 T6` WHEN `Reorder([]string{"T5","T2","T3"})` runs THEN `changed` is true and re-parsing `out` yields the ID order `T1 T5 T2 T4 T3 T6`.
  - GIVEN a plan `T1 T2 T3` WHEN `Reorder([]string{"T3","T1","T2"})` runs THEN the ID order is `T3 T1 T2`, and WHEN `Reorder([]string{"T1","T3"})` runs THEN `changed` is false and `out` is nil.
  - GIVEN any of these calls WHEN the argument list has one ID, or repeats an ID THEN the returned error is a `*ReorderError` with `Kind == ReorderBadArgs`; and WHEN it names an ID absent from the plan, or an ID that matches two blocks, or the plan has no blocks THEN `Kind == ReorderNotFound`.
  - GIVEN an argument list with three unknown IDs WHEN `Reorder` runs THEN the returned `*ReorderError` has all three in `IDs`, and `Error()` names all three.
  - GIVEN a plan whose preamble holds intro text and a "Completed (archived to …)" ledger WHEN a reorder runs THEN the bytes before the first `###` are byte-identical to the input.
  - GIVEN a block followed by a stray note line before the next `###` WHEN that block moves THEN the note moves with it, and the test asserts this pinned behaviour explicitly.
  - GIVEN a plan written with CRLF line endings WHEN a reorder runs THEN each block keeps its CRLF bytes and the blocks are separated by LF, matching `Archive`.
  - GIVEN a plan whose first byte begins a `###` block WHEN a reorder runs THEN `out` does not begin with a blank line.
  - GIVEN the new code WHEN `go test ./internal/plan/...` runs THEN `order_test.go` covers every case above and passes.
- References: `docs/aics/ordering/aic.md`, `docs/adr/0013-order-is-a-slot-permutation.md`, `docs/adr/0014-exit-5-means-self-validation-failed.md`, `skills/plan/references/plan-format.md`.
- Status: DONE.

### ORD-2: Add the `arctool order` command

- WHAT: Wire `(*Plan).Reorder` into a new `arctool order` command with `--dry-run`, the standard selection flags, the exit-code mapping, and the printed result.
- HOW:
  Edit `cmd/arctool/main.go`. Add `"errors"` to the import block (it is not imported today).

  Dispatch: add `case "order": os.Exit(cmdOrder(os.Args[2:]))` next to the existing `take/done/block/todo` case.

  `func cmdOrder(args []string) int`, modelled on `cmdMutate` for argument handling and on `cmdArchive`
  for the dry-run shape:
  - `flags, pos := splitArgs(args, map[string]bool{"plan": true, "aic": true})`. `splitArgs` is required
    so `arctool order T3 T1 --aic ordering` works; the stdlib flag package rejects positionals before flags.
  - Flags: `--plan`, `--aic` (same help strings as `cmdMutate`), and `--dry-run` "print the resulting order without writing".
  - `if len(pos) < 2` print `usage: arctool order <id> <id> [<id>…] [--dry-run] [--aic SLUG | --plan PATH]` to stderr and return 2.
  - `resolvePlan(aicsDir, *planFlag, *aicFlag)` then `loadPlan(planPath)`, propagating their codes unchanged.
  - `out, changed, err := p.Reorder(pos)`. On error, `var rerr *plan.ReorderError; errors.As(err, &rerr)`,
    print `arctool: %v` to stderr, then return 2 for `plan.ReorderBadArgs`, 3 for `plan.ReorderNotFound`,
    and for `plan.ReorderSelfCheck` print a second stderr line `arctool: nothing written` and return 5.
    A non-`*ReorderError` error is impossible today; return 1 for it rather than ignoring it.
  - `if !changed` print `already in that order` on stdout, print the current order (below), return 0. Write nothing.
  - `if *dryRun` print `would reorder <planPath>:`, then the new order, return 0. The self-check has already
    run inside `Reorder`, so a dry run cannot show an order that a real run would refuse.
  - Otherwise `atomicWrite(planPath, out)`; on error print `arctool: write %s: %v` and return 4. Then print
    `reordered <n> task(s) in <planPath>` where n is the count of blocks whose position changed, then the new order.

  Printing helper `func printOrder(p *plan.Plan)`: one line per block,
  `fmt.Printf("%3d  %-16s %-8s %s\n", i+1, t.ID, st, t.Title)`, with `st` falling back to `(none)` for an
  empty status exactly as `cmdList` does. Call it with `plan.Parse(out)` for the new order, and with the
  loaded plan for the no-op case.

  Update the `usage` const in the same file:
  - Add, in the command block after the `take|done|todo` and `block` lines:
    `  arctool order  <id> <id> [<id>…] [--dry-run] [--aic SLUG | --plan PATH]   re-order task blocks`
    followed by an indented note giving the worked example, because "move to first" is the phrase that
    misleads: `slot permutation: named tasks swap among the positions they already hold`, and
    `T1 T2 T3 + "order T3 T1 T2" -> T3 T1 T2; a task you do not name never moves`.
  - Change the exit-code line `  5  archive self-validation failed` to
    `  5  self-validation failed (nothing written)`, per ADR-0014.

  Also change the existing archive failure message at the `VerifyArchive` call site from
  `archive self-validation failed` to `archive: self-validation failed` so both commands read the same way.

  Out of scope: `AGENTS.md` and `README.md` (ORD-3), the skills (ORD-4), the version bumps (ORD-5).
- WHERE:
  Layer `cmd/arctool`: `cmd/arctool/main.go` (imports, dispatch, `cmdOrder`, `printOrder`, the `usage` const, the archive failure message).
  Tests: `cmd/arctool/main_test.go`.
- WHY: The permutation is unreachable without a command, and the exit-code mapping is the part skills branch on.
- Acceptance:
  - GIVEN a temp plan file with blocks `T1 T2 T3` WHEN `cmdOrder([]string{"T3","T1","T2","--plan",path})` runs THEN it returns 0 and the file on disk parses to the ID order `T3 T1 T2`.
  - GIVEN the same file WHEN `cmdOrder([]string{"T3","T1","--dry-run","--plan",path})` runs THEN it returns 0 and the file's bytes are unchanged.
  - GIVEN the same file WHEN the call names one ID, or repeats an ID, THEN the return code is 2; WHEN it names an ID absent from the plan THEN the return code is 3; and in both cases the file's bytes are unchanged.
  - GIVEN a plan file with no `###` blocks WHEN `cmdOrder` runs with two IDs THEN the return code is 3.
  - GIVEN positionals written before flags (`T3 T1 --plan path`) WHEN `cmdOrder` runs THEN the plan is resolved and the command returns 0, proving `splitArgs` is used.
  - GIVEN the built binary WHEN `arctool help` runs THEN the output contains the `arctool order` line, the worked example `T3 T1 T2`, and the exit line `5  self-validation failed (nothing written)`, and does not contain `5  archive self-validation failed`.
  - GIVEN the new code WHEN `go test ./cmd/...` runs THEN the tests above pass and `go vet ./...` is clean.
- References: `docs/aics/ordering/aic.md`, `docs/adr/0013-order-is-a-slot-permutation.md`, `docs/adr/0014-exit-5-means-self-validation-failed.md`.
- Status: DONE.

### ORD-3: Update AGENTS.md and README.md for the new command

- WHAT: Record `arctool order` and the widened meaning of exit code 5 in the two repo-level documents, and fix the README sentence that the new command makes false.
- HOW:
  `AGENTS.md`, in `## Conventions`, change the exit-code line from
  `5 archive self-validation` to `5 self-validation` and keep the rest of that sentence
  (`skills key off them; do not renumber`) exactly as it is.

  `AGENTS.md`, in `## Hard rules`, extend the existing "Status mutations stay byte-preserving and atomic"
  rule so it also covers the second whole-file rewrite. Keep the existing two clauses and add: `order`
  permutes whole blocks, re-parses its own output before writing, and writes nothing on a failed check.
  Do not create a new rule bullet; the invariant is the same one.

  `README.md`, in the `arctool` command block (the fenced `bash` list that ends with the two `archive`
  lines), add two lines:
  `arctool order AIC-3 AIC-1 AIC-2   # re-order task blocks (slot permutation; unnamed tasks never move)`
  and
  `arctool order AIC-3 AIC-1 --dry-run   # preview the new order without writing`.

  `README.md`, the sentence directly under that block currently reads that all mutations are
  "byte-preserving, atomic rewrites of the single status line". That stops being true. Rewrite it to say
  that status changes rewrite only the one status line, that `archive` and `order` rewrite the whole file
  after re-parsing their own output, and that every write is atomic. Keep the closing promise that
  `arctool` never reformats a plan beyond the block spacing those two commands normalise.

  Follow `skills/source-map/source/Writing Style.md`: no long dashes in the prose you add.

  Do not run `arctool sync`: the initiative registry blocks are already current and this task changes
  nothing they read.

  Out of scope: the skills (ORD-4), the version bumps (ORD-5), any code.
- WHERE:
  Docs: `AGENTS.md` (the exit-code convention line, the status-mutation hard rule), `README.md` (the arctool command block and the sentence under it).
- WHY: Exit codes are interface, and a README sentence that promises status-line-only rewrites would be wrong the moment `order` ships.
- Acceptance:
  - GIVEN the edited `AGENTS.md` WHEN `grep -n "archive self-validation" AGENTS.md` runs THEN it finds nothing, and the exit-code line reads `5 self-validation`.
  - GIVEN the edited `AGENTS.md` WHEN the status-mutation hard rule is read THEN it names `order` as a whole-file rewrite that re-parses before writing and writes nothing on a failed check.
  - GIVEN the edited `README.md` WHEN `grep -n "arctool order" README.md` runs THEN it finds both the plain and the `--dry-run` example.
  - GIVEN the edited `README.md` WHEN `grep -n "single status line" README.md` runs THEN it finds nothing.
  - GIVEN both edited files WHEN `grep -nE "—|–" AGENTS.md README.md` runs THEN every remaining hit is a pre-existing line or a generated initiative-registry line, and none is in the prose this task added.
  - GIVEN the repo WHEN `arctool sync --check` runs THEN it exits 0.
- References: `docs/aics/ordering/aic.md`, `docs/adr/0014-exit-5-means-self-validation-failed.md`, `skills/source-map/source/Writing Style.md`.
- Status: DONE.

### ORD-4: Teach `/arcdlc:plan` and `/arcdlc:execute` about ordering

- WHAT: Add the `arctool order` step to `/arcdlc:plan` and the never-applied reorder proposal to `/arcdlc:execute`.
- HOW:
  `skills/plan/SKILL.md`, in `## Step 3 — Write and validate`, beside the existing `arctool validate`
  bullet, add an ordering bullet. It must obey the "arctool is always optional" rule: it hangs off the
  same single `command -v arctool` probe that is already there, and it states the manual fallback, which
  is to move the `###` blocks by hand and re-check that each block kept its `- Status:` line.

  The bullet says three things:
  1. The command is `arctool order <ID> <ID> … --aic <slug>`, and `--dry-run` previews it.
  2. The semantics, with the worked example, because "move to first" misleads: named tasks swap among
     the positions they already hold, so on `T1 T2 T3` the call `arctool order T3 T1` gives `T3 T2 T1`,
     while `arctool order T3 T1 T2` gives `T3 T1 T2`. A task you do not name never moves.
  3. Reorder between `/arcdlc:execute` runs, not during one. Nothing stops a reorder while a task is
     `TAKEN`, and an agent mid-task keeps working from the block it already read.

  Authoring rule 5 in `plan-format.md` stays the source of the obligation. Do not restate the rule here
  and do not edit `plan-format.md` (ADR-0013 and the AIC both keep it out of scope).

  `skills/execute/SKILL.md`, in `## Per-task contract` step 7, extend the existing sentence about a task
  that cannot be completed. Add the one case that is a reorder: the executor is implementing a task and
  finds it needs something a task *below* it will build. Then it blocks as it already does, with a reason
  naming the other task, for example `arctool block <id> -m "needs <OTHER-ID>, which is below it"`, and
  adds one line to its report giving the exact `arctool order …` command that would fix the plan. State
  plainly that it never runs that command: the engineer decides. Keep step 7's existing stop-the-run rule.

  Do not add a pre-flight scan anywhere. Orchestrator mode's "never read source files or diffs, only plan
  state, subagent reports, and commit subjects" rule stands, and a scan of every block's prose would break it.

  Both files must keep their dual path references intact and stay install-agnostic. Do not touch the
  `## Talk simple, write like a human` block in either file. Follow
  `skills/source-map/source/Writing Style.md` for the prose you add.

  Out of scope: `skills/archive/SKILL.md` (its exit-5 wording is already generic), every other skill,
  `plan-format.md`, and the version bumps (ORD-5).
- WHERE:
  Skills: `skills/plan/SKILL.md` (Step 3 validate section), `skills/execute/SKILL.md` (per-task contract step 7).
- WHY: A command no skill names will not be used, and without the "never applies" rule an executor could start rewriting its own queue mid-run.
- Acceptance:
  - GIVEN the edited `skills/plan/SKILL.md` WHEN `grep -n "arctool order" skills/plan/SKILL.md` runs THEN it finds the command, and the surrounding bullet contains both worked examples (`T3 T2 T1` and `T3 T1 T2`) and the manual fallback.
  - GIVEN the edited `skills/plan/SKILL.md` WHEN the ordering bullet is read THEN it says to reorder between `/arcdlc:execute` runs, not during one.
  - GIVEN the edited `skills/execute/SKILL.md` WHEN step 7 is read THEN it names the below-it dependency case, keeps `arctool block <id> -m "<reason>"` as the action, requires the report to carry the exact `arctool order …` command, and states that the executor never runs it.
  - GIVEN the edited `skills/execute/SKILL.md` WHEN `grep -n "pre-flight\|scan every block" skills/execute/SKILL.md` runs THEN it finds nothing.
  - GIVEN both files WHEN `grep -c '^## Talk simple, write like a human$' skills/plan/SKILL.md skills/execute/SKILL.md` runs THEN each reports 1, so the CI skill-layout check still passes.
  - GIVEN both files WHEN `grep -n "arcdlc-plan/references\|../plan/references" skills/execute/SKILL.md` runs THEN both dual path forms are still present.
- References: `docs/aics/ordering/aic.md`, `docs/adr/0013-order-is-a-slot-permutation.md`, `skills/plan/references/plan-format.md`, `skills/source-map/source/Writing Style.md`.
- Status: DONE.

### ORD-5: Bump arctool to 0.11.0 and the plugin bundle to 0.13.0

- WHAT: Bump the CLI version and both plugin manifests, in lockstep, as the last change of the initiative.
- HOW:
  `cmd/arctool/main.go`: change `const version = "0.10.0"` to `"0.11.0"`. The CLI gained a command, so this is a minor bump.

  `.claude-plugin/plugin.json` and `.antigravity-plugin/plugin.json`: change `"version": "0.12.0"` to
  `"0.13.0"` in both. CI asserts the two manifest versions are equal, so they must move together.

  Run this task last, after ORD-1 through ORD-4 are `DONE`, so the version names a bundle that already
  contains everything.

  Do not tag a release: releases are cut by pushing a `v*` tag, which is a human decision outside this plan.

  Out of scope: `.claude-plugin/marketplace.json`, which carries no version field, and any changelog
  file, of which this repo has none.
- WHERE:
  Layer `cmd/arctool`: `cmd/arctool/main.go` (the `version` const).
  Manifests: `.claude-plugin/plugin.json`, `.antigravity-plugin/plugin.json`.
- WHY: AGENTS.md requires bumping whichever component changed, and CI fails if the two plugin manifests drift apart.
- Acceptance:
  - GIVEN the edited source WHEN `go run ./cmd/arctool version` runs THEN it prints `arctool 0.11.0`.
  - GIVEN both manifests WHEN `test "$(jq -r .version .antigravity-plugin/plugin.json)" = "$(jq -r .version .claude-plugin/plugin.json)"` runs THEN it succeeds, and both read `0.13.0`.
  - GIVEN the repo WHEN `go build ./... && go test ./... && gofmt -l . && go vet ./...` runs THEN the build and tests pass and `gofmt -l .` prints nothing.
- References: `docs/aics/ordering/aic.md`.
- Status: TODO.
