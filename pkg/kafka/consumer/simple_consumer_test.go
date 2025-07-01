package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

// mockProcessor simulates message processing.
func mockProcessor(msg *sarama.ConsumerMessage) error {
	time.Sleep(5 * time.Millisecond) // Simulated processing time
	return nil
}

// TestConsumerBenchmark simulates a Kafka consumer with a finite message stream.
func TestConsumerBenchmark(t *testing.T) {
	broker := sarama.NewMockBroker(t, 1)
	defer broker.Close()

	// Create a mock topic with 10 messages
	mockResponses := make([]sarama.MockResponse, 0)
	for i := 0; i < 10; i++ {
		msg := &sarama.ConsumerMessage{
			Topic: "test_topic",
			Value: []byte("mock message"),
		}
		mockResponses = append(mockResponses, sarama.NewMockConsumerMessage(msg))
	}
	broker.SetHandlerByMap(map[string]sarama.MockResponse{
		"Fetch": sarama.NewMockFetchResponse(t, 1).
			SetMessage("test_topic", 0, 0, mockResponses...),
	})

	// Create consumer
	consumer := NewSimpleConsumer([]string{broker.Addr()}, "test_topic", mockProcessor, "roundrobin", true, 5)
	assert.NotNil(t, consumer)

	// Measure execution time
	start := time.Now()
	err := consumer.ConsumeWithRetry("test_topic", 0, sarama.OffsetOldest)
	duration := time.Since(start)

	assert.NoError(t, err)
	t.Logf("Processed 10 messages in %v", duration)
}
