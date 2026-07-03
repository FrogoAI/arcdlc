---
description: Build or update an initiative's architecture document under docs/aics/<slug>/ — AIC by default, or arc42, TOGAF, C4, ADR when given as argument (e.g. /arcdlc:aic arc42). Pass --aic <slug> to name the initiative folder, else it is derived from the interview. Always starts with a mandatory grill-with-docs interview before any document is written. Use when the user runs /arcdlc:aic, invokes arcdlc-aic, or asks to create an architecture document for an initiative.
argument-hint: "[aic|arc42|togaf|c4|adr] [--aic <slug>]"
---

# ArcDLC AIC (/arcdlc:aic)

Produce the architecture document that anchors the ArcDLC delivery pipeline:

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`, with `/arcdlc:examinate` feeding
compliance gaps into the plan at any point.

## Argument: document format

The optional argument selects the format. Resolve the template through the sibling `source-map` skill of this bundle
(from this file: `../source-map/source/` in the plugin layout, `../arcdlc-source-map/source/` in flat installs):

| Argument | Output file | Template in `../source-map/source/` |
| --- | --- | --- |
| *(none)* or `aic` | `docs/aics/<slug>/aic.md` | `AIC Template.md` |
| `arc42` | `docs/aics/<slug>/arc42.md` | `Arc42.md`, `arc42/arc42-template-EN.md` |
| `togaf` | `docs/aics/<slug>/togaf.md` | `TOGAF.md` |
| `c4` | `docs/aics/<slug>/c4.md` | `C4.md` |
| `adr` | `docs/adr/NNNN-<title>.md` (global) | `ADR.md` |
| anything else | `docs/aics/<slug>/<format>.md` | Look it up in the `source-map` table; if no matching source exists, tell the user and list available formats. |

## Argument: initiative folder (`--aic <slug>`)

Each initiative gets its own folder `docs/aics/<slug>/`, holding its architecture document, `plan.md`,
`gap.md`, and `plan-archive.md`. Determine the slug:

- If the user passes `--aic <slug>`, use it (a single kebab-case path segment: no `/` or `..`).
- Otherwise derive a kebab-case slug from the initiative title established in the grill interview, and
  **confirm it with the user** before creating the folder (e.g. "Checkout Redesign" → create
  `docs/aics/checkout-redesign/`?).

Write the architecture document — and later `plan.md` — inside that folder. ADRs stay **global** under
`docs/adr/`; `CONTEXT.md` stays at the repo root (both are cross-cutting, not per-initiative).

## Step 1 — Gather existing context (before asking anything)

Read what already exists so the interview builds on it instead of repeating it:

- `docs/aics/` (list existing initiative folders; and `docs/aics/<slug>/` for this one's AIC, arc42, plan, gap register)
- `CONTEXT.md` / `CONTEXT-MAP.md` (domain glossary)
- `docs/adr/` (prior decisions)
- `AGENTS.md`, `CLAUDE.md`, `README.md` of the target project
- The relevant template from the table above

## Step 2 — MANDATORY: grill the design

Never write the architecture document straight from the request. Run the `grill-with-docs` skill first — the whole
point of `/arcdlc:aic` is that the process is controlled: interview, then document.

- Invoke the `grill-with-docs` skill (which runs a `grilling` session using the `domain-modeling` skill).
- If `grill-with-docs` is not installed in this environment, run the same discipline inline:
  - Interview the user relentlessly about every aspect of the initiative — one question at a time, with your
    recommended answer for each.
  - If a question can be answered by exploring the codebase, explore instead of asking.
  - As decisions crystallise, write them down immediately: glossary terms into `CONTEXT.md`, architectural decisions
    into `docs/adr/NNNN-<slug>.md`.

Cover at minimum: problem/goal, scope boundaries (in/out), key quality attributes, ownership and context boundaries,
data model and storage, communication patterns, deployment model, tech stack deltas, risks, and open questions.

The interview ends only when the user confirms shared understanding or explicitly says to proceed.

## Step 3 — Write the document

- Fill the template section by section at the output path from the table.
- Every significant decision in the document must trace to an interview answer, an existing ADR, or code evidence.
  Do not imagine, invent, or silently assume architecture conclusions (this is the `source-map` rule).
- Anything still undecided goes into an explicit "Open questions" section — never into the body as if decided.
- Link the ADRs created during the interview from the document.

## Step 4 — Review and hand off

- Walk the user through the draft; iterate until approved.
- Tell the user the next step: `/arcdlc:plan` (with the same `--aic <slug>` if several initiatives exist) to
  decompose this document into the executable `docs/aics/<slug>/plan.md`.
