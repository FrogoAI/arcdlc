package plan

import (
	"strings"
	"testing"
)

func TestWithStatusBytePreserving(t *testing.T) {
	p := Parse(cleanPlan) // AIC-1 DONE, AIC-2 TODO
	aic2 := &p.Tasks[1]

	got := string(p.WithStatus(aic2, StatusTAKEN, ""))
	// Only AIC-2's status line changes; AIC-2 is the first (and only) TODO.
	want := strings.Replace(string(cleanPlan), "- Status: TODO.", "- Status: TAKEN.", 1)
	if got != want {
		t.Fatalf("WithStatus not byte-preserving.\n got: %q\nwant: %q", got, want)
	}
	// AIC-1's DONE must be untouched.
	if strings.Count(got, "- Status: DONE.") != 1 {
		t.Errorf("AIC-1 DONE status was disturbed")
	}
}

func TestWithStatusPreservesNoPeriod(t *testing.T) {
	src := lines(
		"### AIC-1: No-period status",
		"- WHAT: x.",
		"- WHERE: internal/x.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: TODO", // no trailing period
	)
	p := Parse(src)
	got := string(p.WithStatus(&p.Tasks[0], StatusDONE, ""))
	if !strings.Contains(got, "- Status: DONE\n") {
		t.Fatalf("no-period not preserved: %q", got)
	}
	if strings.Contains(got, "DONE.") {
		t.Errorf("should not have added a trailing period")
	}
}

func TestWithStatusBlockReasonRoundTrip(t *testing.T) {
	p := Parse(cleanPlan)
	aic2 := &p.Tasks[1]

	blocked := p.WithStatus(aic2, StatusBLOCKED, "waiting on cache dep")
	if !strings.Contains(string(blocked), "- Status: BLOCKED — waiting on cache dep.") {
		t.Fatalf("block reason not rendered: %q", string(blocked))
	}
	// Re-parse: status is BLOCKED and reason is recoverable.
	rp := Parse(blocked)
	if rp.Tasks[1].Status != StatusBLOCKED {
		t.Fatalf("re-parsed status = %v, want BLOCKED", rp.Tasks[1].Status)
	}
	// todo release drops the reason tail.
	back := p2WithStatus(rp, 1, StatusTODO, "")
	if strings.Contains(back, "BLOCKED") || !strings.Contains(back, "- Status: TODO.") {
		t.Fatalf("todo release did not drop reason: %q", back)
	}
}

// p2WithStatus is a tiny helper to rewrite by task index on a re-parsed plan.
func p2WithStatus(p *Plan, idx int, s Status, reason string) string {
	return string(p.WithStatus(&p.Tasks[idx], s, reason))
}

func TestBlockKeepsExistingReasonWhenNoNewOne(t *testing.T) {
	src := lines(
		"### AIC-1: Already blocked",
		"- WHAT: x.",
		"- WHERE: internal/x.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: BLOCKED — original reason.",
	)
	p := Parse(src)
	// Re-block with empty reason preserves the existing tail.
	got := string(p.WithStatus(&p.Tasks[0], StatusBLOCKED, ""))
	if !strings.Contains(got, "BLOCKED — original reason.") {
		t.Fatalf("existing reason not preserved: %q", got)
	}
}

func TestTransition(t *testing.T) {
	type tc struct {
		cmd     string
		cur     Status
		force   bool
		reason  bool
		wantAct Action
		wantTgt Status
	}
	cases := []tc{
		{"take", StatusTODO, false, false, ActWrite, StatusTAKEN},
		{"take", StatusTAKEN, false, false, ActNoop, StatusTAKEN},
		{"take", StatusDONE, false, false, ActRefuse, StatusTAKEN},
		{"take", StatusDONE, true, false, ActWrite, StatusTAKEN},
		{"done", StatusTAKEN, false, false, ActWrite, StatusDONE},
		{"done", StatusTODO, false, false, ActRefuse, StatusDONE},
		{"done", StatusTODO, true, false, ActWrite, StatusDONE},
		{"done", StatusDONE, false, false, ActNoop, StatusDONE},
		{"todo", StatusTAKEN, false, false, ActWrite, StatusTODO},
		{"todo", StatusBLOCKED, false, false, ActWrite, StatusTODO},
		{"todo", StatusTODO, false, false, ActNoop, StatusTODO},
		{"block", StatusTODO, false, false, ActWrite, StatusBLOCKED},
		{"block", StatusBLOCKED, false, false, ActNoop, StatusBLOCKED},
		{"block", StatusBLOCKED, false, true, ActWrite, StatusBLOCKED},
		{"block", StatusDONE, false, false, ActRefuse, StatusBLOCKED},
		{"block", StatusDONE, true, false, ActWrite, StatusBLOCKED},
	}
	for _, c := range cases {
		tgt, act, _ := Transition(c.cmd, c.cur, c.force, c.reason)
		if act != c.wantAct || tgt != c.wantTgt {
			t.Errorf("%s from %v (force=%v reason=%v): got act=%v tgt=%v, want act=%v tgt=%v",
				c.cmd, c.cur, c.force, c.reason, act, tgt, c.wantAct, c.wantTgt)
		}
	}
}
