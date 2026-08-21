package clickhouse

import (
	"common/events"
	"context"
	"database/sql"
	"encoding/json"
)

type Store struct {
	db *sql.DB
}

func NewStore(dsn string) (*Store, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) InsertBatch(ctx context.Context, batch []events.Envelope) error {
	if len(batch) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO events (event_id, event_type, occurred_at, entity_id, payload) VALUES (?, ?, ?, ?, ?)",
	)

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, e := range batch {
		payloadJSON, err := marshalPayload(e.Payload)
		if err != nil {
			continue
		}

		if _, err := stmt.ExecContext(ctx,
			e.EventID, string(e.EventType), e.OccurredAt, e.EntityID, payloadJSON,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func marshalPayload(payload any) (string, error) {
	data, err := json.Marshal(payload)
	return string(data), err
}
