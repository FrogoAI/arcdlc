# 0002 — `arctool sync` maintains the initiative registry via HTML-comment marker blocks

## Status

Accepted (2026-07-03).

## Context

Initiative folders are undiscoverable: nothing in the target project lists what lives under
`docs/aics/`. The requirement is that `AGENTS.md` and `README.md` MUST track every initiative
folder with its title and a short description. Manual registration (how policies are registered
today) drifts the first time a folder is created or deleted outside the skill flow.

## Decision

Add `arctool sync [--check]`:

- Scans `docs/aics/*/` and derives, per initiative:
  - **Title** — the first `# ` H1 of the architecture document. Document precedence within the
    folder: `aic.md`, `arc42.md`, `togaf.md`, `c4.md`, else the first `*.md` alphabetically.
  - **Summary** — a one-line `> ` blockquote directly under that H1. The `aic` skill is required
    to write this line for every architecture document it produces (format contract). Fallbacks:
    first non-empty paragraph truncated to ~120 chars; no document at all → the slug, flagged
    `(no architecture doc)`.
- Rewrites **only** the region between `<!-- arcdlc:initiatives:begin -->` and
  `<!-- arcdlc:initiatives:end -->` in `AGENTS.md` and `README.md` at the project root, via
  temp-file + rename — every byte outside the markers is preserved (same invariant style as
  status flips). Markers absent → append an `## Initiatives` section containing them. File
  absent → create a minimal stub (H1 + the section) and say so.
- Block content: one bullet per initiative, alphabetical by slug:
  `- [<title>](docs/aics/<slug>/<doc>) — <summary>`; with no initiatives the block holds `_none_`.
- `--check` writes nothing and exits non-zero when the blocks are stale (CI-friendly).
- Scope: **initiatives only**. Policies keep their existing skill-driven three-place
  registration; unifying them into sync is deferred.
- The `aic` and `remove` skills run `arctool sync` after creating or deleting a folder (manual
  fallback: edit the marker blocks by hand).

## Consequences

Easier: the registry cannot drift silently (CI can enforce it); removal cleanup comes for free;
agents and humans see every initiative from the repo front door. Harder: arctool now edits two
user-owned root files — mitigated by the markers-only rewrite; the H1 + blockquote pair becomes
contract surface that needs tests; hand edits placed inside the markers are overwritten by design.
