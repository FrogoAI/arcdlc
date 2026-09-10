package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Edit is one file rewrite that takes out the marker comments the register now
// holds. Bytes is the whole new content; the caller writes it atomically.
type Edit struct {
	File    string `json:"file"`    // path relative to the sweep root, slash separated
	Removed int    `json:"removed"` // markers taken out of this file
	Bytes   []byte `json:"-"`

	deleted, cut int // what VerifyEdit must find, and nothing besides
}

// Skip is a marker Strip left in the code, with the reason it stays.
type Skip struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Marker string `json:"marker"`
	Reason string `json:"reason"`
}

// Strip removes the marker comments in found from the files they came from. It
// deletes only the lines (or the one line span) the sweep recorded for each
// finding, verifies each result, and returns the new content per file. Nothing
// is written here: the caller writes the register first, so the marker text is
// never the only copy in flight.
//
// A marker Strip cannot remove safely is returned in skipped and left in the
// code, so the next sweep finds it again.
func Strip(root string, found []Finding) (edits []Edit, skipped []Skip, err error) {
	if root == "" {
		root = "."
	}
	order, byFile := fileOrder(found)
	for _, rel := range order {
		orig, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if readErr != nil {
			return nil, nil, fmt.Errorf("read %s: %w", rel, readErr)
		}
		lines := splitKeepEnds(orig)

		del := map[int]bool{}
		cuts := map[int][2]int{}
		removed := 0
		for _, f := range byFile[rel] {
			if !f.strip.ok {
				skipped = append(skipped, Skip{File: rel, Line: f.Line, Marker: f.Marker, Reason: f.strip.reason})
				continue
			}
			if !stillThere(lines, f) {
				skipped = append(skipped, Skip{File: rel, Line: f.Line, Marker: f.Marker,
					Reason: "the file changed since the sweep"})
				continue
			}
			for n := f.strip.fromLine; n > 0 && n <= f.strip.toLine; n++ {
				del[n] = true
			}
			if f.strip.cutLine > 0 {
				cuts[f.strip.cutLine] = [2]int{f.strip.cutFrom, f.strip.cutTo}
			}
			removed++
		}
		if removed == 0 {
			continue
		}
		collapseBlanks(lines, del)

		out := make([]byte, 0, len(orig))
		for n := 1; n <= len(lines); n++ {
			raw := string(lines[n-1])
			if del[n] {
				continue
			}
			if span, ok := cuts[n]; ok {
				raw = cutSpan(raw, span[0], span[1])
			}
			out = append(out, raw...)
		}
		e := Edit{File: rel, Removed: removed, Bytes: out, deleted: len(del), cut: len(cuts)}
		if verr := VerifyEdit(orig, e); verr != nil {
			return nil, nil, fmt.Errorf("%s: %w", rel, verr)
		}
		edits = append(edits, e)
	}
	return edits, skipped, nil
}

// fileOrder groups findings by file, keeping the order the sweep produced.
func fileOrder(found []Finding) ([]string, map[string][]Finding) {
	var order []string
	byFile := map[string][]Finding{}
	for _, f := range found {
		if _, seen := byFile[f.File]; !seen {
			order = append(order, f.File)
		}
		byFile[f.File] = append(byFile[f.File], f)
	}
	return order, byFile
}

// stillThere reports whether the marker is where the sweep saw it. A file edited
// between the sweep and the strip keeps its comment instead of losing a line of
// code.
func stillThere(lines [][]byte, f Finding) bool {
	n := f.strip.cutLine
	if n == 0 {
		n = f.strip.fromLine
	}
	if n < 1 || n > len(lines) {
		return false
	}
	return strings.Contains(string(lines[n-1]), f.Marker)
}

// collapseBlanks drops one of the two blank lines a whole-line deletion would
// leave behind, so the diff stays clean and gofmt has nothing to say.
func collapseBlanks(lines [][]byte, del map[int]bool) {
	for n := range del {
		if del[n-1] || !del[n] { // only the last line of a deleted run
			continue
		}
		above, below := n-1, n+1
		for del[below] {
			below++
		}
		if isBlankLine(lines, above) && isBlankLine(lines, below) {
			del[below] = true
		}
	}
}

func isBlankLine(lines [][]byte, n int) bool {
	if n < 1 || n > len(lines) {
		return false
	}
	return strings.TrimSpace(string(lines[n-1])) == ""
}

// cutSpan removes [from,to) from a line and trims the whitespace the comment
// left behind, keeping the line terminator.
func cutSpan(raw string, from, to int) string {
	body, term := splitTerm(raw)
	if from < 0 || to > len(body) || from > to {
		return raw
	}
	kept := strings.TrimRight(body[:from], " \t") + body[to:]
	return strings.TrimRight(kept, " \t") + term
}

func splitTerm(raw string) (body, term string) {
	switch {
	case strings.HasSuffix(raw, "\r\n"):
		return raw[:len(raw)-2], "\r\n"
	case strings.HasSuffix(raw, "\n"):
		return raw[:len(raw)-1], "\n"
	}
	return raw, ""
}

// VerifyEdit proves that an edit removed comment lines and nothing else, before
// a byte of it is written. It checks the result independently of how it was
// built: every line of the result must be a line of the original, in order,
// either byte for byte or with a trailing span cut off, and the counts must
// match what the strip plan intended.
func VerifyEdit(orig []byte, e Edit) error {
	o, n := splitKeepEnds(orig), splitKeepEnds(e.Bytes)
	if len(o)-len(n) != e.deleted {
		return fmt.Errorf("result has %d line(s), original has %d, plan deletes %d",
			len(n), len(o), e.deleted)
	}
	i, cuts := 0, 0
	for _, line := range n {
		for i < len(o) && !sameLine(o[i], line) && !isCutOf(o[i], line) {
			i++
		}
		if i >= len(o) {
			return fmt.Errorf("result line %q is not in the original, in order", clipLine(line))
		}
		if !sameLine(o[i], line) {
			cuts++
		}
		i++
	}
	if cuts != e.cut {
		return fmt.Errorf("%d line(s) were shortened, plan shortens %d", cuts, e.cut)
	}
	return nil
}

func sameLine(a, b []byte) bool { return string(a) == string(b) }

// isCutOf reports whether want is orig with exactly one contiguous span taken
// out: the line kept its code and lost its comment, wherever the comment sat.
func isCutOf(orig, want []byte) bool {
	ob, _ := splitTerm(string(orig))
	wb, _ := splitTerm(string(want))
	if len(wb) >= len(ob) || wb == "" {
		return false
	}
	i := 0
	for i < len(wb) && wb[i] == ob[i] {
		i++
	}
	j := 0
	for j < len(wb)-i && wb[len(wb)-1-j] == ob[len(ob)-1-j] {
		j++
	}
	return i+j == len(wb)
}

func clipLine(b []byte) string {
	s, _ := splitTerm(string(b))
	if len(s) > 60 {
		s = s[:60] + "..."
	}
	return s
}
