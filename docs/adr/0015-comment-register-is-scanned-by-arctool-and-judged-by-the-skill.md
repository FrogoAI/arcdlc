# ADR-0015 — The comment register is scanned by `arctool` and judged by the skill

- Status: Accepted
- Date: 2026-09-10
- Initiative: none (bundle-wide, introduced with `/arcdlc:assist`)

## Context

A `// TODO` in the code is work nobody planned. `/arcdlc:assist` turns those markers into plan tasks,
the same way `/arcdlc:examinate` turns compliance gaps into plan tasks.

Finding the markers is the expensive part for an agent and the cheap part for a program. An agent that
greps a repository, reads every hit with its surrounding code, and repeats that on every re-run pays
for the whole sweep in context. A forty-marker repository costs forty file reads. A program does the
same sweep for free and hands back forty lines.

Judging the markers is the opposite. "TODO maybe rethink this" names no work. Deciding whether a
marker is real work, stale, or a wish only the author can explain needs the engineer, one question at
a time. No program can do that.

So the two halves of the job belong to two different owners, and they meet in one file:
`docs/aics/<slug>/comments.md`.

## Decision

`arctool scan` writes the evidence. `/arcdlc:assist` writes the judgement. Each block in the register
has a line-level owner, and neither side writes the other's lines.

| Line | Written by | Rewritten later by |
| --- | --- | --- |
| `### <ID>: <title>` | `scan` writes the ID and a draft title | the skill rewrites the title; the ID never changes |
| `- Marker:` | `scan` | nobody, it is the finding's identity |
| `- Verdict:` | `scan` writes `NEW.` | the skill writes its judgement; `scan` writes `RESOLVED (<date>).` once the marker is gone |
| `- WHAT:` `- HOW:` `- WHY:` `- Acceptance:` | `scan` writes the empty keys | the skill fills them |
| `- WHERE:` | `scan` writes the marker's file and line, plus the code line under it | the skill adds the other files the change touches |

Three invariants make shared ownership safe:

- **`scan` only appends and only ever rewrites a `- Verdict:` line.** It re-parses its own output
  before writing and writes nothing on a failed check (exit `5`).
- **Identity is the file plus the normalized marker text, not the line number.** Line numbers drift on
  every edit; a moved marker keeps its block, its ID, and its judgement.
- **`scan` never touches `plan.md`.** Mirroring an `ACTIONABLE` finding into the plan is the skill's
  step, because it needs the judgement.

`/arcdlc:assist` keeps a full fallback for a machine with no `arctool`: it greps by the same rules and
writes the register by hand.

## Justification

- **The expensive half is the mechanical half.** Moving the sweep into the CLI is what makes the
  workflow affordable to re-run, and re-running is the point: a register that is only ever built once
  goes stale in a week.
- **A deterministic sweep is testable.** Opener rules per language, word boundaries, string-literal
  skipping and the exclusion list are unit-tested in `internal/scan`. An agent's grep is not.
- **The judgement cannot be automated, so it should not be faked.** A task written from a guess costs
  more than no task. The verdict column makes "not judged yet" a visible state rather than a silent
  one.
- **The register keeps history.** A resolved finding stays, marked `RESOLVED (<date>)`. Deleting it
  would let the next sweep plan the same comment twice.

## Trade-offs

- **Two owners of one file need a written boundary.** This ADR and the table in `AGENTS.md` are that
  boundary. An agent that edits a `- Marker:` line makes the next scan register the finding again as
  new. The skill states the rule; nothing enforces it.
- **A sweep without a language parser has limits.** Requiring the marker to open the comment, plus
  per-language openers and a string-literal check, kills the common false positives, but a marker
  quoted inside a multi-line string can still slip through. The engineer sees it as an `UNCLEAR` or
  `STALE` finding and it costs one answer, not a wrong task.
- **The marker set is explicit, not discovered.** `TODO` is the default and anything else has to be
  named. A team using `@todo` or `Todo` gets nothing until it passes `--marker`. Case-insensitive
  matching was rejected: it turns every `Todo` in prose into a finding.
