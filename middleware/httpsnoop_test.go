package middlware

import (
	// without that line, go mod tidy might remove httpsnoop in student repository
	_ "github.com/felixge/httpsnoop"
)
