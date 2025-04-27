package rule

import "github.com/fabiouggeri/page/util"

type RangeRule struct {
	start rune
	end   rune
}

var _ TerminalRule = &RangeRule{}

func (r *RangeRule) Start() rune {
	return r.start
}

func (r *RangeRule) End() rune {
	return r.end
}

func (r *RangeRule) Text() string {
	return r.String()
}

func (r *RangeRule) Size() int32 {
	return 1
}

func (r *RangeRule) CaseSensitive() bool {
	return true
}

func (r *RangeRule) ToText(writer util.TextWriter) {
	writer.WriteRune(r.start).WriteString("-").WriteRune(r.end)
}

func (r *RangeRule) Visit(visitor RuleVisitor) {
	visitor.VisitRangeRule(r)
}

func (r *RangeRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}
