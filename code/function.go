package code

import (
	"strings"

	"github.com/fabiouggeri/page/util"
)

type Function struct {
	name       string
	comment    string
	qualifiers []Qualifier
	params     []*VariableDeclaration
	returnType DataType
	body       *Block
}

var _ Code = &Function{}

func newFunction(name string, params ...*VariableDeclaration) *Function {
	f := &Function{name: name, qualifiers: make([]Qualifier, 0), params: make([]*VariableDeclaration, 0, 4), body: newBlock()}
	f.params = append(f.params, params...)
	return f
}

func (f *Function) Comment(comment string) {
	f.comment = comment
}

func (f *Function) GetComment() string {
	return f.comment
}

func (f *Function) IsEmpty() bool {
	return false
}

func (f *Function) String() string {
	str := strings.Builder{}
	str.WriteString(f.name)
	str.WriteRune('(')
	for _, p := range f.params {
		str.WriteString(p.String())
	}
	str.WriteRune(')')
	if f.returnType != nil {
		str.WriteRune(' ')
		str.WriteString(f.returnType.Name())
	}
	return str.String()
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Qualifiers(qualifiers ...Qualifier) *Function {
	f.qualifiers = append(f.qualifiers, qualifiers...)
	return f
}

func (f *Function) GetQualifiers() []Qualifier {
	return f.qualifiers
}

func (f *Function) Params(params ...*VariableDeclaration) *Function {
	f.params = append(f.params, params...)
	return f
}

func (f *Function) GetParams() []*VariableDeclaration {
	return f.params
}

func (f *Function) ReturnType(returnDatatype DataType) *Function {
	f.returnType = returnDatatype
	return f
}

func (f *Function) Return(expr Expression) *Function {
	f.body.Statements(newReturn(expr))
	return f
}

func (f *Function) GetReturnType() DataType {
	return f.returnType
}

func (f *Function) Body(statements ...Code) *Function {
	f.body.Statements(statements...)
	return f
}

func (f *Function) GetBody() *Block {
	return f.body
}

func (f *Function) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateFunction(f, str)
}

func (f *Function) If(condition Expression) *If {
	return f.body.If(condition)
}

func (f *Function) While(condition Expression) *While {
	return f.body.While(condition)
}

func (f *Function) Declare(datatype DataType, name string) *VariableDeclaration {
	return f.body.Declare(datatype, name)
}

func (f *Function) Var(name string) *Variable {
	return f.body.Var(name)
}

func (f *Function) Call(name string, args ...Expression) *Function {
	f.body.Call(name, args...)
	return f
}

func (f *Function) MethodCall(obj Expression, name string, args ...Expression) *Function {
	f.body.MethodCall(obj, name, args...)
	return f
}

func (f *Function) Assign(variable *Variable, value Expression) *Function {
	f.body.Statements(variable.Assign(value))
	return f
}

func (f *Function) Switch(test Expression) *Switch {
	s := newSwitch(test)
	f.body.Statements(s)
	return s
}
