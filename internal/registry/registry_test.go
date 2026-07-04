package registry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFiles creates each named file under dir with the given content.
func writeFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFindArchDoc(t *testing.T) {
	cases := []struct {
		name  string
		files []string
		want  string
	}{
		{"aic wins over others", []string{"aic.md", "arc42.md", "plan.md"}, "aic.md"},
		{"arc42 over c4", []string{"arc42.md", "c4.md", "plan.md", "gap.md"}, "arc42.md"},
		{"togaf over c4", []string{"togaf.md", "c4.md"}, "togaf.md"},
		{"alphabetical fallback", []string{"notes.md", "zeta.md"}, "notes.md"},
		{"non-arch files excluded from fallback", []string{"plan.md", "gap.md", "plan-archive.md"}, ""},
		{"custom doc alongside plan", []string{"design.md", "plan.md", "gap.md"}, "design.md"},
		{"empty folder", []string{}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			fm := map[string]string{}
			for _, f := range tc.files {
				fm[f] = "# x\n"
			}
			writeFiles(t, dir, fm)
			if got := findArchDoc(dir); got != tc.want {
				t.Errorf("findArchDoc = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseTitle(t *testing.T) {
	cases := []struct{ in, want string }{
		{"# Payments\n\nbody", "Payments"},
		{"stray line\n# Real Title\n", "Real Title"},
		{"## Not H1\n# Yes H1\n", "Yes H1"},
		{"no heading here\n", ""},
	}
	for _, tc := range cases {
		if got := parseTitle([]byte(tc.in)); got != tc.want {
			t.Errorf("parseTitle(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseSummary(t *testing.T) {
	t.Run("blockquote under H1", func(t *testing.T) {
		got := parseSummary([]byte("# Payments\n\n> Take payments end to end.\n\nmore text\n"))
		if got != "Take payments end to end." {
			t.Errorf("summary = %q", got)
		}
	})
	t.Run("first paragraph fallback", func(t *testing.T) {
		got := parseSummary([]byte("# Payments\n\nThis is the opening paragraph.\nSecond line of it.\n\nNext para.\n"))
		if got != "This is the opening paragraph. Second line of it." {
			t.Errorf("summary = %q", got)
		}
	})
	t.Run("truncated to summaryMax with ellipsis", func(t *testing.T) {
		long := strings.Repeat("word ", 60) // ~300 chars, no blockquote
		got := parseSummary([]byte("# T\n\n" + long + "\n"))
		if !strings.HasSuffix(got, "…") {
			t.Errorf("expected ellipsis, got %q", got)
		}
		if n := len([]rune(got)); n > summaryMax+1 {
			t.Errorf("summary %d runes, want <= %d", n, summaryMax+1)
		}
	})
	t.Run("nothing after H1", func(t *testing.T) {
		if got := parseSummary([]byte("# Only a title\n")); got != "" {
			t.Errorf("summary = %q, want empty", got)
		}
	})
}

func TestLoad(t *testing.T) {
	aics := t.TempDir()

	t.Run("aic with blockquote", func(t *testing.T) {
		dir := filepath.Join(aics, "payments")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFiles(t, dir, map[string]string{
			"aic.md":  "# Payments\n\n> One-line summary.\n\nbody\n",
			"plan.md": "# plan\n",
		})
		got := Load(aics, "payments")
		if got.Title != "Payments" || got.Summary != "One-line summary." {
			t.Errorf("got %+v", got)
		}
		if got.DocRelPath != filepath.ToSlash(filepath.Join(aics, "payments", "aic.md")) {
			t.Errorf("DocRelPath = %q", got.DocRelPath)
		}
	})

	t.Run("no architecture doc", func(t *testing.T) {
		dir := filepath.Join(aics, "bare")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFiles(t, dir, map[string]string{"plan.md": "# plan\n", "gap.md": "# gap\n"})
		got := Load(aics, "bare")
		if got.Title != "bare" || got.Summary != "(no architecture doc)" || got.DocRelPath != "" {
			t.Errorf("got %+v, want slug title + no-doc summary + empty DocRelPath", got)
		}
	})
}
