package handlers

import (
	"common/events"
	"context"
	"log"
	"time"
)

type KafkaEventHandler struct {
	kafkaProducer *events.Producer
}

func NewKafkaEventHandler(producer *events.Producer) KafkaEventHandler {
	return KafkaEventHandler{kafkaProducer: producer}
}

func (h *KafkaEventHandler) PublishAsync(eventType events.EventType, entityID int64, payload any) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.kafkaProducer.Publish(ctx, eventType, entityID, payload); err != nil {
			log.Printf("failed to publish %s event (key=%s): %v", eventType, entityID, err)
		}
	}()
}
