# antigravity-cli — Plan

Task format contract: `skills/plan/references/plan-format.md`. Execute with
`/arcdlc:execute antigravity-cli` (one task, one commit, top-to-bottom).

Sequencing: AGY-1 creates the machine-readable manifests (both plugin.json files at the new bundle
version); AGY-2 teaches `install.sh` the `antigravity` agent (prefer `agy plugin install`, else flat
skills), which consumes AGY-1's manifest; AGY-3 extends CI (installer-smoke fallback + manifest
lint), which depends on AGY-1 and AGY-2; AGY-4 updates the human-facing docs. Dependency-ordered.

## Risk Coverage

Reconciliation of every risk in the AIC's "Technical Challenges & Risks" and "Open questions"
(per ADR-0004 / the /arcdlc:plan risk gate):

- **Volatile / contested paths** — covered by AGY-2: all Antigravity paths are defined as single
  named variables (`gemini_dir`, and derived plugin/skills dirs) with one place to change; **and
  accepted (process):** the *correctness* of those paths against a shipping Antigravity is confirmed
  by the manual live-verification note below (CI cannot).
- **Undocumented command namespacing** — accepted: command names come from the `arcdlc-<name>`
  directory prefix (→ `/arcdlc-<name>`), not a plugin namespace; the actually-registered names are
  confirmed in the manual live-verification note.
- **`agy plugin install` source semantics** — covered by AGY-2: the installer assembles a
  purpose-built staging bundle (root `plugin.json` + `skills/arcdlc-<name>/`) as the install source,
  so naming is under our control rather than inferred from the repo layout; **accepted:** that the
  `agy` rung actually registers commands is part of the manual live-verification note.
- **CI cannot exercise the preferred (agy) path** — accepted (AIC H5): CI proves only the flat
  fallback (AGY-3) via the existing `ARCDLC_NO_PLUGIN_CLI=1` force-flat switch; the `agy` path is
  the manual live-verification note.
- **`shellcheck` regressions** — covered by AGY-2 acceptance (`shellcheck install.sh` clean; existing
  CI gate).
- **Double manifest drift** — covered by AGY-1 (both manifests created/bumped to the same version in
  one commit) and AGY-3 (CI lints `.antigravity-plugin/plugin.json` and asserts its `version` equals
  `.claude-plugin/plugin.json`'s — this resolves the AIC open question "assert version equality" in
  favor of yes).
- **Open question — project-scope `.agents/skills/`** — accepted/deferred (out of scope, global
  install only).
- **Open question — README quick-start auto-detect / `agy plugin import gemini` migration docs** —
  covered in part by AGY-4 (manual Antigravity install instructions in README); the migration blurb
  is deferred until the `agy` path is verified live.

Manual live-verification (process, not a code task; perform before calling the initiative done):
in a real Antigravity install, run the installer, confirm `/arcdlc-*` commands register (via `agy`
when present, else from `~/.gemini/config/skills/`), and confirm uninstall clears them.

Completed (archived to docs/aics/antigravity-cli/plan-archive.md):
- AGY-1: Antigravity plugin manifest + bundle version bump
- AGY-2: install.sh — antigravity agent (detect, prefer-plugin-else-flat, uninstall)
- AGY-3: CI — installer-smoke fallback coverage + manifest lint
- AGY-4: Docs — name Antigravity as a fourth agent

### ANTIGRAVITY-CLI-GAP-01 (MISSING): Add Antigravity to plugin descriptions

- WHAT: Update the `description` fields in plugin manifests to enumerate Antigravity.
- HOW:
  Update the `description` in `.claude-plugin/plugin.json` and `.claude-plugin/marketplace.json` (and `.antigravity-plugin/plugin.json` if applicable) to mention that the bundle supports multiple agents, explicitly naming Antigravity (e.g. "for Claude Code, Codex, OpenCode, Cursor, and Antigravity").
- WHERE: `.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json`.
- WHY: "Enumerate Antigravity as a fourth agent in the ... .claude-plugin/plugin.json + marketplace.json descriptions" from `docs/aics/antigravity-cli/aic.md`.
- Acceptance:
  - GIVEN the plugin manifests WHEN `grep -i antigravity .claude-plugin/plugin.json` THEN it matches the updated description.
- References: `docs/aics/antigravity-cli/gap.md`, `docs/aics/antigravity-cli/aic.md`.
- Status: DONE.
