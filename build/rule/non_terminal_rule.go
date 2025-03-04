package rule

import "github.com/fabiouggeri/page/util"

type NonTerminalRule struct {
	id      string
	rule    Rule
	options map[*RuleOption]string
}

var _ SimpleRule = &NonTerminalRule{}

func (r *NonTerminalRule) Id() string {
	return r.id
}

func (r *NonTerminalRule) Rule() Rule {
	return r.rule
}

func (r *NonTerminalRule) SetRule(rule Rule) {
	r.rule = rule
}

func (r *NonTerminalRule) ToText(writer util.TextWriter) {
	writer.WriteString(r.id)
}

func (r *NonTerminalRule) Visit(visitor LexerVisitor) {
	visitor.VisitNonTerminal(r)
}

func (r *NonTerminalRule) String() string {
	str := util.NewStringTextWriter()
	r.ToText(str)
	return str.String()
}

func (r *NonTerminalRule) Option(option *RuleOption, value string) *NonTerminalRule {
	r.options[option] = value
	return r
}

func (r *NonTerminalRule) DelOption(option *RuleOption) *NonTerminalRule {
	delete(r.options, option)
	return r
}

func (r *NonTerminalRule) HasOption(option *RuleOption) bool {
	_, found := r.options[option]
	return found
}

func (r *NonTerminalRule) GetOption(option *RuleOption) (string, bool) {
	value, found := r.options[option]
	return value, found
}

func (r *NonTerminalRule) Options() []*RuleOption {
	options := make([]*RuleOption, 0, len(r.options))
	for k := range r.options {
		options = append(options, k)
	}
	return options
}

func (r *NonTerminalRule) WalkThrough(visit func(r Rule), isVisit func(r Rule) bool) {
	wv := newWalkerVisitor(visit, isVisit)
	r.Visit(wv)
}

func (r *NonTerminalRule) IsLexer() bool {
	if r.HasOption(FRAGMENT) {
		return false
	} else if r.HasOption(ATOMIC) || r.HasOption(TOKEN) {
		return true
	}
	return checkIsLexerRule(make(map[Rule]bool, 64), r.Rule())
}

func compositeIsLexerRule(rulesMap map[Rule]bool, r CompoundRule) bool {
	rules := r.Rules()
	for _, component := range rules {
		switch castRule := component.(type) {
		case *AndRule:
			if allRulesAreLiteral(castRule.Rules()) {
				return false
			} else if !compositeIsLexerRule(rulesMap, castRule) {
				return false
			}
		case *OrRule:
			if !compositeIsLexerRule(rulesMap, castRule) {
				return false
			}
		case *NonTerminalRule:
			if !castRule.HasOption(FRAGMENT) {
				return false
			}
		default:
			if !checkIsLexerRule(rulesMap, castRule) {
				return false
			}
		}
	}
	return true
}

func allRulesAreLiteral(rules []Rule) bool {
	for _, component := range rules {
		switch component.(type) {
		case *StringRule:
		case *CharRule:
		default:
			return false
		}
	}
	return true
}

func checkIsLexerRule(rulesMap map[Rule]bool, r Rule) bool {
	lexerRule, found := rulesMap[r]
	if !found {
		switch castRule := r.(type) {
		case *StringRule,
			*CharRule,
			*RangeRule:
			lexerRule = true
		case *AndRule:
			if allRulesAreLiteral(castRule.Rules()) {
				lexerRule = false
			} else {
				lexerRule = compositeIsLexerRule(rulesMap, castRule)
			}
		case *OrRule:
			lexerRule = compositeIsLexerRule(rulesMap, castRule)
		case *NonTerminalRule:
			if castRule.HasOption(FRAGMENT) || castRule.HasOption(ATOMIC) || castRule.HasOption(TOKEN) {
				lexerRule = true
			} else {
				lexerRule = checkIsLexerRule(rulesMap, castRule.Rule())
			}
		case *OneOrMoreRule:
			lexerRule = checkIsLexerRule(rulesMap, castRule.Rule())
		case *ZeroOrMoreRule:
			lexerRule = checkIsLexerRule(rulesMap, castRule.Rule())
		case *OptionalRule:
			lexerRule = checkIsLexerRule(rulesMap, castRule.Rule())
		case *NotRule:
			lexerRule = checkIsLexerRule(rulesMap, castRule.Rule())
		case *TestRule:
			lexerRule = checkIsLexerRule(rulesMap, castRule.Rule())
		default:
			lexerRule = false
		}
		rulesMap[r] = lexerRule
	}
	return lexerRule
}
