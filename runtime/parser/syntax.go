package parser

import (
	"github.com/fabiouggeri/page/util"
)

type ParserRuleType int

type Syntax struct {
	startRule       int
	lastNonTerminal int
	rulesNames      []string
	rulesTable      [][]int
}

const (
	AND_RULE          ParserRuleType = 0
	OR_RULE           ParserRuleType = 1
	ONE_OR_MORE_RULE  ParserRuleType = 2
	ZERO_OR_MORE_RULE ParserRuleType = 3
	OPTIONAL_RULE     ParserRuleType = 4
	TEST_NOT_RULE     ParserRuleType = 5
	TEST_RULE         ParserRuleType = 6
	TERMINAL_RULE     ParserRuleType = 7
	NON_TERMINAL_RULE ParserRuleType = 8
)

func SyntaxNew(totalRules int, lastNonTerminal int) *Syntax {
	return &Syntax{
		startRule:       -1,
		lastNonTerminal: lastNonTerminal,
		rulesNames:      make([]string, totalRules),
		rulesTable:      make([][]int, totalRules),
	}
}

func (s *Syntax) SetStartRule(index int) {
	s.startRule = index
}

func (s *Syntax) StartRule() int {
	return s.startRule
}

func (s *Syntax) Set(index int, name string, rules []int) {
	s.rulesNames[index] = name
	s.rulesTable[index] = rules
}

func (s *Syntax) RulesCount() int {
	return len(s.rulesNames)
}

func (s *Syntax) Rule(index int) (string, []int) {
	return s.rulesNames[index], s.rulesTable[index]
}

func (s *Syntax) Subrules(index int) []int {
	return s.rulesTable[index]
}

func (s *Syntax) RuleName(index int) string {
	return s.rulesNames[index]
}

func (s *Syntax) IsSubRule(index int) bool {
	return index > s.lastNonTerminal
}

func (s *Syntax) LastNonTerminal() int {
	return s.lastNonTerminal
}

func (s *Syntax) Write(writer *util.StringCodeWriter) {
	ruleName := s.rulesNames[s.startRule]
	writer.WriteString("Start Rule: ").WriteString(ruleName).NewLine()
	for i, rules := range s.rulesTable {
		ruleName = s.rulesNames[i]
		writer.WriteString(ruleName).WriteString(" -> ")
		writer.WriteString(ruleType(rules[0]))
		if ParserRuleType(rules[0]) != TERMINAL_RULE {
			for _, rule := range rules[1:] {
				ruleName = s.rulesNames[rule]
				writer.WriteString(" ").WriteString(ruleName)
			}
		}
		writer.NewLine()
	}
}

func ruleType(ruleType int) string {
	switch ParserRuleType(ruleType) {
	case AND_RULE:
		return "AND"
	case OR_RULE:
		return "OR"
	case ONE_OR_MORE_RULE:
		return "ONE_OR_MORE"
	case ZERO_OR_MORE_RULE:
		return "ZERO_OR_MORE"
	case OPTIONAL_RULE:
		return "OPTIONAL"
	case TEST_NOT_RULE:
		return "TEST_NOT"
	case TEST_RULE:
		return "TEST"
	case TERMINAL_RULE:
		return "TERMINAL"
	case NON_TERMINAL_RULE:
		return "NON_TERMINAL"
	default:
		return "UNKNOWN"
	}
}
