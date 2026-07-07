package auth

import (
	"errors"
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
)

// OptionalUser returns the session user if present, or nil for anonymous
// requests. Used by pages that render differently when logged in but don't
// require it (browse, movie/series detail, the nav bar).
func (s *Sessions) OptionalUser(r *http.Request) *SessionUser {
	return s.Read(r)
}

func (s *Sessions) RequireUser(r *http.Request) (SessionUser, error) {
	u := s.Read(r)
	if u == nil {
		return SessionUser{}, ErrUnauthenticated
	}
	return *u, nil
}

func (s *Sessions) RequireAdmin(r *http.Request) (SessionUser, error) {
	u, err := s.RequireUser(r)
	if err != nil {
		return SessionUser{}, err
	}
	if u.Role != models.RoleAdmin {
		return SessionUser{}, ErrForbidden
	}
	return u, nil
}
