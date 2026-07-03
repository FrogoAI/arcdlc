package plan

import (
	"regexp"
	"strings"
)

// FirstTODO returns a pointer to the first task whose status is TODO (document
// order), or nil when the queue has no pending task.
func (p *Plan) FirstTODO() *Task {
	for i := range p.Tasks {
		if p.Tasks[i].Status == StatusTODO {
			return &p.Tasks[i]
		}
	}
	return nil
}

// ByID returns every task with the given ID. More than one result means the
// plan has a duplicate ID (a validate error) — callers should refuse to guess.
func (p *Plan) ByID(id string) []*Task {
	var out []*Task
	for i := range p.Tasks {
		if p.Tasks[i].ID == id {
			out = append(out, &p.Tasks[i])
		}
	}
	return out
}

// Raw returns the exact bytes of a task's block, unmodified.
func (p *Plan) Raw(t *Task) []byte {
	return p.Bytes[t.RawStart:t.RawEnd]
}

// StatusCounts tallies tasks by status.
func (p *Plan) StatusCounts() map[Status]int {
	m := map[Status]int{}
	for _, t := range p.Tasks {
		m[t.Status]++
	}
	return m
}

// Layer is one parsed WHERE line: a layer/section label and its target paths.
type Layer struct {
	Layer   string   `json:"layer"`
	Targets []string `json:"targets"`
}

var (
	reWhereLayer   = regexp.MustCompile("^\\s*Layer\\s+`([^`]+)`\\s*:\\s*(.+)$")
	reWhereGeneric = regexp.MustCompile(`^\s*([A-Za-z][\w/ .-]*?)\s*:\s*(.+)$`)
)

// ParseWhereLayers derives the layer→targets breakdown from a WHERE body. It is
// best-effort and used only for JSON output: it returns an empty slice when the
// block does not use the "Layer `x`: ..." / "Tests: ..." convention.
func ParseWhereLayers(where string) []Layer {
	out := []Layer{}
	for _, ln := range strings.Split(where, "\n") {
		var name, rest string
		if m := reWhereLayer.FindStringSubmatch(ln); m != nil {
			name, rest = m[1], m[2]
		} else if m := reWhereGeneric.FindStringSubmatch(ln); m != nil {
			name, rest = strings.ToLower(strings.TrimSpace(m[1])), m[2]
		} else {
			continue
		}
		if targets := splitTargets(rest); len(targets) > 0 {
			out = append(out, Layer{Layer: name, Targets: targets})
		}
	}
	return out
}

// splitTargets cleans a comma-separated target list (whitespace, a trailing
// period, and wrapping backticks removed), like parseRefs but for WHERE lines.
func splitTargets(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		p = strings.TrimSuffix(p, ".")
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "`")
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
