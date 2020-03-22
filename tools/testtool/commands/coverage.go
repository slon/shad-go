package commands

import (
	"fmt"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/cover"
)

// coverageCommentPrefix is a prefix of coverage comment.
//
// Coverage comment has the following form:
//
// // min coverage: 80.5%
const coverageCommentPrefix = "min coverage: "

type CoverageRequirements struct {
	Enabled  bool
	Percent  float64
	Packages []string
}

// getCoverageRequirements searches for comment in test files
// that specifies test coverage requirements.
//
// Stops on first matching comment.
func getCoverageRequirements(rootPackage string) *CoverageRequirements {
	files := listTestFiles(rootPackage)

	for _, f := range files {
		if r, _ := searchCoverageComment(f); r.Enabled {
			return r
		}
	}

	return &CoverageRequirements{}
}

// searchCoverageComment searches for the first occurrence of the comment of the form
//
// // min coverage: 80.5%
//
// Stops on the first matching comment.
func searchCoverageComment(fname string) (*CoverageRequirements, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, c := range f.Comments {
		t := c.Text()
		if !strings.HasPrefix(t, coverageCommentPrefix) || !strings.HasSuffix(t, "%\n") {
			continue
		}
		t = strings.TrimPrefix(t, coverageCommentPrefix)
		t = strings.TrimSuffix(t, "%\n")

		parts := strings.Split(t, " ")
		if len(parts) != 2 {
			continue
		}

		percent, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		if percent < 0 || percent > 100.0 {
			continue
		}

		return &CoverageRequirements{
			Enabled:  true,
			Percent:  percent,
			Packages: strings.Split(parts[0], ","),
		}, nil
	}

	return &CoverageRequirements{}, nil
}

// calCoverage calculates coverage percent for given coverage profile.
func calCoverage(fileNames []string) (float64, error) {
	type block struct {
		fileName            string
		startLine, startCol int
		endLine, endCol     int
		numStmt             int
	}
	counts := map[block]int{}

	for _, f := range fileNames {
		profiles, err := cover.ParseProfiles(f)
		if err != nil {
			return 0.0, fmt.Errorf("cannot parse coverage profile file %s: %w", f, err)
		}

		for _, p := range profiles {
			for _, b := range p.Blocks {
				counts[block{
					p.FileName,
					b.StartLine, b.StartCol,
					b.EndLine, b.EndCol,
					b.NumStmt,
				}] += b.Count
			}
		}
	}

	var total, covered int
	for b, count := range counts {
		total += b.numStmt
		if count > 0 {
			covered += b.numStmt
		}
	}

	if total == 0 {
		return 0.0, nil
	}

	return float64(covered) / float64(total) * 100, nil
}
