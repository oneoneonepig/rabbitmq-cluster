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
	"strconv"
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
	interval      string
	msgCount      int64 = 0
)

func init() {
	flag.StringVarP(&username, "username", "u", "admin", "username")
	flag.StringVarP(&password, "password", "p", "admin", "password")
	flag.StringVarP(&host, "host", "h", "localhost:5672", "host address and port")
	flag.StringVarP(&name, "name", "n", "ex", "exchange name")
	flag.IntVarP(&exchangeCount, "count", "c", 1, "number of exchanges")
	flag.StringVarP(&message, "message", "m", "Hello world!", "message body")
	flag.StringVarP(&interval, "interval", "i", "1s", "interval between messages")
	flag.Parse()
}

func main() {

	// Convert interval string to time
	interval_t, err := time.ParseDuration(interval)
	failOnError(err, "Interval conversion failure")

	// Create starting timestamp
	start := time.Now()

	// Create signal handler
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	// --- Start Loop ---
	fmt.Printf("Sending messages to all queues every %s\n", interval)
	go func() {

		conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		for {
			for i := 1; i <= exchangeCount; i++ {
				exchangeName := name + "_" + fmt.Sprintf("%02d", i)
				err = ch.Publish(
					exchangeName, // exchange
					"",           // routing key
					false,        // mandatory
					false,        // immediate
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Body:         []byte(strconv.FormatInt(msgCount, 10) + ":" + message),
					})
				failOnError(err, "Failed to publish a message")

				// Increase count and sleep
				msgCount++
				time.Sleep(interval_t)
			}
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
