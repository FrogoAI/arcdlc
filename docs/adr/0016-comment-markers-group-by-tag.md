# ADR-0016 — Comment markers group by tag, and their block can grow

- Status: Accepted
- Date: 2026-09-10
- Initiative: none (bundle-wide, extends `/arcdlc:assist` and `arctool scan`)
- Amends: [ADR-0015](0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md)

## Context

One change is often written down in several places. An engineer who wants the index rebuild moved out
of `main` leaves a note where it sits, another where the return format has to change with it, and a
third next to the port it renames. ADR-0015 gave each marker its own block and its own plan task, so
one change arrived as three tasks with no link between them. The executor then implements a third of
a change three times, and each task's `Acceptance` can only test a third of it.

The register also had two smaller problems. Every block arrived with `- Verdict: NEW.`, a value
written by a program that judges nothing. And a block comment that closed on a later line was
registered with only its first line of text and left in the code, so the rest of what the engineer
wrote never reached the register.

## Decision

**A marker may carry a group tag, and a tag is one block.** The syntax is `MARKER:TAG` right after
the marker word, ended by a space or the end of the line: `// TODO:G1 move the rebuild into
internal`. The tag is folded to upper case, so `:g1` and `:G1` are one group. It reaches across the
whole sweep, so the same tag in three files is one block, and it belongs to one marker word, so
`FIXME:G1` is a different group from `TODO:G1`. An untagged marker stays a block of its own.

**A group block is named after its tag.** The slot after `-CMT-` holds either an auto number
(`TODO-CMT-01`) or the engineer's tag (`TODO-CMT-G1`). A tag with no letter in it is refused and read
as a plain marker, with a line in the run's report, because `// TODO:01` would otherwise claim the
auto-numbered ID `TODO-CMT-01`.

**The sweep seeds `WHAT` and writes no verdict.** It writes one `- Marker:` line per member, the
`WHERE` list, and a draft `WHAT`: the words the engineer typed, marker word and tag removed, members
joined with `; `. The verdict is judgement, so `- Verdict:` is written empty and only
`/arcdlc:assist` ever fills it. An empty verdict means the block is unjudged.

**A block may grow, and may change in no other way.** When a sweep finds a marker whose tag is
already registered, it appends to that block: a `- Marker:` line, more text on `WHAT`, another entry
under `WHERE`. It never opens a second block for the tag, so a tag maps to one ID forever. `Verify`
enforces the rest: a block nobody added to must be byte identical, a block that grew must keep every
earlier marker line, keep the old `WHAT` as a prefix and the old `WHERE` lines as a prefix, and match
byte for byte everywhere else. The title, `HOW`, `WHY`, `Acceptance` and the verdict a skill wrote
cannot be touched by a program.

**A block comment that closes on a later line is one finding, and is removed whole.** Its text runs
from the marker to the closer, decoration and closer trimmed. A second marker word inside the same
comment does not split it: one comment is one finding, so two groups need two comments. Three shapes
stay in the code and are reported as skipped: a marker on a decoration line of a block another line
opened (`/*` alone, then ` * TODO ...`), a block that never closes, and a block whose closer shares a
line with code.

## Justification

- **The unit of planning is the change, not the comment.** Three notes about one move are one task
  with one set of acceptance criteria. Merging them in the register is the only place it can happen
  cheaply: the sweep already reads every marker, and the engineer already knows which notes belong
  together at the moment of typing them.
- **The tag is the engineer's own label, so it makes a better ID than a number.** `TODO-CMT-G1` is
  recognisable in `plan.md`, in a commit message, and in the code the tag came from. A number would
  need a second lookup to mean anything.
- **Growth beats a second block.** A tag that came back would otherwise open `TODO-CMT-G1-02`, and
  the two blocks would describe one change with no link. Appending keeps one tag, one block, one
  task. It costs the append-only rule, which is why the growth is bounded to three named places and
  checked before anything is written.
- **A program that judges nothing should write no judgement.** `NEW` looked like a verdict and was a
  placeholder. An empty key says the same thing without pretending.
- **The register is the only copy of the marker text.** A block comment left in the code with half
  its text in the register was the one case where the sweep lost words the engineer wrote.

## Trade-offs

- **A grown block can outrun its plan task.** When the task was already `DONE`, the new marker cannot
  be folded into it, so `/arcdlc:assist` files a follow-up task (`TODO-CMT-G1-02`) that references
  the same block. The register and the plan stay in step, at the cost of one task ID that has no
  block of its own.
- **Append-only is now append-and-grow.** Two owners share a file, and one of them may edit three
  lines the other reads. `Verify` is the guard, and a failed check writes nothing anywhere (exit
  `5`), but the rule is harder to state than "nothing is ever rewritten".
- **A tag is only as good as the discipline behind it.** Tagging two unrelated notes with `G1`
  produces one task that does two things. The skill's answer is to judge a group as one unit and to
  grill when the members pull apart, not to split the block, because a split would lose a marker's
  text.
- **`:` after a marker word is a common habit.** `// TODO: fix this` has no tag, because the run of
  tag bytes is empty, and `// TODO:refactor(later)` has none either, because the text runs straight
  into it. Both read as plain markers, which is what they are, but the rule has to be learned.
- **One comment is one finding, whatever it holds.** Two markers inside one block comment cannot go
  to two groups. Splitting on a marker word inside a block would make the strip cut a comment in
  half, which is exactly what ADR-0015 refuses to do.
- **The `--json` shape changed.** `new` and `extended` now carry blocks with a `members` list instead
  of one finding each. The only consumer is `/arcdlc:assist`, and both moved together.
