# ADR-0017 — A marker is `ARCDLC` inside a comment, and cutting it out is asked for

- Status: Accepted
- Date: 2026-09-10
- Initiative: none (bundle-wide, tightens `arctool scan` and `/arcdlc:assist`)
- Amends: [ADR-0015](0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md) and
  [ADR-0016](0016-comment-markers-group-by-tag.md)
- Amended by: [ADR-0018](0018-only-single-line-comments-carry-markers.md): a marker counts in a
  single-line comment only

## Context

The sweep swept `TODO` by default and deleted every comment it registered. Both halves were wrong in
the same direction: they let a program change code that nobody had pointed it at.

Three things went wrong at once, and the third one proved it.

A repository that already uses `TODO` has notes that were never meant for this workflow. Running the
sweep on it planned them all and then took them out of the code.

Deleting the comment was part of the same run as recording it, so an engineer who wanted the list got
the edit whether or not they wanted it.

And the marker was matched by looking for a comment opener on a line, with a per-line check for quotes
around it. That check could not see a string that spans lines. This repository's own `usage` constant
in `cmd/arctool/main.go` is a Go raw string, and one line inside it reads:

```
                 markers sharing a tag are one block: // TODO:G1 in three files, one task
```

The sweep read that as a comment, registered it as a finding, and would have cut it out of the
constant, breaking the tool's own help text. The register would then have held a task invented from a
line of documentation.

## Decision

**The default marker is `ARCDLC`.** `scan.DefaultMarkers` is the bundle's own word, so a sweep over a
repository full of `TODO` notes finds nothing until somebody asks for them: `arctool scan --marker
TODO`, or `/arcdlc:assist <slug> TODO`. Group tags read the same way, `// ARCDLC:T1 ...`, and blocks
are named `ARCDLC-CMT-T1` or `ARCDLC-CMT-01`.

**Cutting a comment out of the code needs `--strip`.** Without it the sweep writes the register and
leaves the working tree exactly as it was, and says so: `left every marker comment in the code: pass
--strip to remove the 4 this run registered`. With it, nothing else changes: the register is still
written first, every edit is still verified, and a failed check still writes nothing anywhere.

**`/arcdlc:assist` asks, and it asks late.** Removing the comments is Step 6 of the skill, after the
register is judged and the tasks are in the plan, and it is one question with a recommended answer.
`--strip` is only passed on a yes. A no leaves every comment in place, which costs nothing: the
register already holds each marker, so the next sweep recognises them and duplicates no block. A
partial answer is done by hand, because `--strip` is all or nothing.

**A marker counts only inside a comment.** `internal/scan` lexes each file line by line, following
string literals as it goes, and carries a pending multi-line delimiter (a Go or JavaScript raw string,
a Python or Kotlin triple quote) from one line to the next. The walk stops at the comment opener, so
quotes inside comment text are never read as string delimiters. When the scanner cannot tell, it
reports no comment.

## Justification

- **A tool that edits code must be pointed at what it may edit.** A default marker of the bundle's own
  word, and an explicit flag for the edit, are the two smallest changes that make the sweep safe to
  run on a repository nobody has prepared for it.
- **`ARCDLC` is a word no codebase writes by accident.** It reads as an instruction to this workflow,
  which is what it is. `TODO` is a note to a human, and a repository's `TODO` notes belong to that
  repository.
- **Recording and editing are separate decisions, so they are separate commands.** The register is
  cheap, reversible, and useful on its own. The edit is neither reversible in place nor useful until
  the work is planned.
- **The question belongs after the plan.** A comment is only redundant once the task exists. Asking at
  that point makes the recommended answer easy to defend, and the run has already shown exactly what
  it would delete.
- **A scanner that follows strings is the honest minimum.** A parser per language is out of the
  question for a standard-library tool that sweeps twenty file types. Following quotes across lines
  catches the case that actually bites: text that documents markers, inside a constant.
- **When the scanner is unsure, it stays quiet.** A missed marker costs one sweep, and the engineer
  sees it in the report. A marker read out of a constant costs an edit to code nobody asked it to
  touch.

## Trade-offs

- **Teams must adopt a new word.** Existing `TODO` notes are invisible until `--marker TODO` is passed,
  which is a surprise for anyone who read the older docs. That is the price of not touching notes that
  were never addressed to this workflow.
- **Two runs instead of one.** A sweep that both records and removes is now `arctool scan` followed by
  `arctool scan --strip`. The second run re-reads the register, recognises every marker, and removes
  only what is registered, so the extra run is cheap and its output is a plain list of what it deleted.
- **The comment can outlive the task.** Answer no and the marker stays in the code while the plan
  carries the work. The register recognises it, so nothing is planned twice, but a reader of that file
  sees a note for work already queued. The skill's report says which comments were left.
- **The scanner is a heuristic, and it can go quiet.** An unbalanced backtick or triple quote in real
  code carries a pending literal to the end of the file and every marker below it is missed. Stopping
  the walk at the comment opener keeps prose out of that, and the direction of the failure is
  deliberate, but it is still a failure the engineer has to notice in the report.
- **Language-specific quoting is not modelled.** A shell heredoc, a Perl `q{}`, a C++ raw string
  literal: none of these are string delimiters to this scanner. A marker written inside one still reads
  as a comment. `--strip` being opt-in is what keeps that from being a code edit.
