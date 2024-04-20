package service

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/examples/kafka/custom-consumer-group/internal/event/model"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"log"
	"sync"
	"time"
)

const (
	ChunkSize  = 300
	WindowTime = 10
)

type EventService struct {
	logger logger.Logger

	ready  chan bool
	stopCh chan struct{}

	windowEnd time.Time
	mux       sync.Mutex

	producer *kafka.SyncProducer
	buffers  map[string][]string
}

func (e *EventService) Setup(session sarama.ConsumerGroupSession) error {
	close(e.ready)
	e.windowEnd = time.Now().Add(WindowTime * time.Second)
	return nil
}

func (e *EventService) Cleanup(session sarama.ConsumerGroupSession) error {
	// Flush any remaining data in the buffers
	e.mux.Lock()
	defer e.mux.Unlock()

	for key, buffer := range e.buffers {
		if len(buffer) > 0 {
			err := e.flushBuffer(key, buffer)
			if err != nil {
				e.logger.Infof("Failed to flush buffer for key %s: %v\n", key, err)
			}
		}
	}
	return nil
}

func (e *EventService) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		key := string(message.Key)

		e.mux.Lock()

		// Initialize the buffer for the key if it doesn't exist
		if _, ok := e.buffers[key]; !ok {
			e.buffers[key] = make([]string, 0)
		}

		// Append the single row to the buffer for the corresponding key
		e.buffers[key] = append(e.buffers[key], string(message.Value))

		// Check if the buffer size reaches the chunk size
		if len(e.buffers[key]) == ChunkSize {
			err := e.flushBuffer(key, e.buffers[key])
			if err != nil {
				log.Printf("Failed to flush buffer for key %s: %v\n", key, err)
				continue
			}

			// Commit offsets for the processed messages
			session.MarkMessage(message, "")

			// Reset the buffer for the key
			e.buffers[key] = make([]string, 0)

			// Update the window end time
			e.windowEnd = time.Now().Add(WindowTime * time.Second)

			session.Commit()
		} else {

			awake := true

			go func() {
				for awake {
					e.mux.Lock()
					for key, buffer := range e.buffers {
						if len(buffer) > 0 && time.Since(e.windowEnd) > WindowTime {
							err := e.flushBuffer(key, buffer)
							if err != nil {
								log.Printf("Failed to flush buffer for key %s: %v\n", key, err)
								continue
							}
							// Commit offsets for the processed messages
							session.MarkMessage(message, "")
							// Reset the buffer for the key
							e.buffers[key] = make([]string, 0)
							// Update the window end time
							e.windowEnd = time.Now().Add(WindowTime * time.Second)

							awake = false

							session.Commit()
						}
					}
					e.mux.Unlock()
				}
			}()
		}

		e.mux.Unlock()

		select {
		case <-e.stopCh:
			return nil
		default:
		}
	}

	return nil
}

type GroupedEvents struct {
	Key  string             `json:"key"`
	Data []model.EmailEvent `json:"data"`
}

func (w *EventService) flushBuffer(key string, buffer []string) error {

	events := make([]model.EmailEvent, 0)
	for i := range buffer {

		var event model.EmailEvent
		err := json.Unmarshal([]byte(buffer[i]), &event)
		if err != nil {
			return fmt.Errorf("failed to unmarshal buffer: %v", err)
		}

		events = append(events, event)
	}
	webhookEvent := GroupedEvents{
		Key:  key,
		Data: events,
	}

	jsonBytes, err := json.Marshal(webhookEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal buffer: %v", err)
	}

	w.producer.PublishMessageWithKey("grouped-events", key, string(jsonBytes))

	return nil
}

func NewEventService(logger logger.Logger, producer *kafka.SyncProducer) *EventService {
	return &EventService{
		logger:    logger,
		ready:     make(chan bool),
		stopCh:    make(chan struct{}),
		windowEnd: time.Time{},
		mux:       sync.Mutex{},
		producer:  producer,
		buffers:   make(map[string][]string),
	}

}
