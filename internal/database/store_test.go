package database

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/ffss92/webhookd/internal/postgres"
	"github.com/google/uuid"
)

var testDB *postgres.TestInstance

func TestMain(m *testing.M) {
	testDB = postgres.MustTestInstance()
	defer func() {
		if err := testDB.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	m.Run()
}

func TestInTx_Succeeds(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	store := New(pool)

	var subID uuid.UUID
	err := store.InTx(t.Context(), func(ctx context.Context, store *Store) error {
		sub := &Subscriber{Name: "test"}
		err := store.SaveSubscriber(ctx, sub)
		if err != nil {
			t.Fatal(err)
		}
		subID = sub.ID

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.GetSubscriber(t.Context(), subID)
	if err != nil {
		t.Fatalf("expected valid tx to create subscriber but got %v", err)
	}
}

func TestInTx_Fails(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	store := New(pool)

	var subID uuid.UUID
	txErr := errors.New("something went wrong")
	err := store.InTx(t.Context(), func(ctx context.Context, store *Store) error {
		sub := &Subscriber{Name: "test"}
		err := store.SaveSubscriber(ctx, sub)
		if err != nil {
			t.Fatal(err)
		}
		subID = sub.ID

		return txErr
	})
	if !errors.Is(err, txErr) {
		t.Fatalf("expected error to be txErr but got %v", err)
	}

	_, err = store.GetSubscriber(t.Context(), subID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected failed tx to not save subscriber: %v", err)
	}
}
