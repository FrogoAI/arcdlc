# CONTEXT — ArcDLC ubiquitous language

Glossary for humans and agents working on ArcDLC. Terms below are used in this exact sense across
skills, `arctool`, ADRs, and architecture documents.

- **Initiative** — one architecture-driven change effort, living entirely in `docs/aics/<slug>/`.
- **Slug** — the initiative's identifier: a single kebab-case path segment (no `/`, no `..`).
  Mandatory first positional argument of every pipeline skill; `--aic SLUG` in `arctool`.
- **Initiative folder** — `docs/aics/<slug>/`, holding the architecture document plus `plan.md`,
  `gap.md`, and `plan-archive.md` (the latter two are always siblings of `plan.md`).
- **Architecture document** — `aic.md`, `arc42.md`, `togaf.md`, or `c4.md` inside the initiative
  folder. Its first `# ` H1 is the initiative **title**; the one-line `> ` blockquote directly
  under the H1 is the initiative **summary** (both are contract, parsed by `arctool sync`).
- **Registry** — the generated list of initiatives between `<!-- arcdlc:initiatives:begin -->`
  and `<!-- arcdlc:initiatives:end -->` in `AGENTS.md` and `README.md`. Owned by `arctool sync`;
  never edit inside the markers by hand.
- **Sync** — `arctool sync [--check]`: regenerates the registry from `docs/aics/*/`; `--check`
  verifies drift without writing.
- **Removal** — `/arcdlc:remove <slug>`: engineer-confirmed deletion of an initiative folder plus
  registry cleanup. Git history is the archive; no graveyard copies in the tree.
- **Pipeline skills** — `aic`, `plan`, `execute`, `examinate`, `archive` (plus the lifecycle
  skills `remove` and `policy`).
- **Plan contract** — the task-block format defined in `skills/plan/references/plan-format.md`,
  parsed mechanically by `internal/plan`.
