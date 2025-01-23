package producer

import (
	"fmt"
	"testing"
	"time"
)

// Order represents an e-commerce order
type Order struct {
	OrderID     string    `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	Items       []string  `json:"items"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// Payment represents payment transaction
type Payment struct {
	PaymentID   string    `json:"payment_id"`
	OrderID     string    `json:"order_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}

// Order processing example
func TestExampleOrderProcessing(t *testing.T) {
	producer, _ := NewTransactionalProducer([]string{"localhost:9092"}, "order-processor")
	defer producer.Close()

	order := Order{
		OrderID: "123",
		Items:   []string{"item1", "item2"},
	}

	payment := Payment{
		OrderID:     "123",
		Amount:      100.00,
		ProcessedAt: time.Now(),
	}

	batch := MessageBatch{
		Messages: []Message{
			{
				Topic: "orders",
				Key:   order.OrderID,
				Value: order,
				Headers: []Header{
					{Key: "type", Value: []byte("order")},
					{Key: "version", Value: []byte("1.0")},
				},
			},
			{
				Topic: "payments",
				Key:   order.OrderID,
				Value: payment,
				Headers: []Header{
					{Key: "type", Value: []byte("payment")},
				},
			},
		},
		Options: &BatchOptions{
			Timeout: 5 * time.Second,
			Retries: 3,
		},
	}

	if err := producer.SendMessageBatch(batch); err != nil {
		fmt.Printf("Failed to process order: %v\n", err)
	}
}

// Event logging example
func TestExampleEventLogging(t *testing.T) {
	producer, _ := NewTransactionalProducer([]string{"localhost:9092"}, "event-logger")
	defer producer.Close()

	event := struct {
		EventType string      `json:"event_type"`
		Timestamp time.Time   `json:"timestamp"`
		Data      interface{} `json:"data"`
	}{
		EventType: "user_login",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"user_id": "user123",
			"ip":      "192.168.1.1",
		},
	}

	msg := Message{
		Topic: "events",
		Key:   "user123",
		Value: event,
		Headers: []Header{
			{Key: "event_type", Value: []byte("user_login")},
		},
	}

	if err := producer.SendMessage(msg); err != nil {
		fmt.Printf("Failed to log event: %v\n", err)
	}
}

// Benchmark the producer
func BenchmarkTransactionalProducer_Throughput(b *testing.B) {
	producer, err := NewTransactionalProducer(brokers, "benchmark-producer")
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(messageSize)

	start := time.Now()
	for i := 0; i < b.N; i++ {

		key := fmt.Sprintf("key-%d", i)
		if err := producer.SendMessage(Message{
			Topic: testTopic,
			Key:   key,
			Value: message,
		}); err != nil {
			b.Fatalf("Failed to send batch: %v", err)
		}
	}

	duration := time.Since(start)

	messagesPerSecond := float64(numMessages) / duration.Seconds()
	b.ReportMetric(messagesPerSecond, "msgs/sec")
	b.ReportMetric(float64(messageSize*numMessages)/duration.Seconds(), "bytes/sec")
}
