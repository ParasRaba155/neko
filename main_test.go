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
				t.Fatalf("got = %q expected = %q", got, test.expected)
			}
		})
	}
}

func TestCreateNumberdLine(t *testing.T) {
	type input struct {
		line string
		num  int
	}
	testCases := map[string]struct {
		input    input
		expected string
	}{
		"23rd line": {
			input: input{
				line: "This is a text containing some words",
				num:  23,
			},
			expected: "    23  This is a text containing some words",
		},
		"1st line": {
			input: input{
				line: "This is a text containing some words",
				num:  1,
			},
			expected: "     1  This is a text containing some words",
		},
		"0th line": {
			input: input{
				line: "This is a text containing some words",
				num:  0,
			},
			expected: "     0  This is a text containing some words",
		},
		"negative line": {
			input: input{
				line: "This is a text containing some words",
				num:  -2,
			},
			expected: "    -2  This is a text containing some words",
		},
		"line number is 1e6": {
			input: input{
				line: "This is a text containing some words\t hello",
				num:  1e6,
			},
			expected: "1000000  This is a text containing some words\t hello",
		},
		"line number is 1e6+234": {
			input: input{
				line: "This is a text containing some words\t hello",
				num:  1e6 + 234,
			},
			expected: "1000234  This is a text containing some words\t hello",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			got := createNumberedLine(test.input.line, test.input.num)
			if got != test.expected {
				t.Fatalf("got = %q expected = %q", got, test.expected)
			}
		})
	}
}
