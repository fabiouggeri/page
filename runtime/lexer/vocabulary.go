package lexer

import (
	"strings"

	"github.com/fabiouggeri/page/util"
)

type Vocabulary struct {
	tokensNames      []string
	transitionsTable [][]int
	tokensTypes      [][]int
}

func NewVocabulary(tokensNames []string, transitionsTable [][]int, tokensTypes [][]int) *Vocabulary {
	return &Vocabulary{
		tokensNames:      tokensNames,
		transitionsTable: transitionsTable,
		tokensTypes:      tokensTypes,
	}
}

func (v *Vocabulary) TokensNames() []string {
	return v.tokensNames
}

func (v *Vocabulary) TransitionsTable() [][]int {
	return v.transitionsTable
}

func (v *Vocabulary) TokenTypes(index int) []int {
	if index < 0 || index >= len(v.tokensTypes) {
		return []int{}
	}
	return v.tokensTypes[index]
}

func (v *Vocabulary) String() string {
	s := strings.Builder{}
	for _, t := range v.tokensNames {
		s.WriteString("Token: " + t + "\n")
	}
	return s.String()
}

func (v *Vocabulary) TokenName(index int) string {
	if index < 0 || index >= len(v.tokensNames) {
		return ""
	}
	return v.tokensNames[index]
}

func (v *Vocabulary) Write(writer util.TextWriter) {
	v.writeTokensNames(writer)
	writer.NewLine()
	v.writeTokensTypes(writer)
	writer.NewLine()
	v.writeTransitionsTable(writer)
}

func (v *Vocabulary) writeTokensNames(writer util.TextWriter) {
	writer.WriteString("Tokens:").NewLine()
	writer.WriteString("=======").NewLine()
	if len(v.tokensNames) == 0 {
		return
	}
	writer.Indent(3)
	for _, t := range v.tokensNames {
		writer.WriteString(t).NewLine()
	}
	writer.Indent(-3)
}

func (v *Vocabulary) writeTokensTypes(writer util.TextWriter) {
	writer.WriteString("Tokens Types:").NewLine()
	writer.WriteString("=============").NewLine()
	if len(v.tokensTypes) == 0 {
		return
	}
	writer.Indent(3)
	for stateIndex, types := range v.tokensTypes {
		writer.WriteF("%d: ", stateIndex)
		for typeIndex, t := range types {
			if typeIndex > 0 {
				writer.WriteRune(',')
			}
			writer.WriteString(v.TokenName(t))
		}
		writer.NewLine()
	}
	writer.Indent(-3)
}

func (v *Vocabulary) writeTransitionsTable(writer util.TextWriter) {
	writer.WriteString("Transitions Table:").NewLine()
	writer.WriteString("==================").NewLine()
	if len(v.transitionsTable) == 0 {
		return
	}
	writer.Indent(3)
	writer.WriteString("State")
	for symbol := range v.transitionsTable[0] {
		v.writeSymbol(writer, rune(symbol))
	}
	writer.NewLine()
	for state, stateTransitions := range v.transitionsTable {
		writer.WriteF("%5d", state)
		for _, symbol := range stateTransitions {
			writer.WriteF(" %3d", symbol)
		}
		writer.NewLine()
	}
	writer.Indent(-3)
}

func (v *Vocabulary) writeSymbol(writer util.TextWriter, symbol rune) {
	if symbol <= 32 || symbol >= 127 {
		writer.WriteF(" %3d", symbol)
	} else {
		writer.WriteF(" '%c'", symbol)
	}
}
