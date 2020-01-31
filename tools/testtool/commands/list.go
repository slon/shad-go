package commands

import (
	"sort"
	"strings"
)

// List all _test.go files in given directory including the ones with "private" build tag.
//
// Returns absolute paths.
func listTestFiles(rootPackage string) []string {
	files := getPackageFiles(rootPackage, []string{"-tags", "private"})
	var tests []string
	for f := range files {
		if strings.HasSuffix(f, "_test.go") {
			tests = append(tests, f)
		}
	}

	sort.Strings(tests)
	return tests
}

// List all .go source files in given directory protected by "!change" build tag.
//
// Returns absolute paths.
func listProtectedFiles(rootPackage string) []string {
	allFiles := getPackageFiles(rootPackage, nil)
	allFilesWithoutProtected := getPackageFiles(rootPackage, []string{"-tags", "change"})

	var protectedFiles []string
	for f := range allFiles {
		if _, ok := allFilesWithoutProtected[f]; !ok {
			protectedFiles = append(protectedFiles, f)
		}
	}

	sort.Strings(protectedFiles)
	return protectedFiles
}

func listPrivateFiles(rootPackage string) []string {
	allFiles := getPackageFiles(rootPackage, []string{})
	allWithPrivate := getPackageFiles(rootPackage, []string{"-tags", "private,solution"})

	var files []string
	for f := range allWithPrivate {
		if _, isPublic := allFiles[f]; !isPublic {
			files = append(files, f)
		}
	}

	sort.Strings(files)
	return files
}