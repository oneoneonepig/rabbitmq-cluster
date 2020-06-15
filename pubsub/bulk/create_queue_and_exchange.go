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
	queueCount    int
)

func init() {
	//	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	//	CommandLine.StringVarP(&username, "username", "u", "admin", "username")
	//	CommandLine.StringVarP(&password, "password", "p", "admin", "password")
	//	CommandLine.StringVarP(&host, "host", "h", "localhost:5672", "host address and port")
	//	CommandLine.StringVarP(&name, "name", "n", "ex", "exchange name")
	//	CommandLine.IntVarP(&exchangeCount, "exchange-count", "e", 1, "number of exchanges")
	//	CommandLine.IntVarP(&queueCount, "queue-count", "q", 1, "number of queues per exchange")
	//	CommandLine.SortFlags = false
	//	CommandLine.Parse(os.Args)
	flag.StringVarP(&username, "username", "u", "admin", "username")
	flag.StringVarP(&password, "password", "p", "admin", "password")
	flag.StringVarP(&host, "host", "h", "localhost:5672", "host address and port")
	flag.StringVarP(&name, "name", "n", "ex", "exchange name")
	flag.IntVarP(&exchangeCount, "exchange-count", "e", 1, "number of exchanges")
	flag.IntVarP(&queueCount, "queue-count", "q", 1, "number of queues per exchange")
	flag.Parse()
}

func main() {
	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Create  exchanges
	for i := 1; i <= exchangeCount; i++ {
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

		for j := 1; j <= queueCount; j++ {
			queueName := exchangeName + "_" + fmt.Sprintf("%02d", j)
			_, err := ch.QueueDeclare(
				queueName, // name
				true,      // durable
				false,     // delete when unused
				false,     // exclusive
				false,     // no-wait
				nil,       // arguments
			)
			failOnError(err, "Failed to declare a queue")

			err = ch.QueueBind(
				queueName,    // queue name
				"",           // routing key
				exchangeName, // exchange
				false,        // no-wait
				nil,
			)
			failOnError(err, "Failed to bind a queue")

			log.Printf("Queue %s created, binded with %s", queueName, exchangeName)
		}
	}

}
