package main

import (
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	username     = flag.StringP("username", "u", "admin", "username")
	password     = flag.StringP("password", "p", "admin", "password")
	host         = flag.StringP("host", "h", "localhost:5672", "host address and port")
	exchangeName = flag.StringP("name", "n", "pubsub", "exchange name")
	message      = flag.StringP("message", "m", "Hello world!", "message body")
)

func init() {
	flag.Parse()
}

func main() {
	conn, err := amqp.Dial("amqp://" + *username + ":" + *password + "@" + *host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		*exchangeName, // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		*exchangeName, // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(*message),
		})
	failOnError(err, "Failed to publish a message")
}
