# ADR-0028 — A finished initiative is closed in place with a `CLOSED.md` note

- Status: Accepted
- Date: 2026-09-25
- Initiative: [aic-tracking](../aics/aic-tracking/aic.md)
- Supersedes: [ADR-0003](0003-initiative-removal-by-skill-not-arctool.md)

## Context

An initiative had two end states: open forever in `docs/aics/<slug>/`, or deleted by
`/arcdlc:remove`. ADR-0003 chose deletion as the way out and git history as the archive, on the
reasoning that a finished initiative kept in the tree loses focus and context.

That cost more than it saved. Five initiatives were deleted on 2026-09-16, and seven ADRs now read
"folder retired, see git history": the designs behind them can only be read with `git log`, so no
later initiative used them. Meanwhile `init`, fully `DONE` and archived, still shows as open, and
nothing stops new tasks from being added to it.

Moving a finished folder to an archive (as OpenSpec does with `changes/archive/`) was designed and
dropped: every link into the folder breaks, and fixing them needs a link rewriter.

## Decision

1. **Closing is the way out, in place.** `/arcdlc:close <slug>` writes `docs/aics/<slug>/CLOSED.md`
   (date, outcome, reason, unfinished tasks, what was left open) and deletes `plan.md`,
   `plan-archive.md`, `plan-human.md`, `gap.md` and `comments.md`. The design documents never move,
   so every link to them keeps working.
2. **A closed initiative is final.** `aic`, `plan`, `examinate` and `assist` stop on a folder
   holding `CLOSED.md`, and no skill deletes it. Follow-up work is a new initiative that links the
   closed design. A later design that reverses it appends a `## Superseded` line to its
   `CLOSED.md`; the design documents never change.
3. **`/arcdlc:remove` deletes a design folder for good,** only after an explicit confirmation that
   states the design will survive only in git history and that references to it may be left
   pointing at nothing.
4. **Close and remove do not handle links.** A link to a deleted file may break; git history holds
   the target.
5. **The registry lists open initiatives only.** `arctool sync` skips a folder holding `CLOSED.md`
   and ends the block with one line counting the closed ones, so `AGENTS.md` stays bounded.
6. **The skills do the work.** Beyond that `sync` rule, `arctool` gains only `arctool status`, a
   read-only listing of every initiative with its phase. It still deletes nothing and runs no git,
   as ADR-0003 decided.

## Consequences

Easier: a finished design stays where every link points, with a note saying it is finished and how.
Done and paused are told apart by one file. Status is read from the files, with no index.

Harder: `docs/aics/` holds open and closed designs side by side, and stays the actual design only
because a design round that reverses a closed design appends a "Superseded" line to its
`CLOSED.md`. A deleted plan or
register can leave a dead link behind.

The full design is in [the aic-tracking architecture document](../aics/aic-tracking/aic.md).
