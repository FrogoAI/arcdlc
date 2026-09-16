# Architecture decision records

Every decision that was hard, surprising, or expensive to reverse. Newest last.

| # | Decision | Status |
|---|---|---|
| 0001 | [Initiative selection is always explicit](0001-initiative-selection-is-always-explicit.md) | Accepted |
| 0002 | [`arctool sync` maintains the initiative registry via HTML-comment marker blocks](0002-registry-sync-via-marker-blocks.md) | Accepted |
| 0003 | [Initiative removal is a skill-side operation; arctool stays non-destructive](0003-initiative-removal-by-skill-not-arctool.md) | Accepted |
| 0004 | [`/arcdlc:plan` enforces risk-mitigation coverage before handoff](0004-plan-enforces-risk-mitigation-coverage.md) | Accepted |
| 0005 | [Antigravity support via a native plugin bundle with a flat-skills fallback](0005-antigravity-support-via-plugin-with-flat-fallback.md) | Accepted |
| 0006 | [Cursor support via flat personal skills](0006-cursor-support-via-flat-personal-skills.md) | Accepted |
| 0007 | [The MDCA standard is normative; the reference library ships one MDCA document](0007-mdca-standard-is-normative-single-document.md) | Accepted |
| 0008 | [Domain imports are judged by purity, not by folder](0008-domain-imports-judged-by-purity-not-folder.md) | Accepted |
| 0009 | [`Arc42.md` must generate a full arc42 document on its own](0009-arc42-source-is-self-sufficient.md) | Superseded by ADR-0010 |
| 0010 | [arc42 ships as the plain upstream skeleton plus one guide that carries all the help](0010-arc42-plain-template-plus-guide.md) | Accepted |
| 0011 | [Tech Stack Canvas is a first-class `/arcdlc:aic` format, and format arguments compose](0011-tsc-is-a-first-class-format-and-formats-compose.md) | Accepted |
| 0012 | [HTML output is a `:html` suffix on the format token, not a separate format](0012-html-output-via-format-suffix.md) | Accepted |
| 0013 | [`arctool order` permutes slots, it does not move tasks](0013-order-is-a-slot-permutation.md) | Accepted |
| 0014 | [Exit code 5 means "self-validation failed", not "archive self-validation failed"](0014-exit-5-means-self-validation-failed.md) | Accepted |
| 0015 | [The comment register is scanned by `arctool` and judged by the skill](0015-comment-register-is-scanned-by-arctool-and-judged-by-the-skill.md) | Accepted |
| 0016 | [Comment markers group by tag, and their block can grow](0016-comment-markers-group-by-tag.md) | Accepted |
| 0017 | [A marker is `ARCDLC` inside a comment, and cutting it out is asked for](0017-a-marker-is-arcdlc-in-a-comment-and-cutting-it-is-asked-for.md) | Accepted |
| 0018 | [Only single-line comments carry markers](0018-only-single-line-comments-carry-markers.md) | Accepted |
| 0019 | [The executor tier is asked for, not guessed](0019-the-executor-tier-is-asked-for-not-guessed.md) | Accepted |
| 0020 | [The reference library is dissolved into the skills that use it](0020-reference-library-dissolved-into-the-skills.md) | Accepted |

A decision that changes an earlier one carries `Amends` or `Supersedes` in its header, and the
earlier record points forward. Nothing here is deleted: a superseded ADR is how you find out why the
current rule exists.
