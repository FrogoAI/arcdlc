# Source Library Cleanup

> Make the bundled reference library agent-grade: redact leaked data, delete docs that contradict the plan contract, merge the MDCA pair, and cut ~25–30k tokens of human-facing prose.

- Status: Approved
- Date: 2026-08-01
- Format: AIC

## Goals

### 🟢 Business Case

`skills/source-map/source/` is ArcDLC's knowledge asset: 30 documents, ~500 KB, routed by
`skills/source-map/SKILL.md` and read on demand into an agent's context. `/arcdlc:aic` generates from
it, `/arcdlc:examinate` audits code against it, `/arcdlc:policy` is governed by it. A full audit of
every file (2026-08-01) found it costs more than it returns in three distinct ways:

1. **It leaks.** `Workflow Policy.md` names nine real team members with their roles and embeds a
   customer's beta and production IP addresses — inside a bundle distributed publicly through the
   plugin marketplace and `install.sh`.
2. **It contradicts the product.** `Task Decomposition and Jira Sync.md` describes `plan.md` as
   Jira-linked phase tables with statuses `TODO/PARTIAL/WIRED/ROUTED/DONE` under `initiatives/<name>/`;
   `Workflow Policy.md` makes the Jira story the unit of execution with `plan.md` an optional
   back-link. An agent that follows either authors a plan `internal/plan` cannot parse.
3. **It is padded.** Between 40% and 50% of the library is prose written for human readers —
   philosophy sections, RACI matrices, meeting cadences, "why this matters" pitches, and one 987-line
   all-placeholder template whose key content is delivered as images an agent cannot read. None of it
   changes a single line an agent writes, and all of it is paid for on every read.

The library is also self-contradicting where it is used most: `/arcdlc:examinate` reads both MDCA
documents in full (806 lines) and they disagree on three mechanically checkable rules, with no
precedence declared.

### 🟢 Functional Overview

Four workstreams, in priority order:

1. **Redaction** — remove named individuals and customer infrastructure addresses from
   `Workflow Policy.md`.
2. **Contradiction removal** — delete `Task Decomposition and Jira Sync.md` and its routing row; cut
   `Workflow Policy.md`'s Jira-as-executor and human-process sections and point its delivery flow at
   `plan.md`; fix `Engineering Principles.md`'s "Done means deployed to production and QA-tested",
   which contradicts the plan contract's definition of `DONE`.
3. **MDCA merge** — one normative `mdca.md` built on the standard's spine, per
   [ADR-0007](../../adr/0007-mdca-standard-is-normative-single-document.md) and
   [ADR-0008](../../adr/0008-domain-imports-judged-by-purity-not-folder.md).
4. **Density and auditability** — per-document trims across the remaining library; `Arc42.md` made
   self-sufficient and `arc42/` retired per
   [ADR-0009](../../adr/0009-arc42-source-is-self-sufficient.md); `Tech Stack Canvas Original.md`
   deleted as a strict subset; `ECS.md`'s and `Go Client.md`'s project-specific halves fenced as
   example-only; the five duplicated diagram-convention blocks consolidated; citable rule IDs
   (`DDD-N`, `SOLID-N`, `ECS-N`, per-clause MDCA anchors) added to every document
   `/arcdlc:examinate` audits against.

Every file added, removed, or renamed carries a matching update to the routing table in
`skills/source-map/SKILL.md` and to the standards list in `skills/examinate/SKILL.md`.

Out of scope: `arctool`, `internal/plan`, the plan format, and any Go code — this initiative changes
no CLI behavior.

### 🟢 Quality Goals

1. **Instructional completeness beats brevity.** No trim may remove a rule an agent needs to produce
   compliant output. A document that is short but no longer sufficient is a regression, not a win.
2. **Auditability.** Every document `/arcdlc:examinate` audits against exposes rules a gap block can
   cite by identifier, not by section name.
3. **Context economy.** Target ~25–30k fewer tokens across the library's audit and generation paths,
   and ~950 KB off every install.

### 🟢 Organizational Constraints

- Single maintainer, who is also the editor of the MDCA standard — doctrine rulings are available
  in-repo rather than requiring external sign-off. The three MDCA conflicts were ruled on during this
  initiative's interview and recorded as ADR-0007 and ADR-0008.
- Execution runs through `/arcdlc:execute` with Sonnet-class executors: every task must carry its
  decisions in `HOW`, because content editing offers a weaker model unlimited room to improvise.
- No release deadline. The only version impact is a plugin-manifest bump (both manifests in lockstep).

### 🟢 Technical Constraints

- **Content has no test suite.** Acceptance criteria must be mechanically checkable — a `grep` that
  must come up empty, a file that must not exist, a heading that must be present, a line count — not
  "the document reads better".
- **Skills stay install-agnostic** (AGENTS.md hard rule). Every reference fixed or added keeps both
  path forms: `../<skill>/…` for the plugin layout and `../arcdlc-<skill>/…` for flat installs.
- **`install.sh` copies `skills/` wholesale**, so anything under `source/` ships to Claude Code,
  Codex, OpenCode, Cursor, and Antigravity alike. Bytes are a distribution cost, not just a context
  cost.
- **The routing table is the index.** AGENTS.md requires a row per document in
  `skills/source-map/SKILL.md`; a document with no row is unreachable, and a row with no document is
  a broken route.
- **Two documents are already at target density** — `KISS.md` and `Tech Stack Canvas.md` — and are
  explicitly excluded from trimming.

### 🟢 Business Context

```
        ┌────────────────────────────────────────────┐
        │        skills/source-map/source/           │
        │      (30 reference documents, ~500 KB)     │
        └────────────────────────────────────────────┘
           ▲            ▲             ▲            ▲
           │            │             │            │
    /arcdlc:aic  /arcdlc:examinate  /arcdlc:policy  ad-hoc agent
    (templates)   (audit rules)     (framework)     (source-map routing)

        distributed by ── install.sh ──▶ Claude Code plugin bundle
                                    └──▶ flat installs: Codex, OpenCode,
                                         Cursor, Antigravity
```

The library has no runtime and no consumers outside agent context windows. Its "users" are the four
pipeline skills above; its distribution boundary is public.

## Architectural Hypotheses

### 🔵 H1 — The library is agent-facing only; human process content does not belong in it

- **Context:** Much of the library was inherited from a company workspace and still addresses humans —
  team responsibilities, RACI matrices, meeting cadences, onboarding, monthly themes, daily routines.
- **Decision:** Treat "does this change what an agent produces?" as the inclusion test. Content that
  fails it is cut, not reworded.
- **Justification:** An agent pays full context price for prose it cannot act on, and human-process
  rules (stand-ups, QA sign-off) actively mislead an executor that works a `plan.md` queue.
- **Trade-offs:** The library stops doubling as onboarding material for people. Accepted: the repo's
  README and policies serve that audience, and git history retains the removed text.

### 🔵 H2 — `mdca_standard.md` is normative; the library ships one MDCA document

- **Context:** Two MDCA files, 806 lines, ~70% overlapping, three conflicting rules, no precedence.
- **Decision:** Per [ADR-0007](../../adr/0007-mdca-standard-is-normative-single-document.md), the
  standard wins and the pair merges into a single `mdca.md` on the standard's spine; layout is §7.2
  verbatim; transactions are opened only by application services (§8.6), and the §9.4 example that
  violates this is rewritten.
- **Justification:** A rule with three documented answers cannot be cited by a gap block, and an audit
  that picks a winner silently is not reproducible.
- **Trade-offs:** The prose on-ramp is lost; `Status: Draft` on the surviving document becomes
  misleading (open question Q1).

### 🔵 H3 — Domain imports are judged by purity, not by folder

- **Context:** `pkg/` is defined as "shared infrastructure" (§7.2:201), so the domain-import rule was
  expressed as a folder ban — which breaks for `pkg/clock`, a package that is not infrastructure.
- **Decision:** Per [ADR-0008](../../adr/0008-domain-imports-judged-by-purity-not-folder.md), the rule
  becomes "the domain must not import infrastructure", where an import is legal if the target's own
  imports are stdlib-only and it performs no I/O and holds no mutable global state.
- **Justification:** Purity is checkable by reading an import list; "discouraged" is not checkable at
  all.
- **Trade-offs:** Purity is transitive and can silently lapse; the clock case needs an explicit note
  that controllable time is a port decision, not a layering one.

### 🔵 H4 — `Arc42.md` becomes self-sufficient before `arc42/` is deleted

- **Context:** The bundled upstream template is 987 placeholder lines plus 948 KB of images an agent
  cannot read; `Arc42.md` is denser and is what `C4.md`/`UML.md` cross-references resolve against.
- **Decision:** Per [ADR-0009](../../adr/0009-arc42-source-is-self-sufficient.md), rewrite `Arc42.md`
  to generate a complete twelve-section arc42 document on its own, **then** remove `arc42/` and
  repoint both routing rows.
- **Justification:** Generation quality is bounded by the source an agent can parse.
- **Trade-offs:** `Arc42.md` grows against the trimming goal, and mis-sequencing would degrade arc42
  generation silently — so the plan must order the deletion strictly after the rewrite.

### 🔵 H5 — Documents that contradict the plan contract are deleted, not reconciled

- **Context:** `Task Decomposition and Jira Sync.md` conflicts with the live contract in five ways at
  once, down to calling MCP tools by names the installed server does not expose.
- **Decision:** Delete it, with its routing row. Do not re-author a Jira flow in this initiative.
- **Justification:** The cheapest correct move: it removes five of the seven worst contradictions in a
  single edit and ~2.9k tokens, and a rewrite would be a feature, not a cleanup.
- **Trade-offs:** The repo loses its Jira integration guidance entirely; re-authoring it against
  `docs/aics/<slug>/plan.md` becomes separate future work.

### 🔵 H6 — Auditable documents carry citable rule identifiers

- **Context:** `/arcdlc:examinate` files gaps as `<PREFIX>-GAP-NN` blocks whose `WHY` should cite the
  violated rule, but `ddd.md`, `solid.md`, and `ECS.md` expose only section names.
- **Decision:** Number the rules — `DDD-N` on the smells table, `SOLID-N` across the five rule tables,
  `ECS-N` on the rules block, per-clause anchors on MDCA's `P1`–`P10`.
- **Justification:** It makes a gap's `WHY` verifiable and lets a re-audit confirm the same rule.
- **Trade-offs:** IDs become a compatibility surface — renumbering later invalidates gap blocks
  already filed in existing `gap.md` registers.

### 🔵 H7 — One document per task

- **Context:** ~15 documents need independent edits, executed by a weaker model.
- **Decision:** Decompose per document (with the MDCA merge and the arc42 rewrite/delete pair as the
  exceptions that span two files), each task naming exactly what must survive.
- **Justification:** Per-file tasks keep `WHERE` to one or two files, make each commit reviewable, and
  contain the blast radius of an over-eager trim.
- **Trade-offs:** Cross-document deduplication (the five diagram-convention blocks) does not fit one
  file and needs its own task ordered after the per-file passes.

## Assessment

### 🔴 Technical Challenges & Risks

- **R1 — The leak is already distributed.** Redaction stops future exposure only: the names and IPs
  remain in git history and in every copy already installed. Whether more is warranted — history
  rewrite, address rotation, notifying the customer — is open question Q2 and is not a code decision.
- **R2 — Over-trimming.** The failure mode of this initiative is a document that no longer teaches
  what an agent needs. Mitigation: acceptance criteria state what must remain (specific headings,
  examples, rule tables), not only what must go.
- **R3 — arc42 sequencing.** Deleting `arc42/` before `Arc42.md` is complete degrades arc42 generation
  with no test to catch it. Mitigation: strict task ordering plus acceptance criteria naming all
  twelve sections.
- **R4 — The MDCA merge may surface further conflicts.** Three are known and ruled on; a 806-line
  merge can expose more. Mitigation: the merge task must re-verify the resolved rules appear exactly
  once, consistently, in the merged file.
- **R5 — Cross-reference rot.** Four documents reference `mdca.md`; deletions touch two skill files
  and the routing table. A missed row leaves a broken route that nothing tests.
- **R6 — Executor drift.** Fifteen files of prose editing offers a weaker model unlimited latitude.
  Mitigation: mechanical acceptance criteria, and `HOW` carrying the section list to cut.
- **R7 — Third-party attribution.** `Conventional Commits.md` derives from a CC BY 3.0 source without
  carrying its attribution, and arc42 material is CC BY-SA. Adopting the template's structure into
  `Arc42.md` inherits that obligation — see open question Q3.

### 🔴 Open questions

- **Q1** — On merge, should `mdca_standard.md`'s `Status: Draft` become `Active` for the surviving
  normative document, and who records that promotion?
- **Q2** — Does the leak require more than redaction: rewriting git history, rotating the exposed
  addresses, or notifying the affected customer?
- **Q3** — Should the library carry explicit attribution for third-party-derived documents
  (CC BY 3.0 for Conventional Commits, CC BY-SA for adopted arc42 structure)?
- **Q4** — `stoic.md` and `CTO Methodology Guide.md` are human-facing by nature. Trim each to its
  agent-usable core (current scope), or drop them from the library entirely?
- **Q5** — `Go Client.md` documents one private codebase's packages. Fence it as a reference
  implementation (like `ECS.md`), or genericize it into a real client-architecture guide?

## References

- [ADR-0007 — The MDCA standard is normative](../../adr/0007-mdca-standard-is-normative-single-document.md)
- [ADR-0008 — Domain imports judged by purity](../../adr/0008-domain-imports-judged-by-purity-not-folder.md)
- [ADR-0009 — `Arc42.md` must be self-sufficient](../../adr/0009-arc42-source-is-self-sufficient.md)
- `AGENTS.md` — hard rules on install-agnostic skills and the routing-table convention
- `skills/source-map/SKILL.md` — the routing table this initiative keeps in sync
- `skills/plan/references/plan-format.md` — the contract the removed documents contradict
