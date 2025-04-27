package ast

type Node struct {
	ruleType   int
	start      int
	end        int
	sibling    *Node
	firstChild *Node
}

func NewNode(ruleType int, start int, end int) *Node {
	return &Node{
		ruleType: ruleType,
		start:    start,
		end:      end,
	}
}

func (n *Node) RuleType() int {
	return n.ruleType
}

func (n *Node) Start() int {
	return n.start
}

func (n *Node) End() int {
	return n.end
}

func (n *Node) Sibling() *Node {
	return n.sibling
}

func (n *Node) SetSibling(sibling *Node) {
	n.sibling = sibling
}

func (n *Node) FirstChild() *Node {
	return n.firstChild
}

func (n *Node) SetFirstChild(firstChild *Node) {
	n.firstChild = firstChild
}
