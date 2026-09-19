# Plan: init

<!-- arcdlc:source docs/aics/init/aic.md sha256:84ca0d7783e0133610230370b708b8d6da9834f69068f41693932a4822c6b973 -->

Format: [`skills/plan/references/plan-format.md`](../../../skills/plan/references/plan-format.md).
Source: [`aic.md`](aic.md), decided 2026-09-19. Terms: [`CONTEXT.md`](../../../CONTEXT.md).
Decision: [ADR-0026](../../adr/0026-a-workspace-keeps-its-docs-in-a-sibling-repository-named-docs.md).

## Risk Coverage

- R1 (a stray write replaces a root link): INIT-1 fixes the one arcdlc writer; INIT-7 ships `make -C docs check` and tells the hub guide to run it when the registry looks wrong.
- R2 (a new skill fails the CI pins or never installs): INIT-6 adds `init` to `bundle.Skills` and `SUBSKILLS` in the same commit as the skill file; INIT-9 updates the four CI loops and the manifests.
- R3 (bundle version skew across machines sharing a hub): INIT-7 records the minimum `arctool` version in the hub guide template and makes `init` check `arctool version` before writing workspace rules.
- R4 (a push retry loop): INIT-3, one fast-forward and one retry, then `BLOCKED` with the git error.
- R5 (migration on a dirty repository): INIT-8 refuses to start on any repository whose `git status --porcelain` is not empty.
- R6 (a moved plan with paths that cannot be found): INIT-8 reports every unprefixed path per task and names `/arcdlc:plan <slug>` as the next step for that plan.
- R7 (a harness with no guide because a link is missing): INIT-7, four links, `check`, and a first line in the hub `README.md` template for a reader at the root.
- OQ1 (a root that is itself a repository holding repositories): INIT-6, `init` names the case and stops. No layout is scaffolded for it.
- OQ2 (which branch the migration commit lands on when a repository is not on its default branch): INIT-8, `init` asks per repository, one question per turn.
- OQ3 (converting MetricAid's `Repo:` lines inside `WHERE`): accepted, not covered. The migration step is for per-repository folders, not for an existing hub; that rewrite is MetricAid's call.
- OQ4 (`sync` run inside the hub writes `_none_`): accepted, not covered. A separate small `arctool` guard, out of this initiative's scope by the AIC.
- Routing (a skill description nobody can test in CI): process mitigation. INIT-9 adds the `init` rows to `docs/routing-checks.md`; the engineer re-runs that file by hand after the release, as `AGENTS.md` requires for every description change.

Completed (archived to docs/aics/init/plan-archive.md):
- INIT-1: Make `arctool sync` write through a symlink and base links on the real file
- INIT-2: Document the `Repo` custom key in the plan format and make `/arcdlc:plan` write it in a workspace
- INIT-3: Add the workspace per-task contract to `/arcdlc:execute`
- INIT-4: Add the hub-write paragraph to `aic`, `plan`, `examinate` and `assist`
- INIT-5: Add the hub-write paragraph to `archive`, `remove`, `policy` and `plan-human`
- INIT-6: Create the `init` skill with layout detection, the interview and the single-repository scaffold
- INIT-7: Write the workspace scaffold step of the `init` skill with the hub templates and Makefile
- INIT-8: Write the migration step of the `init` skill
- INIT-9: Register the eleventh skill in CI, the manifests, the counts and the routing checks
- INIT-10: Document the two layouts in the README and the bundle guide
