package rabbitmq

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	queueName   = "processes_queue"
	consumerTag = ""
	durable     = false
	autoDelete  = false
	exclusive   = false
	noWait      = false
	noLocal     = false
	autoAck     = true
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	consumer   *<-chan amqp.Delivery
)

func GetConnection() *amqp.Connection {
	if connection != nil {
		return connection
	}
	openConnection()
	return connection
}

func GetChannel() *amqp.Channel {
	if channel != nil {
		return channel
	}
	openChannel()
	return channel
}

func GetQueue() *amqp.Queue {
	if queue != nil {
		return queue
	}
	openQueue()
	return queue
}

func GetConsumer() *<-chan amqp.Delivery {
	if consumer != nil {
		return consumer
	}
	openConsumer()
	return consumer
}

func Close() {
	connection.Close()
	channel.Close()
}

func openConnection() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		fmt.Println("error trying to connect with rabbitmq\n[ERROR] ", err)
		return
	}
	connection = conn
}

func openChannel() {
	ch, err := connection.Channel()
	if err != nil {
		fmt.Println("error while creating a new channel\n[ERROR] ", err)
		return
	}
	channel = ch
}

func openQueue() {
	q, err := channel.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		amqp.Table{},
	)
	if err != nil {
		fmt.Println("error trying to declare the queue\n[ERROR] ", err)
		return
	}
	queue = &q
}

func openConsumer() {
	msgs, err := channel.Consume(
		queueName,
		consumerTag,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		amqp.Table{},
	)
	if err != nil {
		fmt.Println("error trying to get the consumer\n[ERROR] ", err)
		return
	}
	consumer = &msgs
}
