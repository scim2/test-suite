package spec

// Extra contains additional compliance checks that go beyond the
// literal RFC requirements. These exercise attribute type combinations,
// extension handling, and edge cases that the standard User/Group
// schemas do not cover.
var Extra = []Requirement{
	// Attribute types.
	{
		ID:       "EXTRA-TYPE-INT",
		Level:    Must,
		Summary:  "Integer attributes must round-trip without fractional parts",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-TYPE-BOOL",
		Level:    Must,
		Summary:  "Boolean attributes must round-trip as JSON true/false",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-TYPE-DT",
		Level:    Must,
		Summary:  "DateTime attributes must round-trip as valid xsd:dateTime strings",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-TYPE-BIN",
		Level:    Must,
		Summary:  "Binary attributes must round-trip as base64-encoded strings",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-TYPE-REF",
		Level:    Must,
		Summary:  "Reference attributes must round-trip as URI strings",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-TYPE-COMPLEX",
		Level:    Must,
		Summary:  "Complex attributes must round-trip with sub-attributes intact",
		Feature:  TestResource,
		Testable: true,
	},

	// Mutability.
	{
		ID:       "EXTRA-MUT-IMMUTABLE",
		Level:    Must,
		Summary:  "Immutable attributes must be accepted at creation and rejected on update",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-MUT-WRITEONLY",
		Level:    Must,
		Summary:  "WriteOnly attributes must be accepted but never returned",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-MUT-READONLY",
		Level:    Must,
		Summary:  "ReadOnly attributes provided by the client must be ignored",
		Feature:  TestResource,
		Testable: true,
	},

	// Returned behavior.
	{
		ID:       "EXTRA-RET-ALWAYS",
		Level:    Must,
		Summary:  "Returned:always attributes must appear even when excluded",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-RET-NEVER",
		Level:    Must,
		Summary:  "Returned:never attributes must never appear in responses",
		Feature:  TestResource,
		Testable: true,
	},

	// Case sensitivity.
	{
		ID:       "EXTRA-CASE-EXACT",
		Level:    Must,
		Summary:  "Filters on caseExact=true attributes must match case-sensitively",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-CASE-INSENSITIVE",
		Level:    Must,
		Summary:  "Filters on caseExact=false attributes must match case-insensitively",
		Feature:  TestResource,
		Testable: true,
	},

	// Extensions.
	{
		ID:       "EXTRA-EXT-CONTAINER",
		Level:    Must,
		Summary:  "Extension attributes must be namespaced under the extension schema URI",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-EXT-ROUNDTRIP",
		Level:    Must,
		Summary:  "Extension attributes must round-trip through create and get",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-EXT-SCHEMAS",
		Level:    Must,
		Summary:  "Resources with extensions must include extension URI in schemas array",
		Feature:  TestResource,
		Testable: true,
	},

	// Multi-valued.
	{
		ID:       "EXTRA-MV-ROUNDTRIP",
		Level:    Must,
		Summary:  "Multi-valued attributes must round-trip preserving all values",
		Feature:  TestResource,
		Testable: true,
	},
	{
		ID:       "EXTRA-MV-COMPLEX-PRIMARY",
		Level:    Must,
		Summary:  "Multi-valued complex attributes must enforce at most one primary=true",
		Feature:  TestResource,
		Testable: true,
	},
}

func init() {
	All = append(All, Extra...)
}
