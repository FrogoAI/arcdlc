# Source Library Cleanup — Plan

Task format: `skills/plan/references/plan-format.md`. Architecture: [aic.md](aic.md).

All paths below are relative to the repo root. Every task edits documents only — no Go code, no
`arctool` behavior. Run `go build ./... && go test ./...` before each commit anyway: the repo must
stay green.

## Risk Coverage

- **R1 — leak already distributed** → covered by SLC-1 (working-tree redaction). Git-history rewrite
  is **accepted/not done** (engineer ruling, 2026-08-01): the exposed addresses are rotated
  out-of-band instead, so no force-push and no invalidated clones.
- **R2 — over-trimming** → covered by process: every trim task's `Acceptance` names what must
  **survive**, not only what is removed, and each `HOW` lists the exact headings to cut.
- **R3 — arc42 sequencing** → covered by SLC-6 → SLC-7 ordering; SLC-6's acceptance names all twelve
  sections, and SLC-7 may not start until SLC-6 is `DONE`.
- **R4 — MDCA merge surfaces further conflicts** → covered by SLC-5, whose acceptance requires each
  resolved rule to appear exactly once in the merged file.
- **R5 — cross-reference rot** → covered per-task (each delete updates its routing row) plus SLC-29,
  the whole-library consistency sweep.
- **R6 — executor drift** → covered by process: mechanical acceptance criteria (grep / file absence /
  heading presence) on every task.
- **R7 — third-party attribution** → covered by SLC-28 (CC BY 3.0 for Conventional Commits) and
  SLC-6 (CC BY-SA note for adopted arc42 structure).
- **Q1 — MDCA status line** → resolved: `Version: 1.0`, `Status: Active`. Covered by SLC-5.
- **Q2 — leak beyond redaction** → resolved: redact + rotate, no history rewrite. See R1.
- **Q3 — attribution** → resolved: yes, carry it. Covered by SLC-28 and SLC-6.
- **Q4 — stoic.md / CTO guide** → resolved: trim to the agent-usable core. Covered by SLC-17, SLC-18.
- **Q5 — Go Client.md** → resolved: fence as a reference implementation. Covered by SLC-14.
- **Deferred, not in this plan:** adding modern-Go coverage (generics, `context`, `slices`/`maps`,
  `errors.Join`) to `Go Best Practice.md`. That is new instruction, not cleanup; file it as its own
  initiative.

---

### SLC-1: Redact people and customer infrastructure from Workflow Policy.md

- WHAT: Remove the nine named individuals and the customer's beta/production IP addresses from `skills/source-map/source/Workflow Policy.md`.
- HOW:
  In the `## Scope` section, the table under "This policy applies to the following team members and roles"
  lists nine `| Name | Role |` rows (Ariel, Greg, Maxim, Kate, Petro, Galya, David, Vince, Mauricio).
  Replace the whole table with a role-only list (CTO, DevOps Engineer, Engineering Team Lead, Fullstack
  Engineer, AQA Engineer, Project Manager, Product Owner, Frontend Engineer, Backend Engineer) — keep the
  roles, drop every personal name.
  In `## Deployment Policy`, the environment table contains `52.60.108.117` (beta) and `52.60.174.103` (prod).
  Replace both values with placeholders (`<beta-host>`, `<prod-host>`); keep the table and its other columns.
  Sweep the rest of the file for any other personal name or customer host and apply the same treatment.
  Out of scope: every structural trim of this file (covered by SLC-3). Change only the identifying values.
- WHERE:
  Docs: `skills/source-map/source/Workflow Policy.md`.
- WHY: The file ships publicly via `install.sh` and the plugin marketplace; it currently discloses employee names and a customer's live infrastructure addresses.
- Acceptance:
  - GIVEN the redacted file WHEN `grep -rnE "52\.60\.(108|174)\.[0-9]+" skills/` runs THEN it returns no matches.
  - GIVEN the redacted file WHEN `grep -nE "\| (Ariel|Greg|Maxim|Kate|Petro|Galya|David|Vince|Mauricio) " "skills/source-map/source/Workflow Policy.md"` runs THEN it returns no matches.
  - GIVEN the redacted file WHEN its `## Scope` section is read THEN the nine roles are still listed and the section still states who the policy applies to.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: DONE.

### SLC-2: Delete Task Decomposition and Jira Sync.md

- WHAT: Remove `skills/source-map/source/Task Decomposition and Jira Sync.md` and its routing row.
- HOW:
  `git rm "skills/source-map/source/Task Decomposition and Jira Sync.md"`.
  In `skills/source-map/SKILL.md`, the row "Delivery, workflow, or Jira" lists three files; drop
  `source/Task Decomposition and Jira Sync.md` from the cell, keeping `source/Workflow Policy.md` and
  `source/tbd.md`, and drop the word "Jira" from that row's label.
  Then grep the repo for any remaining reference to the deleted filename and remove those too.
  Do not re-author a Jira flow — per AIC H5 the capability is dropped, not replaced.
  Out of scope: `Workflow Policy.md` content (SLC-1, SLC-3).
- WHERE:
  Docs: `skills/source-map/source/Task Decomposition and Jira Sync.md` (deleted), `skills/source-map/SKILL.md`.
- WHY: It describes `plan.md` as Jira-linked phase tables with statuses `TODO/PARTIAL/WIRED/ROUTED/DONE` under `initiatives/<name>/`, and calls MCP tools by names the installed server does not expose — an agent following it authors a plan `internal/plan` cannot parse.
- Acceptance:
  - GIVEN the change WHEN `ls "skills/source-map/source/Task Decomposition and Jira Sync.md"` runs THEN it exits non-zero (file absent).
  - GIVEN the change WHEN `grep -rn "Task Decomposition and Jira Sync" .` runs (excluding `.git/` and `docs/aics/source-library-cleanup/`) THEN it returns no matches.
  - GIVEN `skills/source-map/SKILL.md` WHEN its routing table is read THEN the delivery/workflow row still routes to `source/Workflow Policy.md` and `source/tbd.md`.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-3: Trim Workflow Policy.md to its agent-relevant core

- WHAT: Cut the human-process and customer-specific sections of `Workflow Policy.md` and point its delivery flow at `docs/aics/<slug>/plan.md`.
- HOW:
  Delete these whole sections: `## Team Responsibilities`, `## RACI Matrix - End-to-End Workflow`,
  `## Product Owner Operating Guide`, `## Meeting and Communication Cadence`, `## Compliance and Exceptions`,
  `## Consequences of Non-Compliance`, `## Revision History`.
  Inside `## Deployment Policy`, delete the customer-specific migration walkthrough (the subsections from the
  multi-environment frontend migration through the final full cutover, naming `andromeda-v2` / `legacy-proxy`);
  keep the generic deployment-readiness and QA-evidence rules.
  Inside `## Workflow Stages`, delete the TBD subsection and replace it with one line pointing to `source/tbd.md`.
  Keep in full: `## Purpose`, `## Policy Statement`, `## Scope` (as redacted by SLC-1), `## Governing References`,
  `## Definitions`, `## Core Workflow Principles`, the engineering implementation rules, and the QA-evidence and
  deployment-readiness rules.
  Add one paragraph under `## Core Workflow Principles`: the executable unit of work is a task block in
  `docs/aics/<slug>/plan.md`, worked by `/arcdlc:execute`, not a ticket in an external tracker; any tracker is a
  reporting mirror. Where the file still frames a Jira story as the unit of execution, reword to match.
  Out of scope: the identifying values (SLC-1 already handled them).
- WHERE:
  Docs: `skills/source-map/source/Workflow Policy.md`.
- WHY: ~400 lines of meeting cadence, RACI, and one customer's migration cost context on every read and teach a Jira-first flow that contradicts the plan contract.
- Acceptance:
  - GIVEN the trimmed file WHEN `grep -cE "^## " "skills/source-map/source/Workflow Policy.md"` runs THEN the count is at most 10 (was 16).
  - GIVEN the trimmed file WHEN `grep -n "andromeda-v2\|legacy-proxy\|RACI Matrix\|Meeting and Communication Cadence" "skills/source-map/source/Workflow Policy.md"` runs THEN it returns no matches.
  - GIVEN the trimmed file WHEN `grep -n "docs/aics/<slug>/plan.md" "skills/source-map/source/Workflow Policy.md"` runs THEN at least one match states that the plan task block is the executable unit of work.
  - GIVEN the trimmed file WHEN `## Core Workflow Principles`, `## Definitions`, and `## Deployment Policy` are read THEN each is still present with its generic rules intact.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/plan/references/plan-format.md`.
- Status: TODO.

### SLC-4: Fix and scope Engineering Principles.md

- WHAT: Correct the "Done" definition that contradicts the plan contract, add a scoping header, and cut the policy ceremony from `Engineering Principles.md`.
- HOW:
  In `## Procedures`, the line defining a task as Done only after deployment to production and QA testing
  contradicts `plan-format.md`, where `DONE` means implemented, acceptance criteria demonstrably met, validated,
  and committed. Rewrite it to the plan-contract definition, and keep deployment/QA as a separate release gate
  rather than part of task completion.
  Add a note directly under the H1: this document records one organization's engineering practice; where it names
  a specific stack (NATS, Kubernetes, Swagger, event-driven internal communication) the rule applies only to
  projects on that stack — it is not a universal mandate.
  Delete these sections: `## Purpose`, `## Policy Statement`, `## Scope`, `## Roles & Responsibilities`,
  `## Allowed & Prohibited Conduct`, `## Consequences of Non-Compliance`.
  In `## Definitions`, delete the KISS definition and replace it with a pointer to `source/KISS.md`.
  Keep `## Procedures` in full apart from the Done fix.
- WHERE:
  Docs: `skills/source-map/source/Engineering Principles.md`.
- WHY: An executor reading it marks tasks incomplete until a production deploy that `/arcdlc:execute` never performs, and imposes a stack the target project may not use.
- Acceptance:
  - GIVEN the file WHEN its `## Procedures` section is read THEN the completion rule matches `plan-format.md`'s `DONE` (implemented, acceptance met, validated, committed) and does not require production deployment or QA sign-off for task completion.
  - GIVEN the file WHEN `grep -cE "^## " "skills/source-map/source/Engineering Principles.md"` runs THEN the count is at most 3.
  - GIVEN the file WHEN the text under the H1 is read THEN it states the stack-specific rules are conditional.
  - GIVEN the file WHEN `grep -n "KISS.md" "skills/source-map/source/Engineering Principles.md"` runs THEN it returns at least one match.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/plan/references/plan-format.md`.
- Status: TODO.

### SLC-5: Merge the MDCA pair into one normative mdca.md

- WHAT: Replace `mdca.md` + `mdca_standard.md` with a single normative `mdca.md` built on the standard's spine, resolving the three known conflicts.
- HOW:
  Build the new `mdca.md` from `mdca_standard.md`'s normative content: §6 `P1`–`P10`, §7 (layout and layer
  rules), §8 (rules), §11 compliance checklist, §12 anti-patterns, Appendix A import matrix.
  Graft in the only unique content from the old `mdca.md`: the "Server, Client, Library — How MDCA Specializes"
  table, the "MDCA, SOLID, and Clean Code Alignment" table, and the five anti-pattern rows the standard lacks.
  Header: `Version: 1.0`, `Status: Active` (engineer ruling; it is now the sole normative MDCA source).
  Resolve the three conflicts — these are binding, see ADR-0007 and ADR-0008:
  1. Domain imports: state the purity rule from ADR-0008 — the domain must not import infrastructure; an import
     is legal when the target's own imports are stdlib-only (transitively) and it performs no I/O and holds no
     mutable global state. Amend §7.1's "standard library only" to this rule, and change Appendix A's
     `internal/domain/... → pkg/...` row from "discouraged" to "allowed when the target is pure (§7.1);
     forbidden otherwise". Add the note that a domain needing controllable time defines a port — importing a
     concrete clock is a testability decision, not a layering violation.
  2. Transactions: §8.6 verbatim — application services are the only place a transaction is opened. Rewrite the
     §9.4 outbox example, which currently opens one inside `domain/order/service.go`, so the transaction is
     opened by the application service.
  3. Layout: §7.2 verbatim (with `events.go`, without `mocks/`). Drop the old `mdca.md` layout entirely.
  Give `P1`–`P10` per-clause anchors so a gap block can cite `P3.2` rather than "P3"; where a principle bundles
  several checkable facts, split it into numbered clauses.
  Drop from the merged file: §5 "Relationship to Prior Art" including the differential summary, §2 Scope,
  §14 References, the `**Editor:** / **Last updated:**` lines, and the §9 examples that duplicate `ddd.md`
  (§9.1 entities, §9.2 repositories) — keep §9.4 in its corrected form.
  Then `git rm skills/source-map/source/mdca_standard.md`, update the MDCA row in `skills/source-map/SKILL.md`
  to list `source/mdca.md` only, and update `skills/examinate/SKILL.md` where it tells the agent to read
  `mdca.md` + `mdca_standard.md`.
- WHERE:
  Docs: `skills/source-map/source/mdca.md`, `skills/source-map/source/mdca_standard.md` (deleted), `skills/source-map/SKILL.md`, `skills/examinate/SKILL.md`.
- WHY: `/arcdlc:examinate` reads 806 lines with three unresolved conflicts and must silently pick a winner, so MDCA audits are not reproducible.
- Acceptance:
  - GIVEN the merge WHEN `ls skills/source-map/source/mdca_standard.md` runs THEN it exits non-zero (file absent).
  - GIVEN the merged file WHEN `grep -n "Status:" skills/source-map/source/mdca.md | head -1` runs THEN it shows `Status: Active`.
  - GIVEN the merged file WHEN the domain-import rule is searched THEN exactly one rule statement exists, it is the purity test, and no line says "standard library only" or "discouraged" for `domain → pkg/`.
  - GIVEN the merged file WHEN transaction ownership is searched THEN every statement says application services only, and no example opens a transaction inside a `domain/` file.
  - GIVEN the merged file WHEN the folder layout is searched THEN exactly one canonical layout appears, it contains `events.go`, and it does not contain `mocks/`.
  - GIVEN the merged file WHEN `grep -c "^| \`P[0-9]" skills/source-map/source/mdca.md` or the P-clause list is read THEN `P1`–`P10` each carry citable clause anchors.
  - GIVEN the repo WHEN `grep -rn "mdca_standard" skills/` runs THEN it returns no matches.
- References: `docs/aics/source-library-cleanup/aic.md`, `docs/adr/0007-mdca-standard-is-normative-single-document.md`, `docs/adr/0008-domain-imports-judged-by-purity-not-folder.md`.
- Status: TODO.

### SLC-6: Make Arc42.md able to generate a full arc42 document alone

- WHAT: Rewrite `Arc42.md` so it carries everything needed to produce a complete twelve-section arc42 document without the bundled upstream template.
- HOW:
  Read `skills/source-map/source/arc42/arc42-template-EN.md` and adopt into `Arc42.md` any required section
  structure, prompt, or table shape that `Arc42.md` currently lacks — as text, never as an image reference.
  All twelve arc42 sections must be present and numbered (1 Introduction & Goals … 12 Glossary), each with what
  it must contain and what a filled example looks like where `Arc42.md` already has one.
  Do not copy the template's `**Motivation**` and `**Further Information**` blocks, its placeholder prose, or its
  pandoc HTML tables.
  Fix two existing errors while here: the document/section comparison lists `Tasks` as an AIC section — the AIC
  template has no such section, so remove it; and "when the CTO explicitly requests it" is a stale role reference,
  replace it with the `/arcdlc:aic <slug> arc42` invocation.
  Delete the redundant `## Relationship: AIC vs Arc42 vs TOGAF` block (it restates the "When to Use" table above it).
  Add an attribution line: arc42 is published under CC BY-SA; this document adopts its structure.
  Out of scope: deleting `arc42/` (SLC-7 does that, and only after this task is `DONE`).
- WHERE:
  Docs: `skills/source-map/source/Arc42.md`, reading `skills/source-map/source/arc42/arc42-template-EN.md`.
- WHY: `arc42/` is retired next; if `Arc42.md` is not self-sufficient first, arc42 generation silently degrades and no test catches it.
- Acceptance:
  - GIVEN the rewritten file WHEN its section list is read THEN all twelve arc42 sections appear, numbered 1–12, each with its required content described.
  - GIVEN the rewritten file WHEN `grep -n "Tasks" skills/source-map/source/Arc42.md` runs THEN no line lists `Tasks` as a section of the AIC.
  - GIVEN the rewritten file WHEN `grep -n "CTO" skills/source-map/source/Arc42.md` runs THEN it returns no matches.
  - GIVEN the rewritten file WHEN `grep -n "CC BY-SA" skills/source-map/source/Arc42.md` runs THEN it returns at least one match.
  - GIVEN the rewritten file WHEN `grep -n "<img\|\.png" skills/source-map/source/Arc42.md` runs THEN it returns no matches (no content depends on an image).
- References: `docs/aics/source-library-cleanup/aic.md`, `docs/adr/0009-arc42-source-is-self-sufficient.md`.
- Status: TODO.

### SLC-7: Delete the bundled arc42/ directory

- WHAT: Remove `skills/source-map/source/arc42/` (template + 6 PNGs, 984 KB) and repoint every route to `Arc42.md`.
- HOW:
  Verify first that SLC-6 is `DONE` and `Arc42.md` lists all twelve sections; if it does not, `arctool block` this
  task rather than deleting.
  `git rm -r skills/source-map/source/arc42/`.
  In `skills/aic/SKILL.md`, the format table's `arc42` row lists `Arc42.md, arc42/arc42-template-EN.md` — reduce it
  to `Arc42.md`.
  In `skills/source-map/SKILL.md`, the "Architecture documentation" row lists both — reduce it to `source/Arc42.md`.
  Grep the repo for any other reference to `arc42/` and update it.
- WHERE:
  Docs: `skills/source-map/source/arc42/` (deleted), `skills/aic/SKILL.md`, `skills/source-map/SKILL.md`.
- WHY: 984 KB ships to every install on every supported agent, and three of the six images are referenced by nothing while the rest illustrate a file an agent cannot read.
- Acceptance:
  - GIVEN the change WHEN `ls skills/source-map/source/arc42` runs THEN it exits non-zero (directory absent).
  - GIVEN the change WHEN `grep -rn "arc42/arc42-template\|arc42/images" skills/` runs THEN it returns no matches.
  - GIVEN `skills/aic/SKILL.md` WHEN its format table is read THEN the `arc42` row's template cell names `Arc42.md` only.
  - GIVEN `du -sk skills/source-map/source` WHEN compared to before THEN it is at least 900 KB smaller.
- References: `docs/aics/source-library-cleanup/aic.md`, `docs/adr/0009-arc42-source-is-self-sufficient.md`.
- Status: TODO.

### SLC-8: Deduplicate and correct Go Server.md

- WHAT: Remove the 246-line copy of `Go Best Practice.md` from `Go Server.md`, replace the MDCA preamble with a pointer, and fix the document's violations of its own domain-purity rule.
- HOW:
  Delete the whole `## Go Best Practices` section (it duplicates `Go Best Practice.md`: formatting, naming,
  errors, control structures, defer, interfaces, concurrency, data, comments, iota, init, embedding, testing).
  Replace it with one line: general Go rules live in `Go Best Practice.md`; this document covers server
  architecture only.
  Replace the `## Modular Domain-Centric Architecture (MDCA)` preamble with a two-line summary plus a pointer to
  `mdca.md` (merged in SLC-5) — do not restate MDCA here.
  Delete `## Table of Contents` (an agent reads the file whole).
  Fix the self-contradictions, applying ADR-0008's purity rule: the domain-purity statement in
  `## Section Responsibilities` is correct, but the DDD example injects a `pkg/` messaging publisher into a
  domain service and puts a third-party `decimal.Decimal` in a domain model. Rework the example so the domain
  defines a `Publisher` port that an adapter implements, and use a domain-owned money type. In the same example,
  the publish call discards its error — handle it.
  Out of scope: `Go Best Practice.md` itself (SLC-9).
- WHERE:
  Docs: `skills/source-map/source/Go Server.md`.
- WHY: A third of the file is a stale duplicate of another document in the same directory, and its examples contradict the layering rule it states.
- Acceptance:
  - GIVEN the file WHEN `grep -n "^## Go Best Practices" "skills/source-map/source/Go Server.md"` runs THEN it returns no matches.
  - GIVEN the file WHEN `grep -c "" "skills/source-map/source/Go Server.md"` runs THEN the line count is below 560 (was 811).
  - GIVEN the file WHEN its DDD example is read THEN no `domain/` file imports a `pkg/` infrastructure package or a third-party decimal type, and every publish call handles its error.
  - GIVEN the file WHEN `grep -n "mdca.md" "skills/source-map/source/Go Server.md"` runs THEN it returns at least one match.
- References: `docs/aics/source-library-cleanup/aic.md`, `docs/adr/0008-domain-imports-judged-by-purity-not-folder.md`.
- Status: TODO.

### SLC-9: Resolve the interface contradiction and trim Go Best Practice.md

- WHAT: Fix the mutually exclusive interface rules and remove the sections that restate Go syntax.
- HOW:
  In the interfaces table, two rows conflict: "Accept interfaces, return structs" and "Don't export the type if it
  only implements an interface / Export the interface, return it from constructors. Hide the concrete type."
  Keep "accept interfaces, return structs" and delete the contradicting row — a constructor returns the concrete
  type; consumers define the interface they need (which the "Define at the consumer" row already states).
  Delete `## 3. Control Structures` (if/for/switch/type-switch syntax), the basics half of `## 5. Data`
  (`new` vs `make`, slice/map/composite-literal syntax) keeping only the "zero value is useful" rule, and
  `## 12. Printing` (the fmt verb table). These restate the language spec.
  "Zero value is useful" is currently stated twice — keep one occurrence.
  Out of scope: adding generics/`context`/`slices` coverage (deferred, see the plan's Risk Coverage).
- WHERE:
  Docs: `skills/source-map/source/Go Best Practice.md`.
- WHY: An executor reading the interfaces table gets two incompatible instructions, and ~150 lines teach Go syntax the model already knows.
- Acceptance:
  - GIVEN the file WHEN its interfaces section is read THEN exactly one rule governs what a constructor returns, and it is "return structs".
  - GIVEN the file WHEN `grep -n "^## 3\. Control Structures\|^## 12\. Printing" "skills/source-map/source/Go Best Practice.md"` runs THEN it returns no matches.
  - GIVEN the file WHEN `grep -c "zero value is useful" "skills/source-map/source/Go Best Practice.md"` (case-insensitive) runs THEN the count is 1.
  - GIVEN the file WHEN `## 2. Naming`, `## 7. Interfaces`, `## 9. Concurrency`, and `## 10. Errors` are read THEN each is still present in full.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-10: Trim Clean Code.md and fix its contradictions

- WHAT: Cut the philosophy and duplicated smells sections, and fix two statements that contradict other library documents.
- HOW:
  Delete `## Philosophy` (LeBlanc's Law, Boy Scout Rule, the 10:1 reading ratio — motivation, no rules).
  In `## 10. Smells and Heuristics (Quick Reference)`, keep only the entries that are not already stated in
  sections 1–9 or in `Go Best Practice.md` — the Go-specific ones and the handful with no earlier home; delete
  the rest of the table.
  Fix: the passage calling a struct with exported fields plus methods a "hybrid, worst of both worlds"
  contradicts `Go Server.md`'s domain model example. Reword it to what is actually meant — do not expose mutable
  state that invariants depend on — and drop the reference to a `model/` convention that does not exist
  (the layout is `domain/{feature}/model.go`).
  Fix: the testing section says no test framework is required and then blesses `testify`, and its only example
  uses `assert.InDelta`. Rewrite the example with the standard library so an agent copying it adds no dependency;
  keep at most one sentence noting `testify` is acceptable where a project already uses it.
  Flag `golang.org/x/sync/errgroup` as non-stdlib where it is recommended.
- WHERE:
  Docs: `skills/source-map/source/Clean Code.md`.
- WHY: The smells table restates the document back to itself, and two passages contradict the Go architecture documents an agent reads alongside it.
- Acceptance:
  - GIVEN the file WHEN `grep -n "^## Philosophy" "skills/source-map/source/Clean Code.md"` runs THEN it returns no matches.
  - GIVEN the file WHEN its testing section's example is read THEN it compiles against the standard library only, with no `testify` import.
  - GIVEN the file WHEN `grep -n "model/" "skills/source-map/source/Clean Code.md"` runs THEN no line claims a `model/` convention in `Go Server.md`.
  - GIVEN the file WHEN sections `## 1. Naming` through `## 9. Systems` are read THEN all nine are still present.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/source-map/source/Go Server.md`.
- Status: TODO.

### SLC-11: Trim ddd.md and add citable rule IDs

- WHAT: Remove the narrative sections, collapse two "don't use this" patterns into smell rows, number the smells, and fix two inaccuracies.
- HOW:
  Delete `## Philosophy`, `## DDD in Go vs. DDD in Java/C#`, `## DDD and the Workspace Methodologies`,
  `## Further Reading`.
  Collapse `### Factories` and `### Specifications` (both exist only to say "rarely needed in Go") into two rows
  of `## Smells and Anti-Patterns`.
  Number every row of `## Smells and Anti-Patterns` as `DDD-1` … `DDD-N` so a gap block can cite `DDD-4`.
  Fix: the repository rule "one repository per aggregate root, not per table" is unactionable where the codebase
  has no aggregates — restate it as "one repository per aggregate root (or per entity, where no aggregate
  exists)".
  Fix: the reference to `togaf.md` / `arc42.md` means generated initiative documents, not the source files
  `TOGAF.md` / `Arc42.md` — qualify them as `docs/aics/<slug>/togaf.md` and `docs/aics/<slug>/arc42.md`.
  In `## Pragmatic Adoption Checklist`, give the ✅/⚠️/❓/❌ tiers an explicit pass/fail mapping so an auditor can
  decide whether a ⚠️ item is a violation.
- WHERE:
  Docs: `skills/source-map/source/ddd.md`.
- WHY: `/arcdlc:examinate <slug> DDD` can currently cite only section names, and ~70 lines are narrative that changes no output.
- Acceptance:
  - GIVEN the file WHEN `## Smells and Anti-Patterns` is read THEN every row carries a unique `DDD-N` identifier.
  - GIVEN the file WHEN `grep -nE "^## (Philosophy|DDD in Go vs|DDD and the Workspace|Further Reading)" skills/source-map/source/ddd.md` runs THEN it returns no matches.
  - GIVEN the file WHEN the repository rule is read THEN it covers the no-aggregate case explicitly.
  - GIVEN the file WHEN `grep -n "togaf.md\|arc42.md" skills/source-map/source/ddd.md` runs THEN every match is prefixed `docs/aics/<slug>/`.
  - GIVEN the file WHEN `## Strategic DDD` and `## Tactical DDD` are read THEN both are still present in full.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-12: Trim solid.md and add citable rule IDs

- WHAT: Remove the philosophy and cross-method sections, number the rules, and fix an example that does not typecheck against `ddd.md`.
- HOW:
  Delete `## Philosophy` (the rigidity/fragility/immobility/viscosity vocabulary is never used again),
  `## SOLID and MDCA / DDD` (duplicated by the alignment table merged into `mdca.md` in SLC-5), and
  `## Further Reading`.
  Number the rules in the five principles' `### Rules` tables as `SOLID-1` … `SOLID-N`, continuous across the
  document, so a gap block can cite `SOLID-12`. Keep the "Principle Violated" column in
  `## Anti-Patterns Specific to Go` and map each row to its rule ID.
  Fix: the SRP example returns `o.Weight * 2` typed as `Money`, while `ddd.md` defines `Money` as a struct —
  make the example consistent with `ddd.md`'s definition.
  Keep all five `### Example` blocks — they are the document's teeth.
- WHERE:
  Docs: `skills/source-map/source/solid.md`.
- WHY: A SOLID audit can only cite "ISP" today, not which rule, and one example contradicts the DDD document.
- Acceptance:
  - GIVEN the file WHEN the five `### Rules` tables are read THEN every rule carries a unique `SOLID-N` identifier.
  - GIVEN the file WHEN `grep -nE "^## (Philosophy|SOLID and MDCA|Further Reading)" skills/source-map/source/solid.md` runs THEN it returns no matches.
  - GIVEN the file WHEN the SRP example is read THEN its `Money` usage matches the struct definition in `ddd.md`.
  - GIVEN the file WHEN the five principle sections are read THEN each still contains its `### Example` block.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/source-map/source/ddd.md`.
- Status: TODO.

### SLC-13: Fence ECS.md's project-specific half and add rule IDs

- WHAT: Relabel the game-specific content as a non-normative example, number the rules, and add source attribution.
- HOW:
  `## Our Implementation` and `## Adding a New Feature` document one Ebiten game (`PuzzleGrid`, `tray_piece`,
  `CaseConfig`, `BaseScene.Systems`, `CoordsFromScene`, `fsystems/`) — 96 of 152 lines. Move both under a single
  `## Reference Implementation (example only — not a conformance requirement)` heading and state in one line that
  an audit must not file gaps against a project for lacking these types.
  In `## Anti-Patterns`, the rows citing `obj.X * scaleX`, `ebiten.CursorPosition()`, and `s.state = StateXxx`
  have the same problem — either restate them in project-neutral terms or move them under the same fence.
  Number the MUST/MUST NOT blocks in `## Rules` as `ECS-C1`, `ECS-C2`, `ECS-S1`, `ECS-S2`, `ECS-E1` (components,
  systems, entities) so a gap block can cite them.
  Add a source-attribution line under the H1 (every other document in the library has one) and a cross-reference
  to `mdca.md`, which names ECS as the client-side module unit.
- WHERE:
  Docs: `skills/source-map/source/ECS.md`.
- WHY: `/arcdlc:examinate <slug> ECS` on any other codebase files gaps demanding types the target never had — the single largest false-positive source in the library.
- Acceptance:
  - GIVEN the file WHEN `## Rules` is read THEN each MUST/MUST NOT block carries a unique `ECS-*` identifier.
  - GIVEN the file WHEN the project-specific sections are read THEN they sit under a heading that marks them as an example and states they are not audit criteria.
  - GIVEN the file WHEN `grep -n "mdca.md" skills/source-map/source/ECS.md` runs THEN it returns at least one match.
  - GIVEN the file WHEN the text under the H1 is read THEN a source attribution is present.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-14: Fence Go Client.md as a reference implementation

- WHAT: Mark the project-specific packages as an example rather than an API an agent may assume, and trim the navigation cruft.
- HOW:
  Add a paragraph under the H1: this document describes one Ebiten client codebase; `ui.S()`,
  `ui.DrawRoundedRect`, `resources.NewManager()`, `core.NewCamera`, `shapes.NewBox`, `scene.GetCursorPos()`, and
  `systems.Zone` are **that project's own packages**, not Ebiten APIs — in a new project they are contracts to
  build, not imports to call. The architecture guidance (layering, scene lifecycle, resource ownership) stays
  normative; the named symbols do not.
  Add a cross-reference to `ECS.md`, whose Registry/Systems/assemblage vocabulary this document uses.
  Delete `## Table of Contents` and the trailing "Quick Reference" section (it restates the whole document).
  Delete the theming color-variable table and the window-property table (project inventory, not guidance).
  The broken `instruction.md` reference has already been fixed to `Go Server.md`; leave it.
- WHERE:
  Docs: `skills/source-map/source/Go Client.md`.
- WHY: An agent in a fresh project reads these as available APIs and hallucinates imports that exist nowhere.
- Acceptance:
  - GIVEN the file WHEN the text under the H1 is read THEN it states the named packages belong to one project and are contracts to build, not APIs to import.
  - GIVEN the file WHEN `grep -n "^## Table of Contents\|^## Quick Reference" "skills/source-map/source/Go Client.md"` runs THEN it returns no matches.
  - GIVEN the file WHEN `grep -n "ECS.md" "skills/source-map/source/Go Client.md"` runs THEN it returns at least one match.
  - GIVEN the file WHEN its architecture sections are read THEN the layering and scene-lifecycle guidance is unchanged.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/source-map/source/ECS.md`.
- Status: TODO.

### SLC-15: Trim Go Library.md

- WHAT: Remove the boilerplate sections and name the missing dependency behind the config example.
- HOW:
  Delete the Makefile, CI-workflow, `.gitignore`, and closing "Principles" sections (the last restates the
  document's own earlier sections). Keep the linter configuration — its YAML was corrected already.
  The config example uses `envDefault` struct tags and a `GetConfigFromEnv()` helper without naming the library
  that parses them: either name it explicitly or replace the example with a stdlib `os.Getenv` version.
  Reconcile the absolute "No `internal/`" rule with `Clean Code.md`'s "use `internal/` to hide implementation
  details" by stating the context: a published library keeps its API surface in the root package; `internal/` is
  for applications and for library code that must not become API.
- WHERE:
  Docs: `skills/source-map/source/Go Library.md`.
- WHY: ~80 lines are copy-paste boilerplate available from any template, and the config example cannot be followed without guessing a dependency.
- Acceptance:
  - GIVEN the file WHEN `grep -nE "^## .*(Makefile|CI|gitignore)" "skills/source-map/source/Go Library.md"` runs THEN it returns no matches.
  - GIVEN the file WHEN its config section is read THEN either the parsing library is named or the example uses only the standard library.
  - GIVEN the file WHEN the `internal/` rule is read THEN it states the library-versus-application context.
  - GIVEN the file WHEN `grep -c "" "skills/source-map/source/Go Library.md"` runs THEN the line count is below 340.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/source-map/source/Clean Code.md`.
- Status: TODO.

### SLC-16: Trim Twelve-Factor App.md

- WHAT: Cut the paraphrase half of each factor and the language-survey lists, keeping the actionable half.
- HOW:
  Each factor has a `### What It Means` subsection that paraphrases 12factor.net and a `### How to Apply` that
  states what to do. Reduce every `### What It Means` to at most three lines; keep `### How to Apply` in full.
  Delete outright: the per-language dependency-manager list under `## II. Dependencies`, the per-language
  webserver list under `## VII. Port Binding`, and the "historically, there are three gaps" table under
  `## X. Dev/Prod Parity`.
  Keep `## Quick Reference` in full — its 12-row Rule + Violation table is what makes the document auditable, and
  factors are already citable as `12F-III`.
- WHERE:
  Docs: `skills/source-map/source/Twelve-Factor App.md`.
- WHY: The paraphrase halves cost ~140 lines and add nothing an agent acts on; the apply-halves plus the quick reference already carry the whole audit.
- Acceptance:
  - GIVEN the file WHEN each `### What It Means` subsection is read THEN none exceeds three lines.
  - GIVEN the file WHEN all twelve factor sections are read THEN each still has its `### How to Apply` content intact.
  - GIVEN the file WHEN `## Quick Reference` is read THEN all twelve rows are present.
  - GIVEN the file WHEN `grep -c "" "skills/source-map/source/Twelve-Factor App.md"` runs THEN the line count is below 240.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-17: Trim stoic.md to its agent-usable core

- WHAT: Keep the decision protocols and the agent guidance; drop the human daily-practice material.
- HOW:
  Keep: `## The Short Version`, `## The Stoic Standard`, `## Control`, `## The Core Question`, `## Perception`,
  `## Action`, `## Will`, all six protocol sections (`## Decision Protocol`, `## Conflict Protocol`,
  `## Failure Protocol`, `## Anxiety Protocol`, `## Desire Protocol`, `## Anger Protocol`),
  `## Work and Craft`, `## Leadership`, `## What Stoicism Is Not`, `## Agent Response Guide`, and
  `## Agent Checklist Before Sending an Answer`.
  Delete: `## Monthly Themes`, `## Daily Stoic Routine`, `## Grief, Loss, and Change`, `## Relationships`,
  `## Money, Status, and Possessions`, `## Mortality`, `## Compact Reminders`, `## Glossary`,
  `## One-Page Daily Card`.
- WHERE:
  Docs: `skills/source-map/source/stoic.md`.
- WHY: The routed use is conflict, leadership, and decision support for an agent; monthly themes and a daily card are a human practice manual.
- Acceptance:
  - GIVEN the file WHEN `grep -nE "^## (Monthly Themes|Daily Stoic Routine|Compact Reminders|Glossary|One-Page Daily Card)" skills/source-map/source/stoic.md` runs THEN it returns no matches.
  - GIVEN the file WHEN `## Agent Response Guide` and `## Agent Checklist Before Sending an Answer` are read THEN both are present unchanged.
  - GIVEN the file WHEN the six protocol sections are read THEN all six are present.
  - GIVEN the file WHEN `grep -c "" skills/source-map/source/stoic.md` runs THEN the line count is below 350 (was 680).
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-18: Trim CTO Methodology Guide.md to its agent-usable core

- WHAT: Keep the sections that inform technical decisions; drop the people-management playbook.
- HOW:
  Keep: `## 1. Engineering Strategy (Rumelt Framework)`, `## 2. Three-Phase Planning`,
  `## 4. Trust, But Verify`, `## 8. Standards Calibration (Show. Document. Share.)`,
  `## 13. Technology Standardization (Standardize-by-Default)`, `## 14. Priority Framework (Company, Team, Me)`,
  `## 15. Corporate Values (Quality Test)`.
  Delete: sections 3, 5, 6, 7, 9, 10, 11, 12, 16, 17, 18 and `## Key References`.
  Renumber the surviving sections 1–7 so the numbering stays contiguous, and note under the H1 that the document
  is an operating-model reference, not an audit target.
- WHERE:
  Docs: `skills/source-map/source/CTO Methodology Guide.md`.
- WHY: Hiring systems, onboarding, compensation, culture surveys, and hub-office launches change nothing an agent produces.
- Acceptance:
  - GIVEN the file WHEN `grep -cE "^## [0-9]" "skills/source-map/source/CTO Methodology Guide.md"` runs THEN the count is 7.
  - GIVEN the file WHEN its numbered sections are read THEN they run 1–7 with no gaps, and the Rumelt strategy, standardize-by-default, and values-quality-test content is present.
  - GIVEN the file WHEN `grep -nE "Hiring System|Onboarding|Culture Surveys|Hub Office|First 90 Days" "skills/source-map/source/CTO Methodology Guide.md"` runs THEN it returns no matches.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-19: Delete Tech Stack Canvas Original.md

- WHAT: Remove the duplicate canvas and its routing entry.
- HOW:
  `Tech Stack Canvas Original.md` is a strict subset of `Tech Stack Canvas.md`: every section recurs there
  ("Business Goals" appears as "Business Feature Description"), and the kept file adds System Name, Services,
  Architecture Representation, and Sequence Flow. Confirm this by diffing the two section lists, then
  `git rm "skills/source-map/source/Tech Stack Canvas Original.md"`.
  In `skills/source-map/SKILL.md`, the "Tech stack decisions" row lists both files — reduce it to
  `source/Tech Stack Canvas.md`.
  Do not modify `Tech Stack Canvas.md`; it is already at target density.
- WHERE:
  Docs: `skills/source-map/source/Tech Stack Canvas Original.md` (deleted), `skills/source-map/SKILL.md`.
- WHY: It contributes no unique instruction while costing a read whenever the tech-stack row is routed.
- Acceptance:
  - GIVEN the change WHEN `ls "skills/source-map/source/Tech Stack Canvas Original.md"` runs THEN it exits non-zero.
  - GIVEN the change WHEN `grep -rn "Tech Stack Canvas Original" skills/` runs THEN it returns no matches.
  - GIVEN `skills/source-map/SKILL.md` WHEN the tech-stack row is read THEN it routes to `source/Tech Stack Canvas.md` only.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-20: Extract one shared Diagram Conventions document

- WHAT: Replace the five duplicated file-convention/palette blocks with a single routed document, resolving the palette conflict.
- HOW:
  `C4.md`, `UML.md`, `BPMN.md`, `Flowchart.md`, and `Arc42.md` each end with a near-identical `## File Convention`
  block (the dot → png → embed workflow), and two of them additionally define a color scheme.
  Create `skills/source-map/source/Diagram Conventions.md` holding: the dot → png → embed workflow once, the file
  naming/location rule, and **one** palette. Resolve the conflict explicitly — `C4.md` uses a blue family
  (`#08427B`/`#1168BD`/`#438DD5`/`#85BBF0`) while UML/BPMN/Flowchart use start-green/process-cyan; keep C4's
  palette for C4 diagrams and the shared palette for the others, stating which applies where, so a mixed document
  is not self-contradicting.
  Delete the five per-file blocks and the two color-scheme sections (`C4.md` `## Color Scheme Summary`,
  `Flowchart.md` `## Color Scheme`), replacing each with a one-line pointer to the new file.
  Add a routing row to `skills/source-map/SKILL.md` for the new document, and add it to the existing "Diagrams" row's
  reference list.
  Out of scope: the per-file content trims (SLC-21, SLC-22).
- WHERE:
  Docs: `skills/source-map/source/Diagram Conventions.md` (new), `skills/source-map/source/C4.md`, `skills/source-map/source/UML.md`, `skills/source-map/source/BPMN.md`, `skills/source-map/source/Flowchart.md`, `skills/source-map/source/Arc42.md`, `skills/source-map/SKILL.md`.
- WHY: Five copies of one convention cost five reads to learn one rule, and the two palettes collide in any document that mixes notations.
- Acceptance:
  - GIVEN the change WHEN `grep -rn "^## File Convention" skills/source-map/source/` runs THEN it matches only `Diagram Conventions.md`.
  - GIVEN `Diagram Conventions.md` WHEN it is read THEN it states which palette applies to C4 and which to UML/BPMN/Flowchart.
  - GIVEN each of the five documents WHEN its tail is read THEN it points to `Diagram Conventions.md`.
  - GIVEN `skills/source-map/SKILL.md` WHEN its routing table is read THEN `source/Diagram Conventions.md` has a row.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-21: Trim C4.md and UML.md

- WHAT: Remove the duplicated example, the restated routing tables, and the dead diagram-type rows.
- HOW:
  `C4.md`: delete the "Lists Unification" example (a second rendering of the DOT template directly above it) and
  `## Diagram Selection Decision Tree` (restates the `## When to Use C4` table). Keep the four level sections and
  their templates. Where DOT blocks hardcode another project's domain (CBF, estimator, policy-api, bloom.sync),
  replace those names with neutral placeholders so an agent does not cargo-cult them.
  `UML.md`: delete `### Deployment Diagram` (the same Kubernetes topology as `C4.md`'s deployment section) and
  point to `C4.md` instead. In `## Diagram Types Overview`, delete the rows self-labelled rarely needed or not
  applicable (Profile, Composite Structure, Object, Interaction Overview, Timing). Remove the per-subsection
  `**When**:`/`**Use in**:` lines that restate the overview table's "When to Use in Initiatives" column.
  Fix: the document lists `tsc.md` as an architecture document — no format produces it; remove it from that list.
- WHERE:
  Docs: `skills/source-map/source/C4.md`, `skills/source-map/source/UML.md`.
- WHY: Both documents restate their own routing tables and each other's deployment diagram, and UML documents five diagram types it tells the reader never to use.
- Acceptance:
  - GIVEN `C4.md` WHEN `grep -n "^## Diagram Selection Decision Tree" skills/source-map/source/C4.md` runs THEN it returns no matches.
  - GIVEN `C4.md` WHEN `grep -nE "CBF|bloom\.sync|policy-api" skills/source-map/source/C4.md` runs THEN it returns no matches.
  - GIVEN `UML.md` WHEN `grep -n "tsc.md" skills/source-map/source/UML.md` runs THEN it returns no matches.
  - GIVEN `UML.md` WHEN `## Diagram Types Overview` is read THEN no row is labelled "rarely needed" or "not applicable".
  - GIVEN both files WHEN their remaining diagram sections are read THEN every DOT template still compiles as valid DOT.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-22: Trim BPMN.md and Flowchart.md

- WHAT: Remove the notation detail the prescribed output cannot render, and the patterns duplicated in two notations.
- HOW:
  `BPMN.md`: delete `## Event Types (Detail)` and `## Task Types` — both describe glyphs ("◎ with envelope") that
  the document's own DOT output cannot draw; keep only the `<<Service Task>>` label convention from the latter.
  Of the four DOT templates, keep the XOR and Parallel ones and delete the Simple and Error templates they subsume.
  `Flowchart.md`: delete `## Decision Patterns` (all four patterns are re-rendered as DOT in `## DOT Templates`
  below), `## Flowchart vs Other Notations` (same routing rows as `## When to Use`), and `## What is a Flowchart`
  (narrative history).
- WHERE:
  Docs: `skills/source-map/source/BPMN.md`, `skills/source-map/source/Flowchart.md`.
- WHY: Both documents teach the same content twice — once as prose or ASCII, once as DOT — and BPMN taxonomies describe symbols the toolchain cannot emit.
- Acceptance:
  - GIVEN `BPMN.md` WHEN `grep -nE "^## (Event Types \(Detail\)|Task Types)" skills/source-map/source/BPMN.md` runs THEN it returns no matches.
  - GIVEN `BPMN.md` WHEN its DOT templates are read THEN the XOR and Parallel templates are present and valid DOT.
  - GIVEN `Flowchart.md` WHEN `grep -nE "^## (Decision Patterns|Flowchart vs Other Notations|What is a Flowchart)" skills/source-map/source/Flowchart.md` runs THEN it returns no matches.
  - GIVEN `Flowchart.md` WHEN `## Standard Symbols` and `## DOT Templates` are read THEN both are present in full.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-23: Trim TOGAF.md and mark it generator-only

- WHAT: Remove the decorative ADM cycle and state that the document is a template source, not an audit target.
- HOW:
  Delete the `## TOGAF ADM Phases` ASCII cycle — the same eight phases are re-listed immediately below as the
  eight document sections. Reduce `## When to Use` to three lines.
  Add one line under the H1: this document is consumed by `/arcdlc:aic <slug> togaf` to generate
  `docs/aics/<slug>/togaf.md`; it is not a standard `/arcdlc:examinate` audits against.
  In the sample rows of the architecture-requirements section, replace the invented "Pending"/"Active" statuses
  with placeholders so agents stop copying them verbatim into real documents.
  Keep the ArchiMate element and relationship tables and the DOT template — those are the working parts.
- WHERE:
  Docs: `skills/source-map/source/TOGAF.md`.
- WHY: The cycle diagram is pure decoration, and without the generator-only note an agent may try to audit code against TOGAF.
- Acceptance:
  - GIVEN the file WHEN `grep -n "^## TOGAF ADM Phases" skills/source-map/source/TOGAF.md` runs THEN it returns no matches.
  - GIVEN the file WHEN the text under the H1 is read THEN it states the document is a generator template and not an audit target.
  - GIVEN the file WHEN its ArchiMate tables and DOT template are read THEN they are present unchanged.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-24: Make ADR.md usable on its own

- WHAT: Give the ADR document one H1, the repo's filename convention, a closed status vocabulary, and a worked example.
- HOW:
  The file currently has two H1s — a preamble titled after Michael Nygard's template and a second `# Title`
  belonging to the template itself — so an agent copying "the template" emits the preamble into the ADR.
  Restructure to a single H1 plus a clearly delimited template block.
  Delete the `adr-tools` vendor recommendation; keep a one-line attribution to Nygard's original article, using
  the live URL rather than the stale `thinkrelevance.com` form.
  Add: the filename convention `docs/adr/NNNN-<kebab-title>.md` that `skills/aic/SKILL.md` already requires; a
  closed status vocabulary (Proposed, Accepted, Superseded, Deprecated) with the rule that a superseding ADR
  links the one it replaces; and one short worked ADR so the shape is demonstrated, not just described.
  Keep the Context / Decision / Consequences structure as the template's core.
- WHERE:
  Docs: `skills/source-map/source/ADR.md`.
- WHY: It is the least usable document per byte in the library — the repo's own ADRs follow conventions the template never states.
- Acceptance:
  - GIVEN the file WHEN `grep -c "^# " skills/source-map/source/ADR.md` runs THEN the count is 1.
  - GIVEN the file WHEN `grep -n "adr-tools\|thinkrelevance" skills/source-map/source/ADR.md` runs THEN it returns no matches.
  - GIVEN the file WHEN it is read THEN it states the `docs/adr/NNNN-<kebab-title>.md` convention and lists the allowed status values.
  - GIVEN the file WHEN it is read THEN it contains one complete worked example ADR.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/aic/SKILL.md`.
- Status: TODO.

### SLC-25: Fix the dangling reference in AIC Template.md

- WHAT: Give the template's Engineering Principles reference a resolvable path.
- HOW:
  Under `### 🟢 Technical Constraints`, the parenthetical "(at least follow Engineering Principles)" names a
  document without a path, unlike every other cross-reference in the library. Change it to
  `source/Engineering Principles.md`.
  Do not touch the 🟢/🔵/🔴 heading markers — they are the house style used by the existing initiative documents
  under `docs/aics/`.
  Out of scope: adding a filled example (the existing initiative AICs serve that purpose).
- WHERE:
  Docs: `skills/source-map/source/AIC Template.md`.
- WHY: An agent filling the template cannot resolve the one document it is told to follow.
- Acceptance:
  - GIVEN the file WHEN `grep -n "Engineering Principles" "skills/source-map/source/AIC Template.md"` runs THEN the match includes the path `source/Engineering Principles.md`.
  - GIVEN the file WHEN its headings are read THEN all eight canvas sections and their emoji markers are unchanged.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-26: Trim the two policy framework documents and supply the missing template

- WHAT: Cut the policy ceremony from both framework documents, add the policy template `Policy of Policies.md` mandates but never shows, and align `Policy of Initiatives.md` with the product.
- HOW:
  `Policy of Policies.md`: delete the personnel-covered role list, the RACI matrix, the
  "Allowed & Prohibited Conduct" section (restates Procedures), and "Consequences of Non-Compliance".
  Keep the Unique ID classification table. **Add** the official policy template it mandates but never shows — the
  section skeleton `/arcdlc:policy` must fill (header block, Purpose, Policy Statement, Scope, Definitions,
  Procedures, Roles & Responsibilities with RACI, Allowed & Prohibited Conduct, Consequences, Revision History),
  matching Step 3 of `skills/policy/SKILL.md` exactly.
  `Policy of Initiatives.md`: delete the RACI matrix, "Allowed & Prohibited Conduct", "Consequences of
  Non-Compliance", and the grooming choreography. Keep the standard-versus-technical split and the required AIC
  sections. Replace the decomposition step that produces "epics and user stories" with the product's actual flow:
  `/arcdlc:plan <slug>` decomposes an architecture document into task blocks in `docs/aics/<slug>/plan.md`. State
  that an initiative lives in `docs/aics/<slug>/`.
- WHERE:
  Docs: `skills/source-map/source/Policy of Policies.md`, `skills/source-map/source/Policy of Initiatives.md`.
- WHY: `/arcdlc:policy` is told to use a template the framework never provides, and the initiative policy describes a decomposition the product does not perform.
- Acceptance:
  - GIVEN `Policy of Policies.md` WHEN it is read THEN it contains a policy template whose section list matches Step 3 of `skills/policy/SKILL.md`.
  - GIVEN `Policy of Policies.md` WHEN it is read THEN the Unique ID classification table is still present.
  - GIVEN `Policy of Initiatives.md` WHEN `grep -n "epics and user stories" "skills/source-map/source/Policy of Initiatives.md"` runs THEN it returns no matches.
  - GIVEN `Policy of Initiatives.md` WHEN `grep -n "docs/aics/<slug>/" "skills/source-map/source/Policy of Initiatives.md"` runs THEN it returns at least one match.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/policy/SKILL.md`.
- Status: TODO.

### SLC-27: Trim tbd.md

- WHAT: Remove the narrative and cross-reference sections; keep the practice.
- HOW:
  Delete the philosophy/"TBD bargain" opening, the adoption-path section, the "TBD and the workspace
  methodologies" cross-reference section, and "Further Reading".
  Keep everything else — the branching rules, review cadence, feature-flag guidance, and release strategy are the
  document's working content and are not duplicated elsewhere.
- WHERE:
  Docs: `skills/source-map/source/tbd.md`.
- WHY: ~45 lines are motivation and cross-references; the rest is self-contained and actionable.
- Acceptance:
  - GIVEN the file WHEN `grep -nE "^## (Philosophy|The TBD Bargain|Adoption Path|TBD and the Workspace Methodologies|Further Reading)" skills/source-map/source/tbd.md` runs THEN it returns no matches.
  - GIVEN the file WHEN its branching, feature-flag, and release sections are read THEN they are present unchanged.
- References: `docs/aics/source-library-cleanup/aic.md`.
- Status: TODO.

### SLC-28: Trim Conventional Commits.md and carry its attribution

- WHAT: Remove the sections aimed at project maintainers and add the upstream licence attribution.
- HOW:
  Delete `## Why Use Conventional Commits` and the FAQ entries about spec-extension versioning, squash workflows,
  contributor onboarding, initial development phase, and commit-type mistakes. Keep the full specification
  clauses, the examples, and the revert recipe.
  Keep the note added to the provenance header stating ArcDLC's narrowing rules (scope is the initiative slug,
  `Refs:`/`#AI-assisted` footers required, one task = one commit).
  Add the upstream attribution: the specification is published by the Conventional Commits project under
  CC BY 3.0 — name the licence and link the source.
- WHERE:
  Docs: `skills/source-map/source/Conventional Commits.md`.
- WHY: The removed FAQ addresses maintainers of other projects, and one of its answers contradicts the one-commit-per-task rule; the upstream licence requires attribution the file does not yet carry.
- Acceptance:
  - GIVEN the file WHEN `grep -n "^## Why Use Conventional Commits" "skills/source-map/source/Conventional Commits.md"` runs THEN it returns no matches.
  - GIVEN the file WHEN `grep -n "CC BY 3.0" "skills/source-map/source/Conventional Commits.md"` runs THEN it returns at least one match.
  - GIVEN the file WHEN its specification section is read THEN all sixteen numbered clauses are present.
  - GIVEN the file WHEN the revert section is read THEN the `revert:` recipe is present.
- References: `docs/aics/source-library-cleanup/aic.md`, `skills/execute/SKILL.md`.
- Status: TODO.

### SLC-29: Library consistency sweep and version bump

- WHAT: Verify the routing table matches the files on disk, no reference dangles, and bump the plugin version.
- HOW:
  Build the list of files under `skills/source-map/source/` and compare it to the routing table in
  `skills/source-map/SKILL.md`: every document must have exactly one row, and every row must name a file that
  exists. Fix both directions.
  Check `skills/examinate/SKILL.md`'s standards list names only surviving files (MDCA in particular, after SLC-5).
  Grep every `*.md` under `skills/` for markdown links and backticked `*.md` paths, and verify each target exists;
  fix or remove the dangling ones. Preserve the dual-path forms (`../<skill>/…` and `../arcdlc-<skill>/…`) wherever
  a skill references a sibling skill.
  Bump the version in `.claude-plugin/plugin.json` and `.antigravity-plugin/plugin.json` to the same new value
  (they must stay in lockstep — CI asserts equality).
  Run the repo's checks: `go build ./...`, `go test ./...`, `gofmt -l .`, `go vet ./...`, and
  `arctool sync --check`.
- WHERE:
  Docs: `skills/source-map/SKILL.md`, `skills/examinate/SKILL.md`, `.claude-plugin/plugin.json`, `.antigravity-plugin/plugin.json`.
- WHY: This initiative adds, deletes, and renames library files across 28 tasks; a missed routing row leaves a document unreachable or a route broken, and nothing else tests it.
- Acceptance:
  - GIVEN the routing table WHEN each referenced `source/…` path is tested with `test -f` THEN every one exists.
  - GIVEN the files under `skills/source-map/source/` WHEN each is looked up in the routing table THEN every one appears in exactly one row.
  - GIVEN the repo WHEN every backticked or linked `*.md` path under `skills/` is resolved THEN no path is dangling (excluding the intentional flat-install `../arcdlc-*/…` forms).
  - GIVEN both plugin manifests WHEN their `version` fields are compared THEN they are equal and higher than `0.7.0`.
  - GIVEN the repo WHEN `go build ./... && go test ./... && gofmt -l . && go vet ./...` and `arctool sync --check` run THEN all succeed with no output from `gofmt`.
- References: `docs/aics/source-library-cleanup/aic.md`, `AGENTS.md`.
- Status: TODO.
