package automata

import (
	"github.com/fabiouggeri/page/util"
	"golang.org/x/exp/slices"
)

type dfaState struct {
	initial     bool
	final       bool
	nfaStates   *util.Set[*State]
	transitions map[rune]*util.Set[*State]
}

func NFAToDFA(state *State) *State {
	allSymbols := state.AllSymbols()
	dfaStates := make([]*dfaState, 0)
	dfaStates = append(dfaStates, newDFAState(util.NewSet(state.EpsilonClosures().Items()...)))
	for i := 0; i < len(dfaStates); i++ {
		dfaState := dfaStates[i]
		dfaStates = buildDFATransitions(dfaStates, allSymbols, dfaState)
	}
	return minimizeDFA(buildDFA(dfaStates))
}

func buildDFATransitions(dfaStates []*dfaState, allSymbols []rune, dfaState *dfaState) []*dfaState {
	// percorre os estados NFA do novo estado DFA
	for _, nfaState := range dfaState.nfaStates.Items() {
		dfaStates = addNFATransitions(dfaStates, allSymbols, dfaState, nfaState)
	}
	for _, dfaTransition := range dfaState.transitions {
		if !foundDFAState(dfaStates, dfaTransition) {
			dfaStates = append(dfaStates, newDFAState(dfaTransition))
		}
	}
	return dfaStates
}

func addNFATransitions(dfaStates []*dfaState, allSymbols []rune, dfaState *dfaState, nfaState *State) []*dfaState {
	dfaState.initial = dfaState.initial || nfaState.initial
	dfaState.final = dfaState.final || nfaState.final
	for _, symbol := range allSymbols {
		targets := nfaState.Targets(symbol)
		if targets.Length() > 0 {
			dfaTransition, found := dfaState.transitions[symbol]
			if !found {
				dfaTransition = util.NewSet[*State]()
				dfaState.transitions[symbol] = dfaTransition
			}
			dfaTransition.AddAll(targets.Items()...)
		}
	}
	return dfaStates
}

func newDFAState(nfaStates *util.Set[*State]) *dfaState {
	return &dfaState{nfaStates: nfaStates, transitions: make(map[rune]*util.Set[*State])}
}

func foundDFAState(dfaStates []*dfaState, nfaStates *util.Set[*State]) bool {
	for _, d := range dfaStates {
		if d.nfaStates.Equals(nfaStates) {
			return true
		}
	}
	return false
}

func buildDFA(dfaStates []*dfaState) *State {
	var initialState *State
	nextId := int32(0)
	states := make([]*State, 0, len(dfaStates))
	for _, dfaState := range dfaStates {
		s := NewState(nextId, dfaState.initial, dfaState.final)
		states = append(states, s)
		checkIsFinal(s, dfaState)
		nextId++
		if s.initial {
			initialState = s
		}
	}
	for index, dfaState := range dfaStates {
		for symbol, targets := range dfaState.transitions {
			targetIndex := findDfaState(dfaStates, targets)
			states[index].AddTransitions(symbol, states[targetIndex])
		}
	}
	return initialState
}

func checkIsFinal(state *State, dfaState *dfaState) {
	for _, nfaState := range dfaState.nfaStates.Items() {
		for _, closure := range nfaState.EpsilonClosures().Items() {
			if closure.final {
				state.final = true
				state.rulesTypes.AddAll(closure.rulesTypes.Items()...)
			}
		}
	}
}

func findDfaState(allStates []*dfaState, targets *util.Set[*State]) int {
	for index, dfaState := range allStates {
		if dfaState.nfaStates.Equals(targets) {
			return index
		}
	}
	return -1
}

func minimizeDFA(state *State) *State {
	allSymbols := state.AllSymbols()
	allStates := state.AllStates()
	slices.SortFunc(allStates, func(a, b *State) int { return int(a.id - b.id) })

	statesTable := buildStatesTable(allStates)
	checkUnmarkedPairs(statesTable, allStates, allSymbols)
	return combineUnmarkedStates(statesTable, allStates)
}

func combineUnmarkedStates(statesTable [][]bool, allStates []*State) *State {
	var initialState *State
	nextId := int32(0)
	newStates := make(map[int]*State)
	for s1Id := range statesTable {
		s1 := allStates[s1Id]
		combine := false
		row := statesTable[s1Id]
		for s2Id := s1Id + 1; s2Id < len(row); s2Id++ {
			if !row[s2Id] {
				combine = true
				s2 := allStates[s2Id]
				combinedState := NewState(nextId, s1.initial || s2.initial, s1.final || s2.final)
				if combinedState.final {
					combinedState.rulesTypes.AddAll(s1.RulesTypes()...)
					combinedState.rulesTypes.AddAll(s2.RulesTypes()...)
				}
				newStates[s1Id] = combinedState
				newStates[s2Id] = combinedState
			}
		}
		if !combine {
			newState := NewState(nextId, s1.initial, s1.final)
			if newState.final {
				newState.rulesTypes.AddAll(s1.RulesTypes()...)
			}
			newStates[s1Id] = newState
		}
	}
	for _, s := range allStates {
		newState := newStates[int(s.id)]
		for symbol, targets := range s.transitions {
			target := targets.Items()[0]
			newTarget := allStates[target.id]
			newState.AddTransitions(symbol, newTarget)
		}
		if newState.initial {
			initialState = newState
		}
	}
	return initialState
}

func checkUnmarkedPairs(statesTable [][]bool, allStates []*State, allSymbols []Symbol) {
	for s1Id := range statesTable {
		row := statesTable[s1Id]
		for s2Id := s1Id + 1; s2Id < len(row); s2Id++ {
			if !row[s2Id] {
				row[s2Id] = markState(statesTable, allStates, s1Id, s2Id, allSymbols)
			}
		}
	}
}

func markState(statesTable [][]bool, allStates []*State, s1Id, s2Id int, allSymbols []Symbol) bool {
	s1 := allStates[s1Id]
	s2 := allStates[s2Id]
	for _, symbol := range allSymbols {
		targetsS1 := s1.transitions[symbol]
		targetsS2 := s2.transitions[symbol]
		if targetsS1 != nil && targetsS2 != nil {
			targetS1 := targetsS1.Items()[0]
			targetS2 := targetsS2.Items()[0]
			if statesTable[targetS1.id][targetS2.id] {
				return true
			}
		}
	}
	return false
}

func buildStatesTable(allStates []*State) [][]bool {
	statesTable := make([][]bool, len(allStates))
	for row := range statesTable {
		statesTable[row] = make([]bool, len(allStates))
	}
	for s1Id := range statesTable {
		row := statesTable[s1Id]
		for s2Id := s1Id + 1; s2Id < len(row); s2Id++ {
			s1 := allStates[s1Id]
			s2 := allStates[s2Id]
			row[s2Id] = s1.final && !s2.final
		}
	}
	return statesTable
}
