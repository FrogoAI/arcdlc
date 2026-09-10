package scan

import (
	"strings"
	"testing"

	"github.com/FrogoAI/arcdlc/internal/plan"
)

const day = "2026-09-10"

func find(file string, line int, marker, text, code string) Finding {
	return Finding{File: file, Line: line, Marker: marker, Text: text, Code: code}
}

func render(t *testing.T, existing []byte, found []Finding, markers ...string) ([]byte, Result) {
	t.Helper()
	if len(markers) == 0 {
		markers = []string{"TODO"}
	}
	out, res, err := Render(existing, found, markers, day)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if err := Verify(existing, out, res); err != nil {
		t.Fatalf("Verify: %v\n---\n%s", err, out)
	}
	return out, res
}

func TestRenderFirstRunWritesHeaderAndBlock(t *testing.T) {
	out, res := render(t, nil, []Finding{
		find("internal/foo/bar.go", 42, "TODO", "TODO move this into its own package", "func rebuild() error {"),
	})
	if !res.Changed || len(res.New) != 1 || res.NewIDs[0] != "TODO-CMT-01" {
		t.Fatalf("result = %+v", res)
	}
	s := string(out)
	for _, want := range []string{
		"# Comment register",
		"### TODO-CMT-01: Move this into its own package",
		"- Marker: `internal/foo/bar.go:42` `TODO move this into its own package`",
		"- Verdict: NEW.",
		"- WHAT:",
		"- HOW:",
		"- WHERE: internal/foo/bar.go (marker at line 42)",
		"  func rebuild() error {",
		"- WHY:",
		"- Acceptance:",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("register is missing %q:\n%s", want, s)
		}
	}
	// The block must parse as a task block, because the skill mirrors it into the plan.
	tasks := plan.Parse(out).Tasks
	if len(tasks) != 1 || tasks[0].ID != "TODO-CMT-01" || !tasks[0].HeadingOK {
		t.Fatalf("parsed tasks = %+v", tasks)
	}
	if !tasks[0].HasWhat || !tasks[0].HasWhere || !tasks[0].HasWhy || !tasks[0].HasAcceptance {
		t.Fatalf("block is missing keys: %+v", tasks[0])
	}
}

func TestRenderIsIdempotent(t *testing.T) {
	found := []Finding{
		find("a.go", 3, "TODO", "TODO one", "func a() {}"),
		find("b.go", 9, "TODO", "TODO two", "func b() {}"),
	}
	first, _ := render(t, nil, found)
	second, res := render(t, first, found)
	if string(second) != string(first) {
		t.Fatalf("second run changed the file:\n%s", second)
	}
	if res.Changed || len(res.New) != 0 || res.Known != 2 || res.Kept != 2 {
		t.Fatalf("result = %+v, want no change", res)
	}
	// A moved marker keeps its block: identity is the text, not the line.
	moved := []Finding{
		find("a.go", 300, "TODO", "TODO one", "func a() {}"),
		find("b.go", 900, "TODO", "TODO two", "func b() {}"),
	}
	third, res := render(t, first, moved)
	if string(third) != string(first) || res.Changed {
		t.Fatalf("a moved marker must not change the register:\n%s", third)
	}
}

func TestRenderAppendsWithoutTouchingEarlierBlocks(t *testing.T) {
	first, _ := render(t, nil, []Finding{find("a.go", 3, "TODO", "TODO one", "")})
	// The skill fills the judgement in.
	judged := strings.Replace(string(first), "- Verdict: NEW.", "- Verdict: ACTIONABLE.", 1)
	judged = strings.Replace(judged, "- WHAT:", "- WHAT: Split the store.", 1)
	judged = strings.Replace(judged, "### TODO-CMT-01: One", "### TODO-CMT-01: Split the store package", 1)

	out, res := render(t, []byte(judged), []Finding{
		find("a.go", 3, "TODO", "TODO one", ""),
		find("c.go", 7, "TODO", "TODO three", "func c() {}"),
	})
	s := string(out)
	if !strings.HasPrefix(s, judged[:len(judged)-1]) {
		t.Fatalf("earlier content was rewritten:\n%s", s)
	}
	for _, want := range []string{
		"### TODO-CMT-01: Split the store package",
		"- WHAT: Split the store.",
		"- Verdict: ACTIONABLE.",
		"### TODO-CMT-02: Three",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q:\n%s", want, s)
		}
	}
	if len(res.New) != 1 || res.NewIDs[0] != "TODO-CMT-02" || res.Kept != 1 {
		t.Fatalf("result = %+v", res)
	}
}

func TestRenderLeavesABlockWhoseMarkerIsGone(t *testing.T) {
	first, _ := render(t, nil, []Finding{
		find("a.go", 3, "TODO", "TODO one", ""),
		find("b.go", 9, "TODO", "TODO two", ""),
	})
	judged := strings.Replace(string(first), "- Verdict: NEW.", "- Verdict: ACTIONABLE.", 1)

	// The sweep deletes each comment once it is registered, so the next sweep
	// finds nothing. That must not touch a single block.
	out, res := render(t, []byte(judged), nil)
	if string(out) != judged {
		t.Fatalf("an empty sweep rewrote the register:\n%s", out)
	}
	if res.Changed || len(res.New) != 0 || res.Total != 2 {
		t.Fatalf("result = %+v, want no change", res)
	}
	if strings.Contains(string(out), "RESOLVED") {
		t.Error("a verdict was invented for a marker the sweep removed")
	}
}

func TestRenderCountsAKnownMarkerWithoutAppending(t *testing.T) {
	found := []Finding{find("a.go", 3, "TODO", "TODO one", "")}
	first, _ := render(t, nil, found)
	// Somebody wrote the same comment again: the register already holds it, so no
	// second block appears, and the sweep still takes the comment out of the code.
	out, res := render(t, first, found)
	if string(out) != string(first) {
		t.Fatalf("a known marker was registered twice:\n%s", out)
	}
	if res.Known != 1 || len(res.New) != 0 {
		t.Fatalf("result = %+v, want one known finding", res)
	}
}

func TestRenderNumbersPerMarkerAndContinues(t *testing.T) {
	first, _ := render(t, nil, []Finding{
		find("a.go", 1, "TODO", "TODO one", ""),
		find("b.go", 2, "FIXME", "FIXME one", ""),
	}, "TODO", "FIXME")
	out, res := render(t, first, []Finding{
		find("a.go", 1, "TODO", "TODO one", ""),
		find("b.go", 2, "FIXME", "FIXME one", ""),
		find("c.go", 3, "TODO", "TODO two", ""),
		find("d.go", 4, "FIXME", "FIXME two", ""),
	}, "TODO", "FIXME")
	if strings.Join(res.NewIDs, " ") != "TODO-CMT-02 FIXME-CMT-02" {
		t.Fatalf("new IDs = %v", res.NewIDs)
	}
	if got := strings.Count(string(out), "### "); got != 4 {
		t.Fatalf("%d blocks, want 4", got)
	}
}

func TestRenderKeepsTwoIdenticalMarkersInOneFile(t *testing.T) {
	found := []Finding{
		find("a.go", 3, "TODO", "TODO fix", ""),
		find("a.go", 40, "TODO", "TODO fix", ""),
	}
	out, res := render(t, nil, found)
	if len(res.New) != 2 {
		t.Fatalf("result = %+v, want two findings", res)
	}
	second, res := render(t, out, found)
	if string(second) != string(out) || res.Changed {
		t.Fatalf("identical markers are not stable across runs: %+v", res)
	}
}

func TestRenderStripsBackticksFromTheMarkerLine(t *testing.T) {
	out, _ := render(t, nil, []Finding{
		find("a.go", 3, "TODO", "TODO drop `oldName` and use NewName", ""),
	})
	if !strings.Contains(string(out), "- Marker: `a.go:3` `TODO drop 'oldName' and use NewName`") {
		t.Fatalf("marker line broke on a backtick:\n%s", out)
	}
	entries, err := ParseEntries(out)
	if err != nil || len(entries) != 1 {
		t.Fatalf("ParseEntries: %v %+v", err, entries)
	}
}

func TestRenderRejectsABlockWithoutAMarkerLine(t *testing.T) {
	bad := []byte("# Comment register\n\n### TODO-CMT-01: hand written\n\n- WHAT: x.\n")
	if _, _, err := Render(bad, nil, []string{"TODO"}, day); err == nil {
		t.Fatal("Render accepted a block with no Marker line")
	}
}

func TestVerifyCatchesADroppedBlock(t *testing.T) {
	first, _ := render(t, nil, []Finding{
		find("a.go", 3, "TODO", "TODO one", ""),
		find("b.go", 9, "TODO", "TODO two", ""),
	})
	// Simulate a renderer that lost the second block.
	cut := first[:strings.Index(string(first), "### TODO-CMT-02")]
	err := Verify(first, cut, Result{Total: 2, Kept: 2})
	if err == nil {
		t.Fatal("Verify accepted a register that lost a block")
	}
	if !strings.Contains(err.Error(), "block") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderKeepsCRLF(t *testing.T) {
	first, _ := render(t, nil, []Finding{find("a.go", 3, "TODO", "TODO one", "")})
	crlf := []byte(strings.ReplaceAll(string(first), "\n", "\r\n"))
	out, res := render(t, crlf, []Finding{
		find("a.go", 3, "TODO", "TODO one", ""),
		find("b.go", 9, "TODO", "TODO two", ""),
	})
	if len(res.New) != 1 {
		t.Fatalf("result = %+v", res)
	}
	if strings.Contains(strings.ReplaceAll(string(out), "\r\n", ""), "\n") {
		t.Fatalf("mixed line endings:\n%q", out)
	}
}

func TestTitleDraftsAnInstruction(t *testing.T) {
	cases := []struct{ text, want string }{
		{"TODO We should move this function into a separate package, and use another naming.",
			"Move this function into a separate package"},
		{"TODO: fix the retry loop", "Fix the retry loop"},
		{"TODO", "Handle the TODO marker in a.go"},
	}
	for _, c := range cases {
		got := title(Finding{File: "a.go", Marker: "TODO", Text: c.text})
		if got != c.want {
			t.Errorf("title(%q) = %q, want %q", c.text, got, c.want)
		}
	}
}
