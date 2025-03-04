package lexer

import (
	"github.com/fabiouggeri/page/runtime/error"
	"github.com/fabiouggeri/page/runtime/input"
)

type Lexer struct {
	vocabulary *Vocabulary
	input      input.Input
	lastChar   rune
	index      int
	row        int
	col        int
	errors     []error.Error
	tokens     []*Token
}

const LEX_ERROR_EOF = 1
const LEX_ERROR_INVALID_CHAR = 2

const TKN_EOF = 0

var EOF_TOKEN = Token{}

func New(vocabulary *Vocabulary, input input.Input) *Lexer {
	return &Lexer{
		vocabulary: vocabulary,
		input:      input,
		lastChar:   '\x00',
		index:      0,
		row:        1,
		col:        1,
	}
}

func (l *Lexer) Errors() []error.Error {
	return l.errors
}

func (l *Lexer) Tokens() []*Token {
	tkn, _ := l.NextToken()
	for tkn != &EOF_TOKEN {
		tkn, _ = l.NextToken()
	}
	return l.tokens
}

func (l *Lexer) NextToken() (*Token, error.Error) {
	if l.index < len(l.tokens) {
		token := l.tokens[l.index]
		l.index++
		return token, nil
	}
	if l.input.Eof() {
		return &EOF_TOKEN, nil
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
	col := l.col
	row := l.row
	tokenLen := 0
	state := 0
	transitionsTable := l.vocabulary.TransitionsTable()
	start := l.input.Index()
	for {
		var nextState int
		c := l.input.GetChar()
		if c == '\x00' {
			if state == 0 {
				return nil, l.error(LEX_ERROR_EOF, start, l.row, l.col, "Unexpected end of file")
			}
			return &Token{index: start, len: tokenLen, row: row, col: col, types: l.vocabulary.TokenTypes(state)}, nil
		}
		if int(c) >= len(transitionsTable[state]) {
			nextState = transitionsTable[state][0]
			if nextState == 0 {
				l.skipChar(c)
				return nil, l.error(LEX_ERROR_INVALID_CHAR, start, l.row, l.col, "Invalid character '%c'", c)
			}
		} else {
			nextState = transitionsTable[state][c]
		}
		if nextState == 0 {
			if state == 0 {
				l.skipChar(c)
				return nil, l.error(LEX_ERROR_INVALID_CHAR, start, l.row, l.col, "Invalid character '%c'", c)
			}
			return &Token{index: start, len: tokenLen, row: row, col: col, types: l.vocabulary.TokenTypes(state)}, nil
		}
		state = nextState
		l.lastChar = c
		l.skipChar(c)
		tokenLen++
	}
}

func (l *Lexer) skipChar(c rune) {
	if c == '\n' {
		l.row++
		l.col = 1
	} else if c != '\n' {
		l.col++
	}
	l.input.Skip()
}

func (l *Lexer) error(code int, index int, row int, col int, message string, args ...any) error.Error {
	err := newError(index, row, col, code, message, args...)
	l.errors = append(l.errors, err)
	return err
}
