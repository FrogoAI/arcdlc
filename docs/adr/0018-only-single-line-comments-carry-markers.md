# ADR-0018 — Only single-line comments carry markers

- Status: Accepted
- Date: 2026-09-10
- Initiative: none (bundle-wide, narrows `arctool scan`)
- Amends: [ADR-0015](0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md),
  [ADR-0016](0016-comment-markers-group-by-tag.md) and
  [ADR-0017](0017-a-marker-is-arcdlc-in-a-comment-and-cutting-it-is-asked-for.md), each of which
  described block comments as findings

## Context

Block comments cost more than they are worth to a sweep, and the bill arrived in instalments.

ADR-0015 read a marker on a `/*` line and refused to remove it when the block spanned lines, so the
comment came back on every sweep. ADR-0016 taught the sweep to join a block's lines and remove the
whole comment, which added a closer token per language, a rule for the text between the opener and
the closer, and three shapes that had to be refused: a marker on a decorated line of a block another
line opened, a block that never closes, and a block whose closer shares a line with code. Then a
marker written under a bare `/*` turned out not to be found at all, because nothing on that line
opens a comment, which needed block-comment state carried from line to line in the scanner. And the
first request after that was to remove such a comment, which is only safe when the block holds
nothing but the note: another rule, another shape to detect, another way to strand an opener.

Each step was small. Together they made the extent of a note depend on per-language closing tokens,
decoration and nesting, and made the removal a case analysis with a wrong answer that edits code.

A single-line comment has none of that. It ends where its line ends.

## Decision

**A marker counts in a single-line comment and nowhere else.** `// ARCDLC ...` and `# ARCDLC ...` are
markers. `/* ARCDLC ... */`, `<!-- ARCDLC ... -->`, `(* ARCDLC ... *)` and every other block comment
are prose to the sweep, in every language, whatever they hold.

**A note may still span lines, one comment per line.** Consecutive single-line comments continue the
same finding, so a three-line note is one task. A blank line, a line of code, or a second marker ends
the run.

**`openersByExt` lists single-line openers only.** A language whose only comment is a block comment
carries no markers, and its entry is empty rather than absent: `.css`, `.html`, `.htm`, `.xml`,
`.ml`, `.mli`. The sweep reads such a file, finds nothing, and does not report it as a file type it
had to guess at.

**The strip has no shape to refuse.** Removing a marker is a whole-line delete of the comment lines,
or a cut to the end of a line when the comment trails code. Both are safe in every language, so
`stripPlan` lost its `ok` and `reason` fields. The one thing that still stops a removal is a file
edited between the sweep and the edit, which `Strip` reports.

`blockClosers`, `blockSpan`, `blockLine` and the block-comment state in the line scanner are gone.
The scanner still follows string literals across lines, because a marker inside a constant is the one
false positive that matters.

## Justification

- **One rule beats a family of them.** The old behaviour needed a closing token per language, a
  decoration rule, a "is the block only this note" test, and three documented refusals. The new one
  needs a line.
- **The extent of a note stops being a language question.** A single-line comment ends at the newline
  in Go, Ada, Fortran, batch and everything else. That is what makes the sweep's promise the same in
  every file type.
- **A removal that cannot strand an opener cannot break a build.** The worst a whole-line delete or a
  line-tail cut can do is remove a comment. Every dangerous edit the sweep could make involved a
  block comment.
- **The cost falls on the person writing the note, at the moment they write it.** `//` instead of
  `/*` is a keystroke, and the register is one comment style simpler to explain to a team.

## Trade-offs

- **A note in HTML, CSS, XML or OCaml cannot be written at all.** Those languages have no single-line
  comment, so their entries are empty and the sweep finds nothing there. The note has to live in a
  file that can hold one, and `docs/comment-markers.md` says so.
- **A marker already written as a block comment goes quiet.** It is not registered, not removed, and
  not reported: to the sweep it does not exist. `/arcdlc:assist` is told to name this when a marker an
  engineer wrote never reaches the register, but nothing detects it automatically.
- **A long note costs a comment marker per line.** Three lines of `//` instead of one `/* */`. Editors
  do this with one keystroke, and the sweep joins the lines back into one finding.
- **Two commits of block-comment support were deleted.** The support worked, and its tests passed. It
  is gone because the rules it needed kept growing, which is a judgement about the next change rather
  than about the code that exists.
