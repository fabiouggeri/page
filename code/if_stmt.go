package code

import "github.com/fabiouggeri/page/util"

type If struct {
	condition Expression
	thenBody  *Block
	elseIf    []*ElseIf
	elseBody  *Block
}

type ElseIf struct {
	condition Expression
	body      *Block
}

var _ Code = &If{}

func newIf(condition Expression) *If {
	return &If{condition: condition, thenBody: newBlock(), elseIf: make([]*ElseIf, 0, 2), elseBody: newBlock()}
}

func (d *If) IsEmpty() bool {
	return false
}

func (i *If) String() string {
	return "if " + i.condition.String()
}

func (i *If) Condition() Expression {
	return i.condition
}

func (i *If) ThenBody() *Block {
	return i.thenBody
}

func (i *If) ElseBody() *Block {
	return i.elseBody
}

func (i *If) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateIf(i, str)
}

func (i *If) Then(code ...Code) *If {
	i.thenBody.Statements(code...)
	return i
}

func (i *If) Else(code ...Code) *If {
	i.elseBody.Statements(code...)
	return i
}

func (i *If) ElseIf(condition Expression, code ...Code) *If {
	elseIf := &ElseIf{condition: condition, body: newBlock(code...)}
	i.elseIf = append(i.elseIf, elseIf)
	return i
}

func (i *If) ElseIfs() []*ElseIf {
	return i.elseIf
}

func (ei *ElseIf) Condition() Expression {
	return ei.condition
}

func (ei *ElseIf) Body() *Block {
	return ei.body
}
