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
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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
	fmt.Println("Start receiving messages...")

	conn, err := amqp.Dial("amqp://" + *username + ":" + *password + "@" + *host)

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		*queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			msgCount++
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
