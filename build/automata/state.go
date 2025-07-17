package automata

import (
	"fmt"
	"strings"

	"github.com/fabiouggeri/page/util"
	"golang.org/x/exp/slices"
)

type Symbol = int32

const (
	EPSILON Symbol = 0 // '\u03B5'
	ANY     Symbol = 0x7FFFFFFF
)

type State struct {
	id          int32
	initial     bool
	final       bool
	transitions map[Symbol]*util.Set[*State]
	rulesTypes  *util.Set[*RuleType]
}

func NewState(id int32, initial, final bool) *State {
	return &State{id: id,
		initial:     initial,
		final:       final,
		transitions: make(map[Symbol]*util.Set[*State]),
		rulesTypes:  util.NewSet[*RuleType]()}
}

func (s *State) Id() int32 {
	return s.id
}

func (s *State) Initial() bool {
	return s.initial
}

func (s *State) SetInitial(i bool) {
	s.initial = i
}

func (s *State) Final() bool {
	return s.final
}

func (s *State) SetFinal(f bool) {
	s.final = f
}

func (s *State) AddTransitions(sym Symbol, target ...*State) *State {
	set, found := s.transitions[sym]
	if !found {
		set = util.NewSet[*State]()
		s.transitions[sym] = set
	}
	set.AddAll(target...)
	return s
}

func (s *State) SetTransitions(transitions map[Symbol]*util.Set[*State]) *State {
	s.transitions = transitions
	return s
}

func (s *State) visitState(visited map[int32]bool, doThis func(state *State) bool, isMove func(souce *State, sym Symbol, target *State) bool) bool {
	_, found := visited[s.id]
	if found {
		return true
	}
	visited[s.id] = true
	if !doThis(s) {
		return false
	}
	for char, transition := range s.transitions {
		for _, targetState := range transition.Items() {
			if isMove(s, char, targetState) && !targetState.visitState(visited, doThis, isMove) {
				return false
			}
		}
	}
	return true
}

func (s *State) WalkThrough(doThis func(state *State) bool, isMove func(souce *State, sym Symbol, target *State) bool) {
	visited := make(map[int32]bool)
	s.visitState(visited, doThis, isMove)
}

func (s *State) FinalStates() []*State {
	finalStates := make([]*State, 0)
	s.WalkThrough(func(state *State) bool {
		if state.final {
			finalStates = append(finalStates, state)
		}
		return true
	},
		func(souce *State, sym Symbol, target *State) bool {
			return true
		})
	return finalStates
}

func (s *State) SetInitialFinal(s1 *State, s2 *State) {
	s1.AddTransitions(EPSILON, s)
	for _, es := range s.FinalStates() {
		es.AddTransitions(EPSILON, s2)
		es.final = false
	}
	s.initial = false
}

func (s *State) AddRuleType(tt *RuleType) {
	s.rulesTypes.Add(tt)
}

func (s *State) RulesTypes() []*RuleType {
	return s.rulesTypes.Items()
}

func (s *State) RulesTypesCount() int {
	return s.rulesTypes.Length()
}

func (s *State) Transitions() map[Symbol]*util.Set[*State] {
	return s.transitions
}

func (s *State) transitionsToDot(visited map[int32]bool, writer util.TextWriter) {
	_, found := visited[s.id]
	if found {
		return
	}
	visited[s.id] = true
	targets := make(map[*State][]rune)
	for sym, transition := range s.transitions {
		for _, target := range transition.Items() {
			_, found := targets[target]
			if !found {
				targets[target] = make([]rune, 0)
			}
			targets[target] = append(targets[target], sym)
		}
	}
	for target, label := range targets {
		writer.WriteF("%d -> %d [label=\"%s\"]", s.id, target.id, symbolsToLabel(label)).NewLine()
		target.transitionsToDot(visited, writer)
	}
}

func symbolsToLabel(symbols []rune) string {
	str := &strings.Builder{}
	slices.SortFunc(symbols, func(a, b rune) int { return int(a - b) })
	if slices.Index(symbols, ANY) >= 0 {
		lastChar := rune(0)
		str.WriteString("[^")
		for _, c := range symbols {
			if c > lastChar+1 && c != ANY {
				for r := lastChar + 1; r < c; r++ {
					charToLabel(str, r)
				}
			}
			lastChar = c
		}
		str.WriteRune(']')
	} else {
		for _, c := range symbols {
			charToLabel(str, c)
		}

	}
	return str.String()
}

func charToLabel(str *strings.Builder, c rune) {
	switch c {
	case '\n':
		str.WriteString("\\\\n")
	case '\r':
		str.WriteString("\\\\r")
	case '\t':
		str.WriteString("\\\\t")
	case 0x03:
		str.WriteString("EOI")
	case '"':
		str.WriteString("\\\"")
	case '\\':
		str.WriteString("\\\\")
	case EPSILON:
		str.WriteRune('€')
	case ANY:
		str.WriteString("...")
	default:
		if c < 32 {
			str.WriteString(fmt.Sprintf("0x%02x", c))
		} else {
			str.WriteRune(c)
		}
	}
}

func (s *State) ToDot(name string, writer util.TextWriter) {
	writer.WriteString("digraph ").WriteString(name).WriteString(" {").NewLine()
	writer.Indent(3)
	writer.WriteString("fontname=\"Helvetica,Arial,sans-serif\"").NewLine()
	writer.WriteString("node [fontname=\"Helvetica,Arial,sans-serif\"]").NewLine()
	writer.WriteString("edge [fontname=\"Helvetica,Arial,sans-serif\"]").NewLine()
	writer.WriteString("rankdir=LR").NewLine()
	writer.WriteString("node [shape = doublecircle];")
	finalStates := s.FinalStates()
	for _, fs := range finalStates {
		writer.WriteF(" %d", fs.id)
	}
	writer.NewLine()
	writer.WriteString("node [shape = circle]").NewLine()
	for _, fs := range finalStates {
		writer.WriteF("%d [label=\"%d\\n", fs.id, fs.id)
		rulesTypes := fs.rulesTypes.Items()
		for i, tt := range rulesTypes {
			if i > 0 {
				writer.WriteRune(',')
			}
			writer.WriteString(tt.Name())
		}
		writer.WriteString("\"]").NewLine()
	}
	s.transitionsToDot(make(map[int32]bool), writer)
	writer.Indent(-3).WriteRune('}').NewLine()
}

func (s *State) String() string {
	w := util.NewStringTextWriter()
	s.ToDot("Automata", w)
	return w.String()
}

func (s *State) AllStates() []*State {
	states := make([]*State, 0)
	s.WalkThrough(func(state *State) bool {
		states = append(states, state)
		return true
	},
		func(souce *State, sym Symbol, target *State) bool {
			return true
		})
	return states
}

func (s *State) Symbols() []Symbol {
	symbols := make([]Symbol, 0)
	for symbol := range s.transitions {
		symbols = append(symbols, symbol)
	}
	return symbols
}

func (s *State) AllSymbols() []Symbol {
	symbols := make(map[Symbol]bool, 256)
	s.WalkThrough(
		func(state *State) bool {
			return true
		},
		func(souce *State, sym Symbol, target *State) bool {
			if sym != EPSILON {
				symbols[sym] = true
			}
			return true
		},
	)
	allSymbols := make([]Symbol, 0, len(symbols))
	for key := range symbols {
		allSymbols = append(allSymbols, key)
	}
	return allSymbols
}

func (s *State) EpsilonClosures() *util.Set[*State] {
	states := util.NewSet[*State]()
	s.WalkThrough(
		func(state *State) bool {
			states.Add(state)
			return true
		},
		func(souce *State, sym Symbol, target *State) bool {
			return sym == EPSILON
		},
	)
	return states
}

func (s *State) symbolTargets(sym Symbol, symbolTargets *util.Set[*State]) {
	targets, found := s.transitions[sym]
	if found {
		for _, s := range targets.Items() {
			symbolTargets.Add(s)
		}
	}
	targets, found = s.transitions[EPSILON]
	if found {
		for _, epsilonTarget := range targets.Items() {
			epsilonTarget.symbolTargets(sym, symbolTargets)
		}
	}
}

func (s *State) Targets(sym Symbol) *util.Set[*State] {
	states := util.NewSet[*State]()
	s.symbolTargets(sym, states)
	return states
}
