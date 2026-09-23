# TOGAF ADM & ArchiMate Instruction

**Reviewed**: 2026-09-23

**Source**: TOGAF Standard (The Open Group), ArchiMate 3.2 Specification
**Purpose**: Offline instruction for generating TOGAF Architecture Development Method documentation with ArchiMate diagrams. No internet connection required.
**Consumed by**: `/arcdlc:aic <slug> togaf`, which generates `docs/aics/<slug>/togaf.md` from this template. TOGAF is **not** one of the standards `/arcdlc:examinate` audits code against — this is a generator template, not an audit target.

---

## When to Use

Generate `togaf.md` only on explicit request, or recommend it when the initiative is enterprise-wide or
cross-domain (multiple business units affected), spans a phased legacy migration, or is compliance-heavy
(SOC2, PCI-DSS, GDPR) — i.e. when arc42's scope is too narrow and stakeholders require formal EA artifacts.

---

## Document Structure (8 Sections)

The **Ask** lines under a section are the questions `/arcdlc:aic` puts to the engineer to fill it,
and are not copied into the document. Every section ends answered, deferred to §8.4 Open questions
with who answers it and by which phase, or `Not applicable: <reason>`. None is deleted.

### 1. Architecture Vision (Phase A)

Define the problem, stakeholders, principles, and scope.

#### 1.1 Problem Statement
Clear description of the business problem or opportunity driving this architecture work.

**Ask**: what does the problem cost today, in money, time or risk, and where does that number come
from? Who decides that the target state has been reached?

#### 1.2 Stakeholders and Concerns

| Stakeholder | Role | Key Concerns |
|:------------|:-----|:-------------|
| *role* | *what they do* | *what they care about architecturally* |

#### 1.3 Architecture Principles

| # | Principle | Rationale | Implications |
|:--|:----------|:----------|:-------------|
| 1 | *principle name* | *why this matters* | *what this forces/enables* |

**Ask**: which principle would change the design most if it were dropped, and who can drop it?

Reference Engineering Principles (POL-ENG-001) for technology principles. Add initiative-specific principles here.

#### 1.4 Vision Summary
1-2 paragraph narrative of the target state. Include a high-level ArchiMate motivation diagram if helpful.

#### 1.5 Scope

| Dimension | In Scope | Out of Scope |
|:----------|:---------|:-------------|
| Business processes | ... | ... |
| Applications/services | ... | ... |
| Data domains | ... | ... |
| Technology platforms | ... | ... |

**Ask**: what is left out on purpose, and what does leaving it out cost? Which quality goals, taken
from the ISO 25010 list, must the target state meet, and what measurement proves each one?

---

### 2. Business Architecture (Phase B)

Model the business layer: capabilities, processes, services, actors.

#### 2.1 Business Capability Map
Table of capabilities with description and supporting systems.

#### 2.2 Business Processes
For each key process:
- Trigger, Actors, Outcome
- Process flow diagram (BPMN or ASCII)

**Ask**: which business units and partners does each process cross, and what must each of them
change? Which step fails, times out, or runs twice most often, and what happens then?

#### 2.3 ArchiMate Business Layer Diagram

Use DOT with ArchiMate notation. **Color: `#FFFFB5`** (pale yellow) for all business elements.

**ArchiMate Business Elements**:

| Element | Notation | Description |
|---------|----------|-------------|
| Business Actor | `<<Business Actor>>` | Person or organization |
| Business Role | `<<Business Role>>` | Responsibility assigned to an actor |
| Business Process | `<<Business Process>>` | Sequence of business behaviors |
| Business Service | `<<Business Service>>` | Externally visible business behavior |
| Business Object | `<<Business Object>>` | Passive element (data, document) |
| Business Event | `<<Business Event>>` | Something that happens and triggers behavior |

**ArchiMate Relationships**:

| Relationship | Meaning | Arrow |
|-------------|---------|-------|
| Composition | "is part of" | Filled diamond |
| Aggregation | "groups" | Open diamond |
| Assignment | "performs / is allocated to" | Filled circle -> |
| Realization | "realizes" | Dashed open arrow |
| Serving | "serves / used by" | Open arrow |
| Triggering | "triggers" | Filled arrow |
| Flow | "transfers" | Dashed filled arrow |
| Access | "reads/writes" | Dashed arrow |

---

### 3. Information Systems Architecture (Phase C)

Two sub-architectures:

#### 3.1 Application Architecture

Model the application layer. **Color: `#B5FFFF`** (pale cyan) for all application elements.

**ArchiMate Application Elements**:

| Element | Notation | Description |
|---------|----------|-------------|
| Application Component | `<<App Component>>` | Modular, deployable unit of software |
| Application Service | `<<App Service>>` | Externally visible application behavior |
| Application Interface | `<<App Interface>>` | Access point (API endpoint, UI) |
| Application Function | `<<App Function>>` | Internal behavior element |
| Application Event | `<<App Event>>` | Application-level event |
| Application Process | `<<App Process>>` | Sequence of application behaviors |

**What to document**:
- C4 Container diagram adapted to ArchiMate notation
- Service inventory: name, responsibility, APIs, events produced/consumed
- Integration patterns: NATS subjects (per Engineering Principles), REST endpoints

**Ask**: why do the application boundaries sit there, and why not the obvious other split? Does each
new component have a reason that an existing one cannot meet? Which partners can an attacker
control, and which can be slow or absent without warning? When a dependency is slow, and when it is
gone, does the system fail open or closed?

#### 3.2 Data Architecture

- Data entities and their relationships (ER diagram or ArchiMate data objects)
- Data flow between services
- Storage technology mapping (which service uses which DB)
- Data ownership: which service is the source of truth for which entity

**Ask**: how many bytes per record, how many records alive at once, and for how long? Which data is
personal, and where is it allowed to live? How many reads and writes does one business action make?

**ArchiMate Data Elements**:

| Element | Notation | Description |
|---------|----------|-------------|
| Data Object | `<<Data Object>>` | Passive element used by application components |

---

### 4. Technology Architecture (Phase D)

Model the technology/infrastructure layer. **Color: `#C9E7B7`** (pale green) for all technology elements.

**ArchiMate Technology Elements**:

| Element | Notation | Description |
|---------|----------|-------------|
| Node | `<<Node>>` | Computational resource (server, VM, container) |
| Device | `<<Device>>` | Physical hardware |
| System Software | `<<System Software>>` | Software platform (OS, DB engine, runtime) |
| Technology Service | `<<Tech Service>>` | Externally visible infra behavior (DNS, LB) |
| Artifact | `<<Artifact>>` | Physical piece of data (binary, config file, image) |
| Communication Network | `<<Network>>` | Communication medium (VPN, VPC, internet) |
| Technology Interface | `<<Tech Interface>>` | Access point on technology layer (port, socket) |
| Technology Process | `<<Tech Process>>` | Sequence of technology behaviors |
| Technology Event | `<<Tech Event>>` | Technology-level event |

**What to document**:
- Deployment diagram: K8s cluster layout, database nodes, message broker, CDN
- Network topology: VPCs, security groups, ingress/egress
- Infrastructure as Code references
- Align with Twelve-Factor App factors VII (Port Binding), IX (Disposability), X (Dev/Prod Parity)

**Ask**: what load does the target state carry today and at the design target, and where do both
numbers come from? What is the throughput limit of each chosen platform, where was it measured, and
what happens when it is reached? What does it cost per month? Which metrics and alert values tell
the person woken at 03:00 what broke?

---

### 5. Opportunities and Solutions (Phase E)

Gap analysis between current (baseline) and target architectures.

#### 5.1 Gap Analysis

| Component | Baseline | Target | Gap | Action |
|-----------|----------|--------|-----|--------|
| *name* | *current state* | *desired state* | *what's missing* | *build / buy / reuse / migrate* |

**Ask**: for each gap, why build rather than buy or reuse, and what option lost?

#### 5.2 Solution Building Blocks

Map gaps to concrete deliverables:
- New services to build
- Existing services to modify
- Infrastructure to provision
- Data migrations to execute

---

### 6. Migration Planning (Phase F)

Define implementation phases using the Genesis -> Custom -> Product model from Policy of Initiatives (POL-TECH-001).

#### 6.1 Transition Architectures

| Phase | Name | Scope | Duration | Deliverables |
|-------|------|-------|----------|-------------|
| Genesis | *name* | PoC + core foundation | *estimate* | *what ships* |
| Custom | *name* | MVP production-grade | *estimate* | *what ships* |
| Product | *name* | Full feature set | *estimate* | *what ships* |

**Ask**: what condition opens each next phase? How long does the old path run beside the new one,
how fast can a phase be rolled back, and who decides?

#### 6.2 Migration Dependency Diagram

Show dependencies between migration work packages. Which must complete before others can start.

---

### 7. Implementation Governance (Phase G)

#### 7.1 Architecture Compliance

| Principle | Compliance Check | Status |
|-----------|-----------------|--------|
| Performance by Design | Benchmark results for hot paths | *status* |
| KISS | Code review checklist item | *status* |
| Event-Driven | All internal comms via NATS | *status* |

#### 7.2 Decision Log

| ID | Decision | Date | Rationale | Status |
|----|----------|------|-----------|--------|
| ADR-001 | *decision* | *date* | *why* | Accepted / Superseded |

#### 7.3 Risk Register

| # | Risk | Probability | Impact | Mitigation | Owner |
|---|------|-------------|--------|------------|-------|
| 1 | *description* | H/M/L | H/M/L | *action* | *role* |

**Ask**: which one number breaks the design if it is wrong by ten times? Which assumptions does the
design rest on, and what changes if one is false?

---

### 8. Architecture Change Management (Phase H)

#### 8.1 Change Triggers
What events would trigger re-architecture:
- Business model change
- 10x scale requirement
- New compliance regulation
- Technology end-of-life

#### 8.2 Change Process
How architecture changes flow: proposal -> AIC -> review -> approval -> implementation.
Reference Policy of Initiatives (POL-TECH-001) stages.

#### 8.3 Technical Debt Register

| # | Debt Item | Impact | Priority | Remediation Plan |
|---|-----------|--------|----------|-----------------|
| 1 | *description* | *effect on quality* | H/M/L | *plan* |

#### 8.4 Open questions

Every answer the interview deferred. A named unknown is fine; a hidden one is not.

| Question | Who answers | By when (phase) |
|----------|-------------|-----------------|
| *question* | *name or role* | *Genesis / Custom / Product* |

---

## ArchiMate Diagram Conventions

### Color Scheme

| Layer | Color | Hex |
|-------|-------|-----|
| Business | Pale yellow | `#FFFFB5` |
| Application | Pale cyan | `#B5FFFF` |
| Technology | Pale green | `#C9E7B7` |
| Motivation | White | `#FFFFFF` |
| Strategy | Light pink | `#F5DEDC` |

### DOT Template

```dot
digraph ArchiMate_Layered {
    graph [
        label="ArchiMate: <Diagram Title>"
        labelloc=t fontsize=16 fontname="Arial"
        rankdir=TB pad=0.4
    ]
    node [shape=box style="filled,rounded" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    // Business Layer
    subgraph cluster_business {
        label="Business Layer" style=filled fillcolor="#FFFFB520"
        actor1 [label="<<Business Actor>>\nPartner" fillcolor="#FFFFB5"]
        process1 [label="<<Business Process>>\nProcess Name" fillcolor="#FFFFB5"]
        service1 [label="<<Business Service>>\nService Name" fillcolor="#FFFFB5" shape=component]
    }

    // Application Layer
    subgraph cluster_application {
        label="Application Layer" style=filled fillcolor="#B5FFFF20"
        comp1 [label="<<App Component>>\nService Name" fillcolor="#B5FFFF"]
        comp2 [label="<<App Component>>\nService Name" fillcolor="#B5FFFF"]
        data1 [label="<<Data Object>>\nEntity Name" fillcolor="#B5FFFF" shape=note]
    }

    // Technology Layer
    subgraph cluster_technology {
        label="Technology Layer" style=filled fillcolor="#C9E7B720"
        node1 [label="<<Node>>\nK8s Cluster" fillcolor="#C9E7B7"]
        sw1 [label="<<System Software>>\nPostgreSQL" fillcolor="#C9E7B7"]
        artifact1 [label="<<Artifact>>\nDocker Image" fillcolor="#C9E7B7"]
    }

    // Cross-layer relationships
    actor1 -> process1 [label="triggers"]
    process1 -> service1 [label="realizes"]
    service1 -> comp1 [label="realizes" style=dashed]
    comp1 -> comp2 [label="serves"]
    comp1 -> data1 [label="accesses" style=dashed]
    comp1 -> node1 [label="deployed on" style=dashed]
    sw1 -> node1 [label="assigned to"]
}
```

### Render Command

All diagrams go in `images/` inside the initiative folder, in dual format (`.dot` source + `.png` compiled):

```bash
dot -Tpng images/<name>.dot -o images/<name>.png
```

Embed in documents: `![Title](images/<name>.png)`

---

## Relationship to Other Documents

| TOGAF Section | Maps to AIC Field | Maps to Arc42 Section |
|--------------|-------------------|----------------------|
| 1. Architecture Vision | Business Case + Quality Goals | 1. Introduction & Goals |
| 2. Business Architecture | Business Context | 3.1 Business Context |
| 3. IS Architecture | Architectural Hypotheses | 5. Building Block View |
| 4. Technology Architecture | (Technical Constraints) | 7. Deployment View |
| 5. Opportunities & Solutions | Technical Challenges & Risks | 11. Risks & Tech Debt |
| 6. Migration Planning | Functional Overview (stages) | - |
| 7. Implementation Governance | - | 9. Architecture Decisions |
| 8. Change Management | Open questions | 11. Risks & Tech Debt |
