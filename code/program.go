package code

import "github.com/fabiouggeri/page/util"

type Program struct {
	packageName string
	pathName    string
	statements  []Code
}

var _ Code = &Program{}

func NewProgram(pathName string) *Program {
	return &Program{pathName: pathName, statements: make([]Code, 0, 16)}
}

func (p *Program) IsEmpty() bool {
	return len(p.statements) == 0
}

func (p *Program) String() string {
	return p.pathName
}

func (p *Program) Package(name string) *Program {
	p.packageName = name
	return p
}

func (p *Program) GetPackage() string {
	return p.packageName
}

func (p *Program) GetPathName() string {
	return p.pathName
}

func (p *Program) GetStatements() []Code {
	return p.statements
}

func (p *Program) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateProgram(p, str)
}

func (p *Program) Struct(name string, fields ...*VariableDeclaration) *Struct {
	s := NewStruct(name, fields...)
	p.statements = append(p.statements, s)
	return s
}

func (p *Program) Interface(name string, methods ...*Function) *Interface {
	i := NewInterface(name, methods...)
	p.statements = append(p.statements, i)
	return i
}

func (p *Program) Function(name string, params ...*VariableDeclaration) *Function {
	f := newFunction(name, params...)
	p.statements = append(p.statements, f)
	return f
}

func (p *Program) Declare(datatype DataType, name string) *VariableDeclaration {
	v := NewVarDeclaration(name, datatype)
	p.statements = append(p.statements, v)
	return v
}

func (p *Program) Statements(stmt ...Code) *Program {
	p.statements = append(p.statements, stmt...)
	return p
}
