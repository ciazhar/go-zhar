package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMQ struct {
	logger     logger.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(connectionName, username, password, host, port string, logger logger.Logger) *RabbitMQ {

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(connectionName)

	conn, err := amqp.DialConfig(fmt.Sprintf(AmqpUrl, username, password, host, port), config)
	if err != nil {
		logger.Fatalf(ErrConnFailed, err)
	}
	logger.Info(MsgConnSucceed)

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf(ErrChanFailed, err)
	}
	logger.Info(MsgChanCreated)

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
		r.logger.Fatalf(ErrQueueFailed, err)
	}
	r.logger.Infof(MsgQueueCreated, queueName)
}

func (r *RabbitMQ) ConsumeMessages(queueName string, out func(msg string), stop chan struct{}) {
	messages, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		r.logger.Fatalf(ErrConsumerFailed, err)
	}
	r.logger.Infof(MsgConsumerSucceed, queueName)

	go func() {
		for msg := range messages {
			select {
			case <-stop:
				r.logger.Info(MsgConsumerStopped)
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
		r.logger.Infof(ErrProducerFailed, err)
	}

}

func (r *RabbitMQ) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		Expiration:  fmt.Sprintf("%d", ttlMilliseconds),
	}

	if err := r.channel.PublishWithContext(ctx, "", queueName, false, false, publishing); err != nil {
		r.logger.Infof(ErrProducerFailed, err)
	}
}

func (r *RabbitMQ) Close() {
	defer func() {
		err := r.channel.Close()
		if err != nil {
			r.logger.Fatalf(ErrClosingChannel, err)
		}

		err = r.connection.Close()
		if err != nil {
			r.logger.Fatalf(ErrClosingConnection, err)
		}
	}()
}
