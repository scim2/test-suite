package testsuite

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/op"
)

var (
	// More info: https://tools.ietf.org/html/rfc5234#appendix-B
	// alpha = %x41-5A / %x61-7A ; A-Z / a-z
	alpha = op.Or{
		parser.CheckRuneRange('a', 'z'),
		parser.CheckRuneRange('A', 'Z'),
	}
	// digit = %x30-39 ; 0-9
	digit = parser.CheckRuneRange('0', '9')

	// More info: https://tools.ietf.org/html/rfc7643#section-2.1
	// nameChar = "$" / "-" / "_" / DIGIT / ALPHA
	nameChar = op.Or{'$', '-', '_', digit, alpha}
	// attrName = ALPHA *(nameChar)
	attrName = op.And{alpha, op.MinZero(nameChar)}
)

// IsAttributeName checks whether the given string is a valid attribute name.
// More info: https://tools.ietf.org/html/rfc7643#section-2.1
func IsAttributeName(s string) bool {
	if err := is(attrName, s); err != nil {
		return false
	}
	return true
}
