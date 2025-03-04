package source

type Source interface {
	StringAt(start, end uint32) string
	Match(c rune) bool
	MatchIgnoreCase(c rune) bool
	MatchRange(start, end rune) bool
	MatchString(text string) bool
	MatchStringIgnoreCase(text string) bool
	EOI() bool
	Index() uint32
	SetIndex(index uint32)
}
