---
name: arcdlc-remove
description: Remove a completed initiative — delete its docs/aics/<slug>/ folder and clean the initiative registry in AGENTS.md and README.md — after an explicit engineer confirmation. The initiative slug is the required first argument (e.g. /arcdlc:remove payments). Use when the user runs /arcdlc:remove, invokes arcdlc-remove, or asks to delete/retire an initiative and drop it from the registry.
argument-hint: "<slug>"
---

# ArcDLC Remove (/arcdlc:remove)

Retire a finished initiative so it stops costing agent focus and context. This deletes the whole
`docs/aics/<slug>/` folder and refreshes the initiative registry in `AGENTS.md` and `README.md`. Git
history is the archive — nothing is copied into a graveyard folder.

Removal is **destructive** and always requires an explicit human confirmation. `arctool` itself has no
delete command (it stays non-destructive); this skill does the deletion.

## Argument: initiative slug (required, first positional)

The initiative slug is the **first positional argument** and is **required**: `/arcdlc:remove <slug>`.

- If no slug is given, **stop and report the error** — highlight that the slug is missing and list the
  existing initiatives under `docs/aics/` so the user can pick one. Do not guess.
- Resolve the folder `docs/aics/<slug>/`. If it does not exist, stop and say so, listing what does.

## Step 1 — Show what will be removed

Before deleting anything, report the blast radius so the engineer decides with full information:

- The initiative **title** (the first `# ` heading of its architecture document).
- **Task counts by status.** Prefer `arctool` (probe once with `command -v arctool`):
  `arctool list --aic <slug>` prints the `TODO/TAKEN/DONE/BLOCKED` tallies. *Fallback (no `arctool`):
  read `docs/aics/<slug>/plan.md` and count `- Status:` lines by value.*
- The **file list** under `docs/aics/<slug>/` (e.g. `aic.md`, `plan.md`, `gap.md`, `plan-archive.md`).

If any task is **not `DONE`** (`TODO`, `TAKEN`, or `BLOCKED`), warn loudly and name the counts —
removal discards not-yet-completed work. Removal is still allowed (abandoned initiatives also need
cleanup), but the engineer must acknowledge it.

## Step 2 — Require explicit confirmation (always)

Ask the engineer to confirm the deletion, every time — there is no flag that skips this. State the
folder path and, when relevant, the not-`DONE` task count. Proceed only on an explicit yes; on anything
else, stop and change nothing.

## Step 3 — Delete and clean the registry

1. Delete the folder. If it is tracked by git, use `git rm -r docs/aics/<slug>/` (leaves the deletion
   staged for the user to commit). Otherwise remove it from the working tree (`rm -rf
   docs/aics/<slug>/`).
2. Refresh the registry so the removed initiative disappears from `AGENTS.md` and `README.md`. Prefer
   `arctool sync` (writes only the `<!-- arcdlc:initiatives -->` marker blocks). *Fallback (no
   `arctool`): delete the initiative's bullet from those marker blocks by hand; if no initiatives
   remain, leave the block reading `_none_`.*

Do not commit on the user's behalf unless they ask; leave the deletion and registry edit staged so they
can review.

## Step 4 — Report

Confirm what was removed: the slug, the folder, and that the registry was refreshed. If the initiative
had unfinished tasks, restate that they were discarded and note that `git` history still holds them.

The plan format and status contract this skill reads task counts from live in
`../plan/references/plan-format.md` (flat installs: `../arcdlc-plan/references/plan-format.md`).
