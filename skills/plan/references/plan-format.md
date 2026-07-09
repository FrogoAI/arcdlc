# Plan Task Authoring Guide

This guide defines the `docs/aics/<slug>/plan.md` task format. The plan is an executable queue: `/arcdlc:execute` (or
any compatible runner) picks tasks off it mechanically, so the format is a contract, not a style preference. It is
written to be executed by a **less capable model than the one that planned it**: every decision that matters belongs
in the task block, not in the executor's judgment.

Keep format rules in this file. Each `plan.md` contains only the plan content and a short link back to this guide —
no runner instructions. `gap.md` (evidence register) and `plan-archive.md` (archive) are always siblings of `plan.md`
in the same `docs/aics/<slug>/` folder. Initiative selection, folder layout, and the initiative registry are defined
by the skills and `AGENTS.md`, not here.

A plan may open with a free-form `## Risk Coverage` preamble section mapping each architecture-document
risk to the task IDs that cover it (or "accepted" with a reason) — see the `/arcdlc:plan` risk-coverage
gate. It is ordinary prose, not a `### ` task block, so the runner ignores it; only `### ` blocks are tasks.

## Required Task Block Format

```md
### <TASK-ID>: <Short Title>

- WHAT: <Clear implementation scope, one line.>
- HOW:
  <Optional. Implementation decisions the executor must follow: signatures, naming, data shapes,
  edge cases, error handling. End with "Out of scope: …" when adjacent work must NOT be touched.>
- WHERE:
  Layer `domain`: <files/modules>
  Layer `repository`: <files/modules>
  Layer `handler`: <files/modules>
  Tests: <files/modules>
  Swagger/docs: <files/modules>
- WHY: <Why this is required / risk if skipped, one line.>
- Acceptance:
  - GIVEN <precondition> WHEN <action> THEN <observable result>.
  - GIVEN <precondition> WHEN <action> THEN <observable result>.
- References: `<doc/path-1>`, `<doc/path-2>`.
- Status: TODO.
```

### Heading

`### <TASK-ID>: <Short Title>` is the normal form. A parenthetical source-status tag —
`### <TASK-ID> (MISSING): <Short Title>` — is **optional** and only carried by gap-derived tasks
(from `/arcdlc:examinate`); the only valid values are `MISSING`, `PARTIAL`, `DRIFT`
(`arctool validate` warns on anything else). Tasks decomposed from an architecture document take
**no** tag. Do not put status in the heading — status lives on the `- Status:` line.

### Section keys

The parser reads exactly these keys (exact casing). `WHAT`, `WHY`, `References`, and `Status` are
**single-line** — anything on following lines is invisible to the runner. `HOW`, `WHERE`, and
`Acceptance` are **multi-line**: they absorb indented lines until the next key. Keep multi-line
bodies to plain/bulleted lines with inline backticks — a fenced code block ends the section.

- `WHAT` (required) — the scope, one sentence.
- `HOW` (optional) — the design decisions a weaker executor would otherwise have to guess:
  function/interface signatures, naming, data shapes, algorithm choice, edge cases, error handling.
  Also the place for scope fencing: `Out of scope: <thing> (covered by <TASK-ID>).`
- `WHERE` (required) — the exact files/modules expected to change, one `Layer` line per layer. The
  layer names above match the ArcDLC Go server layout; for other projects keep the `Layer` line
  format with that project's real layer/module names.
- `WHY` (required) — motivation / risk if skipped, one sentence.
- `Acceptance` (required) — the definition of done: concrete, testable criteria the executor must
  demonstrate. Prefer `GIVEN … WHEN … THEN …`; a checklist of verifiable outcomes is acceptable when
  a scenario would be contrived. Where a criterion is test-verifiable, **name the runnable check**
  (test file, test name, or command) so "met" is demonstrable, not asserted. A criterion must never
  be a paraphrase of `WHAT`.
- `References` (required) — comma-separated source docs needed for implementation. Each entry must
  stay a clean file path (the runner parses and opens them); when a referenced document is long,
  name the relevant section inside `HOW` (e.g. `see docs/aics/checkout/aic.md §"Data model"`) so
  the executor does not re-derive which part applies.
- `Status` (required) — see the lifecycle below.

## Status Lifecycle (In-Block)

Status is tracked by the `- Status:` line inside each task block (not in heading text).

- `- Status: TODO.` means queued.
- `- Status: TAKEN.` means currently being executed.
- `- Status: DONE.` means completed.
- Optional manual status: `- Status: BLOCKED.`

The executor transitions status as follows:

1. Finds the first `###` block whose `- Status:` is `TODO` (top-to-bottom order).
2. Changes status to `TAKEN` before starting implementation, so a crashed or interrupted session is visible.
3. Implements that one task, following `HOW` when present and staying out of anything it fences off.
4. Verifies **every** `Acceptance` criterion (via tests or observable behavior), then sets status to
   `DONE` only when the task is implemented, all criteria are demonstrably met, validated, and committed.
5. On failure — including any acceptance criterion that cannot be met — reverts status to `TODO`, or sets
   `BLOCKED` with a one-line reason naming the failing criterion, instead of leaving `TAKEN` behind.
6. One task per session/commit; repeats from step 1 for the next task.

Because `Acceptance` lives inside the task block, `/arcdlc:archive` carries it verbatim into
`docs/aics/<slug>/plan-archive.md` — the acceptance criteria become the durable record of what "done" meant.

A block missing its `- Status:` line, or whose status is not exactly `TODO` (after trimming spaces and a trailing
period), is skipped by the executor — check this first when a plan appears to have no pending work.

## Gap Register Sync

When `docs/aics/<slug>/gap.md` is used as an input source (e.g. produced by `/arcdlc:examinate`), every `### ...` gap
must have a matching task block in the same folder's `docs/aics/<slug>/plan.md`.

The `plan.md` copy must preserve:

- The same task ID and heading (including the `(MISSING|PARTIAL|DRIFT)` tag).
- The same `WHAT`, `HOW` (when present), `WHERE`, `WHY`, and `Acceptance` content.
- Executor metadata: `References` and `Status`.

Use `gap.md` as the evidence register and `plan.md` as the executable queue — both within the same initiative folder.

## Authoring Rules

1. Use unique `<TASK-ID>` values (for example: `WA240-VER-03`, `AIC-1`). IDs need only be unique **within one
   initiative's `plan.md`** — each `/arcdlc:execute` run targets a single plan.
2. Keep exactly one task per `###` block, headings at level `###`, section keys exact
   (`WHAT`, `HOW`, `WHERE`, `WHY`, `Acceptance`, `References`, `Status`).
3. Keep `Status` values uppercase (`TODO`, `TAKEN`, `DONE`, `BLOCKED`), preferably with a trailing
   period (`- Status: TODO.`).
4. Size each task so one agent session can implement, test, and commit it: one coherent slice,
   roughly ≤5–6 files in `WHERE`. If `WHERE` spans unrelated modules, split the task.
5. Order blocks by dependency — the queue runs top-to-bottom, so a task may only depend on tasks above it.
6. Make each block self-sufficient: an executor reading only the block plus its referenced sections
   must be able to implement it without asking questions. If you cannot name a file or a decision
   while authoring, resolve it now — do not defer it to the executor.
7. `arctool validate --strict` fails a task with a missing or empty `Acceptance` section, an empty
   `References` list, or a `WHERE` with no concrete file/module.
8. Do not place executor instructions inside `plan.md`; update this file instead.

## Minimal Example

```md
### WA999-VER-01: Add endpoint parity for sample flow

- WHAT: Add `/v2/sampleapi/items/{id}` read endpoint with legacy-compatible response.
- HOW:
  Handler calls `item.Service.Get(ctx, id)`; map `item.ErrNotFound` to 404, everything else to 500.
  Response body mirrors the legacy v1 shape (see AIC §"Item API parity") — do not rename fields.
  Out of scope: list endpoint and pagination (WA999-VER-02).
- WHERE:
  Layer `handler`: `services/sampleapi/internal/handler/item.go`, `router.go`.
  Layer `domain`: `services/sampleapi/internal/domain/item/{port,service}.go`.
  Layer `repository`: `services/sampleapi/internal/repository/item.go`.
  Tests: `services/sampleapi/internal/{handler,domain,repository}/item_test.go`.
  Swagger/docs: `services/sampleapi/docs/swagger.json`.
- WHY: Migration parity and consumer cutover are blocked without this route.
- Acceptance:
  - GIVEN an existing item id WHEN `GET /v2/sampleapi/items/{id}` THEN the response is 200 with the legacy-compatible body.
  - GIVEN an unknown id WHEN the same call is made THEN the response is 404.
  - GIVEN the new code WHEN `go test ./services/sampleapi/...` runs THEN `item_test.go` covers the 200 and 404 paths and passes.
- References: `docs/aics/initiative-99/aic.md`.
- Status: TODO.
```
