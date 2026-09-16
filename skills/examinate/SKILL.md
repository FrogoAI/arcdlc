---
name: arcdlc-examinate
description: Examine existing code for compliance with a named architecture, policy, or design (e.g. /arcdlc:examinate MDCA — also DDD, SOLID, ECS, Modern Go, the Go architecture guides, KISS, or any named standard such as Clean Code or Twelve-Factor), with a project policy authored by /arcdlc:policy (e.g. /arcdlc:examinate docs/policies/log-retention.md), or with the project's own AIC. Records violations as gap blocks in docs/aics/<slug>/gap.md and adds matching TODO tasks to docs/aics/<slug>/plan.md. Use when the user runs /arcdlc:examinate, invokes arcdlc-examinate, or asks for a compliance audit / gap analysis of the codebase.
argument-hint: "<slug> [MDCA|DDD|SOLID|...|policy-path]"
---

# ArcDLC Examinate (/arcdlc:examinate)

Audit the existing codebase against a policy or design, register every gap as evidence, and feed the gaps into the
executable plan so `/arcdlc:execute` can close them.

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

Only the length differs. A reply is short; a file is complete. Never drop a rule, path, decision,
trade-off, or acceptance criterion from a file to save space, and never pad a reply to fill one. The
full standard, with examples and a pre-save check, is `references/Writing Style.md` in the bundle's
`grilling` skill.

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

The initiative slug is the **required first positional argument**: `/arcdlc:examinate <slug> [target]`.
Gaps and their mirrored tasks are filed into `docs/aics/<slug>/` (as `gap.md` and `plan.md`), and the
slug is passed to `arctool` as `--aic <slug>`. If the slug is missing, stop and report the error,
listing the existing initiatives under `docs/aics/` — never guess. If the named initiative folder does
not exist yet (a fresh audit, e.g. `mdca-audit`), confirm the slug with the user and create
`docs/aics/<slug>/`.

## Step 1 — Resolve the standard to audit against

- With a standard as the second argument (e.g. `/arcdlc:examinate <slug> MDCA`): read the matching rule set from `references/` next to this file, in full — e.g. MDCA → `mdca.md`; DDD → `ddd.md`; SOLID → `solid.md`;
  Go architecture → `Go Server.md` / `Go Client.md` / `Go Library.md`; modern Go → `Modern Go.md`.
  For `Modern Go.md`, read the `go` line in the target project's `go.mod` first. A rule that needs a
  newer Go than the project declares is not a violation: report it once in the register intro as an
  upgrade opportunity, and file no gap block for it.
- **With a standard that has no bundled rule set** (Clean Code, Twelve-Factor, trunk-based development,
  a language style guide): audit against it from your own knowledge of that standard. Say in the gap
  register's intro that the rule set was not bundled, and cite each finding by the standard's own named
  principle ("Twelve-Factor III, config in the environment"), never by an invented identifier. Only
  `mdca.md`, `solid.md`, `ddd.md`, `ECS.md` and `Modern Go.md` carry identifiers a gap block may cite,
  as `MDCA-P3.2`, `SOLID-7`, `GO-12` and the like.
- With a project policy path as the second argument (e.g. `/arcdlc:examinate <slug> docs/policies/log-retention.md`,
  typically one authored by `/arcdlc:policy`): read that policy in full and extract its Allowed/Prohibited rules as the
  checkable rule set.
- With no second argument (no standard or policy path): audit against the initiative's own architecture — `docs/aics/<slug>/aic.md` (or other docs in
  that folder), `docs/adr/`, and `CONTEXT.md`. If none of these exist either, stop and ask which policy to audit
  against (or suggest `/arcdlc:aic` first).
- **Precedence — project policy beats the bundled reference.** Before auditing against any bundled reference, check
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
- Do not settle a finding that needs the engineer. Dead-looking code, a contradiction, an ambiguous rule, a risky
  fix, or a violation that may be deliberate all go through Step 2.5 first.

## Step 2.5 — When a finding needs a human, grill (mandatory)

Some findings are not yours to settle. Recording a fix the engineer never approved is worse than
recording no fix: `/arcdlc:execute` will carry it out. Stop and ask the moment you hit one of these:

- **Code that looks dead.** You can find no caller in this repository, which does not mean there is
  none: another service, a plugin, a test harness, or a customer may call it. Deleting it is the
  engineer's call.
- **A contradiction.** The policy contradicts itself, another policy, an ADR, `CONTEXT.md`, or the
  architecture document. Pick nothing until the engineer says which one wins.
- **An ambiguous rule.** The rule can be read two ways and the two readings produce different gaps.
- **A risky or hard-to-reverse fix.** A data migration, a deleted public interface, a changed wire
  format, a dependency swap, a rename that breaks callers outside this repository.
- **A violation that may be deliberate.** The code reads like a conscious exception, so the question is
  whether the exception stands, not how to remove it.

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

Each answer lands in the gap register, one way or the other:

- **Fix approved:** write the engineer's decision into the gap's `HOW`, then mirror the gap into the
  plan as usual.
- **Deviation accepted:** add `- Accepted: <reason>, <date>.` to the gap block. It stays in `gap.md` as
  evidence and gets **no** task in `plan.md`. An accepted deviation with no written reason is not
  accepted, it is forgotten.
- **A decision that outlives the audit:** write an ADR in `docs/adr/` and cite it in the gap's
  `References`.
- **A new term:** add it to `CONTEXT.md`.

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

- `- Accepted: <reason>, <date>.` is optional. It appears only on a gap the engineer accepted in Step 2.5, and it
  keeps that gap out of the plan.
- `<PREFIX>` is the audited policy or initiative (e.g. `MDCA-GAP-01`, `AIC-GAP-03`). Number gaps sequentially,
  continuing from existing entries — never renumber or delete previous gaps.
- If a gap from an earlier examination is now fixed, mark it in `gap.md` (e.g. append `— resolved <date>`) instead of
  removing it.

## Step 4 — Sync gaps into the plan

Per the Register Sync rules in `../plan/references/plan-format.md` (flat installs:
`../arcdlc-plan/references/plan-format.md`), append a matching task block to `docs/aics/<slug>/plan.md` for every new
gap:

- Same task ID and heading as in `gap.md`; same `WHAT`, `HOW` (when present), `WHERE`, `WHY`, and `Acceptance` content.
- Add runner metadata: `References` (must include `docs/aics/<slug>/gap.md`, the policy source, and the architecture
  doc if relevant) and `- Status: TODO.`
- Append after the existing blocks; never modify existing tasks or reuse an ID already present in the plan.
- Skip a gap carrying `- Accepted:` from Step 2.5. It stays in `gap.md` as evidence and gets no task, and the
  `- Accepted:` line never reaches the plan.
- If `docs/aics/<slug>/plan.md` does not exist yet, create it per the format guide (a `/arcdlc:plan` run can merge it
  with architecture-driven tasks later).
- Validate the plan afterwards: `arctool validate --strict --aic <slug>` (probe with `command -v arctool`,
  or install with `make install` from the arcdlc repo root). Without it, say so once and hand-check the
  "Authoring Rules" section of the format guide. Fix every finding before handoff.

## Step 5 — Report

Summarize: policy audited, rules checked, gaps found by status (MISSING/PARTIAL/DRIFT) with severity, and how many
tasks were added to the plan. Name every finding you took to the engineer in Step 2.5 and what was decided, and list
the accepted deviations separately: they are gaps on purpose, and nobody should read them as work left undone. Suggest
`/arcdlc:execute <slug>` to start closing the rest.
