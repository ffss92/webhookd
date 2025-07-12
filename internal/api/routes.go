package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) Routes() http.Handler {
	r := chi.NewMux()

	r.Route("/api/v1/subscribers", func(r chi.Router) {
		r.Post("/", s.handleSubscriberCreate())

		r.Route("/{subID}", func(r chi.Router) {
			r.Use(s.withSubscriber)
			r.Get("/", s.handleSubscriberDetail())
			r.Delete("/", s.handleSubscriberDelete())
			r.Get("/endpoints", s.handleSubscriberEndpointList())
		})
	})

	r.Route("/api/v1/endpoints", func(r chi.Router) {
		r.Post("/", s.handleEndpointCreate())
		r.Get("/{endpointID}", s.handleEndpointDetail())
		r.Delete("/{endpointID}", s.handleEndpointDelete())
	})

	return r
}
