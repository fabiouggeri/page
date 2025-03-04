package language

type Language struct {
	tokensNames []string
	tokensTypes [][]uint32
	lexerTable  [][]uint32
	parserTable [][]uint32
}
