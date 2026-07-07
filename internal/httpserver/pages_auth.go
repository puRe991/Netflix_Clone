package httpserver

import "net/http"

type authFormData struct {
	PageData
}

func (s *Server) LoginPage(w http.ResponseWriter, r *http.Request) error {
	return s.render(w, r, "login.html", authFormData{PageData: s.basePageData(w, r)})
}

func (s *Server) RegisterPage(w http.ResponseWriter, r *http.Request) error {
	return s.render(w, r, "register.html", authFormData{PageData: s.basePageData(w, r)})
}
