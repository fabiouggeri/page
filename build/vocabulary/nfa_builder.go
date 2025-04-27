package vocabulary

import (
	"unicode"

	"github.com/fabiouggeri/page/build/automata"
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/util"
)

type nfaVisitor struct {
	nextId     int32
	states     *util.Deque[*automata.State]
	nextRuleId uint16
	rulesTypes map[string]*automata.RuleType
	maxSymbol  rune
}

var _ rule.RuleVisitor = &nfaVisitor{}

func RulesToNFA(rules ...*rule.NonTerminalRule) *automata.State {
	maxSymbol := rulesMaxSymbol(rules...)
	v := &nfaVisitor{nextId: 0,
		states:     util.NewDeque[*automata.State](),
		nextRuleId: 0,
		rulesTypes: make(map[string]*automata.RuleType),
		maxSymbol:  maxSymbol,
	}
	s1 := v.newInitialState()
	for _, r := range rules {
		ruleType := v.registerRule(r)
		r.Visit(v)
		state := v.pop()
		for _, finalState := range state.FinalStates() {
			finalState.AddRuleType(ruleType)
		}
		state.SetInitial(false)
		s1.AddTransitions(automata.EPSILON, state)
	}
	return s1
}

func rulesMaxSymbol(rules ...*rule.NonTerminalRule) rune {
	ms := newMaxSymbolVisitor()
	for _, r := range rules {
		r.Visit(ms)
	}
	return ms.maxSymbol
}

func (n *nfaVisitor) newInitialState() *automata.State {
	s := automata.NewState(n.nextId, true, false)
	n.nextId++
	return s
}

func (n *nfaVisitor) newFinalState() *automata.State {
	s := automata.NewState(n.nextId, false, true)
	n.nextId++
	return s
}

func (n *nfaVisitor) newState() *automata.State {
	s := automata.NewState(n.nextId, false, false)
	n.nextId++
	return s
}

func (n *nfaVisitor) registerRule(rule *rule.NonTerminalRule) *automata.RuleType {
	tt, found := n.rulesTypes[rule.Id()]
	if !found {
		tt = automata.NewRuleId(n.nextRuleId, rule.Id(), rule)
		n.rulesTypes[rule.Id()] = tt
		n.nextRuleId++
	}
	return tt
}

func (n *nfaVisitor) push(s *automata.State) {
	n.states.Push(s)
}

func (n *nfaVisitor) pop() *automata.State {
	s, err := n.states.Pop()
	if err != nil {
		return nil
	}
	return s
}

func (n *nfaVisitor) VisitNonTerminal(rule *rule.NonTerminalRule) {
	rule.Rule().Visit(n)
}

func (n *nfaVisitor) VisitAndRule(rule *rule.AndRule) {
	rules := rule.Rules()
	for i, r := range rules {
		r.Visit(n)
		if i > 0 {
			s2 := n.pop()
			s1 := n.pop()
			for _, es := range s1.FinalStates() {
				es.AddTransitions(automata.EPSILON, s2)
				es.SetFinal(false)
			}
			s2.SetInitial(false)
			n.push(s1)
		}
	}
}

func (n *nfaVisitor) VisitOrRule(visitedRule *rule.OrRule) {
	s1 := n.newInitialState()
	s2 := n.newFinalState()
	rules := visitedRule.Rules()
	if allCharRules(rules) {
		for _, r := range rules {
			switch castRule := r.(type) {
			case *rule.CharRule:
				addCharTransitions(castRule, s1, s2)
			case *rule.RangeRule:
				for c := castRule.Start(); c <= castRule.End(); c++ {
					s1.AddTransitions(c, s2)
				}
			default:
				// do nothing
			}
		}
	} else {
		for _, r := range rules {
			r.Visit(n)
			n.pop().SetInitialFinal(s1, s2)
		}
	}
	n.push(s1)
}

func allCharRules(rules []rule.Rule) bool {
	for _, r := range rules {
		switch castRule := r.(type) {
		case *rule.CharRule,
			*rule.RangeRule:
			// do nothing
		case *rule.OrRule:
			if !allCharRules(castRule.Rules()) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (n *nfaVisitor) VisitCharRule(rule *rule.CharRule) {
	is := n.newInitialState()
	addCharTransitions(rule, is, n.newFinalState())
	n.push(is)
}

func addCharTransitions(rule *rule.CharRule, is *automata.State, fs *automata.State) {
	if rule.CaseSensitive() {
		low := unicode.ToLower(rule.Char())
		up := unicode.ToUpper(rule.Char())
		if low == up {
			is.AddTransitions(low, fs)
		} else {
			is.AddTransitions(low, fs)
			is.AddTransitions(up, fs)
		}
	} else {
		is.AddTransitions(rule.Char(), fs)
	}
}

func (n *nfaVisitor) VisitOneOrMoreRule(rule *rule.OneOrMoreRule) {
	rule.Rule().Visit(n)
	s := n.pop()
	if s == nil {
		return
	}
	s1 := n.newInitialState()
	s2 := n.newFinalState()
	for _, es := range s.FinalStates() {
		es.AddTransitions(automata.EPSILON, s)
	}
	s.SetInitialFinal(s1, s2)
	n.push(s1)
}

func (n *nfaVisitor) VisitZeroOrMoreRule(rule *rule.ZeroOrMoreRule) {
	rule.Rule().Visit(n)
	s := n.pop()
	if s == nil {
		return
	}
	s1 := n.newInitialState()
	s2 := n.newFinalState()
	for _, es := range s.FinalStates() {
		es.AddTransitions(automata.EPSILON, s)
	}
	s.SetInitialFinal(s1, s2)
	s1.AddTransitions(automata.EPSILON, s2)
	n.push(s1)
}

func (n *nfaVisitor) VisitOptionalRule(rule *rule.OptionalRule) {
	rule.Rule().Visit(n)
	s := n.pop()
	if s == nil {
		return
	}
	s1 := n.newInitialState()
	s2 := n.newFinalState()
	s.SetInitialFinal(s1, s2)
	s1.AddTransitions(automata.EPSILON, s2)
	n.push(s1)
}

func (n *nfaVisitor) VisitRangeRule(rule *rule.RangeRule) {
	s1 := n.newInitialState()
	s2 := n.newFinalState()
	for c := rule.Start(); c <= rule.End(); c++ {
		s1.AddTransitions(c, s2)
	}
	n.push(s1)
}

func (n *nfaVisitor) VisitStringRule(rule *rule.StringRule) {
	var len int32 = 1
	runes := []rune(rule.Text())
	s1 := n.newInitialState()
	s2 := s1
	if rule.CaseSensitive() {
		for _, c := range runes {
			ns := n.newState()
			s2.AddTransitions(c, ns)
			s2 = ns
			if len >= rule.Size() {
				s2.SetFinal(true)
			}
			len++
		}
	} else {
		for _, c := range runes {
			ns := n.newState()
			low := unicode.ToLower(c)
			up := unicode.ToUpper(c)
			if low == up {
				s2.AddTransitions(c, ns)
			} else {
				s2.AddTransitions(low, ns)
				s2.AddTransitions(up, ns)
			}
			s2 = ns
			if len >= rule.Size() {
				s2.SetFinal(true)
			}
			len++
		}
	}
	s2.SetFinal(true)
	n.push(s1)
}

// VisitNotRule implements rules.LexerVisitor.
func (n *nfaVisitor) VisitNotRule(rule *rule.NotRule) {
	rule.Rule().Visit(n)
	state := n.pop()
	if state == nil {
		return
	}
	states := state.AllStates()
	for _, s := range states {
		n.negateTransitions(s)
	}
	n.push(state)
}

func (n *nfaVisitor) negateTransitions(state *automata.State) {
	transitions := state.Transitions()
	targetsMap := make(map[*automata.State]*util.Set[automata.Symbol])
	for sym, targets := range transitions {
		if sym != automata.EPSILON {
			symTargets := targets.Items()
			for _, target := range symTargets {
				allSymbols, found := targetsMap[target]
				if found {
					allSymbols.Add(sym)
				} else {
					targetsMap[target] = util.NewSet[automata.Symbol](sym)
				}
			}
		}
	}
	if len(targetsMap) > 0 {
		state.SetTransitions(n.targetsMapToTransitions(targetsMap))
	}
}

func (n *nfaVisitor) targetsMapToTransitions(targetsMap map[*automata.State]*util.Set[automata.Symbol]) map[automata.Symbol]*util.Set[*automata.State] {
	newTransitions := make(map[automata.Symbol]*util.Set[*automata.State])
	for target, symbols := range targetsMap {
		for i := automata.Symbol(1); i <= n.maxSymbol; i++ {
			if !symbols.Contains(i) {
				symTargets, found := newTransitions[i]
				if found {
					symTargets.Add(target)
				} else {
					newTransitions[i] = util.NewSet(target)
				}
			}
		}
		if !symbols.Contains(automata.ANY) {
			symTargets, found := newTransitions[automata.ANY]
			if found {
				symTargets.Add(target)
			} else {
				newTransitions[automata.ANY] = util.NewSet(target)
			}
		}
	}
	return newTransitions
}

func (n *nfaVisitor) VisitTestRule(rule *rule.TestRule) {
	panic("rule type not supported for lexer")
}
