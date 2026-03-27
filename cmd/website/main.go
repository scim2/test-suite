package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/scim2/test-suite/spec"
)

var pageFooterRe = regexp.MustCompile(`\[Page \d+\]\s*$`)

//go:embed template.gohtml
var tmplHTML string

func main() {
	reportPath := flag.String("report", "", "Path to compliance-report.json")
	outPath := flag.String("out", "compliance-report.html", "Output HTML file")
	version := flag.String("version", "", "Version string (e.g. git hash) to display in the footer")
	commitHash := flag.String("commit", "", "Full commit hash for GitHub source links")
	flag.Parse()

	rfcFiles := map[int]string{
		7642: "spec/testdata/rfc7642.txt",
		7643: "spec/testdata/rfc7643.txt",
		7644: "spec/testdata/rfc7644.txt",
	}
	rfcLines := make(map[string][]string)
	// lineMap maps RFC number -> old 1-based line number -> new 1-based line number.
	// Lines that were removed map to 0.
	lineMap := make(map[int]map[int]int)
	for num, path := range rfcFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %v\n", path, err)
			os.Exit(1)
		}
		lines := strings.Split(string(data), "\n")
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
		reqs = append(reqs, reqData{
			ID:              r.ID,
			Level:           r.Level.String(),
			Summary:         r.Summary,
			RFC:             r.Source.RFC,
			Section:         r.Source.Section,
			StartLine:       startLine,
			EndLine:         endLine,
			SourceStartLine: r.Source.StartLine,
			SourceEndLine:   r.Source.EndLine,
			StartCol:        r.Source.StartCol,
			EndCol:          r.Source.EndCol,
			Feature:         string(r.Feature),
			Testable:        r.Testable,
		})
	}

	reportJSON := "null"
	if *reportPath != "" {
		data, err := os.ReadFile(*reportPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading report: %v\n", err)
			os.Exit(1)
		}
		reportJSON = string(data)
	}

	// Build reverse line map: new 1-based -> old 1-based per RFC.
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
		fmt.Fprintf(os.Stderr, "marshaling RFCs: %v\n", err)
		os.Exit(1)
	}
	reqsJSON, err := json.Marshal(reqs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshaling requirements: %v\n", err)
		os.Exit(1)
	}
	origLinesJSON, err := json.Marshal(origLines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshaling original lines: %v\n", err)
		os.Exit(1)
	}

	tmpl, err := template.New("page").Parse(tmplHTML)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing template: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating output: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "closing output: %v\n", err)
			os.Exit(1)
		}
	}()

	err = tmpl.Execute(f, templateData{
		RFCs:          template.JS(rfcsJSON),
		Requirements:  template.JS(reqsJSON),
		Report:        template.JS(reportJSON),
		OriginalLines: template.JS(origLinesJSON),
		Version:       *version,
		CommitHash:    *commitHash,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Written to %s\n", *outPath)
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

type reqData struct {
	ID              string `json:"id"`
	Level           string `json:"level"`
	Summary         string `json:"summary"`
	RFC             int    `json:"rfc"`
	Section         string `json:"section"`
	StartLine       int    `json:"startLine"`
	EndLine         int    `json:"endLine"`
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
