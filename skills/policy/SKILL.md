---
name: arcdlc-policy
description: Author a governance policy in docs/policies/<name>.md following the ArcDLC "Policy of Policies" framework — mandatory grilled interview first, then the full policy document, then register it in docs/policies/README.md, the project README.md, and AGENTS.md. The policy name is the required first argument (e.g. /arcdlc:policy vacations). Use when the user runs /arcdlc:policy, invokes arcdlc-policy, or asks to create/write a company or engineering policy, SOP, or governance document.
argument-hint: "<name> [POL-GOV|POL-HR|POL-TECH|POL-ENG|POL-SEC|POL-DOC]"
---

# ArcDLC Policy (/arcdlc:policy)

Author a governance policy the same way ArcDLC builds software: a controlled process, not a
straight generation. This is the governance track of the delivery pipeline — it produces a policy
document that `/arcdlc:examinate` can then audit code or process against:

`/arcdlc:policy <name>` → `/arcdlc:examinate <slug> docs/policies/<name>.md` → `/arcdlc:execute <slug>`

The audit and execution steps run against an initiative folder, so they take an initiative slug as
their first argument — the policy path is the **second** argument of `/arcdlc:examinate`, never the first.

There is no plan step here: the policy is the rules, and `/arcdlc:examinate` already files each
violation as a `TODO` task in `docs/aics/<slug>/plan.md` for `/arcdlc:execute` to close. A policy with no
code impact (or a clean audit) ends the track at the document.

The governing framework is `source/Policy of Policies.md` in the sibling `source-map` skill. It is
the contract for how every policy must be created, structured, and managed — read it before writing.

## Talk simple, write like a human

Plain B2 English, everywhere: short sentences, one idea each, active voice, a named actor.

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

- `source/Policy of Policies.md` via the sibling `source-map` skill (from this file:
  `../source-map/source/Policy of Policies.md` in the plugin layout;
  `../arcdlc-source-map/source/Policy of Policies.md` in flat installs). Also read
  `source/Policy of Initiatives.md` and `source/Engineering Principles.md` when relevant.
- The target project's `README.md` and `AGENTS.md` (you will register the policy in both).
- Existing `docs/policies/` — its index and the highest Unique ID number in each class, so the new
  policy continues the sequence instead of colliding.

## Step 1 — MANDATORY: grill the policy

Never write the policy straight from the request. A policy encodes real decisions (scope, owner,
approver, the actual rules), so interview first — this mirrors `/arcdlc:aic`.

The interview runs on this bundle's own grilling skill, `arcdlc-grilling` (`/arcdlc:grilling`), a
sibling of this one. ArcDLC depends on no external grilling skill: do not look for one, and never
stop to report a skill as missing.

1. **Invoke `arcdlc-grilling`** — the normal path. It ships in every ArcDLC install and is
   model-invocable.
2. **Inline**, only if that skill is genuinely not invocable here: run the same protocol yourself,
   reading it from `../grilling/SKILL.md` (plugin layout) or `../arcdlc-grilling/SKILL.md` (flat
   installs). What is mandatory is the grilled interview, not the invocation.

Whichever path runs, all of this holds:

- **One question at a time.** Ask one question, wait for the answer, then ask the next. Never a
  numbered round, never "two quick ones".
- Every question carries your recommended answer and one line of why.
- Facts are yours to find: if the codebase or an existing policy answers it, look it up, do not ask.

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
`../source-map/source/Writing Style.md` (flat installs: `../arcdlc-source-map/source/Writing Style.md`).
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
rather than stating it as settled (the `source-map` rule).

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
