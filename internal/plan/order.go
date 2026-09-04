package plan

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// ReorderErrKind classifies a refused reorder so a caller can map the rejection
// to an exit code without matching on message text.
type ReorderErrKind uint8

const (
	// ReorderBadArgs marks a malformed ID list: fewer than two IDs, or the same
	// ID given more than once. Callers map it to the usage exit code (2).
	ReorderBadArgs ReorderErrKind = iota + 1
	// ReorderNotFound marks an ID that does not resolve to exactly one block, or
	// a plan with no task blocks at all. Callers map it to exit 3.
	ReorderNotFound
	// ReorderSelfCheck marks a failed post-render invariant. The rendered bytes
	// are discarded and must not be written; callers map it to exit 5.
	ReorderSelfCheck
)

// ReorderError is the only error type Reorder returns. Msg is the one-line
// reason; IDs holds every offending task ID — the whole set, never just the
// first, because the caller is usually an agent and one message is one retry.
type ReorderError struct {
	Kind ReorderErrKind
	Msg  string
	IDs  []string
}

func (e *ReorderError) Error() string {
	if len(e.IDs) == 0 {
		return e.Msg
	}
	return e.Msg + ": " + strings.Join(e.IDs, ", ")
}

// Reorder returns the plan bytes with the named tasks permuted among the
// positions they already occupy. It collects the positions of the named blocks,
// sorts those positions ascending, then fills them with the named tasks in the
// order given. Every task nobody named keeps its position, so no unnamed task
// can silently drift past a task it depends on.
//
// Given blocks T1 T2 T3 T4 T5 T6, Reorder([]string{"T5","T2","T3"}) fills
// positions 2, 3 and 5 and yields T1 T5 T2 T4 T3 T6.
//
// changed is false — with nil bytes and no error — when the requested order is
// the order the plan already has. Nothing is written; the caller writes out.
// Every error is a *ReorderError. Rendered bytes are re-parsed and checked
// against the input before they are returned, so a permutation bug cannot reach
// the file: a failed check discards out and yields ReorderSelfCheck.
func (p *Plan) Reorder(ids []string) (out []byte, changed bool, err error) {
	// Validation runs in a fixed order so the message an agent gets is
	// deterministic for an input that trips more than one rule.
	if len(ids) < 2 {
		return nil, false, &ReorderError{Kind: ReorderBadArgs, Msg: "order needs at least two task IDs"}
	}
	seen := make(map[string]int, len(ids))
	var repeated []string
	for _, id := range ids {
		seen[id]++
		if seen[id] == 2 { // report a thrice-given ID once, in first-seen order
			repeated = append(repeated, id)
		}
	}
	if len(repeated) > 0 {
		return nil, false, &ReorderError{Kind: ReorderBadArgs, Msg: "task ID given more than once", IDs: repeated}
	}
	if len(p.Tasks) == 0 {
		return nil, false, &ReorderError{Kind: ReorderNotFound, Msg: "plan has no task blocks"}
	}

	named := make([]int, len(ids)) // index in p.Tasks of each named task
	var unknown, ambiguous []string
	for i, id := range ids {
		switch m := p.ByID(id); len(m) {
		case 0:
			unknown = append(unknown, id)
		case 1:
			named[i] = p.taskIndex(m[0])
		default:
			ambiguous = append(ambiguous, id)
		}
	}
	if len(unknown) > 0 {
		return nil, false, &ReorderError{Kind: ReorderNotFound, Msg: "no such task", IDs: unknown}
	}
	if len(ambiguous) > 0 {
		return nil, false, &ReorderError{
			Kind: ReorderNotFound,
			Msg:  "ambiguous task ID (duplicate in the plan)",
			IDs:  ambiguous,
		}
	}

	// The slots are the named positions, ascending; slot i receives ids[i].
	slots := append([]int(nil), named...)
	sort.Ints(slots)
	newOrder := make([]int, len(p.Tasks))
	for i := range newOrder {
		newOrder[i] = i
	}
	for i, s := range slots {
		newOrder[s] = named[i]
	}
	for i, idx := range newOrder {
		if idx != i {
			changed = true
			break
		}
	}
	if !changed {
		return nil, false, nil
	}

	// Layout follows Archive: the preamble verbatim, then every block separated
	// by exactly one blank line, with a single trailing newline. The ledger is
	// left exactly as found — order never regenerates it.
	var buf bytes.Buffer
	if pre := bytes.TrimRight(p.Bytes[:p.Tasks[0].RawStart], "\n"); len(pre) > 0 {
		buf.Write(pre)
		buf.WriteString("\n")
	}
	want := make([]string, len(newOrder))
	for i, idx := range newOrder {
		if buf.Len() > 0 { // a plan that starts with "###" gains no leading blank line
			buf.WriteString("\n")
		}
		buf.Write(bytes.TrimRight(p.Raw(&p.Tasks[idx]), "\n"))
		buf.WriteString("\n")
		want[i] = p.Tasks[idx].ID
	}
	out = buf.Bytes()

	if err := verifyReorder(p, out, want); err != nil {
		return nil, false, &ReorderError{Kind: ReorderSelfCheck, Msg: "order self-check failed: " + err.Error()}
	}
	return out, true, nil
}

// taskIndex returns the position of t in p.Tasks, or -1 when t is not one of
// p's tasks. Identity, not ID: a plan may hold duplicate IDs.
func (p *Plan) taskIndex(t *Task) int {
	for i := range p.Tasks {
		if &p.Tasks[i] == t {
			return i
		}
	}
	return -1
}

// verifyReorder re-parses rendered order output and checks the safety
// invariants before the bytes are returned: the block count is unchanged, the
// ID sequence is the one the permutation asked for, the per-status counts are
// unchanged, and the multiset of block bodies is unchanged. Returns nil when
// the rewrite is safe. It lives inside Reorder's call path rather than beside
// it (unlike VerifyArchive) because order writes one file, so no caller needs
// the bytes before the check and folding it in makes the check unskippable.
func verifyReorder(orig *Plan, out []byte, want []string) error {
	np := Parse(out)
	if len(np.Tasks) != len(orig.Tasks) {
		return fmt.Errorf("block count changed: %d -> %d", len(orig.Tasks), len(np.Tasks))
	}
	for i := range np.Tasks {
		if np.Tasks[i].ID != want[i] {
			return fmt.Errorf("block %d is %q, want %q", i+1, np.Tasks[i].ID, want[i])
		}
	}
	oc, nc := orig.StatusCounts(), np.StatusCounts()
	for _, s := range []Status{StatusNone, StatusTODO, StatusTAKEN, StatusDONE, StatusBLOCKED} {
		if oc[s] != nc[s] {
			return fmt.Errorf("%s count changed: %d -> %d", s, oc[s], nc[s])
		}
	}
	ob, nb := blockBodies(orig), blockBodies(np)
	for i := range ob {
		if ob[i] != nb[i] {
			return fmt.Errorf("block content changed: %q is not in the original", nb[i])
		}
	}
	return nil
}

// blockBodies returns every block's text with trailing newlines trimmed, sorted
// so two renderings compare as multisets. Content, not ID, is the key: a plan
// may hold duplicate IDs among the blocks nobody named.
func blockBodies(p *Plan) []string {
	out := make([]string, len(p.Tasks))
	for i := range p.Tasks {
		out[i] = string(bytes.TrimRight(p.Raw(&p.Tasks[i]), "\n"))
	}
	sort.Strings(out)
	return out
}
