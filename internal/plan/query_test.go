package plan

import (
	"strings"
	"testing"
)

func TestFirstTODO(t *testing.T) {
	p := Parse(cleanPlan) // AIC-1 DONE, AIC-2 TODO
	got := p.FirstTODO()
	if got == nil || got.ID != "AIC-2" {
		t.Fatalf("FirstTODO = %v, want AIC-2", got)
	}

	// No TODO left → nil.
	allDone := Parse(lines(
		"### AIC-1: Only task",
		"- WHAT: x.",
		"- WHERE: internal/x.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: DONE.",
	))
	if allDone.FirstTODO() != nil {
		t.Fatalf("FirstTODO should be nil when nothing is TODO")
	}
}

func TestByID(t *testing.T) {
	p := Parse(cleanPlan)
	if m := p.ByID("AIC-1"); len(m) != 1 || m[0].Title != "Bootstrap module" {
		t.Fatalf("ByID(AIC-1) = %v", m)
	}
	if m := p.ByID("NOPE"); len(m) != 0 {
		t.Fatalf("ByID(NOPE) should be empty, got %v", m)
	}

	dup := Parse(lines(
		"### AIC-1: First",
		"- Status: TODO.",
		"",
		"### AIC-1: Second",
		"- Status: TODO.",
	))
	if m := dup.ByID("AIC-1"); len(m) != 2 {
		t.Fatalf("ByID on duplicate should return 2, got %d", len(m))
	}
}

func TestRawStartsAtHeading(t *testing.T) {
	p := Parse(cleanPlan)
	t2 := p.FirstTODO()
	raw := string(p.Raw(t2))
	if !strings.HasPrefix(raw, "### AIC-2") {
		t.Fatalf("Raw should start at the AIC-2 heading, got %q", raw)
	}
	if !strings.Contains(raw, "- Status: TODO.") {
		t.Fatalf("Raw should include the status line")
	}
}

func TestStatusCounts(t *testing.T) {
	c := Parse(cleanPlan).StatusCounts()
	if c[StatusDONE] != 1 || c[StatusTODO] != 1 {
		t.Fatalf("counts = %v, want DONE:1 TODO:1", c)
	}
}

func TestParseWhereLayers(t *testing.T) {
	// AIC-2 in cleanPlan has: Layer `handler`: ... and Tests: ...
	p := Parse(cleanPlan)
	layers := ParseWhereLayers(p.Tasks[1].Where)
	if len(layers) != 2 {
		t.Fatalf("want 2 layers, got %d: %+v", len(layers), layers)
	}
	if layers[0].Layer != "handler" || len(layers[0].Targets) != 1 ||
		layers[0].Targets[0] != "internal/h/h.go" {
		t.Errorf("layer[0] = %+v", layers[0])
	}
	if layers[1].Layer != "tests" || layers[1].Targets[0] != "internal/h/h_test.go" {
		t.Errorf("layer[1] = %+v", layers[1])
	}

	// Multi-target line splits on commas.
	ml := ParseWhereLayers("Layer `handler`: internal/h/h.go, router.go.")
	if len(ml) != 1 || len(ml[0].Targets) != 2 || ml[0].Targets[1] != "router.go" {
		t.Errorf("multi-target = %+v", ml)
	}

	// Non-conforming WHERE → empty.
	if got := ParseWhereLayers("somewhere in the code"); len(got) != 0 {
		t.Errorf("free-text WHERE should yield no layers, got %+v", got)
	}
}

func TestFirstTODOEmpty(t *testing.T) {
	if Parse(nil).FirstTODO() != nil {
		t.Fatalf("empty plan FirstTODO should be nil")
	}
}

func TestLineEnd(t *testing.T) {
	p := Parse(cleanPlan)
	for _, tk := range p.Tasks {
		if tk.LineEnd < tk.Line {
			t.Errorf("task %q LineEnd %d < Line %d", tk.ID, tk.LineEnd, tk.Line)
		}
	}
}
