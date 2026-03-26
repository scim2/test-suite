package spec

// RFC7644 contains requirements from RFC 7644 (SCIM Protocol).
var RFC7644 = []Requirement{

	// Section 1.3 - Definitions

	{
		ID:      "RFC7644-1.3-L201",
		Level:   MustNot,
		Summary: "Base URI must not contain a query string",
		Source: Source{
			RFC:       7644,
			Section:   "1.3",
			StartLine: 201,
			EndLine:   201,
		},
		Feature:  Core,
		Testable: false,
	},

	// Section 2 - Authentication and Authorization

	{
		ID:      "RFC7644-2-L307",
		Level:   Shall,
		Summary: "Service provider shall indicate supported HTTP auth schemes via WWW-Authenticate",
		Source: Source{
			RFC:       7644,
			Section:   "2",
			StartLine: 307,
			EndLine:   309,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-2-L311",
		Level:   Must,
		Summary: "Service provider must map authenticated client to access control policy",
		Source: Source{
			RFC:       7644,
			Section:   "2",
			StartLine: 311,
			EndLine:   314,
		},
		Feature:  Core,
		Testable: false,
	},

	// Section 3.1 - Background

	{
		ID:      "RFC7644-3.1-L436",
		Level:   Must,
		Summary: "SCIM message schemas URI prefix must begin with urn:ietf:params:scim:api:",
		Source: Source{
			RFC:       7644,
			Section:   "3.1",
			StartLine: 436,
			EndLine:   437,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.1-L439",
		Level:   ShallNot,
		Summary: "SCIM message schemas shall not be discoverable via /Schemas",
		Source: Source{
			RFC:       7644,
			Section:   "3.1",
			StartLine: 438,
			EndLine:   440,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.2 - SCIM Endpoints and HTTP Methods

	{
		ID:      "RFC7644-3.2-L490",
		Level:   MustNot,
		Summary: "PUT must not be used to create new resources",
		Source: Source{
			RFC:       7644,
			Section:   "3.2",
			StartLine: 489,
			EndLine:   490,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.3 - Creating Resources

	{
		ID:      "RFC7644-3.3-L584",
		Level:   Shall,
		Summary: "readOnly attributes in POST request body shall be ignored",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: 583,
			EndLine:   584,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L602",
		Level:   Shall,
		Summary: "Successful resource creation returns HTTP 201 Created",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: 601,
			EndLine:   602,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L603",
		Level:   Should,
		Summary: "Response body should contain the newly created resource representation",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: 603,
			EndLine:   604,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L605",
		Level:   Shall,
		Summary: "Created resource URI shall be in Location header and meta.location",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: 604,
			EndLine:   607,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L625",
		Level:   Must,
		Summary: "Duplicate resource creation must return 409 Conflict with uniqueness scimType",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: 623,
			EndLine:   627,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.3.1 - Resource Types

	{
		ID:      "RFC7644-3.3.1-L712",
		Level:   Shall,
		Summary: "meta.resourceType shall be set by service provider to match endpoint",
		Source: Source{
			RFC:       7644,
			Section:   "3.3.1",
			StartLine: 711,
			EndLine:   715,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.4.2 - Query Resources

	{
		ID:      "RFC7644-3.4.2-L817",
		Level:   Must,
		Summary: "Query responses must use ListResponse schema URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2",
			StartLine: 817,
			EndLine:   818,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2-L847",
		Level:   Shall,
		Summary: "Empty query result returns 200 with totalResults 0",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2",
			StartLine: 847,
			EndLine:   848,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.4.2.1 - Query Endpoints

	{
		ID:      "RFC7644-3.4.2.1-L912",
		Level:   Shall,
		Summary: "Too many results shall return 400 with scimType tooMany",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.1",
			StartLine: 910,
			EndLine:   914,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.1-L918",
		Level:   Must,
		Summary: "Multi-resource-type queries must process filters per Section 3.4.2.2",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.1",
			StartLine: 916,
			EndLine:   918,
		},
		Feature:  Filter,
		Testable: true,
	},

	// Section 3.4.2.2 - Filtering

	{
		ID:      "RFC7644-3.4.2.2-L944",
		Level:   Must,
		Summary: "Filter parameter must contain at least one valid expression",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: 944,
			EndLine:   946,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1001",
		Level:   Shall,
		Summary: "Boolean/Binary gt comparison shall fail with 400 invalidFilter",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: 1001,
			EndLine:   1004,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1127",
		Level:   Must,
		Summary: "SCIM filters must conform to the ABNF rules",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: 1127,
			EndLine:   1128,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1197",
		Level:   Must,
		Summary: "Complex attributes must use fully qualified sub-attribute in filters",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: 1197,
			EndLine:   1198,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1208",
		Level:   Must,
		Summary: "Unrecognized filter operations must return 400 with invalidFilter",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: 1208,
			EndLine:   1211,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1217",
		Level:   Shall,
		Summary: "String filter comparisons shall respect caseExact characteristic",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: 1216,
			EndLine:   1218,
		},
		Feature:  Filter,
		Testable: true,
	},

	// Section 3.4.2.3 - Sorting

	{
		ID:      "RFC7644-3.4.2.3-L1319",
		Level:   Shall,
		Summary: "sortOrder shall default to ascending when sortBy is provided",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.3",
			StartLine: 1318,
			EndLine:   1319,
		},
		Feature:  Sort,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.3-L1321",
		Level:   Must,
		Summary: "sortOrder must sort according to attribute type (case-insensitive or case-exact)",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.3",
			StartLine: 1321,
			EndLine:   1326,
		},
		Feature:  Sort,
		Testable: true,
	},

	// Section 3.4.2.4 - Pagination

	{
		ID:      "RFC7644-3.4.2.4-L1358",
		Level:   Shall,
		Summary: "startIndex value less than 1 shall be interpreted as 1",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.4",
			StartLine: 1357,
			EndLine:   1359,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.4-L1363",
		Level:   MustNot,
		Summary: "Service provider must not return more results than count specifies",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.4",
			StartLine: 1363,
			EndLine:   1364,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.4.2.5 - Attributes

	{
		ID:      "RFC7644-3.4.2.5-L1438",
		Level:   Must,
		Summary: "attributes and excludedAttributes parameters must be supported",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.5",
			StartLine: 1437,
			EndLine:   1439,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.5-L1449",
		Level:   Shall,
		Summary: "excludedAttributes shall have no effect on returned:always attributes",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.5",
			StartLine: 1449,
			EndLine:   1451,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.4.3 - Querying Resources Using HTTP POST

	{
		ID:      "RFC7644-3.4.3-L1476",
		Level:   Must,
		Summary: "POST search requests must use SearchRequest schema URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.3",
			StartLine: 1476,
			EndLine:   1477,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.5 - Modifying Resources

	{
		ID:      "RFC7644-3.5-L1607",
		Level:   Must,
		Summary: "Implementers must support HTTP PUT",
		Source: Source{
			RFC:       7644,
			Section:   "3.5",
			StartLine: 1607,
			EndLine:   1608,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5-L1609",
		Level:   Should,
		Summary: "Implementers should support HTTP PATCH for partial modifications",
		Source: Source{
			RFC:       7644,
			Section:   "3.5",
			StartLine: 1608,
			EndLine:   1610,
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.1 - Replacing with PUT

	{
		ID:      "RFC7644-3.5.1-L1637",
		Level:   MustNot,
		Summary: "HTTP PUT must not be used to create new resources",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: 1636,
			EndLine:   1637,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1644",
		Level:   Shall,
		Summary: "readWrite/writeOnly values provided in PUT shall replace existing values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: 1644,
			EndLine:   1645,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1660",
		Level:   Must,
		Summary: "Immutable attribute input must match existing values or return 400",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: 1659,
			EndLine:   1663,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1665",
		Level:   Shall,
		Summary: "readOnly values provided in PUT shall be ignored",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: 1665,
			EndLine:   1665,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1667",
		Level:   Must,
		Summary: "Required attributes must be specified in PUT request",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: 1667,
			EndLine:   1668,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.5.2 - Modifying with PATCH

	{
		ID:      "RFC7644-3.5.2-L1805",
		Level:   Must,
		Summary: "PATCH body must contain schemas attribute with PatchOp URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1805,
			EndLine:   1806,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1808",
		Level:   Must,
		Summary: "PATCH body must contain Operations array with one or more operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1808,
			EndLine:   1810,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1810",
		Level:   Must,
		Summary: "Each PATCH operation must have exactly one op member",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1810,
			EndLine:   1813,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1886",
		Level:   Must,
		Summary: "Each PATCH operation must be compatible with attribute mutability and schema",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1886,
			EndLine:   1888,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1888",
		Level:   MustNot,
		Summary: "Client must not modify readOnly or immutable attributes via PATCH",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1888,
			EndLine:   1889,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1927",
		Level:   Shall,
		Summary: "Setting primary to true shall auto-set primary to false for other values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1926,
			EndLine:   1929,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1931",
		Level:   Shall,
		Summary: "PATCH request shall be treated as atomic",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1931,
			EndLine:   1932,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1933",
		Level:   Must,
		Summary: "On PATCH operation error, original resource must be restored",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1932,
			EndLine:   1934,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1939",
		Level:   Must,
		Summary: "Successful PATCH must return 200 with resource body or 204 No Content",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: 1939,
			EndLine:   1944,
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.2.1 - Add Operation

	{
		ID:      "RFC7644-3.5.2.1-L1972",
		Level:   Must,
		Summary: "PATCH add operation must contain a value member",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: 1972,
			EndLine:   1973,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.1-L1988",
		Level:   Shall,
		Summary: "Add to complex attribute shall specify sub-attributes in value parameter",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: 1987,
			EndLine:   1988,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.1-L2004",
		Level:   ShallNot,
		Summary: "Add of already-existing value shall not change the modify timestamp",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: 2001,
			EndLine:   2005,
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.2.2 - Remove Operation

	{
		ID:      "RFC7644-3.5.2.2-L2146",
		Level:   Shall,
		Summary: "Removing single-value attribute shall make it unassigned",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: 2145,
			EndLine:   2147,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.2-L2151",
		Level:   Shall,
		Summary: "Removing all values of multi-valued attribute shall make it unassigned",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: 2149,
			EndLine:   2151,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.2-L2167",
		Level:   Shall,
		Summary: "Removing required or readOnly attribute shall return error with mutability scimType",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: 2166,
			EndLine:   2170,
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.2.3 - Replace Operation

	{
		ID:      "RFC7644-3.5.2.3-L2376",
		Level:   Shall,
		Summary: "Replace on non-existent attribute shall treat as add",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: 2375,
			EndLine:   2376,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.3-L2387",
		Level:   Shall,
		Summary: "Replace with valuePath filter shall replace all matching record values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: 2384,
			EndLine:   2387,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.3-L2396",
		Level:   Shall,
		Summary: "Replace with unmatched valuePath filter shall return 400 noTarget",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: 2394,
			EndLine:   2398,
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.6 - Deleting Resources

	{
		ID:      "RFC7644-3.6-L2642",
		Level:   Must,
		Summary: "Deleted resource must return 404 for all subsequent operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: 2641,
			EndLine:   2644,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.6-L2644",
		Level:   Must,
		Summary: "Deleted resource must be omitted from future query results",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: 2644,
			EndLine:   2645,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.6-L2657",
		Level:   Shall,
		Summary: "Successful DELETE returns 204 No Content",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: 2657,
			EndLine:   2658,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.7 - Bulk Operations

	{
		ID:      "RFC7644-3.7-L2775",
		Level:   Must,
		Summary: "Successful bulk job must return 200 OK",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: 2775,
			EndLine:   2777,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7-L2779",
		Level:   Must,
		Summary: "Service provider must continue performing changes and disregard partial failures",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: 2779,
			EndLine:   2780,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7-L2789",
		Level:   Must,
		Summary: "Service provider must return the same bulkId with newly created resource",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: 2789,
			EndLine:   2790,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7-L2796",
		Level:   Must,
		Summary: "Optimized bulk processing must preserve client intent and achieve same result",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: 2795,
			EndLine:   2798,
		},
		Feature:  Bulk,
		Testable: false,
	},

	// Section 3.7.1 - Circular Reference Processing

	{
		ID:      "RFC7644-3.7.1-L2814",
		Level:   Must,
		Summary: "Service provider must try to resolve circular cross-references in bulk",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.1",
			StartLine: 2814,
			EndLine:   2816,
		},
		Feature:  Bulk,
		Testable: true,
	},

	// Section 3.7.3 - Response and Error Handling

	{
		ID:      "RFC7644-3.7.3-L3201",
		Level:   Must,
		Summary: "Bulk response must include result of all processed operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.3",
			StartLine: 3201,
			EndLine:   3204,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7.3-L3206",
		Level:   Must,
		Summary: "Bulk status must include HTTP response code",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.3",
			StartLine: 3205,
			EndLine:   3210,
		},
		Feature:  Bulk,
		Testable: true,
	},

	// Section 3.7.4 - Maximum Operations

	{
		ID:      "RFC7644-3.7.4-L3495",
		Level:   Must,
		Summary: "Service provider must define max operations and max payload size for bulk",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.4",
			StartLine: 3495,
			EndLine:   3496,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7.4-L3499",
		Level:   Must,
		Summary: "Exceeding bulk limits must return 413 Payload Too Large",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.4",
			StartLine: 3498,
			EndLine:   3501,
		},
		Feature:  Bulk,
		Testable: true,
	},

	// Section 3.8 - Data Input/Output Formats

	{
		ID:      "RFC7644-3.8-L3537",
		Level:   Must,
		Summary: "Servers must accept and return JSON using UTF-8 encoding",
		Source: Source{
			RFC:       7644,
			Section:   "3.8",
			StartLine: 3537,
			EndLine:   3538,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.8-L3551",
		Level:   Must,
		Summary: "Service providers must support Accept: application/scim+json",
		Source: Source{
			RFC:       7644,
			Section:   "3.8",
			StartLine: 3551,
			EndLine:   3552,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.9 - Additional Operation Response Parameters

	{
		ID:      "RFC7644-3.9-L3596",
		Level:   Must,
		Summary: "attributes parameter overrides default list and must include minimum set",
		Source: Source{
			RFC:       7644,
			Section:   "3.9",
			StartLine: 3596,
			EndLine:   3600,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.11 - /Me Authenticated Subject Alias

	{
		ID:      "RFC7644-3.11-L3683",
		Level:   Must,
		Summary: "/Me response Location header must be the permanent resource location",
		Source: Source{
			RFC:       7644,
			Section:   "3.11",
			StartLine: 3682,
			EndLine:   3685,
		},
		Feature:  Me,
		Testable: true,
	},

	// Section 3.12 - HTTP Status and Error Response Handling

	{
		ID:      "RFC7644-3.12-L3707",
		Level:   Must,
		Summary: "Errors must be returned in JSON format in the response body",
		Source: Source{
			RFC:       7644,
			Section:   "3.12",
			StartLine: 3707,
			EndLine:   3709,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.13 - SCIM Protocol Versioning

	{
		ID:      "RFC7644-3.13-L3956",
		Level:   Must,
		Summary: "When version is specified, service provider must use it or reject the request",
		Source: Source{
			RFC:       7644,
			Section:   "3.13",
			StartLine: 3956,
			EndLine:   3957,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.14 - Versioning Resources

	{
		ID:      "RFC7644-3.14-L3967",
		Level:   Must,
		Summary: "When supported, ETags must be specified as an HTTP header",
		Source: Source{
			RFC:       7644,
			Section:   "3.14",
			StartLine: 3967,
			EndLine:   3969,
		},
		Feature:  ETag,
		Testable: true,
	},

	// Section 4 - Service Provider Configuration Endpoints

	{
		ID:      "RFC7644-4-L4068",
		Level:   Shall,
		Summary: "/ServiceProviderConfig shall return ServiceProviderConfig schema",
		Source: Source{
			RFC:       7644,
			Section:   "4",
			StartLine: 4068,
			EndLine:   4070,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-4-L4098",
		Level:   Shall,
		Summary: "/Schemas shall return all supported schemas in ListResponse format",
		Source: Source{
			RFC:       7644,
			Section:   "4",
			StartLine: 4097,
			EndLine:   4099,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-4-L4125",
		Level:   Shall,
		Summary: "Query parameters (filter, sort, pagination) shall be ignored on discovery endpoints",
		Source: Source{
			RFC:       7644,
			Section:   "4",
			StartLine: 4124,
			EndLine:   4125,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 5 - Internationalized Strings

	{
		ID:      "RFC7644-5-L4214",
		Level:   Must,
		Summary: "Service providers must use PRECIS rules for userName and password comparison",
		Source: Source{
			RFC:       7644,
			Section:   "5",
			StartLine: 4213,
			EndLine:   4218,
		},
		Feature:  Core,
		Testable: false,
	},

	// Section 7.2 - TLS Support

	{
		ID:      "RFC7644-7.2-L4348",
		Level:   Must,
		Summary: "Clients and service providers must require transport-layer security",
		Source: Source{
			RFC:       7644,
			Section:   "7.2",
			StartLine: 4347,
			EndLine:   4349,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7644-7.2-L4350",
		Level:   Must,
		Summary: "Service provider must support TLS 1.2",
		Source: Source{
			RFC:       7644,
			Section:   "7.2",
			StartLine: 4349,
			EndLine:   4351,
		},
		Feature:  Core,
		Testable: false,
	},
}
