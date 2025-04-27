package syntax

import (
	"strconv"

	"github.com/fabiouggeri/page/build/grammar"
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/runtime/parser"
	"github.com/fabiouggeri/page/util"
)

type syntaxBuilder struct {
	syntax        *parser.Syntax
	parserRules   map[string]*parserRule
	currentRule   *parserRule
	rulesBuilding *util.Deque[*parserRule]
	vocabulary    *lexer.Vocabulary
	nextId        int
}
type parserRule struct {
	rule  *rule.NonTerminalRule
	id    int
	name  string
	rules []int
	lexer bool
}

var _ rule.RuleVisitor = &syntaxBuilder{}

func FromGrammar(g *grammar.Grammar, vocabulary *lexer.Vocabulary) *parser.Syntax {
	builder := &syntaxBuilder{
		parserRules: make(map[string]*parserRule, 0),
		vocabulary:  vocabulary,
	}
	builder.build(g)
	return builder.syntax
}

func (b *syntaxBuilder) build(g *grammar.Grammar) {
	allRules := make([]*parserRule, 0, len(g.ParserRules())+len(g.LexerRules()))
	grammarParserRules := g.ParserRules()
	for _, r := range grammarParserRules {
		parserRule := &parserRule{
			rule:  r,
			id:    len(b.parserRules),
			name:  r.Id(),
			rules: make([]int, 0),
		}
		b.parserRules[r.Id()] = parserRule
		allRules = append(allRules, parserRule)
	}
	grammarLexerRules := g.LexerRules()
	for _, r := range grammarLexerRules {
		parserRule := &parserRule{
			rule:  r,
			id:    len(b.parserRules),
			name:  r.Id(),
			rules: make([]int, 0),
			lexer: true,
		}
		b.parserRules[r.Id()] = parserRule
		allRules = append(allRules, parserRule)
	}
	lastRuleId := len(allRules) - 1
	for _, r := range allRules {
		b.createSyntax(r, lastRuleId)
	}
	b.syntax = parser.SyntaxNew(len(b.parserRules), lastRuleId)
	for _, v := range b.parserRules {
		b.syntax.Set(v.id, v.name, v.rules)
		if b.isMainRule(g, v.name) {
			b.syntax.SetStartRule(v.id)
		}
	}
}

func (b *syntaxBuilder) isMainRule(g *grammar.Grammar, k string) bool {
	mainRule := g.GetRule(k)
	return mainRule != nil && mainRule.Id() == g.MainRule().Id()
}

func (b *syntaxBuilder) createSyntax(r *parserRule, lastGrammarRuleId int) {
	if r.lexer {
		tokenIndex := b.vocabulary.TokenIndex(r.rule.Id())
		if tokenIndex < 0 {
			panic("Rule not found: " + r.rule.Id())
		}
		r.rules = append(r.rules, int(parser.TERMINAL_RULE), tokenIndex)
		return
	}
	b.currentRule = r
	b.nextId = 1
	b.rulesBuilding = util.NewDeque[*parserRule]()
	r.rule.Rule().Visit(b)
	lastRule, err := b.rulesBuilding.Pop()
	if err != nil {
		panic("Error building syntax: " + err.Error())
	}
	if lastRule.id <= lastGrammarRuleId {
		r.rules = append(r.rules, int(parser.NON_TERMINAL_RULE), lastRule.id)
	} else {
		r.rules = append(r.rules, lastRule.rules...)
		delete(b.parserRules, lastRule.name)
	}
}

// VisitAndRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitAndRule(rule *rule.AndRule) {
	b.createCompoundRule(parser.AND_RULE, rule)
}

// VisitOrRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitOrRule(rule *rule.OrRule) {
	b.createCompoundRule(parser.OR_RULE, rule)
}

func (b *syntaxBuilder) createCompoundRule(ruleType parser.ParserRuleType, rule rule.CompoundRule) {
	for _, r := range rule.Rules() {
		r.Visit(b)
	}
	rules := make([]int, len(rule.Rules())+1)
	rules[0] = int(ruleType)
	index := len(rule.Rules())
	for range rule.Rules() {
		r, err := b.rulesBuilding.Pop()
		if err != nil {
			panic("Error building syntax: " + err.Error())
		}
		rules[index] = r.id
		index--
	}
	newRule := &parserRule{
		id:    len(b.parserRules),
		name:  b.currentRule.name + "#" + strconv.Itoa(b.nextId),
		rules: rules,
	}
	b.rulesBuilding.Push(newRule)
	b.parserRules[newRule.name] = newRule
	b.nextId++
}

// VisitNonTerminal implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitNonTerminal(rule *rule.NonTerminalRule) {
	parserRule, ok := b.parserRules[rule.Id()]
	if !ok {
		panic("Rule not found: " + rule.Id())
	}
	b.rulesBuilding.Push(parserRule)
}

func (b *syntaxBuilder) createSimpleRule(ruleType parser.ParserRuleType, rule rule.SimpleRule) {
	rule.Rule().Visit(b)
	rules := make([]int, 0, 2)
	r, err := b.rulesBuilding.Pop()
	if err != nil {
		panic("Error building syntax: " + err.Error())
	}
	rules = append(rules, int(ruleType), r.id)
	newRule := &parserRule{
		id:    len(b.parserRules),
		name:  b.currentRule.name + "#" + strconv.Itoa(b.nextId),
		rules: rules,
	}
	b.rulesBuilding.Push(newRule)
	b.parserRules[newRule.name] = newRule
	b.nextId++
}

// VisitNotRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitNotRule(rule *rule.NotRule) {
	b.createSimpleRule(parser.TEST_NOT_RULE, rule)
}

// VisitOneOrMoreRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitOneOrMoreRule(rule *rule.OneOrMoreRule) {
	b.createSimpleRule(parser.ONE_OR_MORE_RULE, rule)
}

// VisitOptionalRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitOptionalRule(rule *rule.OptionalRule) {
	b.createSimpleRule(parser.OPTIONAL_RULE, rule)
}

// VisitTestRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitTestRule(rule *rule.TestRule) {
	b.createSimpleRule(parser.TEST_RULE, rule)
}

// VisitZeroOrMoreRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitZeroOrMoreRule(rule *rule.ZeroOrMoreRule) {
	b.createSimpleRule(parser.ZERO_OR_MORE_RULE, rule)
}

// VisitCharRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitCharRule(rule *rule.CharRule) {
	panic("Not a parser rule")
}

// VisitRangeRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitRangeRule(rule *rule.RangeRule) {
	panic("Not a parser rule")
}

// VisitStringRule implements rule.RuleVisitor.
func (b *syntaxBuilder) VisitStringRule(rule *rule.StringRule) {
	panic("Not a parser rule")
}
