package gsm

import (
	"flag"
	"fmt"
	gsmlog "gsm/log"
)

type Configuration struct {
	AmqpQos          int
	ConnectionString string
	DisplayRev       bool
	DisplayVersion   bool
	ExitImmediately  bool
	GemDir           string
	Logger           gsmlog.GsmLogger
}

var (
	revFlag     = flag.Bool("rev", false, "-rev\t\tPrint git revision and exit.")
	versionFlag = flag.Bool("version", false, "-version\tPrint version and exit.")
	uriFlag     = flag.String("uri", "amqp://guest:guest@localhost:5672", "-uri\t\tAMQP uri for consumer.")
	qosFlag     = flag.Int("qos", 10, "-qos\t\tQOS for the AMQP connection")
)

func NewConfigurationFromFlags(logger gsmlog.GsmLogger) *Configuration {
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
		AmqpQos:          *qosFlag,
		ConnectionString: *uriFlag,
		DisplayRev:       *revFlag,
		DisplayVersion:   *versionFlag,
		ExitImmediately:  exitImmediately,
		GemDir:           gemDir,
		Logger:           logger,
	}
}

func printOptions() {
	fmt.Println("Options:")
	flag.VisitAll(func(flag *flag.Flag) {
		fmt.Println(flag.Usage)
	})
}
