# ArcDLC Init: one command that sets up a repository or a multi-repository workspace

> Add `/arcdlc:init`, which detects whether it stands in one repository or in a workspace of several, scaffolds the files the other skills need, and in a workspace sets up a sibling `docs` repository as the single home for every initiative, so cross-repository work runs with the same commands and the same paths as a single repository.

Decided on 2026-09-19 in a grilled interview. Terms used here (workspace, hub, root links, product
repository, Repo key, layout) are defined in [`CONTEXT.md`](../../../CONTEXT.md). The one hard
decision is [ADR-0026](../../adr/0026-a-workspace-keeps-its-docs-in-a-sibling-repository-named-docs.md).

## Goals

### Business Case

ArcDLC today assumes one repository: the code and its `docs/` folder in one checkout. Real delivery is
often several repositories that ship together. FrogoDB is four (`fdb-server`, `fdb-client`,
`fdb-tools`, `fdb-operator`) and most of its initiatives touch two or three of them. MetricAid is
thirty-nine.

Two layouts have been tried for that case, and the evidence is on disk.

The "owning repository" layout, used by FrogoDB: each initiative lives in the repository whose goal it
serves, and its tasks name files in the other repositories. It produced four defects. The workspace
root's `AGENTS.md` (5.7 KB) and `CONTEXT.md` (14 KB) are untracked files nobody can review or commit;
the context file itself says a correction there "can never be reviewed or committed". Ten initiatives
sit in one repository and one in another, with two registries. The root map omits `fdb-operator`.
Three ADR sequences (63, 4 and 6 records) grew in parallel and collide on `0001` to `0006`.

The "docs hub" layout, used by MetricAid since 2026-09-02: a sibling repository checked out as `docs/`
holds `aics/`, `adr/`, `policies/`, the agent guide and the glossary; the root carries symlinks into
it; the agent runs from the root. Thirty initiatives have run through it. Because the repository is
named `docs`, every `docs/aics/<slug>/` path the skills and `arctool` hard-code resolves unchanged,
and no skill needed a mode switch. Its cost is that arcdlc knows nothing about it: the rules live in
MetricAid's own hub guide, two of the root links are missing (`CLAUDE.md` and `CONTEXT.md`, so the
execute skill cannot see the executor pin from the root), and `arctool sync` run at the root replaces
a root symlink with a regular file, so the team hand-edits the registry and bans the command.

This initiative makes the hub layout an arcdlc feature with one entry point, fixes the two defects
found in the field, and gives an existing per-repository project a way in.

### Functional Overview

- **`/arcdlc:init`**, an eleventh skill. It detects the layout, runs a short grilled interview for
  what it cannot read, and scaffolds. In a workspace it clones or creates the hub as `docs/`, writes
  the hub's files, creates the four root links, and optionally migrates existing per-repository
  initiatives, ADRs and glossary terms into the hub.
- **Workspace rules in the sibling skills.** `/arcdlc:execute` gains the two-commit per task
  contract keyed on the task's `Repo` key, with the hub pushed after every write and a fast-forward
  retry on a rejected push. `/arcdlc:plan` writes the `Repo` key in a workspace. `/arcdlc:aic`,
  `/arcdlc:plan`, `/arcdlc:remove`, `/arcdlc:examinate` and `/arcdlc:assist` learn in one paragraph
  each that `docs/` may be a repository of its own which is committed and pushed after every write,
  and that `arctool scan` sweeps one product repository at a time (`--path <repo>`).
- **`arctool sync` becomes symlink-safe.** A target that is a symlink is resolved and the real file
  is written; each registry link is written relative to the real file's directory, so the hub's
  `AGENTS.md` reads `aics/<slug>/aic.md` and a single repository's still reads
  `docs/aics/<slug>/aic.md`.
- **Documentation.** The README's command table and a new "One repository or a workspace" section,
  `AGENTS.md`, `CONTEXT.md` (done in this round), `plan-format.md` documenting `Repo` as a
  recognised custom key, `install.sh`'s sub-skill list and the four CI checks that pin the skill
  layout and the two verbatim blocks across all skills.

Out of scope: any per-skill configuration of where `docs/` is (ADR-0026 locks the name), an
`arctool init` subcommand, parallel execution across repositories (ADR-0024 stands), a source layout
or build files for a new project, and moving a repository's reference documents (`protocol.md`,
`storage.md` and the like) into the hub.

### Quality Goals

| Priority | Quality goal | Measure |
| --- | --- | --- |
| 1 | **Compatibility** (arc42: compatibility). A single repository sees no change. | Every existing skill, `arctool` command and plan runs byte-identical on a single repository before and after; `arctool sync` output on a non-symlinked target is unchanged, verified by the existing tests. |
| 2 | **Operability** (arc42: operability). Setup is one command per layout and one per machine. | `/arcdlc:init` on an empty folder, a single repository and a workspace each end with the other skills runnable; a new machine is `make -C docs init`; `make -C docs check` reports a lost link. |
| 3 | **Integrity** (arc42: integrity). Two runs on one hub cannot silently lose a status write. | A rejected hub push is retried after a fast-forward; a fast-forward conflict on the task's own status line blocks that task with the reason and ends the run; no `--force` anywhere. |

### Organizational Constraints

- The skill bundle rules in `AGENTS.md` bind: both verbatim blocks in every `SKILL.md`, install-agnostic
  paths with the flat-install twin (`../grilling/` and `../arcdlc-grilling/`), `arctool` optional in
  every skill with a shell fallback, `SUBSKILLS` in `install.sh` and the four CI checks updated in
  the same change set, a `**Reviewed**` date on every reference document.
- The hub of a workspace lives on its default branch only: no branches, no pull requests, commit and
  push after every write. A product repository keeps its own guide, branch shape and review flow; a
  run commits there and never pushes (Q4, Q13).
- ADR filenames are kept on migration and cited by full filename; a bare number is not unique in a
  merged hub (Q7).
- The hub is the only home of an initiative. Nothing is mirrored into a product repository, and no
  product repository keeps a `docs/aics/` folder after migration (the root guides of both FrogoDB
  and MetricAid already forbid two copies).

### Technical Constraints

- `arctool` stays pure standard library and static. The `sync` change is `filepath.EvalSymlinks` on
  the target plus a link base computed from the real file's directory; the rename stays atomic.
- `aicsDir` stays the constant `docs/aics` (ADR-0001, ADR-0026). No environment variable, no marker
  file, no flag decides where initiatives live.
- The `Repo` key is a custom key, carried by arctool 0.18.0 and later as `extra.Repo`. The plan
  format's eight defined keys are unchanged; the parser, validator and mutator are untouched.
- Skills are plain Markdown with no agent-specific syntax. Layout detection and every scaffold step
  are shell the skill carries in its own text (`git rev-parse`, `ln -sfn`, `git mv`).
- Symlinks: `docs/CLAUDE.md -> AGENTS.md` is a tracked symlink inside the hub; the four root links
  are untracked. Both must work on Linux and macOS; Windows is not a target of the installer today
  and is not one here.

### Business Context

The system under construction is the `init` skill plus the workspace behaviour of the existing
skills and `arctool`. Its partners:

| Partner | What crosses the boundary |
| --- | --- |
| The engineer | Runs `/arcdlc:init` once per project, answers the interview, runs `make -C docs init` on each later machine. |
| The agent harness (Claude Code, Codex, OpenCode, Cursor, Antigravity) | Loads `CLAUDE.md` or `AGENTS.md` from the working directory; the root links make the hub's guide that file. |
| The hub's git remote | Receives a push after every hub write; rejects a non-fast-forward, which is the collision signal. |
| Product repositories | Receive commits on their own branches, never pushes; their `AGENTS.md` sets the branch shape; their `docs/aics` and `docs/adr` are moved out once. |
| `arctool` | Reads and writes `docs/aics/<slug>/plan.md` from the root; `sync` writes through the root links. |
| The sibling skills | `plan` writes the `Repo` key, `execute` reads it first, `remove` deletes in the hub and pushes, `examinate` and `assist` sweep one repository at a time. |

## Architectural Hypotheses

### H1. The hub directory is named `docs`, and that name is the entire configuration

- **Context.** Every skill and `arctool` address `docs/aics/<slug>/` relative to the working
  directory. Q2.
- **Decision.** The hub is a git repository checked out as `docs/` directly under the workspace root,
  no override. `init` clones it under that name whatever the remote is called. Recorded as
  [ADR-0026](../../adr/0026-a-workspace-keeps-its-docs-in-a-sibling-repository-named-docs.md).
- **Justification.** Zero configuration: nothing to read in ten skills and one binary, nothing to
  fail silently. One fact to teach. The single-repository workflow is untouched because both layouts
  look identical to every skill.
- **Trade-offs.** A product repository cannot be checked out as `docs`. The hub must sit one level
  under the root. Both accepted in the ADR.

### H2. `init` detects the layout and asks only what it cannot read

- **Context.** The skill is the one place that knows there are two layouts; the others stay
  layout-blind. Q1, Q8, Q9.
- **Decision.** Detection, in order: the working directory is a git repository, so single
  repository; it is not one and holds at least one git repository as a direct child, so workspace;
  it is empty or holds no repository at all, so one question, single repository recommended. The
  interview then covers only: the hub remote (existing or to be created), whether to migrate existing
  per-repository folders, the executor tier, and a one-paragraph project description written to the
  Writing Style guide. Facts such as each repository's remote URL and default branch (`origin/HEAD`)
  are read, not asked.
- **Justification.** "Facts are your job, decisions are theirs" is the grilling rule. A repository
  map read from git is correct on the day it is written; one typed by a person is the map that
  omitted `fdb-operator`.
- **Trade-offs.** A root that is itself a git repository with nested repositories (submodules, or a
  monorepo with vendored checkouts) is neither layout. `init` names the case and stops; see Open
  questions.

### H3. The single-repository scaffold is six files and nothing else

- **Context.** An empty folder must become a place the other skills can run. Q8.
- **Decision.** `AGENTS.md` (the project description from the interview plus the registry markers),
  `CLAUDE.md` as a symlink to it, `README.md` (markers), `CONTEXT.md` (empty glossary plus
  `Executor tier: <answer>`), `docs/adr/README.md` and `docs/aics/README.md` (one screen each,
  stating the folder's rules). Never a source layout, a Makefile or a build file. An existing file is
  never overwritten: `AGENTS.md` and `README.md` get the marker block appended when it is missing,
  every other existing file is skipped and named in the report.
- **Justification.** Each of the six is something `aic`, `execute`, `remove` or `sync` reads or
  writes. The code's layout is the project's decision.
- **Trade-offs.** A project with an `AGENTS.md` written in another style keeps it; `init` adds only
  the markers. Accepted: rewriting someone's guide is exactly the overreach temperance forbids.

### H4. The workspace scaffold is the single-repository one inside the hub, plus four root links and a Makefile

- **Context.** MetricAid's hub is the working model, with two links missing. Q6.
- **Decision.** Inside `docs/`: the same six files, with `aics/` and `adr/` at the hub root instead of
  under a `docs/` subfolder, and `docs/CLAUDE.md` a tracked symlink to `docs/AGENTS.md`. Plus a
  `Makefile` with `init` (create or relink the four root links, refuse to clobber a regular file) and
  `check` (report a missing or wrong link, exit non-zero; GNU make exits 2 on a failed recipe and reserves 1 for `-q`, so the contract is pass or fail, never a specific code). At the root: `AGENTS.md`, `CLAUDE.md`,
  `README.md`, `CONTEXT.md`, each a symlink to the file of the same name in `docs/`. The hub's
  `AGENTS.md` carries a repository table (name, remote, default branch, one line of purpose from the
  interview) and the workspace rules from H5 and H6. Nothing is written into a product repository.
- **Justification.** `CLAUDE.md` is what Claude Code auto-loads, `AGENTS.md` is what Codex and the
  others load; both must resolve to the hub's one guide. `CONTEXT.md` at the root is where the
  execute skill reads the executor pin, and today MetricAid's pin in `docs/CONTEXT.md` is invisible to
  it. The Makefile is how every later machine is set up, and it is already proven at MetricAid.
- **Trade-offs.** Four untracked links instead of two. Accepted: `make -C docs check` finds a lost
  one, and the links are the only thing the root holds.

### H5. A workspace task names its repository with a `Repo` custom key and root-relative paths

- **Context.** One task changes one repository; the executor must know which without reading prose.
  Q3.
- **Decision.** Every task in a workspace plan carries `- Repo: <name>` as a top-level custom key
  (`docs` for a hub-only task). `WHERE` paths are relative to the workspace root
  (`fdb-server/internal/migration/...`). `/arcdlc:plan` writes the key; `arctool next --json` returns
  it as `extra.Repo`; `/arcdlc:execute` reads it first and treats a workspace task without one as a
  plan defect to block on. `plan-format.md` documents `Repo` as a recognised custom key.
- **Justification.** It reaches the subagent mechanically today with no parser change. `WHERE` keeps
  its job of naming files, and the repository is readable from the path as well. A `WHERE` first line
  (MetricAid today) is not a field, and a `WHERE` holding only `Repo: —` fails `validate --strict`.
- **Trade-offs.** A custom key is not validated: a misspelled `- Repo:` value is a spawn that cannot
  find its repository, caught by the subagent's block, not by `arctool`. Accepted over a ninth defined
  key, which is a contract change across parser, validator, tests and every skill.

### H6. Two commits per task, hub pushed after every write, fast-forward on rejection, product repository never pushed

- **Context.** `plan.md` and the code are in different repositories, so the one-commit contract
  cannot hold. Q4, Q5.
- **Decision.** In a workspace, `/arcdlc:execute`'s per-task contract becomes: the code commit in the
  task's repository, on the branch shape that repository's own `AGENTS.md` sets, with the standard
  Conventional Commits message and `Refs: <TASK-ID>`; then the status commit in the hub,
  `chore(<slug>): mark <TASK-ID> DONE` (or `BLOCKED`), `Refs: <TASK-ID>`, on the hub's default branch,
  pushed at once. A rejected non-fast-forward push is retried after `git pull --ff-only`; a conflict
  on the fast-forward blocks the task with the reason and ends the run. The product repository is
  never pushed by a run. The same applies to every hub write by any skill: an AIC, an ADR, a
  `CONTEXT.md` term, a `gap.md` entry is committed and pushed when written.
- **Justification.** Pushing the hub at once is what makes a collision visible immediately, which
  was the requirement. Two runs on different tasks touch different status lines, so the fast-forward
  is clean and needs nobody; two runs on the same task conflict on its status line, which is the
  signal wanted. Product repositories have their own review flows, and two of FrogoDB's four already
  ask for unpushed work.
- **Trade-offs.** Two commits and a network round trip per task instead of one commit. An offline
  run stops at the first status commit's push; accepted, the alternative is a hub that lies. A
  status commit and its code commit can be separated by a crash between them; `Refs:` on both makes
  the pair recoverable, and a `TAKEN` task with no hub commit marks the crash as before.

### H7. `arctool sync` writes through a symlink and bases links on the real file

- **Context.** `registry.WriteFile` writes a temp file and renames it onto the target, which replaces
  a symlink with a regular file. MetricAid bans the command and hand-edits. Q10.
- **Decision.** `sync` resolves each target with `filepath.EvalSymlinks` when it is a symlink, writes
  the real file with the same temp-and-rename, and computes every registry link relative to the real
  file's directory. In a single repository the real file is at the root and links read
  `docs/aics/<slug>/aic.md` as today; in a hub they read `aics/<slug>/aic.md`. `--check` stays
  read-only. The CLI version is bumped.
- **Justification.** Every skill already calls `sync` with no arguments, the registry stays generated
  as ADR-0002 requires, and the change sits beside the rename. MetricAid's hand-edited blocks can be
  regenerated from the day the tool ships.
- **Trade-offs.** Registry links inside the hub differ from the links a single repository's file
  shows. Accepted: each file's links are correct from where that file is read, which is the only
  thing a link is for.

### H8. Migration moves initiatives, ADRs and glossary terms, rewrites plans mechanically, and commits once per repository

- **Context.** FrogoDB has ten initiative folders in two repositories, 73 ADRs in three sequences
  and a 14 KB untracked root glossary. Q7, Q11, Q12, Q13.
- **Decision.** For each product repository with `docs/aics/` or `docs/adr/`: `git mv` every
  initiative folder into `docs/aics/` of the hub and every ADR into `docs/adr/`, keeping filenames;
  merge the repository's `CONTEXT.md` terms and any untracked root `CONTEXT.md` into the hub's, one
  entry per term; remove the `<!-- arcdlc:initiatives -->` blocks from the repository's `AGENTS.md`
  and `README.md`; rewrite every moved `plan.md` by adding `- Repo: <repo>` to each task and
  prefixing each `WHERE` path that exists in that repository with `<repo>/`; run
  `arctool validate --strict` on every moved plan and stop on an error; commit the hub with a message
  naming each source repository and its commit, because file history does not follow a move across
  repositories; commit each product repository once, `docs: move initiatives and ADRs to the docs hub`,
  never pushed. Reference documents stay where the repository's `README.md` links them, unless the
  engineer names extra files during the interview. The hub's ADR index states that a bare number is
  not unique and that records are cited by filename.
- **Justification.** A task block must run with no context but itself, and after the move its old
  relative paths are wrong from the root, so an untouched plan is a queue of blocks that fail their
  own contract. Keeping ADR filenames keeps every citation in commits and code true. Reference
  documents describe one repository and are linked from its README, package docs and scripts; moving
  them breaks links for no gain. One commit per repository avoids both the two-copies state and
  unstaged deletions across four checkouts.
- **Trade-offs.** `git log` on a moved file starts at the migration commit; the hub commit names the
  source so the old history can be found. Text mentions of `docs/aics/...` in scripts and test strings
  (FrogoDB has them in a Makefile, two deploy scripts and one chaos test) become stale paths;
  migration lists them in its report and changes none, since they are not links and belong to the
  repository. Source stamps in moved plans keep their path (`docs/aics/<slug>/aic.md` is still right
  from the root) and their hash, so they stay valid.

### H9. `init` is a skill, not an `arctool` subcommand

- **Context.** The deterministic parts are a handful of shell lines; the interview is the substance.
  Q14.
- **Decision.** One new skill, `skills/init/SKILL.md`, carrying the shell for detection, scaffold,
  links and migration in its own steps, with the hub `Makefile` for every later machine. No new
  `arctool` command.
- **Justification.** `arctool` must stay optional in every skill, so the shell path exists anyway; a
  subcommand would be a second code path for the same six files. The one tool change this initiative
  needs is `sync` (H7), which is a fix to an existing command.
- **Trade-offs.** The scaffold's exact bytes are not tested by `go test`; they are covered by the
  installer smoke check in CI (the skill file lands) and by the skill's own acceptance criteria in the
  plan.

## Assessment

### Technical Challenges and Risks

- **R1. A stray write replaces a root link.** Any tool that renames onto the link (an editor's safe
  write, an older `arctool`) turns it into a regular file and the hub copy goes stale. Mitigation:
  H7 fixes the one arcdlc writer; `make -C docs check` detects it; the hub guide says to run `check`
  when the registry looks wrong.
- **R2. Ten skills, one new one, and CI pins two verbatim blocks across all of them.** Adding
  `init` without copying both blocks byte for byte fails CI; adding it without touching `SUBSKILLS`
  in `install.sh` means it never installs on the flat targets. Mitigation: the plan makes the skill
  file, the installer list and the four CI checks one task.
- **R3. Bundle version skew across machines sharing a hub.** The `Repo` key is carried only by
  arctool 0.18.0 and later; an older binary on another machine swallows it into the section above
  (the defect fixed in 0.18.0). Mitigation: the hub's `AGENTS.md` records the minimum `arctool`
  version, and `init` checks `arctool version` and refuses to write a workspace plan rule below it.
  Pinning with `install.sh --ref` is already the backlog's answer.
- **R4. A push retry loop.** A hub that rejects every push (a protected branch, a missing remote
  permission) must not retry forever. Mitigation: one fast-forward and one retry; a second rejection
  blocks the task and ends the run with the git error verbatim.
- **R5. Migration on a dirty repository.** `git mv` and a commit on a repository with unrelated
  staged or modified files sweeps them in (the hazard both FrogoDB and MetricAid memories record).
  Mitigation: migration refuses to start on any repository whose `git status --porcelain` is not
  empty, and names the repository.
- **R6. A moved plan whose `WHERE` paths cannot be found.** A path that exists in no repository gets
  no prefix and the task is left as it was, which `validate --strict` may pass and the executor may
  then fail. Mitigation: migration reports every unprefixed path per task, and `/arcdlc:plan` is the
  named next step for those plans.
- **R7. Two agent harnesses, two auto-loaded files.** Claude Code reads `CLAUDE.md`, the others read
  `AGENTS.md`. If either root link is missing, that harness runs with no guide and no idea the hub
  exists. Mitigation: four links, `make -C docs check`, and a first line in the hub's `README.md`
  telling a reader at the root what to run.

## Open questions

- **OQ1. A root that is itself a git repository containing other repositories.** Submodules, or a
  monorepo with vendored checkouts. `init` reports the case and stops. Whether it should treat the
  outer repository as a single repository or as a workspace has no answer yet; wait for a real case.
- **OQ2. Which product-repository branch shape the migration commit uses.** H8 commits once per
  repository on "its default branch or the branch its guide names". For FrogoDB two repositories sit
  on a feature branch today (`incident-0`). The migration step asks per repository when the checkout
  is not on the default branch.
- **OQ3. Whether MetricAid's hand-written `Repo:` first lines in `WHERE` are converted.** Thirty
  plans carry the MetricAid form. The migration step is for per-repository folders, not for an
  existing hub. A one-off rewrite there is MetricAid's call, not this initiative's.
- **OQ4. Whether `arctool sync` should also be run from inside the hub.** With H7, `sync` from the
  root writes through the links. `sync` run inside `docs/` finds no `docs/aics/` and writes `_none_`,
  as today. Guarding that (refuse to write an empty registry over a non-empty one without a flag) is
  a separate, small `arctool` change and is not decided here.
