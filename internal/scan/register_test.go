package scan

import (
	"strings"
	"testing"

	"github.com/FrogoAI/arcdlc/internal/plan"
)

func find(file string, line int, marker, text, code string) Finding {
	return Finding{File: file, Line: line, Marker: marker, Text: text, Code: code}
}

// grouped is a finding the engineer tagged, so the register holds it with every
// other marker carrying the same tag.
func grouped(file string, line int, marker, group, text, code string) Finding {
	f := find(file, line, marker, text, code)
	f.Group = group
	return f
}

func render(t *testing.T, existing []byte, found []Finding) ([]byte, Result) {
	t.Helper()
	out, res, err := Render(existing, found)
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
	if !res.Changed || len(res.New) != 1 || res.New[0].ID != "TODO-CMT-01" {
		t.Fatalf("result = %+v", res)
	}
	s := string(out)
	for _, want := range []string{
		"# Comment register",
		"### TODO-CMT-01: Move this into its own package",
		"- Marker: `internal/foo/bar.go:42` `TODO move this into its own package`",
		"- Verdict:\n",
		"- WHAT: move this into its own package",
		"- HOW:",
		"- WHERE:\n  internal/foo/bar.go (marker at line 42)\n    func rebuild() error {",
		"- WHY:",
		"- Acceptance:",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("register is missing %q:\n%s", want, s)
		}
	}
	// The sweep judges nothing: the verdict is the skill's to write.
	if strings.Contains(s, "- Verdict: NEW") {
		t.Error("the sweep wrote a verdict")
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
	judged := judge(string(first), "Split the store package", "Split the store.")

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
	if len(res.New) != 1 || res.New[0].ID != "TODO-CMT-02" || res.Kept != 1 {
		t.Fatalf("result = %+v", res)
	}
}

func TestRenderLeavesABlockWhoseMarkerIsGone(t *testing.T) {
	first, _ := render(t, nil, []Finding{
		find("a.go", 3, "TODO", "TODO one", ""),
		find("b.go", 9, "TODO", "TODO two", ""),
	})
	judged := judge(string(first), "Split the store package", "Split the store.")

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
	})
	out, res := render(t, first, []Finding{
		find("a.go", 1, "TODO", "TODO one", ""),
		find("b.go", 2, "FIXME", "FIXME one", ""),
		find("c.go", 3, "TODO", "TODO two", ""),
		find("d.go", 4, "FIXME", "FIXME two", ""),
	})
	if got := blockIDs(res.New); got != "TODO-CMT-02 FIXME-CMT-02" {
		t.Fatalf("new IDs = %v", got)
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

func TestRenderGroupsMarkersByTag(t *testing.T) {
	out, res := render(t, nil, []Finding{
		grouped("internal/a.go", 12, "TODO", "G1", "TODO:G1 move the rebuild into internal", "func rebuild() error {"),
		find("internal/b.go", 20, "TODO", "TODO drop the retry", "func retry() {}"),
		grouped("internal/c.go", 40, "TODO", "G1", "TODO:g1 change the return format", "return out, nil"),
	})
	if len(res.New) != 2 {
		t.Fatalf("result = %+v, want two blocks", res)
	}
	if got := blockIDs(res.New); got != "TODO-CMT-G1 TODO-CMT-01" {
		t.Fatalf("new IDs = %v", got)
	}
	if len(res.New[0].Members) != 2 || res.New[0].Group != "G1" {
		t.Fatalf("group block = %+v, want two members", res.New[0])
	}
	s := string(out)
	for _, want := range []string{
		"### TODO-CMT-G1: Move the rebuild into internal",
		"- Marker: `internal/a.go:12` `TODO:G1 move the rebuild into internal`",
		"- Marker: `internal/c.go:40` `TODO:g1 change the return format`",
		"- WHAT: move the rebuild into internal; change the return format",
		"- WHERE:\n  internal/a.go (marker at line 12)\n    func rebuild() error {\n" +
			"  internal/c.go (marker at line 40)\n    return out, nil\n",
		"### TODO-CMT-01: Drop the retry",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("register is missing %q:\n%s", want, s)
		}
	}
	// The tag is the block's identity, not part of what it asks for, so it stays on
	// the "- Marker:" lines and out of the heading and WHAT.
	for _, line := range strings.Split(s, "\n") {
		text := ""
		switch {
		case strings.HasPrefix(line, "### "):
			text = line[strings.Index(line, ": ")+2:] // the ID keeps the tag, the title must not
		case strings.HasPrefix(line, "- WHAT:"):
			text = line
		}
		if strings.Contains(text, "G1") || strings.Contains(text, "g1") {
			t.Errorf("the tag leaked into %q", line)
		}
	}
	if tasks := plan.Parse(out).Tasks; len(tasks) != 2 || !tasks[0].HasWhere {
		t.Fatalf("parsed tasks = %+v", tasks)
	}
}

func TestRenderKeepsOneTagPerMarkerWord(t *testing.T) {
	_, res := render(t, nil, []Finding{
		grouped("a.go", 1, "TODO", "G1", "TODO:G1 one", ""),
		grouped("b.go", 2, "FIXME", "G1", "FIXME:G1 two", ""),
	})
	if got := blockIDs(res.New); got != "TODO-CMT-G1 FIXME-CMT-G1" {
		t.Fatalf("new IDs = %v, want one block per marker word", got)
	}
}

func TestRenderGrowsTheBlockATagAlreadyHas(t *testing.T) {
	first, _ := render(t, nil, []Finding{
		grouped("internal/a.go", 12, "TODO", "G1", "TODO:G1 move the rebuild into internal", ""),
	})
	// The skill judged the block and the plan carries the task.
	judged := judge(string(first), "Move the index rebuild into internal/index",
		"Move the index rebuild into internal/index.")
	judged = strings.Replace(judged, "- HOW:", "- HOW:\n  Keep the exported name.", 1)
	judged = strings.Replace(judged, "- Acceptance:", "- Acceptance:\n  - GIVEN a build WHEN go test ./... THEN it passes.", 1)

	// Weeks later somebody tags one more site with G1.
	out, res := render(t, []byte(judged), []Finding{
		grouped("internal/a.go", 12, "TODO", "G1", "TODO:G1 move the rebuild into internal", ""),
		grouped("internal/port.go", 8, "TODO", "G1", "TODO:G1 also rename the port", "func Port() int {"),
	})
	if len(res.New) != 0 || len(res.Extended) != 1 || res.Known != 1 || res.Kept != 0 {
		t.Fatalf("result = %+v, want one extended block", res)
	}
	if res.Extended[0].ID != "TODO-CMT-G1" || len(res.Extended[0].Members) != 1 {
		t.Fatalf("extended = %+v, want the one new member", res.Extended[0])
	}
	s := string(out)
	if got := strings.Count(s, "### "); got != 1 {
		t.Fatalf("%d blocks, want the tag to grow its own block:\n%s", got, s)
	}
	for _, want := range []string{
		"### TODO-CMT-G1: Move the index rebuild into internal/index",
		"- Marker: `internal/a.go:12` `TODO:G1 move the rebuild into internal`",
		"- Marker: `internal/port.go:8` `TODO:G1 also rename the port`",
		"- Verdict: ACTIONABLE.",
		"- WHAT: Move the index rebuild into internal/index; also rename the port",
		"  Keep the exported name.",
		"  internal/a.go (marker at line 12)",
		"  internal/port.go (marker at line 8)\n    func Port() int {",
		"  - GIVEN a build WHEN go test ./... THEN it passes.",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("grown block is missing %q:\n%s", want, s)
		}
	}
}

func TestVerifyRejectsAJudgementARunRewrote(t *testing.T) {
	first, _ := render(t, nil, []Finding{
		grouped("a.go", 1, "TODO", "G1", "TODO:G1 one", ""),
	})
	judged := judge(string(first), "Move the store", "Move the store.")
	found := []Finding{
		grouped("a.go", 1, "TODO", "G1", "TODO:G1 one", ""),
		grouped("b.go", 2, "TODO", "G1", "TODO:G1 two", ""),
	}
	out, res, err := Render([]byte(judged), found)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	for name, broken := range map[string]string{
		"a rewritten verdict": strings.Replace(string(out), "- Verdict: ACTIONABLE.", "- Verdict:", 1),
		"a rewritten WHY":     strings.Replace(string(out), "- WHY:", "- WHY: invented", 1),
		"a shortened WHAT":    strings.Replace(string(out), "- WHAT: Move the store; two", "- WHAT: two", 1),
		"a dropped marker": strings.Replace(string(out),
			"- Marker: `a.go:1` `TODO:G1 one`\n", "", 1),
		"a rewritten WHERE line": strings.Replace(string(out),
			"  a.go (marker at line 1)", "  z.go (marker at line 1)", 1),
	} {
		if broken == string(out) {
			t.Fatalf("%s: the test did not change anything:\n%s", name, out)
		}
		if err := Verify([]byte(judged), []byte(broken), res); err == nil {
			t.Errorf("Verify accepted %s:\n%s", name, broken)
		}
	}
	// The edit Render actually made must pass.
	if err := Verify([]byte(judged), out, res); err != nil {
		t.Fatalf("Verify rejected a clean growth: %v", err)
	}
}

func TestRenderNumbersAroundATagBlock(t *testing.T) {
	first, _ := render(t, nil, []Finding{
		grouped("a.go", 1, "TODO", "G1", "TODO:G1 one", ""),
	})
	_, res := render(t, first, []Finding{
		find("b.go", 2, "TODO", "TODO two", ""),
	})
	if got := blockIDs(res.New); got != "TODO-CMT-01" {
		t.Fatalf("new IDs = %v, want numbering to ignore the tag block", got)
	}
}

func TestRenderStripsBackticksFromTheMarkerLine(t *testing.T) {
	out, _ := render(t, nil, []Finding{
		find("a.go", 3, "TODO", "TODO drop `oldName` and use NewName", ""),
	})
	want := "- Marker: `a.go:3` `TODO drop 'oldName' and use NewName`"
	if !strings.Contains(string(out), want) {
		t.Fatalf("missing %q:\n%s", want, out)
	}
	entries, err := ParseEntries(out)
	if err != nil || len(entries) != 1 {
		t.Fatalf("ParseEntries: %v, entries = %+v", err, entries)
	}
}

func TestParseEntriesReadsEveryMarkerOfABlock(t *testing.T) {
	out, _ := render(t, nil, []Finding{
		grouped("a.go", 1, "TODO", "G1", "TODO:G1 one", ""),
		grouped("b.go", 2, "TODO", "G1", "TODO:G1 two", ""),
	})
	entries, err := ParseEntries(out)
	if err != nil {
		t.Fatalf("ParseEntries: %v", err)
	}
	if len(entries) != 1 || len(entries[0].Members) != 2 {
		t.Fatalf("entries = %+v, want one block with two members", entries)
	}
	e := entries[0]
	if e.Marker != "TODO" || e.Group != "G1" {
		t.Fatalf("entry = %+v, want the marker word and the tag off the ID", e)
	}
	if e.Members[1].File != "b.go" || e.Members[1].Line != 2 {
		t.Fatalf("second member = %+v", e.Members[1])
	}
}

func TestRenderRejectsABlockWithoutAMarkerLine(t *testing.T) {
	broken := []byte("# Comment register\n\n### TODO-CMT-01: No marker line\n\n- WHAT: something\n")
	if _, _, err := Render(broken, []Finding{find("a.go", 1, "TODO", "TODO one", "")}); err == nil {
		t.Fatal("Render accepted a register block with no marker line")
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

func TestRenderGrowsABlockInCRLF(t *testing.T) {
	first, _ := render(t, nil, []Finding{grouped("a.go", 1, "TODO", "G1", "TODO:G1 one", "")})
	crlf := []byte(strings.ReplaceAll(string(first), "\n", "\r\n"))
	out, res := render(t, crlf, []Finding{
		grouped("a.go", 1, "TODO", "G1", "TODO:G1 one", ""),
		grouped("b.go", 2, "TODO", "G1", "TODO:G1 two", ""),
	})
	if len(res.Extended) != 1 {
		t.Fatalf("result = %+v", res)
	}
	if strings.Contains(strings.ReplaceAll(string(out), "\r\n", ""), "\n") {
		t.Fatalf("mixed line endings:\n%q", out)
	}
}

func TestTitleDraftsAnInstruction(t *testing.T) {
	cases := []struct {
		text, group, want string
	}{
		{"TODO We should move this function into a separate package, and use another naming.", "",
			"Move this function into a separate package"},
		{"TODO: fix the retry loop", "", "Fix the retry loop"},
		{"TODO", "", "Handle the TODO marker in a.go"},
		{"TODO:G1 move the rebuild into internal", "G1", "Move the rebuild into internal"},
		{"TODO:G1", "G1", "Handle the TODO marker in a.go"},
	}
	for _, c := range cases {
		got := title(Finding{File: "a.go", Marker: "TODO", Group: c.group, Text: c.text})
		if got != c.want {
			t.Errorf("title(%q) = %q, want %q", c.text, got, c.want)
		}
	}
}

// judge stands in for /arcdlc:assist: it writes the verdict, the title and WHAT that
// the sweep must never touch again.
func judge(register, heading, what string) string {
	lines := strings.Split(register, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "### "):
			lines[i] = line[:strings.Index(line, ": ")+2] + heading
		case strings.HasPrefix(line, "- WHAT:"):
			lines[i] = "- WHAT: " + what
		case strings.HasPrefix(line, "- Verdict:"):
			lines[i] = "- Verdict: ACTIONABLE."
		}
	}
	return strings.Join(lines, "\n")
}

func blockIDs(blocks []Block) string {
	var ids []string
	for _, b := range blocks {
		ids = append(ids, b.ID)
	}
	return strings.Join(ids, " ")
}
