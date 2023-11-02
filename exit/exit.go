/*
Package exit provides functions for handling program exit codes and messages in the spotlike application.
*/
package exit

import (
	"fmt"
	"os"
)

// ExitCode represents possible exit codes for the spotlike application.
type ExitCode int

// constants
const (
	// CodeOk is an exit code indicating a successful operation.
	CodeOk ExitCode = iota
	// CodeErrInit is an exit code indicating an initialization error.
	CodeErrInit
	// CodeErrCmd is an exit code indicating a command error.
	CodeErrCmd
)

// GetNumber returns the numeric value of an ExitCode.
func (exitCode ExitCode) GetNumber() int {
	return int(exitCode)
}

// Exit exits the spotlike application with the OK exit code.
func Exit() {
	os.Exit(CodeOk.GetNumber())
}

// ExitWithMessage prints a message and exits the spotlike application with the OK exit code.
func ExitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(CodeOk.GetNumber())
}

// ErrorExit exits the spotlike application with the specified exit code.
func ErrorExit(exitCode ExitCode) {
	os.Exit(exitCode.GetNumber())
}

// ErrorExitWithMessage prints an error message, specifies the exit code, and then exits the spotlike application.
func ErrorExitWithMessage(message string, exitCode ExitCode) {
	fmt.Fprintf(os.Stderr, "Error :\t%s\n", message)
	os.Exit(exitCode.GetNumber())
}
