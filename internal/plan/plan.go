// Package plan parses docs/aics/<slug>/plan.md into a Plan value and back.
//
// It is a pure, side-effect-free, standard-library-only implementation of the
// plan format defined in the arcdlc plan skill's references/plan-format.md.
// The parser records per-block byte offsets so that mutating commands can
// splice single regions while leaving every other byte identical.
package plan

import (
	"regexp"
	"strings"
)

// Status is the normalized value of a task block's "- Status:" line.
type Status uint8

const (
	// StatusNone means the block has no recognizable status token. It is a
	// distinct, first-class state: it is what the runner would silently skip,
	// so it must never be conflated with StatusTODO.
	StatusNone Status = iota
	StatusTODO
	StatusTAKEN
	StatusDONE
	StatusBLOCKED
)

func (s Status) String() string {
	switch s {
	case StatusTODO:
		return "TODO"
	case StatusTAKEN:
		return "TAKEN"
	case StatusDONE:
		return "DONE"
	case StatusBLOCKED:
		return "BLOCKED"
	default:
		return ""
	}
}

// SourceStatus is the optional "(MISSING|PARTIAL|DRIFT)" tag in a heading,
// carried from an examinate gap. Empty when the heading has no parenthetical.
// Unknown values are preserved verbatim so validate can flag them.
type SourceStatus string

const (
	SrcMISSING SourceStatus = "MISSING"
	SrcPARTIAL SourceStatus = "PARTIAL"
	SrcDRIFT   SourceStatus = "DRIFT"
)

// Task is one "### <ID> (<SRC>): <Title>" block and its parsed keys.
type Task struct {
	ID           string
	SourceStatus SourceStatus
	Title        string

	What       string
	How        string // optional: implementation decisions the executor must follow
	Where      string
	Why        string
	References []string
	Acceptance string
	Status     Status
	StatusRaw  string // the status token as written (for casing checks), e.g. "TODO" or "Todo"

	// Presence flags for required-key validation (a key can be present but empty).
	// HasHow is informational only: HOW is never required.
	HasWhat, HasHow, HasWhere, HasWhy, HasRefs, HasStatus, HasAcceptance bool

	HeadingOK bool // heading matched the expected "### <ID>...: <Title>" shape
	Line      int  // 1-based line of the heading (block start), for messages
	LineEnd   int  // 1-based line of the block's last line

	// Byte offsets into Plan.Bytes for byte-preserving edits.
	// Bytes[RawStart:RawEnd] is the exact block text (### line through the byte
	// before the next block/EOF). StatusLine* address just the status line.
	RawStart, RawEnd               int
	StatusLineStart, StatusLineEnd int // -1 when there is no status line
	StatusLine                     int // 1-based line of the status line, 0 if none
}

// Region marks a non-task span (preamble, ledger) preserved as-is by archive.
type Region struct {
	Kind       string // "preamble" | "ledger"
	Start, End int
}

// Plan is a parsed plan file. Bytes is never mutated in place.
type Plan struct {
	Bytes   []byte
	CRLF    bool
	Tasks   []Task
	Regions []Region
}

var (
	reHeading = regexp.MustCompile(`^###\s+(\S+?)\s*(?:\(([^)]*)\))?\s*:\s*(.*?)\s*$`)
	reKey     = regexp.MustCompile(`^-\s+(WHAT|HOW|WHERE|WHY|References|Acceptance|Status):\s?(.*)$`)
	reFence   = regexp.MustCompile("^(```|~~~)")
	rePathExt = regexp.MustCompile(`\.\w{1,6}(\s|$|,|` + "`" + `)`)
)

type srcLine struct {
	start, end int    // byte range; end includes the trailing newline
	text       string // line content without trailing \r or \n
}

func splitLines(b []byte) []srcLine {
	var out []srcLine
	i, n := 0, len(b)
	for i < n {
		j := i
		for j < n && b[j] != '\n' {
			j++
		}
		end := j
		if j < n {
			end = j + 1 // include the newline
		}
		text := strings.TrimSuffix(string(b[i:j]), "\r")
		out = append(out, srcLine{start: i, end: end, text: text})
		i = end
	}
	return out
}

func detectCRLF(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] == '\n' {
			return i > 0 && b[i-1] == '\r'
		}
	}
	return false
}

// Parse turns plan bytes into a Plan. It never errors: malformed content is
// preserved and surfaced by Validate instead.
func Parse(b []byte) *Plan {
	p := &Plan{Bytes: b, CRLF: detectCRLF(b)}
	lines := splitLines(b)

	// First pass: find block-start line indices (### , outside fenced code).
	inFence := false
	var starts []int
	for idx, ln := range lines {
		if reFence.MatchString(strings.TrimSpace(ln.text)) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(ln.text, "### ") {
			starts = append(starts, idx)
		}
	}

	if len(starts) == 0 {
		if len(b) > 0 {
			p.Regions = append(p.Regions, Region{Kind: "preamble", Start: 0, End: len(b)})
		}
		return p
	}

	if lines[starts[0]].start > 0 {
		p.Regions = append(p.Regions, Region{Kind: "preamble", Start: 0, End: lines[starts[0]].start})
	}

	for i, s := range starts {
		endIdx := len(lines)
		if i+1 < len(starts) {
			endIdx = starts[i+1]
		}
		p.Tasks = append(p.Tasks, parseBlock(lines, s, endIdx))
	}
	return p
}

func parseBlock(lines []srcLine, start, end int) Task {
	head := lines[start]
	t := Task{
		Line:            start + 1,
		LineEnd:         end, // 1-based number of the last line (0-based end-1)
		RawStart:        head.start,
		RawEnd:          lines[end-1].end,
		StatusLineStart: -1,
		StatusLineEnd:   -1,
	}
	if m := reHeading.FindStringSubmatch(head.text); m != nil {
		t.HeadingOK = true
		t.ID = m[1]
		t.SourceStatus = SourceStatus(m[2])
		t.Title = strings.TrimSpace(m[3])
	}

	inFence := false
	for i := start + 1; i < end; {
		ln := lines[i]
		if reFence.MatchString(strings.TrimSpace(ln.text)) {
			inFence = !inFence
			i++
			continue
		}
		if inFence {
			i++
			continue
		}
		m := reKey.FindStringSubmatch(ln.text)
		if m == nil {
			i++
			continue
		}
		key, val := m[1], m[2]
		switch key {
		case "WHAT":
			t.What, t.HasWhat = strings.TrimSpace(val), true
			i++
		case "WHY":
			t.Why, t.HasWhy = strings.TrimSpace(val), true
			i++
		case "References":
			t.References, t.HasRefs = parseRefs(val), true
			i++
		case "Status":
			t.Status, t.StatusRaw = normStatus(val)
			t.HasStatus = true
			t.StatusLineStart, t.StatusLineEnd = ln.start, ln.end
			t.StatusLine = i + 1
			i++
		case "HOW", "WHERE", "Acceptance":
			// Multi-line: absorb following lines until the next key / block end.
			body := []string{val}
			j := i + 1
			for j < end {
				nx := lines[j].text
				if reFence.MatchString(strings.TrimSpace(nx)) ||
					reKey.MatchString(nx) || strings.HasPrefix(nx, "### ") {
					break
				}
				body = append(body, nx)
				j++
			}
			v := strings.TrimSpace(strings.Join(body, "\n"))
			switch key {
			case "HOW":
				t.How, t.HasHow = v, true
			case "WHERE":
				t.Where, t.HasWhere = v, true
			default:
				t.Acceptance, t.HasAcceptance = v, true
			}
			i = j
		}
	}
	return t
}

// parseRefs splits a References value into clean paths: comma-separated, with
// surrounding whitespace, a single trailing period, and wrapping backticks removed.
func parseRefs(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		p = strings.TrimSuffix(p, ".")
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "`")
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// normStatus normalizes a status value to an enum, returning the raw token
// (case preserved, trailing period stripped) for casing checks. For BLOCKED
// with a trailing reason, only the first token is considered.
func normStatus(v string) (Status, string) {
	tok := strings.TrimSpace(v)
	if idx := strings.IndexAny(tok, " \t"); idx >= 0 {
		tok = tok[:idx]
	}
	tok = strings.TrimRight(tok, ".")
	switch strings.ToUpper(tok) {
	case "TODO":
		return StatusTODO, tok
	case "TAKEN":
		return StatusTAKEN, tok
	case "DONE":
		return StatusDONE, tok
	case "BLOCKED":
		return StatusBLOCKED, tok
	default:
		return StatusNone, tok
	}
}

// looksLikePath is the loose heuristic for the strict WHERE check: a concrete
// target usually contains a path separator, a backticked token, or a file extension.
func looksLikePath(s string) bool {
	if strings.Contains(s, "/") || strings.Contains(s, "`") {
		return true
	}
	return rePathExt.MatchString(s)
}
