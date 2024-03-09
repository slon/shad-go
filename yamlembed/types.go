package yamlembed

type Foo struct {
	A string
	p int64
}

type Bar struct {
	I      int64
	B      string
	UpperB string
	OI     []string
	F      []any
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}
