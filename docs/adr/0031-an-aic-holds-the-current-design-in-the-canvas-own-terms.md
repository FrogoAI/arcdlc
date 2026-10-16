# ADR-0031 — An AIC holds the current design, in the canvas's own terms

- Status: Accepted
- Date: 2026-10-16
- Initiative: none (bundle-wide, changes the AIC template, `/arcdlc:aic` and the writing standard)

## Context

Engineers said AICs were hard to read: very technical, reasons missing where they were needed,
Functional Overview full of undecided or too-detailed points, Quality Goals padded with text, and
decisions that read like a robot filled in a form. Comparing the template with the
[Architecture Inception Canvas](https://canvas.arc42.org/architecture-inception-canvas) it claims to
follow showed where it drifted:

- The canvas asks for the product box in Functional Overview: the business value and the
  higher-level business functions. The template asked about stages and conditions instead.
- The canvas asks for the three most important quality goals, named. The template asked for a
  measured scenario and per-step latency for each, which pulled numbers into Goals that only the
  design can produce.
- The canvas's middle part is Solution, of which hypotheses are the start. The template had only
  hypotheses, so the whole design was squeezed into numbered blocks of Context, Decision,
  Justification and Trade-offs.
- `/arcdlc:aic` told the agent to note each reversal and its reason in the document body, so a
  finished AIC carried its own interview history.
- The writing standard banned any "narrative line", which also banned walking through one concrete
  case, the clearest way to explain a mechanism.

## Decision

The AIC follows the canvas's three parts: Goals, Solution, Assessment.

- **Goals are short.** Functional Overview is the product box, in business words, with nothing
  undecided. Quality Goals are three attribute names from ISO/IEC 25010 or Q42, highest first; one
  gets a line of reason only when it is not plain and the architect asks for it. A quality goal is an
  input and is replaced in place when the design shows it was not the driver.
- **Solution holds the design.** Architectural Hypotheses are the few large-scale assumptions, two to
  five, each in whatever form reads best. Solution Design, a new section, holds the design built on
  them: diagrams, worked examples, algorithms, the maths behind a limit, API changes, delivery order.
  Every choice there carries its reason unless any engineer would know it.
- **The body holds the current design only.** A reversed hypothesis leaves the body; its reason goes
  to an ADR when it would stop the next reader re-opening it, otherwise to the commit message. A
  rejected option stays as one line only where a reader would ask about it. This amends the
  body-text rule in `/arcdlc:aic` Step 3; ADR supersession is unchanged.
- **A fresh reader tests the draft.** Before handover a subagent with only the document reports what
  it could not understand, missing reasons, water, and unreadable diagrams. Without a subagent the
  agent does that read itself.
- **The writing standard allows a worked example** and adds three sections: write prose, not a form;
  cut the water; give the reason when it is not plain. The shared block in all twelve `SKILL.md`
  files changes with it.
- `references/AIC Example.md` shows a finished AIC, adapted with names changed from an initiative an
  engineer chose as the best of eight.

## Consequences

Easier: a product owner can read Goals in a minute, an engineer finds the design in one section with
pictures, and a finished AIC no longer argues with its own past.

Harder: an architect who wants the history in the document has to open the ADRs or git log. Solution
Design has no fixed shape, so two AICs look less alike, which is the point but makes them harder to
compare section by section. Existing AICs keep the old headings until their next `/arcdlc:aic` run;
nothing parses those headings, so nothing breaks.
