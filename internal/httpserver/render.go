package httpserver

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/web"
)

var templateFuncs = template.FuncMap{
	"formatDuration": func(minutes *int) string {
		if minutes == nil || *minutes <= 0 {
			return "–"
		}
		return fmt.Sprintf("%d Min.", *minutes)
	},
	"formatSeconds": func(seconds int64) string {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	},
	"hasGenre": func(genres []models.Genre, id string) bool {
		for _, g := range genres {
			if g.ID == id {
				return true
			}
		}
		return false
	},
	"joinTags": func(tags []string) string {
		return strings.Join(tags, ", ")
	},
	"add": func(a, b int) int { return a + b },
	"dict": func(pairs ...any) (map[string]any, error) {
		if len(pairs)%2 != 0 {
			return nil, fmt.Errorf("dict requires an even number of arguments")
		}
		m := make(map[string]any, len(pairs)/2)
		for i := 0; i < len(pairs); i += 2 {
			key, ok := pairs[i].(string)
			if !ok {
				return nil, fmt.Errorf("dict keys must be strings")
			}
			m[key] = pairs[i+1]
		}
		return m, nil
	},
}

func loadTemplates() (map[string]*template.Template, error) {
	entries, err := web.FS.ReadDir("templates/pages")
	if err != nil {
		return nil, err
	}

	var pages []string
	for _, e := range entries {
		pages = append(pages, "templates/pages/"+e.Name())
	}

	result := make(map[string]*template.Template, len(pages))
	for _, p := range pages {
		name := path.Base(p)
		t := template.New(name).Funcs(templateFuncs)
		t, err := t.ParseFS(web.FS, "templates/layout.html", "templates/partials/*.html", p)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", name, err)
		}
		result[name] = t
	}
	return result, nil
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, page string, data any) error {
	t, ok := s.templates[page]
	if !ok {
		return fmt.Errorf("template %s not found", page)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteTemplate(w, "layout", data)
}
