package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type ConnectionPool struct {
	mu          sync.Mutex
	connections []*amqp.Connection
}

func NewConnectionPool(username, password, host, port string, size int) *ConnectionPool {
	pool := &ConnectionPool{}
	for i := 0; i < size; i++ {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
		if err != nil {
			log.Fatalf("Failed to create connection: %s", err)
		}
		pool.connections = append(pool.connections, conn)
	}
	return pool
}

func (p *ConnectionPool) Get() (*amqp.Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.connections) == 0 {
		return nil, fmt.Errorf("connection pool is empty")
	}

	conn := p.connections[0]
	p.connections = p.connections[1:]
	return conn, nil
}

func (p *ConnectionPool) Put(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.connections = append(p.connections, conn)
}
