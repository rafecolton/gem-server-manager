package gsm

import (
	"os"
)

import (
	"github.com/streadway/amqp"
	gsmlog "gsm/log"
)

func Orchestrate(delivery amqp.Delivery, logger gsmlog.GsmLogger) string {
	body := string(delivery.Body)
	err := delivery.Ack(false)
	if err != nil {
		logger.Printf("Error acking delivery %+v: %+v\n", delivery, err)
		os.Exit(6)
	}
	return body
}
