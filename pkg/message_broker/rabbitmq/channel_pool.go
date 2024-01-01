package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ciazhar/zhar/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type ChannelPool struct {
	mu         sync.Mutex
	Connection *amqp.Connection
	channels   []*amqp.Channel
}

func NewChannelPool(username, password, host, port string, size int) *ChannelPool {
	pool := &ChannelPool{}
	conn, err := amqp.Dial(fmt.Sprintf(AmqpUrl, username, password, host, port))
	pkg.FailOnError(err, ErrConnFailed)
	for i := 0; i < size; i++ {

		// Create a channel
		ch, err := conn.Channel()
		pkg.FailOnError(err, ErrChanFailed)
		fmt.Println(MsgChanCreated)

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
	pkg.FailOnError(err, ErrGetChannelFromPool)
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
	pkg.FailOnError(err, ErrQueueFailed)
	fmt.Printf(MsgQueueCreated, queueName)
}

func (r *ChannelPool) ConsumeMessages(queueName string, out func(string2 string)) {
	ch, err := r.Get()
	pkg.FailOnError(err, ErrGetChannelFromPool)
	defer r.Put(ch)

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

func (r *ChannelPool) PublishMessage(ctx context.Context, queueName string, message string) {
	ch, err := r.Get()
	pkg.FailOnError(err, ErrGetChannelFromPool)
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
	pkg.FailOnError(err, ErrProducerFailed)
	fmt.Printf(MsgProducerSucceed, message, queueName)
}

func (r *ChannelPool) PublishMessageWithTTL(ctx context.Context, queueName string, message string, ttlMilliseconds int) {
	ch, err := r.Get()
	pkg.FailOnError(err, ErrGetChannelFromPool)
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
	pkg.FailOnError(err, ErrProducerFailed)
	fmt.Printf(MsgProducerSucceed, message, queueName)
}

func (r *ChannelPool) Close() {
	defer func() {
		for i := range r.channels {
			err := r.channels[i].Close()
			pkg.FailOnError(err, ErrClosingChannel)
		}

		err := r.Connection.Close()
		pkg.FailOnError(err, ErrClosingConnection)
	}()
}
