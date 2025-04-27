package rule

import "github.com/fabiouggeri/page/util"

type OrRule struct {
	rules []Rule
}

var _ CompoundRule = &OrRule{}

func (r *OrRule) Rules() []Rule {
	return r.rules
}

func (r *OrRule) SetRule(index int, rule Rule) {
	if index >= 0 && index < len(r.rules) {
		r.rules[index] = rule
	}
}

func (r *OrRule) ToText(writer util.TextWriter) {
	first := true
	writer.WriteRune('(')
	for _, s := range r.rules {
		if first {
			first = false
		} else {
			writer.WriteString(" | ")
		}
		s.ToText(writer)
	}
	writer.WriteRune(')')
}

func (r *OrRule) Visit(visitor RuleVisitor) {
	visitor.VisitOrRule(r)
}

func (r *OrRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
