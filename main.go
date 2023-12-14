/*
Package main is the entry point of the spotlike application and handles command execution and program exit codes.
*/
package main

import (
	"github.com/yanosea/spotlike/cmd"
	"github.com/yanosea/spotlike/exit"
)

// main is the entry point of spotlike
func main() {
	if err := cmd.Init(); err != nil {
		exit.ErrorExit(exit.CodeErrInit)
	}

	if err := cmd.Execute(); err != nil {
		exit.ErrorExit(exit.CodeErrCmd)
	}

	exit.Exit()
}
