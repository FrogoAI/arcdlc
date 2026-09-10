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
