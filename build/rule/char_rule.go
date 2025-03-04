package rule

import "github.com/fabiouggeri/page/util"

type CharRule struct {
	char          rune
	caseSensitive bool
}

var EOI *CharRule = &CharRule{char: '\x03', caseSensitive: false}

var specialChars = map[rune]string{
	'\x03': "EOI",
	'\n':   "'\\n'",
	'\r':   "'\\r'",
	'\t':   "'\\t'",
	'\f':   "'\\f'",
	'\b':   "'\\b'",
	'\\':   "'\\\\'",
}

var _ TerminalRule = &CharRule{}

func (r *CharRule) Char() rune {
	return r.char
}

func (r *CharRule) CaseSensitive() bool {
	return r.caseSensitive
}

func (r *CharRule) Text() string {
	return r.String()
}

func (r *CharRule) Size() int32 {
	return 1
}

func (r *CharRule) ToText(writer util.TextWriter) {
	charStr, found := specialChars[r.char]
	if found {
		writer.WriteString(charStr)
	} else {
		if r.caseSensitive {
			writer.WriteRune('\'')
		} else {
			writer.WriteRune('"')
		}
		writer.WriteRune(r.char)
		if r.caseSensitive {
			writer.WriteRune('\'')
		} else {
			writer.WriteRune('"')
		}
	}
}

func (r *CharRule) Visit(visitor LexerVisitor) {
	visitor.VisitCharRule(r)
}

func (r *CharRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
