# ADR-0033 — A derived initiative keeps a frozen copy of its source

- Status: Accepted
- Date: 2026-11-23

## Context

At FrogoAI an architect writes one AIC for a whole initiative, kept in the company knowledge base. One AIC can hold 3 stages, several teams, and 7 epics. An engineer is assigned one stage or one epic and runs arcdlc in a team workspace.

Running `/arcdlc:aic`, `/arcdlc:plan`, and `/arcdlc:examinate` on the whole architect AIC fails in two ways. `/arcdlc:plan` plans all 7 epics, so the engineer has to cut other people's work out by hand. `/arcdlc:examinate` reports the other 6 epics as gaps. Writing a fresh AIC with `/arcdlc:aic` fails the other way: nothing ties it to the architect's decisions, so it can quietly redesign a contract the architect fixed.

A `--from` flag on `/arcdlc:aic` was considered and dropped. A forgotten flag fails silently: the skill designs from scratch and never sees the source.

## Decision

1. A new skill, `/arcdlc:derive <slug> <source> <scope>`, starts an engineer's part of someone else's design. All three arguments are required on the first run, so a missing source or scope stops the skill before it writes anything.
2. It writes `docs/aics/<slug>/source-aic.md`: the source text unchanged, under a header with the source path or URL, the copy date, the sha256 of the text, and the scope. Only a run of `/arcdlc:derive` writes it.
3. It grills and writes a scoped `docs/aics/<slug>/aic.md` from the AIC Template. A decision the source made is an input that is cited, not asked. The interview covers what the source leaves to the engineer: the Engineer decides and Starting value items, the repositories, the failure cases, test data, and migration steps of this scope.
4. No skill changes a source decision. An answer that would goes to Open questions with the architect as owner. `/arcdlc:examinate` with no standard audits against both `aic.md` and `source-aic.md`, limited to the scope, and files a conflict between them as an `SRC-GAP` whose task is `BLOCKED` until the architect answers.
5. A re-run with the same source compares hashes. A changed source is diffed against the snapshot, each dependent decision is asked about, and the snapshot is replaced.
6. `arctool sync` never takes `source-aic.md` as the architecture document.

## Consequences

- The bundle has 13 skills. One more description sits in every session's context.
- The source stays the source of truth in its own home. The snapshot makes every session read the same text, and a change reaches the plan only through a deliberate refresh.
- Fetching a source from a URL depends on the tools the agent has. When it can't, the engineer exports the page to a file. A direct knowledge base connector is future work.
