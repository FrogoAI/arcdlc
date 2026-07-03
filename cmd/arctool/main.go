// Command arctool is the deterministic companion for the ArcDLC plan
// (docs/aics/plan.md). It covers the full plan lifecycle: read commands
// (validate, next, show, list), status mutation (take, done, block, todo),
// archive, and version.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/FrogoAI/arcdlc/internal/plan"
)

const (
	version     = "0.5.0"
	defaultPlan = "docs/aics/plan.md"
)

const usage = `arctool %s — deterministic companion for the ArcDLC plan (docs/aics/plan.md)

usage:
  arctool next   [--json] [--plan PATH]                       first TODO block (exit 3 if none)
  arctool show   <id> [--json] [--plan PATH]                  one block by task ID
  arctool list   [--status TODO|TAKEN|DONE|BLOCKED] [--json] [--plan PATH]
  arctool take|done|todo <id> [--force] [--plan PATH]         flip status (TODO->TAKEN->DONE / release)
  arctool block  <id> [-m reason] [--force] [--plan PATH]     mark BLOCKED
  arctool validate [--strict] [--json] [--warn-as-error] [--require-acceptance] [--plan PATH]
                 (--strict implies --require-acceptance: every task needs an Acceptance section)
  arctool archive  [--dry-run] [--plan PATH]                  move DONE blocks to plan-archive.md
  arctool version

exit codes:
  0  ok / clean
  1  contract or guardrail failure (validation problems)
  2  usage error
  3  not found / nothing to do (e.g. next with no TODO)
  4  I/O error
  5  archive self-validation failed
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usage, version)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "next":
		os.Exit(cmdNext(os.Args[2:]))
	case "show":
		os.Exit(cmdShow(os.Args[2:]))
	case "list":
		os.Exit(cmdList(os.Args[2:]))
	case "take", "done", "block", "todo":
		os.Exit(cmdMutate(os.Args[1], os.Args[2:]))
	case "validate":
		os.Exit(cmdValidate(os.Args[2:]))
	case "archive":
		os.Exit(cmdArchive(os.Args[2:]))
	case "version", "--version", "-v":
		fmt.Printf("arctool %s\n", version)
	case "help", "--help", "-h":
		fmt.Printf(usage, version)
	default:
		fmt.Fprintf(os.Stderr, "arctool: unknown command %q\n\n", os.Args[1])
		fmt.Fprintf(os.Stderr, usage, version)
		os.Exit(2)
	}
}

// loadPlan reads and parses the plan; on read failure it reports and returns exit code 4.
func loadPlan(path string) (*plan.Plan, int) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "arctool: cannot read %s: %v\n", path, err)
		return nil, 4
	}
	return plan.Parse(b), 0
}

// splitArgs separates flags from positionals so a positional can precede flags
// (e.g. `show AIC-1 --json`), which the stdlib flag package does not allow.
// valueFlags names the flags that consume a following token (e.g. "plan").
func splitArgs(args []string, valueFlags map[string]bool) (flags, pos []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if len(a) > 1 && a[0] == '-' {
			flags = append(flags, a)
			name := strings.TrimLeft(a, "-")
			if !strings.Contains(a, "=") && valueFlags[name] && i+1 < len(args) {
				i++
				flags = append(flags, args[i])
			}
			continue
		}
		pos = append(pos, a)
	}
	return
}

// --- task rendering (shared by next and show) ---

type taskJSON struct {
	ID           string       `json:"id"`
	SourceStatus string       `json:"sourceStatus"`
	Title        string       `json:"title"`
	What         string       `json:"what"`
	Where        string       `json:"where"`
	WhereLayers  []plan.Layer `json:"whereLayers"`
	Why          string       `json:"why"`
	Acceptance   string       `json:"acceptance"`
	References   []string     `json:"references"`
	Status       string       `json:"status"`
	LineStart    int          `json:"lineStart"`
	LineEnd      int          `json:"lineEnd"`
}

func toTaskJSON(t *plan.Task) taskJSON {
	refs := t.References
	if refs == nil {
		refs = []string{}
	}
	return taskJSON{
		ID:           t.ID,
		SourceStatus: string(t.SourceStatus),
		Title:        t.Title,
		What:         t.What,
		Where:        t.Where,
		WhereLayers:  plan.ParseWhereLayers(t.Where),
		Why:          t.Why,
		Acceptance:   t.Acceptance,
		References:   refs,
		Status:       t.Status.String(),
		LineStart:    t.Line,
		LineEnd:      t.LineEnd,
	}
}

func emitTask(p *plan.Plan, t *plan.Task, asJSON bool) int {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(toTaskJSON(t)); err != nil {
			fmt.Fprintf(os.Stderr, "arctool: %v\n", err)
			return 4
		}
		return 0
	}
	// Verbatim block, trailing blank lines trimmed to a single newline.
	os.Stdout.Write(bytes.TrimRight(p.Raw(t), "\n"))
	os.Stdout.Write([]byte("\n"))
	return 0
}

// --- commands ---

func cmdNext(args []string) int {
	fs := flag.NewFlagSet("next", flag.ContinueOnError)
	planPath := fs.String("plan", defaultPlan, "path to the plan file")
	asJSON := fs.Bool("json", false, "emit the task as JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	p, code := loadPlan(*planPath)
	if code != 0 {
		return code
	}
	t := p.FirstTODO()
	if t == nil {
		if *asJSON {
			fmt.Println("null")
		} else {
			fmt.Fprintln(os.Stderr, "arctool: no pending TODO task")
		}
		return 3
	}
	return emitTask(p, t, *asJSON)
}

func cmdShow(args []string) int {
	flags, pos := splitArgs(args, map[string]bool{"plan": true})
	fs := flag.NewFlagSet("show", flag.ContinueOnError)
	planPath := fs.String("plan", defaultPlan, "path to the plan file")
	asJSON := fs.Bool("json", false, "emit the task as JSON")
	if err := fs.Parse(flags); err != nil {
		return 2
	}
	if len(pos) != 1 {
		fmt.Fprintln(os.Stderr, "usage: arctool show <id> [--json] [--plan PATH]")
		return 2
	}
	id := pos[0]
	p, code := loadPlan(*planPath)
	if code != 0 {
		return code
	}
	matches := p.ByID(id)
	switch len(matches) {
	case 0:
		if *asJSON {
			fmt.Println("null")
		} else {
			fmt.Fprintf(os.Stderr, "arctool: no task %q\n", id)
		}
		return 3
	case 1:
		return emitTask(p, matches[0], *asJSON)
	default:
		var atLines []string
		for _, m := range matches {
			atLines = append(atLines, fmt.Sprintf("%d", m.Line))
		}
		fmt.Fprintf(os.Stderr, "arctool: ambiguous id %q (%d matches at lines %s)\n",
			id, len(matches), strings.Join(atLines, ", "))
		return 3
	}
}

func cmdList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	planPath := fs.String("plan", defaultPlan, "path to the plan file")
	statusFilter := fs.String("status", "", "filter by status (TODO/TAKEN/DONE/BLOCKED)")
	asJSON := fs.Bool("json", false, "emit as JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	filter := strings.ToUpper(strings.TrimSpace(*statusFilter))
	switch filter {
	case "", "TODO", "TAKEN", "DONE", "BLOCKED":
	default:
		fmt.Fprintf(os.Stderr, "arctool: invalid --status %q (want TODO/TAKEN/DONE/BLOCKED)\n", *statusFilter)
		return 2
	}

	p, code := loadPlan(*planPath)
	if code != 0 {
		return code
	}
	c := p.StatusCounts()

	type listItem struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		SourceStatus string `json:"sourceStatus"`
		Title        string `json:"title"`
	}
	items := []listItem{}
	for _, t := range p.Tasks {
		if filter != "" && t.Status.String() != filter {
			continue
		}
		items = append(items, listItem{t.ID, t.Status.String(), string(t.SourceStatus), t.Title})
	}

	if *asJSON {
		out := struct {
			Tasks  []listItem     `json:"tasks"`
			Counts map[string]int `json:"counts"`
		}{
			Tasks: items,
			Counts: map[string]int{
				"TODO": c[plan.StatusTODO], "TAKEN": c[plan.StatusTAKEN],
				"DONE": c[plan.StatusDONE], "BLOCKED": c[plan.StatusBLOCKED],
				"total": len(p.Tasks),
			},
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fmt.Fprintf(os.Stderr, "arctool: %v\n", err)
			return 4
		}
		return 0
	}

	for _, it := range items {
		st := it.Status
		if st == "" {
			st = "(none)"
		}
		fmt.Printf("%-16s %-8s %s\n", it.ID, st, it.Title)
	}
	fmt.Printf("(%d TODO, %d TAKEN, %d DONE, %d BLOCKED — %d total)\n",
		c[plan.StatusTODO], c[plan.StatusTAKEN], c[plan.StatusDONE], c[plan.StatusBLOCKED], len(p.Tasks))
	return 0
}

// cmdMutate implements take/done/block/todo: a guarded, byte-preserving,
// atomic rewrite of one task's status line.
func cmdMutate(cmd string, args []string) int {
	valueFlags := map[string]bool{"plan": true}
	if cmd == "block" {
		valueFlags["m"] = true
	}
	flags, pos := splitArgs(args, valueFlags)

	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	planPath := fs.String("plan", defaultPlan, "path to the plan file")
	force := fs.Bool("force", false, "override transition guardrails")
	var reason *string
	if cmd == "block" {
		reason = fs.String("m", "", "one-line reason for BLOCKED")
	}
	if err := fs.Parse(flags); err != nil {
		return 2
	}
	if len(pos) != 1 {
		fmt.Fprintf(os.Stderr, "usage: arctool %s <id> [flags]\n", cmd)
		return 2
	}
	id := pos[0]

	p, code := loadPlan(*planPath)
	if code != 0 {
		return code
	}
	matches := p.ByID(id)
	switch len(matches) {
	case 0:
		fmt.Fprintf(os.Stderr, "arctool: no task %q\n", id)
		return 3
	case 1:
		// ok
	default:
		var at []string
		for _, m := range matches {
			at = append(at, fmt.Sprintf("%d", m.Line))
		}
		fmt.Fprintf(os.Stderr, "arctool: ambiguous id %q (%d matches at lines %s)\n", id, len(matches), strings.Join(at, ", "))
		return 3
	}
	t := matches[0]

	if t.StatusLineStart < 0 {
		fmt.Fprintf(os.Stderr, "arctool: task %q has no \"- Status:\" line; run `arctool validate` and fix it first\n", id)
		return 1
	}

	hasReason := reason != nil && *reason != ""
	target, act, expected := plan.Transition(cmd, t.Status, *force, hasReason)

	switch act {
	case plan.ActNoop:
		fmt.Printf("%s already %s\n", id, t.Status)
		return 0
	case plan.ActRefuse:
		fmt.Fprintf(os.Stderr, "arctool: refusing: %s is %s, %s requires %s (use --force to override)\n",
			id, t.Status, cmd, expected)
		return 1
	}

	r := ""
	if reason != nil {
		r = *reason
	}
	out := p.WithStatus(t, target, r)
	if err := atomicWrite(*planPath, out); err != nil {
		fmt.Fprintf(os.Stderr, "arctool: write %s: %v\n", *planPath, err)
		return 4
	}
	fmt.Printf("%s %s\n", target, id)
	return 0
}

// atomicWrite writes data to path via a temp file in the same directory
// followed by rename, so a crash never leaves a half-written plan.
func atomicWrite(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".arctool-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op after a successful rename
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if fi, err := os.Stat(path); err == nil {
		_ = os.Chmod(tmpName, fi.Mode().Perm())
	}
	return os.Rename(tmpName, path)
}

func cmdValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	planPath := fs.String("plan", defaultPlan, "path to the plan file")
	strict := fs.Bool("strict", false, "enable strict checks and treat warnings as failures")
	asJSON := fs.Bool("json", false, "emit findings as JSON")
	warnAsError := fs.Bool("warn-as-error", false, "treat warnings as failures")
	reqAcc := fs.Bool("require-acceptance", false, "require an Acceptance section per task")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	p, code := loadPlan(*planPath)
	if code != 0 {
		return code
	}
	// --strict implies --require-acceptance at the CLI boundary: an executable
	// plan is only mature when every task has testable success criteria. The
	// core Validate keeps the two opts orthogonal (so existing tests hold).
	findings := p.Validate(plan.ValidateOpts{Strict: *strict, RequireAcceptance: *reqAcc || *strict})
	errs, warns := plan.Counts(findings)
	fail := errs > 0 || ((*strict || *warnAsError) && warns > 0)

	if *asJSON {
		out := struct {
			OK        bool           `json:"ok"`
			TaskCount int            `json:"taskCount"`
			Findings  []plan.Finding `json:"findings"`
		}{OK: !fail, TaskCount: len(p.Tasks), Findings: findings}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fmt.Fprintf(os.Stderr, "arctool: %v\n", err)
			return 4
		}
	} else {
		for _, f := range findings {
			fmt.Printf("%s:%d: %s: %s\n", *planPath, f.Line, f.Severity, f.Message)
		}
		if fail {
			fmt.Printf("FAIL: %d task(s), %d error(s), %d warning(s)\n", len(p.Tasks), errs, warns)
		} else {
			fmt.Printf("ok: %d task(s), %d problem(s)\n", len(p.Tasks), len(findings))
		}
	}

	if fail {
		return 1
	}
	return 0
}

// cmdArchive moves DONE blocks to the archive file and compacts the plan into a
// ledger. It writes the archive first, then the plan, so a crash never loses a
// DONE block; a re-run heals a duplicate.
func cmdArchive(args []string) int {
	fs := flag.NewFlagSet("archive", flag.ContinueOnError)
	planPath := fs.String("plan", defaultPlan, "path to the plan file")
	archivePath := fs.String("archive", "", "archive file path (default: plan-archive.md beside the plan)")
	dryRun := fs.Bool("dry-run", false, "show what would move without writing")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	ap := *archivePath
	if ap == "" {
		ap = filepath.Join(filepath.Dir(*planPath), "plan-archive.md")
	}

	p, code := loadPlan(*planPath)
	if code != 0 {
		return code
	}
	res, ok := p.Archive(ap, time.Now().Format("2006-01-02"))
	if !ok {
		fmt.Println("nothing to archive")
		return 0
	}
	pending := len(p.Tasks) - len(res.Archived)

	if *dryRun {
		fmt.Printf("would archive %d block(s) to %s:\n", len(res.Archived), ap)
		for _, t := range res.Archived {
			fmt.Printf("  %s: %s\n", t.ID, t.Title)
		}
		fmt.Printf("pending after: %d\n", pending)
		return 0
	}

	if err := plan.VerifyArchive(p, res); err != nil {
		fmt.Fprintf(os.Stderr, "arctool: archive self-validation failed: %v; nothing written\n", err)
		return 5
	}

	archiveContent, code := buildArchive(ap, res.Section)
	if code != 0 {
		return code
	}
	// Archive first (additive), then plan (destructive).
	if err := atomicWrite(ap, archiveContent); err != nil {
		fmt.Fprintf(os.Stderr, "arctool: write %s: %v\n", ap, err)
		return 4
	}
	if err := atomicWrite(*planPath, res.NewPlan); err != nil {
		fmt.Fprintf(os.Stderr, "arctool: write %s: %v\n", *planPath, err)
		return 4
	}
	fmt.Printf("archived %d, pending %d\n", len(res.Archived), pending)
	return 0
}

// buildArchive returns the archive file's new content: existing content plus a
// blank line and the new section, or a fresh header + section when absent.
func buildArchive(path string, section []byte) ([]byte, int) {
	existing, err := os.ReadFile(path)
	switch {
	case err == nil:
		buf := existing
		if len(buf) > 0 && buf[len(buf)-1] != '\n' {
			buf = append(buf, '\n')
		}
		buf = append(buf, '\n')
		buf = append(buf, section...)
		return buf, 0
	case os.IsNotExist(err):
		buf := []byte(plan.ArchiveHeader)
		buf = append(buf, '\n')
		buf = append(buf, section...)
		return buf, 0
	default:
		fmt.Fprintf(os.Stderr, "arctool: read %s: %v\n", path, err)
		return nil, 4
	}
}
