package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tree writes files into a fresh directory and returns its path.
func tree(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		full := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func sweep(t *testing.T, dir string, markers ...string) []Finding {
	t.Helper()
	found, err := Sweep(Opts{Root: dir, Markers: markers})
	if err != nil {
		t.Fatalf("Sweep: %v", err)
	}
	return found
}

func TestSweepFindsMarkersBehindEveryOpener(t *testing.T) {
	dir := tree(t, map[string]string{
		"go/a.go":     "package a\n\n// TODO rename this\nfunc A() {}\n",
		"py/b.py":     "# TODO drop the shim\ndef b(): pass\n",
		"sql/c.sql":   "-- TODO add an index\nSELECT 1;\n",
		"lisp/d.el":   "; TODO port this\n(defun d ())\n",
		"tex/e.tex":   "% TODO check the table\n\\table\n",
		"c/f.c":       "/* TODO free the buffer */\nint f(void) { return 0; }\n",
		"go/trail.go": "package a\n\nfunc T() {} // TODO split this\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 7 {
		for _, f := range found {
			t.Logf("%s:%d %q", f.File, f.Line, f.Text)
		}
		t.Fatalf("found %d markers, want 7", len(found))
	}
}

func TestSweepIgnoresNonComments(t *testing.T) {
	dir := tree(t, map[string]string{
		// A marker in a string literal has no comment opener before it.
		"a.go": "package a\n\nvar s = \"TODO not a comment\"\nvar t = a * TODO_CONST\n",
		// The plan's own status lines must never register as findings.
		"docs/aics/x/plan.md": "### X-1: t\n- Status: TODO.\n",
		// Documents are not source code.
		"notes.txt": "// TODO written in prose\n",
		"README.md": "// TODO written in prose\n",
	})
	if found := sweep(t, dir, "TODO"); len(found) != 0 {
		t.Fatalf("found %d markers, want 0: %+v", len(found), found)
	}
}

func TestSweepWordBoundary(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\n// TODOLIST is a type name, not a marker\n// TODO: this one counts\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 {
		t.Fatalf("found %d markers, want 1: %+v", len(found), found)
	}
	if !strings.HasPrefix(found[0].Text, "TODO: this one counts") {
		t.Fatalf("text = %q", found[0].Text)
	}
}

func TestSweepJoinsContinuationLines(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\n" +
			"// TODO We should move this function into a separate package,\n" +
			"// and use another naming. The function must be publicly available.\n" +
			"func rebuild() {}\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 {
		t.Fatalf("found %d markers, want 1 (continuation lines are one finding)", len(found))
	}
	want := "TODO We should move this function into a separate package, and use another naming. The function must be publicly available."
	if found[0].Text != want {
		t.Fatalf("text = %q, want %q", found[0].Text, want)
	}
	if found[0].Code != "func rebuild() {}" {
		t.Fatalf("code = %q", found[0].Code)
	}
	if found[0].Line != 3 {
		t.Fatalf("line = %d, want 3", found[0].Line)
	}
}

func TestSweepSecondMarkerEndsTheFirst(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\n// TODO first thing\n// TODO second thing\nfunc a() {}\n",
	})
	if found := sweep(t, dir, "TODO"); len(found) != 2 {
		t.Fatalf("found %d markers, want 2: %+v", len(found), found)
	}
}

func TestSweepSkipsExcludedDirsBinaryAndOversize(t *testing.T) {
	big := "package a\n// TODO too big\n" + strings.Repeat("x", maxFileBytes)
	dir := tree(t, map[string]string{
		"vendor/v.go":       "// TODO vendored\n",
		"node_modules/n.js": "// TODO dependency\n",
		"bin/b.go":          "// TODO built\n",
		"dist/d.go":         "// TODO built\n",
		"big.go":            big,
		"bin.dat":           "// TODO binary\x00payload\n",
		"keep.go":           "// TODO the only real one\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 || found[0].File != "keep.go" {
		t.Fatalf("found %+v, want only keep.go", found)
	}
}

func TestSweepSeveralMarkersAndStableOrder(t *testing.T) {
	dir := tree(t, map[string]string{
		"b.go": "// FIXME broken\n",
		"a.go": "// TODO later\n// FIXME also here\n",
	})
	found := sweep(t, dir, "TODO", "FIXME")
	var got []string
	for _, f := range found {
		got = append(got, f.File+":"+f.Marker)
	}
	want := "a.go:TODO a.go:FIXME b.go:FIXME"
	if strings.Join(got, " ") != want {
		t.Fatalf("order = %q, want %q", strings.Join(got, " "), want)
	}
}

func TestFingerprintIgnoresLineNumber(t *testing.T) {
	a := Finding{File: "a.go", Line: 10, Marker: "TODO", Text: "TODO  fix   this"}
	b := Finding{File: "a.go", Line: 99, Marker: "TODO", Text: "TODO fix this"}
	if a.Fingerprint() != b.Fingerprint() {
		t.Fatalf("fingerprints differ: %q vs %q", a.Fingerprint(), b.Fingerprint())
	}
}

func TestSweepIgnoresOpenersInsideStringLiterals(t *testing.T) {
	dir := tree(t, map[string]string{
		// A test fixture that quotes a comment must not become a finding.
		"fixture_test.go": "package a\n\nvar files = map[string]string{\n" +
			"\t\"a.go\": \"// TODO quoted, not written\\n\",\n" +
			"\t\"b.go\": \"/* TODO also quoted */\\n\",\n}\n",
		// A real trailing comment on a line that also holds a string stays a finding.
		"real.go": "package a\n\nvar s = \"hello\" // TODO rename s\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 || found[0].File != "real.go" {
		t.Fatalf("found %+v, want only real.go", found)
	}
}

func TestSweepRequiresTheMarkerToOpenTheComment(t *testing.T) {
	dir := tree(t, map[string]string{
		// Prose about markers is not a marker.
		"a.go": "package a\n\n// status token as written, e.g. TODO or Todo\nfunc a() {}\n",
		// A flag that looks like a comment opener in another language.
		"b.go": "package a\n\nvar usage = `arctool next [--json] first TODO block`\n",
		// A format string is not a comment either.
		"c.go": "package a\n\nfunc c() { fmt.Printf(\"want %s (TODO/DONE)\", x) }\n",
		"d.go": "package a\n\n// TODO this one opens the comment\nfunc d() {}\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 || found[0].File != "d.go" {
		t.Fatalf("found %+v, want only d.go", found)
	}
}

func TestSweepReadsDocCommentDecoration(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.rs": "/// TODO rewrite the doc\nfn a() {}\n",
		"b.sh": "#!/bin/sh\n# TODO handle the flag\n",
		"c.c":  "/*\n * TODO free the buffer\n */\nint c(void) { return 0; }\n",
	})
	if found := sweep(t, dir, "TODO"); len(found) != 3 {
		t.Fatalf("found %d, want 3: %+v", len(found), found)
	}
}

func TestSweepRecordsTheCodeATrailingMarkerSitsOn(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\nfunc retry() {} // TODO add a backoff\n\nfunc keep() {}\n",
		"b.c":  "int a = 1 /* TODO widen this */ + 2;\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 2 {
		t.Fatalf("found %d, want 2", len(found))
	}
	if found[0].Code != "func retry() {}" {
		t.Errorf("code = %q, want the line the comment trails", found[0].Code)
	}
	if found[1].Code != "int a = 1 + 2;" {
		t.Errorf("code = %q, want the line without its comment", found[1].Code)
	}
}

func TestSweepReadsTheGroupTag(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\n" +
			"// TODO:G1 move the rebuild into internal\n" +
			"func a() {}\n\n" +
			"// TODO:g1 change the return format\n" +
			"func b() {}\n\n" +
			"// TODO:api-2 rename the handler\n" +
			"func c() {}\n\n" +
			"// TODO drop the retry\n" +
			"func d() {}\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 4 {
		t.Fatalf("found %d markers, want 4", len(found))
	}
	want := []string{"G1", "G1", "API-2", ""}
	for i, w := range want {
		if found[i].Group != w {
			t.Fatalf("finding %d group = %q, want %q (text %q)", i, found[i].Group, w, found[i].Text)
		}
	}
	// The text stays verbatim: file plus text is the finding's identity.
	if found[1].Text != "TODO:g1 change the return format" {
		t.Fatalf("text = %q, want it kept as written", found[1].Text)
	}
}

func TestSweepRejectsTagShapesThatAreNotTags(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\n" +
			"// TODO: drop the retry\n" +
			"func a() {}\n\n" +
			"// TODO:01 renumber this\n" +
			"func b() {}\n\n" +
			"// TODO:refactor(later) split this\n" +
			"func c() {}\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 3 {
		t.Fatalf("found %d markers, want 3", len(found))
	}
	for _, f := range found {
		if f.Group != "" {
			t.Fatalf("%q read as group %q, want no group", f.Text, f.Group)
		}
	}
	// A digits-only tag would claim the auto-numbered ID, so the sweep says so.
	if got := IgnoredTag(found[1]); got != "01" {
		t.Fatalf("IgnoredTag = %q, want %q", got, "01")
	}
	if got := IgnoredTag(found[0]); got != "" {
		t.Fatalf("IgnoredTag on a plain marker = %q, want empty", got)
	}
}

func TestSweepJoinsABlockCommentThatClosesLowerDown(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.c": "/* TODO:G1 free the buffer\n" +
			" * and check the size\n" +
			"before every write\n" +
			"*/\n" +
			"int f(void) { return 0; }\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 {
		t.Fatalf("found %d markers, want 1: %+v", len(found), found)
	}
	f := found[0]
	want := "TODO:G1 free the buffer and check the size before every write"
	if f.Text != want {
		t.Fatalf("text = %q, want %q", f.Text, want)
	}
	if f.Group != "G1" {
		t.Fatalf("group = %q, want G1", f.Group)
	}
	if f.Line != 1 {
		t.Fatalf("line = %d, want 1", f.Line)
	}
	if f.Code != "int f(void) { return 0; }" {
		t.Fatalf("code = %q, want the line under the closer", f.Code)
	}
}

func TestSweepReadsOneBlockCommentAsOneFinding(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.c": "/* TODO free the buffer\nTODO and check the size\n*/\nint f(void) { return 0; }\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 {
		t.Fatalf("found %d markers, want 1: a block comment is one finding", len(found))
	}
	if want := "TODO free the buffer TODO and check the size"; found[0].Text != want {
		t.Fatalf("text = %q, want %q", found[0].Text, want)
	}
}

func TestSweepKeepsTheMarkerLineOfAnUnterminatedBlock(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.c": "/* TODO free the buffer\nand check the size\nint f(void) { return 0; }\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 {
		t.Fatalf("found %d markers, want 1", len(found))
	}
	if want := "TODO free the buffer"; found[0].Text != want {
		t.Fatalf("text = %q, want %q", found[0].Text, want)
	}
}

func TestSweepDropsTheCloserOfASingleLineBlock(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.c": "/* TODO free the buffer */\nint f(void) { return 0; }\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 {
		t.Fatalf("found %d markers, want 1", len(found))
	}
	if want := "TODO free the buffer"; found[0].Text != want {
		t.Fatalf("text = %q, want %q", found[0].Text, want)
	}
}

func TestSweepIgnoresMarkersInsideMultiLineStrings(t *testing.T) {
	dir := tree(t, map[string]string{
		// A Go raw string that spans lines: every line of it is text in a constant,
		// however much it looks like code with comments in it.
		"help.go": "package a\n\nconst usage = `arctool scan\n" +
			"  markers sharing a tag are one block: // TODO:G1 in three files, one task\n" +
			"  // TODO not a marker either\n" +
			"`\n\n// TODO this one is real\nfunc a() {}\n",
		// A Python docstring, same shape with a different delimiter.
		"doc.py": "def f():\n    \"\"\"Docs.\n    # TODO written in a docstring\n    \"\"\"\n    return 1\n\n# TODO this one is real\n",
		// The literal closes, so the comment after it on the same line is a comment.
		"tail.go": "package a\n\nvar s = `raw` // TODO after a closed literal\n",
	})
	found := sweep(t, dir, "TODO")
	var got []string
	for _, f := range found {
		got = append(got, f.File+": "+f.Text)
	}
	want := []string{
		"doc.py: TODO this one is real",
		"help.go: TODO this one is real",
		"tail.go: TODO after a closed literal",
	}
	if len(got) != len(want) {
		t.Fatalf("found %d markers, want %d: %q", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("finding %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSweepIgnoresAMarkerInsideACharLiteralOrString(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\nvar q = '/'\nvar s = \"// TODO quoted\"\nvar t = `// TODO raw`\n" +
			"var u = \"a\\\"b\" // TODO real one\n",
	})
	found := sweep(t, dir, "TODO")
	if len(found) != 1 {
		t.Fatalf("found %d markers, want 1: %+v", len(found), found)
	}
	if found[0].Text != "TODO real one" {
		t.Fatalf("text = %q", found[0].Text)
	}
}

func TestDefaultMarkerIsTheBundlesOwnWord(t *testing.T) {
	dir := tree(t, map[string]string{
		"a.go": "package a\n\n// TODO a note this team already had\nfunc a() {}\n\n" +
			"// ARCDLC:T1 plan this one\nfunc b() {}\n",
	})
	found := sweep(t, dir) // no marker named: the default applies
	if len(found) != 1 {
		t.Fatalf("found %d markers, want only the ARCDLC one: %+v", len(found), found)
	}
	if found[0].Marker != "ARCDLC" || found[0].Group != "T1" {
		t.Fatalf("finding = %+v, want the ARCDLC marker with tag T1", found[0])
	}
}
