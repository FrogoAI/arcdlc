# Cursor Support

> Add Cursor as a supported agent via flat personal skills (~/.cursor/skills/arcdlc-<name>) — installer, CI, and docs only, no skill or arctool change.

Format: AIC (template: `source-map/source/AIC Template.md`). Status: Draft — decompose with
`/arcdlc:plan cursor-support`. Interview date: 2026-07-14.

## Goals

### 🟢 Business Case

ArcDLC's promise is "bring your own agent" — one install-agnostic bundle across coding agents. Today
that means Claude Code, Codex, and OpenCode (with Antigravity in flight). Cursor ships a skill system
that is shape-compatible with ArcDLC's, so users on Cursor currently get no supported install path
even though the skills would run there unchanged. This initiative closes that gap: it makes
`install.sh` install ArcDLC for Cursor and names Cursor as a first-class agent in the documentation,
keeping ArcDLC's cross-agent reach current with the market.

### 🟢 Functional Overview

1. **New installer agent `cursor`.** `install.sh` gains a `cursor` agent key, auto-detected when
   `~/.cursor` exists, included in `--agents all` and the explicit list, and handled in `--uninstall`.
2. **Flat install** ([ADR-0006](../../adr/0006-cursor-support-via-flat-personal-skills.md)). Reuse the
   Codex/OpenCode branch: flatten each sub-skill as `arcdlc-<name>` into `~/.cursor/skills/`. Cursor
   has no plugin namespace or plugin-install CLI for skills, so there is a single flat path — no
   prefer-plugin ladder (contrast Antigravity's ADR-0005).
3. **Version bump.** `.claude-plugin/plugin.json` 0.4.0 → 0.5.0 in the same change set. No new manifest
   (Cursor needs none), and `arctool` (0.8.0) is **not** bumped.
4. **Cosmetic doc mentions.** Enumerate Cursor as a supported agent in the README badge and philosophy
   line, the "bring your own agent" bullet, the Installation section (auto-detect list + a manual
   Cursor snippet), and the AGENTS.md install-agnostic note. No behavioral skill change.
5. **CI coverage.** Extend the installer-smoke test to create `$FAKE/.cursor`, run the installer, and
   assert `~/.cursor/skills/arcdlc-<name>/SKILL.md` present then removed on `--uninstall`.

Out of scope for this initiative (see Open Questions): Cursor **project-scope** skills
(`<root>/.cursor/skills/`), distributing ArcDLC via Cursor **rules**/`AGENTS.md`, and any `arctool`
change.

### 🟢 Quality Goals

1. **Non-regression / isolation** — the change is additive: the Claude/Codex/OpenCode/Antigravity
   install paths, the skills' content, `arctool`, and the plan contract are byte-for-byte unaffected.
   The new agent is one more branch, not a refactor of the installer.
2. **Consistency** — Cursor reuses the existing `arcdlc-<name>` flat layout and its dual-path
   references verbatim; the `cursor` branch is the Codex/OpenCode branch with a third destination, so
   there is one install mental model.
3. **Verifiability** — file placement (flat skills present, uninstall clean) is proven in CI via a
   fake `$HOME`. Because Cursor has no plugin-install CLI, there is no "preferred" path CI cannot
   reach (contrast Antigravity), so the smoke test guards the whole install contract.

### 🟢 Organizational Constraints

- Skills stay install-agnostic: no skill-content behavior change; the flat dual-path references
  (`../arcdlc-plan/...`, `../arcdlc-source-map/...`) must keep resolving under Cursor (they already do
  — Cursor uses the same `arcdlc-<name>` sibling layout).
- `arctool` remains optional and, here, entirely untouched — no new module dependency, still pure
  standard library, static binaries.
- No new sub-skill, so `SUBSKILLS` in `install.sh` and the CI skill-layout enumeration are unchanged;
  the changes that *are* required (installer branch, uninstall, smoke test, version bump, docs) land
  in the same change set.
- Version bump in the same change set: `.claude-plugin/plugin.json` 0.4.0 → 0.5.0. `arctool` (0.8.0)
  is **not** bumped, and no second manifest is introduced.

### 🟢 Technical Constraints

- `install.sh` stays POSIX-ish bash and must pass `shellcheck` (CI gate).
- CI can create a fake `~/.cursor` and assert flat file placement with no external CLI (Cursor has no
  plugin-install CLI), so the entire Cursor install path is CI-verifiable — no manual-only rung.
- The Cursor skills directory (`~/.cursor/skills/`) must be defined as a single named variable (one
  place to change) alongside the existing agent directory variables, not a scattered literal.

### 🟢 Business Context

System under construction: the ArcDLC bundle (skills + `arctool` + installer). This initiative adds
one communication partner — the Cursor agent — to the existing set:

```
 engineer ──/arcdlc:<skill> <slug>──▶ coding agent
                                       ├─ Claude Code   (plugin  → ~/.claude/skills/arcdlc | claude plugin)
                                       ├─ Codex         (flat    → ~/.codex/skills/arcdlc-<name>)
                                       ├─ OpenCode      (flat    → ~/.config/opencode/skills/arcdlc-<name>)
                                       ├─ Antigravity   (plugin | flat → ~/.gemini/...)
                                       └─ Cursor        (NEW, flat → ~/.cursor/skills/arcdlc-<name>)
 install.sh ──detect (~/.cursor)──▶ installs the flat Cursor skills
 CI ◀── installer-smoke (flat path) ── repo
```

## Architectural Hypotheses

### 🔵 H1 — Cursor reuses the flat `arcdlc-<name>` layout unchanged

- **Context:** Cursor skills need only `description` frontmatter; `name` defaults to the directory
  name; skills auto-load from `~/.cursor/skills/`. ArcDLC's flat installs already use `arcdlc-<name>`
  sibling directories with `../arcdlc-*/...` references.
- **Decision:** Do not introduce a new reference-path variant or edit any skill body. Cursor ships the
  existing `arcdlc-<name>` directories; commands come from the directory names.
- **Justification:** The flat layout is already agent-agnostic; reuse keeps one file set for
  Codex/OpenCode/Antigravity/Cursor and zero skill divergence.
- **Trade-offs:** `argument-hint` frontmatter (Claude-specific) is carried as inert extra metadata.

### 🔵 H2 — Cursor has no plugin mechanism → a single flat branch ([ADR-0006](../../adr/0006-cursor-support-via-flat-personal-skills.md))

- **Context:** Cursor auto-discovers skills from a directory; it has no plugin namespace and no
  plugin-install CLI for skills, so there is no way to (and no reason to) register twice.
- **Decision:** A single `cursor` key installs exactly one way — flatten `arcdlc-<name>` into
  `~/.cursor/skills/`. No prefer-plugin ladder.
- **Justification:** Simplest correct branch; literally the Codex/OpenCode branch with a third
  destination. Contrast Antigravity's ADR-0005, which needed a plugin-vs-flat fork.
- **Trade-offs:** Cursor treats the skills as model-invocable by default (no `disable-model-invocation`
  field) — accepted and desirable, since each description already names its triggers.

### 🔵 H3 — Auto-detect on `~/.cursor`; key named `cursor`; personal scope only

- **Context:** `~/.cursor` is Cursor's per-user home; personal skills live under `~/.cursor/skills/`.
- **Decision:** Name the key `cursor` and auto-detect on the presence of `~/.cursor`; install into the
  personal skills directory only. Project-scope (`.cursor/skills/`) is deferred.
- **Justification:** Matches the other agents' global-install model and the simplest reliable presence
  signal.
- **Trade-offs:** `~/.cursor` is present for IDE-only users (no CLI) — accepted, since skills are inert
  until invoked; project-shared skills are not installed by default.

### 🔵 H4 — No new manifest; version tracked in `.claude-plugin/plugin.json`

- **Context:** Cursor needs no plugin manifest (contrast `.antigravity-plugin/plugin.json`).
- **Decision:** Introduce no new manifest; bump the bundle version in `.claude-plugin/plugin.json`
  0.4.0 → 0.5.0.
- **Justification:** Fewer moving parts than the Antigravity branch; nothing to keep in lockstep.
- **Trade-offs:** None beyond the single existing version line.

### 🔵 H5 — CI proves the entire (flat) Cursor path

- **Context:** GitHub runners can create a fake `~/.cursor` and there is no external CLI to stand up.
- **Decision:** Extend the installer-smoke test to create `$FAKE/.cursor`, run the installer, assert
  `~/.cursor/skills/arcdlc-<name>/SKILL.md` present, then assert removal on `--uninstall`.
- **Justification:** Because Cursor has no plugin-install CLI, the flat path is the whole story — CI
  guards the complete install contract (contrast Antigravity, whose plugin path is manual-only).
- **Trade-offs:** None material; mirrors the existing Codex/OpenCode smoke assertions.

## Assessment

### 🔴 Technical Challenges & Risks

- **`shellcheck` regressions.** The new installer branch must pass the CI `shellcheck` gate.
  *Mitigation:* extend the existing `codex|opencode` branch to also handle `cursor` (a third
  destination), keeping new code minimal; run `shellcheck install.sh` locally before commit.
- **Volatile skills directory.** Cursor's skills path could change across versions.
  *Mitigation:* define `~/.cursor/skills/` via a single named variable alongside the existing agent
  directory variables, not a scattered literal.
- **Model-invocation surprise.** Users expecting explicit invocation may see skills auto-trigger.
  *Mitigation:* descriptions already gate their triggers; documented as expected behavior.
- **Broad `~/.cursor` detection.** May install for a user who has only the Cursor IDE (no CLI).
  *Mitigation:* skills are inert until invoked; documented.

### 🔴 Open questions

- Should the installer also wire Cursor **project-scope** skills (`<root>/.cursor/skills/`), not just
  the personal install? (Deferred; personal-only for now.)
- Should ArcDLC additionally ship guidance via Cursor **rules** or a project `AGENTS.md` block, or is
  the skill bundle sufficient? (Deferred; skills-only for now.)
- Should we set `disable-model-invocation` (or similar) to gate auto-invocation for Cursor, at the
  cost of skill-content divergence? (Deferred; keep skills byte-identical.)

## References

- [ADR-0006 — Cursor support via flat personal skills](../../adr/0006-cursor-support-via-flat-personal-skills.md)
- [ADR-0005 — Antigravity support via a native plugin bundle with a flat-skills fallback](../../adr/0005-antigravity-support-via-plugin-with-flat-fallback.md) — the parallel initiative; Cursor is the simpler, flat-only case.
- `install.sh` — the installer this initiative extends (the Codex/OpenCode flat branch is the model).
- `.github/workflows/ci.yml` — installer-smoke test to extend.
- `.claude-plugin/plugin.json` — the manifest whose version this initiative bumps.
- `CONTEXT.md` — ArcDLC ubiquitous language.
