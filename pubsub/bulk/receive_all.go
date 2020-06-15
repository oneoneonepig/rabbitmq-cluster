package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
	// "github.com/pkg/profile"
	"os"
	"os/signal"
	"syscall"
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
	outputDot     bool
	autoAck       bool
	msgCount      int64
)

func init() {
	flag.StringVarP(&username, "username", "u", "admin", "username")
	flag.StringVarP(&password, "password", "p", "admin", "password")
	flag.StringVarP(&host, "host", "h", "localhost:5672", "host address and port")
	flag.StringVarP(&name, "name", "n", "ex", "exchange name")
	flag.IntVarP(&exchangeCount, "exchange-count", "e", 1, "number of exchanges")
	flag.IntVarP(&queueCount, "queue-count", "q", 1, "number of queues per exchange")
	flag.BoolVarP(&outputDot, "dot", "d", false, "print a single dot for every message received")
	flag.BoolVarP(&autoAck, "auto-ack", "a", false, "auto-ack when message consumed.If not set, comsumer will delay a random duration under 1 second to simulate processing messages.")
	flag.Parse()
}

func main() {
	// defer profile.Start(profile.MemProfile).Stop()

	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	for i := 1; i <= exchangeCount; i++ {
		exchangeName := name + "_" + fmt.Sprintf("%02d", i)
		for j := 1; j <= queueCount; j++ {
			queueName := exchangeName + "_" + fmt.Sprintf("%02d", j)
			msgs, err := ch.Consume(
				queueName, // queue
				"",        // consumer
				false,     // auto-ack
				false,     // exclusive
				false,     // no-local
				false,     // no-wait
				nil,       // args
			)
			failOnError(err, "Failed to register a consumer")
			go func() {
				//fmt.Printf("i: %d, j: %d\n", i, j)
				//fmt.Printf("queueName: %s\n", queueName)
				for d := range msgs {
					if autoAck {
						if outputDot {
							fmt.Printf(".")
						} else {
							log.Printf(" [x] %s", d.Body)
						}
					} else {
						rand.Seed(time.Now().UnixNano())
						n := time.Duration(rand.Intn(1000)) * time.Millisecond
						time.Sleep(n)
						if outputDot {
							fmt.Printf(".")
						} else {
							log.Printf(" [x] %s | %s | %s", queueName, n, d.Body)
						}
					}
					msgCount++
					d.Ack(false)
				}
			}()
		}
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-killSignal
	log.Printf("Interrupted\n")
	log.Printf("Messages received: %d\n", msgCount)
}
