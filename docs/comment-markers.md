# Comment markers: the syntax and the comment styles the sweep reads

`arctool scan` reads work notes out of code comments and writes them into an initiative's
`comments.md`. This page is the contract for what it looks for and where it can see it. The table it
comes from lives in `internal/scan/scan.go`, and a test fails when the two disagree.

## Single-line comments only

A marker counts in a single-line comment and nowhere else. A block comment (`/* */`, `<!-- -->`,
`(* *)`, `{- -}`) is prose to the sweep, whatever it holds:

```go
// ARCDLC move the index rebuild into internal/index   ← a marker

/* ARCDLC move the index rebuild into internal/index */  ← not a marker

/*
	ARCDLC move the index rebuild into internal/index      ← not a marker
*/
```

One rule instead of a family of them. A single-line comment ends where its line ends, so the note's
extent needs no per-language knowledge of closing tokens, nesting or decoration. Removing one is a
whole-line delete or a cut to the end of a line, and both are safe in every language. That is why
`arctool scan --strip` can never strand an opener or eat a line of code.

A language whose only comment is a block comment therefore carries no markers at all. Write the note
in a file that can hold one, or add a single-line comment style to the language (there is none to
add for HTML).

## The marker

The marker word is `ARCDLC`, and only `ARCDLC` by default. That is the point: a repository full of
`TODO` and `FIXME` notes is untouched until somebody asks for those, with `arctool scan --marker TODO`
or `/arcdlc:assist <slug> TODO`.

```go
// ARCDLC move the index rebuild into internal/index
// ARCDLC:T1 move the index rebuild into internal/index
```

Three rules decide whether a marker counts:

- **It opens the comment.** `// ARCDLC ...` is a marker. `// see ARCDLC above` is prose about markers.
- **It is a whole word.** `ARCDLCLIST` is a name.
- **It sits inside a comment.** `ARCDLC` inside a string or a long constant is text in that constant.
  The scanner follows string literals, including the ones that run past the end of their line, so a
  line of documentation inside a Go raw string is not a finding.

### Group tags

`ARCDLC:TAG` right after the marker word, ended by a space or the end of the line. Every marker with
the same tag becomes one block and one task, however many files it spans, because it is one change
written down in several places.

```go
// ARCDLC:T1 move the index rebuild into internal/index
// ARCDLC:t1 and return one value instead of two
```

- Tags ignore case: `:t1` and `:T1` are one group, recorded as `T1`.
- A tag reaches across the whole sweep, and belongs to one marker word: `TODO:T1` is a different group.
- A tag needs a letter. `// ARCDLC:01` would claim the auto-numbered block ID `ARCDLC-CMT-01`, so it
  reads as a plain marker and the run says so.
- The block is named after the tag: `ARCDLC-CMT-T1`. An untagged marker gets a number:
  `ARCDLC-CMT-01`.

## Where a marker may sit

```go
// ARCDLC one line

// ARCDLC the note continues on
// every comment line under it, so a three-line note is one task

func f() {} // ARCDLC trailing a line of code
```

A blank line, a line of code, or a second `ARCDLC` ends the run. Doubled and decorated openers are
read as the same style, and the padding in front of the marker word is not part of the note: `///`,
`//!`, `##`, `;;;`, `%%`.

## Comment styles, per file type

| Comment style | File types |
| --- | --- |
| `//` | `.c` `.cc` `.cjs` `.cpp` `.cs` `.d` `.dart` `.dpr` `.fs` `.fsi` `.fsx` `.go` `.gradle` `.groovy` `.h` `.hpp` `.java` `.js` `.json5` `.jsonc` `.jsx` `.kt` `.kts` `.less` `.mjs` `.pas` `.php` `.pp` `.proto` `.rs` `.scala` `.scss` `.sol` `.svelte` `.swift` `.ts` `.tsx` `.v` `.vue` `.zig` |
| `#` | `.awk` `.bazelrc` `.bash` `.bzl` `.cmake` `.cr` `.dockerignore` `.env` `.ex` `.exs` `.fish` `.gemspec` `.gitignore` `.gql` `.graphql` `.jl` `.mk` `.nim` `.nix` `.pl` `.pm` `.ps1` `.psm1` `.py` `.r` `.rake` `.rb` `.sh` `.tcl` `.toml` `.yaml` `.yml` `.zsh` |
| `#` (by file name) | `BUILD` `Brewfile` `CMakeLists.txt` `Dockerfile` `Doxyfile` `Fastfile` `Gemfile` `Justfile` `Makefile` `Podfile` `Rakefile` `Vagrantfile` `WORKSPACE` `justfile` `makefile` |
| `#` `!` | `.properties` |
| `#` `//` | `.hcl` `.tf` `.tfvars` |
| `--` | `.ada` `.adb` `.ads` `.applescript` `.elm` `.hs` `.lua` `.sql` `.vhd` `.vhdl` |
| `;` | `.clj` `.cljs` `.el` `.lisp` `.rkt` `.scm` `.ss` |
| `;` `#` | `.cfg` `.ini` |
| `;` `#` `//` | `.asm` `.s` |
| `%` | `.erl` `.hrl` `.tex` |
| `%` `//` | `.m` (MATLAB and Objective-C share the extension) |
| `!` | `.f` `.f03` `.f90` `.f95` `.for` |
| `'` | `.bas` `.vb` `.vbs` |
| `"` | `.vim` `.vimrc` |
| `::` `REM` | `.bat` `.cmd` (`REM` needs a space after it, so `REMOVE` is not a comment) |
| none | `.css` `.htm` `.html` `.ml` `.mli` `.xml`: block comments only, so no marker can be written |

### A file type that is not listed

It is read with `//` and `#`, and the run says so, so an empty result is never mistaken for a clean
tree:

```
note 5 file type(s) the sweep does not know (.iml, .json, .mod, .svg, LICENSE) were read with the
default // and # comment style
```

Documents are skipped entirely: `.md`, `.markdown`, `.txt`, `.rst`. So is anything under `.git`,
`vendor`, `node_modules`, `dist`, `bin` and `docs`, along with binary files and anything over 1 MB.

To add a language, add its extension to `openersByExt` in `internal/scan/scan.go` and to the table
above. `TestDocumentedCommentStylesMatchTheTable` fails when one of the two is missing.

## Removing the comments

`arctool scan` records and stops. It edits a source file only when asked with `--strip`, and
`/arcdlc:assist` only asks after the work is in the plan. Because every marker sits in a single-line
comment, there is no shape it has to refuse: a standalone comment line and the comment lines that
continue it are deleted, and a comment trailing code is cut off the end of its line.

One thing still stops a removal: a file edited between the sweep and the strip keeps its comment, and
the run reports it. Nothing is written anywhere when a check fails.
