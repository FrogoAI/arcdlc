# Plan: findings

<!-- arcdlc:source docs/aics/findings/aic.md sha256:01a838ee29f8e90a33982e4bc5a5831e6714d4d285cb5723113fd3974e849205 -->

Format: [`skills/plan/references/plan-format.md`](../../../skills/plan/references/plan-format.md).
Source: [`aic.md`](aic.md), decided 2026-09-25. Terms: [`CONTEXT.md`](../../../CONTEXT.md).
Decision: [ADR-0029](../../adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md).

## Risk Coverage

- R1 (a cheap executor files noise): FND-4 puts the definition of a finding and the cap of three per task into the per-task contract and the spawn prompt.
- R2 (every planned task is `DONE` but a finding is still `BLOCKED`, so close cannot write `delivered`): accepted by the engineer as intended; no task. `/arcdlc:close` does not change.
- R3 (a mistyped reason hides a finding): FND-1 enforces the `found during <ID>:` reason and the `<ID>-F<n>` ID in `Plan.Append`; FND-2 exposes it as `arctool add`. The residual (the hand append without `arctool` has no check) fails loudly at Verification, which FND-5 keeps.
- Open questions: none in the architecture document.

### FND-1: Add `Plan.Append` to `internal/plan` with the finding checks

- WHAT: Add `func (p *Plan) Append(block []byte, archived []string) (out []byte, id string, err error)` in a new file `internal/plan/add.go`, plus an exported `func (p *Plan) StatusReason(t *Task) string`, so a caller can append one checked task block to the end of a plan.
- HOW:
  Error type, mirroring `ReorderError` in `internal/plan/order.go`: `type AddErrKind uint8` with constants `AddContract` (maps to exit 1) and `AddSelfCheck` (maps to exit 5); `type AddError struct { Kind AddErrKind; Msg string; Findings []Finding }`; `Error()` returns `Msg`, plus `": "` and the findings' `Message` values joined by `"; "` when `Findings` is not empty. Every error `Append` returns is an `*AddError`.
  `StatusReason(t *Task) string`: returns `""` when `t.StatusLineStart < 0`; otherwise calls the existing `dissectStatusLine(p.Bytes[t.StatusLineStart:t.StatusLineEnd])` in `internal/plan/mutate.go` and returns its `reason` value.
  `Append` checks, in this order, each failure an `AddContract` error with the message given:
  1. `nb := Parse(block)`; `len(nb.Tasks) != 1` gives `"input must be exactly one ### task block"`. Any non-blank line before the block's heading (a preamble region in `nb.Regions` containing non-whitespace bytes) gives the same message.
  2. `t := nb.Tasks[0]`; `!t.HeadingOK || t.ID == "" || t.Title == ""` gives `"malformed heading"`.
  3. `len(p.ByID(t.ID)) > 0` or `t.ID` is in `archived` gives `"task ID <ID> already exists"` (ID quoted with `%q`).
  4. `t.Status` is not `StatusTODO` and not `StatusBLOCKED` gives `"only a TODO or BLOCKED task can be added"`.
  5. `fs := nb.Validate(ValidateOpts{Strict: true, RequireAcceptance: true})`; any finding (error or warning) gives `"block fails the strict task checks"` with `Findings: fs`.
  6. When `t.Status == StatusBLOCKED`: `r := nb.StatusReason(&nb.Tasks[0])` must match `^found during (\S+): (.+)$` (package-level `regexp.MustCompile`), else `"a BLOCKED block needs the reason \"found during <TASK-ID>: <reason>\""`. The captured source ID must be a task ID in `p` (`len(p.ByID(src)) > 0`) or in `archived`, else `"source task <ID> is not in the plan or its archive"`. `t.ID` must match `"^" + regexp.QuoteMeta(src) + "-F[1-9][0-9]*$"`, else `"a finding's ID must be <source>-F<n>, got <ID>"`.
  Rendering: `body := bytes.TrimRight(block, "\r\n")` with leading blank lines trimmed too; when `p.CRLF` is true, convert every `"\n"` in `body` not already preceded by `"\r"` to `"\r\n"`. `nl` is `"\r\n"` when `p.CRLF`, else `"\n"`. `out` = `p.Bytes`; plus `nl` when `p.Bytes` is non-empty and does not end in `"\n"`; plus `nl` (one blank line) when `p.Bytes` is non-empty; plus `body`; plus `nl`. Never modify `p.Bytes` in place (copy first).
  Self-check before returning: `np := Parse(out)`; fail with `AddSelfCheck` and message `"self-check failed: <what>"` unless `bytes.HasPrefix(out, p.Bytes)`, `len(np.Tasks) == len(p.Tasks)+1`, and `np.Tasks[len(np.Tasks)-1].ID == t.ID`. Return `out, t.ID, nil` on success.
  Out of scope: the `arctool add` command and its exit-code mapping (FND-2); any change to `Parse`, `Validate`, or the plan format.
- WHERE:
  Layer `domain`: `internal/plan/add.go`.
  Tests: `internal/plan/add_test.go`.
- WHY: The executor's findings must enter the plan through one checked, append-only write so a malformed or mis-shaped finding never reaches the file (AIC H2, risk R3).
- Acceptance:
  - GIVEN a plan with task `A-1` (TODO) and a well-formed block `### A-1-F1: Fix x` with status `BLOCKED — found during A-1: ready to release.` WHEN `Append` runs with `archived` nil THEN it returns ID `A-1-F1`, `out` starts with the original bytes, and `Parse(out)` has two tasks; checked by `TestAppendFinding` in `internal/plan/add_test.go`.
  - GIVEN the same plan WHEN `Append` gets a block whose ID is `A-1`, a block with status `DONE`, two blocks, a block missing `Acceptance`, a `BLOCKED` block with reason `found in A-1: x`, a `BLOCKED` block with source `Z-9` not in plan or archive, and a `BLOCKED` block `### B-2: ...` found during `A-1` THEN each returns an `*AddError` with `Kind == AddContract`; checked by the table test `TestAppendRefuses`.
  - GIVEN a `BLOCKED` finding whose source ID is only in `archived` WHEN `Append` runs THEN it succeeds; checked by `TestAppendSourceInArchive`.
  - GIVEN a plan using CRLF line endings and a block written with LF WHEN `Append` runs THEN every line of the appended block in `out` ends in `\r\n`; checked by `TestAppendCRLF`.
  - GIVEN a plan whose last byte is not a newline WHEN `Append` runs THEN `out` starts with the original bytes and the new heading starts on its own line after one blank line; checked by `TestAppendNoTrailingNewline`.
  - GIVEN the change WHEN `go test ./internal/plan/...`, `go vet ./...` and `gofmt -l .` run THEN tests pass, vet is clean and gofmt prints nothing.
- References: `docs/aics/findings/aic.md`, `docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`, `internal/plan/order.go`, `internal/plan/mutate.go`, `internal/plan/validate.go`.
- Status: DONE.

### FND-2: Add the `arctool add` command

- WHAT: Add `arctool add [--aic SLUG | --plan PATH]`, which reads one task block on standard input and appends it to the plan through `Plan.Append`, then bump `arctool` to 0.22.0.
- HOW:
  In `cmd/arctool/main.go`: add `case "add": os.Exit(cmdAdd(os.Args[2:]))` to the switch in `main`, after `"block"`/`"todo"`. Add `cmdAdd(args []string) int` built like `cmdArchive`: a `flag.FlagSet` named `add` with `--plan` and `--aic`; a parse error or any positional argument returns 2 with `usage: arctool add [--aic SLUG | --plan PATH] < block.md` on stderr.
  It calls `cmdAddFrom(os.Stdin, planFlag, aicFlag)`, a helper taking an `io.Reader` so tests can pass a `strings.Reader`. The helper: `resolvePlan`; `loadPlan` (its exit 4 stands); read all of the reader (read error: `arctool: cannot read standard input: <err>`, return 4); load the archive IDs from `filepath.Join(filepath.Dir(planPath), "plan-archive.md")` with `plan.Parse` (a missing file means no IDs; any other read error returns 4); call `p.Append(input, ids)`.
  Error mapping via `errors.As` to `*plan.AddError`: print `arctool: <err>` to stderr; `AddContract` returns 1; `AddSelfCheck` also prints `arctool: nothing written` and returns 5.
  On success `atomicWrite(planPath, out)` (failure: `arctool: write <path>: <err>`, return 4), then print `added <ID> to <planPath>` to stdout and return 0.
  Usage text: after the `block` line add `  arctool add    [--aic SLUG | --plan PATH] < block.md   append one task block to the end of the plan` and on the next line `                 a BLOCKED block is a finding: reason "found during <ID>: ...", ID <ID>-F<n>`. Update the package doc comment's command list to include `add`. Set `const version` to `"0.22.0"`.
  Out of scope: `Plan.Append` itself (FND-1); every skill file (FND-4 to FND-7).
- WHERE:
  Layer `cmd`: `cmd/arctool/main.go`.
  Tests: `cmd/arctool/main_test.go`.
- WHY: Executors and their fallback docs need one command that appends a finding atomically instead of hand-editing `plan.md` (AIC H2).
- Acceptance:
  - GIVEN a temp plan with task `A-1` and stdin holding a valid finding block `A-1-F1` WHEN `cmdAddFrom` runs with `--plan` THEN it returns 0, stdout contains `added A-1-F1`, and `arctool list` on the file shows `A-1-F1` as `BLOCKED`; checked by `TestCmdAddAppends`.
  - GIVEN the same plan and a block with a duplicate ID WHEN `cmdAddFrom` runs THEN it returns 1 and the plan file is byte for byte unchanged; checked by `TestCmdAddRefusesAndLeavesTheFileAlone`.
  - GIVEN a `plan-archive.md` beside the plan holding task `A-0` and a finding block `A-0-F1` found during `A-0` WHEN `cmdAddFrom` runs THEN it returns 0; and a block with ID `A-0` returns 1; checked by `TestCmdAddReadsTheArchive`.
  - GIVEN a positional argument WHEN `cmdAdd` runs THEN it returns 2; checked by `TestCmdAddUsage`.
  - GIVEN `fmt.Sprintf(usage, version)` WHEN it is searched THEN it contains `arctool add    [--aic SLUG | --plan PATH] < block.md` and `version` is `0.22.0`; checked by `TestUsageDocumentsAdd`.
  - GIVEN the change WHEN `go test ./...`, `go vet ./...` and `gofmt -l .` run THEN tests pass, vet is clean and gofmt prints nothing.
- References: `docs/aics/findings/aic.md`, `cmd/arctool/main.go`, `internal/plan/add.go`.
- Status: DONE.

### FND-3: Document findings and `arctool add` in the plan format guide

- WHAT: Describe finding blocks, their reason and ID shape, the review verdicts, and `arctool add` in `skills/plan/references/plan-format.md`.
- HOW:
  In `## Status Lifecycle (In-Block)`, after the numbered list, add a subsection `### Findings` saying, in plain prose: an executor that sees a problem outside its task files it as an ordinary task block appended to the end of the plan, never inside another block; its status line is `- Status: BLOCKED — found during <TASK-ID>: <reason>.` (the `—` is the separator `arctool block` writes); its ID is `<TASK-ID>-F<n>`, counting from 1; it carries `WHAT`, `WHERE` with `file:line`, `WHY`, one runnable `Acceptance` check and `References`, and no `HOW`; it is appended with `arctool add --aic <slug> < block.md`, which refuses a block breaking any of these rules with exit 1 and writes nothing, and by hand only when `arctool` is absent. Then list the three verdicts the end-of-run review writes into the reason with `arctool block <id> -m`: `found during <TASK-ID>: sharpen: <answer>`, `found during <TASK-ID>: redesign: <answer>`, `found during <TASK-ID>: dismiss: <answer>`, and say a released finding is simply `TODO` again (`arctool todo <id>`). Link ADR-0029 as `../../../docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`.
  In `## Authoring Rules`, rule 10, add one sentence: the executor's one text write to `plan.md` is appending a finding through `arctool add`.
  Set the `**Reviewed**:` line to `2026-09-25` only after reading the whole file end to end.
  Out of scope: any change to the task block format, the eight keys, or `internal/plan`; the skills (FND-4 to FND-6).
- WHERE:
  Layer `docs`: `skills/plan/references/plan-format.md`.
- WHY: The format guide is the contract every skill reads; a finding shape that only the skills describe would drift from what `arctool add` enforces (AIC H1, H2).
- Acceptance:
  - GIVEN the change WHEN `grep -n 'found during <TASK-ID>' skills/plan/references/plan-format.md` runs THEN it prints at least four lines (the status line and the three verdicts).
  - GIVEN the change WHEN `grep -n 'arctool add' skills/plan/references/plan-format.md` runs THEN it prints at least two lines.
  - GIVEN the change WHEN `grep -n '^\*\*Reviewed\*\*: 2026-09-25' skills/plan/references/plan-format.md` runs THEN it prints one line, and `go test ./internal/bundle/...` passes.
- References: `docs/aics/findings/aic.md`, `docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`, `skills/plan/references/plan-format.md`.
- Status: DONE.

### FND-4: Teach the executor to file findings in `/arcdlc:execute`

- WHAT: Add the filing of findings to the per-task contract and the spawn prompt in `skills/execute/SKILL.md`, and make the "you never edit `plan.md`" lines name the one exception.
- HOW:
  In `## Per-task contract`, insert a new step between step 4 (tests and lint) and step 5 (verify acceptance), and renumber the later steps (5 to 6, 6 to 7, 7 to 8, 8 to 9), updating every in-file reference to those step numbers (for example "step 5 of the contract", "per-task contract, step 7"). The new step 5, "File what you saw outside this task", says: a finding is a defect, risk or inconsistency you saw in a file you read for this task, that this task does not cover and that does not stop it reaching `DONE`; ideas and improvements are not findings, and a problem that stops the task is a block (step 8), not a finding. File at most three, the three that matter most, and put how many you left out in your one line of notes. Each is a task block with ID `<TASK-ID>-F<n>` counting from 1, `WHAT`, `WHERE` with `file:line`, `WHY`, one runnable `Acceptance` check, `References` naming the file you read and `docs/aics/<slug>/plan.md`, no `HOW`, and `- Status: BLOCKED — found during <TASK-ID>: <reason>.` where the reason names what a person must decide, or reads `ready to release` when the block is already mechanical. Append it with `arctool add --aic <slug> < <file>`; exit 1 means the block broke a rule, so fix it and retry, never drop the finding. *Fallback: append the block to the end of `plan.md` by hand, after one blank line.* It is committed with this task (step 7); in a workspace it rides in the hub commit that marks the task `DONE`. Point at `plan-format.md` `### Findings` for the full shape.
  In `## Session strategy`, orchestrator mode step 3, add a bullet to what the prompt must say: the definition of a finding and the cap of three, from the new step 5, word for word.
  Change `(you never edit \`plan.md\`)` in the tier paragraph and `You never edit \`plan.md\`.` in orchestrator step 5 so each adds `other than appending a finding through \`arctool add\``. In `## When a task is unclear`, the bullet "Do not edit another task, and do not rewrite this task's text in `plan.md`" gains: "Appending a finding (per-task contract, step 5) is the one exception."
  Keep the `## Talk simple, write like a human` and `## Judge by the four virtues` blocks byte identical.
  Out of scope: the Verification phase, the review and the report (FND-5); `skills/plan/SKILL.md` (FND-6).
- WHERE:
  Layer `skills`: `skills/execute/SKILL.md`.
- WHY: Without a rule to file it, an executor's finding lives only in the dispatcher's chat and is lost when the session ends (AIC business case, H1, risk R1).
- Acceptance:
  - GIVEN the change WHEN `grep -n 'File what you saw outside this task' skills/execute/SKILL.md` runs THEN it prints one line inside `## Per-task contract`.
  - GIVEN the change WHEN `grep -c 'arctool add' skills/execute/SKILL.md` runs THEN it prints 3 or more.
  - GIVEN the change WHEN `grep -n 'at most three' skills/execute/SKILL.md` runs THEN it prints at least two lines (the contract step and the spawn prompt).
  - GIVEN the change WHEN `go test ./internal/bundle/...` and the `run` scripts of the CI steps `Check the writing-style block is the same in every skill` and `Check the virtues block is the same in every skill` in `.github/workflows/ci.yml` are run from the repository root THEN all three exit 0.
- References: `docs/aics/findings/aic.md`, `docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`, `skills/execute/SKILL.md`, `skills/plan/references/plan-format.md`.
- Status: DONE.

### FND-5: Add the findings review to the end of an `/arcdlc:execute` run

- WHAT: Make Verification part 1 ignore findings, add a `## Review the findings` section after `## Verification phase`, and list findings in the report, in `skills/execute/SKILL.md`.
- HOW:
  Verification part 1, first bullet: `arctool list --status BLOCKED` must be empty except for blocks whose reason starts with `found during`; list them with `arctool list --status BLOCKED` and read each reason with `arctool show <id>`. A `BLOCKED` block with any other reason is still a stopped task and is reported as today.
  New section `## Review the findings (last step of the run)`, placed after `## Verification phase` and before `## Report`. It says: list every `BLOCKED` block whose reason starts with `found during`. For each, one question per turn through the grilling protocol, carrying a recommended answer, offer four answers. Release: `arctool todo <id>`; the next run executes it, never this one. Sharpen: `arctool block <id> -m "found during <TASK-ID>: sharpen: <answer>"`. Redesign: `arctool block <id> -m "found during <TASK-ID>: redesign: <answer>"`, and suggest `/arcdlc:aic <slug>`. Dismiss: `arctool block <id> -m "found during <TASK-ID>: dismiss: <answer>"`. `<TASK-ID>` is the one already in the reason. Each `block` or `todo` is a status write, so in a workspace it gets its hub commit per `## In a workspace` step 2; in a single repository commit `plan.md` once after the review, message `chore(<slug>): record the findings review`, plus `#AI-assisted`. When any verdict is `sharpen` or `dismiss`, end by telling the engineer to run `/arcdlc:plan <slug>`. The review runs at the end of a whole-queue run and at the end of a single-task run. A non-interactive run asks nothing: it leaves every finding `BLOCKED` and lists each, ID and reason, first in its report.
  `## Report`: add a sentence: list every finding filed this run (ID and reason) and the verdict each got, before the per-task summary.
  Keep the `## Talk simple, write like a human` and `## Judge by the four virtues` blocks byte identical.
  Out of scope: filing findings and the spawn prompt (FND-4); `skills/plan/SKILL.md` (FND-6).
- WHERE:
  Layer `skills`: `skills/execute/SKILL.md`.
- WHY: A finding waits as `BLOCKED` until a person decides; the review is where the engineer sees each one once, after the run is verified (AIC H3, ADR-0029 amending ADR-0023).
- Acceptance:
  - GIVEN the change WHEN `grep -n '^## Review the findings' skills/execute/SKILL.md` runs THEN it prints one line, after the line of `^## Verification phase` and before the line of `^## Report`.
  - GIVEN the change WHEN `grep -c 'found during <TASK-ID>: \(sharpen\|redesign\|dismiss\)' skills/execute/SKILL.md` runs THEN it prints 3 or more.
  - GIVEN the change WHEN `grep -n 'found during' skills/execute/SKILL.md` runs THEN at least one match lies inside `## Verification phase`.
  - GIVEN the change WHEN `go test ./internal/bundle/...` and the `run` scripts of the CI steps `Check the writing-style block is the same in every skill` and `Check the virtues block is the same in every skill` in `.github/workflows/ci.yml` are run from the repository root THEN all three exit 0.
- References: `docs/aics/findings/aic.md`, `docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`, `docs/adr/0023-the-dispatcher-verifies-and-the-subagent-never-guesses.md`, `skills/execute/SKILL.md`.
- Status: TODO.

### FND-6: Teach `/arcdlc:plan` to act on reviewed findings

- WHAT: Add a step to `skills/plan/SKILL.md` that turns `sharpen` findings into `TODO` tasks, deletes `dismiss` findings, and leaves `redesign` findings alone.
- HOW:
  Add `## Step 1.5 — Act on reviewed findings` between `## Step 1` and `## Step 2`. It says: when `docs/aics/<slug>/plan.md` exists, find every `BLOCKED` block whose reason starts with `found during` (`arctool list --status BLOCKED --aic <slug>`, then `arctool show <id>`). Reason `found during <X>: sharpen: <answer>`: write `<answer>` into a `HOW` key as the engineer's decision, make the block pass the mechanical check of Step 3, keep its ID, and set its status to `TODO` with `arctool todo <id>`. Reason `found during <X>: dismiss: <answer>`: delete the whole block from `plan.md` and put `Dismissed <id>: <answer>` in the body of the commit that carries the change. Reason `found during <X>: redesign: <answer>`: leave it as it is; it waits for `/arcdlc:aic <slug>`. A finding with no verdict: leave it as it is; it waits for the review in `/arcdlc:execute`. A `sharpen` answer that is not enough to make the block mechanical is a design-readiness gap, handled by Step 2.4. Run `arctool validate --strict --aic <slug>` afterwards.
  Keep the `## Talk simple, write like a human` and `## Judge by the four virtues` blocks byte identical.
  Out of scope: `skills/execute/SKILL.md` (FND-4, FND-5); the format guide (FND-3).
- WHERE:
  Layer `skills`: `skills/plan/SKILL.md`.
- WHY: `/arcdlc:plan` owns plan text, so it is the only skill that may sharpen or delete a finding the engineer reviewed (AIC H3).
- Acceptance:
  - GIVEN the change WHEN `grep -n '^## Step 1.5 — Act on reviewed findings' skills/plan/SKILL.md` runs THEN it prints one line, between the lines of `^## Step 1` and `^## Step 2 `.
  - GIVEN the change WHEN `grep -c 'sharpen\|dismiss\|redesign' skills/plan/SKILL.md` runs THEN it prints 3 or more.
  - GIVEN the change WHEN `go test ./internal/bundle/...` and the `run` scripts of the CI steps `Check the writing-style block is the same in every skill` and `Check the virtues block is the same in every skill` in `.github/workflows/ci.yml` are run from the repository root THEN all three exit 0.
- References: `docs/aics/findings/aic.md`, `docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`, `skills/plan/SKILL.md`.
- Status: TODO.

### FND-7: Update `AGENTS.md`, `README.md` and the bundle version for findings

- WHAT: Record findings and `arctool add` in the hard rules of `AGENTS.md` and in `README.md`, and bump the bundle to 0.35.0 in both plugin manifests.
- HOW:
  `AGENTS.md`, rule "Every write is atomic and byte-preserving outside its own region": add a sentence after the `order` sentence: `add` appends one block after the last byte of the plan, changes no earlier byte, re-parses its own output, and writes nothing on a failed check. Add a new hard rule after "A spawned subagent has nobody to ask, and must never guess in place of asking": **A finding is filed, never dropped and never run unasked.** One paragraph: an executor that sees a problem outside its task appends it as a `BLOCKED` task, `found during <TASK-ID>: <reason>`, ID `<TASK-ID>-F<n>`, at most three per task, through `arctool add`; `next` never returns it; the dispatcher reviews each at the end of the run (release, sharpen, redesign, dismiss, written into the reason); `/arcdlc:plan` acts on the verdicts. Link ADR-0029 as `docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`.
  `README.md`: in the skills table row for `/arcdlc:execute <slug> [TASK-ID]`, add to the description: "Problems seen outside a task are filed as blocked findings at the end of the plan and reviewed with you at the end of the run." In the `### arctool CLI` section, first paragraph, add `appends a checked task block (\`arctool add\`)` to the list of what it does.
  Set `"version"` to `"0.35.0"` in `.claude-plugin/plugin.json` and `.antigravity-plugin/plugin.json`.
  Do not edit inside the `<!-- arcdlc:initiatives -->` markers; run `arctool sync --check` instead.
  Out of scope: `CONTEXT.md` (the `Finding` term was written by `/arcdlc:aic`); every `SKILL.md`.
- WHERE:
  Layer `docs`: `AGENTS.md`, `README.md`.
  Layer `manifest`: `.claude-plugin/plugin.json`, `.antigravity-plugin/plugin.json`.
- WHY: `AGENTS.md` is the rule set agents working on this repository read first, and the manifests carry the bundle version that ships the changed skills.
- Acceptance:
  - GIVEN the change WHEN `grep -n 'A finding is filed, never dropped' AGENTS.md` runs THEN it prints one line.
  - GIVEN the change WHEN `grep -n 'arctool add' README.md AGENTS.md` runs THEN each file has at least one match.
  - GIVEN the change WHEN `grep -h '"version"' .claude-plugin/plugin.json .antigravity-plugin/plugin.json` runs THEN both lines read `"version": "0.35.0",`.
  - GIVEN the change WHEN `arctool sync --check` runs THEN it exits 0.
- References: `docs/aics/findings/aic.md`, `docs/adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md`, `AGENTS.md`, `README.md`.
- Status: TODO.
