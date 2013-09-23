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
	done         = make(chan bool)
	logger       gsmlog.GsmLogger
	processError error
)

func init() {
	flag.Var(&logger, "v", "-v\t\tUse verbose output.")
}

func main() {
	logger.Initialize()
	config := gsm.NewConfigurationFromFlags(logger)

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
	be := gsm.NewBundleExecer(*config)

	go consumer.Consume(deliveries)

	for delivery := range deliveries {
		switch delivery.(type) {
		case nil:
			return
		case error:
			logger.Println("main - something bad happened - error delivery type")
		case amqp.Delivery:
			instructions, err := orc.Orchestrate(delivery.(amqp.Delivery))
			if err != nil {
				logger.Println("main - Unable to determine instructions from message")
				logger.Printf("main - Message body: %s\n", string(delivery.(amqp.Delivery).Body))
			} else {
				go func() {
					err = be.ProcessInstructions(instructions)
					if err != nil {
						if processError != nil {
							logger.Println("main - Message processing erred twice in a row, something is probably wrong.")
							os.Exit(86)
						}
						processError = err
					} else {
						processError = nil
					}
				}()
			}
		default:
			logger.Println("main - something bad happened - unexpected delivery type")
		}
	}
}
