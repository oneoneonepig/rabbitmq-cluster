package main

import (
	"log"
	"os"
	"strings"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	username  string = "admin"
	password  string = "admin"
	host      string = "10.20.131.53"
	port      string = "5672"
	queueName string = "ha.hello.pubsub"
)

func main() {
	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":" + port)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs",     // exchange
		"", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2 || os.Args[1] == "") {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
