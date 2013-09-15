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

	go consumer.Consume(deliveries)

	go func() {
		for delivery := range deliveries {
			switch delivery.(type) {
			case nil:
				done <- true
			case error:
				logger.Println("something bad happened")
			default:
				instructions := gsm.Orchestrate(delivery.(amqp.Delivery))
				gsm.ProcessInstructions(instructions)
				logger.Printf("raw_instructions: %+v\n", instructions)
			}
		}
	}()

	<-done
}
