package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"log"
	// "os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	username      string
	password      string
	host          string
	name          string
	exchangeCount int
	message       string
)

func init() {
	flag.StringVarP(&username, "username", "u", "admin", "username")
	flag.StringVarP(&password, "password", "p", "admin", "password")
	flag.StringVarP(&host, "host", "h", "localhost:5672", "host address and port")
	flag.StringVarP(&name, "name", "n", "ex", "exchange name")
	flag.IntVarP(&exchangeCount, "count", "c", 1, "number of exchanges")
	flag.StringVarP(&message, "message", "m", "Hello world!", "message body")
	flag.Parse()
}

func main() {
	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	for i := 1; i <= exchangeCount; i++ {
		exchangeName := name + "_" + fmt.Sprintf("%02d", i)
		err = ch.Publish(
			exchangeName, // exchange
			"",            // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(message),
			})
		failOnError(err, "Failed to publish a message")
	}
}
