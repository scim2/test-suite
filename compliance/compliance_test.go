package compliance

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/scim2/test-suite/scim"
	"github.com/scim2/test-suite/spec"
	"github.com/scim2/test-suite/testserver"
)

var (
	flagURL    = flag.String("scim.url", "", "SCIM base URL (uses in-memory server if empty)")
	flagToken  = flag.String("scim.token", "", "Bearer token")
	flagUser   = flag.String("scim.user", "", "Basic auth user")
	flagPass   = flag.String("scim.pass", "", "Basic auth password")
	flagForce  = flag.String("scim.force", "", "Comma-separated features to force-enable")
	flagReport = flag.String("scim.report", "compliance-report.txt", "Path to write JSON report")
)

func TestCompliance(t *testing.T) {
	flag.Parse()

	baseURL := *flagURL

	var ts *httptest.Server
	if baseURL == "" {
		ts = httptest.NewServer(testserver.New())
		defer ts.Close()
		baseURL = ts.URL
		fmt.Println("using in-memory test server at", baseURL)
	}

	client := &scim.Client{
		BaseURL:    baseURL,
		HTTPClient: new(http.Client),
	}

	switch {
	case *flagToken != "":
		client.Auth = scim.BearerAuth{Token: *flagToken}
	case *flagUser != "":
		client.Auth = scim.BasicAuth{User: *flagUser, Pass: *flagPass}
	}

	features, err := client.Discover()
	if err != nil {
		t.Fatalf("discovery failed: %v", err)
	}

	report := RunAll(Config{
		Client:   client,
		Features: features,
		Force:    parseForce(*flagForce),
		Verbose:  true,
	})

	writeReport(baseURL, report.Results)

	// Fail the test if any requirement failed.
	for _, r := range report.Results {
		if r.Outcome == Fail {
			t.Errorf("[FAIL] %s/%s: %s",
				strings.Join(r.RequirementIDs, ","), r.TestName, r.Message)
		}
	}
}

func parseForce(s string) map[spec.Feature]bool {
	m := make(map[spec.Feature]bool)
	if s == "" {
		return m
	}
	for f := range strings.SplitSeq(s, ",") {
		m[spec.Feature(strings.TrimSpace(f))] = true
	}
	return m
}

func writeEntry(b *strings.Builder, e reportEntry) {
	fmt.Fprintf(b, "  [%s] %s\n", e.Level, e.ID)
	fmt.Fprintf(b, "    %s\n", e.Summary)
	fmt.Fprintf(b, "    RFC %d, Section %s, Lines %d-%d  (feature: %s)\n",
		e.RFC, e.Section, e.Lines[0], e.Lines[1], e.Feature)
	if e.Message != "" {
		fmt.Fprintf(b, "    -> %s\n", e.Message)
	}
}

func writeReport(baseURL string, results []Result) {
	// Build per-requirement outcome map (first outcome wins for dedup).
	type reqResult struct {
		outcome  Outcome
		message  string
		subtests []subtestEntry
	}
	byReq := make(map[string]*reqResult)
	for _, r := range results {
		for _, id := range r.RequirementIDs {
			rr, exists := byReq[id]
			if !exists {
				rr = &reqResult{}
				byReq[id] = rr
			}
			if !exists || r.Outcome == Fail {
				rr.outcome = r.Outcome
				rr.message = r.Message
			}
			if strings.Contains(r.TestName, "/") {
				parts := strings.SplitN(r.TestName, "/", 2)
				outcome := "pass"
				switch r.Outcome {
				case Fail:
					outcome = "fail"
				case Warn:
					outcome = "warn"
				case Skip:
					outcome = "skip"
				}
				rr.subtests = append(rr.subtests, subtestEntry{
					Name:    parts[1],
					Outcome: outcome,
					Message: r.Message,
				})
			}
		}
	}

	// Build level counts.
	byLevel := make(map[string]*levelCounts)
	for id, rr := range byReq {
		req := spec.ByID(id)
		if req == nil {
			continue
		}
		level := req.Level.String()
		c, ok := byLevel[level]
		if !ok {
			c = &levelCounts{}
			byLevel[level] = c
		}
		switch rr.outcome {
		case Pass:
			c.Pass++
		case Fail:
			c.Fail++
		case Warn:
			c.Warn++
		}
	}

	// Categorize all testable requirements.
	var failures, warnings, passed, untested []reportEntry

	testable := 0
	for _, r := range spec.All {
		if !r.Testable {
			continue
		}
		testable++

		entry := reportEntry{
			ID:      r.ID,
			Level:   r.Level.String(),
			Summary: r.Summary,
			RFC:     r.Source.RFC,
			Section: r.Source.Section,
			Lines:   [2]int{r.Source.StartLine, r.Source.EndLine},
			Feature: string(r.Feature),
		}

		rr, tested := byReq[r.ID]
		if !tested {
			entry.Outcome = "untested"
			untested = append(untested, entry)
			continue
		}

		entry.Subtests = rr.subtests

		switch rr.outcome {
		case Pass:
			entry.Outcome = "pass"
			passed = append(passed, entry)
		case Fail:
			entry.Outcome = "fail"
			entry.Message = rr.message
			failures = append(failures, entry)
		case Warn:
			entry.Outcome = "warn"
			entry.Message = rr.message
			warnings = append(warnings, entry)
		}
	}

	// Marshal level counts to plain map.
	levelMap := make(map[string]levelCounts, len(byLevel))
	for k, v := range byLevel {
		levelMap[k] = *v
	}

	coverPct := 0.0
	if testable > 0 {
		coverPct = float64(len(byReq)) / float64(testable) * 100
	}

	rpt := report{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		BaseURL:   baseURL,
		ByLevel:   levelMap,
		Failures:  failures,
		Warnings:  warnings,
		Passed:    passed,
		Untested:  untested,
		Coverage: reportCoverage{
			Tested:   len(byReq),
			Testable: testable,
			Percent:  coverPct,
		},
	}

	// Write report file.
	if *flagReport != "" {
		var b strings.Builder

		fmt.Fprintf(&b, "SCIM 2.0 Compliance Report\n")
		fmt.Fprintf(&b, "Generated: %s\n", rpt.Timestamp)
		fmt.Fprintf(&b, "Target:    %s\n", rpt.BaseURL)
		fmt.Fprintf(&b, "Coverage:  %d/%d testable requirements (%.0f%%)\n",
			rpt.Coverage.Tested, rpt.Coverage.Testable, rpt.Coverage.Percent)

		fmt.Fprintf(&b, "\nSummary\n")
		for _, level := range []string{"MUST", "MUST NOT", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "MAY"} {
			c, ok := levelMap[level]
			if !ok {
				continue
			}
			total := c.Pass + c.Fail
			line := fmt.Sprintf("  %-12s %d/%d passed, %d failed",
				level, c.Pass, total, c.Fail)
			if c.Warn > 0 {
				line += fmt.Sprintf(", %d warned", c.Warn)
			}
			fmt.Fprintln(&b, line)
		}

		if len(failures) > 0 {
			fmt.Fprintf(&b, "\nFailures (%d)\n", len(failures))
			for _, f := range failures {
				writeEntry(&b, f)
			}
		}

		if len(warnings) > 0 {
			fmt.Fprintf(&b, "\nWarnings (%d) - feature not advertised\n", len(warnings))
			for _, w := range warnings {
				writeEntry(&b, w)
			}
		}

		if len(passed) > 0 {
			fmt.Fprintf(&b, "\nPassed (%d)\n", len(passed))
			for _, p := range passed {
				fmt.Fprintf(&b, "  [PASS] %-30s [%-10s] %s\n", p.ID, p.Level, p.Summary)
			}
		}

		if len(untested) > 0 {
			fmt.Fprintf(&b, "\nUntested (%d)\n", len(untested))
			for _, u := range untested {
				fmt.Fprintf(&b, "  [ -- ] %-30s [%-10s] %s\n", u.ID, u.Level, u.Summary)
			}
		}

		if err := os.WriteFile(*flagReport, []byte(b.String()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write report: %v\n", err)
		} else {
			fmt.Printf("\nReport written to %s\n", *flagReport)
		}
	}

	// Write JSON report alongside text.
	if *flagReport != "" {
		jsonPath := strings.TrimSuffix(*flagReport, filepath.Ext(*flagReport)) + ".json"
		jsonData, err := json.MarshalIndent(rpt, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal JSON report: %v\n", err)
		} else if err := os.WriteFile(jsonPath, jsonData, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write JSON report: %v\n", err)
		} else {
			fmt.Printf("JSON report written to %s\n", jsonPath)
		}
	}

	// Console summary (short).
	fmt.Println("\nSCIM 2.0 Compliance Report")
	for _, level := range []string{"MUST", "MUST NOT", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "MAY"} {
		c, ok := levelMap[level]
		if !ok {
			continue
		}
		total := c.Pass + c.Fail
		line := fmt.Sprintf("  %-12s %d/%d passed, %d failed",
			level, c.Pass, total, c.Fail)
		if c.Warn > 0 {
			line += fmt.Sprintf(", %d warned", c.Warn)
		}
		fmt.Println(line)
	}

	if len(failures) > 0 {
		fmt.Printf("\n  Failures (%d):\n", len(failures))
		for _, f := range failures {
			fmt.Printf("    %s: %s\n", f.ID, f.Message)
		}
	}

	fmt.Printf("\n  Coverage: %d/%d (%.0f%%) - see %s for details\n",
		rpt.Coverage.Tested, rpt.Coverage.Testable, rpt.Coverage.Percent, *flagReport)
}

type levelCounts struct {
	Pass int `json:"pass"`
	Fail int `json:"fail"`
	Warn int `json:"warn"`
}

type report struct {
	Timestamp string                 `json:"timestamp"`
	BaseURL   string                 `json:"baseURL"`
	ByLevel   map[string]levelCounts `json:"byLevel"`
	Failures  []reportEntry          `json:"failures"`
	Warnings  []reportEntry          `json:"warnings"`
	Passed    []reportEntry          `json:"passed"`
	Untested  []reportEntry          `json:"untested"`
	Coverage  reportCoverage         `json:"coverage"`
}

type reportCoverage struct {
	Tested   int     `json:"tested"`
	Testable int     `json:"testable"`
	Percent  float64 `json:"percent"`
}

type reportEntry struct {
	ID       string         `json:"id"`
	Level    string         `json:"level"`
	Summary  string         `json:"summary"`
	RFC      int            `json:"rfc"`
	Section  string         `json:"section"`
	Lines    [2]int         `json:"lines"`
	Feature  string         `json:"feature"`
	Outcome  string         `json:"outcome"`
	Message  string         `json:"message,omitempty"`
	Subtests []subtestEntry `json:"subtests,omitempty"`
}

type subtestEntry struct {
	Name    string `json:"name"`
	Outcome string `json:"outcome"`
	Message string `json:"message,omitempty"`
}
