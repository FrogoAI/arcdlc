# ArcDLC Close: a finished initiative keeps its design and stops taking tasks

> Add `/arcdlc:close <slug>`, which writes a `CLOSED.md` note into the initiative's folder and deletes its plan and registers, so a finished design stays where every link already points and no new work lands on it by accident.

Interview started on 2026-09-25. Terms are defined in [`CONTEXT.md`](../../../CONTEXT.md).

## Goals

### Business Case

An initiative today has two end states: open forever in `docs/aics/<slug>/`, or deleted by
`/arcdlc:remove`. There is no third state for "finished, kept as a reference", and that gap causes
three symptoms.

- **Done cannot be told from paused.** The `init` initiative has 10 tasks, all `DONE` and archived
  to `plan-archive.md`, and its `plan.md` holds no task blocks. The registry in `AGENTS.md` and
  `README.md` lists it exactly as it would list an initiative started yesterday.
- **Finished work still takes new tasks.** Nothing stops `/arcdlc:plan`, `/arcdlc:examinate` or
  `/arcdlc:assist` from appending to a finished initiative's queue.
- **Finished designs leave the tree.** Five initiatives were deleted with `/arcdlc:remove` on
  2026-09-16. Seven ADRs (0005, 0007 to 0010, 0012 to 0014) now read "folder retired, see git
  history", so the designs behind them can only be read with `git log`. The next initiative cannot
  use them as a reference.

The fix is one new end state, not a tracking system. The engineer asked for it to follow KISS: no
new index files, nothing moves, and a `CLOSED.md` file in the initiative's own folder is its state.

### Functional Overview

The initiative gives an initiative a finished state that keeps its design in place, and a way to see
where every initiative stands.

- `/arcdlc:close <slug>` writes `docs/aics/<slug>/CLOSED.md` and deletes the plan, its board stories
  and both registers (H1 to H7).
- Every skill that changes an initiative refuses a closed one; a closed initiative is final (H10).
- A later design that reverses a closed one appends a "Superseded" line to its `CLOSED.md` (H14).
- `/arcdlc:remove` deletes a design folder for good, after a confirmation (H11).
- `arctool` gains one read-only command listing every initiative with its phase (H8), and `sync`
  leaves closed initiatives out of the registry (H9); nothing else (H6). Close and remove do not
  read or rewrite links (H5).

Left out on purpose: restoring the five initiatives deleted on 2026-09-16 (`antigravity-cli`,
`cursor-support`, `initiative-lifecycle`, `ordering`, `source-library-cleanup`). Their documents
stay in git history (commit `a0cb40e`), closed designs accumulate only from this change on, and the
seven ADRs reading "folder retired, see git history" keep that line. The cost: those designs stay
out of reach of a reader browsing `docs/aics/`.

Closing `init` is not part of this initiative's plan. Closing needs a person's yes (H6), so it
cannot be a mechanical task, and writing that yes into a task would bypass the confirmation. After
release, the engineer runs `/arcdlc:close init` as the first real close; that opens the next stage.
The cost: the first real close happens outside the plan's verification.

### Quality Goals

1. **Integrity: nothing is lost silently.** Only the five files H4 names can be deleted, and each
   one is named in the confirmation the engineer answered. Pass: the skill's last step lists the
   folder and `git status`, and the report shows the design documents unchanged and exactly the
   named files deleted.
2. **Honesty: no silent state.** A closed design always holds `CLOSED.md`, and no closed design
   holds `plan.md`. Pass: after close, `ls docs/aics/<slug>/` shows `CLOSED.md` and none of the four
   files.
3. **Reliability: resumable.** A close interrupted between writing `CLOSED.md` and the deletion
   finishes when re-run. Pass: the skill's first step detects `CLOSED.md` beside a leftover plan
   file and goes straight to the confirmation and the deletion (H6).

Speed and scale are not goals: a repository holds tens of initiatives, not thousands. Link integrity
is not a goal: close and remove do not handle links (H5).

### Organizational Constraints

The arcdlc maintainer owns the initiative. The team is one engineer plus executor agents on the
`sonnet` pin in `CONTEXT.md`. There is no deadline and no budget line. It ships as one release,
`arctool` 0.21.0 and bundle 0.34.0, through a pull request with CI green. The engineer's rule of no
new index files is binding; only the maintainer can lift it.

### Technical Constraints

All from `AGENTS.md`; none is challenged by this initiative.

- `arctool` changes by one read-only listing command (H8) and the closed-initiative rule in `sync`
  (H9), and nothing else: it stays pure standard
  library, and its exit codes keep their meaning. `/arcdlc:close` may use existing commands
  (`arctool list --aic <slug>` for the counts, `arctool sync` for the registry), probing
  `command -v arctool` first and describing the manual fallback, as every skill does.
- The skill works as `/arcdlc:close` (plugin) and `arcdlc-close` (flat install), with dual paths
  (`../grilling/...` and `../arcdlc-grilling/...`) where it reaches a sibling.
- `## Talk simple, write like a human` and `## Judge by the four virtues` are copied verbatim into
  the new `SKILL.md`; CI pins both.
- Adding `close` as the twelfth skill touches, in one change set: `SUBSKILLS` in `install.sh`,
  `Skills` in `internal/bundle/bundle.go`, the four skill loops in `.github/workflows/ci.yml`,
  `docs/routing-checks.md`, and every place that says "eleven" skills (`AGENTS.md`, `CONTEXT.md`, CI
  comments). Changing `remove` (H11, R2) touches its `SKILL.md`, its description, and the
  `/arcdlc:remove` mentions in `README.md`, `AGENTS.md`, `CONTEXT.md` and the `aic`, `init`,
  `examinate` and `assist` skills.
- Version bumps: `const version` in `cmd/arctool/main.go` (0.20.0 today) for the listing command,
  and both plugin manifests (0.33.0 today), moved together.
- The plan format is not changed: close reads task statuses with the existing parser.
- In a workspace, the hub rules of ADR-0026 apply (H12).
- The plan contract path stays `docs/aics/<slug>/plan.md`: nothing in ADR-0001 or ADR-0026 changes.

### Business Context

| Partner | Flow | Change it needs |
| --- | --- | --- |
| Engineer | Runs `/arcdlc:close`, answers the outcome and confirmation questions, reads `arctool status`. | None. |
| `/arcdlc:aic`, `plan`, `examinate`, `assist` | Change an initiative. | Stop on a folder holding `CLOSED.md` (H10). |
| `/arcdlc:execute`, `archive`, `plan-human` | Read an existing plan. | None: after close `plan.md` is gone, so they stop as they do today. |
| `/arcdlc:aic` | Starts or amends a design. | Reads closed designs' `CLOSED.md` in Step 1 (H9); appends a "Superseded" line to the `CLOSED.md` of a closed design it reverses (H14). |
| `/arcdlc:remove` | Takes a slug. | Deletes `docs/aics/<slug>/` after a confirmation (H11). |
| `/arcdlc:init` | Scaffolds `docs/`. | None. |
| `arctool` | Reads and writes `docs/aics/<slug>/plan.md`, syncs the registry. | `status` (H8) and closed initiatives left out of `sync` (H9); nothing else (H6). |
| git and GitHub | Commit the change, render `CLOSED.md` when opened. | None. |

No partner is controlled by an attacker: every input is a file in the engineer's own repository.

## Architectural Hypotheses

### H1. A closed initiative stays in `docs/aics/<slug>/`, marked by `CLOSED.md`, with its plan and registers deleted

Close changes nothing but two things in the initiative's folder: it writes `CLOSED.md` (H2) and
deletes `plan.md`, `plan-archive.md`, `plan-human.md`, `gap.md` and `comments.md` (H4). The
architecture documents and their images stay exactly where they are. `/arcdlc:remove` deletes a
design folder for good, after a confirmation (H11).

Justification: nothing moves, so every link to a design (ADR-0026 links `../aics/init/aic.md`
today) keeps working with no rewrite, and the plan contract path `docs/aics/<slug>/plan.md` stays as
it is. Close shrinks to a note plus a delete of the files H4 chose to drop. OpenSpec, AWS AI-DLC and
SPDD were compared on 2026-09-25: all three keep finished designs in the tree, and the two that move
folders (OpenSpec's `changes/archive/`) never fix the links a move breaks.

Trade-offs: work in flight is seen through `arctool status` (H8), not through a folder listing.
`docs/aics/` holds open and closed designs side by side, told apart by `CLOSED.md`, and it stays
the actual design only through H10 and H14.

Both `docs/adr/` and `docs/aics/` stay, at different grain. An ADR records one decision that is hard
to reverse, surprising, and a real trade-off; a design holds the goals, numbers, context, every
hypothesis with its rejected options, and the risks. Most reasoning never reaches an ADR: `init` has
9 hypotheses and 1 ADR, `aic-tracking` 14 and 1. Rejected: deleting designs on close and keeping
only ADRs (keeps about one decision in ten, the loss this initiative exists to fix), and dropping
ADRs into the designs (loses the global decision list `/arcdlc:examinate` audits against).

Reversals, both in this interview. The first design moved the whole folder to
`docs/history/<date>_<slug>/` on close and rewrote every link into it. At plan time the engineer
proposed splitting design from work (`docs/aics/` for designs, `docs/plans/` for plans), because the
move and its link rewriting were the most complex part. The engineer then proposed this simpler
form: it keeps the split's benefit (no move, no link rewriting) and drops its largest cost, moving
the plan contract to a new path across `arctool`, every skill, ADR-0001, ADR-0026 and `init`.

### H2. `/arcdlc:close` writes a closing note as `CLOSED.md` in the closed folder

The close date, the outcome and what was left open go into a new `docs/aics/<slug>/CLOSED.md` in
a fixed short shape:

- `# <title>`, copied from the architecture document's H1.
- `Closed: YYYY-MM-DD`.
- `Outcome:` one of `delivered`, `partly delivered`, `stopped`.
- The task count.
- For an outcome other than `delivered`: the engineer's one line of reason (H7), and each unfinished
  task by ID and title (H3), with the `- Marker:` text and `- WHERE:` locations of any
  `ARCDLC-CMT-*` task (R3).
- What was left open: the architecture document's unanswered Open questions and its deferred risks,
  copied as they stand.
- A relative link to each architecture document in the folder.
- Later, a `## Superseded` section, appended by other initiatives' design rounds (H14). Close never
  writes it.

An initiative with no architecture document (an audit from `/arcdlc:examinate` or a sweep from
`/arcdlc:assist`, holding only a plan and a register) is closed the same way: `CLOSED.md` takes the
slug as its title, lists `Documents: none`, and the folder stays as a one-file record that the work
happened and how it ended. One rule for every initiative beats a special case; rejected: refusing
such a close and pointing to `/arcdlc:remove`.

The name is always `CLOSED.md`, never `README.md`: an initiative folder may already hold a
`README.md` of its own, and a fixed name that no other file uses means zero collisions and one path
every reader and tool can rely on (H8, H10). The architecture documents stay exactly as they were
decided: the note records what happened, the documents record what was decided.

Trade-offs: GitHub does not render `CLOSED.md` on opening the folder, as it would a `README.md`; a
reader clicks it. One small generated file per closed initiative. It lives inside the initiative's
own folder, so it is not a shared index file and it cannot drift from other initiatives. Rejected: a
`## Closed` section added to the architecture document (edits a decided design, and with several
formats it is unclear which document takes it), and no note at all with the date read from git (the
outcome and what was left open would be lost). Reversal: the first round named the note `README.md`,
for GitHub's folder rendering. At plan time a collision with a folder's own `README.md` surfaced,
and the engineer chose one fixed name, `CLOSED.md`, over refusing, prepending, or switching names
only on a collision.

### H3. Close refuses only while a task is `TAKEN`, and shows unfinished work before it asks

`/arcdlc:close` stops, changing nothing, while any task in `plan.md` is `TAKEN`, because an agent is
working on it at that moment.

A `TODO` or `BLOCKED` task does not stop it. Before asking for confirmation, the skill highlights
how many tasks are `TODO` and how many are `BLOCKED`, by count and with each ID and title, and
closes only on an explicit yes. Each unfinished task is then listed in `CLOSED.md` as not done, and
`Outcome` is `partly delivered` or `stopped`. With every task `DONE` (or no plan at all), `Outcome`
is `delivered`.

Justification: a stopped design is still a useful reference ("we tried X, it stopped here, for this
reason"). If close refused unfinished work, the only exit would be `/arcdlc:remove`, which recreates
the data loss this initiative fixes.

Trade-off: closed designs include work that was not finished, and `Outcome` in `CLOSED.md` is the
only thing that tells a reader which is which. Rejected: a strict close that needs every task
`DONE`.

### H4. Close deletes the queue, its board stories and both registers in place; everything else stays

Deleted: `plan.md`, `plan-archive.md`, `plan-human.md`, `gap.md` and `comments.md`. Everything else
in the folder, every architecture document, image and other file, stays untouched.

Justification: `plan-human.md` is the plan rewritten as board stories by `/arcdlc:plan-human`, so it
goes with the plan. `gap.md` and `comments.md` are the working notes of analysis
(`/arcdlc:examinate`, `/arcdlc:assist`). Anything important they found has already changed the
architecture document or an ADR. What is left in them is too detailed and too noisy for a reference,
and git history keeps it. OpenSpec, AI-DLC and SPDD all keep their plans after finishing; this
design deletes them on purpose, and `CLOSED.md` keeps what matters from them (the counts, the
unfinished tasks, the markers).

Trade-off: `comments.md` can be the only copy of a marker's text once `--strip` removed the comment
from the code. For a `stopped` initiative, that text leaves the working tree and is recoverable only
through git, except what R3 copies into `CLOSED.md`. The confirmation (H3) names every deleted file
so the engineer accepts that loss knowingly. Rejected: keeping both registers.

### H5. Close and remove do not read or rewrite links

A link to one of the five deleted files, or into a folder `/arcdlc:remove` deletes, is left as it is
and may break. Neither skill searches for links.

Justification: nothing moves on close, so the only links that can break point at the plan and its
registers, which designs and ADRs rarely link to, or at a design someone chose to delete. KISS wins
over a link rewriter for those cases.

Trade-off: a broken link is found by whoever clicks it; git history holds the target. Rejected: an
`arctool relink` command rewriting links to a `(removed at close on <date>, see git history)` note,
designed earlier in this interview (H13).

### H6. `/arcdlc:close` does the whole close itself; `arctool` has no close command

`/arcdlc:close <slug>`:

1. Reads the task counts from `docs/aics/<slug>/plan.md` (`arctool list --aic <slug>` when present,
   else by counting `- Status:` lines) and the archived count from `plan-archive.md`.
2. Stops, changing nothing, while any task is `TAKEN` (H3).
3. Shows the blast radius: the counts, each unfinished task by ID and title, the files that stay,
   and the files that will be deleted.
4. Asks the outcome questions when H7 needs them, then asks for an explicit yes.
5. Writes `docs/aics/<slug>/CLOSED.md` (H2).
6. Deletes the five files with `git rm` (`rm` when untracked), then runs `arctool sync` (by hand
   without it).

Justification: every step is a read, one small file written, and a delete. None needs a new tool, so
the initiative ships as skill and documentation changes only.

Trade-off: nothing mechanical guarantees the `CLOSED.md` shape or the counts; the skill's text does.
A crash between steps 5 and 6 leaves a half-closed folder; a re-run finishes it, because step 1 sees
`CLOSED.md` beside a leftover file and goes straight to the confirmation and the deletion. Rejected:
an `arctool close` command doing the checks and writing the note, designed earlier in this
interview.

Reversals, all in this interview: the first version moved the folder to
`docs/history/<date>_<slug>/` and re-linked every reference; the next had `arctool close` check and
write the note in place. The engineer then chose no `arctool` part in close or remove, so the
initiative stays small; its `arctool` changes are the status listing (H8) and the `sync` rule (H9).

### H7. The outcome is derived when it can be, and asked through the grilling protocol when it cannot

With every task `DONE`, or no plan at all, the outcome is `delivered` and nobody is asked. With any
task `TODO` or `BLOCKED`, the engineer picks `partly delivered` or `stopped` and gives one line of
reason, which goes into `CLOSED.md`. The skill never writes `delivered` while an unfinished task
exists.

Every question `/arcdlc:close` puts to the engineer (the outcome, the reason, the confirmation
before deletion) runs through `arcdlc-grilling` (`/arcdlc:grilling`), in the same style as
`/arcdlc:aic`: one question per turn, each with a recommended answer and one line of why. Following
the bundle rule that no skill hard-depends on another, the skill invokes the sibling first and falls
back to reading `../grilling/SKILL.md` (flat installs: `../arcdlc-grilling/SKILL.md`) and running
the protocol inline. `close` therefore joins the skills that restate the shared one-question
paragraph for that fallback.

Rejected: always asking, even when every task is `DONE` (a question with only one honest answer),
and always deriving the outcome from counts (loses the reason, which is the most useful line for the
next reader).

### H8. Status is read from the files by one read-only `arctool` command, never stored

The phase of an initiative follows from its files:

- `designing`: no `plan.md` and no `CLOSED.md`.
- `in progress`: `plan.md` has at least one task that is not `DONE`.
- `ready to close`: `plan.md` exists and every task is `DONE`.
- `closed`: `CLOSED.md` exists, with the outcome read from it.

`arctool status [--json]` prints one row per folder in `docs/aics/`, open and closed together: the
slug, the phase, the task counts (`TODO`, `TAKEN`, `BLOCKED`, `DONE`, with tasks in
`plan-archive.md` counted as `DONE`), and for a closed initiative the close date and outcome read
from `CLOSED.md`. It takes no `--aic`, only reads, exits 0, and prints only the header when there
are no initiatives. Besides the `sync` rule of H9, it
is the only `arctool` change; close, remove and the refusal of a closed design stay in the skills.

Justification: a status computed when asked cannot drift, and nothing new is written. The engineer's
KISS rule, no new index files, holds. OpenSpec's `openspec list` derives status from task progress
the same way; AWS AI-DLC keeps an `intents.json` registry, the kind of index file ruled out here.

Trade-off: someone reading `README.md` on GitHub sees the initiatives, not their phase. Rejected:
writing the phase into the registry lines (changes `AGENTS.md` on the first and last task of every
initiative and ties `/arcdlc:plan` and `/arcdlc:execute` to `arctool sync`), and no command at all
(status only by reading folders by hand), and reusing `arctool list` with no `--aic` (one command
meaning two things depending on a flag, where today it exits 2).

### H9. The registry lists open initiatives only; closed designs are read at design time

`arctool sync` skips every folder in `docs/aics/` that holds `CLOSED.md`, so a close removes the
initiative from the registry in `AGENTS.md` and `README.md` (its `sync` step, H6). When at least one
closed initiative exists, the block ends with one line, after a blank line: `N closed initiatives:
run arctool status, or see CLOSED.md in each folder.` (`1 closed initiative: ...` for one). Without
`arctool`, the skill deletes the bullet and writes the line by hand.

`/arcdlc:aic` Step 1 gains one item: read the title and outcome in each `docs/aics/*/CLOSED.md` and
open the related designs, so a closed initiative is used as a reference when a new design starts.

Justification: the registry is where people and agents look, so it must not show finished work as
active, and `AGENTS.md` is loaded into every agent session, so a list that grew with every close
would cost context forever. The line changes only on close, which already runs `sync`.

Trade-off: this is a second `arctool` change, inside `sync`, beyond the one command of H8; the
engineer accepted it because close must remove the initiative from `AGENTS.md`. Rejected: leaving
`sync` alone (closed designs listed as active) and marking closed bullets with a date (the list
grows without limit).

### H10. A closed initiative is final: no skill changes it, and no skill deletes `CLOSED.md`

`/arcdlc:aic`, `/arcdlc:plan`, `/arcdlc:examinate` and `/arcdlc:assist` stop, changing nothing, on a
folder holding `CLOSED.md`, and say that the initiative was closed on its date. No skill deletes
`CLOSED.md`. The one write a closed folder ever takes is a `## Superseded` line appended to
`CLOSED.md` by another initiative's design round (H14); the design documents never change. Follow-up
work on the same area is a new initiative with its own slug, whose design links the closed one as
its starting point. The slug of a closed initiative stays taken, because its folder stays.

Justification: a finished initiative is not changed without an explicit decision, and the engineer
chose the simplest form of that rule: no command changes it at all, rather than every command asking
for a yes. A closed design stays exactly as it was when the work ended, which is what makes it a
trustworthy reference.

Trade-off: two designs can describe the same area, the closed one and its follow-up; H14's line in
the older `CLOSED.md` is how a reader of it learns that. Rejected: every command reopening after an
explicit yes, and a single door back in through `/arcdlc:aic` that deleted `CLOSED.md`.

Reversals, both in this interview: the first design made a closed slug final and handled a later
initiative with the same name through dated `docs/history/<date>_<slug>/` folders; the next let
`/arcdlc:aic` reopen a closed design by deleting `CLOSED.md` on the engineer's yes. The engineer
then ruled that `/arcdlc:aic` must act as `plan`, `examinate` and `assist` do.

### H11. `/arcdlc:remove` deletes a design folder for good, after a confirmation

`/arcdlc:remove <slug>` deletes `docs/aics/<slug>/` and runs `arctool sync`, after it shows the
title, the phase (H8) and the file list, and gets an explicit yes. The confirmation states that the
design leaves the working tree and stays only in git history. There is no flag that skips it.
It does not rewrite links into the folder (H5), and its confirmation says so plainly: other files
may still reference the removed design, and those references will point at nothing. The engineer
decides with that known, and can cancel to handle the references first.

It refuses only while a task is `TAKEN`, because an agent is working on it at that moment, the same
rule as close (H3). An initiative with `TODO` or `BLOCKED` tasks can be removed: the skill shows the
unfinished counts loudly before the confirmation. Rejected: refusing while `plan.md` exists, which
would force a close whose `CLOSED.md` is deleted a moment later.

Reversals, all in this interview: the first round kept `remove` for open initiatives; the second
retired it; the third limited it to `docs/history/`. With the in-place close there is no
`docs/history/`, and `remove` returns to deleting a design folder, as the engineer first asked.
[ADR-0028](../../adr/0028-a-finished-initiative-is-closed-in-place.md) records the result and
supersedes ADR-0003.

### H12. In a workspace, close lands as one hub commit

`CLOSED.md`, the five deletions and the registry refresh go into one hub commit,
`docs(<slug>): close the initiative` plus `#AI-assisted`, pushed at once and retried once after
`git -C docs pull --ff-only`. The skill deletes with `git -C docs rm aics/<slug>/<file>`. Product
repositories are never read or written, as
[ADR-0026](../../adr/0026-a-workspace-keeps-its-docs-in-a-sibling-repository-named-docs.md)
requires.

Justification: every other skill already follows this rule in a workspace.

### H13. `arctool relink` is withdrawn

Designed earlier in this interview as one link command for close, remove and a person. Withdrawn on
2026-09-25: nothing moves, and close and remove do not handle links (H5).

### H14. A closed initiative that a later design reverses gets a line in its `CLOSED.md`

When a `/arcdlc:aic` round of one initiative reverses a decision that a closed initiative's design
records, the same round appends one line to the closed initiative's `CLOSED.md`, under a
`## Superseded` heading it creates when missing:
`- Superseded in part by [<title>](../<name>/aic.md) on <YYYY-MM-DD>: <which sections>.` A design
replaced entirely gets `- Superseded by [<title>](../<name>/aic.md) on <YYYY-MM-DD>.` instead. The
closed design documents are never touched, and no `/arcdlc:aic` run is started on the closed
initiative. `/arcdlc:aic` Step 1 already reads `docs/aics/` and every `CLOSED.md` (H9), which is how
the round finds the design it reverses.

When the reversed design belongs to an initiative that is still open, the round does not write into
it: two open designs that disagree are a decision for the engineer, and the round raises it as a
question through the grilling protocol.

Justification: without it, a reader of an old design has no signal that a later one changed it,
and "`docs/aics/` holds the actual design" is untrue from the first reversal on. `CLOSED.md` is the
closed initiative's record of what happened after the design was decided, so a later reversal
belongs there, and the design documents stay exactly as decided, keeping H10's rule whole.

Trade-off: the signal sits in `CLOSED.md`, one click from the design, not in the design itself.
Rejected: a line under the design's summary (edits a closed design), an exception to H10 asked with
an explicit yes, a forward reference in the new design only (the old one reads as current), and
dropping the rule (reversals would live only in ADRs).

Reversal: the first answer wrote the line under the older design's summary. Once H10 made a closed
initiative final, the engineer moved it into `CLOSED.md`, which is the one file in a closed folder
that records what happened after the close.

## Assessment

### Technical Challenges and Risks

- **R1. Links broken by a deletion.** Links to the five deleted files, or into a folder
  `/arcdlc:remove` deletes, are not rewritten (H5). Likelihood low (designs and ADRs rarely link a
  plan or a register; remove is rare). Impact low (a dead link; git history holds the target).
  Mitigation: accepted by the engineer on 2026-09-25, to keep close and remove free of link
  handling.
- **R2. Routing: a request to finish an initiative lands on `/arcdlc:remove`.** "We are done with
  X", "retire that initiative" and "delete the old plan" route to `/arcdlc:remove` today, and no CI
  can test routing. Likelihood high (the phrases are in `remove`'s description and in people's
  habits). Impact medium (a wrong route deletes a design, though only after its confirmation).
  Mitigation, reduce by design with a safe fallback: the finishing phrases move to `close`'s
  description ("we are done with X", "wrap up", "close the initiative", "mark it finished", "retire
  that initiative"); `remove`'s description narrows to deleting a design for good ("delete the
  design", "purge the initiative"); and `docs/routing-checks.md` gains rows for both skills, re-run
  by hand after release as `AGENTS.md` requires. Rejected: relying on the refusal message alone.
- **R3. Markers left in code after a `stopped` close.** An unfinished task that came from
  `/arcdlc:assist` (`ARCDLC-CMT-*`) may leave its `ARCDLC` comment in the code, while `comments.md`,
  the register that explained it, is dropped (H4). A later `/arcdlc:assist` sweep finds the marker
  again and judges it from scratch. Likelihood low (it needs an unfinished `assist` task and an
  unstripped marker at close time). Impact low (only the earlier judgement is lost). Mitigation,
  reduce by design: for each unfinished `ARCDLC-CMT-*` task, the not-done list in `CLOSED.md` (H3)
  carries the register block's `- Marker:` text and `- WHERE:` locations, copied as they stand, so
  the skill reads `comments.md` as well as `plan.md` before it deletes them. Rejected:
  accepting the loss with git history as the only copy.
- **R4. A closed initiative is changed anyway.** Nothing mechanical stops a skill, or an agent
  outside one, from writing into a folder that holds `CLOSED.md`. Likelihood low (four skills check
  first). Impact medium (a new `plan.md` beside `CLOSED.md` leaves the folder in two states).
  Mitigation, reduce by design: `aic`, `plan`, `examinate` and `assist` check for `CLOSED.md` before
  writing and stop (H10). The cost: one paragraph in four `SKILL.md` files, kept in step by hand.

### Open questions

None. Every section was answered in the interview of 2026-09-25; nothing was deferred.

| Question | Who answers | By when (phase or stage) |
| --- | --- | --- |
