---
name: arcdlc-aic
description: Build or update an initiative's architecture document under docs/aics/<slug>/. The initiative slug is the required first argument (e.g. /arcdlc:aic payments); an optional second argument picks the format (AIC by default, or arc42, tsc, TOGAF, C4, ADR), accepts a comma-separated list to produce several from one interview (e.g. /arcdlc:aic payments arc42,tsc), and takes a :html suffix to emit HTML instead of Markdown (e.g. /arcdlc:aic payments arc42:html). Always starts with a mandatory grilled interview before any document is written. Use when the user runs /arcdlc:aic, invokes arcdlc-aic, or asks to create an architecture document for an initiative.
argument-hint: "<slug> [aic|arc42|tsc|togaf|c4|adr, comma-separated, :html for HTML]"
---

# ArcDLC AIC (/arcdlc:aic)

Produce the architecture document that anchors the ArcDLC delivery pipeline:

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`, with `/arcdlc:examinate` feeding
compliance gaps into the plan at any point.

## Talk simple, write like a human

Plain B2 English, everywhere: short sentences, one idea each, active voice, a named actor. Cut
empty intensifiers: `honestly`, `genuinely`, `truly`, `clearly`, `obviously` add nothing.

- **Replies to the user:** bullets, not paragraphs. No filler, no praise, no restating the request.
  Say what you did, what you found, what comes next.
- **Files you write:** no AI filler ("Furthermore", "In conclusion", "It is important to note",
  "delve", "leverage", "robust", "seamless", "In today's fast-paced world"), no warm-up opener, no
  invented summary. Vary sentence length. Concrete names, numbers, and paths, never "significantly
  improves performance".
- **No long dashes.** Use a full stop, a comma, a colon, or brackets instead of `—` and `–`.
  Hyphens, flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.

Short talk, full content. Brevity is for your replies, never for the files: never drop a rule, path,
decision, trade-off, or acceptance criterion to save space. The full standard, with examples and a
pre-save check, is `source/Writing Style.md` in the bundle's `source-map` skill.

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
| `arc42` | `docs/aics/<slug>/arc42.md` | `Arc42 Guide.md` (start here — the full instruction) and `arc42-template-EN.md` (the upstream Markdown skeleton it tells you to copy). |
| `arc42:html` | `docs/aics/<slug>/arc42.html` | `Arc42 Guide.md` and `arc42-template.html` (the upstream Asciidoctor skeleton). Follow the guide's *Procedure — HTML output*. |
| `tsc` | `docs/aics/<slug>/tsc.md` | `Tech Stack Canvas.md` |
| `togaf` | `docs/aics/<slug>/togaf.md` | `TOGAF.md` |
| `c4` | `docs/aics/<slug>/c4.md` | `C4.md` |
| `adr` | `docs/adr/NNNN-<title>.md` (global) | `ADR.md` |
| anything else | `docs/aics/<slug>/<format>.md` | Look it up in the `source-map` table; if no matching source exists, tell the user and list available formats. |

### Output format: Markdown by default, HTML on request

Every format writes Markdown unless the engineer asks for HTML. Append `:html` to the format token —
`/arcdlc:aic payments arc42:html` — which changes only the extension of the output path and the
template used. Accept the natural phrasings too and normalise them to the same thing: "arc42 in
html", "arc42 html", "arc42, as html". If the engineer asks for HTML for a format that has no HTML
template in `../source-map/source/`, say so and offer the Markdown one rather than inventing markup.

`arctool sync` reads HTML architecture documents as well as Markdown ones (first `<h1>` for the
title, first `<p>` after it for the summary), so an HTML-only initiative still registers — but only
if the generated document replaces the template's logo `<h1>`, which the guide requires.

### Several formats at once

The argument accepts a **comma-separated list** — `/arcdlc:aic payments arc42,tsc`, and `:html`
composes with it (`arc42:html,tsc`) — and then every named format is produced from the **same single
interview**, in the order given:

- Run Step 2 (the grilled interview) **once**, covering what all the requested formats need. Never
  interview per format.
- Write one file per format at its own output path, each with its own `# <Title>` and `> ` summary
  line (Step 3). Same decisions, different lenses.
- Where two formats cover the same ground (e.g. arc42 §5/§7 and the TSC *Solution* / *Infrastructure*
  blocks), state it in full in one document and cross-link from the other with a relative link. Never
  fork the wording — two divergent copies is the failure this rule exists to prevent.
- The registry takes only one document per initiative. `arctool sync` picks it by precedence:
  `aic`, `arc42`, `togaf`, `c4`, `tsc` by format rank, `.md` before `.html` within a format, else the
  first `*.md` alphabetically, else the first `*.html`. Make sure the winning document's H1 and
  summary describe the whole initiative.
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

Never write the architecture document straight from the request — the whole point of `/arcdlc:aic` is
that the process is controlled: interview first, document second.

The interview runs on this bundle's own grilling skill, `arcdlc-grilling` (`/arcdlc:grilling`), a
sibling of this one. ArcDLC depends on no external grilling skill: do not look for one, and never
stop to report a skill as missing.

1. **Invoke `arcdlc-grilling`** — the normal path. It ships in every ArcDLC install and is
   model-invocable.
2. **Inline**, only if that skill is not invocable here: run the same protocol yourself,
   reading it from `../grilling/SKILL.md` (plugin layout) or `../arcdlc-grilling/SKILL.md` (flat
   installs). What is mandatory is the grilled interview, not the invocation.

Whichever path runs, all of this holds:

- **One question at a time.** Ask one question, wait for the answer, then ask the next. Never a
  numbered round, never "two quick ones", never "and also…" tacked onto a question.
- Every question carries your recommended answer and one line of why.
- Facts are yours to find: if the codebase, config, or an existing doc can answer it, look it up
  instead of asking.
- Decisions are written down the moment they settle — glossary terms into `CONTEXT.md`, hard and
  surprising trade-offs into `docs/adr/NNNN-<slug>.md`.

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

### Write it for a human reader (every format, Markdown and HTML)

The document is read by people. A new engineer must understand it without asking anyone. That is
the bar, and the `## Talk simple, write like a human` rules at the top of this skill apply to every section of
every format produced here, Markdown and HTML alike. Read
`../source-map/source/Writing Style.md` (flat installs: `../arcdlc-source-map/source/Writing Style.md`)
before the first section and run its pre-save check before you hand the draft over.

Four points matter most in an architecture document:

- Say the thing, then the reason, in that order: "We keep sessions in Postgres because we already
  run it and the volume is small."
- Break up stacked noun phrases ("install-agnostic cross-agent skill bundle distribution") into a
  sentence a person can say out loud.
- Expand every acronym on its first use in each document. Keep a domain term when it is this
  project's or this domain's own word, and define it once in plain words, here or in `CONTEXT.md`.
  Do not replace it with a vaguer everyday word.
- The `# <Title>` and the `> ` summary line follow the same rules. The summary is one plain
  sentence that tells a stranger what this initiative does, not a compressed spec.

Plain words, full content: simple English never means less detail. Keep every decision, constraint,
trade-off, risk, path, and acceptance criterion the template asks for. Just say it in words a person
reads once and gets.

## Step 4 — Register and hand off

- Walk the user through the draft; iterate until approved.
- Register the initiative so it is discoverable. Probe once with `command -v arctool`; if present, run
  `arctool sync` to refresh the `<!-- arcdlc:initiatives -->` blocks in `AGENTS.md` and `README.md`
  from this new folder. *Fallback (no `arctool`): add the initiative by hand to those marker blocks —
  `- [<Title>](docs/aics/<slug>/<doc>) — <summary>`.*
- Tell the user the next step: `/arcdlc:plan <slug>` to decompose this document into the executable
  `docs/aics/<slug>/plan.md`.
