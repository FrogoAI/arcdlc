# Plan: aic-tracking

<!-- arcdlc:source docs/aics/aic-tracking/aic.md sha256:7674bb1f10f5bd7a3d8aae61ae367ffdc503050bb887fa877b12db25759e88c6 -->

Format: [`skills/plan/references/plan-format.md`](../../../skills/plan/references/plan-format.md).
Source: [`aic.md`](aic.md), decided 2026-09-25. Terms: [`CONTEXT.md`](../../../CONTEXT.md).
Decision: [ADR-0028](../../adr/0028-a-finished-initiative-is-closed-in-place.md).

## Risk Coverage

- R1 (links broken by a deletion): accepted by the engineer on 2026-09-25; close and remove do not handle links (AIC H5). TRK-4 makes `remove`'s confirmation say that references may be left pointing at nothing.
- R2 (a request to finish an initiative lands on `remove`): TRK-3 puts the finishing phrases in `close`'s description and adds its routing rows; TRK-4 narrows `remove`'s description and re-points its routing rows. Process mitigation: the engineer re-runs `docs/routing-checks.md` by hand after the release, as `AGENTS.md` requires for every description change.
- R3 (marker text lost when `comments.md` is deleted): TRK-3, the `## Not done` section of `CLOSED.md` carries each unfinished comment task's `- Marker:` and `- WHERE:` lines.
- R4 (a closed initiative is changed anyway): TRK-5 adds the closed check to `aic`, `plan`, `examinate` and `assist`.
- Open questions: none in the AIC.

### TRK-1: Add the `arctool status` command

- WHAT: Add a read-only `arctool status [--json]` command that lists every initiative folder under `docs/aics/` with its phase, task counts, and, for a closed one, its close date and outcome.
- HOW:
  New file `cmd/arctool/status.go`. Add `case "status": os.Exit(cmdStatus(os.Args[2:]))` to the switch in `main` (`cmd/arctool/main.go`), and one usage line after the `list` line: `arctool status [--json]   every initiative under docs/aics/ with its phase: designing, in progress, ready to close, closed`.
  `cmdStatus(args []string) int` parses `--json` with a `flag.FlagSet` named `status`; a parse error or any positional argument returns 2. It calls `runStatus(aicsDir, asJSON, os.Stdout, os.Stderr)`.
  `runStatus(dir string, asJSON bool, out, errw io.Writer) int` reads `os.ReadDir(dir)`. A missing `dir` is not an error: it prints the header (text) or `{"initiatives":[]}` (JSON) and returns 0. It keeps entries that are directories and whose name does not start with `.`, in `os.ReadDir` order (already sorted by name).
  For each folder build a `statusRow{Slug, Phase string; Counts statusCounts; Closed, Outcome string}` where `statusCounts{TODO, TAKEN, BLOCKED, DONE int}`, with JSON tags `slug`, `phase`, `counts`, `closed`, `outcome` and `TODO`, `TAKEN`, `BLOCKED`, `DONE`.
  Counts: when `plan.md` exists, parse it with `plan.Parse` and add `StatusCounts()` values; when `plan-archive.md` exists, parse it the same way and add only its `DONE` count to `DONE`. A read error on either file prints `arctool: cannot read <path>: <err>` to `errw` and returns 4.
  Phase, first match wins: `closed` when `CLOSED.md` exists; `designing` when `plan.md` does not exist; `in progress` when TODO + TAKEN + BLOCKED > 0; else `ready to close`.
  Closed and Outcome: only for `closed`. Read `CLOSED.md` line by line; the first line that, after `strings.TrimSpace`, starts with `- Closed:` gives `Closed` (the rest, trimmed), and the first that starts with `- Outcome:` gives `Outcome`. A missing line leaves the field empty. For an open initiative both stay empty.
  Text output: a header then one line per row, both formatted with `%-28s %-14s %5s %5s %7s %5s  %-10s %s` (header cells `SLUG PHASE TODO TAKEN BLOCKED DONE CLOSED OUTCOME`; row cells use `strconv.Itoa` for counts and `-` for an empty Closed or Outcome). JSON output: `{"initiatives":[...]}` via `json.MarshalIndent` with two-space indent, and an empty list encodes as `[]`, never `null`.
  Returns 0. Out of scope: `sync` (TRK-2), any `--aic` flag, any write.
- WHERE:
  Layer `cmd`: `cmd/arctool/status.go`, `cmd/arctool/main.go`.
  Tests: `cmd/arctool/status_test.go`.
- WHY: The engineer needs one place that shows which initiatives are designing, in progress, ready to close or closed, read from the files with no index to drift (AIC H8).
- Acceptance:
  - GIVEN a temp `docs/aics/` holding `a/` (only `aic.md`), `b/` (`plan.md` with one `TODO` and one `DONE` task), `c/` (`plan.md` with two `DONE` tasks and `plan-archive.md` with three `DONE` tasks) and `d/` (`aic.md` and a `CLOSED.md` containing `- Closed: 2026-09-25` and `- Outcome: stopped`) WHEN `runStatus` runs with `asJSON` true THEN the rows are, in order, `a designing`, `b in progress` (TODO 1, DONE 1), `c ready to close` (DONE 5) and `d closed` with `closed` `2026-09-25` and `outcome` `stopped`; checked by `TestRunStatusPhases` in `cmd/arctool/status_test.go`.
  - GIVEN a `dir` that does not exist WHEN `runStatus` runs with `asJSON` true THEN it returns 0 and prints `"initiatives": []`; checked by `TestRunStatusMissingDir`.
  - GIVEN the same fixture WHEN `runStatus` runs with `asJSON` false THEN the first line starts with `SLUG` and the `d` line contains `2026-09-25` and `stopped`; checked by `TestRunStatusText`.
  - GIVEN the change WHEN `go test ./cmd/arctool/...`, `go vet ./...` and `gofmt -l .` run THEN tests pass, vet is clean and gofmt prints nothing.
- References: `docs/aics/aic-tracking/aic.md`, `cmd/arctool/main.go`, `internal/plan/query.go`.
- Status: DONE.

### TRK-2: Leave closed initiatives out of the registry and count them in one line

- WHAT: Make `arctool sync` skip every initiative folder holding `CLOSED.md` and end the registry block with one line counting those folders, then bump `arctool` to 0.21.0.
- HOW:
  In `cmd/arctool/main.go`, change `scanInitiatives(dir string) []registry.Initiative` to `scanInitiatives(dir string) (inits []registry.Initiative, closed int)`: a folder that holds a file named `CLOSED.md` is counted in `closed` and never loaded; every other folder behaves as today. Update its one caller, `runSync`, to pass `closed` on.
  In `internal/registry/registry.go`, add a `closed int` parameter as the last parameter of `Render`, `Preview` and `WriteFile`. `Render` builds the body exactly as today (bullets, or `_none_` when there are none), then, when `closed > 0`, appends `"\n\n"` plus the count line. The count line is ``1 closed initiative: run `arctool status`, or see `CLOSED.md` in each folder.`` when `closed == 1`, and ``N closed initiatives: run `arctool status`, or see `CLOSED.md` in each folder.`` otherwise, with N the number. `Preview` passes `closed` to `Render`; `WriteFile` passes it to `Preview`. Update every existing call site and test to pass `0` where no closed count applies. `Rebase`, `Splice` and `section` keep their signatures.
  In the usage text, add one line under the `sync` entry: `folders holding CLOSED.md are left out and counted in one line`. Set `const version` to `"0.21.0"`.
  Out of scope: `arctool status` (TRK-1), any change to `Load`, `findArchDoc` or `resolvePlan`.
- WHERE:
  Layer `cmd`: `cmd/arctool/main.go`.
  Layer `registry`: `internal/registry/registry.go`.
  Tests: `internal/registry/registry_test.go`, `cmd/arctool/main_test.go`.
- WHY: The registry is loaded into every agent session, so it must not list finished work as active or grow with every close (AIC H9).
- Acceptance:
  - GIVEN `Render(nil, 0)` WHEN called THEN it returns `_none_`; GIVEN `Render(nil, 2)` THEN it returns `_none_` followed by a blank line and the line starting `2 closed initiatives:`; GIVEN one initiative and `closed == 1` THEN the result ends with the line starting `1 closed initiative:`; checked by `TestRender` and a new `TestRenderClosedCount` in `internal/registry/registry_test.go`.
  - GIVEN a temp root with `docs/aics/open/aic.md` and `docs/aics/done/aic.md` plus `docs/aics/done/CLOSED.md` WHEN `runSync` writes `AGENTS.md` THEN the block holds the `open` bullet, no `done` bullet, and the `1 closed initiative:` line; and `runSync` with `check` true right after returns 0; checked by a new `TestRunSyncSkipsClosed` in `cmd/arctool/main_test.go`.
  - GIVEN the change WHEN `go build ./...`, `go test ./...`, `go vet ./...` and `gofmt -l .` run THEN all pass and gofmt prints nothing; `go run ./cmd/arctool version` prints `arctool 0.21.0`.
- References: `docs/aics/aic-tracking/aic.md`, `docs/adr/0002-registry-sync-via-marker-blocks.md`, `internal/registry/registry.go`, `cmd/arctool/main.go`.
- Status: DONE.

### TRK-3: Create the `/arcdlc:close` skill and register it as the twelfth skill

- WHAT: Write `skills/close/SKILL.md`, which closes an initiative in place by writing `CLOSED.md` and deleting its five work files, and register `close` everywhere the bundle lists its skills.
- HOW:
  Frontmatter: `name: arcdlc-close`; `argument-hint: "<slug>"`; `description: Close a finished initiative in place: write a CLOSED.md note with the outcome, delete its plan, board stories and registers, and stop it taking new work, after showing what is unfinished and asking for a yes. Use when someone says we are done with X, wrap up X, close the initiative, mark it finished, or retire that initiative, or runs /arcdlc:close <slug>, or invokes arcdlc-close.`
  Body, in this order:
  1. `# ArcDLC Close (/arcdlc:close)`, then a short intro: close keeps the design in `docs/aics/<slug>/` where every link points, marks it with `CLOSED.md`, deletes `plan.md`, `plan-archive.md`, `plan-human.md`, `gap.md` and `comments.md`, and leaves the registry through `arctool sync`. A closed initiative is final: no skill changes it or deletes `CLOSED.md`; follow-up work is a new initiative. The chain line: `/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive` → `/arcdlc:close`, with `/arcdlc:remove` beside it for deleting a design for good.
  2. The `## Talk simple, write like a human` block and the `## Judge by the four virtues` block, each copied byte for byte from `skills/aic/SKILL.md` (from its heading line to the line before the next `## ` heading).
  3. `## Argument: initiative slug (required, first positional)`: no slug means stop and list the folders under `docs/aics/`; a missing folder means stop and list what exists. Never guess.
  4. `## In a workspace`: the workspace paragraph of `skills/remove/SKILL.md`, with these facts in place of remove's: the writes are `CLOSED.md`, the deleted files, `AGENTS.md` and `README.md`; the files are deleted with `git -C docs rm aics/<slug>/<file>`; everything lands in one hub commit `docs(<slug>): close the initiative` plus `#AI-assisted`, pushed at once, retried once after `git -C docs pull --ff-only`.
  5. `## Ask through the grilling protocol`: the paragraph of `skills/aic/SKILL.md` that begins `Prefer this bundle's own`, copied byte for byte to its end.
  6. `## Step 1 — Read the state`. When `CLOSED.md` exists and none of the five files remain: stop, say the initiative was closed on the date in its `- Closed:` line. When `CLOSED.md` exists and some of the five remain (an interrupted close): skip to Step 4, then Step 6. Otherwise read the task counts: `arctool list --aic <slug>` when `command -v arctool` succeeds, else count the `- Status:` lines of `plan.md` by value; the archived count is the number of `### ` blocks in `plan-archive.md`. No `plan.md` means zero tasks. Any `TAKEN` task: stop, change nothing, name each `TAKEN` task ID, because an agent is working on it.
  7. `## Step 2 — Show what closing does`: the counts; each `TODO` and `BLOCKED` task by ID, title and status, highlighted when there are any; the files that stay; the files of the five that exist and will be deleted. State that links to deleted files are not rewritten.
  8. `## Step 3 — Settle the outcome`: every task `DONE`, or no plan, means `delivered` and no question. Otherwise ask, one question, `partly delivered` or `stopped` (recommend `partly delivered` when any task is `DONE`, else `stopped`), then ask for one line of reason. Never write `delivered` while a task is `TODO` or `BLOCKED`.
  9. `## Step 4 — Confirm`: one question naming the folder, the files to delete, and that git history keeps them. Proceed only on an explicit yes; anything else stops with nothing changed.
  10. `## Step 5 — Write CLOSED.md` with this exact shape. Line 1 `# <title>`, where title is the first `# ` heading of the architecture document (`aic.md`, then `arc42.md`, `togaf.md`, `c4.md`, `tsc.md`, then the first other `.md` alphabetically that is not one of the five files), or the slug when there is none. Then a blank line and three list lines: `- Closed: YYYY-MM-DD` (today), `- Outcome: <delivered | partly delivered | stopped>`, `- Tasks: <total> (<done> done, <todo> todo, <blocked> blocked)` where total and done include the archived tasks. When the outcome is not `delivered`: `## Reason` with the one line; `## Not done` with one bullet per `TODO` or `BLOCKED` task, `` - `<ID>`: <title> (<STATUS>) ``, and when `comments.md` holds a `### ` block with the same ID, indented sub-bullets copying each of that block's `- Marker:` lines and its `- WHERE:` value as they stand. Then `## Left open` with the body under the architecture document's `Open questions` heading copied as it stands, or `None recorded.`. Then `## Documents` with one bullet `- [<file>](<file>)` per `.md` or `.html` file left in the folder other than `CLOSED.md`, or `Documents: none` when there are none. Close never writes a `## Superseded` section.
  11. `## Step 6 — Delete and refresh the registry`: delete each of the five files that exists with `git rm` (plain `rm` when untracked), then run `arctool sync`. Without `arctool`: delete the initiative's bullet from the `<!-- arcdlc:initiatives -->` blocks in `AGENTS.md` and `README.md` and write or update the closed-count line exactly as `arctool sync` does. Leave the changes staged; commit only when asked (in a workspace, per the workspace section).
  12. `## Step 7 — Report`: the `CLOSED.md` path, the outcome, the deleted files, the registry refresh, and that follow-up work is a new initiative.
  Registration: add `"close"` after `"assist"` in `Skills` in `internal/bundle/bundle.go`; add `close` after `assist` in `SUBSKILLS` in `install.sh` and in the four `for` lists of `.github/workflows/ci.yml` (skill layout, writing-style block, virtues block, installer smoke test); change `all eleven must match` in `ci.yml` to `all twelve must match`. Add three rows to the table in `docs/routing-checks.md`: `we are done with the payments initiative, close it` → `close`; `wrap up the checkout work` → `close`; `mark init as finished` → `close`.
  Out of scope: changes to `remove` (TRK-4), the closed check in other skills (TRK-5), `AGENTS.md`, `README.md`, `CONTEXT.md` and version bumps (TRK-7).
- WHERE:
  Skills: `skills/close/SKILL.md`.
  Bundle: `internal/bundle/bundle.go`, `install.sh`, `.github/workflows/ci.yml`.
  Docs: `docs/routing-checks.md`.
- WHY: Closing is the new end state that keeps a finished design and stops new work on it (AIC H1 to H7, H12, R2, R3).
- Acceptance:
  - GIVEN the new skill WHEN `go test ./internal/bundle/...` runs THEN it passes (the skill is on disk and in `Skills`, carries both shared blocks, and every path it names resolves in both layouts).
  - GIVEN `skills/close/SKILL.md` WHEN each shared block is extracted with the `awk` commands of the two block checks in `.github/workflows/ci.yml` and compared with `diff` against the same block of `skills/aic/SKILL.md` THEN both diffs are empty.
  - GIVEN the change WHEN `grep -c close install.sh` and `grep -c 'assist close examinate' .github/workflows/ci.yml` run THEN the first is at least 1 and the second is 4; `shellcheck install.sh` is clean.
  - GIVEN `skills/close/SKILL.md` WHEN read THEN it contains the headings `## Step 1` to `## Step 7`, the strings `- Closed:`, `- Outcome:`, `## Not done`, `## Left open`, `## Documents`, `plan-human.md` and `command -v arctool`; checked with `grep -c`.
- References: `docs/aics/aic-tracking/aic.md`, `skills/aic/SKILL.md`, `skills/remove/SKILL.md`, `skills/plan/references/plan-format.md`, `docs/adr/0028-a-finished-initiative-is-closed-in-place.md`.
- Status: DONE.

### TRK-4: Rewrite `/arcdlc:remove` to delete a design for good with a reference warning

- WHAT: Update `skills/remove/SKILL.md` so it deletes `docs/aics/<slug>/` for good, refuses only while a task is `TAKEN`, and warns that references to the removed design may be left pointing at nothing, and re-point its routing rows.
- HOW:
  New description: `Delete an initiative's folder for good, design included, after showing what is lost and an explicit confirmation. Use when someone says delete the design for good, purge the initiative, or throw away that initiative entirely, or runs /arcdlc:remove <slug>, or invokes arcdlc-remove.` Replace the intro: remove deletes the whole `docs/aics/<slug>/` folder, open or closed, and refreshes the registry; git history is the only copy afterwards; finishing an initiative is `/arcdlc:close`, which keeps the design.
  Step 1 (show what will be removed) keeps title, file list and task counts, and adds: when `CLOSED.md` exists, show its `- Closed:` and `- Outcome:` lines. Any `TAKEN` task: stop, change nothing, name the task IDs. `TODO` or `BLOCKED` tasks: warn loudly with the counts, as today.
  Step 2 (confirmation) must state, every time: the folder path; that the design leaves the working tree and stays only in git history; that other files may still reference this design and those references will point at nothing, because remove does not rewrite links; and that the engineer can cancel to handle the references first. Proceed only on an explicit yes.
  Step 3 (delete and clean the registry) keeps today's steps; `arctool sync` also refreshes the closed-count line.
  Keep the `## Talk simple`, `## Judge by the four virtues` and `## In a workspace` sections unchanged.
  In `docs/routing-checks.md`: change the expected skill of the row `we are done with cursor-support, retire it` to `close`; replace the row `delete the old plan folder` with `delete the payments design for good` → `remove`; add `purge the old spike initiative` → `remove`.
  Out of scope: any link search or rewrite (AIC H5), the `close` skill (TRK-3).
- WHERE:
  Skills: `skills/remove/SKILL.md`.
  Docs: `docs/routing-checks.md`.
- WHY: Finishing now goes through `close`; `remove` stays the deliberate way to delete a design, and the engineer must know references can be left behind (AIC H11, R1, R2).
- Acceptance:
  - GIVEN the new `skills/remove/SKILL.md` WHEN `grep -n "point at nothing" skills/remove/SKILL.md` and `grep -n "TAKEN" skills/remove/SKILL.md` run THEN both print at least one line.
  - GIVEN the file WHEN `grep -n "docs/history\|graveyard" skills/remove/SKILL.md` runs THEN it prints nothing about `docs/history`.
  - GIVEN `docs/routing-checks.md` WHEN `grep -n "retire it" docs/routing-checks.md` runs THEN the row ends in `` `close` |``, and `grep -c '`remove` |' docs/routing-checks.md` prints 2.
  - GIVEN the change WHEN `go test ./internal/bundle/...` runs THEN it passes.
- References: `docs/aics/aic-tracking/aic.md`, `skills/remove/SKILL.md`, `docs/adr/0028-a-finished-initiative-is-closed-in-place.md`.
- Status: DONE.

### TRK-5: Add the closed-initiative check to `aic`, `plan`, `examinate` and `assist`

- WHAT: Add one identical paragraph to four skills so each stops on an initiative folder that holds `CLOSED.md`.
- HOW:
  The paragraph, identical in all four files, is the text between the two lines of three dashes below, with its inline code spans as written:
  ---
  **A closed initiative is final.** If `docs/aics/<slug>/CLOSED.md` exists, stop and change nothing: say that `<slug>` was closed on the date in its `- Closed:` line, and that follow-up work is a new initiative with its own slug whose design links `docs/aics/<slug>/` as its starting point. No skill edits or deletes a closed initiative's files. The one exception is the `## Superseded` line that `/arcdlc:aic` appends to its `CLOSED.md` when another initiative's design reverses it.
  ---
  Placement: `skills/aic/SKILL.md` at the end of the section `## Argument: initiative slug (required, first positional)`; `skills/plan/SKILL.md` at the end of `## Initiative selection`; `skills/examinate/SKILL.md` and `skills/assist/SKILL.md` at the end of `## Initiative selection`. The check runs before any file is written, including the fresh-folder case of `examinate` and `assist`.
  Out of scope: the design-time read and the supersede rule in `aic` (TRK-6), `execute`, `archive` and `plan-human`, which already stop when `plan.md` is gone.
- WHERE:
  Skills: `skills/aic/SKILL.md`, `skills/plan/SKILL.md`, `skills/examinate/SKILL.md`, `skills/assist/SKILL.md`.
- WHY: A closed initiative is final, and these four are the skills that can write into an initiative folder (AIC H10, R4).
- Acceptance:
  - GIVEN the four files WHEN `grep -c "A closed initiative is final" skills/aic/SKILL.md skills/plan/SKILL.md skills/examinate/SKILL.md skills/assist/SKILL.md` runs THEN each prints 1.
  - GIVEN the four files WHEN the line holding `A closed initiative is final` is extracted from each with `grep` and compared with `diff` THEN all four are identical.
  - GIVEN the change WHEN `go test ./internal/bundle/...` runs THEN it passes.
- References: `docs/aics/aic-tracking/aic.md`, `docs/adr/0028-a-finished-initiative-is-closed-in-place.md`.
- Status: TODO.

### TRK-6: Teach `/arcdlc:aic` to read closed designs and record a reversal in their `CLOSED.md`

- WHAT: Add the closed-design read to `/arcdlc:aic` Step 1, the rule that a round reversing a closed design appends a `## Superseded` line to its `CLOSED.md`, and `close` to the chain line.
- HOW:
  In `skills/aic/SKILL.md`, Step 1 list, add this bullet after the `docs/aics/` bullet, with its inline code spans as written: - `docs/aics/*/CLOSED.md`: closed initiatives. Read each title and `- Outcome:` line and open the designs related to this one; they are references and are never edited.
  In Step 3, after the paragraph list that ends with the `CONTEXT.md` bullet of "A reversed decision is superseded, never deleted", add a bullet `Recorded in a closed initiative's design` with these rules: append one line to that initiative's `CLOSED.md` under a `## Superseded` heading, creating the heading at the end of the file when missing; the line is `- Superseded in part by [<title>](../<name>/aic.md) on <YYYY-MM-DD>: <which sections>.`, or `- Superseded by [<title>](../<name>/aic.md) on <YYYY-MM-DD>.` when the whole design is replaced, where `<name>` and `<title>` are this initiative's slug and H1; never touch the closed design documents, and never start an `/arcdlc:aic` run on the closed initiative. When the reversed design belongs to an initiative that is still open (no `CLOSED.md`), write nothing into it: raise the disagreement as one grilling question, because two open designs that disagree are the engineer's decision.
  Chain line near the top: `/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive` → `/arcdlc:close`, with `/arcdlc:remove` named beside it for deleting a design for good.
  Out of scope: the closed check paragraph (TRK-5), other skills.
- WHERE:
  Skills: `skills/aic/SKILL.md`.
- WHY: Closed designs are the references new designs start from, and a reader of an old design must learn when a later one changed it without the old design being edited (AIC H9, H14).
- Acceptance:
  - GIVEN `skills/aic/SKILL.md` WHEN `grep -n "CLOSED.md" skills/aic/SKILL.md` runs THEN it shows the Step 1 bullet, the `## Superseded` rule, and the TRK-5 paragraph.
  - GIVEN the file WHEN `grep -n "Superseded in part by" skills/aic/SKILL.md` and `grep -n "/arcdlc:close" skills/aic/SKILL.md` run THEN both print at least one line.
  - GIVEN the change WHEN `go test ./internal/bundle/...` runs THEN it passes.
- References: `docs/aics/aic-tracking/aic.md`, `skills/aic/SKILL.md`.
- Status: TODO.

### TRK-7: Update `AGENTS.md`, `README.md`, `CONTEXT.md`, the `init` chain line and the bundle version

- WHAT: Bring the repository guides, the glossary, the `init` chain line and both plugin manifests in line with `close`, the in-place close, `arctool status`, and twelve skills.
- HOW:
  `AGENTS.md` (edit only outside the `<!-- arcdlc:initiatives -->` markers): in "What this repo is", add `close` to the workflow list and change `Eleven skills` to `Twelve skills`; change every other `eleven` meaning the skill count to `twelve`; in the hard rule on the one-question interview, replace the sentence that begins "Six skills restate this in one shared paragraph" so it names eight skills (`aic`, `policy`, `examinate`, `assist`, `plan-human`, `execute`, `init`, `close`) and says to change it in all nine together (the eight plus `grilling`); in "Initiatives are folders", add `plan-human.md` if missing and add that a folder holding `CLOSED.md` is closed and final; in "The initiative registry is generated", replace the `/arcdlc:remove` sentence with one that says: `/arcdlc:close <slug>` writes `CLOSED.md` and deletes the plan files; `sync` leaves closed initiatives out of the registry and counts them in one line; `/arcdlc:remove <slug>` deletes a folder after an explicit confirmation; `arctool` itself deletes nothing. Link ADR-0028 (`docs/adr/0028-a-finished-initiative-is-closed-in-place.md`) beside the ADR-0002 and ADR-0003 links.
  `README.md` (outside the markers): add a table row after the `/arcdlc:archive` row: command `/arcdlc:close <slug>`; description "Close a finished initiative in place: a `CLOSED.md` note with the outcome, the plan files deleted, no new work after."; output `docs/aics/<slug>/CLOSED.md`, refreshed registry; change the `/arcdlc:remove` row to deleting a design for good, after a confirmation that warns about references left behind; in "Initiatives live in folders", replace the clause "and `/arcdlc:remove <slug>` retires a finished one" with one saying that `/arcdlc:close <slug>` closes a finished one in place and `arctool status` lists every initiative with its phase; in the layout tree, `eleven skills` becomes `twelve skills` and `close/` joins the `policy/  execute/  assist/  remove/  archive/` line; in the examples block, add `/arcdlc:close payments              # close the finished initiative, keep its design` above the `remove` line and change the `remove` comment to `# delete the initiative for good (after confirming)`.
  `CONTEXT.md`: in "Pipeline skills", add `close` to the lifecycle skills and change `Eleven in total` to `Twelve in total`; change `all eleven` in "The four virtues" to `all twelve`; in "Initiative folder", add `plan-human.md` and `CLOSED.md` to the files listed; in "Sync", add that folders holding `CLOSED.md` are left out and counted in one line.
  `skills/init/SKILL.md`: the chain line becomes `/arcdlc:init` → `/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive` → `/arcdlc:close`.
  `.claude-plugin/plugin.json` and `.antigravity-plugin/plugin.json`: `version` to `0.34.0`.
  Out of scope: anything inside the registry markers (`arctool sync` owns them), the ADR files, the skills changed by TRK-3 to TRK-6.
- WHERE:
  Docs: `AGENTS.md`, `README.md`, `CONTEXT.md`.
  Skills: `skills/init/SKILL.md`.
  Manifests: `.claude-plugin/plugin.json`, `.antigravity-plugin/plugin.json`.
- WHY: The guides and counts are read by every agent and every new user; they must name the new skill, the new end state and the new status command (AIC Technical Constraints).
- Acceptance:
  - GIVEN the change WHEN `grep -rn -i "eleven" AGENTS.md CONTEXT.md README.md` runs THEN it prints no line about the skill count.
  - GIVEN the change WHEN `grep -n "arcdlc:close" README.md AGENTS.md skills/init/SKILL.md` and `grep -n "arctool status" README.md` run THEN each file prints at least one line.
  - GIVEN both manifests WHEN `jq -r .version .claude-plugin/plugin.json .antigravity-plugin/plugin.json` runs THEN both print `0.34.0`.
  - GIVEN the change WHEN `arctool sync --check` (or `go run ./cmd/arctool sync --check`) runs THEN it exits 0.
  - GIVEN the change WHEN `go test ./...` runs THEN it passes.
- References: `docs/aics/aic-tracking/aic.md`, `docs/adr/0028-a-finished-initiative-is-closed-in-place.md`, `AGENTS.md`.
- Status: TODO.
