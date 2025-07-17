package parser

import (
	"github.com/fabiouggeri/page/runtime/ast"
	"github.com/fabiouggeri/page/runtime/error"
	"github.com/fabiouggeri/page/runtime/lexer"
)

type Parser struct {
	lexer       *lexer.Lexer
	syntax      *Syntax
	currentNode *ast.Node
	errors      []error.Error
	ignore      bool
}

func New(l *lexer.Lexer, s *Syntax) *Parser {
	return &Parser{
		lexer:  l,
		syntax: s,
		errors: make([]error.Error, 0),
		ignore: false,
	}
}

func (p *Parser) Execute() *ast.Node {
	startRule := p.syntax.StartRule()
	if startRule < 0 || startRule >= p.syntax.RulesCount() {
		panic("undefined start rule")
	}
	p.currentNode = ast.NewNode(-1, 0, 0)
	if p.parseRule(startRule) {
		return p.currentNode
	}
	return nil
}

func (p *Parser) Errors() []error.Error {
	return p.errors
}

func (p *Parser) parseRule(ruleId int) bool {
	previousIgnore := p.ignore
	if p.syntax.HasOption(ruleId, IGNORE) {
		p.ignore = true
	}
	//fmt.Printf("Rule (%d, %d): %s\n", p.lexer.Row(), p.lexer.Col(), p.syntax.rulesNames[ruleId])
	lastNode := p.currentNode
	index := p.lexer.Index()
	match := false
	//row, col := p.lexer.Row(), p.lexer.Col()
	rules := p.syntax.Subrules(ruleId)
	switch ParserRuleType(rules[0]) {
	case AND_RULE:
		match = p.parseAndRule(rules)
	case OR_RULE:
		match = p.parseOrRule(rules)
	case ONE_OR_MORE_RULE:
		match = p.parseOneOrMoreRule(rules)
	case ZERO_OR_MORE_RULE:
		match = p.parseZeroOrMoreRule(rules)
	case OPTIONAL_RULE:
		match = p.parseOptionalRule(rules)
	case TEST_NOT_RULE:
		match = p.parseTestNotRule(rules)
	case TEST_RULE:
		match = p.parseTestRule(rules)
	case TERMINAL_RULE:
		match = p.parseTerminalRule(rules)
	case NON_TERMINAL_RULE:
		match = p.parseNonTerminalRule(rules)
	default:
		panic("undefined rule type")
	}
	if match && !p.syntax.IsSubRule(ruleId) && !p.syntax.HasOption(ruleId, SKIP_NODE) && !p.ignore {
		p.createNode(ruleId, index, lastNode)
	}
	p.ignore = previousIgnore
	return match
}

func (p *Parser) createNode(ruleId int, index int, lastNode *ast.Node) {
	startTkn, _ := p.lexer.Token(index)
	endTkn, _ := p.lexer.Token(p.lexer.Index() - 1)
	p.currentNode = ast.NewNode(ruleId, startTkn.Index(), endTkn.Index()+endTkn.Len())
	p.currentNode.SetFirstChild(lastNode.Sibling())
	lastNode.SetSibling(p.currentNode)
}

func (p *Parser) parseAndRule(rules []int) bool {
	index := p.lexer.Index()
	for _, sub := range rules[1:] {
		if !p.parseRule(sub) {
			p.lexer.SetIndex(index)
			return false
		}
	}
	return true
}

func (p *Parser) parseOrRule(rules []int) bool {
	index := p.lexer.Index()
	for _, sub := range rules[1:] {
		if p.parseRule(sub) {
			return true
		}
		p.lexer.SetIndex(index)
	}
	return false
}

func (p *Parser) parseOneOrMoreRule(rules []int) bool {
	index := p.lexer.Index()
	if !p.parseRule(rules[1]) {
		p.lexer.SetIndex(index)
		return false
	}
	index = p.lexer.Index()
	for p.parseRule(rules[1]) {
		index = p.lexer.Index()
	}
	p.lexer.SetIndex(index)
	return true
}

func (p *Parser) parseZeroOrMoreRule(rules []int) bool {
	index := p.lexer.Index()
	for p.parseRule(rules[1]) {
		index = p.lexer.Index()
	}
	p.lexer.SetIndex(index)
	return true
}

func (p *Parser) parseOptionalRule(rules []int) bool {
	index := p.lexer.Index()
	if !p.parseRule(rules[1]) {
		p.lexer.SetIndex(index)
	}
	return true
}

func (p *Parser) parseTestRule(rules []int) bool {
	index := p.lexer.Index()
	if p.parseRule(rules[1]) {
		p.lexer.SetIndex(index)
		return true
	}
	p.lexer.SetIndex(index)
	return false
}

func (p *Parser) parseTestNotRule(rules []int) bool {
	index := p.lexer.Index()
	if p.parseRule(rules[1]) {
		p.lexer.SetIndex(index)
		return false
	}
	p.lexer.SetIndex(index)
	return true
}

func (p *Parser) parseNonTerminalRule(rules []int) bool {
	index := p.lexer.Index()
	if p.parseRule(rules[1]) {
		return true
	}
	p.lexer.SetIndex(index)
	return false
}

func (p *Parser) parseTerminalRule(rules []int) bool {
	index := p.lexer.Index()
	tkn, err := p.lexer.NextToken()
	if err != nil {
		p.LexError(err)
		p.lexer.SetIndex(index)
		return false
	}
	for p.lexer.IsIgnored(tkn) {
		tkn, err = p.lexer.NextToken()
		if err != nil {
			p.LexError(err)
			p.lexer.SetIndex(index)
			return false
		}
	}
	if tkn.IsType(rules[1]) {
		return true
	}
	p.lexer.SetIndex(index)
	return false
}

func (p *Parser) LexError(lexError error.Error) {
	p.errors = append(p.errors, lexError)
}

func (p *Parser) Error(errorCode int, row int, col int, message string) {
	err := &ParserError{
		code:    errorCode,
		row:     row,
		col:     col,
		message: message,
	}
	p.errors = append(p.errors, err)
}
