package gswat

import (
	"flag"
	"fmt"
)

type Configuration struct {
	DisplayVersion  bool
	DisplayRev      bool
	ExitImmediately bool
}

var (
	revFlag     = flag.Bool("rev", false, "-rev\tPrint git revision and exit.")
	versionFlag = flag.Bool("version", false, "-version\tPrint version and exit.")
)

func NewConfigurationFromFlags() *Configuration {
	flag.Usage = func() {
		fmt.Println("Usage: gswat-server [options]")
		printOptions()
	}
	flag.Parse()

	exitImmediately := *revFlag || *versionFlag

	return &Configuration{
		DisplayVersion:  *versionFlag,
		DisplayRev:      *revFlag,
		ExitImmediately: exitImmediately,
	}
}

func printOptions() {
	fmt.Println("Options:")
	flag.VisitAll(func(flag *flag.Flag) {
		fmt.Println(flag.Usage)
	})
}
