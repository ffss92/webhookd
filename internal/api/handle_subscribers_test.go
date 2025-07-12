package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ffss92/webhookd/internal/database"
	"github.com/google/uuid"
)

func TestHandleSubscriberCreate(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	api := &Server{
		pool:  pool,
		store: database.New(pool),
	}

	srv := httptest.NewServer(api.Routes())
	defer srv.Close()

	testCases := []struct {
		name   string
		req    *CreateSubscriberRequest
		status int
	}{
		{
			name: "valid request",
			req: &CreateSubscriberRequest{
				Name: "Test",
			},
			status: http.StatusCreated,
		},
		{
			name: "invalid request",
			req: &CreateSubscriberRequest{
				Name: "",
			},
			status: http.StatusUnprocessableEntity,
		},
		{
			name:   "invalid request (missing body)",
			req:    nil,
			status: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}

			client := srv.Client()

			req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/subscribers", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != tt.status {
				t.Errorf("expected status to be %d but got %d", tt.status, res.StatusCode)
			}
			if res.StatusCode == http.StatusCreated {
				var subscriber Subscriber
				err := json.NewDecoder(res.Body).Decode(&subscriber)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestHandleSubscriberDetail(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	api := &Server{
		pool:  pool,
		store: database.New(pool),
	}

	srv := httptest.NewServer(api.Routes())
	defer srv.Close()

	sub := &database.Subscriber{
		Name: "Test",
	}
	err := api.store.SaveSubscriber(t.Context(), sub)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		subID  string
		status int
	}{
		{
			name:   "valid id",
			subID:  sub.ID.String(),
			status: http.StatusOK,
		},
		{
			name:   "non existing id",
			subID:  uuid.NewString(),
			status: http.StatusNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			client := srv.Client()

			path := fmt.Sprintf("/api/v1/subscribers/%s", tt.subID)
			req, err := http.NewRequest(http.MethodGet, srv.URL+path, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != tt.status {
				t.Fatalf("expected status %d but got %d", tt.status, res.StatusCode)
			}
			if res.StatusCode == http.StatusOK {
				var sub Subscriber
				err := json.NewDecoder(res.Body).Decode(&sub)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestHandleSubscriberDelete(t *testing.T) {
	t.Parallel()

	pool := testDB.NewPool(t)
	api := &Server{
		pool:  pool,
		store: database.New(pool),
	}

	srv := httptest.NewServer(api.Routes())
	defer srv.Close()

	sub := &database.Subscriber{
		Name: "Test",
	}
	err := api.store.SaveSubscriber(t.Context(), sub)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		subID  string
		status int
	}{
		{
			name:   "valid id",
			subID:  sub.ID.String(),
			status: http.StatusNoContent,
		},
		{
			name:   "non existing id",
			subID:  uuid.NewString(),
			status: http.StatusNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			client := srv.Client()

			path := fmt.Sprintf("/api/v1/subscribers/%s", tt.subID)
			req, err := http.NewRequest(http.MethodDelete, srv.URL+path, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != tt.status {
				t.Fatalf("expected status %d but got %d", tt.status, res.StatusCode)
			}
		})
	}
}
