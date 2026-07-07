package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pure991/streamflix/internal/models"
)

const CookieName = "streamflix_session"
const sessionTTL = 7 * 24 * time.Hour

// SessionUser is the identity carried in the signed session cookie.
// Subscription status is deliberately NOT part of it: it must always be
// re-read from the database at the point it matters, so an expiring/renewed
// Stripe subscription takes effect immediately instead of waiting out a
// stale 7-day-old claim.
type SessionUser struct {
	ID    string      `json:"id"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
	Name  *string     `json:"name,omitempty"`
}

type sessionClaims struct {
	SessionUser
	jwt.RegisteredClaims
}

type Sessions struct {
	secret []byte
	secure bool
}

func NewSessions(secret string, secureCookies bool) *Sessions {
	return &Sessions{secret: []byte(secret), secure: secureCookies}
}

func (s *Sessions) Create(w http.ResponseWriter, u SessionUser) error {
	claims := sessionClaims{
		SessionUser: u,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sessionTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
	return nil
}

func (s *Sessions) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// Read returns the session user from the request cookie, or nil if there is
// no valid session (missing, malformed, expired, or bad signature).
func (s *Sessions) Read(r *http.Request) *SessionUser {
	cookie, err := r.Cookie(CookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}

	claims := &sessionClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		return nil
	}

	return &claims.SessionUser
}
