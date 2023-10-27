/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package main

import (
	"github.com/yanosea/spotlike/cmd"
	"github.com/yanosea/spotlike/log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.ErrorExit(log.ErrCmd)
	}
}
