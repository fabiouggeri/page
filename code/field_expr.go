package code

import "github.com/fabiouggeri/page/util"

type FieldReference struct {
	object Expression
	name   string
}

var _ Expression = &FieldReference{}

func newFieldReference(object Expression, name string) *FieldReference {
	return &FieldReference{object: object, name: name}
}

func (f *FieldReference) Object() Expression {
	return f.object
}

func (f *FieldReference) Name() string {
	return f.name
}

func (f *FieldReference) Assign(other Expression) Expression {
	return newDoubleOperatorExpression(f, ASSIGN, other)
}

func (f *FieldReference) Append(other Expression) Expression {
	return newDoubleOperatorExpression(f, APPEND, other)
}

func (f *FieldReference) Add(other Expression) Expression {
	return newDoubleOperatorExpression(f, PLUS, other)
}

func (f *FieldReference) And(other Expression) Expression {
	return newDoubleOperatorExpression(f, APPEND, other)
}

func (f *FieldReference) Call(name string, args ...Expression) Expression {
	return newMethodCall(f, name, args...)
}

func (f *FieldReference) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(f, DIFF, other)
}

func (f *FieldReference) Div(other Expression) Expression {
	return newDoubleOperatorExpression(f, DIV, other)
}

func (f *FieldReference) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(f, EQUALS, other)
}

func (f *FieldReference) Field(name string) LeftValue {
	return newFieldReference(f, name)
}

func (f *FieldReference) Greater(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER, other)
}

func (f *FieldReference) GreaterEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER_EQUAL, other)
}

func (f *FieldReference) IsEmpty() bool {
	return false
}

func (f *FieldReference) Less(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS, other)
}

func (f *FieldReference) LessEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS_EQUAL, other)
}

func (f *FieldReference) Minus(other Expression) Expression {
	return newDoubleOperatorExpression(f, MINUS, other)
}

func (f *FieldReference) Mult(other Expression) Expression {
	return newDoubleOperatorExpression(f, MULT, other)
}

func (f *FieldReference) Not() Expression {
	return newSingleOperatorExpression(NOT, f)
}

func (f *FieldReference) Or(other Expression) Expression {
	return newDoubleOperatorExpression(f, OR, other)
}

func (f *FieldReference) Index(other Expression) Expression {
	return newDoubleOperatorExpression(f, INDEX, other)
}

func (f *FieldReference) Dec() Expression {
	return newSingleOperatorExpression(DEC, f)
}

func (f *FieldReference) Inc() Expression {
	return newSingleOperatorExpression(INC, f)
}

func (f *FieldReference) String() string {
	return f.object.String() + "." + f.name
}

func (f *FieldReference) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateFieldReference(f, str)
}
