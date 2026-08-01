---
name: arcdlc-source-map
description: Use when a user asks for architecture, architecture decision records or ADRs, engineering governance, delivery workflow, Jira decomposition, commit message conventions or Conventional Commits, Go architecture, Go best practices, clean code, KISS, simple design, modularity, DDD, SOLID, ECS, MDCA, TOGAF, ArchiMate, AIC, arc42, C4, UML, BPMN, flowcharts, Twelve-Factor App guidance, tech stack decisions, CTO methodology, or asks which ArcDLC source reference file to read.
---

# ArcDLC Source Map

Use the `source/` directory next to this `SKILL.md` as the reference library. Read only the files relevant to the user's request.

This skill is routing guidance only. Do not infer facts about the user's current project from these references unless the user explicitly asks you to apply them.

Do not imagine, invent, or silently assume important architecture conclusions or decisions. If a request requires deciding architecture, ownership boundaries, tech stack, data model, communication pattern, deployment model, or governance tradeoff, and the project context is missing or ambiguous, ask the user for the missing decision or context before presenting the conclusion as final. If work must proceed before confirmation, label the result as a tentative option and list the assumptions explicitly.

If a project has its own `AGENTS.md`, `CLAUDE.md`, or README guidance, follow that project guidance together with this source map.

## Delivery Workflow Commands

The delivery pipeline is handled by the sibling skills of this `arcdlc` bundle. On Claude Code they are plugin commands (`/arcdlc:<name>`); on Codex, OpenCode, Cursor, and Antigravity the same skills are installed flat as `arcdlc-<name>` (e.g. `arcdlc-aic`) — identical behavior, invoked by skill name instead of a slash command. ArcDLC is a universal delivery tool: it builds **applications** and authors **policies**, and both feed the same executable plan queue (`docs/aics/<slug>/plan.md`): architecture is decomposed into it, audit findings are filed into it, and `/arcdlc:execute` works it off task by task. When a user asks to design, plan, audit, implement, or govern an initiative end-to-end, route them through this pipeline instead of improvising:

- **Application track:** `/arcdlc:aic <slug>` → `/arcdlc:plan <slug>` → `/arcdlc:execute <slug>` → `/arcdlc:archive <slug>`
- **Governance track:** `/arcdlc:policy <name>` → `/arcdlc:examinate <slug> docs/policies/<name>.md` → `/arcdlc:execute <slug>`

Each initiative lives in its own folder `docs/aics/<slug>/` (holding the architecture doc, `plan.md`, `gap.md`, `plan-archive.md`). **Selection is mandatory and never inferred:** every pipeline skill takes the slug as its required first positional argument, and a missing slug is an error that lists the initiatives under `docs/aics/` instead of guessing one. `arctool` mirrors that rule with `--aic <slug>` (or `--plan PATH`); given neither it lists the initiatives and exits `2`. The legacy flat `docs/aics/plan.md` is reachable only via `--plan`.

The governance track has no separate plan step: a policy is rules, not work to decompose. `/arcdlc:examinate` files each violation as a `TODO` task directly into `docs/aics/<slug>/plan.md`, and `/arcdlc:execute` closes them; with no findings (or a process-only policy) the track ends at the policy document.

| Stage | Command | Output |
| --- | --- | --- |
| Architecture document (grilled interview first, then AIC / arc42 / TOGAF / C4 / ADR) | `/arcdlc:aic <slug> [format]` | `docs/aics/<slug>/<format>.md`, ADRs, `CONTEXT.md` |
| Governance policy (grilled interview first, per the Policy of Policies framework) | `/arcdlc:policy <name>` | `docs/policies/<name>.md`, `docs/policies/README.md` index, README/AGENTS refs |
| Decompose the document into the executable task queue | `/arcdlc:plan <slug>` | `docs/aics/<slug>/plan.md` |
| Examine code for policy/design compliance and register gaps as plan tasks | `/arcdlc:examinate <slug> [policy]` | `docs/aics/<slug>/gap.md`, new TODO blocks in `docs/aics/<slug>/plan.md` |
| Implement plan tasks (all, or one by ID) under the plan status contract | `/arcdlc:execute <slug> [TASK-ID]` | code, tests, one Conventional Commits commit per task |
| Archive DONE tasks to keep the plan small | `/arcdlc:archive <slug>` | `docs/aics/<slug>/plan-archive.md` |
| Retire a finished initiative (always confirmed first) | `/arcdlc:remove <slug>` | deleted `docs/aics/<slug>/`, refreshed registry |

The plan task format lives in `../plan/references/plan-format.md` (flat installs: `../arcdlc-plan/references/plan-format.md`).
Every plan task carries testable `Acceptance` criteria, which `/arcdlc:execute` must demonstrate before a task is `DONE`.

Optional accelerator: `arctool` is a stdlib-only Go CLI (source in `cmd/arctool/` and `internal/plan/` at the root of the arcdlc repository) that implements the
`plan-format.md` contract mechanically, covering the full lifecycle: `validate` (`--strict` also requires an
`Acceptance` section per task), `next`, `show`, `list`, status flips (`take`/`done`/`block`/`todo`), `archive`, and
`sync` (regenerates the initiative registry in `AGENTS.md`/`README.md`; `--check` fails on drift).
It is always optional — every pipeline skill probes `command -v arctool` and falls back to manual markdown handling
when it is absent.

## Source Map

| User request                                                  | Reference files                                                                                       |
|---------------------------------------------------------------|-------------------------------------------------------------------------------------------------------|
| AIC or Architecture Inception Canvas                          | `source/AIC Template.md`                                                                              |
| ADR or architecture decision record                           | `source/ADR.md`                                                                                       |
| Architecture documentation                                    | `source/Arc42.md`, `source/arc42/arc42-template-EN.md`                                                |
| Enterprise architecture, TOGAF, or ArchiMate                  | `source/TOGAF.md`                                                                                     |
| Tech stack decisions                                          | `source/Tech Stack Canvas.md`, `source/Tech Stack Canvas Original.md`                                 |
| Go server architecture                                        | `source/Go Server.md`                                                                                 |
| Go client architecture                                        | `source/Go Client.md`                                                                                 |
| Go libraries                                                  | `source/Go Library.md`                                                                                |
| Go best practices or clean code                               | `source/Go Best Practice.md`, `source/Clean Code.md`                                                  |
| KISS, simple design, modularity, or avoiding over-engineering | `source/KISS.md`, `source/Clean Code.md`, `source/Go Best Practice.md`                                |
| MDCA                                                          | `source/mdca.md`, `source/mdca_standard.md`                                                           |
| DDD                                                           | `source/ddd.md`                                                                                       |
| SOLID                                                         | `source/solid.md`                                                                                     |
| ECS                                                           | `source/ECS.md`                                                                                       |
| Diagrams                                                      | `source/C4.md`, `source/UML.md`, `source/BPMN.md`, `source/Flowchart.md`                              |
| Delivery, workflow, or Jira                                   | `source/Workflow Policy.md`, `source/Task Decomposition and Jira Sync.md`, `source/tbd.md`            |
| Commit messages, changelogs, or semantic version bumps        | `source/Conventional Commits.md`                                                                      |
| Executable plan format (docs/aics/<slug>/plan.md)             | `../plan/references/plan-format.md` (flat: `../arcdlc-plan/references/plan-format.md`)               |
| Engineering governance                                        | `source/Engineering Principles.md`, `source/Policy of Policies.md`, `source/Policy of Initiatives.md` |
| Create / author a policy (docs/policies/<name>.md)            | `/arcdlc:policy` skill, which applies `source/Policy of Policies.md`                                 |
| App methodology                                               | `source/Twelve-Factor App.md`                                                                         |
| CTO or operating model                                        | `source/CTO Methodology Guide.md`                                                                     |
| Leadership, Human-Behavior, Philosophie, Conflict-Solving     | `source/stoic.md`                                                                                     |
