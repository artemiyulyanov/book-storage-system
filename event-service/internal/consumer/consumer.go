package consumer

import (
	"common/events"
	"context"
	"encoding/json"
	"errors"
	"event-service/internal/clickhouse"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var topics = []string{
	string(events.BookCreated),
	string(events.BookUpdated),
	string(events.BookDeleted),
	string(events.UserRegistered),
	string(events.UserUpdated),
	string(events.UserLoggedIn),
	string(events.UserDeleted),
}

const (
	batchSize     = 100
	flushInterval = 5 * time.Second
)

func Run(ctx context.Context, brokers []string, groupID string, store *clickhouse.Store) {
	for _, topic := range topics {
		go consumeTopic(ctx, brokers, groupID, topic, store)
	}

	<-ctx.Done()
}

func consumeTopic(ctx context.Context, brokers []string, groupID, topic string, store *clickhouse.Store) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	defer reader.Close()

	var batch []events.Envelope
	var rawMessages []kafka.Message

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := store.InsertBatch(ctx, batch); err != nil {
			log.Printf("[%s] failed to insert batch: %v", topic, err)
			return
		}
		if err := reader.CommitMessages(ctx, rawMessages...); err != nil {
			log.Printf("[%s] failed to commit offsets: %v", topic, err)
		}
		batch = batch[:0]
		rawMessages = rawMessages[:0]
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return
		case <-ticker.C:
			flush()
		default:
			fetchCtx, cancel := context.WithTimeout(ctx, flushInterval)
			m, err := reader.FetchMessage(fetchCtx)
			cancel()

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				log.Printf("[%s] fetch error: %v", topic, err)
				continue
			}

			var envelope events.Envelope
			if err := json.Unmarshal(m.Value, &envelope); err != nil {
				log.Printf("[%s] failed to unmarshal message: %v", topic, err)
				continue
			}

			batch = append(batch, envelope)
			rawMessages = append(rawMessages, m)

			if len(batch) >= batchSize {
				flush()
			}
		}
	}
}
