package database

import (
	"crypto/rand"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestEndpointLifecycle(t *testing.T) {
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

	endpoint := &Endpoint{
		Label:        "test-endpoint",
		URL:          "http://endpoint.com",
		Secret:       rand.Text(),
		SubscriberID: sub.ID,
		FilterTypes:  []string{"foo.bar", "foo.bar"},
	}
	err = store.SaveEndpoint(t.Context(), endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoint.FilterTypes) == 2 {
		t.Fatal("expected duplicate filter type to be removed on save")
	}

	read, err := store.GetEndpoint(t.Context(), endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(endpoint, read); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	endpoint.Disabled = true
	endpoint.FilterTypes = []string{"test.created", "test.created"}
	err = store.UpdateEndpoint(t.Context(), endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoint.FilterTypes) == 2 {
		t.Fatal("expected duplicate filter type to be removed on update")
	}

	read, err = store.GetEndpoint(t.Context(), endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(endpoint, read); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	err = store.DeleteEndpoint(t.Context(), endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.GetEndpoint(t.Context(), endpoint.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted endpoint to not be found but got %v", err)
	}
}

func TestGetEndpoint_NotFound(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	store := New(pool)

	_, err := store.GetEndpoint(t.Context(), uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound but got %v", err)
	}
}

func ptr[T any](value T) *T {
	return &value
}

func TestListEndpoints(t *testing.T) {
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

	create := []*Endpoint{
		{
			Label:        "test-1",
			URL:          "http://test-1.com",
			SubscriberID: sub.ID,
			FilterTypes:  []string{"test.created", "test.updated"},
		},
		{
			Label:        "test-3",
			URL:          "http://test-3.com",
			SubscriberID: sub.ID,
		},
		{
			Label:        "test-2",
			URL:          "http://test-2.com",
			SubscriberID: sub.ID,
			Disabled:     true,
			FilterTypes:  []string{"test.deleted"},
		},
	}

	for _, endpoint := range create {
		err = store.SaveEndpoint(t.Context(), endpoint)
		if err != nil {
			t.Fatal(err)
		}
	}

	testCases := []struct {
		name     string
		expected []*Endpoint
		params   ListEndpointsParams
	}{
		{
			name:     "no filters",
			expected: create,
			params: ListEndpointsParams{
				SubscriberID: sub.ID,
			},
		},
		{
			name: "enabled",
			expected: []*Endpoint{
				create[0],
				create[1],
			},
			params: ListEndpointsParams{
				SubscriberID: sub.ID,
				Disabled:     ptr(false),
			},
		},
		{
			name: "disabled",
			expected: []*Endpoint{
				create[2],
			},
			params: ListEndpointsParams{
				SubscriberID: sub.ID,
				Disabled:     ptr(true),
			},
		},
		{
			name: "filter type (test.created)",
			expected: []*Endpoint{
				create[0],
				create[1],
			},
			params: ListEndpointsParams{
				SubscriberID: sub.ID,
				FilterType:   ptr("test.created"),
			},
		},
		{
			name: "filter type (test.deleted)",
			expected: []*Endpoint{
				create[1],
				create[2],
			},
			params: ListEndpointsParams{
				SubscriberID: sub.ID,
				FilterType:   ptr("test.deleted"),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := store.ListEndpoints(t.Context(), tt.params)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
