# ADR-0010 — arc42 ships as the plain upstream skeleton plus one guide that carries all the help

- Status: Accepted
- Date: 2026-08-31
- Initiative: [source-library-cleanup](../aics/source-library-cleanup/aic.md)
- Supersedes: [ADR-0009](0009-arc42-source-is-self-sufficient.md)

## Context

[ADR-0009](0009-arc42-source-is-self-sufficient.md) made a bundle-authored `Arc42.md` self-sufficient
and retired the upstream template, on the reasoning that a single agent-readable source beats a
template whose key content is locked in images. Upstream material has since been re-added, and both
upstream editions were read side by side against `Arc42.md`:

| | `Arc42.md` (363 lines) | annotated edition (986 lines) | plain edition (237 lines) |
|---|---|---|---|
| Numbering | `1.`…`12.` with subsections | unnumbered, empty `# ` on line 1 | unnumbered, empty `# ` on line 1 |
| Fill-in slots | absent | present | present |
| Teaching apparatus | none | 65 bold labels (17 Motivation, 17 Form, 18 Contents/Content, 11 Further Information, 2 Structure/Examples) | none |
| Worked examples | 9 filled tables | 0 | 0 |
| Key content as images | none | 3 figures (ISO 25010 §1.2, building-block hierarchy §5, concept topics §8) | none |

Three findings drove the decision:

1. The two upstream editions have an **identical heading tree**. The annotated one adds only help
   text, which upstream itself says is for familiarization — the plain edition is what you document a
   real system in.
2. The bundle briefly shipped all three files. That cost ~1540 lines of reading per
   `/arcdlc:aic <slug> arc42` and left the agent to strip 65 teaching blocks at generation time — with
   a live trap, since §5 and §7 also contain *unbolded* `Motivation` / `Contained Building Blocks` /
   `Important Interfaces` lines that are content slots, not teaching.
3. Refreshing to a new arc42 release must stay a copy-paste. That is only true for a file kept
   verbatim.

## Decision

Ship exactly **two** files, and no third source:

- **`arc42-template-EN.md` — the upstream *plain* edition, verbatim.** The skeleton the output
  document is copied from: the twelve sections in order with every fill-in slot, no help. Refreshing
  to a newer arc42 release means pasting the new **plain** edition over this file and re-checking the
  numbering map. The annotated edition is never pasted here.
- **`Arc42 Guide.md` — the complete instruction.** It absorbs the annotated edition's help, converted
  from descriptive prose into an imperative Section Catalogue: per section a **Write** entry (what
  belongs there, from upstream `**Contents**`) and a **Shape** entry (the form it takes, from upstream
  `**Form**`). On top of that it carries what no upstream edition has:
  - a **slot decoder** — which markup is a slot to fill, which unbolded labels in §5/§7 are content
    headings that must be kept, which slots are "repeat as many as you need", and the two pandoc HTML
    tables to rewrite as markdown;
  - the heading→number map (1–12 with subsections), so the `C4.md` / `UML.md` / `BPMN.md` citations
    resolve against headings that carry no numbers upstream;
  - the `# <Title>` + `> ` summary contract that neither edition has and `arctool sync` requires;
  - the three figures reproduced as text, and a worked example per section.
- The guide **names no other document format**. Choosing arc42 over AIC, TOGAF, C4, or TSC is the
  `aic` skill's job; a format guide that compares formats is the mess this rule prevents.
- Upstream's `**Motivation**` and `**Further Information**` blocks are **not** carried over. They are
  rationale for the template and links to `docs.arc42.org`; neither shapes a generated document.

## Justification

- **Halves the cost.** 706 lines read per arc42 generation instead of ~1540, with no loss of section
  content — the plain edition keeps every slot, and every `**Contents**` / `**Form**` point is in the
  catalogue.
- **Mechanical instead of interpretive.** The agent copies a skeleton and looks up what fills each
  slot. It never has to strip 65 teaching blocks and never risks deleting a content slot that happens
  to share a word with one.
- **Upgrades stay copy-paste.** The one file that must track upstream is verbatim and unannotated.
- **Imperative beats descriptive.** Upstream help is written to teach a human architect
  ("You should know the quality goals of your most important stakeholders…"). The catalogue restates
  the same requirements as instructions an agent can execute and a reviewer can check.
- **The numbered cross-references keep resolving**, which ADR-0009 could achieve only by owning the
  whole document.

## Trade-offs

- **The help text is paraphrased, not verbatim.** A future upstream revision of the *explanations*
  must be re-derived into the catalogue by hand; only the skeleton refreshes by copy-paste. Accepted:
  the guide was already bundle-owned and CC BY-SA permits adaptation with attribution, which the file
  carries. Section names, structure, and guidance are credited to upstream in its header.
- **Two sources can drift.** If a new plain edition renames or adds a section, the numbering map and
  the catalogue go stale. Mitigation: the guide states the refresh procedure and names re-checking the
  map as part of it; the map is a single table.
- **The bundle no longer ships arc42's figures.** With the with-help edition gone, nothing rendered
  the four PNGs, so `source/images/` is deleted — 808 KB that `install.sh` had been copying to every
  agent on every install. Their content survives as text in the catalogue (that is what made them
  droppable), and the originals stay available at arc42.org. Accepted: a reader who wants the drawn
  version goes upstream for it.
- **ADR-0009's install-size argument is given up** — though far less than it appears: the bundle now
  ships 706 lines of arc42 source against `Arc42.md`'s 363, and gets upstream fidelity plus complete
  scaffolding for it.
