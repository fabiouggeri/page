package code

import "github.com/fabiouggeri/page/util"

type DoubleOperatorExpression struct {
	left     Expression
	operator Operator
	right    Expression
}

var _ Expression = &DoubleOperatorExpression{}

func newDoubleOperatorExpression(left Expression, op Operator, right Expression) *DoubleOperatorExpression {
	return &DoubleOperatorExpression{left: left, operator: op, right: right}
}

func (d *DoubleOperatorExpression) IsEmpty() bool {
	return false
}

func (d *DoubleOperatorExpression) String() string {
	return d.left.String() + d.operator.String() + d.right.String()
}

func (d *DoubleOperatorExpression) Left() Expression {
	return d.left
}

func (d *DoubleOperatorExpression) Operator() Operator {
	return d.operator
}

func (d *DoubleOperatorExpression) Right() Expression {
	return d.right
}

func (d *DoubleOperatorExpression) Append(other Expression) Expression {
	return newDoubleOperatorExpression(d, APPEND, other)
}

func (d *DoubleOperatorExpression) Add(other Expression) Expression {
	return newDoubleOperatorExpression(d, PLUS, other)
}

func (d *DoubleOperatorExpression) Minus(other Expression) Expression {
	return newDoubleOperatorExpression(d, MINUS, other)
}

func (d *DoubleOperatorExpression) Mult(other Expression) Expression {
	return newDoubleOperatorExpression(d, MULT, other)
}

func (d *DoubleOperatorExpression) Div(other Expression) Expression {
	return newDoubleOperatorExpression(d, DIV, other)
}

func (d *DoubleOperatorExpression) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(d, EQUALS, other)
}

func (d *DoubleOperatorExpression) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(d, DIFF, other)
}

func (d *DoubleOperatorExpression) And(other Expression) Expression {
	return newDoubleOperatorExpression(d, AND, other)
}

func (d *DoubleOperatorExpression) Or(other Expression) Expression {
	return newDoubleOperatorExpression(d, OR, other)
}

func (d *DoubleOperatorExpression) Not() Expression {
	return newSingleOperatorExpression(NOT, d)
}

func (d *DoubleOperatorExpression) Call(name string, args ...Expression) Expression {
	return newMethodCall(d, name, args...)
}

func (f *DoubleOperatorExpression) Greater(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER, other)
}

func (f *DoubleOperatorExpression) GreaterEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER_EQUAL, other)
}

func (f *DoubleOperatorExpression) Less(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS, other)
}

func (f *DoubleOperatorExpression) LessEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS_EQUAL, other)
}

func (f *DoubleOperatorExpression) Index(other Expression) Expression {
	return newDoubleOperatorExpression(f, INDEX, other)
}

func (d *DoubleOperatorExpression) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateDoubleOpExpr(d, str)
}

func (d *DoubleOperatorExpression) Field(name string) LeftValue {
	return newFieldReference(d, name)
}

func (d *DoubleOperatorExpression) Dec() Expression {
	return newSingleOperatorExpression(DEC, d)
}

func (d *DoubleOperatorExpression) Inc() Expression {
	return newSingleOperatorExpression(INC, d)
}
