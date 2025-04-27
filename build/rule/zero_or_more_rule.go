package rule

import "github.com/fabiouggeri/page/util"

type ZeroOrMoreRule struct {
	rule Rule
}

var _ SimpleRule = &ZeroOrMoreRule{}

func (r *ZeroOrMoreRule) Rule() Rule {
	return r.rule
}

func (r *ZeroOrMoreRule) SetRule(rule Rule) {
	r.rule = rule
}

func (r *ZeroOrMoreRule) ToText(writer util.TextWriter) {
	r.rule.ToText(writer)
	writer.WriteRune('*')
}

func (r *ZeroOrMoreRule) Visit(visitor RuleVisitor) {
	visitor.VisitZeroOrMoreRule(r)
}

func (r *ZeroOrMoreRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
