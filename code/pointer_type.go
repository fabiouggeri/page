package code

import "github.com/fabiouggeri/page/util"

type Pointer struct {
	dataType DataType
}

var _ DataType = &Pointer{}
var _ Code = &Pointer{}

func newPointer(datatype DataType) *Pointer {
	return &Pointer{dataType: datatype}
}

func (p *Pointer) IsEmpty() bool {
	return false
}

func (p *Pointer) String() string {
	return "*" + p.dataType.Name()
}

func (p *Pointer) Name() string {
	return "*" + p.dataType.Name()
}

func (p *Pointer) IsPrimitive() bool {
	return true
}

func (p *Pointer) IsArray() bool {
	return false
}

func (p *Pointer) Declare(varName string, qualifiers ...Qualifier) *VariableDeclaration {
	return NewVarDeclaration(varName, p).Qualifiers(qualifiers...)
}

func (p *Pointer) Type() DataType {
	return p.dataType
}

func (p *Pointer) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GeneratePointer(p, str)
}
