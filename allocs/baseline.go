// +build !solution
// +build !change

package allocs

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
)

type Counter interface {
	Count(r io.Reader) error
	String() string
}

type BaselineCounter struct {
	counts map[string]int
}

func NewBaselineCounter() Counter {
	return BaselineCounter{counts: map[string]int{}}
}

func (c BaselineCounter) Count(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	dataStr := string(data)
	for _, line := range strings.Split(dataStr, "\n") {
		for _, word := range strings.Split(line, " ") {
			c.counts[word]++
		}
	}
	return nil
}

func (c BaselineCounter) String() string {
	keys := make([]string, 0, 0)
	for word := range c.counts {
		keys = append(keys, word)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	result := ""
	for _, key := range keys {
		line := fmt.Sprintf("word '%s' has %d occurrences\n", key, c.counts[key])
		result += line
	}
	return result
}
