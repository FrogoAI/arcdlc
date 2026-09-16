package plan

import "testing"

// TestCustomKeysArePreserved pins the contract the engineer set: the seven
// defined keys are the format, a custom key is allowed, and the whole task
// reaches the executor. Before this, an unknown key was swallowed into the
// multi-line section above it (corrupting that section too) or vanished
// entirely, and validate stayed silent either way.
func TestCustomKeysArePreserved(t *testing.T) {
	src := "### D-1: Do the thing\n" +
		"- WHAT: Do the thing.\n" +
		"- HOW: Use the existing helper.\n" +
		"- BUDGET: 3 retries maximum, then dead-letter.\n" +
		"- WHERE: internal/api/thing.go\n" +
		"- WHY: It is not done.\n" +
		"- Acceptance:\n" +
		"  - GIVEN it WHEN `go test ./...` runs THEN it passes.\n" +
		"- Rollout plan:\n" +
		"  - Behind a flag for one week.\n" +
		"  - Then default on.\n" +
		"- References: docs/aics/demo/aic.md\n" +
		"- Status: TODO.\n"

	p := Parse([]byte(src))
	if len(p.Tasks) != 1 {
		t.Fatalf("want 1 task, got %d", len(p.Tasks))
	}
	task := p.Tasks[0]

	// The defined key must not absorb the custom one that follows it.
	if task.How != "Use the existing helper." {
		t.Errorf("HOW was corrupted by the following custom key:\n  got %q", task.How)
	}
	// Every defined key still parses.
	for name, got := range map[string]string{
		"WHAT": task.What, "WHERE": task.Where, "WHY": task.Why, "Acceptance": task.Acceptance,
	} {
		if got == "" {
			t.Errorf("%s is empty; a custom key broke normal parsing", name)
		}
	}
	if task.Status != StatusTODO {
		t.Errorf("Status = %q", task.Status)
	}

	if len(task.Extra) != 2 {
		t.Fatalf("want 2 custom keys, got %d: %+v", len(task.Extra), task.Extra)
	}
	if task.Extra[0].Key != "BUDGET" || task.Extra[0].Value != "3 retries maximum, then dead-letter." {
		t.Errorf("first custom key = %+v", task.Extra[0])
	}
	// A custom key carries its indented body, like a defined multi-line key.
	if task.Extra[1].Key != "Rollout plan" {
		t.Errorf("second custom key = %q", task.Extra[1].Key)
	}
	if want := "- Behind a flag for one week.\n  - Then default on."; task.Extra[1].Value != want {
		t.Errorf("custom multi-line body:\n  got  %q\n  want %q", task.Extra[1].Value, want)
	}
}

// TestCustomKeysDoNotFailValidation pins that a custom key is allowed, not
// merely tolerated: --strict must stay silent about it.
func TestCustomKeysDoNotFailValidation(t *testing.T) {
	src := "### D-1: Do the thing\n" +
		"- WHAT: Do the thing.\n" +
		"- BUDGET: 3 retries maximum.\n" +
		"- WHERE: internal/api/thing.go\n" +
		"- WHY: It is not done.\n" +
		"- Acceptance:\n" +
		"  - GIVEN it WHEN `go test ./...` runs THEN it passes.\n" +
		"- References: docs/aics/demo/aic.md\n" +
		"- Status: TODO.\n"

	for _, f := range Parse([]byte(src)).Validate(ValidateOpts{Strict: true, RequireAcceptance: true}) {
		t.Errorf("custom key produced a finding: %s: %s", f.Rule, f.Message)
	}
}

// TestCustomKeyWithNoPrecedingMultiLineKey covers the other half of the old
// bug: with nothing above it to be swallowed into, the key used to vanish.
func TestCustomKeyWithNoPrecedingMultiLineKey(t *testing.T) {
	src := "### D-1: Do the thing\n" +
		"- WHAT: Do the thing.\n" +
		"- BUDGET: 3 retries maximum.\n" +
		"- Status: TODO.\n"

	task := Parse([]byte(src)).Tasks[0]
	if len(task.Extra) != 1 || task.Extra[0].Key != "BUDGET" {
		t.Fatalf("custom key lost: %+v", task.Extra)
	}
}

// TestStatusMutationPreservesCustomKeys pins that the byte-preserving status
// rewrite does not disturb a custom key sitting next to it.
func TestStatusMutationPreservesCustomKeys(t *testing.T) {
	src := "### D-1: Do the thing\n" +
		"- WHAT: Do the thing.\n" +
		"- BUDGET: 3 retries maximum.\n" +
		"- Status: TODO.\n"

	p := Parse([]byte(src))
	out := p.WithStatus(&p.Tasks[0], StatusTAKEN, "")
	task := Parse(out).Tasks[0]
	if len(task.Extra) != 1 || task.Extra[0].Value != "3 retries maximum." {
		t.Errorf("custom key damaged by a status flip: %+v", task.Extra)
	}
	if task.Status != StatusTAKEN {
		t.Errorf("status = %q", task.Status)
	}
}
