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

Completed (archived to docs/aics/source-library-cleanup/plan-archive.md):
- SLC-1: Redact people and customer infrastructure from Workflow Policy.md
- SLC-2: Delete Task Decomposition and Jira Sync.md
- SLC-3: Trim Workflow Policy.md to its agent-relevant core
- SLC-4: Fix and scope Engineering Principles.md
- SLC-5: Merge the MDCA pair into one normative mdca.md
- SLC-6: Make Arc42.md able to generate a full arc42 document alone
- SLC-7: Delete the bundled arc42/ directory
- SLC-8: Deduplicate and correct Go Server.md
- SLC-9: Resolve the interface contradiction and trim Go Best Practice.md
- SLC-10: Trim Clean Code.md and fix its contradictions
- SLC-11: Trim ddd.md and add citable rule IDs
- SLC-12: Trim solid.md and add citable rule IDs
- SLC-13: Fence ECS.md's project-specific half and add rule IDs
- SLC-14: Fence Go Client.md as a reference implementation
- SLC-15: Trim Go Library.md
- SLC-16: Trim Twelve-Factor App.md
- SLC-17: Trim stoic.md to its agent-usable core
- SLC-18: Trim CTO Methodology Guide.md to its agent-usable core
- SLC-19: Delete Tech Stack Canvas Original.md
- SLC-20: Extract one shared Diagram Conventions document
- SLC-21: Trim C4.md and UML.md
- SLC-22: Trim BPMN.md and Flowchart.md
- SLC-23: Trim TOGAF.md and mark it generator-only
- SLC-24: Make ADR.md usable on its own
- SLC-25: Fix the dangling reference in AIC Template.md
- SLC-26: Trim the two policy framework documents and supply the missing template
- SLC-27: Trim tbd.md
- SLC-28: Trim Conventional Commits.md and carry its attribution
- SLC-30: Restore the TBD applicability rule
- SLC-29: Library consistency sweep and version bump
