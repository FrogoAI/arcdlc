# Modern Go

**Reviewed**: 2026-09-16

**Purpose**: the Go rules a current model applies inconsistently, because the language changed after
most of the Go code it learned from was written. This is an audit target: every rule has an
identifier a gap block can cite, like `GO-7`.

**Not here on purpose**: `gofmt`, naming, early return, `defer`, pointer versus value receivers, the
zero value, and the rest of [Effective Go](https://go.dev/doc/effective_go). Those have not changed
and need no bundled copy. Project structure lives in `Go Server.md`, `Go Client.md` and
`Go Library.md`; layering lives in `mdca.md`.

**Check the `go` line in `go.mod` first.** Each rule names the version it needs. A rule above that
line is not a finding: report it as an upgrade opportunity instead, once, in the register intro.

---

## Types and generics

| ID | Rule | Since |
|----|------|-------|
| `GO-1` | Write `any`, never `interface{}`. They are the same type and one of them is legible. | 1.18 |
| `GO-2` | Reach for a type parameter only when it removes real duplication that an interface cannot. Generic code that has exactly one instantiation is a worse version of the concrete function. | 1.18 |
| `GO-3` | Take constraints from `cmp.Ordered` and the `constraints` you actually need. Do not hand-roll a constraint that the standard library already names. | 1.21 |

## Errors

| ID | Rule | Since |
|----|------|-------|
| `GO-4` | Compare with `errors.Is`, not `==`. Extract with `errors.As`, not a type assertion. A caller that compares directly breaks the first time anyone wraps. | 1.13 |
| `GO-5` | Wrap with `%w` when the caller should be able to inspect the cause, and with `%v` when it should not. Wrapping everything makes every internal error part of your API. | 1.13 |
| `GO-6` | Join independent failures with `errors.Join` rather than concatenating strings or returning only the first. Validation and cleanup are the usual cases. | 1.20 |
| `GO-7` | One `%w` per format string carries one cause. Use `errors.Join` for more than one, not several `%w` verbs. | 1.20 |

## Slices, maps, and the builtins

| ID | Rule | Since |
|----|------|-------|
| `GO-8` | Use `slices.Contains`, `slices.Index`, `slices.Equal`, `slices.Sort` and `slices.SortFunc` instead of writing the loop. A hand-written linear search is now a finding. | 1.21 |
| `GO-9` | Sort comparators return an `int` from `cmp.Compare`, not a `bool`. `slices.SortFunc` takes the comparator, `slices.SortStableFunc` when ties must hold. | 1.21 |
| `GO-10` | Use `min`, `max` and `clear`. They are builtins, so the hand-written helper and the `for k := range m { delete(m, k) }` loop both go. | 1.21 |
| `GO-11` | Prefer `maps.Keys` and `maps.Values` over building the slice by hand. They return iterators, so pair them with `slices.Collect` or `slices.Sorted`. Before 1.23 the standard `maps` package has only `Clone`, `Copy`, `Equal` and `DeleteFunc`, so the hand-written loop is correct there. | 1.23 |

## Loops

| ID | Rule | Since |
|----|------|-------|
| `GO-12` | Loop variables are per iteration. Delete the `x := x` shadow line: it is now noise, and its presence says the code was written for an older Go. | 1.22 |
| `GO-13` | A goroutine or a `defer` inside a loop captures that iteration's variable correctly. The old capture bug is gone, so do not add a parameter just to work around it. | 1.22 |
| `GO-14` | Count with `for i := range 10`, not `for i := 0; i < 10; i++`. | 1.22 |
| `GO-15` | Return an `iter.Seq` or `iter.Seq2` when a caller should walk a sequence it does not own. Prefer it to a callback parameter or to materialising a slice. | 1.23 |

## Context

| ID | Rule | Since |
|----|------|-------|
| `GO-16` | `ctx context.Context` is the first parameter. Never a struct field, never stored, never `nil`. | 1.7 |
| `GO-17` | `context.TODO()` marks unfinished work and must not reach a release. `context.Background()` belongs at the composition root, not in business logic. | 1.7 |
| `GO-18` | Detach with `context.WithoutCancel` when work must outlive the request that started it. Spawning a `Background()` context loses the values too. | 1.21 |
| `GO-19` | Every `context.WithCancel` and `WithTimeout` gets its `cancel` deferred, on every path. | 1.7 |

## Concurrency

| ID | Rule | Since |
|----|------|-------|
| `GO-20` | `sync.OnceFunc` and `sync.OnceValue` replace a `sync.Once` plus a closure plus a package-level result variable. | 1.21 |
| `GO-21` | Start a goroutine only when its lifetime and its stopping condition are both obvious at the call site. An unbounded `go` inside a loop is a finding. | any |

## Logging

| ID | Rule | Since |
|----|------|-------|
| `GO-22` | Structured logging is `log/slog` in the standard library. A new third-party logging dependency needs a written reason. | 1.21 |
| `GO-23` | Log key and value pairs, not formatted sentences. `slog.Info("retry failed", "attempt", n, "err", err)`, never `slog.Info(fmt.Sprintf(...))`. | 1.21 |

## Strings

| ID | Rule | Since |
|----|------|-------|
| `GO-24` | `strings.Cut` replaces the `Index` plus slice pair. `strings.CutPrefix` and `CutSuffix` replace `HasPrefix` plus `TrimPrefix`. | 1.18, 1.20 |
| `GO-25` | Build with `strings.Builder`, not `+=` in a loop. | 1.10 |

## Tests

| ID | Rule | Since |
|----|------|-------|
| `GO-26` | Register teardown with `t.Cleanup`, set environment with `t.Setenv`, and take a scratch directory from `t.TempDir`. Do not hand-roll any of the three. | 1.14, 1.17 |
| `GO-27` | Table-driven tests name each case and run it through `t.Run`, so a failure names itself. | any |
| `GO-28` | Take the test's context from `t.Context()` rather than `context.Background()`, so cancellation follows the test. | 1.24 |
