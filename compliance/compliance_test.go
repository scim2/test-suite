package compliance

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
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

var (
	client   *scim.Client
	features *scim.Features
	force    map[spec.Feature]bool

	mu      sync.Mutex
	results []Result
)

// TestMain initializes the SCIM client, discovers features, runs all
// tests, and writes a compliance report.
func TestMain(m *testing.M) {
	flag.Parse()

	baseURL := *flagURL

	var ts *httptest.Server
	if baseURL == "" {
		ts = httptest.NewServer(testserver.New())
		baseURL = ts.URL
		fmt.Println("using in-memory test server at", baseURL)
	}

	client = &scim.Client{
		BaseURL:    baseURL,
		HTTPClient: new(http.Client),
	}

	switch {
	case *flagToken != "":
		client.Auth = scim.BearerAuth{Token: *flagToken}
	case *flagUser != "":
		client.Auth = scim.BasicAuth{User: *flagUser, Pass: *flagPass}
	}

	force = parseForce(*flagForce)

	var err error
	features, err = client.Discover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "discovery failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	writeReport(baseURL)

	if ts != nil {
		ts.Close()
	}

	os.Exit(code)
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

func writeReport(baseURL string) {
	mu.Lock()
	defer mu.Unlock()

	// Build per-requirement outcome map (first outcome wins for dedup).
	type reqResult struct {
		outcome Outcome
		message string
	}
	byReq := make(map[string]reqResult)
	for _, r := range results {
		for _, id := range r.RequirementIDs {
			if _, exists := byReq[id]; !exists || r.Outcome == Fail {
				byReq[id] = reqResult{r.Outcome, r.Message}
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
	Pass, Fail, Warn int
}

type report struct {
	Timestamp string
	BaseURL   string
	ByLevel   map[string]levelCounts
	Failures  []reportEntry
	Warnings  []reportEntry
	Passed    []reportEntry
	Untested  []reportEntry
	Coverage  reportCoverage
}

type reportCoverage struct {
	Tested   int
	Testable int
	Percent  float64
}

type reportEntry struct {
	ID      string
	Level   string
	Summary string
	RFC     int
	Section string
	Lines   [2]int
	Feature string
	Outcome string
	Message string
}
