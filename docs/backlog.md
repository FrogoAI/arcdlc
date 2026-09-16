# Backlog

Questions that were raised, deliberately deferred, and are still open. Each was accepted as
out of scope by the initiative that found it, and each needs its own interview before it becomes work.

Harvested on 2026-09-16 from the five initiatives retired that day. Git history holds the full
initiative folders; the ADR named beside an entry holds the reasoning.

One entry has already been closed: the deferred "add modern Go coverage" question was answered on
2026-09-16 by `skills/examinate/references/Modern Go.md`, 28 citable rules covering what changed in
Go after most published Go code was written.

## Plan and ordering

- **Should `/arcdlc:examinate` place gap tasks by dependency order instead of appending them?**
  Today it appends. Changing it alters a different skill's behaviour, so it needs its own interview.
- **Should a task block gain a `- DEPENDS:` key, with `arctool order --auto` deriving the order from
  it?** The command line designed for `arctool order` does not change if this is added later. See
  [ADR-0013](adr/0013-order-is-a-slot-permutation.md).

  **Raised again 2026-09-16, with a stronger reason.** Orchestrator mode spawns strictly one subagent at
  a time, because the queue is dependency-ordered and commits must not interleave. On a 150-task queue
  that serialises 150 spawns, and most of those tasks are almost certainly independent of each other.
  The ordering is implicit today, so nothing can tell which. `- DEPENDS:` is the precondition for ever
  running independent tasks in parallel, and a long queue is where that would pay.

- **Where should the executor tier pin live?** [ADR-0019](adr/0019-the-executor-tier-is-asked-for-not-guessed.md)
  put it in `CONTEXT.md`, which is global: one value for the whole project.
  [ADR-0021](adr/0021-a-planned-task-must-be-mechanical.md) then measured 55 real tasks and found tier
  fitness varies per initiative, not per project: `source-library-cleanup` had 30 repetitive tasks with a
  median `HOW` of 10 lines and suits a cheap tier, while `ordering` had 5 design-heavy ones with a median
  of 33 and suits the same model at low effort. Those two want different tiers and today must share one
  pin.

  The precedence ladder already has the finest grain (a task's `HOW`) and the coarsest (the global pin,
  then asking). The missing rung is per-initiative, which is exactly where the data says the answer
  changes. A marker in `plan.md`, like the source stamp, is the natural place. Changing the plan format is
  a contract change, so this needs a decision rather than a patch.
- **Should `arctool validate` warn on irregular block spacing?** Both `order` and `archive` normalise
  spacing silently today, which is consistent, so nothing is broken while this stays open.

## Install and distribution

- **Project-scope skills.** The installer writes personal skills only: `~/.cursor/skills/`,
  `~/.codex/skills/`, `~/.gemini/config/skills/`. Project-scope equivalents (`.cursor/skills/`,
  `.agents/skills/`) were deferred by both the Cursor and Antigravity initiatives.
- **Should Cursor skills set `disable-model-invocation`?** Left unset so the skills stay byte
  identical across agents, which the install-agnostic rule requires. Cursor may therefore invoke them
  from their descriptions as well as explicitly. See
  [ADR-0006](adr/0006-cursor-support-via-flat-personal-skills.md).
- **Should ArcDLC also distribute through Cursor rules or a project `AGENTS.md`?** Skills only so far.
- **Antigravity migration documentation.** A README note on `agy plugin import gemini` was deferred
  until the `agy` install path is verified on a live Antigravity install. See
  [ADR-0005](adr/0005-antigravity-support-via-plugin-with-flat-fallback.md).

## Closed, with the reasoning kept

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
