---
name: arcdlc-examinate
description: Examine existing code for compliance with a named architecture, policy, or design (e.g. /arcdlc:examinate MDCA — also DDD, SOLID, Clean Code, Go Server, Twelve-Factor, ECS), with a project policy authored by /arcdlc:policy (e.g. /arcdlc:examinate docs/policies/log-retention.md), or with the project's own AIC. Records violations as gap blocks in docs/aics/<slug>/gap.md and adds matching TODO tasks to docs/aics/<slug>/plan.md. Use when the user runs /arcdlc:examinate, invokes arcdlc-examinate, or asks for a compliance audit / gap analysis of the codebase.
argument-hint: "<slug> [MDCA|DDD|SOLID|...|policy-path]"
---

# ArcDLC Examinate (/arcdlc:examinate)

Audit the existing codebase against a policy or design, register every gap as evidence, and feed the gaps into the
executable plan so `/arcdlc:execute` can close them.

## Talk simple, write like a human

Plain B2 English, everywhere: short sentences, one idea each, active voice, a named actor. Cut
empty intensifiers: `honestly`, `genuinely`, `truly`, `clearly`, `obviously` add nothing.

- **Replies to the user:** bullets, not paragraphs. No filler, no praise, no restating the request.
  Say what you did, what you found, what comes next.
- **Files you write:** no AI filler ("Furthermore", "In conclusion", "It is important to note",
  "delve", "leverage", "robust", "seamless", "In today's fast-paced world"), no warm-up opener, no
  invented summary. Vary sentence length. Concrete names, numbers, and paths, never "significantly
  improves performance".
- **Be direct, name the thing.** The fewest words that carry the fact; a word that only adds tone
  comes out. Never a metaphor, a narrative line, or a question. Say what must happen: "The alerter
  retries three times, then dead-letters." Every title that names work is an instruction, a verb
  plus its object: "Implement the alerter service", never "One message out, and the three places it
  goes".
- **No long dashes.** Use a full stop, a comma, a colon, or brackets instead of `—` and `–`.
  Hyphens, flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.

Short talk, full content. Brevity is for your replies, never for the files: never drop a rule, path,
decision, trade-off, or acceptance criterion to save space. The full standard, with examples and a
pre-save check, is `source/Writing Style.md` in the bundle's `source-map` skill.

## Initiative selection

The initiative slug is the **required first positional argument**: `/arcdlc:examinate <slug> [target]`.
Gaps and their mirrored tasks are filed into `docs/aics/<slug>/` (as `gap.md` and `plan.md`), and the
slug is passed to `arctool` as `--aic <slug>`. If the slug is missing, stop and report the error,
listing the existing initiatives under `docs/aics/` — never guess. If the named initiative folder does
not exist yet (a fresh audit, e.g. `mdca-audit`), confirm the slug with the user and create
`docs/aics/<slug>/`.

## Step 1 — Resolve the standard to audit against

- With a standard as the second argument (e.g. `/arcdlc:examinate <slug> MDCA`): look the policy up in the sibling `source-map` skill's table
  (from this file: `../source-map/source/` in the plugin layout, `../arcdlc-source-map/source/` in flat installs)
  and read every listed reference in full — e.g. MDCA → `mdca.md`; DDD → `ddd.md`; SOLID → `solid.md`;
  Go architecture → `Go Server.md` / `Go Client.md` / `Go Library.md`.
- With a project policy path as the second argument (e.g. `/arcdlc:examinate <slug> docs/policies/log-retention.md`,
  typically one authored by `/arcdlc:policy`): read that policy in full and extract its Allowed/Prohibited rules as the
  checkable rule set.
- With no second argument (no standard or policy path): audit against the initiative's own architecture — `docs/aics/<slug>/aic.md` (or other docs in
  that folder), `docs/adr/`, and `CONTEXT.md`. If none of these exist either, stop and ask which policy to audit
  against (or suggest `/arcdlc:aic` first).
- **Precedence — project policy beats the bundled reference.** Before auditing against any `source/` file, check
  whether the project redefines the same subject in `docs/policies/*.md`, `AGENTS.md` / `CLAUDE.md`, or `docs/adr/`.
  Where it does, the project document is normative: audit against its rule and drop the bundled one for every rule
  they both cover, keeping the bundled reference only for subjects the project leaves unspecified. Code that follows a
  project rule is never a gap, even when the bundled reference says otherwise. Note the override once in the gap
  register's intro (which document won, for which rule); if the conflict looks unintended, flag it to the user instead
  of filing gaps either way.
- Extract the concrete, checkable rules from the resolved rule set before looking at code, so findings cite a rule,
  not a feeling.

## Step 2 — Examine the code

- Sweep the codebase systematically against each rule: package/module layout, layer boundaries, dependency
  direction, naming, error handling, tests — whatever the policy governs.
- Every finding needs evidence: `file:line` (or package/module) plus the rule it violates. No evidence, no gap.
- Classify each finding's source status: `MISSING` (required element absent), `PARTIAL` (present but incomplete),
  or `DRIFT` (present but violates the policy).
- Do not report style preferences that the policy does not mandate.

## Step 3 — Write the gap register

Write or update `docs/aics/<slug>/gap.md` — the evidence register. One gap per `###` block, using the plan heading
format so it can be mirrored into the plan verbatim. Write each block for a **less capable executor**: the mirrored
task will be implemented by whatever model runs `/arcdlc:execute`, reading only the block and its references — so
the block must carry the fix decision, not just the complaint. Word every field the way
`## Talk simple, write like a human` says: name the file, the rule, and the effect, in plain
sentences with no hedging and no long dashes.

Use this block format:

```md
### <PREFIX>-GAP-NN (<MISSING|PARTIAL|DRIFT>): <Short Title>

- WHAT: <What must change to become compliant, one line.>
- HOW:
  <The concrete fix you already know from the audit: target structure, naming, what moves where.
  Out of scope: <adjacent non-violations the executor must leave alone>.>
- WHERE: <Exact files/modules with the violation, with file:line evidence.>
- WHY: <The violated rule, quoted or paraphrased, with the policy source path.>
- Acceptance:
  - GIVEN <the violating code> WHEN <the concrete re-check: the named test, lint rule, grep, or build command> THEN <the observable compliant result>.
```

Precision rules for each gap block:

- `<Short Title>` is an instruction that names the thing: an imperative verb plus its object, the real
  component named ("Add the retry limit to the alerter"). Never a turn of phrase, never a question.
- `WHAT` names the change, not the finding ("Move DB access behind a repository port", not "handler
  violates layering").
- `HOW` (optional but preferred) records the fix decisions the audit already surfaced — where the code
  should land, what the compliant shape looks like — so the executor does not re-derive them. Fence
  adjacent non-violations with `Out of scope:` to prevent over-fixing.
- `WHERE` cites evidence as exact `file:line` (or package/module) — the executor edits only these.
- `Acceptance` must be a runnable re-check wherever one exists (the lint rule, a named test, a grep
  that must come up empty, a build that must pass) — "re-examine and the violation is gone" is the
  fallback, not the default. This also lets the mirrored plan task pass `arctool validate --strict`
  (which requires an `Acceptance` section).

- `<PREFIX>` is the audited policy or initiative (e.g. `MDCA-GAP-01`, `AIC-GAP-03`). Number gaps sequentially,
  continuing from existing entries — never renumber or delete previous gaps.
- If a gap from an earlier examination is now fixed, mark it in `gap.md` (e.g. append `— resolved <date>`) instead of
  removing it.

## Step 4 — Sync gaps into the plan

Per the Gap Register Sync rules in `../plan/references/plan-format.md` (flat installs:
`../arcdlc-plan/references/plan-format.md`), append a matching task block to `docs/aics/<slug>/plan.md` for every new
gap:

- Same task ID and heading as in `gap.md`; same `WHAT`, `HOW` (when present), `WHERE`, `WHY`, and `Acceptance` content.
- Add runner metadata: `References` (must include `docs/aics/<slug>/gap.md`, the policy source, and the architecture
  doc if relevant) and `- Status: TODO.`
- Append after the existing blocks; never modify existing tasks or reuse an ID already present in the plan.
- If `docs/aics/<slug>/plan.md` does not exist yet, create it per the format guide (a `/arcdlc:plan` run can merge it
  with architecture-driven tasks later).
- After updating the plan, validate it. Prefer `arctool validate --strict --aic <slug>` (probe once with
  `command -v arctool`, or install it from the arcdlc repo root: `make install`); fix any findings before handoff.
  If `arctool` is unavailable, say so once and hand-check unique IDs, present/uppercase `Status`, and required keys per
  `../plan/references/plan-format.md`.

## Step 5 — Report

Summarize: policy audited, rules checked, gaps found by status (MISSING/PARTIAL/DRIFT) with severity, and how many
tasks were added to the plan. Suggest `/arcdlc:execute <slug>` to start closing them.
