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

## Open obligations

- **Antigravity live verification, never recorded as done.** The `antigravity-cli` initiative required
  a manual check on a real Antigravity install: run the installer, confirm the `/arcdlc-*` commands
  register (through `agy` when present, else from `~/.gemini/config/skills/`), and confirm uninstall
  clears them. CI proves only the flat fallback, because it cannot run `agy`. Nothing records that
  this was ever performed, so treat Antigravity support as unverified against a shipping build.
