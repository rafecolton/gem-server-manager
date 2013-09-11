package main

import (
	"fmt"
	"gswat"
	"os"
)

func main() {
	config := gswat.NewConfigurationFromFlags()

	if config.DisplayVersion {
		fmt.Println(gswat.ProgVersion())
	}

	if config.DisplayRev {
		fmt.Println(gswat.Rev)
	}

	if config.ExitImmediately {
		os.Exit(1)
	}

	gswat.Serve()
}
