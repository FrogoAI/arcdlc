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

Completed (archived to docs/aics/aic-tracking/plan-archive.md):
- TRK-1: Add the `arctool status` command
- TRK-2: Leave closed initiatives out of the registry and count them in one line
- TRK-3: Create the `/arcdlc:close` skill and register it as the twelfth skill
- TRK-4: Rewrite `/arcdlc:remove` to delete a design for good with a reference warning
- TRK-5: Add the closed-initiative check to `aic`, `plan`, `examinate` and `assist`
- TRK-6: Teach `/arcdlc:aic` to read closed designs and record a reversal in their `CLOSED.md`
- TRK-7: Update `AGENTS.md`, `README.md`, `CONTEXT.md`, the `init` chain line and the bundle version
