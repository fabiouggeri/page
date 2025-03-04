package lexer

type Token struct {
	index int
	len   int
	row   int
	col   int
	types []int
}

func NewToken(index, len, row, col int, types []int) *Token {
	return &Token{
		index: index,
		len:   len,
		row:   row,
		col:   col,
		types: types,
	}
}

func (t *Token) Index() int {
	return t.index
}

func (t *Token) Len() int {
	return t.len
}

func (t *Token) Row() int {
	return t.row
}

func (t *Token) Col() int {
	return t.col
}

func (t *Token) Types() []int {
	return t.types
}
