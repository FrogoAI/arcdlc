# 0004 — `/arcdlc:plan` enforces risk-mitigation coverage before handoff

## Status

Accepted (2026-07-03).

## Context

The AIC template's "Technical Challenges & Risks" (and "Open questions") section is where an
architecture interview records what could go wrong. Today `/arcdlc:plan` decomposes the document
into tasks but nothing checks that those risks are actually addressed — a plan can hand off to
`/arcdlc:execute` with named risks that no task mitigates and no one has consciously accepted. That
contradicts ArcDLC's "audited every step, not assumed done" philosophy: the risks were surfaced,
then silently dropped.

## Decision

`/arcdlc:plan` gains a mandatory risk-coverage gate, run after decomposition and before final
validation/handoff:

1. Read the architecture document's "Technical Challenges & Risks" and "Open questions" sections.
2. For each risk, determine whether it is **covered** — either addressed by at least one plan task
   (or an explicit process mitigation recorded in the plan) or consciously **accepted/deferred**
   with a rationale.
3. If any risk is neither covered nor accepted, **run a grilling session with the engineer**
   focused on the uncovered risks, decide a mitigation for each, and write it into the plan — as a
   new task when it needs implementation, or as an accepted-risk note when it does not.
4. Record the outcome as a **Risk Coverage** mapping in the plan preamble (risk → task IDs /
   accepted), so the check is demonstrable rather than asserted.
5. The gate is a **hard gate**: handoff is blocked until every risk is either covered by a task or
   explicitly accepted. When the architecture document has no risks section (e.g. a pure
   gap-driven plan with no AIC), the gate is a no-op.

## Consequences

Easier: risks raised in the interview cannot silently evaporate; "why isn't there a task for X?"
is answered in the plan itself; the mitigation conversation happens before code, not after an
incident. Harder: `/arcdlc:plan` is no longer a pure one-shot decomposition — it may pause for a
grill; a vague or unfalsifiable risk can stall handoff until it is sharpened or explicitly
accepted (this friction is intended). "Accepted risk" is an escape valve that must stay honest —
it records a decision, not a way to skip the gate.
