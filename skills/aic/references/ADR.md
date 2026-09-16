# ADR — Architecture Decision Record

**Source**: Michael Nygard, ["Documenting Architecture Decisions"](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) (2011) — the Context / Decision / Consequences structure and the sequential-numbering rule are his; the metadata block, status vocabulary, and naming rules below are this bundle's.
**Purpose**: Offline instruction for writing one architecture decision record. The template block is the whole file to be written; everything outside it is guidance, not content to copy.
**Consumed by**: `/arcdlc:aic <slug> adr`, and by any skill that records a decision reached during an interview. ADRs are **global** — they live in `docs/adr/`, never inside an initiative folder.

---

## When to write one

An ADR records the choice of one alternative over others against stated criteria — not a task, not a plan, not a
description of how something works. Write one when the choice is expensive to reverse, constrains later work, or
will otherwise be re-litigated by everyone who forgets why it was made: technology and dependency choices, boundary
and ownership rules, protocol and format contracts, deliberate deviations from a standard the project follows.

Skip it for choices that are local, cheap to undo, or already implied by an accepted ADR. One decision per file:
two decisions that can be accepted independently are two ADRs, cross-linked.

## File name and title

- Path: `docs/adr/NNNN-<kebab-title>.md`.
- `NNNN` is a zero-padded four-digit sequence number, assigned monotonically — take the highest number already in
  `docs/adr/` and add one. Numbers are never reused, not even by an ADR that was superseded or withdrawn.
- `<kebab-title>` is the title in lower-case kebab-case, e.g.
  `0008-domain-imports-judged-by-purity-not-folder.md`.
- The file carries exactly **one** H1, `# ADR-NNNN — <Title>`; every other heading is `##`.
- The title states the decision as a claim, in the present tense — "Domain imports are judged by purity, not by
  folder", not "Decide how to judge domain imports" and not "Domain import rules".

## Status

Exactly one of four values. The vocabulary is closed:

| Status | Meaning |
| --- | --- |
| `Proposed` | Written, not yet decided. Nothing may depend on it. |
| `Accepted` | In force, and binding on the code and documents it names. |
| `Superseded` | Replaced by a later ADR. The file stays — it is the record of what used to be true, and why. |
| `Deprecated` | No longer in force and not replaced; the situation it governed is gone. |

- Status only moves forward: `Proposed` → `Accepted` → (`Superseded` | `Deprecated`). A proposal that is turned
  down is simply not merged — there is no rejected state to record.
- An accepted ADR is never rewritten into a different decision. To change one, write a new ADR that supersedes it.
- **A superseding ADR links the ADR it replaces, and the replaced ADR links back** — relative paths within
  `docs/adr/`, so the chain is one hop in either direction:
  - in the new ADR, a metadata line: `- Supersedes: [ADR-0004](0004-plan-enforces-risk-mitigation-coverage.md)`
  - in the old ADR, the status line is the only edit:
    `- Status: Superseded by [ADR-0021](0021-plan-risk-coverage-is-advisory.md)`

## Template

Copy the block below into `docs/adr/NNNN-<kebab-title>.md`. It is the entire file — nothing goes above the H1.

  ```markdown
  # ADR-NNNN — <Decision, stated as a claim>

  - Status: Proposed | Accepted | Superseded by [ADR-NNNN](NNNN-<kebab-title>.md) | Deprecated
  - Date: YYYY-MM-DD
  - Initiative: [<slug>](../aics/<slug>/aic.md)

  ## Context

  The forces at work: what is happening that makes a decision necessary, with the facts, numbers, and constraints
  that bound it. Neutral — a reader must be able to reach a different conclusion from the same context.

  ## Decision

  What is being done, in the active voice and stated so it can be checked. Name the alternatives that were
  rejected and the criterion that rejected them.

  ## Consequences

  What becomes easier and what becomes harder — both, always, including the costs accepted knowingly.
  ```

Drop the `- Initiative:` line when the decision belongs to no initiative. When the rationale and the costs each need
more than a paragraph, `## Consequences` may be split into `## Justification` (why this option won) and
`## Trade-offs` (what it costs, and why that is accepted).

## Worked example

  ```markdown
  # ADR-0042 — List membership is checked in-process with a counting Bloom filter

  - Status: Accepted
  - Date: 2026-03-14
  - Initiative: [scoring-latency](../aics/scoring-latency/aic.md)

  ## Context

  Scoring checks list membership on every request — 500 TPS — and each check is an Aerospike network call costing
  ~2 ms, which is 60% of the scoring budget. The lists are read millions of times a day and changed a few hundred
  times: reads dominate writes by four orders of magnitude, and a stale read is tolerable for seconds.

  ## Decision

  Hold each list in-process in a counting Bloom filter sized for a 0.001% false-positive rate (~114 MB per
  instance), and propagate additions and removals over NATS. A negative answer is authoritative; a positive is
  confirmed against Aerospike before it is acted on. Rejected: a local LRU cache (unbounded worst case, and no
  answer for the misses that dominate this workload) and a co-located Aerospike replica (removes the network hop
  but not the ~1 ms lookup, and doubles the cluster).

  ## Consequences

  Easier: a membership check drops from ~2 ms to ~850 ns, and scoring keeps serving when Aerospike is slow or
  unreachable. Harder: +114 MB of heap per instance; the filter must be rebuilt at start-up, so start-up is no
  longer instant; and a missed NATS message shows up as a stale answer, bounded by a nightly rebuild.
  ```
