# Plan Task Authoring Guide

This guide defines the `docs/aics/plan.md` task format. The plan is an executable queue: `/arcdlc:execute` (or any
compatible runner) picks tasks off it mechanically, so the format is a contract, not a style preference.

Keep format rules in this file. `docs/aics/plan.md` should only contain the plan content and a short link back to this
guide.

## Required Task Block Format

```md
### <TASK-ID> (<SOURCE-STATUS>): <Short Title>

- WHAT: <Clear implementation scope.>
- WHERE:
  Layer `domain`: <files/modules>
  Layer `repository`: <files/modules>
  Layer `handler`: <files/modules>
  Tests: <files/modules>
  Swagger/docs: <files/modules>
- WHY: <Why this is required / risk if skipped.>
- Acceptance:
  - GIVEN <precondition> WHEN <action> THEN <observable result>.
  - GIVEN <precondition> WHEN <action> THEN <observable result>.
- References: `<doc/path-1>`, `<doc/path-2>`, `<doc/path-3>`.
- Status: TODO.
```

`Acceptance` is the task's definition of done: concrete, testable success criteria the executor
must be able to demonstrate. Prefer `GIVEN … WHEN … THEN …` scenarios; a plain checklist of
verifiable outcomes is acceptable when a scenario would be contrived. Every criterion must be
something a test or an observable behavior can confirm — not a restatement of `WHAT`.

The `WHERE` layer names above match the ArcDLC Go server layout (`domain`, `repository`, `handler`). For projects with
a different structure, keep the `Layer` line format but use that project's real layer/module names.

## Status Lifecycle (In-Block)

Status is tracked by the `- Status:` line inside each task block (not in heading text).

- `- Status: TODO.` means queued.
- `- Status: TAKEN.` means currently being executed.
- `- Status: DONE.` means completed.
- Optional manual status: `- Status: BLOCKED.`

The executor transitions status as follows:

1. Finds the first `###` block whose `- Status:` is `TODO` (top-to-bottom order).
2. Changes status to `TAKEN` before starting implementation, so a crashed or interrupted session is visible.
3. Implements that one task, with the full plan file as context.
4. Verifies **every** `Acceptance` criterion (via tests or observable behavior), then sets status to
   `DONE` only when the task is implemented, all criteria are demonstrably met, validated, and committed.
5. On failure — including any acceptance criterion that cannot be met — reverts status to `TODO`, or sets
   `BLOCKED` with a one-line reason naming the failing criterion, instead of leaving `TAKEN` behind.
6. One task per session/commit; repeats from step 1 for the next task.

Because `Acceptance` lives inside the task block, `/arcdlc:archive` carries it verbatim into
`docs/aics/plan-archive.md` — the acceptance criteria become the durable record of what "done" meant.

A block missing its `- Status:` line, or whose status is not exactly `TODO` (after trimming spaces and a trailing
period), is skipped by the executor — check this first when a plan appears to have no pending work.

## Gap Register Sync

When `docs/aics/gap.md` is used as an input source (e.g. produced by `/arcdlc:examinate`), every `### ...` gap must
have a matching task block in `docs/aics/plan.md`.

The `plan.md` copy must preserve:

- The same task ID and heading.
- The same `WHAT`, `WHERE`, `WHY`, and `Acceptance` content.
- Executor metadata: `References` and `Status`.

Use `docs/aics/gap.md` as the evidence register and `docs/aics/plan.md` as the executable queue.

## Authoring Rules

1. Use unique `<TASK-ID>` values (for example: `WA240-VER-03`, `AIC-1`).
2. Keep exactly one task per `###` block.
3. Keep headings at level `###`. Do not prefix with `TODO` — the status is tracked by the `- Status:` line inside the block.
4. Keep section keys exact: `WHAT`, `WHERE`, `WHY`, `Acceptance`, `References`, `Status`.
5. Keep `Status` values uppercase: `TODO`, `TAKEN`, `DONE`, `BLOCKED`.
6. In `WHERE`, list exact files/modules expected to change.
7. In `Acceptance`, give at least one testable success criterion; `arctool validate --strict` fails a task
   with no (or an empty) `Acceptance` section.
8. In `References`, include all source docs required for implementation.
9. Prefer ending `Status` with a trailing period for consistency (`- Status: TODO.`).
10. Do not place executor instructions inside `plan.md`; update this file instead.

## Minimal Example

```md
### WA999-VER-01 (MISSING): Add endpoint parity for sample flow

- WHAT: Add `/v2/sampleapi/items/{id}` read endpoint with legacy-compatible response.
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
  - GIVEN the new handler WHEN the test suite runs THEN `item_test.go` covers the 200 and 404 paths.
- References: `docs/epics/initiative-99/aic.md`, `docs/epics/initiative-99/arc42.md`.
- Status: TODO.
```
