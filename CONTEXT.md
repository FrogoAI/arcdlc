# CONTEXT — ArcDLC ubiquitous language

Glossary for humans and agents working on ArcDLC. Terms below are used in this exact sense across
skills, `arctool`, ADRs, and architecture documents.

- **Initiative** — one architecture-driven change effort, living entirely in `docs/aics/<slug>/`.
- **Slug** — the initiative's identifier: a single kebab-case path segment (no `/`, no `..`).
  Mandatory first positional argument of every pipeline skill; `--aic SLUG` in `arctool`.
- **Initiative folder** — `docs/aics/<slug>/`, holding the architecture document plus, while the
  initiative is open, `plan.md`, `gap.md`, `comments.md`, `plan-human.md`, and `plan-archive.md`
  (the last four are always siblings of `plan.md`); once closed, `CLOSED.md` replaces all five.
- **Architecture document** — `aic`, `arc42`, `togaf`, `c4`, or `tsc` inside the initiative folder,
  as `.md` or (when asked for with `:html`) `.html`. One initiative may hold several; `arctool sync`
  picks one by format rank, `.md` before `.html`. Its first `# ` H1 is the initiative **title**; the one-line `> ` blockquote directly
  under the H1 is the initiative **summary** (both are contract, parsed by `arctool sync`).
- **Registry** — the generated list of initiatives between `<!-- arcdlc:initiatives:begin -->`
  and `<!-- arcdlc:initiatives:end -->` in `AGENTS.md` and `README.md`. Owned by `arctool sync`;
  never edit inside the markers by hand.
- **Sync** — `arctool sync [--check]`: regenerates the registry from `docs/aics/*/`, leaving out any
  folder that holds `CLOSED.md` and counting the closed ones in one line; `--check` verifies drift
  without writing.
- **Close** — `/arcdlc:close <slug>`: marks a finished initiative by writing `CLOSED.md` into
  `docs/aics/<slug>/` and deleting `plan.md`, `plan-archive.md`, `plan-human.md`, `gap.md` and
  `comments.md`. The design never moves. A closed initiative is final: no skill changes it or
  deletes `CLOSED.md`. See [ADR-0028](docs/adr/0028-a-finished-initiative-is-closed-in-place.md).
- **Removal** — `/arcdlc:remove <slug>`: engineer-confirmed deletion of `docs/aics/<slug>/` plus
  registry cleanup. Git history keeps the deleted folder.
- **Closed initiative** — an initiative folder holding `CLOSED.md`: the design documents, their
  images, and the note with the close date, the **outcome** (`delivered`, `partly delivered` or
  `stopped`) and what was left open. Kept as a reference for later designs.
- **Phase** — the status of an initiative, read from its files by `arctool status`, never stored:
  `designing` (no `plan.md`, no `CLOSED.md`), `in progress` (a task not `DONE`), `ready to close`
  (every task `DONE`), `closed` (`CLOSED.md`).
- **Pipeline skills** — `init`, `aic`, `plan`, `execute`, `examinate`, `assist`, `archive`, plus the
  lifecycle skills `close`, `remove` and `policy`, `plan-human` beside the chain, and `grilling` under
  them all. Twelve in total.
- **Comment register** — `docs/aics/<slug>/comments.md`: one block per code comment marker found in
  the repository. `arctool scan` writes the evidence; it removes the comment from the code only with
  `--strip`, which `/arcdlc:assist` passes after the engineer says yes. Once stripped, the register is
  the only copy of that marker text. `/arcdlc:assist` writes the judgement, and every block whose
  verdict is `ACTIONABLE` is mirrored into `plan.md` as a task. The file is append-only: blocks are
  filled in, never rewritten or removed. See [ADR-0015](docs/adr/0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md)
  and [ADR-0017](docs/adr/0017-a-marker-is-arcdlc-in-a-comment-and-cutting-it-is-asked-for.md).
- **Marker** — a word that opens a code comment to mark unfinished work. The default is `ARCDLC`,
  the bundle's own word, so a sweep never touches the `TODO` and `FIXME` notes a repository already
  had; `--marker TODO` is how a team opts in. Only a single-line comment carries one: `// ARCDLC ...`
  counts, `/* ARCDLC ... */` does not, in every language. A word that merely appears inside a comment,
  or inside a string, is not a marker. See [ADR-0018](docs/adr/0018-only-single-line-comments-carry-markers.md).
- **Group tag** — `ARCDLC:T1` in three files is one record and one task, named after the tag
  (`ARCDLC-CMT-T1`). Tags fold to upper case and must contain a letter. See [ADR-0016](docs/adr/0016-comment-markers-group-by-tag.md).
- **Executor tier** — the model and the effort level a `/arcdlc:execute` task subagent runs at. Never
  below the capability floor: shell commands, file edits, the project's test and lint commands, a
  commit. `/arcdlc:execute` never picks it alone. It asks once per run, unless a task's `Executor`
  line or the pin below already names one. Defined in `skills/execute/SKILL.md`.
- **Executor pin** — the optional `Executor tier: <name>` line in this file. When it is there,
  `/arcdlc:execute` uses it exactly as written and asks nothing. It may name an effort level too, for
  example `Executor tier: opus, low effort`. It covers every task that has no task pin.
- **Task pin** — the optional `- Executor:` key on one task in `plan.md`, for example
  `- Executor: opus, high effort.` or `- Executor: sonnet`. It binds that task only, above the executor
  pin and above the question, and `/arcdlc:plan` writes it only when the engineer asked for that task
  to run there. `arctool next --json` returns it as `executor`. See
  [ADR-0025](docs/adr/0025-a-task-pins-its-executor-tier-with-an-executor-key.md).
- **Plan contract** — the task-block format defined in `skills/plan/references/plan-format.md`,
  parsed mechanically by `internal/plan`.
- **Slot permutation** — the semantics of `arctool order`. The named tasks keep the positions they
  already hold in `plan.md`, and only their contents are permuted among those positions. A task that
  is not named never moves. See [ADR-0013](docs/adr/0013-order-is-a-slot-permutation.md).
- **Reference document** — a bundled file in the `references/` folder of the skill that consumes it.
  It earns its place only by changing what an agent produces, so it is either a **generator template**
  (copied by `/arcdlc:aic` to produce an architecture document) or an **audit target** (a rule set
  `/arcdlc:examinate` checks code against, exposing identifiers a gap block can cite). A rule the agent
  must always follow is not a reference document: it goes inline in the `SKILL.md`, because a skill is
  loaded while a reference is only read if the agent opens it. See [ADR-0020](docs/adr/0020-reference-library-dissolved-into-the-skills.md).
- **The four virtues** — wisdom, courage, justice, temperance: how a skill decides, not how it writes.
  `## Judge by the four virtues` is verbatim in all twelve `SKILL.md` files and pinned by CI.
- **Workspace** — several git repositories checked out side by side under one directory, the
  **workspace root**, which is not itself a git repository. The agent, every `/arcdlc:*` command and
  every `arctool` call run from the root. See [ADR-0026](docs/adr/0026-a-workspace-keeps-its-docs-in-a-sibling-repository-named-docs.md).
- **Hub** — the git repository checked out as `docs/` directly under the workspace root. It holds
  `aics/`, `adr/`, `policies/`, `AGENTS.md`, `CLAUDE.md` (a symlink to `AGENTS.md`), `README.md` and
  `CONTEXT.md`, laid out exactly as a single repository's `docs/` folder plus those four files. The
  name `docs` is the whole configuration: `docs/aics/<slug>/` resolves from the root unchanged. The hub
  lives on its default branch only, with no branches and no pull requests, and is pushed after every
  write. Nothing is ever linked or copied into a product repository.
- **Root links** — the four symlinks at the workspace root, `AGENTS.md`, `CLAUDE.md`, `README.md` and
  `CONTEXT.md`, each pointing at the file of the same name in the hub. Untracked, created and checked
  by `make -C docs init` and `make -C docs check`, and the only thing the root holds.
- **Product repository** — any repository in the workspace other than the hub. It keeps its own
  `AGENTS.md` and its own branch and review flow; a run commits there and never pushes.
- **Repo key** — the custom plan key `- Repo: <name>` a task carries in a workspace, naming the one
  repository its `WHERE` files live in (`docs` for a hub-only task). Carried through by `arctool` under
  `extra` as the entry whose `key` is `Repo`; `WHERE` paths stay relative to the workspace root. Written by `/arcdlc:plan` and
  `/arcdlc:init`'s migration step, read first by `/arcdlc:execute`.
- **Layout** — what `/arcdlc:init` detects: **single repository** (the working directory is a git
  repository, its `docs/` is a folder) or **workspace** (the working directory is not a git repository
  and holds at least one). Every other skill is layout-blind by construction.

## Executor tier pin

`/arcdlc:execute` reads this file for one optional line and uses it for every task subagent it
spawns, except a task whose block carries its own `- Executor:` line, which wins for that task. The
line stands on its own, at the start of a line, and reads like `Executor tier: haiku` or
`Executor tier: opus, low effort`.

The pin is how you stop being asked. Either way the run names the tier it used in its report.
`arctool sync` rewrites `AGENTS.md` and `README.md` only, so a pin here survives every sync.

Executor tier: sonnet

Set on 2026-09-19 for the `init` initiative: the harness in use (Claude Code) exposes a model per
subagent and no effort dial, and every block carries its decisions, so the cheapest model that clears
the capability floor runs the queue.
