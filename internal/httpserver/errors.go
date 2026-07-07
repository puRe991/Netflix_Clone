package httpserver

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/store"
)

// ValidationError carries a user-facing message for a bad request body.
type ValidationError struct{ Message string }

func (e ValidationError) Error() string { return e.Message }

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	isAPI := strings.HasPrefix(r.URL.Path, "/api/")

	status := http.StatusInternalServerError
	message := "Internal Server Error"

	var ve ValidationError
	switch {
	case errors.Is(err, auth.ErrUnauthenticated):
		status, message = http.StatusUnauthorized, "Authentication required"
	case errors.Is(err, auth.ErrForbidden):
		status, message = http.StatusForbidden, "Forbidden"
	case errors.Is(err, auth.ErrCSRF):
		status, message = http.StatusForbidden, "Invalid or missing CSRF token"
	case errors.Is(err, store.ErrNotFound):
		status, message = http.StatusNotFound, "Not found"
	case errors.As(err, &ve):
		status, message = http.StatusUnprocessableEntity, ve.Message
	default:
		log.Printf("unhandled error on %s %s: %v", r.Method, r.URL.Path, err)
	}

	if isAPI {
		JSON(w, status, map[string]string{"error": message})
		return
	}

	if status == http.StatusUnauthorized || (status == http.StatusForbidden && errors.Is(err, auth.ErrForbidden)) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte("<!doctype html><html><body style=\"background:#050505;color:#f8fafc;font-family:sans-serif;padding:2rem\"><h1>" + message + "</h1><p><a href=\"/\" style=\"color:#dc2626\">Zur Startseite</a></p></body></html>"))
}
