// Package httpserver wires routes, templates and handlers together.
package httpserver

import (
	"html/template"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/billing"
	"github.com/pure991/streamflix/internal/store"
)

type Server struct {
	Store     *store.Store
	Sessions  *auth.Sessions
	Billing   *billing.Billing
	AppEnv    string
	AppURL    string
	templates map[string]*template.Template
}

func New(st *store.Store, sessions *auth.Sessions, bill *billing.Billing, appEnv, appURL string) (*Server, error) {
	tmpls, err := loadTemplates()
	if err != nil {
		return nil, err
	}
	return &Server{Store: st, Sessions: sessions, Billing: bill, AppEnv: appEnv, AppURL: appURL, templates: tmpls}, nil
}
