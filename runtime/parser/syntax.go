package parser

import (
	"strings"

	"github.com/fabiouggeri/page/util"
)

type ParserRuleType int
type ParserRuleOption uint16

type Syntax struct {
	startRule       int
	lastNonTerminal int
	rulesNames      []string
	rulesTable      [][]int
	rulesOptions    []ParserRuleOption
	firstTable      [][]int
	followTables    []FollowTable
}

type FollowTable struct {
	rulesFollow []RuleFollow
}
type RuleFollow struct {
	rule        int
	rulesFollow []int
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

const (
	SKIP_NODE ParserRuleOption = 0x0001
	MEMOIZE   ParserRuleOption = 0x0002
	IGNORE    ParserRuleOption = 0x0004
)

func SyntaxNew(totalRules int, lastNonTerminal int) *Syntax {
	return &Syntax{
		startRule:       -1,
		lastNonTerminal: lastNonTerminal,
		rulesNames:      make([]string, totalRules),
		rulesTable:      make([][]int, totalRules),
		rulesOptions:    make([]ParserRuleOption, totalRules),
		firstTable:      make([][]int, totalRules),
		followTables:    make([]FollowTable, totalRules),
	}
}

func NewRuleFollow(ruleId int, followIds []int) RuleFollow {
	return RuleFollow{
		rule:        ruleId,
		rulesFollow: followIds,
	}
}

func (s *Syntax) SetStartRule(ruleId int) {
	s.startRule = ruleId
}

func (s *Syntax) StartRule() int {
	return s.startRule
}

func (s *Syntax) Set(ruleId int, name string, rules []int) {
	s.rulesNames[ruleId] = name
	s.rulesTable[ruleId] = rules
}

func (s *Syntax) RulesCount() int {
	return len(s.rulesNames)
}

func (s *Syntax) Rule(ruleId int) (string, []int) {
	return s.rulesNames[ruleId], s.rulesTable[ruleId]
}

func (s *Syntax) Subrules(ruleId int) []int {
	return s.rulesTable[ruleId]
}

func (s *Syntax) RuleName(ruleId int) string {
	return s.rulesNames[ruleId]
}

func (s *Syntax) RuleId(name string) int {
	for i, ruleName := range s.rulesNames {
		if strings.EqualFold(ruleName, name) {
			return i
		}
	}
	return -1
}

func (s *Syntax) SetFollow(ruleId int, follow []RuleFollow) {
	s.followTables[ruleId] = FollowTable{
		rulesFollow: follow,
	}
}

func (s *Syntax) SetFirst(ruleId int, firstIds []int) {
	s.firstTable[ruleId] = firstIds
}

func (s *Syntax) IsSubRule(index int) bool {
	return index > s.lastNonTerminal
}

func (s *Syntax) LastNonTerminal() int {
	return s.lastNonTerminal
}

func (s *Syntax) SetOption(index int, option ParserRuleOption) {
	s.rulesOptions[index] |= option
}

func (s *Syntax) Options(index int) ParserRuleOption {
	return s.rulesOptions[index]
}

func (s *Syntax) HasOption(index int, option ParserRuleOption) bool {
	return s.rulesOptions[index]&option != 0
}

func (s *Syntax) Write(writer *util.StringCodeWriter) {
	ruleName := s.rulesNames[s.startRule]
	writer.WriteString("Start Rule: ").
		WriteString(ruleName).
		NewLine().
		NewLine()
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
	writer.NewLine().WriteString("First:").
		NewLine().
		WriteString("======").
		NewLine()
	for i := range s.firstTable {
		if len(s.firstTable[i]) == 0 {
			continue
		}
		ruleName = s.rulesNames[i]
		writer.WriteString(ruleName).WriteString(" -> ")
		for j, firstRule := range s.firstTable[i] {
			if j > 0 {
				writer.WriteString(", ")
			}
			writer.WriteString(s.rulesNames[firstRule]).WriteString(" ")
		}
		writer.NewLine()
	}
	writer.NewLine().WriteString("Follow:").
		NewLine().
		WriteString("=======").
		NewLine()
	for i := range s.followTables {
		ruleName = s.rulesNames[i]
		if len(s.followTables[i].rulesFollow) == 0 {
			continue
		}
		writer.WriteString(ruleName).WriteString(" -> ")
		for j, follow := range s.followTables[i].rulesFollow {
			if j > 0 {
				writer.WriteString(", ")
			}
			writer.WriteString(s.rulesNames[follow.rule]).WriteString(" : {")
			for k, followRule := range follow.rulesFollow {
				if k > 0 {
					writer.WriteString(", ")
				}
				writer.WriteString(s.rulesNames[followRule])
			}
			writer.WriteRune('}')
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
