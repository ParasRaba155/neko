package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

// flag vars
var (
	numberNonBlank  bool
	showLineNumber  bool
	showEnds        bool
	showNonPrinting bool
	showTabs        bool
)

var (
	lineNum         = 0
	numsOfEmptyLine = 0
)

const usage = `Usage: neko [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

  -b, --number-nonblank    number nonempty output lines, overrides -n
  -e, --show-ends          display $ at end of each line
  -n, --number             number all output lines
  -t, --show-tabs          display TAB characters as ^I
  -v, --show-nonprinting   use ^ and M- notation, except for LFD and TAB
  -h, --help               display this help and exit

Examples:
  neko f - g  Output f's contents, then standard input, then g's contents.
  neko        Copy standard input to standard output.
`

func printContent(r io.Reader, filename string) {
	sc := bufio.NewScanner(r)

	buf := make([]byte, 1024)
	sc.Buffer(buf, 512)

	for sc.Scan() {
		line := sc.Text()
		lineNum++
		var b strings.Builder

		if line == "" {
			numsOfEmptyLine++
		}

		if numberNonBlank {
			if line != "" {
				currLineNum := lineNum - numsOfEmptyLine
				b.WriteString(createNumberedLine(line, currLineNum))
			}
		}

		if showLineNumber && !numberNonBlank {
			b.WriteString(createNumberedLine(line, lineNum))
		}

		if showTabs {
			if b.Len() == 0 {
				b.WriteString(strings.ReplaceAll(line, "\t", "^I"))
			} else {
				temp := b.String()
				b.Reset()
				b.WriteString(strings.ReplaceAll(temp, "\t", "^I"))
			}
		}

		if showEnds {
			if b.Len() == 0 {
				b.WriteString(line)
			}
			b.WriteRune('$')
		}

		if showNonPrinting {
			if b.Len() == 0 {
				b.WriteString(convertNonPrintin(line))
			} else {
				temp := b.String()
				b.Reset()
				b.WriteString(convertNonPrintin(temp))
			}
		}

		anyOpts := numberNonBlank || showLineNumber || showEnds || showNonPrinting || showTabs
		if !anyOpts {
			b.WriteString(line)
		}
		fmt.Fprintln(os.Stdout, b.String())
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "neko %s: %s\n", filename, errors.Unwrap(err))
	}
}

func leftPad(str string, size int, char rune) string {
	var b strings.Builder
	toBepadded := size - len(str)
	for i := 0; i < toBepadded; i++ {
		b.WriteRune(char)
	}
	b.WriteString(str)
	return b.String()
}

func convertNonPrintin(line string) string {
	var result strings.Builder
	for _, ch := range line {
		// from 32 to 127 where common day English ASCII resides
		if ch < 32 {
			// handle the tab char
			if ch == '\t' && !showTabs {
				result.WriteRune('\t')
				continue
			}
			// for non tab char write carrot char and move on
			result.WriteRune('^')
			continue
		}

		// now handle the English ASCII chars
		if ch < 127 {
			result.WriteRune(ch)
			continue
		}

		// if char is delete char i.e. 127
		if ch == '\x7F' {
			result.WriteRune('^')
			result.WriteRune('?')
			continue
		}

		// handle non ASCII characters
		result.WriteRune('M')
		result.WriteRune('-')

		// range 128 to 159 is to be mapped to lower ASCII
		if ch < 128+32 {
			result.WriteRune(ch - 128)
			continue
		}

		// 255 is treated as special delete character
		if ch == 128+127 {
			result.WriteRune('^')
			result.WriteRune('?')
			continue
		}

		// convert other control characters to carrot
		result.WriteRune('^')
		result.WriteRune(ch - 128 + 64)
	}
	return result.String()
}

func createNumberedLine(line string, num int) string {
	var b strings.Builder
	numStr := strconv.FormatInt(int64(num), 10)
	b.WriteString(leftPad(numStr, 6, ' '))
	b.WriteString("  ")
	b.WriteString(line)
	return b.String()
}

func init() {
	// flags
	flag.BoolVarP(&numberNonBlank, "number-nonblank", "b", false, "number nonempty output lines, overrides -n")
	flag.BoolVarP(&showLineNumber, "number", "n", false, "number all output lines")
	flag.BoolVarP(&showEnds, "show-ends", "e", false, "display $ at end of each line")
	flag.BoolVarP(&showNonPrinting, "show-nonpriting", "v", false, "use ^ and M- notation, except for LFD and TAB")
	flag.BoolVarP(&showTabs, "show-tabs", "t", false, "display TAB characters as ^I")
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	// TODO: make flag parsing posix compliant without the `pflag' dependency
	flag.Parse()

	if len(flag.Args()) == 0 {
		printContent(os.Stdin, "stdin")
	}

	for _, arg := range flag.Args() {
		if arg == "-" {
			printContent(os.Stdin, "stdin")
			continue
		}
		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stdout, "neko: %s: %s\n", arg, errors.Unwrap(err))
		} else {
			printContent(file, file.Name())
			err = file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "neko: %s: %s\n", arg, errors.Unwrap(err))
			}
		}
	}
}
