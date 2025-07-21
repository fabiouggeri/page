package parser

import "strings"

type ASTNode struct {
	ruleType   int
	startToken int
	endToken   int
	sibling    *ASTNode
	firstChild *ASTNode
}

func NewASTNode(ruleType int, start int, end int) *ASTNode {
	return &ASTNode{
		ruleType:   ruleType,
		startToken: start,
		endToken:   end,
	}
}

func (n *ASTNode) RuleType() int {
	return n.ruleType
}

func (n *ASTNode) StartToken() int {
	return n.startToken
}

func (n *ASTNode) EndToken() int {
	return n.endToken
}

func (n *ASTNode) Sibling() *ASTNode {
	return n.sibling
}

func (n *ASTNode) SetSibling(sibling *ASTNode) {
	n.sibling = sibling
}

func (n *ASTNode) FirstChild() *ASTNode {
	return n.firstChild
}

func (n *ASTNode) SetFirstChild(firstChild *ASTNode) {
	n.firstChild = firstChild
}

func (n *ASTNode) Find(syntax *Syntax, pathToNode string) *ASTNode {
	return n.findNode(syntax, strings.Split(pathToNode, "/"))
}

func (n *ASTNode) findNode(syntax *Syntax, rulesNames []string) *ASTNode {
	if len(rulesNames) == 0 {
		return nil
	}
	child := n.firstChild
	for child != nil {
		if strings.EqualFold(syntax.RuleName(child.ruleType), rulesNames[0]) {
			if len(rulesNames) == 1 {
				return child
			} else {
				return child.findNode(syntax, rulesNames[1:])
			}
		}
		child = child.sibling
	}
	return nil
}

func (n *ASTNode) List(syntax *Syntax, pathToNode string) []*ASTNode {
	rulesNames := strings.Split(pathToNode, "/")
	if len(rulesNames) == 0 {
		return []*ASTNode{}
	}
	var startNode *ASTNode
	var ruleName string
	if len(rulesNames) > 1 {
		startNode = n.findNode(syntax, rulesNames[:len(rulesNames)-1])
		if startNode == nil {
			return []*ASTNode{}
		}
		ruleName = rulesNames[len(rulesNames)-1]
	} else {
		startNode = n
		ruleName = rulesNames[0]
	}
	subnodes := make([]*ASTNode, 0)
	child := startNode.firstChild
	for child != nil {
		if strings.EqualFold(syntax.RuleName(child.ruleType), ruleName) {
			subnodes = append(subnodes, child)
		}
		child = child.sibling
	}
	return subnodes
}
