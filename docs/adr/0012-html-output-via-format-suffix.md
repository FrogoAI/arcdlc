# ADR-0012 — HTML output is a `:html` suffix on the format token, not a separate format

- Status: Accepted
- Date: 2026-08-31
- Initiative: [source-library-cleanup](../aics/source-library-cleanup/aic.md)

## Context

arc42 publishes its plain template in several renderings. The bundle carries the Markdown one; the
Asciidoctor HTML one (`arc42-template.html`) covers the same twelve sections and differs in three
ways that matter:

- it is **already numbered** (`1.`, `1.1.`, and one level deeper for repeat slots: `5.1.1`, `7.2.1`),
- it carries a table of contents with 31 anchor links that must stay consistent with the headings,
- its `<h1>` is the arc42 logo image, and `<title>` is the literal string `Template`.

Engineers asked for it by output shape — "arc42, but in HTML" — not as a different document format.
The architecture, the sections, and the interview behind them are identical.

## Decision

- **HTML is a modifier, not a format.** The format token takes a `:html` suffix:
  `/arcdlc:aic payments arc42:html`. It changes the template read and the output extension, nothing
  else. It composes with the comma list (`arc42:html,tsc`). Markdown stays the default for every
  format. The skill also accepts the natural phrasings ("arc42 in html", "arc42 html") and normalises
  them to the same token.
- **One Section Catalogue serves both.** `Arc42 Guide.md` keeps a single per-section *Write* / *Shape*
  catalogue and splits only the mechanical procedure into two, because only the markup differs. The
  slot decoder gains an HTML half: `<em>&lt;…&gt;</em>` slots, `<dt class="hdlist1">` content headings
  (the same white box template as Markdown's definition lists), `<div class="sect3">` repeat slots,
  and `<a class="anchor">` elements that the TOC depends on.
- **The generated HTML must replace the logo `<h1>` and the `<title>`.** This serves the registry
  contract and removes the file's only `<img>` in one step, so a generated document has no image
  dependency and `images/` is never copied into an initiative folder.
- **`arctool` reads HTML architecture documents.** `findArchDoc` accepts `.html`, and the registry
  parses the first `<h1>` as title and the first `<p>` after it as summary, with tags stripped and
  entities unescaped. Precedence is format rank first, `.md` before `.html` within a format, then
  alphabetical `.md`, then alphabetical `.html`. Parsing is stdlib-only (`strings` + `html`).

## Justification

- **A suffix keeps the format axis honest.** arc42-in-HTML is arc42. Making it a separate format
  token would double every format row, and `arctool sync`'s precedence list would have to rank a
  rendering against a format, which is not a meaningful comparison.
- **Without registry support the feature would be a trap.** `findArchDoc` filtered on `.md`, so an
  HTML-only initiative would have registered as "(no architecture doc)" — silently, since sync
  reports no error for it.
- **`.md` before `.html`** because a Markdown document is the one the rest of the pipeline reads
  cheaply; when both exist for the same format, they are the same content and Markdown is the better
  input for `/arcdlc:plan`.

## Trade-offs

- **A hand-rolled HTML text extractor.** No `golang.org/x/net/html` — that would break the
  stdlib-only rule — so tag stripping is a small state machine. It is deliberately minimal and would
  mis-handle a `<` inside an attribute value. Accepted: it reads documents this bundle's own skills
  generate from a known template, and it is covered by table-driven tests including the
  logo-left-in-place case.
- **The TOC is the agent's job.** Deleting a section means deleting its `<li>`, and nothing checks
  that. Accepted: the guide makes it an explicit step; a stale anchor is a visible broken link, not a
  silent data error.
- **Two templates to refresh per arc42 release** instead of one. Accepted: each is still a verbatim
  copy-paste, and the shared Section Catalogue does not change.
