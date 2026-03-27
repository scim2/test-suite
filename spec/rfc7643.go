package spec

// RFC7643 contains requirements from RFC 7643 (SCIM Core Schema).
var RFC7643 = []Requirement{

	// Section 1.1 - Requirements Notation and Conventions

	{
		ID:    "RFC7643-1.1-L232",
		Level: MustNot,
		Summary: "Quoted values in the spec must not include the quotes " +
			"in protocol messages",
		Source: Source{
			RFC:       7643,
			Section:   "1.1",
			StartLine: 231,
			EndLine:   233,
		},
		Feature:  Core,
		Testable: false,
	},

	// Section 2.1 - Attributes

	{
		ID:      "RFC7643-2.1-L366",
		Level:   Must,
		Summary: "Resources must be JSON and specify schema via the schemas attribute",
		Source: Source{
			RFC:       7643,
			Section:   "2.1",
			StartLine: 365,
			EndLine:   367,
			StartCol:  25,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.1-L369",
		Level:   Must,
		Summary: "Attribute names must conform to ATTRNAME ABNF rules",
		Source: Source{
			RFC:       7643,
			Section:   "2.1",
			StartLine: 369,
			EndLine:   374,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3 - Attribute Data Types

	{
		ID:      "RFC7643-2.3-L443",
		Level:   ShouldNot,
		Summary: "Extensions should not introduce new data types",
		Source: Source{
			RFC:       7643,
			Section:   "2.3",
			StartLine: 443,
			EndLine:   444,
		},
		Feature:  Core,
		Testable: false,
	},

	// Section 2.3.4 - Integer

	{
		ID:      "RFC7643-2.3.4-L521",
		Level:   MustNot,
		Summary: "Integer values must not contain fractional or exponent parts",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.4",
			StartLine: 519,
			EndLine:   522,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.5 - DateTime

	{
		ID:      "RFC7643-2.3.5-L527",
		Level:   Must,
		Summary: "DateTime values must be valid xsd:dateTime with both date and time",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.5",
			StartLine: 526,
			EndLine:   529,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.5-L531",
		Level:   Must,
		Summary: "DateTime JSON values must conform to XML constraints",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.5",
			StartLine: 531,
			EndLine:   533,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.6 - Binary

	{
		ID:      "RFC7643-2.3.6-L537",
		Level:   Must,
		Summary: "Binary attribute values must be base64 encoded per RFC 4648",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.6",
			StartLine: 536,
			EndLine:   538,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.7 - Reference

	{
		ID:      "RFC7643-2.3.7-L552",
		Level:   Must,
		Summary: "Reference values must be absolute or relative URIs",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: 551,
			EndLine:   553,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.7-L555",
		Level:   Must,
		Summary: "Base URI for relative URI resolution must include all components up to endpoint URI",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: 553,
			EndLine:   570,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.7-L577",
		Level:   Must,
		Summary: "Reference URI must be HTTP-addressable",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: 577,
			EndLine:   578,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.7-L578",
		Level:   Must,
		Summary: "GET on a reference URI must return the target resource or appropriate HTTP code",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: 578,
			EndLine:   580,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.8 - Complex

	{
		ID:      "RFC7643-2.3.8-L592",
		Level:   MustNot,
		Summary: "Servers and clients must not require or expect attributes in any specific order",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.8",
			StartLine: 591,
			EndLine:   594,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.8-L595",
		Level:   MustNot,
		Summary: "Complex attributes must not contain sub-attributes that are themselves complex",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.8",
			StartLine: 594,
			EndLine:   596,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.4 - Multi-Valued Attributes

	{
		ID:      "RFC7643-2.4-L608",
		Level:   Shall,
		Summary: "Multi-valued attribute objects with sub-attributes shall be considered complex",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: 606,
			EndLine:   608,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-2.4-L612",
		Level:   Must,
		Summary: "Predefined multi-valued sub-attributes must be used with the meanings defined in spec",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: 608,
			EndLine:   612,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-2.4-L634",
		Level:   Must,
		Summary: "The primary attribute value true must appear no more than once",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: 633,
			EndLine:   635,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.4-L655",
		Level:   Should,
		Summary: "Service providers should canonicalize multi-valued attribute values when returning",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: 655,
			EndLine:   658,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.4-L663",
		Level:   ShouldNot,
		Summary: "Should not return the same (type, value) combination more than once per attribute",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: 660,
			EndLine:   664,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.5 - Unassigned and Null Values

	{
		ID:      "RFC7643-2.5-L682",
		Level:   Shall,
		Summary: "Unassigned attributes, null, and empty array are equivalent in state",
		Source: Source{
			RFC:       7643,
			Section:   "2.5",
			StartLine: 680,
			EndLine:   684,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3 - SCIM Resources

	{
		ID:      "RFC7643-3-L709",
		Level:   Must,
		Summary: "All schema representations must include a non-empty schemas array",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: 709,
			EndLine:   711,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3-L711",
		Level:   Must,
		Summary: "Schemas attribute must only contain values defined as schema and schemaExtensions",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: 711,
			EndLine:   713,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3-L713",
		Level:   MustNot,
		Summary: "Duplicate values must not be included in the schemas attribute",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: 713,
			EndLine:   714,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3-L714",
		Level:   MustNot,
		Summary: "Value order in schemas attribute must not impact behavior",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: 714,
			EndLine:   715,
		},
		Feature:  Core,
		Testable: false,
	},

	// Section 3.1 - Common Attributes

	{
		ID:      "RFC7643-3.1-L852",
		Level:   Must,
		Summary: "Common attributes must be defined for all resources (except ServiceProviderConfig and ResourceType)",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 848,
			EndLine:   853,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L855",
		Level:   Must,
		Summary: "Service provider must assign id and meta attributes upon creation",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 853,
			EndLine:   857,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L866",
		Level:   Must,
		Summary: "Each resource representation must include a non-empty id value",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 865,
			EndLine:   867,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L867",
		Level:   Must,
		Summary: "Resource id must be unique across the entire set of resources",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 867,
			EndLine:   868,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L868",
		Level:   Must,
		Summary: "Resource id must be stable and non-reassignable",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 868,
			EndLine:   870,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L872",
		Level:   MustNot,
		Summary: "Resource id must not be specified by the client",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 871,
			EndLine:   872,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L873",
		Level:   MustNot,
		Summary: "The string bulkId must not be used within any unique identifier value",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 872,
			EndLine:   875,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L890",
		Level:   MustNot,
		Summary: "externalId must not be specified by the service provider",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 889,
			EndLine:   891,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L891",
		Level:   Must,
		Summary: "Service provider must interpret externalId as scoped to the provisioning domain",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 891,
			EndLine:   892,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-3.1-L911",
		Level:   Shall,
		Summary: "The meta attribute shall be ignored when provided by clients",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 910,
			EndLine:   912,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L920",
		Level:   Must,
		Summary: "meta.created must be a DateTime",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 918,
			EndLine:   920,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L925",
		Level:   Must,
		Summary: "meta.lastModified must equal meta.created if resource was never modified",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 921,
			EndLine:   925,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L927",
		Level:   Must,
		Summary: "meta.location must match Content-Location header",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 926,
			EndLine:   929,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.1-L940",
		Level:   Must,
		Summary: "Weak entity-tags must be prefixed with W/",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: 935,
			EndLine:   942,
		},
		Feature:  ETag,
		Testable: true,
	},

	// Section 3.3 - Attribute Extensions to Resources

	{
		ID:      "RFC7643-3.3-L976",
		Level:   Must,
		Summary: "The schemas attribute must contain at least one value (the base schema)",
		Source: Source{
			RFC:       7643,
			Section:   "3.3",
			StartLine: 975,
			EndLine:   978,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.3-L981",
		Level:   Shall,
		Summary: "Extension schema URI shall be used as JSON container for extension attributes",
		Source: Source{
			RFC:       7643,
			Section:   "3.3",
			StartLine: 979,
			EndLine:   985,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 4.1.1 - User Singular Attributes

	{
		ID:      "RFC7643-4.1.1-L1035",
		Level:   Must,
		Summary: "Each User must include a non-empty userName value",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: 1035,
			EndLine:   1038,
		},
		Feature:  Users,
		Testable: true,
	},

	// Section 4.1.1 - Password

	{
		ID:      "RFC7643-4.1.1-L1186",
		Level:   ShallNot,
		Summary: "Cleartext or hashed password values shall not be returnable",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: 1184,
			EndLine:   1187,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.1-L1222",
		Level:   MustNot,
		Summary: "Password value must not be returned in any form (returned is never)",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: 1220,
			EndLine:   1224,
		},
		Feature:  Users,
		Testable: true,
	},

	// Section 4.1.2 - User Multi-Valued Attributes

	{
		ID:      "RFC7643-4.1.2-L1244",
		Level:   Should,
		Summary: "Email values should be specified per RFC 5321",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: 1243,
			EndLine:   1246,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1281",
		Level:   Must,
		Summary: "Photo URI resource must be an image file, not a web page",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: 1280,
			EndLine:   1283,
		},
		Feature:  Users,
		Testable: false,
	},
	{
		ID:      "RFC7643-4.1.2-L1322",
		Level:   Must,
		Summary: "Country value must be in ISO 3166-1 alpha-2 code format",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: 1320,
			EndLine:   1323,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1351",
		Level:   Must,
		Summary: "Group value sub-attribute must be the id of the Group resource",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: 1342,
			EndLine:   1356,
		},
		Feature:  Groups,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1354",
		Level:   Must,
		Summary: "Group membership changes must be applied via the Group resource",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: 1353,
			EndLine:   1356,
		},
		Feature:  Groups,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1378",
		Level:   Must,
		Summary: "X.509 certificate values must be base64 encoded DER",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: 1375,
			EndLine:   1379,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1379",
		Level:   MustNot,
		Summary: "A single certificate value must not contain multiple certificates",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: 1379,
			EndLine:   1381,
		},
		Feature:  Users,
		Testable: false,
	},

	// Section 5 - Service Provider Configuration Schema

	{
		ID:      "RFC7643-5-L1581",
		Level:   Should,
		Summary: "authenticationSchemes should be publicly accessible without prior authentication",
		Source: Source{
			RFC:       7643,
			Section:   "5",
			StartLine: 1580,
			EndLine:   1583,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 6 - ResourceType Schema

	{
		ID:      "RFC7643-6-L1619",
		Level:   Must,
		Summary: "Service providers must specify the resource type name",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: 1619,
			EndLine:   1621,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-6-L1641",
		Level:   Must,
		Summary: "ResourceType schema URI must equal the id of the associated Schema resource",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: 1639,
			EndLine:   1642,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-6-L1650",
		Level:   Must,
		Summary: "Extension schema URI must equal the id of a Schema resource",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: 1648,
			EndLine:   1651,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-6-L1655",
		Level:   Must,
		Summary: "If schema extension is required, resource must include it and its required attributes",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: 1652,
			EndLine:   1658,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 7 - Schema Definition

	{
		ID:      "RFC7643-7-L1691",
		Level:   Must,
		Summary: "Service providers must specify the schema URI",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: 1689,
			EndLine:   1695,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1700",
		Level:   Must,
		Summary: "Service providers must specify the schema name when applicable",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: 1699,
			EndLine:   1701,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1748",
		Level:   Shall,
		Summary: "Server shall use case sensitivity when evaluating filters per caseExact",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: 1747,
			EndLine:   1749,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1750",
		Level:   Shall,
		Summary: "Server shall preserve case for case-exact attribute values",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: 1749,
			EndLine:   1751,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1759",
		Level:   ShallNot,
		Summary: "readOnly attributes shall not be modified",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: 1759,
			EndLine:   1759,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1766",
		Level:   ShallNot,
		Summary: "Immutable attributes shall not be updated after creation",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: 1762,
			EndLine:   1766,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1769",
		Level:   ShallNot,
		Summary: "writeOnly attribute values shall not be returned",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: 1767,
			EndLine:   1770,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 9.2 - Passwords and Other Sensitive Security Data

	{
		ID:      "RFC7643-9.2-L5117",
		Level:   MustNot,
		Summary: "Password values must not be stored in cleartext form",
		Source: Source{
			RFC:       7643,
			Section:   "9.2",
			StartLine: 5114,
			EndLine:   5117,
		},
		Feature:  Core,
		Testable: false,
	},
}
