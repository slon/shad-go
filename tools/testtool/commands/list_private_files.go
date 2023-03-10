package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listPrivateFilesCmd = &cobra.Command{
	Use:   "list-private-files",
	Short: "list private files",
	Run:   runListPrivateFiles,
}

func init() {
	rootCmd.AddCommand(listPrivateFilesCmd)
}

func doListPrivateFiles() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	privateFiles := listPrivateFiles(".")
	for _, f := range privateFiles {
		rel, err := filepath.Rel(cwd, f)
		if err != nil {
			return err
		}

		fmt.Println(rel)
	}

	return nil
}

func runListPrivateFiles(cmd *cobra.Command, args []string) {
	if err := doListPrivateFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "testtool: %v\n", err)
		os.Exit(1)
	}
}
