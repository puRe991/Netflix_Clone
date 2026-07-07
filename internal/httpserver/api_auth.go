package httpserver

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/models"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	email, err := requireEmail(r.FormValue("email"))
	if err != nil {
		return err
	}
	password, err := requireMinLen("Passwort", r.FormValue("password"), 8)
	if err != nil {
		return err
	}
	name, err := optionalMinMaxLen("Name", r.FormValue("name"), 2, 80)
	if err != nil {
		return err
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	profileName := "Profil"
	if name != nil {
		profileName = *name
	}

	userID, err := s.Store.CreateUserWithProfile(r.Context(), email, passwordHash, name, profileName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return validationErr("Diese E-Mail-Adresse ist bereits registriert")
		}
		return err
	}

	if err := s.Sessions.Create(w, auth.SessionUser{ID: userID, Email: email, Role: models.RoleUser, Name: name}); err != nil {
		return err
	}

	redirect(w, r, "/profiles")
	return nil
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	email, err := requireEmail(r.FormValue("email"))
	if err != nil {
		return err
	}
	password := r.FormValue("password")
	if password == "" {
		return validationErr("Passwort erforderlich")
	}

	user, err := s.Store.GetUserByEmail(r.Context(), email)
	if err != nil || !auth.VerifyPassword(password, user.PasswordHash) {
		JSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return nil
	}

	if err := s.Sessions.Create(w, auth.SessionUser{ID: user.ID, Email: user.Email, Role: user.Role, Name: user.Name}); err != nil {
		return err
	}

	redirect(w, r, "/browse")
	return nil
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) error {
	s.Sessions.Clear(w)
	redirect(w, r, "/")
	return nil
}
