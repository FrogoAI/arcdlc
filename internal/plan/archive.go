package plan

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// ArchiveHeader is written to a fresh archive file before the first section.
const ArchiveHeader = "# Plan Archive\n\n" +
	"Archived task blocks from [plan.md](plan.md). Managed by `arctool archive`; do not hand-edit.\n"

var (
	reLedgerHead = regexp.MustCompile(`^Completed \(archived to .+\):\s*$`)
	reLedgerItem = regexp.MustCompile(`^-\s+\S+:\s+.+$`)
)

// ArchiveResult is the computed (but not yet written) outcome of an archive run.
type ArchiveResult struct {
	Archived []*Task // DONE tasks moved out, in document order
	NewPlan  []byte  // rewritten plan.md: DONE blocks removed, ledger extended
	Section  []byte  // "## Archived <date>\n\n<blocks>" to append to the archive file
}

// Archive computes the plan rewrite and the archive section for the current
// DONE blocks. archiveRel is the path shown in the ledger line (e.g.
// "docs/aics/<slug>/plan-archive.md"); date is the section date (YYYY-MM-DD). It
// returns (nil, false) when there is nothing to archive. Kept (non-DONE) blocks
// keep their exact content; only inter-block spacing is normalized to one blank
// line. Nothing is written — the caller writes the archive first, then the plan.
func (p *Plan) Archive(archiveRel, date string) (*ArchiveResult, bool) {
	var done, keep []*Task
	for i := range p.Tasks {
		if p.Tasks[i].Status == StatusDONE {
			done = append(done, &p.Tasks[i])
		} else {
			keep = append(keep, &p.Tasks[i])
		}
	}
	if len(done) == 0 {
		return nil, false
	}

	firstStart := len(p.Bytes)
	if len(p.Tasks) > 0 {
		firstStart = p.Tasks[0].RawStart
	}
	intro, existing := splitPreamble(p.Bytes[:firstStart])

	entries := append([]string{}, existing...)
	for _, t := range done {
		entries = append(entries, fmt.Sprintf("- %s: %s", t.ID, t.Title))
	}
	ledger := "Completed (archived to " + archiveRel + "):\n" + strings.Join(entries, "\n") + "\n"

	var plan bytes.Buffer
	if s := strings.TrimRight(intro, "\n"); s != "" {
		plan.WriteString(s)
		plan.WriteString("\n\n")
	}
	plan.WriteString(ledger)
	for _, t := range keep {
		plan.WriteString("\n")
		plan.Write(bytes.TrimRight(p.Raw(t), "\n"))
		plan.WriteString("\n")
	}

	var sec bytes.Buffer
	sec.WriteString("## Archived " + date + "\n")
	for _, t := range done {
		sec.WriteString("\n")
		sec.Write(bytes.TrimRight(p.Raw(t), "\n"))
		sec.WriteString("\n")
	}

	return &ArchiveResult{Archived: done, NewPlan: plan.Bytes(), Section: sec.Bytes()}, true
}

// VerifyArchive re-parses the rewritten plan and checks the safety invariants
// before anything is written: no DONE block remains, the pending TODO/TAKEN/
// BLOCKED counts are unchanged, and every archived ID appears exactly once in
// the new ledger and once in the archive section. Returns nil when safe.
func VerifyArchive(orig *Plan, res *ArchiveResult) error {
	np := Parse(res.NewPlan)
	for i := range np.Tasks {
		if np.Tasks[i].Status == StatusDONE {
			return fmt.Errorf("DONE block %q still present in plan after archive", np.Tasks[i].ID)
		}
	}
	oc, nc := orig.StatusCounts(), np.StatusCounts()
	for _, s := range []Status{StatusTODO, StatusTAKEN, StatusBLOCKED} {
		if oc[s] != nc[s] {
			return fmt.Errorf("pending %s count changed: %d -> %d", s, oc[s], nc[s])
		}
	}

	prEnd := len(res.NewPlan)
	if len(np.Tasks) > 0 {
		prEnd = np.Tasks[0].RawStart
	}
	_, ledger := splitPreamble(res.NewPlan[:prEnd])
	sec := Parse(res.Section)
	for _, t := range res.Archived {
		if n := countLedgerID(ledger, t.ID); n != 1 {
			return fmt.Errorf("archived %q appears %d time(s) in the ledger (want 1)", t.ID, n)
		}
		if n := len(sec.ByID(t.ID)); n != 1 {
			return fmt.Errorf("archived %q appears %d time(s) in the archive section (want 1)", t.ID, n)
		}
	}
	return nil
}

// splitPreamble separates plan-preamble bytes into intro text and any existing
// ledger entries (the "- id: title" bullets under a "Completed (...)" header).
func splitPreamble(pre []byte) (intro string, ledger []string) {
	lines := strings.Split(string(pre), "\n")
	var introLines []string
	for i := 0; i < len(lines); i++ {
		if reLedgerHead.MatchString(strings.TrimSpace(lines[i])) {
			for i+1 < len(lines) && reLedgerItem.MatchString(strings.TrimSpace(lines[i+1])) {
				i++
				ledger = append(ledger, strings.TrimSpace(lines[i]))
			}
			continue // drop the header line too; it is regenerated
		}
		introLines = append(introLines, lines[i])
	}
	return strings.TrimRight(strings.Join(introLines, "\n"), "\n"), ledger
}

func countLedgerID(items []string, id string) int {
	n := 0
	for _, it := range items {
		if ledgerID(it) == id {
			n++
		}
	}
	return n
}

func ledgerID(item string) string {
	s := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(item), "-"))
	if i := strings.Index(s, ":"); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}
