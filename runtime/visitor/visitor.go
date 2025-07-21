package visitor

import (
	"fmt"

	"github.com/fabiouggeri/page/runtime/parser"
)

type RuleVisitor struct {
	syntax             *parser.Syntax
	enterRuleCallbacks map[int]func(parser *parser.Parser, node *parser.ASTNode)
	exitRuleCallbacks  map[int]func(parser *parser.Parser, node *parser.ASTNode)
}

func New(s *parser.Syntax) *RuleVisitor {
	return &RuleVisitor{
		syntax:             s,
		enterRuleCallbacks: make(map[int]func(parser *parser.Parser, node *parser.ASTNode)),
		exitRuleCallbacks:  make(map[int]func(parser *parser.Parser, node *parser.ASTNode)),
	}
}

func (l *RuleVisitor) EnterRule(ruleId int, callback func(parser *parser.Parser, node *parser.ASTNode)) error {
	if ruleId < 0 || ruleId > l.syntax.LastNonTerminal() {
		return fmt.Errorf("rule id %d not found", ruleId)
	}
	l.enterRuleCallbacks[ruleId] = callback
	return nil
}

func (l *RuleVisitor) EnterRuleName(ruleName string, callback func(parser *parser.Parser, node *parser.ASTNode)) error {
	ruleId := l.syntax.RuleId(ruleName)
	if ruleId < 0 {
		return fmt.Errorf("rule '%s' not found", ruleName)
	}
	return l.EnterRule(ruleId, callback)
}

func (l *RuleVisitor) ExitRule(ruleId int, callback func(parser *parser.Parser, node *parser.ASTNode)) error {
	if ruleId < 0 || ruleId > l.syntax.LastNonTerminal() {
		return fmt.Errorf("rule id %d not found", ruleId)
	}
	l.exitRuleCallbacks[ruleId] = callback
	return nil
}

func (l *RuleVisitor) ExitRuleName(ruleName string, callback func(parser *parser.Parser, node *parser.ASTNode)) error {
	ruleId := l.syntax.RuleId(ruleName)
	if ruleId < 0 {
		return fmt.Errorf("rule '%s' not found", ruleName)
	}
	return l.ExitRule(ruleId, callback)
}
