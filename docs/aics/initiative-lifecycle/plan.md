# initiative-lifecycle — Plan

Task format contract: `skills/plan/references/plan-format.md`. Execute with
`/arcdlc:execute initiative-lifecycle` (one task, one commit, top-to-bottom).

Sequencing note: the arctool 0.6.0 per-folder work is uncommitted in this working tree and this plan
builds on top of it. Land 0.6.0 first, then work this queue. Tasks IL-1..IL-4 change `arctool`;
IL-5 updates the contract; IL-6..IL-9 update the skills; IL-10 adds the plan risk-coverage gate;
IL-11 updates docs and bumps versions.

## Risk Coverage

Every risk in the AIC's "Technical Challenges & Risks" and "Open questions" is covered by a task or
explicitly accepted (per ADR-0004; this plan dogfoods the gate IL-10 introduces):

- **In-flight collision** — accepted (process): the sequencing note above; land 0.6.0 first, don't interleave.
- **Editing user-owned root files** — covered by IL-3 (byte-preserving splice + tests).
- **Summary parsing on arbitrary agent-written docs** — covered by IL-2 (table-driven fallback tests, incl. arc42/TOGAF).
- **Breaking invocation change** — covered by IL-5 (contract), IL-6/IL-7 (skills), IL-11 (all docs slug-first).
- **Slug hygiene** — covered by IL-1 (`validSlug` retained and reused by `sync`/`remove`).
- **Open question — `arcdlc:policies` block later** — accepted/deferred (AIC H6): out of scope for this initiative.
- **Open question — CI runs `arctool sync --check` here** — covered by IL-11 acceptance (`sync --check` exits 0); wiring it as a CI gate is deferred.

Completed (archived to docs/aics/initiative-lifecycle/plan-archive.md):
- IL-1: arctool — make initiative selection mandatory (remove auto-detect)
- IL-2: arctool — initiative title/summary parsing (pure functions)
- IL-3: arctool — marker-block splice + stub creation (byte-preserving)
- IL-4: arctool — wire the `sync [--check]` subcommand
- IL-5: plan-format contract — mandatory selection + title/summary + registry
- IL-6: aic skill — slug-first mandatory + summary blockquote + post-create sync
- IL-7: plan/execute/examinate/archive skills — slug-first mandatory
- IL-8: policy skill — require the policy name argument
- IL-9: /arcdlc:remove skill + install/CI registration
- IL-10: plan skill — risk-mitigation coverage gate
- IL-11: docs + registry blocks + version bumps
