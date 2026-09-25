# Routing checks

A skill only fires on words in its `description`. Nothing in CI can test that, because testing it
means running a real agent. This file is the manual substitute: say the prompt, see which skill the
agent reaches for, compare against the expected column.

Re-run it after any change to a `description`. `stoic.md` is why this file exists: its content was
fine, its trigger words were misspelled and missing from the frontmatter, and it was never opened
once.

| Prompt | Expect |
|---|---|
| we've got TODOs scattered all over the code, turn them into tickets | `assist` |
| clean up the FIXMEs in internal/ | `assist` |
| what is marked unfinished in the code? | `assist` |
| design the architecture for our payments service | `aic` |
| spec out the system before we build it | `aic` |
| write an arc42 doc for checkout | `aic` |
| break this down into tasks | `plan` |
| what are the steps to build this? | `plan` |
| turn the architecture into work | `plan` |
| start implementing the plan | `execute` |
| do the next task | `execute` |
| pick up the backlog and ship it | `execute` |
| check we follow DDD | `examinate` |
| find our tech debt | `examinate` |
| where are we off-architecture? | `examinate` |
| review the codebase for compliance | `examinate` |
| make Jira tickets from the plan | `plan-human` |
| put this work on the board | `plan-human` |
| I need something to hand a developer | `plan-human` |
| write our on-call policy | `policy` |
| draft an SOP for incident response | `policy` |
| we need a standard for code review | `policy` |
| grill me on this idea | `grilling` |
| poke holes in my design before I build it | `grilling` |
| ask me what I am missing | `grilling` |
| the plan is getting long | `archive` |
| clean up the finished tasks | `archive` |
| we are done with the payments initiative, close it | `close` |
| wrap up the checkout work | `close` |
| mark init as finished | `close` |
| we are done with cursor-support, retire it | `close` |
| delete the payments design for good | `remove` |
| purge the old spike initiative | `remove` |
| set up arcdlc in this repo | `init` |
| we have five repos, set up a docs hub | `init` |
| bootstrap the project for arcdlc | `init` |

## Known collisions

- **"tickets"** pulls both `assist` and `plan-human`. Disambiguated by source: "tickets from the
  code" or "from the TODOs" is `assist`, "tickets from the plan" is `plan-human`. If a prompt names
  neither source, ask rather than guess.
- **"build it"** is `execute` when a plan exists and `aic` when nothing is designed yet. The skill
  that fires should check for `docs/aics/<slug>/plan.md` before assuming.
- **"set up"** pulls both `init` and `policy`. "Set up" with a policy word ("set up our on-call
  process") is `policy`; "set up" naming a repository or a workspace is `init`.

## Rules for writing a description

1. Lead with what the user wants done, in their words, not with what the skill does internally.
2. List the phrasings a person actually types. "find our tech debt" routes; "gap analysis" alone
   does not.
3. Keep mechanics out. The slug being the first positional argument, the output paths, and the
   format list belong in the body, which is read only after the skill fires.
4. Name both trigger forms, `/arcdlc:<name>` and `arcdlc-<name>`. `internal/bundle` fails the build
   if either is missing.
