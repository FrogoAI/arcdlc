package bundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot walks up from the test's working directory to the module root.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod found above the test directory")
		}
		dir = parent
	}
}

// TestBundleIsConsistent is the gate. It fails when a SKILL.md names a file
// that does not exist, drops a shared block, lets two copies of a block drift,
// forgets a trigger word in its description, or reaches outside skills/.
func TestBundleIsConsistent(t *testing.T) {
	probs, err := Check(repoRoot(t))
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	for _, p := range probs {
		t.Errorf("%s", p)
	}
}

// TestSkillsMatchInstaller pins the sub-skill list against install.sh, because
// a skill missing from SUBSKILLS is built but never installed.
func TestSkillsMatchInstaller(t *testing.T) {
	root := repoRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	var line string
	for _, l := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(l, "SUBSKILLS=") {
			line = l
			break
		}
	}
	if line == "" {
		t.Fatal("install.sh has no SUBSKILLS assignment")
	}
	got := strings.Fields(strings.Trim(strings.TrimPrefix(line, "SUBSKILLS="), `"`))
	if len(got) != len(Skills) {
		t.Fatalf("install.sh SUBSKILLS has %d skills, bundle.Skills has %d:\n  install.sh: %v\n  bundle:     %v",
			len(got), len(Skills), got, Skills)
	}
	for i := range got {
		if got[i] != Skills[i] {
			t.Errorf("SUBSKILLS[%d] = %q, bundle.Skills[%d] = %q", i, got[i], i, Skills[i])
		}
	}
}

// TestRetiredSkillsAreSwept pins that every skill the bundle no longer ships is
// named in LEGACY_SUBSKILLS, so an upgrade removes it instead of stranding it.
func TestRetiredSkillsAreSwept(t *testing.T) {
	root := repoRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "LEGACY_SUBSKILLS=") {
		t.Fatal("install.sh has no LEGACY_SUBSKILLS; a retired skill would be left behind on upgrade")
	}
	// Both the install and the uninstall path must sweep it.
	body := string(raw)
	if !strings.Contains(body, "$SUBSKILLS $LEGACY_SUBSKILLS") {
		t.Error("uninstall does not iterate LEGACY_SUBSKILLS")
	}
	if strings.Count(body, "for s in $LEGACY_SUBSKILLS") < 2 {
		t.Error("the flat install paths do not both sweep LEGACY_SUBSKILLS")
	}
}
