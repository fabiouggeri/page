package parser

import (
	"fmt"

	"github.com/fabiouggeri/page/runtime/error"
)

type ParserError struct {
	col     int
	row     int
	code    int
	message string
}

const LEXER_ERROR = 1

var _ error.Error = &ParserError{}

// Code implements error.Error.
func (p *ParserError) Code() int {
	return p.code
}

// Col implements error.Error.
func (p *ParserError) Col() int {
	return p.col
}

// Message implements error.Error.
func (p *ParserError) Message() string {
	return p.message
}

// Row implements error.Error.
func (p *ParserError) Row() int {
	return p.row
}

// String implements error.Error.
func (p *ParserError) String() string {
	return fmt.Sprintf("Error %d: %s at row %d, col %d", p.code, p.message, p.row, p.col)
}
