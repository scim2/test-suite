package testserver

import (
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

const (
	testResourceSchemaURI  = "urn:ietf:params:scim:schemas:test:2.0:TestResource"
	testExtensionSchemaURI = "urn:ietf:params:scim:schemas:extension:test:2.0:TestResource"
)

// testExtensionSchema defines an extension schema for TestResource.
func testExtensionSchema() schema.Schema {
	return schema.Schema{
		ID:          testExtensionSchemaURI,
		Name:        optional.NewString("TestExtension"),
		Description: optional.NewString("Extension schema for compliance testing"),
		Attributes: schema.Attributes{
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name: "extString",
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:     "extRequired",
				Required: true,
			})),
			schema.SimpleCoreAttribute(schema.SimpleNumberParams(schema.NumberParams{
				Name: "extInteger",
				Type: schema.AttributeTypeInteger(),
			})),
		},
	}
}

// testResourceSchema defines a schema that covers every attribute type,
// mutability, returned, and uniqueness combination so the compliance
// suite can exercise requirements that User/Group alone cannot reach.
func testResourceSchema() schema.Schema {
	return schema.Schema{
		ID:          testResourceSchemaURI,
		Name:        optional.NewString("TestResource"),
		Description: optional.NewString("Synthetic resource for compliance testing"),
		Attributes: schema.Attributes{
			// required string (like userName for Users)
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "identifier",
				Required:   true,
				Uniqueness: schema.AttributeUniquenessServer(),
			})),

			// string, caseExact=true
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:      "caseExactString",
				CaseExact: true,
			})),

			// string, caseExact=false (default)
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name: "caseInsensitiveString",
			})),

			// integer
			schema.SimpleCoreAttribute(schema.SimpleNumberParams(schema.NumberParams{
				Name: "integerAttr",
				Type: schema.AttributeTypeInteger(),
			})),

			// boolean
			schema.SimpleCoreAttribute(schema.SimpleBooleanParams(schema.BooleanParams{
				Name: "booleanAttr",
			})),

			// dateTime
			schema.SimpleCoreAttribute(schema.SimpleDateTimeParams(schema.DateTimeParams{
				Name: "dateTimeAttr",
			})),

			// binary
			schema.SimpleCoreAttribute(schema.SimpleBinaryParams(schema.BinaryParams{
				Name: "binaryAttr",
			})),

			// reference (URI)
			schema.SimpleCoreAttribute(schema.SimpleReferenceParams(schema.ReferenceParams{
				Name:           "referenceAttr",
				ReferenceTypes: []schema.AttributeReferenceType{schema.AttributeReferenceTypeURI},
			})),

			// readOnly string (server-assigned)
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "readOnlyAttr",
				Mutability: schema.AttributeMutabilityReadOnly(),
			})),

			// immutable string (set at creation, cannot be changed)
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "immutableAttr",
				Mutability: schema.AttributeMutabilityImmutable(),
			})),

			// writeOnly string (like password, returned:never)
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "writeOnlyAttr",
				Mutability: schema.AttributeMutabilityWriteOnly(),
				Returned:   schema.AttributeReturnedNever(),
			})),

			// returned:request (only returned when explicitly requested)
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:     "requestOnlyAttr",
				Returned: schema.AttributeReturnedRequest(),
			})),

			// returned:always
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:     "alwaysReturnedAttr",
				Returned: schema.AttributeReturnedAlways(),
			})),

			// multi-valued strings
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:        "multiStrings",
				MultiValued: true,
			})),

			// complex attribute with sub-attributes
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name: "complexAttr",
				SubAttributes: []schema.SimpleParams{
					schema.SimpleStringParams(schema.StringParams{
						Name: "sub1",
					}),
					schema.SimpleNumberParams(schema.NumberParams{
						Name: "sub2",
						Type: schema.AttributeTypeInteger(),
					}),
					schema.SimpleBooleanParams(schema.BooleanParams{
						Name: "sub3",
					}),
				},
			}),

			// multi-valued complex (like emails) with value/type/primary
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name:        "multiComplex",
				MultiValued: true,
				SubAttributes: []schema.SimpleParams{
					schema.SimpleStringParams(schema.StringParams{
						Name: "value",
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name: "type",
					}),
					schema.SimpleBooleanParams(schema.BooleanParams{
						Name: "primary",
					}),
				},
			}),

			// server-unique string (for uniqueness testing)
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "uniqueString",
				Uniqueness: schema.AttributeUniquenessServer(),
			})),
		},
	}
}
