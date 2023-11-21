package exit

import (
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

// ErrorExit exits the spotlike application with the specified exit code.
func ErrorExit(exitCode ExitCode) {
	os.Exit(exitCode.GetNumber())
}
