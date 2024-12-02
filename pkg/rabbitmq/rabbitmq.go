package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ciazhar/go-start-small/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

const (
	maxRetryAttempts = 5
	retryDelay       = 2 * time.Second
)

func New(connectionName, username, password, host, port string) *RabbitMQ {

	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for i := 0; i < maxRetryAttempts; i++ {
		config := amqp.Config{Properties: amqp.NewConnectionProperties()}
		config.Properties.SetClientConnectionName(connectionName)

		conn, err = amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port), config)
		if err != nil {
			logger.LogError(context.Background(), err, "Failed to connect to RabbitMQ", map[string]interface{}{
				"username": username,
				"host":     host,
				"port":     port,
				"attempt":  i + 1,
			})
			time.Sleep(retryDelay)
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			logger.LogError(context.Background(), err, "Failed to open a channel", map[string]interface{}{
				"username": username,
				"host":     host,
				"port":     port,
				"attempt":  i + 1,
			})
			time.Sleep(retryDelay)
			conn.Close() // Close the connection if we fail to open a channel
			continue
		}

		logger.LogInfo(context.Background(), "Connected to RabbitMQ", nil)
		break
	}

	if conn == nil || ch == nil {
		logger.LogFatal(context.Background(), err, "Failed to connect to RabbitMQ", map[string]interface{}{
			"username": username,
			"host":     host,
			"port":     port,
			"attempt":  maxRetryAttempts,
		})
	}

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
	}
}

func (r *RabbitMQ) CreateQueue(queueName string) {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to declare queue", map[string]interface{}{"queue": queueName})
	}
	logger.LogInfo(context.Background(), "Created queue", map[string]interface{}{"queue": queueName})
}

func (r *RabbitMQ) CreateQueueDelay(exchange, queue, routingKey string) {
	args := amqp.Table{
		"x-delayed-type": "direct",
	}
	err := r.channel.ExchangeDeclare(
		exchange,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to declare exchange", map[string]interface{}{
			"exchange":   exchange,
			"queue":      queue,
			"routingKey": routingKey,
		})

	}

	if _, err = r.channel.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		logger.LogFatal(context.Background(), err, "Failed to declare queue", map[string]interface{}{
			"queue":      queue,
			"exchange":   exchange,
			"routingKey": routingKey,
		})
	}

	err = r.channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to bind queue", map[string]interface{}{
			"queue":      queue,
			"exchange":   exchange,
			"routingKey": routingKey,
		})
	}
}

func (r *RabbitMQ) CreateRoutingKey(queue, routingKey, exchange string) {
	err := r.channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to bind queue", map[string]interface{}{"queue": queue, "routingKey": routingKey, "exchange": exchange})
	}
}

func (r *RabbitMQ) ConsumeMessages(ctx context.Context, queueName string, handler func(msg string)) {
	for i := 0; i < maxRetryAttempts; i++ {

		messages, err := r.channel.Consume(
			queueName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logger.LogFatal(context.Background(), err, "Failed to register consumer", map[string]interface{}{
				"queue":   queueName,
				"attempt": i + 1,
			})
			time.Sleep(retryDelay)
			continue
		}
		logger.LogInfo(context.Background(), "Consumer registered on queue", map[string]interface{}{"queue": queueName})

		for {
			select {
			case msg, ok := <-messages:
				if !ok {
					logger.LogInfo(context.Background(), "Consumer channel for queue closed", map[string]interface{}{"queue": queueName})
					break
				}
				handler(string(msg.Body))
			case <-ctx.Done():
				logger.LogInfo(context.Background(), "Context done, stopping consumer for queue", map[string]interface{}{"queue": queueName})
				return
			}
		}
	}
	logger.LogError(context.Background(), nil, "Failed to consume from queue", map[string]interface{}{"queue": queueName, "attempt": maxRetryAttempts})
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, queueName string, message string, ttlMilliseconds ...int) {
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	if len(ttlMilliseconds) > 0 {
		publishing.Expiration = fmt.Sprintf("%d", ttlMilliseconds[0])
	}

	if err := r.retryPublish(ctx, "", queueName, publishing); err != nil {
		logger.LogError(context.Background(), err, "Failed to publish message to queue", map[string]interface{}{"queue": queueName})
	}
}

func (r *RabbitMQ) retryPublish(ctx context.Context, exchange, key string, publishing amqp.Publishing) error {
	var err error
	for i := 0; i < maxRetryAttempts; i++ {
		err = r.channel.PublishWithContext(ctx, exchange, key, false, false, publishing)
		if err == nil {
			return nil
		}
		logger.LogError(context.Background(), err, "Failed to publish message to exchange", map[string]interface{}{
			"exchange":   exchange,
			"routingKey": key,
			"attempt":    i + 1,
		})
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("failed to publish message after %d attempts: %w", maxRetryAttempts, err)
}

func (r *RabbitMQ) PublishDelayedMessage(ctx context.Context, routingKey, message, exchange string, delay time.Duration) {
	headers := amqp.Table{
		"x-delay": int64(delay / time.Millisecond),
	}
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		Headers:     headers,
	}

	if err := r.retryPublishDelayed(ctx, routingKey, exchange, routingKey, publishing); err != nil {
		logger.LogError(context.Background(), err, "Failed to publish delayed message to exchange", map[string]interface{}{"exchange": exchange, "routingKey": routingKey})
	}
}

func (r *RabbitMQ) retryPublishDelayed(ctx context.Context, routingKey, exchange, key string, publishing amqp.Publishing) error {
	var err error
	for i := 0; i < maxRetryAttempts; i++ {
		err = r.channel.PublishWithContext(ctx, exchange, routingKey, false, false, publishing)
		if err == nil {
			return nil
		}
		logger.LogError(context.Background(), err, "Failed to publish message to exchange", map[string]interface{}{
			"exchange":   exchange,
			"routingKey": key,
			"attempt":    i + 1,
		})
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("failed to publish message after %d attempts: %w", maxRetryAttempts, err)
}

func (r *RabbitMQ) Close() {
	if err := r.channel.Close(); err != nil {
		logger.LogError(context.Background(), err, "Failed to close channel", nil)
	}
	if err := r.connection.Close(); err != nil {
		logger.LogError(context.Background(), err, "Failed to close connection", nil)
	}
}

type ConsumerConfig struct {
	Queue   string
	Handler func(message string)
}

func (r *RabbitMQ) StartConsumers(ctx context.Context, consumers []ConsumerConfig, wg *sync.WaitGroup) {
	for _, config := range consumers {
		wg.Add(1)
		go func(config ConsumerConfig) {
			defer wg.Done()
			r.ConsumeMessages(ctx, config.Queue, config.Handler)
		}(config)
		logger.LogInfo(ctx, "Started consumer for queue", map[string]interface{}{"queue": config.Queue})
	}
}
