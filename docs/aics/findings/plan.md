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

Completed (archived to docs/aics/findings/plan-archive.md):
- FND-1: Add `Plan.Append` to `internal/plan` with the finding checks
- FND-2: Add the `arctool add` command
- FND-3: Document findings and `arctool add` in the plan format guide
- FND-4: Teach the executor to file findings in `/arcdlc:execute`
- FND-5: Add the findings review to the end of an `/arcdlc:execute` run
- FND-6: Teach `/arcdlc:plan` to act on reviewed findings
- FND-7: Update `AGENTS.md`, `README.md` and the bundle version for findings
