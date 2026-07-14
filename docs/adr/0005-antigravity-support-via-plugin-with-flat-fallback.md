# ADR-0005 — Antigravity support via a native plugin bundle with a flat-skills fallback

- Status: Accepted
- Date: 2026-07-14
- Initiative: [antigravity-cli](../aics/antigravity-cli/aic.md)

## Context

ArcDLC ships as an install-agnostic skill bundle plus the optional `arctool` binary, today
targeting Claude Code, Codex, and OpenCode. Google's Antigravity CLI is a fourth agent whose
skill format is shape-compatible with ArcDLC's: a directory with `SKILL.md` (only `description`
is required in frontmatter; `name` defaults to the directory name) that auto-registers as a
slash command. Antigravity additionally has a **plugin** mechanism — a namespaced bundle
(`plugin.json` + `skills/`, `agents/`, `rules/`, `hooks.json`, `mcp_config.json`) installed via
`agy plugin install` and auto-crawled from its plugins directory.

Two install shapes are therefore possible for Antigravity, and doing both at once would register
the same `/arcdlc-*` commands twice (once from a loose global skills directory, once from the
plugin), causing duplicate/colliding commands.

This mirrors a decision already made for Claude Code, which prefers the official `claude plugin`
CLI and falls back to dropping the bundle into `~/.claude/skills/` (see `install.sh`).

## Decision

Add a single new installer agent key, **`antigravity`**, auto-detected when `~/.gemini` exists,
that installs exactly one way at a time, preferring the richer mechanism:

1. **Preferred — native plugin.** Stage the ArcDLC bundle as an Antigravity plugin
   (`.antigravity-plugin/plugin.json` + the `arcdlc-<name>` skill directories) into the plugins
   directory (`~/.gemini/antigravity-cli/plugins/arcdlc/`), via `agy plugin install` when `agy`
   is on PATH, otherwise by direct copy (the CLI auto-crawls the folder on startup).
2. **Fallback — flat skills.** When the plugin mechanism cannot be established, flatten each
   sub-skill as `arcdlc-<name>` into the shared global skills directory
   (`~/.gemini/config/skills/`), identical to the Codex/OpenCode layout.

Both shapes reuse the **same flattened `arcdlc-<name>` skill directories**, so the existing flat
dual-path references (`../arcdlc-plan/...`, `../arcdlc-source-map/...`) resolve unchanged. No new
reference-path variant, no skill-content behavior change, and `arctool` is untouched.

The Antigravity `plugin.json` lives in a new top-level `.antigravity-plugin/` directory, parallel
to `.claude-plugin/`, and is versioned in lockstep with the skill bundle.

## Justification

- **No duplicate registration.** Exactly one install path runs, so a given machine sees each
  `/arcdlc-*` command once.
- **Consistency with the existing design.** The prefer-rich-mechanism-else-drop-files ladder is
  structurally identical to the Claude branch, so the installer keeps one mental model.
- **Zero skill divergence.** Reusing the `arcdlc-<name>` layout means both Antigravity shapes and
  the Codex/OpenCode installs share one set of files and one set of relative references.

## Trade-offs

- **Undocumented mechanics accepted as risk.** Antigravity is young and community sources disagree
  on the exact plugins/skills paths and on plugin command-namespacing; the chosen paths and the
  `arcdlc-<name>` naming for command prefixes are verified manually against a live install, not in
  CI (CI cannot run Antigravity or `agy`).
- **A second plugin manifest to maintain.** `.antigravity-plugin/plugin.json` must be version-bumped
  alongside `.claude-plugin/plugin.json`; the AGENTS.md version-bump rule is extended to cover it.
- **Shared `~/.gemini`.** Because `~/.gemini/config/skills/` is shared with Gemini CLI, the fallback
  install also exposes `/arcdlc-*` to Gemini CLI — an intended, documented side effect.
