# 0001 — Initiative selection is always explicit

## Status

Accepted (2026-07-03). Supersedes the auto-detect resolver designed for arctool 0.6.0
(unreleased; still in the working tree at decision time).

## Context

arctool 0.6.0 introduced per-initiative folders (`docs/aics/<slug>/`) with an auto-detect
resolver: with neither `--plan` nor `--aic` given, exactly one initiative → selected implicitly;
several → exit 2; none → exit 3. On review the engineer overrode this: implicit selection hides
which plan a command is about to mutate, and behavior silently changes when a second initiative
appears — the same command that worked yesterday exits 2 today.

## Decision

Selection is mandatory and explicit everywhere:

- **arctool** requires `--aic SLUG` or `--plan PATH` on every plan-addressing command. Neither
  given → exit 2 (usage) with a message that highlights the missing selection and lists the
  initiatives found under `docs/aics/`. A named slug that does not exist → exit 3 (not found).
  The auto-detect branch is removed.
- **Pipeline skills** (`aic`, `plan`, `execute`, `examinate`, `archive`, plus `policy` for its
  name) take the slug as the **first positional argument** and stop with an error when it is
  missing, highlighting the omission and listing existing initiatives.
- The legacy flat `docs/aics/plan.md` is no longer implicitly addressable. Skills that find one
  tell the user to migrate it into `docs/aics/<slug>/`; `arctool --plan PATH` remains the escape
  hatch for any path.

## Consequences

Easier: command targets are predictable and self-documenting; behavior no longer shifts with the
number of initiatives; every error teaches the available slugs. Harder: one extra token to type
in single-initiative repos; the 0.6.0 auto-detect code is deleted before it ever shipped in a
release; the "Keep the legacy flat layout working" hard rule is narrowed to the `--plan` escape
hatch.
