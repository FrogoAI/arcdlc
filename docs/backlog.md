# Backlog

Questions that were raised, deliberately deferred, and are still open. Each was accepted as
out of scope by the initiative that found it, and each needs its own interview before it becomes work.

Harvested on 2026-09-16 from the five initiatives retired that day. Git history holds the full
initiative folders; the ADR named beside an entry holds the reasoning.

One entry has already been closed: the deferred "add modern Go coverage" question was answered on
2026-09-16 by `skills/examinate/references/Modern Go.md`, 28 citable rules covering what changed in
Go after most published Go code was written.

## Plan and ordering

- **Should `arctool validate` warn on irregular block spacing?** Both `order` and `archive` normalise
  spacing silently today, which is consistent, so nothing is broken while this stays open.

## Closed, with the reasoning kept

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
  design, nothing catches a drifting bundle, and an engineer on an older `arctool` can parse a plan
  written by a newer one and silently miss a key. The cheap fix is pinning, which the installer already
  supports (`./install.sh --ref v0.27.0`) plus recording the expected version in the consuming project.
  Filed below as its own question rather than left implied.


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

- **Where should the executor tier pin live?** **Decided 2026-09-16: `CONTEXT.md` is fine, keep it.**
  ADR-0021 showed tier fitness varies per initiative, which argued for a per-initiative pin, but the
  precedence ladder already covers it without a format change: a run with no pin asks once, so an
  initiative that wants a different tier gets answered differently that run, and a task that wants one
  specifically says so in its `HOW`. Residual, accepted: a global pin is used silently, so an initiative
  needing a different tier relies on the engineer editing the one line or overriding in `HOW`. A plan
  format change was judged too much machinery for that.

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
