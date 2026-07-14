# Plan Archive

Archived task blocks from [plan.md](plan.md). Managed by `arctool archive`; do not hand-edit.

## Archived 2026-07-14

### CUR-1: Add the `cursor` agent to install.sh

- WHAT: Add a `cursor` installer agent key that flat-installs `arcdlc-<name>` into `~/.cursor/skills/`, with auto-detection and uninstall.
- HOW:
  Introduce `cursor_dir="$HOME/.cursor"` beside the existing `claude_dir`/`codex_dir`/`opencode_dir` declarations (never a scattered literal — the AIC Technical Constraints require one named variable).
  In `resolve_agents()`: under `auto`, append `cursor` when `[ -d "$cursor_dir" ]`; under `all`, output `claude codex opencode cursor`.
  Reuse the flat layout (ADR-0006, H1/H2): fold `cursor` into the existing flat branch that today reads `codex|opencode)`. Replace its `[ "$agent" = codex ] && root=... || root=...` two-way pick with a `case "$agent"` that maps `codex`→`$codex_dir/skills`, `opencode`→`$opencode_dir/skills`, `cursor`→`$cursor_dir/skills`; keep the rest of the branch (mkdir, per-`$SUBSKILLS` `rm -rf` + `cp -R` to `$root/$PLUGIN-$s`, info line) unchanged so behavior is identical to Codex/OpenCode.
  In the `--uninstall` block, add `"$cursor_dir/skills/$PLUGIN-$s"` to the existing per-`$SUBSKILLS` `rm -rf` list.
  Update the enumerations that name agents so they stay truthful: the `--agents LIST` header comment (line ~11), and the "no agent directories found (~/.claude, ~/.codex, ~/.config/opencode)" warn (line ~126) — add `~/.cursor`.
  Out of scope: Cursor project-scope `.cursor/skills/` (deferred), any docs/README/AGENTS text (CUR-4), CI (CUR-2), the version bump (CUR-3). Do not touch the Claude plugin branch or the `arctool` section. Antigravity is a separate initiative — do not add it to `all` here.
- WHERE:
  Installer: `install.sh` (agent dir vars, `resolve_agents`, the flat install branch, the uninstall loop, the two agent-enumerating comments).
- WHY: Without an installer branch, Cursor users have no supported install path — the whole point of the initiative.
- Acceptance:
  - GIVEN a fake `$HOME` with `.cursor/` present WHEN `HOME=$FAKE ./install.sh --agents cursor --skills-only` runs THEN `$FAKE/.cursor/skills/arcdlc-<name>/SKILL.md` exists for every name in `$SUBSKILLS`.
  - GIVEN those files exist WHEN `HOME=$FAKE ./install.sh --uninstall --skills-only` runs THEN `$FAKE/.cursor/skills/arcdlc-aic` no longer exists.
  - GIVEN a fake `$HOME` containing only `.cursor/` WHEN `resolve_agents` runs under `--agents auto` THEN `cursor` is selected (verifiable via the smoke test in CUR-2).
  - GIVEN the edited script WHEN `shellcheck install.sh` runs THEN it reports no findings.
- References: `docs/aics/cursor-support/aic.md`, `docs/adr/0006-cursor-support-via-flat-personal-skills.md`.
- Status: DONE.

### CUR-2: Cover the Cursor install path in CI

- WHAT: Extend the CI installer-smoke test to assert the Cursor flat install and its clean uninstall.
- HOW:
  In the "Installer smoke test" step, add `"$FAKE/.cursor"` to the `mkdir -p` that seeds fake agent homes (so auto-detect picks Cursor).
  In the per-subskill assertion loop that already checks `.codex` and `.config/opencode`, add `test -f "$FAKE/.cursor/skills/arcdlc-$s/SKILL.md"`.
  After the `--uninstall` invocation, alongside the existing post-uninstall assertions add `test ! -e "$FAKE/.cursor/skills/arcdlc-aic"`.
  No manifest lint is added (Cursor introduces no manifest — H4). Do not change the Claude assertions or the `arctool version` check.
  Out of scope: the `install.sh` logic itself (CUR-1) — this task only asserts it.
- WHERE:
  CI: `.github/workflows/ci.yml` (the `Installer smoke test` step: the `mkdir -p` line, the subskill assertion loop, the post-uninstall assertions).
- WHY: Cursor has no plugin CLI, so the flat path is the entire install contract — leaving it unguarded would let a regression ship silently (AIC H5).
- Acceptance:
  - GIVEN the workflow WHEN the smoke test runs with a fake `$HOME` seeded with `.cursor/` THEN it asserts `~/.cursor/skills/arcdlc-<name>/SKILL.md` present for the subskill list after install.
  - GIVEN the uninstall step ran WHEN the post-uninstall assertions run THEN they include a check that `~/.cursor/skills/arcdlc-aic` is absent.
- References: `docs/aics/cursor-support/aic.md`, `.github/workflows/ci.yml`.
- Status: DONE.

### CUR-3: Bump the bundle version to 0.5.0

- WHAT: Bump `.claude-plugin/plugin.json` `version` from `0.4.0` to `0.5.0`.
- HOW:
  Edit only the `version` field in `.claude-plugin/plugin.json` to `0.5.0`. Leave `.claude-plugin/marketplace.json` untouched (it carries no version). No new manifest is created (H4). `arctool` is not bumped (its version lives in `cmd/arctool/main.go` and is out of scope).
  Out of scope: any second manifest, `cmd/arctool/main.go`.
- WHERE:
  Manifest: `.claude-plugin/plugin.json`.
- WHY: The bundle's install surface changed (a new supported agent); the AGENTS.md version-bump rule requires bumping the component that changed.
- Acceptance:
  - GIVEN the edited manifest WHEN `jq -e '.version == "0.5.0"' .claude-plugin/plugin.json` runs THEN it exits 0.
  - GIVEN the edited manifest WHEN `jq -e '.name == "arcdlc"' .claude-plugin/plugin.json` runs THEN it still exits 0 (CI manifest lint unaffected).
- References: `docs/aics/cursor-support/aic.md`, `.claude-plugin/plugin.json`.
- Status: DONE.

### CUR-4: Name Cursor as a supported agent in the docs

- WHAT: Enumerate Cursor as a first-class supported agent across `README.md` and `AGENTS.md`, and add a manual Cursor install snippet.
- HOW:
  In `README.md`: add `Cursor` to the agents badge (the `agents-...` shields line), the philosophy block line (`agent-native across Claude Code, Codex, and OpenCode`), the "Bring your own agent" bullet, the Quick Start "Installs the skills into every agent it detects (…)" line, and the Installation "detects which agents you have (…)" bullet. Add a manual install snippet for Cursor — clone into `~/.cursor/skills/arcdlc-<name>` mirroring the existing Codex/OpenCode manual block (its `skills_root` is `~/.cursor/skills`); either extend the "Codex / OpenCode (manual)" heading to include Cursor or add a short sibling subsection. State the two accepted behaviors from ADR-0006: Cursor auto-triggers skills from their description (model-invocation), and `~/.cursor` presence is the detection signal.
  In `AGENTS.md`: extend the "Skills must stay install-agnostic" hard rule so `arcdlc-<name>` is named for `Codex/OpenCode/Cursor` (edit `AGENTS.md` only — `CLAUDE.md` is a symlink to it).
  Do NOT edit inside the `<!-- arcdlc:initiatives -->` marker blocks (already synced by `arctool sync`).
  Out of scope: installer logic (CUR-1), CI (CUR-2), version (CUR-3), Cursor rules / project `AGENTS.md` distribution (deferred).
- WHERE:
  Docs: `README.md` (badge, philosophy, BYO bullet, Quick Start, Installation bullet, manual install snippet), `AGENTS.md` (install-agnostic hard rule).
- WHY: "Bring your own agent" is a headline promise; an installed-but-undocumented agent leaves users unaware Cursor is supported.
- Acceptance:
  - GIVEN `README.md` WHEN grepped THEN the agents badge, philosophy line, BYO bullet, and Installation text each name `Cursor`, and a manual snippet references `~/.cursor/skills`.
  - GIVEN `AGENTS.md` WHEN grepped THEN the install-agnostic rule lists `Cursor` alongside Codex/OpenCode.
  - GIVEN the edits WHEN the `<!-- arcdlc:initiatives:begin -->…:end -->` region is compared THEN it is byte-unchanged (only prose outside the markers changed), verifiable via `./bin/arctool sync --check` exiting 0.
- References: `docs/aics/cursor-support/aic.md`, `docs/adr/0006-cursor-support-via-flat-personal-skills.md`, `README.md`, `AGENTS.md`.
- Status: DONE.
