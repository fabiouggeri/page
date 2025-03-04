package code

import "github.com/fabiouggeri/page/util"

type Comment struct {
	comment string
}

var _ Code = &Comment{}

func newComment(comment string) *Comment {
	return &Comment{comment: comment}
}

func (c *Comment) Comment() string {
	return c.comment
}

func (c *Comment) IsEmpty() bool {
	return c.comment == ""
}

func (c *Comment) String() string {
	return "/*" + c.comment + "*/"
}

func (c *Comment) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateComment(c, str)
}
