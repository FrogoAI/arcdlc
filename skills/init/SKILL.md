---
name: arcdlc-init
description: Set up this repository or a multi-repository workspace for ArcDLC, detecting which one it is, interviewing for the one paragraph and the executor tier only a person can give, and scaffolding the files every other skill reads. Use when someone says set up arcdlc here, bootstrap this repo for arcdlc, initialise the project, set up a multi-repo workspace, create the docs hub, or runs /arcdlc:init [--migrate], or invokes arcdlc-init.
argument-hint: "[--migrate]"
---

# ArcDLC Init (/arcdlc:init)

Set up a repository, or a workspace of several repositories, so the rest of the ArcDLC pipeline has a
place to run. Every other skill reads `docs/aics/<slug>/`, `CONTEXT.md`, and the registry markers in
`AGENTS.md` and `README.md`. `/arcdlc:init` is what creates them, before the first initiative exists.

`/arcdlc:init` → `/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive` → `/arcdlc:remove`

Run it once. It never rewrites a file that already exists, so running it again only adds what is
missing. Pass `--migrate` once a hub already exists, to also run Step 5 and move existing
per-repository initiatives, ADRs, and glossary terms into it.

## Talk simple, write like a human

Every rule here governs everything you emit, chat messages exactly as much as the files you write.
Plain English: short sentences, common words, one idea each, active voice, a named actor.

- **No AI filler, anywhere.** Not "Furthermore", "In conclusion", "It is important to note", "delve",
  "leverage", "robust", "seamless". No warm-up opener, no praise, no restating the request back, no
  invented summary. Concrete names, numbers, and paths, never "significantly improves performance".
- **Be direct, name the thing.** The fewest words that carry the fact. Cut empty intensifiers:
  `honestly`, `genuinely`, `truly`, `clearly`, `obviously`. Never a metaphor, a narrative line, or a
  question. Every title that names work is a verb plus its object: "Implement the alerter service".
- **No long dashes.** A full stop, a comma, a colon, or brackets instead of `—` and `–`. Hyphens,
  flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.
- **Shape.** In chat, bullets rather than paragraphs: what you did, what you found, what comes next.
  In files, vary sentence length so the prose does not read as a list.

Nothing is dropped to save space, in a reply as much as in a file. Never leave out a rule, path,
decision, trade-off, open question, or acceptance criterion because the answer is getting long: if it
bears on what the reader does next, it goes in. Short means no padding, never less content. Cut filler
words, repetition, and throat-clearing; never cut a fact. When completeness makes a reply long, let it
be long. The full standard, with examples and a pre-save check, is
`../grilling/references/Writing Style.md` (flat installs:
`../arcdlc-grilling/references/Writing Style.md`).

## Judge by the four virtues

- **Wisdom.** Unclear is a question, not a guess. A contradiction, a missing decision, code that looks
  dead, a change that is risky or hard to reverse: grill it, never pick for the engineer.
- **Courage.** Say the hard thing. A task that is not mechanical is a plan defect, not a reason to
  raise the tier. Never stop to report a helper skill as missing.
- **Justice.** Write every answer where the next session reads it, never only in the chat. Name the
  tier you actually used. Never report a same-tier spawn as a cheaper run.
- **Temperance.** Touch only what you were pointed at. Cutting a comment out of the code is asked for,
  not assumed. No invented summary, no scope you were not given.

## Step 1: detect the layout

The working directory decides which layout Step 3 or Step 4 scaffolds (the `Layout` term in
`CONTEXT.md`). Detection reads git; it asks only in the one case git cannot answer.

1. Run `git rev-parse --show-toplevel`.
2. **It fails.** The working directory is not inside a git repository. Check every direct child
   directory with `git -C <child> rev-parse --is-inside-work-tree`.
   - **At least one child answers `true`.** Workspace. Continue to Step 2.
   - **No child does.** Ask, in Step 2: is the working directory a single repository that has not
     had `git init` run in it yet, or the root of a workspace to create. Recommend single
     repository: it is the common case, and nothing about it has to be undone if the answer changes
     later.
3. **It succeeds and the output differs from the working directory.** The working directory sits
   inside a repository, but not at its root. Stop, and report: "run `/arcdlc:init` from the
   repository root `<path>`", naming the path `show-toplevel` printed. Do nothing else.
4. **It succeeds and the output equals the working directory.** The working directory is a git
   repository. Check every direct child directory with `git -C <child> rev-parse
   --is-inside-work-tree`.
   - **At least one child answers `true`.** The repository holds a nested repository of its own (a
     submodule, or a vendored checkout). This is neither layout: report the case by name, say that
     no layout is scaffolded for it, and stop. Open question OQ1 in `docs/aics/init/aic.md` covers
     it; nothing here invents one.
   - **No child does.** Single repository. Continue to Step 2.

## Step 2: interview

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

Ask only what git and the filesystem cannot tell you:

1. **The project's one-paragraph description.** What it is, in the reader's terms. Write it the way
   `../grilling/references/Writing Style.md` (flat installs: `../arcdlc-grilling/references/Writing
   Style.md`) requires: plain English, a named actor, no filler. It becomes the body of `AGENTS.md`
   and `README.md`.
2. **The executor tier.** The model `/arcdlc:execute` should run its task subagents at, and an
   optional effort level. Recommend `sonnet`, the cheapest tier that clears the capability floor
   `docs/adr/0019-the-executor-tier-is-asked-for-not-guessed.md` sets. Record the answer as
   `Executor tier: <model>` or `Executor tier: <model>, <effort>`, the same shape `CONTEXT.md`'s
   executor-tier pin already uses.
3. **In a workspace only, the hub remote.** An existing URL, or none. With none, Step 4 runs
   `git init` in `docs/` and tells the engineer a remote has to be added before the hub is ever
   pushed.
4. **In a workspace only, whether to migrate existing per-repository folders.** Yes runs Step 5 after
   the scaffold; no leaves every repository's `docs/aics/` and `docs/adr/` exactly where they are.

Read, never ask, the test and lint commands: a `Makefile` with `test` and `lint` targets, or the
matching `package.json` scripts. When neither exists, the templates in Step 3 carry the literal text
`none found` instead.

## Step 3: scaffold a single repository

Six files, each one another ArcDLC skill reads or writes; never a source layout, a Makefile, or a
build file (`docs/aics/init/aic.md` H3). Never overwrite a file that already exists: append to it
when the rule below says to, otherwise leave it alone and name it in the report.

The templates are in `references/Scaffold Templates.md`, one fenced block per file, with
`<Project>`, `<description>`, `<tier>`, `<test-command>`, and `<lint-command>` filled from Step 2 and
from the `Makefile`/`package.json` read.

1. `AGENTS.md`. Missing: write it from the `AGENTS.md` template. Present: read it. Already holds
   `<!-- arcdlc:initiatives:begin -->`: leave it untouched, and report it as skipped. Does not:
   append a `## Initiatives` section holding the two markers with `_none_` between them, keeping
   every existing byte above it unchanged.
2. `README.md`. The same rule as `AGENTS.md`, with the `README.md` template.
3. `CLAUDE.md`. Missing: `ln -s AGENTS.md CLAUDE.md`. Present as a regular file, or as a symlink to
   anything else: leave it, and report it as skipped. Codex, OpenCode, and Cursor read `AGENTS.md`
   directly; only Claude Code needs the symlink.
4. `CONTEXT.md`. Missing: write it from the `CONTEXT.md` template, `## Execution` holding
   `Executor tier: <answer>` from Step 2, `## Terms` empty. Present: leave it, and report it as
   skipped.
5. `docs/adr/README.md`. Missing: write it from the template. Present: leave it, and report it as
   skipped.
6. `docs/aics/README.md`. Missing: write it from the template. Present: leave it, and report it as
   skipped.

Then, when `command -v arctool` succeeds, run `arctool sync` so the registry reflects whatever is
already under `docs/aics/` (empty on a fresh scaffold, which writes `_none_`, the same text a newly
written or newly appended file already holds). *Fallback, no `arctool`: leave `_none_` between the
markers.*

## Step 4: scaffold a workspace

Runs only when Step 1 detected a workspace.

1. **Check out the hub.** `docs/` is the hub (the `Hub` term in `CONTEXT.md`).
   - `docs/` does not exist and Step 2 gave a remote: `git clone <url> docs`.
   - `docs/` does not exist and Step 2 gave none: `mkdir docs && git -C docs init`.
   - `docs/` exists and is a git work tree (`git -C docs rev-parse --is-inside-work-tree` succeeds):
     use it as is.
   - `docs/` exists and is not a git work tree: stop and report it. A plain folder named `docs` at a
     workspace root is the one collision `docs/adr/0026-a-workspace-keeps-its-docs-in-a-sibling-repository-named-docs.md`
     does not allow.
2. **Check `arctool`.** When `command -v arctool` succeeds, run `arctool version`. Below `0.20.0`,
   stop and report that the workspace rules need the `sync` fix shipped in `0.20.0`, naming
   `install.sh` as how to get it.
3. **Write the hub files**, from `references/Scaffold Templates.md`, filled the way Step 3 fills its
   templates, from Step 2's answers and the repository table read in the next step. Never overwrite a
   file that already exists.
   1. `docs/AGENTS.md`. Missing: write it from the `docs/AGENTS.md` template, including the
      repository table. Present: read it. Already holds `<!-- arcdlc:initiatives:begin -->`: leave it
      untouched, and report it as skipped. Does not: append a `## Initiatives` section holding the
      two markers with `_none_` between them, keeping every existing byte above it unchanged.
   2. `docs/README.md`. The same rule as `docs/AGENTS.md`, with the `docs/README.md` template.
   3. `docs/CLAUDE.md`. Missing: `ln -s AGENTS.md docs/CLAUDE.md`, then `git -C docs add CLAUDE.md`
      so the symlink is tracked. Present as a regular file, or as a symlink to anything else: leave
      it, and report it as skipped.
   4. `docs/CONTEXT.md`. Missing: write it from the `docs/CONTEXT.md` template, `## Execution`
      holding `Executor tier: <answer>` from Step 2, `## Repos` holding the same rows as the
      `docs/AGENTS.md` table minus the remote and default-branch columns, `## Terms` empty. Present:
      leave it, and report it as skipped.
   5. `docs/adr/README.md`. Missing: write it from the `docs/adr/README.md` template (Step 3's
      template, unchanged, placed at the hub root instead of under a repository's `docs/`). Present:
      leave it, and report it as skipped.
   6. `docs/aics/README.md`. Missing: write it from the `docs/aics/README.md` template (Step 3's
      template, unchanged, same placement rule). Present: leave it, and report it as skipped.
   7. `docs/Makefile`. Missing: write it from the `docs/Makefile` template. Present: leave it, and
      report it as skipped.
4. **Read the repository table.** For every direct child of the workspace root that is a git
   repository, other than `docs/` itself: `git -C <repo> remote get-url origin` for the remote; the
   short name of `git -C <repo> symbolic-ref refs/remotes/origin/HEAD` for the default branch when it
   resolves, otherwise the output of `git -C <repo> branch --show-current` marked `(no origin/HEAD)`;
   and a one-line purpose, asked in the interview, one question per turn. These rows fill the
   `docs/AGENTS.md` table (name, remote, default branch, purpose) and the `docs/CONTEXT.md` `## Repos`
   table (name, purpose).
5. **Create the root links and commit.** Run `make -C docs init` from the workspace root. When
   `command -v arctool` succeeds, run `arctool sync` from the root so the registry reflects whatever
   the hub already holds (fallback, no `arctool`: leave `_none_` between the markers, which a freshly
   written file already holds). Commit the hub: `git -C docs commit -m "docs: set up the docs hub" -m
   "#AI-assisted"`. When a remote exists, `git -C docs push origin <default>`; on a non-fast-forward
   rejection, retry once after `git -C docs pull --ff-only origin <default>`; a second rejection, or a
   fast-forward that itself fails, stops the run and reports the git error verbatim.

Report additions, for a workspace: the hub path, its remote or "no remote yet", the four root links,
and the files left alone (as Step 3 reports them). When Step 1 found `docs/aics/` in a child
repository and the interview's migration answer was no, offer `/arcdlc:init --migrate` as the next
step instead of naming an initiative to start.

## Step 5: migrate

Written by INIT-8.

## Report

State the layout detected, every file created, every file left alone and why it was left, and the
next step: `/arcdlc:aic <slug>` to start the first initiative.
