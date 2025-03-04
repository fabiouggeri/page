package code

import "github.com/fabiouggeri/page/util"

type DataType interface {
	Code
	Name() string
	IsPrimitive() bool
	IsArray() bool
	Declare(varName string, qualifiers ...Qualifier) *VariableDeclaration
}

type Type struct {
	name string
}

type PrimitiveDataType int

const (
	Int8 PrimitiveDataType = iota
	Int16
	Int32
	Int64
	Float32
	Float64
	Char
	String
	Boolean
)

var _ DataType = Int8
var _ DataType = &Type{}

func (p PrimitiveDataType) Name() string {
	switch p {
	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case String:
		return "string"
	case Boolean:
		return "boolean"
	default:
		return "unknown"
	}
}

func (p PrimitiveDataType) IsPrimitive() bool {
	return true
}

func (t PrimitiveDataType) IsArray() bool {
	return false
}

func (p PrimitiveDataType) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateDataType(p, str)
}

func (p PrimitiveDataType) IsEmpty() bool {
	return false
}

func (p PrimitiveDataType) Declare(varName string, qualifiers ...Qualifier) *VariableDeclaration {
	return NewVarDeclaration(varName, p).Qualifiers(qualifiers...)
}

func (p PrimitiveDataType) String() string {
	return p.Name()
}

func NewType(name string) *Type {
	return &Type{name: name}
}

func (t *Type) Name() string {
	return t.name
}

func (t *Type) IsEmpty() bool {
	return false
}

func (t *Type) IsPrimitive() bool {
	return false
}

func (t *Type) IsArray() bool {
	return false
}

func (t *Type) Declare(varName string, qualifiers ...Qualifier) *VariableDeclaration {
	return NewVarDeclaration(varName, t).Qualifiers(qualifiers...)
}

func (t *Type) Pointer() *Pointer {
	return newPointer(t)
}

func (t *Type) String() string {
	return "type " + t.name
}

func (t *Type) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateDataType(t, str)
}
