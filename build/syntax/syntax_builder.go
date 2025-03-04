package syntax

import (
	"github.com/fabiouggeri/page/build/grammar"
	"github.com/fabiouggeri/page/runtime/parser"
)

func FromGrammar(g *grammar.Grammar) *parser.Syntax {
	return nil // FromDFA(dfa.NFAToDFA(dfa.RulesToNFA(g.ParserRules()...)))
}
