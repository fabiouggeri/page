package automata

import (
	"strconv"

	"github.com/fabiouggeri/page/build/rule"
)

type RuleType struct {
	id   uint16
	name string
	rule *rule.NonTerminalRule
}

func NewRuleId(id uint16, name string, rule *rule.NonTerminalRule) *RuleType {
	return &RuleType{id: id,
		name: name,
		rule: rule}
}

func (t *RuleType) Id() uint16 {
	return t.id
}

func (t *RuleType) Name() string {
	return t.name
}

func (t *RuleType) Rule() *rule.NonTerminalRule {
	return t.rule
}

func (t *RuleType) String() string {
	return strconv.Itoa(int(t.id)) + " - " + t.name + ": " + t.rule.String()
}
