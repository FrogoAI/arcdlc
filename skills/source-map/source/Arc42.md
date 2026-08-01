# Arc42 Architecture Documentation

**Source**: [arc42.org](https://arc42.org) - Template Version 9.0, by Dr. Peter Hruschka & Dr. Gernot Starke
**Purpose**: Instruction for generating full arc42 architecture documentation. This is the comprehensive version (12 sections). For lightweight inception use the AIC Template.
**Attribution**: arc42 is published under CC BY-SA by Dr. Peter Hruschka and Dr. Gernot Starke. This document adopts the arc42 section structure, numbering, and terminology; the explanations, tables, and examples are this bundle's own. The original template is at <https://arc42.org>.
**Self-contained**: everything needed to write all twelve sections is in this file. No second template is required.

---

## When to Use

| Document | When | Sections |
|----------|------|----------|
| **aic.md** (AIC Canvas) | Every initiative -- quick inception document | Business Case, Functional Overview, Quality Goals, Organizational Constraints, Technical Constraints, Business Context, Architectural Hypotheses, Technical Challenges & Risks |
| **arc42.md** (Full Arc42) | Complex/critical systems needing comprehensive architecture documentation | All 12 sections below |
| **togaf.md** (TOGAF ADM) | Enterprise-wide or compliance-heavy initiatives | See `source/TOGAF.md` |

Generate `arc42.md` when the engineer invokes `/arcdlc:aic <slug> arc42` (flat installs: `arcdlc-aic <slug> arc42`), or recommend that format when:
- The system has multiple interacting services/components
- Multiple teams will build/maintain the system
- Deployment is non-trivial (multi-region, hybrid, complex infrastructure)
- Long-term maintainability and onboarding matter

---

## Document Skeleton

A complete arc42 document has exactly these twelve top-level sections, in this order and with these numbers:

1. Introduction and Goals
2. Architecture Constraints
3. Context and Scope
4. Solution Strategy
5. Building Block View
6. Runtime View
7. Deployment View
8. Cross-cutting Concepts
9. Architecture Decisions
10. Quality Requirements
11. Risks and Technical Debt
12. Glossary

The numbering is part of the contract: other references in this library cite arc42 by number (for example "Arc42 §5 Level 1" in `C4.md`, "Arc42 Section 6" in `UML.md`). Never renumber or reorder. A section that genuinely does not apply stays in place with an explicit "Not applicable -- <reason>" rather than being dropped.

Write for the engineer who will maintain the system, not for a reviewer who will approve it: state decisions and their consequences, keep prose short, prefer tables and diagrams over paragraphs.

---

## Template Structure

### 1. Introduction and Goals

Describes the driving forces behind the architecture: the underlying business goals, the essential features and functional requirements, the quality goals for the architecture, and the relevant stakeholders with their expectations. Three subsections:

#### 1.1 Requirements Overview

Short description of the key functional requirements. Link to external requirements documents if they exist (name the version and where to find it).

**What to write**: 5-10 most important requirements as a prioritized list or use-case table. Not the full backlog -- just the architecturally significant ones. Keep the extract short: balance readability of this document against redundancy with the requirements documents.

#### 1.2 Quality Goals

The top 3-5 quality goals for the architecture (not project goals -- the two are not necessarily identical). Be concrete and avoid buzzwords: these goals drive fundamental architectural decisions and define how the architecture will be judged. Use ISO 25010 categories:

| Category | Examples |
|----------|---------|
| Performance Efficiency | Response time, throughput, resource utilization |
| Reliability | Availability, fault tolerance, recoverability |
| Security | Confidentiality, integrity, authentication |
| Maintainability | Modularity, reusability, modifiability, testability |
| Portability | Adaptability, installability |
| Compatibility | Interoperability, co-existence |
| Usability | Learnability, operability, accessibility |
| Functional Suitability | Correctness, completeness, appropriateness |

**Format**: Table with quality goal + concrete scenario (measurable), ordered by priority.

```markdown
| Quality Goal | Scenario |
|-------------|----------|
| Low latency | API responds within 50ms at p99 under 1000 RPS |
| High availability | System remains operational with 99.9% uptime (8.7h downtime/year) |
| Modifiability | New payment provider integrates in < 1 week |
```

#### 1.3 Stakeholders

Everyone who should know the architecture, has to be convinced of it, has to work with it or with the code, needs its documentation for their work, or has to decide about the system or its development. Stakeholders determine the extent and level of detail of the documentation.

**Format**: Table with Role/Name, Contact, and expectations with respect to the architecture and its documentation.

```markdown
| Role/Name | Contact | Expectations |
|-----------|---------|--------------|
| Product owner / J. Doe | jdoe@example.com | Knows what is cheap and what is expensive to change |
| SRE on-call | #sre-oncall | Runbook-level clarity: what fails, what pages, how to roll back |
| New backend engineer | -- | Onboards to a service in < 1 day using sections 3, 5, 8 |
```

---

### 2. Architecture Constraints

Constraints that limit architectural freedom -- in design, in implementation, or in the development process. Some originate outside the system and hold for the whole organization. Every constraint must be dealt with; some are negotiable, and saying which is part of the documentation. Three categories:

| Category | Examples |
|----------|---------|
| **Technical** | Must use Go, must run on Kubernetes, must use NATS for internal messaging, must follow Engineering Principles (POL-ENG-001) |
| **Organizational** | Team size, timeline, budget, skill availability |
| **Conventions** | API versioning `/v1/`, branch naming `epic/{TASK}_{TITLE}`, Twelve-Factor compliance |

**Format**: Simple table with constraint and explanation -- one table per category when the list is long.

---

### 3. Context and Scope

Delimits the system from all its communication partners (neighbouring systems and users) and thereby specifies its external interfaces. Two views:

#### 3.1 Business Context

The system as a **black box**. Shows **all** communication partners (users, external systems) and what data/messages flow in and out. Focus on **domain-level** inputs and outputs.

**Diagram**: C4 Context diagram (DOT format). Include:
- All human actors
- All external systems
- All data flows with labels describing the domain content

**Table alternative** (instead of, or alongside, the diagram) -- the table is titled with the system name:

```markdown
| Communication Partner | Input (to the system) | Output (from the system) |
|-----------------------|-----------------------|--------------------------|
| Merchant portal user | Payment intent, refund request | Transaction status, receipt |
| Acquiring bank | Authorization result | Authorization request, capture |
```

#### 3.2 Technical Context

Maps domain interfaces to technical channels: protocols (REST, gRPC, NATS), formats (JSON, Protobuf), hardware (load balancer, VPN, CDN).

**Diagram**: Deployment-aware context diagram showing protocols and channels.

**Mapping table**: every domain input/output from 3.1 mapped to the channel that carries it. This mapping is the point of the subsection -- a channel diagram without it is incomplete.

```markdown
| Domain I/O | Channel | Protocol / Format |
|------------|---------|-------------------|
| Payment intent | Public API gateway | HTTPS REST / JSON |
| Authorization request | Bank VPN tunnel | ISO 8583 over TCP |
```

---

### 4. Solution Strategy

Fundamental decisions and solution approach. Short and strategic:

- Technology decisions (language, framework, database)
- Top-level decomposition (monorepo, microservices, monolith) -- including the architectural or design pattern chosen
- Quality goal strategies (how each goal from 1.2 is addressed)
- Organizational decisions (team structure, development process, work delegated to third parties)

**Format**: 4-8 bullet points or a table mapping quality goals to solution approaches. State what was decided and why, grounded in the problem statement, the quality goals of 1.2, and the constraints of section 2. Refer to details in later sections rather than repeating them.

---

### 5. Building Block View

**Static decomposition** of the system into building blocks (modules, components, subsystems, packages, libraries, layers, ...) and their dependencies. This view is **mandatory** in every arc42 document -- by analogy to a house, it is the floor plan.

The view is a hierarchy of black boxes and white boxes, refined level by level:

- **Level 1** -- white box description of the overall system, together with black box descriptions of every building block it contains.
- **Level 2** -- white box descriptions of selected Level 1 blocks, each with black box descriptions of the blocks inside it.
- **Level 3+** -- the same refinement applied to selected Level 2 blocks, and so on.

#### Level 1: Overall System (White Box)

**White box template** -- every white box, at any level, carries these four parts:

| Part | Content |
|------|---------|
| Overview diagram | C4 Container or component diagram of the contained blocks |
| Motivation | Why the system is decomposed this way |
| Contained building blocks | Black box description of each contained block |
| Important interfaces | Interfaces needed to understand the white box that no single black box owns |

**Black box description** (per block):

| Field | Description |
|-------|-------------|
| Purpose/Responsibility | What it does |
| Interface(s) | APIs, NATS subjects, events consumed/produced |
| Quality/Performance | SLA, throughput, latency |
| Directory/File Location | Path in the monorepo (if applicable) |
| Fulfilled Requirements | Traceability back to 1.1 (optional) |
| Open Issues/Risks | Known problems |

For a short, pragmatic overview, use the compact form instead -- one row per block, name and responsibility only:

```markdown
| Name | Responsibility |
|------|----------------|
| api-gateway | Terminates TLS, authenticates callers, routes to services |
| scoring | Evaluates transactions against rule sets and list membership |
```

#### Level 2: Zoom Into Important Blocks

White-box descriptions of selected Level 1 blocks. Same white box template, one level deeper.

Decide which blocks earn this detail: **prefer relevance over completeness**. Describe the important, surprising, risky, complex, or volatile blocks; leave out the normal, simple, boring, or standardized ones.

#### Level 3+ (if needed)

Further zoom into complex subsystems, using the same white box template again. Use sparingly -- prefer relevance over completeness.

**Go-specific**: Map to monorepo structure from Go Server.md:
- Level 1 = services in `cmd/`
- Level 2 = internal packages per service in `internal/<service>/`
- Level 3 = domain modules within packages

---

### 6. Runtime View

**Dynamic behavior**: How building blocks interact at runtime. Document architecturally significant scenarios:

- Important use cases / features -- how do the building blocks execute them?
- Critical external interface interactions -- how does the system cooperate with users and neighbouring systems?
- Operation and administration: startup / shutdown sequences
- Error and exception scenarios

**Format**: Sequence diagrams (ASCII or DOT), numbered step lists, activity or flow charts, BPMN, state machines. Pick the notation that best fits the scenario.

**Per scenario**: a numbered heading naming it (6.1, 6.2, ...), the diagram or the numbered steps, and a short note on the notable aspects of the interaction -- ordering guarantees, timeouts, retries, what happens when a step fails.

**How many**: 3-5 key scenarios. Don't document all flows -- focus on architecturally relevant ones (complex, risky, critical-path).

---

### 7. Deployment View

**Infrastructure mapping**: the technical infrastructure the system executes on -- locations, environments, machines, processors, channels, network topology -- and the mapping of building blocks onto it. Document every environment that differs architecturally (development, test, production).

#### Level 1: Infrastructure Overview

- Deployment diagram showing nodes (K8s cluster, databases, message broker, CDN) and the physical connections between them
- Mapping of building blocks to infrastructure elements
- Justification for this deployment structure
- Quality and performance features of the infrastructure (capacity, redundancy, inter-zone latency)

#### Level 2: Detailed Infrastructure (if needed)

Zoom into specific infrastructure elements (e.g., K8s namespace layout, database replication topology), repeating the Level 1 structure for each element selected.

**Align with Twelve-Factor App**:
- Factor X (Dev/Prod Parity): Document all environments
- Factor V (Build, Release, Run): Show the pipeline
- Factor VII (Port Binding): Document port assignments

---

### 8. Cross-cutting Concepts

Concepts -- practices, patterns, regulations, solution ideas -- that apply across multiple building blocks and give the architecture its conceptual integrity. Pick ONLY the relevant ones; do not attempt to cover every topic listed:

| Concept | When to Include |
|---------|----------------|
| Domain model | Always for DDD systems |
| Error handling | When strategy is non-obvious |
| Logging/Monitoring | How structured logging works, which metrics |
| Security | Authentication, authorization, encryption at rest/transit |
| Persistence | ORM strategy, migration approach, connection pooling |
| Communication patterns | NATS subject conventions, REST backoff strategy |
| Testing strategy | Unit vs integration vs E2E boundaries |
| Configuration | How Twelve-Factor config is managed |
| Build/Deploy | CI/CD pipeline description |

**Format**: one numbered sub-section per concept (8.1, 8.2, ...) with explanation and examples -- a short concept note, an example implementation, or a cross-cutting model excerpt. Cross-reference Engineering Principles and Twelve-Factor App where applicable.

---

### 9. Architecture Decisions

Important, expensive, large-scale, or risky architectural decisions with rationale. A decision is the choice of one alternative over others against stated criteria.

**Format**: Architecture Decision Records (ADR):

```markdown
### ADR-001: Use Counting Bloom Filter for list membership

**Status**: Accepted
**Context**: Scoring checks list membership at 500 TPS via Aerospike network calls (~2ms each).
**Decision**: Replace with in-memory Counting Bloom Filters (~850ns lookup, 0.001% FPR).
**Consequences**: +4000x faster lookups, +114.5 MB memory per instance, requires CBF sync via NATS.
```

Alternatives when a full ADR per decision is too heavy: a list or table ordered by importance and consequences, or one sub-section per decision. A decision that concerns exactly one building block may instead live in that block's white box description in section 5.

Avoid redundancy with Section 4 (Solution Strategy). Section 4 = high-level "what"; Section 9 = detailed "why" for specific decisions.

---

### 10. Quality Requirements

Detailed quality requirements beyond the top 3-5 from Section 1.2 -- including the lesser ones that are nice-to-have and carry no high risk if missed. Reference the goals already stated in 1.2 rather than restating them.

#### 10.1 Quality Requirements Overview

Table or tree of all quality categories with descriptions (the ISO 25010 categories from 1.2 make a serviceable spine; a quality attribute tree works as well). If these summaries are already specific and measurable, 10.2 may be skipped.

#### 10.2 Quality Scenarios

Concrete, measurable acceptance criteria. Two kinds are especially useful:

- **Usage scenarios** -- the system's runtime reaction to a stimulus ("responds to a user request within one second").
- **Change scenarios** -- the effect, effort, or duration of a modification or extension ("a new payment provider is added in under a week").

Short form -- **Context** (which component, which situation), **Source/Stimulus** (who or what triggers it), **Metric/Acceptance Criteria** (the measured response):

```markdown
| ID | Category | Scenario | Metric |
|----|----------|----------|--------|
| QS-1 | Performance | Under 1000 RPS, API response time | < 50ms p99 |
| QS-2 | Reliability | Single node failure | No data loss, < 30s recovery |
| QS-3 | Security | Unauthorized API access attempt | Rejected with 401, logged, alerted |
```

Long form, for the few scenarios that need it: Scenario ID, Scenario Name, Source, Stimulus, Environment, Artifact, Response, Response Measure.

---

### 11. Risks and Technical Debt

Ordered list of identified risks and technical debt, highest priority first, each with the measure that would minimize, mitigate, or avoid it:

```markdown
| # | Risk / Debt | Probability | Impact | Mitigation |
|---|-------------|------------|--------|------------|
| 1 | CBF false positive rate increases with list size | Medium | Low | Monitor FPR metric, auto-rebuild at threshold |
| 2 | NATS delivery failure causes stale CBF | Low | High | Periodic full rebuild (daily), staleness alert |
```

---

### 12. Glossary

The domain and technical terms stakeholders use when discussing the system, defined once so that nobody relies on synonyms or homonyms. Essential for onboarding and cross-team communication.

```markdown
| Term | Definition |
|------|-----------|
| CBF | Counting Bloom Filter -- probabilistic data structure for set membership |
| AIC | Architecture Inception Canvas -- lightweight arc42 subset for initiative initiation |
```

---

## Diagram Conventions

All diagrams live in the `images/` folder inside the initiative directory.

**Dual format** (always both):
- Raw source: `images/<name>.dot`
- Compiled image: `images/<name>.png`
- Compile: `dot -Tpng images/<name>.dot -o images/<name>.png`
- Embed in docs: `![Title](images/<name>.png)`

**Naming**: lowercase, underscores, descriptive:
- `c4_context.dot` / `.png` — system in its environment
- `c4_container.dot` / `.png` — services/modules decomposition
- `archimate_business.dot` / `.png` — ArchiMate business layer
- `deployment.dot` / `.png` — infrastructure mapping
- `sequence_<flow>.dot` / `.png` — runtime scenario (or ASCII inline for simple flows)

**Tooling**:
- **Primary**: Graphviz DOT for C4, ArchiMate, component, deployment diagrams
- **Fallback**: Markdown ASCII for simple sequence flows (inline in docs, no image file needed)
