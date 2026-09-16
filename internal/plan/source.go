package plan

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// A source stamp records which architecture document a plan was decomposed
// from, and what that document said at the time.
//
// It exists because a design matures by re-running /arcdlc:aic, and a plan
// written from round two keeps executing happily after round five moved the
// design underneath it. Execution is mechanical by contract (ADR-0021), so the
// executor has no context to notice: it builds the superseded design faithfully
// and reports success. The stamp turns that silent drift into a failed check.
//
// Format, in an HTML comment so it never renders and never parses as a task:
//
//	<!-- arcdlc:source docs/aics/checkout/aic.md sha256:3f786850e387550f… -->
const sourceMarkerPrefix = "<!-- arcdlc:source "

var reSourceStamp = regexp.MustCompile(`^<!--\s*arcdlc:source\s+(\S+)\s+sha256:([0-9a-f]{64})\s*-->$`)

// Source is one stamped architecture document.
type Source struct {
	Path   string // as written in the plan, relative to the repository root
	Sum    string // lowercase hex sha256 recorded when the plan was written
	Line   int    // 1-based line in the plan
	Offset int    // byte offset of the line start, for rewriting
	Length int    // byte length of the line, excluding its newline
}

// Sources returns every source stamp in the plan, in document order.
func (p *Plan) Sources() []Source {
	var out []Source
	for i, ln := range splitLines(p.Bytes) {
		t := strings.TrimSpace(ln.text)
		if !strings.HasPrefix(t, sourceMarkerPrefix) {
			continue
		}
		m := reSourceStamp.FindStringSubmatch(t)
		if m == nil {
			continue
		}
		out = append(out, Source{
			Path:   m[1],
			Sum:    m[2],
			Line:   i + 1,
			Offset: ln.start,
			Length: len(ln.text),
		})
	}
	return out
}

// SumFile returns the lowercase hex sha256 of a file's bytes.
func SumFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

// CheckSources compares every stamp against the file it names, resolving paths
// against root. It does the file reading that Validate deliberately does not,
// so the validator stays a pure function of the plan text.
//
// A missing file is an error: the plan points at a document that is gone. A
// changed file is a warning, which --strict turns into a failure, because the
// plan may still be correct and only the engineer can say.
func CheckSources(p *Plan, root string) []Finding {
	f := []Finding{}
	for _, s := range p.Sources() {
		got, err := SumFile(filepath.Join(root, s.Path))
		if err != nil {
			f = append(f, Finding{SevError, "", "source-missing",
				fmt.Sprintf("plan was decomposed from %s, which cannot be read: %v", s.Path, err), s.Line})
			continue
		}
		if got != s.Sum {
			f = append(f, Finding{SevWarning, "", "source-changed",
				fmt.Sprintf("%s changed since this plan was written; re-run /arcdlc:plan, or re-stamp with `arctool stamp` if the change does not affect the tasks", s.Path), s.Line})
		}
	}
	return f
}

// Stamp rewrites every existing stamp to the current hash of the file it names,
// and returns the new plan bytes plus the paths it refreshed. Stamps whose file
// cannot be read are left alone and reported as an error, so a typo is visible
// rather than silently dropped.
//
// It only ever updates a line that is already there. Adding the first stamp for
// a document is AddSource's job, because deciding which document a plan derives
// from is a judgement, not a rewrite.
func Stamp(p *Plan, root string) (out []byte, refreshed []string, err error) {
	srcs := p.Sources()
	out = append([]byte(nil), p.Bytes...)

	// Rewrite back to front so earlier offsets stay valid.
	for i := len(srcs) - 1; i >= 0; i-- {
		s := srcs[i]
		got, e := SumFile(filepath.Join(root, s.Path))
		if e != nil {
			return nil, nil, fmt.Errorf("stamp %s: %w", s.Path, e)
		}
		if got == s.Sum {
			continue
		}
		line := fmt.Sprintf("<!-- arcdlc:source %s sha256:%s -->", s.Path, got)
		out = append(out[:s.Offset], append([]byte(line), out[s.Offset+s.Length:]...)...)
		refreshed = append(refreshed, s.Path)
	}
	// Reverse so the report reads in document order.
	for i, j := 0, len(refreshed)-1; i < j; i, j = i+1, j-1 {
		refreshed[i], refreshed[j] = refreshed[j], refreshed[i]
	}
	return out, refreshed, nil
}

// AddSource inserts a stamp for path if the plan does not already carry one,
// placing it immediately after the plan's H1 (or at the top when there is none)
// so it sits in the preamble and never inside a task block.
func AddSource(p *Plan, root, path string) ([]byte, bool, error) {
	for _, s := range p.Sources() {
		if s.Path == path {
			return p.Bytes, false, nil
		}
	}
	sum, err := SumFile(filepath.Join(root, path))
	if err != nil {
		return nil, false, fmt.Errorf("stamp %s: %w", path, err)
	}
	nl := "\n"
	if p.CRLF {
		nl = "\r\n"
	}
	line := fmt.Sprintf("<!-- arcdlc:source %s sha256:%s -->", path, sum)

	at := 0
	for _, ln := range splitLines(p.Bytes) {
		if strings.HasPrefix(ln.text, "# ") {
			at = ln.start + len(ln.text) + len(nl)
			break
		}
	}
	if at > len(p.Bytes) {
		at = len(p.Bytes)
	}
	out := append([]byte(nil), p.Bytes[:at]...)
	if at == 0 {
		out = append(out, []byte(line+nl+nl)...)
	} else {
		out = append(out, []byte(nl+line+nl)...)
	}
	out = append(out, p.Bytes[at:]...)
	return out, true, nil
}
