package main

import (
	"flag"
	"fmt"
	"gsm"
	gsmlog "gsm/log"
	"os"
)

import (
	"github.com/streadway/amqp"
)

var (
	done   = make(chan bool)
	logger gsmlog.GsmLogger
)

func init() {
	flag.Var(&logger, "v", "-v\t\tUse verbose output.")
}

func main() {
	config := gsm.NewConfigurationFromFlags(logger)
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

	deliveries := make(chan interface{})
	consumer := gsm.NewConsumer(*config)
	orc := gsm.NewOrchestrator(*config)

	go consumer.Consume(deliveries)

	go func() {
		for delivery := range deliveries {
			switch delivery.(type) {
			case nil:
				done <- true
			case error:
				logger.Println("something bad happened")
			case amqp.Delivery:
				instructions, err := orc.Orchestrate(delivery.(amqp.Delivery))
				if err != nil {
					logger.Println("Unable to determine instructions from message")
					logger.Printf("Message body: %s\n", string(delivery.(amqp.Delivery).Body))
				} else {
					go gsm.ProcessInstructions(instructions)
				}
			default:
				logger.Println("something bad happened")
			}
		}
	}()

	<-done
}
