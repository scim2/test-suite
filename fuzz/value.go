package fuzz

import (
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/elimity-com/scim/schema"
)

func defaultBinaryValue(_ schema.CoreAttribute, r *rand.Rand) any {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(r.IntN(256))
	}
	return base64.StdEncoding.EncodeToString(b)
}

func defaultBooleanValue(_ schema.CoreAttribute, r *rand.Rand) any {
	return r.IntN(2) == 0
}

func defaultDateTimeValue(_ schema.CoreAttribute, r *rand.Rand) any {
	// Random timestamp within the last year.
	now := time.Now().UTC()
	offset := time.Duration(r.IntN(365*24)) * time.Hour
	t := now.Add(-offset)
	return t.Format(time.RFC3339)
}

func defaultDecimalValue(_ schema.CoreAttribute, r *rand.Rand) any {
	return r.Float64() * 100.0
}

func defaultIntegerValue(_ schema.CoreAttribute, r *rand.Rand) any {
	return r.IntN(1000)
}

func defaultReferenceValue(_ schema.CoreAttribute, r *rand.Rand) any {
	return fmt.Sprintf("https://example.com/%08x", r.Uint32())
}

func defaultStringValue(attr schema.CoreAttribute, r *rand.Rand) any {
	if cv := attr.CanonicalValues(); len(cv) > 0 {
		return cv[r.IntN(len(cv))]
	}
	n := 8 + r.IntN(9) // 8-16 chars
	b := make([]byte, n)
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range b {
		b[i] = chars[r.IntN(len(chars))]
	}
	return string(b)
}
