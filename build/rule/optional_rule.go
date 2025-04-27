package rule

import "github.com/fabiouggeri/page/util"

type OptionalRule struct {
	rule Rule
}

var _ SimpleRule = &OptionalRule{}

func (r *OptionalRule) Rule() Rule {
	return r.rule
}

func (r *OptionalRule) SetRule(rule Rule) {
	r.rule = rule
}

func (r *OptionalRule) ToText(writer util.TextWriter) {
	r.rule.ToText(writer)
	writer.WriteRune('?')
}

func (r *OptionalRule) Visit(visitor RuleVisitor) {
	visitor.VisitOptionalRule(r)
}

func (r *OptionalRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
