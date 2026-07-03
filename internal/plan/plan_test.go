package plan

import (
	"strings"
	"testing"
)

// lines joins body lines into a plan file (trailing newline included).
// Backticks are literal inside double-quoted Go strings, so paths read naturally.
func lines(ss ...string) []byte {
	return []byte(strings.Join(ss, "\n") + "\n")
}

var cleanPlan = lines(
	"Plan format: see plan-format.md",
	"",
	"### AIC-1: Bootstrap module",
	"",
	"- WHAT: Create the module.",
	"- WHERE:",
	"  Layer `domain`: internal/x/x.go",
	"- WHY: Needed to start.",
	"- References: `docs/aics/checkout/aic.md`.",
	"- Status: DONE.",
	"",
	"### AIC-2 (MISSING): Add endpoint",
	"",
	"- WHAT: Add the read endpoint.",
	"- WHERE:",
	"  Layer `handler`: internal/h/h.go",
	"  Tests: internal/h/h_test.go",
	"- WHY: Parity is blocked without it.",
	"- References: `docs/aics/checkout/aic.md`, `docs/aics/checkout/gap.md`.",
	"- Status: TODO.",
)

func find(fs []Finding, rule string) *Finding {
	for i := range fs {
		if fs[i].Rule == rule {
			return &fs[i]
		}
	}
	return nil
}

func TestParseCleanPlan(t *testing.T) {
	p := Parse(cleanPlan)
	if len(p.Tasks) != 2 {
		t.Fatalf("want 2 tasks, got %d", len(p.Tasks))
	}

	a := p.Tasks[0]
	if a.ID != "AIC-1" || a.Title != "Bootstrap module" {
		t.Errorf("task0 heading: id=%q title=%q", a.ID, a.Title)
	}
	if a.Status != StatusDONE {
		t.Errorf("task0 status: got %v", a.Status)
	}
	if !a.HasWhat || !a.HasWhere || !a.HasWhy || !a.HasRefs || !a.HasStatus {
		t.Errorf("task0 missing a key flag: %+v", a)
	}
	if len(a.References) != 1 || a.References[0] != "docs/aics/checkout/aic.md" {
		t.Errorf("task0 refs: %v", a.References)
	}

	b := p.Tasks[1]
	if b.SourceStatus != SrcMISSING {
		t.Errorf("task1 source status: %q", b.SourceStatus)
	}
	if b.Status != StatusTODO {
		t.Errorf("task1 status: %v", b.Status)
	}
	if len(b.References) != 2 {
		t.Errorf("task1 refs: %v", b.References)
	}
	if !strings.Contains(b.Where, "Tests: internal/h/h_test.go") {
		t.Errorf("task1 WHERE did not absorb multi-line body: %q", b.Where)
	}
}

func TestValidateCleanNoFindings(t *testing.T) {
	p := Parse(cleanPlan)
	if f := p.Validate(ValidateOpts{}); len(f) != 0 {
		t.Fatalf("clean plan (default) should have no findings, got %+v", f)
	}
	if f := p.Validate(ValidateOpts{Strict: true}); len(f) != 0 {
		t.Fatalf("clean plan (strict) should have no findings, got %+v", f)
	}
}

func TestOffsetsTileFile(t *testing.T) {
	p := Parse(cleanPlan)
	// preamble.End must meet task0.RawStart; each task must abut the next; last must reach EOF.
	if len(p.Regions) != 1 || p.Regions[0].Kind != "preamble" {
		t.Fatalf("want one preamble region, got %+v", p.Regions)
	}
	if p.Regions[0].End != p.Tasks[0].RawStart {
		t.Errorf("preamble end %d != task0 start %d", p.Regions[0].End, p.Tasks[0].RawStart)
	}
	if p.Tasks[0].RawEnd != p.Tasks[1].RawStart {
		t.Errorf("task0 end %d != task1 start %d", p.Tasks[0].RawEnd, p.Tasks[1].RawStart)
	}
	if p.Tasks[1].RawEnd != len(cleanPlan) {
		t.Errorf("task1 end %d != EOF %d", p.Tasks[1].RawEnd, len(cleanPlan))
	}
	// Each block slice starts with its heading; each status slice is the status line.
	for _, tk := range p.Tasks {
		if !strings.HasPrefix(string(p.Bytes[tk.RawStart:tk.RawEnd]), "### ") {
			t.Errorf("task %q raw range does not start at heading", tk.ID)
		}
		got := string(p.Bytes[tk.StatusLineStart:tk.StatusLineEnd])
		if !strings.Contains(got, "Status:") {
			t.Errorf("task %q status slice = %q", tk.ID, got)
		}
	}
}

func TestValidateRules(t *testing.T) {
	cases := []struct {
		name   string
		strict bool
		rule   string
		sev    Severity
		plan   []byte
	}{
		{
			name: "missing-status", rule: "missing-status", sev: SevError,
			plan: lines(
				"### AIC-1: No status here",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References: `a`.",
			),
		},
		{
			name: "invalid-status", rule: "invalid-status", sev: SevError,
			plan: lines(
				"### AIC-1: Bad status",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References: `a`.",
				"- Status: WIP.",
			),
		},
		{
			name: "duplicate-id", rule: "duplicate-id", sev: SevError,
			plan: lines(
				"### AIC-1: First",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References: `a`.",
				"- Status: TODO.",
				"",
				"### AIC-1: Second",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References: `a`.",
				"- Status: TODO.",
			),
		},
		{
			name: "missing-key", rule: "missing-key", sev: SevError,
			plan: lines(
				"### AIC-1: No why",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- References: `a`.",
				"- Status: TODO.",
			),
		},
		{
			name: "malformed-heading", rule: "malformed-heading", sev: SevError,
			plan: lines(
				"### no colon in this heading",
				"- Status: TODO.",
			),
		},
		{
			name: "status-casing", rule: "status-casing", sev: SevWarning,
			plan: lines(
				"### AIC-1: Lower status",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References: `a`.",
				"- Status: Todo.",
			),
		},
		{
			name: "unknown-source-status", rule: "unknown-source-status", sev: SevWarning,
			plan: lines(
				"### AIC-1 (WEIRD): Bad source tag",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References: `a`.",
				"- Status: TODO.",
			),
		},
		{
			name: "empty-references", strict: true, rule: "empty-references", sev: SevError,
			plan: lines(
				"### AIC-1: Empty refs",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References:",
				"- Status: TODO.",
			),
		},
		{
			name: "where-no-path", strict: true, rule: "where-no-path", sev: SevError,
			plan: lines(
				"### AIC-1: Vague where",
				"- WHAT: x.",
				"- WHERE: somewhere in the code",
				"- WHY: y.",
				"- References: `a`.",
				"- Status: TODO.",
			),
		},
		{
			name: "dep-order", strict: true, rule: "dep-order", sev: SevError,
			plan: lines(
				"### AIC-1: Refers forward",
				"- WHAT: x.",
				"- WHERE: internal/x.go",
				"- WHY: y.",
				"- References: AIC-2.",
				"- Status: TODO.",
				"",
				"### AIC-2: Later task",
				"- WHAT: x.",
				"- WHERE: internal/y.go",
				"- WHY: y.",
				"- References: `a`.",
				"- Status: TODO.",
			),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := Parse(c.plan)
			f := p.Validate(ValidateOpts{Strict: c.strict})
			got := find(f, c.rule)
			if got == nil {
				t.Fatalf("expected rule %q, findings: %+v", c.rule, f)
			}
			if got.Severity != c.sev {
				t.Errorf("rule %q severity = %q, want %q", c.rule, got.Severity, c.sev)
			}
		})
	}
}

func TestRequireAcceptance(t *testing.T) {
	p := Parse(cleanPlan) // clean plan has no Acceptance keys
	if f := p.Validate(ValidateOpts{RequireAcceptance: true}); find(f, "missing-acceptance") == nil {
		t.Fatalf("expected missing-acceptance when required, got %+v", f)
	}
	if f := p.Validate(ValidateOpts{}); find(f, "missing-acceptance") != nil {
		t.Fatalf("missing-acceptance should be off by default")
	}
}

// acceptancePlan is a strict-clean plan where every task carries a multi-line
// Acceptance section (Given/When/Then).
var acceptancePlan = lines(
	"Plan format: see plan-format.md",
	"",
	"### AIC-1: Add read endpoint",
	"",
	"- WHAT: Add the read endpoint.",
	"- WHERE:",
	"  Layer `handler`: internal/h/h.go",
	"- WHY: Parity is blocked without it.",
	"- Acceptance:",
	"  - GIVEN an existing id WHEN GET /v2/items/{id} THEN 200 + legacy body.",
	"  - GIVEN a missing id WHEN GET THEN 404.",
	"- References: `docs/aics/checkout/aic.md`.",
	"- Status: TODO.",
)

func TestAcceptanceParsedAndStrictClean(t *testing.T) {
	p := Parse(acceptancePlan)
	if len(p.Tasks) != 1 {
		t.Fatalf("want 1 task, got %d", len(p.Tasks))
	}
	tk := p.Tasks[0]
	if !tk.HasAcceptance {
		t.Fatalf("acceptance not detected: %+v", tk)
	}
	if !strings.Contains(tk.Acceptance, "GIVEN an existing id") ||
		!strings.Contains(tk.Acceptance, "GIVEN a missing id WHEN GET THEN 404") {
		t.Errorf("acceptance did not absorb multi-line body: %q", tk.Acceptance)
	}
	// A plan with acceptance on every task must be clean even when acceptance is required.
	if f := p.Validate(ValidateOpts{Strict: true, RequireAcceptance: true}); len(f) != 0 {
		t.Fatalf("acceptance plan (strict+require) should have no findings, got %+v", f)
	}
}

func TestEmptyAcceptanceStrict(t *testing.T) {
	// Acceptance key present but with no body: strict must flag it.
	p := Parse(lines(
		"### AIC-1: Add read endpoint",
		"- WHAT: x.",
		"- WHERE: internal/x.go",
		"- WHY: y.",
		"- Acceptance:",
		"- References: `a`.",
		"- Status: TODO.",
	))
	if f := p.Validate(ValidateOpts{Strict: true}); find(f, "empty-acceptance") == nil {
		t.Fatalf("expected empty-acceptance under strict, got %+v", f)
	}
}

func TestFencedCodeIgnored(t *testing.T) {
	// A ### inside a fenced block must not be parsed as a task.
	p := Parse(lines(
		"Here is an example in a fence:",
		"```md",
		"### EXAMPLE-1: not a real task",
		"- Status: TODO.",
		"```",
		"### AIC-1: Real task",
		"- WHAT: x.",
		"- WHERE: internal/x.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: TODO.",
	))
	if len(p.Tasks) != 1 || p.Tasks[0].ID != "AIC-1" {
		t.Fatalf("fenced ### leaked into tasks: %d tasks", len(p.Tasks))
	}
}

func TestEmptyPlan(t *testing.T) {
	p := Parse(nil)
	if len(p.Tasks) != 0 {
		t.Fatalf("empty plan should have no tasks")
	}
	if f := p.Validate(ValidateOpts{Strict: true}); len(f) != 0 {
		t.Fatalf("empty plan should have no findings, got %+v", f)
	}
}

func TestBlockedWithReason(t *testing.T) {
	p := Parse(lines(
		"### AIC-1: Blocked one",
		"- WHAT: x.",
		"- WHERE: internal/x.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: BLOCKED — waiting on cache dep.",
	))
	if p.Tasks[0].Status != StatusBLOCKED {
		t.Fatalf("want BLOCKED, got %v (raw %q)", p.Tasks[0].Status, p.Tasks[0].StatusRaw)
	}
	if f := p.Validate(ValidateOpts{}); len(f) != 0 {
		t.Fatalf("blocked-with-reason should be clean, got %+v", f)
	}
}
