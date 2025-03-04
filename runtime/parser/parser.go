package parser

import (
	"github.com/fabiouggeri/page/runtime/ast"
	"github.com/fabiouggeri/page/runtime/lexer"
)

type Parser struct {
	lexer  *lexer.Lexer
	syntax *Syntax
}

func New(l *lexer.Lexer, s *Syntax) *Parser {
	return &Parser{
		lexer:  l,
		syntax: s,
	}
}

func (p *Parser) Execute() *ast.Node {
	return nil
}
