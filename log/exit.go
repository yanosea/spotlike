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

// GetNumber : get number of exit code
func (exitCode ExitCode) GetNumber() int {
	return int(exitCode)
}

const (
	// ExitCodeOK : ok
	ExitCodeOK ExitCode = iota
	// ExitCodeErrInit : init error
	ExitCodeErrInit
	// ExitCodeErrCmd : command error
	ExitCodeErrCmd
)

// Exit : output log and exit command
func Exit(message string) {
	fmt.Println(message)
	os.Exit(ExitCodeOK.GetNumber())
}

// ErrorExit : output error log and exit command
func ErrorExit(message string, exitCode ExitCode) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	os.Exit(exitCode.GetNumber())
}
