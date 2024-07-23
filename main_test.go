package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		t.Run(name, func(t *testing.T) {
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := createNumberedLine(test.input.line, test.input.num)
			if got != test.expected {
				t.Errorf("got = %q expected = %q", got, test.expected)
			}
		})
	}
}

func getAllPrintableASCII() []byte {
	result := make([]byte, 0, 127-32+1)
	var i byte
	for i = 32; i < 127; i++ {
		result = append(result, i)
	}
	return result
}

func getAllCharsLessThan32() [32]byte {
	var result [32]byte
	var i byte
	for i = 0; i < 32; i++ {
		result[i] = i
	}
	return result
}

func TestConvertNonPritnin(t *testing.T) {
	t.Parallel()
	t.Run("Common English ASCII Chars", func(t *testing.T) {
		input, expected := getAllPrintableASCII(), string(getAllPrintableASCII())
		got := convertNonPrintin(string(input), false)
		if !cmp.Equal(got, expected) {
			diff := cmp.Diff(got, expected)
			t.Errorf("convertNonPrintin() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Characters with ASCII less than 32 and showTabs", func(t *testing.T) {
		input := getAllCharsLessThan32()
		var s strings.Builder
		for range input {
			s.WriteRune('^')
		}

		got := convertNonPrintin(string(input[:]), true)
		if !cmp.Equal(got, s.String()) {
			diff := cmp.Diff(got, s.String())
			t.Errorf("convertNonPrintin() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Characters with ASCII less than 32 and !showTabs", func(t *testing.T) {
		input := getAllCharsLessThan32()

		var s strings.Builder
		for i := range input {
			if input[i] == '\t' {
				s.WriteRune('\t')
				continue
			}
			s.WriteRune('^')
		}

		got := convertNonPrintin(string(input[:]), false)
		if !cmp.Equal(got, s.String()) {
			diff := cmp.Diff(got, s.String())
			t.Errorf("convertNonPrintin() mismatch (-want +got):\n%s", diff)
		}
	})
}
