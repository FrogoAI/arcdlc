# AGENTS.md

Guidance for AI coding agents working **on this repository**. (If you are looking for how to
*use* ArcDLC in another project, read [README.md](README.md) and the skills under `skills/`.)

## What this repo is

ArcDLC is two deliverables in one repo, and they share a contract:

1. **A skill bundle** (`skills/`, packaged by `.claude-plugin/`) — the `/arcdlc:*` delivery
   workflow (aic, policy, plan, plan-human, examinate, assist, execute, remove, archive), the
   `grilling` interview skill they all run on, plus the `source-map` reference library.
2. **The `arctool` CLI** (`cmd/arctool`, `internal/plan`, `internal/registry`, `internal/scan`) — a
   deterministic runner for the plan format those skills produce and consume, the initiative registry
   sync, and the code comment sweep (`arctool scan`) that writes the comment register.

The shared contract is `skills/plan/references/plan-format.md`. It is parsed mechanically by
`internal/plan`; treat it as an API, not prose.

## Initiatives

Active initiatives in this repo (kept in sync by `arctool sync`; do not edit inside the markers):

<!-- arcdlc:initiatives:begin -->
- [Antigravity CLI](docs/aics/antigravity-cli/aic.md) — Add Google Antigravity as a fourth supported agent — a native plugin bundle with a flat-skills fallback.
- [Cursor Support](docs/aics/cursor-support/aic.md) — Add Cursor as a supported agent via flat personal skills (~/.cursor/skills/arcdlc-<name>) — installer, CI, and docs only…
- [Initiative Lifecycle](docs/aics/initiative-lifecycle/aic.md) — Mandatory slug-first selection, an arctool-synced initiative registry, and an always-confirmed removal flow.
- [Task Ordering](docs/aics/ordering/aic.md) — Add an arctool command that re-orders task blocks in plan.md, because /arcdlc:execute runs them top to bottom.
- [Source Library Cleanup](docs/aics/source-library-cleanup/aic.md) — Make the bundled reference library agent-grade: redact leaked data, delete docs that contradict the plan contract, merge…
<!-- arcdlc:initiatives:end -->

## Build, test, verify

```bash
go build ./...            # must always compile
go test ./...             # must always be green
gofmt -l .                # must print nothing
go vet ./...              # must be clean
make build                # bin/arctool
make release              # dist/ binaries for linux/darwin × amd64/arm64
```

CI (`.github/workflows/ci.yml`) enforces all of the above plus plugin-manifest and skill-layout
checks. Do not merge with a red pipeline.

## Hard rules

- **`arctool` stays pure standard library.** Do not add module dependencies; release binaries must
  remain static (`CGO_ENABLED=0`).
- **The plan format is a contract.** Any change to `plan-format.md` requires matching changes in
  `internal/plan` (parser/validator/mutator/archiver), its tests, and the skills that reference
  the format — in the same change set.
- **Skills must stay install-agnostic.** Every SKILL.md must work both as a Claude Code plugin
 command (`/arcdlc:<name>`) and as a flat skill (`arcdlc-<name>` on Codex/OpenCode/Cursor/Antigravity). Keep the
 dual path references (`../plan/...` and `../arcdlc-plan/...`) intact when editing.
- **`arctool` is always optional in skills.** Every skill that uses it must probe
  `command -v arctool` and describe the manual fallback. Never make a skill hard-depend on the CLI.
- **The bundle depends on nothing outside itself.** A skill that delegates delegates to a *sibling
  skill in this bundle*, never to an external one. The interview is `skills/grilling` (`arcdlc-grilling`);
  `/arcdlc:aic`, `/arcdlc:policy`, and `/arcdlc:plan` name it and no other. Do not reintroduce
  references to third-party skills (`grilling`, `grill-with-docs`, `domain-modeling` on the user's
  machine) — that dependency is what `skills/grilling` exists to remove.
- **No skill hard-depends on another skill.** Even for a sibling, a delegating skill states a
  preference ladder that leads with the model-invocable skill and ends in an inline fallback (read
  the sibling's `SKILL.md` and run its protocol yourself), and never stops to report a helper skill
  as missing. What is mandatory is the behaviour (the grilled interview), never the invocation.
- **The interview is one question at a time.** `skills/grilling` asks exactly one question per turn,
  waits for the answer, then asks the next — never a numbered round. Every question carries a
  recommended answer. Seven skills restate this rule for their inline fallback (`aic`, `policy`,
  `plan`, `plan-human`, `examinate`, `assist`, `execute`); change it in all eight together.
- **Unclear is a question, not a guess.** A skill that hits a contradiction, a missing decision,
  code that looks dead, or a change that is risky or hard to reverse escalates to the grilled
  interview instead of picking for the engineer. `aic`, `policy` and `plan` grill up front; `examinate`
  (Step 2.5), `assist` (Step 3), `plan-human` and `execute` grill mid-flight. Every answer is recorded
  where the next session reads it (an ADR, `CONTEXT.md`, the gap's `HOW`, the story, the commit body),
  never only in the chat. An answer that changes the plan stops the run: `/arcdlc:plan` owns the plan
  text.
- **Everything the skills write is written for humans.** Every file every skill produces — the
  architecture documents of `/arcdlc:aic` in any format, Markdown or HTML, plus policies, ADRs,
  `CONTEXT.md`, plan tasks, gap blocks, and commit messages — is plain English: short active
  sentences, common words, expanded acronyms, no stacked noun phrases, no AI filler, no empty intensifiers
  (`honestly`, `genuinely`, `truly`, `clearly`, `obviously`), and no long dashes (`—`, `–`) in prose. Domain terms survive: a word like *provenance* is defined once in plain words, never
  swapped for a vaguer one. Plain words never mean less content: every decision, constraint,
  trade-off, and open question the template asks for still has to be there. The full standard is
  `skills/source-map/source/Writing Style.md`; the `## Talk simple, write like a human` block in every
  `SKILL.md` is its short form and must stay self-sufficient, because a skill is loaded but a
  reference is only read if the agent opens it. A long dash a format contract owns (the `— ` separator
  in generated initiative-registry lines) is exempt; that is `internal/registry`'s output, not prose.
- **Status mutations stay byte-preserving and atomic.** `take`/`done`/`block`/`todo` rewrite only
  the one `- Status:` line via temp-file + rename; `archive` writes the archive before compacting
  the plan; `order` permutes whole blocks, re-parses its own output before writing, and writes
  nothing on a failed check. `scan` is append-only on `comments.md` (a block is never
  rewritten or removed) and is the one command that edits source files: it deletes the marker comment
  lines it just registered, nothing else, never `plan.md`. It re-parses the register and re-checks
  every source edit before writing, and writes nothing on a failed check. Preserve these invariants.
- **The comment register has two owners.** `arctool scan` owns each block's task ID and its
  `- Marker:` line (the finding's identity: file plus marker text); `/arcdlc:assist` owns the title and
  every judgement key (`WHAT`, `HOW`, `WHY`, `Acceptance`, the verdict). Neither side writes the
  other's lines. The register is the only copy of a marker's text once the sweep has removed the
  comment, so a block is never deleted. See
  [ADR-0015](docs/adr/0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md).
- **Initiatives are folders; selection is mandatory and explicit.** Each initiative lives in
  `docs/aics/<slug>/` (holding the architecture doc, `plan.md`, `gap.md`, `comments.md`,
  `plan-archive.md` — the last three are always siblings of `plan.md`). Selection is always named, never inferred:
  skills take the slug as their first positional argument (missing → error listing initiatives), and
  `arctool` requires `--aic <slug>` or `--plan PATH` (neither → lists initiatives, exit 2). The
  resolver lives in `cmd/arctool` (`resolvePlan`); keep the skills' manual fallback describing the
  same rule. A folder is created lazily: `atomicWrite` makes the parent directory, so a write into a
  slug that has no folder yet (`arctool scan --aic <new-slug>`) lands and reports the folder it
  created. Creation is cheap and visible; deletion stays with `/arcdlc:remove`. The legacy flat
  `docs/aics/plan.md` is reachable only via `--plan`. Task IDs are unique per plan, not globally. ADRs (`docs/adr/`) and `CONTEXT.md` stay global, not per-initiative.
- **The initiative registry is generated.** `arctool sync` keeps the initiative list (title +
  summary, parsed from each arch doc's `# ` H1 and the `> ` blockquote under it, per `internal/registry`)
  inside the `<!-- arcdlc:initiatives -->` marker blocks in `AGENTS.md` and `README.md`, rewriting
  only that region (byte-preserving elsewhere, atomic). `sync --check` fails on drift for CI. Never
  hand-edit inside the markers. `/arcdlc:remove <slug>` deletes an initiative folder (always after an
  explicit confirmation) and re-syncs; `arctool` itself performs no deletion.
- **Version bumps:** the CLI version lives in `cmd/arctool/main.go` (`const version`); the plugin
  version lives in `.claude-plugin/plugin.json`, with `.antigravity-plugin/plugin.json` as a second
  plugin manifest kept in lockstep with it. Bump whichever component you changed (both plugin
  manifests together). Releases are cut by pushing a `v*` tag.

## Conventions

- One skill per directory under `skills/`, entry file always `SKILL.md`, YAML frontmatter with a
  `description` that names its triggers (the `/arcdlc:<name>` command and the `arcdlc-<name>`
  flat form).
- **Every `SKILL.md` carries the same `## Talk simple, write like a human` block**, verbatim, placed
  after the intro and before the first step. It covers both audiences in one place: how the agent
  talks to the user while the skill runs (plain English, bullets, no filler) and how it writes the
  files the skill produces (no AI filler, no long dashes, varied rhythm, concrete facts, domain terms
  kept and defined once). Brevity applies to replies, never to files: the block never lets a rule,
  path, decision, or acceptance criterion be dropped. Keep it short but self-sufficient — every rule
  an agent must follow stays inline, and only the tables, examples, and the pre-save grep live in the
  long form, `skills/source-map/source/Writing Style.md`. CI checks that all eleven blocks are byte
  identical, so copy the block when adding a skill, and change all eleven plus the long form
  together.
- Reference documents belong in `skills/source-map/source/` and are routed via the table in
  `skills/source-map/SKILL.md` — add a row when adding a document.
- Adding or renaming a sub-skill requires updating the `SUBSKILLS` list in `install.sh` and the
  skill-layout / installer-smoke checks in `.github/workflows/ci.yml` in the same change set.
- Exit codes of `arctool` are part of its interface (0 ok, 1 contract failure, 2 usage, 3 not
  found/empty, 4 I/O, 5 self-validation) — skills key off them; do not renumber.
- The Antigravity plugin manifest lives in `.antigravity-plugin/` (alongside the Claude Code manifest
  in `.claude-plugin/`).
- `CLAUDE.md` is a symlink to this file; edit `AGENTS.md` only.
