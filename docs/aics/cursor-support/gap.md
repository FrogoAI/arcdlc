# Cursor Support — Gap Register

Audit date: 2026-07-14. Standard audited against: the initiative's own architecture
(`docs/aics/cursor-support/aic.md` + `docs/adr/0006-cursor-support-via-flat-personal-skills.md`).
Question: **do the ArcDLC skills load/run natively in Cursor, and is `arctool` usable there?**

## Rules checked (extracted from AIC H1–H3, H5 + ADR-0006)

- **R1 — Native skill discovery.** Cursor auto-discovers flat `arcdlc-<name>` skills from
  `~/.cursor/skills/`, accepting the SKILL.md frontmatter as-is (only `description` required; `name`
  defaults to the directory; description ≤ 1024 chars; name lowercase/hyphen ≤ 64).
- **R2 — Dual-path references resolve.** The flat in-body references (`../arcdlc-plan/...`,
  `../arcdlc-source-map/...`) resolve as siblings under the same skills root.
- **R3 — `arctool` optional and available.** Skills probe `command -v arctool` and degrade to a manual
  fallback; the binary is installed agent-agnostically on `PATH` and is therefore usable from Cursor's
  shell exactly as for the other agents.

## Findings — 0 violations (no MISSING / PARTIAL / DRIFT)

Both audited claims hold, with **live, on-machine evidence**:

- **R1 satisfied (native).** `~/.cursor/skills/` already contains all eight flat skills
  (`arcdlc-aic … arcdlc-source-map`), and this very `/arcdlc-examinate` invocation is executing from a
  flat `arcdlc-<name>` skill inside Cursor — direct proof the flat SKILL.md format loads and runs
  natively. Frontmatter conforms: no skill declares `name` (so it defaults to its directory, which is
  the intended install-agnostic behavior — hard-coding `name` would break the Claude-plugin layout
  where the directory is `aic`, not `arcdlc-aic`); every `description` is ≤ 559 chars (well under
  Cursor's 1024 limit, `skills/*/SKILL.md`); every derived name is lowercase/hyphen and ≤ 17 chars
  (< 64). `argument-hint` is carried as inert metadata (AIC H1 trade-off), not rejected.
- **R2 satisfied.** All eight `skills/*/SKILL.md` contain the `../arcdlc-plan/…` /
  `../arcdlc-source-map/…` dual-path references; under `~/.cursor/skills/` these are siblings, so they
  resolve unchanged.
- **R3 satisfied.** Seven of eight skills probe `command -v arctool` (source-map is a reference
  library and correctly uses none); `arctool` resolves on `PATH` at `~/.local/bin/arctool`, so it is
  usable from Cursor's terminal like any other agent. `arctool` remains untouched by this initiative.

## Result

**COMPLIANT.** The skills work natively in Cursor and `arctool` is supported there. No gap blocks and
no plan tasks are added. One accepted, non-gap behavior (already recorded in ADR-0006): with no
`disable-model-invocation` field, Cursor treats the skills as model-invocable (auto-trigger from
description) and invokes them as `arcdlc-<name>` rather than a `/arcdlc:<name>` plugin namespace.
