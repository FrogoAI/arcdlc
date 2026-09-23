# AIC Template

**Reviewed**: 2026-09-23

Each section says what to write, then lists the questions `/arcdlc:aic` asks to fill it. The
questions are the interview agenda, so an engineer thinking about one side of the design does not
lose the others. Every section ends one of three ways, and the engineer chooses: answered;
deferred, with an entry under Open questions naming who answers it and by which phase; or one line,
`Not applicable: <reason>`. Never delete a section: a reader cannot tell a missing answer from an
answer of "no". The Ask lines guide the interview and are not copied into the document.

## [Goals](https://canvas.arc42.org/architecture-inception-canvas)

### 🟢 Business Case

Brief description of the business case or economic driver behind the software system.

Ask:

- What does the problem cost today, in money, time or risk, and where does that number come from?
- What does solving it save or earn?
- Who decides that the system works: an engineer, an analyst, or the invoice?

### 🟢 Functional Overview

The most important functional requirements at a high level: the 5 to 10 that shape the architecture,
not the backlog.

Ask:

- What does the system do, in one sentence, without naming a technology?
- What is in the first stage, what is left out on purpose, and what does leaving it out cost?
- What condition opens the next stage?

### 🟢 Quality Goals

The three most essential [quality goals](https://docs.arc42.org/section-1/) for the architecture, which have the highest priority for the most important stakeholders.
Take each goal from the ISO 25010 list and give it a scenario that can be measured.

Ask:

- How many requests or user actions per second today, how many at the design target, and where do
  both numbers come from?
- What is the end-to-end latency target, split per step, and which step is the slowest?
- What measurement proves each goal is met, and what number counts as a pass?

### 🟢 Organizational Constraints

Any organizational requirement that limits the software architects' freedom of decision.

Ask:

- Who owns the system, and what are the budget, schedule, team size and partners?
- Which constraint, if lifted, would change the design most, and who can lift it?

### 🟢 Technical Constraints

Any technical requirement that restricts the software architects' freedom of decision (at least
follow `Engineering Principles.md` in the policy skill's references: `../../policy/references/`,
flat installs `../../arcdlc-policy/references/`).

Ask:

- Which platform, language, transport and stores already run and cannot be changed here?
- Which conventions does the work inherit: API versioning, naming, Twelve-Factor?
- Which of these constraints can be challenged, and by whom?

### 🟢 Business Context

Keep your system under construction as a black box, separate from all its communication partners. Communication partners include neighboring external systems and users (including charts, if needed).

Ask:

- Who and what talks to the system, and what data flows in and out of each, in domain words?
- Which protocol, format, endpoint or subject carries each flow?
- What must each partner change on its side, and how much work is it?
- Which partners can an attacker control? Which can be slow, or absent, without warning?

## [Architectural Hypotheses](https://canvas.arc42.org/architecture-inception-canvas)

### 🔵 Architectural Hypotheses

Resulting in architectural hypotheses and essential, expensive, large-scale, or risky architectural decisions, including justifications (include charts if needed).
Number every hypothesis and include: Context, Decision, Justification, and Trade-offs.

Ask, one theme at a time:

- **Structure.** Which language, stores, transport and split into services or modules? Why do the
  boundaries sit there, and why not the obvious other split? Does each new service have a reason
  that a package inside an existing service cannot meet?
- **Behaviour.** How is data saved, and how is it found again, told on concrete values? What happens
  when a step fails, times out, or arrives twice, and what does the caller receive?
- **Numbers.** How many calls to other services, and how many store reads and writes, per user
  action? How many bytes per stored item, how many items alive at once, and for how long? What is
  the chosen store's throughput limit, where was it measured, and what happens when it is reached?
  What does it cost per month?
- **Operation.** Where does each part run, and which environments differ? When a dependency is slow,
  and when it is gone, does the system fail open or closed? How does it go live, and how fast does
  it go back? Which metrics and alert values tell the person woken at 03:00 what broke?
- **Security.** Who may call what, how is the caller authenticated, what is encrypted, and which data
  is personal?
- **Decisions.** Which options lost, and by what criterion? Where does a rejected option stay
  reachable later without a rewrite?

## [Assessment](https://canvas.arc42.org/architecture-inception-canvas)

### 🔴 Technical Challenges & Risks

Identified current known challenges and technical risks, highest first, each with probability, impact
and the measure that reduces it.

Ask:

- Which one number breaks the design if it is wrong by ten times?
- What is left rough on purpose, and what will it cost to fix?
- Which assumptions does the design rest on, and what changes if one is false?

### Open questions

Every answer the interview deferred. A named unknown is fine; a hidden one is not.

| Question | Who answers | By when (phase or stage) |
| --- | --- | --- |
