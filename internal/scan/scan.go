// Package scan sweeps a source tree for code comment markers (TODO, FIXME and
// friends) and renders them into the comment register, docs/aics/<slug>/comments.md.
//
// The split of work is deliberate: this package finds evidence (file, line, the
// marker text, the code line under it) and writes only the lines it owns. The
// judgement (title, WHAT, HOW, WHY, Acceptance, verdict) belongs to the
// /arcdlc:assist skill, which fills it in after it grills the engineer.
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
	File   string `json:"file"`   // path relative to the sweep root, slash separated
	Line   int    `json:"line"`   // 1-based line of the marker
	Marker string `json:"marker"` // the marker word, e.g. "TODO"
	Text   string `json:"text"`   // the comment from the marker on, continuation lines joined
	Code   string `json:"code"`   // first code line under the comment, "" when there is none
}

// Opts configures a sweep.
type Opts struct {
	Root    string   // directory to walk; "" means "."
	Markers []string // marker words to look for; empty means DefaultMarkers
	Exclude []string // directory names to skip; nil means DefaultExclude
}

// DefaultMarkers is what a sweep looks for when the caller names nothing.
var DefaultMarkers = []string{"TODO"}

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
	var out []Finding
	for i := 0; i < len(lines); i++ {
		marker, rest, ok := markerOnLine(lines[i], markers, ops)
		if !ok {
			continue
		}
		text := rest
		// Consecutive comment lines below the marker continue the same finding,
		// so a three-line TODO is one task, not three.
		j := i + 1
		for ; j < len(lines); j++ {
			cont, contOK := continuationOnLine(lines[j], markers, ops)
			if !contOK {
				break
			}
			if cont != "" {
				text += " " + cont
			}
		}
		out = append(out, Finding{
			File:   rel,
			Line:   i + 1,
			Marker: marker,
			Text:   Normalize(text, maxTextLen),
			Code:   codeUnder(lines, j, ops),
		})
		i = j - 1
	}
	return out
}

// markerOnLine reports the marker that opens the comment on this line, and
// returns the comment text from the marker on. The marker must be the comment's
// first word: a sentence that merely mentions TODO is prose about markers, not a
// marker.
func markerOnLine(line string, markers []string, ops []string) (marker, rest string, ok bool) {
	body, found := commentBody(line, ops)
	if !found {
		return "", "", false
	}
	body = trimDecoration(body)
	for _, m := range markers {
		if m == "" || !strings.HasPrefix(body, m) {
			continue
		}
		after := body[len(m):]
		if after != "" && isWordByte(after[0]) {
			continue // TODOLIST is a name, not a marker
		}
		return m, strings.TrimSpace(body), true
	}
	return "", "", false
}

// trimDecoration removes the padding a doc comment puts between the opener and
// the text: "///", "//!", "#!", "* ", "<!--".
func trimDecoration(body string) string {
	return strings.TrimLeft(body, "/*!<-= \t")
}

// continuationOnLine reports whether the line is a plain comment line that
// continues the marker above it. A line that starts its own marker ends the
// previous finding, and so does anything that is not a comment.
func continuationOnLine(line string, markers []string, ops []string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return "", false
	}
	if !startsWithOpener(trimmed, ops) {
		return "", false
	}
	body, _ := commentBody(line, ops)
	body = strings.TrimSpace(body)
	if body == "" || body == "/" { // an empty comment or a closing */ ends the run
		return "", true
	}
	if _, _, isNew := markerOnLine(line, markers, ops); isNew {
		return "", false // a second marker starts its own finding
	}
	return body, true
}

// commentBody returns the text after the first comment opener on the line.
func commentBody(line string, ops []string) (string, bool) {
	at, width := -1, 0
	for _, op := range ops {
		i := strings.Index(line, op)
		if i < 0 {
			continue
		}
		// A bare "*" only opens a comment at the start of a line: it is a
		// multiplication sign anywhere else.
		if op == "*" && strings.TrimSpace(line[:i]) != "" {
			continue
		}
		if insideString(line, i) {
			continue // a comment quoted in a string literal, as test fixtures do
		}
		if at < 0 || i < at || (i == at && len(op) > width) {
			at, width = i, len(op)
		}
	}
	if at < 0 {
		return "", false
	}
	return strings.TrimLeft(line[at+width:], " \t"), true
}

// insideString reports whether position at sits inside a string literal that
// opens earlier on the same line. Counting unescaped quotes is enough for the
// one case that matters: source that quotes a comment, which every scanner test
// fixture does.
func insideString(line string, at int) bool {
	dquote, backtick := 0, 0
	for i := 0; i < at; i++ {
		switch line[i] {
		case '\\':
			i++ // skip the escaped byte
		case '"':
			dquote++
		case '`':
			backtick++
		}
	}
	return dquote%2 == 1 || backtick%2 == 1
}

// startsWithOpener reports whether an already-trimmed line begins a comment.
func startsWithOpener(trimmed string, ops []string) bool {
	for _, op := range ops {
		if strings.HasPrefix(trimmed, op) {
			return true
		}
	}
	return false
}

func isWordByte(c byte) bool {
	return c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// codeUnder returns the first line at or after idx that is neither blank nor a
// comment, trimmed and capped. It is the code the marker sits on.
func codeUnder(lines []string, idx int, ops []string) string {
	for ; idx < len(lines); idx++ {
		trimmed := strings.TrimSpace(lines[idx])
		if trimmed == "" || startsWithOpener(trimmed, ops) {
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
