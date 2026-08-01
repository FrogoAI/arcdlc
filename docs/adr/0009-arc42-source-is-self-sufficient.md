# ADR-0009 — `Arc42.md` must generate a full arc42 document on its own

- Status: Accepted
- Date: 2026-08-01
- Initiative: [source-library-cleanup](../aics/source-library-cleanup/aic.md)

## Context

`/arcdlc:aic <slug> arc42` routes to two sources: `Arc42.md` (317 lines — numbered sections, worked
examples, an ISO-25010 table, a Go/monorepo mapping) and `arc42/arc42-template-EN.md` (987 lines —
the upstream template, all placeholders, no examples).

The bundled upstream template is the weaker of the two in an agent context:

- Its key content — the ISO-25010 quality categories and the building-block hierarchy — is delivered
  only as `<img src>` figures (lines 73–76, 319–322, 685–688), which an agent cannot read. `Arc42.md`
  supplies the same material as text.
- Its headings are unnumbered, so the cross-references "Arc42 §5 Level 1" and "Arc42 Section 6" in
  `C4.md:27` and `UML.md:30–31` resolve only against `Arc42.md`.
- Line 1 is an empty `# ` heading; the file has no title.
- It carries ~17 "Motivation" and 11 "Further Information" blocks explaining why each section exists.

The directory also costs 984 KB — six PNGs, three of them referenced by nothing — and `install.sh`
copies `skills/` wholesale, so every install on every supported agent pays for it.

## Decision

`Arc42.md` becomes **self-sufficient**: it must contain everything needed to generate a complete,
standard-conformant arc42 document — all twelve sections with their required content — without any
second source. Where the upstream template carries structure `Arc42.md` lacks, that structure is
adopted and copied into `Arc42.md`.

Once `Arc42.md` is self-sufficient, the bundled `arc42/` directory is removed and the `arc42` row in
`skills/aic/SKILL.md` and the "Architecture documentation" row in `skills/source-map/SKILL.md` point
at `Arc42.md` alone.

**The order is a hard constraint:** the rewrite lands first, the deletion second. The bundle must
never be in a state where the arc42 format has no complete source.

## Justification

- **One source, fully readable.** Generation quality is set by the source an agent can actually
  parse; text beats figures it cannot see.
- **Cross-references already point at `Arc42.md`.** `C4.md` and `UML.md` cite numbered arc42
  sections that only `Arc42.md` provides, so it is already the de facto canonical source.
- **984 KB in every install** buys nothing an agent can use — three of the six images are orphans,
  and the rest illustrate the file being retired.

## Trade-offs

- **Loss of verbatim upstream text.** The bundle stops shipping the official arc42 template as
  published. Accepted: fidelity is preserved by adopting its structure into `Arc42.md`, and the
  upstream document remains available at arc42.org for anyone who wants the original.
- **`Arc42.md` grows.** Absorbing the missing structure works against the initiative's trimming goal.
  Accepted deliberately: the growth is bounded by dropping the template's motivation blocks,
  placeholder prose, and pandoc HTML tables, and it removes a 987-line file in exchange.
- **Sequencing risk.** If the rewrite is under-scoped and the deletion still lands, arc42 generation
  degrades silently. Mitigation: the plan orders the deletion task after the rewrite task and its
  acceptance criteria name the twelve sections.
