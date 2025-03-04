package parser

type Syntax struct {
	nonTerminalNames []string
	nonTerminalTypes []int
	transitionsTable [][]int
}
