# arc42 — How to Build the Document

**What this is**: the complete instruction for turning `arc42-template-EN.md` (this directory) into a
finished arc42 document. Two files, one job each: the template is the **skeleton you copy**, this
guide is **everything you need to fill it**. No third source is required.

**The template file** is the upstream *plain* edition, kept verbatim. To move to a newer arc42
release, paste the new **plain** edition over `arc42-template-EN.md` and re-check the numbering map
below. Never paste the "with help" edition there — its explanations live in this guide instead.

**Attribution**: arc42 is created and © by Dr. Peter Hruschka, Dr. Gernot Starke and contributors,
published under CC BY-SA. See <https://arc42.org>. Section names, structure, and the guidance in the
Section Catalogue derive from the upstream template; the examples, mappings, and house rules are this
bundle's own.

Reach for arc42 when the system has multiple interacting services, several teams will build or
maintain it, deployment is non-trivial (multi-region, hybrid), or long-term maintainability and
onboarding matter.

---

## Procedure

1. **Copy** `arc42-template-EN.md` to the output path. It already carries the twelve sections, in
   order, with every fill-in slot.
2. **Number the headings** per the map below. The plain edition ships unnumbered.
3. **Clear the pandoc leftovers** (next section).
4. **Write the header**: `# <Title>` plus a one-line `> ` summary blockquote directly under it.
   ```markdown
   # Payments Platform — arc42

   > Card payment authorization and capture for the merchant portal, split into gateway, scoring, and ledger services.
   ```
   This is a hard contract: `arctool sync` parses the H1 as the initiative title and the `> ` line as
   its summary for the registry in `AGENTS.md` and `README.md`. The template has neither — it opens
   with an empty `# `.
5. **Fill every section** from the Section Catalogue below, deleting the slots you do not use.

A section that genuinely does not apply stays in place with an explicit `Not applicable — <reason>`
rather than being dropped.

## Reading the template's slots

The template is pandoc output. Four kinds of markup appear in it, and they are not the same thing:

| In the file | What it is | What to do |
|-------------|-----------|-----------|
| `&lt;Name black box 1&gt;`, `*<white box template>*`, `<Diagram or Table>`, `…​` | A slot for your content. | Fill it, or delete the whole slot. Never ship the angle brackets. |
| Unbolded label on its own line followed by an indented italic slot — `Motivation`, `Contained Building Blocks`, `Important Interfaces`, `Quality and/or Performance Features`, `Mapping of Building Blocks to Infrastructure` | **Content headings.** These are parts of arc42's white box and infrastructure templates (§5, §7). The line ends with two trailing spaces — it is a definition list. | **Keep them** and write your text under each. Deleting them destroys arc42's white box template. |
| `<Name black box 2>`, `<black box n>`, `<Runtime Scenario n>`, `<Concept n>`, `<Infrastructure Element n>`, `White Box <building block m>` | Repeat slots — "as many as you need". | Keep one per real element, renamed after it. Delete the rest. |
| `<table><colgroup>…</table>` (1.3 Stakeholders, 12 Glossary) | Pandoc HTML tables. | Rewrite as plain markdown tables. |

Also delete the empty `# ` on line 1 and the **About arc42** preamble — your header from step 4
replaces them.

## Numbering (contract)

Other references in this library cite arc42 by number: "Arc42 §5 Level 1" in `C4.md`, "Arc42 Section
6" in `UML.md` and `BPMN.md`. Never renumber or reorder.

| # | Heading | Numbered subsections |
|---|---------|----------------------|
| 1 | Introduction and Goals | 1.1 Requirements Overview · 1.2 Quality Goals · 1.3 Stakeholders |
| 2 | Architecture Constraints | — |
| 3 | Context and Scope | 3.1 Business Context · 3.2 Technical Context |
| 4 | Solution Strategy | — |
| 5 | Building Block View | 5.1 Whitebox Overall System · 5.2 Level 2 · 5.3 Level 3 |
| 6 | Runtime View | 6.1 … 6.n, one per scenario, named after it |
| 7 | Deployment View | 7.1 Infrastructure Level 1 · 7.2 Infrastructure Level 2 |
| 8 | Cross-cutting Concepts | 8.1 … 8.n, one per concept, named after it |
| 9 | Architecture Decisions | — |
| 10 | Quality Requirements | 10.1 Quality Requirements Overview · 10.2 Quality Scenarios |
| 11 | Risks and Technical Debts | — |
| 12 | Glossary | — |

---

# Section Catalogue

Each entry gives **Write** (what belongs in the section) and **Shape** (the form it takes). The
upstream *with-help* edition delivers three of its points as figures — the ISO 25010 categories
(§1.2), the building-block hierarchy (§5), and the cross-cutting concept topics (§8). All three are
reproduced below as text, so no image is needed to write any section. The originals are at
<https://arc42.org>.

## 1 Introduction and Goals

The driving forces the team must consider: underlying business goals, essential features, essential
functional requirements, quality goals for the architecture, and the relevant stakeholders with their
expectations. Three subsections.

### 1.1 Requirements Overview

**Write**: a short description of the functional requirements and driving forces — an extract, not the
backlog. Link to the requirements documents if they exist, naming the version and where to find them.

**Shape**: short text, usually a tabular use-case format. Keep the extract as short as possible:
balance readability of this document against redundancy with the requirements documents. 5–10
architecturally significant requirements is the right order of magnitude.

### 1.2 Quality Goals

**Write**: the top three (maximum five) quality goals **for the architecture** whose fulfillment
matters most to the major stakeholders. These are not project goals — the two are not necessarily
identical. Be very concrete and avoid buzzwords: these goals drive fundamental architectural
decisions and define how the architecture will be judged.

**Shape**: a table of quality goal + concrete measurable scenario, ordered by priority.

The upstream figure here is the ISO 25010:2023 category set. Choose from:

| Category | Examples |
|----------|---------|
| Functional Suitability | Correctness, completeness, appropriateness |
| Performance Efficiency | Response time, throughput, resource utilization |
| Compatibility | Interoperability, co-existence |
| Interaction Capability (Usability) | Learnability, operability, accessibility |
| Reliability | Availability, fault tolerance, recoverability |
| Security | Confidentiality, integrity, authentication |
| Maintainability | Modularity, reusability, modifiability, testability |
| Flexibility (Portability) | Adaptability, installability, scalability |
| Safety | Fail-safe behaviour, hazard warning |

```markdown
| Quality Goal | Scenario |
|-------------|----------|
| Low latency | API responds within 50ms at p99 under 1000 RPS |
| High availability | System remains operational with 99.9% uptime (8.7h downtime/year) |
| Modifiability | New payment provider integrates in < 1 week |
```

### 1.3 Stakeholders

**Write**: an explicit overview of every person, role, or organization that should know the
architecture, has to be convinced of it, has to work with it or with the code, needs its
documentation for their work, or has to decide about the system or its development. They determine
the extent and level of detail of the documentation.

**Shape**: a table of role/name, contact, and expectations with respect to the architecture and its
documentation. Replaces the template's pandoc HTML table.

```markdown
| Role/Name | Contact | Expectations |
|-----------|---------|--------------|
| Product owner / J. Doe | jdoe@example.com | Knows what is cheap and what is expensive to change |
| SRE on-call | #sre-oncall | Runbook-level clarity: what fails, what pages, how to roll back |
| New backend engineer | — | Onboards to a service in < 1 day using sections 3, 5, 8 |
```

## 2 Architecture Constraints

**Write**: every requirement that limits freedom in design decisions, implementation decisions, or
decisions about the development process. Some originate outside the system and hold for the whole
organization. Every constraint must be dealt with; some are negotiable, and saying which is part of
the documentation.

**Shape**: simple tables of constraint + explanation, subdivided when the list is long into technical
constraints, organizational and political constraints, and conventions (programming or versioning
guidelines, documentation or naming conventions).

| Category | Examples |
|----------|---------|
| **Technical** | Must use Go, must run on Kubernetes, must use NATS for internal messaging, must follow Engineering Principles (POL-ENG-001) |
| **Organizational** | Team size, timeline, budget, skill availability |
| **Conventions** | API versioning `/v1/`, branch naming `epic/{TASK}_{TITLE}`, Twelve-Factor compliance |

## 3 Context and Scope

**Write**: the delimitation of the system (the scope) from all its communication partners —
neighbouring systems and users (the context) — thereby specifying the external interfaces. The domain
and technical interfaces to communication partners are among the system's most critical aspects; make
sure they are completely understood.

**Shape**: context diagrams, and/or lists of communication partners and their interfaces. Split into
business context (domain-specific I/O) and technical context (channels, protocols, hardware).

### 3.1 Business Context

**Write**: **all** communication partners (users, IT systems) with their domain-specific inputs and
outputs or interfaces. Optionally domain-specific formats or communication protocols. Every
stakeholder should understand which data are exchanged with the system's environment.

**Shape**: any diagram that shows the system as a black box and specifies the domain interfaces.
Alternatively or additionally a table — its title is the name of your system, its three columns are
communication partner, inputs, outputs.

```markdown
| Communication Partner | Input (to the system) | Output (from the system) |
|-----------------------|-----------------------|--------------------------|
| Merchant portal user | Payment intent, refund request | Transaction status, receipt |
| Acquiring bank | Authorization result | Authorization request, capture |
```

### 3.2 Technical Context

**Write**: the technical interfaces — channels and transmission media — linking the system to its
environment, **plus** a mapping of the domain-specific I/O from 3.1 to those channels, saying which
I/O uses which channel.

**Shape**: e.g. a UML deployment diagram of channels to neighbouring systems, together with the
mapping table. The mapping is the point of this subsection — a channel diagram without it is
incomplete, and every domain I/O from 3.1 must appear.

```markdown
| Domain I/O | Channel | Protocol / Format |
|------------|---------|-------------------|
| Payment intent | Public API gateway | HTTPS REST / JSON |
| Authorization request | Bank VPN tunnel | ISO 8583 over TCP |
```

## 4 Solution Strategy

**Write**: a short summary of the fundamental decisions and solution strategies that shape the
architecture — technology decisions; the top-level decomposition, including the architectural or
design pattern used; how each key quality goal from 1.2 is achieved; and relevant organizational
decisions such as the development process or work delegated to third parties. These are the
cornerstones the later detailed decisions rest on.

**Shape**: keep the explanations short — 4–8 bullets or a table mapping quality goals to approaches.
Say what was decided and why, grounded in the problem statement, the quality goals of 1.2, and the
constraints of section 2. Refer to details in later sections rather than repeating them.

## 5 Building Block View

**Write**: the static decomposition of the system into building blocks — modules, components,
subsystems, classes, interfaces, packages, libraries, frameworks, layers, partitions, tiers,
functions, macros, operations, data structures — and their dependencies. This view is **mandatory in
every arc42 document**; by analogy to a house, it is the floor plan. It keeps the source code
understandable through abstraction and lets you talk to stakeholders without disclosing
implementation details.

**Shape**: a hierarchy of black boxes and white boxes, refined level by level. The upstream figure
here shows that hierarchy; in words:

- **Level 1** — the white box description of the overall system, together with black box descriptions
  of **all** contained building blocks.
- **Level 2** — white box descriptions of *selected* Level 1 blocks, each with black box descriptions
  of the blocks inside it.
- **Level 3+** — the same refinement applied to selected Level 2 blocks, and so on.

**Go monorepo mapping** (see `Go Server.md`): Level 1 = services in `cmd/`; Level 2 = internal
packages per service in `internal/<service>/`; Level 3 = domain modules within those packages.

### 5.1 Whitebox Overall System

**Write**: the white box template, whose four parts are already slots in the template file —

- an **overview diagram**;
- **Motivation** — why the system is decomposed this way;
- **Contained Building Blocks** — a black box description of each contained block;
- **Important Interfaces** (optional) — interfaces not covered by any single block's black box
  description but needed to understand the white box.

There is no fixed template for an interface. In the worst case you must specify syntax, semantics,
protocols, error handling, restrictions, versions, qualities, and necessary compatibilities; in the
best case an example or a simple signature is enough.

**Shape — two alternatives for the contained blocks.** Either *one* table for a short, pragmatic
overview of all blocks and their interfaces:

```markdown
| Name | Responsibility |
|------|----------------|
| api-gateway | Terminates TLS, authenticates callers, routes to services |
| scoring | Evaluates transactions against rule sets and list membership |
```

…or a list of per-block black box descriptions, one `###` heading per block, filling the black box
template: Purpose/Responsibility; Interface(s), which may include qualities and performance
characteristics; (optional) quality/performance characteristics such as availability and run-time
behaviour; (optional) directory/file location; (optional) fulfilled requirements, when you need
traceability back to 1.1; (optional) open issues, problems, and risks.

### 5.2 Level 2

**Write**: the inner structure of *some* Level 1 building blocks, as white boxes, using the same white
box template one level deeper.

**Shape**: decide which blocks are important enough to justify the detail — **prefer relevance over
completeness**. Specify the important, surprising, risky, complex, or volatile blocks. Leave out the
normal, simple, boring, or standardized parts of the system.

### 5.3 Level 3

**Write**: the inner structure of *some* Level 2 blocks, same template again. When you need more
levels than three, copy this part of the structure for each additional level. Use sparingly.

## 6 Runtime View

**Write**: the concrete behavior and interactions of the building blocks, as scenarios drawn from —

- important use cases or features: how do the building blocks execute them?
- interactions at critical external interfaces: how do the building blocks cooperate with users and
  neighbouring systems?
- operation and administration: launch, start-up, stop;
- error and exception scenarios.

The main criterion for choosing scenarios is **architectural relevance**. It is *not* important to
describe a large number of them — document a representative selection, typically 3–5.

**Shape**: one numbered subsection per scenario, named after it, each holding the diagram or the
numbered steps **plus** a description of the notable aspects of the interactions between the building
block instances shown — ordering guarantees, timeouts, retries, what happens when a step fails.
Notations: a numbered list of steps in natural language, activity diagrams or flow charts, sequence
diagrams, BPMN or event process chains, state machines.

## 7 Deployment View

**Write**: (1) the technical infrastructure the system executes on — geographical locations,
environments, computers, processors, channels, network topologies and other infrastructure elements —
and (2) the mapping of the software building blocks onto those elements. Document every environment
that differs architecturally (development, test, production). A deployment view matters especially
when the software runs distributed across more than one computer, processor, server, or container, or
when you design your own hardware. From a software perspective it is enough to capture the
infrastructure elements needed to show where the building blocks run.

**Shape**: 3.2 may already hold a top-level deployment diagram with your own infrastructure as ONE
black box; here you zoom into it. UML deployment diagrams, nested when the infrastructure is complex,
or any notation able to show nodes and channels.

Align with `Twelve-Factor App.md`: Factor X (Dev/Prod Parity) — document all environments; Factor V
(Build, Release, Run) — show the pipeline; Factor VII (Port Binding) — document port assignments.

### 7.1 Infrastructure Level 1

**Write**, filling the template's slots: an overview diagram of the distribution across locations,
environments, computers and processors with the physical connections between them; **Motivation** —
the justification for this deployment structure; **Quality and/or Performance Features** of the
infrastructure (capacity, redundancy, inter-zone latency); and **Mapping of Building Blocks to
Infrastructure**. For multiple environments or alternative deployments, copy and adapt this section
per environment.

### 7.2 Infrastructure Level 2

**Write**: the internal structure of *some* Level 1 infrastructure elements — e.g. K8s namespace
layout, database replication topology. Repeat the Level 1 structure for each element you select.

## 8 Cross-cutting Concepts

**Write**: the concepts — practices, patterns, regulations, solution ideas — that relate to multiple
building blocks and give the architecture its conceptual integrity. Some concern **all** elements of
the system, others only a few: logging typically concerns every component, security only some.

**Shape**: concept papers of any structure, example implementations (especially for technical
concepts), or cross-cutting model excerpts and scenarios using the notations of the architecture
views. Pick **only** the most-needed topics and give each its own subsection (8.1, 8.2, …). Do **not**
attempt to cover every candidate topic.

The upstream figure here is a menu of candidates:

| Concept | When to Include |
|---------|----------------|
| Domain model | Always for DDD systems |
| Error handling | When the strategy is non-obvious |
| Logging/Monitoring | How structured logging works, which metrics |
| Security | Authentication, authorization, encryption at rest/transit |
| Persistence | ORM strategy, migration approach, connection pooling |
| Communication patterns | NATS subject conventions, REST backoff strategy |
| Testing strategy | Unit vs integration vs E2E boundaries |
| Configuration | How Twelve-Factor config is managed |
| Build/Deploy | CI/CD pipeline description |

Cross-reference `Engineering Principles.md` and `Twelve-Factor App.md` where applicable.

## 9 Architecture Decisions

**Write**: the important, expensive, large-scale, or risky architecture decisions, with rationales. A
decision means selecting one alternative over others against given criteria; stakeholders must be
able to retrace them. Use judgement about *where* a decision belongs: one that concerns a single
building block may instead live in that block's white box description in section 5. Avoid redundancy
with section 4 — section 4 is the high-level *what*, section 9 the detailed *why*.

**Shape**: an ADR per important decision; or a list or table ordered by importance and consequences;
or a separate subsection per decision.

```markdown
### 9.1 ADR-001: Use Counting Bloom Filter for list membership

**Status**: Accepted
**Context**: Scoring checks list membership at 500 TPS via Aerospike network calls (~2ms each).
**Decision**: Replace with in-memory Counting Bloom Filters (~850ns lookup, 0.001% FPR).
**Consequences**: +4000x faster lookups, +114.5 MB memory per instance, requires CBF sync via NATS.
```

Decisions recorded as full ADRs under `docs/adr/` are linked from here, not copied — see `ADR.md`.

## 10 Quality Requirements

**Write**: all relevant quality requirements. The most important ones are already in 1.2 — only
*reference* them here. This section is where the lesser requirements go: the nice-to-haves that
create no high risk when they are not fully achieved.

### 10.1 Quality Requirements Overview

**Write**: a summary of the quality requirements. There are often dozens or hundreds of detailed ones,
so summarize by category or topic (the ISO 25010:2023 categories from 1.2, or the Q42 model at
<https://quality.arc42.org>). If these summaries are already precise, specific, and measurable, you
may skip 10.2 entirely.

**Shape**: a simple table, one line per category or topic with a short description. A mindmap or a
quality attribute tree — "quality" as the root, refined tree-like, called a *Quality Attribute Utility
Tree* by [Bass+21] — works as well.

### 10.2 Quality Scenarios

**Write**: concrete, measurable scenarios that turn the quality requirements into acceptance criteria.
Two kinds are especially useful:

- **Usage scenarios** — the system's runtime reaction to a stimulus, including efficiency and
  performance. *Example: the system reacts to a user's request within one second.*
- **Change scenarios** — the desired effect of a modification or extension of the system or its
  immediate environment, with the effort or duration measured. *Example: a new payment provider is
  added in under a week.*

**Shape — short form** (favoured by Q42): Context/Background — what kind of system or component, what
situation; Source/Stimulus — who or what triggers the behaviour; Metric/Acceptance Criteria — a
response including a measure.

```markdown
| ID | Category | Scenario | Metric |
|----|----------|----------|--------|
| QS-1 | Performance | Under 1000 RPS, API response time | < 50ms p99 |
| QS-2 | Reliability | Single node failure | No data loss, < 30s recovery |
| QS-3 | Security | Unauthorized API access attempt | Rejected with 401, logged, alerted |
```

**Long form** (favoured by the SEI and [Bass+21]), for the few scenarios that need it: Scenario ID,
Scenario Name, Source, Stimulus, Environment, Artifact, Response, Response Measure.

[Bass+21]: Len Bass, Paul Clements, Rick Kazman, *Software Architecture in Practice*, 4th edition,
Addison-Wesley, 2021. Worked examples of quality requirements: <https://quality.arc42.org>.

## 11 Risks and Technical Debts

**Write**: the identified technical risks and technical debts, ordered by priority, each with the
measure that would minimize, mitigate, or avoid it. Management stakeholders — project managers,
product owners — need this as part of the overall risk analysis and measurement planning.

**Shape**: an ordered list or table, highest priority first, including the suggested measures.

```markdown
| # | Risk / Debt | Probability | Impact | Mitigation |
|---|-------------|------------|--------|------------|
| 1 | CBF false positive rate increases with list size | Medium | Low | Monitor FPR metric, auto-rebuild at threshold |
| 2 | NATS delivery failure causes stale CBF | Low | High | Periodic full rebuild (daily), staleness alert |
```

## 12 Glossary

**Write**: the most important domain and technical terms stakeholders use when discussing the system,
defined once so that everyone has an identical understanding and nobody relies on synonyms or
homonyms. In multi-language teams it doubles as the source for translations.

**Shape**: a table with Term and Definition columns — plus more columns if you need translations.
Replaces the template's pandoc HTML table.

```markdown
| Term | Definition |
|------|-----------|
| CBF | Counting Bloom Filter — probabilistic data structure for set membership |
| Capture | Settling a previously authorized payment so funds actually move |
```

---

## House style

Write for the engineer who will maintain the system, not for a reviewer who will approve it: state
decisions and their consequences, keep prose short, prefer tables and diagrams over paragraphs. Every
significant statement must trace to an interview answer, an existing ADR, or code evidence — never
invent architecture conclusions.

Diagrams embedded in an arc42 document follow the shared DOT → PNG workflow, `images/` naming rule,
tooling, and palettes — see `Diagram Conventions.md`.
