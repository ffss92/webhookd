package database

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestSubscriberLifecycle(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	store := New(pool)

	sub := &Subscriber{
		Name: "test",
	}
	err := store.SaveSubscriber(t.Context(), sub)
	if err != nil {
		t.Fatal(err)
	}

	read, err := store.GetSubscriber(t.Context(), sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(sub, read); diff != "" {
		t.Fatalf("mismatch (-expected, +got):\n%s", diff)
	}

	sub.Name = "test-updated"
	err = store.UpdateSubscriber(t.Context(), sub)
	if err != nil {
		t.Fatal(err)
	}

	read, err = store.GetSubscriber(t.Context(), sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(sub, read); diff != "" {
		t.Fatalf("mismatch (-expected, +got):\n%s", diff)
	}

	err = store.DeleteSubscriber(t.Context(), sub.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.GetSubscriber(t.Context(), sub.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("expected deleted subscriber to not be found")
	}
}

func TestGetSubscriber_NotFound(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	store := New(pool)

	_, err := store.GetSubscriber(t.Context(), uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound but got %v", err)
	}
}
