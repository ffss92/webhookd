package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Message struct {
	ID           uuid.UUID
	Type         string
	Data         json.RawMessage
	Tags         []string
	SubscriberID uuid.UUID
	CreatedAt    time.Time
}

func (s Store) SaveMessage(ctx context.Context, msg *Message) error {
	if msg.Tags == nil {
		msg.Tags = make([]string, 0)
	}

	query := `
	INSERT INTO messages (type, data, tags, subscriber_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`
	args := []any{msg.Type, msg.Data, msg.Tags, msg.SubscriberID}
	err := s.pool.QueryRow(ctx, query, args...).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

func (s Store) GetMessage(ctx context.Context, msgID uuid.UUID) (*Message, error) {
	query := `
	SELECT id, type, data, tags, subscriber_id, created_at
	FROM messages
	WHERE id = $1`
	var msg Message
	err := s.pool.QueryRow(ctx, query, msgID).Scan(
		&msg.ID,
		&msg.Type,
		&msg.Data,
		&msg.Tags,
		&msg.SubscriberID,
		&msg.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get message: %w", err)
		}
	}
	return &msg, nil
}

func (s Store) DeleteMessage(ctx context.Context, msgID uuid.UUID) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, msgID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}
