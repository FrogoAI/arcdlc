# ADR-0027 — The template is the interview agenda

- Status: Accepted
- Date: 2026-09-23
- Initiative: none (bundle-wide, changes the `/arcdlc:aic` interview)

## Context

`/arcdlc:aic` grilled the design against a fixed list of ten topics ("Cover at minimum: problem/goal,
scope boundaries, ...") and then filled the template. The list and the template were two views of
the same thing, and they did not match. The list never asked for load, fan-out, store limits, cost,
failure behaviour or alerts. The AIC template gave one line of hint per section, so nothing else
asked either.

An engineer designing a system thinks along one line, usually the one the request came from. The
numbers, the failure cases and the operating concerns are the parts that go missing, and they go
missing silently: a section with nothing in it reads the same as a section whose answer is "no".

A 60-item criteria list with a 0 to 3 score, checked after each run and triggering another
interview on a low score, was considered and rejected. The model that wrote the document would grade
its own work. Chasing a score would push it to invent sourced numbers where none exist, which breaks
"do not invent a decision". And a score puts the model in judgement over an architect's answer on a
complex design, which is the architect's call.

## Decision

**The sections of the chosen template are the interview agenda.** `/arcdlc:aic` hands the interview
every section of every requested template, in order. The questions come from each section's own
guidance: its description and new **Ask** lines, added to the AIC, arc42, Tech Stack Canvas, TOGAF
and C4 templates. The questions that were in the criteria list now live in the section they
fill.

**Every section closes one of three ways, and the engineer chooses:** answered; deferred, under Open
questions with who answers it and by which phase; or `Not applicable: <reason>`. A recommended
answer may be a deferral or not applicable. The engineer's answer stands, even a weak one: the risk
is said once, then the answer is recorded.

**A coverage check runs before handover.** Every template section holds content, a not-applicable
line, or an Open questions entry. A section with none was skipped and is asked about. The check
counts sections and grades nothing.

**Inputs are rewritten, decisions are superseded.** The goals, requirements, quality goals,
constraints and context (the AIC Goals group, arc42 §1 to §3) describe the problem as it is known
today, and gathering data is expected to change them. They are rewritten in place, with no reversal
note; git history keeps the old text. The hypotheses, arc42 §4 to §11 and the ADRs are decisions and
keep the supersession rule. A changed input puts every decision resting on it in question, and each
one is asked about: still holds, amend, or reverse.

**The Ask lines start a section, grilling finishes it.** Each answer is followed down until the
section is settled; the Ask lines guarantee breadth, the grilling protocol still supplies depth.

**A risk is grilled when it appears, not at plan time.** Whoever raises it, the engineer or the
agent, it goes into the risks section and is asked about one question at a time: likelihood and
impact, then the mitigation. The engineer picks one of five: reduce it by design, detect it with a
metric and an alert value, accept it with a reason and an owner, defer it as an open question, or
dismiss it with a reason. The coverage check also requires a chosen mitigation on every risk.
`/arcdlc:plan` Step 2.5 is unchanged and now mostly maps the recorded mitigations to tasks. A
mitigation is often a design decision, a fallback or a limit, and finding it at plan time meant a
trip back to `/arcdlc:aic`.

**Each section is written the moment it closes.** The interview used to write only glossary terms
and ADRs as it went; every other answer lived in the chat until the end. A context compaction
summarises the chat and can drop a figure, a deferral's owner, or which sections were already
closed. So each closed section goes into the document at once, and the document is the progress
record: after a compaction or in a new session, filled sections are closed, empty template slots
are open, and the interview resumes at the first open one. The engineer's review and `arctool sync`
still wait for Step 4.

**A re-run starts with the deferred items.** Open questions are the queue from one phase of an
initiative to the next.

## Justification

- **One source of truth.** The template already defines what a finished document holds. A separate
  criteria file would drift from it.
- **Architect-centric.** The AI makes sure nothing is forgotten. It does not decide whether an answer
  is good enough.
- **Mechanical.** "Every section closed" is checkable without judgement; a quality score is not.
- **Honest about unknowns.** A deferral with an owner and a phase is a legitimate answer, so there is
  no pressure to invent one.

## Trade-offs

- **Longer interviews.** More sections are asked about. Asking whether a theme applies before its
  sections, and filling what the code answers, keeps the count down.
- **A rewritten goal leaves no trail in the document.** Only git history shows what a goal used to
  say. This was accepted because the goal is not a decision to defend: the decisions it fed are, and
  they keep their record.
- **A half-written document sits in the initiative folder during the interview.** A plan run on it
  stops at `/arcdlc:plan` Step 2.4 on the empty sections, so no draft marker was added.
- **A weak answer passes.** Nothing grades it. `/arcdlc:plan` Step 2.4 still stops on a design too
  thin to plan, which is where a weak answer shows its cost.
- **The AIC template is longer.** It is still one page of sections; the Ask lines are guidance for
  the interview and are not copied into the document.
