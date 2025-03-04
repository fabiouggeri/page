package code

import "github.com/fabiouggeri/page/util"

type While struct {
	condition Expression
	body      *Block
}

var _ Code = &While{}

func newWhile(condition Expression) *While {
	return &While{condition: condition, body: newBlock()}
}

func (w *While) IsEmpty() bool {
	return false
}

func (w *While) String() string {
	return "while " + w.condition.String()
}

func (w *While) Condition() Expression {
	return w.condition
}

func (w *While) Body() *Block {
	return w.body
}

func (w *While) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateWhile(w, str)
}

func (w *While) Do(code ...Code) *While {
	w.body.Statements(code...)
	return w
}

func (w *While) If(condition Expression) *If {
	i := newIf(condition)
	w.body.Statements(i)
	return i
}

func (w *While) While(condition Expression) *While {
	while := newWhile(condition)
	w.body.Statements(while)
	return while
}

func (w *While) Var(name string) *Variable {
	v := NewVar(name)
	w.body.Statements(v)
	return v
}

func (w *While) Switch(test Expression) *Switch {
	s := newSwitch(test)
	w.body.Statements(s)
	return s
}

func (w *While) Assign(variable *Variable, value Expression) *While {
	w.body.Statements(variable.Assign(value))
	return w
}
