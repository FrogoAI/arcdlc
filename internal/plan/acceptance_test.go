package plan

import "testing"

// TestLooksDemonstrable pins what counts as a checkable acceptance criterion.
// A task block is run by a weaker model than the one that planned it, so a
// criterion it cannot check is a criterion it will declare met.
func TestLooksDemonstrable(t *testing.T) {
	demonstrable := []string{
		"- GIVEN an empty plan WHEN `arctool next` runs THEN it exits 3.",
		"- `go test ./internal/scan/...` passes.",
		"- TestDocumentedCommentStylesMatchTheTable passes.",
		"- The file internal/scan/strip.go contains no TODO marker.",
		"- Running the binary exits with code 2 on a missing slug.",
		"- `arctool validate --strict` reports no findings.",
		"- cmd/arctool/main.go declares version 0.20.0.",
	}
	for _, s := range demonstrable {
		if !looksDemonstrable(s) {
			t.Errorf("should be demonstrable but is not:\n  %s", s)
		}
	}

	vague := []string{
		"- The feature works correctly.",
		"- The code is clean and maintainable.",
		"- Everything is implemented properly.",
		"- The user is happy with the result.",
		"- Performance is significantly improved.",
		"- It behaves as described above.",
	}
	for _, s := range vague {
		if looksDemonstrable(s) {
			t.Errorf("should be flagged as unverifiable but passed:\n  %s", s)
		}
	}
}

// TestUnverifiableAcceptanceIsStrictOnly checks the rule fires under --strict
// and stays silent otherwise, so existing plans are not broken by an upgrade.
func TestUnverifiableAcceptanceIsStrictOnly(t *testing.T) {
	src := "### T-1: Do the thing\n" +
		"- WHAT: Do the thing.\n" +
		"- WHERE: internal/thing/thing.go\n" +
		"- WHY: The thing is not done.\n" +
		"- Acceptance:\n" +
		"  - The feature works correctly.\n" +
		"- References: docs/aics/x/aic.md\n" +
		"- Status: TODO.\n"

	p := Parse([]byte(src))

	has := func(fs []Finding, rule string) bool {
		for _, f := range fs {
			if f.Rule == rule {
				return true
			}
		}
		return false
	}

	if has(p.Validate(ValidateOpts{}), "unverifiable-acceptance") {
		t.Error("rule fired without --strict; it must not break existing plans")
	}
	strict := p.Validate(ValidateOpts{Strict: true})
	if !has(strict, "unverifiable-acceptance") {
		t.Fatalf("rule did not fire under --strict; findings: %+v", strict)
	}
	for _, f := range strict {
		if f.Rule == "unverifiable-acceptance" && f.Severity != SevWarning {
			t.Errorf("want severity %q, got %q", SevWarning, f.Severity)
		}
	}
}

// TestRunnableAcceptancePassesStrict is the other half: a good criterion must
// not be flagged.
func TestRunnableAcceptancePassesStrict(t *testing.T) {
	src := "### T-1: Do the thing\n" +
		"- WHAT: Do the thing.\n" +
		"- WHERE: internal/thing/thing.go\n" +
		"- WHY: The thing is not done.\n" +
		"- Acceptance:\n" +
		"  - GIVEN the change WHEN `go test ./internal/thing/...` runs THEN it passes.\n" +
		"- References: docs/aics/x/aic.md\n" +
		"- Status: TODO.\n"

	p := Parse([]byte(src))
	for _, f := range p.Validate(ValidateOpts{Strict: true}) {
		if f.Rule == "unverifiable-acceptance" {
			t.Errorf("flagged a runnable criterion: %s", f.Message)
		}
	}
}
