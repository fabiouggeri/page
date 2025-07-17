package rule

import "github.com/fabiouggeri/page/util"

type TokenId string

type Rule interface {
	ToText(writer util.TextWriter)
	String() string
	Visit(visitor RuleVisitor)
}

type TerminalRule interface {
	Rule
	Text() string
	Size() int32
	CaseSensitive() bool
}

type SimpleRule interface {
	Rule
	Rule() Rule
	SetRule(rule Rule)
}

type CompoundRule interface {
	Rule
	Rules() []Rule
	SetRule(index int, rule Rule)
}

type RuleVisitor interface {
	VisitNonTerminal(rule *NonTerminalRule)
	VisitAndRule(rule *AndRule)
	VisitOrRule(rule *OrRule)
	VisitZeroOrMoreRule(rule *ZeroOrMoreRule)
	VisitOneOrMoreRule(rule *OneOrMoreRule)
	VisitOptionalRule(rule *OptionalRule)
	VisitCharRule(rule *CharRule)
	VisitRangeRule(rule *RangeRule)
	VisitStringRule(rule *StringRule)
	VisitTestRule(rule *TestRule)
	VisitNotRule(rule *NotRule)
}

func New(id string, rule Rule) *NonTerminalRule {
	return &NonTerminalRule{
		id:      id,
		rule:    rule,
		options: make(map[*RuleOption]string),
	}
}

func String(text string) TerminalRule {
	if len(text) == 1 {
		return &CharRule{char: rune(text[0]), caseSensitive: true}
	} else {
		return &StringRule{text: text, size: int32(len(text)), caseSensitive: true}
	}
}

func StringPartial(text string, size int32) TerminalRule {
	if len(text) == 1 {
		return &CharRule{char: rune(text[0]), caseSensitive: true}
	} else {
		return &StringRule{text: text, size: size, caseSensitive: true}
	}
}

func StringI(text string) TerminalRule {
	if len(text) == 1 {
		return &CharRule{char: rune(text[0]), caseSensitive: false}
	} else {
		return &StringRule{text: text, size: int32(len(text)), caseSensitive: false}
	}
}

func StringPartialI(text string, size int32) TerminalRule {
	if len(text) == 1 {
		return &CharRule{char: rune(text[0]), caseSensitive: false}
	} else {
		return &StringRule{text: text, size: size, caseSensitive: false}
	}
}

func Char(char rune) *CharRule {
	return &CharRule{char: char, caseSensitive: true}
}

func CharI(char rune) *CharRule {
	return &CharRule{char: char, caseSensitive: false}
}

func And(rules ...Rule) *AndRule {
	rule := &AndRule{rules: make([]Rule, 0, len(rules))}
	rule.rules = append(rule.rules, rules...)
	return rule
}

func Or(rules ...Rule) *OrRule {
	rule := &OrRule{rules: make([]Rule, 0, len(rules))}
	rule.rules = append(rule.rules, rules...)
	return rule
}

func Range(start rune, end rune) *RangeRule {
	return &RangeRule{start: start, end: end}
}

func OneOrMore(rule Rule) *OneOrMoreRule {
	return &OneOrMoreRule{rule: rule}
}

func ZeroOrMore(rule Rule) *ZeroOrMoreRule {
	return &ZeroOrMoreRule{rule: rule}
}

func Optional(rule Rule) *OptionalRule {
	return &OptionalRule{rule: rule}
}

func Test(rule Rule) *TestRule {
	return &TestRule{rule: rule}
}

func Not(rule Rule) *NotRule {
	return &NotRule{rule: rule}
}
