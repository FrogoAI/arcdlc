# Policy of Initiatives (POL-TECH-001)

## Purpose

The purpose of this policy is to establish a standardized, transparent, and efficient framework for proposing, defining, architecting, and managing all new business and technical initiatives. By mandating the use of the Architecture Inception Canvas (AIC) template, this policy ensures that every initiative is consistently documented, architecturally sound, and aligned with strategic goals before development work begins. This process is designed to foster collaboration, mitigate risk, and ensure clarity and accountability throughout the initiative lifecycle.

---

## Policy Statement

All new initiatives, regardless of their origin (Business or Technical), must be documented using the official AIC Template. No initiative shall be considered valid or approved for development until it has progressed through the procedures outlined in this document. This policy is the single source of truth for the initiative management lifecycle and supersedes any previous informal processes.

---

## Scope

This policy applies to all full-time employees across all departments and locations involved in the conception, design, approval, and implementation of new features, systems, or technical improvements. It governs the entire lifecycle of all documents related to initiatives, including the mandatory AIC and any subsequent proof-of-concept (PoC) or design artifacts.

---

## Definitions

* **Initiative**: A formal proposal for creating a new feature, system, or significant technical change. All initiatives are categorized as either a "Standard Initiative" or a "Technical Initiative." An initiative lives in its own folder, `docs/aics/<slug>/`, which holds its architecture document and the task plan derived from it.
* **Architecture Inception Canvas (AIC)**: The mandatory document, based on the Arc42 framework, used to describe and define an initiative. It is written to `docs/aics/<slug>/aic.md`.
* **Standard Initiative**: An initiative with a direct business impact, originating from either Business or Technical teams. It requires the full AIC documentation.
* **Technical Initiative**: A small-scale initiative with no direct business impact (e.g., refactoring, service renaming). It follows a simplified documentation process.
* **Initiator**: The individual who first proposes an initiative and is responsible for its initial documentation.
* **Architect**: The individual(s) responsible for defining the architectural solution. The CTO, Software Architect, or Cloud Architect fulfills this role.
* **Initiative Stages**:
  * **Genesis**: The initial phase, covering the creation of the AIC and an optional Proof of Concept (PoC).
  * **Custom**: The development of core, production-grade functionality (MVP).
  * **Product**: The implementation of the full-featured initiative as designed.
* **Business Context**: A schema definition of which part of the system will be affected by changes. Can include text description and schema/diagram. More details on how to write Business Context can be found in the [arc42 official documentation](https://docs.arc42.org/section-3/).

---

## Procedures

### Stage 1: Initiative Initiation

1. **Standard Initiative Initiation**:
   * An initiative can be proposed by authorized personnel from either the Business or Technical teams.
     * **Business Initiators**: Engineering Manager, Product Owner, Delivery Manager, System Analyst, Business Analyst.
     * **Technical Initiators**: CTO, Product Owner, Software Architect, Cloud Architect, Team Lead, and Engineering Manager.
   * The Initiator must create a new initiative folder, `docs/aics/<slug>/`, and write its architecture document using the official AIC Template.
   * The Initiator is responsible for completing the "Collect Business Requirements" sections, including the **Business Case**, **Functional Overview**, and **Quality Goals**. Quality Goals must be selected from the [official Arc42 list](https://docs.arc42.org/section-1/#12-quality-goals); custom or undefined Quality Goals are not permitted.

2. **Technical Initiative Initiation**:
   * Any team member may propose a Technical Initiative.
   * The Initiator must use the AIC Template, but is only required to complete: **Header (Meta)**, **Business Context**, **Architectural hypotheses**, and **Tasks**.
   * The Business Case section must be filled with: "*No impact*."

### Stage 2: Architectural Definition

1. Upon initiation, the initiative is assigned to an Architect (CTO, Software Architect, or Cloud Architect).
2. The Architect is responsible for completing the technical sections of the AIC, with a primary focus on **Architectural Hypotheses** and **Technical Challenges & Risks**.
3. The Architect must supplement the 'Architectural Hypotheses' section with relevant diagrams. It is strongly recommended to use standard notations, such as the [C4 model](https://c4model.com/diagrams), [BPMN](https://en.wikipedia.org/wiki/Business_Process_Model_and_Notation), [UML](https://en.wikipedia.org/wiki/Unified_Modeling_Language), or [flowcharts](https://en.wikipedia.org/wiki/Flowchart), where applicable.
4. The Initiator may define the vision and initial proposal, but the final architectural decision rests with the assigned Architect.

### Stage 3: Review and Refinement

1. The completed AIC is presented for review to the development team(s) whose context the initiative touches, so that feedback and open issues surface before implementation. The Architect refines the document until that feedback is resolved.
2. For a Technical Initiative, the AIC must be reviewed and approved by at least one of: CTO, Software Architect, or Cloud Architect. This can occur asynchronously or in a meeting.

### Stage 4: Task Separation and Implementation

1. Once the review feedback is incorporated, the initiative is considered architecturally approved.
2. The approved architecture document is then decomposed into an executable task queue: `/arcdlc:plan <slug>` turns `docs/aics/<slug>/aic.md` into task blocks in `docs/aics/<slug>/plan.md`.
3. Every task block must carry its own Acceptance criteria, concrete enough for QA and for the executor to verify; the plan is the shared contract between architecture, implementation, and testing.
4. The work then proceeds through the defined Initiative Stages (Genesis, Custom, Product).
