package code

import "github.com/fabiouggeri/page/util"

type Null struct{}

var _ Expression = &Null{}

func newNull() *Null {
	return &Null{}
}

func (n *Null) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateNull(n, str)
}

func (n *Null) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(n, EQUALS, other)
}

func (n *Null) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(n, DIFF, other)
}

func (n *Null) IsEmpty() bool {
	return false
}

func (n *Null) String() string {
	return "null"
}

func (n *Null) Append(other Expression) Expression {
	panic("invalid operation in null")
}

func (n *Null) Add(other Expression) Expression {
	panic("invalid operation in null")
}

func (n *Null) And(other Expression) Expression {
	panic("invalid operation in null")
}

func (n *Null) Div(other Expression) Expression {
	panic("invalid operation in null")
}

func (n *Null) Minus(other Expression) Expression {
	panic("invalid operation in null")
}

func (n *Null) Mult(other Expression) Expression {
	panic("invalid operation in null")
}

func (n *Null) Not() Expression {
	panic("invalid operation in null")
}

func (n *Null) Or(other Expression) Expression {
	panic("invalid operation in null")
}

func (d *Null) Call(name string, args ...Expression) Expression {
	panic("invalid operation in null")
}

func (d *Null) Greater(other Expression) Expression {
	panic("invalid operation in null")
}

func (d *Null) GreaterEqual(other Expression) Expression {
	panic("invalid operation in null")
}

func (d *Null) Less(other Expression) Expression {
	panic("invalid operation in null")
}

func (d *Null) LessEqual(other Expression) Expression {
	panic("invalid operation in null")
}

func (d *Null) Index(other Expression) Expression {
	panic("invalid operation in null")
}

func (d *Null) Field(name string) LeftValue {
	panic("invalid operation in null")
}

func (d *Null) Dec() Expression {
	panic("invalid operation in null")
}

func (d *Null) Inc() Expression {
	panic("invalid operation in null")
}
