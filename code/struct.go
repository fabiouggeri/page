package code

import "github.com/fabiouggeri/page/util"

type Struct struct {
	comment string
	name    string
	fields  []*VariableDeclaration
	methods []*Function
}

var _ DataType = &Struct{}
var _ Code = &Struct{}

func NewStruct(name string, fields ...*VariableDeclaration) *Struct {
	s := &Struct{name: name, fields: make([]*VariableDeclaration, 0, 8), methods: make([]*Function, 0, 4)}
	s.fields = append(s.fields, fields...)
	return s
}

func (s *Struct) Comment(comment string) {
	s.comment = comment
}

func (s *Struct) GetComment() string {
	return s.comment
}

func (s *Struct) IsEmpty() bool {
	return false
}

func (s *Struct) String() string {
	return "struct " + s.name
}

func (s *Struct) Name() string {
	return s.name
}

func (s *Struct) IsPrimitive() bool {
	return false
}

func (s *Struct) IsArray() bool {
	return false
}

func (s *Struct) Declare(varName string, qualifiers ...Qualifier) *VariableDeclaration {
	return NewVarDeclaration(varName, s).Qualifiers(qualifiers...)
}

func (s *Struct) GetFields() []*VariableDeclaration {
	return s.fields
}

func (s *Struct) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateStruct(s, str)
}

func (s *Struct) Field(name string, datatype DataType) *VariableDeclaration {
	v := NewVarDeclaration(name, datatype)
	s.fields = append(s.fields, v)
	return v
}

func (s *Struct) Method(name string, params ...*VariableDeclaration) *Function {
	f := newFunction(name, params...)
	s.methods = append(s.methods, f)
	return f
}

func (s *Struct) Methods(methods ...*Function) *Struct {
	s.methods = append(s.methods, methods...)
	return s
}

func (s *Struct) GetMethods() []*Function {
	return s.methods
}

func (s *Struct) Fields(fields ...*VariableDeclaration) *Struct {
	s.fields = append(s.fields, fields...)
	return s
}

func (s *Struct) Type() *Type {
	return NewType(s.name)
}

func (s *Struct) Pointer() *Pointer {
	return newPointer(s)
}
