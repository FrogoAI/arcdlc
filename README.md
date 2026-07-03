# ArcDLC

[![CI](https://github.com/FrogoAI/arcdlc/actions/workflows/ci.yml/badge.svg)](https://github.com/FrogoAI/arcdlc/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**ArcDLC** (Architecture Development Life Cycle) is an open, agent-native delivery workflow:
a skill bundle for AI coding agents (Claude Code, Codex, OpenCode) that takes an initiative from
architecture interview to audited, task-by-task implementation — plus `arctool`, an optional
zero-dependency Go CLI that drives the executable plan mechanically.

ArcDLC is a universal delivery tool: it builds **applications** and authors **policies**, and both
feed the same audit → plan → execute machinery.

- **Application track:** `/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`
- **Governance track:** `/arcdlc:policy` → `/arcdlc:examinate docs/policies/<name>.md` → `/arcdlc:plan` → `/arcdlc:execute`

## Commands

| Command | What it does | Output |
| --- | --- | --- |
| `/arcdlc:aic [aic\|arc42\|togaf\|c4\|adr]` | Build the initiative's architecture document (AIC by default). Always runs a grilled interview first. | `docs/aics/<format>.md`, ADRs, `CONTEXT.md` |
| `/arcdlc:policy [name]` | Author a governance policy per the Policy of Policies framework — grilled interview first. | `docs/policies/<name>.md` + index |
| `/arcdlc:plan` | Decompose the approved architecture document into the executable task queue. | `docs/aics/plan.md` |
| `/arcdlc:examinate [policy]` | Examine existing code for compliance with a named policy or design (`MDCA`, `DDD`, `SOLID`, …; default: the project's own AIC) and register gaps as plan tasks. | `docs/aics/gap.md`, new TODO blocks in `docs/aics/plan.md` |
| `/arcdlc:execute [TASK-ID]` | Implement all pending plan tasks (or one by ID): status `TODO→TAKEN→DONE`, tests/lint, one commit per task. | code, tests, commits |
| `/arcdlc:archive` | Move `DONE` task blocks into `docs/aics/plan-archive.md`, keeping the plan small. | compacted plan + archive |
| `source-map` skill | Routing table into the bundled architecture & engineering reference library (AIC, arc42, TOGAF, C4, ADR, DDD, SOLID, MDCA, Go guides, Twelve-Factor, …). | reference guidance |

Every plan task carries testable `Acceptance` criteria; `/arcdlc:execute` must demonstrate them
before a task may be marked `DONE`. The full contract lives in
[`skills/plan/references/plan-format.md`](skills/plan/references/plan-format.md).

## Repository Layout

```
arcdlc/
├── .claude-plugin/          # plugin.json + marketplace.json (Claude Code plugin metadata)
├── skills/                  # one skill per directory (SKILL.md each)
│   ├── source-map/          # reference library (SKILL.md + source/)
│   ├── aic/  policy/  plan/  examinate/  execute/  archive/
│   └── plan/references/plan-format.md   # the executable-plan contract
├── cmd/arctool/               # arctool CLI entry point
├── internal/plan/           # plan parser, validator, mutator, archiver (+ tests)
├── Makefile                 # build / install / test / release
├── install.sh               # one-line installer (skills + arctool, all agents)
└── .github/workflows/       # CI (lint+test+cross-compile) and tag-driven releases
```

## Installation

### One-line install (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/FrogoAI/arcdlc/main/install.sh | bash
```

The installer:

- detects which agents you have (Claude Code, Codex, OpenCode) and installs the skills for each
  — via the official `claude plugin` CLI when available (Claude Code ≥ 2.1.157), otherwise into
  the agent's skills directory;
- installs the `arctool` binary to `~/.local/bin`: a checksum-verified release binary for
  **linux/amd64, linux/arm64, darwin/amd64, darwin/arm64**, falling back to a source build when
  Go is installed and no release binary is reachable;
- is idempotent — re-running upgrades everything in place.

Options go after `bash -s --` (or as flags to a local `./install.sh`):

```bash
... | bash -s -- --agents claude,codex   # explicit agent list (default: auto-detect)
... | bash -s -- --bindir ~/bin          # custom arctool location
... | bash -s -- --skills-only           # skip arctool
... | bash -s -- --tool-only             # skip the skills
... | bash -s -- --uninstall             # remove everything it installed
```

Piping scripts to bash requires trust — the script is short, dependency-free (`curl` + `tar`),
and worth the read: [`install.sh`](install.sh).

### Claude Code (manual)

Via the plugin marketplace — non-interactive from the shell (Claude Code ≥ 2.1.157), or with the
in-app `/plugin` equivalents:

```bash
claude plugin marketplace add FrogoAI/arcdlc
claude plugin install arcdlc@arcdlc
```

Commands appear namespaced as `/arcdlc:<name>`.

Alternative — clone into your skills directory; Claude Code auto-loads the plugin:

```bash
git clone https://github.com/FrogoAI/arcdlc ~/.claude/skills/arcdlc
```

### Codex / OpenCode (manual)

These agents have no plugin namespace, so install each sub-skill flattened as `arcdlc-<name>`
(identical behavior, invoked by skill name instead of a slash command):

```bash
git clone https://github.com/FrogoAI/arcdlc /tmp/arcdlc
skills_root=~/.codex/skills          # OpenCode: ~/.config/opencode/skills
mkdir -p "$skills_root"
for d in /tmp/arcdlc/skills/*/; do
  cp -r "$d" "$skills_root/arcdlc-$(basename "$d")"
done
```

### `arctool` CLI (optional, recommended)

`arctool` is the deterministic companion for `docs/aics/plan.md`: it validates the plan contract,
picks the next task, and flips task status atomically so the agent never hand-edits status lines.
It is pure Go standard library — the binaries are static and need no runtime.

Every ArcDLC skill probes `command -v arctool` and falls back to manual markdown handling when it
is absent, so the CLI is always optional.

```bash
# via go install
go install github.com/FrogoAI/arcdlc/cmd/arctool@latest

# or from a release: download the binary for your platform from
# https://github.com/FrogoAI/arcdlc/releases and put it on PATH
```

## Building from Source

Requires Go ≥ 1.22.

```bash
git clone https://github.com/FrogoAI/arcdlc
cd arcdlc
make build      # local binary at bin/arctool
make install    # install into ~/.local/bin (override with BINDIR=...)
make test       # go test ./...
make release    # static cross-compiled binaries in dist/ (linux/darwin × amd64/arm64)
```

CI runs `gofmt`, `go vet`, `go test`, validates the plugin manifests and skill layout, and
cross-compiles all platforms on every push and pull request. Pushing a `v*` tag builds the
release binaries with SHA256 checksums and publishes them as a GitHub release.

## Usage Examples

### End-to-end application flow

```
/arcdlc:aic                  # grilled interview → docs/aics/aic.md (+ ADRs, CONTEXT.md)
/arcdlc:plan                 # decompose the document → docs/aics/plan.md task queue
/arcdlc:execute              # implement every TODO task, one commit per task
/arcdlc:archive              # move DONE blocks to docs/aics/plan-archive.md
```

Run a single task, or audit an existing codebase:

```
/arcdlc:execute AIC-3        # implement only task AIC-3
/arcdlc:examinate MDCA       # audit code against MDCA, gaps become plan tasks
/arcdlc:aic arc42            # produce an arc42 document instead of an AIC
```

### Governance flow

```
/arcdlc:policy log-retention                       # grilled interview → docs/policies/log-retention.md
/arcdlc:examinate docs/policies/log-retention.md   # audit the repo against the policy
/arcdlc:plan                                       # fold gaps into the executable plan
/arcdlc:execute                                    # close the gaps task by task
```

### Driving the plan with `arctool`

A plan task block looks like this (full contract in
[`plan-format.md`](skills/plan/references/plan-format.md)):

```md
### AIC-1 (MISSING): Add health endpoint

- WHAT: Add `GET /healthz` returning build version.
- WHERE:
  Layer `handler`: `internal/handler/health.go`, `router.go`.
  Tests: `internal/handler/health_test.go`.
- WHY: Deploys are unverifiable without a liveness probe.
- Acceptance:
  - GIVEN a running server WHEN `GET /healthz` THEN the response is 200 with the build version.
- References: `docs/aics/aic.md`.
- Status: TODO.
```

```bash
arctool validate --strict      # enforce the contract (unique IDs, statuses, Acceptance per task)
arctool list                   # all tasks + status counts
arctool next --json            # first TODO block as JSON (exit 3 when queue is empty)
arctool take AIC-1             # claim it: TODO → TAKEN (refuses non-TODO without --force)
arctool done AIC-1             # complete it: TAKEN → DONE
arctool block AIC-1 -m "vendor API returns 500 on staging"
arctool archive --dry-run      # preview which DONE blocks would move to plan-archive.md
arctool archive                # move them (archive written first — crash-safe)
```

All mutations are guarded, byte-preserving, atomic rewrites of the single status line —
`arctool` never reformats your plan.

## License

[MIT](LICENSE) — © FrogoAI.
