package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/FrogoAI/arcdlc/internal/plan"
)

// statusFmt is the column layout shared by the header and every row of the
// text form of `arctool status`.
const statusFmt = "%-28s %-14s %5s %5s %7s %5s  %-10s %s\n"

// statusCounts tallies one initiative's tasks by status. plan-archive.md
// contributes only to DONE: an archived task is still done, and it carries no
// other status.
type statusCounts struct {
	TODO    int `json:"TODO"`
	TAKEN   int `json:"TAKEN"`
	BLOCKED int `json:"BLOCKED"`
	DONE    int `json:"DONE"`
}

// statusRow is one line of `arctool status`: an initiative folder, its phase
// derived from the files it holds (AIC H8), and, once closed, the date and
// outcome read from CLOSED.md.
type statusRow struct {
	Slug    string       `json:"slug"`
	Phase   string       `json:"phase"`
	Counts  statusCounts `json:"counts"`
	Closed  string       `json:"closed"`
	Outcome string       `json:"outcome"`
}

// cmdStatus is the read-only companion to list: it shows every initiative
// folder under docs/aics/, not just one selected plan, so it takes no --aic.
func cmdStatus(args []string) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "emit as JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "usage: arctool status [--json]")
		return 2
	}
	return runStatus(aicsDir, *asJSON, os.Stdout, os.Stderr)
}

// runStatus lists every initiative folder under dir with its phase, task
// counts, and, for a closed one, its close date and outcome. It only reads,
// never writes. A missing dir is not an error: it is read the same as one
// with no folders, so the header (text) or an empty list (JSON) is printed
// and it still exits 0.
func runStatus(dir string, asJSON bool, out, errw io.Writer) int {
	entries, _ := os.ReadDir(dir)

	rows := []statusRow{}
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		row, code := buildStatusRow(dir, e.Name(), errw)
		if code != 0 {
			return code
		}
		rows = append(rows, row)
	}

	if asJSON {
		emitStatusJSON(rows, out)
		return 0
	}
	emitStatusText(rows, out)
	return 0
}

// buildStatusRow computes one initiative's row from the files in
// dir/slug/: plan.md and plan-archive.md for the counts, CLOSED.md for the
// phase, close date and outcome.
func buildStatusRow(dir, slug string, errw io.Writer) (statusRow, int) {
	folder := filepath.Join(dir, slug)
	row := statusRow{Slug: slug}

	hasPlan := fileExists(filepath.Join(folder, "plan.md"))
	if hasPlan {
		c, code := readStatusCounts(filepath.Join(folder, "plan.md"), errw)
		if code != 0 {
			return statusRow{}, code
		}
		row.Counts = c
	}
	if fileExists(filepath.Join(folder, "plan-archive.md")) {
		c, code := readStatusCounts(filepath.Join(folder, "plan-archive.md"), errw)
		if code != 0 {
			return statusRow{}, code
		}
		row.Counts.DONE += c.DONE
	}

	closedPath := filepath.Join(folder, "CLOSED.md")
	switch {
	case fileExists(closedPath):
		row.Phase = "closed"
		row.Closed, row.Outcome = readClosedNote(closedPath)
	case !hasPlan:
		row.Phase = "designing"
	case row.Counts.TODO+row.Counts.TAKEN+row.Counts.BLOCKED > 0:
		row.Phase = "in progress"
	default:
		row.Phase = "ready to close"
	}
	return row, 0
}

// readStatusCounts parses one plan file (plan.md or plan-archive.md) and
// tallies its tasks by status. A read error is reported to errw and yields
// exit code 4.
func readStatusCounts(path string, errw io.Writer) (statusCounts, int) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(errw, "arctool: cannot read %s: %v\n", path, err)
		return statusCounts{}, 4
	}
	c := plan.Parse(b).StatusCounts()
	return statusCounts{
		TODO:    c[plan.StatusTODO],
		TAKEN:   c[plan.StatusTAKEN],
		BLOCKED: c[plan.StatusBLOCKED],
		DONE:    c[plan.StatusDONE],
	}, 0
}

// readClosedNote reads CLOSED.md line by line for its "- Closed:" and
// "- Outcome:" lines. The first match of each wins; a missing line leaves
// that field empty.
func readClosedNote(path string) (closed, outcome string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", ""
	}
	for _, ln := range strings.Split(string(b), "\n") {
		ln = strings.TrimSpace(ln)
		switch {
		case closed == "" && strings.HasPrefix(ln, "- Closed:"):
			closed = strings.TrimSpace(strings.TrimPrefix(ln, "- Closed:"))
		case outcome == "" && strings.HasPrefix(ln, "- Outcome:"):
			outcome = strings.TrimSpace(strings.TrimPrefix(ln, "- Outcome:"))
		}
	}
	return closed, outcome
}

// fileExists reports whether path exists and is a regular file (not a
// directory).
func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// emitStatusJSON writes {"initiatives": [...]}, two-space indented; rows is
// never nil, so an empty result encodes as [], never null.
func emitStatusJSON(rows []statusRow, out io.Writer) {
	payload := struct {
		Initiatives []statusRow `json:"initiatives"`
	}{Initiatives: rows}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		// statusRow holds only strings, ints and a slice: this cannot fail.
		return
	}
	out.Write(b)
	out.Write([]byte("\n"))
}

// emitStatusText writes the header and one line per row, both in the same
// fixed column layout. An empty Closed or Outcome prints as "-".
func emitStatusText(rows []statusRow, out io.Writer) {
	fmt.Fprintf(out, statusFmt, "SLUG", "PHASE", "TODO", "TAKEN", "BLOCKED", "DONE", "CLOSED", "OUTCOME")
	for _, r := range rows {
		closed := orDash(r.Closed)
		outcome := orDash(r.Outcome)
		fmt.Fprintf(out, statusFmt,
			r.Slug, r.Phase,
			strconv.Itoa(r.Counts.TODO), strconv.Itoa(r.Counts.TAKEN),
			strconv.Itoa(r.Counts.BLOCKED), strconv.Itoa(r.Counts.DONE),
			closed, outcome)
	}
}

// orDash returns s, or "-" when it is empty.
func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
