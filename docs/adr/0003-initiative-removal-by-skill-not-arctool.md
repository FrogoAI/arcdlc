# 0003 — Initiative removal is a skill-side operation; arctool stays non-destructive

## Status

Accepted (2026-07-03).

## Context

Completed initiatives keep their folders forever; the engineer's position is that in-tree history
of a finished initiative loses focus and context. Removal must exist, must always re-ask the
engineer, and must clean the registry in `AGENTS.md`/`README.md`. arctool's write invariants are
single-line byte-preserving status mutations and (per ADR-0002) markers-only block rewrites —
recursive tree deletion does not fit that contract.

## Decision

Add a new skill `/arcdlc:remove <slug>` (flat form: `arcdlc-remove`); arctool gets **no** delete
command. The skill:

1. Resolves `docs/aics/<slug>/` and reports what removal means: initiative title, task counts by
   status from `plan.md`, and the file list.
2. Warns loudly when any task is not `DONE` (removal is still allowed — abandoned initiatives
   also need cleanup).
3. **Always** requires explicit engineer confirmation, every time — no flag skips it.
4. Deletes the folder (`git rm -r` when tracked, `rm -rf` otherwise) and runs `arctool sync` to
   clean the registry (manual fallback: hand-edit the marker blocks).

Git history is the recovery path; no graveyard copy is kept in the working tree.

## Consequences

Easier: the repo keeps only live initiatives; the registry stays truthful; the deterministic tool
keeps a small, auditable blast radius. Harder: removal is not available as a one-shot CLI command
— intentional, it must pass a human; unarchived task history leaves the working tree and is
recoverable only through git.
