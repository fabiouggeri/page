package vocabulary

import (
	"github.com/fabiouggeri/page/build/automata"
	"github.com/fabiouggeri/page/build/grammar"
	runtime "github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/util"
)

type vocabularyBuilder struct {
	maxSymbol   rune
	tokensTypes *util.Set[string]
	tokensMap   map[string]int
	dfa         *automata.State
}

func FromGrammar(grammar *grammar.Grammar) *runtime.Vocabulary {
	return FromDFA(automata.NFAToDFA(RulesToNFA(grammar.LexerRules()...)))
}

func FromDFA(dfa *automata.State) *runtime.Vocabulary {
	vb := &vocabularyBuilder{
		maxSymbol:   rune(0),
		tokensTypes: util.NewSet[string](),
		tokensMap:   make(map[string]int),
		dfa:         dfa,
	}
	return vb.build()
}

func (vb *vocabularyBuilder) build() *runtime.Vocabulary {
	vb.dfa.WalkThrough(vb.visitState, func(souce *automata.State, char rune, target *automata.State) bool { return true })
	tokenId := 1
	tokensNames := make([]string, 0, vb.dfa.RulesTypesCount())
	tokensTypes := vb.tokensTypes.Items()
	tokensNames = append(tokensNames, "eof")
	vb.tokensMap["eof"] = 0
	for _, tokenType := range tokensTypes {
		tokensNames = append(tokensNames, tokenType)
		vb.tokensMap[tokenType] = tokenId
		tokenId++
	}
	v := runtime.NewVocabulary(tokensNames, vb.buildTransitionTable(), vb.buildTokensTable())
	return v
}

func (vb *vocabularyBuilder) visitState(state *automata.State) bool {
	for _, tt := range state.RulesTypes() {
		vb.tokensTypes.Add(tt.Name())
	}
	for _, symbol := range state.Symbols() {
		if symbol > vb.maxSymbol && symbol != automata.ANY {
			vb.maxSymbol = symbol
		}
	}
	return true
}

func (vb *vocabularyBuilder) buildTransitionTable() [][]int {
	states := vb.dfa.AllStates()
	transitionTable := make([][]int, len(states))
	for _, s := range states {
		entry := createTransitionsTableEntry(int(vb.maxSymbol) + 1)
		for symbol, targets := range s.Transitions() {
			items := targets.Items()
			if symbol != automata.ANY {
				entry[int(symbol)] = int(items[0].Id()) // DFA must have only one target for each symbol
			} else {
				entry[0] = int(items[0].Id()) // DFA must have only one target for each symbol
			}
		}
		transitionTable[s.Id()] = entry
	}
	return transitionTable
}

func createTransitionsTableEntry(size int) []int {
	a := make([]int, size)
	for i := 0; i < size; i++ {
		a[i] = 0
	}
	return a
}

func (vb *vocabularyBuilder) buildTokensTable() [][]int {
	states := vb.dfa.AllStates()
	tokensTable := make([][]int, len(states))
	for _, s := range states {
		tokenTypes := s.RulesTypes()
		a := make([]int, 0, len(tokenTypes))
		for _, tt := range tokenTypes {
			a = append(a, vb.tokenId(tt.Name()))
		}
		tokensTable[s.Id()] = a
	}
	return tokensTable
}

func (vb *vocabularyBuilder) tokenId(tokenName string) int {
	if name, found := vb.tokensMap[tokenName]; found {
		return name
	}
	return 0
}
