---
description: Author a governance policy in docs/policies/<name>.md following the ArcDLC "Policy of Policies" framework — mandatory grilled interview first, then the full policy document, then register it in docs/policies/README.md, the project README.md, and AGENTS.md. The policy name is the required first argument (e.g. /arcdlc:policy vacations). Use when the user runs /arcdlc:policy, invokes arcdlc-policy, or asks to create/write a company or engineering policy, SOP, or governance document.
argument-hint: "<name> [POL-GOV|POL-HR|POL-TECH|POL-ENG|POL-SEC|POL-DOC]"
---

# ArcDLC Policy (/arcdlc:policy)

Author a governance policy the same way ArcDLC builds software: a controlled process, not a
straight generation. This is the governance track of the delivery pipeline — it produces a policy
document that `/arcdlc:examinate` can then audit code or process against:

`/arcdlc:policy` → `/arcdlc:examinate docs/policies/<name>.md` → `/arcdlc:execute`

There is no plan step here: the policy is the rules, and `/arcdlc:examinate` already files each
violation as a `TODO` task in `docs/aics/<slug>/plan.md` for `/arcdlc:execute` to close. A policy with no
code impact (or a clean audit) ends the track at the document.

The governing framework is `source/Policy of Policies.md` in the sibling `source-map` skill. It is
the contract for how every policy must be created, structured, and managed — read it before writing.

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

- Invoke the `grill-with-docs` skill (which runs a `grilling` session).
- If `grill-with-docs` is not installed, run the same discipline inline: interview the user
  relentlessly, one question at a time, each with your recommended answer; explore the codebase or
  existing policies to answer questions instead of asking when you can.

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

Use the mandatory structure from `Policy of Policies.md`. Write in simple, direct, active voice;
avoid jargon and passive voice. Include every section:

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
   working in this repo and can be audited with `/arcdlc:examinate docs/policies/<name>.md`. List
   the policy with its Unique ID and path.

Keep all three in sync: the index row, the README link, and the AGENTS entry must reference the same
ID and path.

## Step 5 — Review and hand off

- Walk the user through the draft; iterate until approved. Remind them the policy is `Draft` until a
  Policy Approver signs off (per the framework, it may auto-activate two weeks after approval if
  reviewers do not respond).
- Suggest the next step: `/arcdlc:examinate docs/policies/<name>.md` to audit the codebase or
  process against the new policy, feeding any gaps into `docs/aics/<slug>/plan.md` for `/arcdlc:execute`.
