package plan

import (
	"regexp"
	"strings"
)

// Action is the outcome a status command should take for a given current state.
type Action uint8

const (
	ActRefuse Action = iota // disallowed transition (guardrail)
	ActNoop                 // already in the requested state; do nothing, succeed
	ActWrite                // rewrite the status line
)

// Transition resolves a status command against the current status. cmd is one
// of "take", "done", "block", "todo". force relaxes the guardrails; hasReason
// is true when `block` was given an explicit -m reason. On ActRefuse, expected
// is a short human hint of the state the command needs.
func Transition(cmd string, cur Status, force, hasReason bool) (target Status, act Action, expected string) {
	switch cmd {
	case "take":
		target = StatusTAKEN
		switch {
		case cur == StatusTAKEN:
			act = ActNoop
		case cur == StatusTODO || force:
			act = ActWrite
		default:
			act, expected = ActRefuse, "TODO"
		}
	case "done":
		target = StatusDONE
		switch {
		case cur == StatusDONE:
			act = ActNoop
		case cur == StatusTAKEN || force:
			act = ActWrite
		default:
			act, expected = ActRefuse, "TAKEN"
		}
	case "todo":
		target = StatusTODO
		switch {
		case cur == StatusTODO:
			act = ActNoop
		case cur == StatusTAKEN || cur == StatusBLOCKED || force:
			act = ActWrite
		default:
			act, expected = ActRefuse, "TAKEN or BLOCKED"
		}
	case "block":
		target = StatusBLOCKED
		switch {
		case cur == StatusBLOCKED:
			if hasReason { // updating the reason is a real change
				act = ActWrite
			} else {
				act = ActNoop
			}
		case cur == StatusDONE && !force:
			act, expected = ActRefuse, "not DONE"
		default:
			act = ActWrite
		}
	}
	return
}

var reStatusLine = regexp.MustCompile(`^(\s*)-\s*Status:\s*(.*?)\s*$`)

// WithStatus returns plan bytes with task t's "- Status:" line rewritten to s,
// preserving the line's leading indentation, trailing-period style, and line
// terminator. Every other byte is identical. For StatusBLOCKED an explicit
// reason ("" keeps any existing reason) is rendered as " — reason". The caller
// must ensure t.StatusLineStart >= 0 (a block with no status line cannot be
// rewritten — that is a validate error, not an implicit insert).
func (p *Plan) WithStatus(t *Task, s Status, reason string) []byte {
	if t.StatusLineStart < 0 {
		return append([]byte(nil), p.Bytes...) // no status line: nothing to rewrite
	}
	orig := p.Bytes[t.StatusLineStart:t.StatusLineEnd]
	indent, hadPeriod, term, existingReason := dissectStatusLine(orig)

	value := s.String()
	if s == StatusBLOCKED {
		r := reason
		if r == "" {
			r = existingReason
		}
		if r != "" {
			value += " — " + r
		}
	}
	if hadPeriod {
		value += "."
	}
	newLine := indent + "- Status: " + value + term

	out := make([]byte, 0, len(p.Bytes)-len(orig)+len(newLine))
	out = append(out, p.Bytes[:t.StatusLineStart]...)
	out = append(out, newLine...)
	out = append(out, p.Bytes[t.StatusLineEnd:]...)
	return out
}

// dissectStatusLine splits a raw "- Status: ..." line into its indentation,
// whether it ends with a period, its line terminator, and any BLOCKED reason.
func dissectStatusLine(line []byte) (indent string, hadPeriod bool, term, reason string) {
	s := string(line)
	switch {
	case strings.HasSuffix(s, "\r\n"):
		term, s = "\r\n", strings.TrimSuffix(s, "\r\n")
	case strings.HasSuffix(s, "\n"):
		term, s = "\n", strings.TrimSuffix(s, "\n")
	}
	m := reStatusLine.FindStringSubmatch(s)
	if m == nil {
		return "", false, term, ""
	}
	indent = m[1]
	val := strings.TrimSpace(m[2])
	if strings.HasSuffix(val, ".") {
		hadPeriod = true
		val = strings.TrimSpace(strings.TrimRight(val, "."))
	}
	if fields := strings.Fields(val); len(fields) > 0 && strings.ToUpper(fields[0]) == "BLOCKED" {
		rest := strings.TrimSpace(val[len(fields[0]):])
		reason = strings.TrimSpace(strings.TrimLeft(rest, "—- \t"))
	}
	return
}
