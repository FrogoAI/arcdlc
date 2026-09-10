package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stripFile writes one source file, sweeps it, strips it, and returns the new
// content plus whatever the strip refused to touch.
func stripFile(t *testing.T, name, source string, markers ...string) (string, []Skip) {
	t.Helper()
	if len(markers) == 0 {
		markers = []string{"TODO"}
	}
	dir := tree(t, map[string]string{name: source})
	found, err := Sweep(Opts{Root: dir, Markers: markers})
	if err != nil {
		t.Fatalf("Sweep: %v", err)
	}
	edits, skipped, err := Strip(dir, found)
	if err != nil {
		t.Fatalf("Strip: %v", err)
	}
	if len(edits) == 0 {
		return source, skipped
	}
	if len(edits) != 1 || edits[0].File != name {
		t.Fatalf("edits = %+v, want one for %s", edits, name)
	}
	return string(edits[0].Bytes), skipped
}

func TestStripRemovesAStandaloneComment(t *testing.T) {
	got, skipped := stripFile(t, "a.go", "package a\n\n// TODO rename this\nfunc A() {}\n")
	want := "package a\n\nfunc A() {}\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if len(skipped) != 0 {
		t.Fatalf("skipped = %+v", skipped)
	}
}

func TestStripRemovesContinuationLines(t *testing.T) {
	got, _ := stripFile(t, "a.go", "package a\n\n"+
		"// TODO We should move this into a separate package,\n"+
		"// and use another naming.\n"+
		"func rebuild() {}\n")
	want := "package a\n\nfunc rebuild() {}\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripKeepsTheCodeOnATrailingComment(t *testing.T) {
	got, _ := stripFile(t, "a.go", "package a\n\nfunc A() {} // TODO split this\n")
	want := "package a\n\nfunc A() {}\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripKeepsTheCodeAroundABlockComment(t *testing.T) {
	got, _ := stripFile(t, "a.c", "int a = 1 /* TODO widen this */ + 2;\n")
	want := "int a = 1 + 2;\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripRemovesASingleLineBlockComment(t *testing.T) {
	got, _ := stripFile(t, "a.c", "/* TODO free the buffer */\nint f(void) { return 0; }\n")
	want := "int f(void) { return 0; }\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripLeavesMultiLineBlockComments(t *testing.T) {
	source := "/*\n * TODO free the buffer\n */\nint f(void) { return 0; }\n"
	got, skipped := stripFile(t, "a.c", source)
	if got != source {
		t.Fatalf("the code was edited inside a block comment:\n%s", got)
	}
	if len(skipped) != 1 || !strings.Contains(skipped[0].Reason, "block comment") {
		t.Fatalf("skipped = %+v", skipped)
	}
	if skipped[0].Line != 2 || skipped[0].Marker != "TODO" {
		t.Fatalf("skip = %+v", skipped[0])
	}
}

func TestStripLeavesAnUnterminatedBlockOpener(t *testing.T) {
	source := "/* TODO free the buffer\n   and check the size\nint f(void) { return 0; }\n"
	got, skipped := stripFile(t, "a.c", source)
	if got != source {
		t.Fatalf("an unterminated block comment was cut:\n%s", got)
	}
	if len(skipped) != 1 || !strings.Contains(skipped[0].Reason, "does not close") {
		t.Fatalf("skipped = %+v", skipped)
	}
}

func TestStripRemovesABlockCommentThatClosesLowerDown(t *testing.T) {
	got, skipped := stripFile(t, "a.c", "/* TODO free the buffer\nand check the size\n*/\nint f(void) { return 0; }\n")
	want := "int f(void) { return 0; }\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if len(skipped) != 0 {
		t.Fatalf("skipped = %+v, want none", skipped)
	}
}

func TestStripKeepsTheCodeBeforeABlockCommentThatClosesLowerDown(t *testing.T) {
	source := "int f(void) { /* TODO free the buffer\nand check the size\n*/\n  return 0; }\n"
	got, _ := stripFile(t, "a.c", source)
	want := "int f(void) {\n  return 0; }\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripLeavesABlockCommentWithCodeAfterItsCloser(t *testing.T) {
	// Cutting the closer off that line would take the code's indentation with it.
	source := "/* TODO free the buffer\nand check the size\n*/ int f(void) { return 0; }\n"
	got, skipped := stripFile(t, "a.c", source)
	if got != source {
		t.Fatalf("the code line was cut:\n%s", got)
	}
	if len(skipped) != 1 || !strings.Contains(skipped[0].Reason, "code follows the closing */") {
		t.Fatalf("skipped = %+v", skipped)
	}
}

func TestStripCollapsesTheBlankLineItWouldDouble(t *testing.T) {
	got, _ := stripFile(t, "a.go", "package a\n\nfunc A() {}\n\n// TODO drop B\n\nfunc B() {}\n")
	want := "package a\n\nfunc A() {}\n\nfunc B() {}\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripRemovesSeveralMarkersFromOneFile(t *testing.T) {
	got, _ := stripFile(t, "a.py", "# TODO one\ndef a(): pass\n\n# TODO two\ndef b(): pass\n")
	want := "def a(): pass\n\ndef b(): pass\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripKeepsCRLF(t *testing.T) {
	got, _ := stripFile(t, "a.go", "package a\r\n\r\n// TODO rename this\r\nfunc A() {}\r\n")
	want := "package a\r\n\r\nfunc A() {}\r\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStripSkipsAFileThatChangedSinceTheSweep(t *testing.T) {
	dir := tree(t, map[string]string{"a.go": "package a\n\n// TODO rename this\nfunc A() {}\n"})
	found, err := Sweep(Opts{Root: dir, Markers: []string{"TODO"}})
	if err != nil {
		t.Fatal(err)
	}
	// Somebody edits the file between the sweep and the strip.
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n\nfunc A() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	edits, skipped, err := Strip(dir, found)
	if err != nil {
		t.Fatalf("Strip: %v", err)
	}
	if len(edits) != 0 {
		t.Fatalf("edits = %+v, want none", edits)
	}
	if len(skipped) != 1 || !strings.Contains(skipped[0].Reason, "changed since the sweep") {
		t.Fatalf("skipped = %+v", skipped)
	}
}

func TestVerifyEditCatchesACorruptedResult(t *testing.T) {
	orig := []byte("package a\n\n// TODO rename this\nfunc A() {}\n")
	cases := []struct {
		name string
		e    Edit
	}{
		{"a line of code went missing", Edit{Bytes: []byte("package a\n\n"), deleted: 1}},
		{"a line was added", Edit{
			Bytes:   []byte("package a\n\nfunc A() {}\nfunc B() {}\n"),
			deleted: 1,
		}},
		{"a line was rewritten", Edit{Bytes: []byte("package a\n\nfunc Z() {}\n"), deleted: 1}},
		{"the lines were reordered", Edit{Bytes: []byte("func A() {}\n\npackage a\n"), deleted: 1}},
	}
	for _, c := range cases {
		if err := VerifyEdit(orig, c.e); err == nil {
			t.Errorf("VerifyEdit accepted %s", c.name)
		}
	}
	// The real edit passes.
	good := Edit{Bytes: []byte("package a\n\nfunc A() {}\n"), deleted: 1}
	if err := VerifyEdit(orig, good); err != nil {
		t.Errorf("VerifyEdit rejected a correct edit: %v", err)
	}
}
