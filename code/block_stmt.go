package code

import "github.com/fabiouggeri/page/util"

type Block struct {
	statements []Code
}

var _ Code = &Block{}

func newBlock(statements ...Code) *Block {
	b := &Block{statements: make([]Code, 0)}
	b.statements = append(b.statements, statements...)
	return b
}

func (b *Block) String() string {
	return "{}"
}

func (b *Block) Statements(statements ...Code) *Block {
	b.statements = append(b.statements, statements...)
	return b
}

func (b *Block) GetStatements() []Code {
	return b.statements
}

func (b *Block) IsEmpty() bool {
	return len(b.statements) == 0
}

func (b *Block) Clear() {
	b.statements = make([]Code, 0)
}

func (b *Block) If(condition Expression) *If {
	i := newIf(condition)
	b.statements = append(b.statements, i)
	return i
}

func (b *Block) While(condition Expression) *While {
	w := newWhile(condition)
	b.statements = append(b.statements, w)
	return w
}

func (b *Block) Declare(datatype DataType, name string) *VariableDeclaration {
	v := NewVarDeclaration(name, datatype)
	b.statements = append(b.statements, v)
	return v
}

func (b *Block) Var(name string) *Variable {
	v := NewVar(name)
	b.statements = append(b.statements, v)
	return v
}

func (b *Block) Break() *Block {
	b.Statements(newBreak())
	return b
}

func (b *Block) Continue() *Block {
	b.Statements(newContinue())
	return b
}

func (b *Block) Call(name string, args ...Expression) *Block {
	b.Statements(newFunctionCall(name, args...))
	return b
}

func (b *Block) MethodCall(obj Expression, name string, args ...Expression) *Block {
	b.Statements(newMethodCall(obj, name, args...))
	return b
}

func (b *Block) Switch(test Expression) *Block {
	b.Statements(newSwitch(test))
	return b
}

func (b *Block) Assign(variable LeftValue, value Expression) *Block {
	b.Statements(variable.Assign(value))
	return b
}

func (b *Block) Comment(comment string) *Block {
	b.Statements(newComment(comment))
	return b
}

func (b *Block) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateBlock(b, str)
}
