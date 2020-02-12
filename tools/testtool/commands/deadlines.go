package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const timeFormat = "02-01-2006 15:04"

type (
	Task struct {
		Name  string `yaml:"task"`
		Score int    `yaml:"score"`
	}

	Group struct {
		Name     string `yaml:"group"`
		Start    string `yaml:"start"`
		Deadline string `yaml:"deadline"`
		Tasks    []Task `yaml:"tasks"`
	}

	Deadlines []Group
)

func (g Group) IsOpen() bool {
	t, _ := time.Parse(timeFormat, g.Start)
	return time.Until(t) < 0
}

func (d Deadlines) FindTask(name string) (*Group, *Task) {
	for _, g := range d {
		for _, t := range g.Tasks {
			if t.Name == name {
				return &g, &t
			}
		}
	}

	return nil, nil
}

func loadDeadlines(filename string) (Deadlines, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var d Deadlines
	if err := yaml.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("error reading deadlines: %w", err)
	}

	for _, g := range d {
		if _, err := time.Parse(timeFormat, g.Start); err != nil {
			return nil, fmt.Errorf("invalid time format in task %q: %w", g.Name, err)
		}

		if _, err := time.Parse(timeFormat, g.Deadline); err != nil {
			return nil, fmt.Errorf("invalid time format in task %q: %w", g.Name, err)
		}
	}

	return d, nil
}

func findChangedTasks(d Deadlines, files []string) []string {
	tasks := map[string]struct{}{}

	for _, f := range files {
		components := strings.Split(f, string(filepath.Separator))
		if len(components) == 0 {
			continue
		}

		_, task := d.FindTask(components[0])
		if task == nil {
			continue
		}

		tasks[task.Name] = struct{}{}
	}

	var l []string
	for t := range tasks {
		l = append(l, t)
	}

	sort.Strings(l)
	return l
}
