# ADR-0034: A review records findings and leaves the plan alone

- Status: Accepted
- Date: 2026-11-25

## Context

`/arcdlc:examinate` audits the whole codebase and mirrors every gap into `plan.md` in the same run. An engineer who wants a reviewer's read of one branch before a pull request gets the wrong thing twice: findings on code the branch never touched, and tasks in the queue before anyone decided the findings are worth fixing.

## Decision

1. A new skill, `/arcdlc:review <slug> [standards] [--base <ref>|--range <a>..<b>] [--repo <name>]`, reviews one change set: the diff from `git merge-base HEAD <default branch>` to the working tree plus untracked files, or the exact `--range`. A problem in untouched code is reported in one line and never filed.
2. The standards argument is a comma-separated list of standard names and policy paths, resolved exactly as `/arcdlc:examinate` Step 1 resolves one. With none named, the review uses the initiative's architecture document, ADRs and `CONTEXT.md`. In every case it also applies `docs/policies/*.md` and `CONVENTION.md` when they exist, and checks correctness, security and tests, because a review that passes a bug has not reviewed.
3. Findings are `REV-NN` blocks in the gap-block format, with the plan contract's `MISSING|PARTIAL|DRIFT` tag and a register-only `- Severity: BLOCKER|MAJOR|MINOR.` line. Each run appends a `## Review` section with a verdict, `CHANGES REQUESTED` or `APPROVED`, to `docs/aics/<slug>/review.md`, and shows the same in chat. IDs continue across runs. A re-run marks fixed findings `- Resolved:`.
4. The review never writes `plan.md`. `/arcdlc:plan <slug>` mirrors every finding that is neither `Accepted` nor `Resolved`, drops the register-only lines, and builds a plan from `review.md` alone when the folder has no architecture document.
5. `/arcdlc:examinate` reads `CONVENTION.md` too, so both skills apply the same project rules.

## Consequences

- The bundle has 14 skills. One more description sits in every session's context, and "review" now routes two ways: `docs/routing-checks.md` records the split.
- `review.md` is a third register beside `gap.md` and `comments.md`. `arctool sync` never takes it as the architecture document, and `/arcdlc:close` deletes it with the other registers.
- The heading tag stays inside the plan contract, so `internal/plan` needed no change. Severity lives in a line the plan copy drops.
- A finding sits in `review.md` until the engineer runs `/arcdlc:plan`. That gap between review and plan is deliberate: the engineer decides which review becomes work.
