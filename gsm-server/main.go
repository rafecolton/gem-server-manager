package main

import (
	"fmt"
	"gsm"
	"os"
)

func main() {
	config := gsm.NewConfigurationFromFlags()

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
