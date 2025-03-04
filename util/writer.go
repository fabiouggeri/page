package util

import (
	"fmt"
	"strings"
)

type TextWriter interface {
	WriteF(format string, args ...any) TextWriter
	WriteString(s string) TextWriter
	WriteRune(r rune) TextWriter
	Indent(i int) TextWriter
	NewLine() TextWriter
	HasError() bool
	Error() error
	CleanError()
	String() string
}

type StringCodeWriter struct {
	text        strings.Builder
	indentation int
	newLine     bool
	err         error
}

var _ TextWriter = &StringCodeWriter{}

func NewStringTextWriter() *StringCodeWriter {
	return &StringCodeWriter{indentation: 0, newLine: true, err: nil}
}

func (w *StringCodeWriter) indentIfNewLine() bool {
	if w.err != nil {
		return false
	} else if w.newLine {
		_, w.err = w.text.WriteString(strings.Repeat(" ", w.indentation))
		w.newLine = false
	}
	return w.err == nil
}

func (w *StringCodeWriter) String() string {
	return w.text.String()
}

func (w *StringCodeWriter) HasError() bool {
	return w.err != nil
}

func (w *StringCodeWriter) Error() error {
	return w.err
}

func (w *StringCodeWriter) CleanError() {
	w.err = nil
}

func (w *StringCodeWriter) WriteF(format string, args ...any) TextWriter {
	if !w.indentIfNewLine() {
		return w
	}
	_, w.err = w.text.WriteString(fmt.Sprintf(format, args...))
	return w
}

func (w *StringCodeWriter) WriteString(s string) TextWriter {
	if !w.indentIfNewLine() {
		return w
	}
	_, w.err = w.text.WriteString(s)
	return w
}

func (w *StringCodeWriter) WriteRune(r rune) TextWriter {
	if !w.indentIfNewLine() {
		return w
	}
	_, w.err = w.text.WriteRune(r)
	return w
}

func (w *StringCodeWriter) Indent(i int) TextWriter {
	w.indentation += i
	if w.indentation < 0 {
		w.indentation = 0
	}
	return w
}

func (w *StringCodeWriter) NewLine() TextWriter {
	if w.err != nil {
		return w
	}
	_, w.err = w.text.WriteRune('\n')
	w.newLine = true
	return w
}
