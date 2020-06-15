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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	// Declare flags variables
	username     = flag.StringP("username", "u", "admin", "username")
	password     = flag.StringP("password", "p", "admin", "password")
	host         = flag.StringP("host", "h", "localhost:5672", "host address and port")
	exchangeName = flag.StringP("name", "n", "pubsub", "exchange name")
	message      = flag.StringP("message", "m", "Hello world!", "message body")
	interval     = flag.StringP("interval", "i", "1s", "interval between messages")
	// Declare normal variables
	// Use int64 to match types when calculating average elapsed time
	msgCount int64 = 0
)

func init() {
	flag.Parse()
}

func main() {
	// Convert interval string to time
	interval_t, err := time.ParseDuration(*interval)
	failOnError(err, "Interval conversion failure")

	// Create starting timestamp
	start := time.Now()

	// Create signal handler
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	// --- Start Loop ---
	fmt.Printf("Sending messages to queue %s every %s\n", *exchangeName, *interval)
	go func() {
		for {
			// Dial connection
			conn, err := amqp.Dial("amqp://" + *username + ":" + *password + "@" + *host)
			failOnError(err, "Failed to connect to RabbitMQ")

			// Open channel
			ch, err := conn.Channel()
			failOnError(err, "Failed to open a channel")

			// Declare exchange
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

			// Publish message
			err = ch.Publish(
				*exchangeName, // exchange
				"",            // routing key
				true,         // mandatory
				false,         // immediate
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(*message),
				})
			failOnError(err, "Failed to publish a message")

			// Add count to msgCount
			msgCount++

			// Close channel and connection
			ch.Close()
			conn.Close()

			// Sleep
			time.Sleep(interval_t)

		}
	}()

	// --- End Loop ---
	<-killSignal
	fmt.Printf("Interrupted\n")

	// Skip and end if no messages were sent
	if msgCount == 0 {
		log.Panicf("No messages were sent. Probably the connection is not established?\n")
	}

	// Print message count, elapsed time and average time per message
	elapsed := time.Since(start)
	avgElapsedInt64 := int64(elapsed) / msgCount
	avgElapsedDuration := time.Duration(avgElapsedInt64)

	fmt.Printf("Messages sent: %d\n", msgCount)
	fmt.Printf("Time elapsed: %s\n", elapsed.Truncate(time.Millisecond).String())
	fmt.Printf("Average per message: %s\n", avgElapsedDuration.Truncate(time.Millisecond).String())

}
