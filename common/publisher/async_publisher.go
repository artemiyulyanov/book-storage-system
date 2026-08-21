package publisher

import (
	"common/events"
	"context"
	"log"
	"time"
)

type KafkaAsyncPublisher struct {
	producer *events.Producer
}

func NewKafkaAsyncPublisher(producer *events.Producer) *KafkaAsyncPublisher {
	return &KafkaAsyncPublisher{producer: producer}
}

func (publisher *KafkaAsyncPublisher) PublishAsync(eventType events.EventType, entityID int64, payload any) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := publisher.producer.Publish(ctx, eventType, entityID, payload); err != nil {
			log.Printf("failed to publish %s event (key=%s): %v", eventType, entityID, err)
		}
	}()
}
