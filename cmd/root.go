/*
Copyright © 2023 Paras Raba
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
	numberNonBlank bool
	showLineNumber bool
    showEnds bool
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
		lineNum += 1
		line := sc.Text()
		if line == "" {
			numsOfEmptyLine += 1
		}
		if numberNonBlank {
			if line == "" {
				fmt.Println(line)
				continue
			}
            currentNonEmptyLine := lineNum - numsOfEmptyLine
            line = createNumberedLine(line, currentNonEmptyLine)
            fmt.Println(line)
			continue
		}
		if showLineNumber {
            line = createNumberedLine(line, lineNum)
			fmt.Println(line)
		}
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

func appendDollar(line string) string {
    return line + "$"
}

func replaceTabs(line string) string {
    return strings.ReplaceAll(line, "\t", "^I")
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
}
