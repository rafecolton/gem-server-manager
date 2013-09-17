package gsm

import (
	"flag"
	"fmt"
	gsmlog "gsm/log"
)

type Configuration struct {
	AmqpQos          int
	AuthToken        string
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
		fmt.Println(`Usage: gsm-server [options] <gemdir> <auth_token>

Required Args:
gemdir     - the directory in which to place the retrieved gems
auth_token - the GitHub authorization token
`)
		printOptions()
	}
	flag.Parse()

	var gemDir string
	var authToken string

	exitImmediately := *revFlag || *versionFlag

	if flag.NArg() < 2 {
		gemDir = ""
		authToken = ""
		exitImmediately = true
	} else {
		gemDir = flag.Arg(0)
		authToken = flag.Arg(1)
	}

	return &Configuration{
		AmqpQos:          *qosFlag,
		AuthToken:        authToken,
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
