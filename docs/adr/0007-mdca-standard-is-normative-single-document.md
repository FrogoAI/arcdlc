# ADR-0007 — The MDCA standard is normative; the reference library ships one MDCA document

- Status: Accepted
- Date: 2026-08-01
- Initiative: [source-library-cleanup](../aics/source-library-cleanup/aic.md)

## Context

`skills/source-map/source/` ships two MDCA documents, and `skills/examinate/SKILL.md` instructs the
agent to read **both in full** (806 lines) for any MDCA audit:

- `mdca.md` (337 lines) — prose, no rule identifiers, line 4 declares itself "the **canonical
  reference** for MDCA".
- `mdca_standard.md` (469 lines) — normative: `P1`–`P10` (§6), §7 layout, §8 rules, §11 compliance
  checklist, §12 anti-patterns, and Appendix A's import allow/deny matrix. Its header reads
  `Version: 1.0`, `Status: Draft`.

Roughly 70% of the content is duplicated at different levels of rigor, neither file declares
precedence over the other, and they give conflicting answers on three mechanically checkable rules:

1. **`domain → pkg/`** — `mdca.md:133` forbids and permits it in the same sentence; §7.1 says
   "standard library only"; Appendix A allows "pure utilities (e.g. `pkg/clock`)".
2. **Transaction ownership** — §8.6 says application services **MUST** be the only place a
   transaction is opened; `mdca.md:282` also allows adapters; the standard's own §9.4 example opens
   one inside `domain/order/service.go`.
3. **Canonical folder layout** — §7.2 has `events.go` and no `mocks/`; `mdca.md:249–263` has `mocks/`
   and no `events.go`. Both are labelled canonical.

An audit therefore ingests three unresolved conflicts and silently picks a winner.

## Decision

**`mdca_standard.md` is normative.** Where the two documents disagree, the standard wins, and
`mdca.md` is aligned to it — not the reverse.

The library ships **one** MDCA document, built on the standard's spine (§6 `P1`–`P10`, §7, §8, §11,
§12, Appendix A), absorbing only what is unique to `mdca.md`: the Server/Client/Library
specialization table, the MDCA–SOLID–Clean Code alignment table, and the five anti-pattern rows the
standard lacks.

The merged file keeps the filename **`mdca.md`**, because `ddd.md` cross-references that name in four
places while `mdca_standard.md` is referenced only by two skill files — renaming would ripple further
than merging.

The three conflicts are resolved as:

1. `domain → pkg/` — restated as a purity rule; see [ADR-0008](0008-domain-imports-judged-by-purity-not-folder.md).
2. **Transactions** — §8.6 stands verbatim: application services are the only place a transaction is
   opened. "Domain service" means no DB logic in the domain layer — no SQL, no low-level database
   manipulation. The §9.4 outbox example violates this and is rewritten.
3. **Layout** — §7.2 verbatim. `mocks/` is dropped from the guidance.

## Justification

- **`/arcdlc:examinate` needs one answer per rule.** A gap block cites a rule; a rule with three
  documented answers cannot be cited, and an audit that picks silently is not reproducible.
- **The standard is the only auditable artifact.** It already carries per-rule identifiers and an
  import matrix; `mdca.md` states the same content as prose that nothing can reference.
- **Halving the read cost.** One merged document of roughly 280 lines replaces 806, saving ~8–9k
  tokens on every MDCA audit — the single largest saving available in the library.

## Trade-offs

- **The prose on-ramp is lost.** `mdca.md`'s narrative introduction disappears; the merged document
  is normative in tone throughout. Accepted: the audience is an agent applying rules, not a reader
  being persuaded.
- **`Status: Draft` becomes misleading** on a document that is now the sole normative source. The
  status line is left as an open question for the initiative rather than silently promoted.
- **Cross-references to `mdca_standard.md`** in `skills/examinate/SKILL.md` and
  `skills/source-map/SKILL.md` must be updated in the same change set, or an audit will read a path
  that no longer exists.
