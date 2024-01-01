package basic

import (
	"fmt"
	error2 "github.com/ciazhar/zhar/pkg/error"
	amqp "github.com/rabbitmq/amqp091-go"
)

func New(username, password, host, port string) (*amqp.Connection, *amqp.Channel) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	error2.FailOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Connected to RabbitMQ")

	// Create a channel
	ch, err := conn.Channel()
	error2.FailOnError(err, "Failed to open a channel")
	fmt.Println("Channel created")

	return conn, ch
}

func CreateQueue(ch *amqp.Channel, queueName string) {
	// Declare a queue
	_, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	error2.FailOnError(err, "Failed to declare a queue")
	fmt.Printf("Queue '%s' declared\n", queueName)
}

func ConsumeMessages(ch *amqp.Channel, queueName string, out func(string2 string)) {
	// Consume messages from the queue
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	error2.FailOnError(err, "Failed to register a consumer")

	fmt.Printf("Waiting for messages on queue '%s'. To exit press CTRL+C\n", queueName)

	// Use a goroutine to process incoming messages
	go func() {
		for d := range msgs {
			out(string(d.Body))
		}
	}()
}

func PublishMessage(ch *amqp.Channel, queueName string, message string) {
	// Publish a message to the queue
	err := ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	error2.FailOnError(err, "Failed to publish a message")
	fmt.Printf("Message '%s' published to queue '%s'\n", message, queueName)
}

func PublishMessageWithTTL(ch *amqp.Channel, queueName string, message string, ttlMilliseconds int) {
	// Publish a message to the queue
	err := ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Expiration:  fmt.Sprintf("%d", ttlMilliseconds),
		})
	error2.FailOnError(err, "Failed to publish a message")
	fmt.Printf("Message '%s' published to queue '%s'\n", message, queueName)
}
