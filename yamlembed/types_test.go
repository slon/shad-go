package yamlembed

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestFoo_unmarshal(t *testing.T) {
	data := []byte(`
aa: hello
p: 42
`)

	var foo Foo
	require.NoError(t, yaml.Unmarshal(data, &foo))
	require.Equal(t, Foo{A: "hello"}, foo)
}

func TestFoo_marshal(t *testing.T) {
	foo := Foo{
		A: "hello",
		p: 42,
	}

	expected := `aa: hello
`

	data, err := yaml.Marshal(foo)
	require.NoError(t, err)
	require.Equal(t, expected, string(data))
}

func TestBar_unmarshal(t *testing.T) {
	data := []byte(`
i: 29
b: world
oi:
- pool
- tree
f:
- 3
- data
`)

	expected := Bar{
		B:      "world",
		UpperB: "WORLD",
		OI:     []string{"pool", "tree"},
		F:      []any{3, "data"},
	}

	var bar Bar
	require.NoError(t, yaml.Unmarshal(data, &bar))
	require.Equal(t, expected, bar)
}

func TestBar_marshal(t *testing.T) {
	bar := Bar{
		B:      "world",
		UpperB: "WORLD",
		OI:     []string{},
		F:      []any{3, "data"},
	}

	expected := `b: world
f: [3, data]
`

	data, err := yaml.Marshal(bar)
	require.NoError(t, err)
	require.Equal(t, expected, string(data))
}

func TestBaz_unmarshal(t *testing.T) {
	data := []byte(`
aa: hello
p: 42
i: 29
b: world
oi:
- pool
- tree
f:
- 3
- data
`)

	var baz Baz
	require.NoError(t, yaml.Unmarshal(data, &baz))
	require.Equal(t, Baz{
		Foo{
			A: "hello",
		},
		Bar{
			B:      "world",
			UpperB: "WORLD",
			OI:     []string{"pool", "tree"},
			F:      []any{3, "data"},
		},
	}, baz)
}

func TestBaz_marshal(t *testing.T) {
	baz := Baz{
		Foo: Foo{
			A: "hello",
		},
		Bar: Bar{
			B:      "world",
			UpperB: "WORLD",
			OI:     []string{},
			F:      []any{3, "data"},
		},
	}

	expected := `aa: hello
b: world
f: [3, data]
`

	data, err := yaml.Marshal(baz)
	require.NoError(t, err)
	require.Equal(t, expected, string(data))
}
