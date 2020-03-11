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
	interval  = flag.StringP("interval", "i", "1s", "interval between messages")

	// Declare normal variables
	// Use int64 to match types when calculating average elapsed time
	msgCount int64 = 0
)

func main() {
	// Create starting timestamp
	start := time.Now()

	// Parse flags
	flag.Parse()
	interval_t, err := time.ParseDuration(*interval)
	failOnError(err, "Failed to convert interval to time.Duration")

	// Create signal handler
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	// --- Start Loop ---
	fmt.Printf("Sending messages to queue %s every %s\n", *queueName, *interval)
	go func() {
		for {
			// Dial connection
			conn, err := amqp.Dial("amqp://" + *username + ":" + *password + "@" + *host)
			failOnError(err, "Failed to connect to RabbitMQ")
			defer conn.Close()

			// Open channel
			ch, err := conn.Channel()
			failOnError(err, "Failed to open a channel")
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

			// Add count to msgCount
			msgCount++

			// Sleep
			time.Sleep(interval_t)
		}
	}()

	// --- End Loop ---

	<-killSignal
	fmt.Printf("Interrupted\n")

	// Skip and end if no messages were sent
	if msgCount == 0 {
		log.Fatalf("No messages were sent. Probably the connection is not established?\n")
	}

	// Print message count, elapsed time and average time per message
	elapsed := time.Since(start)
	avgElapsedInt64 := int64(elapsed) / msgCount
	avgElapsedDuration := time.Duration(avgElapsedInt64)

	fmt.Printf("Messages sent: %d\n", msgCount)
	fmt.Printf("Time elapsed: %s\n", elapsed.Truncate(time.Millisecond).String())
	fmt.Printf("Average per message: %s\n", avgElapsedDuration.Truncate(time.Millisecond).String())
}
