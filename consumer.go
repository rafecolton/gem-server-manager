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
	return &Consumer{
		Configuration: config,
	}
}

func (me *Consumer) Consume(deliveries chan interface{}) {
	if err := me.establishConnection(); err != nil {
		me.Logger.Println("amqp - Error establishing a connection & channel: %+v", err)
		os.Exit(7)
	}

	// set qos
	err := me.channel.Qos(me.AmqpQos, 0, false)
	if err != nil {
		me.Logger.Printf("amqp - Error setting QOS: %+v\n", err)
		os.Exit(3)
	}
	me.Logger.Printf("amqp - Channel QOS set to %d\n", me.AmqpQos)

	uuidBytes, err := exec.Command("uuidgen").Output()
	if err != nil {
		me.Logger.Printf("amqp - Error calling uuidgen: %+v\n", err)
		os.Exit(4)
	}

	consumerTag = string(uuidBytes)

	err = me.channel.ExchangeDeclare(me.Exchange, "topic", true, false, false, true, nil)
	if err != nil {
		me.Logger.Printf("amqp - Error declaring exchange %s\n", me.Exchange)
		os.Exit(13)
	}

	_, err = me.channel.QueueDeclare(me.Queue, true, false, false, true, nil)
	if err != nil {
		me.Logger.Printf("amqp - Error declaring queue %s\n", me.Queue)
		os.Exit(17)
	}

	err = me.channel.QueueBind(me.Queue, me.Binding, me.Exchange, true, nil)
	if err != nil {
		me.Logger.Printf("amqp - Error binding queue %s to exchange %s using binding %s\n",
			me.Queue,
			me.Exchange,
			me.Binding)
		os.Exit(19)
	}

	/*
		autoAck = false (must manually Ack)
		exclusive = true (so we only try to read from one consumer at a time)
		noLocal = true (so that if this consumer republishes to the specified
		queue, it will not pick the message back up)
		noWait = true
	*/
	consumerChan, err := me.channel.Consume(me.Queue, consumerTag, false, true, true, true, nil)
	if err != nil {
		me.Logger.Printf("amqp - Error establishing consume channel: %+v", err)
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
		me.Logger.Printf("amqp - Error creating connection: %+v\n", err)
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
				me.Logger.Println("amqp - Error establishing a connection & channel: %+v", err)
				os.Exit(11)
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
