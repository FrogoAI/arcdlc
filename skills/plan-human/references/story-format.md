# Story format: what a board ticket must hold

This is the contract for `docs/aics/<slug>/plan-human.md`. A story that follows it can be pasted into a tracker and worked on without asking anyone a question.

## The file

```
# <Initiative name>: engineer stories

<Two or three lines: what is being built, and where it runs.>

<One line naming the architecture document as the source of the reasoning,
 and the plan as the task-level detail.>

**Everything from the first story heading down is ticket text.**
<One line: a story is copied into its ticket as it stands, so it names no file
 in this repository and no ADR.>

<How many stories, in how many phases, which need another team.>

## What we are building
<A table of the deployable units and one line each on what they do.>

## Where everything runs
## Out of scope
## Rules the code must follow
<The limits that come from tools and decisions nobody in this initiative can
 change. Numbered, so a story can say "rule 6".>

---

# Phase 1: <what this phase delivers, in the reader's terms>

## Story 1: <title>
...

---

# Phase 2: <what this phase delivers>

## Story N: <title>
```

The heading section is repository-only. Everything below the first story heading is ticket text.

## The title

An imperative verb plus its object, naming the real component.

| Do not write | Write |
| --- | --- |
| One message out, and the three places it goes | Implement the alerter service |
| See it working, and prove it holds at full speed | Run the phase 1 load test at 6 000 bets per second |
| Turn the counters into something a person and a machine can both read | Implement the alerter service |
| The calculator: from a bet to a counted alert | Implement the calculator service |
| Settings and channels in Backoffice | Add the Backoffice settings screen and the channel owner field |

Never a metaphor, a narrative line, or a question. Read the title with the body covered: if you cannot say what changes and where, rename it.

## The six parts

In this order. Skip a part only when it does not apply.

```
## Story <n>: <title>

*Tracker: <ID>. Plan tasks: <IDs>.<  Needs DevOps. | Touches <repo>.>*

<Two or three lines, present tense: what the thing does when it works.>

**WHAT TO DO**

1. <One instruction per step. Imperative. Numbers and names, not adjectives.>
...
N. <The last steps add the metrics and the tests.>

**<THE DATA SHAPE>**            <- only where there is one

<A fenced block with the real shape, and a line on why a field is there.>

**EXAMPLE**

<One worked case on concrete values. Real numbers, real ids.>

**ACCEPTANCE CRITERIA**

- I <do this>, I <get that>.
...

**OUT OF SCOPE**

<What a reader would reasonably assume is included and is not, and which
 story covers it instead.>
```

### The pointer line

`*Tracker: SC-1367. Plan tasks: CGH-08.*` It links the story to its ticket and its tasks, and it is the one line that is **not** copied into the ticket, because the ticket is the thing it points at.

### WHAT TO DO

- One instruction per numbered step. An engineer works down the list.
- State required behaviour as fact: "The alerter deletes the records it took, then publishes." Not "should consider deleting".
- Where a decision looks like a bug, say why it is not, in the step: "Deleting first is deliberate and must not be fixed: a failed publish loses that batch."
- Put the numbers in: an interval, a cap, a limit, a timeout, a TTL.
- End with the tests. A story whose steps never mention tests will not get any.

### The data shape

Where a story produces or consumes a structure, show it. A JSON body, a record layout, a table of fields, a table of metrics with names, types and labels. Then one line on any field a reader would otherwise drop as redundant.

This is the part most often replaced with prose, and prose does not survive: "expose every metric the design depends on" cannot be implemented, and a dashboard cannot reference it.

### EXAMPLE

One worked case on real values, showing input and result together.

> 30 records are waiting: 10 for organization A, 14 for organization B, 6 for organization C. One run publishes 3 messages, one per organization, with 10, 14 and 6 elements.

An example is how a reader checks they understood the steps. A design that cannot be shown on one case is not ready.

### ACCEPTANCE CRITERIA

Written as input and output, from the point of view of whoever checks it.

```
Bad:  GIVEN 100 waiting records WHEN one run happens THEN 30 are taken and 70 remain
Good: I put 100 waiting records for one organization in the store and run one cycle,
      I get 30 taken and deleted and 70 left.

Bad:  The endpoint is disabled by default.
Good: I call the endpoint with the flag unset, I get 404 and no route is registered.

Bad:  Secrets are handled correctly.
Good: I post a body carrying password, key, purse and balance, none of the four
      values appears in the record, the answer, the log or any error text.
```

Every `Acceptance` line from the matching plan tasks appears here, reworded and never dropped. Add the ones a person can only check by looking, which a plan task may have left implicit.

### OUT OF SCOPE

Two or three lines. What a reader would assume is included, and which story has it instead. This is where a board stops two people building the same thing.

## Self-contained

A ticket lives where this repository does not exist.

- **Never** name a file in this repository, and never cite an ADR number. If the story needs the field mapping, the envelope shape, or the subject name, write it into the story.
- Say "the repository README" or "an operations document in the repository", not a path, when the story's output is a document.
- **Code identifiers stay.** `cmd/alerter`, `internal/decoder/`, `PARTNER_ID`, `cgh_alerts_pending`, a set name, a subject name. These are the instruction, not a pointer.
- Refer to sibling work by story number, not by task id: "which is story 7".

## Worked example

```
## Story 3: Implement the alerter service

*Tracker: SC-1367. Plan tasks: CGH-08.*

The alerter wakes on a timer every 5 minutes, reads up to 30 alert records,
deletes them, and publishes one envelope per organization to NATS.

**WHAT TO DO**

1. Create `cmd/alerter` as a **single copy** on a timer. Two copies would take
   the same records. The interval, the record cap of 30 and the size budget are
   all configuration.
2. Read `cgh_alerts_{partner_id}` with a cursor. Take at most the cap, and stop
   earlier if the envelope would cross the size budget.
3. Skip records below their alert's minimum count: do not read them and do not
   delete them. They expire by their own TTL.
4. Group what you read **by organization**. **Delete the records you took, then
   build and publish.** Deleting first is deliberate and must not be "fixed": a
   failed publish loses that batch. The same conditions fire again in minutes.
5. Publish to **`cgh.alerts.v1.<org_id>`**, one subject per organization.
6. Add metrics: `cgh_alerts_taken_total`, `cgh_envelopes_sent_total`,
   `cgh_alerts_in_store`, `cgh_alerts_pending`.
7. Render **no text**. Rendering belongs to an output service.
8. Add tests for the drain, the grouping and the envelope shape.

**THE ENVELOPE**

{ "cycle_id": ..., "org_id": ..., "generated_at": ...,
  "alerts": [ { "alert_name", "group", "group_id", "count",
                "first_seen", "variables", "triggers" } ] }

`group` carries the group fields as named fields, and `group_id` the same thing
as one readable string. Both are needed: a dashboard filters on `group.game_id`
without parsing `g:730`.

**EXAMPLE**

30 records are waiting: 10 for organization A, 14 for B, 6 for C. One run
publishes 3 messages, one per organization, with 10, 14 and 6 elements.

**ACCEPTANCE CRITERIA**

- I put 100 waiting records for one organization in the store and run one cycle,
  I get 30 taken and deleted and 70 left.
- I put records for three organizations and run one cycle, I get 3 envelopes on
  three subjects, and no envelope holds two organizations.
- I put a record below its minimum count and run one cycle, it is neither read
  nor deleted.
- I read a published envelope, I get all seven fields per alert.

**OUT OF SCOPE**

Telegram, which is story 7. Elasticsearch, which is story 8. Partner
destinations, which are story 16.
```

## The mapping table

Keep this in `docs/aics/<slug>/CONTEXT.md`, under the delivery section. It is how the next session finds which ticket holds which work.

```
| Story | Title | Plan tasks | Tracker | Phase |
| --- | --- | --- | --- | --- |
| 3 | Implement the alerter service | CGH-08 | SC-1367 | 1 |
| 4 | Configure the four alerts and prove each one fires | CGH-27, CGH-13 | SC-1361 | 1 |
```

Record what happened to the cut as well: which stories split, which tickets were created, which were retired and why. A board that changed shape without a written reason gets changed back.

## Checks before saving

1. Read every title with the body covered. Can you say what changes and where?
2. Every task id from `plan.md` appears in exactly one story.
3. Every plan acceptance criterion reaches a story, reworded into "I do this, I get that".
4. Every story has acceptance criteria.
5. Below the first story heading: no repository path, no ADR number.
6. No long dash anywhere in the file.
7. Every set of things (metrics, fields, endpoints, keys) is a named table, not prose.
8. Every story has an example, unless it produces nothing observable.
