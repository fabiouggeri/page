package code

import "github.com/fabiouggeri/page/util"

type VariableDeclaration struct {
	comment    string
	datatype   DataType
	name       string
	qualifiers []Qualifier
	value      Expression
}

var _ Code = &VariableDeclaration{}

func NewVarDeclaration(name string, datatype DataType) *VariableDeclaration {
	return &VariableDeclaration{name: name, datatype: datatype, qualifiers: make([]Qualifier, 0)}
}

func (v *VariableDeclaration) Comment(comment string) {
	v.comment = comment
}

func (v *VariableDeclaration) GetComment() string {
	return v.comment
}

func (v *VariableDeclaration) IsEmpty() bool {
	return false
}

func (v *VariableDeclaration) String() string {
	return "var " + v.name + " " + v.datatype.Name()
}

func (v *VariableDeclaration) Name() string {
	return v.name
}

func (v *VariableDeclaration) DataType() DataType {
	return v.datatype
}

func (v *VariableDeclaration) GetValue() Expression {
	return v.value
}

func (v *VariableDeclaration) Value(value Expression) *VariableDeclaration {
	v.value = value
	return v
}

func (v *VariableDeclaration) Var() *Variable {
	return NewVar(v.name)
}

func (v *VariableDeclaration) Qualifiers(qualifiers ...Qualifier) *VariableDeclaration {
	v.qualifiers = append(v.qualifiers, qualifiers...)
	return v
}

func (v *VariableDeclaration) GetQualifiers() []Qualifier {
	return v.qualifiers
}

func (v *VariableDeclaration) HasQualifier(qualifier Qualifier) bool {
	for _, varQualifier := range v.qualifiers {
		if varQualifier == qualifier {
			return true
		}
	}
	return false
}

func (v *VariableDeclaration) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateVarDeclaration(v, str)
}
