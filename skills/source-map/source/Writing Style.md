# Writing Style: documents a person wants to read

## What this is

Every file ArcDLC writes is read by a human. A new engineer opens it on their second day and
has to understand the system without asking anyone. That is the bar.

Machine-flavoured prose fails that bar. It is smooth, long, evenly balanced, and you finish a
paragraph without learning one fact you can act on. This file is the writing standard that keeps
it out.

It applies to everything the ArcDLC skills produce: architecture documents (AIC, arc42, Tech Stack
Canvas, TOGAF, C4), ADRs, `CONTEXT.md`, policies under `docs/policies/`, plan tasks, gap findings,
commit messages, and the HTML variants of all of them.

Two rules stand above the rest:

1. Write plain English. Short sentences, common words, a named actor.
2. Never trade content for simplicity. Plain words, full content, always.

## Scope: what this does not touch

Do not "simplify" any of these:

- **Code and identifiers.** Paths, flags, config keys, type names, slugs. `--plan`, `kebab-case`,
  `plan-format.md` stay exactly as they are.
- **Lines a format contract owns.** The initiative registry line
  `- [Title](docs/aics/<slug>/aic.md) — summary` uses a long dash by contract, and `arctool sync`
  writes it. The plan format's keys are fixed too. Contract beats style.
- **Headings and field names a template defines.** arc42 section 5 is "Building Block View" even
  though that is a stacked noun phrase. Rename nothing the template names.
- **Quoted material.** Error text, log lines, a standard's wording, the user's own words. Quote
  exactly, then explain in your own plain words.
- **Domain terms.** See "Keep the domain words" below.

## Delete the filler

**Transitions and connectives.** Cut them. Start the sentence with its subject.

> Furthermore, Moreover, Additionally, In addition, That said, That being said, In conclusion,
> To summarize, Ultimately, Overall, It is important to note, It is worth noting, It should be
> noted, As mentioned earlier, Needless to say, At the end of the day.

**Openers that say nothing.** A document never warms up. First sentence, first fact.

> In today's fast-paced world, In the ever-evolving landscape of, As we navigate, At its core,
> Let's dive in, Simply put, In this document we will explore, This section aims to.

**Vocabulary that marks a machine wrote it.**

| Do not write | Write |
| --- | --- |
| delve into | look at, dig into |
| leverage, utilise, utilize | use |
| facilitate, enable (as filler) | let, allow, do |
| robust, seamless, cutting-edge | name what it does |
| holistic, comprehensive | drop it or list what is covered |
| myriad, plethora, a wealth of | a number, or the number |
| tapestry, beacon, realm, journey, landscape | the real noun |
| underscore, testament to, pivotal, crucial, vital | important, or say why |
| in order to | to |
| a variety of, a range of | the actual list |
| honestly, genuinely, truly, really | delete the word |
| clearly, obviously, essentially, basically | delete the word |

**Empty intensifiers.** Cut `honestly`, `genuinely`, `truly`, `really`, `simply`, `clearly`,
`obviously`, `essentially`, `basically`. Every one of them adds emphasis where a fact belongs,
and the sentence is stronger with the word gone. Two earn a special mention:

- `honestly` implies the sentences around it were not honest. Say the thing and let it stand.
- `genuinely` claims sincerity instead of showing evidence. "A genuinely hard trade-off" tells
  the reader nothing. "Both options cost us something: Redis adds a service to run, an
  in-process cache loses us multi-node reads" tells them everything.

`clearly` and `obviously` do extra damage. When the point is not obvious to the reader, the word
tells them they are slow. Delete it and explain instead.

One exception: `actually` may stay when it marks a real contrast, as in "code you actually read"
(read, not assumed) or "what the diff actually does" (not what the task title claimed). If it is
not doing that job, it is filler.

**Praise and hedging.** No "great question", no "this powerful approach", no "we believe this
elegant solution". A document states decisions and reasons. Nothing else.

**Invented endings.** If the template asks for a summary section, write it. Do not bolt a
concluding paragraph onto a section that did not ask for one, and never end with a restatement of
what the reader just read.

## Delete the AI sentence shapes

These patterns read as generated even when every word is fine:

- "It is not just X, it is Y." Also "This is not about X. It is about Y."
- Three of everything. Three adjectives, three clauses, three bullets, every time.
- Rows of twins: five sentences in a row of the same length and the same build.
- Perfectly mirrored clauses: "Fast to write, easy to read, simple to maintain."
- A rhetorical question as an opener: "So what does this mean for the team?"
- Every bullet in a list built to the identical mould.
- Bold on half the words in a paragraph.
- Emoji, unless the template already uses them.

## No long dashes

Do not use `—` (em dash) or `–` (en dash) as punctuation. They are the loudest tell there is.

Replace one with the punctuation that fits the thought:

- A full stop, when the two halves are two thoughts. Usually the right answer.
- A comma, when it is an aside.
- A colon, when what follows explains what came before.
- Brackets, when the text is an aside the sentence could drop.

```
Bad:  We store sessions in Postgres — we already run it — and the volume is small.
Good: We store sessions in Postgres. We already run it, and the volume is small.

Bad:  The parser is strict — every task needs an Acceptance section.
Good: The parser is strict: every task needs an Acceptance section.
```

Hyphens are untouched. `install-agnostic`, `--strict`, `arcdlc-plan`, `2026-09-04` are all fine.
So is a long dash inside a line a format contract owns (see "Scope" above).

## Vary the rhythm

Real writing is uneven. Generated writing is not.

- Mix lengths on purpose. A five-word sentence next to a twenty-five-word one.
- Let a paragraph be one line when one line is the whole point.
- Break the pattern before the third bullet that starts the same way.
- Read the paragraph out loud in your head. If it sounds like a metronome, rewrite it.

## Say who does what

- Active voice, present tense, a named actor. "The API writes to Postgres", not "Persistence is
  realised via the relational store".
- Say the thing first, the reason second. "We keep sessions in Postgres because we already run it."
- Break stacked nouns into a sentence a person can say out loud. "Install-agnostic cross-agent
  skill bundle distribution" becomes "One bundle installs on every supported agent."
- Contractions are fine. Use them where they sound natural, skip them where they do not.
- Keep paragraphs under about five lines. Past that, use a list, a table, or a diagram.

## Concrete beats abstract

Every vague claim is a place where a fact is missing.

```
Bad:  This significantly improves performance.
Good: This cuts p99 latency from 800 ms to 120 ms.

Bad:  We considered several alternatives.
Good: We looked at Redis and at an in-process cache. Redis lost on operational cost.

Bad:  Robust error handling is in place.
Good: A failed write retries three times, then goes to the dead-letter queue.
```

Numbers, names, paths, versions, dates. If you do not have one, say what is unknown and put it in
"Open questions" instead of dressing it up.

## Keep the domain words

Simple English is about the sentence, not the vocabulary. A precise term carries meaning that a
common word throws away, so it stays.

Keep terms like provenance, idempotent, backpressure, quorum, eventual consistency, tenancy, saga,
byte-preserving, atomic. Keep this project's own words too, whatever they are.

The rule is:

- Define the term once, in plain words, at first use in the document, or in `CONTEXT.md`. Then use
  it freely for the rest of the document.
- Expand every acronym on first use in each document.
- Never swap a domain term for a loose everyday word. "Provenance" is not "where the data came
  from, roughly". The definition is plain. The term is exact. You need both.
- Never drop a distinction because explaining it takes a sentence. That sentence is the document.

## Plain words, full content

This is the rule people break first. Simple English never means less material.

Keep every decision, constraint, trade-off, risk, dependency, path, exit code, and acceptance
criterion the template asks for. If a section gets shorter, it is because the filler went, not
because the content did.

Trace every claim to a source: an interview answer, an ADR, or code you actually read. If nothing
backs it, it belongs in "Open questions", not in the body dressed as a decision.

## The pass before you save

Run this on every document, every time. It takes a minute.

1. Read the first sentence of every section. Does each one carry a fact? Delete the warm-up ones.
2. Search the text for a long dash. Replace every one that is not contract-owned.
3. Scan for the ban list above. Rewrite each hit, or cut the sentence.
4. Check the shape: three sentences in a row of the same length is a rewrite.
5. Pick three claims at random. Can you name the source of each? If not, move or cut them.
6. Check the terms: is every acronym expanded once, and every domain term defined once?
7. Last read: would a new engineer act on this document without asking a question that the
   document should have answered?

A rough grep for the mechanical part, run from the repo root:

```bash
grep -nEi "—|–|\b(delve|leverage|utili[sz]e|robust|seamless|holistic|myriad|plethora|tapestry|beacon|realm|landscape|underscore|testament|pivotal|crucial|honestly|genuinely|truly|clearly|obviously|essentially|basically)\b|\b(furthermore|moreover|in conclusion|it is important to note|it'?s important to note|in today'?s|ever-evolving|at its core|needless to say)\b|not just .* but" docs/aics/<slug>/*.md docs/policies/*.md
```

Judge every hit. Some are legitimate: a quoted error message, a contract line, a domain term that
happens to be on the list. The grep finds candidates, you decide.

## Worked examples

**Architecture document**

```
Bad:  In today's rapidly evolving payments landscape, it is crucial to leverage a robust and
      seamless integration layer — one that facilitates comprehensive observability across the
      entire transaction lifecycle.

Good: Payments go through one integration service. It owns retries, idempotency keys, and the
      audit trail, so the rest of the system never talks to a provider directly. We accept the
      extra hop because provider outages then have one place to fail.
```

**Plan task**

```
Bad:  Goal: Enhance the validation subsystem to ensure robust handling of a variety of edge cases.
Good: Goal: Reject a task block that has no Acceptance section. `arctool validate --strict` exits 1
      and names the task ID.
```

**Gap finding**

```
Bad:  The current implementation exhibits significant deviations from the prescribed architectural
      guidance, which could potentially impact maintainability going forward.
Good: `internal/plan/mutator.go` writes the status line with a full file rewrite. The policy
      requires temp-file plus rename. A crash mid-write truncates the plan.
```
