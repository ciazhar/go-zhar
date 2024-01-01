package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/zhar/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func New(username, password, host, port string) *RabbitMQ {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(fmt.Sprintf(AmqpUrl, username, password, host, port))
	pkg.FailOnError(err, ErrConnFailed)
	fmt.Println(MsgConnSucceed)

	// Create a channel
	ch, err := conn.Channel()
	pkg.FailOnError(err, ErrChanFailed)
	fmt.Println(MsgChanCreated)

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}
}

func (r *RabbitMQ) CreateQueue(queueName string) {
	// Declare a queue
	_, err := r.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	pkg.FailOnError(err, ErrQueueFailed)
	fmt.Printf(MsgQueueCreated, queueName)
}

func (r *RabbitMQ) ConsumeMessages(queueName string, out func(string2 string)) {
	// Consume messages from the queue
	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	pkg.FailOnError(err, ErrConsumerFailed)
	fmt.Printf(MsgConsumerSucceed, queueName)

	// Use a goroutine to process incoming messages
	go func() {
		for d := range msgs {
			out(string(d.Body))
		}
	}()
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, queueName string, message string) {
	// Publish a message to the queue
	err := r.Channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: pkg.TextPlain,
			Body:        []byte(message),
		})
	pkg.FailOnError(err, ErrProducerFailed)
	fmt.Printf(MsgProducerSucceed, message, queueName)
}

func (r *RabbitMQ) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	// Publish a message to the queue
	err := r.Channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: pkg.TextPlain,
			Body:        []byte(message),
			Expiration:  fmt.Sprintf("%d", ttlMilliseconds),
		})
	pkg.FailOnError(err, ErrProducerFailed)
	fmt.Printf(MsgProducerSucceed, message, queueName)
}

func (r *RabbitMQ) Close() {
	defer func() {
		err := r.Channel.Close()
		pkg.FailOnError(err, ErrClosingChannel)

		err = r.Connection.Close()
		pkg.FailOnError(err, ErrClosingConnection)
	}()
}
