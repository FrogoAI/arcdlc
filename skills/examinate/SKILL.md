---
description: Examine existing code for compliance with a named architecture, policy, or design (e.g. /arcdlc:examinate MDCA — also DDD, SOLID, Clean Code, Go Server, Twelve-Factor, ECS), with a project policy authored by /arcdlc:policy (e.g. /arcdlc:examinate docs/policies/log-retention.md), or with the project's own AIC. Records violations as gap blocks in docs/aics/gap.md and adds matching TODO tasks to docs/aics/plan.md. Use when the user runs /arcdlc:examinate, invokes arcdlc-examinate, or asks for a compliance audit / gap analysis of the codebase.
argument-hint: "[MDCA|DDD|SOLID|arc42|...|(default: project AIC)]"
---

# ArcDLC Examinate (/arcdlc:examinate)

Audit the existing codebase against a policy or design, register every gap as evidence, and feed the gaps into the
executable plan so `/arcdlc:execute` can close them.

## Step 1 — Resolve the standard to audit against

- With an argument (e.g. `/arcdlc:examinate MDCA`): look the policy up in the sibling `source-map` skill's table
  (from this file: `../source-map/source/` in the plugin layout, `../arcdlc-source-map/source/` in flat installs)
  and read every listed reference in full — e.g. MDCA → `mdca.md` +
  `mdca_standard.md`; DDD → `ddd.md`; SOLID → `solid.md`; Go architecture → `Go Server.md` / `Go Client.md` /
  `Go Library.md`.
- With a project policy path (e.g. `/arcdlc:examinate docs/policies/log-retention.md`, typically one authored by
  `/arcdlc:policy`): read that policy in full and extract its Allowed/Prohibited rules as the checkable rule set.
- Without an argument: audit against the project's own architecture — `docs/aics/aic.md` (or other docs in
  `docs/aics/`), `docs/adr/`, and `CONTEXT.md`. If none of these exist either, stop and ask which policy to audit
  against (or suggest `/arcdlc:aic` first).
- Extract the concrete, checkable rules from the reference before looking at code, so findings cite a rule, not a
  feeling.

## Step 2 — Examine the code

- Sweep the codebase systematically against each rule: package/module layout, layer boundaries, dependency
  direction, naming, error handling, tests — whatever the policy governs.
- Every finding needs evidence: `file:line` (or package/module) plus the rule it violates. No evidence, no gap.
- Classify each finding's source status: `MISSING` (required element absent), `PARTIAL` (present but incomplete),
  or `DRIFT` (present but violates the policy).
- Do not report style preferences that the policy does not actually mandate.

## Step 3 — Write the gap register

Write or update `docs/aics/gap.md` — the evidence register. One gap per `###` block, using the plan heading format
so it can be mirrored into the plan verbatim:

```md
### <PREFIX>-GAP-NN (<MISSING|PARTIAL|DRIFT>): <Short Title>

- WHAT: <What must change to become compliant.>
- WHERE: <Files/modules with the violation, with file:line evidence.>
- WHY: <The violated rule, quoted or paraphrased, with the policy source path.>
- Acceptance:
  - GIVEN the audited code WHEN re-examined against <rule> THEN the violation is gone (no MISSING/PARTIAL/DRIFT finding).
```

Derive the `Acceptance` criterion from the violated rule, so "compliant" is testable rather than asserted — this
also lets the mirrored plan task pass `arctool validate --strict` (which requires an `Acceptance` section).

- `<PREFIX>` is the audited policy or initiative (e.g. `MDCA-GAP-01`, `AIC-GAP-03`). Number gaps sequentially,
  continuing from existing entries — never renumber or delete previous gaps.
- If a gap from an earlier examination is now fixed, mark it in `gap.md` (e.g. append `— resolved <date>`) instead of
  removing it.

## Step 4 — Sync gaps into the plan

Per the Gap Register Sync rules in `../plan/references/plan-format.md` (flat installs:
`../arcdlc-plan/references/plan-format.md`), append a matching task block to `docs/aics/plan.md` for every new gap:

- Same task ID and heading as in `gap.md`; same `WHAT`, `WHERE`, `WHY`, and `Acceptance` content.
- Add runner metadata: `References` (must include `docs/aics/gap.md`, the policy source, and the architecture doc if
  relevant) and `- Status: TODO.`
- Append after the existing blocks; never modify existing tasks or reuse an ID already present in the plan.
- If `docs/aics/plan.md` does not exist yet, create it per the format guide (a `/arcdlc:plan` run can merge it with
  architecture-driven tasks later).
- After updating the plan, validate it. Prefer `arctool validate --strict --plan docs/aics/plan.md` (probe once with
  `command -v arctool`, or install it from the arcdlc repo root: `make install`); fix any findings before handoff.
  If `arctool` is unavailable, say so once and hand-check unique IDs, present/uppercase `Status`, and required keys per
  `../plan/references/plan-format.md`.

## Step 5 — Report

Summarize: policy audited, rules checked, gaps found by status (MISSING/PARTIAL/DRIFT) with severity, and how many
tasks were added to the plan. Suggest `/arcdlc:execute` to start closing them.
