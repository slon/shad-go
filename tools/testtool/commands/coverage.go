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
	Enabled bool
	Percent float64
}

// getCoverageRequirements searches for comment in test files
// that specifies test coverage requirements.
//
// Stops on first matching comment.
func getCoverageRequirements(rootPackage string) *CoverageRequirements {
	files := listTestFiles(rootPackage)

	r := &CoverageRequirements{}
	for _, f := range files {
		r, _ := searchCoverageComment(f)
		if r.Enabled {
			return r
		}
	}

	return r
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
		percent, err := strconv.ParseFloat(t, 64)
		if err != nil {
			continue
		}
		if percent < 0 || percent > 100.0 {
			continue
		}
		return &CoverageRequirements{Enabled: true, Percent: percent}, nil
	}

	return &CoverageRequirements{}, nil
}

// calCoverage calculates coverage percent for given coverage profile.
func calCoverage(profile string) (float64, error) {
	profiles, err := cover.ParseProfiles(profile)
	if err != nil {
		return 0.0, fmt.Errorf("cannot parse coverage profile file %s: %w", profile, err)
	}

	var total, covered int
	for _, p := range profiles {
		for _, block := range p.Blocks {
			total += block.NumStmt
			if block.Count > 0 {
				covered += block.NumStmt
			}
		}
	}

	if total == 0 {
		return 0.0, nil
	}

	return float64(covered) / float64(total) * 100, nil
}
