# KISS Principle

## Meaning

In ArcDLC guidance, **KISS** means **Keep It Simple and Smart**.

KISS means efficient, smart, readable code - not simplistic code. Solve the
actual problem with the smallest design that remains clear, fast, modular,
testable, and ready to extend when the domain requires it.

KISS is not an excuse to under-design. It is a guardrail against accidental
complexity.

## What KISS Requires

- Solve the actual problem, not an imagined future problem.
- Keep the main path easy to read from input to output.
- Prefer domain names and clear responsibilities over generic names.
- Add abstractions only when they remove real complexity, protect a boundary,
  isolate volatility, or make testing clearer.
- Keep code efficient by design. Simplicity includes avoiding unnecessary work,
  allocations, network calls, database calls, and synchronization.
- Keep optional or additional behavior in independent modules when that keeps
  the core code smaller and easier to reason about.
- Make modules easy to append functionality to without forcing unrelated code to
  change.

## What KISS Rejects

- Adding `service`, `manager`, `processor`, `helper`, `common`, or `util`
  packages/types when they do not own a clear responsibility.
- Creating pass-through layers that only call another layer with the same data.
- Hiding complexity behind vague names instead of removing or isolating it.
- Splitting code into many files or packages before there is a real boundary.
- Collapsing distinct responsibilities into one large function or package in
  the name of "simplicity".
- Choosing a simplistic implementation that is slower, harder to read, harder to
  test, or harder to extend than a slightly more explicit modular design.

## Service and Manager Rule

Do not introduce a `Service`, `Manager`, `Processor`, or similar type by default.
Use one only when it has a cohesive responsibility that cannot be named more
precisely.

A layer is justified when it owns meaningful work, such as:

- enforcing a domain rule or invariant
- coordinating a real use case across multiple collaborators
- isolating an external dependency behind a stable boundary
- separating optional behavior from the core flow
- reducing unavoidable complexity in a way that makes the caller simpler

If the type only forwards data, renames another call, or exists because "every
project has a service layer", remove it or inline the call.

## Modular Simplicity

Simple code is modular code. A module should have a clear reason to exist, a
small API, and low coupling to unrelated behavior.

When additional functionality is needed, prefer an independent module, adapter,
or package if it lets the existing core stay readable and stable. Do not bury
optional behavior inside the core flow just to avoid one more module.

## Decision Checklist

Before adding a package, type, interface, or layer, answer these questions:

- What actual problem does this solve now?
- Can a reader follow the main flow without jumping through empty layers?
- Does the abstraction remove real duplication or isolate a real boundary?
- Is the performance behavior obvious and efficient enough for the use case?
- Will the design stay readable when the next expected feature is added?
- Can the new functionality be isolated in an independent module instead of
  making the core more complex?

If the answer is unclear, keep the implementation direct and name the real
responsibility first.
