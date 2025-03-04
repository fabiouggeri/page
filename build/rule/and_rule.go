package rule

import "github.com/fabiouggeri/page/util"

type AndRule struct {
	rules []Rule
}

var _ CompoundRule = &AndRule{}

func (r *AndRule) Rules() []Rule {
	return r.rules
}

func (r *AndRule) SetRule(index int, rule Rule) {
	if index >= 0 && index < len(r.rules) {
		r.rules[index] = rule
	}
}

func (r *AndRule) ToText(writer util.TextWriter) {
	first := true
	writer.WriteRune('(')
	for _, s := range r.rules {
		if first {
			first = false
		} else {
			writer.WriteRune(' ')
		}
		s.ToText(writer)
	}
	writer.WriteRune(')')
}

func (r *AndRule) Visit(visitor LexerVisitor) {
	visitor.VisitAndRule(r)
}

func (r *AndRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
