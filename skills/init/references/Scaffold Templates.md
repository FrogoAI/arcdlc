# Scaffold Templates

**Reviewed**: 2026-09-19

The five files `/arcdlc:init` writes when they are missing: Step 3 for a single repository, Step 4
for a workspace hub. One fenced block per file. Placeholders: `<Project>` (the repository or
workspace name), `<description>` (the interview's one-paragraph answer), `<tier>` (the interview's
executor-tier answer), `<test-command>` and `<lint-command>` (a `Makefile` target or a `package.json`
script, or the literal text `none found` when neither exists). Every heading, marker, and table row
shown here is copied exactly; only the placeholders change. `CLAUDE.md` carries no template: it is
always `ln -s AGENTS.md CLAUDE.md`.

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
