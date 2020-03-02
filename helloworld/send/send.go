package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	username  string = "guest"
	password  string = "guest"
	host      string = "localhost"
	port      string = "5672"
	queueName string = "ha.hello"
)

func main() {
	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":" + port)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}
