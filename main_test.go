package main

import "testing"

func TestLeftPad(t *testing.T) {
	type input struct {
		str  string
		size int
		char rune
	}
	testCases := map[string]struct {
		input    input
		expected string
	}{
		"spaces": {
			input: input{
				str:  "1",
				size: 5,
				char: ' ',
			},
			expected: "    1",
		},
		"at the rate": {
			input: input{
				str:  "1",
				size: 5,
				char: '@',
			},
			expected: "@@@@1",
		},
		"padded input": {
			input: input{
				str:  "12345",
				size: 5,
				char: '"',
			},
			expected: "12345",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			got := leftPad(test.input.str, test.input.size, test.input.char)
			if got != test.expected {
				t.Fatalf("got = %s expected = %s", got, test.expected)
			}
		})
	}
}
