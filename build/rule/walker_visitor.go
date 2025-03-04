package rule

type walkerVisitor struct {
	doVisit func(r Rule)
	isVisit func(r Rule) bool
	visited map[Rule]bool
}

var _ LexerVisitor = &walkerVisitor{}

func newWalkerVisitor(visit func(r Rule), isVisit func(r Rule) bool) *walkerVisitor {
	wv := &walkerVisitor{
		doVisit: visit,
		isVisit: isVisit,
		visited: make(map[Rule]bool),
	}
	return wv
}

func (w *walkerVisitor) shouldVisit(rule Rule) bool {
	return w.isVisit == nil || w.isVisit(rule)
}

// VisitAndRule implements LexerVisitor.
func (w *walkerVisitor) VisitAndRule(rule *AndRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	rules := rule.Rules()
	for _, r := range rules {
		if w.shouldVisit(r) {
			r.Visit(w)
		}
	}
}

// VisitOrRule implements LexerVisitor.
func (w *walkerVisitor) VisitOrRule(rule *OrRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	rules := rule.Rules()
	for _, r := range rules {
		if w.shouldVisit(r) {
			r.Visit(w)
		}
	}
}

// VisitCharRule implements LexerVisitor.
func (w *walkerVisitor) VisitCharRule(rule *CharRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
}

// VisitNonTerminal implements LexerVisitor.
func (w *walkerVisitor) VisitNonTerminal(rule *NonTerminalRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	if w.shouldVisit(rule.Rule()) {
		rule.Rule().Visit(w)
	}
}

// VisitNotRule implements LexerVisitor.
func (w *walkerVisitor) VisitNotRule(rule *NotRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	if w.shouldVisit(rule.Rule()) {
		rule.Rule().Visit(w)
	}
}

// VisitOneOrMoreRule implements LexerVisitor.
func (w *walkerVisitor) VisitOneOrMoreRule(rule *OneOrMoreRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	if w.shouldVisit(rule.Rule()) {
		rule.Rule().Visit(w)
	}
}

// VisitOptionalRule implements LexerVisitor.
func (w *walkerVisitor) VisitOptionalRule(rule *OptionalRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	if w.shouldVisit(rule.Rule()) {
		rule.Rule().Visit(w)
	}
}

// VisitRangeRule implements LexerVisitor.
func (w *walkerVisitor) VisitRangeRule(rule *RangeRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
}

// VisitStringRule implements LexerVisitor.
func (w *walkerVisitor) VisitStringRule(rule *StringRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
}

// VisitTestRule implements LexerVisitor.
func (w *walkerVisitor) VisitTestRule(rule *TestRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	if w.shouldVisit(rule.Rule()) {
		rule.Rule().Visit(w)
	}
}

// VisitZeroOrMoreRule implements LexerVisitor.
func (w *walkerVisitor) VisitZeroOrMoreRule(rule *ZeroOrMoreRule) {
	if _, found := w.visited[rule]; found {
		return
	}
	w.visited[rule] = true
	w.doVisit(rule)
	if w.shouldVisit(rule.Rule()) {
		rule.Rule().Visit(w)
	}
}
