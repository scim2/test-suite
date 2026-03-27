package spec

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed testdata/rfc7642.txt
var rfc7642Text string

//go:embed testdata/rfc7643/*.txt
var rfc7643FS embed.FS

var rfc7643Sections = []section{
	{"00-front.txt", 1},
	{"01-section-1.txt", 148},
	{"02-section-1.1.txt", 207},
	{"03-section-1.2.txt", 240},
	{"04-section-2.txt", 316},
	{"05-section-2.1.txt", 355},
	{"06-section-2.2.txt", 399},
	{"07-section-2.3.txt", 438},
	{"08-section-2.4.txt", 598},
	{"09-section-2.5.txt", 679},
	{"10-section-3.txt", 689},
	{"11-section-3.1.txt", 847},
	{"12-section-3.2.txt", 959},
	{"13-section-3.3.txt", 968},
	{"14-section-4.txt", 1015},
	{"15-section-4.1.txt", 1021},
	{"16-section-4.2.txt", 1383},
	{"17-section-4.3.txt", 1428},
	{"18-section-5.txt", 1477},
	{"19-section-6.txt", 1602},
	{"20-section-7.txt", 1660},
	{"21-section-8.txt", 1855},
	{"22-section-9.txt", 5103},
	{"23-back.txt", 5215},
}

//go:embed testdata/rfc7644/*.txt
var rfc7644FS embed.FS

var rfc7644Sections = []section{
	{"00-front.txt", 1},
	{"01-section-1.txt", 148},
	{"02-section-1.1.txt", 161},
	{"03-section-1.2.txt", 175},
	{"04-section-1.3.txt", 194},
	{"05-section-2.txt", 231},
	{"06-section-2.1.txt", 355},
	{"07-section-2.2.txt", 380},
	{"08-section-3.txt", 399},
	{"09-section-3.1.txt", 401},
	{"10-section-3.2.txt", 467},
	{"11-section-3.3.txt", 574},
	{"12-section-3.4.txt", 717},
	{"13-section-3.5.txt", 1604},
	{"14-section-3.6.txt", 2639},
	{"15-section-3.7.txt", 2695},
	{"16-section-3.8.txt", 3535},
	{"17-section-3.9.txt", 3571},
	{"18-section-3.10.txt", 3647},
	{"19-section-3.11.txt", 3667},
	{"20-section-3.12.txt", 3703},
	{"21-section-3.13.txt", 3948},
	{"22-section-3.14.txt", 3961},
	{"23-section-4.txt", 4060},
	{"24-section-5.txt", 4207},
	{"25-section-6.txt", 4222},
	{"26-section-7.txt", 4330},
	{"27-back.txt", 4560},
}

// RFCText returns the full text of the given RFC, reassembled from
// per-section files. The result is byte-identical to the original
// monolithic file.
func RFCText(rfc int) (string, error) {
	switch rfc {
	case 7642:
		return rfc7642Text, nil
	case 7643:
		return assembleFS(rfc7643FS, "testdata/rfc7643")
	case 7644:
		return assembleFS(rfc7644FS, "testdata/rfc7644")
	default:
		return "", fmt.Errorf("unknown RFC %d", rfc)
	}
}

// SectionFile returns the testdata file path and local line numbers for a
// given RFC and global line range. This is used to generate source links
// that point to the correct section file.
func SectionFile(rfc, startLine, endLine int) (path string, localStart, localEnd int) {
	var sections []section
	var dir string
	switch rfc {
	case 7642:
		return "spec/testdata/rfc7642.txt", startLine, endLine
	case 7643:
		sections = rfc7643Sections
		dir = "spec/testdata/rfc7643"
	case 7644:
		sections = rfc7644Sections
		dir = "spec/testdata/rfc7644"
	default:
		return "", 0, 0
	}

	// Find the section that contains startLine.
	idx := len(sections) - 1
	for i := len(sections) - 1; i >= 0; i-- {
		if sections[i].startLine <= startLine {
			idx = i
			break
		}
	}

	s := sections[idx]
	return dir + "/" + s.file, startLine - s.startLine + 1, endLine - s.startLine + 1
}

func assembleFS(fs embed.FS, dir string) (string, error) {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, e := range entries {
		data, err := fs.ReadFile(dir + "/" + e.Name())
		if err != nil {
			return "", err
		}
		b.Write(data)
	}
	return b.String(), nil
}

// lineMapper converts local line numbers within a section file to
// global line numbers in the original monolithic RFC text.
type lineMapper func(local int) int

// rfcLines returns a lineMapper for the given section file within an RFC.
func rfcLines(sections []section, file string) lineMapper {
	for _, s := range sections {
		if s.file == file {
			base := s.startLine - 1
			return func(local int) int { return base + local }
		}
	}
	panic("unknown section file: " + file)
}

// section records a section file's name and the global line where it starts.
type section struct {
	file      string
	startLine int // 1-based global line number
}
