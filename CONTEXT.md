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
- **Pipeline skills** — `aic`, `plan`, `execute`, `examinate`, `assist`, `archive` (plus the
  lifecycle skills `remove` and `policy`, and `plan-human` beside the chain).
- **Comment register** — `docs/aics/<slug>/comments.md`: one block per code comment marker (`TODO`,
  `FIXME` and friends) found in the repository. `arctool scan` writes the evidence, `/arcdlc:assist`
  writes the judgement, and every block whose verdict is `ACTIONABLE` is mirrored into `plan.md` as a
  task. See [ADR-0015](docs/adr/0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md).
- **Marker** — a word that opens a code comment to mark unfinished work: `TODO` by default, also
  `FIXME`, `HACK`, `XXX`, `BUG`. A word that merely appears inside a comment is not a marker.
- **Plan contract** — the task-block format defined in `skills/plan/references/plan-format.md`,
  parsed mechanically by `internal/plan`.
- **Slot permutation** — the semantics of `arctool order`. The named tasks keep the positions they
  already hold in `plan.md`, and only their contents are permuted among those positions. A task that
  is not named never moves. See [ADR-0013](docs/adr/0013-order-is-a-slot-permutation.md).
- **Reference library** — `skills/source-map/source/`: the bundled reference documents, reached only
  through the routing table in `skills/source-map/SKILL.md` (one row per document). Agent-facing: a
  document earns its place by changing what an agent produces.
- **Source document** — one file in the reference library. Either a **generator template** (consumed
  by `/arcdlc:aic` to produce an architecture document) or an **audit target** (a rule set
  `/arcdlc:examinate` checks code against); an audit target must expose rules a gap block can cite by
  identifier.
