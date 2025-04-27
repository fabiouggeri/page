package lexer

import (
	"fmt"

	"github.com/fabiouggeri/page/runtime/error"
)

type lexerError struct {
	index   int
	row     int
	col     int
	code    int
	message string
}

const LEX_ERROR_EOF = 1
const LEX_ERROR_INVALID_CHAR = 2

var _ error.Error = &lexerError{}

func newError(index, row, col, code int, message string, args ...any) *lexerError {
	return &lexerError{
		index:   index,
		row:     row,
		col:     col,
		code:    code,
		message: fmt.Sprintf(message, args...),
	}
}

func (e *lexerError) Index() int {
	return e.index
}

func (e *lexerError) Row() int {
	return e.row
}

func (e *lexerError) Col() int {
	return e.col
}

func (e *lexerError) Code() int {
	return e.code
}

func (e *lexerError) Message() string {
	return e.message
}

func (e *lexerError) String() string {
	return fmt.Sprintf("Error %d: %s at row %d, col %d", e.code, e.message, e.row, e.col)
}
