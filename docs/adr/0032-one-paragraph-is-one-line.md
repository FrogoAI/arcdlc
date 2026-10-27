# ADR-0032 — A paragraph in a generated document is one line

- Status: Accepted
- Date: 2026-10-27
- Initiative: none (bundle-wide, changes the shared writing-style block and the templates the skills copy)

## Context

Documents the skills wrote broke sentences in the middle: a line ended at about 100 columns and the sentence went on in the next one. In a Markdown viewer that is invisible, but pasted into a ticket, a chat, or an HTML page the break shows as a broken sentence. Engineers asked for a fix in session after session, because nothing in the bundle said not to do it.

The agent copied the habit from the bundle itself. Every `SKILL.md`, the AIC template, the AIC example, the ADR template, and the bundled policies were hard-wrapped at 100 columns, and an agent matches the layout of the text it reads.

## Decision

Every paragraph, list item, and acceptance criterion a skill writes is one physical line, however long.

- The shared block `## Talk simple, write like a human` gains a bullet, "One paragraph, one line", in all twelve `SKILL.md` files. It says the bundle's own wrapping is not a pattern to copy.
- `skills/grilling/references/Writing Style.md` gains a section with the reasons and the exceptions (code blocks, tables, Graphviz sources, contract lines, and commit messages, which keep the git convention of a body wrapped at 72 columns), a step in the pass before saving, and a grep that finds hard-wrapped prose.
- The references whose text the skills copy into their output are unwrapped: `AIC Template.md`, `AIC Example.md`, `ADR.md`, `Tech Stack Canvas.md`, `TOGAF.md`, `arc42-template-EN.md`, `Scaffold Templates.md`, and `Engineering Principles.md`. The `SKILL.md` files and the rule-set references stay wrapped: they are this repository's source, read by the agent, never copied.

## Consequences

Easier: a generated document pastes cleanly into a tracker, a chat, or a rendered page, and nobody has to ask for the fix again.

Harder: a diff of a long paragraph shows the whole line as changed rather than one wrapped line, so a reviewer relies on word diff (`git diff --word-diff`). Documents written before this change keep their wrapping until a skill next amends them.
