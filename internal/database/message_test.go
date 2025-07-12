package database

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestGetMessage_NotFound(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	store := New(pool)

	_, err := store.GetMessage(t.Context(), uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound but got %v", err)
	}
}

func TestMessageLifecycle(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	store := New(pool)

	sub := &Subscriber{
		Name: "test",
	}
	err := store.SaveSubscriber(t.Context(), sub)
	if err != nil {
		t.Fatalf("failed to save subscriber: %v", err)
	}

	testCases := []struct {
		name string
		msg  *Message
	}{
		{
			name: "with tags",
			msg: &Message{
				Type:         "test",
				Data:         json.RawMessage(`{"key": "value"}`),
				Tags:         []string{"tag1", "tag2"},
				SubscriberID: sub.ID,
			},
		},
		{
			name: "without tags",
			msg: &Message{
				Type:         "test-1",
				Data:         json.RawMessage(`{"key": "value"}`),
				SubscriberID: sub.ID,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := store.SaveMessage(t.Context(), tc.msg)
			if err != nil {
				t.Fatalf("failed to save message: %v", err)
			}

			got, err := store.GetMessage(t.Context(), tc.msg.ID)
			if err != nil {
				t.Fatalf("failed to get message: %v", err)
			}

			if diff := cmp.Diff(tc.msg, got); diff != "" {
				t.Errorf("message mismatch (-want +got):\n%s", diff)
			}

			err = store.DeleteMessage(t.Context(), tc.msg.ID)
			if err != nil {
				t.Fatalf("failed to delete message: %v", err)
			}

			_, err = store.GetMessage(t.Context(), tc.msg.ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("expected ErrNotFound after deletion, got %v", err)
			}
		})
	}
}
