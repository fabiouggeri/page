package syntax

import (
	"github.com/fabiouggeri/page/build/automata"
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/util"
)

type nfaVisitor struct {
	nextId int32
	states *util.Deque[*automata.State]
}

var _ rule.RuleVisitor = &nfaVisitor{}

func createNFA(parserRule *parserRule) *automata.State {
	visitor := &nfaVisitor{}
	parserRule.rule.Visit(visitor)
	//s1 := visitor.newInitialState()
	return nil
}

func (n *nfaVisitor) newInitialState() *automata.State {
	s := automata.NewState(n.nextId, true, false)
	n.nextId++
	return s
}

func (n *nfaVisitor) newFinalState() *automata.State {
	s := automata.NewState(n.nextId, false, true)
	n.nextId++
	return s
}

func (n *nfaVisitor) VisitNonTerminal(rule *rule.NonTerminalRule) {
}

func (n *nfaVisitor) VisitAndRule(rule *rule.AndRule) {
	rules := rule.Rules()
	for _, r := range rules {
		r.Visit(n)
	}
}

func (n *nfaVisitor) VisitOrRule(rule *rule.OrRule) {
}

func (n *nfaVisitor) VisitOneOrMoreRule(rule *rule.OneOrMoreRule) {
}

func (n *nfaVisitor) VisitOptionalRule(rule *rule.OptionalRule) {
}

func (n *nfaVisitor) VisitZeroOrMoreRule(rule *rule.ZeroOrMoreRule) {
}

func (n *nfaVisitor) VisitNotRule(rule *rule.NotRule) {
}

func (n *nfaVisitor) VisitTestRule(rule *rule.TestRule) {
}

func (n *nfaVisitor) VisitCharRule(rule *rule.CharRule) {
}

func (n *nfaVisitor) VisitRangeRule(rule *rule.RangeRule) {
}

func (n *nfaVisitor) VisitStringRule(rule *rule.StringRule) {
}
