package gsm

import (
	"os"
	"os/exec"
	"sync"
)

import (
	"github.com/streadway/amqp"
)

var (
	consumerTag string
	conn        *amqp.Connection
)

type Consumer struct {
	Configuration
	channel     *amqp.Channel
	conn        *amqp.Connection
	notifyClose chan *amqp.Error
	sync.Mutex
}

func NewConsumer(config Configuration) *Consumer {
	c := &Consumer{
		Configuration: config,
	}
	return c
}

func (me *Consumer) Consume(deliveries chan interface{}) {
	if err := me.establishConnection(); err != nil {
		me.Logger.Println("Error establishing a connection & channel: %+v", err)
		os.Exit(7)
	}

	// set qos
	qosValue := 10 // make configurable?
	err := me.channel.Qos(qosValue, 0, false)
	if err != nil {
		me.Logger.Printf("Error setting QOS: %+v\n", err)
		os.Exit(3)
	}
	me.Logger.Printf("Channel QOS set to %d\n", qosValue)

	uuidBytes, err := exec.Command("uuidgen").Output()
	if err != nil {
		me.Logger.Printf("Error calling uuidgen: %+v\n", err)
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
	consumerChan, err := me.channel.Consume("firehose", consumerTag, false, true, true, true, nil)
	if err != nil {
		me.Logger.Printf("Error establishing consume channel: %+v", err)
		os.Exit(5)
	}

	for delivery := range consumerChan {
		deliveries <- delivery
	}
}

func (me *Consumer) establishConnection() (err error) {
	me.Lock()
	defer me.Unlock()

	//establish connection
	if me.conn != nil {
		return nil
	}
	me.conn, err = amqp.Dial(me.ConnectionString)
	if err != nil {
		me.Logger.Printf("Error creating connection: %+v\n", err)
		return err
	}
	me.Logger.Println("amqp - connection established")

	// open channel
	me.channel, err = me.conn.Channel()
	if err != nil {
		return err
	}
	me.Logger.Println("amqp - channel opened")

	if err = me.channel.Confirm(false); err != nil {
		return err
	}
	me.Logger.Println("amqp - confirm mode set")

	go func() {
		me.notifyClose = me.channel.NotifyClose(make(chan *amqp.Error))

		select {
		case e := <-me.notifyClose:
			me.Logger.Printf("amqp - The channel opened with RabbitMQ has been closed. %d: %s", e.Code, e.Reason)
			me.disconnect()
			if err := me.establishConnection(); err != nil {
				me.Logger.Println("Error establishing a connection & channel: %+v", err)
				os.Exit(7)
			}
		}
	}()

	return nil
}

func (me *Consumer) disconnect() {
	me.Lock()
	defer me.Unlock()

	if me.channel != nil {
		me.channel.Close()
		me.channel = nil
	}

	if me.conn != nil {
		me.conn.Close()
		me.conn = nil
	}
}
