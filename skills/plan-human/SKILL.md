---
name: arcdlc-plan-human
description: Turn the executable docs/aics/<slug>/plan.md task queue into docs/aics/<slug>/plan-human.md, the engineer stories that go on a board, one story per service, each with numbered instructions and technical acceptance criteria, and each self-contained so a ticket needs no file from the repository. The initiative slug is the required first argument (e.g. /arcdlc:plan-human payments). Use when the user runs /arcdlc:plan-human, invokes arcdlc-plan-human, or asks for engineer stories, board tickets, or a human-readable plan.
argument-hint: "<slug>"
---

# ArcDLC Human Plan (/arcdlc:plan-human)

Turn the executable task queue into `docs/aics/<slug>/plan-human.md`: the stories an engineer picks up from a board.

Pipeline position: it sits **beside** `/arcdlc:plan`, not inside the chain. The pipeline stays `/arcdlc:aic` → `/arcdlc:plan` → `/arcdlc:execute` → `/arcdlc:archive`, and `/arcdlc:execute` never waits for this skill. Run it when the work also has to appear on a board.

`plan.md` and `plan-human.md` describe the same work at two sizes, on purpose.

| File | Reader | Unit | Sized for |
| --- | --- | --- | --- |
| `plan.md` | `/arcdlc:execute` | a task | one agent session |
| `plan-human.md` | a person on a board | a story | one thing an engineer can pick up and demo |

The story format is defined in `references/story-format.md`. Read it before writing anything. It is the contract, and a story that breaks it cannot be pasted into a tracker.

## Talk simple, write like a human

Plain English everywhere: short sentences, common words, one idea each, active voice, a named actor.

- **Replies to the user.** Bullets, not paragraphs. Say what you did, what you found, what comes next.
  No filler, no praise, no restating the request.
- **Files you write.** No AI filler ("Furthermore", "In conclusion", "It is important to note",
  "delve", "leverage", "robust", "seamless"), no warm-up opener, no invented summary. Vary sentence
  length. Concrete names, numbers, and paths, never "significantly improves performance".
- **Be direct, name the thing.** The fewest words that carry the fact. Cut empty intensifiers:
  `honestly`, `genuinely`, `truly`, `clearly`, `obviously`. Never a metaphor, a narrative line, or a
  question. Every title that names work is a verb plus its object: "Implement the alerter service".
- **No long dashes.** A full stop, a comma, a colon, or brackets instead of `—` and `–`. Hyphens,
  flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.

Short talk, full content. Brevity is for your replies, never for the files: never drop a rule, path,
decision, trade-off, or acceptance criterion to save space. The full standard, with examples and a
pre-save check, is `references/Writing Style.md` in the bundle's `grilling` skill.

## Judge by the four virtues

- **Wisdom.** Unclear is a question, not a guess. A contradiction, a missing decision, code that looks
  dead, a change that is risky or hard to reverse: grill it, never pick for the engineer.
- **Courage.** Say the hard thing. A plan a weaker model cannot execute is a plan defect, not a reason
  to raise the tier. Never stop to report a helper skill as missing.
- **Justice.** Write every answer where the next session reads it, never only in the chat. Name the
  tier you actually used. Never report a same-tier spawn as a cheaper run.
- **Temperance.** Touch only what you were pointed at. Cutting a comment out of the code is asked for,
  not assumed. No invented summary, no scope you were not given.

## Initiative selection

The slug is the **required first positional argument**: `/arcdlc:plan-human <slug>`. If it is missing, stop and list the initiatives under `docs/aics/`. Never guess.

## Step 1: locate inputs

- `docs/aics/<slug>/plan.md` is **required**. If it does not exist, stop and tell the user to run `/arcdlc:plan` first. Never write stories from a verbal description: the task queue is the source of truth for what the work is.
- `docs/aics/<slug>/aic.md` for the goals, the phase split, the alert or feature catalogue, and the constraints an engineer needs.
- `CONTEXT.md` at the repo root for the project's own words. Use them. Do not invent a synonym for a term the project already has.
- `docs/adr/` for the decisions a story has to state as fact.

Read all of `plan.md` before cutting a single story. A story that covers half a task is worse than no story.

## Step 2: cut the stories

**One story is one service, or one thing you can demo.** That is the whole rule. Two failures come from breaking it, and both are expensive.

- **Too coarse:** one story holding four services. Nobody can pick it up, review it, or say it is done. If a story would deliver more than one deployable unit, split it by unit.
- **Too fine:** one story per plan task. The board fills with slices of one piece of work, and two tickets read almost the same. If a task cannot be demonstrated on its own, it belongs inside a story with its neighbours.

Practical guidance:

- A story that names a service delivers **that whole service**, however many tasks that takes. Five tasks building one calculator are one story.
- Bundle two tasks into one story only when they are **one demo**: create a repository and deploy it, add the metrics and run the load test, one screen plus the field it needs.
- A capability that only exists so another story can be proven is not its own story. Put it in the story it serves, as its first steps.
- Take the **phase split from `plan.md`**. Do not invent one. If the phases look wrong, grill the engineer about it and stop; moving a phase boundary changes the architecture document, the diagrams and the tracker, and that is a decision for the engineer, not a side effect of writing stories.
- Group stories under a phase heading, and keep story numbers stable across runs. When a phase boundary moves, move the heading, not the numbers: a renumbered board loses every reference anyone has written down.

## When the plan is unclear, grill (mandatory)

A ticket is read by somebody who cannot ask you anything, so a story may not carry a guess. Stop and
ask the engineer the moment you hit one of these:

- A task leaves the story without a value it needs: a field name, a JSON shape, a subject name, an
  interval, a cap, a limit.
- Two tasks contradict each other, or a task contradicts the architecture document, an ADR, or
  `CONTEXT.md`.
- An `Acceptance` line cannot be reworded into "I do this, I get that" without inventing the
  observable result.
- A task names no deployable unit, so you cannot tell which story it belongs to.
- The phase split in `plan.md` does not match the architecture document.

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

Where the answer goes:

- A concrete value: into the story, written out in full, because the ticket cannot reach this
  repository.
- A new term: into `CONTEXT.md`, in the project's own words.
- An answer that changes the plan itself (a task, a phase boundary, an acceptance criterion): stop and
  say so. Moving a phase boundary changes the architecture document, the diagrams and the tracker, so
  it is the engineer's decision and `/arcdlc:plan` writes it, not this skill.

## Step 3: write each story

Follow `references/story-format.md` exactly: the title rule, the six parts, and the acceptance-criteria form.

Two rules decide whether the result is usable.

**Self-contained.** A ticket lives in a tracker, where nothing from this repository exists. So a story names no file in this repository and no ADR. If a story needs a field mapping, a JSON shape, or a subject name, **write it into the story**. Code identifiers stay: `cmd/alerter`, `internal/decoder/`, a config key and a metric name are the instruction itself, not a pointer to something to go and read.

**Named, not described.** A list in prose is not implementable. If the work is a set of metrics, fields, endpoints, or configuration keys, give them as a table with their exact names, types and labels. "Expose every metric the design depends on" is a defect; a table of seventeen metric names is a task. This is the single most common way a story looks finished and is not.

## Step 4: coverage gate

This is a **hard gate**. Do not hand off until it passes, and prove it with counts rather than asserting it.

1. **Every task appears in exactly one story.** No task covered twice, none missing. List the task ids per story and compare the set against `plan.md`.
2. **Every acceptance criterion survives.** Each `Acceptance` line in `plan.md` reaches some story's `ACCEPTANCE CRITERIA`, reworded from `GIVEN … WHEN … THEN …` into "I do this, I get that". Reworded, never dropped. Count them on both sides.
3. **Every story has acceptance criteria.** A story without them is not a story.
4. **No repository path and no ADR reference below the first story heading.** Grep for it.
5. **No long dash anywhere.** Grep for it.

Report the four counts to the user: stories, tasks covered, acceptance criteria carried, and violations found.

## Step 5: write the file

- Write `docs/aics/<slug>/plan-human.md` in the order `references/story-format.md` gives: the heading section, then the phase headings and their stories.
- **Mark the heading section as not ticket text.** Everything from the first story heading down is copied into a tracker as it stands. The section above it explains the initiative and is for the repository only. Say that in one line, or somebody will paste the whole file into one ticket.
- Record the story to task to ticket mapping in `CONTEXT.md` at the repo root, under the delivery section: story number, title, plan tasks, tracker id, phase. That table is how the next session knows which ticket holds which work.

## Step 6: sync a tracker, only when asked

Writing the file is this skill's job. Touching a tracker is not, unless the user asks for it in this session.

- **Never create, edit or delete a ticket without an explicit go.** Show the user the exact list of summaries first: what is rewritten, what is created, what changes phase. A tracker is outward-facing, and other people get notified.
- **Reuse ticket ids.** When a story splits, keep the existing id for one part and create the rest. Never delete a ticket to make the numbering tidy: check first for an assignee, comments, attachments, logged time, votes, and links from other issues.
- **Summary convention:** whatever the project already uses, taken from the tickets that exist. Follow it exactly, including any team prefix and any phase suffix.
- **Write created ids back** into `plan-human.md` and the `CONTEXT.md` mapping table in the same session. An id that lives only in a chat log is lost.
- If the project has a ticket that already covers a story, map the story to it and leave the ticket alone. Do not paste a specification into a ticket that is already in review.

## Ending

- Report the counts from Step 4, the stories per phase, and anything you had to leave open.
- Name the next step: `/arcdlc:execute <slug>` implements the tasks; the stories are what the board shows while it happens.
