package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	Task struct {
		Name  string   `yaml:"task"`
		Watch []string `yaml:"watch"`
	}

	Group struct {
		Name  string `yaml:"group"`
		Tasks []Task `yaml:"tasks"`
	}

	Deadlines []Group
)

func (d Deadlines) Tasks() []*Task {
	var tasks []*Task
	for _, g := range d {
		for i := range g.Tasks {
			tasks = append(tasks, &g.Tasks[i])
		}
	}
	return tasks
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
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var m struct {
		Deadlines struct {
			Schedule Deadlines `yaml:"schedule"`
		} `yaml:"deadlines"`
	}

	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("error reading deadlines: %w", err)
	}

	d := m.Deadlines.Schedule
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
		if task != nil {
			tasks[task.Name] = struct{}{}
		}

		for _, task := range d.Tasks() {
			for _, path := range task.Watch {
				if strings.HasPrefix(f, path) {
					tasks[task.Name] = struct{}{}
				}
			}
		}
	}

	var l []string
	for t := range tasks {
		l = append(l, t)
	}

	sort.Strings(l)
	return l
}
