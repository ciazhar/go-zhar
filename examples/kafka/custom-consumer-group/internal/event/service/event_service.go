package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/examples/kafka/custom-consumer-group/internal/event/model"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"sync"
	"time"
)

const (
	ChunkSize  = 300
	WindowTime = 5
)

type WebhookEventConsumer struct {
	ctx context.Context

	mux sync.Mutex

	db       *pgxpool.Pool
	producer *kafka.SyncProducer
	buffers  map[string][]string
	logger   *logger.Logger
}

func NewWebhookEventConsumer(
	ctx context.Context,
	producer *kafka.SyncProducer,
	logger *logger.Logger,
) *WebhookEventConsumer {
	return &WebhookEventConsumer{
		ctx:      ctx,
		producer: producer,
		buffers:  make(map[string][]string),
		logger:   logger,
	}
}

func (w *WebhookEventConsumer) Setup(_ sarama.ConsumerGroupSession) (err error) {
	return
}

func (w *WebhookEventConsumer) Cleanup(_ sarama.ConsumerGroupSession) (err error) {
	// Flush any remaining data in the buffers
	w.mux.Lock()
	defer w.mux.Unlock()

	for key, buffer := range w.buffers {
		if len(buffer) > 0 {
			err := w.flushBuffer(key, buffer)
			if err != nil {
				w.logger.Infof("Failed to flush buffer for key %s: %v\n", key, err)
			}
		}
	}
	return
}

func (w *WebhookEventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	ticker := time.NewTicker(WindowTime * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return
			}

			key := string(message.Key)

			w.mux.Lock()
			buffer := w.buffers[key]
			buffer = append(buffer, string(message.Value))
			w.buffers[key] = buffer
			w.mux.Unlock()

			if len(buffer) == ChunkSize {
				w.flushAndCommit(key, buffer, message, session)
			} else {
				select {
				case <-ticker.C:
					w.mux.Lock()
					if len(w.buffers[key]) > 0 {
						w.flushAndCommit(key, w.buffers[key], message, session)
					}
					w.mux.Unlock()
				default:
				}
			}

		case <-w.ctx.Done():
			return
		}
	}
}

func (w *WebhookEventConsumer) flushAndCommit(key string, buffer []string, message *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {

	w.logger.Infof("Flushing buffer for key %s\n", key)

	err := w.flushBuffer(key, buffer)
	if err != nil {
		w.logger.Infof("Failed to flush buffer for key %s: %v\n", key, err)
		return
	}

	if message != nil {
		session.MarkMessage(message, "")
	}
	w.buffers[key] = nil
	session.Commit()
}

func (w *WebhookEventConsumer) flushBuffer(key string, buffer []string) (err error) {
	events := make([]model.EmailEvent, 0, len(buffer))
	for _, data := range buffer {
		var event model.EmailEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			w.logger.Infof("Failed to unmarshal buffer for key %s: %v\n", key, err)
			return
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return
	}

	if err := w.sendWebhook(events, key); err != nil {
		w.logger.Infof("Failed to send webhook for key %s: %v\n", key, err)
		return
	}

	return
}

func (w *WebhookEventConsumer) sendWebhook(events []model.EmailEvent, key string) (err error) {

	webhookEvent := GroupedEvents{
		Key:  key,
		Data: events,
	}

	jsonBytes, err := json.Marshal(webhookEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal buffer: %v", err)
	}

	w.producer.PublishMessageWithKey("grouped-events", key, string(jsonBytes))

	return
}

type GroupedEvents struct {
	Key  string             `json:"key"`
	Data []model.EmailEvent `json:"data"`
}
