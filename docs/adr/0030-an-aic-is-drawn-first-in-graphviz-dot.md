# ADR-0030 — An AIC is drawn first, in Graphviz DOT

- Status: Accepted
- Date: 2026-10-16
- Initiative: none (bundle-wide, changes `/arcdlc:aic` and its references)

## Context

Engineers reading AICs reported that most had no diagram at all, and that they expected a C4 diagram
and a BPMN diagram or flowchart in almost every one. The cause was in the bundle. The AIC template
mentioned pictures only as "include charts if needed". `SKILL.md` never pointed at
`Diagram Conventions.md`, and that file opened with "Draw a diagram only when it shows something the
reader cannot get from the text". The BPMN templates there were boxes in swimlanes, without the
events, gateways, pools and message flows that make a diagram BPMN.

## Decision

Every AIC carries three diagrams: a C4 system context diagram in Business Context, a C4 container
diagram and a flowchart or BPMN diagram of the main process in Solution Design. A missing one needs a
`Not applicable: <reason>` line, which the coverage check counts like a section. Beyond the three, a
paragraph that explains structure, order, branching or state is drawn instead.

Diagrams are written in Graphviz DOT, compiled to PNG, with the source committed beside the image.
After compiling, the agent opens the PNG and looks at it, against a checklist in
`Diagram Conventions.md`, because a diagram that compiles can still be unreadable. BPMN is drawn with
real BPMN 2.0 elements, each mapped to a DOT form, with rules for pools, lanes and flows, and a lane
layout technique that makes lanes full-width bands.

Rejected: Mermaid as the default. Generated Mermaid often fails to compile or draws overlapping lines,
in the engineers' experience every time it was tried, and it does not produce a PNG without extra
tooling. It is used only when the engineer asks. PlantUML is allowed for sequence diagrams when it is
installed, because DOT has no sequence notation.

Rejected as the generated format, kept for hand editing: draw.io. Its file is open, plain XML (the
mxGraph model), and a `.drawio.svg` is both the picture and its editable source, with real BPMN 2.0
shapes. But every node carries absolute coordinates and there is no automatic layout in the file, so
a model writing it guesses the layout, the same failure as Mermaid. A diagram an engineer wants to
maintain by hand is converted once from DOT with `graphviz2drawio`, and the `.drawio.svg` replaces
the DOT source from then on.

## Consequences

Easier: a reader sees the system, its containers and its main process before reading a paragraph,
and the text shrinks to what a picture cannot carry. A DOT source is diffable and stays editable.

Harder: an AIC run needs Graphviz installed; without it the agent writes the source and asks the
engineer to compile. Looking at every PNG costs a little time per diagram. DOT cannot draw BPMN event
markers or boundary events, so an event's trigger type is written in its label and a boundary event
is modelled another way.
