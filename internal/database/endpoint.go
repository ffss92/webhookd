package database

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Endpoint struct {
	ID           uuid.UUID
	Label        string
	URL          string
	Secret       string
	Disabled     bool
	FilterTypes  []string
	SubscriberID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func removeDuplicates[T comparable](values []T) []T {
	if values == nil {
		return make([]T, 0)
	}
	unique := make(map[T]struct{})
	for _, value := range values {
		unique[value] = struct{}{}
	}
	return slices.Collect(maps.Keys(unique))
}

func (s Store) SaveEndpoint(ctx context.Context, endpoint *Endpoint) error {
	endpoint.FilterTypes = removeDuplicates(endpoint.FilterTypes)

	query := `
	INSERT INTO endpoints (label, url, secret, filter_types, disabled, subscriber_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`
	args := []any{
		endpoint.Label,
		endpoint.URL,
		endpoint.Secret,
		endpoint.FilterTypes,
		endpoint.Disabled,
		endpoint.SubscriberID,
	}

	err := s.pool.QueryRow(ctx, query, args...).Scan(
		&endpoint.ID,
		&endpoint.CreatedAt,
		&endpoint.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save endpoint: %w", err)
	}
	return nil
}

func (s Store) GetEndpoint(ctx context.Context, endpointID uuid.UUID) (*Endpoint, error) {
	query := `
	SELECT
		id, label, url, secret, disabled, 
		filter_types, subscriber_id, created_at, updated_at
	FROM endpoints
	WHERE id = $1`

	endpoint, err := scanEndpoint(s.pool.QueryRow(ctx, query, endpointID))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return endpoint, nil
}

type ListEndpointsParams struct {
	SubscriberID uuid.UUID
	Disabled     *bool
	FilterType   *string
}

func (s Store) ListEndpoints(ctx context.Context, params ListEndpointsParams) ([]*Endpoint, error) {
	query := `
	SELECT
		id, label, url, secret, disabled, 
		filter_types, subscriber_id, created_at, updated_at
	FROM endpoints
	WHERE subscriber_id = $1
	AND (disabled = $2 OR $2 IS NULL)
	AND (
		$3::TEXT IS NULL
		OR filter_types = '{}'
		OR filter_types @> ARRAY[$3]
	)`
	args := []any{params.SubscriberID, params.Disabled, params.FilterType}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	endpoints := make([]*Endpoint, 0)
	for rows.Next() {
		endpoint, err := scanEndpoint(rows)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return endpoints, nil
}

func scanEndpoint(row pgx.Row) (*Endpoint, error) {
	var endpoint Endpoint
	err := row.Scan(
		&endpoint.ID, &endpoint.Label, &endpoint.URL, &endpoint.Secret, &endpoint.Disabled,
		&endpoint.FilterTypes, &endpoint.SubscriberID, &endpoint.CreatedAt, &endpoint.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (s Store) UpdateEndpoint(ctx context.Context, endpoint *Endpoint) error {
	endpoint.FilterTypes = removeDuplicates(endpoint.FilterTypes)

	query := `
	UPDATE endpoints SET
		label = $2,
		url = $3,
		disabled = $4,
		filter_types = $5,
		secret = $6,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1
	RETURNING updated_at`
	args := []any{
		endpoint.ID,
		endpoint.Label,
		endpoint.URL,
		endpoint.Disabled,
		endpoint.FilterTypes,
		endpoint.Secret,
	}
	err := s.pool.QueryRow(ctx, query, args...).Scan(&endpoint.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update endpoint: %w", err)
	}
	return nil
}

func (s Store) DeleteEndpoint(ctx context.Context, endpointID uuid.UUID) error {
	query := `DELETE FROM endpoints WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, endpointID)
	if err != nil {
		return fmt.Errorf("failed to delete endpoint: %w", err)
	}
	return nil
}
