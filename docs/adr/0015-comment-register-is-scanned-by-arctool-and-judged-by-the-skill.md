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
same sweep for free and hands back forty lines. Taking the comments back out of the code is the same
trade: forty small edits an agent would otherwise make file by file.

Judging the markers is the opposite. "TODO maybe rethink this" names no work. Deciding whether a
marker is real work, stale, or a wish only the author can explain needs the engineer, one question at
a time. No program can do that.

So the two halves of the job belong to two different owners, and they meet in one file:
`docs/aics/<slug>/comments.md`.

## Decision

`arctool scan` writes the evidence and removes the comment. `/arcdlc:assist` writes the judgement.
Each block in the register has a line-level owner, and neither side writes the other's lines.

| Line | Written by | Rewritten later by |
| --- | --- | --- |
| `### <ID>: <title>` | `scan` writes the ID and a draft title | the skill rewrites the title; the ID never changes |
| `- Marker:` | `scan` | nobody, it is the finding's identity |
| `- Verdict:` | `scan` writes `NEW.` | the skill writes its judgement; `scan` never writes it again |
| `- WHAT:` `- HOW:` `- WHY:` `- Acceptance:` | `scan` writes the empty keys | the skill fills them |
| `- WHERE:` | `scan` writes the marker's file and line, plus the code line under it | the skill adds the other files the change touches |

The marker comment does not stay in both places. Once a block is in the register, `arctool scan`
deletes the marker comment from the source: the register is where a marker lives, the code is where the
work lives. That is why the register is append-only and why a block is never deleted, not even for a
finding judged `STALE`: it is the only copy of the text.

Four invariants make this safe:

- **The register is written first, the code second.** The marker text is never in flight without a
  home. A failed source write leaves the marker in the code, and the next sweep retries it.
- **`scan` deletes only the lines it recorded.** A standalone comment line (with its continuation
  lines) goes; a comment trailing code has its span cut and the code kept. A marker inside a
  multi-line block comment is left alone and reported, because removing one of its lines can leave a
  dangling opener.
- **Every edit is verified before anything is written.** Each result must be the original file with
  whole lines deleted or one span cut, in order, with the counts the plan intended. A failed check
  writes nothing anywhere (exit `5`).
- **Identity is the file plus the normalized marker text, not the line number.** Line numbers drift on
  every edit; a marker that was already registered is recognised rather than duplicated.
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
- **A marker in two places is a marker planned twice.** Leaving the comment in the code invites a
  second sweep, a second block, and a second task for work already queued. Moving it makes the
  register the single answer to "what did the team leave behind".

## Trade-offs

- **The sweep edits the working tree.** That is a real cost: an engineer who runs `/arcdlc:assist` gets
  a code diff along with a register. `--dry-run` shows both halves and writes nothing, the register is
  written before any source file, and every edit is verified first. Beyond that, the recovery is the
  one every repository already has: `git diff` and `git checkout`.
- **A deferred marker is no longer visible in the code.** Somebody reading that file will not see the
  note. The register holds it, and the reason has to be written well enough to survive the move. This
  is deliberate: a marker nobody planned is not a plan, and the register is where unplanned work is
  tracked.

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
