# Scaffold Templates

**Reviewed**: 2026-09-19

The files `/arcdlc:init` writes when they are missing: Step 3 for a single repository, Step 4 for a
workspace hub. A hub reuses the `docs/adr/README.md` and `docs/aics/README.md` templates below
unchanged, at the hub root instead of under a repository's `docs/`, and adds its own `docs/AGENTS.md`,
`docs/README.md`, `docs/CONTEXT.md`, and `docs/Makefile`. One fenced block per file. Placeholders:
`<Project>` (the repository or workspace name), `<description>` (the interview's one-paragraph
answer), `<tier>` (the interview's executor-tier answer), `<test-command>` and `<lint-command>` (a
`Makefile` target or a `package.json` script, or the literal text `none found` when neither exists;
single repository only, a hub has no build of its own). Every heading, marker, and table row shown
here is copied exactly; only the placeholders change, and a repository-table row is filled once per
child repository read at scaffold time. `CLAUDE.md` carries no template in either layout: it is
always `ln -s AGENTS.md CLAUDE.md` (`ln -s AGENTS.md docs/CLAUDE.md`, tracked, in a hub).

## AGENTS.md

```
# <Project> Agent Guide

<description>

## Build, test, verify

- Test: `<test-command>`
- Lint: `<lint-command>`

## Initiatives

<!-- arcdlc:initiatives:begin -->
_none_
<!-- arcdlc:initiatives:end -->
```

## README.md

```
# <Project>

<description>

## Initiatives

<!-- arcdlc:initiatives:begin -->
_none_
<!-- arcdlc:initiatives:end -->
```

## CONTEXT.md

```
# CONTEXT: <Project> vocabulary

Glossary for humans and agents working on this project.

## Execution

Executor tier: <tier>

## Terms
```

## docs/adr/README.md

```
# Architecture decision records

A record of every decision that was hard, surprising, or expensive to reverse. Name it
`NNNN-<title>.md`: a four-digit number, the next one free, then a short kebab-case title. Numbers are
never reused and never renumbered; a superseded record stays, its `Status` line saying what replaced
it.

| # | Decision | Status |
|---|---|---|
```

## docs/aics/README.md

```
# Initiatives

One folder per initiative, `docs/aics/<slug>/`: the architecture document, `plan.md`, and, once a
skill has written them, `gap.md`, `comments.md`, and `plan-archive.md`. The slug is a single
kebab-case path segment: no `/`, no `..`.

`/arcdlc:aic <slug>` starts one. `arctool sync` keeps the list in `AGENTS.md` and `README.md` in
step with what is here.
```

## docs/AGENTS.md

```
# <Project> Agent Guide

<description>

## Repos

| Repo | Remote | Default branch | Purpose |
|---|---|---|---|

## How to work

Run the agent, and every `/arcdlc:*` and `arctool` command, from the workspace root: one level above
this repository. The root holds four symlinks into this repository: `AGENTS.md`, `CLAUDE.md`,
`README.md`, and `CONTEXT.md`. Run `make -C docs init` once on a new machine to create them. Run
`make -C docs check` whenever the initiative registry looks wrong or a link seems missing.

## Git rules

This repository stays on its default branch. No feature branches, no pull requests. Commit and push
after every write. A push a remote rejects gets one retry, after `git pull --ff-only`. A second
rejection stops the run and reports the git error.

## Workspace tasks

A task in a workspace plan makes two commits: the code commit in the repository its `- Repo:` key
names, then the status commit here. `plan-format.md` documents the `Repo` key. No skill writes into a
product repository except that one code commit.

`arctool` below `0.20.0` drops the `Repo` key silently, so `/arcdlc:init` refuses to write a workspace
plan rule with an older binary. Keep every machine on `0.20.0` or later.

## ADR numbering

Each product repository keeps its own ADR sequence until it moves into `docs/adr/` here. Once two
sequences share this folder, a bare number is not unique: cite a record by its filename, never by
number alone.

## Initiatives

<!-- arcdlc:initiatives:begin -->
_none_
<!-- arcdlc:initiatives:end -->
```

## docs/README.md

```
# <Project>

This is the docs hub of the workspace above. Run `make -C docs init` and start your agent from the
workspace root.

<description>

## Initiatives

<!-- arcdlc:initiatives:begin -->
_none_
<!-- arcdlc:initiatives:end -->
```

## docs/CONTEXT.md

```
# CONTEXT: <Project> vocabulary

Glossary for humans and agents working on this workspace.

## Execution

Executor tier: <tier>

## Repos

| Repo | Purpose |
|---|---|

## Terms
```

## docs/Makefile

```
# Docs hub for a workspace. The workspace root is the parent of this directory. Run the agent, and
# every /arcdlc:* and arctool command, from there. AGENTS.md, CLAUDE.md, README.md, and CONTEXT.md at
# the root are symlinks into this repository.

HUB   := $(notdir $(CURDIR))
ROOT  := $(realpath $(CURDIR)/..)
LINKS := AGENTS.md CLAUDE.md README.md CONTEXT.md

.PHONY: init check

## init: create or repair the four root symlinks
init:
	@if [ -z "$(ROOT)" ]; then echo "error: cannot resolve the workspace root" >&2; exit 1; fi
	@for f in $(LINKS); do \
		target="$(HUB)/$$f"; path="$(ROOT)/$$f"; \
		if [ ! -e "$(CURDIR)/$$f" ]; then \
			echo "error: $(CURDIR)/$$f does not exist" >&2; exit 1; \
		fi; \
		if [ -L "$$path" ]; then \
			if [ "$$(readlink "$$path")" = "$$target" ]; then \
				echo "ok       $$f -> $$target"; continue; \
			fi; \
			echo "relink   $$f (was -> $$(readlink "$$path"))"; \
		elif [ -e "$$path" ]; then \
			echo "SKIP $$f exists as a regular file. Remove it by hand to replace it with the symlink." >&2; \
			continue; \
		else \
			echo "create   $$f -> $$target"; \
		fi; \
		ln -sfn "$$target" "$$path"; \
	done
	@echo
	@echo "Start agent sessions and run arcdlc from $(ROOT)"

## check: report whether the four root symlinks are correct
check:
	@rc=0; for f in $(LINKS); do \
		target="$(HUB)/$$f"; path="$(ROOT)/$$f"; \
		if [ -L "$$path" ] && [ "$$(readlink "$$path")" = "$$target" ]; then \
			echo "ok       $$f -> $$target"; \
		else \
			echo "MISSING  $$f (run: make init)"; rc=1; \
		fi; \
	done; exit $$rc
```
