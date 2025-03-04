package rule

import (
	"strings"

	"github.com/fabiouggeri/page/util"
)

// StringRule represents a rule that matches a specific string.
// It contains the text to match, its size, and whether the match is case-sensitive.
type StringRule struct {
	text          string
	size          int32
	caseSensitive bool
}

var _ TerminalRule = &StringRule{}

// CaseSensitive returns whether the string match is case-sensitive.
func (r *StringRule) CaseSensitive() bool {
	return r.caseSensitive
}

// Text returns the text associated with the StringRule.
func (r *StringRule) Text() string {
	return r.text
}

// Size returns the size of the text associated with the StringRule.
func (r *StringRule) Size() int32 {
	return r.size
}

// ToText writes the StringRule to the provided TextWriter.
// If the match is case-sensitive, the text is enclosed in double quotes.
// Otherwise, it is enclosed in single quotes.
// If the size of the text does not match the length of the text, the size is also written.
func (r *StringRule) ToText(writer util.TextWriter) {
	if r.caseSensitive {
		writer.WriteRune('\'').WriteString(escapeString(r.text)).WriteRune('\'')
	} else {
		writer.WriteRune('"').WriteString(escapeString(r.text)).WriteRune('"')
	}
	if r.size != int32(len(r.text)) {
		writer.WriteF(":%d", r.size)
	}
}

func escapeString(s string) string {
	str := strings.Builder{}
	for _, c := range s {
		switch c {
		case '\x03':
			str.WriteString("EOI")
		case '\n':
			str.WriteString("\\n")
		case '\r':
			str.WriteString("\\r")
		case '\t':
			str.WriteString("\\t")
		case '\f':
			str.WriteString("\\f")
		case '\b':
			str.WriteString("\\b")
		case '\\':
			str.WriteString("\\\\")
		default:
			str.WriteRune(c)
		}
	}
	return str.String()
}

// Visit accepts a LexerVisitor and calls its VisitStringRule method.
func (r *StringRule) Visit(visitor LexerVisitor) {
	visitor.VisitStringRule(r)
}

// String returns a string representation of the StringRule.
func (r *StringRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
