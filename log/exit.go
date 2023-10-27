/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package log

import (
	"fmt"
	"os"
)

// ExitCode : exit code of spotlike
type ExitCode int

const (
	// Ok : ok
	Ok ExitCode = iota
	// ErrInit : init error
	ErrInit
	// ErrCmd : command error
	ErrCmd
)

// GetNumber : get number of exit code
func (exitCode ExitCode) GetNumber() int {
	return int(exitCode)
}

// Exit : output message and exit spotlike
func Exjt(message string) {
	fmt.Println(message)
	os.Exit(Ok.GetNumber())
}

// ErrorExit : output exit code and exit spotlike
func ErrorExit(exitCode ExitCode) {
	os.Exit(exitCode.GetNumber())
}

// ErrorExitWithMessage : output error message and exit code, then exit spotlike
func ErrorExitWithMessage(message string, exitCode ExitCode) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	os.Exit(exitCode.GetNumber())
}
