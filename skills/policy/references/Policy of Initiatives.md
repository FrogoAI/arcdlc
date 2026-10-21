# Policy of Initiatives (POL-TECH-001)

**Reviewed**: 2026-10-22

## Purpose

This policy sets how a business or technical initiative is proposed, designed, approved, and handed to delivery. Every initiative gets the same documents and the same steps, so a design is written down and reviewed before anyone builds it, and one place always shows where the initiative stands.

## Policy Statement

Every new initiative, business or technical, is documented in an AIC (Architecture Inception Canvas) written from the official AIC Template. No initiative is approved for development until it has gone through the procedures in this policy. One architecture task in the task tracker carries the initiative's status from the start of design to verification. This policy replaces any earlier informal process.

## Scope

Employees and contractors who propose, design, approve, or build new features, systems, or technical improvements. It covers every document of an initiative: the AIC, its Meta, the Tech Stack Canvas when one is needed, and any proof of concept (PoC).

## Definitions

* **Initiative:** a formal proposal for a new feature, system, or significant technical change, either a Standard Initiative or a Technical Initiative. An initiative lives in its own folder, `docs/aics/<slug>/`, which holds its architecture document and the task plan derived from it. Where the organization keeps documentation in a knowledge base, that folder is published there as the initiative's space.
* **Architecture Inception Canvas (AIC):** the mandatory design document, based on arc42, written to `docs/aics/<slug>/aic.md` from the arcdlc AIC Template.
* **Meta:** a separate document next to the AIC, never part of it. It holds the initiative type, the Architect, the link to the architecture task, the links to the delivery epics, and the reviewers with their review status and date.
* **Architecture task:** one issue in the task tracker, titled with the initiative name. Its state is the only status of the initiative's architecture. Weekly status notes and the verification comment go on it.
* **Delivery epic:** an epic in the task tracker for part of the initiative, for example one per team or service. Every delivery epic links to the architecture task.
* **Standard Initiative:** an initiative with direct business impact, from the business or a technical team. It needs the full AIC.
* **Technical Initiative:** a small initiative with no direct business impact, for example a refactoring or a service rename. It follows a shorter process.
* **Initiator:** the person who proposes an initiative and writes its business sections.
* **Architect:** the person who designs the solution: the CTO, a Software Architect, or a Cloud Architect.
* **Initiative Stages:**
  * **Genesis:** the AIC, and a proof of concept when one is needed. A PoC or MVP that serves real traffic is Genesis, never a spike.
  * **Custom:** core production-grade functionality (MVP).
  * **Product:** the full initiative as designed.
* **Business Context:** the part of the system the change affects, in text and a diagram. See the [arc42 documentation](https://docs.arc42.org/section-3/).

## Procedures

### Initiation

1. **Standard Initiative:**
   * An initiative can be proposed from the business or from a technical team.
     * **Business Initiators:** Engineering Manager, Product Owner, Delivery Manager, System Analyst, Business Analyst.
     * **Technical Initiators:** CTO, Product Owner, Software Architect, Cloud Architect, Team Lead, Engineering Manager.
   * The Initiator writes the business sections of the AIC: **Business Case**, **Functional Overview**, and **Quality Goals**. Quality Goals come from the [official arc42 list](https://docs.arc42.org/section-1/#12-quality-goals); custom or undefined Quality Goals are not allowed.
   * Before an Architect is assigned, the Initiator may write these sections as a draft. The draft uses the AIC Template format, and the Architect moves it into the AIC.
2. **Technical Initiative:**
   * Any team member may propose a Technical Initiative.
   * The Initiator uses the AIC Template, but only **Business Context** and the **Solution** part (**Architectural Hypotheses** and **Solution Design**) are required.
   * The Business Case reads: `Not applicable: no direct business impact.` Every other section is answered, deferred, or marked not applicable as described in Architectural definition.

### Architectural definition

1. The initiative is assigned to an Architect (CTO, Software Architect, or Cloud Architect).
2. At the start, the Architect creates the architecture task, the initiative folder `docs/aics/<slug>/` with the AIC and the Meta, and links the Meta and the architecture task both ways. The Architect does not create the delivery epics.
3. The Architect writes the technical sections of the AIC, mainly the **Solution** part and **Technical Challenges & Risks**.
4. Architecture documents come in two formats only: the AIC and, when the initiative brings a new technology stack, a Tech Stack Canvas (`/arcdlc:aic <slug> tsc`).
5. The Architect draws the AIC's diagrams in standard notations. Three are expected in every AIC: a [C4](https://c4model.com/diagrams) system context diagram, a C4 container diagram, and a [flowchart](https://en.wikipedia.org/wiki/Flowchart) or [BPMN](https://en.wikipedia.org/wiki/Business_Process_Model_and_Notation) diagram of the main process. A missing one carries `Not applicable: <reason>`. [UML](https://en.wikipedia.org/wiki/Unified_Modeling_Language) and further diagrams are added where they explain the design better than text.
6. The Initiator may set the vision and the first proposal, but the final architectural decision is the Architect's.
7. Every section of the AIC ends one of three ways, chosen by the Architect: answered; deferred, listed under **Open questions** with who will answer it and by which Initiative Stage; or one line, `Not applicable: <reason>`. No section is deleted. `/arcdlc:aic <slug>` walks the sections as its interview and can be re-run at any stage to answer a deferred question, reopen a section, or add detail.

### Architecture task states

The architecture task's state is the initiative's status. The Architect moves it.

| State | Meaning |
| :---- | :---- |
| New | The Architect created the task on assignment and has not started yet. |
| Plan | The Architect sets up the initiative folder, the AIC, the Meta, and the reviewers. |
| Design | The Architect writes the AIC, with a PoC when the initiative is in Genesis. |
| Review | The team reviews the AIC, and the Architect presents it. |
| Development | Entering this state is the approval point. Epics and tickets are created, reviewed, and built. |
| Verified | The verification comment is posted, and the AIC is closed. |
| On Hold | Paused, with the CTO's agreement. The Architect writes the reason and the expected restart date in a comment. Deadlines and weekly status notes stop. The CTO reviews every On Hold task at least every 30 days. |

While the task is not On Hold, the Architect posts a status note on it every week, until it is Verified.

### Review and refinement

1. The Architect presents the completed AIC to the development team(s) whose context the initiative touches, so feedback and open issues surface before implementation. The Architect refines the AIC until that feedback is resolved.
2. For a Technical Initiative, at least one of the CTO, a Software Architect, or a Cloud Architect approves the AIC, in a meeting or asynchronously.

### Tasks and implementation

1. Once the review feedback is incorporated, the initiative is architecturally approved, and the architecture task moves to Development.
2. The Engineer creates the delivery epics and the tickets in the task tracker from the AIC and the Architect's recommendation, and links every epic to the architecture task. `/arcdlc:plan <slug>` turns `docs/aics/<slug>/aic.md` into task blocks in `docs/aics/<slug>/plan.md`, and `/arcdlc:plan-human <slug>` turns those into tickets. Tasks are never written into the AIC itself.
3. Every task carries its own Acceptance criteria, concrete enough for QA and for the executor to verify. The Test Manager writes the test tasks for each delivery epic.
4. Genesis code in a new, isolated service may merge without a peer code review. It still needs the Architect's review and green CI, including security scans. Before the initiative moves from Genesis to Custom, the Team Lead's reviewers do one full code review of that service.
5. The work then goes through the Initiative Stages (Genesis, Custom, Product). After the last ticket reaches production, the Architect posts the verification comment on the architecture task, listing every hypothesis as Confirmed or Rejected with evidence, and moves the task to Verified.

## Roles & Responsibilities

| Rule/Step | A | R | C | I |
| :---- | :---- | :---- | :---- | :---- |
| Write the business sections of the AIC | Product Owner | Initiator | Architect | Engineering Manager |
| Create the architecture task, the initiative folder, and the Meta | CTO | Architect | | Initiator, Product Owner |
| Write the technical sections of the AIC | CTO | Architect | Team Lead | Engineer |
| Move the architecture task between states and post weekly notes | CTO | Architect | | Product Owner, Team Lead |
| Create the delivery epics and tickets | Team Lead | Engineer | Architect | Test Manager |
| Write the test tasks | Test Manager | Test Manager | Architect, Team Lead | Engineer |
| Verify the built system against the AIC | CTO | Architect | Team Lead, Engineer | Product Owner, Test Manager |
