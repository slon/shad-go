package commands

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	problemFlag     = "problem"
	studentRepoFlag = "student-repo"
	privateRepoFlag = "private-repo"

	testdataDir      = "testdata"
	moduleImportPath = "gitlab.com/slon/shad-go"
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
	tmpRepo, err := ioutil.TempDir("/tmp", problem+"-")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Chmod(tmpRepo, 0755); err != nil {
		log.Fatal(err)
	}

	defer func() { _ = os.RemoveAll(tmpRepo) }()
	log.Printf("testing submission in %s", tmpRepo)

	// Path to private problem folder.
	privateProblem := path.Join(privateRepo, problem)

	// Copy student repo files to temp dir.
	log.Printf("copying student repo")
	copyContents(studentRepo, ".", tmpRepo)

	// Copy tests from private repo to temp dir.
	log.Printf("copying tests")
	tests := listTestFiles(privateProblem)
	copyFiles(privateRepo, relPaths(privateRepo, tests), tmpRepo)

	// Copy !change files from private repo to temp dir.
	log.Printf("copying !change files")
	protected := listProtectedFiles(privateProblem)
	copyFiles(privateRepo, relPaths(privateRepo, protected), tmpRepo)

	// Copy testdata directory from private repo to temp dir.
	log.Printf("copying testdata directory")
	copyDir(privateRepo, path.Join(problem, testdataDir), tmpRepo)

	// Copy go.mod and go.sum from private repo to temp dir.
	log.Printf("copying go.mod and go.sum")
	copyFiles(privateRepo, []string{"go.mod", "go.sum"}, tmpRepo)

	// Run tests.
	log.Printf("running tests")
	runTests(tmpRepo, problem)
}

// copyDir recursively copies src directory to dst.
func copyDir(baseDir, src, dst string) {
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		return
	}

	cmd := exec.Command("rsync", "-prR", src, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = baseDir

	if err := cmd.Run(); err != nil {
		log.Fatalf("directory copying failed: %s", err)
	}
}

// copyContents recursively copies src contents to dst.
func copyContents(baseDir, src, dst string) {
	copyDir(baseDir, src+"/", dst)
}

// copyFiles copies files preserving directory structure relative to baseDir.
//
// Existing files get replaced.
func copyFiles(baseDir string, relPaths []string, dst string) {
	for _, p := range relPaths {
		cmd := exec.Command("rsync", "-prR", p, dst)
		cmd.Dir = baseDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("file copying failed: %s", err)
		}
	}
}

func randomName() string {
	var raw [8]byte
	_, _ = rand.Read(raw[:])
	return hex.EncodeToString(raw[:])
}

// runTests runs all tests in directory with race detector.
func runTests(testDir, problem string) {
	binCache, err := ioutil.TempDir("/tmp", "bincache")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Chmod(binCache, 0755); err != nil {
		log.Fatal(err)
	}

	runGo := func(arg ...string) {
		log.Printf("> go %s", strings.Join(arg, " "))

		cmd := exec.Command("go", arg...)
		cmd.Env = append(os.Environ(), "GOFLAGS=")
		cmd.Dir = testDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	binaries := map[string]string{}
	testBinaries := map[string]string{}

	binPkgs, testPkgs := listTestsAndBinaries(testDir, []string{"-tags", "private", "-mod", "readonly"})
	for binaryPkg := range binPkgs {
		binPath := filepath.Join(binCache, randomName())
		binaries[binaryPkg] = binPath
		runGo("build", "-mod", "readonly", "-tags", "private", "-o", binPath, binaryPkg)
	}

	binariesJSON, _ := json.Marshal(binPkgs)

	for testPkg := range testPkgs {
		binPath := filepath.Join(binCache, randomName())
		testBinaries[testPkg] = binPath
		runGo("test", "-mod", "readonly", "-tags", "private", "-c", "-o", binPath, testPkg)
	}

	for testPkg, testBinary := range testBinaries {
		relPath := strings.TrimPrefix(testPkg, moduleImportPath)

		cmd := exec.Command(testBinary)
		if currentUserIsRoot() {
			if err := sandbox(cmd); err != nil {
				log.Fatal(err)
			}
		}

		cmd.Dir = filepath.Join(testDir, relPath)
		cmd.Env = []string{"TESTTOOL_BINARIES=" + string(binariesJSON)}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
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
