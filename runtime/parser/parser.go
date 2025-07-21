package parser

import (
	"github.com/fabiouggeri/page/runtime/error"
	"github.com/fabiouggeri/page/runtime/lexer"
)

type memorizedRule struct {
	node  *ASTNode
	start int
	end   int
}

type Parser struct {
	lexer       *lexer.Lexer
	syntax      *Syntax
	currentNode *ASTNode
	errors      []error.Error
	memorized   []*memorizedRule
	ignore      bool
}

func New(l *lexer.Lexer, s *Syntax) *Parser {
	return &Parser{
		lexer:     l,
		syntax:    s,
		errors:    make([]error.Error, 0),
		memorized: make([]*memorizedRule, len(s.rulesNames)),
		ignore:    false,
	}
}

func (p *Parser) Lexer() *lexer.Lexer {
	return p.lexer
}

func (p *Parser) Syntax() *Syntax {
	return p.syntax
}

func (p *Parser) Execute() *ASTNode {
	startRule := p.syntax.StartRule()
	if startRule < 0 || startRule >= p.syntax.RulesCount() {
		panic("undefined start rule")
	}
	p.currentNode = NewASTNode(-1, 0, 0)
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
	mem := p.memorized[ruleId]
	if mem != nil && mem.start == index {
		if mem.start <= mem.end {
			p.lexer.SetIndex(mem.end)
			return true
		} else {
			return false
		}
	}
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
	if match && !p.ignore && !p.syntax.IsSubRule(ruleId) && !p.syntax.HasOption(ruleId, SKIP_NODE) {
		p.createNode(ruleId, index, lastNode)
	} else if mem != nil {
		mem.start = index
		mem.end = -1
		mem.node = nil
	}
	p.ignore = previousIgnore
	return match
}

func (p *Parser) createNode(ruleId int, index int, lastNode *ASTNode) {
	p.currentNode = NewASTNode(ruleId, index, p.lexer.Index()-1)
	p.currentNode.SetFirstChild(lastNode.Sibling())
	lastNode.SetSibling(p.currentNode)
	if p.memorized[ruleId] == nil {
		p.memorized[ruleId] = &memorizedRule{
			node:  p.currentNode,
			start: index,
			end:   p.lexer.Index(),
		}
	} else {
		p.memorized[ruleId].node = p.currentNode
		p.memorized[ruleId].start = index
		p.memorized[ruleId].end = p.lexer.Index()
	}
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

func (p *Parser) NodeTokens(node *ASTNode) (*lexer.Token, *lexer.Token) {
	startToken, _ := p.lexer.Token(node.StartToken())
	endToken, _ := p.lexer.Token(node.EndToken())
	return startToken, endToken
}

func (p *Parser) NodeText(node *ASTNode) string {
	startToken, endToken := p.NodeTokens(node)
	if startToken == nil || endToken == nil {
		return ""
	}
	return p.lexer.Input().GetText(startToken.Index(), endToken.Index()+endToken.Len())
}

func (p *Parser) Position(node *ASTNode) (int, int) {
	token, _ := p.lexer.Token(node.StartToken())
	if token == nil {
		return 0, 0
	}
	return token.Row(), token.Col()
}
