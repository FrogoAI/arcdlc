package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FrogoAI/arcdlc/internal/plan"
)

// mkPlan creates an empty plan.md at dir/rel, making parent directories as needed.
func mkPlan(t *testing.T, dir, rel string) string {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte("# plan\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return full
}

// captureStderr runs fn with os.Stderr redirected to a pipe and returns what it wrote.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// captureStdout runs fn with os.Stdout redirected to a pipe and returns what it wrote.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestResolvePlan(t *testing.T) {
	t.Run("explicit --plan wins over everything", func(t *testing.T) {
		dir := t.TempDir()
		mkPlan(t, dir, "plan.md")          // a flat plan also present
		mkPlan(t, dir, "checkout/plan.md") // and a folder
		got, code := resolvePlan(dir, "some/explicit/path.md", "checkout")
		if code != 0 || got != "some/explicit/path.md" {
			t.Fatalf("got (%q,%d), want (some/explicit/path.md,0)", got, code)
		}
	})

	t.Run("--aic maps to the folder plan", func(t *testing.T) {
		dir := t.TempDir()
		want := filepath.Join(dir, "pay", "plan.md")
		got, code := resolvePlan(dir, "", "pay")
		if code != 0 || got != want {
			t.Fatalf("got (%q,%d), want (%q,0)", got, code, want)
		}
	})

	t.Run("--plan path is used verbatim", func(t *testing.T) {
		got, code := resolvePlan(t.TempDir(), "any/where/plan.md", "")
		if code != 0 || got != "any/where/plan.md" {
			t.Fatalf("got (%q,%d), want (any/where/plan.md,0)", got, code)
		}
	})

	t.Run("bad slug is a usage error", func(t *testing.T) {
		dir := t.TempDir()
		for _, bad := range []string{"../evil", "a/b", "..", ".", `a\b`} {
			if got, code := resolvePlan(dir, "", bad); code != 2 {
				t.Errorf("slug %q: got (%q,%d), want exit 2", bad, got, code)
			}
		}
	})

	t.Run("no selection with one initiative -> exit 2 and lists it", func(t *testing.T) {
		dir := t.TempDir()
		mkPlan(t, dir, "checkout/plan.md")
		var got string
		var code int
		stderr := captureStderr(t, func() { got, code = resolvePlan(dir, "", "") })
		if code != 2 {
			t.Fatalf("got (%q,%d), want exit 2 (selection is mandatory, never auto-detected)", got, code)
		}
		if !strings.Contains(stderr, "checkout") {
			t.Errorf("stderr does not list the available slug 'checkout':\n%s", stderr)
		}
	})

	t.Run("no selection with only a legacy flat plan -> exit 2", func(t *testing.T) {
		dir := t.TempDir()
		mkPlan(t, dir, "plan.md")
		if got, code := resolvePlan(dir, "", ""); code != 2 {
			t.Fatalf("got (%q,%d), want exit 2 (flat plan reachable only via --plan)", got, code)
		}
	})

	t.Run("no selection with multiple initiatives -> exit 2", func(t *testing.T) {
		dir := t.TempDir()
		mkPlan(t, dir, "checkout/plan.md")
		mkPlan(t, dir, "pay/plan.md")
		if got, code := resolvePlan(dir, "", ""); code != 2 {
			t.Fatalf("got (%q,%d), want exit 2", got, code)
		}
	})

	t.Run("no selection and nothing found -> exit 2 (not 3)", func(t *testing.T) {
		dir := t.TempDir()
		if got, code := resolvePlan(dir, "", ""); code != 2 {
			t.Fatalf("got (%q,%d), want exit 2", got, code)
		}
	})
}

func TestListInitiatives(t *testing.T) {
	dir := t.TempDir()
	mkPlan(t, dir, "pay/plan.md")
	mkPlan(t, dir, "checkout/plan.md")
	mkPlan(t, dir, "plan.md") // legacy flat plan has no slug, must be excluded
	got := listInitiatives(dir)
	want := []string{"checkout", "pay"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("listInitiatives = %v, want %v (sorted, flat excluded)", got, want)
	}
}

func TestValidSlug(t *testing.T) {
	good := []string{"checkout", "pay-v2", "AIC_1", "a.b", "123"}
	bad := []string{"", ".", "..", "a/b", "../x", "a\\b", "x..y"}
	for _, s := range good {
		if !validSlug(s) {
			t.Errorf("validSlug(%q) = false, want true", s)
		}
	}
	for _, s := range bad {
		if validSlug(s) {
			t.Errorf("validSlug(%q) = true, want false", s)
		}
	}
}

func TestSlugOf(t *testing.T) {
	if s := slugOf("docs/aics", "docs/aics/plan.md"); s != "" {
		t.Errorf("flat plan slug = %q, want empty", s)
	}
	if s := slugOf("docs/aics", "docs/aics/checkout/plan.md"); s != "checkout" {
		t.Errorf("folder slug = %q, want checkout", s)
	}
}

func TestRunSync(t *testing.T) {
	root := t.TempDir()
	aics := filepath.Join(root, "docs", "aics")
	agents := filepath.Join(root, "AGENTS.md")
	readme := filepath.Join(root, "README.md")
	targets := []string{agents, readme}

	mkInit := func(slug, title, summary string) {
		dir := filepath.Join(aics, slug)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		body := "# " + title + "\n\n> " + summary + "\n"
		if err := os.WriteFile(filepath.Join(dir, "aic.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mkInit("pay", "Payments", "take money")
	mkInit("checkout", "Checkout", "buy flow")

	if code := runSync(aics, targets, false, io.Discard, io.Discard); code != 0 {
		t.Fatalf("sync write exit=%d, want 0", code)
	}
	for _, f := range targets {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, "[Checkout]") || !strings.Contains(s, "[Payments]") {
			t.Errorf("%s missing an initiative:\n%s", f, s)
		}
		if strings.Index(s, "[Checkout]") > strings.Index(s, "[Payments]") {
			t.Errorf("%s not sorted by slug", f)
		}
	}

	// Fresh tree: --check reports no drift.
	if code := runSync(aics, targets, true, io.Discard, io.Discard); code != 0 {
		t.Fatalf("--check on fresh tree exit=%d, want 0", code)
	}
	// Remove an initiative: --check now reports drift (non-zero).
	if err := os.RemoveAll(filepath.Join(aics, "pay")); err != nil {
		t.Fatal(err)
	}
	if code := runSync(aics, targets, true, io.Discard, io.Discard); code == 0 {
		t.Fatalf("--check after change exit=0, want non-zero")
	}
}

func TestRunSyncHTMLOnlyInitiative(t *testing.T) {
	// An initiative whose only architecture document is HTML (written by
	// /arcdlc:aic <slug> arc42:html) must still reach the registry: scanInitiatives
	// once gated on the folder holding a .md file, which skipped these silently.
	root := t.TempDir()
	aics := filepath.Join(root, "docs", "aics")
	agents := filepath.Join(root, "AGENTS.md")
	dir := filepath.Join(aics, "payments")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := `<!DOCTYPE html><html><body>` +
		`<h1>Payments Platform &#8212; arc42</h1>` +
		`<div class="paragraph"><p>Card authorization and capture.</p></div>` +
		`</body></html>`
	if err := os.WriteFile(filepath.Join(dir, "arc42.html"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}

	if code := runSync(aics, []string{agents}, false, io.Discard, io.Discard); code != 0 {
		t.Fatalf("sync exit=%d, want 0", code)
	}
	b, err := os.ReadFile(agents)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{
		"[Payments Platform — arc42]",
		"docs/aics/payments/arc42.html",
		"Card authorization and capture.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("registry missing %q:\n%s", want, got)
		}
	}
	if code := runSync(aics, []string{agents}, true, io.Discard, io.Discard); code != 0 {
		t.Fatalf("--check after write exit=%d, want 0 (not idempotent)", code)
	}
}

func TestRunSyncNoInitiativesStub(t *testing.T) {
	root := t.TempDir()
	aics := filepath.Join(root, "docs", "aics")
	if err := os.MkdirAll(aics, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := filepath.Join(root, "README.md") // missing -> stub created
	if code := runSync(aics, []string{readme}, false, io.Discard, io.Discard); code != 0 {
		t.Fatalf("exit=%d, want 0", code)
	}
	b, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "_none_") {
		t.Errorf("no-initiative registry should render _none_:\n%s", b)
	}
}

// --- arctool order ---

// orderPlan writes a plan.md holding one minimal task block per id, in the
// given order, into a fresh temp dir and returns its path.
func orderPlan(t *testing.T, ids ...string) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("# plan\n")
	for _, id := range ids {
		fmt.Fprintf(&b, "\n### %s: Task %s\n- WHAT: x.\n- WHERE: internal/%s.go\n- WHY: y.\n- References: `a`.\n- Status: TODO.\n", id, id, id)
	}
	path := filepath.Join(t.TempDir(), "plan.md")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// runOrder calls cmdOrder with both standard streams captured.
func runOrder(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	stderr = captureStderr(t, func() {
		stdout = captureStdout(t, func() { code = cmdOrder(args) })
	})
	return code, stdout, stderr
}

// readFile returns the file's bytes as a string, so a test can assert that a
// refused or dry run left them untouched.
func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// planIDs returns the task IDs of the plan file at path, in file order.
func planIDs(t *testing.T, path string) string {
	t.Helper()
	var ids []string
	for _, task := range plan.Parse([]byte(readFile(t, path))).Tasks {
		ids = append(ids, task.ID)
	}
	return strings.Join(ids, " ")
}

func TestCmdOrderWritesThePermutation(t *testing.T) {
	path := orderPlan(t, "T1", "T2", "T3")
	code, stdout, stderr := runOrder(t, "T3", "T1", "T2", "--plan", path)
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if got := planIDs(t, path); got != "T3 T1 T2" {
		t.Fatalf("order on disk = %q, want %q", got, "T3 T1 T2")
	}
	if want := "reordered 3 task(s) in " + path; !strings.Contains(stdout, want) {
		t.Errorf("stdout missing %q:\n%s", want, stdout)
	}
	for _, want := range []string{"  1  T3", "  2  T1", "  3  T2"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing order line %q:\n%s", want, stdout)
		}
	}
}

func TestCmdOrderDryRunWritesNothing(t *testing.T) {
	path := orderPlan(t, "T1", "T2", "T3")
	before := readFile(t, path)
	code, stdout, stderr := runOrder(t, "T3", "T1", "--dry-run", "--plan", path)
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if after := readFile(t, path); after != before {
		t.Fatalf("--dry-run rewrote the file:\n%s", after)
	}
	if want := "would reorder " + path + ":"; !strings.Contains(stdout, want) {
		t.Errorf("stdout missing %q:\n%s", want, stdout)
	}
	// Slot permutation (ADR-0013): T3 and T1 swap positions 1 and 3, T2 stays.
	for _, want := range []string{"  1  T3", "  2  T2", "  3  T1"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("preview missing order line %q:\n%s", want, stdout)
		}
	}
}

func TestCmdOrderNoOpPrintsTheCurrentOrder(t *testing.T) {
	path := orderPlan(t, "T1", "T2", "T3")
	before := readFile(t, path)
	// T1 already holds the first named slot, so naming T1 T3 changes nothing.
	code, stdout, stderr := runOrder(t, "T1", "T3", "--plan", path)
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if after := readFile(t, path); after != before {
		t.Fatalf("a no-op rewrote the file:\n%s", after)
	}
	if !strings.Contains(stdout, "already in that order") {
		t.Errorf("stdout missing the no-op notice:\n%s", stdout)
	}
	for _, want := range []string{"  1  T1", "  2  T2", "  3  T3"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing order line %q:\n%s", want, stdout)
		}
	}
}

func TestCmdOrderRefusesAndLeavesTheFileAlone(t *testing.T) {
	cases := []struct {
		name string
		ids  []string
		want int // 2 = usage error, 3 = not found
	}{
		{"one id", []string{"T1"}, 2},
		{"repeated id", []string{"T1", "T1"}, 2},
		{"unknown id", []string{"T3", "T9"}, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := orderPlan(t, "T1", "T2", "T3")
			before := readFile(t, path)
			args := append(append([]string{}, c.ids...), "--plan", path)
			code, _, stderr := runOrder(t, args...)
			if code != c.want {
				t.Fatalf("exit=%d, want %d (stderr: %s)", code, c.want, stderr)
			}
			if after := readFile(t, path); after != before {
				t.Fatalf("a refused order rewrote the file:\n%s", after)
			}
			if stderr == "" {
				t.Error("a refused order said nothing on stderr")
			}
		})
	}
}

func TestCmdOrderEmptyPlanIsNotFound(t *testing.T) {
	path := mkPlan(t, t.TempDir(), "plan.md") // "# plan\n": no "###" blocks
	code, _, stderr := runOrder(t, "T1", "T2", "--plan", path)
	if code != 3 {
		t.Fatalf("exit=%d, want 3 (stderr: %s)", code, stderr)
	}
}

func TestCmdOrderAcceptsPositionalsBeforeFlags(t *testing.T) {
	// Proves splitArgs is used: the stdlib flag package stops at the first
	// positional, which would leave --plan unparsed and fail selection with 2.
	path := orderPlan(t, "T1", "T2", "T3")
	code, _, stderr := runOrder(t, "T3", "T1", "--plan", path)
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if got := planIDs(t, path); got != "T3 T2 T1" {
		t.Fatalf("order on disk = %q, want %q", got, "T3 T2 T1")
	}
}

func TestUsageDocumentsOrderAndTheWidenedExitFive(t *testing.T) {
	help := fmt.Sprintf(usage, version) // exactly what `arctool help` prints
	for _, want := range []string{
		"arctool order  <id> <id> [<id>…] [--dry-run] [--aic SLUG | --plan PATH]",
		"slot permutation: named tasks swap among the positions they already hold",
		`T1 T2 T3 + "order T3 T1 T2" -> T3 T1 T2; a task you do not name never moves`,
		"5  self-validation failed (nothing written)",
	} {
		if !strings.Contains(help, want) {
			t.Errorf("help output missing %q:\n%s", want, help)
		}
	}
	if bad := "5  archive self-validation failed"; strings.Contains(help, bad) {
		t.Errorf("help output still carries the pre-ADR-0014 line %q", bad)
	}
}

// --- scan ---

// scanTree writes a source file with markers plus an initiative folder, and
// returns the tree root, the plan path that selects the initiative, and the
// register path scan should write beside it.
func scanTree(t *testing.T, source string) (root, planPath, register string) {
	t.Helper()
	root = t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	planPath = mkPlan(t, root, filepath.Join("docs", "aics", "demo", "plan.md"))
	return root, planPath, filepath.Join(filepath.Dir(planPath), "comments.md")
}

func runScan(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	stderr = captureStderr(t, func() {
		stdout = captureStdout(t, func() { code = cmdScan(args) })
	})
	return code, stdout, stderr
}

const twoMarkers = "package a\n\n// TODO split the store\nfunc a() {}\n\n// FIXME retry never backs off\nfunc b() {}\n"

func TestCmdScanNoSelectionIsUsageError(t *testing.T) {
	root, _, _ := scanTree(t, twoMarkers)
	code, _, stderr := runScan(t, "--path", root)
	if code != 2 {
		t.Fatalf("exit=%d, want 2", code)
	}
	if !strings.Contains(stderr, "no initiative selected") {
		t.Errorf("stderr missing the selection error:\n%s", stderr)
	}
}

func TestCmdScanWritesTheRegister(t *testing.T) {
	root, planPath, register := scanTree(t, twoMarkers)
	code, stdout, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register)
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "TODO-CMT-01") || !strings.Contains(stdout, "1 new") {
		t.Errorf("stdout does not report the new block:\n%s", stdout)
	}
	got := readFile(t, register)
	for _, want := range []string{
		"# Comment register",
		"### TODO-CMT-01: Split the store",
		"- Marker: `a.go:3` `TODO split the store`",
		"- Verdict: NEW.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("register missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "FIXME") {
		t.Errorf("FIXME was swept without being asked for:\n%s", got)
	}
}

func TestCmdScanDryRunWritesNothing(t *testing.T) {
	root, planPath, register := scanTree(t, twoMarkers)
	code, stdout, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register, "--dry-run")
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "would write") {
		t.Errorf("stdout does not say it wrote nothing:\n%s", stdout)
	}
	if _, err := os.Stat(register); !os.IsNotExist(err) {
		t.Fatalf("dry run created %s", register)
	}
}

func TestCmdScanNoMarkerIsNotFound(t *testing.T) {
	root, planPath, register := scanTree(t, "package a\n\nfunc a() {}\n")
	code, _, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register)
	if code != 3 {
		t.Fatalf("exit=%d, want 3", code)
	}
	if !strings.Contains(stderr, "no TODO marker found") {
		t.Errorf("stderr missing the not-found line:\n%s", stderr)
	}
	if _, err := os.Stat(register); !os.IsNotExist(err) {
		t.Fatalf("an empty sweep created %s", register)
	}
}

func TestCmdScanSecondRunFindsNothingAndKeepsTheRegister(t *testing.T) {
	root, planPath, register := scanTree(t, twoMarkers)
	if code, _, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register,
		"--marker", "TODO,FIXME"); code != 0 {
		t.Fatalf("first run exit=%d (stderr: %s)", code, stderr)
	}
	first := readFile(t, register)
	if strings.Count(first, "### ") != 2 {
		t.Fatalf("register does not hold both findings:\n%s", first)
	}
	// Both comments are gone from the code now, so there is nothing left to sweep
	// and the register must not move a byte.
	code, _, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register,
		"--marker", "TODO,FIXME")
	if code != 3 {
		t.Fatalf("second run exit=%d, want 3 (stderr: %s)", code, stderr)
	}
	if got := readFile(t, register); got != first {
		t.Fatalf("the second run rewrote the register:\n%s", got)
	}
}

func TestCmdScanRemovesTheMarkerFromTheCode(t *testing.T) {
	root, planPath, register := scanTree(t, twoMarkers)
	code, stdout, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register)
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "removed 1 comment(s) from 1 file(s)") {
		t.Errorf("stdout does not report the removal:\n%s", stdout)
	}
	if !strings.Contains(stdout, "review it with git diff") {
		t.Errorf("stdout does not point at the code change:\n%s", stdout)
	}
	got := readFile(t, filepath.Join(root, "a.go"))
	if strings.Contains(got, "TODO") {
		t.Fatalf("the marker is still in the code:\n%s", got)
	}
	for _, want := range []string{"package a", "func a() {}", "// FIXME retry never backs off", "func b() {}"} {
		if !strings.Contains(got, want) {
			t.Errorf("the strip removed more than the marker comment, %q is gone:\n%s", want, got)
		}
	}
	if !strings.Contains(readFile(t, register), "TODO split the store") {
		t.Error("the register does not hold the removed marker text")
	}
	// A second run has nothing left to sweep.
	code, stdout, _ = runScan(t, "--path", root, "--plan", planPath, "--comments", register)
	if code != 3 {
		t.Fatalf("second run exit=%d, want 3: %s", code, stdout)
	}
}

func TestCmdScanDryRunLeavesTheCodeAlone(t *testing.T) {
	root, planPath, register := scanTree(t, twoMarkers)
	before := readFile(t, filepath.Join(root, "a.go"))
	code, stdout, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register, "--dry-run")
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "would remove 1 comment(s)") {
		t.Errorf("stdout does not say what it would remove:\n%s", stdout)
	}
	if readFile(t, filepath.Join(root, "a.go")) != before {
		t.Fatal("a dry run edited the code")
	}
	if _, err := os.Stat(register); !os.IsNotExist(err) {
		t.Fatal("a dry run wrote the register")
	}
}

func TestCmdScanJSONShape(t *testing.T) {
	root, planPath, register := scanTree(t, twoMarkers)
	code, stdout, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register, "--json")
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	var got scanJSON
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, stdout)
	}
	if got.Found != 1 || got.Total != 1 || !got.Changed || len(got.New) != 1 {
		t.Fatalf("payload = %+v", got)
	}
	if got.New[0].ID != "TODO-CMT-01" || got.New[0].Line != 3 || got.New[0].Code != "func a() {}" {
		t.Fatalf("finding = %+v", got.New[0])
	}
	if len(got.Edited) != 1 || got.Edited[0].Removed != 1 {
		t.Fatalf("edited = %+v", got.Edited)
	}
	if got.Skipped == nil {
		t.Error("skipped should be an empty list, not null")
	}
}

func TestCmdScanRejectsARegisterWithoutMarkerLines(t *testing.T) {
	root, planPath, register := scanTree(t, twoMarkers)
	if err := os.WriteFile(register, []byte("# Comment register\n\n### TODO-CMT-01: hand written\n\n- WHAT: x.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, register)
	code, _, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register)
	if code != 1 {
		t.Fatalf("exit=%d, want 1", code)
	}
	if !strings.Contains(stderr, "Marker:") {
		t.Errorf("stderr does not name the broken line:\n%s", stderr)
	}
	if readFile(t, register) != before {
		t.Error("a broken register was rewritten")
	}
}

func TestUsageDocumentsScan(t *testing.T) {
	help := fmt.Sprintf(usage, version)
	for _, want := range []string{
		"arctool scan   [--marker LIST] [--path DIR] [--exclude LIST] [--comments PATH]",
		"sweep source comments for markers (default TODO) into comments.md",
	} {
		if !strings.Contains(help, want) {
			t.Errorf("help output missing %q:\n%s", want, help)
		}
	}
}

func TestCmdScanCreatesAMissingInitiativeFolder(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte(twoMarkers), 0o644); err != nil {
		t.Fatal(err)
	}
	// The initiative has no folder yet: no docs/, no docs/aics/, no slug dir.
	planPath := filepath.Join(root, "docs", "aics", "review", "plan.md")
	register := filepath.Join(filepath.Dir(planPath), "comments.md")

	code, stdout, stderr := runScan(t, "--path", root, "--plan", planPath, "--comments", register)
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "new initiative folder") {
		t.Errorf("stdout does not report the folder it created:\n%s", stdout)
	}
	if !isDir(filepath.Dir(register)) {
		t.Fatalf("%s was not created", filepath.Dir(register))
	}
	if !strings.Contains(readFile(t, register), "### TODO-CMT-01") {
		t.Error("the register was not written into the new folder")
	}
	// A second run finds the folder and says nothing about creating it.
	_, stdout, _ = runScan(t, "--path", root, "--plan", planPath, "--comments", register)
	if strings.Contains(stdout, "initiative folder") {
		t.Errorf("stdout reports a folder that already existed:\n%s", stdout)
	}
}

func TestCmdScanDryRunCreatesNoFolder(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte(twoMarkers), 0o644); err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(root, "docs", "aics", "review", "plan.md")
	code, stdout, stderr := runScan(t, "--path", root, "--plan", planPath, "--dry-run")
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "would create") {
		t.Errorf("stdout does not say the folder would be created:\n%s", stdout)
	}
	if isDir(filepath.Dir(planPath)) {
		t.Fatalf("dry run created %s", filepath.Dir(planPath))
	}
}

func TestCmdScanCreatesTheFolderFromTheSlug(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte(twoMarkers), 0o644); err != nil {
		t.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatal(err)
		}
	}()

	// Exactly what the request names: arctool scan --marker TODO --aic review,
	// with docs/aics/review/ absent.
	code, stdout, stderr := runScan(t, "--marker", "TODO", "--aic", "review")
	if code != 0 {
		t.Fatalf("exit=%d, want 0 (stderr: %s)", code, stderr)
	}
	if !strings.Contains(stdout, "created docs/aics/review/ (new initiative folder)") {
		t.Errorf("stdout does not name the created folder:\n%s", stdout)
	}
	if !strings.Contains(readFile(t, filepath.Join("docs", "aics", "review", "comments.md")), "TODO-CMT-01") {
		t.Error("the register is not in docs/aics/review/")
	}
}
