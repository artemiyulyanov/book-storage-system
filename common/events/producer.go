package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	mu      sync.Mutex
	writers map[EventType]*kafka.Writer
	brokers []string
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writers: make(map[EventType]*kafka.Writer),
		brokers: brokers,
	}
}

func (p *Producer) writerFor(eventType EventType) *kafka.Writer {
	p.mu.Lock()
	defer p.mu.Unlock()

	if w, ok := p.writers[eventType]; ok {
		return w
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(p.brokers...),
		Topic:    string(eventType),
		Balancer: &kafka.LeastBytes{},
	}

	p.writers[eventType] = w
	return w
}

func (p *Producer) Publish(ctx context.Context, eventType EventType, entityID int64, payload any) error {
	envelope := Envelope{
		EventID:    uuid.NewString(),
		EventType:  eventType,
		OccurredAt: time.Now().UTC(),
		EntityID:   entityID,
		Payload:    payload,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	return p.writerFor(eventType).WriteMessages(ctx, kafka.Message{
		Key:   []byte(envelope.EventID),
		Value: data,
	})
}

func (p *Producer) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, w := range p.writers {
		w.Close()
	}
}
