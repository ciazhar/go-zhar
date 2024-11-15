package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	// "github.com/rcrowley/go-metrics"
)

type BenchmarkConfig struct {
	Brokers       []string
	Topic         string
	MessageSize   int
	MessageCount  int
	ProducerCount int
	ConsumerCount int
	ConsumerGroup string
	RunDuration   time.Duration
	IsRedPanda    bool
}

type BenchmarkResults struct {
	MessagesProduced   uint64
	MessagesConsumed   uint64
	Duration           time.Duration
	ProducerThroughput float64
	ConsumerThroughput float64
	AverageLatency     time.Duration
}

func main() {
	// Parse command line flags
	brokers := flag.String("brokers", "localhost:9092", "Kafka/Redpanda brokers list")
	topic := flag.String("topic", "benchmark", "Topic to use for benchmark")
	messageSize := flag.Int("message-size", 1024, "Size of each message in bytes")
	messageCount := flag.Int("message-count", 1000000, "Number of messages to produce")
	producerCount := flag.Int("producers", 1, "Number of concurrent producers")
	consumerCount := flag.Int("consumers", 1, "Number of concurrent consumers")
	consumerGroup := flag.String("consumer-group", "benchmark-group", "Consumer group ID")
	duration := flag.Duration("duration", 1*time.Minute, "Duration of the benchmark")
	isRedPanda := flag.Bool("redpanda", false, "Use Redpanda-specific configurations")
	flag.Parse()

	config := BenchmarkConfig{
		Brokers:       []string{*brokers},
		Topic:         *topic,
		MessageSize:   *messageSize,
		MessageCount:  *messageCount,
		ProducerCount: *producerCount,
		ConsumerCount: *consumerCount,
		ConsumerGroup: *consumerGroup,
		RunDuration:   *duration,
		IsRedPanda:    *isRedPanda,
	}

	results := runBenchmark(config)
	printResults(results)
}

func runBenchmark(config BenchmarkConfig) BenchmarkResults {
	var (
		wg               sync.WaitGroup
		messagesProduced uint64
		messagesConsumed uint64
		startTime        = time.Now()
	)

	// Setup signals for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), config.RunDuration)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Create Kafka/Redpanda configuration
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5

	if config.IsRedPanda {
		// Redpanda-specific optimizations
		kafkaConfig.Net.MaxOpenRequests = 1
		kafkaConfig.Producer.Compression = sarama.CompressionSnappy
	}

	// Start producers
	for i := 0; i < config.ProducerCount; i++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()
			producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
			if err != nil {
				log.Printf("Failed to create producer %d: %v", producerID, err)
				return
			}
			defer producer.Close()

			message := make([]byte, config.MessageSize)
			for i := range message {
				message[i] = byte(i % 256)
			}

			for {
				select {
				case <-ctx.Done():
					return
				default:
					msg := &sarama.ProducerMessage{
						Topic:     config.Topic,
						Value:     sarama.ByteEncoder(message),
						Timestamp: time.Now(),
					}

					if _, _, err := producer.SendMessage(msg); err != nil {
						log.Printf("Failed to send message: %v", err)
						continue
					}

					atomic.AddUint64(&messagesProduced, 1)
				}
			}
		}(i)
	}

	// Start consumers
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	for i := 0; i < config.ConsumerCount; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()

			group, err := sarama.NewConsumerGroup(config.Brokers, config.ConsumerGroup, consumerConfig)
			if err != nil {
				log.Printf("Failed to create consumer %d: %v", consumerID, err)
				return
			}
			defer group.Close()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					err := group.Consume(ctx, []string{config.Topic}, &ConsumerGroupHandler{
						ready:            make(chan bool),
						messagesConsumed: &messagesConsumed,
					})
					if err != nil {
						log.Printf("Error from consumer: %v", err)
					}
				}
			}
		}(i)
	}

	// Wait for completion or timeout
	select {
	case <-signals:
		cancel()
	case <-ctx.Done():
	}

	wg.Wait()
	duration := time.Since(startTime)

	return BenchmarkResults{
		MessagesProduced:   messagesProduced,
		MessagesConsumed:   messagesConsumed,
		Duration:           duration,
		ProducerThroughput: float64(messagesProduced) / duration.Seconds(),
		ConsumerThroughput: float64(messagesConsumed) / duration.Seconds(),
	}
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	ready            chan bool
	messagesConsumed *uint64
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		atomic.AddUint64(h.messagesConsumed, 1)
		session.MarkMessage(message, "")
	}
	return nil
}

func printResults(results BenchmarkResults) {
	fmt.Printf("\nBenchmark Results:\n")
	fmt.Printf("================\n")
	fmt.Printf("Duration: %v\n", results.Duration)
	fmt.Printf("Messages Produced: %d\n", results.MessagesProduced)
	fmt.Printf("Messages Consumed: %d\n", results.MessagesConsumed)
	fmt.Printf("Producer Throughput: %.2f msgs/sec\n", results.ProducerThroughput)
	fmt.Printf("Consumer Throughput: %.2f msgs/sec\n", results.ConsumerThroughput)
}
