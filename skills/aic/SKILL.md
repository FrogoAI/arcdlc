---
name: arcdlc-aic
description: Design an initiative's architecture and write it down. Produces an AIC by default, or arc42, Tech Stack Canvas, TOGAF, C4 or ADR, in Markdown or HTML, and always grills the design one question at a time before writing anything. Use when someone says design the architecture, write an architecture doc, spec out the system, draw up the design, or do a technical design for X, or runs /arcdlc:aic <slug> [format], or invokes arcdlc-aic.
argument-hint: "<slug> [aic|arc42|tsc|togaf|c4|adr, comma-separated, :html for HTML]"
---

# ArcDLC AIC (/arcdlc:aic)

Produce the architecture document that anchors the ArcDLC delivery pipeline. This is the entry point
of the application track; `/arcdlc:policy` is the entry point of the governance track.

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive` → `/arcdlc:remove`

Feeding the same plan queue at any point: `/arcdlc:examinate` (compliance gaps), `/arcdlc:assist`
(code comment markers). Beside the chain: `/arcdlc:plan-human` turns the queue into board stories.

Every one of these takes the initiative slug as its first argument, and every one writes into
`docs/aics/<slug>/`. When a user asks to design, plan, audit, implement, or govern something
end-to-end, route them through this pipeline instead of improvising.

## Talk simple, write like a human

Every rule here governs everything you emit, chat messages exactly as much as the files you write.
Plain English: short sentences, common words, one idea each, active voice, a named actor.

- **No AI filler, anywhere.** Not "Furthermore", "In conclusion", "It is important to note", "delve",
  "leverage", "robust", "seamless". No warm-up opener, no praise, no restating the request back, no
  invented summary. Concrete names, numbers, and paths, never "significantly improves performance".
- **Be direct, name the thing.** The fewest words that carry the fact. Cut empty intensifiers:
  `honestly`, `genuinely`, `truly`, `clearly`, `obviously`. Never a metaphor, a narrative line, or a
  question. Every title that names work is a verb plus its object: "Implement the alerter service".
- **No long dashes.** A full stop, a comma, a colon, or brackets instead of `—` and `–`. Hyphens,
  flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.
- **Shape.** In chat, bullets rather than paragraphs: what you did, what you found, what comes next.
  In files, vary sentence length so the prose does not read as a list.

Nothing is dropped to save space, in a reply as much as in a file. Never leave out a rule, path,
decision, trade-off, open question, or acceptance criterion because the answer is getting long: if it
bears on what the reader does next, it goes in. Short means no padding, never less content. Cut filler
words, repetition, and throat-clearing; never cut a fact. When completeness makes a reply long, let it
be long. The full standard, with examples and a pre-save check, is
`../grilling/references/Writing Style.md` (flat installs:
`../arcdlc-grilling/references/Writing Style.md`).

## Judge by the four virtues

- **Wisdom.** Unclear is a question, not a guess. A contradiction, a missing decision, code that looks
  dead, a change that is risky or hard to reverse: grill it, never pick for the engineer.
- **Courage.** Say the hard thing. A task that is not mechanical is a plan defect, not a reason to
  raise the tier. Never stop to report a helper skill as missing.
- **Justice.** Write every answer where the next session reads it, never only in the chat. Name the
  tier you actually used. Never report a same-tier spawn as a cheaper run.
- **Temperance.** Touch only what you were pointed at. Cutting a comment out of the code is asked for,
  not assumed. No invented summary, no scope you were not given.

## Argument: initiative slug (required, first positional)

The initiative slug is the **first positional argument** and is **required**:
`/arcdlc:aic <slug> [formats]` (e.g. `/arcdlc:aic payments`, `/arcdlc:aic payments arc42`, or
`/arcdlc:aic payments arc42,tsc`).

- If no slug is given, stop and report the error, listing the existing initiatives under `docs/aics/`
  so the user can pick one or name a new one. Never guess or derive a slug.
- The slug is a single kebab-case path segment (no `/` or `..`). Its folder is `docs/aics/<slug>/`,
  holding the architecture document, `plan.md`, `gap.md`, and `plan-archive.md`.

Write the architecture document — and later `plan.md` — inside that folder. ADRs stay **global** under
`docs/adr/`; `CONTEXT.md` stays at the repo root (both are cross-cutting, not per-initiative).

## Argument: document format(s)

The optional **second** positional argument selects the format (the first is the slug above). Resolve
the template from `references/` next to this file:

| Argument | Output file | Template in `references/` |
| --- | --- | --- |
| *(none)* or `aic` | `docs/aics/<slug>/aic.md` | `AIC Template.md` |
| `arc42` | `docs/aics/<slug>/arc42.md` | `Arc42 Guide.md` (start here — the full instruction) and `arc42-template-EN.md` (the upstream Markdown skeleton it tells you to copy). |
| `arc42:html` | `docs/aics/<slug>/arc42.html` | `Arc42 Guide.md` and `arc42-template.html` (the upstream Asciidoctor skeleton, copied with `cp`, never read into context). Follow the guide's *Procedure — HTML output*. |
| `tsc` | `docs/aics/<slug>/tsc.md` | `Tech Stack Canvas.md` |
| `togaf` | `docs/aics/<slug>/togaf.md` | `TOGAF.md` |
| `c4` | `docs/aics/<slug>/c4.md` | `C4.md` |
| `adr` | `docs/adr/NNNN-<title>.md` (global) | `ADR.md` |
| anything else | `docs/aics/<slug>/<format>.md` | Look for a matching template in `references/`; if there is none, tell the user and list the formats above. |

### Output format: Markdown by default, HTML on request

Every format writes Markdown unless the engineer asks for HTML. Append `:html` to the format token —
`/arcdlc:aic payments arc42:html` — which changes only the extension of the output path and the
template used. Accept the natural phrasings too and normalise them to the same thing: "arc42 in
html", "arc42 html", "arc42, as html". If the engineer asks for HTML for a format that has no HTML
template in `references/`, say so and offer the Markdown one rather than inventing markup.

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
- If any entry in the list is not a known format and has no template in `references/`, **write nothing**:
  report the unknown name and list the supported formats.

## In a workspace

The working directory is a workspace when it is not a git work tree and `git -C docs rev-parse
--show-toplevel` resolves to `<cwd>/docs` (see the `Workspace` and `Hub` terms in `CONTEXT.md`). In a
workspace, `docs/` is a repository of its own. Every write this skill makes there, an architecture
document, an ADR, a `CONTEXT.md` term, `gap.md`, `comments.md`, `plan.md`, is committed in the hub with
`git -C docs commit -- <hub-relative paths>` under the message `docs(<slug>): <what changed>` plus
`#AI-assisted`, and pushed to the hub's default branch at once. A rejected push is retried once after
`git -C docs pull --ff-only`. Nothing is written into a product repository by these skills. `arctool
sync` runs from the workspace root and writes through the root links.

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

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

Cover at minimum: problem/goal, scope boundaries (in/out), key quality attributes, ownership and context boundaries,
data model and storage, communication patterns, deployment model, tech stack deltas, risks, and open questions.

**Drive the interview strategically, in this order.** An architecture document that lists technology
without naming the problem it solves is a wish list, not a strategy.

1. **Diagnosis.** What is actually happening, and what evidence says so? Name the root cause, not the
   symptom. Refuse a goal stated only as a solution ("we need Kafka"): ask what breaks today.
2. **Guiding policy.** The approach that addresses the diagnosis, with its trade-offs stated out loud.
   A policy that costs nothing is not a policy. Test each one: can it navigate real trade-offs, will
   people actually be held to it, does it multiply effort rather than add to it?
3. **Coherent actions.** Concrete steps that follow from the policy, each one something
   `/arcdlc:plan` can turn into a task. An action that does not follow from the policy is scope creep.

When two answers compete, rank them company first, then team, then the individual preference. Record
the ranking you used, so the plan inherits the reason and not just the outcome.

The interview ends only when the user confirms shared understanding or explicitly says to proceed.

## Step 3 — Write the document

**If the document already exists, amend it. Never regenerate it.** Re-running `/arcdlc:aic <slug>` is
the normal way a design matures: each round adds detail, tests an idea, or reverses a call. So:

- Keep every section this interview did not touch, word for word. A detail added three rounds ago and
  not discussed today is still the decision. Rewriting the file from the template silently drops it,
  which is the one failure that makes iterating unsafe.
- Change only what the interview changed, and say which sections changed in your closing report.
- **A reversed decision is superseded, never deleted.** When this round rejects something an earlier
  round settled, both records have to move or the project now holds two answers:
  - Recorded as an ADR: write the new ADR, set `- Supersedes:` on it, and set the old one's
    `- Status:` to `Superseded by [ADR-NNNN](...)`. Never leave two ADRs that disagree:
    `/arcdlc:examinate` audits against `docs/adr/` and will file gaps for the decision you reversed.
  - Recorded only in the document body: replace the text, and note the reversal and its reason in the
    document, so the next reader does not re-litigate it.
  - Recorded in `CONTEXT.md`: a term whose meaning changed is redefined in place, not duplicated.
- **If a plan already exists** at `docs/aics/<slug>/plan.md`, this round may have invalidated tasks
  derived from what you changed. Say so in the report, name the affected sections, and tell the
  engineer to re-run `/arcdlc:plan <slug>`. Never edit `plan.md` yourself: `/arcdlc:plan` owns it.

For a new document:

- Fill the template section by section at the output path from the table, once per requested format.
- The document must open with a level-1 heading (`# <Title>`) and a one-line summary blockquote
  (`> …`) directly under it. This title and summary are a contract: `arctool sync` parses them into
  the initiative registry in `AGENTS.md`/`README.md`, so keep the summary to a single informative line.
- Every significant decision in the document must trace to an interview answer, an existing ADR, or code evidence.
  Do not imagine, invent, or silently assume architecture conclusions.
- Anything still undecided goes into an explicit "Open questions" section — never into the body as if decided.
- Link the ADRs created during the interview from the document.

### Write it for a human reader (every format, Markdown and HTML)

The bar is a new engineer understanding it without asking anyone. `## Talk simple, write like a human`
applies to every section of every format here. Four points matter most in an architecture document:

- Say the thing, then the reason, in that order: "We keep sessions in Postgres because we already
  run it and the volume is small."
- Break up stacked noun phrases ("install-agnostic cross-agent skill bundle distribution") into a
  sentence a person can say out loud.
- Expand every acronym on its first use. Keep a domain term when it is this project's own word and
  define it once, here or in `CONTEXT.md`. Never swap it for a vaguer everyday word.
- The `# <Title>` and `> ` summary follow the same rules. The summary is one plain sentence telling a
  stranger what this initiative does, not a compressed spec.

Run the pre-save check in `../grilling/references/Writing Style.md` (flat installs:
`../arcdlc-grilling/references/Writing Style.md`) before handing the draft over.

## Step 4 — Register and hand off

- Walk the user through the draft; iterate until approved.
- Register the initiative so it is discoverable. Probe once with `command -v arctool`; if present, run
  `arctool sync` to refresh the `<!-- arcdlc:initiatives -->` blocks in `AGENTS.md` and `README.md`
  from this new folder. *Fallback (no `arctool`): add the initiative by hand to those marker blocks —
  `- [<Title>](docs/aics/<slug>/<doc>) — <summary>`.*
- Tell the user the next step: `/arcdlc:plan <slug>` to decompose this document into the executable
  `docs/aics/<slug>/plan.md`.
