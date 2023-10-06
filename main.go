/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package main

import (
	"github.com/yanosea/spotlike/cmd"
	"github.com/yanosea/spotlike/log"
)

func main() {
	// get new command
	command := cmd.New()

	// init command
	if err := command.Init(); err != nil {
		log.ErrorExit(err.Error(), log.ExitCodeErrInit)
	}

	// execute command
	if err := command.Execute(); err != nil {
		log.ErrorExit(err.Error(), log.ExitCodeErrCmd)
	}
}
