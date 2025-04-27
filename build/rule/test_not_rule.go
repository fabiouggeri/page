package rule

import "github.com/fabiouggeri/page/util"

type NotRule struct {
	rule Rule
}

var _ SimpleRule = &NotRule{}

func (r *NotRule) Rule() Rule {
	return r.rule
}

func (r *NotRule) SetRule(rule Rule) {
	r.rule = rule
}

func (r *NotRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}

func (r *NotRule) ToText(writer util.TextWriter) {
	r.rule.ToText(writer)
	writer.WriteRune('!')
}

func (r *NotRule) Visit(visitor RuleVisitor) {
	visitor.VisitNotRule(r)
}
