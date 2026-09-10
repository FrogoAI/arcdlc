// Package scan sweeps a source tree for code comment markers (TODO, FIXME and
// friends) and renders them into the comment register, docs/aics/<slug>/comments.md.
//
// The split of work is deliberate: this package finds evidence (file, line, the
// marker text, the code line under it) and writes only the lines it owns. The
// judgement (title, HOW, WHY, Acceptance, verdict) belongs to the
// /arcdlc:assist skill, which fills it in after it grills the engineer. WHAT is
// seeded here from the marker text and rewritten there.
//
// A marker may carry a group tag, "// TODO:G1 move the rebuild out of main". Every
// marker with the same tag is one block of the register, however many files it
// spans, so one change written down in three places is one task.
//
// It is a pure, standard-library-only implementation.
package scan

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Finding is one marker found in one file.
type Finding struct {
	File   string `json:"file"`            // path relative to the sweep root, slash separated
	Line   int    `json:"line"`            // 1-based line of the marker
	Marker string `json:"marker"`          // the marker word, e.g. "TODO"
	Group  string `json:"group,omitempty"` // group tag, upper cased, "" when the marker stands alone
	Text   string `json:"text"`            // the comment from the marker on, continuation lines joined
	Code   string `json:"code"`            // first code line under the comment, "" when there is none

	// strip is how Strip removes this comment from its file: which whole lines to
	// delete, and which spans of a line to cut. It is derived by the sweep, which
	// is the only place that knows the comment's exact shape.
	strip stripPlan
}

// stripPlan records the edit that deletes one marker comment and nothing else.
// Line numbers are 1-based and inclusive.
type stripPlan struct {
	ok               bool
	reason           string    // why the comment must stay, when ok is false
	fromLine, toLine int       // whole lines to delete; 0 when there are none
	cuts             []lineCut // spans to cut out of a line that keeps its code
}

// lineCut is one span to cut out of one line, leaving the rest of that line alone.
type lineCut struct {
	line     int // 1-based
	from, to int // byte range to cut
}

// span is the shape of one comment the sweep matched: the lines it covers, where
// it opens, and, for a block comment that closes on a later line, where that
// closer ends.
type span struct {
	from, to int    // 0-based line indices, inclusive
	opener   string // the comment opener on line from
	at       int    // column of that opener
	closeAt  int    // column just past the closer on line to; 0 unless the block closes lower down
	unclosed bool   // a block comment with no closer anywhere below
}

// Opts configures a sweep.
type Opts struct {
	Root    string   // directory to walk; "" means "."
	Markers []string // marker words to look for; empty means DefaultMarkers
	Exclude []string // directory names to skip; nil means DefaultExclude
}

// DefaultMarkers is what a sweep looks for when the caller names nothing. ARCDLC is
// the bundle's own word, so a sweep never touches the TODO and FIXME notes a team
// already had. A team that wants those swept asks for them: --marker TODO.
var DefaultMarkers = []string{"ARCDLC"}

// DefaultExclude lists directory names a sweep never enters. docs is on the
// list because the register itself lives there: a plan task or a gap block is
// not a code comment.
var DefaultExclude = []string{".git", "vendor", "node_modules", "dist", "bin", "docs"}

// Comment openers are chosen per file type, because a marker must sit behind a
// comment start that the language actually has. Reading "--" or "%" as a
// comment in a Go file turns every "--json" flag and every "%s" format string
// into a false finding.
var (
	cFamily  = []string{"//", "/*", "*"}
	hashOnly = []string{"#"}
	dashOnly = []string{"--"}

	openersByExt = map[string][]string{
		".go": cFamily, ".c": cFamily, ".h": cFamily, ".cc": cFamily, ".cpp": cFamily,
		".hpp": cFamily, ".cs": cFamily, ".java": cFamily, ".js": cFamily, ".jsx": cFamily,
		".mjs": cFamily, ".cjs": cFamily, ".ts": cFamily, ".tsx": cFamily, ".swift": cFamily,
		".kt": cFamily, ".kts": cFamily, ".rs": cFamily, ".scala": cFamily, ".php": cFamily,
		".dart": cFamily, ".proto": cFamily, ".zig": cFamily, ".groovy": cFamily,
		".css": cFamily, ".scss": cFamily, ".less": cFamily,

		".py": hashOnly, ".rb": hashOnly, ".sh": hashOnly, ".bash": hashOnly, ".zsh": hashOnly,
		".fish": hashOnly, ".yml": hashOnly, ".yaml": hashOnly, ".toml": hashOnly,
		".tf": hashOnly, ".pl": hashOnly, ".pm": hashOnly, ".r": hashOnly, ".mk": hashOnly,
		".cmake": hashOnly, ".env": hashOnly, ".gemspec": hashOnly, ".rake": hashOnly,

		".sql": {"--", "/*", "*"}, ".lua": dashOnly, ".hs": dashOnly, ".elm": dashOnly,
		".el": {";"}, ".lisp": {";"}, ".clj": {";"}, ".scm": {";"},
		".ini": {";", "#"}, ".cfg": {";", "#"},
		".tex": {"%"}, ".erl": {"%"}, ".m": {"%", "//"},
		".html": {"<!--"}, ".htm": {"<!--"}, ".xml": {"<!--"}, ".vue": {"<!--", "//", "/*"},
		".svelte": {"<!--", "//", "/*"},
	}

	openersByName = map[string][]string{
		"Makefile": hashOnly, "makefile": hashOnly, "Dockerfile": hashOnly,
		"Justfile": hashOnly, "justfile": hashOnly, "Rakefile": hashOnly,
	}

	// defaultOpeners covers a file type this list does not name.
	defaultOpeners = []string{"//", "#"}
)

// openersFor returns the comment openers of one file.
func openersFor(path string) []string {
	if ops, ok := openersByName[filepath.Base(path)]; ok {
		return ops
	}
	if ops, ok := openersByExt[strings.ToLower(filepath.Ext(path))]; ok {
		return ops
	}
	return defaultOpeners
}

// docExt lists document extensions a sweep skips. The register is about source
// code, so a marker written in prose is not a finding.
var docExt = map[string]bool{".md": true, ".markdown": true, ".txt": true, ".rst": true}

const (
	maxFileBytes = 1 << 20 // skip anything larger; source files are not this big
	sniffBytes   = 8 << 10 // read this much to decide whether a file is binary
	maxTextLen   = 400     // cap one finding's joined comment text
	maxCodeLen   = 160     // cap the code line kept as context
)

// Sweep walks o.Root and returns every marker it finds, sorted by file then
// line so two runs on an unchanged tree produce the same order.
func Sweep(o Opts) ([]Finding, error) {
	root := o.Root
	if root == "" {
		root = "."
	}
	markers := o.Markers
	if len(markers) == 0 {
		markers = DefaultMarkers
	}
	exclude := o.Exclude
	if exclude == nil {
		exclude = DefaultExclude
	}
	skipDir := make(map[string]bool, len(exclude))
	for _, d := range exclude {
		if d = strings.TrimSpace(d); d != "" {
			skipDir[d] = true
		}
	}

	var found []Finding
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // an unreadable entry is skipped, never fatal
		}
		if d.IsDir() {
			if path != root && skipDir[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() || docExt[strings.ToLower(filepath.Ext(path))] {
			return nil
		}
		if fi, statErr := d.Info(); statErr == nil && fi.Size() > maxFileBytes {
			return nil
		}
		b, readErr := os.ReadFile(path)
		if readErr != nil || isBinary(b) {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		found = append(found, scanFile(filepath.ToSlash(rel), b, markers, openersFor(path))...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(found, func(i, j int) bool {
		if found[i].File != found[j].File {
			return found[i].File < found[j].File
		}
		return found[i].Line < found[j].Line
	})
	return found, nil
}

// isBinary reports whether the head of b holds a NUL byte.
func isBinary(b []byte) bool {
	head := b
	if len(head) > sniffBytes {
		head = head[:sniffBytes]
	}
	return bytes.IndexByte(head, 0) >= 0
}

// scanFile finds every marker in one file's bytes.
func scanFile(rel string, b []byte, markers []string, ops []string) []Finding {
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	lcs := lexFile(lines, ops)
	var out []Finding
	for i := 0; i < len(lines); i++ {
		marker, group, rest, ok := markerHit(lines[i], lcs[i], markers)
		if !ok {
			continue
		}
		opener, at := lcs[i].opener, lcs[i].at
		s := span{from: i, to: i, opener: opener, at: at}
		text := rest
		closer, isBlock := blockClosers[opener]
		switch {
		case isBlock && strings.Contains(lines[i][at+len(opener):], closer):
			// A block comment that opens and closes on this line: its words are the
			// finding, its closer is not.
			body := lines[i][at+len(opener):]
			text, s.to = lineRun(lines, lcs, i, markers, blockLine(body[:strings.Index(body, closer)]))
		case isBlock:
			// A block comment that closes further down is one finding, from its
			// opener to its closer, whatever the lines between look like.
			text, s.to, s.closeAt = blockSpan(lines, i, at+len(opener), closer, rest)
			s.unclosed = s.closeAt == 0
		default:
			text, s.to = lineRun(lines, lcs, i, markers, rest)
		}
		plan := planStrip(lines, s)
		// A comment that trails code sits on the line it is about; a standalone
		// comment is about the first code line under it.
		code := codeUnder(lines, lcs, s.to+1)
		if len(plan.cuts) > 0 && plan.cuts[0].line == i+1 {
			code = cutSpan(lines[i], plan.cuts[0].from, plan.cuts[0].to)
		}
		out = append(out, Finding{
			File:   rel,
			Line:   i + 1,
			Marker: marker,
			Group:  group,
			Text:   Normalize(text, maxTextLen),
			Code:   Normalize(code, maxCodeLen),
			strip:  plan,
		})
		i = s.to
	}
	return out
}

// lineRun joins the comment lines under the marker's own line. Consecutive comment
// lines continue the same finding, so a three-line TODO is one task, not three.
func lineRun(lines []string, lcs []lineComment, from int, markers []string, first string) (text string, to int) {
	text, to = first, from
	for n := from + 1; n < len(lines); n++ {
		cont, ok := continuationOnLine(lines[n], lcs[n], markers)
		if !ok {
			break
		}
		if cont != "" {
			text += " " + cont
		}
		to = n
	}
	return text, to
}

// blockSpan joins a block comment that closes on a later line: every line after the
// marker's joins the text with its decoration trimmed, up to the line that holds the
// closer, and the closer itself is dropped. One block comment is one finding, so a
// second marker word inside the span does not split it.
//
// A block that closes nowhere below returns the marker line alone and closeAt 0, so
// the strip leaves the comment in a file it cannot read to the end.
func blockSpan(lines []string, from, bodyAt int, closer, first string) (text string, to, closeAt int) {
	text = first
	for n := from; n < len(lines); n++ {
		start := 0
		if n == from {
			start = bodyAt
		}
		k := strings.Index(lines[n][start:], closer)
		if k < 0 {
			if n > from {
				if part := blockLine(lines[n]); part != "" {
					text += " " + part
				}
			}
			continue
		}
		if n > from {
			if part := blockLine(lines[n][:start+k]); part != "" {
				text += " " + part
			}
		}
		return text, n, start + k + len(closer)
	}
	return first, from, 0
}

// blockLine trims one line of a block comment down to its words: the decoration a
// block puts at the start of a line ("* ") and the space around it go.
func blockLine(line string) string {
	return strings.TrimSpace(trimDecoration(strings.TrimSpace(line)))
}

// planStrip works out how to delete one marker comment without touching a byte
// of code.
//
// Three shapes are removable, and nothing else is:
//   - a standalone comment line (plus its continuation lines): the lines go.
//   - a comment trailing code on one line: the comment span is cut, the code stays.
//   - a block comment that closes on a later line: the whole span goes, and code
//     before the opener is cut free of it first.
//
// A marker inside a block comment that another line opened is left alone, because
// deleting one of its lines can leave a dangling opener or an empty comment. So is a
// block comment that never closes (that file is broken and the sweep will not guess)
// and one whose closer shares a line with code.
func planStrip(lines []string, s span) stripPlan {
	if s.opener == "*" {
		return stripPlan{reason: "inside a multi-line block comment"}
	}
	if s.unclosed {
		return stripPlan{reason: "opens a block comment that does not close"}
	}
	line := lines[s.from]
	standalone := strings.TrimSpace(line[:s.at]) == ""

	if s.closeAt > 0 { // a block comment that closes on a later line
		if strings.TrimSpace(lines[s.to][s.closeAt:]) != "" {
			// Cutting the closer off this line would take the code's indentation
			// with it, so the comment stays and the next sweep sees it again.
			return stripPlan{reason: "code follows the closing " + blockClosers[s.opener]}
		}
		p := stripPlan{ok: true}
		first := s.from + 1 // 1-based
		if !standalone {
			p.cuts = append(p.cuts, lineCut{line: first, from: s.at, to: len(line)})
			first++
		}
		if first <= s.to+1 {
			p.fromLine, p.toLine = first, s.to+1
		}
		return p
	}

	end := len(line) // a line comment runs to the end of the line
	if closer, isBlock := blockClosers[s.opener]; isBlock {
		k := strings.Index(line[s.at+len(s.opener):], closer)
		if k < 0 {
			return stripPlan{reason: "opens a block comment that does not close on the same line"}
		}
		end = s.at + len(s.opener) + k + len(closer)
	}

	if standalone && end >= len(strings.TrimRight(line, " \t")) {
		return stripPlan{ok: true, fromLine: s.from + 1, toLine: s.to + 1}
	}
	p := stripPlan{ok: true, cuts: []lineCut{{line: s.from + 1, from: s.at, to: end}}}
	if s.to > s.from { // continuation lines below the code line are whole-line deletes
		p.fromLine, p.toLine = s.from+2, s.to+1
	}
	return p
}

// blockClosers maps a block-comment opener to the token that ends it.
var blockClosers = map[string]string{"/*": "*/", "<!--": "-->"}

// markerHit reports the marker that opens the comment on this line, its group tag,
// and the comment text from the marker on. The marker must be the comment's first
// word: a sentence that merely mentions ARCDLC is prose about markers, not a marker.
func markerHit(line string, lc lineComment, markers []string) (marker, group, rest string, ok bool) {
	if lc.at < 0 {
		return "", "", "", false
	}
	body := trimDecoration(commentBody(line, lc))
	for _, m := range markers {
		if m == "" || !strings.HasPrefix(body, m) {
			continue
		}
		after := body[len(m):]
		if after != "" && isWordByte(after[0]) {
			continue // ARCDLCLIST is a name, not a marker
		}
		return m, groupTag(after), strings.TrimSpace(body), true
	}
	return "", "", "", false
}

// groupTag reads the group tag off a marker: a ":" straight after the marker word,
// then a run of tag bytes that ends at a space or at the end of the line. The tag is
// folded to upper case, so ":g1" and ":G1" are the same group.
//
// Two shapes look like a tag and are not one. A tag with no letter in it
// ("// TODO:01 fix this") would claim the auto-numbered block ID TODO-CMT-01, so it
// reads as a plain marker. So does a tag the text runs straight into
// ("// TODO:refactor(later) x"), which is prose, not a label.
func groupTag(after string) string {
	if !strings.HasPrefix(after, ":") {
		return ""
	}
	end, letter := 1, false
	for ; end < len(after) && isTagByte(after[end]); end++ {
		if !isDigitByte(after[end]) && after[end] != '_' && after[end] != '-' {
			letter = true
		}
	}
	if end == 1 || !letter {
		return ""
	}
	if end < len(after) && after[end] != ' ' && after[end] != '\t' {
		return ""
	}
	return strings.ToUpper(after[1:end])
}

// IgnoredTag returns the tag-shaped text the sweep refused on this finding, or ""
// when there was nothing to refuse. The sweep reports it, so a tag that quietly
// became a plain marker is visible to the engineer who typed it.
func IgnoredTag(f Finding) string {
	if f.Group != "" || !strings.HasPrefix(f.Text, f.Marker) {
		return ""
	}
	after := f.Text[len(f.Marker):]
	if !strings.HasPrefix(after, ":") {
		return ""
	}
	end := 1
	for ; end < len(after) && isTagByte(after[end]); end++ {
	}
	if end == 1 || (end < len(after) && after[end] != ' ') {
		return ""
	}
	return after[1:end]
}

// trimDecoration removes the padding a doc comment puts between the opener and
// the text: "///", "//!", "#!", "* ", "<!--".
func trimDecoration(body string) string {
	return strings.TrimLeft(body, "/*!<-= \t")
}

// continuationOnLine reports whether the line is a plain comment line that
// continues the marker above it. A line that starts its own marker ends the
// previous finding, and so does anything that is not a comment.
func continuationOnLine(line string, lc lineComment, markers []string) (string, bool) {
	if !onlyComment(line, lc) {
		return "", false
	}
	body := strings.TrimSpace(commentBody(line, lc))
	if body == "" || body == "/" { // an empty comment or a closing */ ends the run
		return "", true
	}
	if m, _, _, isNew := markerHit(line, lc, markers); isNew && m != "" {
		return "", false // a second marker starts its own finding
	}
	return body, true
}

// lineComment says where a comment opens on one line, after the sweep has followed
// the string literals of the file. at is -1 when the line carries no comment.
type lineComment struct {
	at     int
	opener string
}

// commentBody returns the text after the comment opener on this line.
func commentBody(line string, lc lineComment) string {
	return strings.TrimLeft(line[lc.at+len(lc.opener):], " \t")
}

// onlyComment reports whether the line holds nothing but a comment, so that it can
// continue the marker above it and is not code a marker sits on.
func onlyComment(line string, lc lineComment) bool {
	return lc.at >= 0 && strings.TrimSpace(line[:lc.at]) == ""
}

// lexFile works out, line by line, where a comment opens, following the string
// literals of the file as it goes.
//
// This is what keeps the sweep off constants. A marker counts only inside a comment,
// and "// ARCDLC ..." written inside a string is text in that string. The single-line
// case needs no state, but a Go or JavaScript raw string and a Python triple quote run
// past the end of their line, so the pending delimiter is carried to the next one.
//
// It is a scanner, not a parser, and it is blunt in one direction on purpose: when it
// cannot tell, it reports no comment. A missed marker costs one sweep. A marker read
// out of a constant costs an edit to code nobody asked it to touch.
func lexFile(lines []string, ops []string) []lineComment {
	out := make([]lineComment, len(lines))
	open := ""
	for i, line := range lines {
		at, opener, still := lexLine(line, ops, open)
		out[i] = lineComment{at: at, opener: opener}
		open = still
	}
	return out
}

// multiQuotes are the string delimiters that can run past the end of a line: a Go or
// JavaScript raw string, and a Python or Kotlin triple quote.
var multiQuotes = []string{`"""`, "'''", "`"}

// lexLine walks one line and returns where its comment opens, plus the multi-line
// string delimiter still waiting to close when the line ends. open is that delimiter
// carried from the line before, "" when nothing is pending.
//
// The walk stops at the comment opener, so a quote in comment text is never read as a
// string delimiter: a lone backtick in prose cannot swallow the lines below it.
func lexLine(line string, ops []string, open string) (at int, opener, still string) {
	i := 0
	if open != "" {
		k := strings.Index(line, open)
		if k < 0 {
			return -1, "", open // the whole line sits inside the literal
		}
		i = k + len(open)
	}
	for i < len(line) {
		if op, ok := openerAt(line, i, ops); ok {
			return i, op, ""
		}
		if q := quoteAt(line, i); q != "" {
			k := closeQuote(line, i+len(q), q)
			if k < 0 {
				if isMultiQuote(q) {
					return -1, "", q // the literal runs on into the next line
				}
				return -1, "", "" // an unterminated quote ends with its line
			}
			i = k + len(q)
			continue
		}
		i++
	}
	return -1, "", ""
}

// openerAt reports the comment opener that starts at this column, the longest one
// when two of them match. A bare "*" only opens a comment at the start of a line: it
// is a multiplication sign anywhere else.
func openerAt(line string, i int, ops []string) (string, bool) {
	best := ""
	for _, op := range ops {
		if len(op) <= len(best) || !strings.HasPrefix(line[i:], op) {
			continue
		}
		if op == "*" && strings.TrimSpace(line[:i]) != "" {
			continue
		}
		best = op
	}
	return best, best != ""
}

// quoteAt reports the string delimiter that starts at this column.
func quoteAt(line string, i int) string {
	for _, q := range multiQuotes {
		if strings.HasPrefix(line[i:], q) {
			return q
		}
	}
	if line[i] == '"' || line[i] == '\'' {
		return line[i : i+1]
	}
	return ""
}

func isMultiQuote(q string) bool {
	for _, m := range multiQuotes {
		if q == m {
			return true
		}
	}
	return false
}

// closeQuote returns the column where the delimiter closes, or -1 when it does not
// close on this line. A backslash escapes the byte after it.
func closeQuote(line string, from int, q string) int {
	for i := from; i < len(line); i++ {
		if line[i] == '\\' {
			i++
			continue
		}
		if strings.HasPrefix(line[i:], q) {
			return i
		}
	}
	return -1
}

func isTagByte(c byte) bool {
	return c == '_' || c == '-' || isDigitByte(c) || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigitByte(c byte) bool { return c >= '0' && c <= '9' }

func isWordByte(c byte) bool {
	return c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// codeUnder returns the first line at or after idx that is neither blank nor a
// comment, trimmed and capped. It is the code the marker sits on.
func codeUnder(lines []string, lcs []lineComment, idx int) string {
	for ; idx < len(lines); idx++ {
		trimmed := strings.TrimSpace(lines[idx])
		if trimmed == "" || onlyComment(lines[idx], lcs[idx]) {
			continue
		}
		return Normalize(trimmed, maxCodeLen)
	}
	return ""
}

// Normalize collapses whitespace, drops backticks (they would break the code
// span the register writes the text into), and caps the length. Both writing
// and reading a register go through it, so a finding's fingerprint is stable.
func Normalize(s string, max int) string {
	s = strings.Join(strings.Fields(strings.ReplaceAll(s, "`", "'")), " ")
	if max > 0 && len(s) > max {
		s = strings.TrimSpace(s[:max])
	}
	return s
}

// Fingerprint identifies a finding across runs. Line numbers drift on every
// edit, so identity is the file plus the normalized comment text.
func (f Finding) Fingerprint() string {
	return f.File + "\x00" + Normalize(f.Text, maxTextLen)
}
