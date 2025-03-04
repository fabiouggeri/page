package grammar

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/util"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Grammar struct {
	name        string
	options     GrammarOptions
	encode      *charmap.Charmap
	mainRule    *rule.NonTerminalRule
	rules       map[string]*rule.NonTerminalRule
	lexerRules  *util.Set[*rule.NonTerminalRule]
	parserRules *util.Set[*rule.NonTerminalRule]
	errors      []error
}

type GrammarOptions struct {
	lexerName  string
	parserName string
}

var replacements = map[rune]rune{
	'á': 'a', 'à': 'a', 'ä': 'a', 'â': 'a', 'ã': 'a', 'å': 'a',
	'é': 'e', 'è': 'e', 'ë': 'e', 'ê': 'e',
	'í': 'i', 'ì': 'i', 'ï': 'i', 'î': 'i',
	'ó': 'o', 'ò': 'o', 'ö': 'o', 'ô': 'o', 'õ': 'o',
	'ú': 'u', 'ù': 'u', 'ü': 'u', 'û': 'u',
	'ç': 'c', 'ñ': 'n',
	'Á': 'A', 'À': 'A', 'Ä': 'A', 'Â': 'A', 'Ã': 'A', 'Å': 'A',
	'É': 'E', 'È': 'E', 'Ë': 'E', 'Ê': 'E',
	'Í': 'I', 'Ì': 'I', 'Ï': 'I', 'Î': 'I',
	'Ó': 'O', 'Ò': 'O', 'Ö': 'O', 'Ô': 'O', 'Õ': 'O',
	'Ú': 'U', 'Ù': 'U', 'Ü': 'U', 'Û': 'U',
	'Ç': 'C', 'Ñ': 'N',
}

var namedChars = map[rune]string{
	'(':  "open_par",
	')':  "close_par",
	'{':  "open_brace",
	'}':  "close_brace",
	'[':  "open_bracket",
	']':  "close_bracket",
	'<':  "open_angle",
	'>':  "close_angle",
	'.':  "dot",
	',':  "comma",
	';':  "semicolon",
	':':  "colon",
	'=':  "equal",
	'!':  "exclamation",
	'?':  "question",
	'+':  "plus",
	'-':  "minus",
	'*':  "asterisk",
	'/':  "slash",
	'%':  "percent",
	'&':  "ampersand",
	'|':  "pipe",
	'^':  "caret",
	'~':  "tilde",
	'@':  "at",
	'#':  "hash",
	'$':  "dollar",
	'_':  "underscore",
	'\\': "backslash",
	'`':  "backtick",
	'"':  "double_quote",
	'\'': "single_quote",
	' ':  "space",
	'\t': "tab",
	'\n': "newline",
	'\r': "carriage_return",
	'\f': "form_feed",
	'\v': "vertical_tab",
	'\b': "backspace",
}

func New(name string) *Grammar {
	g := &Grammar{name: name,
		rules: make(map[string]*rule.NonTerminalRule, 256),
	}
	g.options.lexerName = name + "Lexer"
	g.options.parserName = name + "Parser"
	return g
}

func FromBuffer(grammar []byte) (*Grammar, error) {
	return parseGrammar(grammar)
}

func FromString(text string) (*Grammar, error) {
	return FromBuffer([]byte(text))
}

func FromFile(pathname string) (*Grammar, error) {
	buffer, err := os.ReadFile(pathname)
	if err != nil {
		return nil, err
	}
	return FromBuffer(buffer)
}

func FromFileEncode(filePathname string, encode *charmap.Charmap) (*Grammar, error) {
	file, err := os.Open(filePathname)
	if err != nil {
		return nil, err
	}
	reader := transform.NewReader(file, encode.NewDecoder())
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	g, err := FromBuffer(buffer)
	if err != nil {
		return nil, err
	}
	g.encode = encode
	return g, nil
}

func (g *Grammar) Name() string {
	return g.name
}

func (g *Grammar) Options() *GrammarOptions {
	return &g.options
}

func (o *GrammarOptions) LexerName() string {
	return o.lexerName
}

func (o *GrammarOptions) ParserName() string {
	return o.parserName
}

func (g *Grammar) Rules(rules ...*rule.NonTerminalRule) error {
	for _, r := range rules {
		if _, found := g.rules[r.Id()]; !found {
			g.rules[r.Id()] = r
		} else {
			return fmt.Errorf("rule %s already defined", r.Id())
		}
	}
	g.lexerRules = nil
	g.parserRules = nil
	g.errors = nil
	return nil
}

func (g *Grammar) mapRules() {
	if g.lexerRules != nil && g.parserRules != nil {
		return
	}
	g.errors = make([]error, 0)
	g.lexerRules = util.NewSet[*rule.NonTerminalRule]()
	g.parserRules = util.NewSet[*rule.NonTerminalRule]()
	for _, r := range g.rules {
		if r.IsLexer() {
			g.lexerRules.Add(r)
		} else if !r.HasOption(rule.FRAGMENT) {
			g.parserRules.Add(r)
		}
	}
	lexerRulesMap := g.mapLexerRulesByName()
	for _, r := range g.parserRules.Items() {
		g.mapAnonymousRules(lexerRulesMap, r.Rule())
	}
	g.validateRules()
}

func (g *Grammar) validateRules() {
	lexerRules := g.lexerRules.Items()
	for _, r := range lexerRules {
		parserRules := g.parserRulesReferences(r)
		if len(parserRules) > 0 {
			for _, pr := range parserRules {
				g.errors = append(g.errors, fmt.Errorf("lexer rule '%s' references the parser rule '%s'", r.Id(), pr.Id()))
			}
		}
	}
}

func (g *Grammar) parserRulesReferences(r *rule.NonTerminalRule) []*rule.NonTerminalRule {
	nonTerminalRules := make([]*rule.NonTerminalRule, 0)
	r.WalkThrough(func(r rule.Rule) {
		switch castRule := r.(type) {
		case *rule.NonTerminalRule:
			nonTerminalRules = append(nonTerminalRules, castRule)
		default:
			// do nothing
		}
	}, nil)
	parserRules := make([]*rule.NonTerminalRule, 0, len(nonTerminalRules))
	for _, ntr := range nonTerminalRules {
		if g.parserRules.Contains(ntr) {
			parserRules = append(parserRules, ntr)
		}
	}
	return parserRules
}

func (g *Grammar) mapLexerRulesByName() map[string]rule.Rule {
	lexerRulesMap := make(map[string]rule.Rule, g.lexerRules.Length())
	for _, r := range g.lexerRules.Items() {
		switch castRule := r.Rule().(type) {
		case rule.TerminalRule:
			lexerRulesMap[castRule.String()] = r
		default:
			// do nothing
		}
	}
	return lexerRulesMap
}

func (g *Grammar) mapAnonymousRules(lexerRulesMap map[string]rule.Rule, r rule.Rule) rule.Rule {
	switch castRule := r.(type) {
	case *rule.AndRule:
		for i, r := range castRule.Rules() {
			castRule.SetRule(i, g.mapAnonymousRules(lexerRulesMap, r))
		}
	case *rule.OrRule:
		for i, r := range castRule.Rules() {
			castRule.SetRule(i, g.mapAnonymousRules(lexerRulesMap, r))
		}
	case *rule.NotRule:
		castRule.SetRule(g.mapAnonymousRules(lexerRulesMap, castRule.Rule()))
	case *rule.TestRule:
		castRule.SetRule(g.mapAnonymousRules(lexerRulesMap, castRule.Rule()))
	case *rule.OptionalRule:
		castRule.SetRule(g.mapAnonymousRules(lexerRulesMap, castRule.Rule()))
	case *rule.ZeroOrMoreRule:
		castRule.SetRule(g.mapAnonymousRules(lexerRulesMap, castRule.Rule()))
	case *rule.OneOrMoreRule:
		castRule.SetRule(g.mapAnonymousRules(lexerRulesMap, castRule.Rule()))
	case *rule.CharRule:
		return g.mapCharRule(lexerRulesMap, castRule)
	case *rule.RangeRule:
		return g.mapRangeRule(lexerRulesMap, castRule)
	case *rule.StringRule:
		return g.mapStringRule(lexerRulesMap, castRule)
	default:
		// do nothing
	}
	return r
}

func (g *Grammar) mapStringRule(lexerRulesMap map[string]rule.Rule, r *rule.StringRule) rule.Rule {
	mappedRule, found := lexerRulesMap[r.String()]
	if found {
		return mappedRule
	}
	ruleName := stringRuleName(lexerRulesMap, r.Text(), r.CaseSensitive())
	newRule := rule.New(ruleName, r)
	lexerRulesMap[r.String()] = newRule
	g.lexerRules.Add(newRule)
	return newRule
}

func charName(char rune) string {
	if name, found := namedChars[char]; found {
		return name
	}
	return strconv.FormatInt(int64(char), 16)
}

func stringRuleName(lexerRulesMap map[string]rule.Rule, str string, caseSensitive bool) string {
	if len(str) == 1 {
		return charRuleName(lexerRulesMap, rune(str[0]), caseSensitive)
	}
	ruleName := strings.Builder{}
	if caseSensitive {
		ruleName.WriteString("str_")
	} else {
		ruleName.WriteString("stri_")
	}
	for _, char := range str {
		if replacement, found := replacements[char]; found {
			ruleName.WriteRune(replacement)
		} else if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			ruleName.WriteRune(char)
		} else {
			ruleName.WriteString(charName(char))
		}
	}
	if _, found := lexerRulesMap[ruleName.String()]; found {
		ruleName.WriteRune('_')
		ruleName.WriteString(strconv.FormatInt(int64(len(lexerRulesMap)), 10))
	}
	return ruleName.String()
}

func (g *Grammar) mapRangeRule(lexerRulesMap map[string]rule.Rule, r *rule.RangeRule) rule.Rule {
	mappedRule, found := lexerRulesMap[r.String()]
	if found {
		return mappedRule
	}
	ruleName := rangeRuleName(lexerRulesMap, r.Start(), r.End())
	newRule := rule.New(ruleName, r)
	lexerRulesMap[r.String()] = newRule
	g.lexerRules.Add(newRule)
	return newRule
}

func rangeRuleName(lexerRulesMap map[string]rule.Rule, start, end rune) string {
	ruleName := strings.Builder{}
	ruleName.WriteString("range_")
	if replacement, found := replacements[start]; found {
		ruleName.WriteRune(replacement)
	} else if (start >= 'A' && start <= 'Z') || (start >= 'a' && start <= 'z') || (start >= '0' && start <= '9') {
		ruleName.WriteRune(start)
	} else {
		ruleName.WriteString(charName(start))
	}
	ruleName.WriteRune('_')
	if replacement, found := replacements[end]; found {
		ruleName.WriteRune(replacement)
	} else if (end >= 'A' && end <= 'Z') || (end >= 'a' && end <= 'z') || (end >= '0' && end <= '9') {
		ruleName.WriteRune(end)
	} else {
		ruleName.WriteString(charName(end))
	}
	if _, found := lexerRulesMap[ruleName.String()]; found {
		ruleName.WriteRune('_')
		ruleName.WriteString(strconv.FormatInt(int64(len(lexerRulesMap)), 10))
	}
	return ruleName.String()
}

func (g *Grammar) mapCharRule(lexerRulesMap map[string]rule.Rule, r *rule.CharRule) rule.Rule {
	mappedRule, found := lexerRulesMap[r.String()]
	if found {
		return mappedRule
	}
	ruleName := charRuleName(lexerRulesMap, r.Char(), r.CaseSensitive())
	if _, found = lexerRulesMap[ruleName]; found {
		ruleName = ruleName + "_" + strconv.FormatInt(int64(len(lexerRulesMap)), 10)
	}
	newRule := rule.New(ruleName, r)
	lexerRulesMap[r.String()] = newRule
	g.lexerRules.Add(newRule)
	return newRule
}

func charRuleName(lexerRulesMap map[string]rule.Rule, char rune, caseSensitive bool) string {
	ruleName := strings.Builder{}
	if char == rule.EOI.Char() {
		return "EOI"
	}
	if caseSensitive {
		ruleName.WriteString("chr_")
	} else {
		ruleName.WriteString("chri_")
	}
	if replacement, found := replacements[char]; found {
		ruleName.WriteRune(replacement)
	} else if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
		ruleName.WriteRune(char)
	} else {
		ruleName.WriteString(charName(char))
	}
	if _, found := lexerRulesMap[ruleName.String()]; found {
		ruleName.WriteRune('_')
		ruleName.WriteString(strconv.FormatInt(int64(len(lexerRulesMap)), 10))
	}
	return ruleName.String()
}

func (g *Grammar) HasError() bool {
	return len(g.Errors()) > 0
}

func (g *Grammar) Errors() []error {
	if g.errors == nil {
		g.mapRules()
	}
	return g.errors
}

func (g *Grammar) LexerRules() []*rule.NonTerminalRule {
	if g.lexerRules == nil {
		g.mapRules()
	}
	return g.lexerRules.Items()
}

func (g *Grammar) ParserRules() []*rule.NonTerminalRule {
	if g.parserRules == nil {
		g.mapRules()
	}
	return g.parserRules.Items()
}

func (g *Grammar) ToText(writer util.TextWriter) {
	lexerRules := g.LexerRules()
	writer.WriteString("/* ").WriteString(util.PadC(g.name+" grammar", 76, ' ')).WriteString("*/").NewLine().NewLine()
	if len(lexerRules) > 0 {
		writer.WriteString("/******************************************************************************/").NewLine()
		writer.WriteString("/*                                LEXER RULES                                 */").NewLine()
		writer.WriteString("/******************************************************************************/").NewLine()
		for _, r := range lexerRules {
			writer.WriteString(r.Id()).WriteString(": ")
			r.Rule().ToText(writer)
			writer.WriteRune(';').NewLine()
		}
	}
	parserRules := g.ParserRules()
	if len(parserRules) > 0 {
		writer.WriteString("/******************************************************************************/").NewLine()
		writer.WriteString("/*                                PARSER RULES                                */").NewLine()
		writer.WriteString("/******************************************************************************/").NewLine()
		for _, r := range parserRules {
			writer.WriteString(r.Id()).WriteString(": ")
			r.Rule().ToText(writer)
			writer.WriteRune(';').NewLine()
		}
	}
}

func (g *Grammar) GetRule(name string) *rule.NonTerminalRule {
	return g.rules[name]
}

func (g *Grammar) Validate() error {
	for _, r := range g.rules {
		if r.Rule() == nil {
			return fmt.Errorf("rule '%s' not defined", r.Id())
		}
	}
	return nil
}
