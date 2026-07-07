package httpserver

import "net/http"

type homeData struct {
	PageData
}

func (s *Server) HomePage(w http.ResponseWriter, r *http.Request) error {
	return s.render(w, r, "home.html", homeData{PageData: s.basePageData(w, r)})
}

type pricingData struct {
	PageData
}

func (s *Server) PricingPage(w http.ResponseWriter, r *http.Request) error {
	return s.render(w, r, "pricing.html", pricingData{PageData: s.basePageData(w, r)})
}

type legalData struct {
	PageData
	Title string
}

func (s *Server) ImpressumPage(w http.ResponseWriter, r *http.Request) error {
	return s.render(w, r, "legal.html", legalData{PageData: s.basePageData(w, r), Title: "Impressum"})
}

func (s *Server) DatenschutzPage(w http.ResponseWriter, r *http.Request) error {
	return s.render(w, r, "legal.html", legalData{PageData: s.basePageData(w, r), Title: "Datenschutz"})
}
