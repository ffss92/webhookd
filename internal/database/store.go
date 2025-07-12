package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound = errors.New("not found")
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Store struct {
	pool DBTX
}

func New(pool DBTX) *Store {
	return &Store{
		pool: pool,
	}
}

func (s *Store) InTx(ctx context.Context, fn func(ctx context.Context, store *Store) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}

	store := New(tx)
	if err := fn(ctx, store); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return fmt.Errorf("failed to rollback tx (%v): %w", txErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
