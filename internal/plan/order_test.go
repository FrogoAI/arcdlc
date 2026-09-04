package plan

import (
	"errors"
	"strings"
	"testing"
)

// blockLines renders one minimal, valid task block plus the blank line that
// separates it from the next one.
func blockLines(id string) []string {
	return []string{
		"### " + id + ": Task " + id,
		"- WHAT: x.",
		"- WHERE: internal/" + id + ".go",
		"- WHY: y.",
		"- References: `a`.",
		"- Status: TODO.",
		"",
	}
}

// planOf builds a plan with no preamble whose first byte starts a "###" block.
func planOf(ids ...string) []byte {
	var ls []string
	for _, id := range ids {
		ls = append(ls, blockLines(id)...)
	}
	return lines(ls[:len(ls)-1]...) // drop the trailing separator blank line
}

func crlf(b []byte) []byte {
	return []byte(strings.ReplaceAll(string(b), "\n", "\r\n"))
}

func assertOrder(t *testing.T, out []byte, want string) {
	t.Helper()
	if got := strings.Join(ids2(Parse(out).Tasks), " "); got != want {
		t.Fatalf("order = %q, want %q\n---\n%s", got, want, out)
	}
}

func reorderErr(t *testing.T, err error) *ReorderError {
	t.Helper()
	var re *ReorderError
	if !errors.As(err, &re) {
		t.Fatalf("error = %v (%T), want *ReorderError", err, err)
	}
	return re
}

func TestReorderSlotPermutation(t *testing.T) {
	p := Parse(planOf("T1", "T2", "T3", "T4", "T5", "T6"))
	out, changed, err := p.Reorder([]string{"T5", "T2", "T3"})
	if err != nil {
		t.Fatalf("Reorder: %v", err)
	}
	if !changed {
		t.Fatal("changed = false, want true")
	}
	assertOrder(t, out, "T1 T5 T2 T4 T3 T6")
}

func TestReorderFullPermutationAndNoOp(t *testing.T) {
	p := Parse(planOf("T1", "T2", "T3"))

	out, changed, err := p.Reorder([]string{"T3", "T1", "T2"})
	if err != nil || !changed {
		t.Fatalf("Reorder = (changed %v, err %v), want (true, nil)", changed, err)
	}
	assertOrder(t, out, "T3 T1 T2")

	// Naming a subset that already sits in the requested order is a no-op:
	// T1 already holds the first named slot.
	out, changed, err = p.Reorder([]string{"T1", "T3"})
	if err != nil {
		t.Fatalf("Reorder: %v", err)
	}
	if changed {
		t.Error("changed = true, want false for an order the plan already has")
	}
	if out != nil {
		t.Errorf("out = %q, want nil", out)
	}
}

func TestReorderBadArgs(t *testing.T) {
	p := Parse(planOf("T1", "T2", "T3"))
	cases := []struct {
		name    string
		ids     []string
		wantIDs string
	}{
		{"no ids", nil, ""},
		{"one id", []string{"T1"}, ""},
		{"repeated id", []string{"T1", "T2", "T1"}, "T1"},
		{"id given three times", []string{"T1", "T2", "T1", "T1"}, "T1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, changed, err := p.Reorder(c.ids)
			re := reorderErr(t, err)
			if re.Kind != ReorderBadArgs {
				t.Errorf("Kind = %d, want ReorderBadArgs", re.Kind)
			}
			if got := strings.Join(re.IDs, ", "); got != c.wantIDs {
				t.Errorf("IDs = %q, want %q", got, c.wantIDs)
			}
			if out != nil || changed {
				t.Errorf("out = %q, changed = %v; want nil, false", out, changed)
			}
		})
	}
}

func TestReorderNotFound(t *testing.T) {
	dup := Parse(planOf("T1", "T2", "T2", "T3"))
	empty := Parse(lines("Task format: see plan-format.md.", "", "No blocks yet."))
	cases := []struct {
		name    string
		plan    *Plan
		ids     []string
		wantMsg string
		wantIDs string
	}{
		{"unknown id", Parse(planOf("T1", "T2")), []string{"T1", "T9"}, "no such task", "T9"},
		{"duplicate in the plan", dup, []string{"T2", "T1"}, "ambiguous task ID (duplicate in the plan)", "T2"},
		{"unknown wins over ambiguous", dup, []string{"T2", "T9"}, "no such task", "T9"},
		{"plan has no blocks", empty, []string{"T1", "T2"}, "plan has no task blocks", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, changed, err := c.plan.Reorder(c.ids)
			re := reorderErr(t, err)
			if re.Kind != ReorderNotFound {
				t.Errorf("Kind = %d, want ReorderNotFound", re.Kind)
			}
			if re.Msg != c.wantMsg {
				t.Errorf("Msg = %q, want %q", re.Msg, c.wantMsg)
			}
			if got := strings.Join(re.IDs, ", "); got != c.wantIDs {
				t.Errorf("IDs = %q, want %q", got, c.wantIDs)
			}
			if out != nil || changed {
				t.Errorf("out = %q, changed = %v; want nil, false", out, changed)
			}
		})
	}
}

func TestReorderReportsEveryUnknownID(t *testing.T) {
	p := Parse(planOf("T1", "T2", "T3"))
	_, _, err := p.Reorder([]string{"T7", "T1", "T8", "T9"})
	re := reorderErr(t, err)
	if got := strings.Join(re.IDs, ", "); got != "T7, T8, T9" {
		t.Errorf("IDs = %q, want all three unknown IDs", got)
	}
	msg := re.Error()
	if msg != "no such task: T7, T8, T9" {
		t.Errorf("Error() = %q, want it to name all three", msg)
	}
}

func TestReorderPreservesPreambleAndLedger(t *testing.T) {
	ls := []string{
		"Task format: see plan-format.md.",
		"",
		"Completed (archived to docs/aics/ordering/plan-archive.md):",
		"- T0: Bootstrap",
		"",
	}
	for _, id := range []string{"T1", "T2", "T3"} {
		ls = append(ls, blockLines(id)...)
	}
	src := lines(ls[:len(ls)-1]...)

	p := Parse(src)
	out, changed, err := p.Reorder([]string{"T3", "T1"})
	if err != nil || !changed {
		t.Fatalf("Reorder = (changed %v, err %v), want (true, nil)", changed, err)
	}
	assertOrder(t, out, "T3 T2 T1")

	want := string(src[:p.Tasks[0].RawStart])
	np := Parse(out)
	if got := string(out[:np.Tasks[0].RawStart]); got != want {
		t.Errorf("preamble changed:\ngot  %q\nwant %q", got, want)
	}
	if n := strings.Count(string(out), "Completed (archived"); n != 1 {
		t.Errorf("ledger header count = %d, want 1 (order must not regenerate it)", n)
	}
}

// A stray line after a block's keys belongs to that block, so it travels with
// the block when the block moves. Pinned deliberately: it is the behaviour that
// makes a hand-annotated plan survive a reorder.
func TestReorderCarriesTrailingNoteWithItsBlock(t *testing.T) {
	ls := append(blockLines("T1"), blockLines("T2")...)
	ls = append(ls, "Note: T2 needs the staging key.", "")
	ls = append(ls, blockLines("T3")...)
	src := lines(ls[:len(ls)-1]...)

	out, changed, err := Parse(src).Reorder([]string{"T3", "T2"})
	if err != nil || !changed {
		t.Fatalf("Reorder = (changed %v, err %v), want (true, nil)", changed, err)
	}
	assertOrder(t, out, "T1 T3 T2")

	body := string(out)
	head := strings.Index(body, "### T2:")
	note := strings.Index(body, "Note: T2 needs the staging key.")
	if head < 0 || note < 0 {
		t.Fatalf("block or note missing:\n%s", body)
	}
	if note < head {
		t.Fatalf("note did not move with T2:\n%s", body)
	}
	if strings.Contains(body[head+len("### T2:"):note], "### ") {
		t.Fatalf("note landed under another block:\n%s", body)
	}
}

func TestReorderCRLF(t *testing.T) {
	src := crlf(planOf("T1", "T2", "T3"))
	p := Parse(src)
	if !p.CRLF {
		t.Fatal("fixture is not CRLF")
	}
	out, changed, err := p.Reorder([]string{"T3", "T1"})
	if err != nil || !changed {
		t.Fatalf("Reorder = (changed %v, err %v), want (true, nil)", changed, err)
	}
	assertOrder(t, out, "T3 T2 T1")

	body := string(out)
	if !strings.Contains(body, "### T3: Task T3\r\n- WHAT: x.\r\n") {
		t.Errorf("block lost its CRLF bytes:\n%q", body)
	}
	// Blocks are joined with LF even in a CRLF plan, exactly as Archive does:
	// a block's own trailing bytes are kept and one LF separates the blocks.
	if !strings.Contains(body, "- Status: TODO.\r\n\n### T2:") {
		t.Errorf("blocks are not separated by LF:\n%q", body)
	}
	if !strings.Contains(body, "- Status: TODO.\r\n\r\n\n### T1:") {
		t.Errorf("a block's own trailing blank line was not kept:\n%q", body)
	}
	if !strings.HasSuffix(body, "- Status: TODO.\r\n\r\n") {
		t.Errorf("last block did not keep its CRLF tail:\n%q", body)
	}
}

func TestReorderNoLeadingBlankLine(t *testing.T) {
	out, changed, err := Parse(planOf("T1", "T2", "T3")).Reorder([]string{"T3", "T1"})
	if err != nil || !changed {
		t.Fatalf("Reorder = (changed %v, err %v), want (true, nil)", changed, err)
	}
	if !strings.HasPrefix(string(out), "### T3:") {
		t.Errorf("output gained a leading blank line:\n%q", string(out[:20]))
	}
	if !strings.HasSuffix(string(out), "- Status: TODO.\n") || strings.HasSuffix(string(out), "\n\n") {
		t.Errorf("output does not end in exactly one newline:\n%q", out)
	}
}

func TestVerifyReorderRejectsCorruptedOutput(t *testing.T) {
	p := Parse(planOf("T1", "T2", "T3"))
	out, _, err := p.Reorder([]string{"T3", "T1"})
	if err != nil {
		t.Fatalf("Reorder: %v", err)
	}
	if err := verifyReorder(p, out, []string{"T3", "T2", "T1"}); err != nil {
		t.Fatalf("valid output rejected: %v", err)
	}
	if err := verifyReorder(p, out, []string{"T1", "T2", "T3"}); err == nil {
		t.Error("wrong ID sequence accepted")
	}
	short := out[:strings.Index(string(out), "### T1:")]
	if err := verifyReorder(p, short, []string{"T3", "T2"}); err == nil {
		t.Error("dropped block accepted")
	}
	edited := []byte(strings.Replace(string(out), "- WHY: y.", "- WHY: z.", 1))
	if err := verifyReorder(p, edited, []string{"T3", "T2", "T1"}); err == nil {
		t.Error("edited block body accepted")
	}
}
