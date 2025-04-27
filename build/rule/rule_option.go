package rule

type RuleOption struct {
	name          string
	parameterized bool
	mandatory     bool
}

var (
	MAIN      *RuleOption = &RuleOption{name: "Main", parameterized: false, mandatory: false}
	TOKEN     *RuleOption = &RuleOption{name: "Token", parameterized: false, mandatory: false}
	ATOMIC    *RuleOption = &RuleOption{name: "Atomic", parameterized: false, mandatory: false}
	SKIP_NODE *RuleOption = &RuleOption{name: "SkipNode", parameterized: false, mandatory: false}
	FRAGMENT  *RuleOption = &RuleOption{name: "Fragment", parameterized: false, mandatory: false}
	NAME      *RuleOption = &RuleOption{name: "Name", parameterized: true, mandatory: true}
	MEMOIZE   *RuleOption = &RuleOption{name: "Memoize", parameterized: false, mandatory: false}
	IGNORE    *RuleOption = &RuleOption{name: "Ignore", parameterized: false, mandatory: false}
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
