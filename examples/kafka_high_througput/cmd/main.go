package main

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/pkg/kafka/consumer"
	"github.com/ciazhar/go-start-small/pkg/kafka/producer"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

const (
	//brokerAddress = "localhost:807,localhost:8098,localhost:8099"
	brokerAddress = "localhost:9092,localhost:9093,localhost:9094"
	topic         = "example-topic"
	batchSize     = 1000
	flushTimeout  = 2 * time.Second
	groupID       = "example-group"
	messageCount  = 1_000_000 // Total messages to send
)

func main() {

	// Start time
	startTime := time.Now()

	// Signal handler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	// Start producer
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		produceMessages()
	}()

	// Start consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeMessages(ctx)
	}()

	elapsed := time.Since(startTime)
	log.Printf("Finished sending %d messages in %v", messageCount, elapsed)

	wg.Wait()
}

func produceMessages() {
	bp, err := producer.NewBatchProducer(strings.Split(brokerAddress, ","), batchSize, flushTimeout)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	go func() {
		for range bp.Producer.Errors() {
		}
	}()

	go func() {
		for range bp.Producer.Successes() {
		}
	}()

	for i := 0; i < messageCount; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("Message %d", i+1)),
		}
		bp.Messages <- msg
	}

	close(bp.Messages)
	bp.Wg.Wait()
}

func consumeMessages(ctx context.Context) {
	cg, err := consumer.NewConsumerGroup(strings.Split(brokerAddress, ","), groupID, []string{topic})
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	go func() {
		if err := cg.Consume(ctx); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()

	<-ctx.Done()
}

func handleSignals(cancel context.CancelFunc) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
	log.Println("Interrupt signal received. Shutting down...")
	cancel()
}
