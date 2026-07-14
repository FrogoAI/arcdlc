### ANTIGRAVITY-CLI-GAP-01 (MISSING): Add Antigravity to plugin descriptions

- WHAT: Update the `description` fields in plugin manifests to enumerate Antigravity.
- HOW:
  Update the `description` in `.claude-plugin/plugin.json` and `.claude-plugin/marketplace.json` (and `.antigravity-plugin/plugin.json` if applicable) to mention that the bundle supports multiple agents, explicitly naming Antigravity (e.g. "for Claude Code, Codex, OpenCode, Cursor, and Antigravity").
- WHERE: `.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json`.
- WHY: "Enumerate Antigravity as a fourth agent in the ... .claude-plugin/plugin.json + marketplace.json descriptions" from `docs/aics/antigravity-cli/aic.md`.
- Acceptance:
  - GIVEN the plugin manifests WHEN `grep -i antigravity .claude-plugin/plugin.json` THEN it matches the updated description.
