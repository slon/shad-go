package testtool

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// CheckForbiddenImport checks that the project does not use forbidden imports.
func CheckForbiddenImport(t *testing.T, forbiddenPackage string) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("Cannot get current file path")
	}
	srcDir := filepath.Dir(filename)

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(importPath, forbiddenPackage) {
				position := fset.Position(imp.Path.Pos())
				t.Errorf("Forbidden %s package import found in %s at line %d", forbiddenPackage, path, position.Line)
			}
		}

		return nil
	})

	if err != nil {
		t.Errorf("Failed to walk through the source directory: %v", err)
	}
}
