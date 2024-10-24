package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/kafka"
)

// Benchmark the performance of the synchronous producer
func BenchmarkSyncProducer(producer sarama.SyncProducer, topic string, messageCount int) time.Duration {
	start := time.Now()
	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Sync Message %d", i)
		kafka.SendMessage(producer, topic, message)
	}
	elapsed := time.Since(start)
	return elapsed
}

// Benchmark the performance of the asynchronous producer
func BenchmarkAsyncProducer(producer sarama.AsyncProducer, topic string, messageCount int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup

	// Handle async responses
	go func() {
		for success := range producer.Successes() {
			_ = success // Success, but we don't need to log for performance benchmarking
		}
	}()

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send async message: %v", err.Err)
		}
	}()

	wg.Add(messageCount)
	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Async Message %d", i)
		kafka.SendAsyncMessage(producer, topic, message)
		wg.Done() // Simulate completion for benchmarking
	}

	wg.Wait() // Ensure all messages are sent
	elapsed := time.Since(start)
	return elapsed
}

func main() {
	brokers := []string{"localhost:9092"}
	topic := "benchmark-topic-"
	messageCount := 1000

	// Sync Producer Benchmark
	syncProducer := kafka.CreateProducer(brokers)
	defer syncProducer.Close()

	fmt.Println("Benchmarking Sync Producer...")
	syncDuration := BenchmarkSyncProducer(syncProducer, topic, messageCount)
	fmt.Printf("Sync Producer took: %v\n", syncDuration)

	// Async Producer Benchmark
	asyncProducer:= kafka.CreateAsyncProducer(brokers)
	defer asyncProducer.Close()

	fmt.Println("Benchmarking Async Producer...")
	asyncDuration := BenchmarkAsyncProducer(asyncProducer, topic, messageCount)
	fmt.Printf("Async Producer took: %v\n", asyncDuration)

	// Compare the two
	if syncDuration > asyncDuration {
		fmt.Println("Async producer was faster.")
	} else {
		fmt.Println("Sync producer was faster.")
	}
}
