package plan

import "testing"

// TestExecutorKeyPinsOneTask pins the per-task tier: an optional single-line
// Executor key, read exactly as written apart from surrounding whitespace and
// one trailing period, absent on tasks that do not carry it. It sits above the
// CONTEXT.md pin in the dispatcher's resolution order, so the parser has to
// hand it over as a field rather than leave it in prose or in Extra.
func TestExecutorKeyPinsOneTask(t *testing.T) {
	p := Parse(lines(
		"### AIC-1: Migrate the ledger table",
		"- WHAT: Migrate it.",
		"- HOW:",
		"  Batch by 1000 rows; keep the old table until AIC-3.",
		"- Executor: opus, high effort.",
		"- WHERE: internal/ledger/migrate.go",
		"- WHY: y.",
		"- Acceptance:",
		"  - GIVEN the table WHEN `go test ./internal/ledger/...` runs THEN it passes.",
		"- References: `docs/aics/demo/aic.md`.",
		"- Status: TODO.",
		"",
		"### AIC-2: Rename a flag",
		"- WHAT: Rename it.",
		"- Executor: sonnet",
		"- WHERE: cmd/x/main.go",
		"- WHY: y.",
		"- Acceptance:",
		"  - GIVEN the flag WHEN `go build ./...` runs THEN it compiles.",
		"- References: `docs/aics/demo/aic.md`.",
		"- Status: TODO.",
		"",
		"### AIC-3: No pin here",
		"- WHAT: x.",
		"- WHERE: internal/y.go",
		"- WHY: y.",
		"- Acceptance:",
		"  - GIVEN z WHEN run THEN ok.",
		"- References: `docs/aics/demo/aic.md`.",
		"- Status: TODO.",
	))
	if len(p.Tasks) != 3 {
		t.Fatalf("want 3 tasks, got %d", len(p.Tasks))
	}
	a, b, c := p.Tasks[0], p.Tasks[1], p.Tasks[2]

	if !a.HasExecutor || a.Executor != "opus, high effort" {
		t.Errorf("AIC-1 executor = %q (has=%v), want %q with the period trimmed", a.Executor, a.HasExecutor, "opus, high effort")
	}
	// The single-line key must end the HOW body above it, not be absorbed by it.
	if a.How != "Batch by 1000 rows; keep the old table until AIC-3." {
		t.Errorf("HOW absorbed the Executor line: %q", a.How)
	}
	if !b.HasExecutor || b.Executor != "sonnet" {
		t.Errorf("AIC-2 executor = %q (has=%v)", b.Executor, b.HasExecutor)
	}
	if c.HasExecutor || c.Executor != "" {
		t.Errorf("AIC-3 should carry no pin: %+v", c)
	}
	// Executor is a defined key now, so it must never show up as a custom one.
	for _, tk := range p.Tasks {
		if len(tk.Extra) != 0 {
			t.Errorf("task %s: Executor leaked into Extra: %+v", tk.ID, tk.Extra)
		}
	}
	// The key is optional and a present one is fine: strict must stay silent.
	if f := p.Validate(ValidateOpts{Strict: true, RequireAcceptance: true}); len(f) != 0 {
		t.Fatalf("executor plan should have no findings, got %+v", f)
	}
}

// TestEmptyExecutorIsAStrictError pins that a pin with nothing on it is a
// defect under --strict. Silently treating it as "no pin" would let a typo
// hand the task to the wrong tier without a word.
func TestEmptyExecutorIsAStrictError(t *testing.T) {
	src := lines(
		"### AIC-1: x",
		"- WHAT: x.",
		"- Executor:",
		"- WHERE: internal/x.go",
		"- WHY: y.",
		"- Acceptance:",
		"  - GIVEN z WHEN `go test ./...` runs THEN ok.",
		"- References: `docs/aics/demo/aic.md`.",
		"- Status: TODO.",
	)
	p := Parse(src)
	if !p.Tasks[0].HasExecutor || p.Tasks[0].Executor != "" {
		t.Fatalf("empty Executor mis-parsed: %+v", p.Tasks[0])
	}
	if f := p.Validate(ValidateOpts{}); len(f) != 0 {
		t.Errorf("non-strict should not report an empty Executor, got %+v", f)
	}
	f := p.Validate(ValidateOpts{Strict: true, RequireAcceptance: true})
	if len(f) != 1 || f[0].Rule != "empty-executor" || f[0].Severity != SevError {
		t.Fatalf("want one empty-executor error, got %+v", f)
	}
}

// TestStatusMutationPreservesExecutor pins that the byte-preserving status
// rewrite leaves the pin exactly as written.
func TestStatusMutationPreservesExecutor(t *testing.T) {
	src := lines(
		"### AIC-1: x",
		"- WHAT: x.",
		"- Executor: opus, high effort.",
		"- Status: TODO.",
	)
	p := Parse(src)
	out := p.WithStatus(&p.Tasks[0], StatusTAKEN, "")
	got := Parse(out).Tasks[0]
	if got.Executor != "opus, high effort" || got.Status != StatusTAKEN {
		t.Errorf("status flip disturbed the pin: %+v", got)
	}
}
