package util

import "strings"

func EscapeString(str string) string {
	var builder strings.Builder
	for _, r := range str {
		switch r {
		case '\v':
			builder.WriteString("\\v")
		case '\a':
			builder.WriteString("\\a")
		case '\b':
			builder.WriteString("\\b")
		case '\f':
			builder.WriteString("\\f")
		case '\t':
			builder.WriteString("\\t")
		case '\r':
			builder.WriteString("\\r")
		case '\n':
			builder.WriteString("\\n")
		default:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func PadL(str string, padLen int, char rune) string {
	if padLen > len(str) {
		return strings.Repeat(string(char), padLen-len(str)) + str
	} else if padLen < len(str) {
		return str[:padLen]
	} else {
		return str
	}
}

func PadC(str string, padLen int, char rune) string {
	if padLen > len(str) {
		return strings.Repeat(string(char), (padLen-len(str))/2) + str + strings.Repeat(string(char), (padLen-len(str)+1)/2)
	} else if padLen < len(str) {
		padLR := (len(str) - padLen) / 2
		return str[padLR : padLR+padLen]
	} else {
		return str
	}
}

func ToTitle(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
