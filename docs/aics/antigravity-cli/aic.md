# Antigravity CLI

> Add Google Antigravity as a fourth supported agent — a native plugin bundle with a flat-skills fallback.

Format: AIC (template: `source-map/source/AIC Template.md`). Status: Draft — decompose with
`/arcdlc:plan antigravity-cli`. Interview date: 2026-07-14.

## Goals

### 🟢 Business Case

ArcDLC's promise is "bring your own agent" — one install-agnostic bundle across coding agents.
Today that means Claude Code, Codex, and OpenCode. Google's **Antigravity** (agentic IDE + CLI)
ships a skill/plugin system that is shape-compatible with ArcDLC's skills, so users on Antigravity
currently get no supported install path even though the skills would run there unchanged. This
initiative closes that gap: it makes `install.sh` install ArcDLC for Antigravity and updates the
documentation to name it as a first-class agent, keeping ArcDLC's cross-agent reach current with
the market.

### 🟢 Functional Overview

1. **New installer agent `antigravity`.** `install.sh` gains an `antigravity` agent key,
   auto-detected when `~/.gemini` exists, included in `--agents all` and the explicit list, and
   handled in `--uninstall`.
2. **Prefer-plugin-else-flat install** ([ADR-0005](../../adr/0005-antigravity-support-via-plugin-with-flat-fallback.md)).
   - *Preferred:* stage the ArcDLC bundle as an Antigravity plugin (a committed
     `.antigravity-plugin/plugin.json` + the `arcdlc-<name>` skill directories) into
     `~/.gemini/antigravity-cli/plugins/arcdlc/`, via `agy plugin install` when `agy` is on PATH,
     otherwise by direct copy (auto-crawled on startup).
   - *Fallback:* when the plugin mechanism cannot be established, flatten each sub-skill as
     `arcdlc-<name>` into `~/.gemini/config/skills/` — identical to the Codex/OpenCode layout.
3. **New Antigravity manifest.** A committed `.antigravity-plugin/plugin.json` (`name`, `version`,
   `description`, `skills`) parallel to `.claude-plugin/`, version-bumped in lockstep with the
   bundle.
4. **Cosmetic doc mentions.** Enumerate Antigravity as a fourth agent in the README badge and
   philosophy line, `.claude-plugin/plugin.json` + `marketplace.json` descriptions, and the
   AGENTS.md install-agnostic note. No behavioral skill change.
5. **CI coverage.** Extend the installer-smoke test to exercise the Antigravity **fallback** path
   (CI has no `agy`) and add a manifest lint for `.antigravity-plugin/plugin.json`.

Out of scope for this initiative (see Open Questions): Antigravity **project-scope** skills
(`<root>/.agents/skills/`), and any `arctool` change.

### 🟢 Quality Goals

1. **Non-regression / isolation** — the change is additive: Claude/Codex/OpenCode install paths,
   the skills' content, `arctool`, and the plan contract are byte-for-byte unaffected. The new
   agent is one more branch, not a refactor of the installer.
2. **Consistency** — Antigravity reuses the existing `arcdlc-<name>` flat layout and its dual-path
   references verbatim; the prefer-rich-else-drop-files ladder is structurally identical to the
   Claude branch, so there is one install mental model.
3. **Verifiability** — file placement (plugin bundle staged, fallback skills present, uninstall
   clean) is proven in CI via a fake `$HOME`; the parts CI cannot reach (`agy`, live command
   registration) are enumerated as explicit manual verification steps.

### 🟢 Organizational Constraints

- Skills stay install-agnostic: no skill-content behavior change; the flat dual-path references
  (`../arcdlc-plan/...`, `../arcdlc-source-map/...`) must keep resolving under both Antigravity
  shapes (they already do — Antigravity uses the same `arcdlc-<name>` sibling layout).
- `arctool` remains optional and, here, entirely untouched — no new module dependency, still pure
  standard library, static binaries.
- No new sub-skill, so `SUBSKILLS` in `install.sh` and the CI skill-layout enumeration are
  unchanged; the changes that *are* required (installer branch, uninstall, smoke test, manifest
  lint) land in the same change set.
- Version bump in the same change set: `.claude-plugin/plugin.json` 0.4.0 → 0.5.0, and the new
  `.antigravity-plugin/plugin.json` created at the matching bundle version (0.5.0). `arctool`
  (0.8.0) is **not** bumped. The AGENTS.md version-bump rule is extended to name the second
  manifest.

### 🟢 Technical Constraints

- `install.sh` stays POSIX-ish bash and must pass `shellcheck` (CI gate).
- The Antigravity `plugin.json` schema is taken from Antigravity's documented required fields
  (`name`, `version`, `description`; optional `skills`, `rules`); it must be valid JSON and lint
  under `jq` in CI.
- CI cannot install Antigravity or `agy`; the smoke test exercises only the flat-fallback file
  placement (forced, analogous to `ARCDLC_NO_PLUGIN_CLI=1` for Claude), never the `agy` path.
- Installer directory paths for Antigravity are community-reported and version-sensitive; they must
  be defined as single named variables (one place to change) with an override, not scattered
  literals.

### 🟢 Business Context

System under construction: the ArcDLC bundle (skills + `arctool` + installer). This initiative adds
one communication partner — the Antigravity agent — to the existing set:

```
 engineer ──/arcdlc:<skill> <slug>──▶ coding agent
                                       ├─ Claude Code   (plugin  → ~/.claude/skills/arcdlc | claude plugin)
                                       ├─ Codex         (flat    → ~/.codex/skills/arcdlc-<name>)
                                       ├─ OpenCode      (flat    → ~/.config/opencode/skills/arcdlc-<name>)
                                       └─ Antigravity   (NEW)
                                            ├─ preferred: plugin → ~/.gemini/antigravity-cli/plugins/arcdlc/
                                            │             (agy plugin install | direct copy)
                                            └─ fallback : flat   → ~/.gemini/config/skills/arcdlc-<name>
 install.sh ──detect (~/.gemini)──▶ chooses one Antigravity path
 CI ◀── installer-smoke (fallback path) + jq manifest lint ── repo
```

## Architectural Hypotheses

### 🔵 H1 — Antigravity reuses the flat `arcdlc-<name>` layout unchanged

- **Context:** Antigravity skills need only `description` frontmatter; `name` defaults to the
  directory name; skills auto-register as slash commands. ArcDLC's flat installs already use
  `arcdlc-<name>` sibling directories with `../arcdlc-*/...` references.
- **Decision:** Do not introduce a new reference-path variant or edit any skill body. Both
  Antigravity shapes ship the existing `arcdlc-<name>` directories; `/arcdlc-<name>` commands come
  from the directory names.
- **Justification:** The flat layout is already agent-agnostic; reuse keeps one file set for
  Codex/OpenCode/Antigravity and zero skill divergence.
- **Trade-offs:** `argument-hint` frontmatter (Claude-specific) is carried as inert extra metadata;
  bare command names depend on the `arcdlc-` directory prefix rather than a plugin namespace.

### 🔵 H2 — One agent key, prefer-plugin-else-flat ([ADR-0005](../../adr/0005-antigravity-support-via-plugin-with-flat-fallback.md))

- **Context:** Installing both a loose global skills copy and the plugin registers the same
  commands twice.
- **Decision:** A single `antigravity` key installs exactly one way — prefer the native plugin
  (`agy plugin install` or direct-stage the `.antigravity-plugin` bundle), fall back to flat skills
  in `~/.gemini/config/skills/` only when the plugin mechanism is unavailable.
- **Justification:** No duplicate registration; mirrors the Claude "prefer `claude plugin`, else
  drop into `~/.claude/skills/`" ladder exactly.
- **Trade-offs:** The installer branch is slightly more complex than Codex/OpenCode's single copy;
  the `agy` path is unverifiable in CI.

### 🔵 H3 — Auto-detect on `~/.gemini`; key named `antigravity`

- **Context:** `~/.gemini` is the shared Gemini-ecosystem home; `~/.gemini/config/skills/` is shared
  with Gemini CLI.
- **Decision:** Name the key `antigravity` (matching the request) and auto-detect on the presence of
  `~/.gemini`; document that the fallback install therefore also exposes `/arcdlc-*` to Gemini CLI.
- **Justification:** Matches user intent and the simplest reliable presence signal; the shared skills
  directory makes broader Gemini exposure a harmless, documented bonus.
- **Trade-offs:** May install for a user who has only Gemini CLI (no Antigravity) — accepted, since
  skills are inert until invoked; the key name slightly understates the shared reach.

### 🔵 H4 — Antigravity manifest lives in `.antigravity-plugin/`, versioned with the bundle

- **Context:** `.claude-plugin/plugin.json` already holds the Claude manifest and the bundle version.
- **Decision:** Commit `.antigravity-plugin/plugin.json` parallel to `.claude-plugin/`; create it at
  the current bundle version (0.5.0) and bump it in lockstep with `.claude-plugin/plugin.json`;
  extend the AGENTS.md version-bump rule to name it.
- **Justification:** Two clearly separated, in-repo, reviewable, CI-lintable manifests; lockstep
  versioning avoids a second version line drifting.
- **Trade-offs:** A second manifest to remember on every bundle bump (mitigated by the AGENTS.md rule
  and a `jq` CI lint).

### 🔵 H5 — CI proves the fallback path only; the plugin/`agy` path is manual

- **Context:** GitHub runners cannot install Antigravity or `agy`.
- **Decision:** Extend the installer-smoke test to create `$FAKE/.gemini`, run the installer forcing
  the flat fallback, and assert `~/.gemini/config/skills/arcdlc-<name>/SKILL.md` present then removed
  on uninstall; add a `jq` lint of `.antigravity-plugin/plugin.json`. The `agy plugin install` path
  and live `/arcdlc-*` registration are a documented manual checklist.
- **Justification:** Keeps CI deterministic and dependency-free while still guarding the file-placement
  contract; matches the existing Claude `ARCDLC_NO_PLUGIN_CLI=1` smoke pattern.
- **Trade-offs:** The preferred (plugin) path has no automated regression guard.

## Assessment

### 🔴 Technical Challenges & Risks

- **Volatile / contested paths.** Community sources disagree on Antigravity's skills and plugins
  directories (`~/.gemini/config/skills/` vs `~/.gemini/antigravity-cli/skills`; plugins under
  `~/.gemini/antigravity-cli/plugins/` vs `~/.gemini/config/plugins/`), and Antigravity is young.
  *Mitigation:* define each path as a single named variable with an override, document the chosen
  values, and confirm against a live install during execution; do not scatter path literals.
- **Undocumented command namespacing.** Whether plugin skills register namespaced (`/arcdlc:aic`) or
  bare (`/aic`) is not clearly documented; a real example plugin registers bare commands.
  *Mitigation:* rely on the `arcdlc-<name>` directory prefix (→ `/arcdlc-<name>`) rather than a
  namespace; verify the actual registered names manually before claiming the initiative done.
- **`agy plugin install` source semantics.** How `agy plugin install <src>` derives skill/command
  names (from the source `skills/` layout) is unconfirmed; it may not reproduce the `arcdlc-<name>`
  prefix if pointed at the repo's own `skills/aic|plan|...` directories.
  *Mitigation:* stage a purpose-built `.antigravity-plugin` bundle with `arcdlc-<name>` skill dirs
  (direct copy is authoritative); treat `agy plugin install` as the optional top rung.
- **CI cannot exercise the preferred path.** No `agy`/Antigravity on runners → the plugin path is
  unguarded. *Mitigation:* smoke-test the fallback; keep a written manual verification checklist.
- **`shellcheck` regressions.** New installer branches must pass the CI `shellcheck` gate.
  *Mitigation:* mirror the existing `codex|opencode` branch structure; run `shellcheck install.sh`
  locally before commit.
- **Double manifest drift.** Forgetting to bump `.antigravity-plugin/plugin.json` with the bundle.
  *Mitigation:* AGENTS.md rule update + `jq` lint (and optionally assert version equality with
  `.claude-plugin/plugin.json`).

### 🔴 Open questions

- Should the installer also wire Antigravity **project-scope** skills (`<root>/.agents/skills/`),
  not just the global install? (Deferred; global-only for now.)
- Should CI additionally assert that `.antigravity-plugin/plugin.json`'s `version` equals
  `.claude-plugin/plugin.json`'s, to enforce lockstep mechanically?
- Once verified live, should the README "Quick Start" auto-detect line and the `agy plugin import
  gemini` migration be documented for existing Gemini CLI users?

## References

- [ADR-0005 — Antigravity support via a native plugin bundle with a flat-skills fallback](../../adr/0005-antigravity-support-via-plugin-with-flat-fallback.md)
- `install.sh` — the installer this initiative extends (Claude prefer-plugin-else-flat branch is the model).
- `.github/workflows/ci.yml` — installer-smoke test and manifest lint to extend.
- `.claude-plugin/plugin.json` — the existing manifest; the new `.antigravity-plugin/plugin.json` mirrors its versioning.
- `CONTEXT.md` — ArcDLC ubiquitous language.
