/*
Copyright	 © 2023 Paras Raba
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "neko [OPTION]... [FILE]...",
	Short: "neko a chip clone of cat command",
	Long:  `neko a tool that tries and fails to mimic the fraction of power of GNU cat`,
	Run: func(cmd *cobra.Command, args []string) {
		for i := range args {
			file, err := os.Open(args[i])
			if err != nil {
				fmt.Printf("neko %s: %s\n", args[i], errors.Unwrap(err))
			} else {
				printFileContent(file)
			}
		}
	},
}

func printFileContent(f *os.File) {
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		lineNum++
		var result string

		if line == "" {
			numsOfEmptyLine++
		}

		if numberNonBlank {
			if line != "" {
				currLineNum := lineNum - numsOfEmptyLine
				result = createNumberedLine(line, currLineNum)
			}
		}

		if showLineNumber && !numberNonBlank {
			result = createNumberedLine(line, lineNum)
		}

		if showTabs {
			result = strings.ReplaceAll(line, "\t", "^I")
		}

		if showEnds {
			result = line + "$"
		}

		if showNonPrinting {
			result = convertNonPrintin(line)
		}
		anyOpts := numberNonBlank || showLineNumber || showEnds || showNonPrinting || showTabs
		if !anyOpts {
			result = line
		}
		fmt.Println(result)
	}
	if err := sc.Err(); err != nil {
		fmt.Printf("neko %s: %s\n", f.Name(), errors.Unwrap(err))
	}
}

func leftPad(str string, size int, char string) string {
	toBepadded := size - len(str)
	for i := 0; i < toBepadded; i++ {
		str = char + str
	}
	return str
}

func convertNonPrintin(line string) string {
	var result string
	for _, ch := range line {
		if ch >= 32 {
			if ch < 127 {
				result += string(ch)
			} else if ch == 127 {
				result += "^?"
			} else {
				result += "M-"
				if ch < 128+32 {
					if ch < 128+127 {
						result += string(ch - 128)
					} else {
						result += "^?"
					}
				} else {
					result += "^"
					result += string(ch - 128 + 64)
				}
			}
		} else if ch == '\t' && !showTabs {
			result += "\t"
		} else {
			result += "^"
			result += string(ch + 64)
		}
	}
	return result
}

func createNumberedLine(line string, num int) string {
	paddedNum := leftPad(fmt.Sprintf("%d", num), 6, " ")
	line = paddedNum + "  " + line
	return line
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().BoolVarP(&numberNonBlank, "number-nonblank", "b", false, "number nonempty output lines, overrides -n")
	rootCmd.Flags().BoolVarP(&showLineNumber, "number", "n", false, "number all output lines")
	rootCmd.Flags().BoolVarP(&showEnds, "show-ends", "E", false, "display $ at end of each line")
	rootCmd.Flags().BoolVarP(&showNonPrinting, "show-nonprinting", "v", false, "use ^ and M- notation, except for LFD and TAB")
	rootCmd.Flags().BoolVarP(&showTabs, "show-tabs", "t", false, "display TAB characters as ^I")
}
