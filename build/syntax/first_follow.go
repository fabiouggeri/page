package syntax

import (
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/util"
)

type firstFollow struct {
	vocabulary     *lexer.Vocabulary
	rules          []*parserRule
	firstRules     map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule]
	followingRules map[*rule.NonTerminalRule]*followingRules
}

type followingRules struct {
	followingRules map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule]
}

var (
	EMPTY_RULE = &rule.NonTerminalRule{}
)

func computeFirstFollow(vocabulary *lexer.Vocabulary, rules []*parserRule) *firstFollow {
	ff := &firstFollow{
		vocabulary:     vocabulary,
		rules:          rules,
		firstRules:     make(map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule]),
		followingRules: make(map[*rule.NonTerminalRule]*followingRules),
	}
	ff.computeRulesFirst()
	ff.computeRulesFollow()
	return ff

}

func (ff *firstFollow) computeRulesFollow() {
	for _, nonTerminal := range ff.rules {
		following := &followingRules{
			followingRules: make(map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule]),
		}
		ff.followingRules[nonTerminal.rule] = following
		if andRule, ok := nonTerminal.rule.Rule().(*rule.AndRule); ok {
			ff.computeFollow(following, andRule.Rules())
		}
	}
}

func (ff *firstFollow) computeFollow(following *followingRules, rules []rule.Rule) {
	for i := 0; i < len(rules)-1; i++ {
		nextRules := util.NewSet[*rule.NonTerminalRule]()
		currentRules := ff.endRules(rules[i])
		j := i + 1
		firstRules := ff.initialRules(rules[j])
		nextRules.AddAll(firstRules.Items()...)
		for firstRules.Contains(EMPTY_RULE) && j+1 < len(rules) {
			j++
			firstRules = ff.initialRules(rules[j])
			nextRules.AddAll(firstRules.Items()...)
		}
		ff.addFollowRules(following, currentRules, nextRules)
	}
}

func (ff *firstFollow) addFollowRules(following *followingRules, currentRules *util.Set[*rule.NonTerminalRule], nextRules *util.Set[*rule.NonTerminalRule]) {
	for _, subrule := range currentRules.Items() {
		for _, nextRule := range nextRules.Items() {
			ff.addFirstNextRule(following, subrule, nextRule)
		}
	}
}

func (ff *firstFollow) addFirstNextRule(following *followingRules, subrule *rule.NonTerminalRule, nextRule *rule.NonTerminalRule) {
	if subrule != EMPTY_RULE && nextRule != EMPTY_RULE {
		if _, found := following.followingRules[subrule]; !found {
			following.followingRules[subrule] = util.NewSet[*rule.NonTerminalRule]()
		}
		if firstRules, found := ff.firstRules[nextRule]; found {
			for _, nextSubRule := range firstRules.Items() {
				following.followingRules[subrule].Add(nextSubRule)
			}
		} else if ff.vocabulary.TokenIndex(nextRule.Id()) >= 0 {
			following.followingRules[subrule].Add(nextRule)
		}
	}
}

func (ff *firstFollow) initialRules(r rule.Rule) *util.Set[*rule.NonTerminalRule] {
	switch castRule := r.(type) {
	case *rule.NonTerminalRule:
		return util.NewSet(castRule)
	case *rule.AndRule:
		result := util.NewSet[*rule.NonTerminalRule]()
		rules := castRule.Rules()
		index := 0
		firstRules := ff.initialRules(rules[index])
		result.AddAll(firstRules.Items()...)
		for firstRules.Contains(EMPTY_RULE) && index < len(rules)-1 {
			index++
			firstRules = ff.initialRules(rules[index])
			result.AddAll(firstRules.Items()...)
		}
		return result
	case *rule.OrRule:
		result := util.NewSet[*rule.NonTerminalRule]()
		for _, r := range castRule.Rules() {
			result.AddAll(ff.initialRules(r).Items()...)
		}
		return result
	case *rule.OneOrMoreRule:
		return ff.initialRules(castRule.Rule())
	case *rule.ZeroOrMoreRule:
		result := ff.initialRules(castRule.Rule())
		result.Add(EMPTY_RULE)
		return result
	case *rule.OptionalRule:
		result := ff.initialRules(castRule.Rule())
		result.Add(EMPTY_RULE)
		return result
	default:
		return util.NewSet[*rule.NonTerminalRule]()
	}
}

func (ff *firstFollow) endRules(r rule.Rule) *util.Set[*rule.NonTerminalRule] {
	switch castRule := r.(type) {
	case *rule.NonTerminalRule:
		//		if ff.vocabulary.TokenIndex(castRule.Id()) < 0 {
		return util.NewSet(castRule)
		// } else {
		// 	return util.NewSet[*rule.NonTerminalRule]()
		// }
	case *rule.AndRule:
		result := util.NewSet[*rule.NonTerminalRule]()
		rules := castRule.Rules()
		lastIndex := len(rules) - 1
		lastRules := ff.endRules(rules[lastIndex])
		result.AddAll(lastRules.Items()...)
		for lastRules.Contains(EMPTY_RULE) && lastIndex > 0 {
			lastIndex--
			lastRules = ff.endRules(rules[lastIndex])
			result.AddAll(lastRules.Items()...)
		}
		return result
	case *rule.OrRule:
		result := util.NewSet[*rule.NonTerminalRule]()
		for _, r := range castRule.Rules() {
			result.AddAll(ff.endRules(r).Items()...)
		}
		return result
	case *rule.OneOrMoreRule:
		return ff.endRules(castRule.Rule())
	case *rule.ZeroOrMoreRule:
		result := ff.endRules(castRule.Rule())
		result.Add(EMPTY_RULE)
		return result
	case *rule.OptionalRule:
		result := ff.endRules(castRule.Rule())
		result.Add(EMPTY_RULE)
		return result
	default:
		return util.NewSet[*rule.NonTerminalRule]()
	}
}

func (ff *firstFollow) computeRulesFirst() {
	for _, parserRule := range ff.rules {
		if ff.vocabulary.TokenIndex(parserRule.rule.Id()) >= 0 {
			ff.firstRules[parserRule.rule] = util.NewSet(parserRule.rule)
		} else {
			firstVisitor := &firstVisitor{
				visited:    make(map[*rule.NonTerminalRule]struct{}),
				vocabulary: ff.vocabulary,
				firstRules: util.NewSet[*rule.NonTerminalRule](),
			}
			parserRule.rule.Rule().Visit(firstVisitor)
			ff.firstRules[parserRule.rule] = firstVisitor.firstRules
		}
	}
}

func (ff *firstFollow) following(nonTerminalRule *rule.NonTerminalRule) map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule] {
	if rules, found := ff.followingRules[nonTerminalRule]; found {
		return rules.followingRules
	}
	return make(map[*rule.NonTerminalRule]*util.Set[*rule.NonTerminalRule])
}

func (ff *firstFollow) firsts(nonTerminalRule *rule.NonTerminalRule) *util.Set[*rule.NonTerminalRule] {
	if rules, found := ff.firstRules[nonTerminalRule]; found {
		return rules
	}
	return util.NewSet[*rule.NonTerminalRule]()
}
