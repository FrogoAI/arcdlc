package plan

import (
	"strings"
	"testing"
)

var archivePlan = lines(
	"Task format: see plan-format.md.",
	"",
	"### AIC-1: Domain model",
	"- WHAT: x.",
	"- WHERE: internal/a.go",
	"- WHY: y.",
	"- References: `a`.",
	"- Status: DONE.",
	"",
	"### AIC-2: Repository",
	"- WHAT: x.",
	"- WHERE: internal/b.go",
	"- WHY: y.",
	"- References: `a`.",
	"- Status: TODO.",
	"",
	"### AIC-3: Handler",
	"- WHAT: x.",
	"- WHERE: internal/c.go",
	"- WHY: y.",
	"- References: `a`.",
	"- Status: DONE.",
)

func TestArchiveBasic(t *testing.T) {
	p := Parse(archivePlan)
	res, ok := p.Archive("docs/aics/plan-archive.md", "2026-07-03")
	if !ok {
		t.Fatal("expected something to archive")
	}
	if len(res.Archived) != 2 || res.Archived[0].ID != "AIC-1" || res.Archived[1].ID != "AIC-3" {
		t.Fatalf("archived = %v", ids(res.Archived))
	}

	np := Parse(res.NewPlan)
	if len(np.Tasks) != 1 || np.Tasks[0].ID != "AIC-2" || np.Tasks[0].Status != StatusTODO {
		t.Fatalf("new plan tasks = %v", ids2(np.Tasks))
	}
	body := string(res.NewPlan)
	if strings.Contains(body, "Status: DONE") {
		t.Errorf("new plan still contains a DONE status:\n%s", body)
	}
	for _, want := range []string{
		"Completed (archived to docs/aics/plan-archive.md):",
		"- AIC-1: Domain model",
		"- AIC-3: Handler",
		"### AIC-2: Repository\n- WHAT: x.", // kept block preserved contiguously
	} {
		if !strings.Contains(body, want) {
			t.Errorf("new plan missing %q:\n%s", want, body)
		}
	}

	sec := Parse(res.Section)
	if len(sec.Tasks) != 2 || sec.ByID("AIC-1") == nil || sec.ByID("AIC-3") == nil {
		t.Errorf("section tasks = %v", ids2(sec.Tasks))
	}
	if !strings.HasPrefix(string(res.Section), "## Archived 2026-07-03") {
		t.Errorf("section header wrong: %q", string(res.Section)[:20])
	}

	if err := VerifyArchive(p, res); err != nil {
		t.Fatalf("VerifyArchive: %v", err)
	}
}

func TestArchiveNothing(t *testing.T) {
	p := Parse(lines(
		"### AIC-1: Only pending",
		"- WHAT: x.",
		"- WHERE: internal/a.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: TODO.",
	))
	if _, ok := p.Archive("x.md", "2026-07-03"); ok {
		t.Fatal("nothing should be archived when no DONE blocks exist")
	}
}

func TestArchiveExtendsExistingLedger(t *testing.T) {
	src := lines(
		"Task format: see plan-format.md.",
		"",
		"Completed (archived to docs/aics/plan-archive.md):",
		"- AIC-0: Bootstrap",
		"",
		"### AIC-1: Domain model",
		"- WHAT: x.",
		"- WHERE: internal/a.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: DONE.",
		"",
		"### AIC-2: Repo",
		"- WHAT: x.",
		"- WHERE: internal/b.go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: TODO.",
	)
	p := Parse(src)
	res, ok := p.Archive("docs/aics/plan-archive.md", "2026-07-03")
	if !ok {
		t.Fatal("expected archive")
	}
	body := string(res.NewPlan)
	if n := strings.Count(body, "Completed (archived"); n != 1 {
		t.Fatalf("want exactly one ledger header, got %d:\n%s", n, body)
	}
	for _, want := range []string{"- AIC-0: Bootstrap", "- AIC-1: Domain model"} {
		if !strings.Contains(body, want) {
			t.Errorf("ledger missing %q:\n%s", want, body)
		}
	}
	// Existing ledger entries are not re-parsed as tasks.
	if np := Parse(res.NewPlan); len(np.Tasks) != 1 || np.Tasks[0].ID != "AIC-2" {
		t.Errorf("new plan tasks = %v", ids2(np.Tasks))
	}
	if err := VerifyArchive(p, res); err != nil {
		t.Fatalf("VerifyArchive: %v", err)
	}
}

func ids(ts []*Task) []string {
	out := make([]string, len(ts))
	for i, t := range ts {
		out[i] = t.ID
	}
	return out
}

func ids2(ts []Task) []string {
	out := make([]string, len(ts))
	for i := range ts {
		out[i] = ts[i].ID
	}
	return out
}
