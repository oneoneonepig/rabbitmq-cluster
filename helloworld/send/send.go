package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

var (
	username  = flag.StringP("username", "u", "admin", "username")
	password  = flag.StringP("password", "p", "admin", "password")
	host      = flag.StringP("host", "h", "localhost:5672", "host address and port")
	queueName = flag.StringP("name", "n", "two.hello", "queue name")
)

func main() {
	// Create starting timestamp
	start := time.Now()

	// Parse flags
	flag.Parse()

	// Dial connection
	conn, err := amqp.Dial("amqp://" + *username + ":" + *password + "@" + *host)
	failOnError(err)
	defer conn.Close()

	// Open channel
	ch, err := conn.Channel()
	failOnError(err)
	defer ch.Close()

	// Declare queue
	q, err := ch.QueueDeclare(
		*queueName, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err)

	// Publish message
	body := "Hello! It's " + time.Now().Format(time.UnixDate)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err)

	// Print information and elasped time
	fmt.Printf("Message sent to queue %s\n", *queueName)
	elapsed := time.Since(start)
	fmt.Printf("Time elapsed: %s\n", elapsed.Truncate(time.Millisecond).String())
}
