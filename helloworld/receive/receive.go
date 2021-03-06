package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Helper function: error handling
func failOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

var (
	// Declare flags variables
	username  = flag.StringP("username", "u", "admin", "username")
	password  = flag.StringP("password", "p", "admin", "password")
	host      = flag.StringP("host", "h", "localhost:5672", "host address and port")
	queueName = flag.StringP("name", "n", "two.hello", "queue name")

	// Declare normal variables
	// Use int64 to match types when calculating average elapsed time
	msgCount int64 = 0
)

func main() {

	// Create starting timestamp
	start := time.Now()

	// Parse flags
	flag.Parse()

	// Create signal handler
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	// --- Start receiving messages ---
	fmt.Printf("Receiving messages from queue %s\n", *queueName)

	conn, err := amqp.Dial("amqp://" + *username + ":" + *password + "@" + *host)

	failOnError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		*queueName, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			msgCount++
			d.Ack(false)
		}
	}()

	// --- End receiving messages ---

	<-killSignal
	fmt.Printf("Interrupted\n")

	// Print message count and elapsed time
	elapsed := time.Since(start)

	fmt.Printf("Messages received: %d\n", msgCount)
	fmt.Printf("Time elapsed: %s\n", elapsed.Truncate(time.Millisecond).String())
}
