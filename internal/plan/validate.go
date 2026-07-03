package plan

import (
	"fmt"
	"strings"
)

// Severity is the level of a validation finding.
type Severity string

const (
	SevError   Severity = "error"
	SevWarning Severity = "warning"
)

// Finding is one validation problem. JSON tags match the arctool spec schema.
type Finding struct {
	Severity Severity `json:"severity"`
	ID       string   `json:"id,omitempty"`
	Rule     string   `json:"rule"`
	Message  string   `json:"message"`
	Line     int      `json:"line"`
}

// ValidateOpts selects which rule set runs.
type ValidateOpts struct {
	Strict            bool // enable the strict-only rules
	RequireAcceptance bool // enable the (otherwise off) Acceptance-present rule
}

var validSrc = map[SourceStatus]bool{SrcMISSING: true, SrcPARTIAL: true, SrcDRIFT: true}

// Validate runs the format-contract checks and returns findings in document
// order. The returned slice is always non-nil (so JSON emits [] not null).
// Failure/exit decisions are the caller's: any error fails; warnings fail only
// under strict / warn-as-error.
func (p *Plan) Validate(opts ValidateOpts) []Finding {
	f := []Finding{}
	seen := map[string]int{} // id -> first-seen heading line

	label := func(t Task) string {
		if t.ID != "" {
			return fmt.Sprintf("%q", t.ID)
		}
		return fmt.Sprintf("(line %d)", t.Line)
	}

	for _, t := range p.Tasks {
		// 5 — heading shape
		if !t.HeadingOK || t.ID == "" || t.Title == "" {
			f = append(f, Finding{SevError, t.ID, "malformed-heading",
				fmt.Sprintf(`malformed heading (want "### <ID> (<SRC>): <Title>") at line %d`, t.Line), t.Line})
		}

		// 3 — unique IDs
		if t.ID != "" {
			if first, ok := seen[t.ID]; ok {
				f = append(f, Finding{SevError, t.ID, "duplicate-id",
					fmt.Sprintf("duplicate task ID %q (first seen line %d)", t.ID, first), t.Line})
			} else {
				seen[t.ID] = t.Line
			}
		}

		// 1 — status present (the flagship "silently skipped" check)
		switch {
		case !t.HasStatus:
			f = append(f, Finding{SevError, t.ID, "missing-status",
				fmt.Sprintf(`task %s has no "- Status:" line (runner will silently skip it)`, label(t)), t.Line})
		case t.Status == StatusNone:
			// 2 — status token recognized
			f = append(f, Finding{SevError, t.ID, "invalid-status",
				fmt.Sprintf("task %s has invalid status %q (want TODO/TAKEN/DONE/BLOCKED)", label(t), t.StatusRaw), t.StatusLine})
		case t.StatusRaw != strings.ToUpper(t.StatusRaw):
			// 7 — status casing (warning)
			f = append(f, Finding{SevWarning, t.ID, "status-casing",
				fmt.Sprintf("task %s status %q should be uppercase %q", label(t), t.StatusRaw, strings.ToUpper(t.StatusRaw)), t.StatusLine})
		}

		// 4 — required keys present
		var missing []string
		for _, k := range []struct {
			name string
			has  bool
		}{
			{"WHAT", t.HasWhat}, {"WHERE", t.HasWhere}, {"WHY", t.HasWhy},
			{"References", t.HasRefs}, {"Status", t.HasStatus},
		} {
			if !k.has {
				missing = append(missing, k.name)
			}
		}
		if len(missing) > 0 {
			f = append(f, Finding{SevError, t.ID, "missing-key",
				fmt.Sprintf("task %s missing required key(s): %s", label(t), strings.Join(missing, ", ")), t.Line})
		}

		// 6 — source-status tag known (warning)
		if t.SourceStatus != "" && !validSrc[t.SourceStatus] {
			f = append(f, Finding{SevWarning, t.ID, "unknown-source-status",
				fmt.Sprintf("task %s has unknown source-status %q", label(t), t.SourceStatus), t.Line})
		}

		// 11 — Acceptance present (off unless requested)
		if opts.RequireAcceptance && !t.HasAcceptance {
			f = append(f, Finding{SevError, t.ID, "missing-acceptance",
				fmt.Sprintf(`task %s has no "- Acceptance:" criteria`, label(t)), t.Line})
		}

		// strict-only per-task rules
		if opts.Strict {
			if t.HasRefs && len(t.References) == 0 {
				f = append(f, Finding{SevError, t.ID, "empty-references",
					fmt.Sprintf("task %s has empty References", label(t)), t.Line})
			}
			if t.HasWhere && !looksLikePath(t.Where) {
				f = append(f, Finding{SevError, t.ID, "where-no-path",
					fmt.Sprintf("task %s WHERE lists no concrete file/module", label(t)), t.Line})
			}
			if t.HasAcceptance && strings.TrimSpace(t.Acceptance) == "" {
				f = append(f, Finding{SevError, t.ID, "empty-acceptance",
					fmt.Sprintf("task %s has an empty Acceptance section", label(t)), t.Line})
			}
		}
	}

	// 10 — dependency order (strict): a task must not reference a later task ID.
	if opts.Strict {
		pos := map[string]int{}
		for i, t := range p.Tasks {
			if t.ID != "" {
				if _, ok := pos[t.ID]; !ok {
					pos[t.ID] = i
				}
			}
		}
		for i, t := range p.Tasks {
			for _, r := range t.References {
				if j, ok := pos[r]; ok && j > i {
					f = append(f, Finding{SevError, t.ID, "dep-order",
						fmt.Sprintf("task %q references later task %q (queue runs top-to-bottom)", t.ID, r), t.Line})
				}
			}
		}
	}

	return f
}

// Counts returns the number of error and warning findings.
func Counts(findings []Finding) (errs, warns int) {
	for _, f := range findings {
		if f.Severity == SevError {
			errs++
		} else {
			warns++
		}
	}
	return
}
