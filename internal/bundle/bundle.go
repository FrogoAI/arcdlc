// Package bundle checks the skill bundle against the rules AGENTS.md states in
// prose: every skill carries the shared blocks, every path a skill names
// resolves in both install layouts, and no skill reaches outside the bundle.
//
// The checks exist because a wrong instruction to an agent is worse than a
// missing one, and nothing else in the pipeline reads a SKILL.md.
package bundle

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Skills is the bundle's sub-skill list. It must match SUBSKILLS in install.sh
// and the loops in .github/workflows/ci.yml.
var Skills = []string{
	"aic", "archive", "assist", "close", "examinate", "execute",
	"grilling", "init", "plan", "plan-human", "policy", "remove",
}

// SharedBlocks are the sections every SKILL.md carries verbatim.
var SharedBlocks = []string{
	"## Talk simple, write like a human",
	"## Judge by the four virtues",
}

// Problem is one rule violation, named by the file it was found in.
type Problem struct {
	File string
	Line int // 0 when the problem is about the file as a whole
	Msg  string
}

func (p Problem) String() string {
	if p.Line > 0 {
		return fmt.Sprintf("%s:%d: %s", p.File, p.Line, p.Msg)
	}
	return fmt.Sprintf("%s: %s", p.File, p.Msg)
}

// Paths are only recognised inside backticks or markdown link parentheses.
// Those delimiters are what make a filename with a space in it (this bundle has
// several) unambiguous, which whitespace alone cannot do.
var (
	crossSkillTick  = regexp.MustCompile("`\\.\\./([A-Za-z0-9_-]+)/([^`]+)`")
	crossSkillParen = regexp.MustCompile(`\(\.\./([A-Za-z0-9_-]+)/([^)]+)\)`)
	localRefTick    = regexp.MustCompile("`(references/[^`]+)`")
	escapesBundle   = regexp.MustCompile("[`(]\\.\\./\\.\\./")
)

// Check reads the skill bundle under root and returns every problem found,
// sorted by file. An empty slice means the bundle is consistent.
func Check(root string) ([]Problem, error) {
	var probs []Problem
	skillsDir := filepath.Join(root, "skills")

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("read skills dir: %w", err)
	}
	onDisk := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() {
			onDisk[e.Name()] = true
		}
	}
	for _, s := range Skills {
		if !onDisk[s] {
			probs = append(probs, Problem{File: "skills/" + s, Msg: "listed in Skills but missing on disk"})
		}
		delete(onDisk, s)
	}
	for extra := range onDisk {
		probs = append(probs, Problem{File: "skills/" + extra, Msg: "on disk but not in Skills; add it to install.sh SUBSKILLS and ci.yml too"})
	}

	// The reference copy of each shared block comes from the first skill.
	refBlocks := map[string]string{}

	for _, name := range Skills {
		rel := filepath.Join("skills", name, "SKILL.md")
		raw, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			probs = append(probs, Problem{File: rel, Msg: "cannot read: " + err.Error()})
			continue
		}
		text := string(raw)
		probs = append(probs, checkFrontmatter(rel, name, text)...)
		probs = append(probs, checkBlocks(rel, text, refBlocks)...)
		probs = append(probs, checkPaths(root, rel, name, text)...)
	}

	sort.Slice(probs, func(i, j int) bool {
		if probs[i].File != probs[j].File {
			return probs[i].File < probs[j].File
		}
		return probs[i].Line < probs[j].Line
	})
	return probs, nil
}

// checkFrontmatter enforces the naming and trigger rules: a skill is invoked
// either as /arcdlc:<name> or as arcdlc-<name>, and its description is the only
// text an agent sees when deciding whether to fire it, so both forms must
// appear there.
func checkFrontmatter(rel, name, text string) []Problem {
	var probs []Problem
	if !strings.HasPrefix(text, "---\n") {
		return []Problem{{File: rel, Msg: "no YAML frontmatter"}}
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return []Problem{{File: rel, Msg: "unterminated YAML frontmatter"}}
	}
	fm := text[4 : 4+end]

	want := "name: arcdlc-" + name
	if !strings.Contains(fm, want) {
		probs = append(probs, Problem{File: rel, Msg: "frontmatter must declare " + want})
	}
	desc := ""
	for _, line := range strings.Split(fm, "\n") {
		if strings.HasPrefix(line, "description:") {
			desc = line
		}
	}
	switch {
	case desc == "":
		probs = append(probs, Problem{File: rel, Msg: "frontmatter has no description"})
	default:
		for _, trigger := range []string{"/arcdlc:" + name, "arcdlc-" + name} {
			if !strings.Contains(desc, trigger) {
				probs = append(probs, Problem{File: rel,
					Msg: "description must name the trigger " + trigger + "; a skill only fires on words in its description"})
			}
		}
	}
	return probs
}

// checkBlocks enforces that every skill carries each shared block and that all
// copies are byte identical.
func checkBlocks(rel, text string, ref map[string]string) []Problem {
	var probs []Problem
	for _, heading := range SharedBlocks {
		body, ok := section(text, heading)
		if !ok {
			probs = append(probs, Problem{File: rel, Msg: "missing the block " + heading})
			continue
		}
		if first, seen := ref[heading]; !seen {
			ref[heading] = body
		} else if first != body {
			probs = append(probs, Problem{File: rel,
				Msg: "the block " + heading + " differs from the other skills; all copies must be byte identical"})
		}
	}
	return probs
}

// section returns the text of a "## " section, from its heading to the next
// "## " heading or end of file.
func section(text, heading string) (string, bool) {
	i := strings.Index(text, "\n"+heading+"\n")
	if i < 0 {
		if !strings.HasPrefix(text, heading+"\n") {
			return "", false
		}
		i = -1
	}
	rest := text[i+1+len(heading):]
	if j := strings.Index(rest, "\n## "); j >= 0 {
		return rest[:j], true
	}
	return rest, true
}

// checkPaths enforces the install-agnostic rule. A skill may name a file it
// owns under references/, or a sibling skill's file, and a sibling reference
// must spell both layouts: ../<name>/ for the plugin layout and
// ../arcdlc-<name>/ for a flat install. Nothing may reach outside skills/,
// because only skills/ is copied onto a user machine.
func checkPaths(root, rel, name, text string) []Problem {
	var probs []Problem
	dir := filepath.Join(root, "skills", name)

	for _, m := range lineMatches(text, localRefTick) {
		p := m.groups[0]
		if strings.ContainsAny(p, "<*…") {
			continue // a placeholder, not a path
		}
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			probs = append(probs, Problem{File: rel, Line: m.line, Msg: "names " + p + ", which does not exist"})
		}
	}

	for _, m := range lineMatches(text, escapesBundle) {
		probs = append(probs, Problem{File: rel, Line: m.line,
			Msg: "path escapes skills/ (" + m.text + "); only skills/ is copied to a user machine"})
	}

	cross := append(lineMatches(text, crossSkillTick), lineMatches(text, crossSkillParen)...)
	for _, m := range cross {
		target, rest := m.groups[0], m.groups[1]
		plain := strings.TrimPrefix(target, "arcdlc-")
		if !contains(Skills, plain) {
			probs = append(probs, Problem{File: rel, Line: m.line,
				Msg: "references unknown skill " + target})
			continue
		}
		if strings.ContainsAny(rest, "<*…") {
			continue
		}
		// The plugin-layout half must resolve on disk.
		if _, err := os.Stat(filepath.Join(root, "skills", plain, rest)); err != nil {
			probs = append(probs, Problem{File: rel, Line: m.line,
				Msg: "names ../" + plain + "/" + rest + ", which does not exist"})
		}
		// Both layouts must be spelled somewhere in the file.
		other := "../arcdlc-" + plain + "/"
		if strings.HasPrefix(target, "arcdlc-") {
			other = "../" + plain + "/"
		}
		if !strings.Contains(text, other) {
			probs = append(probs, Problem{File: rel, Line: m.line,
				Msg: "cross-skill path " + m.text + " has no " + other + " counterpart; both install layouts must be spelled"})
		}
	}
	return probs
}

type match struct {
	line   int
	text   string
	groups []string
}

func lineMatches(text string, re *regexp.Regexp) []match {
	var out []match
	for i, line := range strings.Split(text, "\n") {
		for _, g := range re.FindAllStringSubmatch(line, -1) {
			out = append(out, match{line: i + 1, text: g[0], groups: g[1:]})
		}
	}
	return out
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// reReviewed matches the freshness marker every reference document carries, in
// either the bold or the metadata-list style the documents already use.
var reReviewed = regexp.MustCompile(`(?m)^(?:\*\*Reviewed\*\*:|- Reviewed:)\s*(\d{4})-(\d{2})-(\d{2})\s*$`)

// Reference is one bundled reference document and the date it was last read
// end to end.
type Reference struct {
	Path     string
	Reviewed time.Time
}

// References walks every skill's references/ directory and returns each
// document with its review date.
//
// The marker exists because a reference rots silently. Go Best Practice.md sat
// in this bundle describing pre-1.18 Go, with nothing on its face to say so,
// until someone read all 428 lines to find out. A date turns that from
// archaeology into a glance.
func References(root string) ([]Reference, []Problem, error) {
	var refs []Reference
	var probs []Problem

	for _, skill := range Skills {
		dir := filepath.Join(root, "skills", skill, "references")
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("read %s: %w", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			rel := filepath.Join("skills", skill, "references", e.Name())
			raw, err := os.ReadFile(filepath.Join(root, rel))
			if err != nil {
				probs = append(probs, Problem{File: rel, Msg: "cannot read: " + err.Error()})
				continue
			}
			m := reReviewed.FindSubmatch(raw)
			if m == nil {
				probs = append(probs, Problem{File: rel,
					Msg: "no `**Reviewed**: YYYY-MM-DD` marker; a reference with no date rots unnoticed"})
				continue
			}
			t, err := time.Parse("2006-01-02", string(m[1])+"-"+string(m[2])+"-"+string(m[3]))
			if err != nil {
				probs = append(probs, Problem{File: rel, Msg: "unparseable Reviewed date: " + err.Error()})
				continue
			}
			refs = append(refs, Reference{Path: rel, Reviewed: t})
		}
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Reviewed.Before(refs[j].Reviewed) })
	return refs, probs, nil
}
