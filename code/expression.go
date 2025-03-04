package code

type Expression interface {
	Code
	Add(other Expression) Expression
	Append(other Expression) Expression
	Minus(other Expression) Expression
	Mult(other Expression) Expression
	Div(other Expression) Expression
	Equals(other Expression) Expression
	Diff(other Expression) Expression
	And(other Expression) Expression
	Or(other Expression) Expression
	Not() Expression
	Call(name string, args ...Expression) Expression
	Less(other Expression) Expression
	LessEqual(other Expression) Expression
	Greater(other Expression) Expression
	GreaterEqual(other Expression) Expression
	Index(other Expression) Expression
	Inc() Expression
	Dec() Expression
	Field(name string) LeftValue
}

type LeftValue interface {
	Expression
	Assign(other Expression) Expression
}

type Operator int

const (
	PLUS Operator = iota
	MINUS
	MULT
	DIV
	EQUALS
	ASSIGN
	DIFF
	AND
	OR
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	NOT
	INC
	DEC
	ADDRESS_OF
	INDEX
	APPEND
)

func (op Operator) String() string {
	switch op {
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case MULT:
		return "*"
	case DIV:
		return "/"
	case EQUALS:
		return "=="
	case ASSIGN:
		return "="
	case DIFF:
		return "!="
	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"
	case INC:
		return "++"
	case DEC:
		return "--"
	case ADDRESS_OF:
		return "&"
	case INDEX:
		return "[]"
	default:
		return "unknown"
	}
}
