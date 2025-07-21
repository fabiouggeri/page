package visitor

import "github.com/fabiouggeri/page/runtime/parser"

type ASTWalker struct {
	parser  *parser.Parser
	visitor *RuleVisitor
}

func NewWalker(parser *parser.Parser, visitor *RuleVisitor) *ASTWalker {
	return &ASTWalker{
		parser:  parser,
		visitor: visitor,
	}
}

func (w *ASTWalker) Walk(node *parser.ASTNode) {
	if node == nil {
		return
	}

	ruleId := node.RuleType()
	if callback, exists := w.visitor.enterRuleCallbacks[ruleId]; exists {
		callback(w.parser, node)
	}

	child := node.FirstChild()
	for child != nil {
		w.Walk(child)
		child = child.Sibling()
	}

	if callback, exists := w.visitor.exitRuleCallbacks[ruleId]; exists {
		callback(w.parser, node)
	}
}
