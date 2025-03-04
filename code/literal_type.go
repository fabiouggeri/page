package code

import (
	"fmt"

	"github.com/fabiouggeri/page/util"
)

type LiteralDataType interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string | ~bool
}

type Literal[T LiteralDataType] struct {
	datatype DataType
	value    T
}

type ArrayLiteral struct {
	dataType DataType
	items    []Expression
}

var _ Expression = &Literal[int8]{}
var _ Expression = &ArrayLiteral{}

var TRUE = &Literal[bool]{datatype: Boolean, value: true}
var FALSE = &Literal[bool]{datatype: Boolean, value: false}
var ZERO = &Literal[int32]{datatype: Int32, value: 0}
var ONE = &Literal[int32]{datatype: Int32, value: 1}

func Value[T LiteralDataType](value T) *Literal[T] {
	return &Literal[T]{datatype: datatypeValue(value), value: value}
}

func Int(value int32) *Literal[int32] {
	return Value(value)
}

func NewChar(value rune) *Literal[rune] {
	return &Literal[rune]{datatype: Char, value: value}
}

func datatypeValue(value any) DataType {
	switch value.(type) {
	case bool:
		return Boolean
	case int8:
		return Int8
	case int16:
		return Int16
	case int32:
		return Int32
	case int64:
		return Int64
	case float32:
		return Float32
	case float64:
		return Float64
	case string:
		return String
	// case int:
	// case complex64:
	// case complex128:
	// case uint:
	// case uint8:
	// case uint16:
	// case uint32:
	// case uint64:
	// case uintptr:
	default:
		panic("Invalid literal datatype " + fmt.Sprint(value))
	}
}

func (l *Literal[T]) IsEmpty() bool {
	return false
}

func (l *Literal[T]) Value() T {
	return l.value
}

func (l *Literal[T]) String() string {
	return fmt.Sprint(l.value)
}

func (l *Literal[T]) Add(other Expression) Expression {
	return newDoubleOperatorExpression(l, PLUS, other)
}

func (l *Literal[T]) Append(other Expression) Expression {
	return newDoubleOperatorExpression(l, APPEND, other)
}

func (l *Literal[T]) Minus(other Expression) Expression {
	return newDoubleOperatorExpression(l, MINUS, other)
}

func (l *Literal[T]) Mult(other Expression) Expression {
	return newDoubleOperatorExpression(l, MULT, other)
}

func (l *Literal[T]) Div(other Expression) Expression {
	return newDoubleOperatorExpression(l, DIV, other)
}

func (l *Literal[T]) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(l, EQUALS, other)
}

func (l *Literal[T]) Diff(other Expression) Expression {
	return newDoubleOperatorExpression(l, DIFF, other)
}

func (l *Literal[T]) And(other Expression) Expression {
	return newDoubleOperatorExpression(l, AND, other)
}

func (l *Literal[T]) Or(other Expression) Expression {
	return newDoubleOperatorExpression(l, OR, other)
}

func (l *Literal[T]) Not() Expression {
	return newSingleOperatorExpression(NOT, l)
}

func (l *Literal[T]) Call(name string, args ...Expression) Expression {
	return newMethodCall(l, name, args...)
}

func (l *Literal[T]) Greater(other Expression) Expression {
	return newDoubleOperatorExpression(l, GREATER, other)
}

func (l *Literal[T]) GreaterEqual(other Expression) Expression {
	return newDoubleOperatorExpression(l, GREATER_EQUAL, other)
}

func (l *Literal[T]) Less(other Expression) Expression {
	return newDoubleOperatorExpression(l, LESS, other)
}

func (l *Literal[T]) LessEqual(other Expression) Expression {
	return newDoubleOperatorExpression(l, LESS_EQUAL, other)
}

func (l *Literal[T]) Index(other Expression) Expression {
	return newDoubleOperatorExpression(l, INDEX, other)
}

func (l *Literal[T]) Field(name string) LeftValue {
	return newFieldReference(l, name)
}

func (l *Literal[T]) Dec() Expression {
	panic("invalid operation to literal value")
}

func (l *Literal[T]) Inc() Expression {
	panic("invalid operation to literal value")
}

func (l *Literal[T]) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateLiteral(l.datatype, l.value, str)
}

func NewArrayLiteral(dataType DataType, items ...Expression) *ArrayLiteral {
	return &ArrayLiteral{dataType: dataType, items: items}
}

func (a *ArrayLiteral) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateArrayLiteral(a, str)
}

func (a *ArrayLiteral) AppendAll(items ...Expression) Expression {
	a.items = append(a.items, items...)
	return a
}

func (a *ArrayLiteral) Append(other Expression) Expression {
	a.items = append(a.items, other)
	return a
}

func (a *ArrayLiteral) Get(index int) Expression {
	return a.items[index]
}

func (a *ArrayLiteral) Set(index int, value Expression) *ArrayLiteral {
	a.items[index] = value
	return a
}

func (a *ArrayLiteral) DataType() DataType {
	return a.dataType
}

func (a *ArrayLiteral) Items() []Expression {
	return a.items
}

func (a *ArrayLiteral) Length() int {
	return len(a.items)
}

func (a *ArrayLiteral) Index(other Expression) Expression {
	return newDoubleOperatorExpression(a, INDEX, other)
}

func (a *ArrayLiteral) Equals(other Expression) Expression {
	return newDoubleOperatorExpression(a, EQUALS, other)
}

func (a *ArrayLiteral) IsEmpty() bool {
	return false
}

func (a *ArrayLiteral) String() string {
	return fmt.Sprint(a.items)
}

func (a *ArrayLiteral) Add(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) And(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Call(name string, args ...Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Dec() Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Diff(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Div(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Field(name string) LeftValue {
	panic("unimplemented")
}

func (a *ArrayLiteral) Greater(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) GreaterEqual(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Inc() Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Less(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) LessEqual(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Minus(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Mult(other Expression) Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Not() Expression {
	panic("unimplemented")
}

func (a *ArrayLiteral) Or(other Expression) Expression {
	panic("unimplemented")
}
