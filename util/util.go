package util

import (
	"fmt"
)

// FormatIndent indents the provided message by adding two spaces at the beginning.
func FormatIndent(message string) string {
	return "  " + message
}

// PrintlnWithBlankLineBelow prints the provided message with a blank line below it.
func PrintlnWithBlankLineBelow(message string) {
	fmt.Println(fmt.Sprintf("%s\n", message))
}

// PrintlnWithBlankLineAbove prints the provided message with a blank line above it.
func PrintlnWithBlankLineAbove(message string) {
	fmt.Println(fmt.Sprintf("\n%s", message))
}

// PrintBetweenBlankLine prints the provided message between two blank lines.
func PrintBetweenBlankLine(message string) {
	fmt.Println(fmt.Sprintf("\n%s\n", message))
}
