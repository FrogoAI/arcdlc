# Policy of Policies (POL-GOV-001)

**Reviewed**: 2026-10-22

## Purpose

This policy sets how every policy in the organization is written, approved, published, changed, and retired. Without it, each policy takes its own shape, nobody can tell which version is in force, and two policies can say different things about the same rule.

## Policy Statement

Every policy, procedure, and governance document is created, reviewed, approved, published, changed, and retired as this policy says. A policy or a change that skips these steps is not in force. This policy replaces every earlier, informal, or department-specific way of making policies.

## Scope

**Personnel covered**

Employees and contractors, unless a policy says otherwise. It bears most directly on anyone who writes, reviews, approves, publishes, or owns a policy.

**Documents covered**

* **Vision and values documents:** the company's mission and ethical principles, for example a code of conduct.
* **Corporate policies:** organization-wide rules on strategic areas.
* **Functional and operational policies:** rules for one department or function.
* **Procedures and SOPs:** step-by-step instructions that carry out a policy.
* **Supporting documents:** forms, guidelines, and checklists that support a procedure.

**Precedence over bundled references**

A policy the organization writes and approves under this framework outranks any reference document shipped with a tooling bundle (for ArcDLC: the `references/` folders inside each skill). Bundled references are defaults for subjects the organization has not decided yet. Where an approved internal policy covers the same subject, it is the normative rule for audits, reviews, and tooling, and the bundled reference applies only to what the internal policy leaves unspecified.

**Exclusions**

This policy does not replace laws, regulations, or external standards. Every policy written under it complies with them.

## Definitions

* **Policy:** a rule that says what we do and why.
* **Procedure:** step-by-step instructions on how to carry out a policy.
* **Policy Owner:** the senior leader accountable for a policy's content, relevance, and review.
* **Policy Approver:** the C-level person with the final authority to approve a policy. Every policy names one.
* **Policy Reviewer:** a stakeholder who reviews a draft for practicality and clarity. Reviewers are listed in the policy's Stakeholders table.
* **Policy Writer:** the subject matter expert who writes and revises the text.
* **Policy Administrator:** the role that keeps the index and reminds Policy Owners when a review is due.
* **Editorial change:** a change that leaves every rule's meaning as it was: wording, typos, dates, links, formatting, role names.
* **Rule change:** any change that adds, removes, or changes what a rule requires or allows.
* **Accountable (A):** the one person who answers for the result of a step.
* **Responsible (R):** the people who do the step.
* **Consulted (C):** the people whose input is needed before the step is done.
* **Informed (I):** the people told about the step's progress or result.

Other policies refer to these definitions instead of repeating them.

## Procedures

### Stage 1: Need

A new or changed policy starts when someone sees a gap, a new risk, a changed regulation, or a rule that no longer fits. The Policy Writer collects the facts and agrees the objective and the scope with the stakeholders.

### Stage 2: Drafting

The Policy Writer drafts the policy from the official template (Policy layout and Policy Template below). The text is plain and direct, in active voice, with a named actor for every rule.

### Policy layout

Every policy has these parts, in this order. Sections are not numbered.

1. **Title line:** `# <Policy Name> (<Unique ID>)`.
2. **Header table** with these fields: Policy Name, Unique ID, Status, Policy Owner, Policy Approver, Author, Creation Date, Effective Date, Approval Date, Next Review Date, Pending changes (only while a rule change waits for its reviewer window), and Related documents.
3. **Stakeholders:** a table with Role, Person, Review Status, and Review Date. It lists every Policy Reviewer and the Policy Approver.
4. **Purpose:** the problem, risk, or gap the policy addresses.
5. **Policy Statement:** the binding rule in one or two short paragraphs, and a note that it replaces informal earlier practice.
6. **Scope:** personnel covered, documents and systems covered, and exclusions.
7. **Definitions:** the terms the policy uses. Policy roles and RACI letters mean what this policy says.
8. **Procedures:** the steps or stages the policy governs, one sub-heading each.
9. **Roles & Responsibilities:** the Policy Owner, Policy Approver, Policy Writer, and Policy Reviewers, then a RACI table for the key steps.
10. **Allowed & Prohibited Conduct:** concrete rules that can be checked, one rule per bullet.
11. **Consequences of Non-Compliance:** what happens after a breach, and who decides.
12. **Open questions:** only when there are any, as a table with the question, the owner, and the date it is needed by.
13. **Revision History:** a table with Version, Date, Author, and Summary of changes.

Rules for every policy:

* All dates are written as `YYYY-MM-DD`.
* Rule text, RACI tables, and the header use roles, never people's names. A rule bound to a role survives a hire or a leaver. People's names appear only in Stakeholders and Revision History, because those two record who reviewed, approved, or changed the policy, and when.
* The version of a policy is the Version in its Revision History (v1.0, v1.1, ...). Every change adds a row.
* A cross-reference to another policy names its Unique ID and the sub-section title, for example "POL-ENG-005, Procedures: Lifecycle".
* The scope covers employees and contractors unless the policy says otherwise.
* Consequences of Non-Compliance point to the organization's conduct and disciplinary policy for people. For a contractor, a breach leads to the measures the contract allows, not employment discipline.

### Policy Approvers

Every policy names its Policy Approver in the header and in Stakeholders, and the index lists it. The Policy Approver is always a C-level person. The organization decides once which policies the CEO approves, typically governance (this policy), HR and pay, and top-level security, and writes that list here. The CTO, or another C-level owner of the area, approves the rest.

### Stage 3: Review and approval

The Policy Reviewers listed in Stakeholders review the draft, and the Policy Approver approves it.

* A policy cannot become `Active` until its Policy Approver has approved it.
* After the Policy Approver approves, the reviewers have two weeks. If the reviewers have not reviewed it after more than two weeks, the policy becomes `Active` anyway.
* Silence from a reviewer is recorded as `Need Review` in Stakeholders. It is never recorded as an approval.

### Stage 4: Publishing

The official copy of every policy is in the knowledge base, together with the policy index. The index lists every policy with its Unique ID, title, status, owner, approver, next review date, and pending changes.

* The editable source of every policy (for ArcDLC: `docs/policies/<name>.md`, indexed in `docs/policies/README.md`) is published to the knowledge base on the same day it changes.
* The Policy Owner tells the people in the policy's scope that the policy is published or changed.
* A policy that is not in the index is not in force.

### Stage 5: Maintenance and review

* The Policy Owner reviews every policy every six months. A review that changes nothing still gets a Revision History row and a new Next Review Date.
* Every change gets a new Version and a Revision History row that says what changed and why.
* An **editorial change** needs only the Policy Approver's approval. The policy stays `Active`, and the change applies at once.
* A **rule change** needs the Policy Approver's approval and the reviewers' two-week window from Stage 3. The change takes effect when the window ends. Until then, the previous version applies: the header's Pending changes row says so, and the previous text stays available to every reader.
* When a new policy's first two-week window has not ended yet, it still runs, also for any change made during it.

### Stage 6: Retirement

When a policy is no longer needed or another policy replaces it, the Policy Owner and the Policy Approver retire it.

* A retired policy keeps its file as a record, with status `Retired` and the policy that replaces it.
* The index lists it under retired policies.
* Its Unique ID is never reused.

### Unique ID classification

An ID stays with its policy for life. Take the next free number in the class.

| Identifier | Description |
| :---- | :---- |
| POL-GOV | **Governance Policy.** Overall corporate or organizational governance. |
| POL-HR | **Human Resources Policy.** Principles, standards, and procedures for managing the workforce. |
| POL-TECH | **Technology Policy.** Using, developing, and managing technology systems. |
| POL-ENG | **Engineering Policy.** Engineering practices, standards, and procedures. |
| POL-SEC | **Security Policy.** Protecting information, infrastructure, and systems. |
| POL-DOC | **Documentation or Records Management Policy.** Creating, keeping, storing, and disposing of documents and records. |
| POL-INT | **Integration Policy.** Onboarding partners and connecting external systems. |

## Roles & Responsibilities

* **Policy Owner:** accountable for a policy's content and its six-monthly review. Publishes changes and tells the people in scope.
* **Policy Writer:** writes and revises the text.
* **Policy Reviewer:** reviews drafts and changes within the two-week window.
* **Policy Approver:** a C-level person who approves the policy and every change to it.
* **Policy Administrator:** keeps the index and reminds Policy Owners when a review is due.

Each policy maps its own key steps to these roles in a RACI table (see the template below).

## Allowed & Prohibited Conduct

### Allowed

* Draft a new policy from the official template.
* Make an editorial change with the Policy Approver's approval only.
* Propose a change to any policy to its Policy Owner.

### Prohibited

* Use an outdated or unofficial template for a governance document.
* Create, enforce, or distribute a policy that has not been approved as this policy says.
* Treat a reviewer's silence as approval.
* Apply a rule change before its reviewer window has ended.
* Keep the official copy of a policy anywhere other than the knowledge base, or leave it out of the index.
* Put people's names into rule text, RACI tables, or the header.

## Consequences of Non-Compliance

A policy, or a change, that skipped these steps is not in force: the previous approved version applies until the steps are done. Breaches by people are handled under the organization's conduct and disciplinary policy. For a contractor, a breach leads to the measures the contract allows.

## Policy Template

This is the official Policy Template that Stage 2 mandates. Copy the block below into `docs/policies/<name>.md`; it is the entire file, and nothing goes above the H1. Every section is mandatory except Open questions: a section that does not apply is still written, stating why it does not apply, rather than being dropped.

  ```markdown
  # <Policy Name> (<POL-CLASS>-NNN)

  | Field | Value |
  | :---- | :---- |
  | Policy Name | <Policy Name> |
  | Unique ID | <POL-CLASS>-NNN |
  | Status | Draft |
  | Policy Owner | <role> |
  | Policy Approver | <C-level role, see POL-GOV-001, Procedures: Policy Approvers> |
  | Author | <role> |
  | Creation Date | YYYY-MM-DD |
  | Effective Date | <YYYY-MM-DD, or "After <approver> approval"> |
  | Approval Date | TBD |
  | Next Review Date | <Creation Date + 6 months> |
  | Pending changes | <only while a rule change waits for its reviewer window; delete the row otherwise> |
  | Related documents | <POL-ID Title, ...; external standards by name> |

  ## Stakeholders

  Everyone below reviews this policy. The <approver> approves it.

  | Role | Person | Review Status | Review Date |
  | :---- | :---- | :---- | :---- |
  | <role> | <name> | Need Review | |

  ## Purpose

  Why this policy exists: the problem, risk, or gap it addresses.

  ## Policy Statement

  The binding rule, in one or two short paragraphs, plus a note that it replaces informal earlier
  practice on this topic.

  ## Scope

  Personnel covered (employees and contractors, unless the policy says otherwise), documents and
  systems covered, and explicit exclusions.

  ## Definitions

  Policy roles and RACI letters mean what POL-GOV-001 says. Define here only the other terms this
  policy uses.

  ## Procedures

  The steps or stages this policy governs, one sub-heading each. Other policies refer to them as
  "<POL-ID>, Procedures: <sub-heading>".

  ## Roles & Responsibilities

  Name the Policy Owner, Policy Approver, Policy Writer, and Policy Reviewers by role, then map the
  policy's key steps to them:

  | Rule/Step | A | R | C | I |
  | :---- | :---- | :---- | :---- | :---- |
  | <key step> | <accountable> | <responsible> | <consulted> | <informed> |

  ## Allowed & Prohibited Conduct

  ### Allowed

  * <Concrete, auditable behavior this policy permits or requires.>

  ### Prohibited

  * <Concrete, auditable behavior this policy forbids.>

  ## Consequences of Non-Compliance

  What happens when the policy is breached, built into the workflow first. Then: breaches by people
  are handled under the conduct and disciplinary policy; for a contractor, the contract measures.

  ## Open questions

  Only when there are any; delete the section otherwise.

  | # | Question | Owner | Needed by |
  | :---- | :---- | :---- | :---- |

  ## Revision History

  | Version | Date | Author | Summary of changes |
  | :---- | :---- | :---- | :---- |
  | v1.0 | YYYY-MM-DD | <name> | Initial version. |
  ```

Fill `<POL-CLASS>` from the Unique ID classification table above, and `NNN` with the next free number in that class.
