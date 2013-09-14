package main

import (
	"flag"
	"fmt"
	"gsm"
	gsmlog "gsm/log"
	"os"
)

var (
	done   = make(chan bool)
	logger gsmlog.GsmLogger
)

func init() {
	flag.Var(&logger, "v", "-v\t\tUse verbose output.")
}

func main() {
	config := gsm.NewConfigurationFromFlags()
	logger.Initialize()

	if config.DisplayVersion {
		fmt.Println(gsm.ProgVersion())
	}

	if config.DisplayRev {
		fmt.Println(gsm.Rev)
	}

	if config.ExitImmediately {
		os.Exit(1)
	}

	gsm.Serve()
}
