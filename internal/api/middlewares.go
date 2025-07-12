package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/ffss92/webhookd/internal/database"
)

const (
	subscriberKey contextKey = iota
)

type contextKey int

func getSubscriber(ctx context.Context) *database.Subscriber {
	sub, ok := ctx.Value(subscriberKey).(*database.Subscriber)
	if !ok {
		panic("subscriber not present in context")
	}
	return sub
}

func (s *Server) withSubscriber(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), subscriberKey, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
