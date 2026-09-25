package plan

import (
	"bytes"
	"errors"
	"testing"
)

// validFinding renders one well-formed task block that passes the strict task
// checks, with the given ID and raw status value (no trailing period, no
// "- Status: " prefix), for example "TODO", "DONE", or
// "BLOCKED — found during A-1: ready to release".
func validFinding(id, status string) []byte {
	return lines(
		"### "+id+": Fix x",
		"- WHAT: Fix x.",
		"- WHERE: `internal/plan/add.go:10`.",
		"- WHY: it breaks y.",
		"- Acceptance: GIVEN x WHEN y THEN z; checked by `TestFoo`.",
		"- References: `internal/plan/add.go`.",
		"- Status: "+status+".",
	)
}

// findingNoAcceptance is like validFinding but omits the Acceptance key, so it
// fails the RequireAcceptance strict check.
func findingNoAcceptance(id, status string) []byte {
	return lines(
		"### "+id+": Fix x",
		"- WHAT: Fix x.",
		"- WHERE: `internal/plan/add.go:10`.",
		"- WHY: it breaks y.",
		"- References: `internal/plan/add.go`.",
		"- Status: "+status+".",
	)
}

func addErr(t *testing.T, err error) *AddError {
	t.Helper()
	var ae *AddError
	if !errors.As(err, &ae) {
		t.Fatalf("error = %v (%T), want *AddError", err, err)
	}
	return ae
}

func TestAppendFinding(t *testing.T) {
	p := Parse(planOf("A-1"))
	block := validFinding("A-1-F1", "BLOCKED — found during A-1: ready to release")

	out, id, err := p.Append(block, nil)
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "A-1-F1" {
		t.Errorf("id = %q, want %q", id, "A-1-F1")
	}
	if !bytes.HasPrefix(out, p.Bytes) {
		t.Errorf("out does not start with the original bytes")
	}
	np := Parse(out)
	if len(np.Tasks) != 2 {
		t.Fatalf("len(Tasks) = %d, want 2\n---\n%s", len(np.Tasks), out)
	}
	if np.Tasks[1].ID != "A-1-F1" {
		t.Errorf("second task = %q, want %q", np.Tasks[1].ID, "A-1-F1")
	}
}

func TestAppendRefuses(t *testing.T) {
	p := Parse(planOf("A-1"))
	cases := []struct {
		name  string
		block []byte
	}{
		{"duplicate ID", validFinding("A-1", "TODO")},
		{"DONE status", validFinding("A-1-F1", "DONE")},
		{"two blocks", append(append([]byte{}, validFinding("A-1-F1", "TODO")...), validFinding("A-1-F2", "TODO")...)},
		{"missing Acceptance", findingNoAcceptance("A-1-F1", "TODO")},
		{"bad reason wording", validFinding("A-1-F1", "BLOCKED — found in A-1: x")},
		{"source not in plan or archive", validFinding("Z-9-F1", "BLOCKED — found during Z-9: x")},
		{"finding ID does not match source", validFinding("B-2", "BLOCKED — found during A-1: x")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, id, err := p.Append(c.block, nil)
			if err == nil {
				t.Fatalf("Append: got nil error, out = %q, id = %q, want *AddError", out, id)
			}
			if ae := addErr(t, err); ae.Kind != AddContract {
				t.Errorf("Kind = %v, want AddContract", ae.Kind)
			}
			if out != nil {
				t.Errorf("out = %q, want nil", out)
			}
		})
	}
}

func TestAppendSourceInArchive(t *testing.T) {
	p := Parse(planOf("A-1")) // A-0 lives only in the archive, not in this plan
	block := validFinding("A-0-F1", "BLOCKED — found during A-0: ready to release")

	out, id, err := p.Append(block, []string{"A-0"})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "A-0-F1" {
		t.Errorf("id = %q, want %q", id, "A-0-F1")
	}
	np := Parse(out)
	if len(np.Tasks) != 2 {
		t.Fatalf("len(Tasks) = %d, want 2", len(np.Tasks))
	}
}

func TestAppendCRLF(t *testing.T) {
	p := Parse(crlf(planOf("A-1")))
	if !p.CRLF {
		t.Fatal("p.CRLF = false, want true")
	}
	block := validFinding("A-1-F1", "BLOCKED — found during A-1: ready to release") // LF-only input

	out, _, err := p.Append(block, nil)
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	for i, b := range out {
		if b == '\n' && (i == 0 || out[i-1] != '\r') {
			t.Fatalf("bare LF at byte %d in out:\n%q", i, out)
		}
	}
}

func TestAppendNoTrailingNewline(t *testing.T) {
	orig := bytes.TrimRight(planOf("A-1"), "\n")
	p := Parse(orig)
	block := validFinding("A-1-F1", "TODO")

	out, _, err := p.Append(block, nil)
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if !bytes.HasPrefix(out, p.Bytes) {
		t.Fatalf("out does not start with the original bytes")
	}
	rest := out[len(p.Bytes):]
	if !bytes.HasPrefix(rest, []byte("\n\n### A-1-F1")) {
		t.Fatalf("rest = %q, want to start with %q (one blank line then the heading)", rest, "\n\n### A-1-F1")
	}
}
