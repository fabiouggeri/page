package code

import "github.com/fabiouggeri/page/util"

type SingleOperatorExpression struct {
	operator   Operator
	expression Expression
}

var _ Expression = &SingleOperatorExpression{}

func newSingleOperatorExpression(op Operator, expr Expression) *SingleOperatorExpression {
	return &SingleOperatorExpression{operator: op, expression: expr}
}

func (d *SingleOperatorExpression) IsEmpty() bool {
	return false
}

func (d *SingleOperatorExpression) String() string {
	return d.operator.String() + d.expression.String()
}

func (d *SingleOperatorExpression) Expression() Expression {
	return d.expression
}

func (d *SingleOperatorExpression) Operator() Operator {
	return d.operator
}

func (s *SingleOperatorExpression) Append(other Expression) Expression {
	return newDoubleOperatorExpression(s, APPEND, other)
}

func (s *SingleOperatorExpression) Add(other Expression) Expression {
	return newDoubleOperatorExpression(s, PLUS, other)
}

func (s *SingleOperatorExpression) Minus(other Expression) Expression {
	return newDoubleOperatorExpression(s, MINUS, other)
}

func (s *SingleOperatorExpression) Mult(other Expression) Expression {
	return newDoubleOperatorExpression(s, MULT, other)
}

func (s *SingleOperatorExpression) Div(other Expression) Expression {
	return newDoubleOperatorExpression(s, DIV, other)
}

func (s *SingleOperatorExpression) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(s, EQUALS, other)
}

func (s *SingleOperatorExpression) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(s, DIFF, other)
}

func (s *SingleOperatorExpression) And(other Expression) Expression {
	return newDoubleOperatorExpression(s, AND, other)
}

func (s *SingleOperatorExpression) Or(other Expression) Expression {
	return newDoubleOperatorExpression(s, OR, other)
}

func (s *SingleOperatorExpression) Not() Expression {
	return newSingleOperatorExpression(NOT, s)
}

func (s *SingleOperatorExpression) Call(name string, args ...Expression) Expression {
	return newMethodCall(s, name, args...)
}

func (s *SingleOperatorExpression) Greater(other Expression) Expression {
	return newDoubleOperatorExpression(s, GREATER, other)
}

func (s *SingleOperatorExpression) GreaterEqual(other Expression) Expression {
	return newDoubleOperatorExpression(s, GREATER_EQUAL, other)
}

func (s *SingleOperatorExpression) Less(other Expression) Expression {
	return newDoubleOperatorExpression(s, LESS, other)
}

func (s *SingleOperatorExpression) LessEqual(other Expression) Expression {
	return newDoubleOperatorExpression(s, LESS_EQUAL, other)
}

func (s *SingleOperatorExpression) Index(other Expression) Expression {
	return newDoubleOperatorExpression(s, INDEX, other)
}

func (s *SingleOperatorExpression) Field(name string) LeftValue {
	return newFieldReference(s, name)
}

func (s *SingleOperatorExpression) Dec() Expression {
	return newSingleOperatorExpression(DEC, s)
}

func (s *SingleOperatorExpression) Inc() Expression {
	return newSingleOperatorExpression(INC, s)
}

func (s *SingleOperatorExpression) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateSingleOpExpr(s, str)
}
