package fuzz

import (
	"math/rand/v2"
	"testing"

	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

func TestDeterministic(t *testing.T) {
	s := testSchema()
	r1 := rand.New(rand.NewPCG(99, 0))
	r2 := rand.New(rand.NewPCG(99, 0))

	res1 := Generate(s, WithRand(r1))
	res2 := Generate(s, WithRand(r2))

	// Required string should match.
	if res1["requiredStr"] != res2["requiredStr"] {
		t.Errorf("not deterministic: %v != %v", res1["requiredStr"], res2["requiredStr"])
	}
}

func TestGenerateDefault(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	res := Generate(testSchema(), WithRand(r))

	// Must have schemas key.
	schemas, ok := res["schemas"].([]string)
	if !ok || len(schemas) != 1 || schemas[0] != "urn:test:Schema" {
		t.Fatalf("schemas = %v, want [urn:test:Schema]", res["schemas"])
	}

	// Required field must be present.
	if _, ok := res["requiredStr"]; !ok {
		t.Error("requiredStr missing")
	}

	// ReadOnly field must be absent by default.
	if _, ok := res["readOnlyStr"]; ok {
		t.Error("readOnlyStr should be absent by default")
	}

	// All non-readOnly optional fields should be present (default is include all).
	for _, name := range []string{"optionalStr", "intAttr", "decAttr", "boolAttr", "dtAttr", "binAttr", "refAttr", "multiStr", "complexAttr", "canonicalStr"} {
		if _, ok := res[name]; !ok {
			t.Errorf("%s missing", name)
		}
	}

	// Multi-valued should be a slice.
	multi, ok := res["multiStr"].([]any)
	if !ok {
		t.Fatalf("multiStr is %T, want []any", res["multiStr"])
	}
	if len(multi) < 1 || len(multi) > 3 {
		t.Errorf("multiStr len = %d, want 1-3", len(multi))
	}

	// Complex should be a map.
	cplx, ok := res["complexAttr"].(map[string]any)
	if !ok {
		t.Fatalf("complexAttr is %T, want map[string]any", res["complexAttr"])
	}
	if _, ok := cplx["sub1"]; !ok {
		t.Error("complexAttr.sub1 missing")
	}
	if _, ok := cplx["sub2"]; !ok {
		t.Error("complexAttr.sub2 missing")
	}

	// Canonical values should pick from the list.
	cv, ok := res["canonicalStr"].(string)
	if !ok {
		t.Fatalf("canonicalStr is %T, want string", res["canonicalStr"])
	}
	valid := map[string]bool{"alpha": true, "beta": true, "gamma": true}
	if !valid[cv] {
		t.Errorf("canonicalStr = %q, want one of alpha/beta/gamma", cv)
	}
}

func TestGenerateValidatesAgainstSchema(t *testing.T) {
	s := testSchema()
	r := rand.New(rand.NewPCG(42, 0))
	res := Generate(s, WithRand(r))

	// The schema.Validate method expects []interface{} for schemas.
	res["schemas"] = []any{s.ID}

	if _, err := s.Validate(res); err != nil {
		t.Fatalf("generated resource failed schema validation: %v", err)
	}
}

func TestIncludeReadOnly(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	res := Generate(testSchema(), WithRand(r), IncludeReadOnly())

	if _, ok := res["readOnlyStr"]; !ok {
		t.Error("readOnlyStr should be present with IncludeReadOnly")
	}
}

func TestRequiredOnly(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	res := Generate(testSchema(), WithRand(r), RequiredOnly())

	if _, ok := res["requiredStr"]; !ok {
		t.Error("requiredStr missing")
	}
	if _, ok := res["optionalStr"]; ok {
		t.Error("optionalStr should be absent with RequiredOnly")
	}
}

func TestWithFilter(t *testing.T) {
	res := Generate(testSchema(),
		WithFilter(func(attr schema.CoreAttribute) bool {
			return attr.Name() == "requiredStr" || attr.Name() == "intAttr"
		}),
	)

	if _, ok := res["requiredStr"]; !ok {
		t.Error("requiredStr missing")
	}
	if _, ok := res["intAttr"]; !ok {
		t.Error("intAttr missing")
	}
	if _, ok := res["optionalStr"]; ok {
		t.Error("optionalStr should be filtered out")
	}
}

func TestWithGenerator(t *testing.T) {
	res := Generate(testSchema(),
		WithGenerator("string", func(_ schema.CoreAttribute, _ *rand.Rand) any {
			return "custom"
		}),
	)

	if v := res["requiredStr"]; v != "custom" {
		t.Errorf("requiredStr = %v, want custom", v)
	}
}

func TestWithValue(t *testing.T) {
	res := Generate(testSchema(), WithValue("requiredStr", "forced"))

	if v := res["requiredStr"]; v != "forced" {
		t.Errorf("requiredStr = %v, want forced", v)
	}
}

func TestWithValueNested(t *testing.T) {
	res := Generate(testSchema(), WithValue("complexAttr.sub1", "nested-override"))

	cplx, ok := res["complexAttr"].(map[string]any)
	if !ok {
		t.Fatalf("complexAttr is %T, want map[string]any", res["complexAttr"])
	}
	if v := cplx["sub1"]; v != "nested-override" {
		t.Errorf("complexAttr.sub1 = %v, want nested-override", v)
	}
}

func TestWithValueOverridesReadOnly(t *testing.T) {
	res := Generate(testSchema(), WithValue("readOnlyStr", "override"))

	if v := res["readOnlyStr"]; v != "override" {
		t.Errorf("readOnlyStr = %v, want override", v)
	}
}

func testSchema() schema.Schema {
	return schema.Schema{
		ID:          "urn:test:Schema",
		Name:        optional.NewString("Test"),
		Description: optional.NewString("Test schema for fuzzer"),
		Attributes: schema.Attributes{
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:     "requiredStr",
				Required: true,
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name: "optionalStr",
			})),
			schema.SimpleCoreAttribute(schema.SimpleNumberParams(schema.NumberParams{
				Name: "intAttr",
				Type: schema.AttributeTypeInteger(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleNumberParams(schema.NumberParams{
				Name: "decAttr",
			})),
			schema.SimpleCoreAttribute(schema.SimpleBooleanParams(schema.BooleanParams{
				Name: "boolAttr",
			})),
			schema.SimpleCoreAttribute(schema.SimpleDateTimeParams(schema.DateTimeParams{
				Name: "dtAttr",
			})),
			schema.SimpleCoreAttribute(schema.SimpleBinaryParams(schema.BinaryParams{
				Name: "binAttr",
			})),
			schema.SimpleCoreAttribute(schema.SimpleReferenceParams(schema.ReferenceParams{
				Name:           "refAttr",
				ReferenceTypes: []schema.AttributeReferenceType{schema.AttributeReferenceTypeURI},
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "readOnlyStr",
				Mutability: schema.AttributeMutabilityReadOnly(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:        "multiStr",
				MultiValued: true,
			})),
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name: "complexAttr",
				SubAttributes: []schema.SimpleParams{
					schema.SimpleStringParams(schema.StringParams{Name: "sub1"}),
					schema.SimpleNumberParams(schema.NumberParams{
						Name: "sub2",
						Type: schema.AttributeTypeInteger(),
					}),
				},
			}),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:            "canonicalStr",
				CanonicalValues: []string{"alpha", "beta", "gamma"},
			})),
		},
	}
}
