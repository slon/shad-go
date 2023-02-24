package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	description string
	input       []string
	expected    []int
	error       bool
}

var testCases = []testCase{
	{
		description: "push numbers",
		input:       []string{"1 2 3 4 5"},
		expected:    []int{1, 2, 3, 4, 5},
	},
	{
		description: "add",
		input:       []string{"1 2 +"},
		expected:    []int{3},
	},
	{
		description: "nothing to add",
		input:       []string{"+"},
		error:       true,
	},
	{
		description: "add arity",
		input:       []string{"1 +"},
		error:       true,
	},
	{
		description: "sub",
		input:       []string{"3 4 -"},
		expected:    []int{-1},
	},
	{
		description: "nothing to sub",
		input:       []string{"-"},
		error:       true,
	},
	{
		description: "sub arity",
		input:       []string{"1 -"},
		error:       true,
	},
	{
		description: "mul",
		input:       []string{"2 4 *"},
		expected:    []int{8},
	},
	{
		description: "nothing to mul",
		input:       []string{"*"},
		error:       true,
	},
	{
		description: "mul arity",
		input:       []string{"1 *"},
		error:       true,
	},
	{
		description: "div",
		input:       []string{"12 3 /"},
		expected:    []int{4},
	},
	{
		description: "integer division",
		input:       []string{"8 3 /"},
		expected:    []int{2},
	},
	{
		description: "division by zero",
		input:       []string{"4 0 /"},
		error:       true,
	},
	{
		description: "nothing to div",
		input:       []string{"/"},
		error:       true,
	},
	{
		description: "div arity",
		input:       []string{"1 /"},
		error:       true,
	},
	{
		description: "add sub",
		input:       []string{"1 2 + 4 -"},
		expected:    []int{-1},
	},
	{
		description: "mul div",
		input:       []string{"2 4 * 3 /"},
		expected:    []int{2},
	},
	{
		description: "dup",
		input:       []string{"1 dup"},
		expected:    []int{1, 1},
	},
	{
		description: "dup top",
		input:       []string{"1 2 dup"},
		expected:    []int{1, 2, 2},
	},
	{
		description: "nothing to dup",
		input:       []string{"dup"},
		error:       true,
	},
	{
		description: "drop",
		input:       []string{"1 drop"},
		expected:    []int{},
	},
	{
		description: "drop top",
		input:       []string{"1 2 drop"},
		expected:    []int{1},
	},
	{
		description: "nothing to drop",
		input:       []string{"drop"},
		error:       true,
	},
	{
		description: "swap",
		input:       []string{"1 2 swap"},
		expected:    []int{2, 1},
	},
	{
		description: "swap top",
		input:       []string{"1 2 3 swap"},
		expected:    []int{1, 3, 2},
	},
	{
		description: "nothing to swap",
		input:       []string{"swap"},
		error:       true,
	},
	{
		description: "swap arity",
		input:       []string{"1 swap"},
		error:       true,
	},
	{
		description: "over",
		input:       []string{"1 2 over"},
		expected:    []int{1, 2, 1},
	},
	{
		description: "over2",
		input:       []string{"1 2 3 over"},
		expected:    []int{1, 2, 3, 2},
	},
	{
		description: "nothing to over",
		input:       []string{"over"},
		error:       true,
	},
	{
		description: "over arity",
		input:       []string{"1 over"},
		error:       true,
	},
	{
		description: "user-defined",
		input:       []string{": dup-twice dup dup ;", "1 dup-twice"},
		expected:    []int{1, 1, 1},
	},
	{
		description: "user-defined order",
		input:       []string{": countup 1 2 3 ;", "countup"},
		expected:    []int{1, 2, 3},
	},
	{
		description: "user-defined override",
		input:       []string{": foo dup ;", ": foo dup dup ;", "1 foo"},
		expected:    []int{1, 1, 1},
	},
	{
		description: "built-in override",
		input:       []string{": swap dup ;", "1 swap"},
		expected:    []int{1, 1},
	},
	{
		description: "built-in operator override",
		input:       []string{": + * ;", "3 4 +"},
		expected:    []int{12},
	},
	{
		description: "no redefinition",
		input:       []string{": foo 5 ;", ": bar foo ;", ": foo 6 ;", "bar foo"},
		expected:    []int{5, 6},
	},
	{
		description: "reuse in definition",
		input:       []string{": foo 10 ;", ": foo foo 1 + ;", "foo"},
		expected:    []int{11},
	},
	{
		description: "redefine numbers",
		input:       []string{": 1 2 ;"},
		error:       true,
	},
	{
		description: "non-existent word",
		input:       []string{"foo"},
		error:       true,
	},
	{
		description: "DUP case insensitivity",
		input:       []string{"1 DUP Dup dup"},
		expected:    []int{1, 1, 1, 1},
	},
	{
		description: "DROP case insensitivity",
		input:       []string{"1 2 3 4 DROP Drop drop"},
		expected:    []int{1},
	},
	{
		description: "SWAP case insensitivity",
		input:       []string{"1 2 SWAP 3 Swap 4 swap"},
		expected:    []int{2, 3, 4, 1},
	},
	{
		description: "OVER case insensitivity",
		input:       []string{"1 2 OVER Over over"},
		expected:    []int{1, 2, 1, 2, 1},
	},
	{
		description: "user-defined case insensitivity",
		input:       []string{": foo dup ;", "1 FOO Foo foo"},
		expected:    []int{1, 1, 1, 1},
	},
	{
		description: "definition case insensitivity",
		input:       []string{": SWAP DUP Dup dup ;", "1 swap"},
		expected:    []int{1, 1, 1, 1},
	},
	{
		description: "redefine of builtin after define user function on it",
		input:       []string{": foo dup ;", ": dup 1 ;", "2 foo"},
		expected:    []int{2, 2},
	},
}

func TestEval(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			stack, err := eval(tc.input)
			if tc.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, stack)
			}
		})
	}
}

func eval(input []string) ([]int, error) {
	e := NewEvaluator()
	var stack []int
	for _, row := range input {
		var err error
		stack, err = e.Process(row)
		if err != nil {
			return nil, err
		}
	}
	return stack, nil
}
