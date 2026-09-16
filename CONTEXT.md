# CONTEXT — ArcDLC ubiquitous language

Glossary for humans and agents working on ArcDLC. Terms below are used in this exact sense across
skills, `arctool`, ADRs, and architecture documents.

- **Initiative** — one architecture-driven change effort, living entirely in `docs/aics/<slug>/`.
- **Slug** — the initiative's identifier: a single kebab-case path segment (no `/`, no `..`).
  Mandatory first positional argument of every pipeline skill; `--aic SLUG` in `arctool`.
- **Initiative folder** — `docs/aics/<slug>/`, holding the architecture document plus `plan.md`,
  `gap.md`, `comments.md`, and `plan-archive.md` (the last three are always siblings of `plan.md`).
- **Architecture document** — `aic`, `arc42`, `togaf`, `c4`, or `tsc` inside the initiative folder,
  as `.md` or (when asked for with `:html`) `.html`. One initiative may hold several; `arctool sync`
  picks one by format rank, `.md` before `.html`. Its first `# ` H1 is the initiative **title**; the one-line `> ` blockquote directly
  under the H1 is the initiative **summary** (both are contract, parsed by `arctool sync`).
- **Registry** — the generated list of initiatives between `<!-- arcdlc:initiatives:begin -->`
  and `<!-- arcdlc:initiatives:end -->` in `AGENTS.md` and `README.md`. Owned by `arctool sync`;
  never edit inside the markers by hand.
- **Sync** — `arctool sync [--check]`: regenerates the registry from `docs/aics/*/`; `--check`
  verifies drift without writing.
- **Removal** — `/arcdlc:remove <slug>`: engineer-confirmed deletion of an initiative folder plus
  registry cleanup. Git history is the archive; no graveyard copies in the tree.
- **Pipeline skills** — `aic`, `plan`, `execute`, `examinate`, `assist`, `archive`, plus the
  lifecycle skills `remove` and `policy`, `plan-human` beside the chain, and `grilling` under them all.
  Ten in total.
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
  commit. `/arcdlc:execute` never picks it alone. It asks once per run, unless a task's `HOW` or the
  pin below already names one. Defined in `skills/execute/SKILL.md`.
- **Executor pin** — the optional `Executor tier: <name>` line in this file. When it is there,
  `/arcdlc:execute` uses it exactly as written and asks nothing. It may name an effort level too, for
  example `Executor tier: opus, low effort`.
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
  `## Judge by the four virtues` is verbatim in all ten `SKILL.md` files and pinned by CI.

## Executor tier pin

`/arcdlc:execute` reads this file for one optional line and uses it for every task subagent it
spawns. The line stands on its own, at the start of a line, and reads like `Executor tier: haiku` or
`Executor tier: opus, low effort`.

The pin is how you stop being asked. No pin is set for this repository, so every whole-queue run asks
one question before its first spawn: which model and which effort level should run the queue. Answer
it there, or write the pin here and the question goes away. Either way the run names the tier it used
in its report. `arctool sync` rewrites `AGENTS.md` and `README.md` only, so a pin here survives every
sync.
