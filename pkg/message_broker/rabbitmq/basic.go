package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(username, password, host, port string) *RabbitMQ {

	conn, err := amqp.Dial(fmt.Sprintf(AmqpUrl, username, password, host, port))
	if err != nil {
		log.Fatal(ErrConnFailed, err)
	}
	log.Println(MsgConnSucceed)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(ErrChanFailed, err)
	}
	log.Println(MsgChanCreated)

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
	}
}

func (r *RabbitMQ) CreateQueue(queueName string) {
	// Declare a queue
	if _, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		log.Fatal(ErrQueueFailed, err)
	}
	log.Printf(MsgQueueCreated, queueName)
}

func (r *RabbitMQ) ConsumeMessages(queueName string, out func(msg string), stop chan struct{}) {
	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatal(ErrConsumerFailed, err)
	}
	log.Printf(MsgConsumerSucceed, queueName)

	go func() {
		for msg := range msgs {
			select {
			case <-stop:
				log.Println(MsgConsumerStopped)
				// Perform any cleanup logic here
				time.Sleep(2 * time.Second) // Simulate cleanup
				close(stop)
				return
			default:
				out(string(msg.Body))
			}
		}
	}()
	<-stop
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, queueName string, message string) {

	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	if err := r.channel.PublishWithContext(ctx, "", queueName, false, false, publishing); err != nil {
		log.Fatal(ErrProducerFailed, err)
	}

}

func (r *RabbitMQ) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		Expiration:  fmt.Sprintf("%d", ttlMilliseconds),
	}

	if err := r.channel.PublishWithContext(ctx, "", queueName, false, false, publishing); err != nil {
		log.Fatal(ErrProducerFailed, err)
	}
}

func (r *RabbitMQ) Close() {
	defer func() {
		err := r.channel.Close()
		pkg.FailOnError(err, ErrClosingChannel)

		err = r.connection.Close()
		pkg.FailOnError(err, ErrClosingConnection)
	}()
}
