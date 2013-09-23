package gsm

import (
	"flag"
	"fmt"
	gsmlog "gsm/log"
)

type Configuration struct {
	AmqpQos          int
	AuthToken        string
	Binding          string
	ConnectionString string
	DisplayRev       bool
	DisplayVersion   bool
	Exchange         string
	ExitImmediately  bool
	GemDir           string
	GibHost          string
	Logger           gsmlog.GsmLogger
	Queue            string
	ScriptLoc        string
}

var (
	revFlag     = flag.Bool("rev", false, "-rev\t\tPrint git revision and exit.")
	versionFlag = flag.Bool("version", false, "-version\tPrint version and exit.")
	uriFlag     = flag.String("uri", "amqp://guest:guest@localhost:5672", fmt.Sprintf(
		"-uri\t\t%s %s",
		"AMQP uri for consumer.",
		"Default: amqp://guest:guest@localhost:5672/"))
	qosFlag   = flag.Int("qos", 10, "-qos\t\tQOS for the AMQP connection. Default: 10.")
	queue     = flag.String("queue", "firehose", "-queue\t\tAMQP queue from which to consume. Default: firehose")
	exchange  = flag.String("exchange", "", "-exchange\tAMQP exchange to which to bind queue.")
	binding   = flag.String("binding", "#", "-binding\tAMQP binding for messages to be routed to the queue. Default: #")
	scriptLoc = flag.String("script", "/tmp/retrieve-gems", fmt.Sprintf(
		"-script\t\t%s %s",
		"Location for gem-retrieval script",
		"Default: /tmp/retrieve-gems"))
	gibHost = flag.String("gibhost", "http://localhost:9292/", fmt.Sprintf(
		"-gibhost\t%s %s",
		"Geminabox host",
		"Default: http://localhost:9292/"))
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
		Binding:          *binding,
		ConnectionString: *uriFlag,
		DisplayRev:       *revFlag,
		DisplayVersion:   *versionFlag,
		Exchange:         *exchange,
		ExitImmediately:  exitImmediately,
		GemDir:           gemDir,
		GibHost:          *gibHost,
		Logger:           logger,
		Queue:            *queue,
		ScriptLoc:        *scriptLoc,
	}
}

func printOptions() {
	fmt.Println("Options:")
	flag.VisitAll(func(flag *flag.Flag) {
		fmt.Println(flag.Usage)
	})
}
