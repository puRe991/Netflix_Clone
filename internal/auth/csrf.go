package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
)

const csrfCookieName = "streamflix_csrf"

func randomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CSRFToken returns the token to embed in forms/headers on this request,
// creating and setting the cookie on first use.
func (s *Sessions) CSRFToken(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(csrfCookieName); err == nil && c.Value != "" {
		return c.Value
	}

	token := randomToken()
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
	return token
}

var ErrCSRF = errCSRF{}

type errCSRF struct{}

func (errCSRF) Error() string { return "invalid csrf token" }

// VerifyCSRF checks the submitted token (form field or X-CSRF-Token header)
// against the cookie using a constant-time comparison.
func (s *Sessions) VerifyCSRF(r *http.Request, submitted string) error {
	c, err := r.Cookie(csrfCookieName)
	if err != nil || c.Value == "" || submitted == "" {
		return ErrCSRF
	}
	if subtle.ConstantTimeCompare([]byte(c.Value), []byte(submitted)) != 1 {
		return ErrCSRF
	}
	return nil
}
