# Backlog

Questions that were raised, deliberately deferred, and then decided. Each was accepted as out of
scope by the initiative that found it, and each was settled in a `/arcdlc:grilling` interview on
2026-09-16.

**Nothing here is open.** New questions go under a `## Open` heading above this line; everything below
is the record of what was decided and why, kept so a closed question is not re-opened on reasoning
that was already weighed. Git history holds the five retired initiative folders that raised most of
these; the ADR named beside an entry holds the fuller reasoning.

The first to close was the deferred "add modern Go coverage" question, answered by
`skills/examinate/references/Modern Go.md`: 28 citable rules covering what changed in Go after most
published Go code was written.

## Closed, with the reasoning kept

- **Should `arctool validate` warn on irregular block spacing?** **Decided 2026-09-16: no.** Spacing
  changes nothing the runner does: a plan with zero blank lines in one gap and three in another parses
  identically, and all tasks come back correctly. A warning would report something that is not a defect,
  and every warning that is not a defect trains people to skip warnings, including the two that matter,
  `source-changed` and `unverifiable-acceptance`.

  `order` and `archive` already normalise to one blank line whenever they rewrite, so a plan converges on
  tidy spacing by being worked rather than by being nagged.

  Residual, accepted: a hand-edited plan that is never reordered or archived keeps whatever spacing it was
  given, so a later `arctool` rewrite shows spacing churn mixed with the real change in review. If that
  becomes annoying, `arctool fmt` is the honest fix, because it gives you something to run rather than a
  message to ignore.

- **What happens to a key the plan format does not define?** **Decided 2026-09-16: custom keys are
  allowed, preserved, and passed to the executor.** The seven defined keys are the contract; anything
  else is part of the task and is carried through untouched. `validate` says nothing about them.

  This started as a version-skew question (an older `arctool` reading a newer plan) and testing turned
  it into a bug report. An unknown key placed after a multi-line key was **swallowed into it**, so
  `HOW` came back as `"Use the existing helper.\n- BUDGET: 3 retries..."`, corrupting the defined key
  as well as losing the custom one. With no multi-line key above it, the line vanished entirely.
  `validate --strict` exited 0 in both cases. Fixed in arctool 0.18.0: absorption now stops at any
  key-shaped line, custom keys are captured in document order with their indented bodies, and they
  appear under `extra` in `show` and `next --json`.

  The version-skew worry it came from is largely answered by this: an older tool meeting a newer key
  now carries it to the executor rather than eating it. What remains unhandled is a *changed* meaning
  for an existing key, which no marker would catch either.

- **Should `/arcdlc:examinate` place gap tasks by dependency order instead of appending them?**
  **Decided 2026-09-16: no, keep appending.** Appending is the only placement that is always safe. A
  gap task inserted mid-plan changes the run order of tasks already reasoned about, and `examinate`
  audits code rather than re-deriving the plan's dependency reasoning, so it has no basis for knowing
  its task belongs before `AIC-7` rather than after.

  [ADR-0024](adr/0024-the-plan-is-a-queue-not-a-graph.md) sharpens this beyond the reason it was
  deferred with: the plan is an ordered list a person reads, and `arctool order` is a deliberate act. A
  skill silently reordering someone's queue is the implicit-ordering behaviour that decision rejected.

  Residual, accepted: appended gap tasks run last, so a violation blocking work that a later task needs
  is found late. The engineer reorders with `arctool order` when that happens.

- **Project-scope skills, and distribution through Cursor rules or a project `AGENTS.md`.**
  **Decided 2026-09-16: no to both. Global install is the intended model.** `install.sh` writes only
  under `$HOME` (`~/.claude`, `~/.codex`, `~/.config/opencode`, `~/.cursor`, `~/.gemini`), and that is
  correct as it stands: you install ArcDLC once and use it in any repository or folder. Vendoring the
  624 KB bundle into each project would add copies that drift from each other and from upstream, and
  three of the four project-scope paths cannot be verified from here anyway.

  One real problem sits underneath this and is **not** solved by closing it: a team on different agents
  can be on different bundle versions while sharing one `plan.md`, and the plan format is a contract.
  [ADR-0022](adr/0022-a-plan-records-the-design-it-came-from.md)'s source stamp catches a drifting
  design, nothing catches a drifting bundle. The fix is pinning with `./install.sh --ref vX.Y.Z` plus
  recording the expected version in the consuming project.

  **Corrected 2026-09-16:** when that advice was given, `--ref` pinned only the skills. `arctool` was
  always taken from `releases/latest`, so a pinned install produced exactly the skew the pin was meant
  to prevent. Fixed in the same session: a ref beginning with `v` now takes the tool from that release,
  and the installer warns when the skills tree expects a different `arctool` than the one on disk.

- **Should Cursor skills set `disable-model-invocation`?** **Decided 2026-09-16: no, leave it unset.**
  [ADR-0006](adr/0006-cursor-support-via-flat-personal-skills.md) already judged auto-invocation
  desirable rather than merely tolerated, and this confirms it. Two further reasons found while
  deciding: the installer copies skills verbatim and never mutates them per agent, so setting the field
  means setting it everywhere including Claude Code; and six skills prefer to model-invoke
  `arcdlc-grilling` as a sibling, which the field would silently downgrade to the inline fallback.

  Residual, accepted: the skill descriptions were rewritten on 2026-09-16 to be trigger-rich, because
  `execute` and `grilling` could not be reached by natural language at all. Unwanted firing is therefore
  more likely than when ADR-0006 was written. If a skill starts pulling itself into unrelated
  conversations, the fix is a tighter description, not a global switch.

- **Antigravity migration documentation.** **Closed 2026-09-16: not written.** Antigravity is the one
  supported agent nobody here uses, and the `agy` install path has never been confirmed against a
  shipping build. Writing migration instructions for a path we cannot test would be inventing
  confidence. The limit is documented in `README.md` under Agent support.

- **Should a task block gain a `- DEPENDS:` key, and should independent tasks run in parallel?**
  **Decided 2026-09-16: no, to both.** A plan is an ordered list that runs top to bottom, one task at a
  time. A queue is resumable with one pointer (`arctool next`), a graph is not; order stays readable by a
  person rather than computed by a tool; commits stay linear and bisectable; and a concurrency failure
  cannot be expressed in a mechanical block. The cost is wall-clock time, and it is accepted: where
  throughput matters, run separate initiatives, which are already independent. Re-open only with a
  concrete failure that ordering cannot express, not with a wish for speed. See
  [ADR-0024](adr/0024-the-plan-is-a-queue-not-a-graph.md).

- **Where should the executor tier pin live?** **Decided 2026-09-16: `CONTEXT.md` is fine, keep it.
  Re-opened and amended 2026-09-18: the run pin stays in `CONTEXT.md`, and a task pins its own tier
  with an `Executor` key in `plan.md`.** The first decision reasoned that a task wanting a different
  tier could say so in its `HOW`, so a format change was too much machinery. That residual did not
  hold up: a tier in `HOW` is prose the dispatcher has to read for, `arctool next --json` could not
  show it, and the route invited a planner to write "run on opus" past a block that was not
  mechanical. The engineer asked for an explicit per-task override of the `CONTEXT.md` pin, for
  example `- Executor: opus, high effort.` on one migration and `- Executor: sonnet` on a rename. It is
  one optional single-line key, first in the resolution order, written only when the engineer asks for
  it, and `HOW` no longer names a tier. A per-initiative pin is still not added: an initiative whose
  every task wants one tier writes it on each task, or the engineer answers the question that run. See
  [ADR-0025](adr/0025-a-task-pins-its-executor-tier-with-an-executor-key.md).

- **Should rarely-used sections move out of a `SKILL.md` into `references/`?** Raised because
  `execute` is 20 KB and loads in full to run one task, and because ADR-0020 weakened the premise of
  the self-sufficiency rule: references now ship in the same directory as the skill, so the risk of a
  missing file is gone.

  Measured on 2026-09-16 and **declined**. The 70-line executor-tier section breaks down as 42 lines
  of binding rule, 16 lines of the main orchestrator path, 9 lines of the in-session path, and 3
  lines of genuinely rare branch. There is nothing to move. A binding rule cannot go behind a read
  the agent might skip, and the orchestrator loop is not a rare branch, it is what every whole-queue
  run does. Moving it would trade 5 KB for an agent improvising the queue loop when it does not open
  the file.

  The self-sufficiency rule earns its keep. Re-open only with a section that is genuinely rare and
  carries no binding rule.

## Multi-agent support is structural, not tested

Closed as an obligation on 2026-09-16 and rewritten as a documented limit. The `antigravity-cli`
initiative asked for a manual check against a real Antigravity install, and nobody has one. Chasing
that verification was the wrong shape of task: it singled out one agent when the same gap applies to
four.

The bundle is install-agnostic by construction, and CI proves the files land correctly, the installer
is idempotent, uninstall is clean, and a retired skill is swept. None of that proves a skill *behaves*
correctly on an agent, because proving that means running a real agent, which CI cannot do.

`README.md` now states this plainly in its Agent support section: developed and exercised on Claude
Code, install path covered by CI for Codex, OpenCode and Cursor, and only the flat fallback covered
for Antigravity. The plugin manifests say the same. Nothing is known broken; most of it is simply
unverified, and that is now written down rather than implied away.

Re-open per agent, with evidence, when someone actually runs the bundle on one.
