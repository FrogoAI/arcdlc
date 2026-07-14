Plan format: see `skills/plan/references/plan-format.md`.

## Risk Coverage

Reconciled against `docs/aics/cursor-support/aic.md` → Technical Challenges & Risks and Open questions.

- **`shellcheck` regressions** → covered by CUR-1 (extend the existing `codex|opencode` branch rather than add novel shell; CI `shellcheck install.sh` is the gate).
- **Volatile Cursor skills directory** → covered by CUR-1 (`~/.cursor/skills/` reached only via a single `cursor_dir` variable, mirroring the other agent dir vars).
- **Model-invocation surprise (auto-trigger)** → accepted; documented as expected behavior in CUR-4 and ADR-0006. No code change (keeping skills byte-identical).
- **Broad `~/.cursor` detection (IDE-only users)** → accepted; skills are inert until invoked. Noted in CUR-4 docs.
- **Open question — project-scope skills (`.cursor/skills/`)** → accepted/deferred; personal-only this initiative.
- **Open question — distribute via Cursor rules / `AGENTS.md`** → accepted/deferred; skills-only this initiative.
- **Open question — set `disable-model-invocation`** → accepted/deferred; keep skills byte-identical (honors the install-agnostic rule).

## Tasks

Completed (archived to docs/aics/cursor-support/plan-archive.md):
- CUR-1: Add the `cursor` agent to install.sh
- CUR-2: Cover the Cursor install path in CI
- CUR-3: Bump the bundle version to 0.5.0
- CUR-4: Name Cursor as a supported agent in the docs
