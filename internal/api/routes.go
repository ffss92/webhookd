package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) Routes() http.Handler {
	r := chi.NewMux()

	r.Route("/api/v1/subscribers", func(r chi.Router) {
		r.Post("/", s.handleSubscriberCreate())
		r.Get("/{subID}", s.handleSubscriberDetail())
		r.Delete("/{subID}", s.handleSubscriberDelete())
		r.Get("/{subID}/endpoints", s.handleSubscriberEndpointList())
	})

	r.Route("/api/v1/endpoints", func(r chi.Router) {
		r.Post("/", s.handleEndpointCreate())
		r.Get("/{endpointID}", s.handleEndpointDetail())
		r.Delete("/{endpointID}", s.handleEndpointDelete())
	})

	return r
}
