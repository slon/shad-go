package size

import "testing"

func TestSize(t *testing.T) {
	type Test struct {
		in  int
		out string
	}

	var tests = []Test{
		{-1, "negative"},
		{5, "small"},
	}

	for i, test := range tests {
		size := Size(test.in)
		if size != test.out {
			t.Errorf("#%d: Size(%d)=%s; want %s", i, test.in, size, test.out)
		}
	}
}
