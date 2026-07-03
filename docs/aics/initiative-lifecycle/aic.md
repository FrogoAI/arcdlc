# Initiative Lifecycle

> Explicit slug-first initiative selection across all ArcDLC skills and arctool, an arctool-synced
> initiative registry in AGENTS.md/README.md, and an always-confirmed initiative removal flow.

Format: AIC (template: `source-map/source/AIC Template.md`). Status: Draft — decompose with
`/arcdlc:plan initiative-lifecycle`. Interview date: 2026-07-03.

## Goals

### 🟢 Business Case

arctool 0.6.0 (in the working tree, unreleased) made multiple initiatives per repo possible, but
left the lifecycle incomplete: selection is implicit (auto-detect), initiatives are undiscoverable
(nothing lists `docs/aics/*/`), and finished initiatives linger forever, costing agent focus and
context. This initiative closes the loop — create → track → remove — so a repo can run many
initiatives without ambiguity or accumulating noise.

### 🟢 Functional Overview

1. **Explicit selection.** The initiative slug is the mandatory first positional argument of every
   pipeline skill: `/arcdlc:aic <slug> [format]`, `/arcdlc:plan <slug>`,
   `/arcdlc:execute <slug> [TASK-ID]`, `/arcdlc:examinate <slug> [target]`,
   `/arcdlc:archive <slug>`. Missing slug → the skill stops with an error that highlights the
   omission and lists existing initiatives. arctool requires `--aic SLUG` or `--plan PATH` on every
   plan-addressing command (auto-detect removed): neither given → exit 2 + list of initiatives;
   unknown slug → exit 3.
2. **Registry.** New `arctool sync [--check]` maintains the list of initiatives (title, link,
   one-line summary) inside `<!-- arcdlc:initiatives:begin/end -->` marker blocks in `AGENTS.md`
   and `README.md`. Only the block is rewritten; missing files get minimal stubs; `--check` fails
   on drift.
3. **Removal.** New `/arcdlc:remove <slug>` skill: report task counts and files, warn loudly if
   anything is not DONE, always re-confirm with the engineer, delete `docs/aics/<slug>/`, run
   `arctool sync`.
4. **Policy alignment.** `/arcdlc:policy` keeps today's design (file `docs/policies/<name>.md`,
   `POL-<CLASS>-NNN` ID in header + index, three-place registration); the only change is that a
   missing name argument now stops with an error.
5. **Risk-mitigation gate in `/arcdlc:plan`.** After decomposition and before handoff, the plan
   skill checks that every risk in the architecture document's "Technical Challenges & Risks" (and
   "Open questions") is either covered by a task or explicitly accepted; uncovered risks trigger a
   grilling session with the engineer whose outcomes are written back into the plan, recorded as a
   Risk Coverage mapping.

### 🟢 Quality Goals

1. **Uniformity / determinism** — one selection rule everywhere (skills and CLI behave
   identically); exit-code meanings unchanged (0/1/2/3/4/5); registry content derived by a fixed
   parsing rule, never by judgment.
2. **Discoverability** — every initiative is visible from the repo front door (README.md) and the
   agent entry point (AGENTS.md); drift is CI-detectable via `sync --check`.
3. **Safety** — arctool writes stay byte-preserving and atomic (status lines; now also
   markers-only block rewrites via temp-file + rename); destruction lives only in a skill behind
   mandatory human confirmation; git history is the recovery path.

### 🟢 Organizational Constraints

- Skills stay install-agnostic: every change works as `/arcdlc:<name>` (plugin) and
  `arcdlc-<name>` (flat), keeping dual path references intact.
- arctool remains optional in skills: each skill probes `command -v arctool` and documents the
  manual fallback (including hand-editing the registry blocks).
- The new `remove` skill must be added to `SUBSKILLS` in `install.sh` and to the CI skill-layout /
  installer-smoke checks in the same change set.
- Version bumps in the same change set: `const version` → 0.7.0, `.claude-plugin/plugin.json` →
  0.3.0.
- Coordination: 0.6.0 is uncommitted in this working tree; this initiative modifies the same
  files. Land 0.6.0 first (or fold both into one 0.7.0 change) — do not interleave.

### 🟢 Technical Constraints

- arctool stays pure standard library, `CGO_ENABLED=0`, static release binaries.
- `skills/plan/references/plan-format.md` is a mechanical contract: the new H1 + `> ` summary
  rule and the selection rule change there, in `internal/plan`/`cmd/arctool` (+ tests), and in the
  skills — one change set.
- Exit codes are interface: 0 ok, 1 contract failure, 2 usage, 3 not found/empty, 4 I/O,
  5 archive self-validation. Missing selection maps to 2; unknown slug to 3. No renumbering.

### 🟢 Business Context

The system under construction is the ArcDLC bundle (skills + arctool). Communication partners:

```
 engineer ──/arcdlc:<skill> <slug>──▶ coding agent (Claude Code / Codex / OpenCode)
                                          │ runs SKILL.md flow
                                          ▼
                                       arctool ──reads/writes──▶ docs/aics/<slug>/{aic,plan,gap,plan-archive}.md
                                          │
                                          └──markers-only rewrite──▶ AGENTS.md / README.md (registry block)
 engineer ◀──confirmation prompts (remove; slug-missing errors)── skills
 CI ◀──arctool sync --check / validate── repo
```

## Architectural Hypotheses

### 🔵 H1 — Mandatory explicit selection ([ADR-0001](../../adr/0001-initiative-selection-is-always-explicit.md))

- **Context:** 0.6.0 auto-detects a single initiative; behavior shifts when a second appears.
- **Decision:** Slug/`--plan` required everywhere; auto-detect removed; errors list available
  initiatives.
- **Justification:** Predictable targeting; the command's blast radius is visible in the command.
- **Trade-offs:** Extra token in single-initiative repos; deletes freshly written 0.6.0 code.

### 🔵 H2 — Slug is always the first positional in skills; arctool keeps flags

- **Context:** `aic` used its positional for format, `execute` for TASK-ID, `examinate` for the
  audit target; 0.6.0 bolted on an `--aic` flag.
- **Decision:** `/arcdlc:<skill> <slug> [former-positional]` uniformly; `examinate` with no target
  audits the initiative's own AIC; skills drop the `--aic` flag from their docs; arctool keeps
  `--aic SLUG | --plan PATH` (flags are CLI idiom).
- **Justification:** One shape to remember; matches the natural reading
  (`/arcdlc:execute payments`).
- **Trade-offs:** Breaking change to documented skill invocations; all SKILL.md files change.

### 🔵 H3 — Legacy flat layout demoted to the `--plan` escape hatch ([ADR-0001](../../adr/0001-initiative-selection-is-always-explicit.md))

- **Context:** A flat `docs/aics/plan.md` has no slug, so it cannot satisfy mandatory selection.
- **Decision:** Skills detect a flat plan and instruct migration into `docs/aics/<slug>/`;
  `arctool --plan PATH` keeps any path workable; the AGENTS.md hard-rule text is updated.
- **Justification:** Two first-class layouts forever contradicts uniformity; nothing is bricked.
- **Trade-offs:** One-time manual `git mv` for existing flat repos.

### 🔵 H4 — Registry via `arctool sync` marker blocks ([ADR-0002](../../adr/0002-registry-sync-via-marker-blocks.md))

- **Context:** Initiatives MUST be tracked in AGENTS.md and README.md; manual registration drifts.
- **Decision:** `arctool sync [--check]` rewrites only the
  `<!-- arcdlc:initiatives:begin/end -->` region (temp-file + rename); appends an
  `## Initiatives` section when markers are absent; creates a minimal stub when a file is
  missing; bullet list sorted by slug: `- [<title>](docs/aics/<slug>/<doc>) — <summary>`; block
  holds `_none_` when empty. `aic` and `remove` skills run it after changing folders.
- **Justification:** Determinism and drift-resistance beat skill-side hand edits; markers bound
  the blast radius in user-owned files.
- **Trade-offs:** Hand edits inside markers are overwritten by design; a README.md stub in a repo
  without one may surprise (a note is printed).

### 🔵 H5 — Title/summary parsing contract

- **Context:** `sync` needs title + description with no LLM at sync time.
- **Decision:** Title = first `# ` H1 of the architecture doc (precedence: `aic.md`, `arc42.md`,
  `togaf.md`, `c4.md`, else first `*.md` alphabetically). Summary = one-line `> ` blockquote
  directly under the H1 — the `aic` skill must write it for every format it produces (this
  document dogfoods the rule). Fallbacks: first non-empty paragraph truncated to ~120 chars; no
  doc → slug + `(no architecture doc)`.
- **Justification:** Deterministic, testable, and the summary lives inside the document it
  describes — no sidecar metadata to drift.
- **Trade-offs:** New contract surface (parser + tests); pre-existing docs without the blockquote
  fall back to their first paragraph, which may read poorly until touched.

### 🔵 H6 — Sync covers initiatives only

- **Context:** Policies are registered by hand in three places today — the same drift risk.
- **Decision:** One `arcdlc:initiatives` block now; a policies block is explicitly deferred.
- **Justification:** Smallest change that satisfies the mandate; avoids coupling sync to the
  policy header format in the same release.
- **Trade-offs:** Policy registration keeps its manual drift risk for now.

### 🔵 H7 — Removal is a skill, arctool stays non-destructive ([ADR-0003](../../adr/0003-initiative-removal-by-skill-not-arctool.md))

- **Context:** Finished initiatives are noise; deletion must always pass a human; tree deletion
  does not fit arctool's write invariants.
- **Decision:** `/arcdlc:remove <slug>`: report title + task counts by status + files; loud
  warning when any task ≠ DONE; mandatory confirmation every time; `git rm -r` (or `rm -rf` if
  untracked); then `arctool sync`. No arctool delete command; git history is the archive.
- **Justification:** Keeps the deterministic tool auditable and small; confirmation UX belongs in
  the agent layer.
- **Trade-offs:** No one-shot CLI removal (intentional); unarchived history leaves the tree.

### 🔵 H8 — Policies unchanged except mandatory name

- **Context:** The interview considered `{code}-<name>.md` filenames and rejected them — keep the
  earlier policy design.
- **Decision:** `docs/policies/<name>.md` + `POL-<CLASS>-NNN` header IDs + three-place
  registration stay; `/arcdlc:policy` without a name stops with an error.
- **Justification:** The existing scheme already works and is indexed; only strictness was missing.
- **Trade-offs:** Filenames stay decoupled from IDs (the index maps them).

### 🔵 H9 — `/arcdlc:plan` enforces risk-mitigation coverage ([ADR-0004](../../adr/0004-plan-enforces-risk-mitigation-coverage.md))

- **Context:** Risks named in the AIC's "Technical Challenges & Risks" can be silently dropped
  during decomposition — the plan hands off with unmitigated, unacknowledged risks.
- **Decision:** After decomposing, `/arcdlc:plan` checks each risk (and open question) for coverage
  (a task) or explicit acceptance; uncovered risks trigger a grill with the engineer, and the
  mitigations are written into the plan (task or accepted-risk note) plus a Risk Coverage mapping.
  Hard gate: no handoff until every risk is covered or accepted; no-op when the doc has no risks
  section.
- **Justification:** "Audited every step, not assumed done" — risks surfaced in the interview must
  survive into the executable queue instead of evaporating.
- **Trade-offs:** `/arcdlc:plan` may pause for a grill; a vague risk can stall handoff until it is
  sharpened or explicitly accepted (intended friction).

## Assessment

### 🔴 Technical Challenges & Risks

- **In-flight collision:** 0.6.0 sits uncommitted in this working tree and this initiative edits
  the same files (`cmd/arctool/main.go`, all pipeline SKILL.md files, `plan-format.md`,
  `AGENTS.md`, `README.md`). Sequence the work: land 0.6.0, then execute this plan on top (or
  merge both into a single 0.7.0 release) — never interleave two agents in these files.
- **Editing user-owned root files:** the markers-only invariant must be enforced with the same
  rigor as status-line flips (tests for: markers present/absent/malformed, file missing, content
  outside markers byte-identical, `--check` behavior).
- **Summary parsing on arbitrary agent-written docs:** the fallback chain (blockquote → first
  paragraph → slug) needs table-driven tests, including arc42/TOGAF docs that predate the rule.
- **Breaking invocation change:** every documented example (`README.md`, skill descriptions,
  `argument-hint`) must move to slug-first in the same change set, or the docs will teach a shape
  that errors.
- **Slug hygiene:** slug validation (single kebab-case segment, no `/` or `..`) already exists in
  `resolvePlan`; it must survive the resolver rewrite since `sync` and `remove` also consume slugs.

### 🔴 Open questions

- Should `sync` later grow an `arcdlc:policies` block (deferred by H6)?
- Should this repo's own CI run `arctool sync --check` once `docs/aics/` exists here (dogfood)?

## References

- [ADR-0001 — Initiative selection is always explicit](../../adr/0001-initiative-selection-is-always-explicit.md)
- [ADR-0002 — arctool sync via marker blocks](../../adr/0002-registry-sync-via-marker-blocks.md)
- [ADR-0003 — Removal by skill, not arctool](../../adr/0003-initiative-removal-by-skill-not-arctool.md)
- [ADR-0004 — /arcdlc:plan enforces risk-mitigation coverage](../../adr/0004-plan-enforces-risk-mitigation-coverage.md)
- `skills/plan/references/plan-format.md` — the mechanical contract this initiative extends.
- `CONTEXT.md` — glossary for the terms used above.
