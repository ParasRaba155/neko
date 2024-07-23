package main

import (
	"testing"
)

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
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := leftPad(test.input.str, test.input.size, test.input.char)
			if got != test.expected {
				t.Errorf("got = %q expected = %q", got, test.expected)
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
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := createNumberedLine(test.input.line, test.input.num)
			if got != test.expected {
				t.Errorf("got = %q expected = %q", got, test.expected)
			}
		})
	}
}

func TestConvertNonPrintin(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		// NOTE: Go strings are sequences of bytes, and when you use non-ASCII values
		// they are interpreted using the Latin-1 encoding for byte values between 128 and 255 (VIA OpenAI)
		"Handling ASCII 128 and 255 Delete Chars": {
			input:    "\x7fÿ", // '\x7f' is ASCII 127 and 'ÿ' is ASCII 255
			expected: "^?M-^?",
		},
		"Handling ASCII less than 32": {
			input:    "\x01\x02\x04\t",
			expected: "^^^\t",
		},
		"Handling tab character": {
			input:    "Hello \tWorld",
			expected: "Hello \tWorld",
		},
		"Handling the ASCII between 128 to 159": {
			// ASCII 128, 131, 159
			input:    "\u0080\u0083\u009f",
			expected: "M-\x00M-\x03M-\x1f",
		},
		"Handling the ASCII between 160 to 254": {
			// ASCII 160, 172, 250, 254
			input:    "\u00a0¬úþ",
			expected: "M-^`M-^lM-^ºM-^¾",
		},
	}

	for name, tt := range tests {
		tt := tt // capture range variable
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := convertNonPrintingChars(tt.input)
			if result != tt.expected {
				t.Errorf("got = %q expected = %q", result, tt.expected)
			}
		})
	}
}
