package spec

// Feature groups requirements by optional server capability.
// Tests for a feature are skipped when the feature is not supported.
type Feature string

const (
	// Core is always tested (required endpoints, basic CRUD, errors).
	Core Feature = "core"

	// Discovered via /ServiceProviderConfig.
	Filter         Feature = "filter"
	Sort           Feature = "sort"
	ChangePassword Feature = "changePassword"
	Patch          Feature = "patch"
	Bulk           Feature = "bulk"
	ETag           Feature = "etag"

	// Resource types discovered via /ResourceTypes.
	Users  Feature = "users"
	Groups Feature = "groups"

	// Custom test resource for exercising all attribute types.
	TestResource Feature = "testResource"

	// The /Me endpoint alias.
	Me Feature = "me"
)

// Level represents an RFC 2119 compliance keyword.
type Level int

const (
	Must Level = iota
	MustNot
	Shall
	ShallNot
	Should
	ShouldNot
	May
)

func (l Level) String() string {
	switch l {
	case Must:
		return "MUST"
	case MustNot:
		return "MUST NOT"
	case Shall:
		return "SHALL"
	case ShallNot:
		return "SHALL NOT"
	case Should:
		return "SHOULD"
	case ShouldNot:
		return "SHOULD NOT"
	case May:
		return "MAY"
	default:
		return "UNKNOWN"
	}
}

// Requirement is a single testable statement from the SCIM RFCs.
type Requirement struct {
	// ID is a stable identifier.
	// Format: RFC{num}-{section}-L{line}, where {line} is the line
	// containing the RFC 2119 keyword. This may differ from
	// Source.StartLine when the sentence begins on an earlier line.
	ID string

	// Level is the RFC 2119 compliance level.
	Level Level

	// Summary is a short, test-oriented description of what to verify.
	Summary string

	// Source is the location in the RFC txt file.
	Source Source

	// Feature determines when this requirement is tested.
	Feature Feature

	// Testable indicates whether this can be verified black-box.
	// Some requirements (e.g. internal password hashing) are not
	// externally observable and are marked false.
	Testable bool
}

// Source pinpoints where a requirement comes from in the RFC txt.
type Source struct {
	// RFC number: 7642, 7643, or 7644.
	RFC int
	// Section identifier, e.g. "3.5.1".
	Section string
	// StartLine is the starting line number in the plain-text RFC file.
	StartLine int
	// EndLine is the ending line number (inclusive). Equal to StartLine
	// when the requirement fits on a single line.
	EndLine int
	// StartCol is the 1-based column where highlighting begins on
	// StartLine. Zero means highlight from the beginning of the line.
	StartCol int
	// EndCol is the 1-based column where highlighting ends (inclusive)
	// on EndLine. Zero means highlight to the end of the line.
	EndCol int
}
