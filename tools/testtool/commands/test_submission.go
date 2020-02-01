package commands

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

const (
	problemFlag     = "problem"
	studentRepoFlag = "student-repo"
	privateRepoFlag = "private-repo"

	testdataDir = "testdata"
)

var testSubmissionCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"check", "test-submission", "check-submission"},
	Short:   "test submission",
	Long:    `run solution on private and private tests`,
	Run: func(cmd *cobra.Command, args []string) {
		problem, err := cmd.Flags().GetString(problemFlag)
		if err != nil {
			log.Fatal(err)
		}

		studentRepo := mustParseDirFlag(studentRepoFlag, cmd)
		if !problemDirExists(studentRepo, problem) {
			log.Fatalf("%s does not have %s directory", studentRepo, problem)
		}

		privateRepo := mustParseDirFlag(privateRepoFlag, cmd)
		if !problemDirExists(privateRepo, problem) {
			log.Fatalf("%s does not have %s directory", privateRepo, problem)
		}

		testSubmission(studentRepo, privateRepo, problem)
	},
}

func init() {
	rootCmd.AddCommand(testSubmissionCmd)

	testSubmissionCmd.Flags().String(problemFlag, "", "problem directory name (required)")
	_ = testSubmissionCmd.MarkFlagRequired(problemFlag)

	testSubmissionCmd.Flags().String(studentRepoFlag, ".", "path to student repo root")
	testSubmissionCmd.Flags().String(privateRepoFlag, ".", "path to shad-go-private repo root")
}

// mustParseDirFlag parses string directory flag with given name.
//
// Exits on any error.
func mustParseDirFlag(name string, cmd *cobra.Command) string {
	dir, err := cmd.Flags().GetString(name)
	if err != nil {
		log.Fatal(err)
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// Check that repo dir contains problem subdir.
func problemDirExists(repo, problem string) bool {
	info, err := os.Stat(path.Join(repo, problem))
	if err != nil {
		return false
	}
	return info.IsDir()
}

func testSubmission(studentRepo, privateRepo, problem string) {
	// Create temp directory to store all files required to test the solution.
	tmpDir, err := ioutil.TempDir("/tmp", problem+"-")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()
	log.Printf("testing submission in %s", tmpDir)

	// Path to student's problem folder.
	studentProblem := path.Join(studentRepo, problem)
	// Path to private problem folder.
	privateProblem := path.Join(privateRepo, problem)

	// Copy submission files to temp dir.
	log.Printf("copying student solution")
	copyContents(studentProblem, tmpDir)

	// Copy tests from private repo to temp dir.
	log.Printf("copying tests")
	tests := listTestFiles(privateProblem)
	copyFiles(privateProblem, relPaths(privateProblem, tests), tmpDir)

	// Copy !change files from private repo to temp dir.
	log.Printf("copying !change files")
	protected := listProtectedFiles(privateProblem)
	copyFiles(privateProblem, relPaths(privateProblem, protected), tmpDir)

	// Copy testdata directory from private repo to temp dir.
	log.Printf("copying testdata directory")
	copyDir(path.Join(privateProblem, testdataDir), tmpDir)

	// Copy go.mod and go.sum from private repo to temp dir.
	log.Printf("copying go.mod and go.sum")
	copyFiles(privateRepo, []string{"go.mod", "go.sum"}, tmpDir)

	// Run tests.
	log.Printf("running tests")
	runTests(tmpDir)
}

// copyDir recursively copies src directory to dst.
func copyDir(src, dst string) {
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		return
	}

	cmd := exec.Command("rsync", "-r", src, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("directory copying failed: %s", err)
	}
}

// copyContents recursively copies src contents to dst.
func copyContents(src, dst string) {
	copyDir(src+"/", dst)
}

// copyFiles copies files preserving directory structure relative to baseDir.
//
// Existing files get replaced.
func copyFiles(baseDir string, relPaths []string, dst string) {
	for _, p := range relPaths {
		cmd := exec.Command("rsync", "-rR", p, dst)
		cmd.Dir = baseDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("file copying failed: %s", err)
		}
	}
}

// runTests runs all tests in directory with race detector.
func runTests(testDir string) {
	cmd := exec.Command("go", "test", "-v", "-mod", "readonly", "-tags", "private", "-race", "./...")
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	cmd.Dir = testDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("go test command failed: %s", err)
	}
}

// getPackageFiles returns absolute paths for all files in rootPackage and it's subpackages
// including tests and non-go files.
func getPackageFiles(rootPackage string, buildFlags []string) map[string]struct{} {
	cfg := &packages.Config{
		Dir:        rootPackage,
		Mode:       packages.NeedFiles,
		BuildFlags: buildFlags,
		Tests:      true,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatalf("unable to load packages %s: %s", rootPackage, err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	files := make(map[string]struct{})
	for _, p := range pkgs {
		for _, f := range p.GoFiles {
			files[f] = struct{}{}
		}
		for _, f := range p.OtherFiles {
			files[f] = struct{}{}
		}
	}

	return files
}

// relPaths converts paths to relative (to the baseDir) ones.
func relPaths(baseDir string, paths []string) []string {
	ret := make([]string, len(paths))
	for i, p := range paths {
		relPath, err := filepath.Rel(baseDir, p)
		if err != nil {
			log.Fatal(err)
		}
		ret[i] = relPath
	}
	return ret
}
