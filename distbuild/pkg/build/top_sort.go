package build

// TopSort sorts jobs in topological order assuming dependency graph contains no cycles.
func TopSort(jobs []Job) []Job {
	var sorted []Job
	visited := make([]bool, len(jobs))

	jobIDIndex := map[ID]int{}
	for i, j := range jobs {
		jobIDIndex[j.ID] = i
	}

	var visit func(jobIndex int)
	visit = func(jobIndex int) {
		if visited[jobIndex] {
			return
		}

		visited[jobIndex] = true
		for _, dep := range jobs[jobIndex].Deps {
			visit(jobIDIndex[dep])
		}
		sorted = append(sorted, jobs[jobIndex])
	}

	for i := range jobs {
		visit(i)
	}

	return sorted
}
