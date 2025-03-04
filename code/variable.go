package code

import (
	"github.com/fabiouggeri/page/util"
)

type Variable struct {
	name string
}

var _ Expression = &Variable{}

func NewVar(name string) *Variable {
	return &Variable{name: name}
}

func (v *Variable) Address() *SingleOperatorExpression {
	return newSingleOperatorExpression(ADDRESS_OF, v)
}

func (v *Variable) IsEmpty() bool {
	return false
}

func (v *Variable) String() string {
	return v.name
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Declare(datatype DataType) *VariableDeclaration {
	return NewVarDeclaration(v.name, datatype)
}

func (v *Variable) Assign(expr Expression) Expression {
	return newDoubleOperatorExpression(v, ASSIGN, expr)
}

func (v *Variable) Append(other Expression) Expression {
	return newDoubleOperatorExpression(v, APPEND, other)
}

func (v *Variable) Add(other Expression) Expression {
	return newDoubleOperatorExpression(v, PLUS, other)
}

func (v *Variable) Minus(other Expression) Expression {
	return newDoubleOperatorExpression(v, MINUS, other)
}

func (v *Variable) Mult(other Expression) Expression {
	return newDoubleOperatorExpression(v, MULT, other)
}

func (v *Variable) Div(other Expression) Expression {
	return newDoubleOperatorExpression(v, DIV, other)
}

func (v *Variable) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(v, EQUALS, other)
}

func (v *Variable) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(v, DIFF, other)
}

func (v *Variable) And(other Expression) Expression {
	return newDoubleOperatorExpression(v, AND, other)
}

func (v *Variable) Or(other Expression) Expression {
	return newDoubleOperatorExpression(v, OR, other)
}

func (v *Variable) Not() Expression {
	return newSingleOperatorExpression(NOT, v)
}

func (v *Variable) Call(name string, args ...Expression) Expression {
	return newMethodCall(v, name, args...)
}

func (v *Variable) Greater(other Expression) Expression {
	return newDoubleOperatorExpression(v, GREATER, other)
}

func (v *Variable) GreaterEqual(other Expression) Expression {
	return newDoubleOperatorExpression(v, GREATER_EQUAL, other)
}

func (v *Variable) Less(other Expression) Expression {
	return newDoubleOperatorExpression(v, LESS, other)
}

func (v *Variable) LessEqual(other Expression) Expression {
	return newDoubleOperatorExpression(v, LESS_EQUAL, other)
}

func (v *Variable) Index(expr Expression) Expression {
	return newDoubleOperatorExpression(v, INDEX, expr)
}

func (v *Variable) Field(name string) LeftValue {
	return newFieldReference(v, name)
}

func (v *Variable) Dec() Expression {
	return newSingleOperatorExpression(DEC, v)
}

func (v *Variable) Inc() Expression {
	return newSingleOperatorExpression(INC, v)
}

func (v *Variable) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateVar(v, str)
}
