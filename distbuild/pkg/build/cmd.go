package build

import (
	"fmt"
	"strings"
	"text/template"
)

type JobContext struct {
	SourceDir string
	OutputDir string
	Deps      map[ID]string
}

// Render replaces variable references with their real value.
func (c *Cmd) Render(ctx JobContext) (*Cmd, error) {
	var errs []error

	var fixedCtx struct {
		SourceDir string
		OutputDir string
		Deps      map[string]string
	}
	fixedCtx.SourceDir = ctx.SourceDir
	fixedCtx.OutputDir = ctx.OutputDir
	fixedCtx.Deps = map[string]string{}

	for k, v := range ctx.Deps {
		fixedCtx.Deps[k.String()] = v
	}

	render := func(str string) string {
		t, err := template.New("").Parse(str)
		if err != nil {
			errs = append(errs, err)
			return ""
		}

		var b strings.Builder
		if err := t.Execute(&b, fixedCtx); err != nil {
			errs = append(errs, err)
			return ""
		}

		return b.String()
	}

	renderList := func(l []string) []string {
		var result []string
		for _, in := range l {
			result = append(result, render(in))
		}
		return result
	}

	var rendered Cmd

	rendered.CatOutput = render(c.CatOutput)
	rendered.CatTemplate = render(c.CatTemplate)
	rendered.WorkingDirectory = render(c.WorkingDirectory)
	rendered.Exec = renderList(c.Exec)
	rendered.Environ = renderList(c.Environ)

	if len(errs) != 0 {
		return nil, fmt.Errorf("error rendering cmd: %w", errs[0])
	}

	return &rendered, nil
}
