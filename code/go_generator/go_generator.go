package go_generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiouggeri/page/code"
	"github.com/fabiouggeri/page/util"
)

type GoGenerator struct {
	inVarDeclaration bool
	inFunctionBody   bool
	arrayLevel       int
}

var _ code.CodeGenerator = &GoGenerator{}

func New() *GoGenerator {
	return &GoGenerator{inVarDeclaration: true, inFunctionBody: false}
}

func (g *GoGenerator) GenerateProgram(program *code.Program, str util.TextWriter) error {
	str.WriteString("package ").WriteString(program.GetPackage()).NewLine().NewLine()
	for _, s := range program.GetStatements() {
		if err := s.Generate(g, str); err != nil {
			return err
		}
		str.NewLine()
	}
	return str.Error()
}

func (g *GoGenerator) GenerateBlock(body *code.Block, str util.TextWriter) error {
	previousInBody := g.inFunctionBody
	g.inFunctionBody = true
	defer func() { g.inFunctionBody = previousInBody }()
	str.WriteRune('{').NewLine().Indent(3)
	for _, c := range body.GetStatements() {
		if err := c.Generate(g, str); err != nil {
			return err
		}
		str.NewLine()
	}
	return str.Indent(-3).WriteRune('}').Error()
}

func (g *GoGenerator) GenerateIf(ifstmt *code.If, str util.TextWriter) error {
	str.WriteString("if ")
	if err := ifstmt.Condition().Generate(g, str); err != nil {
		return err
	}
	str.WriteRune(' ')
	if err := ifstmt.ThenBody().Generate(g, str); err != nil {
		return err
	}
	for _, ei := range ifstmt.ElseIfs() {
		str.WriteString(" else if ")
		ei.Condition().Generate(g, str)
		ei.Body().Generate(g, str)
	}
	if !ifstmt.ElseBody().IsEmpty() {
		str.WriteString(" else ")
		ifstmt.ElseBody().Generate(g, str)
	}
	return str.Error()
}

func (g *GoGenerator) GenerateWhile(whilestmt *code.While, str util.TextWriter) error {
	str.WriteString("for ")
	if err := whilestmt.Condition().Generate(g, str); err != nil {
		return err
	}
	str.WriteRune(' ')
	return whilestmt.Body().Generate(g, str)
}

func (g *GoGenerator) GenerateFunction(fun *code.Function, str util.TextWriter) error {
	return g.generateFunction("func ", fun, str)
}

func (g *GoGenerator) generateFunction(funType string, fun *code.Function, str util.TextWriter) error {
	writeComment(fun.GetComment(), str, true)
	str.WriteString(funType).WriteString(fun.Name()).WriteRune('(')
	first := true
	previousInVarDeclaration := g.inVarDeclaration
	g.inVarDeclaration = false
	defer func() { g.inVarDeclaration = previousInVarDeclaration }()
	for _, param := range fun.GetParams() {
		if first {
			first = false
		} else {
			str.WriteString(", ")
		}
		if err := param.Generate(g, str); err != nil {
			return err
		}
	}
	str.WriteString(") ")
	if fun.GetReturnType() != nil {
		g.GenerateDataType(fun.GetReturnType(), str)
		str.WriteRune(' ')
	}
	if !fun.GetBody().IsEmpty() {
		fun.GetBody().Generate(g, str)
		str.NewLine()
	}
	return str.Error()
}

func (g *GoGenerator) GenerateStruct(stru *code.Struct, str util.TextWriter) error {
	writeComment(stru.GetComment(), str, true)
	previousInVarDeclaration := g.inVarDeclaration
	g.inVarDeclaration = false
	defer func() { g.inVarDeclaration = previousInVarDeclaration }()
	if stru.Name() != "" {
		str.WriteF("type %s struct {", stru.Name()).NewLine().Indent(3)
	} else {
		str.WriteString("struct {").NewLine().Indent(3)
	}
	for _, c := range stru.GetFields() {
		if err := c.Generate(g, str); err != nil {
			return err
		}
		str.NewLine()
	}
	str.Indent(-3).WriteRune('}').NewLine().NewLine()
	selfPrefix := "func (self *" + stru.Name() + ") "
	for _, f := range stru.GetMethods() {
		if err := g.generateFunction(selfPrefix, f, str); err != nil {
			return err
		}
		str.NewLine()
	}
	return str.Error()
}

func (g *GoGenerator) GenerateInterface(interf *code.Interface, str util.TextWriter) error {
	writeComment(interf.GetComment(), str, true)
	previousInVarDeclaration := g.inVarDeclaration
	g.inVarDeclaration = false
	defer func() { g.inVarDeclaration = previousInVarDeclaration }()
	if interf.Name() != "" {
		str.WriteF("type %s interface {", interf.Name()).NewLine().Indent(3)
	} else {
		str.WriteString("interface {").NewLine().Indent(3)
	}
	for _, c := range interf.GetMethods() {
		if err := g.generateFunction("", c, str); err != nil {
			return err
		}
		str.NewLine()
	}
	str.Indent(-3).WriteRune('}').NewLine().NewLine()
	return str.Error()
}

func (g *GoGenerator) GenerateVar(variable *code.Variable, str util.TextWriter) error {
	return str.WriteString(variable.Name()).Error()
}

func (g *GoGenerator) GenerateVarDeclaration(variable *code.VariableDeclaration, str util.TextWriter) error {
	var err error
	writeComment(variable.GetComment(), str, true)
	if g.inFunctionBody {
		if variable.GetValue() != nil {
			str.WriteString(variable.Name()).WriteString(" := ")
			err = variable.GetValue().Generate(g, str)
		} else {
			if variable.HasQualifier(code.CONST) {
				str.WriteString("const ")
			} else {
				str.WriteString("var ")
			}
			str.WriteString(variable.Name()).WriteRune(' ')
			err = g.GenerateDataType(variable.DataType(), str)
		}
	} else if g.inVarDeclaration {
		if variable.HasQualifier(code.CONST) {
			str.WriteString("const ")
		} else {
			str.WriteString("var ")
		}
		str.WriteString(variable.Name())
		if variable.GetValue() == nil {
			str.WriteRune(' ')
			err = g.GenerateDataType(variable.DataType(), str)
		} else {
			str.WriteString(" = ")
			err = variable.GetValue().Generate(g, str)
		}
	} else {
		str.WriteString(variable.Name()).WriteRune(' ')
		err = g.GenerateDataType(variable.DataType(), str)
	}
	return err
}

func (g *GoGenerator) GenerateDataType(dataType code.DataType, str util.TextWriter) error {
	switch dataType {
	case code.Boolean:
		str.WriteString("bool")
	case code.Int8:
		str.WriteString("int8")
	case code.Int16:
		str.WriteString("int16")
	case code.Int32:
		str.WriteString("int32")
	case code.Int64:
		str.WriteString("int64")
	case code.Float32:
		str.WriteString("float32")
	case code.Float64:
		str.WriteString("float64")
	case code.Char:
		str.WriteString("rune")
	case code.String:
		str.WriteString("string")
	default:
		str.WriteString(dataType.Name())
	}
	return str.Error()
}

func (g *GoGenerator) GenerateDoubleOpExpr(doubleOp *code.DoubleOperatorExpression, str util.TextWriter) error {
	doubleOp.Left().Generate(g, str)
	switch doubleOp.Operator() {
	case code.APPEND:
		str.WriteString(" = append(")
		doubleOp.Left().Generate(g, str)
		str.WriteString(", ")
		doubleOp.Right().Generate(g, str)
		str.WriteRune(')')
	case code.INDEX:
		str.WriteRune('[')
		doubleOp.Right().Generate(g, str)
		str.WriteRune(']')
	default:
		str.WriteString(operator(doubleOp.Operator()))
		doubleOp.Right().Generate(g, str)
	}
	return nil
}

func operator(operator code.Operator) string {
	switch operator {
	case code.PLUS:
		return " + "
	case code.MINUS:
		return " - "
	case code.MULT:
		return " * "
	case code.DIV:
		return " / "
	case code.EQUALS:
		return " == "
	case code.ASSIGN:
		return " = "
	case code.DIFF:
		return " != "
	case code.AND:
		return " && "
	case code.OR:
		return " || "
	case code.NOT:
		return "!"
	case code.INC:
		return "++"
	case code.DEC:
		return "--"
	case code.GREATER:
		return " > "
	case code.GREATER_EQUAL:
		return " >= "
	case code.LESS:
		return " < "
	case code.LESS_EQUAL:
		return " <= "
	case code.ADDRESS_OF:
		return "&"
	default:
		return ""
	}
}

func (g *GoGenerator) GenerateSingleOpExpr(singleOp *code.SingleOperatorExpression, str util.TextWriter) error {
	switch singleOp.Operator() {
	case code.INC:
		singleOp.Expression().Generate(g, str)
		str.WriteString("++")
	case code.DEC:
		singleOp.Expression().Generate(g, str)
		str.WriteString("--")
	default:
		str.WriteString(operator(singleOp.Operator()))
		singleOp.Expression().Generate(g, str)
	}
	return str.Error()
}

func (g *GoGenerator) GenerateArrayLiteral(array *code.ArrayLiteral, str util.TextWriter) error {
	if g.arrayLevel == 0 {
		str.WriteString("[]")
		if err := g.GenerateDataType(array.DataType(), str); err != nil {
			return err
		}
	} else {
		str.NewLine()
		str.Indent(3)
	}
	g.arrayLevel++
	str.WriteRune('{')
	for index, item := range array.Items() {
		if index > 0 {
			str.WriteRune(',')
		}
		item.Generate(g, str)
	}
	str.WriteRune('}')
	if g.arrayLevel > 0 {
		str.Indent(-3)
	}
	g.arrayLevel--
	return nil
}

func (g *GoGenerator) GenerateLiteral(dataType code.DataType, value any, str util.TextWriter) error {
	var err error
	switch dataType {
	case code.Int8, code.Int16, code.Int32, code.Int64, code.Float32, code.Float64, code.Boolean:
		err = str.WriteString(fmt.Sprint(value)).Error()
	case code.Char:
		if v, ok := value.(rune); ok {
			err = str.WriteF("%s", strconv.QuoteRune(v)).Error()
		} else {
			err = fmt.Errorf("invalid value type for Char: %v", value)
		}
	case code.String:
		err = str.WriteF("\"%s\"", value).Error()
	default:
		if v, ok := value.(*code.Struct); ok {
			err = g.GenerateStruct(v, str)
		} else {
			err = fmt.Errorf("invalid value type for Char: %v", value)
		}
	}
	return err
}

func (g *GoGenerator) GeneratePointer(ptr *code.Pointer, str util.TextWriter) error {
	str.WriteRune('*')
	return ptr.Type().Generate(g, str)
}

func (g *GoGenerator) GenerateArray(array *code.Array, str util.TextWriter) error {
	str.WriteString("[]")
	return array.Type().Generate(g, str)
}

func (g *GoGenerator) GenerateReturn(ret *code.Return, str util.TextWriter) error {
	if ret.GetValue() != nil {
		str.WriteString("return ")
		return ret.GetValue().Generate(g, str)
	} else {
		return str.WriteString("return").Error()
	}
}

func (g *GoGenerator) GenerateAllocation(allocation *code.Allocation, str util.TextWriter) error {
	datatype := allocation.GetDataType()
	g.GenerateDataType(datatype, str)
	str.WriteRune('{')
	if stru, ok := datatype.(*code.Struct); ok {
		first := true
		values := allocation.Values()
		for i, f := range stru.GetFields() {
			var value code.Expression
			if i < len(values) {
				value = values[i]
			} else {
				value = f.GetValue()
			}
			if value != nil {
				if first {
					first = false
				} else {
					str.WriteString(", ")
				}
				str.WriteString(f.Name()).WriteString(": ")
				value.Generate(g, str)
			}
		}
	}
	str.WriteRune('}')
	return str.Error()
}

func (g *GoGenerator) GenerateNull(null *code.Null, str util.TextWriter) error {
	return str.WriteString("nil").Error()
}

func (g *GoGenerator) GenerateBreak(brk *code.Break, str util.TextWriter) error {
	return str.WriteString("break").Error()
}

func (g *GoGenerator) GenerateSelf(self *code.Self, str util.TextWriter) error {
	return str.WriteString("self").Error()
}

func (g *GoGenerator) GenerateContinue(cont *code.Continue, str util.TextWriter) error {
	return str.WriteString("continue").Error()
}

func (g *GoGenerator) GenerateFunctionCall(call *code.FunctionCall, str util.TextWriter) error {
	str.WriteString(call.FunctionName()).WriteRune('(')
	first := true
	for _, arg := range call.Args() {
		if first {
			first = false
		} else {
			str.WriteString(", ")
		}
		arg.Generate(g, str)
	}
	str.WriteRune(')')
	return str.Error()
}

func (g *GoGenerator) GenerateMethodCall(call *code.MethodCall, str util.TextWriter) error {
	call.Object().Generate(g, str)
	str.WriteRune('.').WriteString(call.MethodName()).WriteRune('(')
	first := true
	for _, arg := range call.Args() {
		if first {
			first = false
		} else {
			str.WriteString(", ")
		}
		arg.Generate(g, str)
	}
	str.WriteRune(')')
	return str.Error()
}

func (g *GoGenerator) GenerateSwitch(swtch *code.Switch, str util.TextWriter) error {
	str.WriteString("switch ")
	swtch.GetTest().Generate(g, str)
	str.WriteString(" {").NewLine()
	str.Indent(3)
	g.generateSwitchOptions(swtch, str)
	if len(swtch.GetDefault()) > 0 {
		str.WriteString("default:").NewLine()
		str.Indent(3)
		for _, s := range swtch.GetDefault() {
			s.Generate(g, str)
			str.NewLine()
		}
		str.Indent(-3)
	}
	str.Indent(-3).WriteRune('}').NewLine()
	return str.Error()
}

func (g *GoGenerator) generateSwitchOptions(swtch *code.Switch, str util.TextWriter) {
	for _, o := range swtch.GetOptions() {
		str.WriteString("case ")
		o.GetOption().Generate(g, str)
		str.WriteString(": ").NewLine()
		str.Indent(3)
		for _, s := range o.GetStatements() {
			s.Generate(g, str)
			str.NewLine()
		}
		str.Indent(-3)
	}
}

func (g *GoGenerator) GenerateFieldReference(field *code.FieldReference, str util.TextWriter) error {
	field.Object().Generate(g, str)
	return str.WriteRune('.').WriteString(field.Name()).Error()
}

func (g *GoGenerator) GenerateComment(comment *code.Comment, str util.TextWriter) error {
	writeComment(comment.Comment(), str, false)
	return str.Error()
}

func writeComment(comment string, str util.TextWriter, newLine bool) {
	if comment == "" {
		return
	}
	lines := strings.Split(comment, "\n")
	if len(lines) > 1 {
		str.WriteString("/* ").WriteString(lines[0])
		for i := 1; i < len(lines); i++ {
			str.NewLine()
			str.WriteString(lines[i])
		}
		str.WriteString(" */")
		if newLine {
			str.NewLine()
		}
	} else if len(lines) == 1 {
		str.WriteString("// ").WriteString(lines[0])
		if newLine {
			str.NewLine()
		}
	}
}
