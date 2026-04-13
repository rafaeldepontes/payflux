package rabbitmq

import (
	"fmt"
	"os"
	"sync"
	"time"

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
	connOnce   sync.Once
	connection *amqp.Connection

	chanOnce sync.Once
	channel  *amqp.Channel

	queueOnce sync.Once
	queue     *amqp.Queue

	consOnce sync.Once
	consumer *<-chan amqp.Delivery
)

type Broker struct{}

func (b Broker) Publish(body []byte) error {
	return Publish(body)
}

func GetConnection() *amqp.Connection {
	connOnce.Do(func() {
		openConnection()
	})
	return connection
}

func GetChannel() *amqp.Channel {
	chanOnce.Do(func() {
		conn := GetConnection()
		if conn != nil {
			openChannel()
		}
	})
	return channel
}

func GetQueue() *amqp.Queue {
	queueOnce.Do(func() {
		ch := GetChannel()
		if ch != nil {
			openQueue()
		}
	})
	return queue
}

func GetConsumer() *<-chan amqp.Delivery {
	consOnce.Do(func() {
		// Ensure queue is declared before consuming
		q := GetQueue()
		if q != nil {
			openConsumer()
		}
	})
	return consumer
}

func Close() {
	if channel != nil {
		_ = channel.Close()
	}
	if connection != nil {
		_ = connection.Close()
	}
}

func Publish(body []byte) error {
	ch := GetChannel()
	if ch == nil {
		return fmt.Errorf("rabbitmq channel is nil")
	}

	q := GetQueue()
	if q == nil {
		return fmt.Errorf("rabbitmq queue is nil")
	}

	return ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func openConnection() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		fmt.Println("RABBITMQ_URL is empty")
		return
	}

	var conn *amqp.Connection
	var err error

	// Retry logic for connection
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			connection = conn
			return
		}
		fmt.Printf("Waiting for rabbitmq... (attempt %d/10): %v\n", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		fmt.Printf("could not connect to rabbitmq after 10 attempts: %v\n", err)
	}
}

func openChannel() {
	if connection == nil {
		return
	}
	ch, err := connection.Channel()
	if err != nil {
		fmt.Printf("error while creating a new channel\n[ERROR] %v\n", err)
		return
	}
	channel = ch
}

func openQueue() {
	if channel == nil {
		return
	}
	q, err := channel.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		amqp.Table{},
	)
	if err != nil {
		fmt.Printf("error trying to declare the queue\n[ERROR] %v\n", err)
		return
	}
	queue = &q
}

func openConsumer() {
	if channel == nil {
		return
	}
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
		fmt.Printf("error trying to get the consumer\n[ERROR] %v\n", err)
		return
	}
	consumer = &msgs
}
