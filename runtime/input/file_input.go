package input

import (
	"bufio"
	"os"
)

type FileInput struct {
	file   *os.File
	reader *bufio.Reader
	index  int
	buffer []rune
	eof    bool
}

const (
	initialBufferSize = 8 * 1024
)

var _ Input = &FileInput{}

// NewFileInput creates a new FileInput.
func NewFileInput(filePathName string) (*FileInput, error) {
	return NewFileInputSize(filePathName, initialBufferSize)
}

// NewFileInputSize creates a new FileInput with initial file buffer capacity.
func NewFileInputSize(filePathName string, initialFileBufferSize int) (*FileInput, error) {
	var bufferCapacity int
	file, err := os.OpenFile(filePathName, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		bufferCapacity = initialFileBufferSize
	} else {
		bufferCapacity = int(fi.Size())
	}
	return &FileInput{
		file:   file,
		reader: bufio.NewReaderSize(file, initialFileBufferSize),
		index:  0,
		buffer: make([]rune, 0, bufferCapacity),
		eof:    false,
	}, nil
}

func (f *FileInput) Eof() bool {
	return f.eof
}

func (f *FileInput) GetChar() rune {
	if !f.readChar() {
		return '\x00'
	}
	return f.buffer[f.index]
}

func (f *FileInput) readChar() bool {
	if f.eof {
		return false
	}
	if f.index >= len(f.buffer) {
		c, _, err := f.reader.ReadRune()
		if err != nil {
			f.eof = true
			return false
		}
		f.buffer = append(f.buffer, c)
	}
	return true
}

func (f *FileInput) Index() int {
	return f.index
}

func (i *FileInput) SetIndex(index int) {
	if index > i.index {
		for i.index < index {
			if !i.Skip() {
				return
			}
		}
	} else if index >= 0 {
		i.index = index
		return
	}
}

func (f *FileInput) Skip() bool {
	if f.eof {
		return false
	}
	f.index++
	if !f.readChar() {
		f.index--
		return false
	}
	return true
}

func (f *FileInput) Close() {
	if f.file != nil {
		f.file.Close()
	}
	f.file = nil
	f.reader = nil
}

func (f *FileInput) GetText(start int, end int) string {
	for !f.Eof() && len(f.buffer) < end {
		f.Skip()
	}
	if end > len(f.buffer) {
		return ""
	}
	return string(f.buffer[start:end])
}
