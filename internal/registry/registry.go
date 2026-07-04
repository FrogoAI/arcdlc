// Package registry derives the initiative registry that arctool keeps in sync
// inside AGENTS.md and README.md. This file covers the deterministic parsing
// half: given an initiative folder, it finds the architecture document and
// extracts the initiative's title and one-line summary. It reads folders and
// files but never writes; the marker-block rewriting lives alongside it.
package registry

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// summaryMax bounds a fallback summary (a first-paragraph excerpt) in runes.
const summaryMax = 120

// archDocPrecedence lists the recognised architecture-document filenames in the
// order arctool prefers them when a folder holds more than one.
var archDocPrecedence = []string{"aic.md", "arc42.md", "togaf.md", "c4.md"}

// nonArchDocs are initiative files that are never the architecture document, so
// they are excluded from the alphabetical fallback in findArchDoc.
var nonArchDocs = map[string]bool{"plan.md": true, "gap.md": true, "plan-archive.md": true}

// Initiative is one entry in the registry.
type Initiative struct {
	Slug       string
	Title      string
	Summary    string
	DocRelPath string // repo-relative slash path to the arch doc; "" when there is none
}

// Load reads the initiative folder <aicsDir>/<slug> and derives its registry
// entry. Title is the first H1 of the architecture document (chosen by
// findArchDoc); Summary is the one-line blockquote under that H1, else the first
// paragraph truncated. A folder with no architecture document yields the slug as
// title and an explicit "(no architecture doc)" summary.
func Load(aicsDir, slug string) Initiative {
	dir := filepath.Join(aicsDir, slug)
	doc := findArchDoc(dir)
	if doc == "" {
		return Initiative{Slug: slug, Title: slug, Summary: "(no architecture doc)"}
	}
	content, err := os.ReadFile(filepath.Join(dir, doc))
	if err != nil {
		return Initiative{Slug: slug, Title: slug, Summary: "(no architecture doc)"}
	}
	title := parseTitle(content)
	if title == "" {
		title = slug
	}
	summary := parseSummary(content)
	if summary == "" {
		summary = "(no summary)"
	}
	return Initiative{
		Slug:       slug,
		Title:      title,
		Summary:    summary,
		DocRelPath: filepath.ToSlash(filepath.Join(aicsDir, slug, doc)),
	}
}

// findArchDoc returns the architecture-document filename in dir: the first of
// archDocPrecedence that is present, else the first *.md alphabetically that is
// not a known non-architecture file (plan/gap/plan-archive). "" when none.
func findArchDoc(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	present := map[string]bool{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			present[e.Name()] = true
		}
	}
	for _, p := range archDocPrecedence {
		if present[p] {
			return p
		}
	}
	var others []string
	for name := range present {
		if !nonArchDocs[name] {
			others = append(others, name)
		}
	}
	sort.Strings(others)
	if len(others) > 0 {
		return others[0]
	}
	return ""
}

// parseTitle returns the text of the first level-1 ATX heading ("# Title"), or "".
func parseTitle(content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(s, "# ") {
			return strings.TrimSpace(s[2:])
		}
	}
	return ""
}

// parseSummary returns the initiative summary: the one-line blockquote ("> ...")
// immediately following the first H1, or the first non-empty paragraph after it
// truncated to summaryMax runes. Returns "" when nothing follows the H1 (or there
// is no H1).
func parseSummary(content []byte) string {
	lines := strings.Split(string(content), "\n")
	i := 0
	for ; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "# ") {
			i++
			break
		}
	}
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	if i >= len(lines) {
		return ""
	}
	if first := strings.TrimSpace(lines[i]); strings.HasPrefix(first, ">") {
		return truncate(strings.TrimSpace(strings.TrimPrefix(first, ">")))
	}
	var para []string
	for ; i < len(lines); i++ {
		s := strings.TrimSpace(lines[i])
		if s == "" {
			break
		}
		para = append(para, s)
	}
	return truncate(strings.Join(para, " "))
}

// truncate shortens s to at most summaryMax runes, appending an ellipsis when cut.
func truncate(s string) string {
	r := []rune(s)
	if len(r) <= summaryMax {
		return s
	}
	return strings.TrimRight(string(r[:summaryMax]), " ") + "…"
}
