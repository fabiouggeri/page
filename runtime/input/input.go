package input

type Input interface {
	GetChar() rune
	Skip() bool
	Eof() bool
	Index() int
	Close()
}

type StringInput struct {
	input string
	index int
}

func NewStringInput(input string) *StringInput {
	return &StringInput{
		input: input,
		index: 0,
	}
}

func (i *StringInput) GetChar() rune {
	if i.index >= len(i.input) {
		return '\x00'
	}
	c := rune(i.input[i.index])
	return c
}

func (i *StringInput) Skip() bool {
	if i.index >= len(i.input) {
		return false
	}
	i.index++
	return true
}

func (i *StringInput) Eof() bool {
	return i.index >= len(i.input)
}

func (i *StringInput) Index() int {
	return i.index
}

func (i *StringInput) Close() {
	// do nothing
}
