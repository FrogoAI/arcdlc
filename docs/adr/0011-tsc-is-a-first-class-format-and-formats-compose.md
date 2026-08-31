# ADR-0011 — Tech Stack Canvas is a first-class `/arcdlc:aic` format, and format arguments compose

- Status: Accepted
- Date: 2026-08-31
- Initiative: [source-library-cleanup](../aics/source-library-cleanup/aic.md)

## Context

`Tech Stack Canvas.md` has been in the reference library since the beginning, routed from the
`source-map` table under "Tech stack decisions", but it was reachable from `/arcdlc:aic` only through
the table's `anything else` fallback row. Nothing named it, nothing said where its output goes, and
nothing told the agent it must carry the H1 + `> ` summary the registry parses.

It also answers a different question than the other formats. AIC, arc42, TOGAF, and C4 describe the
architecture; TSC describes the **technology picture** — services, stack, integrations,
infrastructure. Teams commonly want both from one design conversation, and re-interviewing per format
would produce two documents built from two different recollections of the same decisions.

## Decision

- `tsc` becomes a **named format**: `/arcdlc:aic <slug> tsc` → `docs/aics/<slug>/tsc.md`, generated
  from `Tech Stack Canvas.md`.
- `Tech Stack Canvas.md` gains a short **how to fill it** preamble (keep the four groups and every
  heading with its colour marker; `Not applicable — <reason>` for a field that does not apply; open
  with `# <Title>` and a one-line `> ` summary). Its field definitions are untouched.
- **The format argument accepts a comma-separated list**: `/arcdlc:aic payments arc42,tsc`. Every
  named format is produced from **one** interview, in the order given, each document with its own H1
  and summary. Where two formats cover the same ground, it is stated in full in one and cross-linked
  from the other — never forked into two copies that can diverge.
- An unrecognised entry anywhere in the list writes **nothing**: the skill reports the unknown name
  and lists the supported formats.
- `tsc.md` joins `archDocPrecedence` in `internal/registry` **last**
  (`aic.md`, `arc42.md`, `togaf.md`, `c4.md`, `tsc.md`), so a folder holding both `arc42.md` and
  `tsc.md` still registers under the arc42 document.

## Justification

- **One interview, many lenses.** The expensive part of `/arcdlc:aic` is the grilled interview, not
  the writing. Running it once and projecting the answers into several documents is strictly cheaper
  and strictly more consistent than running it twice.
- **Precedence last is the safe position.** TSC is a technology view, not the initiative's
  architectural narrative, so it should never outrank arc42 or TOGAF in the registry. Placing it in
  the list at all still beats the alphabetical fallback, which would pick `notes.md` over `tsc.md`.
- **Explicit beats fallback.** A named row states the output path and the source; the `anything else`
  row states neither, so two agents could reasonably write two different files.

## Trade-offs

- **Cross-document consistency is a rule, not a mechanism.** Nothing validates that `arc42.md` and
  `tsc.md` agree. Accepted: the skill mandates single-statement-plus-cross-link, and `/arcdlc:plan`
  reads one named document, so a divergence cannot silently reach the plan.
- **A longer format argument to parse.** `arc42,tsc` is one more shape a skill must handle correctly,
  and a typo inside a list now fails the whole invocation rather than one document. Accepted: failing
  closed and naming the bad entry is better than writing a partial set.
- **`Tech Stack Canvas.md` is no longer purely upstream.** The added preamble is this bundle's own.
  Accepted: the canvas's own field definitions are unchanged, and the preamble is what makes the
  output parseable by `arctool sync`.
