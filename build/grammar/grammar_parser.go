package grammar

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fabiouggeri/page/build/rule"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type grammarParser struct {
	grammar          *Grammar
	index            int
	buffer           []byte
	line             uint32
	col              uint32
	explicitMainRule bool
	options          map[*rule.RuleOption]string
}

type GrammarError struct {
	Line    uint32
	Col     uint32
	Message string
}

const (
	NO_COMMENT    = 0
	LINE_COMMENT  = 1
	BLOCK_COMMENT = 2
)

func newParser(grammar *Grammar, content []byte) *grammarParser {
	return &grammarParser{
		grammar:          grammar,
		index:            0,
		col:              1,
		line:             1,
		explicitMainRule: false,
		buffer:           content,
		options:          make(map[*rule.RuleOption]string),
	}
}

func parseGrammar(content []byte) (*Grammar, error) {
	grammar := New("")
	parser := newParser(grammar, content)
	if err := parser.parse(false); err != nil {
		return nil, err
	}
	return grammar, nil
}

func normalizePathname(fileName string) string {
	var normalPathname strings.Builder
	withExtension := false
	for i := 0; i < len(fileName); i++ {
		c := fileName[i]
		switch c {
		case '/', '\\':
			normalPathname.WriteRune(os.PathSeparator)
		case '.':
			withExtension = true
		default:
			normalPathname.WriteByte(c)
		}
	}
	if !withExtension {
		normalPathname.WriteString(".gy")
	}
	return normalPathname.String()
}

func (p *grammarParser) importGrammar(importPathname string) error {
	var buffer []byte
	grammarPathname := normalizePathname(importPathname)

	file, err := os.Open(grammarPathname)
	if err != nil {
		return err
	}
	if p.grammar.encode != nil {
		reader := transform.NewReader(file, p.grammar.encode.NewDecoder())
		buffer, err = io.ReadAll(reader)
	} else {
		buffer, err = io.ReadAll(file)
	}
	if err != nil {
		return err
	}
	return newParser(p.grammar, buffer).parse(true)
}

func (e GrammarError) Error() string {
	return e.Message
}

func (l *grammarParser) containsOption(option *rule.RuleOption) bool {
	_, found := l.options[option]
	return found
}

func (l *grammarParser) addOption(option *rule.RuleOption, value string) {
	l.options[option] = value
}

func (l *grammarParser) clearOptions() {
	l.options = make(map[*rule.RuleOption]string)
}

func (l *grammarParser) parse(importing bool) error {
	for l.hasNext() {
		err := l.grammarEntry(importing)
		if err != nil {
			return err
		}
		l.skipSpaces()
	}
	return l.grammar.Validate()
}

func (l *grammarParser) hasNext() bool {
	return l.index < len(l.buffer)
}

func (l *grammarParser) error(msg string, args ...any) error {
	return GrammarError{Line: l.line, Col: l.col, Message: fmt.Sprintf("%d, %d: %s", l.line, l.col, fmt.Sprintf(msg, args...))}
}

func (l *grammarParser) grammarEntry(importing bool) error {
	var err error
	l.skipSpaces()
	char := l.currentChar()
	if unicode.IsLetter(char) {
		identifier := l.consumeIdentifier()
		if identifier == "grammar" {
			err = l.grammarNameEntry(importing)
		} else if identifier == "import" {
			err = l.importGrammarEntry()
		} else if identifier == "charset" {
			err = l.charsetEntry(importing)
		} else {
			err = l.nonTerminalEntry(identifier, importing)
		}
	} else if char == '@' {
		err = l.optionEntry()
	} else {
		err = l.error("unknown character found.")
	}
	return err
}

func (l *grammarParser) grammarNameEntry(importing bool) error {
	if l.grammar.name == "" || importing {
		l.skipSpaces()
		if unicode.IsLetter(l.currentChar()) {
			name := l.consumeIdentifier()
			if !importing {
				l.grammar.name = name
			}
			l.skipSpaces()
			if l.currentChar() == ';' {
				l.advanceIndex()
			} else {
				return l.error("; not found after grammar name!")
			}
		} else {
			return l.error("Grammar name not found!")
		}
	} else {
		return l.error("Grammar name already defined!")
	}
	return nil
}

func (l *grammarParser) importGrammarEntry() error {
	var err error
	l.skipSpaces()
	if unicode.IsLetter(l.currentChar()) {
		importName := l.consumeUp(';')
		l.importGrammar(importName)
		if l.currentChar() == ';' {
			l.advanceIndex()
		} else {
			err = l.error("; not found after grammar name!")
		}
	} else {
		err = l.error("Grammar name not found!")
	}
	return err
}

func (l *grammarParser) findEncoder(charset string) (*charmap.Charmap, error) {
	for _, enc := range charmap.All {
		cmap, ok := enc.(*charmap.Charmap)
		if ok && cmap.String() == charset {
			return cmap, nil
		}
	}
	return nil, l.error("Charset encoder not found")
}

func (l *grammarParser) charsetEntry(importing bool) error {
	var err error
	if l.grammar.encode == nil || importing {
		l.skipSpaces()
		if unicode.IsLetter(l.currentChar()) {
			charset := l.consumeIdentifier()
			if !importing {
				var enc *charmap.Charmap
				enc, err = l.findEncoder(charset)
				if err != nil {
					return err
				}
				l.grammar.encode = enc
			}
			l.skipSpaces()
			if l.currentChar() == ';' {
				l.advanceIndex()
			} else {
				err = l.error("; not found after charset!")
			}
		} else {
			err = l.error("Charset not found!")
		}
	} else {
		err = l.error("Charset already defined!")
	}
	return err
}

func (l *grammarParser) nonTerminalEntry(ruleName string, importing bool) error {

	if ruleName == "EOI" {
		return l.error("EOI is a reserved rule name.")
	}

	currentRule := l.grammar.GetRule(ruleName)

	if currentRule == nil {
		currentRule = rule.New(ruleName, nil)
		l.grammar.Rules(currentRule)
	} else if currentRule.Rule() != nil {
		fmt.Printf("Rule %s redefined in line %d.", currentRule.Id(), l.line)
	}
	l.skipSpaces()
	if l.currentChar() == ':' {
		var execRule rule.Rule
		var err error
		l.advanceIndex()
		execRule, err = l.orRule()
		if execRule != nil {
			currentRule.SetRule(execRule)
			for option, value := range l.options {
				currentRule.Option(option, value)
			}
			if !importing {
				if l.containsOption(rule.MAIN) {
					if !l.explicitMainRule {
						l.grammar.mainRule = currentRule
						l.explicitMainRule = true
					} else {
						return l.error("Rule '%s' is already defined as main rule.", l.grammar.mainRule.Id())
					}
				} else if l.grammar.mainRule == nil {
					l.grammar.mainRule = currentRule
				}
			}
			l.clearOptions()
			l.skipSpaces()
			if l.currentChar() == ';' {
				l.advanceIndex()
			} else {
				return l.error("; not found after rule definition!")
			}
		} else if err != nil {
			return err
		} else {
			return l.error("unknown %c found!", l.currentChar())
		}
	} else {
		return l.error(": not found after rule name.")
	}
	return nil
}

func (l *grammarParser) orRule() (rule.Rule, error) {
	rules := make([]rule.Rule, 0)
	currentRule, err := l.andRule()
	for currentRule != nil {
		rules = append(rules, currentRule)
		l.skipSpaces()
		c := l.currentChar()
		if c == '|' {
			l.advanceIndex()
			currentRule, err = l.andRule()
		} else {
			currentRule = nil
		}
	}

	if len(rules) > 1 {
		currentRule = rule.Or(rules...)
	} else if len(rules) == 1 {
		currentRule = rules[0]
	}
	return currentRule, err
}

func (l *grammarParser) andRule() (rule.Rule, error) {
	rules := make([]rule.Rule, 0)
	currentRule, err := l.postFixedRule()
	for currentRule != nil {
		rules = append(rules, currentRule)
		currentRule, err = l.postFixedRule()
	}
	if len(rules) > 1 {
		currentRule = rule.And(rules...)
	} else if len(rules) == 1 {
		currentRule = rules[0]
	}
	return currentRule, err
}

func (l *grammarParser) postFixedRule() (rule.Rule, error) {
	currentRule, err := l.simpleRule()
	if err == nil && currentRule != nil {
		for {
			l.skipSpaces()
			c := l.currentChar()
			switch c {
			case '?':
				currentRule = rule.Optional(currentRule)
				l.advanceIndex()
			case '+':
				currentRule = rule.OneOrMore(currentRule)
				l.advanceIndex()
			case '*':
				currentRule = rule.ZeroOrMore(currentRule)
				l.advanceIndex()
			case '&':
				currentRule = rule.Test(currentRule)
				l.advanceIndex()
			case '!':
				currentRule = rule.Not(currentRule)
				l.advanceIndex()
			default:
				return currentRule, err
			}
		}
	}
	return currentRule, err
}

func (l *grammarParser) simpleRule() (rule.Rule, error) {
	var currentRule rule.Rule
	var err error
	l.skipSpaces()
	c := l.currentChar()
	switch c {
	case '"':
		currentRule, err = l.ignoreCaseLiteralRule()
	case '\'':
		currentRule, err = l.literalRule()
	case '(':
		currentRule, err = l.groupedRule()
	case '[':
		currentRule, err = l.charRangeRule()
	// case '.':
	// 	currentRule = rule.NewAnyCharRule()
	// 	l.advanceIndex()
	default:
		if unicode.IsLetter(c) {
			currentRule = l.identifierRule()
		}
	}
	return currentRule, err
}

func (l *grammarParser) appendEscapedChar(literal *strings.Builder, c rune) {
	switch c {
	case 'n':
		literal.WriteByte('\n')
	case 'r':
		literal.WriteByte('\r')
	case 't':
		literal.WriteByte('\t')
	case 'b':
		literal.WriteByte('\b')
	case 'f':
		literal.WriteByte('\f')
	case 'u':
		literal.WriteString("\\u")
		literal.WriteString(string(l.buffer[l.index : l.index+4]))
		l.advanceIndex()
		l.advanceIndex()
		l.advanceIndex()
		l.advanceIndex()
	default:
		literal.WriteRune(c)
	}
}

func (l *grammarParser) ignoreCaseLiteralRule() (rule.Rule, error) {
	var literal strings.Builder
	var currentRule rule.Rule
	var err error

	lastChar := rune(0)
	l.advanceIndex()
	for l.hasNext() {
		c := l.currentChar()
		l.advanceIndex()
		if lastChar == '\\' {
			l.appendEscapedChar(&literal, c)
			lastChar = rune(0)
		} else if c == '"' {
			break
		} else if c == '\\' {
			lastChar = c
		} else {
			literal.WriteRune(c)
		}
	}
	l.skipSpaces()
	if literal.Len() > 0 {
		if l.currentChar() == ':' {
			l.advanceIndex()
			l.skipSpaces()
			if unicode.IsDigit(l.currentChar()) {
				number := l.consumeNumber()
				len, _ := strconv.Atoi(number)
				if len < literal.Len() {
					currentRule = rule.StringPartialI(literal.String(), int32(len))
				} else {
					err = l.error("Partial match length must be smaller than literal length!")
				}
			} else {
				err = l.error("Literal partial match value not found!")
			}
		} else if literal.Len() > 1 {
			buf := literal.String()
			if buf[0] == '\\' && buf[1] == 'u' {
				runeCode, runeErr := strconv.Atoi(buf[2:])
				if runeErr != nil {
					return nil, l.error("Invalid rune code: %s", buf)
				}
				currentRule = rule.CharI(rune(runeCode))
			} else {
				currentRule = rule.StringI(literal.String())
			}
		} else {
			buf := []rune(literal.String())
			currentRule = rule.CharI(buf[0])
		}
	} else {
		err = l.error("Found empty literal!")
	}
	return currentRule, err
}

func (l *grammarParser) literalRule() (rule.Rule, error) {
	var literal strings.Builder
	var currentRule rule.Rule
	var err error

	lastChar := rune(0)
	l.advanceIndex()
	for l.hasNext() {
		c := l.currentChar()
		l.advanceIndex()
		if lastChar == '\\' {
			l.appendEscapedChar(&literal, c)
			lastChar = rune(0)
		} else if c == '\'' {
			break
		} else if c == '\\' {
			lastChar = c
		} else {
			literal.WriteRune(c)
		}
	}
	l.skipSpaces()
	if literal.Len() > 0 {
		if l.currentChar() == ':' {
			l.advanceIndex()
			l.skipSpaces()
			if unicode.IsDigit(l.currentChar()) {
				number := l.consumeNumber()
				if literal.Len() > 1 {
					len, _ := strconv.Atoi(number)
					currentRule = rule.StringPartial(literal.String(), int32(len))
				} else {
					err = l.error("Partial match not allowed in literal of length one!")
				}
			} else {
				err = l.error("Literal partial match value not found!")
			}
		} else if literal.Len() > 1 {
			buf := literal.String()
			if buf[0] == '\\' && buf[1] == 'u' {
				runeCode, runeErr := strconv.Atoi(buf[2:])
				if runeErr != nil {
					return nil, l.error("Invalid rune code: %s", buf)
				}
				currentRule = rule.Char(rune(runeCode))
			} else {
				currentRule = rule.String(literal.String())
			}
		} else {
			buf := []rune(literal.String())
			currentRule = rule.Char(buf[0])
		}
	} else {
		err = l.error("Found empty literal!")
	}
	return currentRule, err
}

func (l *grammarParser) groupedRule() (rule.Rule, error) {
	l.advanceIndex()
	currentRule, err := l.orRule()
	if err != nil {
		return nil, err
	}
	l.skipSpaces()
	if l.currentChar() != ')' {
		return nil, l.error("Closing parenthesis not found.")
	}
	l.advanceIndex()
	return currentRule, nil
}

func isHexStart(firstChar, secondChar rune) bool {
	return firstChar == '0' && (secondChar == 'x' || secondChar == 'X')
}

func (l *grammarParser) consumeHexValue() (rune, error) {
	var sb strings.Builder
	currChar := l.currentChar()
	for (currChar >= '0' && currChar <= '9') ||
		(currChar >= 'a' && currChar <= 'f') ||
		(currChar >= 'A' && currChar <= 'F') {
		sb.WriteRune(currChar)
		l.advanceIndex()
		currChar = l.currentChar()
	}
	val, err := strconv.Atoi(sb.String())
	return rune(val), err
}

func (l *grammarParser) charRangeRule() (rule.Rule, error) {
	var err error
	l.advanceIndex()
	l.skipSpaces()
	startChar := l.currentChar()
	l.advanceIndex()
	if isHexStart(startChar, l.currentChar()) {
		l.advanceIndex()
		startChar, err = l.consumeHexValue()
		if err != nil {
			return nil, err
		}
	}
	l.skipSpaces()
	if l.currentChar() != '-' {
		return nil, l.error("Character - not found between start and end character in range rule.")
	}
	l.advanceIndex()
	l.skipSpaces()
	endChar := l.currentChar()
	l.advanceIndex()
	if isHexStart(endChar, l.currentChar()) {
		l.advanceIndex()
		endChar, err = l.consumeHexValue()
		if err != nil {
			return nil, err
		}
	}
	l.skipSpaces()
	if l.currentChar() != ']' {
		return nil, l.error("Closing brackets not found.")
	}
	l.advanceIndex()
	return rule.Range(startChar, endChar), nil
}

func (l *grammarParser) identifierRule() rule.Rule {
	var currentRule rule.Rule
	id := l.consumeIdentifier()
	if id == "EOI" {
		currentRule = rule.EOI
	} else {
		nonTermRule := l.grammar.GetRule(id)
		if nonTermRule == nil {
			nonTermRule = rule.New(id, nil)
			l.grammar.Rules(nonTermRule)
		}
		currentRule = nonTermRule
	}
	return currentRule
}

func (l *grammarParser) consumeNumber() string {
	var number strings.Builder

	l.skipSpaces()
	number.WriteRune(l.currentChar())
	l.advanceIndex()
	for l.hasNext() {
		c := l.currentChar()
		if unicode.IsDigit(c) {
			number.WriteRune(c)
			l.advanceIndex()
		} else {
			break
		}
	}
	return number.String()
}

func (l *grammarParser) optionEntry() error {
	foundOption := false
	l.advanceIndex()
	if unicode.IsLetter(l.currentChar()) {
		identifier := l.consumeIdentifier()
		for _, option := range rule.AllOptions {
			if option.Name() == identifier {
				if option.Parameterized() {
					l.skipSpaces()
					if l.currentChar() == '(' {
						l.advanceIndex()
						optionValue := l.consumeUp(')')
						l.skipSpaces()
						if l.currentChar() == ')' {
							l.advanceIndex()
							l.addOption(option, optionValue)
							foundOption = true
							break
						} else {
							return l.error("expected ) not found on option parameter.")
						}
					} else if !option.ParameterMandatory() {
						l.addOption(option, "")
						foundOption = true
						break
					} else {
						return l.error("Expected option parameter not found.")
					}
				} else {
					l.addOption(option, "")
					foundOption = true
					break
				}
			}
		}
		if !foundOption {
			return l.error("Unknown option specified.")
		}
	} else {
		return l.error("Invalid character after option marker.")
	}
	return nil
}

func (l *grammarParser) currentChar() rune {
	r := rune(l.buffer[l.index])
	if r >= utf8.RuneSelf {
		r, _ = utf8.DecodeRune(l.buffer[l.index:])
	}
	return r
}

func (l *grammarParser) charAt(index int) rune {
	r := rune(l.buffer[index])
	if r >= utf8.RuneSelf {
		r, _ = utf8.DecodeRune(l.buffer[index:])
	}
	return r
}

func (l *grammarParser) skipSpaces() {
	state := NO_COMMENT
	for l.index < len(l.buffer) {
		char := l.currentChar()
		switch char {
		case '\t', '\v', '\f', ' ', 0x85, 0xA0:
			l.col++
		case '\r':
			// ignore
		case '\n':
			if state == LINE_COMMENT {
				state = NO_COMMENT
			}
			l.line++
			l.col = 1
		case '*':
			if state == BLOCK_COMMENT {
				if l.index+1 < len(l.buffer) && l.charAt(l.index+1) == '/' {
					l.index++
					l.col++
					state = NO_COMMENT
				}
				l.col++
			} else if state == LINE_COMMENT {
				l.col++
			} else {
				return
			}
		case '/':
			if l.index+1 < len(l.buffer) {
				nextChar := l.charAt(l.index + 1)
				if nextChar == '*' {
					state = BLOCK_COMMENT
					l.index++
					l.col++
				} else if nextChar == '/' {
					state = LINE_COMMENT
					l.index++
					l.col++
				} else {
					return
				}
			}
		default:
			if state == NO_COMMENT {
				return
			} else {
				l.col++
			}
		}
		l.index++
	}
}

func (l *grammarParser) advanceIndex() {
	r, size := rune(l.buffer[l.index]), 1
	if r >= utf8.RuneSelf {
		_, size = utf8.DecodeRune(l.buffer[l.index:])
	}
	l.index += size
	l.col++
}

func (l *grammarParser) consumeUp(endChar rune) string {
	var text strings.Builder
	lastChar := rune(0)

	l.skipSpaces()
	for l.hasNext() {
		char := l.currentChar()
		if lastChar == '\\' {
			text.WriteRune(char)
			lastChar = rune(0)
		} else {
			if char == endChar {
				break
			} else if char != '\\' {
				text.WriteRune(char)
			}
			lastChar = char
		}
		l.advanceIndex()
	}
	return text.String()
}

func (l *grammarParser) consumeIdentifier() string {
	start := l.index
	l.advanceIndex()
	for l.hasNext() {
		char := l.currentChar()
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' || char == '$' {
			l.advanceIndex()
		} else {
			break
		}
	}
	return string(l.buffer[start:l.index])
}
