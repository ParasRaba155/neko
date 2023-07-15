package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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
      --help     display this help and exit

Examples:
  neko f - g  Output f's contents, then standard input, then g's contents.
`

func printFileContent(r io.Reader, filename string) {
	sc := bufio.NewScanner(r)
	// sc.Split(bufio.ScanBytes)
	buf := make([]byte, 1024)
	sc.Buffer(buf, 512)

	stdOut := bufio.NewWriter(os.Stdout)
	defer func() {
		err := stdOut.Flush()
		if err != nil {
			fmt.Fprintf(os.Stderr, "neko: %s could not flush the writer: %s", filename, err)
		}
	}()

	stdErr := bufio.NewWriter(os.Stderr)
	defer func() {
		err := stdErr.Flush()
		if err != nil {
			fmt.Fprintf(os.Stderr, "neko: %s could not flush the writer: %s", filename, err)
		}
	}()

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
				res := b.String()
				b.Reset()
				b.WriteString(convertNonPrintin(res))
			}
		}
		anyOpts := numberNonBlank || showLineNumber || showEnds || showNonPrinting || showTabs
		if !anyOpts {
			b.WriteString(line)
		}
		fmt.Fprintln(stdOut, b.String())
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(stdErr, "neko %s: %s\n", filename, errors.Unwrap(err))
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
		if ch >= 32 {
			if ch < 127 {
				result.WriteRune(ch)
			} else if ch == 127 {
				result.WriteRune('^')
				result.WriteRune('?')
			} else {
				result.WriteRune('M')
				result.WriteRune('-')
				if ch < 128+32 {
					if ch < 128+127 {
						result.WriteRune(ch - 128)
					} else {
						result.WriteRune('^')
						result.WriteRune('?')
					}
				} else {
					result.WriteRune('^')
					result.WriteRune(ch - 128 + 64)
				}
			}
		} else if ch == '\t' && !showTabs {
			result.WriteRune('\t')
		} else {
			result.WriteRune('^')
		}
	}
	return result.String()
}

func createNumberedLine(line string, num int) string {
	// paddedNum := leftPad(fmt.Sprintf("%d", num), 6, ' ')
	// line = paddedNum + "  " + line
	var b strings.Builder
	b.WriteString(leftPad(fmt.Sprintf("%d", num), 6, ' '))
	b.WriteString("  ")
	b.WriteString(line)
	// return fmt.Sprintf("%s  %s", leftPad(fmt.Sprintf("%d", num), 6, ' '), line)
	// return line
	return b.String()
}

func init() {
	// flags
	flag.BoolVar(&numberNonBlank, "number-nonblank", false, "number nonempty output lines, overrides -n")
	flag.BoolVar(&numberNonBlank, "b", false, "number nonempty output lines, overrides -n")
	flag.BoolVar(&showLineNumber, "number", false, "number all output lines")
	flag.BoolVar(&showLineNumber, "n", false, "number all output lines")
	flag.BoolVar(&showEnds, "show-ends", false, "display $ at end of each line")
	flag.BoolVar(&showEnds, "e", false, "display $ at end of each line")
	flag.BoolVar(&showNonPrinting, "show-nonpriting", false, "use ^ and M- notation, except for LFD and TAB")
	flag.BoolVar(&showNonPrinting, "v", false, "use ^ and M- notation, except for LFD and TAB")
	flag.BoolVar(&showTabs, "show-tabs", false, "display TAB characters as ^I")
	flag.BoolVar(&showTabs, "t", false, "display TAB characters as ^I")
}

func main() {
	if len(os.Args) == 1 {
		// the first argument by default is the name of the build file
		// TODO:DEAL WITH NO ARGS
		fmt.Fprintln(os.Stderr, "TODO: DEAL WITH no args")
		os.Exit(1)
	}

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, usage)
	}

	// FIXME: make flag parsing posix compliant
	flag.Parse()
	for _, arg := range flag.Args() {
		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stdout, "neko: %s: %s\n", arg, errors.Unwrap(err))
		} else {
			printFileContent(file, file.Name())
			err = file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "neko: %s: %s\n", arg, errors.Unwrap(err))
			}
		}
	}
}
