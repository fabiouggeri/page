package code

import "github.com/fabiouggeri/page/util"

type Return struct {
	value Expression
}

type Break struct{}
type Continue struct{}

var _ Code = &Return{}
var _ Code = &Break{}
var _ Code = &Continue{}

func newReturn(expr Expression) *Return {
	return &Return{value: expr}
}

func (r *Return) IsEmpty() bool {
	return false
}

func (r *Return) Value(expr Expression) *Return {
	r.value = expr
	return r
}

func (r *Return) GetValue() Expression {
	return r.value
}

func (r *Return) String() string {
	if r.value != nil {
		return "return " + r.value.String()
	}
	return "return"
}

func (r *Return) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateReturn(r, str)
}

func newBreak() *Break {
	return &Break{}
}

func (b *Break) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateBreak(b, str)
}

func (b *Break) IsEmpty() bool {
	return false
}

func (b *Break) String() string {
	return "break"
}

func newContinue() *Continue {
	return &Continue{}
}

func (c *Continue) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateContinue(c, str)
}

func (c *Continue) IsEmpty() bool {
	return false
}

func (c *Continue) String() string {
	return "continue"
}
