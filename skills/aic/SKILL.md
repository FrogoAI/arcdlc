---
name: arcdlc-aic
description: Build or update an initiative's architecture document under docs/aics/<slug>/. The initiative slug is the required first argument (e.g. /arcdlc:aic payments); an optional second argument picks the format (AIC by default, or arc42, tsc, TOGAF, C4, ADR) and accepts a comma-separated list to produce several from one interview (e.g. /arcdlc:aic payments arc42,tsc). Always starts with a mandatory grill-with-docs interview before any document is written. Use when the user runs /arcdlc:aic, invokes arcdlc-aic, or asks to create an architecture document for an initiative.
argument-hint: "<slug> [aic|arc42|tsc|togaf|c4|adr, comma-separated]"
---

# ArcDLC AIC (/arcdlc:aic)

Produce the architecture document that anchors the ArcDLC delivery pipeline:

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`, with `/arcdlc:examinate` feeding
compliance gaps into the plan at any point.

## Talk simple and short

While this skill runs, keep replies to the user simple and brief:

- Plain English, common words (B2). Short sentences, one idea each.
- Bullets over paragraphs. No filler, no praise, no repeating the request.
- Say only what matters: what you did, what you found, what comes next.
- Keep a technical term only when it is this project's own term.

Brief talk, full content: never drop a rule, path, acceptance criterion, or decision
to save words. Files this skill writes keep their required detail.

## Argument: initiative slug (required, first positional)

The initiative slug is the **first positional argument** and is **required**:
`/arcdlc:aic <slug> [formats]` (e.g. `/arcdlc:aic payments`, `/arcdlc:aic payments arc42`, or
`/arcdlc:aic payments arc42,tsc`).

- If no slug is given, **stop and report the error** — highlight that the slug is missing and list the
  existing initiatives under `docs/aics/` so the user can pick a name or reuse one. Do not guess or
  silently derive a slug.
- The slug is a single kebab-case path segment (no `/` or `..`). Its folder is `docs/aics/<slug>/`,
  holding the architecture document, `plan.md`, `gap.md`, and `plan-archive.md`.

Write the architecture document — and later `plan.md` — inside that folder. ADRs stay **global** under
`docs/adr/`; `CONTEXT.md` stays at the repo root (both are cross-cutting, not per-initiative).

## Argument: document format(s)

The optional **second** positional argument selects the format (the first is the slug above). Resolve
the template through the sibling `source-map` skill of this bundle (from this file:
`../source-map/source/` in the plugin layout, `../arcdlc-source-map/source/` in flat installs):

| Argument | Output file | Source in `../source-map/source/` |
| --- | --- | --- |
| *(none)* or `aic` | `docs/aics/<slug>/aic.md` | `AIC Template.md` |
| `arc42` | `docs/aics/<slug>/arc42.md` | `Arc42 Guide.md` (start here — the full instruction) and `arc42-template-EN.md` (the upstream skeleton it tells you to copy). Two files, no third source. |
| `tsc` | `docs/aics/<slug>/tsc.md` | `Tech Stack Canvas.md` |
| `togaf` | `docs/aics/<slug>/togaf.md` | `TOGAF.md` |
| `c4` | `docs/aics/<slug>/c4.md` | `C4.md` |
| `adr` | `docs/adr/NNNN-<title>.md` (global) | `ADR.md` |
| anything else | `docs/aics/<slug>/<format>.md` | Look it up in the `source-map` table; if no matching source exists, tell the user and list available formats. |

### Several formats at once

The argument accepts a **comma-separated list** — `/arcdlc:aic payments arc42,tsc` — and then every
named format is produced from the **same single interview**, in the order given:

- Run Step 2 (the grilled interview) **once**, covering what all the requested formats need. Never
  interview per format.
- Write one file per format at its own output path, each with its own `# <Title>` and `> ` summary
  line (Step 3). Same decisions, different lenses.
- Where two formats cover the same ground (e.g. arc42 §5/§7 and the TSC *Solution* / *Infrastructure*
  blocks), state it in full in one document and cross-link from the other with a relative link. Never
  fork the wording — two divergent copies is the failure this rule exists to prevent.
- The registry takes only one document per initiative. `arctool sync` picks it by precedence:
  `aic.md`, `arc42.md`, `togaf.md`, `c4.md`, `tsc.md`, else the first `*.md` alphabetically. Make sure
  the winning document's H1 and summary describe the whole initiative.
- If any entry in the list is not a known format and has no `source-map` match, **write nothing**:
  report the unknown name and list the supported formats.

## Step 1 — Gather existing context (before asking anything)

Read what already exists so the interview builds on it instead of repeating it:

- `docs/aics/` (list existing initiative folders; and `docs/aics/<slug>/` for this one's AIC, arc42, plan, gap register)
- `CONTEXT.md` / `CONTEXT-MAP.md` (domain glossary)
- `docs/adr/` (prior decisions)
- `AGENTS.md`, `CLAUDE.md`, `README.md` of the target project
- The relevant template(s) from the table above — for `arc42`, both files, starting with `Arc42 Guide.md`

## Step 2 — MANDATORY: grill the design

Never write the architecture document straight from the request. Run the `grill-with-docs` skill first — the whole
point of `/arcdlc:aic` is that the process is controlled: interview, then document.

- Invoke the `grill-with-docs` skill (which runs a `grilling` session using the `domain-modeling` skill).
- If `grill-with-docs` cannot be used here — not installed, or installed but not model-invocable (e.g. marked
  `disable-model-invocation`, in which case ask the user to run it) — run the same discipline inline:
  - Interview the user relentlessly about every aspect of the initiative — one question at a time, with your
    recommended answer for each.
  - If a question can be answered by exploring the codebase, explore instead of asking.
  - As decisions crystallise, write them down immediately: glossary terms into `CONTEXT.md`, architectural decisions
    into `docs/adr/NNNN-<slug>.md`.

Cover at minimum: problem/goal, scope boundaries (in/out), key quality attributes, ownership and context boundaries,
data model and storage, communication patterns, deployment model, tech stack deltas, risks, and open questions.

The interview ends only when the user confirms shared understanding or explicitly says to proceed.

## Step 3 — Write the document

- Fill the template section by section at the output path from the table, once per requested format.
- The document must open with a level-1 heading (`# <Title>`) and a one-line summary blockquote
  (`> …`) directly under it. This title and summary are a contract: `arctool sync` parses them into
  the initiative registry in `AGENTS.md`/`README.md`, so keep the summary to a single informative line.
- Every significant decision in the document must trace to an interview answer, an existing ADR, or code evidence.
  Do not imagine, invent, or silently assume architecture conclusions (this is the `source-map` rule).
- Anything still undecided goes into an explicit "Open questions" section — never into the body as if decided.
- Link the ADRs created during the interview from the document.

## Step 4 — Register and hand off

- Walk the user through the draft; iterate until approved.
- Register the initiative so it is discoverable. Probe once with `command -v arctool`; if present, run
  `arctool sync` to refresh the `<!-- arcdlc:initiatives -->` blocks in `AGENTS.md` and `README.md`
  from this new folder. *Fallback (no `arctool`): add the initiative by hand to those marker blocks —
  `- [<Title>](docs/aics/<slug>/<doc>) — <summary>`.*
- Tell the user the next step: `/arcdlc:plan <slug>` to decompose this document into the executable
  `docs/aics/<slug>/plan.md`.
