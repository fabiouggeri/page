package code

import "github.com/fabiouggeri/page/util"

type Array struct {
	dataType DataType
}

var _ DataType = &Array{}
var _ Code = &Array{}

func NewArray(datatype DataType) *Array {
	return &Array{dataType: datatype}
}

func (a *Array) IsEmpty() bool {
	return false
}

func (a *Array) IsArray() bool {
	return true
}

func (a *Array) String() string {
	return "[]" + a.dataType.Name()
}

func (a *Array) Name() string {
	return "[]" + a.dataType.Name()
}

func (a *Array) IsPrimitive() bool {
	return true
}

func (a *Array) Declare(varName string, qualifiers ...Qualifier) *VariableDeclaration {
	return NewVarDeclaration(varName, a).Qualifiers(qualifiers...)
}

func (a *Array) Type() DataType {
	return a.dataType
}

func (a *Array) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateArray(a, str)
}
