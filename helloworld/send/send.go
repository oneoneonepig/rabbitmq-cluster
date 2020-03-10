package main

import (
	"flag"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	username  = flag.String("username", "admin", "username")
	password  = flag.String("password", "admin", "password")
	host      = flag.String("host", "10.20.131.54", "server address")
	port      = flag.String("port", "5672", "server port")
	queueName = flag.String("name", "two.hello", "queue name")
)

func main() {
	// Create starting timestamp
	start := time.Now()

	// Parse flags
	flag.Parse()

	// Dial connection
	conn, err := amqp.Dial("amqp://" + *username + ":" + *password + "@" + *host + ":" + *port)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare queue
	q, err := ch.QueueDeclare(
		*queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Publish message
	body := "Hello! It's " + time.Now().String()
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

	// Print elasped time
	elapsed := time.Since(start).String()
	log.Println("Time elapsed: " + elapsed)
}
