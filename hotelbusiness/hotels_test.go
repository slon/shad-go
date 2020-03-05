package hotelbusiness

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeLoad_basic(t *testing.T) {
	for _, tc := range []struct {
		title  string
		guests []Guest
		result []Load
	}{
		{
			title:  "empty input",
			guests: []Guest{},
			result: []Load{},
		},
		{
			title:  "one guest",
			guests: []Guest{{1, 2}},
			result: []Load{{1, 1}, {2, 0}},
		},
		{
			title:  "two guests, one-by-one, without any gaps",
			guests: []Guest{{1, 2}, {2, 3}},
			result: []Load{{1, 1}, {3, 0}},
		},
		{
			title:  "two guests, one-by-one, with a gap",
			guests: []Guest{{1, 2}, {3, 4}},
			result: []Load{{1, 1}, {2, 0}, {3, 1}, {4, 0}},
		},
		{
			title:  "two guests, together",
			guests: []Guest{{1, 2}, {1, 2}},
			result: []Load{{1, 2}, {2, 0}},
		},
		{
			title:  "overlapping",
			guests: []Guest{{1, 3}, {3, 5}, {2, 4}},
			result: []Load{{1, 1}, {2, 2}, {4, 1}, {5, 0}},
		},
		{
			title:  "stairs",
			guests: []Guest{{1, 6}, {2, 5}, {3, 4}},
			result: []Load{{1, 1}, {2, 2}, {3, 3}, {4, 2}, {5, 1}, {6, 0}},
		},
		{
			title:  "starting late",
			guests: []Guest{{3, 7}, {5, 7}},
			result: []Load{{3, 1}, {5, 2}, {7, 0}},
		},
		{
			title:  "unordered",
			guests: []Guest{{4, 7}, {2, 4}, {2, 3}},
			result: []Load{{2, 2}, {3, 1}, {7, 0}},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			require.Equal(t, tc.result, ComputeLoad(tc.guests))
		})
	}
}

func TestComputeLoad_stress1(t *testing.T) {
	n := 1000000
	g := make([]Guest, 0, 1000000)
	for i := 0; i < n; i++ {
		g = append(g, Guest{1, 2})
	}
	l := ComputeLoad(g)
	require.Equal(t, []Load{{1, n}, {2, 0}}, l)
}

func TestComputeLoad_stress2(t *testing.T) {
	n := 1000000
	g := make([]Guest, 0, 1000000)
	for i := 0; i < n; i++ {
		g = append(g, Guest{i, i + 1})
	}
	l := ComputeLoad(g)
	require.Equal(t, []Load{{0, 1}, {n, 0}}, l)
}
