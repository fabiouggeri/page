package lexer

import (
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/runtime/error"
	"github.com/fabiouggeri/page/runtime/input"
)

type Lexer struct {
	vocabulary  *Vocabulary
	input       input.Input
	index       int
	row         int
	col         int
	tokensLine  int
	errors      []error.Error
	tokens      []*Token
	eof         bool
	onlyIgnored bool
}

type lexerState struct {
	index       int
	state       int
	row         int
	col         int
	tokensLine  int
	onlyIgnored bool
}

const TKN_EOF = 0

func New(vocabulary *Vocabulary, input input.Input) *Lexer {
	return &Lexer{
		vocabulary: vocabulary,
		input:      input,
		index:      0,
		row:        1,
		col:        1,
	}
}

func (l *Lexer) Errors() []error.Error {
	return l.errors
}

func (l *Lexer) Input() input.Input {
	return l.input
}

func (l *Lexer) Index() int {
	return l.index
}

func (l *Lexer) InputIndex() int {
	return l.input.Index()
}

func (l *Lexer) SetIndex(newIndex int) {
	if newIndex >= 0 && newIndex < len(l.tokens) {
		l.index = newIndex
	}
}

func (l *Lexer) Row() int {
	return l.row
}

func (l *Lexer) Col() int {
	return l.col
}

func (l *Lexer) Tokens() []*Token {
	tkn, _ := l.NextToken()
	for tkn != nil {
		tkn, _ = l.NextToken()
	}
	return l.tokens
}

func (l *Lexer) Token(index int) (*Token, error.Error) {
	for !l.eof && index >= len(l.tokens) {
		token, err := l.readNextToken()
		if err != nil {
			return nil, err
		}
		l.tokens = append(l.tokens, token)
	}
	if index < len(l.tokens) {
		token := l.tokens[index]
		return token, nil
	}
	return nil, l.error(LEX_ERROR_EOF, l.input.Index(), l.row, l.col, "Unexpected end of file")
}

func (l *Lexer) NextToken() (*Token, error.Error) {
	if l.index < len(l.tokens) {
		token := l.tokens[l.index]
		l.index++
		return token, nil
	}
	if l.input.Eof() {
		if !l.eof {
			l.eof = true
			eofTkn := &Token{
				index: l.input.Index(),
				len:   0,
				row:   l.row,
				col:   l.col,
				types: []int{TKN_EOF},
			}
			l.tokens = append(l.tokens, eofTkn)
			l.index++
			return eofTkn, nil
		} else {
			return nil, l.error(LEX_ERROR_EOF, l.input.Index(), l.row, l.col, "Unexpected end of file")
		}
	}
	token, err := l.readNextToken()
	if err != nil {
		return nil, err
	}
	l.tokens = append(l.tokens, token)
	l.index++
	return token, nil
}

func (l *Lexer) readNextToken() (*Token, error.Error) {
	var lastValidState *lexerState
	col := l.col
	row := l.row
	state := 0
	transitionsTable := l.vocabulary.TransitionsTable()
	start := l.input.Index()
	for {
		var nextState int
		c := l.input.GetChar()

		if int(c) >= len(transitionsTable[state]) {
			nextState = transitionsTable[state][0]
		} else if int(c) > 0 {
			nextState = transitionsTable[state][c]
		} else if state == 0 {
			return nil, l.error(LEX_ERROR_EOF, start, l.row, l.col, "Unexpected end of file")
		} else if tt := l.validTokensTypes(state); len(tt) > 0 {
			if l.row > row {
				l.onlyIgnored = true
				l.tokensLine = 0
			} else {
				l.onlyIgnored = l.onlyIgnored && l.onlyIgnoredTypes(tt)
				l.tokensLine++
			}
			return &Token{index: start, len: l.input.Index() - start, row: row, col: col, types: tt}, nil
		}
		if nextState == 0 {
			tt := l.validTokensTypes(state)
			if len(tt) == 0 {
				// has a previous valid state, return it
				if lastValidState != nil {
					tt := l.validTokensTypes(lastValidState.state)
					l.input.SetIndex(lastValidState.index + 1)
					l.tokensLine = lastValidState.tokensLine
					l.onlyIgnored = lastValidState.onlyIgnored
					l.row = lastValidState.row
					l.col = lastValidState.col
					return &Token{index: start, len: l.input.Index() - start, row: row, col: col, types: tt}, nil
				}
				l.skipChar(c)
				return nil, l.error(LEX_ERROR_INVALID_CHAR, start, l.row, l.col, "Invalid character '%c'", c)
			}
			if l.row > row {
				l.onlyIgnored = true
				l.tokensLine = 0
			} else {
				l.onlyIgnored = l.onlyIgnored && l.onlyIgnoredTypes(tt)
				l.tokensLine++
			}
			return &Token{index: start, len: l.input.Index() - start, row: row, col: col, types: tt}, nil
		}
		state = nextState
		l.skipChar(c)
		// store the last valid state if it is a final state
		if l.vocabulary.IsFinalState(state) && !l.vocabulary.AllTokensTypesHasOption(state, rule.IGNORE) {
			lastValidState = &lexerState{
				index:       start,
				state:       state,
				row:         l.row,
				col:         l.col,
				tokensLine:  l.tokensLine,
				onlyIgnored: l.onlyIgnored,
			}
		}
	}
}

func (l *Lexer) validTokensTypes(state int) []int {
	tokensTypes := l.vocabulary.TokenTypes(state)
	validTokens := make([]int, 0, len(tokensTypes))
	for _, tokenType := range tokensTypes {
		if !l.vocabulary.HasOptions(tokenType) {
			validTokens = append(validTokens, tokenType)
		} else if l.vocabulary.HasOption(tokenType, rule.START_LINE) {
			if l.tokensLine == 0 {
				validTokens = append(validTokens, tokenType)
			}
		} else if l.vocabulary.HasOption(tokenType, rule.ONLY_IGNORED) {
			if l.onlyIgnored {
				validTokens = append(validTokens, tokenType)
			}
		} else {
			validTokens = append(validTokens, tokenType)
		}
	}
	return validTokens
}

func (l *Lexer) onlyIgnoredTypes(tokenTypes []int) bool {
	for _, t := range tokenTypes {
		if !l.vocabulary.HasOption(t, rule.IGNORE) {
			return false
		}
	}
	return true
}

func (l *Lexer) skipChar(c rune) {
	switch c {
	case '\n':
		l.row++
		l.col = 1
	case '\r':
		// do nothing
	default:
		l.col++
	}
	l.input.Skip()
}

func (l *Lexer) error(code int, index int, row int, col int, message string, args ...any) error.Error {
	err := newError(index, row, col, code, message, args...)
	l.errors = append(l.errors, err)
	return err
}

func (l *Lexer) IsIgnored(tkn *Token) bool {
	for _, tt := range tkn.types {
		if l.vocabulary.HasOption(tt, rule.IGNORE) {
			return true
		}
	}
	return false
}

func (l *Lexer) Vocabulary() *Vocabulary {
	return l.vocabulary
}
