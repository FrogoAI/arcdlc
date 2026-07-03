---
description: Build or update an initiative's architecture document in docs/ — AIC by default, or arc42, TOGAF, C4, ADR when given as argument (e.g. /arcdlc:aic arc42). Always starts with a mandatory grill-with-docs interview before any document is written. Use when the user runs /arcdlc:aic, invokes arcdlc-aic, or asks to create an architecture document for an initiative.
argument-hint: "[aic|arc42|togaf|c4|adr]"
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
| *(none)* or `aic` | `docs/aics/aic.md` | `AIC Template.md` |
| `arc42` | `docs/aics/arc42.md` | `Arc42.md`, `arc42/arc42-template-EN.md` |
| `togaf` | `docs/aics/togaf.md` | `TOGAF.md` |
| `c4` | `docs/aics/c4.md` | `C4.md` |
| `adr` | `docs/adr/NNNN-<slug>.md` | `ADR.md` |
| anything else | `docs/aics/<format>.md` | Look it up in the `source-map` table; if no matching source exists, tell the user and list available formats. |

If the project keeps architecture docs elsewhere (e.g. `docs/epics/initiative-NN/`), follow the project convention
and keep the same file names.

## Step 1 — Gather existing context (before asking anything)

Read what already exists so the interview builds on it instead of repeating it:

- `docs/aics/` (any existing AIC, arc42, plan, gap register)
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
- Tell the user the next step: `/arcdlc:plan` to decompose this document into the executable `docs/aics/plan.md`.
