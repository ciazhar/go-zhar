package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg"
	"github.com/ciazhar/go-zhar/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type ConnectionPool struct {
	logger      logger.Logger
	mu          sync.Mutex
	connections []*amqp.Connection
}

func NewConnectionPool(connectionName, username, password, host, port string, size int, logger logger.Logger) *ConnectionPool {

	pool := &ConnectionPool{}
	for i := 0; i < size; i++ {

		config := amqp.Config{Properties: amqp.NewConnectionProperties()}
		config.Properties.SetClientConnectionName(connectionName + "-" + fmt.Sprintf("%d", i))

		conn, err := amqp.Dial(fmt.Sprintf(AmqpUrl, username, password, host, port))
		if err != nil {
			logger.Fatalf(ErrConnFailed, err)
		}
		pool.connections = append(pool.connections, conn)
	}
	return pool
}

func (r *ConnectionPool) Get() (*amqp.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.connections) == 0 {
		return nil, fmt.Errorf(MsgConnectionPoolEmpty)
	}

	conn := r.connections[0]
	r.connections = r.connections[1:]
	return conn, nil
}

func (r *ConnectionPool) Put(conn *amqp.Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.connections = append(r.connections, conn)
}

func (r *ConnectionPool) CreateQueue(queueName string) {
	conn, err := r.Get()
	if err != nil {
		r.logger.Fatalf(ErrGetChannelFromPool, err)
	}
	defer r.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		r.logger.Fatalf(ErrChanFailed, err)
	}
	r.logger.Info(MsgChanCreated)

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

func (r *ConnectionPool) ConsumeMessages(queueName string, out func(string2 string)) {
	conn, err := r.Get()
	if err != nil {
		r.logger.Fatalf(ErrGetChannelFromPool, err)
	}
	defer r.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		r.logger.Fatalf(ErrChanFailed, err)
	}
	r.logger.Info(MsgChanCreated)

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
	if err != nil {
		r.logger.Fatalf(ErrConsumerFailed, err)
	}
	r.logger.Infof(MsgConsumerSucceed, queueName)

	// Use a goroutine to process incoming messages
	go func() {
		for d := range msgs {
			out(string(d.Body))
		}
	}()
}

func (r *ConnectionPool) PublishMessage(ctx context.Context, queueName string, message string) {
	conn, err := r.Get()
	if err != nil {
		r.logger.Fatalf(ErrGetChannelFromPool, err)
	}
	defer r.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		r.logger.Fatalf(ErrChanFailed, err)
	}
	r.logger.Info(MsgChanCreated)

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
		r.logger.Fatalf(ErrProducerFailed, err)
	}
	r.logger.Infof(MsgProducerSucceed, message, queueName)
}

func (r *ConnectionPool) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	conn, err := r.Get()
	if err != nil {
		r.logger.Fatalf(ErrGetChannelFromPool, err)
	}
	defer r.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		r.logger.Fatalf(ErrChanFailed, err)
	}
	r.logger.Info(MsgChanCreated)

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
		r.logger.Fatalf(ErrProducerFailed, err)
	}
	r.logger.Infof(MsgProducerSucceed, message, queueName)
}

func (r *ConnectionPool) Close() {
	defer func() {
		for i := range r.connections {
			err := r.connections[i].Close()
			if err != nil {
				r.logger.Fatalf(ErrClosingChannel, err)
			}
		}
	}()
}
