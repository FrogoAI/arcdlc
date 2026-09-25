package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mkFile writes a small text file at dir/rel, making parent directories as
// needed.
func mkFile(t *testing.T, dir, rel, body string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// task returns a minimal, parseable task block with the given ID and status.
func task(id, status string) string {
	return "### " + id + ": task\n\n- Status: " + status + ".\n"
}

// statusFixture builds the a/b/c/d tree the acceptance criteria describe:
//   - a/ only aic.md (designing)
//   - b/ plan.md with one TODO and one DONE task (in progress)
//   - c/ plan.md with two DONE tasks and plan-archive.md with three DONE tasks
//     (ready to close, DONE 5)
//   - d/ aic.md and a CLOSED.md (closed)
func statusFixture(t *testing.T) string {
	t.Helper()
	aics := t.TempDir()

	mkFile(t, aics, "a/aic.md", "# A\n")

	mkFile(t, aics, "b/plan.md", task("B-1", "TODO")+"\n"+task("B-2", "DONE"))

	mkFile(t, aics, "c/plan.md", task("C-1", "DONE")+"\n"+task("C-2", "DONE"))
	mkFile(t, aics, "c/plan-archive.md", task("C-3", "DONE")+"\n"+task("C-4", "DONE")+"\n"+task("C-5", "DONE"))

	mkFile(t, aics, "d/aic.md", "# D\n")
	mkFile(t, aics, "d/CLOSED.md", "# D\n\n- Closed: 2026-09-25\n- Outcome: stopped\n")

	return aics
}

func TestRunStatusPhases(t *testing.T) {
	aics := statusFixture(t)

	var out bytes.Buffer
	if code := runStatus(aics, true, &out, &out); code != 0 {
		t.Fatalf("runStatus exit=%d, want 0; output:\n%s", code, out.String())
	}

	var got struct {
		Initiatives []statusRow `json:"initiatives"`
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	rows := got.Initiatives
	if len(rows) != 4 {
		t.Fatalf("got %d rows, want 4: %+v", len(rows), rows)
	}

	wantSlugs := []string{"a", "b", "c", "d"}
	for i, s := range wantSlugs {
		if rows[i].Slug != s {
			t.Fatalf("row %d slug = %q, want %q (order: %+v)", i, rows[i].Slug, s, rows)
		}
	}

	a, b, c, d := rows[0], rows[1], rows[2], rows[3]

	if a.Phase != "designing" {
		t.Errorf("a phase = %q, want designing", a.Phase)
	}

	if b.Phase != "in progress" {
		t.Errorf("b phase = %q, want in progress", b.Phase)
	}
	if b.Counts.TODO != 1 || b.Counts.DONE != 1 {
		t.Errorf("b counts = %+v, want TODO=1 DONE=1", b.Counts)
	}

	if c.Phase != "ready to close" {
		t.Errorf("c phase = %q, want ready to close", c.Phase)
	}
	if c.Counts.DONE != 5 {
		t.Errorf("c DONE = %d, want 5", c.Counts.DONE)
	}

	if d.Phase != "closed" {
		t.Errorf("d phase = %q, want closed", d.Phase)
	}
	if d.Closed != "2026-09-25" {
		t.Errorf("d closed = %q, want 2026-09-25", d.Closed)
	}
	if d.Outcome != "stopped" {
		t.Errorf("d outcome = %q, want stopped", d.Outcome)
	}
}

func TestRunStatusMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	var out bytes.Buffer
	if code := runStatus(missing, true, &out, &out); code != 0 {
		t.Fatalf("runStatus exit=%d, want 0; output:\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), `"initiatives": []`) {
		t.Errorf("output = %q, want it to contain \"initiatives\": []", out.String())
	}
}

func TestRunStatusText(t *testing.T) {
	aics := statusFixture(t)

	var out bytes.Buffer
	if code := runStatus(aics, false, &out, &out); code != 0 {
		t.Fatalf("runStatus exit=%d, want 0; output:\n%s", code, out.String())
	}

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "SLUG") {
		t.Fatalf("first line = %q, want it to start with SLUG", lines[0])
	}

	var dLine string
	for _, ln := range lines[1:] {
		if fields := strings.Fields(ln); len(fields) > 0 && fields[0] == "d" {
			dLine = ln
		}
	}
	if !strings.Contains(dLine, "2026-09-25") || !strings.Contains(dLine, "stopped") {
		t.Errorf("d line = %q, want it to contain 2026-09-25 and stopped", dLine)
	}
}
