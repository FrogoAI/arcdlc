# ADR-0006 — Cursor support via flat personal skills

- Status: Accepted
- Date: 2026-07-14
- Initiative: [cursor-support](../aics/cursor-support/aic.md)

## Context

ArcDLC ships as an install-agnostic skill bundle plus the optional `arctool` binary, today targeting
Claude Code, Codex, and OpenCode (with Antigravity in flight, see [ADR-0005](0005-antigravity-support-via-plugin-with-flat-fallback.md)).
Cursor is another coding agent whose skill format is shape-compatible with ArcDLC's: a directory with
`SKILL.md` (only `description` is required in frontmatter; `name` defaults to the directory name) that
Cursor auto-discovers. Cursor loads skills from `~/.cursor/skills/<name>/` (personal, all projects) and
`.cursor/skills/<name>/` (project, repo-shared).

Unlike Claude Code (`claude plugin`) and Antigravity (`agy plugin install`), Cursor has **no plugin
namespace and no plugin-install CLI for skills** — skills are simply dropped into a skills directory
and picked up. There is therefore no plugin-vs-flat fork to resolve, and no risk of double
registration.

## Decision

Add a single new installer agent key, **`cursor`**, auto-detected when `~/.cursor` exists, that
installs exactly like the existing Codex/OpenCode branch: flatten each sub-skill as `arcdlc-<name>`
into the personal skills directory (`~/.cursor/skills/`).

- **Same flattened `arcdlc-<name>` skill directories** as Codex/OpenCode/Antigravity, so the existing
  flat dual-path references (`../arcdlc-plan/...`, `../arcdlc-source-map/...`) resolve unchanged. No
  new reference-path variant, no skill-content change, and `arctool` is untouched.
- **No new manifest.** Cursor needs no plugin manifest (contrast `.antigravity-plugin/plugin.json`);
  the bundle version continues to live only in `.claude-plugin/plugin.json`.
- **Personal scope only** (`~/.cursor/skills/`). Project-scope installs (`.cursor/skills/`) are
  deferred (see the initiative's Open Questions).

## Justification

- **Zero skill divergence.** Cursor's skill model is identical in shape to the flat installs ArcDLC
  already supports, so reuse keeps one file set and one set of relative references.
- **Simplest correct branch.** With no plugin mechanism, Cursor is literally the Codex/OpenCode branch
  with a third destination directory — no prefer-rich-mechanism ladder (contrast ADR-0005), so there
  is no path CI cannot exercise.
- **No second manifest to maintain.** Nothing to version-bump in lockstep beyond the existing
  `.claude-plugin/plugin.json`.

## Trade-offs

- **Model-invocation by default.** Without a `disable-model-invocation` field, Cursor treats the
  skills as model-invocable (auto-triggering from the `description`). Accepted, and desirable: each
  ArcDLC description already names its triggers.
- **Claude-specific `argument-hint` frontmatter** is carried as inert extra metadata, exactly as in
  the other flat installs.
- **`~/.cursor` detection is broad.** It is present for users of the Cursor IDE even without the CLI;
  the skills are inert until invoked, so an install for an IDE-only user is harmless.
- **Personal-scope only** means repo-shared project skills are not installed by default; deferred as
  an open question.
