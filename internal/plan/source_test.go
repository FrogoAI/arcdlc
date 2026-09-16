package plan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func findRule(fs []Finding, rule string) *Finding {
	for i := range fs {
		if fs[i].Rule == rule {
			return &fs[i]
		}
	}
	return nil
}

const planBody = "### T-1: Do the thing\n" +
	"- WHAT: Do the thing.\n" +
	"- WHERE: internal/thing/thing.go\n" +
	"- WHY: The thing is not done.\n" +
	"- Acceptance:\n" +
	"  - GIVEN it WHEN `go test ./...` runs THEN it passes.\n" +
	"- References: docs/aics/x/aic.md\n" +
	"- Status: TODO.\n"

// TestSourceStampRoundTrip is the whole point: a plan stamped against a
// document stays clean until that document changes, and then says so.
func TestSourceStampRoundTrip(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/aics/x/aic.md", "# X\n\n> A thing.\n\nOriginal design.\n")

	p := Parse([]byte("# x — Plan\n\n" + planBody))
	stamped, added, err := AddSource(p, root, "docs/aics/x/aic.md")
	if err != nil {
		t.Fatalf("AddSource: %v", err)
	}
	if !added {
		t.Fatal("AddSource reported no change on a plan with no stamp")
	}

	p = Parse(stamped)
	srcs := p.Sources()
	if len(srcs) != 1 {
		t.Fatalf("want 1 source, got %d in:\n%s", len(srcs), stamped)
	}
	if srcs[0].Path != "docs/aics/x/aic.md" {
		t.Errorf("path = %q", srcs[0].Path)
	}
	if len(p.Tasks) != 1 {
		t.Errorf("stamp broke task parsing: want 1 task, got %d", len(p.Tasks))
	}
	if f := findRule(CheckSources(p, root), "source-changed"); f != nil {
		t.Errorf("fresh stamp reported stale: %s", f.Message)
	}

	// The design moves underneath the plan.
	writeFile(t, root, "docs/aics/x/aic.md", "# X\n\n> A thing.\n\nRevised design.\n")
	f := findRule(CheckSources(p, root), "source-changed")
	if f == nil {
		t.Fatal("changed source not reported")
	}
	if f.Severity != SevWarning {
		t.Errorf("want warning, got %q", f.Severity)
	}

	// Re-stamping clears it.
	restamped, refreshed, err := Stamp(p, root)
	if err != nil {
		t.Fatalf("Stamp: %v", err)
	}
	if len(refreshed) != 1 {
		t.Fatalf("want 1 refreshed path, got %v", refreshed)
	}
	if fs := CheckSources(Parse(restamped), root); len(fs) != 0 {
		t.Errorf("still stale after re-stamp: %+v", fs)
	}
}

// TestAddSourceIsIdempotent pins that re-running the planner does not stack
// duplicate stamps.
func TestAddSourceIsIdempotent(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/aics/x/aic.md", "# X\n")
	p := Parse([]byte("# x — Plan\n\n" + planBody))
	b, _, err := AddSource(p, root, "docs/aics/x/aic.md")
	if err != nil {
		t.Fatal(err)
	}
	b2, added, err := AddSource(Parse(b), root, "docs/aics/x/aic.md")
	if err != nil {
		t.Fatal(err)
	}
	if added {
		t.Error("AddSource added a second stamp for the same document")
	}
	if string(b) != string(b2) {
		t.Error("AddSource rewrote a plan that already carried the stamp")
	}
}

// TestMissingSourceIsAnError separates "the design moved" from "the design is
// gone": the first is a warning the engineer judges, the second is broken.
func TestMissingSourceIsAnError(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/aics/x/aic.md", "# X\n")
	p := Parse([]byte("# x — Plan\n\n" + planBody))
	b, _, err := AddSource(p, root, "docs/aics/x/aic.md")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "docs/aics/x/aic.md")); err != nil {
		t.Fatal(err)
	}
	f := findRule(CheckSources(Parse(b), root), "source-missing")
	if f == nil {
		t.Fatal("missing source not reported")
	}
	if f.Severity != SevError {
		t.Errorf("want error, got %q", f.Severity)
	}
}

// TestUnstampedPlanIsNotChecked keeps the feature backward compatible: a plan
// written before stamps exist must not start failing.
func TestUnstampedPlanIsNotChecked(t *testing.T) {
	root := t.TempDir()
	p := Parse([]byte("# x — Plan\n\n" + planBody))
	if len(p.Sources()) != 0 {
		t.Fatal("found a stamp in a plan that has none")
	}
	if fs := CheckSources(p, root); len(fs) != 0 {
		t.Errorf("unstamped plan reported findings: %+v", fs)
	}
}

// TestStampPreservesCRLF pins the line-ending discipline the rest of the
// package keeps.
func TestStampPreservesCRLF(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/aics/x/aic.md", "# X\n")
	src := strings.ReplaceAll("# x — Plan\n\n"+planBody, "\n", "\r\n")
	b, _, err := AddSource(Parse([]byte(src)), root, "docs/aics/x/aic.md")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ReplaceAll(string(b), "\r\n", ""), "\n") {
		t.Error("stamp introduced a bare LF into a CRLF plan")
	}
	if len(Parse(b).Sources()) != 1 {
		t.Error("stamp not found after CRLF insert")
	}
}
