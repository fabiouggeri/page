package code

import "github.com/fabiouggeri/page/util"

type Qualifier int

const (
	PRIVATE Qualifier = iota
	PUBLIC
	CONST
)

type Code interface {
	String() string
	Generate(generator CodeGenerator, str util.TextWriter) error
	IsEmpty() bool
}

type EmptyCode struct{}

type CodeGenerator interface {
	GenerateIf(ifStmt *If, str util.TextWriter) error
	GenerateWhile(whileStmt *While, str util.TextWriter) error
	GenerateFunction(fun *Function, str util.TextWriter) error
	GenerateProgram(prog *Program, str util.TextWriter) error
	GenerateStruct(stru *Struct, str util.TextWriter) error
	GenerateInterface(interf *Interface, str util.TextWriter) error
	GenerateBlock(body *Block, str util.TextWriter) error
	GenerateVar(variable *Variable, str util.TextWriter) error
	GenerateVarDeclaration(variable *VariableDeclaration, str util.TextWriter) error
	GenerateLiteral(datatype DataType, value any, str util.TextWriter) error
	GenerateArrayLiteral(array *ArrayLiteral, str util.TextWriter) error
	GenerateDataType(dataType DataType, str util.TextWriter) error
	GeneratePointer(ptr *Pointer, str util.TextWriter) error
	GenerateDoubleOpExpr(dobleOp *DoubleOperatorExpression, str util.TextWriter) error
	GenerateSingleOpExpr(singleOp *SingleOperatorExpression, str util.TextWriter) error
	GenerateReturn(ret *Return, str util.TextWriter) error
	GenerateAllocation(allocation *Allocation, str util.TextWriter) error
	GenerateNull(null *Null, str util.TextWriter) error
	GenerateBreak(brk *Break, str util.TextWriter) error
	GenerateContinue(cont *Continue, str util.TextWriter) error
	GenerateFunctionCall(call *FunctionCall, str util.TextWriter) error
	GenerateMethodCall(call *MethodCall, str util.TextWriter) error
	GenerateSwitch(swtch *Switch, str util.TextWriter) error
	GenerateSelf(self *Self, str util.TextWriter) error
	GenerateFieldReference(field *FieldReference, str util.TextWriter) error
	GenerateComment(comment *Comment, str util.TextWriter) error
	GenerateArray(array *Array, str util.TextWriter) error
}

type CodeBuilder struct {
}

var EMPTY_CODE EmptyCode = EmptyCode{}

func NewBuilder() *CodeBuilder {
	return &CodeBuilder{}
}

func (e EmptyCode) String() string {
	return ""
}

func (e EmptyCode) Generate(generator CodeGenerator, str util.TextWriter) error {
	return nil
}

func (e EmptyCode) IsEmpty() bool {
	return true
}

func (cb *CodeBuilder) Program(pathName string) *Program {
	return NewProgram(pathName)
}

func (cb *CodeBuilder) Struct(name string) *Struct {
	return NewStruct(name)
}

func (cb *CodeBuilder) If(condition Expression) *If {
	return newIf(condition)
}

func (cb *CodeBuilder) While(condition Expression) *While {
	return newWhile(condition)
}

func (cb *CodeBuilder) Function(name string) *Function {
	return newFunction(name)
}

func (cb *CodeBuilder) Pointer(datatype DataType) *Pointer {
	return newPointer(datatype)
}

func (cb *CodeBuilder) Declare(datatype DataType, name string) *VariableDeclaration {
	return NewVarDeclaration(name, datatype)
}

func (cb *CodeBuilder) Var(name string) *Variable {
	return NewVar(name)
}

func (cb *CodeBuilder) Not(expr Expression) Expression {
	return newSingleOperatorExpression(NOT, expr)
}

func (cb *CodeBuilder) True() *Literal[bool] {
	return TRUE
}

func (cb *CodeBuilder) False() *Literal[bool] {
	return FALSE
}

func (cb *CodeBuilder) Int8(value int8) *Literal[int8] {
	return Value(value)
}

func (cb *CodeBuilder) Int16(value int16) *Literal[int16] {
	return Value(value)
}

func (cb *CodeBuilder) Int32(value int32) *Literal[int32] {
	return Value(value)
}

func (cb *CodeBuilder) Int64(value int64) *Literal[int64] {
	return Value(value)
}

func (cb *CodeBuilder) Float32(value float32) *Literal[float32] {
	return Value(value)
}

func (cb *CodeBuilder) Float64(value float64) *Literal[float64] {
	return Value(value)
}

func (cb *CodeBuilder) Char(value rune) *Literal[rune] {
	return NewChar(value)
}

func (cb *CodeBuilder) Str(value string) *Literal[string] {
	return Value(value)
}

func (cb *CodeBuilder) Block(code ...Code) *Block {
	return newBlock(code...)
}

func (cb *CodeBuilder) Null() *Null {
	return newNull()
}

func (cb *CodeBuilder) Allocate(datatype DataType, values ...Expression) *Allocation {
	return newAllocation(datatype, values...)
}

func (cb *CodeBuilder) Break() *Break {
	return newBreak()
}

func (cb *CodeBuilder) Continue() *Continue {
	return newContinue()
}

func (cb *CodeBuilder) Switch(test Expression) *Switch {
	return newSwitch(test)
}

func (cb *CodeBuilder) Call(name string, args ...Expression) *FunctionCall {
	return newFunctionCall(name, args...)
}

func (cb *CodeBuilder) MethodCall(obj Expression, name string, args ...Expression) *MethodCall {
	return newMethodCall(obj, name, args...)
}

func (cb *CodeBuilder) Self() *Self {
	return newSelf()
}

func (cb *CodeBuilder) Type(name string) *Type {
	return NewType(name)
}

func (cb *CodeBuilder) ArrayOf(datatype DataType) *Array {
	return NewArray(datatype)
}

func (cb *CodeBuilder) Array(dataType DataType, items ...Expression) *ArrayLiteral {
	return NewArrayLiteral(dataType, items...)
}

func (cb *CodeBuilder) ArrayInit(dataType DataType, size int, value Expression) *ArrayLiteral {
	items := make([]Expression, size)
	for i := 0; i < size; i++ {
		items[i] = value
	}
	return NewArrayLiteral(dataType, items...)
}

func (cb *CodeBuilder) ArrayLen(dataType DataType, size int) *ArrayLiteral {
	return NewArrayLiteral(dataType, make([]Expression, size)...)
}

func (cb *CodeBuilder) Return(expr Expression) *Return {
	return newReturn(expr)
}

func (cb *CodeBuilder) AddressOf(value Expression) *SingleOperatorExpression {
	return newSingleOperatorExpression(ADDRESS_OF, value)
}
