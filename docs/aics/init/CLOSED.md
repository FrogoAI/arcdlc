# ArcDLC Init: one command that sets up a repository or a multi-repository workspace

- Closed: 2026-09-25
- Outcome: delivered
- Tasks: 10 (10 done, 0 todo, 0 blocked)

## Left open

- **OQ1. A root that is itself a git repository containing other repositories.** Submodules, or a
  monorepo with vendored checkouts. `init` reports the case and stops. Whether it should treat the
  outer repository as a single repository or as a workspace has no answer yet; wait for a real case.
- **OQ2. Which product-repository branch shape the migration commit uses.** H8 commits once per
  repository on "its default branch or the branch its guide names". For FrogoDB two repositories sit
  on a feature branch today (`incident-0`). The migration step asks per repository when the checkout
  is not on the default branch.
- **OQ3. Whether MetricAid's hand-written `Repo:` first lines in `WHERE` are converted.** Thirty
  plans carry the MetricAid form. The migration step is for per-repository folders, not for an
  existing hub. A one-off rewrite there is MetricAid's call, not this initiative's.
- **OQ4. Whether `arctool sync` should also be run from inside the hub.** With H7, `sync` from the
  root writes through the links. `sync` run inside `docs/` finds no `docs/aics/` and writes `_none_`,
  as today. Guarding that (refuse to write an empty registry over a non-empty one without a flag) is
  a separate, small `arctool` change and is not decided here.

## Documents

- [aic.md](aic.md)
