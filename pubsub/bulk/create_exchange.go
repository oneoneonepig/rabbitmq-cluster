package main

import (
	"fmt"
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
	username string
	password string
	host     string
	name     string
	count    int
)

func init() {
	flag.StringVarP(&username, "username", "u", "admin", "username")
	flag.StringVarP(&password, "password", "p", "admin", "password")
	flag.StringVarP(&host, "host", "h", "localhost:5672", "host address and port")
	flag.StringVarP(&name, "name", "n", "ex", "exchange name")
	flag.IntVarP(&count, "count", "c", 1, "number of exchanges to create")
	flag.Parse()
}

func main() {
	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	for i := 1; i <= count; i++ {
		exchangeName := name + "_" + fmt.Sprintf("%02d", i)
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
		log.Println("Exchange " + exchangeName + " created.")
	}

}
