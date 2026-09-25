package plan

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// AddErrKind classifies a refused append so a caller can map the rejection to
// an exit code without matching on message text.
type AddErrKind uint8

const (
	// AddContract marks a block that fails one of Append's contract checks:
	// wrong shape, a duplicate ID, a status the queue cannot take, a strict
	// validation failure, or a malformed finding. Callers map it to exit 1.
	AddContract AddErrKind = iota + 1
	// AddSelfCheck marks a failed post-render invariant. The rendered bytes
	// are discarded and must not be written; callers map it to exit 5.
	AddSelfCheck
)

// AddError is the only error type Append returns. Msg is the one-line reason;
// Findings carries the strict-validation findings when the block failed that
// check, so a caller can print each one rather than just the summary.
type AddError struct {
	Kind     AddErrKind
	Msg      string
	Findings []Finding
}

func (e *AddError) Error() string {
	if len(e.Findings) == 0 {
		return e.Msg
	}
	msgs := make([]string, len(e.Findings))
	for i, f := range e.Findings {
		msgs[i] = f.Message
	}
	return e.Msg + ": " + strings.Join(msgs, "; ")
}

// reFindingReason matches the reason a BLOCKED finding's status line must
// carry: "found during <TASK-ID>: <reason>". Group 1 is the source task ID.
var reFindingReason = regexp.MustCompile(`^found during (\S+): (.+)$`)

// StatusReason returns the reason carried on t's "- Status:" line (the text
// after "BLOCKED — "), or "" when t has no status line or no reason.
func (p *Plan) StatusReason(t *Task) string {
	if t.StatusLineStart < 0 {
		return ""
	}
	_, _, _, reason := dissectStatusLine(p.Bytes[t.StatusLineStart:t.StatusLineEnd])
	return reason
}

// Append returns the plan bytes with one checked task block added at the end.
// block must hold exactly one well-formed "### <ID>: <Title>" block, TODO or
// BLOCKED, whose ID is not already used by p or by archived, and which passes
// the strict task checks. A BLOCKED block is a finding: its reason must read
// "found during <TASK-ID>: <reason>" naming a task in p or in archived, and
// its ID must be "<TASK-ID>-F<n>".
//
// Every error Append returns is an *AddError. Rendered bytes are re-parsed
// and checked against the input before they are returned, so a rendering bug
// cannot reach the file: a failed check discards out and yields
// AddSelfCheck. p.Bytes is never modified in place.
func (p *Plan) Append(block []byte, archived []string) (out []byte, id string, err error) {
	nb := Parse(block)
	if len(nb.Tasks) != 1 || hasNonBlankPreamble(nb) {
		return nil, "", &AddError{Kind: AddContract, Msg: "input must be exactly one ### task block"}
	}
	t := nb.Tasks[0]
	if !t.HeadingOK || t.ID == "" || t.Title == "" {
		return nil, "", &AddError{Kind: AddContract, Msg: "malformed heading"}
	}
	if len(p.ByID(t.ID)) > 0 || containsID(archived, t.ID) {
		return nil, "", &AddError{Kind: AddContract, Msg: fmt.Sprintf("task ID %q already exists", t.ID)}
	}
	if t.Status != StatusTODO && t.Status != StatusBLOCKED {
		return nil, "", &AddError{Kind: AddContract, Msg: "only a TODO or BLOCKED task can be added"}
	}
	if fs := nb.Validate(ValidateOpts{Strict: true, RequireAcceptance: true}); len(fs) > 0 {
		return nil, "", &AddError{Kind: AddContract, Msg: "block fails the strict task checks", Findings: fs}
	}
	if t.Status == StatusBLOCKED {
		r := nb.StatusReason(&nb.Tasks[0])
		m := reFindingReason.FindStringSubmatch(r)
		if m == nil {
			return nil, "", &AddError{Kind: AddContract,
				Msg: `a BLOCKED block needs the reason "found during <TASK-ID>: <reason>"`}
		}
		src := m[1]
		if len(p.ByID(src)) == 0 && !containsID(archived, src) {
			return nil, "", &AddError{Kind: AddContract,
				Msg: fmt.Sprintf("source task %s is not in the plan or its archive", src)}
		}
		wantID := regexp.MustCompile("^" + regexp.QuoteMeta(src) + "-F[1-9][0-9]*$")
		if !wantID.MatchString(t.ID) {
			return nil, "", &AddError{Kind: AddContract,
				Msg: fmt.Sprintf("a finding's ID must be %s-F<n>, got %s", src, t.ID)}
		}
	}

	body := bytes.TrimRight(block, "\r\n")
	body = bytes.TrimLeft(body, "\r\n")
	nl := "\n"
	if p.CRLF {
		nl = "\r\n"
		body = toCRLF(body)
	}

	buf := make([]byte, 0, len(p.Bytes)+2*len(nl)+len(body)+len(nl))
	buf = append(buf, p.Bytes...)
	if len(p.Bytes) > 0 {
		if !bytes.HasSuffix(p.Bytes, []byte("\n")) {
			buf = append(buf, nl...)
		}
		buf = append(buf, nl...)
	}
	buf = append(buf, body...)
	buf = append(buf, nl...)
	out = buf

	if err := verifyAppend(p, out, t.ID); err != nil {
		return nil, "", &AddError{Kind: AddSelfCheck, Msg: "self-check failed: " + err.Error()}
	}
	return out, t.ID, nil
}

// hasNonBlankPreamble reports whether nb carries a preamble region (bytes
// before the block's heading) that holds more than whitespace.
func hasNonBlankPreamble(nb *Plan) bool {
	for _, r := range nb.Regions {
		if r.Kind == "preamble" && len(bytes.TrimSpace(nb.Bytes[r.Start:r.End])) > 0 {
			return true
		}
	}
	return false
}

// containsID reports whether ids holds id.
func containsID(ids []string, id string) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

// toCRLF converts every "\n" in b not already preceded by "\r" to "\r\n".
func toCRLF(b []byte) []byte {
	out := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		if b[i] == '\n' && (i == 0 || b[i-1] != '\r') {
			out = append(out, '\r', '\n')
			continue
		}
		out = append(out, b[i])
	}
	return out
}

// verifyAppend re-parses rendered append output and checks the safety
// invariants before the bytes are returned: the output starts with the
// original bytes, exactly one task was added, and it is the last one.
// Returns nil when the rewrite is safe.
func verifyAppend(orig *Plan, out []byte, id string) error {
	if !bytes.HasPrefix(out, orig.Bytes) {
		return fmt.Errorf("output does not start with the original bytes")
	}
	np := Parse(out)
	if len(np.Tasks) != len(orig.Tasks)+1 {
		return fmt.Errorf("task count changed: %d -> %d, want %d", len(orig.Tasks), len(np.Tasks), len(orig.Tasks)+1)
	}
	if got := np.Tasks[len(np.Tasks)-1].ID; got != id {
		return fmt.Errorf("last task is %q, want %q", got, id)
	}
	return nil
}
