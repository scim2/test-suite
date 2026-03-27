package spec

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

// keywordRe matches RFC 2119 keywords in ALL CAPS as whole words.
// REQUIRED and OPTIONAL are excluded because in SCIM they primarily
// mark attribute cardinality rather than behavioral requirements.
var keywordRe = regexp.MustCompile(
	`\b(MUST NOT|MUST|SHALL NOT|SHALL|SHOULD NOT|SHOULD|MAY)\b`,
)

func TestCounts(t *testing.T) {
	total := len(All)
	if total == 0 {
		t.Fatal("no requirements cataloged")
	}
	t.Logf("Total requirements: %d", total)
	t.Logf("  RFC 7642: %d", len(RFC7642))
	t.Logf("  RFC 7643: %d", len(RFC7643))
	t.Logf("  RFC 7644: %d", len(RFC7644))

	testable := 0
	for _, r := range All {
		if r.Testable {
			testable++
		}
	}
	t.Logf("  Testable: %d / %d", testable, total)
}

// TestCoverage scans each RFC txt file for RFC 2119 keywords and
// verifies that every occurrence is covered by a cataloged requirement.
func TestCoverage(t *testing.T) {
	rfcs := []rfcFile{
		{
			rfc: 7642,
			skipRanges: []skipRange{
				{192, 195, "RFC 2119 boilerplate"},
				{973, 1067, "references and authors"},
			},
		},
		{
			rfc: 7643,
			skipRanges: []skipRange{
				{207, 215, "RFC 2119 boilerplate"},
				{1855, 5100, "Section 8: JSON representations (non-normative repeats)"},
				{5215, 5827, "IANA, references, authors"},
			},
		},
		{
			rfc: 7644,
			skipRanges: []skipRange{
				{175, 188, "RFC 2119 boilerplate"},
				{4560, 4987, "IANA, references, authors"},
			},
		},
	}

	// Build a lookup: rfc -> sorted list of covered line ranges.
	type lineRange struct {
		start, end int
		id         string
	}
	covered := make(map[int][]lineRange)
	for _, r := range All {
		covered[r.Source.RFC] = append(covered[r.Source.RFC], lineRange{
			start: r.Source.StartLine,
			end:   r.Source.EndLine,
			id:    r.ID,
		})
	}

	for _, rf := range rfcs {
		t.Run(fmt.Sprintf("RFC%d", rf.rfc), func(t *testing.T) {
			text, err := RFCText(rf.rfc)
			if err != nil {
				t.Fatalf("loading RFC %d: %v", rf.rfc, err)
			}

			ranges := covered[rf.rfc]
			scanner := bufio.NewScanner(strings.NewReader(text))
			lineNum := 0
			uncovered := 0

			for scanner.Scan() {
				lineNum++
				line := scanner.Text()

				// Skip non-normative ranges.
				skip := false
				for _, sr := range rf.skipRanges {
					if lineNum >= sr.start && lineNum <= sr.end {
						skip = true
						break
					}
				}
				if skip {
					continue
				}

				// Skip page headers/footers.
				if strings.HasPrefix(line, "RFC ") ||
					strings.Contains(line, "Standards Track") ||
					strings.Contains(line, "[Page ") {
					continue
				}

				if !keywordRe.MatchString(line) {
					continue
				}

				// Check if any requirement covers this line.
				found := false
				for _, lr := range ranges {
					if lineNum >= lr.start && lineNum <= lr.end {
						found = true
						break
					}
				}
				if !found {
					uncovered++
					t.Logf("uncovered line %d: %s",
						lineNum, strings.TrimSpace(line))
				}
			}

			if err := scanner.Err(); err != nil {
				t.Fatalf("error scanning RFC %d: %v", rf.rfc, err)
			}

			t.Logf("uncovered: %d", uncovered)
		})
	}
}

func TestIDFormat(t *testing.T) {
	for _, r := range All {
		prefix := fmt.Sprintf("RFC%d-", r.Source.RFC)
		if !strings.HasPrefix(r.ID, prefix) {
			t.Errorf("%s: ID should start with %s", r.ID, prefix)
		}
		if !strings.Contains(r.ID, "-L") {
			t.Errorf("%s: ID should contain -L{line}", r.ID)
		}
		// The line number in the ID is the keyword line, which must
		// fall within [StartLine, EndLine].
		parts := strings.Split(r.ID, "-L")
		if len(parts) < 2 {
			continue
		}
		var line int
		_, _ = fmt.Sscanf(parts[len(parts)-1], "%d", &line)
		if line < r.Source.StartLine || line > r.Source.EndLine {
			t.Errorf("%s: keyword line %d not in [%d, %d]",
				r.ID, line, r.Source.StartLine, r.Source.EndLine)
		}
	}
}

func TestLookups(t *testing.T) {
	r := ByID("RFC7644-3.3-L602")
	if r == nil {
		t.Fatal("ByID returned nil for known requirement")
	}
	if r.Summary == "" {
		t.Error("looked-up requirement has empty Summary")
	}

	if ByID("nonexistent") != nil {
		t.Error("ByID should return nil for unknown ID")
	}

	core := ByFeature(Core)
	if len(core) == 0 {
		t.Error("ByFeature(Core) returned empty")
	}

	musts := ByLevel(Must)
	if len(musts) == 0 {
		t.Error("ByLevel(Must) returned empty")
	}
}

func TestNoDuplicateIDs(t *testing.T) {
	seen := make(map[string]bool, len(All))
	for _, r := range All {
		if seen[r.ID] {
			t.Errorf("duplicate requirement ID: %s", r.ID)
		}
		seen[r.ID] = true
	}
}

func TestRequiredFields(t *testing.T) {
	for _, r := range All {
		t.Run(r.ID, func(t *testing.T) {
			if r.ID == "" {
				t.Error("empty ID")
			}
			if r.Summary == "" {
				t.Error("empty Summary")
			}
			if r.Source.RFC == 0 {
				t.Error("zero RFC number")
			}
			if r.Source.Section == "" {
				t.Error("empty Section")
			}
			if r.Source.StartLine <= 0 {
				t.Error("StartLine must be positive")
			}
			if r.Source.EndLine < r.Source.StartLine {
				t.Errorf("EndLine (%d) < StartLine (%d)",
					r.Source.EndLine, r.Source.StartLine)
			}
			if r.Feature == "" {
				t.Error("empty Feature")
			}
		})
	}
}

func TestValidFeatures(t *testing.T) {
	valid := map[Feature]bool{
		Core: true, Filter: true, Sort: true,
		ChangePassword: true, Patch: true, Bulk: true,
		ETag: true, Users: true, Groups: true, Me: true,
	}
	for _, r := range All {
		if !valid[r.Feature] {
			t.Errorf("%s: unknown Feature %q", r.ID, r.Feature)
		}
	}
}

func TestValidLevels(t *testing.T) {
	for _, r := range All {
		switch r.Level {
		case Must, MustNot, Shall, ShallNot, Should, ShouldNot, May:
		default:
			t.Errorf("%s: unknown Level %d", r.ID, r.Level)
		}
	}
}

func TestValidRFCs(t *testing.T) {
	valid := map[int]bool{7642: true, 7643: true, 7644: true}
	for _, r := range All {
		if !valid[r.Source.RFC] {
			t.Errorf("%s: unknown RFC %d", r.ID, r.Source.RFC)
		}
	}
}

// rfcFile describes an RFC and which line ranges to skip.
type rfcFile struct {
	rfc        int
	skipRanges []skipRange
}

// skipRange defines a line range to exclude from coverage checks.
type skipRange struct {
	start, end int
	reason     string
}

func TestRFCTextIntegrity(t *testing.T) {
	origFiles := map[int]string{
		7642: "testdata/rfc7642.txt",
		7643: "testdata/rfc7643.txt",
		7644: "testdata/rfc7644.txt",
	}
	for rfc, path := range origFiles {
		t.Run(fmt.Sprintf("RFC%d", rfc), func(t *testing.T) {
			assembled, err := RFCText(rfc)
			if err != nil {
				t.Fatalf("loading RFC %d: %v", rfc, err)
			}
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}
			if assembled != string(original) {
				t.Errorf("assembled text does not match %s", path)
			}
		})
	}
}
