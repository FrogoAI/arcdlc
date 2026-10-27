---
name: arcdlc-aic
description: Design an initiative's architecture and write it down. Produces an AIC by default, or arc42, Tech Stack Canvas, TOGAF, C4 or ADR, in Markdown or HTML, and always grills the design one question at a time before writing anything. Use when someone says design the architecture, write an architecture doc, spec out the system, draw up the design, or do a technical design for X, or runs /arcdlc:aic <slug> [format], or invokes arcdlc-aic.
argument-hint: "<slug> [aic|arc42|tsc|togaf|c4|adr, comma-separated, :html for HTML]"
---

# ArcDLC AIC (/arcdlc:aic)

Produce the architecture document that anchors the ArcDLC delivery pipeline. This is the entry point
of the application track; `/arcdlc:policy` is the entry point of the governance track.

`/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive` → `/arcdlc:close`, with
`/arcdlc:remove` named beside it for deleting a design for good.

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
  `honestly`, `genuinely`, `truly`, `clearly`, `obviously`. Never a metaphor or a rhetorical
  question: a worked example with real values beats both. Every title that names work is a verb plus
  its object: "Implement the alerter service".
- **No long dashes.** A full stop, a comma, a colon, or brackets instead of `—` and `–`. Hyphens,
  flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.
- **Shape.** In chat, bullets rather than paragraphs: what you did, what you found, what comes next.
  In files, vary sentence length and form, and draw what a picture shows better, so a document never
  reads like a list or a form to fill in.
- **One paragraph, one line.** In every file and reply you write, a paragraph, a list item, or an
  acceptance criterion is one physical line, however long. Never break a line inside a sentence to
  fit a width: viewers wrap on their own, and a hard break shows up as a broken sentence in a ticket,
  a chat, or a rendered page. This bundle's own instruction files are wrapped at 100 columns; do not
  copy that into what you write. Code blocks, tables, and contract lines keep their own layout.

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

**A closed initiative is final.** If `docs/aics/<slug>/CLOSED.md` exists, stop and change nothing: say
that `<slug>` was closed on the date in its `- Closed:` line, and that follow-up work is a new
initiative with its own slug whose design links `docs/aics/<slug>/` as its starting point. No skill
edits or deletes a closed initiative's files. The one exception is the `## Superseded` line that
`/arcdlc:aic` appends to its `CLOSED.md` when another initiative's design reverses it.

## Argument: document format(s)

The optional **second** positional argument selects the format (the first is the slug above). Resolve
the template from `references/` next to this file:

| Argument | Output file | Template in `references/` |
| --- | --- | --- |
| *(none)* or `aic` | `docs/aics/<slug>/aic.md` | `AIC Template.md`, with `AIC Example.md` (a finished AIC to match) and `Diagram Conventions.md` (how to draw its diagrams). |
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
- `docs/aics/*/CLOSED.md`: closed initiatives. Read each title and `- Outcome:` line and open the
  designs related to this one; they are references and are never edited.
- `CONTEXT.md` / `CONTEXT-MAP.md` (domain glossary)
- `docs/adr/` (prior decisions)
- `docs/policies/` and its index (the organization's own engineering principles and standards). These
  are the technical constraints the design must follow. When there are none, ask the engineer where
  the organization's principles live before falling back on the bundled example
- `AGENTS.md`, `CLAUDE.md`, `README.md` of the target project
- The relevant template(s) from the table above — for `arc42`, both files, starting with `Arc42 Guide.md`;
  for `aic`, the template, the example and `Diagram Conventions.md`

## Step 2 — MANDATORY: grill the design

Never write the architecture document straight from the request — the whole point of `/arcdlc:aic` is
that the process is controlled: every section is answered in the interview before it is written.

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

**The template is the agenda.** An engineer designing a system thinks along one line, usually the
one the request came from, and loses the rest: the numbers, the failure cases, the person woken at
03:00. The template of each requested format names every part a finished document needs, so walk it:

- Hand the interview the sections of every requested template, in template order, as its topic
  list. When several formats are requested, ask each topic once, where two sections cover the same
  ground (arc42 §4 and the AIC hypotheses).
- Build the questions from each section's own guidance: its description and its **Ask** lines,
  which every template carries (for arc42, in the Section Catalogue of `Arc42 Guide.md`). What the
  code, the config or an ADR already answers, fill in, show, and move on.
- The Ask lines are where a section starts, not where it ends. Follow each answer down until the
  section is settled, as the grilling protocol does, before moving to the next section.
- **Grill a risk the moment it appears**, whether the engineer names it or you spot it in an answer.
  Write it into the risks section, then ask, one question per turn, each with your recommended
  answer: how likely it is and how bad it would be, then how it is mitigated. The engineer picks the
  mitigation: reduce it by design (a new or amended hypothesis), detect it (a metric and the alert
  value), accept it (the reason and who accepts it), defer it (an Open questions entry with an owner
  and a phase), or not a real risk (the reason). Never write a mitigation nobody chose.
- Before a theme that may not apply, ask whether it applies at all ("Does this initiative keep data
  of its own?"), so one answer closes several sections.
- Every section closes one of three ways, and the engineer chooses which:
  1. **Answered.** The decision goes into that section.
  2. **Deferred.** It goes under Open questions with who answers it and by which phase or stage.
     "Later" with no owner and no phase is not a deferral: ask for both.
  3. **Not applicable.** The section keeps one line, `Not applicable: <reason>`.

  Your recommended answer may itself be a deferral or not applicable. The engineer's answer stands,
  even one you would not give: say the risk once, in the recommendation, then record the answer.
  The architect decides; the agenda only makes sure nothing is decided by being forgotten.
- On a re-run, start with what the document already lists under Open questions, then the sections
  this round is about. That is how an initiative moves from one phase to the next. Any closed
  section can be reopened: a `Not applicable` line becomes an answer when the initiative grows into
  it, and an answered section gains detail. Ask the engineer which sections this round is about when
  the request does not say.
- **Write each section into the document the moment it closes**, before the next question: the
  answer, the deferral under Open questions, the `Not applicable` line, the risk with its chosen
  mitigation, and the section's diagram, compiled and looked at. Nothing goes in that the engineer
  did not answer. On a new initiative, create the file
  from the template at the first closed section, with the `# <Title>` and `> ` summary as far as they
  are known; on an existing one, amend it. Either way the rules of Step 3 govern every write. With
  several formats, write the section into each file that holds it.
- **The file is the interview's progress record.** After a context compaction, or in a new session,
  re-read the document: a filled section is closed, an empty template slot is open. Resume at the
  first open section. Never hold an answer only in the chat, because a compaction summary keeps the
  gist and can drop the figure.

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

Step 2 has already written each section as it closed. The rules here govern those writes and the
rest of the document alike.

**If the document already exists, amend it. Never regenerate it.** Re-running `/arcdlc:aic <slug>` is
the normal way a design matures: each round adds detail, tests an idea, or reverses a call. A document
holds two kinds of content, and they change differently:

- **Inputs: the goals, requirements and constraints.** The AIC Goals group (Business Case, Functional
  Overview, Quality Goals, both Constraints sections, Business Context), arc42 §1 to §3, and the
  matching sections of the other formats. They describe the problem as it is known today, and
  gathering data is expected to change them. When the interview changes one, rewrite it in place:
  no reversal note, git history keeps the old text.
- **Decisions: the hypotheses and everything built on them.** The AIC Solution part (Architectural
  Hypotheses and Solution Design) and Risks, arc42 §4 to §11, the ADRs. These are amended and superseded as the rules below say.
- **A changed input puts the decisions on it in question.** Name each hypothesis or ADR that rests on
  the input you changed and ask about it, one at a time: still holds, amend, or reverse. Never keep
  one silently, never change one silently.

For both kinds:

- Keep every section this interview did not touch, word for word. A detail added three rounds ago and
  not discussed today is still the decision. Rewriting the file from the template silently drops it,
  which is the one failure that makes iterating unsafe.
- Change only what the interview changed, and say which sections changed in your closing report.
- **A reversed decision is superseded, never deleted.** When this round rejects something an earlier
  round settled, both records have to move or the project now holds two answers:
  - Recorded as an ADR: write the new ADR, set `- Supersedes:` on it, and set the old one's
    `- Status:` to `Superseded by [ADR-NNNN](...)`. Never leave two ADRs that disagree:
    `/arcdlc:examinate` audits against `docs/adr/` and will file gaps for the decision you reversed.
  - Recorded only in the document body: replace the text. The body shows the current design only,
    so the old hypothesis leaves it. When the reason for the reversal would stop the next reader
    re-opening it, write an ADR for it; otherwise the commit message carries it. A rejected option
    may stay as one line where a reader would otherwise ask "why not X?".
  - Recorded in `CONTEXT.md`: a term whose meaning changed is redefined in place, not duplicated.
  - Recorded in a closed initiative's design: append one line to that initiative's `CLOSED.md` under
    a `## Superseded` heading, creating the heading at the end of the file when missing. The line is
    `- Superseded in part by [<title>](../<name>/aic.md) on <YYYY-MM-DD>: <which sections>.`, or
    `- Superseded by [<title>](../<name>/aic.md) on <YYYY-MM-DD>.` when the whole design is replaced,
    where `<name>` and `<title>` are this initiative's slug and H1. Never touch the closed design
    documents, and never start an `/arcdlc:aic` run on the closed initiative. When the reversed
    design belongs to an initiative that is still open (no `CLOSED.md`), write nothing into it:
    raise the disagreement as one grilling question, because two open designs that disagree are the
    engineer's decision.
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

### Check coverage before handover

Walk the template once more against the document, new or amended. Every section must hold content,
a `Not applicable: <reason>` line, or an Open questions entry with an owner and a phase, every
expected diagram must be there or carry its own `Not applicable` line, and every risk must carry a
mitigation the engineer chose. A section or a risk with none was skipped: ask about
it now, one question at a time, before Step 4. The check counts, it never grades an answer. List the deferred items in your closing report.

### Draw before you write

A diagram goes in wherever it can show what a paragraph would describe: structure, order, branching,
state, who hands what to whom. The AIC expects three: a C4 system context diagram in Business
Context, a C4 container diagram and a flowchart or BPMN diagram of the main process in Solution
Design. Each one missing needs a `Not applicable: <reason>` line. Draw in Graphviz DOT by the rules
in `references/Diagram Conventions.md` and `references/C4.md`, keep the `.dot` source next to the
PNG in `docs/aics/<slug>/images/`, then open the PNG and look at it: a diagram that compiles can
still be unreadable, and that file's checklist says what to fix. The text under a diagram adds only
what the picture cannot show.

### Write it for a human reader (every format, Markdown and HTML)

The bar is a new engineer understanding it without asking anyone, and a product owner trusting it.
`## Talk simple, write like a human` applies to every section of every format here. These points
matter most in an architecture document:

- Say the thing, then the reason, in that order: "We keep sessions in Postgres because we already
  run it and the volume is small."
- Break up stacked noun phrases ("install-agnostic cross-agent skill bundle distribution") into a
  sentence a person can say out loud.
- Expand every acronym on its first use. Keep a domain term when it is this project's own word and
  define it once, here or in `CONTEXT.md`. Never swap it for a vaguer everyday word.
- The `# <Title>` and `> ` summary follow the same rules. The summary is one plain sentence telling a
  stranger what this initiative does, not a compressed spec.
- **Goals stay short, Solution carries the reasons.** Functional Overview is the product box, in
  business words. Quality Goals are three attribute names; one gets a line of reason only when it is
  not plain and the engineer asks for it. Numbers and mechanics belong in Solution Design, and every
  choice there says why unless any engineer on the team would know.
- **The body holds the current design only.** No rejected hypothesis, no superseded option, no "this
  round", "reversed", or "settled on <date>": the ADRs and git history keep that.
- **No water.** Cut a sentence that restates its heading, copies a rule the repository already states
  (link it), explains a choice nobody would question, or walks through how the interview went.
- **Vary the form.** A hypothesis, a decision, or an explanation is written in whatever shape reads
  best: a sentence, a paragraph, a table, a worked example with real values. Never a fixed stack of
  labelled fields repeated for every item; that is what makes a document read like a form.
  `references/AIC Example.md` shows the target.

Run the pre-save check in `../grilling/references/Writing Style.md` (flat installs:
`../arcdlc-grilling/references/Writing Style.md`) before handing the draft over.

### Test the draft on a fresh reader

You wrote the document, so you cannot see what it leaves out. Before Step 4, give it to a reader with
no context. When you can spawn a subagent, spawn one with only the document's path and its images,
and ask it to return a list, quoting each line it means:

- what it could not understand without asking someone,
- a choice in Solution whose reason is missing and not plain,
- a sentence that carries no fact (water), or a passage a diagram would show better,
- a diagram it could not read, or one that disagrees with the text.

Without a subagent, do the same read yourself, section by section, as a new engineer would. Fix what
the list finds; a point only the engineer can settle is one grilling question, one per turn. Do not
paste the list into the document.

## Step 4 — Register and hand off

- Walk the user through the draft; iterate until approved.
- Register the initiative so it is discoverable. Probe once with `command -v arctool`; if present, run
  `arctool sync` to refresh the `<!-- arcdlc:initiatives -->` blocks in `AGENTS.md` and `README.md`
  from this new folder. *Fallback (no `arctool`): add the initiative by hand to those marker blocks —
  `- [<Title>](docs/aics/<slug>/<doc>) — <summary>`.*
- Tell the user the next step: `/arcdlc:plan <slug>` to decompose this document into the executable
  `docs/aics/<slug>/plan.md`.
