# ArcDLC Findings: an executor files what it saw outside its task as a blocked task

> Let an `/arcdlc:execute` executor append a problem it noticed outside its task to the end of `plan.md` as a `BLOCKED` task, found during its own, and have the run take each one to the engineer at its end, so nothing it saw is lost when the session ends.

## Goals

### 🟢 Business Case

A problem an `/arcdlc:execute` executor sees outside its task is lost today when its session ends.
Its step 3 tells it to leave the problem alone, and the one line of notes it returns lives only in the
dispatcher's chat. The problem surfaces again only when it breaks something, when it costs more to
find. Nobody tracks a money figure for this; the measure is qualitative: every problem an executor
saw is either in the plan or dismissed by the engineer on the record.

### 🟢 Functional Overview

A **finding** is a defect, risk or inconsistency an executor sees in files it read for its task,
that its own task does not cover and that does not stop its task reaching `DONE`. Ideas and
improvements are not findings. A problem that stops the task is not a finding either: it keeps
today's path, a `BLOCKED` question or defect. The two paths never overlap. The term is in
`CONTEXT.md`.

### 🟢 Quality Goals

1. **Reliability (ISO 25010: fault tolerance).** `arctool add` never corrupts `plan.md`. Scenario: for
   any input, valid or not, the bytes of every existing block are unchanged after the call, checked by
   a test that compares them; a refused input leaves the file byte for byte as it was.
2. **Usability (ISO 25010: operability).** The engineer reviews every finding in one place, at the end
   of the run, one question per turn. Scenario: after a run that filed findings, every block whose
   reason starts with `found during` was either asked about or listed in the report.
3. **Maintainability (ISO 25010: modifiability).** The plan format contract does not change. Scenario:
   every existing plan in `internal/plan` tests parses and validates exactly as before.

### 🟢 Organizational Constraints

The owner of the bundle builds this alone, with no schedule; it ships in the next bundle and CLI
release.

### 🟢 Technical Constraints

- `arctool` stays pure standard library, released as static binaries (`CGO_ENABLED=0`).
- Exit codes keep their meanings: 0 ok, 1 contract failure, 2 usage, 3 not found or empty, 4 I/O,
  5 self-validation.
- Every write is atomic and byte-preserving outside its own region.
- Skills stay install-agnostic, and every skill that uses `arctool` describes the fallback without it.
- A change touching the plan format means `plan-format.md`, `internal/plan`, its tests and every skill
  that references the format move in one change set.

### 🟢 Business Context

- **The engineer** answers the end-of-run review: release, sharpen, redesign, or dismiss.
- **`/arcdlc:execute`** writes findings (through its executors, with `arctool add`) and runs the
  review.
- **`/arcdlc:plan`** sharpens a finding into a mechanical task, or deletes a dismissed one.
- **`/arcdlc:aic`** takes a finding that moves the design.
- **`/arcdlc:close`** refuses the outcome `delivered` while a finding is still `BLOCKED`.
- **`arctool`** gains `add`; `next`, `list`, `show`, `todo`, `validate` and `archive` handle a finding
  as any other block.

## Architectural Hypotheses

### 🔵 Architectural Hypotheses

Decision record: [ADR-0029](../../adr/0029-a-finding-is-filed-as-a-blocked-task-at-the-end-of-the-plan.md).

**H1. A finding is a task block appended to the end of `plan.md`, `BLOCKED`, with the reason
"found during <TASK-ID>".**

- Context: an executor's notes today go only to the dispatcher's chat, and are gone when the session
  ends. The engineer already reads `plan.md`, and every tool in the bundle already reads task blocks.
- Decision: the executor writes each finding as a full task block in the plan format and appends it
  after the last block of the same `plan.md`, with `- Status: BLOCKED — found during <TASK-ID>:
  <reason>.` (the `—` is the separator `arctool block` already writes). The plan format contract
  does not change: no new section, no new key, no parser change.
- Contents: only facts the executor saw. `WHAT` names the change the problem calls for, `WHERE` the
  exact `file:line` it read, `WHY` what breaks if nobody acts, `Acceptance` one runnable re-check (a
  test, a grep, a command), `References` the file it read and the plan. No `HOW`: an executor that
  knew the fix decisions would be designing, and the review or `/arcdlc:plan` settles them. The
  reason names what a person must decide, or says `ready to release` when the block is already
  mechanical.
- Justification: full reuse, and one state machine. `arctool next` returns only `TODO` tasks, so the
  running queue never picks up a finding by itself: a person releases it first. `list`, `show`,
  `todo`, `order`, `validate`, `archive`, `/arcdlc:plan-human` and `/arcdlc:close` already handle a
  `BLOCKED` block. `/arcdlc:close` never writes `delivered` while a task is `BLOCKED`, so an open
  finding cannot be closed away unnoticed.
- Trade-offs: a task written by an executor, not by `/arcdlc:plan`, is not guaranteed to be
  mechanical, which ADR-0021 requires of every planned task; the review in H3 is where that is
  checked. `BLOCKED` now means two things, a task that stopped and a finding waiting for review, and
  the reason prefix `found during` is what tells them apart.
- Reversed in this round, by the engineer, three times: a `## Findings` section at the end of
  `plan.md` with its own block shape (it needed a parser and archiver change); a task block appended
  as `TODO` (it put unreviewed work straight into the running queue); a gap block in `gap.md` with
  a register-only question line (it gave `gap.md` a second writer and a second meaning). Also
  rejected: a new sibling register `findings.md`, one more file nobody opens.

**H2. A new command, `arctool add`, appends a task block atomically.**

- Context: no `arctool` command adds a task today, and `/arcdlc:execute` never edits `plan.md` text,
  because a hand edit can corrupt the file the status writes depend on.
- Decision: `arctool add --aic <slug>` reads one task block on standard input and appends it after the
  last block of `plan.md`, then prints the added ID. It refuses with exit 1 (contract failure) when the
  input is not exactly one `### ` block with a well-formed heading, when the ID is already in
  `plan.md` or `plan-archive.md` (so a finding never reuses an archived ID), when the status is `TAKEN`
  or `DONE` (only `TODO` or `BLOCKED` can be added), or when the block fails the task checks of
  `validate --strict`. For a `BLOCKED` block it also refuses a reason that is not exactly
  `found during <ID>: <text>`, with `<ID>` a task in `plan.md` or `plan-archive.md`, and an own ID that
  is not `<ID>-F<n>`. It writes through temp file plus rename, changes no byte before the end of the
  file, and re-parses its own output before the rename, exiting 5 on a mismatch, as `order` and
  `archive` do. The executor that saw the problem runs it, while it still has the detail, and the new
  block is part of that task's own commit with its status change (in a workspace, the hub commit
  that marks the task `DONE`). Without `arctool`, the executor appends the block by hand, as the
  fallback.
- Task ID: a finding is `<TASK-ID>-F<n>`, the ID of the task it was found during plus a counter from 1
  (`INIT-4-F1`, `INIT-4-F2`). The executor writes it; `arctool add` catching a clash is what keeps the
  rule mechanical.
- Justification: the one new text write in `execute` goes through `arctool`, for the same reason
  status writes do.
- Trade-offs: a second writer of task text beside `/arcdlc:plan`, limited to appending. Rejected: the
  subagent returns findings in its report and the dispatcher appends them, which grows the
  dispatcher's context on every task.

**H3. The dispatcher reviews the findings at the end of the run.**

- Context: a finding waits as `BLOCKED` until a person decides what happens to it.
- Decision: the review is the last step of the run, after the Verification phase. Verification part 1,
  which today requires no `BLOCKED` task, ignores the blocks whose reason starts with `found during`
  and checks the rest as before. The dispatcher then lists every such block and takes each to the
  engineer, one question per turn: release it (`arctool todo <id>`), sharpen it with
  `/arcdlc:plan <slug>`, take it to `/arcdlc:aic <slug>` when it moves the design, or dismiss it. A
  released finding is executed by the next run, never by this one.
- Recording the answer: the review writes each verdict into the finding's own reason with
  `arctool block <id> -m "found during <TASK-ID>: <verdict>: <answer>"`, where `<verdict>` is
  `sharpen`, `redesign` or `dismiss`; `block` on a `BLOCKED` task already rewrites only the reason,
  atomically. A released finding needs no verdict: `arctool todo` is the record. `/arcdlc:plan` then
  turns every `sharpen` block into a `TODO` task with the answer in `HOW`, deletes every `dismiss`
  block with the answer in its commit body, and leaves every `redesign` block for the next
  `/arcdlc:aic` round. Settled at plan time, 2026-09-25.
- Dismiss: the dispatcher never edits `plan.md` text, so a dismissed finding stays `BLOCKED` and is
  handed to `/arcdlc:plan <slug>`, which deletes the block (it owns plan text) and records the
  engineer's reason in its commit body. Rejected: marking it `DONE`, which would archive work that was
  never done; and a new `arctool drop` command, one more writer of plan text.
- Run modes: every mode files findings the same way (orchestrator subagent, in-session, single-task).
  The review runs wherever the engineer is in the session: at the end of a whole-queue run and at the
  end of a single-task run. A non-interactive run (`claude -p`) leaves the findings `BLOCKED` and lists
  each one, ID and reason, first in its report. It never decides for the engineer.
- Justification: the engineer sees every finding once, with the whole run already verified, and the
  running queue never grows mid-run.

## Assessment

### 🔴 Technical Challenges & Risks

- **R1. A cheap executor files noise.** Style remarks, guesses, one problem split into five, and the
  end-of-run review drowns the engineer. Likelihood medium, impact medium. Mitigation, reduce by
  design: at most three findings per task, and the spawn prompt repeats the definition of a finding
  (a defect, risk or inconsistency seen in a file the executor read; ideas and improvements are not
  findings). An executor that saw more than three records the three that matter most and says in its
  one line of notes how many it left out.

- **R2. Every planned task is `DONE`, but a finding is still `BLOCKED`.** `/arcdlc:close` then offers
  only `partly delivered` or `stopped` and lists the finding under `Not done`, and `arctool status`
  shows `in progress`, not `ready to close`. Likelihood high, impact low. Not a real risk, by the
  engineer's call: a finding nobody decided on is work left open, and `Not done` in `CLOSED.md` is
  where it stays visible once the plan is deleted. The engineer who wants `delivered` releases and
  finishes, or dismisses, the findings first. Close does not change.

- **R3. A mistyped reason hides a finding.** The reason `found during <TASK-ID>:` is the only thing
  that tells a finding from a stopped task. A typo (`Found in INIT-4`) makes the Verification phase
  treat it as a real block and stop, and the review misses it. Likelihood medium with a cheap executor,
  impact low: the failure is loud, not silent. Mitigation, reduce by design: `arctool add` enforces the
  shape (see H2), so the rule lives in code rather than in the executor's reading. Residual: the hand
  append used without `arctool` has no check, and a typo there still fails loudly at Verification.

### Open questions

None. Every section was answered in the interview of 2026-09-25; nothing was deferred.

| Question | Who answers | By when (phase or stage) |
| --- | --- | --- |
