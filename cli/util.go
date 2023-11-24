package cli

import (
	"fmt"
)

// formatResult returns the formatted result
func formatResult(topic string, detail string) string {
	return fmt.Sprintf("%s\t:\t%s", topic, detail)
}

// formatErrorResult returns the formatted error result
func formatErrorResult(error error) string {
	return fmt.Sprintf("Error:\n\t%s", error)
}
