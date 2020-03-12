// +build private

package brokentest

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// Read tests from csv file.
func readTestCases(filename string) ([]*testCase, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var tests []*testCase
	for _, r := range records {
		a, _ := strconv.ParseInt(r[0], 10, 64)
		b, _ := strconv.ParseInt(r[1], 10, 64)
		sum, _ := strconv.ParseInt(r[2], 10, 64)
		tests = append(tests, &testCase{a: a, b: b, sum: sum})
	}

	return tests, nil
}

func TestSumPrivate(t *testing.T) {
	tests, err := readTestCases("./testdata/tests.csv")
	require.NoError(t, err)

	for _, tc := range tests {
		s := Sum(tc.a, tc.b)
		require.Equal(t, tc.sum, s, "%d + %d == %d != %d", tc.a, tc.b, s, tc.sum)
	}
}
