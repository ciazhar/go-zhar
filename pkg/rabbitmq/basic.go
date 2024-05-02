package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
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

	conn, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port), config)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("Failed to open a channel: %s", err)
	}
	logger.Info("Connected to RabbitMQ")

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
		logger:     logger,
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
		r.logger.Fatalf("Failed to declare a queue: %s", err)
	}
	r.logger.Infof("Queue %s created", queueName)
}

func (r *RabbitMQ) CreateQueueDelay(exchange, queue, routingKey string) {
	args := amqp.Table{
		"x-delayed-type": "direct",
	}
	err := r.channel.ExchangeDeclare(
		exchange,            // exchange name
		"x-delayed-message", // exchange type
		true,                // durable
		false,               // auto-delete
		false,               // internal
		false,               // no-wait
		args,                // arguments
	)
	if err != nil {
		r.logger.Fatalf("Failed to declare an exchange: %s", err)
	}

	// Declare the delayed queue
	_, err = r.channel.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		r.logger.Fatalf("Failed to declare a queue: %s", err)
	}

	err = r.channel.QueueBind(
		queue,      // queue name
		routingKey, // routing key (converting integer to string)
		exchange,   // exchange name
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		r.logger.Fatalf("Failed to bind queue: %s", err)
	}
}

func (r *RabbitMQ) CreateRoutingKey(queue, routingKey, exchange string) {
	err := r.channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		r.logger.Fatalf("Failed to bind queue: %s", err)
	}
}

func (r *RabbitMQ) ConsumeMessages(ctx context.Context, queueName string, out func(msg string)) {
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
		r.logger.Fatalf("Failed to register a consumer: %s", err)
	}
	r.logger.Infof("Consumer registered on queue %s", queueName)

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				r.logger.Infof("Consumer for queue %s closed", queueName)
				return
			}
			out(string(msg.Body))
		case <-ctx.Done():
			r.logger.Infof("Stopping consumer for queue %s", queueName)
			return
		}
	}
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, queueName string, message string) {

	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	if err := r.channel.PublishWithContext(ctx, "", queueName, false, false, publishing); err != nil {
		r.logger.Infof("Failed to publish a message: %s", err)
	}
}

func (r *RabbitMQ) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		Expiration:  fmt.Sprintf("%d", ttlMilliseconds),
	}

	if err := r.channel.PublishWithContext(ctx, "", queueName, false, false, publishing); err != nil {
		r.logger.Infof("Failed to publish a message: %s", err)
	}
}

func (r *RabbitMQ) PublishDelayedMessage(ctx context.Context, routingKey string, message string, exchange string, delay time.Duration) {

	r.logger.Infof("Publishing delayed message: %s", message)

	//Calculate the delay in milliseconds
	delayMillis := int64(delay / time.Millisecond)

	// Set the AMQP headers with x-delay
	headers := amqp.Table{
		"x-delay": delayMillis,
	}

	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		Headers:     headers,
	}

	if err := r.channel.PublishWithContext(ctx, exchange, routingKey, false, false, publishing); err != nil {
		r.logger.Infof("Failed to publish a message: %s", err)
	}
}

func (r *RabbitMQ) Close() {
	defer func() {
		err := r.channel.Close()
		if err != nil {
			r.logger.Fatalf("Failed to close channel: %s", err)
		}

		err = r.connection.Close()
		if err != nil {
			r.logger.Fatalf("Failed to close connection: %s", err)
		}
	}()
}

type ConsumerConfig struct {
	Queue   string
	Handler func(message string)
}

func (r *RabbitMQ) StartConsumers(
	ctx context.Context,
	consumers []ConsumerConfig,
	wg *sync.WaitGroup,
	logger logger.Logger,
) {
	// Iterate over the consumers map
	for _, config := range consumers {

		wg.Add(1)

		go func(config ConsumerConfig) {
			defer wg.Done()
			r.ConsumeMessages(ctx, config.Queue, config.Handler)
		}(config)

		logger.Infof("Starting consumer queue %s", config.Queue)
	}
}
