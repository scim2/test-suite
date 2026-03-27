package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/scim2/test-suite/spec"
)

var pageFooterRe = regexp.MustCompile(`\[Page \d+\]\s*$`)

var rfcNums = []int{7642, 7643, 7644}

func fileHashes(paths []string) map[string][32]byte {
	m := make(map[string][32]byte, len(paths))
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		m[p] = sha256.Sum256(data)
	}
	return m
}

func main() {
	reportPath := flag.String("report", "", "Path to compliance-report.json")
	outPath := flag.String("out", "compliance-report.html", "Output HTML file")
	version := flag.String("version", "", "Version string (e.g. git hash) to display in the footer")
	commitHash := flag.String("commit", "", "Full commit hash for GitHub source links")
	watch := flag.Bool("watch", false, "Watch template, RFC text, and report files for changes and re-render")
	tmplPath := flag.String("template", "", "Path to template.gohtml (default: next to this binary)")
	flag.Parse()

	if *tmplPath == "" {
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "finding executable path: %v\n", err)
			os.Exit(1)
		}
		*tmplPath = filepath.Join(filepath.Dir(exe), "template.gohtml")
		// Fall back to source tree location relative to cwd.
		if _, err := os.Stat(*tmplPath); err != nil {
			*tmplPath = "cmd/website/template.gohtml"
		}
	}

	cfg := renderConfig{
		tmplPath:   *tmplPath,
		reportPath: *reportPath,
		outPath:    *outPath,
		version:    *version,
		commitHash: *commitHash,
	}

	if err := render(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !*watch {
		return
	}

	// Collect files to watch. RFC text is embedded at compile time
	// so only the template and report need watching.
	watched := []string{cfg.tmplPath}
	if cfg.reportPath != "" {
		watched = append(watched, cfg.reportPath)
	}

	fmt.Println("Watching for changes...")
	hashes := fileHashes(watched)
	for {
		time.Sleep(500 * time.Millisecond)
		cur := fileHashes(watched)
		changed := false
		for f, h := range cur {
			if h != hashes[f] {
				changed = true
				break
			}
		}
		if !changed {
			continue
		}
		hashes = cur
		if err := render(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "render error: %v\n", err)
		}
	}
}

func render(cfg renderConfig) error {
	tmplHTML, err := os.ReadFile(cfg.tmplPath)
	if err != nil {
		return fmt.Errorf("reading template: %w", err)
	}

	rfcLines := make(map[string][]string)
	lineMap := make(map[int]map[int]int)
	for _, num := range rfcNums {
		text, err := spec.RFCText(num)
		if err != nil {
			return fmt.Errorf("loading RFC %d: %w", num, err)
		}
		lines := strings.Split(text, "\n")
		filtered, lm := stripPageBreaks(lines)
		rfcLines[fmt.Sprintf("%d", num)] = filtered
		lineMap[num] = lm
	}

	reqs := make([]reqData, 0, len(spec.All))
	for _, r := range spec.All {
		startLine := r.Source.StartLine
		endLine := r.Source.EndLine
		if lm, ok := lineMap[r.Source.RFC]; ok {
			if v := lm[startLine]; v > 0 {
				startLine = v
			}
			if v := lm[endLine]; v > 0 {
				endLine = v
			}
		}
		srcFile, srcStart, srcEnd := spec.SectionFile(
			r.Source.RFC, r.Source.StartLine, r.Source.EndLine,
		)
		reqs = append(reqs, reqData{
			ID:              r.ID,
			Level:           r.Level.String(),
			Summary:         r.Summary,
			RFC:             r.Source.RFC,
			Section:         r.Source.Section,
			StartLine:       startLine,
			EndLine:         endLine,
			SourceFile:      srcFile,
			SourceStartLine: srcStart,
			SourceEndLine:   srcEnd,
			StartCol:        r.Source.StartCol,
			EndCol:          r.Source.EndCol,
			Feature:         string(r.Feature),
			Testable:        r.Testable,
		})
	}

	reportJSON := "null"
	if cfg.reportPath != "" {
		data, err := os.ReadFile(cfg.reportPath)
		if err != nil {
			return fmt.Errorf("reading report: %w", err)
		}
		reportJSON = string(data)
	}

	origLines := make(map[string]map[int]int)
	for num, lm := range lineMap {
		rev := make(map[int]int)
		for old, new_ := range lm {
			rev[new_] = old
		}
		origLines[fmt.Sprintf("%d", num)] = rev
	}

	rfcsJSON, err := json.Marshal(rfcLines)
	if err != nil {
		return fmt.Errorf("marshaling RFCs: %w", err)
	}
	reqsJSON, err := json.Marshal(reqs)
	if err != nil {
		return fmt.Errorf("marshaling requirements: %w", err)
	}
	origLinesJSON, err := json.Marshal(origLines)
	if err != nil {
		return fmt.Errorf("marshaling original lines: %w", err)
	}

	tmpl, err := template.New("page").Parse(string(tmplHTML))
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	f, err := os.Create(cfg.outPath)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}

	err = tmpl.Execute(f, templateData{
		RFCs:          template.JS(rfcsJSON),
		Requirements:  template.JS(reqsJSON),
		Report:        template.JS(reportJSON),
		OriginalLines: template.JS(origLinesJSON),
		Version:       cfg.version,
		CommitHash:    cfg.commitHash,
	})
	if closeErr := f.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	fmt.Printf("Written to %s\n", cfg.outPath)
	return nil
}

// stripPageBreaks removes RFC page break blocks (surrounding blanks, footer,
// and header lines) and returns the filtered lines plus a mapping from old
// 1-based line numbers to new 1-based line numbers.
func stripPageBreaks(lines []string) ([]string, map[int]int) {
	remove := make(map[int]bool)
	for i, line := range lines {
		if !pageFooterRe.MatchString(line) {
			continue
		}
		remove[i] = true
		// Remove blank lines before the footer.
		for j := i - 1; j >= 0 && strings.TrimSpace(lines[j]) == ""; j-- {
			remove[j] = true
		}
		// Remove blank line + header + blank lines after the footer.
		for j := i + 1; j < len(lines) && (strings.TrimSpace(lines[j]) == "" || strings.HasPrefix(lines[j], "RFC ")); j++ {
			remove[j] = true
		}
	}

	// Also trim leading blank lines.
	leadingDone := false
	filtered := make([]string, 0, len(lines)-len(remove))
	lm := make(map[int]int, len(lines))
	for i, line := range lines {
		if remove[i] {
			continue
		}
		if !leadingDone && strings.TrimSpace(line) == "" {
			continue
		}
		leadingDone = true
		filtered = append(filtered, line)
		lm[i+1] = len(filtered) // old 1-based -> new 1-based
	}
	return filtered, lm
}

type renderConfig struct {
	tmplPath   string
	reportPath string
	outPath    string
	version    string
	commitHash string
}

type reqData struct {
	ID              string `json:"id"`
	Level           string `json:"level"`
	Summary         string `json:"summary"`
	RFC             int    `json:"rfc"`
	Section         string `json:"section"`
	StartLine       int    `json:"startLine"`
	EndLine         int    `json:"endLine"`
	SourceFile      string `json:"sourceFile"`
	SourceStartLine int    `json:"sourceStartLine"`
	SourceEndLine   int    `json:"sourceEndLine"`
	StartCol        int    `json:"startCol"`
	EndCol          int    `json:"endCol"`
	Feature         string `json:"feature"`
	Testable        bool   `json:"testable"`
}

type templateData struct {
	RFCs          template.JS
	Requirements  template.JS
	Report        template.JS
	OriginalLines template.JS
	Version       string
	CommitHash    string
}
