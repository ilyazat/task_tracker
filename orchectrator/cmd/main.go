package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	GeneralEventStorageName = "general"
	GoQueueDefaultSize      = 100
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Press Ctrl+C to exit...")
		sig := <-sigs
		fmt.Println()
		fmt.Printf("Received signal: %v\n", sig)
		cancel()
	}()

	connectRabbitMQ, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	_, err = channelRabbitMQ.QueueDeclare(
		GeneralEventStorageName, // queue name
		true,                    // durable
		false,                   // auto delete
		false,                   // exclusive
		false,                   // no wait
		nil,                     // arguments
	)
	if err != nil {
		panic(err)
	}

	//message := amqp.Publishing{
	//	ContentType: "text/plain",
	//	Body:        []byte("FUCK IMPERIALISM"),
	//}
	//
	//if err := channelRabbitMQ.PublishWithContext(ctx,
	//	"",              // exchange
	//	"QueueService1", // queue name
	//	false,           // mandatory
	//	false,           // immediate
	//	message,         // message to publish
	//); err != nil {
	//	log.Fatal(err)
	//}
	messages, err := channelRabbitMQ.Consume(
		GeneralEventStorageName, // queue name
		"",                      // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no local
		false,                   // no wait
		nil,                     // arguments
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	go func() {
		for message := range messages {
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	<-ctx.Done()
	fmt.Println("Exiting main goroutine.")
}
