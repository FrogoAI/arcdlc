package scan

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/FrogoAI/arcdlc/internal/plan"
)

// Result reports what one Render call changed.
type Result struct {
	New      []Finding // findings appended by this run, in register order
	NewIDs   []string  // their task IDs, same order as New
	Resolved []string  // IDs whose verdict this run flipped to RESOLVED
	Kept     int       // blocks left untouched
	Total    int       // blocks in the rendered register
	Changed  bool      // false when the register is already correct, byte for byte
}

// Entry is one block already in the register.
type Entry struct {
	ID          string
	Marker      string // marker word taken from the ID prefix
	File        string
	Line        int
	Text        string
	Verdict     string
	Fingerprint string

	rawStart, rawEnd         int
	verdictStart, verdictEnd int // -1 when the block carries no verdict line
}

var (
	reMarkerLine  = regexp.MustCompile("^-\\s+Marker:\\s+`([^`]*)`\\s+`(.*)`\\s*$")
	reVerdictLine = regexp.MustCompile(`^(\s*)-\s*Verdict:\s*(.*?)\s*$`)
	reID          = regexp.MustCompile(`^([A-Za-z0-9_]+)-CMT-(\d+)$`)
)

const (
	// VerdictNew is what scan writes for a finding nobody has judged yet.
	VerdictNew = "NEW"
	// VerdictResolved is the only verdict scan ever writes over another one,
	// and only when the marker has gone from the code.
	VerdictResolved = "RESOLVED"
)

const registerHeader = "# Comment register\n" + `
This file lists every code comment marker ` + "`arctool scan`" + ` found in this repository.
One block is one marker.

` + "`arctool scan`" + ` owns three things in a block: the task ID, the ` + "`- Marker:`" + ` line, and the
flip to ` + "`- Verdict: RESOLVED (<date>).`" + ` once the marker is gone from the code. Everything
else is written by ` + "`/arcdlc:assist`" + ` after it grills the engineer, so do not expect the
empty keys below to stay empty.

A finding whose verdict is ACTIONABLE is mirrored into ` + "`plan.md`" + ` as a task, without its
` + "`- Marker:`" + ` and ` + "`- Verdict:`" + ` lines. The mirroring rules live in the plan format guide,
under register sync.

Verdicts: NEW (not judged yet), ACTIONABLE (mirrored into the plan), UNCLEAR (needs the
engineer), STALE (the code already does it), DEFERRED (real work, not planned now),
RESOLVED (the marker is gone).
`

// ParseEntries reads the blocks of an existing register. It fails when a block
// carries no usable "- Marker:" line, because that line is the only identity a
// finding has.
func ParseEntries(b []byte) ([]Entry, error) {
	p := plan.Parse(b)
	seen := map[string]int{}
	out := make([]Entry, 0, len(p.Tasks))
	for i := range p.Tasks {
		t := &p.Tasks[i]
		e := Entry{ID: t.ID, rawStart: t.RawStart, rawEnd: t.RawEnd, verdictStart: -1, verdictEnd: -1}
		if m := reID.FindStringSubmatch(t.ID); m != nil {
			e.Marker = m[1]
		}
		body := b[t.RawStart:t.RawEnd]
		at := t.RawStart
		for _, raw := range splitKeepEnds(body) {
			line := strings.TrimRight(string(raw), "\r\n")
			if m := reMarkerLine.FindStringSubmatch(line); m != nil && e.File == "" {
				e.File, e.Line = splitFileLine(m[1])
				e.Text = Normalize(m[2], maxTextLen)
			}
			if m := reVerdictLine.FindStringSubmatch(line); m != nil && e.verdictStart < 0 {
				e.Verdict = strings.TrimSuffix(strings.TrimSpace(m[2]), ".")
				e.verdictStart, e.verdictEnd = at, at+len(raw)
			}
			at += len(raw)
		}
		if e.File == "" {
			return nil, fmt.Errorf("block %q (line %d) has no usable %q line", t.ID, t.Line, "- Marker:")
		}
		key := e.File + "\x00" + e.Text
		e.Fingerprint = key + "\x00" + strconv.Itoa(seen[key])
		seen[key]++
		out = append(out, e)
	}
	return out, nil
}

// Render merges the findings of one sweep into an existing register and returns
// the new bytes. Existing blocks are preserved byte for byte, except that a
// block whose marker has gone from the code gets its verdict line rewritten.
// New findings are appended after the last block.
func Render(existing []byte, found []Finding, markers []string, date string) ([]byte, Result, error) {
	entries, err := ParseEntries(existing)
	if err != nil {
		return nil, Result{}, err
	}
	crlf := plan.Parse(existing).CRLF
	nl := "\n"
	if crlf {
		nl = "\r\n"
	}

	swept := make(map[string]bool, len(markers))
	for _, m := range markers {
		swept[m] = true
	}

	byFingerprint := make(map[string]*Entry, len(entries))
	for i := range entries {
		byFingerprint[entries[i].Fingerprint] = &entries[i]
	}

	// Findings the register does not know yet, in sweep order.
	seen := map[string]int{}
	var res Result
	nextNum := nextNumbers(entries)
	var blocks []string
	for _, f := range found {
		key := f.File + "\x00" + Normalize(f.Text, maxTextLen)
		fp := key + "\x00" + strconv.Itoa(seen[key])
		seen[key]++
		if _, ok := byFingerprint[fp]; ok {
			continue
		}
		nextNum[f.Marker]++
		id := fmt.Sprintf("%s-CMT-%02d", f.Marker, nextNum[f.Marker])
		res.New = append(res.New, f)
		res.NewIDs = append(res.NewIDs, id)
		blocks = append(blocks, renderBlock(id, f, nl))
	}

	// Markers the register knows that the sweep no longer finds.
	live := make(map[string]bool, len(found))
	seen = map[string]int{}
	for _, f := range found {
		key := f.File + "\x00" + Normalize(f.Text, maxTextLen)
		live[key+"\x00"+strconv.Itoa(seen[key])] = true
		seen[key]++
	}
	type splice struct {
		start, end int
		text       string
	}
	var splices []splice
	for i := range entries {
		e := &entries[i]
		if !swept[e.Marker] || live[e.Fingerprint] || strings.HasPrefix(e.Verdict, VerdictResolved) {
			res.Kept++
			continue
		}
		if e.verdictStart < 0 {
			res.Kept++ // no verdict line to rewrite: leave the block alone
			continue
		}
		orig := string(existing[e.verdictStart:e.verdictEnd])
		indent := ""
		if m := reVerdictLine.FindStringSubmatch(strings.TrimRight(orig, "\r\n")); m != nil {
			indent = m[1]
		}
		term := lineTerm(orig, nl)
		splices = append(splices, splice{e.verdictStart, e.verdictEnd,
			fmt.Sprintf("%s- Verdict: %s (%s).%s", indent, VerdictResolved, date, term)})
		res.Resolved = append(res.Resolved, e.ID)
	}
	sort.Slice(splices, func(i, j int) bool { return splices[i].start < splices[j].start })

	out := make([]byte, 0, len(existing)+len(blocks)*400)
	if len(strings.TrimSpace(string(existing))) == 0 {
		out = append(out, []byte(withTerm(registerHeader, nl))...)
	} else {
		at := 0
		for _, s := range splices {
			out = append(out, existing[at:s.start]...)
			out = append(out, []byte(s.text)...)
			at = s.end
		}
		out = append(out, existing[at:]...)
	}
	if len(blocks) > 0 {
		out = ensureBlankTail(out, nl)
		out = append(out, []byte(strings.Join(blocks, nl))...)
	}

	res.Total = len(entries) + len(res.New)
	res.Changed = len(res.New) > 0 || len(res.Resolved) > 0 || string(out) != string(existing)
	return out, res, nil
}

// Verify re-reads what Render produced before anything is written. It proves
// that no block was lost, that no untouched block changed, and that a resolved
// block changed on its verdict line only.
func Verify(existing, out []byte, res Result) error {
	before, err := ParseEntries(existing)
	if err != nil {
		return err
	}
	after, err := ParseEntries(out)
	if err != nil {
		return fmt.Errorf("rendered register does not parse: %w", err)
	}
	if len(after) != res.Total {
		return fmt.Errorf("rendered %d block(s), expected %d", len(after), res.Total)
	}
	byID := make(map[string][]Entry, len(after))
	for _, e := range after {
		byID[e.ID] = append(byID[e.ID], e)
	}
	for id, got := range byID {
		if len(got) != 1 {
			return fmt.Errorf("task ID %q appears %d times", id, len(got))
		}
	}
	resolved := make(map[string]bool, len(res.Resolved))
	for _, id := range res.Resolved {
		resolved[id] = true
	}
	for _, b := range before {
		got, ok := byID[b.ID]
		if !ok {
			return fmt.Errorf("block %q went missing", b.ID)
		}
		a := got[0]
		if a.Fingerprint != b.Fingerprint {
			return fmt.Errorf("block %q lost its marker line", b.ID)
		}
		oldRaw := blockText(existing, b.rawStart, b.rawEnd)
		newRaw := blockText(out, a.rawStart, a.rawEnd)
		if oldRaw == newRaw {
			continue
		}
		if !resolved[b.ID] {
			return fmt.Errorf("block %q changed but was not resolved", b.ID)
		}
		if strip(oldRaw, b.verdictStart, b.verdictEnd, b.rawStart) !=
			strip(newRaw, a.verdictStart, a.verdictEnd, a.rawStart) {
			return fmt.Errorf("block %q changed outside its verdict line", b.ID)
		}
	}
	for _, id := range res.NewIDs {
		if _, ok := byID[id]; !ok {
			return fmt.Errorf("new block %q is not in the rendered register", id)
		}
	}
	return nil
}

// blockText returns a block's bytes without the blank line that separates it
// from the next one, so appending a block never counts as changing the one
// above it.
func blockText(src []byte, start, end int) string {
	return strings.TrimRight(string(src[start:end]), " \t\r\n")
}

// strip removes the verdict line from a block's raw text, so two versions of a
// resolved block can be compared on everything else.
func strip(raw string, vStart, vEnd, rawStart int) string {
	if vStart < 0 || vEnd-rawStart > len(raw) {
		return raw
	}
	return raw[:vStart-rawStart] + raw[vEnd-rawStart:]
}

// renderBlock writes the lines scan owns and leaves the judgement empty.
func renderBlock(id string, f Finding, nl string) string {
	// The WHERE line avoids a bare "file:line": plan.ParseWhereLayers reads any
	// "name: targets" line there as a layer, and "internal.go: 3" is not one.
	where := fmt.Sprintf("%s (marker at line %d)", f.File, f.Line)
	var b strings.Builder
	b.WriteString("### " + id + ": " + title(f) + nl)
	b.WriteString(nl)
	b.WriteString("- Marker: `" + fmt.Sprintf("%s:%d", f.File, f.Line) + "` `" + Normalize(f.Text, maxTextLen) + "`" + nl)
	b.WriteString("- Verdict: " + VerdictNew + "." + nl)
	b.WriteString("- WHAT:" + nl)
	b.WriteString("- HOW:" + nl)
	b.WriteString("- WHERE: " + where + nl)
	if code := Normalize(f.Code, maxCodeLen); code != "" {
		b.WriteString("  " + code + nl)
	}
	b.WriteString("- WHY:" + nl)
	b.WriteString("- Acceptance:" + nl)
	return b.String()
}

// leadIn lists the openings a marker text often starts with. Dropping them
// leaves a draft title that reads closer to the instruction the skill will write.
var leadIn = []string{"we should ", "we need to ", "we must ", "we have to ", "should ",
	"need to ", "must ", "please ", "maybe ", "consider ", "it would be good to "}

// title drafts a heading from the marker text. The skill rewrites it into an
// instruction; scan only makes the block readable in the meantime.
func title(f Finding) string {
	t := strings.TrimSpace(strings.TrimPrefix(Normalize(f.Text, maxTextLen), f.Marker))
	t = strings.TrimLeft(t, ":-() \t")
	for _, lead := range leadIn {
		if len(t) > len(lead) && strings.EqualFold(t[:len(lead)], lead) {
			t = t[len(lead):]
			break
		}
	}
	t = firstSentence(t)
	if t == "" {
		return fmt.Sprintf("Handle the %s marker in %s", f.Marker, f.File)
	}
	// A first clause makes a better title than a whole sentence.
	if i := strings.Index(t, ","); i >= 20 {
		t = t[:i]
	}
	const limit = 68
	if len(t) > limit {
		cut := strings.LastIndex(t[:limit], " ")
		if cut < 20 {
			cut = limit
		}
		t = t[:cut]
	}
	t = strings.TrimRight(t, " ,.;:")
	return strings.ToUpper(t[:1]) + t[1:]
}

func firstSentence(s string) string {
	if i := strings.IndexAny(s, ".!?"); i > 20 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

// nextNumbers returns the highest block number already used per marker.
func nextNumbers(entries []Entry) map[string]int {
	out := map[string]int{}
	for _, e := range entries {
		m := reID.FindStringSubmatch(e.ID)
		if m == nil {
			continue
		}
		if n, err := strconv.Atoi(m[2]); err == nil && n > out[m[1]] {
			out[m[1]] = n
		}
	}
	return out
}

// splitKeepEnds splits b into lines, each keeping its line terminator.
func splitKeepEnds(b []byte) [][]byte {
	var out [][]byte
	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == '\n' {
			out = append(out, b[start:i+1])
			start = i + 1
		}
	}
	if start < len(b) {
		out = append(out, b[start:])
	}
	return out
}

func splitFileLine(s string) (string, int) {
	if i := strings.LastIndex(s, ":"); i > 0 {
		if n, err := strconv.Atoi(s[i+1:]); err == nil {
			return s[:i], n
		}
	}
	return s, 0
}

func lineTerm(line, fallback string) string {
	switch {
	case strings.HasSuffix(line, "\r\n"):
		return "\r\n"
	case strings.HasSuffix(line, "\n"):
		return "\n"
	}
	return fallback
}

func withTerm(s, nl string) string {
	if nl == "\r\n" {
		return strings.ReplaceAll(s, "\n", "\r\n")
	}
	return s
}

// ensureBlankTail makes the register end with exactly one blank line, so an
// appended block always sits one empty line below the previous one.
func ensureBlankTail(out []byte, nl string) []byte {
	s := strings.TrimRight(string(out), "\r\n ")
	if s == "" {
		return out
	}
	return []byte(s + nl + nl)
}
