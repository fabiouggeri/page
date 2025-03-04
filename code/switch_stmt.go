package code

import "github.com/fabiouggeri/page/util"

type SwitchOption struct {
	option     Expression
	statements []Code
}

type Switch struct {
	test          Expression
	options       []*SwitchOption
	defaultOption []Code
}

var _ Code = &Switch{}

func newSwitch(test Expression) *Switch {
	return &Switch{test: test, options: make([]*SwitchOption, 0), defaultOption: make([]Code, 0)}
}

func (s *Switch) Option(label Expression, statements ...Code) *Switch {
	s.options = append(s.options, &SwitchOption{option: label, statements: statements})
	return s
}

func (s *Switch) Default(statements ...Code) *Switch {
	s.defaultOption = append(s.defaultOption, statements...)
	return s
}

func (s *Switch) GetTest() Expression {
	return s.test
}

func (s *Switch) GetOptions() []*SwitchOption {
	return s.options
}

func (s *Switch) GetDefault() []Code {
	return s.defaultOption
}

func (s *Switch) Generate(generator CodeGenerator, str util.TextWriter) error {
	return generator.GenerateSwitch(s, str)
}

func (s *Switch) IsEmpty() bool {
	return false
}

func (s *Switch) String() string {
	return "switch"
}

func (o *SwitchOption) GetOption() Expression {
	return o.option
}

func (o *SwitchOption) GetStatements() []Code {
	return o.statements
}
