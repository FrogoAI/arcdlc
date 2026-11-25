---
name: arcdlc-review
description: Review only the changes on this branch, as a reviewer, against the standards or policies you name (MDCA, DDD, SOLID, Modern Go, a docs/policies file) or the project's own architecture. Shows the findings and records them, never touches the plan; /arcdlc:plan turns them into tasks later. Use when someone says review my changes, review this branch, review the diff before I open a PR, check my changes against SOLID, or code review against our policies, or runs /arcdlc:review <slug> [standards], or invokes arcdlc-review.
argument-hint: "<slug> [MDCA,SOLID,...|policy-path] [--base <ref>|--range <a>..<b>] [--repo <name>]"
---

# ArcDLC Review (/arcdlc:review)

Review the change set, and only the change set, the way a reviewer reads a pull request: against the
standards and policies the engineer names, or the initiative's own architecture when none are named.
Show the result in chat and record it in `docs/aics/<slug>/review.md`. Never write `plan.md`: turning
findings into tasks is `/arcdlc:plan <slug>`, run when the engineer decides to.

`/arcdlc:execute` (or any commits) → `/arcdlc:review` → `/arcdlc:plan` → `/arcdlc:execute`

`/arcdlc:examinate` audits the whole codebase and files tasks at once. This skill reads one diff and
files nothing into the queue.

## Talk simple, write like a human

Every rule here governs everything you emit, chat messages exactly as much as the files you write.
Plain English: short sentences, common words, one idea each, active voice, a named actor.

- **No AI filler, anywhere.** Not "Furthermore", "In conclusion", "It is important to note", "delve",
  "leverage", "robust", "seamless". No warm-up opener, no praise, no restating the request back, no
  invented summary. Concrete names, numbers, and paths, never "significantly improves performance".
- **Be direct, name the thing.** The fewest words that carry the fact. Cut empty intensifiers:
  `honestly`, `genuinely`, `truly`, `clearly`, `obviously`. Never a metaphor or a rhetorical
  question: a worked example with real values beats both. Every title that names work is a verb plus
  its object: "Implement the alerter service".
- **No long dashes.** A full stop, a comma, a colon, or brackets instead of `—` and `–`. Hyphens,
  flags, and slugs stay. A format contract that requires `—` wins.
- **Domain terms stay.** Provenance, idempotent, backpressure, this project's own words: define each
  once in plain words, then use it. Simple English is about the sentence, not the term.
- **Shape.** In chat, bullets rather than paragraphs: what you did, what you found, what comes next.
  In files, vary sentence length and form, and draw what a picture shows better, so a document never
  reads like a list or a form to fill in.
- **One paragraph, one line.** In every file and reply you write, a paragraph, a list item, or an
  acceptance criterion is one physical line, however long. Never break a line inside a sentence to
  fit a width: viewers wrap on their own, and a hard break shows up as a broken sentence in a ticket,
  a chat, or a rendered page. This bundle's own instruction files are wrapped at 100 columns; do not
  copy that into what you write. Code blocks, tables, and contract lines keep their own layout.

Nothing is dropped to save space, in a reply as much as in a file. Never leave out a rule, path,
decision, trade-off, open question, or acceptance criterion because the answer is getting long: if it
bears on what the reader does next, it goes in. Short means no padding, never less content. Cut filler
words, repetition, and throat-clearing; never cut a fact. When completeness makes a reply long, let it
be long. The full standard, with examples and a pre-save check, is
`../grilling/references/Writing Style.md` (flat installs:
`../arcdlc-grilling/references/Writing Style.md`).

## Judge by the four virtues

- **Wisdom.** Unclear is a question, not a guess. A contradiction, a missing decision, code that looks
  dead, a change that is risky or hard to reverse: grill it, never pick for the engineer.
- **Courage.** Say the hard thing. A task that is not mechanical is a plan defect, not a reason to
  raise the tier. Never stop to report a helper skill as missing.
- **Justice.** Write every answer where the next session reads it, never only in the chat. Name the
  tier you actually used. Never report a same-tier spawn as a cheaper run.
- **Temperance.** Touch only what you were pointed at. Cutting a comment out of the code is asked for,
  not assumed. No invented summary, no scope you were not given.

## Initiative selection

The initiative slug is the **required first positional argument**: `/arcdlc:review <slug> [standards]`.
The review is recorded in `docs/aics/<slug>/review.md`, and the slug is passed to `arctool` as
`--aic <slug>`. If the slug is missing, stop and report the error, listing the existing initiatives
under `docs/aics/`. Never guess. If the named initiative folder does not exist yet (a review of work
that belongs to no initiative, e.g. `auth-hotfix`), confirm the slug with the engineer and create
`docs/aics/<slug>/`.

**A closed initiative is final.** If `docs/aics/<slug>/CLOSED.md` exists, stop and change nothing: say
that `<slug>` was closed on the date in its `- Closed:` line, and that follow-up work is a new
initiative with its own slug whose design links `docs/aics/<slug>/` as its starting point. No skill
edits or deletes a closed initiative's files. The one exception is the `## Superseded` line that
`/arcdlc:aic` appends to its `CLOSED.md` when another initiative's design reverses it.

## In a workspace

The working directory is a workspace when it is not a git work tree and `git -C docs rev-parse
--show-toplevel` resolves to `<cwd>/docs` (see the `Workspace` and `Hub` terms in `CONTEXT.md`). In a
workspace, `docs/` is a repository of its own. Every write this skill makes there, `review.md`, an
ADR, a `CONTEXT.md` term, is committed in the hub with `git -C docs commit -- <hub-relative paths>`
under the message `docs(<slug>): <what changed>` plus `#AI-assisted`, and pushed to the hub's default
branch at once. A rejected push is retried once after `git -C docs pull --ff-only`. Nothing is
written into a product repository by this skill.

The change set lives in one product repository. `--repo <name>` names it. Without it, list the
sibling repositories that have a change set (Step 1) and ask which one; review two repositories as
two runs. Every finding then carries `- Repo: <name>` directly after `WHAT`, and every `WHERE` path
carries that repository's prefix (`fdb-server/internal/migration/...`), so `/arcdlc:plan` mirrors it
without a rewrite.

## Step 1: Resolve the change set

Review only what changed. Code the change set does not touch is out of scope, however bad it is.

- **Default.** Find the default branch: `git symbolic-ref --short refs/remotes/origin/HEAD` with the
  `origin/` prefix removed, else `main`, else `master`. The base is `git merge-base HEAD <default>`.
  The change set is `git diff <base>` (committed, staged and unstaged work together) plus the
  untracked files from `git ls-files --others --exclude-standard`.
- **`--base <ref>`** replaces the default branch as the base. **`--range <a>..<b>`** reviews exactly
  `git diff <a>..<b>` and ignores uncommitted work.
- **Nothing to review.** On the default branch with a clean tree, or with an empty diff, stop and say
  so. Never widen the scope to "the last few commits" by yourself.
- Record the review header now: branch, head commit (`git rev-parse --short HEAD`), base ref and
  commit, and the size from `git diff --shortstat <base>` plus the untracked count.
- Read each changed file in full, not only its hunks. A hunk alone hides the caller it breaks.
  Generated files, lock files, and vendored code are listed in the header and not reviewed line by
  line.

## Step 2: Resolve the standards

The second argument is a comma-separated list. Each item is a standard name or a policy path, and
the review applies all of them together: `/arcdlc:review <slug> SOLID,MDCA,docs/policies/logging.md`.

- **Project rules always apply**, with or without a named standard, when they exist: every policy
  under `docs/policies/` (written by `/arcdlc:policy`, listed in its index) and `CONVENTION.md` at
  the repository root. In a workspace, read the hub's `docs/policies/` and `docs/CONVENTION.md`, then
  the reviewed repository's own `CONVENTION.md`; the repository's file wins where both cover the same
  rule. Name each file in the review header, or write `none found`.

- **A named standard or a policy path** resolves exactly as `/arcdlc:examinate` Step 1 says: read
  `../examinate/SKILL.md` (flat installs: `../arcdlc-examinate/SKILL.md`) and follow its rules for
  the bundled rule sets in its `references/` folder, the `go.mod` check for Modern Go, a standard with
  no bundled rule set, a project policy, and project-beats-bundled precedence. Only `mdca.md`,
  `solid.md`, `ddd.md`, `ECS.md` and `Modern Go.md` carry identifiers a finding may cite
  (`MDCA-P3.2`, `SOLID-7`, `GO-12`).
- **No standard named.** Review against the initiative's own design: the architecture document in
  `docs/aics/<slug>/` (and `source-aic.md` when `/arcdlc:derive` wrote one, limited to its
  `- Scope:`), `docs/adr/`, `CONTEXT.md`, `docs/policies/*.md`, and the hard rules in `AGENTS.md` or
  `CLAUDE.md`. When none of these exist, review for correctness alone and say so in the header.
- **Correctness is always checked**, whatever was named, because a reviewer who passes a bug has not
  reviewed: a wrong result, a crash, a race, a lost error, a broken contract with a caller, a security
  hole, a changed behaviour with no test, a test that cannot fail. Cite these as `correctness`,
  `security` or `test`.
- Extract the checkable rules before reading the diff, so every finding cites a rule, not a feeling.
  Style the named rules do not mandate is not a finding.

## Step 3: Review the diff

- Walk the change set file by file against every rule from Step 2.
- Every finding needs evidence: `file:line` in the changed code plus the rule it breaks. No evidence,
  no finding.
- A problem in a line the change set did not touch is not a finding. When the change makes it worse
  (a new caller of a broken function), it is a finding on the new line. Otherwise mention it once in
  the report under "Seen outside the change" and file nothing.
- Classify each finding twice. The source status is the plan contract's tag: `MISSING` (a required
  element is absent, such as a test or an error check), `PARTIAL` (present but incomplete), `DRIFT`
  (present but breaks the rule). The severity is the review's own:
  - `BLOCKER`: the change must not merge as it is. A bug, data loss, a security hole, a broken
    contract or public interface, a failing build or test.
  - `MAJOR`: breaks a named standard, policy, ADR, or the architecture document, or changes behaviour
    with no test.
  - `MINOR`: a named rule is broken with no effect on behaviour (naming, layout, a missing doc
    comment the standard requires).
- Run what exists and is cheap: the project's build, test and lint commands from `AGENTS.md`,
  `CLAUDE.md` or the `Makefile`. A failure the diff caused is a `BLOCKER`. A failure that exists on the
  base too goes under "Seen outside the change". Change no file to make a check pass.
- Never fix the code. The review reports; `/arcdlc:plan` and `/arcdlc:execute` change things.

## Step 3.5: When a finding needs a human, grill (mandatory)

Some findings are not yours to settle, because `/arcdlc:plan` will turn the recorded fix into a task
and `/arcdlc:execute` will carry it out. Stop and ask the moment you hit one of these:

- **A violation that may be deliberate.** The diff reads like a conscious exception to the rule.
- **A contradiction.** The diff follows one document and breaks another (a policy against an ADR, the
  AIC against `source-aic.md`). Pick nothing until the engineer says which one wins.
- **An ambiguous rule.** The rule reads two ways and the two readings give different findings.
- **A risky fix.** The right fix is a migration, a changed wire format, a deleted public interface, a
  dependency swap, or a rename that breaks callers outside this repository.
- **Code the change made dead.** You find no caller left, which does not mean there is none.

Prefer this bundle's own `arcdlc-grilling` skill (`/arcdlc:grilling`). If it cannot be invoked here,
read `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and run its protocol
inline. What is mandatory is the grilled interview, never the invocation: never stop to report a
helper skill as missing, and never look for a grilling skill outside this bundle. **One question per
turn**: ask, wait for the answer, then ask the next, each carrying your recommended answer and one
line of why. Never a numbered round. Facts are yours to find: if the code, config, or an existing doc
answers it, look it up instead of asking. Write each decision down the moment it settles, glossary
terms into `CONTEXT.md` and hard trade-offs into `docs/adr/NNNN-<slug>.md`, never only in the chat.

Each answer lands in `review.md`:

- **Fix approved:** write the engineer's decision into the finding's `HOW`.
- **Deviation accepted:** add `- Accepted: <reason>, <date>.` to the finding. It stays in `review.md`
  as evidence, does not count against the verdict, and `/arcdlc:plan` gives it no task.
- **A decision that outlives the review:** write an ADR in `docs/adr/` and cite it in the finding's
  `WHY`. **A new term:** add it to `CONTEXT.md`.

## Step 4: Write the review register

Write or update `docs/aics/<slug>/review.md`. Never touch `plan.md`, `gap.md` or `comments.md`. Each
run appends one `## Review` section below the earlier ones, so the file is the history of every
review of this initiative:

```md
# Review register: <slug>

## Review <YYYY-MM-DD>: <branch> at <head> against <base-ref> at <base>

- Standards: <each rule set, with its path; "correctness only" when nothing else applied>
- Change set: <N files, +A -D, U untracked; files listed but not reviewed and why>
- Verdict: <CHANGES REQUESTED | APPROVED> (<n> BLOCKER, <n> MAJOR, <n> MINOR, <n> accepted)

### REV-NN (<MISSING|PARTIAL|DRIFT>): <Short Title>

- Severity: <BLOCKER|MAJOR|MINOR>.
- WHAT: <What must change, one line.>
- HOW:
  <The fix decision the review already knows. Out of scope: <adjacent code the executor must leave alone>.>
- WHERE: <Exact file:line in the change set.>
- WHY: <The broken rule, quoted or paraphrased, with its source path or identifier.>
- Acceptance:
  - GIVEN <the changed code> WHEN <the named test, lint rule, grep, or build command> THEN <the observable result>.
```

- The verdict is `CHANGES REQUESTED` when any finding that is not accepted is a `BLOCKER` or a
  `MAJOR`, and `APPROVED` otherwise. A review with no finding writes the header with `APPROVED` and no
  blocks.
- `REV-NN` numbers continue from the highest ID already in `review.md`, across sections. Never
  renumber, reuse, or delete a finding.
- Write each block to be **mechanical**, exactly as `/arcdlc:examinate` Step 3 requires of a gap: the
  title is a verb plus its object, `WHAT` names the change, `HOW` carries the decision, `WHERE` is the
  only place the executor edits, and `Acceptance` is a runnable re-check wherever one exists. In a
  workspace, `- Repo: <name>` follows `WHAT`.
- `- Severity:`, `- Accepted:` and `- Resolved:` are register lines. `/arcdlc:plan` drops them when
  it mirrors a finding.
- **On a re-run**, check every earlier finding that is neither accepted nor resolved. When the new
  change set fixes it, add `- Resolved: <date>, <head>.` to its block; `/arcdlc:plan` gives a resolved
  finding no task. When it still holds, leave it and list its ID in the new section's verdict line
  (`still open: REV-03, REV-05`). Never copy it into the new section.
- Validate the blocks against the plan contract, so `/arcdlc:plan` mirrors them without a fix: when
  `command -v arctool` finds it, run `arctool validate --strict --plan docs/aics/<slug>/review.md`
  and fix every error except the two a register always gets: `has no "- Status:" line` and `missing
  required key(s): References, Status`. A register carries neither, so the command exits 1 on a clean
  register too. Without `arctool`, hand-check the "Authoring Rules" section of
  `../plan/references/plan-format.md` (flat installs: `../arcdlc-plan/references/plan-format.md`).

## Step 5: Report

Show the review in chat, in the order a reviewer reads it:

- The verdict line and the counts by severity.
- Every finding, `BLOCKER` first: ID, severity, title, `file:line`, the rule, and the fix in one line.
- Every finding you took to the engineer in Step 3.5 and what was decided. Accepted deviations are
  listed apart, so nobody reads them as work left undone.
- Earlier findings resolved or still open, by ID.
- "Seen outside the change": problems in untouched code, one line each, with no ID.
- What was not reviewed and why (generated files, a check you could not run).
- Next step: `/arcdlc:plan <slug>` turns every open finding into a task. Say it, never run it.
