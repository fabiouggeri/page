package code

import "github.com/fabiouggeri/page/util"

type FunctionCall struct {
	functionName string
	args         []Expression
}

type MethodCall struct {
	object     Expression
	methodName string
	args       []Expression
}

var _ Expression = &FunctionCall{}
var _ Expression = &MethodCall{}

func newFunctionCall(name string, args ...Expression) *FunctionCall {
	return &FunctionCall{functionName: name, args: args}
}

func (f *FunctionCall) FunctionName() string {
	return f.functionName
}

func (f *FunctionCall) Args() []Expression {
	return f.args
}

func (f *FunctionCall) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateFunctionCall(f, str)
}

func (f *FunctionCall) IsEmpty() bool {
	return false
}

func (f *FunctionCall) Append(other Expression) Expression {
	panic("invalid operation in function call")
}

func (f *FunctionCall) Add(other Expression) Expression {
	return newDoubleOperatorExpression(f, PLUS, other)
}

func (f *FunctionCall) And(other Expression) Expression {
	return newDoubleOperatorExpression(f, AND, other)
}

func (f *FunctionCall) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(f, DIFF, other)
}

func (f *FunctionCall) Div(other Expression) Expression {
	return newDoubleOperatorExpression(f, DIV, other)
}

func (f *FunctionCall) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(f, EQUALS, other)
}

func (f *FunctionCall) Minus(other Expression) Expression {
	return newDoubleOperatorExpression(f, MINUS, other)
}

func (f *FunctionCall) Mult(other Expression) Expression {
	return newDoubleOperatorExpression(f, MULT, other)
}

func (f *FunctionCall) Not() Expression {
	return newSingleOperatorExpression(NOT, f)
}

func (f *FunctionCall) Or(other Expression) Expression {
	return newDoubleOperatorExpression(f, OR, other)
}

func (f *FunctionCall) Call(name string, args ...Expression) Expression {
	return newMethodCall(f, name, args...)
}

func (f *FunctionCall) Greater(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER, other)
}

func (f *FunctionCall) GreaterEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER_EQUAL, other)
}

func (f *FunctionCall) Less(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS, other)
}

func (f *FunctionCall) LessEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS_EQUAL, other)
}

func (f *FunctionCall) Index(other Expression) Expression {
	return newDoubleOperatorExpression(f, INDEX, other)
}

func (f *FunctionCall) Field(name string) LeftValue {
	return newFieldReference(f, name)
}

func (f *FunctionCall) String() string {
	return f.functionName + "()"
}

func (f *FunctionCall) Dec() Expression {
	panic("invalid operation in function call")
}

func (f *FunctionCall) Inc() Expression {
	panic("invalid operation in function call")
}

// //////////////////////////// Method ///////////////////////////////
func newMethodCall(obj Expression, name string, args ...Expression) *MethodCall {
	return &MethodCall{object: obj, methodName: name, args: args}
}

func (m *MethodCall) MethodName() string {
	return m.methodName
}

func (m *MethodCall) Object() Expression {
	return m.object
}

func (m *MethodCall) Args() []Expression {
	return m.args
}

func (m *MethodCall) Append(other Expression) Expression {
	panic("invalid operation in Method call")
}

func (m *MethodCall) Add(other Expression) Expression {
	return newDoubleOperatorExpression(m, PLUS, other)
}

func (m *MethodCall) And(other Expression) Expression {
	return newDoubleOperatorExpression(m, AND, other)
}

func (m *MethodCall) Call(name string, args ...Expression) Expression {
	return newMethodCall(m, name, args...)
}

func (m *MethodCall) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(m, DIFF, other)
}

func (m *MethodCall) Div(other Expression) Expression {
	return newDoubleOperatorExpression(m, DIV, other)
}

func (m *MethodCall) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(m, EQUALS, other)
}

func (m *MethodCall) Minus(other Expression) Expression {
	return newDoubleOperatorExpression(m, MINUS, other)
}

func (m *MethodCall) Mult(other Expression) Expression {
	return newDoubleOperatorExpression(m, MULT, other)
}

func (m *MethodCall) Not() Expression {
	return newSingleOperatorExpression(NOT, m)
}

func (m *MethodCall) Or(other Expression) Expression {
	return newDoubleOperatorExpression(m, OR, other)
}

func (f *MethodCall) Greater(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER, other)
}

func (f *MethodCall) GreaterEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, GREATER_EQUAL, other)
}

func (f *MethodCall) Less(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS, other)
}

func (f *MethodCall) LessEqual(other Expression) Expression {
	return newDoubleOperatorExpression(f, LESS_EQUAL, other)
}

func (f *MethodCall) Index(other Expression) Expression {
	return newDoubleOperatorExpression(f, INDEX, other)
}

func (m *MethodCall) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateMethodCall(m, str)
}

func (m *MethodCall) IsEmpty() bool {
	return false
}

func (m *MethodCall) String() string {
	return m.object.String() + "." + m.methodName + "()"
}

func (m *MethodCall) Field(name string) LeftValue {
	return newFieldReference(m, name)
}

func (m *MethodCall) Dec() Expression {
	panic("invalid operation in method call")
}

func (m *MethodCall) Inc() Expression {
	panic("invalid operation in method call")
}
