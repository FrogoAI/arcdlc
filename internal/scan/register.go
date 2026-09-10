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
	New      []Block // blocks appended by this run, in register order
	Extended []Block // blocks that gained members, carrying only the members added
	Known    int     // findings the sweep saw that the register already held
	Kept     int     // blocks this run left untouched
	Total    int     // blocks in the rendered register
	Changed  bool    // false when the register was already correct, byte for byte
}

// Block is one block of the register: one marker on its own, or every marker that
// shares one group tag.
type Block struct {
	ID      string    `json:"id"`
	Group   string    `json:"group,omitempty"`
	Members []Finding `json:"members"`
}

// Entry is one block already in the register.
type Entry struct {
	ID      string
	Marker  string // the marker word, taken from the ID prefix
	Group   string // the group tag, taken from the ID, "" when the block is ungrouped
	Members []Member
	Verdict string

	rawStart, rawEnd   int
	markerEnd          int // just past the last "- Marker:" line
	whatStart, whatEnd int // the "- WHAT:" line; -1 when the block has none
	whereEnd           int // just past the last line of the WHERE section; -1 when it has none
}

// Member is one marker recorded in a block. A block holds several once markers are
// grouped by tag.
type Member struct {
	File        string
	Line        int
	Text        string
	Fingerprint string
}

var (
	reMarkerLine  = regexp.MustCompile("^-\\s+Marker:\\s+`([^`]*)`\\s+`(.*)`\\s*$")
	reVerdictLine = regexp.MustCompile(`^(\s*)-\s*Verdict:\s*(.*?)\s*$`)
	reWhatLine    = regexp.MustCompile(`^-\s+WHAT:`)
	reWhereLine   = regexp.MustCompile(`^-\s+WHERE:`)
	reSectionKey  = regexp.MustCompile(`^-\s+(WHAT|HOW|WHERE|WHY|References|Acceptance|Status):`)
	reID          = regexp.MustCompile(`^([A-Za-z0-9_]+)-CMT-([A-Za-z0-9_-]+)$`)
)

const registerHeader = "# Comment register\n" + `
This file lists every code comment marker ` + "`arctool scan`" + ` found in this repository.
One block is one marker, or one group of markers that share a tag: ` + "`// ARCDLC:T1 ...`" + ` in
three files is one block, because it is one change.

` + "`arctool scan`" + ` owns the task ID, the ` + "`- Marker:`" + ` lines that are the finding's
identity, the ` + "`- WHERE:`" + ` list, and a first draft of ` + "`- WHAT:`" + `. Everything else is
written by ` + "`/arcdlc:assist`" + ` after it grills the engineer, so do not expect the empty keys
below to stay empty. The verdict is the skill's alone: the sweep never writes one.

The sweep leaves the code alone. When the engineer asks for it (` + "`arctool scan --strip`" + `, which
` + "`/arcdlc:assist`" + ` only runs after a yes), each registered marker comment is deleted from the code,
and this file becomes the only place that marker text lives. A block is never removed, and a sweep never
rewrites the judgement in one. A marker whose group tag is already registered is appended to that block:
another ` + "`- Marker:`" + ` line, its words on the end of WHAT, its file on the end of WHERE.

A finding whose verdict is ACTIONABLE is mirrored into ` + "`plan.md`" + ` as a task, without its
` + "`- Marker:`" + ` and ` + "`- Verdict:`" + ` lines. The mirroring rules live in the plan format guide,
under register sync.

Verdicts, all written by the skill: ACTIONABLE (mirrored into the plan as a task), UNCLEAR (needs
the engineer), STALE (the code already does it), DEFERRED (real work, not planned now). An empty
verdict means nobody has judged the block yet.
`

// ParseEntries reads the blocks of an existing register. It fails when a block
// carries no usable "- Marker:" line, because those lines are the only identity a
// block has.
func ParseEntries(b []byte) ([]Entry, error) {
	p := plan.Parse(b)
	seen := map[string]int{}
	out := make([]Entry, 0, len(p.Tasks))
	for i := range p.Tasks {
		t := &p.Tasks[i]
		e := Entry{ID: t.ID, rawStart: t.RawStart, rawEnd: t.RawEnd,
			markerEnd: -1, whatStart: -1, whatEnd: -1, whereEnd: -1}
		if m := reID.FindStringSubmatch(t.ID); m != nil {
			e.Marker = m[1]
			if !allDigits(m[2]) {
				e.Group = strings.ToUpper(m[2])
			}
		}
		at, inWhere := t.RawStart, false
		for _, raw := range splitKeepEnds(b[t.RawStart:t.RawEnd]) {
			line := strings.TrimRight(string(raw), "\r\n")
			switch {
			case reMarkerLine.MatchString(line):
				m := reMarkerLine.FindStringSubmatch(line)
				mem := Member{Text: Normalize(m[2], maxTextLen)}
				mem.File, mem.Line = splitFileLine(m[1])
				key := mem.File + "\x00" + mem.Text
				mem.Fingerprint = key + "\x00" + strconv.Itoa(seen[key])
				seen[key]++
				e.Members = append(e.Members, mem)
				e.markerEnd, inWhere = at+len(raw), false
			case reVerdictLine.MatchString(line):
				if e.Verdict == "" {
					m := reVerdictLine.FindStringSubmatch(line)
					e.Verdict = strings.TrimSuffix(strings.TrimSpace(m[2]), ".")
				}
				inWhere = false
			case reWhatLine.MatchString(line):
				if e.whatStart < 0 {
					e.whatStart, e.whatEnd = at, at+len(raw)
				}
				inWhere = false
			case reWhereLine.MatchString(line):
				e.whereEnd, inWhere = at+len(raw), true
			case reSectionKey.MatchString(line):
				inWhere = false
			default:
				if inWhere && strings.TrimSpace(line) != "" {
					e.whereEnd = at + len(raw)
				}
			}
			at += len(raw)
		}
		if len(e.Members) == 0 {
			return nil, fmt.Errorf("block %q (line %d) has no usable %q line", t.ID, t.Line, "- Marker:")
		}
		out = append(out, e)
	}
	return out, nil
}

// Render merges the findings of one sweep into an existing register and returns the
// new bytes. A finding the register already holds changes nothing. A finding whose
// group tag is already registered grows that block. Everything else is appended as a
// new block after the last one.
func Render(existing []byte, found []Finding) ([]byte, Result, error) {
	entries, err := ParseEntries(existing)
	if err != nil {
		return nil, Result{}, err
	}
	nl := "\n"
	if plan.Parse(existing).CRLF {
		nl = "\r\n"
	}

	known := map[string]bool{}
	byID := make(map[string]*Entry, len(entries))
	taken := make(map[string]bool, len(entries))
	for i := range entries {
		for _, m := range entries[i].Members {
			known[m.Fingerprint] = true
		}
		byID[entries[i].ID] = &entries[i]
		taken[entries[i].ID] = true
	}

	var res Result
	var blocks []string
	var edits []splice
	nextNum := nextNumbers(entries)
	seen := map[string]int{}

	for _, g := range groupFindings(found) {
		var fresh []Finding
		for _, f := range g.members {
			key := f.File + "\x00" + Normalize(f.Text, maxTextLen)
			fp := key + "\x00" + strconv.Itoa(seen[key])
			seen[key]++
			if known[fp] {
				res.Known++ // already registered: the sweep still removes the comment
				continue
			}
			fresh = append(fresh, f)
		}
		if len(fresh) == 0 {
			continue
		}
		if g.tag != "" {
			if e, ok := byID[groupID(g.marker, g.tag)]; ok {
				edits = append(edits, growBlock(existing, e, fresh, nl)...)
				res.Extended = append(res.Extended, Block{ID: e.ID, Group: g.tag, Members: fresh})
				continue
			}
		}
		id := groupID(g.marker, g.tag) // "" for an untagged marker
		if id == "" || taken[id] {
			id = autoID(g.marker, nextNum, taken)
		}
		taken[id] = true
		res.New = append(res.New, Block{ID: id, Group: g.tag, Members: fresh})
		blocks = append(blocks, renderBlock(id, fresh, nl))
	}

	out := make([]byte, 0, len(existing)+len(blocks)*400)
	if len(strings.TrimSpace(string(existing))) == 0 {
		out = append(out, []byte(withTerm(registerHeader, nl))...)
	} else {
		out = append(out, applySplices(existing, edits)...)
	}
	if len(blocks) > 0 {
		out = ensureBlankTail(out, nl)
		out = append(out, []byte(strings.Join(blocks, nl))...)
	}

	res.Kept = len(entries) - len(res.Extended)
	res.Total = len(entries) + len(res.New)
	res.Changed = len(res.New) > 0 || len(res.Extended) > 0 || string(out) != string(existing)
	return out, res, nil
}

// group is the sweep's findings for one block: a marker word plus its group tag, or
// one untagged marker on its own.
type group struct {
	marker  string
	tag     string
	members []Finding
}

// groupFindings collects the sweep by marker word and group tag, keeping the order
// the sweep found them in. An untagged marker is a group of one, so the caller has
// one shape to handle. FIXME:G1 is not TODO:G1: the marker word leads the block ID,
// and one run registers one kind of debt.
func groupFindings(found []Finding) []group {
	var out []group
	at := map[string]int{}
	for _, f := range found {
		if f.Group == "" {
			out = append(out, group{marker: f.Marker, members: []Finding{f}})
			continue
		}
		key := f.Marker + "\x00" + f.Group
		if i, ok := at[key]; ok {
			out[i].members = append(out[i].members, f)
			continue
		}
		at[key] = len(out)
		out = append(out, group{marker: f.Marker, tag: f.Group, members: []Finding{f}})
	}
	return out
}

// splice is one edit to the register's bytes: replace [start,end) with text. An
// insertion has start == end.
type splice struct {
	start, end int
	text       string
}

// applySplices rewrites only the ranges the edits name and copies everything else
// byte for byte.
func applySplices(existing []byte, edits []splice) []byte {
	sort.SliceStable(edits, func(i, j int) bool { return edits[i].start < edits[j].start })
	out := make([]byte, 0, len(existing)+256*len(edits))
	at := 0
	for _, e := range edits {
		if e.start < at {
			continue // overlapping edit: Verify rejects the result, nothing is written
		}
		out = append(out, existing[at:e.start]...)
		out = append(out, e.text...)
		at = e.end
	}
	return append(out, existing[at:]...)
}

// growBlock is the edit that adds members to a block the register already holds: a
// "- Marker:" line each, their words on the end of WHAT, their files on the end of
// WHERE. Nothing else in the block is touched, so the title, HOW, WHY, Acceptance and
// the verdict the skill wrote all stand.
func growBlock(existing []byte, e *Entry, fresh []Finding, nl string) []splice {
	var markers, where strings.Builder
	for _, f := range fresh {
		markers.WriteString(markerLine(f, nl))
		where.WriteString(whereLines(f, nl))
	}
	edits := []splice{{start: e.markerEnd, end: e.markerEnd, text: markers.String()}}

	whatAt := e.markerEnd
	if e.whatStart >= 0 {
		line := strings.TrimRight(string(existing[e.whatStart:e.whatEnd]), "\r\n")
		edits = append(edits, splice{start: e.whatStart, end: e.whatEnd, text: appendWhat(line, fresh) + nl})
		whatAt = e.whatEnd
	} else if what := whatText(fresh); what != "" {
		edits = append(edits, splice{start: whatAt, end: whatAt, text: "- WHAT: " + what + nl})
	}
	if e.whereEnd >= 0 {
		edits = append(edits, splice{start: e.whereEnd, end: e.whereEnd, text: where.String()})
	} else {
		edits = append(edits, splice{start: whatAt, end: whatAt, text: "- WHERE:" + nl + where.String()})
	}
	return edits
}

// appendWhat puts the new members' words on the end of a WHAT line, keeping whatever
// is already there. WHAT stays one line: the runner reads no further. The full stop
// the last sentence ended on goes, so the line reads as one list.
func appendWhat(line string, fresh []Finding) string {
	add := whatText(fresh)
	if add == "" {
		return line
	}
	head := strings.TrimRight(line, " \t")
	if i := strings.Index(head, ":"); i >= 0 && strings.TrimSpace(head[i+1:]) == "" {
		return head + " " + add
	}
	return strings.TrimRight(head, " .;") + "; " + add
}

// Verify re-reads what Render produced before anything is written. It proves that no
// block and no marker was lost, that a block nobody added to is byte identical, and
// that a block which gained members changed in three places only.
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
	grew := make(map[string]int, len(res.Extended))
	for _, b := range res.Extended {
		grew[b.ID] = len(b.Members)
	}
	for _, b := range before {
		got, ok := byID[b.ID]
		if !ok {
			return fmt.Errorf("block %q went missing", b.ID)
		}
		if err := checkKept(existing, out, b, got[0], grew[b.ID]); err != nil {
			return err
		}
	}
	for _, b := range res.New {
		if _, ok := byID[b.ID]; !ok {
			return fmt.Errorf("new block %q is not in the rendered register", b.ID)
		}
	}
	return nil
}

// checkKept proves that a block the register already held came through this run
// whole. A block nobody added to must not move a byte. A block that gained members
// may differ in three places and nowhere else: the "- Marker:" lines it gained, more
// text on the end of WHAT, more lines on the end of WHERE. So a sweep cannot rewrite
// the judgement the skill wrote.
func checkKept(existing, out []byte, b, a Entry, added int) error {
	if added == 0 {
		if blockText(existing, b.rawStart, b.rawEnd) != blockText(out, a.rawStart, a.rawEnd) {
			return fmt.Errorf("block %q changed; only a block that gains a marker may change", b.ID)
		}
		return nil
	}
	if len(a.Members) != len(b.Members)+added {
		return fmt.Errorf("block %q holds %d marker(s), expected %d",
			b.ID, len(a.Members), len(b.Members)+added)
	}
	for i, m := range b.Members {
		if a.Members[i].Fingerprint != m.Fingerprint {
			return fmt.Errorf("block %q lost the marker %s:%d", b.ID, m.File, m.Line)
		}
	}
	wasSkeleton, wasWhere, wasWhat := blockParts(existing, b)
	isSkeleton, isWhere, isWhat := blockParts(out, a)
	if wasSkeleton != isSkeleton {
		return fmt.Errorf("block %q changed outside its markers, WHAT and WHERE", b.ID)
	}
	// A sweep may only add to WHAT, and may drop the full stop it writes behind.
	if !strings.HasPrefix(isWhat, strings.TrimRight(wasWhat, " .;")) {
		return fmt.Errorf("block %q rewrote WHAT; a sweep may only add to it", b.ID)
	}
	if len(isWhere) < len(wasWhere) {
		return fmt.Errorf("block %q dropped a WHERE line", b.ID)
	}
	for i, line := range wasWhere {
		if isWhere[i] != line {
			return fmt.Errorf("block %q rewrote the WHERE line %q", b.ID, line)
		}
	}
	return nil
}

// blockParts splits one block into the pieces Verify compares: the content lines of
// the WHERE section, the text on the WHAT line, and a skeleton, which is every other
// line except the "- Marker:" lines. A sweep may add to the first two. The skeleton
// must not move.
func blockParts(src []byte, e Entry) (skeleton string, where []string, what string) {
	var keep []string
	inWhere := false
	for _, raw := range splitKeepEnds(src[e.rawStart:e.rawEnd]) {
		line := strings.TrimRight(string(raw), "\r\n")
		switch {
		case reMarkerLine.MatchString(line):
			inWhere = false
		case reWhatLine.MatchString(line):
			what = strings.TrimSpace(line[strings.Index(line, ":")+1:])
			inWhere = false
		case reWhereLine.MatchString(line):
			inWhere = true
			keep = append(keep, line)
		case reSectionKey.MatchString(line):
			inWhere = false
			keep = append(keep, line)
		default:
			if inWhere && strings.TrimSpace(line) != "" {
				where = append(where, strings.TrimSpace(line))
				continue
			}
			keep = append(keep, line)
		}
	}
	return strings.TrimRight(strings.Join(keep, "\n"), " \t\n"), where, what
}

// blockText returns a block's bytes without the blank line that separates it from
// the next one, so appending a block never counts as changing the one above it.
func blockText(src []byte, start, end int) string {
	return strings.TrimRight(string(src[start:end]), " \t\r\n")
}

// renderBlock writes the lines the scan owns and leaves the judgement empty.
func renderBlock(id string, members []Finding, nl string) string {
	var b strings.Builder
	b.WriteString("### " + id + ": " + title(members[0]) + nl)
	b.WriteString(nl)
	for _, f := range members {
		b.WriteString(markerLine(f, nl))
	}
	b.WriteString("- Verdict:" + nl)
	if what := whatText(members); what != "" {
		b.WriteString("- WHAT: " + what + nl)
	} else {
		b.WriteString("- WHAT:" + nl)
	}
	b.WriteString("- HOW:" + nl)
	b.WriteString("- WHERE:" + nl)
	for _, f := range members {
		b.WriteString(whereLines(f, nl))
	}
	b.WriteString("- WHY:" + nl)
	b.WriteString("- Acceptance:" + nl)
	return b.String()
}

// markerLine writes one member's identity: where the marker was, and what it said.
func markerLine(f Finding, nl string) string {
	return "- Marker: `" + fmt.Sprintf("%s:%d", f.File, f.Line) + "` `" +
		Normalize(f.Text, maxTextLen) + "`" + nl
}

// whereLines writes one member's location, and the code line under its marker. It
// avoids a bare "file:line", because plan.ParseWhereLayers reads any "name: targets"
// line as a layer, and "internal.go: 3" is not one.
func whereLines(f Finding, nl string) string {
	s := fmt.Sprintf("  %s (marker at line %d)%s", f.File, f.Line, nl)
	if code := Normalize(f.Code, maxCodeLen); code != "" {
		s += "    " + code + nl
	}
	return s
}

// whatText seeds WHAT with what the engineer wrote: the marker text without the
// marker word and its tag, members joined in sweep order. The skill rewrites it into
// an instruction.
func whatText(members []Finding) string {
	var parts []string
	for _, f := range members {
		if t := body(f); t != "" {
			parts = append(parts, t)
		}
	}
	return Normalize(strings.Join(parts, "; "), maxTextLen)
}

// body is the marker text without the marker word and the group tag: the words the
// engineer actually typed.
func body(f Finding) string {
	t := strings.TrimPrefix(Normalize(f.Text, maxTextLen), f.Marker)
	if f.Group != "" {
		head, rest := t, ""
		if i := strings.IndexAny(t, " \t"); i > 0 {
			head, rest = t[:i], t[i:]
		}
		if strings.EqualFold(strings.TrimPrefix(head, ":"), f.Group) {
			t = rest
		}
	}
	return strings.TrimSpace(strings.TrimLeft(t, ":-() \t"))
}

// leadIn lists the openings marker text often starts with. Dropping them leaves a
// draft title that reads closer to the instruction the skill will write.
var leadIn = []string{"we should ", "we need to ", "we must ", "we have to ", "should ",
	"need to ", "must ", "please ", "maybe ", "consider ", "it would be good to "}

// title drafts a heading from the marker text. The skill rewrites it into an
// instruction; the scan only makes the block readable in the meantime.
func title(f Finding) string {
	t := body(f)
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

// groupID is the block ID for a group tag: the tag takes the slot an ungrouped block
// fills with a number.
func groupID(marker, tag string) string {
	if tag == "" {
		return ""
	}
	return marker + "-CMT-" + tag
}

// autoID is the next free numbered ID for a marker word. It steps over an ID the
// register already holds, so a block is never numbered on top of another.
func autoID(marker string, next map[string]int, taken map[string]bool) string {
	for {
		next[marker]++
		id := fmt.Sprintf("%s-CMT-%02d", marker, next[marker])
		if !taken[id] {
			return id
		}
	}
}

// nextNumbers returns the highest block number already used per marker. A block
// named by a group tag has no number to count.
func nextNumbers(entries []Entry) map[string]int {
	out := map[string]int{}
	for _, e := range entries {
		m := reID.FindStringSubmatch(e.ID)
		if m == nil || !allDigits(m[2]) {
			continue
		}
		if n, err := strconv.Atoi(m[2]); err == nil && n > out[m[1]] {
			out[m[1]] = n
		}
	}
	return out
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isDigitByte(s[i]) {
			return false
		}
	}
	return s != ""
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

func withTerm(s, nl string) string {
	if nl == "\r\n" {
		return strings.ReplaceAll(s, "\n", "\r\n")
	}
	return s
}

// ensureBlankTail leaves the register ending in exactly one blank line, so an
// appended block always sits one empty line below the previous one.
func ensureBlankTail(out []byte, nl string) []byte {
	s := strings.TrimRight(string(out), "\r\n ")
	if s == "" {
		return out
	}
	return []byte(s + nl + nl)
}
