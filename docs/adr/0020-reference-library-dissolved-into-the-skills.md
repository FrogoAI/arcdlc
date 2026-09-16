# ADR-0020 — The reference library is dissolved into the skills that use it

- Status: Accepted
- Date: 2026-09-16
- Initiative: none (bundle-wide, removes the `source-map` skill)
- Amends: [ADR-0007](0007-mdca-standard-is-normative-single-document.md), whose consequence required
  the `skills/source-map/SKILL.md` routing table to be updated alongside any MDCA path change. There
  is no routing table now, so that obligation ends. The decision itself, one normative MDCA document,
  stands.

## Context

The `source-map` skill shipped 31 documents, about 500 KB, reached through a keyword routing table.
It was built so an agent could answer without searching the internet.

Two things stopped being true.

A current model already knows SOLID, DDD, Clean Code, the Twelve-Factor App, trunk-based development,
the Conventional Commits spec, and the UML, BPMN and flowchart notations. Shipping those as prose
bought nothing an agent would not produce anyway.

Worse, the documents carrying rules the bundle actually wants followed were unreachable. `stoic.md`
held the four virtues, and its routing keywords were "Leadership, Human-Behavior, Philosophie,
Conflict-Solving": misspelled, and absent from the skill's own `description:` frontmatter, so the
skill could not fire on them. `CTO Methodology Guide.md` existed to make the agent think
strategically and was reachable the same broken way. `tbd.md` described how these projects branch and
no skill mentioned it. Each was read by nothing.

That is the same failure the writing-style rules already avoid by sitting inline in every `SKILL.md`
rather than only in `Writing Style.md`. A reference is read only if something opens it.

## Decision

Apply one test to every bundled document: **does this change what the agent produces, in a way the
model would not do anyway?** Then place it by kind, and delete what fails.

| Kind | Home |
|------|------|
| A rule the agent must always follow | inline in the `SKILL.md`, loaded every invocation |
| A template a skill copies | `references/` inside that skill, because it is data, not reading |
| A rule set whose identifiers get cited | `references/` inside the skill that cites them |
| Public knowledge nothing cites | deleted |

Concretely:

- The four virtues become `## Judge by the four virtues`, a block in all ten `SKILL.md` files, pinned
  byte-identical by CI exactly as the writing-style block is. Each virtue names the concrete ArcDLC
  rule it governs, so it reads as operating instruction.
- Trunk-based development moves into `skills/execute/SKILL.md`, which is the skill that branches and
  commits.
- Rumelt's diagnosis, guiding policy and coherent actions move into `/arcdlc:aic` Step 2, which is
  the interview where strategy is actually decided.
- Generator templates move to `skills/aic/references/`, audit rule sets to
  `skills/examinate/references/`, the policy framework to `skills/policy/references/`, and the shared
  `Writing Style.md` to `skills/grilling/references/`.
- `skills/source-map/` is deleted, and with it the routing table.

## Consequences

- The bundle ships ten skills, not eleven. `install.sh` carries a `LEGACY_SUBSKILLS` sweep so an
  upgrade removes a stranded `arcdlc-source-map` directory instead of leaving it behind, and the CI
  installer smoke test plants one to prove the sweep runs.
- Most dual-path boilerplate is gone. Twenty-two `../source-map/...` plus `../arcdlc-source-map/...`
  pairs became local `references/...` paths that are identical in both install layouts. Only
  `grilling` (writing style) and `plan` (the plan format) are still reached across skills.
- Discovery changes. There is no keyword table routing a general question to a document. `/arcdlc:aic`
  and `/arcdlc:policy` carry the pipeline overview instead, because they are where a user enters.
- `mdca.md`, `solid.md`, `ddd.md` and `ECS.md` survive untrimmed. Their numbered rule identifiers are
  cited by gap blocks already filed, so they are a compatibility surface, not prose.
- Anyone who invoked `arcdlc-source-map` by name loses it. That is breaking, and the version bump to
  0.19.0 carries it.

## Alternatives considered

- **Keep the skill and repair the routing.** Fixes discoverability but not the core problem: a rule
  the agent must always follow cannot depend on the agent choosing to open a file.
- **Delete the three unread documents outright.** Rejected once their intent was known. The content
  was wanted; only its delivery was broken.
- **Trim the public-knowledge documents again.** A 2026-08 cleanup already did that, and `tbd.md` had
  to be partly restored afterwards. Trimming does not fix a document nothing reads.
