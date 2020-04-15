package reversemap

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseMap(t *testing.T) {
	data := []struct {
		forward  interface{}
		backward interface{}
	}{
		{
			forward: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			backward: map[string]string{
				"v1": "k1",
				"v2": "k2",
				"v3": "k3",
			},
		},
		{
			forward: map[int]string{
				1: "v1",
				2: "v2",
				3: "v3",
			},
			backward: map[string]int{
				"v1": 1,
				"v2": 2,
				"v3": 3,
			},
		},
		{
			forward: map[int]int{
				1: 4,
				2: 5,
				3: 6,
			},
			backward: map[int]int{
				4: 1,
				5: 2,
				6: 3,
			},
		},
	}
	for _, d := range data {
		t.Run(reflect.TypeOf(d.forward).String(), func(t *testing.T) {
			assert.Equal(t, d.backward, ReverseMap(d.forward))
			assert.Equal(t, d.forward, ReverseMap(d.backward))
		})
	}
}

func TestReverseInt(t *testing.T) {
	assert.Panics(t, func() {
		ReverseMap(new(int))
	})
}
