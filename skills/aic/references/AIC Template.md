# AIC Template

**Reviewed**: 2026-10-16

The AIC follows the [Architecture Inception Canvas](https://canvas.arc42.org/architecture-inception-canvas) by arc42. It has three parts: **Goals** (what should the software do), **Solution** (how we achieve it) and **Assessment** (how we judge where we stand). A finished AIC is the one source of truth for the architect, the engineer and the product owner. [AIC Example.md](AIC Example.md) shows a finished one.

Each section says what to write, then lists the questions `/arcdlc:aic` asks to fill it. The questions are the interview agenda, so an engineer thinking about one side of the design does not lose the others. Every section ends one of three ways, and the engineer chooses: answered; deferred, with an entry under Open questions naming who answers it and by which phase; or one line, `Not applicable: <reason>`. Never delete a section: a reader cannot tell a missing answer from an answer of "no". The Ask lines guide the interview and are not copied into the document.

## How a finished AIC reads

- **A picture before text.** Where a diagram can show it, draw it, and let the text add only what the picture cannot: the reason, a number, an edge case. Three diagrams are expected in every AIC: a C4 system context diagram in Business Context, a C4 container diagram in Solution Design, and a flowchart or a BPMN diagram of the main process in Solution Design. BPMN is the right choice when the process passes between people or systems. A missing expected diagram needs a `Not applicable: <reason>` line, like a missing section. How to draw them: [Diagram Conventions.md](Diagram Conventions.md) and [C4.md](C4.md).
- **It holds the current design only.** A rejected option, a superseded hypothesis and the history of the interview are not in the body. They live in the ADRs and in git history. The body may name a rejected option in one line, when a reader would otherwise ask "why not X?".
- **Goals are short.** The Goals part states the problem in the words of the business, as the canvas asks. Reasoning, numbers and mechanics belong in Solution.
- **Solution choices carry their reason.** Every choice in Solution says why, unless the reason is plain to any engineer on the team.
- **No water.** No sentence that restates its heading, copies a rule the repository already states elsewhere (link it instead), explains a choice nobody would question, or narrates how the interview went.

## [Goals](https://canvas.arc42.org/architecture-inception-canvas)

### 🟢 Business Case

The economic driver: why someone pays for this. Every team member should be able to explain it in two sentences. A few sentences, with a number when one exists.

Ask:

- What does the problem cost today, in money, time or risk, and where does that number come from?
- What does solving it save or earn?
- Who decides that the system works: an engineer, an analyst, or the invoice?

### 🟢 Functional Overview

What would be printed on the product box: the business value the software delivers and its higher-level business functions, in three to ten bullets a product owner would sign. One line, `Not on the box:`, may name what is left out on purpose.

It does not hold endpoints, field names, data structures, algorithms, UI widgets, or anything still undecided. Mechanics go to Solution; an undecided point goes to Open questions. The test for each line: would the person paying for it understand it and care?

Ask:

- What business value does the software deliver?
- What are the higher-level business functions, the ones you would print on the box?
- What is left off the box on purpose, for now?

### 🟢 Quality Goals

The three most important quality goals of the most important stakeholders, highest first. Name each with one attribute from [ISO/IEC 25010](https://iso25000.com/index.php/en/iso-25000-standards/iso-25010) or the [arc42 quality model (Q42)](https://quality.arc42.org/), in one to three words. Every team member should know them, because they shape the architecture.

No explanation by default. Add one line of reason to a goal only when the choice is not plain and the architect asks for it. Numbers, scenarios and measurements belong in Solution Design, where the design proves them.

```markdown
1. Usability
2. Safety: a wrong branch silently changes production decisions.
3. Maintainability
```

A quality goal is an input, and the design may show it was the wrong one. A goal picked as "performance" on first look is not a driver once the design finds 17 requests per second and one database query each. Replace it in place, then ask about each decision that rested on it.

Ask:

- Which three qualities, if they failed, would hurt the main stakeholder most? Offer your three.
- Is any of the three not plain to the team, so that it needs one line of reason?

### 🟢 Organizational Constraints

What limits the freedom to decide, from outside the code: time, money, team, partners, and a required method or process (for example Scrum or SAFe). Bullets. A reason only when it is not plain.

Ask:

- Who owns the system, and what are the budget, schedule, team size and partners?
- Which constraint, if lifted, would change the design most, and who can lift it?

### 🟢 Technical Constraints

Technical limits the design must accept: mandated platforms, languages and stores, and the organization's own engineering principles. The principles bind first: its policies in `docs/policies/`, or the standard the project's `AGENTS.md` or `README.md` links to. Name the document and its ID and link it; do not copy its rules here. Only when the organization has none, the worked example `Engineering Principles.md` in the policy skill's references (`../../policy/references/`, flat installs `../../arcdlc-policy/references/`) is a starting point, and only if the engineer adopts it.

Ask:

- Which engineering principles or standards does the organization hold this work to, and where are they written down?
- Which platform, language, transport and stores already run and cannot be changed here?
- Which conventions does the work inherit: API versioning, naming, Twelve-Factor?
- Which of these constraints can be challenged, and by whom?

### 🟢 Business Context

The system as a black box among its users and neighbouring systems, and what flows between them in domain words. Draw it as a C4 system context diagram (`images/c4_context.png`). Under the diagram, write a line for a partner only where the picture leaves something out: what the partner must change on its side, or that it can be slow, absent, or controlled by an attacker.

Ask:

- Who and what talks to the system, and what data flows in and out of each, in domain words?
- Which protocol, format, endpoint or subject carries each flow?
- What must each partner change on its side, and how much work is it?
- Which partners can an attacker control? Which can be slow, or absent, without warning?

## [Solution](https://canvas.arc42.org/architecture-inception-canvas)

How the goals are met. It starts with the few hypotheses the design rests on, then shows the design built on them in as much depth as an engineer needs to build it and a product owner needs to trust it.

### 🔵 Architectural Hypotheses

The large-scale assumptions and choices the design rests on, from what the team knows today: expensive to reverse, risky, or shaping everything else. Usually two to five. Number them H1, H2, and so on. Each says what is assumed or chosen, why, and what it costs, in the form that reads best for that one: a sentence, a short paragraph, or a table. There is no fixed set of labelled fields. Link the ADR of a hard one. Only hypotheses that stand are listed: a rejected or superseded one leaves the document, and its record is the ADR.

Ask:

- What does the design assume that, if it turned out false, would change it?
- Which choice is expensive to reverse, and what did it win against?

### 🔵 Solution Design

The design itself: how it works, why, what changes, where, and when. Lead with the C4 container diagram and the flowchart or BPMN diagram of the main process. Then add what this initiative needs to be understood, for example:

- a worked example with real values,
- the algorithm, or the maths behind a limit or a size,
- a state machine, a sequence or a data model diagram,
- what changes in each service, and the API changes, in a short table,
- the numbers that prove the quality goals,
- the order of delivery: stages, migrations, feature switches, how to roll back.

Say what is chosen, then why. A rejected alternative gets one line only when a reader would ask about it; the full comparison is in the ADR. The order of delivery is an outline, not a task list: `/arcdlc:plan` turns it into tasks.

Ask, one theme at a time:

- **Structure.** Which language, stores, transport and split into services or modules? Why do the boundaries sit there, and why not the obvious other split? Does each new service have a reason that a package inside an existing service cannot meet?
- **Behaviour.** How is data saved, and how is it found again, told on concrete values? What happens when a step fails, times out, or arrives twice, and what does the caller receive?
- **Numbers.** How many calls to other services, and how many store reads and writes, per user action? How many bytes per stored item, how many items alive at once, and for how long? What is the chosen store's throughput limit, where was it measured, and what happens when it is reached? What does it cost per month?
- **Operation.** Where does each part run, and which environments differ? When a dependency is slow, and when it is gone, does the system fail open or closed? How does it go live, and how fast does it go back? Which metrics and alert values tell the person woken at 03:00 what broke?
- **Security.** Who may call what, how is the caller authenticated, what is encrypted, and which data is personal?
- **Delivery.** In what order does it ship, what has to migrate, and how is a stage rolled back?

## [Assessment](https://canvas.arc42.org/architecture-inception-canvas)

### 🔴 Technical Challenges & Risks

Known challenges and technical risks, highest first, each with likelihood, impact and the mitigation the engineer chose. A table reads best once there are more than two.

| Risk | Likelihood | Impact | Mitigation |
| --- | --- | --- | --- |

Ask:

- Which one number breaks the design if it is wrong by ten times?
- What is left rough on purpose, and what will it cost to fix?
- Which assumptions does the design rest on, and what changes if one is false?

### Open questions

Every answer the interview deferred. A named unknown is fine; a hidden one is not.

| Question | Who answers | By when (phase or stage) |
| --- | --- | --- |
