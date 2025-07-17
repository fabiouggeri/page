package rule

type RuleOption struct {
	code          int
	name          string
	parameterized bool
	mandatory     bool
}

var (
	MAIN         *RuleOption = &RuleOption{code: 0x0001, name: "Main", parameterized: false, mandatory: false}
	TOKEN        *RuleOption = &RuleOption{code: 0x0002, name: "Token", parameterized: false, mandatory: false}
	ATOMIC       *RuleOption = &RuleOption{code: 0x0004, name: "Atomic", parameterized: false, mandatory: false}
	SKIP_NODE    *RuleOption = &RuleOption{code: 0x0010, name: "SkipNode", parameterized: false, mandatory: false}
	FRAGMENT     *RuleOption = &RuleOption{code: 0x0020, name: "Fragment", parameterized: false, mandatory: false}
	NAME         *RuleOption = &RuleOption{code: 0x0040, name: "Name", parameterized: true, mandatory: true}
	MEMOIZE      *RuleOption = &RuleOption{code: 0x0080, name: "Memoize", parameterized: false, mandatory: false}
	IGNORE       *RuleOption = &RuleOption{code: 0x0100, name: "Ignore", parameterized: false, mandatory: false}
	START_LINE   *RuleOption = &RuleOption{code: 0x0200, name: "StartLine", parameterized: false, mandatory: false}
	ONLY_IGNORED *RuleOption = &RuleOption{code: 0x0400, name: "OnlyIgnoredInLine", parameterized: false, mandatory: false}
)

var AllOptions = []*RuleOption{
	MAIN,
	TOKEN,
	ATOMIC,
	SKIP_NODE,
	FRAGMENT,
	NAME,
	MEMOIZE,
	IGNORE,
	START_LINE,
	ONLY_IGNORED,
}

func (o *RuleOption) Code() int {
	return o.code
}

func (o *RuleOption) String() string {
	return o.name
}

func (o *RuleOption) Name() string {
	return o.name
}

func (o *RuleOption) Parameterized() bool {
	return o.parameterized
}

func (o *RuleOption) ParameterMandatory() bool {
	return o.mandatory
}
