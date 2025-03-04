package code

import (
	"github.com/fabiouggeri/page/util"
)

type Allocation struct {
	datatype DataType
	values   []Expression
}

var _ Expression = &Allocation{}

func newAllocation(datatype DataType, values ...Expression) *Allocation {
	return &Allocation{datatype: datatype, values: values}
}

func (a *Allocation) GetDataType() DataType {
	return a.datatype
}

func (a *Allocation) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateAllocation(a, str)
}

func (a *Allocation) IsEmpty() bool {
	return false
}

func (a *Allocation) Values() []Expression {
	return a.values
}

func (a *Allocation) Address() *SingleOperatorExpression {
	return newSingleOperatorExpression(ADDRESS_OF, a)
}

func (a *Allocation) String() string {
	return "new"
}

func (a *Allocation) Append(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Add(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) And(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Diff(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Div(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Equals(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Minus(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Mult(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Not() Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Or(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Call(name string, args ...Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Greater(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) GreaterEqual(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Less(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) LessEqual(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Index(other Expression) Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Field(name string) LeftValue {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Dec() Expression {
	panic("invalid operation in type allocation")
}

func (a *Allocation) Inc() Expression {
	panic("invalid operation in type allocation")
}
