package gsm

import (
	"flag"
	"fmt"
)

type Configuration struct {
	DisplayVersion  bool
	DisplayRev      bool
	ExitImmediately bool
	UriFlag         string
	GemDir          string
}

var (
	revFlag     = flag.Bool("rev", false, "-rev\t\tPrint git revision and exit.")
	versionFlag = flag.Bool("version", false, "-version\tPrint version and exit.")
	uriFlag     = flag.String("uri", "amqp://guest:guest@localhost:5672", "-uri\t\tAMQP uri for consumer.")
)

func NewConfigurationFromFlags() *Configuration {
	flag.Usage = func() {
		fmt.Println("Usage: gsm-server [options] <gemdir>")
		printOptions()
	}
	flag.Parse()

	var gemDir string

	exitImmediately := *revFlag || *versionFlag

	if flag.NArg() < 1 {
		gemDir = ""
		exitImmediately = true
	} else {
		gemDir = flag.Arg(0)
	}

	return &Configuration{
		DisplayVersion:  *versionFlag,
		DisplayRev:      *revFlag,
		ExitImmediately: exitImmediately,
		UriFlag:         *uriFlag,
		GemDir:          gemDir,
	}
}

func printOptions() {
	fmt.Println("Options:")
	flag.VisitAll(func(flag *flag.Flag) {
		fmt.Println(flag.Usage)
	})
}
