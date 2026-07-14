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

### AGY-1: Antigravity plugin manifest + bundle version bump

- WHAT: Add the committed Antigravity plugin manifest and bump the skill bundle to 0.6.0.
- HOW:
  Create `.antigravity-plugin/plugin.json` — valid JSON, keys `name` = `"arcdlc"`, `version` =
  `"0.6.0"`, `description` (one line naming ArcDLC and that it targets Antigravity CLI/IDE). Keep it
  minimal: rely on Antigravity auto-discovering the bundled `skills/` directory at install time
  (per Antigravity docs) — do NOT hand-list a `skills` array (avoids drift against `SUBSKILLS`).
  Bump `.claude-plugin/plugin.json` `"version"` `0.5.0` → `0.6.0` (0.5.0 was taken by the
  already-merged cursor-support initiative — CUR-3). Do NOT edit any `description` field: the
  manifests do not enumerate agents (the agent list lives in the README, handled by AGY-4), and
  `.claude-plugin/marketplace.json` has no version, so it is not touched here. Match the existing
  2-space JSON indentation.
  Out of scope: install.sh, CI, README, AGENTS.md (AGY-2/3/4); marketplace.json.
- WHERE:
  Manifests: `.antigravity-plugin/plugin.json` (new), `.claude-plugin/plugin.json`.
- WHY: The manifest is the install source for the `agy` path and the CI lint target; lockstep versioning is the drift guard (AIC H4).
- Acceptance:
  - GIVEN the repo WHEN `jq -e '.name=="arcdlc" and .version=="0.6.0"' .antigravity-plugin/plugin.json` runs THEN it exits 0.
  - GIVEN the repo WHEN `jq -e '.version=="0.6.0"' .claude-plugin/plugin.json` runs THEN it exits 0.
  - GIVEN the repo WHEN `[ "$(jq -r .version .antigravity-plugin/plugin.json)" = "$(jq -r .version .claude-plugin/plugin.json)" ]` runs THEN it exits 0 (versions equal).
- References: `docs/aics/antigravity-cli/aic.md`, `docs/adr/0005-antigravity-support-via-plugin-with-flat-fallback.md`.
- Status: DONE.

### AGY-2: install.sh — antigravity agent (detect, prefer-plugin-else-flat, uninstall)

- WHAT: Add an `antigravity` agent to `install.sh`: auto-detect, install (prefer `agy plugin install`, else flat), and uninstall.
- HOW:
  Add one path variable beside the existing agent dirs (~line 68–70): `gemini_dir="$HOME/.gemini"`.
  Derive the two Antigravity targets from it where used: plugin dir `"$gemini_dir/antigravity-cli/plugins/$PLUGIN"`,
  flat skills dir `"$gemini_dir/config/skills"`.
  `resolve_agents`: in the `auto)` branch append antigravity when present —
  `[ -d "$gemini_dir" ] && found="$found antigravity"`; in `all)` output
  `"claude codex opencode antigravity"`.
  Install loop: add a new `antigravity)` case that mirrors the Claude prefer-plugin-else-drop
  structure, gated by the SAME `ARCDLC_NO_PLUGIN_CLI` switch Claude uses:
    - Preferred rung: `if [ -z "${ARCDLC_NO_PLUGIN_CLI:-}" ] && command -v agy >/dev/null 2>&1` then
      assemble a staging bundle in a temp dir — copy `"$src/.antigravity-plugin/plugin.json"` to
      `"$stage/"`, `mkdir -p "$stage/skills"`, and for each `s` in `$SUBSKILLS`
      `cp -R "$src/skills/$s" "$stage/skills/$PLUGIN-$s"` — then `agy plugin install "$stage"`
      (redirect output; `rm -rf "$stage"` afterward). On success `info` that it installed via
      `agy plugin install` (commands `/$PLUGIN-<name>`).
    - Fallback rung (else, i.e. no `agy` or switch set): `root="$gemini_dir/config/skills"; mkdir -p "$root"`,
      and for each `s` in `$SUBSKILLS`: `rm -rf "${root:?}/$PLUGIN-$s"; cp -R "$src/skills/$s" "$root/$PLUGIN-$s"`;
      `info "antigravity: $root/$PLUGIN-<name> (flat fallback; invoke by skill name)"`.
  Uninstall block (~line 87–95): inside the existing `for s in $SUBSKILLS` loop also
  `rm -rf "$gemini_dir/config/skills/$PLUGIN-$s"`; after the loop `rm -rf "$gemini_dir/antigravity-cli/plugins/$PLUGIN"`
  (and best-effort `command -v agy >/dev/null 2>&1 && agy plugin uninstall "$PLUGIN" >/dev/null 2>&1 || true`).
  Keep the branch shellcheck-clean (quote expansions; `${var:?}` before `rm -rf` on derived paths).
  Out of scope: SUBSKILLS list is unchanged (no new sub-skill); CI (AGY-3); docs (AGY-4).
- WHERE:
  Installer: `install.sh` (path vars, `resolve_agents`, install `case`, uninstall block).
- WHY: Without an installer branch there is no supported way to install ArcDLC for Antigravity — the deliverable of this initiative.
- Acceptance:
  - GIVEN `FAKE=$(mktemp -d); mkdir -p "$FAKE/.gemini"` WHEN `HOME="$FAKE" ARCDLC_NO_PLUGIN_CLI=1 ./install.sh --agents antigravity --skills-only` runs THEN `"$FAKE/.gemini/config/skills/arcdlc-$s/SKILL.md"` exists for every `s` in `aic archive examinate execute plan policy remove source-map`.
  - GIVEN that install WHEN `HOME="$FAKE" ./install.sh --uninstall --skills-only` runs THEN `"$FAKE/.gemini/config/skills/arcdlc-aic"` and `"$FAKE/.gemini/antigravity-cli/plugins/arcdlc"` no longer exist.
  - GIVEN only `~/.gemini` present in a fake HOME WHEN `--agents auto` resolves THEN the antigravity flat install runs (its SKILL.md files appear), proving auto-detect includes it.
  - GIVEN the change WHEN `shellcheck install.sh` runs THEN it reports no findings.
- References: `docs/aics/antigravity-cli/aic.md`, `docs/adr/0005-antigravity-support-via-plugin-with-flat-fallback.md`, `install.sh`.
- Status: DONE.

### AGY-3: CI — installer-smoke fallback coverage + manifest lint

- WHAT: Extend `.github/workflows/ci.yml` to lint the Antigravity manifest and assert the flat-fallback install/uninstall.
- HOW:
  In the "Validate plugin manifests" step add:
  `jq -e '.name == "arcdlc"' .antigravity-plugin/plugin.json` and a version-equality assertion —
  `test "$(jq -r .version .antigravity-plugin/plugin.json)" = "$(jq -r .version .claude-plugin/plugin.json)"`.
  In the "Installer smoke test" step: add `"$FAKE/.gemini"` to the `mkdir -p` line so auto-detect
  selects antigravity; the existing run already sets `ARCDLC_NO_PLUGIN_CLI=1`, which forces the flat
  fallback (no `agy` on runners). After the existing codex/opencode assertions add, in the same
  `for s in ...` loop, `test -f "$FAKE/.gemini/config/skills/arcdlc-$s/SKILL.md"`. After the uninstall
  run add `test ! -e "$FAKE/.gemini/config/skills/arcdlc-aic"`.
  Leave the "Check skill layout" step unchanged (no new sub-skill).
  Out of scope: the `agy` preferred path (unreachable in CI — accepted per AIC H5).
- WHERE:
  CI: `.github/workflows/ci.yml` ("Validate plugin manifests" and "Installer smoke test" steps).
- WHY: Guards the manifest contract and the file-placement contract of the new install path against regression.
- Acceptance:
  - GIVEN the repo WHEN the "Validate plugin manifests" commands are run locally THEN the `.antigravity-plugin/plugin.json` name check and the version-equality check both exit 0.
  - GIVEN a local run of the extended smoke-test snippet (a fake `$HOME` with `.gemini`, `ARCDLC_NO_PLUGIN_CLI=1`, `./install.sh`) WHEN it executes THEN every `arcdlc-$s/SKILL.md` under `$FAKE/.gemini/config/skills/` is present, and absent after `--uninstall`.
  - GIVEN the workflow file WHEN parsed (e.g. `python -c 'import yaml,sys; yaml.safe_load(open(".github/workflows/ci.yml"))'`) THEN it is valid YAML.
- References: `docs/aics/antigravity-cli/aic.md`, `.github/workflows/ci.yml`, `install.sh`.
- Status: DONE.

### AGY-4: Docs — name Antigravity as a fourth agent

- WHAT: Add Antigravity to the human-facing agent enumerations and the AGENTS.md conventions/version rule.
- HOW:
  README.md — add `Antigravity` to: the agents badge (`agents-Claude%20Code%20·%20Codex%20·%20OpenCode…`
  → append `%20·%20Antigravity`), the philosophy line "agent-native across Claude Code, Codex, and
  OpenCode", the Quick Start "Installs the skills into every agent it detects (…)" line, and the
  Installation "detects which agents you have (…)" line. In the "Codex / OpenCode (manual)" section,
  document the Antigravity flat manual path (`skills_root=~/.gemini/config/skills`, same
  `arcdlc-<name>` copy loop) and mention `agy plugin install` as the native alternative.
  AGENTS.md — extend the install-agnostic hard rule ("…flat skill (`arcdlc-<name>` on Codex/OpenCode)")
  to include Antigravity; extend the version-bump hard rule to name `.antigravity-plugin/plugin.json`
  as the second plugin manifest kept in lockstep; add a one-line conventions note that the Antigravity
  manifest lives in `.antigravity-plugin/`. Do NOT edit `CLAUDE.md` (it is a symlink to `AGENTS.md`).
  Keep all edits to prose/strings — no behavioral claims beyond what AGY-1..3 implement.
  Out of scope: the `<!-- arcdlc:initiatives -->` registry blocks (owned by `arctool sync`); manifest
  files (AGY-1).
- WHERE:
  Docs: `README.md`, `AGENTS.md`.
- WHY: The docs currently teach three agents; without this they omit the newly supported one and contradict the installer.
- Acceptance:
  - GIVEN the repo WHEN `grep -ci antigravity README.md` runs THEN the count is ≥ 4 (badge, philosophy, quick-start/installation, manual section).
  - GIVEN the repo WHEN `grep -i "antigravity" AGENTS.md` runs THEN both the install-agnostic rule and the version-bump rule (naming `.antigravity-plugin/plugin.json`) match.
  - GIVEN the repo WHEN `git diff --name-only` is inspected THEN `CLAUDE.md` is not modified independently of `AGENTS.md` (symlink).
  - GIVEN the docs WHEN `gofmt -l .` / `go build ./...` run THEN they are unaffected (docs-only change compiles/formats clean).
- References: `docs/aics/antigravity-cli/aic.md`, `README.md`, `AGENTS.md`.
- Status: DONE.
