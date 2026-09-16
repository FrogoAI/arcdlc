---
name: arcdlc-policy
description: Author a governance policy in docs/policies/<name>.md following the ArcDLC "Policy of Policies" framework — mandatory grilled interview first, then the full policy document, then register it in docs/policies/README.md, the project README.md, and AGENTS.md. The policy name is the required first argument (e.g. /arcdlc:policy vacations). Use when the user runs /arcdlc:policy, invokes arcdlc-policy, or asks to create/write a company or engineering policy, SOP, or governance document.
argument-hint: "<name> [POL-GOV|POL-HR|POL-TECH|POL-ENG|POL-SEC|POL-DOC]"
---

# ArcDLC Policy (/arcdlc:policy)

Author a governance policy. This is the entry point of the governance track, and it produces a
document `/arcdlc:examinate` can audit code or process against:

`/arcdlc:policy <name>` → `/arcdlc:examinate <slug> docs/policies/<name>.md` → `/arcdlc:execute <slug>`

The audit and execution steps run against an initiative folder, so the slug is their first argument
and the policy path is the **second** argument of `/arcdlc:examinate`, never the first.

There is no plan step: the policy is the rules, and `/arcdlc:examinate` files each violation straight
into `docs/aics/<slug>/plan.md`. A clean audit, or a policy with no code impact, ends at the document.

Read `references/Policy of Policies.md` next to this file before writing. It is the contract for how
every policy is created, structured, and managed.

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

## Argument

- `name` (**required**, first positional) — the policy topic/slug; the output file is
  `docs/policies/<name>.md` (e.g. `/arcdlc:policy log-retention` → `docs/policies/log-retention.md`).
  If no name is given, **stop and report the error**, highlighting that the policy name is missing, and
  list the existing policies under `docs/policies/` so the user can pick a name or reuse one. Do not
  guess a name.
- optional `POL-*` class — the Unique ID class (see the classification table below). If omitted,
  choose it during the interview.

If the project keeps governance docs elsewhere, follow the project convention but keep the same
file structure and the central index.

## Step 0 — Gather context (before asking anything)

Read what already exists so the interview builds on it instead of repeating it:

- `references/Policy of Policies.md` next to this file. Also read `references/Policy of Initiatives.md`
  and `references/Engineering Principles.md` when relevant.
- The target project's `README.md` and `AGENTS.md` (you will register the policy in both).
- Existing `docs/policies/` — its index and the highest Unique ID number in each class, so the new
  policy continues the sequence instead of colliding.

## Step 1 — MANDATORY: grill the policy

Never write the policy straight from the request. A policy encodes real decisions (scope, owner,
approver, the actual rules), so interview first — this mirrors `/arcdlc:aic`.

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

Cover at minimum, mapping each answer to a template section:

- **Purpose & trigger** — what problem or risk this policy addresses.
- **Policy Statement** — the single binding rule in one or two sentences.
- **Scope** — personnel covered, documents/systems covered, and explicit exclusions.
- **Unique ID class** — `POL-GOV | POL-HR | POL-TECH | POL-ENG | POL-SEC | POL-DOC`.
- **Roles** — Policy Owner (accountable senior leader), Policy Approver (a C-level), reviewers.
- **Concrete rules** — the Allowed and Prohibited conduct, specific enough to audit against.
- **Lifecycle** — effective date, review cadence (default: semi-annual / every 6 months),
  consequences of non-compliance.

The interview ends only when the user confirms shared understanding or says to proceed.

## Step 2 — Assign the Unique ID

Pick the class from the classification table, then the next sequential number by scanning existing
`docs/policies/` entries in that class (e.g. `POL-ENG-003`).

| Identifier | Covers |
| --- | --- |
| `POL-GOV` | Overall corporate or organizational governance. |
| `POL-HR` | Managing the workforce. |
| `POL-TECH` | Using, developing, and managing technology systems. |
| `POL-ENG` | Engineering-specific practices, standards, and procedures. |
| `POL-SEC` | Protecting information, infrastructure, and systems. |
| `POL-DOC` | Documentation / records management. |

## Step 3 — Write `docs/policies/<name>.md`

Use the mandatory structure from `Policy of Policies.md`. The `## Talk simple, write like a human` rules above
cover every line of the policy: plain English, active voice, no filler, no long dashes, and domain
terms kept and defined once. Full rules and the pre-save check live in
`../grilling/references/Writing Style.md` (flat installs: `../arcdlc-grilling/references/Writing Style.md`).
Include every section:

- **Header block**: Policy Name, Unique ID, Author, Creation Date, `Status: Draft`, Effective Date,
  `Approval Date: TBD`, and Next Review Date (Creation Date + 6 months).
- **Purpose** — why the policy exists.
- **Policy Statement** — the binding rule; note it supersedes informal prior practice.
- **Scope** — personnel covered, documents/systems covered, exclusions.
- **Definitions** — key terms (reuse the framework's role definitions where they apply).
- **Procedures** — the step-by-step lifecycle or process the policy governs.
- **Roles & Responsibilities** — including a RACI matrix (Accountable / Responsible / Consulted /
  Informed) for the policy's key steps.
- **Allowed & Prohibited Conduct** — the concrete rules gathered in the interview.
- **Consequences of Non-Compliance**.
- **Revision History** — a table starting with v1.0 (this creation).

Every rule must trace to an interview answer or an existing source; do not invent scope, owners, or
consequences — if something is undecided, ask, or record it in an explicit "Open questions" note
rather than stating it as settled.

## Step 4 — Register the policy (required)

A policy is not effective until it is discoverable. Register it in three places:

1. **`docs/policies/README.md`** — the central table of contents (create if missing). One row per
   policy:

   ```md
   | ID | Policy | Status | Owner | Next Review |
   | --- | --- | --- | --- | --- |
   | POL-ENG-003 | [Log Retention](log-retention.md) | Draft | <owner> | 2027-01-03 |
   ```

2. **The project `README.md`** — add or extend a `## Policies` section linking to the new policy and
   the index (`docs/policies/README.md`). Do not duplicate the policy body; link to it.

3. **`AGENTS.md`** — add or extend a `## Policies` note stating the policy is binding on agents
   working in this repo and can be audited with `/arcdlc:examinate <slug> docs/policies/<name>.md`. List
   the policy with its Unique ID and path.

Keep all three in sync: the index row, the README link, and the AGENTS entry must reference the same
ID and path.

## Step 5 — Review and hand off

- Walk the user through the draft; iterate until approved. Remind them the policy is `Draft` until a
  Policy Approver signs off (per the framework, it may auto-activate two weeks after approval if
  reviewers do not respond).
- Suggest the next step: `/arcdlc:examinate <slug> docs/policies/<name>.md` to audit the codebase or
  process against the new policy, feeding any gaps into `docs/aics/<slug>/plan.md` for `/arcdlc:execute`.
