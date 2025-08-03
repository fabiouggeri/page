package syntax

import (
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/util"
)

type firstVisitor struct {
	visited    map[*rule.NonTerminalRule]struct{}
	vocabulary *lexer.Vocabulary
	firstRules *util.Set[*rule.NonTerminalRule]
}

var _ rule.RuleVisitor = &firstVisitor{}

// VisitNonTerminal implements rule.RuleVisitor.
func (f *firstVisitor) VisitNonTerminal(rule *rule.NonTerminalRule) {
	if f.vocabulary.TokenIndex(rule.Id()) >= 0 {
		f.firstRules.Add(rule)
		return
	}
	if _, found := f.visited[rule]; found {
		return
	}
	f.visited[rule] = struct{}{}
	rule.Rule().Visit(f)
}

// VisitAndRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitAndRule(rule *rule.AndRule) {
	for _, r := range rule.Rules() {
		r.Visit(f)
		if !f.firstRules.Contains(EMPTY_RULE) {
			return
		}
	}
}

// VisitOrRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitOrRule(rule *rule.OrRule) {
	for _, r := range rule.Rules() {
		r.Visit(f)
	}
}

// VisitOneOrMoreRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitOneOrMoreRule(rule *rule.OneOrMoreRule) {
	rule.Rule().Visit(f)
}

// VisitOptionalRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitOptionalRule(rule *rule.OptionalRule) {
	f.firstRules.Add(EMPTY_RULE)
	rule.Rule().Visit(f)
}

// VisitZeroOrMoreRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitZeroOrMoreRule(rule *rule.ZeroOrMoreRule) {
	f.firstRules.Add(EMPTY_RULE)
	rule.Rule().Visit(f)
}

// VisitNotRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitNotRule(rule *rule.NotRule) {
}

// VisitCharRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitCharRule(rule *rule.CharRule) {
}

// VisitRangeRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitRangeRule(rule *rule.RangeRule) {
}

// VisitStringRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitStringRule(rule *rule.StringRule) {
}

// VisitTestRule implements rule.RuleVisitor.
func (f *firstVisitor) VisitTestRule(rule *rule.TestRule) {
}
