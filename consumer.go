package gsm

import (
	gsmlog "gsm/log"
	"os"
	"os/exec"
)

import (
	"github.com/streadway/amqp"
)

var (
	consumerTag string
)

func Consume(connectionUri string, deliveries chan interface{}, logger gsmlog.GsmLogger) {
	defer func() { deliveries <- nil; close(deliveries) }()

	// establish connection
	conn, err := amqp.Dial(connectionUri)
	if err != nil {
		logger.Printf("Error creating connection: %+v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	logger.Println("Connection established")

	// open channel
	channel, err := conn.Channel()
	if err != nil {
		logger.Printf("Error opening channel: %+v\n", err)
		os.Exit(2)
	}

	// set qos
	qosValue := 10 // make configurable?
	err = channel.Qos(qosValue, 0, false)
	if err != nil {
		logger.Printf("Error setting QOS: %+v\n", err)
		os.Exit(3)
	}
	logger.Printf("Channel QOS set to %d\n", qosValue)

	uuidBytes, err := exec.Command("uuidgen").Output()
	if err != nil {
		logger.Printf("Error calling uuidgen: %+v\n", err)
		os.Exit(4)
	}

	consumerTag = string(uuidBytes)

	/*
		  autoAck = false (must manually Ack)
		  exclusive = true (so we only try to read from one consumer at a time)
		  noLocal = true (so that if this consumer republishes to the specified
			queue, it will not pick the message back up)
		  noWait = true
	*/
	consumerChan, err := channel.Consume("firehose", consumerTag, false, true, true, true, nil)
	if err != nil {
		logger.Printf("Error establishing consume channel: %+v", err)
		os.Exit(5)
	}

	for {
		select {
		case delivery := <-consumerChan:
			deliveries <- delivery
		}
	}

	//TODO: handle channel/connection closing
}
