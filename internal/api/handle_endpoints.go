package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ffss92/webhookd/internal/database"
	"github.com/ffss92/webhookd/internal/validator"
	"github.com/ffss92/webhookd/internal/webhook"
	"github.com/google/uuid"
)

type Endpoint struct {
	ID           uuid.UUID `json:"id"`
	Label        string    `json:"label"`
	URL          string    `json:"url"`
	Disabled     bool      `json:"disabled"`
	FilterTypes  []string  `json:"filter_types"`
	SubscriberID uuid.UUID `json:"subscriber_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func mapEndpoint(record *database.Endpoint) *Endpoint {
	return &Endpoint{
		ID:           record.ID,
		Label:        record.Label,
		URL:          record.URL,
		Disabled:     record.Disabled,
		FilterTypes:  record.FilterTypes,
		SubscriberID: record.SubscriberID,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}
}

type CreateEndpointRequest struct {
	Label        string    `json:"label"`
	URL          string    `json:"url"`
	FilterTypes  []string  `json:"filter_types"`
	SubscriberID uuid.UUID `json:"subscriber_id"`

	validator.Validator `json:"-"`
}

func (s *Server) handleEndpointCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input CreateEndpointRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			s.badRequest(w, r, err)
			return
		}

		input.Check(validator.NotBlank(input.Label), "label", "Must be provided")
		input.Check(validator.MaxLength(input.Label, 255), "label", "Must have at most 255 characters")
		input.Check(validator.NotBlank(input.URL), "url", "Must be provided")
		input.Check(validator.HTTPUrl(input.URL), "url", "Must be a valid http or https url")
		input.Check(input.SubscriberID != uuid.Nil, "subscriber_id", "Must not be an empty uuid")
		if !input.IsValid() {
			s.validationError(w, r, input.FieldErrors)
			return
		}

		sub, err := s.store.GetSubscriber(r.Context(), input.SubscriberID)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				input.SetFieldError("subscriber_id", "Invalid subscriber_id value")
				s.validationError(w, r, input.FieldErrors)
			default:
				s.serverError(w, r, err)
			}
			return
		}

		endpoint := &database.Endpoint{
			Label:        input.Label,
			URL:          input.URL,
			FilterTypes:  input.FilterTypes,
			Secret:       webhook.NewSecret(),
			SubscriberID: sub.ID,
		}
		err = s.store.SaveEndpoint(r.Context(), endpoint)
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		res := mapEndpoint(endpoint)
		s.writeJSON(w, r, http.StatusCreated, res)
	}
}
func (s *Server) handleEndpointDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointID, err := uuidParam(r, "endpointID")
		if err != nil {
			s.notFound(w, r)
			return
		}

		endpoint, err := s.store.GetEndpoint(r.Context(), endpointID)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				s.notFound(w, r)
			default:
				s.serverError(w, r, err)
			}
			return
		}

		s.writeJSON(w, r, http.StatusOK, endpoint)
	}
}

func (s *Server) handleEndpointDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointID, err := uuidParam(r, "endpointID")
		if err != nil {
			s.notFound(w, r)
			return
		}

		endpoint, err := s.store.GetEndpoint(r.Context(), endpointID)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				s.notFound(w, r)
			default:
				s.serverError(w, r, err)
			}
			return
		}

		err = s.store.DeleteEndpoint(r.Context(), endpoint.ID)
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
