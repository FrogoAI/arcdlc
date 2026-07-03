package main

import (
	"os"
	"path/filepath"
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

	t.Run("bad slug is a usage error", func(t *testing.T) {
		// Note: an empty --aic is indistinguishable from "flag absent" and falls
		// through to auto-detect, so it is not tested here (see the exit-3 case).
		dir := t.TempDir()
		for _, bad := range []string{"../evil", "a/b", "..", ".", `a\b`} {
			if got, code := resolvePlan(dir, "", bad); code != 2 {
				t.Errorf("slug %q: got (%q,%d), want exit 2", bad, got, code)
			}
		}
	})

	t.Run("auto-detect: single folder is used", func(t *testing.T) {
		dir := t.TempDir()
		want := mkPlan(t, dir, "checkout/plan.md")
		got, code := resolvePlan(dir, "", "")
		if code != 0 || got != want {
			t.Fatalf("got (%q,%d), want (%q,0)", got, code, want)
		}
	})

	t.Run("auto-detect: legacy flat plan is used", func(t *testing.T) {
		dir := t.TempDir()
		want := mkPlan(t, dir, "plan.md")
		got, code := resolvePlan(dir, "", "")
		if code != 0 || got != want {
			t.Fatalf("got (%q,%d), want (%q,0)", got, code, want)
		}
	})

	t.Run("auto-detect: multiple initiatives is ambiguous (exit 2)", func(t *testing.T) {
		dir := t.TempDir()
		mkPlan(t, dir, "checkout/plan.md")
		mkPlan(t, dir, "pay/plan.md")
		if got, code := resolvePlan(dir, "", ""); code != 2 {
			t.Fatalf("got (%q,%d), want exit 2", got, code)
		}
	})

	t.Run("auto-detect: flat + folder together is ambiguous (exit 2)", func(t *testing.T) {
		dir := t.TempDir()
		mkPlan(t, dir, "plan.md")
		mkPlan(t, dir, "pay/plan.md")
		if got, code := resolvePlan(dir, "", ""); code != 2 {
			t.Fatalf("got (%q,%d), want exit 2", got, code)
		}
	})

	t.Run("auto-detect: nothing found (exit 3)", func(t *testing.T) {
		dir := t.TempDir()
		if got, code := resolvePlan(dir, "", ""); code != 3 {
			t.Fatalf("got (%q,%d), want exit 3", got, code)
		}
	})
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
