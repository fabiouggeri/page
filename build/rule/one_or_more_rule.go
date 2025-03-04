package rule

import "github.com/fabiouggeri/page/util"

type OneOrMoreRule struct {
	rule Rule
}

var _ SimpleRule = &OneOrMoreRule{}

func (r *OneOrMoreRule) Rule() Rule {
	return r.rule
}

func (r *OneOrMoreRule) SetRule(rule Rule) {
	r.rule = rule
}

func (r *OneOrMoreRule) ToText(writer util.TextWriter) {
	r.rule.ToText(writer)
	writer.WriteRune('+')
}

func (r *OneOrMoreRule) Visit(visitor LexerVisitor) {
	visitor.VisitOneOrMoreRule(r)
}

func (r *OneOrMoreRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
