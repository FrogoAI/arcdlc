# initiative-lifecycle — Plan

Task format contract: `skills/plan/references/plan-format.md`. Execute with
`/arcdlc:execute initiative-lifecycle` (one task, one commit, top-to-bottom).

Sequencing note: the arctool 0.6.0 per-folder work is uncommitted in this working tree and this plan
builds on top of it. Land 0.6.0 first, then work this queue. Tasks IL-1..IL-4 change `arctool`;
IL-5 updates the contract; IL-6..IL-9 update the skills; IL-10 adds the plan risk-coverage gate;
IL-11 updates docs and bumps versions.

## Risk Coverage

Every risk in the AIC's "Technical Challenges & Risks" and "Open questions" is covered by a task or
explicitly accepted (per ADR-0004; this plan dogfoods the gate IL-10 introduces):

- **In-flight collision** — accepted (process): the sequencing note above; land 0.6.0 first, don't interleave.
- **Editing user-owned root files** — covered by IL-3 (byte-preserving splice + tests).
- **Summary parsing on arbitrary agent-written docs** — covered by IL-2 (table-driven fallback tests, incl. arc42/TOGAF).
- **Breaking invocation change** — covered by IL-5 (contract), IL-6/IL-7 (skills), IL-11 (all docs slug-first).
- **Slug hygiene** — covered by IL-1 (`validSlug` retained and reused by `sync`/`remove`).
- **Open question — `arcdlc:policies` block later** — accepted/deferred (AIC H6): out of scope for this initiative.
- **Open question — CI runs `arctool sync --check` here** — covered by IL-11 acceptance (`sync --check` exits 0); wiring it as a CI gate is deferred.

### IL-1 (MISSING): arctool — make initiative selection mandatory (remove auto-detect)

- WHAT: Drop the auto-detect branch from `resolvePlan`. When neither `--plan` nor `--aic` is given,
  print an error that highlights the missing selection and lists the initiative slugs found under
  `docs/aics/`, then return exit code `2` (usage). Keep: explicit `--plan PATH` wins; `--aic <slug>`
  → `docs/aics/<slug>/plan.md`; an invalid slug → `2`. A well-formed `--aic <slug>` whose folder has
  no `plan.md` resolves to the path as today (downstream open yields not-found `3`). Update the
  `usage` const to state selection is required (drop the "auto-detects" wording) and keep exit-code
  docs 0/1/2/3/4/5 unchanged.
- WHERE:
  Layer `cli`: `cmd/arctool/main.go` (`resolvePlan`, the `usage` const; `slugOf`/`validSlug` retained
  and reused by later tasks).
  Tests: `cmd/arctool/main_test.go` (rewrite the auto-detect cases: single-initiative-no-flag now
  expects exit `2`; multi and none already expect `2`/`3`).
- WHY: ADR-0001 — implicit selection hides which plan a command mutates and silently changes behavior
  when a second initiative appears. Foundational: every other selection change keys off this rule.
- Acceptance:
  - GIVEN a repo with exactly one `docs/aics/<slug>/plan.md` WHEN `arctool next` runs with no `--aic`/`--plan` THEN it exits `2` and stderr lists the available slug(s).
  - GIVEN no initiative under `docs/aics/` WHEN a plan command runs with no flags THEN it exits `2` (usage), not `3`.
  - GIVEN `--aic payments` for an existing folder WHEN any plan command runs THEN it resolves `docs/aics/payments/plan.md`; GIVEN `--plan <path>` THEN that path is used verbatim.
  - GIVEN an invalid slug (`../x`, `a/b`) WHEN passed to `--aic` THEN it exits `2`.
  - GIVEN the change WHEN `go test ./...` runs THEN `cmd/arctool` tests pass and no test still asserts auto-detect.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0001-initiative-selection-is-always-explicit.md`, `skills/plan/references/plan-format.md`.
- Status: DONE.

### IL-2 (MISSING): arctool — initiative title/summary parsing (pure functions)

- WHAT: Add a new stdlib-only package that, given an initiative folder, derives its title and summary
  deterministically. Title = the first `# ` H1 of the architecture document, chosen by precedence
  `aic.md`, `arc42.md`, `togaf.md`, `c4.md`, else the first `*.md` alphabetically. Summary = the
  one-line `> ` blockquote immediately following that H1; fallback to the first non-empty paragraph
  truncated to ~120 chars; if the folder has no `.md` architecture doc, title = slug and summary =
  `(no architecture doc)`. No I/O beyond reading the folder; no third-party imports.
- WHERE:
  Layer `internal`: `internal/registry/registry.go` (new; `findArchDoc`, `parseTitle`, `parseSummary`,
  and an `Initiative` struct capturing slug/title/summary/docRelPath).
  Tests: `internal/registry/registry_test.go` (table-driven, using `t.TempDir()` fixtures).
- WHY: ADR-0002 / AIC H5 — `arctool sync` needs title+description with no LLM at sync time; the summary
  must live inside the document it describes (no sidecar). This is the parsing half, isolated so it is
  unit-testable before any file rewriting.
- Acceptance:
  - GIVEN a folder with `aic.md` starting `# Payments` and a `> one-liner` under it WHEN parsed THEN title=`Payments`, summary=`one-liner`.
  - GIVEN both `arc42.md` and `c4.md` present WHEN parsed THEN the doc precedence picks `arc42.md` over `c4.md`; GIVEN only `notes.md` THEN the first alphabetical `.md` is used.
  - GIVEN an H1 with no following blockquote WHEN parsed THEN summary = first paragraph truncated to ≤120 chars; GIVEN a folder with no `.md` THEN title=slug and summary=`(no architecture doc)`.
  - GIVEN the change WHEN `go test ./internal/registry/...` runs THEN all cases pass; `go vet ./...` and `gofmt -l .` are clean.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0002-registry-sync-via-marker-blocks.md`, `CONTEXT.md`.
- Status: DONE.

### IL-3 (MISSING): arctool — marker-block splice + stub creation (byte-preserving)

- WHAT: Add helpers that render the registry block body and splice it into a file between
  `<!-- arcdlc:initiatives:begin -->` and `<!-- arcdlc:initiatives:end -->`. Rules: rewrite ONLY the
  region between the markers, leaving every byte outside them identical; if the markers are absent,
  append an `## Initiatives` section containing them; if the file does not exist, create a minimal
  stub (H1 for README-style + the section); write atomically via temp-file + rename (reuse the pattern
  in `internal/plan/mutate.go`). Block body: one bullet per initiative sorted by slug —
  `- [<title>](docs/aics/<slug>/<docRelPath>) — <summary>` — or the single line `_none_` when empty.
- WHERE:
  Layer `internal`: `internal/registry/registry.go` (extend: `renderBlock`, `spliceBlock`,
  `ensureFile`/atomic write helper).
  Tests: `internal/registry/registry_test.go` (extend).
- WHY: ADR-0002 — arctool now edits two user-owned root files; the markers-only, atomic invariant must
  hold with the same rigor as the single-line status flips.
- Acceptance:
  - GIVEN a file with content before and after a marker block WHEN spliced THEN the bytes outside the markers are identical and only the inner region changes.
  - GIVEN a file without markers WHEN spliced THEN an `## Initiatives` section with both markers is appended and prior content is preserved.
  - GIVEN a missing file path WHEN spliced THEN a stub file is created containing the section; GIVEN an empty initiative set THEN the block body is `_none_`.
  - GIVEN an already-current file WHEN spliced again THEN the result is byte-identical (idempotent).
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0002-registry-sync-via-marker-blocks.md`, `internal/plan/mutate.go`.
- Status: DONE.

### IL-4 (MISSING): arctool — wire the `sync [--check]` subcommand

- WHAT: Add an `arctool sync [--check]` subcommand. It scans `docs/aics/*/` (each folder with any
  `.md` is an initiative), builds the `Initiative` list via IL-2, and updates `AGENTS.md` and
  `README.md` at the repo root via IL-3. Default mode writes the files (creating stubs as needed) and
  reports what changed. `--check` writes nothing and exits non-zero when either file's block is stale
  (drift), `0` when both are current. Register the subcommand in `main`'s dispatch and document it in
  the `usage` const. Scope is initiatives only (no policies block).
- WHERE:
  Layer `cli`: `cmd/arctool/main.go` (new `cmdSync`, dispatch case, `usage` entry).
  Tests: `cmd/arctool/main_test.go` (end-to-end in `t.TempDir()`: fixtures with 0/1/2 initiatives).
- WHY: AIC H4/H6 — initiatives MUST be tracked in AGENTS.md and README.md; a single deterministic
  command makes the registry drift-resistant and CI-enforceable.
- Acceptance:
  - GIVEN a temp repo with two initiative folders WHEN `arctool sync` runs THEN both `AGENTS.md` and `README.md` contain a marker block with both initiatives sorted by slug.
  - GIVEN the files are already current WHEN `arctool sync --check` runs THEN it exits `0` and writes nothing; GIVEN a folder is then added/removed THEN `--check` exits non-zero.
  - GIVEN `README.md` is missing WHEN `arctool sync` runs THEN a stub is created and the run reports it.
  - GIVEN no `docs/aics/*/` initiatives WHEN `arctool sync` runs THEN the block body is `_none_`.
  - GIVEN the change WHEN `go test ./...` runs THEN it is green and `arctool version` still prints.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0002-registry-sync-via-marker-blocks.md`, `skills/plan/references/plan-format.md`.
- Status: DONE.

### IL-5 (DRIFT): plan-format contract — mandatory selection + title/summary + registry

- WHAT: Rewrite the "Initiative folders" section of the format guide to match the new contract:
  selection is mandatory and explicit everywhere (skills take the slug as the first positional;
  `arctool` requires `--aic SLUG` or `--plan PATH`; no auto-detect); the legacy flat
  `docs/aics/plan.md` is reachable only via `--plan`. Document the architecture-doc contract
  (`# ` H1 = initiative title, one-line `> ` blockquote directly under it = summary) and the
  `arctool sync` registry marker blocks in `AGENTS.md`/`README.md`. Remove every "auto-detect"/
  "honored when it is the only plan" phrase.
- WHERE:
  Layer `docs`: `skills/plan/references/plan-format.md` ("Initiative folders" section; touch the
  minimal example only if a path needs it).
- WHY: The plan format is the mechanical contract; ADR-0001/0002 change it, so the guide must change in
  the same set as the code (IL-1..IL-4) that enforces it.
- Acceptance:
  - GIVEN the format guide WHEN read THEN the "Initiative folders" section states selection is mandatory/explicit and describes the slug-first (skills) and `--aic`/`--plan` (arctool) rules with no auto-detect wording.
  - GIVEN the guide WHEN read THEN it documents the H1-title + `> `-summary architecture-doc contract and the `arctool sync` marker-block registry.
  - GIVEN the guide WHEN `grep -i "auto-detect" skills/plan/references/plan-format.md` runs THEN there are no matches.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0001-initiative-selection-is-always-explicit.md`, `docs/adr/0002-registry-sync-via-marker-blocks.md`.
- Status: DONE.

### IL-6 (DRIFT): aic skill — slug-first mandatory + summary blockquote + post-create sync

- WHAT: Update `skills/aic/SKILL.md` so the initiative slug is the mandatory FIRST positional argument
  (`/arcdlc:aic <slug> [format]`); if it is missing, stop with an error that highlights the omission
  and lists existing initiatives under `docs/aics/`. Drop the `--aic` flag and the derive-slug-from-
  title path from the doc. Require the writer to emit a one-line `> ` summary blockquote directly under
  the document H1 (the format contract for `arctool sync`). After creating the folder, run
  `arctool sync` (probe `command -v arctool`; manual fallback: edit the `<!-- arcdlc:initiatives -->`
  blocks by hand). Update the frontmatter `description` and `argument-hint`, keeping dual-path
  references (`../source-map/...` and `../arcdlc-source-map/...`) intact.
- WHERE:
  Layer `skills`: `skills/aic/SKILL.md` (frontmatter + "Argument: initiative folder", Step 1/3/4).
- WHY: AIC H2/H5 — uniform slug-first invocation and the summary contract; the aic skill is where new
  initiatives (and their summaries) are born and where the registry is first populated.
- Acceptance:
  - GIVEN `skills/aic/SKILL.md` WHEN read THEN it specifies the slug as the mandatory first positional and the missing-slug error (highlight + list), with no `--aic` flag or slug-derivation path remaining.
  - GIVEN the skill WHEN read THEN Step 3 requires a one-line `> ` summary blockquote under the H1, and Step 4 runs `arctool sync` after folder creation with a `command -v arctool` probe and a manual marker-block fallback.
  - GIVEN the frontmatter WHEN read THEN `description`/`argument-hint` show `/arcdlc:aic <slug> [format]` and both `/arcdlc:aic` and `arcdlc-aic` triggers.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0001-initiative-selection-is-always-explicit.md`, `docs/adr/0002-registry-sync-via-marker-blocks.md`, `skills/plan/references/plan-format.md`.
- Status: TODO.

### IL-7 (DRIFT): plan/execute/examinate/archive skills — slug-first mandatory

- WHAT: In each of the four pipeline skills, make the initiative slug the mandatory FIRST positional
  argument, error+list on omission, and remove all auto-detect language. Keep the former positional
  after the slug: `/arcdlc:execute <slug> [TASK-ID]`, `/arcdlc:plan <slug> [format]`,
  `/arcdlc:examinate <slug> [policy]` (no policy → audit the initiative's own AIC), `/arcdlc:archive
  <slug>`. Skills still pass `--aic <slug>` to `arctool` (its CLI idiom) but no longer document a
  user-facing `--aic` flag. When a skill finds a legacy flat `docs/aics/plan.md`, instruct migrating it
  into `docs/aics/<slug>/`. Update each frontmatter `description`/`argument-hint`; keep dual-path
  references intact.
- WHERE:
  Layer `skills`: `skills/plan/SKILL.md`, `skills/execute/SKILL.md`, `skills/examinate/SKILL.md`,
  `skills/archive/SKILL.md` (frontmatter + each "Initiative selection"/argument section).
- WHY: AIC H2/H3 — one selection shape across the whole pipeline; flat layout demoted to the `--plan`
  escape hatch handled by `arctool`.
- Acceptance:
  - GIVEN each of the four SKILL.md files WHEN read THEN it shows the slug as mandatory first positional and the missing-slug error, with no "auto-detect"/"if several exist ask" wording remaining.
  - GIVEN `execute` and `examinate` WHEN read THEN the former positional follows the slug (`<slug> [TASK-ID]`, `<slug> [policy]`) and `examinate` with no policy audits the initiative's own AIC.
  - GIVEN any of the four WHEN it encounters a flat `docs/aics/plan.md` THEN the doc instructs migration into `docs/aics/<slug>/`.
  - GIVEN `grep -ri "auto-detect" skills/{plan,execute,examinate,archive}/SKILL.md` WHEN run THEN there are no matches.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0001-initiative-selection-is-always-explicit.md`, `skills/plan/references/plan-format.md`.
- Status: TODO.

### IL-8 (DRIFT): policy skill — require the policy name argument

- WHAT: Update `skills/policy/SKILL.md` so a missing `name` argument stops with an error that
  highlights the omission (mirroring the pipeline skills). Keep everything else as designed: file stays
  `docs/policies/<name>.md`, the `POL-<CLASS>-NNN` Unique ID stays in the header/index, and the
  three-place registration is unchanged (no `{code}-<name>.md` filename change). Update the frontmatter
  `description`/`argument-hint` to show the name as required.
- WHERE:
  Layer `skills`: `skills/policy/SKILL.md` (frontmatter + "Argument" section).
- WHY: AIC H8 — the engineer chose to keep the existing policy design and add only mandatory-name
  enforcement, consistent with slug-first strictness.
- Acceptance:
  - GIVEN `skills/policy/SKILL.md` WHEN read THEN the "Argument" section states the name is required and describes the missing-name error.
  - GIVEN the skill WHEN read THEN the output path is still `docs/policies/<name>.md` and no `{code}-<name>.md` renaming is introduced.
  - GIVEN the frontmatter WHEN read THEN `argument-hint` shows the name as required (not bracketed-optional).
- References: `docs/aics/initiative-lifecycle/aic.md`.
- Status: TODO.

### IL-9 (MISSING): /arcdlc:remove skill + install/CI registration

- WHAT: Add a new `remove` skill that deletes a completed initiative and cleans the registry, per
  ADR-0003. `skills/remove/SKILL.md` (with frontmatter naming both `/arcdlc:remove` and `arcdlc-remove`
  triggers) must: take the slug as the mandatory first positional (error+list if missing); report the
  initiative title, task counts by status (prefer `arctool list --aic <slug>`, fallback: read
  `plan.md`), and the file list; warn loudly if any task is not `DONE`; ALWAYS require explicit
  engineer confirmation (no skip flag); then delete `docs/aics/<slug>/` (`git rm -r` when tracked, else
  `rm -rf`) and run `arctool sync` (fallback: hand-edit the marker blocks). Keep dual-path references.
  Register the skill: add `remove` to `SUBSKILLS` in `install.sh` and to BOTH skill lists in
  `.github/workflows/ci.yml` (the "Check skill layout" loop and the "Installer smoke test" loop).
- WHERE:
  Layer `skills`: `skills/remove/SKILL.md` (new).
  Layer `install`: `install.sh` (`SUBSKILLS` line).
  Layer `ci`: `.github/workflows/ci.yml` (skill-layout loop and installer-smoke loop).
- WHY: AIC H7 / ADR-0003 — finished initiatives are context noise; removal must always pass a human and
  keep the registry truthful, while `arctool` stays non-destructive. Adding a sub-skill requires the
  matching `install.sh` + CI updates in the same change set (repo hard rule).
- Acceptance:
  - GIVEN `skills/remove/SKILL.md` WHEN read THEN it takes the slug as mandatory first positional, always requires confirmation, warns on non-`DONE` tasks, deletes the folder, and runs `arctool sync` with a `command -v arctool` probe + manual fallback.
  - GIVEN `install.sh` WHEN inspected THEN `SUBSKILLS` includes `remove`; GIVEN `.github/workflows/ci.yml` THEN both skill loops include `remove`.
  - GIVEN the installer smoke test WHEN run (`./install.sh` into a fake HOME) THEN `arcdlc-remove/SKILL.md` is present under the codex and opencode skill roots and uninstall removes it.
  - GIVEN the frontmatter WHEN read THEN its `description` names the `/arcdlc:remove` and `arcdlc-remove` triggers.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0003-initiative-removal-by-skill-not-arctool.md`, `install.sh`, `.github/workflows/ci.yml`.
- Status: TODO.

### IL-10 (MISSING): plan skill — risk-mitigation coverage gate

- WHAT: Add a mandatory risk-coverage gate to `skills/plan/SKILL.md`, running after decomposition
  (Step 2) and before final validation/handoff (Step 3). The skill must: read the architecture
  document's "Technical Challenges & Risks" and "Open questions" sections; for each risk determine
  whether it is covered by at least one plan task (or an explicit process mitigation) or consciously
  accepted/deferred with a rationale; if any risk is neither, run a grilling session with the
  engineer (invoke the `grilling` skill, or the inline discipline if absent) focused on the
  uncovered risks, then write each resulting mitigation into the plan — a new task when it needs
  implementation, or an accepted-risk note when it does not. Record the result as a `## Risk
  Coverage` mapping (risk → task IDs / accepted) in the plan preamble. The gate is hard: do not hand
  off until every risk is covered or accepted; it is a no-op when the document has no risks section
  (e.g. a gap-only plan). Note the `## Risk Coverage` preamble convention in the format guide (it is
  free-form preamble, not a `###` task block, so the runner ignores it). Keep dual-path references
  intact.
- WHERE:
  Layer `skills`: `skills/plan/SKILL.md` (new "Step 2.5 — Risk coverage gate" plus a Step 3 handoff
  precondition; frontmatter description touch if needed).
  Layer `docs`: `skills/plan/references/plan-format.md` (one short paragraph documenting the
  optional `## Risk Coverage` preamble section).
- WHY: ADR-0004 / AIC H9 — risks surfaced in the AIC interview must survive into the executable
  queue instead of silently evaporating ("audited every step, not assumed done").
- Acceptance:
  - GIVEN `skills/plan/SKILL.md` WHEN read THEN it defines a risk-coverage gate that reads the AIC's "Technical Challenges & Risks"/"Open questions", requires each risk to be covered by a task or explicitly accepted, and blocks handoff otherwise.
  - GIVEN a risk with no covering task WHEN the gate runs THEN the skill runs a `grilling` session (with an inline fallback when the skill is absent) and writes the agreed mitigation into the plan as a task or an accepted-risk note.
  - GIVEN the skill WHEN read THEN it requires a `## Risk Coverage` mapping in the plan preamble and states the gate is a no-op when the document has no risks section.
  - GIVEN `plan-format.md` WHEN read THEN it documents the optional `## Risk Coverage` preamble section as free-form (non-`###`) content the runner ignores.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0004-plan-enforces-risk-mitigation-coverage.md`, `skills/plan/references/plan-format.md`.
- Status: TODO.

### IL-11 (DRIFT): docs + registry blocks + version bumps

- WHAT: Bring the prose and versions in line with the new contract and cut the release surface.
  In `README.md`: make every command example slug-first (`See it in action`, `Quick Start`,
  `Commands` table, `Usage Examples`), remove auto-detect wording, rewrite the "Initiatives live in
  folders" section for mandatory explicit selection, add a `/arcdlc:remove` row to the command table,
  and add the `<!-- arcdlc:initiatives:begin/end -->` marker block. In `AGENTS.md`: update the
  "Initiatives are folders; selection is uniform" hard rule to describe mandatory explicit selection
  (no auto-detect), add the `remove` skill and `arctool sync` to the relevant hard rules/conventions,
  and add the marker block. Bump `const version` in `cmd/arctool/main.go` to `0.7.0` and
  `.claude-plugin/plugin.json` `version` to `0.3.0`. Rebuild arctool from source, then run
  `arctool sync` on this repo so the two marker blocks list `initiative-lifecycle`.
- WHERE:
  Layer `docs`: `README.md`, `AGENTS.md` (note: `CLAUDE.md` is a symlink to `AGENTS.md` — edit
  `AGENTS.md` only).
  Layer `cli`: `cmd/arctool/main.go` (`const version`).
  Layer `manifest`: `.claude-plugin/plugin.json` (`version`).
- WHY: The invocation change is breaking; every documented example must teach slug-first or the docs
  contradict the tools. Version bumps are required for the changed components (repo hard rule);
  populating the registry on this repo dogfoods `arctool sync`.
- Acceptance:
  - GIVEN `README.md` WHEN read THEN all `/arcdlc:*` examples are slug-first, no "auto-detect" wording remains, the command table includes `/arcdlc:remove`, and a populated `<!-- arcdlc:initiatives -->` block lists `initiative-lifecycle`.
  - GIVEN `AGENTS.md` WHEN read THEN the selection hard rule states mandatory explicit selection (no auto-detect) and mentions the `remove` skill and `arctool sync`, and it carries a populated marker block.
  - GIVEN the build WHEN `arctool version` runs THEN it prints `0.7.0`, and `.claude-plugin/plugin.json` `version` is `0.3.0`.
  - GIVEN this repo after rebuilding arctool WHEN `arctool sync --check` runs THEN it exits `0`; and `go build ./...`, `go test ./...`, `gofmt -l .`, `go vet ./...` are all clean.
- References: `docs/aics/initiative-lifecycle/aic.md`, `docs/adr/0001-initiative-selection-is-always-explicit.md`, `docs/adr/0002-registry-sync-via-marker-blocks.md`, `docs/adr/0003-initiative-removal-by-skill-not-arctool.md`.
- Status: TODO.
