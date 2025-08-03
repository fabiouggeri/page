package vocabulary

import (
	"unicode"

	"github.com/fabiouggeri/page/build/rule"
)

type maxSymbolVisitor struct {
	maxSymbol   rune
	visitedRule map[rule.Rule]struct{}
}

var _ rule.RuleVisitor = &maxSymbolVisitor{}

func newMaxSymbolVisitor() *maxSymbolVisitor {
	return &maxSymbolVisitor{
		maxSymbol:   0,
		visitedRule: make(map[rule.Rule]struct{}, 0),
	}
}

// VisitAndRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitAndRule(rule *rule.AndRule) {
	rules := rule.Rules()
	for _, r := range rules {
		r.Visit(m)
	}
}

// VisitOrRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitOrRule(rule *rule.OrRule) {
	rules := rule.Rules()
	for _, r := range rules {
		r.Visit(m)
	}
}

// VisitNonTerminal implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitNonTerminal(rule *rule.NonTerminalRule) {
	_, found := m.visitedRule[rule]
	if !found {
		m.visitedRule[rule] = struct{}{}
		rule.Rule().Visit(m)
	}
}

// VisitNotRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitNotRule(rule *rule.NotRule) {
	rule.Rule().Visit(m)
}

// VisitOneOrMoreRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitOneOrMoreRule(rule *rule.OneOrMoreRule) {
	rule.Rule().Visit(m)
}

// VisitOptionalRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitOptionalRule(rule *rule.OptionalRule) {
	rule.Rule().Visit(m)
}

// VisitTestRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitTestRule(rule *rule.TestRule) {
	rule.Rule().Visit(m)
}

// VisitZeroOrMoreRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitZeroOrMoreRule(rule *rule.ZeroOrMoreRule) {
	rule.Rule().Visit(m)
}

// VisitCharRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitCharRule(rule *rule.CharRule) {
	if rule.CaseSensitive() {
		if unicode.ToLower(rule.Char()) > m.maxSymbol {
			m.maxSymbol = unicode.ToLower(rule.Char())
		}
		if unicode.ToUpper(rule.Char()) > m.maxSymbol {
			m.maxSymbol = unicode.ToUpper(rule.Char())
		}
	} else if rule.Char() > m.maxSymbol {
		m.maxSymbol = rule.Char()
	}
}

// VisitRangeRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitRangeRule(rule *rule.RangeRule) {
	for c := rule.Start(); c <= rule.End(); c++ {
		if c > m.maxSymbol {
			m.maxSymbol = c
		}
	}
}

// VisitStringRule implements rule.LexerVisitor.
func (m *maxSymbolVisitor) VisitStringRule(rule *rule.StringRule) {
	runes := []rune(rule.Text())
	if rule.CaseSensitive() {
		for _, c := range runes {
			if unicode.ToLower(c) > m.maxSymbol {
				m.maxSymbol = unicode.ToLower(c)
			}
			if unicode.ToUpper(c) > m.maxSymbol {
				m.maxSymbol = unicode.ToUpper(c)
			}
		}
	} else {
		for _, c := range runes {
			if c > m.maxSymbol {
				m.maxSymbol = c
			}
		}
	}
}
