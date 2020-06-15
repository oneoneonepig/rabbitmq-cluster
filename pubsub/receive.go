package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	username     string
	password     string
	host         string
	exchangeName string
	queueName    string
	outputDot    bool
	autoAck      bool
)

func init() {
	flag.StringVarP(&username, "username", "u", "admin", "username")
	flag.StringVarP(&password, "password", "p", "admin", "password")
	flag.StringVarP(&host, "host", "h", "localhost:5672", "host address and port")
	flag.StringVarP(&exchangeName, "exchange", "e", "pubsub", "exchange name")
	flag.StringVarP(&queueName, "queue", "q", "", "queue name, generate new one if not provided")
	flag.BoolVarP(&outputDot, "dot", "d", false, "print a single dot for every message received")
	flag.BoolVarP(&autoAck, "auto-ack", "a", false, "auto-ack when message consumed")
	flag.Parse()
}

func main() {
	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host)

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	if queueName == "" {
		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		failOnError(err, "Failed to declare a queue")
		queueName = q.Name
		err = ch.QueueBind(
			queueName,    // queue name
			"",           // routing key
			exchangeName, // exchange
			false,
			nil,
		)
		failOnError(err, "Failed to bind a queue")
	}
	log.Printf("Consuming queue: %s", queueName)
	log.Printf("Queue binded to: %s", exchangeName)

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if autoAck {
				if outputDot {
					fmt.Printf(".")
				} else {
					log.Printf(" [x] %s", d.Body)
				}
			} else {
				rand.Seed(time.Now().UnixNano())
				n := time.Duration(rand.Intn(1000)) * time.Millisecond
				time.Sleep(n)
				if outputDot {
					fmt.Printf(".")
				} else {
					log.Printf(" [x] %s | %s", n, d.Body)
				}
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
