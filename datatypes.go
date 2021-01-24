package testsuite

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	"strings"
)

// This file is not necessary since Golang will do this for us.\

// is checks whether the given string (s) complies with the given parser rules (i).
func is(i interface{}, s string) error {
	p, err := ast.New([]byte(s))
	if err != nil {
		return err
	}
	if _, err := p.Expect(i); err != nil {
		return err
	}
	// Check whether we parsed the whole string.
	if _, err := p.Expect(parser.EOD); err != nil {
		return err
	}
	return nil
}

var (
	// More info: https://tools.ietf.org/html/rfc7159#section-7
	// str = quote *char quote
	str = op.And{quote, op.MinZero(char), quote}
	// char = unescaped /
	//        escape (
	//            %x22 /          ; "    quotation mark
	//            %x5C /          ; \    reverse solidus
	//            %x2F /          ; /    solidus
	//            %x62 /          ; b    backspace
	//            %x66 /          ; f    form feed
	//            %x6E /          ; n    line feed
	//            %x72 /          ; r    carriage return
	//            %x74 /          ; t    tab
	//            %x75 hexdigit ) ; uXXXX
	char = op.Or{
		unescaped,
		op.And{
			escape,
			op.Or{
				'"', '\\', '/',
				0x0062, 0x0066, 0x006E, 0x0072, 0x0074,
				op.And{'u', op.Repeat(4, hexdigit)},
			},
		},
	}
	// escape = %x5C ; \
	escape = '\\'
	// quote = %x22 ; "
	quote = '"'
	// unescaped = %x20-21 / %x23-5B / %x5D-10FFFF
	unescaped = op.Or{
		parser.CheckRuneRange(0x0020, 0x0021),
		parser.CheckRuneRange(0x0023, 0x005B),
		parser.CheckRuneRange(0x005D, 0x0010FFFF),
	}

	// More info: https://tools.ietf.org/html/rfc5234#appendix-B
	// hexdigit = digit / "A" / "B" / "C" / "D" / "E" / "F"
	hexdigit = op.Or{
		parser.CheckRuneRange('0', '9'),
		parser.CheckRuneRange('A', 'F'),
	}
)

// IsQuotedString checks whether the given string is a valid (quoted) JSON string.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.1
func IsQuotedString(s string) bool {
	if err := is(str, s); err != nil {
		return false
	}
	return true
}

// IsString checks whether the given string is a valid unquoted JSON string. A string is a sequence of zero or more
// Unicode characters encoded using UTF-8.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.1
func IsString(s string) bool {
	if err := is(op.MinZero(char), s); err != nil {
		return false
	}
	return true
}

var (
	// More info: https://tools.ietf.org/html/rfc7159#section-3
	// f = %x66.61.6c.73.65 ; false
	f = "false"
	// t = %x74.72.75.65 ; true
	t = "true"
	// n = %x6e.75.6c.6c ; null
	n = "null"
)

// IsBoolean checks whether the given string is a boolean value. A boolean has no case sensitivity.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.2
func IsBoolean(s string) bool {
	s = strings.ToLower(s) // For case sensitivity purposes.
	if err := is(op.Or{t, f}, s); err != nil {
		return false
	}
	return true
}

var (
	// More info: https://tools.ietf.org/html/rfc7159#section-6
	// decimal and integer are CUSTOM definition based on the JSON decimal format.
	// decimal = [ minus ] digits frac [ exp ]
	decimal = op.And{op.Optional(minus), digits, frac, op.Optional(exp)}
	// integer = [ minus ] digits
	integer = op.And{op.Optional(minus), digits}
	// dec = %x2E ; .
	dec = '.'
	// nonZero = %x31-39 ; 1-9
	nonZero = parser.CheckRuneRange('1', '9')
	// e = %x65 / %x45 ; e E
	e = op.Or{'e', 'E'}
	// exp = e [ minus / plus ] 1*digit
	exp = op.And{e, op.Optional(op.Or{minus, plus}), op.MinOne(digit)}
	// frac = dec 1*digit
	frac = op.And{dec, op.MinOne(digit)}
	// digits = zero / ( nonZero *digit )
	digits = op.Or{zero, op.And{nonZero, op.MinZero(digit)}}
	// minus = %x2D ; -
	minus = '-'
	// plus = %x2B ; +
	plus = '+'
	// zero = %x30 ; 0
	zero = '0'
)

// IsDecimal checks whether the given string is a decimal value. A decimal is real number with at least one digit to
// the left and right of the period.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.3
func IsDecimal(s string) bool {
	if err := is(decimal, s); err != nil {
		return false
	}
	return true
}

// IsInteger checks whether the given string is a integer value. An integer is a whole number with no fractional digits
// or exponent parts.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.4
func IsInteger(s string) bool {
	if err := is(integer, s); err != nil {
		return false
	}
	return true
}

var (
	// More info: https://www.w3.org/TR/xmlschema11-2/#dateTime
	// XML Schema Definition Language (XSD) 1.1
	datetime = op.And{
		year, '-', month, '-', day,
		op.Or{'T', 't'},
		hour, ':', minute, ':', second,
		offset,
	}
	year  = op.And{op.Optional('-'), op.Repeat(4, digit)}
	month = op.And{parser.CheckRuneRange('0', '1'), digit}
	day   = op.Or{
		op.And{parser.CheckRuneRange('0', '2'), digit},
		op.And{'3', parser.CheckRuneRange('0', '1')},
	}
	hour = op.Or{
		op.And{parser.CheckRuneRange('0', '1'), digit},
		op.And{'3', parser.CheckRuneRange('0', '3')},
	}
	minute = op.And{parser.CheckRuneRange('0', '5'), digit}
	second = op.And{
		parser.CheckRuneRange('0', '5'), digit,
		op.Optional(op.And{'.', op.MinOne(digit)}),
	}
	offset = op.Or{
		op.And{op.Optional('-'), hour, ':', minute},
		op.Or{'z', 'Z'},
	}
)

// IsDateTime checks whether the given string is a dateTime value. A date time format has no case sensitivity.
// e.g., 2008-01-23T04:56:22Z.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.5
func IsDateTime(s string) bool {
	if err := is(datetime, s); err != nil {
		return false
	}
	return true
}

var (
	// More info: https://tools.ietf.org/html/rfc4648#section-4
	base64 = op.And{
		op.MinZero(op.Repeat(4, alphabet)),
		op.Optional(op.Or{
			op.And{op.Repeat(3, alphabet), '='},
			op.And{op.Repeat(2, alphabet), op.Repeat(2, '=')},
			op.And{alphabet, op.Repeat(2, '=')},
		}),
		op.MinZero(op.And{op.Repeat(4, '=')}),
	}
	alphabet = op.Or{alpha, digit, '+', '/'}
	// More info: https://tools.ietf.org/html/rfc4648#section-5
	base64safe = op.And{
		op.MinZero(op.Repeat(4, alphabetSafe)),
		op.Optional(op.Or{
			op.And{op.Repeat(3, alphabetSafe), '='},
			op.And{op.Repeat(2, alphabetSafe), op.Repeat(2, '=')},
			op.And{alphabetSafe, op.Repeat(2, '=')},
		}),
		op.MinZero(op.And{op.Repeat(4, '=')}),
	}
	alphabetSafe = op.Or{alpha, digit, '-', '_'}
)

// IsBase64 checks whether the given string is arbitrary binary data. The attribute value MUST be base64, also supports
// URL-safe encoding.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.6
func IsBase64(s string) bool {
	if len(s) == 0 {
		return true
	}
	if err := is(base64, s); err == nil {
		return true
	}
	if err := is(base64safe, s); err == nil {
		return true
	}
	return false
}

// More info: https://tools.ietf.org/html/rfc7159#section-4
// We need this extra indirection, otherwise there is a loop.
//	object refers to
//	member refers to
//	object
// 	etc...
func object(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{
		ws, '{', ws,
		op.Optional(op.And{
			member,
			op.MinZero(op.And{ws, ',', ws, member}),
		}),
		ws, '}', ws,
	})
}

func member(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{str, ws, ':', ws, op.Or{
		f, n, t,
		object,
		array,
		decimal, integer, str,
	}})
}

// More info: https://tools.ietf.org/html/rfc7159#section-5
func array(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{
		ws, '[', ws,
		op.Optional(op.And{
			value,
			op.MinZero(op.And{ws, ',', ws, value}),
		}),
		ws, ']', ws,
	})
}

// About multi-valued attributes: https://tools.ietf.org/html/rfc7643#section-2.4
// - primitive values
// - objects
// - NO arrays

// More info: https://tools.ietf.org/html/rfc7159#section-3
// value is a CUSTOM definition based on the JSON array format.
func value(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{f, n, t, object, decimal, integer, str})
}

var ws = op.MinZero(op.Or{0x20, 0x09, 0x0A, 0x0D})

func IsObject(s string) bool {
	if err := is(object, s); err != nil {
		return false
	}
	return true
}
