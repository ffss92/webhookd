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

type Subscriber struct {
	ID        uuid.UUID
	Name      string
	Metadata  json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s Store) SaveSubscriber(ctx context.Context, sub *Subscriber) error {
	if sub.Metadata == nil {
		sub.Metadata = json.RawMessage(`null`)
	}

	query := `
	INSERT INTO subscribers (name, metadata)
	VALUES ($1, $2)
	RETURNING id, created_at, updated_at`
	args := []any{sub.Name, sub.Metadata}

	err := s.pool.QueryRow(ctx, query, args...).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save subscriber: %w", err)
	}
	return nil
}

func (s Store) GetSubscriber(ctx context.Context, subID uuid.UUID) (*Subscriber, error) {
	query := `
	SELECT id, name, metadata, created_at, updated_at
	FROM subscribers
	WHERE id = $1`

	var subscriber Subscriber
	err := s.pool.QueryRow(ctx, query, subID).Scan(
		&subscriber.ID,
		&subscriber.Name,
		&subscriber.Metadata,
		&subscriber.CreatedAt,
		&subscriber.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &subscriber, nil
}

func (s Store) UpdateSubscriber(ctx context.Context, subscriber *Subscriber) error {
	query := `
	UPDATE subscribers SET
		name = $2,
		metadata = $3,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1
	RETURNING updated_at`
	args := []any{subscriber.ID, subscriber.Name, subscriber.Metadata}

	err := s.pool.QueryRow(ctx, query, args...).Scan(&subscriber.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) DeleteSubscriber(ctx context.Context, subID uuid.UUID) error {
	query := `DELETE FROM subscribers WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, subID)
	if err != nil {
		return fmt.Errorf("failed to delete subscriber: %w", err)
	}
	return nil
}
