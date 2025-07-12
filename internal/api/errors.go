package api

import (
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"detail,omitempty"`
}

func (s *Server) serverError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Error(
		"failed to write json",
		slog.String("method", r.Method),
		slog.String("url", r.URL.RequestURI()),
		slog.String("err", err.Error()),
	)
	s.writeJSON(w, r, http.StatusInternalServerError, ErrorResponse{
		Message: "Something went wrong",
	})
}

func (s *Server) badRequest(w http.ResponseWriter, r *http.Request, _ error) {
	s.writeJSON(w, r, http.StatusBadRequest, ErrorResponse{
		Message: "Invalid or malformed request",
	})
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, r, http.StatusNotFound, ErrorResponse{
		Message: "Resource not found",
	})
}

func (s *Server) validationError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	s.writeJSON(w, r, http.StatusUnprocessableEntity, ErrorResponse{
		Message: "Validation failed",
		Errors:  errors,
	})
}
