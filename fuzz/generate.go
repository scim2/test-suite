package fuzz

import (
	"math/rand/v2"
	"strings"

	"github.com/elimity-com/scim/schema"
)

// Generate produces a map[string]any resource for the given schema.
// With no options it generates a valid client-side resource: required
// non-readOnly attributes get random values, all optional non-readOnly
// attributes are included, and readOnly attributes are skipped.
func Generate(s schema.Schema, opts ...Option) map[string]any {
	c := newConfig(opts)
	result := map[string]any{
		"schemas": []string{s.ID},
	}
	generateAttrs(s.Attributes, c, "", result)
	return result
}

func generateAttrs(attrs schema.Attributes, c *config, prefix string, out map[string]any) {
	for _, attr := range attrs {
		name := attr.Name()
		path := strings.ToLower(name)
		if prefix != "" {
			path = prefix + "." + path
		}

		// Check for a full-attribute override.
		if v, ok := c.overrides[path]; ok {
			out[name] = v
			continue
		}

		if !shouldInclude(attr, c) {
			continue
		}

		out[name] = generateValue(attr, c, path)
	}
}

func generateComplex(attr schema.CoreAttribute, c *config, path string) any {
	sub := map[string]any{}
	for _, sa := range attr.SubAttributes() {
		saName := sa.Name()
		saPath := path + "." + strings.ToLower(saName)

		if v, ok := c.overrides[saPath]; ok {
			sub[saName] = v
			continue
		}

		if !shouldInclude(sa, c) {
			continue
		}

		sub[saName] = generateSingular(sa, c, saPath)
	}
	return sub
}

func generateSingular(attr schema.CoreAttribute, c *config, path string) any {
	typ := attr.AttributeType()

	if fn, ok := c.generators[typ]; ok {
		return fn(attr, c.rand)
	}

	switch typ {
	case "string":
		return defaultStringValue(attr, c.rand)
	case "integer":
		return defaultIntegerValue(attr, c.rand)
	case "decimal":
		return defaultDecimalValue(attr, c.rand)
	case "boolean":
		return defaultBooleanValue(attr, c.rand)
	case "binary":
		return defaultBinaryValue(attr, c.rand)
	case "dateTime":
		return defaultDateTimeValue(attr, c.rand)
	case "reference":
		return defaultReferenceValue(attr, c.rand)
	case "complex":
		return generateComplex(attr, c, path)
	default:
		return defaultStringValue(attr, c.rand)
	}
}

func generateValue(attr schema.CoreAttribute, c *config, path string) any {
	if attr.MultiValued() {
		n := 1 + c.rand.IntN(3)
		elems := make([]any, n)
		for i := range elems {
			elems[i] = generateSingular(attr, c, path)
		}
		// At most one element may have primary=true (RFC 7643 Section 2.4).
		if attr.HasPrimarySubAttr() {
			primarySet := false
			for _, elem := range elems {
				m, ok := elem.(map[string]any)
				if !ok {
					continue
				}
				if primarySet {
					m["primary"] = false
				} else if p, _ := m["primary"].(bool); p {
					primarySet = true
				}
			}
		}
		return elems
	}
	return generateSingular(attr, c, path)
}

func shouldInclude(attr schema.CoreAttribute, c *config) bool {
	if attr.Mutability() == "readOnly" && !c.includeReadOnly {
		return false
	}
	if !attr.Required() && c.requiredOnly {
		return false
	}
	if c.filter != nil && !c.filter(attr) {
		return false
	}
	return true
}

// Option configures resource generation.
type Option func(*config)

// AllOptional forces all optional attributes to be included.
func AllOptional() Option {
	return func(c *config) { c.allOptional = true }
}

// IncludeReadOnly includes readOnly attributes with generated values.
// By default they are skipped.
func IncludeReadOnly() Option {
	return func(c *config) { c.includeReadOnly = true }
}

// RequiredOnly limits generation to required attributes only.
func RequiredOnly() Option {
	return func(c *config) { c.requiredOnly = true }
}

// WithFilter provides a custom predicate controlling whether an attribute
// is included. It runs after the built-in readOnly/required logic.
// Return true to include, false to skip.
func WithFilter(fn func(schema.CoreAttribute) bool) Option {
	return func(c *config) { c.filter = fn }
}

// WithGenerator replaces the value generator for a given attribute type.
// The typ parameter matches CoreAttribute.AttributeType() return values:
// "string", "integer", "boolean", "decimal", "binary", "dateTime",
// "reference".
func WithGenerator(typ string, fn func(schema.CoreAttribute, *rand.Rand) any) Option {
	return func(c *config) { c.generators[typ] = fn }
}

// WithRand sets the random source for deterministic output.
func WithRand(r *rand.Rand) Option {
	return func(c *config) { c.rand = r }
}

// WithValue forces a specific attribute to a given value. The name is
// matched case-insensitively. Use dot notation for nested attributes
// (e.g. "emails.value"). The attribute is always included regardless
// of other options.
func WithValue(name string, value any) Option {
	return func(c *config) {
		c.overrides[strings.ToLower(name)] = value
	}
}

type config struct {
	rand            *rand.Rand
	includeReadOnly bool
	requiredOnly    bool
	allOptional     bool
	overrides       map[string]any
	generators      map[string]func(schema.CoreAttribute, *rand.Rand) any
	filter          func(schema.CoreAttribute) bool
}

func newConfig(opts []Option) *config {
	c := &config{
		rand:       rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
		overrides:  make(map[string]any),
		generators: make(map[string]func(schema.CoreAttribute, *rand.Rand) any),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}
