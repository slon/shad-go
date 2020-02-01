package commands

import (
	"sort"
	"strings"
)

// listTestFiles returns absolute paths for all _test.go files of the package
// including the ones with "private" build tag.
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

// listProtectedFiles returns absolute paths for all files of the package
// protected by "!change" build tag.
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

// listPrivateFiles returns absolute paths for all files of the package
// protected by "private,solution" build tag.
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
