# ADR-0026 — A workspace keeps its docs in a sibling repository named `docs`

- Status: Accepted
- Date: 2026-09-19
- Initiative: [init](../aics/init/aic.md)
- Amends: nothing. ADR-0001 (initiatives live in `docs/aics/<slug>/`, selection is explicit) is
  unchanged; this record says where that `docs/` folder may come from.

## Context

Every skill and the `arctool` binary address an initiative as `docs/aics/<slug>/`, relative to the
directory the agent runs in (`aicsDir` in `cmd/arctool/main.go`). Until now that assumed one
repository: the code and its `docs/` folder in the same checkout.

Real work is often spread over several repositories that ship together. FrogoDB is four (`fdb-server`,
`fdb-client`, `fdb-tools`, `fdb-operator`); MetricAid is thirty-nine. One initiative there touches two
to four of them. Putting the initiative in "the repository that owns the goal" was tried on FrogoDB
and failed in a known way: the workspace root's `AGENTS.md` and `CONTEXT.md` are untracked files that
nobody can review or commit, the initiative registry is split across repositories, and the root map
went stale (`fdb-operator` is missing from it).

MetricAid has run the other layout since 2026-09-02, on thirty initiatives: the workspace root is not
a git repository, a sibling repository named `docs` holds `aics/`, `adr/`, `CONTEXT.md` and the agent
guide, and the root carries symlinks into it. Because that repository is named `docs`, the path
`docs/aics/<slug>/` resolves from the root without any setting. Nothing in the skills or in
`arctool` had to change for the paths.

The choice was whether to lock that name or to add a setting (an environment variable or a marker
file at the root) that points `docs/aics/` elsewhere.

## Decision

**A workspace's documentation hub is a git repository checked out as `docs/` directly under the
workspace root, and `docs` is the only name it may have.** No environment variable, no marker file,
no flag. `/arcdlc:init` clones it under that name whatever the remote is called
(`git clone <url> docs`), and a hub checked out under another name is reported as the mistake it is.

**The hub's tree is the tree a single repository has under `docs/`, and nothing more.** `aics/`,
`adr/`, `policies/`, reference documents beside them, plus the files the root links to: `AGENTS.md`,
`CLAUDE.md` (a symlink to `AGENTS.md`), `README.md`, `CONTEXT.md`. There is no folder per
repository inside the hub. A document names the repository it describes; it is not filed under one.

**The agent runs from the workspace root**, where `AGENTS.md`, `CLAUDE.md`, `README.md` and
`CONTEXT.md` are symlinks into `docs/`. Every `/arcdlc:*` command and every `arctool` call runs there
too. Nothing is ever linked or copied into a product repository; each keeps its own guide and knows
nothing about the hub.

## Justification

- **Zero configuration is the feature.** The layout works because the directory name is the path
  every skill already hard-codes. A setting would have to be read by ten skills and one binary, and a
  missing setting fails silently: the tool writes into a fresh `docs/` at the root and the engineer
  finds two registries a week later.
- **One fact to teach.** "Clone the hub as `docs`" is the whole setup. `make -C docs init` creates the
  links, and a new machine is ready.
- **The single-repository workflow is unchanged.** A repository with its own `docs/` folder, and a
  workspace with a `docs/` repository, look identical to every skill. That is what lets one skill
  bundle serve both without a mode switch inside each command.
- **Proven, not designed.** MetricAid's hub has carried thirty initiatives across thirty-nine
  repositories with this exact rule, and the defects found there (a `sync` that overwrote the root
  links, two missing links) are fixed by this initiative rather than by a rename.

## Trade-offs

- **A repository cannot be named `docs` for another reason.** A workspace whose product repository is
  itself called `docs` has to clone the hub under that name and the product under another. Accepted:
  the case has not come up, and a product repository can be checked out under any name.
- **The hub must sit directly under the root.** A hub two levels down, or a root that is itself a git
  repository with submodules, is not this layout and `init` does not scaffold it. Accepted: the
  layout's only requirement is one directory level, and a nested case can be modelled as a workspace
  one level lower.
- **Symlinks at the root are untracked.** The root is not a repository, so the links live only on the
  machine that made them and can be lost by a stray write. Accepted, with `make -C docs check` to
  detect it and `make -C docs init` to repair it. The links are the only thing the root holds.
