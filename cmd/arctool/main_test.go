package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
