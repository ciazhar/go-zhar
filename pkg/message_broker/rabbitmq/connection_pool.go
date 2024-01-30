package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type ConnectionPool struct {
	mu          sync.Mutex
	connections []*amqp.Connection
}

func NewConnectionPool(username, password, host, port string, size int) *ConnectionPool {
	pool := &ConnectionPool{}
	for i := 0; i < size; i++ {
		conn, err := amqp.Dial(fmt.Sprintf(AmqpUrl, username, password, host, port))
		pkg.FailOnError(err, ErrConnFailed)
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
	pkg.FailOnError(err, ErrGetChannelFromPool)
	defer r.Put(conn)

	ch, err := conn.Channel()
	pkg.FailOnError(err, ErrChanFailed)
	fmt.Println(MsgChanCreated)

	// Declare a queue
	_, err = ch.QueueDeclare(
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

func (r *ConnectionPool) ConsumeMessages(queueName string, out func(string2 string)) {
	conn, err := r.Get()
	pkg.FailOnError(err, ErrGetChannelFromPool)
	defer r.Put(conn)

	ch, err := conn.Channel()
	pkg.FailOnError(err, ErrChanFailed)
	fmt.Println(MsgChanCreated)

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
	pkg.FailOnError(err, ErrConsumerFailed)
	fmt.Printf(MsgConsumerSucceed, queueName)

	// Use a goroutine to process incoming messages
	go func() {
		for d := range msgs {
			out(string(d.Body))
		}
	}()
}

func (r *ConnectionPool) PublishMessage(ctx context.Context, queueName string, message string) {
	conn, err := r.Get()
	pkg.FailOnError(err, ErrGetChannelFromPool)
	defer r.Put(conn)

	ch, err := conn.Channel()
	pkg.FailOnError(err, ErrChanFailed)
	fmt.Println(MsgChanCreated)

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
	pkg.FailOnError(err, ErrProducerFailed)
	fmt.Printf(MsgProducerSucceed, message, queueName)
}

func (r *ConnectionPool) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	conn, err := r.Get()
	pkg.FailOnError(err, ErrGetChannelFromPool)
	defer r.Put(conn)

	ch, err := conn.Channel()
	pkg.FailOnError(err, ErrChanFailed)
	fmt.Println(MsgChanCreated)

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
	pkg.FailOnError(err, ErrProducerFailed)
	fmt.Printf(MsgProducerSucceed, message, queueName)
}

func (r *ConnectionPool) Close() {
	defer func() {
		for i := range r.connections {
			err := r.connections[i].Close()
			pkg.FailOnError(err, ErrClosingConnection)
		}
	}()
}
