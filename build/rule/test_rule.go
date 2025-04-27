package rule

import "github.com/fabiouggeri/page/util"

type TestRule struct {
	rule Rule
}

var _ SimpleRule = &TestRule{}

func (r *TestRule) Rule() Rule {
	return r.rule
}

func (r *TestRule) SetRule(rule Rule) {
	r.rule = rule
}

func (r *TestRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}

func (r *TestRule) ToText(writer util.TextWriter) {
	r.rule.ToText(writer)
	writer.WriteRune('&')
}

func (r *TestRule) Visit(visitor RuleVisitor) {
	visitor.VisitTestRule(r)
}
