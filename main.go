// Package main is the entry point of spotlike.
package main

import (
	"os"

	"github.com/yanosea/spotlike/cmd"
)

func main() {
	g := cmd.NewGlobalOption(fmtProxy, osProxy, strconvProxy)
	os.Exit(cmd.Execute())
}
