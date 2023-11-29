package util

import (
	"fmt"
)

func FormatIndent(message string) string {
	return "  " + message
}

func PrintlnWithBlankLineBelow(message string) {
	fmt.Println(fmt.Sprintf("%s\n", message))
}
func PrintlnWithBlankLineAbove(message string) {
	fmt.Println(fmt.Sprintf("\n%s", message))
}

func PrintBetweenBlankLine(message string) {
	fmt.Println(fmt.Sprintf("\n%s\n", message))
}
