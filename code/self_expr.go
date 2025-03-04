package code

import "github.com/fabiouggeri/page/util"

type Self struct{}

var _ Expression = &Self{}

func newSelf() *Self {
	return &Self{}
}

func (s *Self) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateSelf(s, str)
}

func (s *Self) Call(name string, args ...Expression) Expression {
	return newMethodCall(s, name, args...)
}

func (s *Self) Field(name string) LeftValue {
	return newFieldReference(s, name)
}

func (s *Self) String() string {
	return "self"
}

func (s *Self) Append(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Add(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) And(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Diff(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Div(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Equals(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Greater(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) GreaterEqual(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) IsEmpty() bool {
	panic("invalid operation in self")
}

func (s *Self) Less(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) LessEqual(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Minus(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Mult(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Index(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Not() Expression {
	panic("invalid operation in self")
}

func (s *Self) Or(other Expression) Expression {
	panic("invalid operation in self")
}

func (s *Self) Dec() Expression {
	panic("invalid operation in self")
}

func (s *Self) Inc() Expression {
	panic("invalid operation in self")
}
