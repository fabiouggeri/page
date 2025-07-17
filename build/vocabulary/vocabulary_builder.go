package vocabulary

import (
	"github.com/fabiouggeri/page/build/automata"
	"github.com/fabiouggeri/page/build/grammar"
	"github.com/fabiouggeri/page/build/rule"
	runtime "github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/util"
)

type vocabularyBuilder struct {
	maxSymbol     rune
	tokensTypes   *util.Set[string]
	tokensOptions map[string]*util.Set[*rule.RuleOption]
	tokensMap     map[string]int
	dfa           *automata.State
}

func FromGrammar(grammar *grammar.Grammar) *runtime.Vocabulary {
	rules := grammar.LexerRules()
	return FromDFA(automata.NFAToDFA(RulesToNFA(rules...)))
}

func FromDFA(dfa *automata.State) *runtime.Vocabulary {
	vb := &vocabularyBuilder{
		maxSymbol:     rune(0),
		tokensTypes:   util.NewSet[string](),
		tokensOptions: make(map[string]*util.Set[*rule.RuleOption]),
		tokensMap:     make(map[string]int),
		dfa:           dfa,
	}
	return vb.build()
}

func (vb *vocabularyBuilder) build() *runtime.Vocabulary {
	vb.dfa.WalkThrough(vb.visitState, func(souce *automata.State, char rune, target *automata.State) bool { return true })
	tokenId := 1
	tokensTypes := vb.tokensTypes.Items()
	tokensNames := make([]string, 0, len(tokensTypes)+1)
	tokensOptions := make([]int, 0, cap(tokensNames))
	tokensNames = append(tokensNames, "EOI")
	tokensOptions = append(tokensOptions, 0)
	vb.tokensMap["EOI"] = 0
	for _, tokenType := range tokensTypes {
		if tokenType != "EOI" {
			tokensNames = append(tokensNames, tokenType)
			optionsSet := 0
			if options, found := vb.tokensOptions[tokenType]; found {
				for _, option := range options.Items() {
					optionsSet |= option.Code()
				}
			}
			tokensOptions = append(tokensOptions, optionsSet)
			vb.tokensMap[tokenType] = tokenId
			tokenId++
		}
	}
	return runtime.NewVocabulary(tokensNames, tokensOptions, vb.buildTransitionTable(), vb.buildTokensTable())
}

func (vb *vocabularyBuilder) visitState(state *automata.State) bool {
	for _, tt := range state.RulesTypes() {
		vb.tokensTypes.Add(tt.Name())
		vb.addTokenOptions(tt.Name(), tt.Rule().Options()...)
	}
	for _, symbol := range state.Symbols() {
		if symbol > vb.maxSymbol && symbol != automata.ANY {
			vb.maxSymbol = symbol
		}
	}
	return true
}

func (vb *vocabularyBuilder) addTokenOptions(tokenName string, optionsToSet ...*rule.RuleOption) {
	options, found := vb.tokensOptions[tokenName]
	if !found {
		options = util.NewSet[*rule.RuleOption]()
		vb.tokensOptions[tokenName] = options
	}
	options.AddAll(optionsToSet...)
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
