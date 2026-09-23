# AGENTS.md

Guidance for AI coding agents working **on this repository**. (If you are looking for how to
*use* ArcDLC in another project, read [README.md](README.md) and the skills under `skills/`.)

## What this repo is

ArcDLC is two deliverables in one repo, and they share a contract:

1. **A skill bundle** (`skills/`, packaged by `.claude-plugin/`) — the `/arcdlc:*` delivery
   workflow (init, aic, policy, plan, plan-human, examinate, assist, execute, remove, archive) plus the
   `grilling` interview skill they all run on. Eleven skills. Each owns the reference documents it
   consumes, in its own `references/` folder.
2. **The `arctool` CLI** (`cmd/arctool`, `internal/plan`, `internal/registry`, `internal/scan`) — a
   deterministic runner for the plan format those skills produce and consume, the initiative registry
   sync, and the code comment sweep (`arctool scan`) that writes the comment register.

The shared contract is `skills/plan/references/plan-format.md`. It is parsed mechanically by
`internal/plan`; treat it as an API, not prose.

Every decision that shaped either deliverable is in [docs/adr/](docs/adr/README.md). Read the ADR
before changing a rule that links to one.

## Initiatives

Active initiatives in this repo (kept in sync by `arctool sync`; do not edit inside the markers).
Decided questions and the reasoning behind them live in [docs/backlog.md](docs/backlog.md); read it
before re-opening one:

<!-- arcdlc:initiatives:begin -->
- [ArcDLC Init: one command that sets up a repository or a multi-repository workspace](docs/aics/init/aic.md) — Add `/arcdlc:init`, which detects whether it stands in one repository or in a workspace of several, scaffolds the files…
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

- **`arctool` stays pure standard library.** No module dependencies; release binaries stay static
  (`CGO_ENABLED=0`).
- **The plan format is a contract.** A change to `skills/plan/references/plan-format.md` requires
  matching changes in `internal/plan` (parser, validator, mutator, archiver), its tests, and every
  skill that references the format, in the same change set.
- **A planned task must be mechanical.** Executable by a model with no context but the block and the
  files it names: no judgement call, no open decision. `HOW` records decisions, never code. A task
  that cannot be made mechanical is too big or rests on an unfinished design, and goes back to
  `/arcdlc:plan` or `/arcdlc:aic`. It is never a reason to raise the executor tier, which only hides
  the defect. See [ADR-0021](docs/adr/0021-a-planned-task-must-be-mechanical.md).
- **A plan records the design it came from.** `/arcdlc:plan` stamps the architecture document's
  sha256 into `plan.md` as `<!-- arcdlc:source <path> sha256:… -->`; `arctool validate` warns when
  that document has changed and errors when it is gone; `/arcdlc:execute` stops on the warning rather
  than building a superseded design. `arctool stamp` refreshes it, which is how an engineer says they
  read the diff and the tasks still hold. An unstamped plan is not checked, so older plans keep
  working.
- **A thin design surfaces at plan time, and is reported as a design gap.** `/arcdlc:plan` Step 2.4
  names each blocker as a missing decision rather than a hard task, groups blockers by the decision
  they need (five tasks stuck on one question are one gap, not five), answers what the engineer can
  settle in an interview, and hands anything that would move several sections back to `/arcdlc:aic`.
  It never lowers the bar to get past the gate.
- **An architecture document is amended, never regenerated.** Re-running `/arcdlc:aic <slug>` is how a
  design matures: each round adds detail, tests an idea, or reverses a call. Sections the interview did
  not touch survive word for word. Inputs (goals, requirements, quality goals, constraints, context)
  are rewritten in place when gathering data changes them, and each decision resting on a changed
  input is asked about, never kept or changed silently. Decisions (hypotheses, arc42 §4 to §11, ADRs)
  are superseded: a reversed one moves in the same run, in both places it lives: the new ADR carries
  `- Supersedes:` and the old one's `- Status:` becomes `Superseded by ...`. Two ADRs that disagree
  become gaps `/arcdlc:examinate` files against code nobody wants changed. The sections of the
  template are the interview agenda, each closed as answered, deferred, or not applicable. See
  [ADR-0027](docs/adr/0027-the-template-is-the-interview-agenda.md).
- **Skills must stay install-agnostic.** Every `SKILL.md` must work as a Claude Code plugin command
  (`/arcdlc:<name>`) and as a flat skill (`arcdlc-<name>` on Codex, OpenCode, Cursor, Antigravity).
  Where a skill still reaches across to a sibling, keep both paths (`../grilling/...` and
  `../arcdlc-grilling/...`). Files a skill owns live in its own `references/`, which needs no dual path.
- **`arctool` is always optional in skills.** Every skill that uses it probes `command -v arctool` and
  describes the manual fallback. Never make a skill hard-depend on the CLI.
- **The bundle depends on nothing outside itself.** A skill delegates only to a sibling in this bundle.
  The interview is `skills/grilling` (`arcdlc-grilling`) and no other. Never reintroduce references to
  third-party skills on the user's machine; removing that dependency is why `skills/grilling` exists.
- **No skill hard-depends on another skill.** A delegating skill leads with the model-invocable sibling
  and ends in an inline fallback (read the sibling's `SKILL.md`, run its protocol yourself). It never
  stops to report a helper skill as missing. The behaviour is mandatory, the invocation is not.
- **The interview is one question at a time.** `skills/grilling` asks one question per turn, waits, then
  asks the next, each carrying a recommended answer. Never a numbered round. Six skills restate this in
  one shared paragraph for their inline fallback (`aic`, `policy`, `examinate`, `assist`, `plan-human`,
  `execute`); change it in all seven together.
- **The eight task keys are the contract; a custom key is carried, not judged.** A key-shaped line the
  format does not define is preserved with its indented body, surfaced under `extra` in `arctool`'s
  JSON, and ignored by `validate`. Absorption of a multi-line key stops at *any* key, defined or not:
  before that, an unknown key was swallowed into the section above it, corrupting that section as well
  as losing itself. A task reaches the executor whole.
- **The plan is a queue, not a graph.** An ordered list, run top to bottom, one task at a time. No
  `DEPENDS` key, no derived ordering, no parallel execution: a queue resumes from one pointer, a person
  reads the order instead of computing it, commits stay linear, and a concurrency failure cannot be
  expressed in a mechanical block. The cost is wall-clock time, and it is accepted. See
  [ADR-0024](docs/adr/0024-the-plan-is-a-queue-not-a-graph.md).
- **A task verifies itself; the queue is verified once, at the end.** A subagent reaches `DONE` only when
  every `Acceptance` criterion holds, and nothing re-checks a `DONE` task, so `DONE` is its assertion. The
  dispatcher does not inspect finished work: re-running the same criteria catches only a subagent that
  lied, and costs that output on every task. What no task can check is whether the tasks agree with each
  other, so the Verification phase does that, and it is the real gate. A contradiction found there is a
  plan defect that goes back to `/arcdlc:plan`, never a fix-up commit.
- **A spawned subagent has nobody to ask, and must never guess in place of asking.** It blocks the task
  with the question or the defect as the reason and returns it unanswered. **A blocked task is a correct
  outcome for a subagent, not a failure**, and its spawn prompt says so, because a cheaper model will
  otherwise guess to look helpful. `BLOCKED` is the only thing the dispatcher engages with: it holds the
  engineer, so it grills, then re-spawns with the answer or hands back to `/arcdlc:plan`.
- **The dispatcher stays thin, and is resumable.** In the loop it reads plan state, one line of notes per
  task, and a block reason. Not reports of successful work, not diffs, not test output. If it fills up
  anyway it stops at a task boundary and the engineer re-runs: the plan carries every status and the
  commits carry every change, so a 150-task queue is any number of sessions over an unchanged plan rather
  than one heroic run.
- **Unclear is a question, not a guess.** A skill that hits a contradiction, a missing decision, code
  that looks dead, or a risky change escalates to the grilled interview instead of choosing for the
  engineer. `aic`, `policy` and `plan` grill up front; `examinate` (Step 2.5), `assist` (Step 3),
  `plan-human` and `execute` grill mid-flight. Every answer is written where the next session reads it
  (an ADR, `CONTEXT.md`, the gap's `HOW`, the story, the commit body), never only in the chat. An answer
  that changes the plan stops the run: `/arcdlc:plan` owns the plan text.
- **Everything the skills write or say is written for humans.** Plain English, no AI filler, no long
  dashes in prose, acronyms expanded, domain terms kept and defined once. The rules cover chat replies
  as much as saved files. Plain words never mean less content, and neither does a short reply: every
  decision, constraint, trade-off and open question still has to be there, in chat as much as in a
  file. Brevity cuts filler, never facts.
  The `— ` separator in generated initiative-registry lines is exempt; that is `internal/registry`
  output, not prose. See the Conventions section below for where the rules live.
- **Judgement follows the four virtues.** Wisdom (ask, do not guess), courage (say the hard thing),
  justice (write the answer down, name the tier you used), temperance (touch only what you were pointed
  at). `## Judge by the four virtues` is verbatim in all eleven `SKILL.md` files and CI pins it byte
  identical, exactly like the writing-style block.
- **Every write is atomic and byte-preserving outside its own region.** `take`/`done`/`block`/`todo`
  rewrite only the one `- Status:` line via temp-file plus rename. `archive` writes the archive before
  compacting the plan. `order` re-parses its own output and writes nothing on a failed check. `scan`
  appends to `comments.md` and may grow one block, never more, and never removes or renumbers one.
  `scan` is also the only command that edits source files, and only with `--strip`, which deletes just
  the marker comment lines it registered. See [ADR-0013](docs/adr/0013-order-is-a-slot-permutation.md)
  and [ADR-0014](docs/adr/0014-exit-5-means-self-validation-failed.md).
- **The comment register has two owners.** `arctool scan` owns the task ID, the `- Marker:` lines, the
  `- WHERE:` list and a first draft of `- WHAT:`; `/arcdlc:assist` owns the title and every judgement
  key. Neither writes the other's lines, the sweep writes no verdict, and a block is never deleted: it
  is the only copy of a marker's text once the comment is stripped. See
  [ADR-0015](docs/adr/0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md).
- **A group tag makes one block.** `// ARCDLC:T1` in three files is one record and one task
  (`ARCDLC-CMT-T1`). Tags fold to upper case and must contain a letter. See
  [ADR-0016](docs/adr/0016-comment-markers-group-by-tag.md).
- **Only a single-line comment carries a marker.** `// ARCDLC ...` counts, `/* ARCDLC ... */` does not,
  in every language. A language whose only comment is a block comment carries no markers, and
  `openersByExt` says so with an empty entry. See
  [ADR-0018](docs/adr/0018-only-single-line-comments-carry-markers.md).
- **The default marker is `ARCDLC`, and removing a comment is asked for.** A sweep never touches the
  TODO and FIXME notes a repository already had; `--marker TODO` is how a team opts in. Cutting a
  comment needs `--strip`, which `/arcdlc:assist` passes only after the engineer answers yes (its
  Step 6). A marker counts only inside a comment, and the scanner re-reads a file line by line rather
  than losing everything below when it cannot tell. See
  [ADR-0017](docs/adr/0017-a-marker-is-arcdlc-in-a-comment-and-cutting-it-is-asked-for.md).
- **The comment styles the sweep reads are documented, and the document is checked.**
  [docs/comment-markers.md](docs/comment-markers.md) is the contract. Adding a language means two edits
  in one change set: the table in `internal/scan/scan.go` and the table on that page.
  `TestDocumentedCommentStylesMatchTheTable` fails when they disagree, in either direction.
- **Initiatives are folders; selection is mandatory and explicit.** Each lives in `docs/aics/<slug>/`
  (the architecture doc plus `plan.md`, `gap.md`, `comments.md`, `plan-archive.md`, `plan-human.md`).
  Skills take the slug as their first positional argument and `arctool` requires `--aic <slug>` or
  `--plan PATH`; missing means an error listing the initiatives, never a guess. The resolver is
  `resolvePlan` in `cmd/arctool`. A folder is created lazily by `atomicWrite`. Task IDs are unique per
  plan, not globally. ADRs and `CONTEXT.md` stay global. See
  [ADR-0001](docs/adr/0001-initiative-selection-is-always-explicit.md).
- **A workspace keeps its docs in a sibling repository named `docs`.** No environment variable, no
  marker file: the name is the whole configuration, and `/arcdlc:init` clones or creates the hub under
  it. Four symlinks at the workspace root, `AGENTS.md`, `CLAUDE.md`, `README.md`, and `CONTEXT.md`,
  point into the hub; `arctool sync` writes through them into the hub's own files, never replacing a
  link with a plain one. The hub lives on its default branch only, with a commit and a push after every
  write. A workspace task carries a `- Repo: <name>` key naming the repository its `WHERE` files live
  in, and `/arcdlc:execute` makes two commits per task: one in the hub, one in that repository. See
  [ADR-0026](docs/adr/0026-a-workspace-keeps-its-docs-in-a-sibling-repository-named-docs.md).
- **The initiative registry is generated.** `arctool sync` owns the `<!-- arcdlc:initiatives -->` blocks
  in `AGENTS.md` and `README.md`; `sync --check` fails on drift. Never hand-edit inside the markers.
  `/arcdlc:remove <slug>` deletes a folder after an explicit confirmation and re-syncs; `arctool` itself
  deletes nothing. See [ADR-0002](docs/adr/0002-registry-sync-via-marker-blocks.md) and
  [ADR-0003](docs/adr/0003-initiative-removal-by-skill-not-arctool.md).
- **The executor tier is asked for, never guessed.** A cheap tier is safe because the plan is
  mechanical, not the other way round. A tier is a model plus an effort level. Order, first
  answer wins: the task's optional `- Executor:` line (`- Executor: opus, high effort.`), then the
  `Executor tier:` pin in `CONTEXT.md`, then one question to the engineer before the first spawn, then
  in-session mode. A task pin binds that one task and is written by `/arcdlc:plan` only when the
  engineer asked for it; `HOW` never names a tier. `arctool` returns it as `executor` in `next --json`
  and `show --json`, and `validate --strict` reports `empty-executor` for a line with nothing on it.
  Any tier must clear the capability floor: shell, file edits, the project's test and lint commands, a
  commit. A block a weaker model cannot execute is a plan defect, never a reason to raise the tier or
  to add a pin. See [ADR-0019](docs/adr/0019-the-executor-tier-is-asked-for-not-guessed.md) and
  [ADR-0025](docs/adr/0025-a-task-pins-its-executor-tier-with-an-executor-key.md).
- **The reference library is dissolved.** Bundled documents live in the `references/` folder of the skill
  that consumes them, and a document earns its place only by changing what an agent produces. There is no
  `source-map` skill and no routing table. See
  [ADR-0020](docs/adr/0020-reference-library-dissolved-into-the-skills.md).
- **Version bumps:** the CLI version is `const version` in `cmd/arctool/main.go`; the bundle version lives
  in `.claude-plugin/plugin.json` and `.antigravity-plugin/plugin.json`, which move together. Bump
  whichever component you changed. Releases are cut by pushing a `v*` tag.

## Conventions

- One skill per directory under `skills/`, entry file always `SKILL.md`, YAML frontmatter with a
  `description` naming its triggers (the `/arcdlc:<name>` command and the `arcdlc-<name>` flat form).
- **The description is the routing contract, and it is the only always-on cost.** All eleven sit in
  context every session and nothing else decides whether a skill fires. Write the phrasings a person
  actually types, not a summary of what the skill does, and keep mechanics (arguments, output paths,
  format lists) in the body, which is read only after the skill fires. After changing one, re-run
  [docs/routing-checks.md](docs/routing-checks.md) by hand: no CI can test routing, because testing
  it means running a real agent.
- **Two blocks are verbatim in all eleven `SKILL.md` files**, after the intro and before the first step:
  `## Talk simple, write like a human` (how the agent talks and how it writes the files it produces)
  and `## Judge by the four virtues` (how it decides). Brevity applies to replies, never to files:
  neither block ever lets a rule, path, decision, or acceptance criterion be dropped. Keep both short
  but self-sufficient, because a skill is loaded while a reference is only read if the agent opens it.
  Only tables, examples, and the pre-save grep live in the long form,
  `skills/grilling/references/Writing Style.md`. CI checks both blocks byte identical across all eleven,
  so copy them when adding a skill and change all eleven plus the long form together.
- Reference documents live in the `references/` folder of the skill that consumes them. A document
  earns its place by changing what an agent produces: a template a skill copies, or a rule set whose
  identifiers a gap block cites. Public knowledge the model already has does not belong here.
  See [ADR-0020](docs/adr/0020-reference-library-dissolved-into-the-skills.md).
- **Every reference carries `**Reviewed**: YYYY-MM-DD`**, the date someone last read it end to end.
  Re-date it when you read it, not when you touch it. `internal/bundle` fails on a missing or
  unparseable marker and logs anything older than a year. The marker exists because a reference rots
  in silence: `Go Best Practice.md` described pre-1.18 Go with nothing on its face to say so, and it
  took reading all 428 lines to find out.
- Adding or renaming a sub-skill requires updating `SUBSKILLS` in `install.sh` and the skill-layout,
  style-block, virtues-block and installer-smoke checks in `.github/workflows/ci.yml`, in the same
  change set. Retiring one means adding its name to `LEGACY_SUBSKILLS` so an upgrade sweeps it.
- Exit codes of `arctool` are part of its interface (0 ok, 1 contract failure, 2 usage, 3 not
  found/empty, 4 I/O, 5 self-validation). Skills key off them; do not renumber.
- The Antigravity plugin manifest lives in `.antigravity-plugin/`, beside `.claude-plugin/`.
- `CLAUDE.md` is a symlink to this file; edit `AGENTS.md` only.
