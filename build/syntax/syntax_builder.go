package syntax

import (
	"strconv"

	"github.com/fabiouggeri/page/build/grammar"
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/runtime/parser"
	"github.com/fabiouggeri/page/util"
	"golang.org/x/exp/maps"
)

type syntaxBuilder struct {
	syntax            *parser.Syntax
	parserRules       map[string]*parserRule
	currentRule       *parserRule
	rulesBuilding     *util.Deque[*parserRule]
	vocabulary        *lexer.Vocabulary
	nextId            int
	lastGrammarRuleId int
}

type parserRule struct {
	rule        *rule.NonTerminalRule
	firstRules  *util.Set[*rule.NonTerminalRule]
	followRules map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule]
	id          int
	name        string
	rules       []int
	lexer       bool
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
	allRules := b.grammarRules(g)
	b.lastGrammarRuleId = len(allRules) - 1
	for _, parserRule := range allRules {
		b.createSyntax(parserRule)
	}

	ff := computeFirstFollow(b.vocabulary, maps.Values(b.parserRules))
	for _, parserRule := range b.parserRules {
		parserRule.firstRules = ff.firsts(parserRule.rule)
		parserRule.followRules = ff.following(parserRule.rule)
		if !parserRule.lexer {
			b.createTranstionTable(parserRule)
		}
	}
	b.syntax = parser.SyntaxNew(len(b.parserRules), b.lastGrammarRuleId)
	for _, parserRule := range b.parserRules {
		b.syntax.Set(parserRule.id, parserRule.name, parserRule.rules)
		b.syntax.SetFirst(parserRule.id, b.firstRulesToId(parserRule.firstRules))
		b.syntax.SetFollow(parserRule.id, b.followRulesToId(parserRule.followRules))
		b.setOptions(parserRule)
		if b.isMainRule(g, parserRule.name) {
			b.syntax.SetStartRule(parserRule.id)
		}
	}
}

func (b *syntaxBuilder) followRulesToId(rules map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule]) []parser.RuleFollow {
	if len(rules) == 0 {
		return []parser.RuleFollow{}
	}
	followingIds := make([]parser.RuleFollow, 0, len(rules))
	for nonTerminal, following := range rules {
		parserRule, found := b.parserRules[nonTerminal.Id()]
		if !found {
			panic("Rule not found: " + nonTerminal.Id())
		}
		rulesIds := make([]int, 0, following.Length())
		for _, f := range following.Items() {
			if f != EMPTY_RULE {
				followingRule, found := b.parserRules[f.Id()]
				if !found {
					panic("Rule not found: " + f.Id())
				}
				rulesIds = append(rulesIds, followingRule.id)
			}
		}
		followingIds = append(followingIds, parser.NewRuleFollow(parserRule.id, rulesIds))
	}
	return followingIds
}

func (b *syntaxBuilder) firstRulesToId(rules *util.Set[*rule.NonTerminalRule]) []int {
	if rules == nil {
		return []int{}
	}
	ids := make([]int, 0, rules.Length())
	for _, first := range rules.Items() {
		if first != EMPTY_RULE {
			parserRule, found := b.parserRules[first.Id()]
			if !found {
				panic("Rule not found: " + first.Id())
			}
			ids = append(ids, parserRule.id)
		}
	}
	return ids
}

func (b *syntaxBuilder) createTranstionTable(parserRule *parserRule) {
	createNFA(parserRule)
}

func (b *syntaxBuilder) grammarRules(g *grammar.Grammar) []*parserRule {
	allRules := make([]*parserRule, 0, len(g.ParserRules())+len(g.LexerRules()))
	grammarLexerRules := g.LexerRules()
	for _, grammarRule := range grammarLexerRules {
		parserRule := &parserRule{
			rule:  grammarRule,
			id:    len(b.parserRules),
			name:  grammarRule.Id(),
			rules: make([]int, 0),
			lexer: true,
		}
		b.parserRules[grammarRule.Id()] = parserRule
		allRules = append(allRules, parserRule)
	}
	grammarParserRules := g.ParserRules()
	for _, grammarRule := range grammarParserRules {
		parserRule := &parserRule{
			rule:  grammarRule,
			id:    len(b.parserRules),
			name:  grammarRule.Id(),
			rules: make([]int, 0),
		}
		b.parserRules[grammarRule.Id()] = parserRule
		allRules = append(allRules, parserRule)
	}
	return allRules
}

func (b *syntaxBuilder) setOptions(v *parserRule) {
	if v.rule == nil {
		return
	}
	for _, option := range v.rule.Options() {
		switch option {
		case rule.IGNORE:
			b.syntax.SetOption(v.id, parser.IGNORE)
		case rule.SKIP_NODE:
			b.syntax.SetOption(v.id, parser.SKIP_NODE)
		case rule.MEMOIZE:
			b.syntax.SetOption(v.id, parser.MEMOIZE)
		default:
			// do nothing
		}
	}
}

func (b *syntaxBuilder) isMainRule(g *grammar.Grammar, k string) bool {
	mainRule := g.GetRule(k)
	return mainRule != nil && mainRule.Id() == g.MainRule().Id()
}

func (b *syntaxBuilder) createSyntax(r *parserRule) {
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
	if lastRule.id <= b.lastGrammarRuleId {
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

func (b *syntaxBuilder) createCompoundRule(ruleType parser.ParserRuleType, compoundRule rule.CompoundRule) {
	for _, r := range compoundRule.Rules() {
		r.Visit(b)
	}
	rules := make([]int, len(compoundRule.Rules())+1)
	rules[0] = int(ruleType)
	index := len(compoundRule.Rules())
	for range compoundRule.Rules() {
		r, err := b.rulesBuilding.Pop()
		if err != nil {
			panic("Error building syntax: " + err.Error())
		}
		rules[index] = r.id
		index--
	}
	newRule := b.findAuxiliarRule(rules)
	if newRule != nil {
		b.rulesBuilding.Push(newRule)
		return
	}
	ruleName := b.currentRule.name + "#" + strconv.Itoa(b.nextId)
	newRule = &parserRule{
		rule:  rule.New(ruleName, compoundRule),
		id:    len(b.parserRules),
		name:  ruleName,
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

func (b *syntaxBuilder) createSimpleRule(ruleType parser.ParserRuleType, simpleRule rule.SimpleRule) {
	simpleRule.Rule().Visit(b)
	rules := make([]int, 0, 2)
	r, err := b.rulesBuilding.Pop()
	if err != nil {
		panic("Error building syntax: " + err.Error())
	}
	rules = append(rules, int(ruleType), r.id)
	newRule := b.findAuxiliarRule(rules)
	if newRule != nil {
		b.rulesBuilding.Push(newRule)
		return
	}
	ruleName := b.currentRule.name + "#" + strconv.Itoa(b.nextId)
	newRule = &parserRule{
		rule:  rule.New(ruleName, simpleRule),
		id:    len(b.parserRules),
		name:  ruleName,
		rules: rules,
	}
	b.rulesBuilding.Push(newRule)
	b.parserRules[newRule.name] = newRule
	b.nextId++
}

func (b *syntaxBuilder) findAuxiliarRule(rules []int) *parserRule {
	for _, parserRule := range b.parserRules {
		if len(parserRule.rules) == len(rules) && parserRule.rules[0] == rules[0] && parserRule.id > b.lastGrammarRuleId {
			match := true
			for i, rule := range rules[1:] {
				if parserRule.rules[i+1] != rule {
					match = false
					break
				}
			}
			if match {
				return parserRule
			}
		}
	}
	return nil
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
