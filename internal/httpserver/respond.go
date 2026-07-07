package httpserver

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func redirect(w http.ResponseWriter, r *http.Request, to string) {
	http.Redirect(w, r, to, http.StatusSeeOther)
}

// handler adapts a (*Server, ResponseWriter, *Request) error func into a
// standard http.HandlerFunc, funneling any returned error through the
// central error mapper instead of letting it become a bare 500.
func (s *Server) handler(h func(*Server, http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(s, w, r); err != nil {
			handleError(w, r, err)
		}
	}
}
