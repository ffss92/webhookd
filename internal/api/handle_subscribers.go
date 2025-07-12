package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ffss92/webhookd/internal/database"
	"github.com/ffss92/webhookd/internal/validator"
	"github.com/google/uuid"
)

type Subscriber struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func mapSubscriber(record *database.Subscriber) (*Subscriber, error) {
	var metadata map[string]any
	err := json.Unmarshal(record.Metadata, &metadata)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		ID:        record.ID,
		Name:      record.Name,
		Metadata:  metadata,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.CreatedAt,
	}, nil
}

type CreateSubscriberRequest struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata"`

	validator.Validator `json:"-"`
}

func (s *Server) handleSubscriberCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input CreateSubscriberRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			s.badRequest(w, r, err)
			return
		}

		input.Name = strings.TrimSpace(input.Name)
		if input.Metadata == nil {
			input.Metadata = make(map[string]any)
		}

		input.Check(validator.NotBlank(input.Name), "name", "Must be provided")
		input.Check(validator.MaxLength(input.Name, 255), "name", "Must have at most 255 characters")
		if !input.IsValid() {
			s.validationError(w, r, input.FieldErrors)
			return
		}

		metadata, err := json.Marshal(input.Metadata)
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		sub := &database.Subscriber{
			Name:     input.Name,
			Metadata: metadata,
		}

		err = s.store.SaveSubscriber(r.Context(), sub)
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		res, err := mapSubscriber(sub)
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		s.writeJSON(w, r, http.StatusCreated, res)
	}
}

func (s *Server) handleSubscriberDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subID, err := uuidParam(r, "subID")
		if err != nil {
			s.notFound(w, r)
			return
		}

		sub, err := s.store.GetSubscriber(r.Context(), subID)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				s.notFound(w, r)
			default:
				s.serverError(w, r, err)
			}
			return
		}

		res, err := mapSubscriber(sub)
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		s.writeJSON(w, r, http.StatusOK, res)
	}
}

func (s *Server) handleSubscriberDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subID, err := uuidParam(r, "subID")
		if err != nil {
			s.notFound(w, r)
			return
		}

		sub, err := s.store.GetSubscriber(r.Context(), subID)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				s.notFound(w, r)
			default:
				s.serverError(w, r, err)
			}
			return
		}

		err = s.store.DeleteSubscriber(r.Context(), sub.ID)
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleSubscriberEndpointList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subID, err := uuidParam(r, "subID")
		if err != nil {
			s.notFound(w, r)
			return
		}

		sub, err := s.store.GetSubscriber(r.Context(), subID)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				s.notFound(w, r)
			default:
				s.serverError(w, r, err)
			}
			return
		}

		endpoints, err := s.store.ListEndpoints(r.Context(), database.ListEndpointsParams{
			SubscriberID: sub.ID,
		})
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		res := make([]*Endpoint, 0, len(endpoints))
		for _, endpoint := range endpoints {
			res = append(res, mapEndpoint(endpoint))
		}
		s.writeJSON(w, r, http.StatusOK, res)
	}
}
