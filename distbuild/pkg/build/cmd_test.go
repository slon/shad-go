package build

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdRender(t *testing.T) {
	tmpl := Cmd{
		CatOutput:   "{{.OutputDir}}/import.map",
		CatTemplate: `bytes={{index .Deps "6100000000000000000000000000000000000000"}}/lib.a`,
	}

	ctx := JobContext{
		OutputDir: "/distbuild/jobs/b",
		Deps: map[ID]string{
			{'a'}: "/distbuild/jobs/a",
		},
	}

	result, err := tmpl.Render(ctx)
	require.NoError(t, err)

	expected := &Cmd{
		CatOutput:   "/distbuild/jobs/b/import.map",
		CatTemplate: "bytes=/distbuild/jobs/a/lib.a",
	}

	require.Equal(t, expected, result)
}
