package code

import "github.com/fabiouggeri/page/util"

type Interface struct {
	comment string
	name    string
	methods []*Function
}

func NewInterface(name string, methods ...*Function) *Interface {
	i := &Interface{name: name, methods: make([]*Function, 0, 4)}
	i.methods = append(i.methods, methods...)
	return i
}

func (s *Interface) Comment(comment string) {
	s.comment = comment
}

func (s *Interface) GetComment() string {
	return s.comment
}

func (s *Interface) IsEmpty() bool {
	return false
}

func (s *Interface) String() string {
	return "interface " + s.name
}

func (s *Interface) Name() string {
	return s.name
}

func (s *Interface) IsPrimitive() bool {
	return false
}

func (s *Interface) IsArray() bool {
	return false
}

func (s *Interface) Declare(varName string, qualifiers ...Qualifier) *VariableDeclaration {
	return NewVarDeclaration(varName, s).Qualifiers(qualifiers...)
}

func (s *Interface) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateInterface(s, str)
}

func (s *Interface) Method(name string, params ...*VariableDeclaration) *Function {
	f := newFunction(name, params...)
	s.methods = append(s.methods, f)
	return f
}

func (s *Interface) Methods(methods ...*Function) *Interface {
	s.methods = append(s.methods, methods...)
	return s
}

func (s *Interface) GetMethods() []*Function {
	return s.methods
}

func (s *Interface) Type() *Type {
	return NewType(s.name)
}

func (s *Interface) Pointer() *Pointer {
	return newPointer(s)
}
