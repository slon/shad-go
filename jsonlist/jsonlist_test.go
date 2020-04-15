package jsonlist

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type S struct {
	A string
	B int
	C interface{}
}

func TestJsonList(t *testing.T) {
	for _, test := range []struct {
		js    string
		value interface{}
	}{
		{
			js:    `1 2 3`,
			value: []int{1, 2, 3},
		},
		{
			js:    `"A" "B" "C"`,
			value: []string{"A", "B", "C"},
		},
		{
			js:    `{"A": "A"} {"B": 2} {"C": 0.5}`,
			value: []S{{A: "A"}, {B: 2}, {C: 0.5}},
		},
		{
			js:    `"A" 2`,
			value: []interface{}{"A", 2.0},
		},
	} {
		t.Run(test.js, func(t *testing.T) {
			emptySlice := reflect.New(reflect.TypeOf(test.value))

			require.NoError(t, Unmarshal(bytes.NewBufferString(test.js), emptySlice.Interface()))
			require.Equal(t, test.value, emptySlice.Elem().Interface())

			var buf bytes.Buffer
			require.NoError(t, Marshal(&buf, test.value))

			emptySlice = reflect.New(reflect.TypeOf(test.value))

			require.NoError(t, Unmarshal(bytes.NewBufferString(test.js), emptySlice.Interface()))
			require.Equal(t, test.value, emptySlice.Elem().Interface())
		})
	}
}
