# CTO Methodology Guide

**Source**: Will Larson, "The Engineering Executive's Primer" (2024)
**Purpose**: Extracted methodologies, frameworks, and decision models for use in initiative planning and engineering leadership.
**Scope**: This is an operating-model reference, not an audit target — nothing in it is a rule that code or documents are checked against.

---

## 1. Engineering Strategy (Rumelt Framework)

**Framework source**: Richard Rumelt, "Good Strategy Bad Strategy".
**When**: Within first 6 months of role, update annually.

A strategy has three parts:

| Part | Purpose | Key Questions |
|------|---------|---------------|
| **Diagnosis** | Theory describing the root cause of the problem | What is actually happening? What data confirms it? |
| **Guiding Policy** | Approach to address diagnosed problems, with explicit trade-offs | How are resources allocated? What rules must all teams follow? How are decisions made? |
| **Coherent Actions** | Concrete steps to implement the policy | What enforcement mechanisms exist? What transformations are needed? |

**Quality tests for Guiding Policy**:
- **Applicability**: Can it navigate real-world trade-offs?
- **Enforceability**: Will teams actually be held to it?
- **Effectiveness**: Does it create a multiplier effect?

**Process** (10 steps):
1. Write it yourself (do not delegate)
2. Write for engineering leaders (staff+ engineers, senior managers)
3. Define stakeholder list, select 3-5 for a working group
4. Draft diagnosis first, iterate one-on-one with working group
5. Draft guiding policy, get feedback from 2-3 external peers
6. Share diagnosis + policy with all stakeholders
7. Draft coherent actions, iterate with working group
8. Talk individually to likely dissenters
9. Present to full engineering org, 1-week feedback window
10. Finalize, promise to evaluate in 2 months, update annually

---

## 2. Three-Phase Planning

**When**: Annual and quarterly planning cycles.

| Phase | Scope | Key Rule |
|-------|-------|----------|
| **1. Financial Plan** | P&L, budget, headcount | Keep fixed for the year; frequent revisions destroy accountability |
| **2. Resource Allocation** | Distribute capacity (e.g., product 63%, infra 25%, DevEx 12%) | Keep stable; frequent changes destroy morale |
| **3. Roadmap Alignment** | Cross-functional agreement on deliverables | Separate proven (80%) from unproven (20%) projects |

**Org design math**: Teams of ~8, groups of 4-6 teams, continue until 5-7 groups report to you.

**Anti-patterns**:
- Planning as checkbox exercise (plans never revisited)
- Headcount as universal cure
- Only "exciting" work gets prioritized
- Over-detailed plans removing team autonomy

---

## 3. Trust, But Verify

**Principle**: Blind trust is not a management technique. It prevents distinguishing good processes with bad outcomes from bad processes with good outcomes.

**Four verification instruments**:

| Instrument | Method |
|-----------|--------|
| **Review meetings** | Weekly/monthly metric reviews (quarterly is too infrequent) |
| **Primary source investigation** | Talk directly to people doing the work, bypass summary layers |
| **Direct data work** | Maintain your own small data sources, cross-reference |
| **Intolerance for inconsistencies** | When something doesn't add up, dig in immediately |

**Cycle**: Trust -> Verify -> Return with findings -> Solve together.

---

## 4. Standards Calibration (Show. Document. Share.)

**When**: Your standards exceed those of peers and you want to influence without authority.

1. **Show** -- Demonstrate the desired standard yourself through several iterations
2. **Document** -- Create a clear document explaining how to replicate, framing benefits from the audience's perspective
3. **Share** -- Send to teams you want to influence; engage the interested, don't pressure the rest

**Key insight**: Telling the CEO about a peer's low standards is really telling them they failed to solve a known problem. Escalate cautiously.

---

## 5. Technology Standardization (Standardize-by-Default)

**Rule**: Standardize by default. Research only when the improvement is **at least 10x** on one dimension without degrading others.

- Fewer technologies = deeper investment in each
- Cap concurrent research experiments (2-3 for ~1000 engineers)
- Binary outcome per research: Continue (migrate from old) or Stop (abandon cleanly)
- Over-standardization creates dead ends; over-research creates abandoned half-migrations

---

## 6. Priority Framework (Company, Team, Me)

**Default**: Prioritize company > team > personal interests.
**Evolution**: "Eventual Quid Pro Quo" -- mostly follow the hierarchy, but periodically prioritize energy-restoring work. If depletion persists >1 year, the role needs changing.

**Test**: Energy-restoring work should be "orthogonal but not opposite" to company needs (1-2 conference talks/year is fine; 8-10 is too many).

---

## 7. Corporate Values (Quality Test)

A useful value must pass three tests:

| Test | Question |
|------|----------|
| **Has an alternative** | Can you invert it and still get something reasonable? |
| **Practical significance** | Can you use it to navigate real trade-offs? |
| **Actually practiced** | Does it describe how people actually behave? |

**Recommended engineering values**:
- "Create new opportunities rather than fight over existing ones"
- "Always go to vendors unless it's a core competency"
- "Follow existing patterns unless there's a 10x improvement"
- "Be curious in conflict resolution"

**Process**: Wait 6+ months before introducing. First change behavior, then document as a value.
