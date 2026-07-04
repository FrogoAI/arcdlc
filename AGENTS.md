# AGENTS.md

Guidance for AI coding agents working **on this repository**. (If you are looking for how to
*use* ArcDLC in another project, read [README.md](README.md) and the skills under `skills/`.)

## What this repo is

ArcDLC is two deliverables in one repo, and they share a contract:

1. **A skill bundle** (`skills/`, packaged by `.claude-plugin/`) — the `/arcdlc:*` delivery
   workflow (aic, policy, plan, examinate, execute, remove, archive) plus the `source-map` reference
   library.
2. **The `arctool` CLI** (`cmd/arctool`, `internal/plan`, `internal/registry`) — a deterministic
   runner for the plan format those skills produce and consume, plus the initiative registry sync.

The shared contract is `skills/plan/references/plan-format.md`. It is parsed mechanically by
`internal/plan`; treat it as an API, not prose.

## Initiatives

Active initiatives in this repo (kept in sync by `arctool sync`; do not edit inside the markers):

<!-- arcdlc:initiatives:begin -->
- [Initiative Lifecycle](docs/aics/initiative-lifecycle/aic.md) — Mandatory slug-first selection, an arctool-synced initiative registry, and an always-confirmed removal flow.
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
  command (`/arcdlc:<name>`) and as a flat skill (`arcdlc-<name>` on Codex/OpenCode). Keep the
  dual path references (`../plan/...` and `../arcdlc-plan/...`) intact when editing.
- **`arctool` is always optional in skills.** Every skill that uses it must probe
  `command -v arctool` and describe the manual fallback. Never make a skill hard-depend on the CLI.
- **Status mutations stay byte-preserving and atomic.** `take`/`done`/`block`/`todo` rewrite only
  the one `- Status:` line via temp-file + rename; `archive` writes the archive before compacting
  the plan. Preserve these invariants.
- **Initiatives are folders; selection is mandatory and explicit.** Each initiative lives in
  `docs/aics/<slug>/` (holding the architecture doc, `plan.md`, `gap.md`, `plan-archive.md` — the
  latter two are always siblings of `plan.md`). Selection is always named, never inferred:
  skills take the slug as their first positional argument (missing → error listing initiatives), and
  `arctool` requires `--aic <slug>` or `--plan PATH` (neither → lists initiatives, exit 2). The
  resolver lives in `cmd/arctool` (`resolvePlan`); keep the skills' manual fallback describing the
  same rule. The legacy flat `docs/aics/plan.md` is reachable only via `--plan`. Task IDs are unique
  per plan, not globally. ADRs (`docs/adr/`) and `CONTEXT.md` stay global, not per-initiative.
- **The initiative registry is generated.** `arctool sync` keeps the initiative list (title +
  summary, parsed from each arch doc's `# ` H1 and the `> ` blockquote under it, per `internal/registry`)
  inside the `<!-- arcdlc:initiatives -->` marker blocks in `AGENTS.md` and `README.md`, rewriting
  only that region (byte-preserving elsewhere, atomic). `sync --check` fails on drift for CI. Never
  hand-edit inside the markers. `/arcdlc:remove <slug>` deletes an initiative folder (always after an
  explicit confirmation) and re-syncs; `arctool` itself performs no deletion.
- **Version bumps:** the CLI version lives in `cmd/arctool/main.go` (`const version`); the plugin
  version lives in `.claude-plugin/plugin.json`. Bump whichever component you changed. Releases
  are cut by pushing a `v*` tag.

## Conventions

- One skill per directory under `skills/`, entry file always `SKILL.md`, YAML frontmatter with a
  `description` that names its triggers (the `/arcdlc:<name>` command and the `arcdlc-<name>`
  flat form).
- Reference documents belong in `skills/source-map/source/` and are routed via the table in
  `skills/source-map/SKILL.md` — add a row when adding a document.
- Adding or renaming a sub-skill requires updating the `SUBSKILLS` list in `install.sh` and the
  skill-layout / installer-smoke checks in `.github/workflows/ci.yml` in the same change set.
- Exit codes of `arctool` are part of its interface (0 ok, 1 contract failure, 2 usage, 3 not
  found/empty, 4 I/O, 5 archive self-validation) — skills key off them; do not renumber.
- `CLAUDE.md` is a symlink to this file; edit `AGENTS.md` only.
