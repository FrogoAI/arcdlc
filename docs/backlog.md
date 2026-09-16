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
