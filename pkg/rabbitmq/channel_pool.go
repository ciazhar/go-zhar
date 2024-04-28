package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg"
	"github.com/ciazhar/go-zhar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type ChannelPool struct {
	logger     logger.Logger
	mu         sync.Mutex
	Connection *amqp.Connection
	channels   []*amqp.Channel
}

func NewChannelPool(connectionName, username, password, host, port string, size int, logger logger.Logger) *ChannelPool {
	pool := &ChannelPool{}

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(connectionName)

	conn, err := amqp.DialConfig(fmt.Sprintf(AmqpUrl, username, password, host, port), config)
	if err != nil {
		logger.Fatalf(ErrConnFailed, err)
	}
	for i := 0; i < size; i++ {

		// Create a channel
		ch, err := conn.Channel()
		if err != nil {
			logger.Fatalf(ErrChanFailed, err)
		}
		logger.Info(MsgChanCreated)

		pool.channels = append(pool.channels, ch)
	}
	return pool
}

func (r *ChannelPool) Get() (*amqp.Channel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.channels) == 0 {
		return nil, fmt.Errorf(MsgConnectionPoolEmpty)
	}

	conn := r.channels[0]
	r.channels = r.channels[1:]
	return conn, nil
}

func (r *ChannelPool) Put(conn *amqp.Channel) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.channels = append(r.channels, conn)
}

func (r *ChannelPool) CreateQueue(queueName string) {
	ch, err := r.Get()
	if err != nil {
		r.logger.Fatalf(ErrGetChannelFromPool, err)
	}
	defer r.Put(ch)

	// Declare a queue
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		r.logger.Fatalf(ErrQueueFailed, err)
	}
	r.logger.Infof(MsgQueueCreated, queueName)
}

func (r *ChannelPool) ConsumeMessages(queueName string, out func(msg string), stop chan struct{}) {
	ch, err := r.Get()
	if err != nil {
		r.logger.Fatalf(ErrGetChannelFromPool, err)
	}
	defer r.Put(ch)

	messages, err := ch.Consume(
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

func (r *ChannelPool) PublishMessage(ctx context.Context, queueName string, message string) {
	ch, err := r.Get()
	if err != nil {
		r.logger.Infof(ErrGetChannelFromPool, err)
	}
	defer r.Put(ch)

	// Publish a message to the queue
	err = ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: pkg.TextPlain,
			Body:        []byte(message),
		})
	if err != nil {
		r.logger.Infof(ErrProducerFailed, err)
		return
	}
	r.logger.Infof(MsgProducerSucceed, message, queueName)
}

func (r *ChannelPool) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	ch, err := r.Get()
	if err != nil {
		r.logger.Fatalf(ErrGetChannelFromPool, err)
	}
	defer r.Put(ch)

	// Publish a message to the queue
	err = ch.PublishWithContext(
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
	if err != nil {
		r.logger.Infof(ErrProducerFailed, err)
		return
	}
	r.logger.Infof(MsgProducerSucceed, message, queueName)
}

func (r *ChannelPool) Close() {
	defer func() {
		for i := range r.channels {
			err := r.channels[i].Close()
			if err != nil {
				r.logger.Fatalf(ErrClosingChannel, err)
			}
		}

		err := r.Connection.Close()
		if err != nil {
			r.logger.Fatalf(ErrClosingConnection, err)
		}
	}()
}
