package rabbitmq

import (
	"fmt"
	error2 "github.com/ciazhar/zhar/pkg/error"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type ChannelPool struct {
	mu       sync.Mutex
	channels []*amqp.Channel
}

func NewChannelPool(username, password, host, port string, size int) *ChannelPool {
	pool := &ChannelPool{}
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	if err != nil {
		log.Fatalf("Failed to create connection: %s", err)
	}
	for i := 0; i < size; i++ {

		// Create a channel
		ch, err := conn.Channel()
		error2.FailOnError(err, "Failed to open a channel")
		fmt.Println("Channel created")

		pool.channels = append(pool.channels, ch)
	}
	return pool
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.channels) == 0 {
		return nil, fmt.Errorf("connection pool is empty")
	}

	conn := p.channels[0]
	p.channels = p.channels[1:]
	return conn, nil
}

func (p *ChannelPool) Put(conn *amqp.Channel) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.channels = append(p.channels, conn)
}
