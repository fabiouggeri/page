package source

import (
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type BufferdSource struct {
	index  uint32
	buffer []byte
}

var (
	defaultBufferSize uint32 = 8192
)

func FromBuffer(source []byte) Source {
	s := &BufferdSource{
		index:  0,
		buffer: source}
	return s
}

func FromString(source string) Source {
	return FromBuffer([]byte(source))
}

func FromFile(filePathname string) (Source, error) {
	handle, err := os.Open(filePathname)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	data, err := io.ReadAll(handle)
	if err != nil {
		return nil, err
	}
	return FromBuffer(data), nil
}

func SetDefaultBufferSize(size uint32) {
	defaultBufferSize = size
}

func GetDefaultBufferSize() uint32 {
	return defaultBufferSize
}

func getRune(buffer []byte, index uint32) (rune, uint32) {
	r, size := rune(buffer[index]), 1
	if r >= utf8.RuneSelf {
		r, size = utf8.DecodeRune(buffer[index:])
	}
	return r, uint32(size)
}

func (s *BufferdSource) StringAt(start, end uint32) string {
	if end >= start && end < uint32(len(s.buffer)) {
		return string(s.buffer[start:end])
	}
	return ""
}

func (s *BufferdSource) Match(character rune) bool {
	if s.index < uint32(len(s.buffer)) {
		r, size := getRune(s.buffer, s.index)
		if r == character {
			s.index += size
			return true

		}
	}
	return false
}

func (s *BufferdSource) MatchIgnoreCase(character rune) bool {
	if s.index < uint32(len(s.buffer)) {
		r, size := getRune(s.buffer, s.index)
		if unicode.ToUpper(r) == unicode.ToUpper(character) {
			s.index += size
			return true
		}
	}
	return false
}

func (s *BufferdSource) MatchRange(start, end rune) bool {
	if s.index < uint32(len(s.buffer)) {
		r, size := getRune(s.buffer, s.index)
		if r >= start && r <= end {
			s.index += size
			return true

		}
	}
	return false
}

func (s *BufferdSource) MatchString(text string) bool {
	if s.index < uint32(len(s.buffer)) {
		strBuffer := string(s.buffer[s.index:])
		if len(strBuffer) >= len(text) && strBuffer[len(strBuffer)-len(text):] == text {
			s.index += uint32(len(text))
			return true
		}
	}
	return false
}

func (s *BufferdSource) MatchStringIgnoreCase(text string) bool {
	if s.index < uint32(len(s.buffer)) {
		strBuffer := string(s.buffer[s.index:])
		if len(strBuffer) >= len(text) && strings.EqualFold(strBuffer[len(strBuffer)-len(text):], text) {
			s.index += uint32(len(text))
			return true
		}
	}
	return false
}

func (s *BufferdSource) EOI() bool {
	return s.index >= uint32(len(s.buffer))
}

func (s *BufferdSource) Index() uint32 {
	return s.index
}

func (s *BufferdSource) SetIndex(index uint32) {
	if index < uint32(len(s.buffer)) {
		s.index = index
	} else {
		s.index = uint32(len(s.buffer))
	}
}
