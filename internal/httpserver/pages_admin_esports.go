package httpserver

import (
	"net/http"
	"strings"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

type adminTeamsData struct {
	PageData
	Teams []models.Team
}

func (s *Server) AdminTeamsPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	teams, err := s.Store.ListTeams(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_teams.html", adminTeamsData{PageData: s.basePageData(w, r), Teams: teams})
}

func (s *Server) AdminTeamCreate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	var in store.TeamInput
	var err error
	if in.Name, err = requireMinMaxLen("Name", r.FormValue("name"), 2, 60); err != nil {
		return err
	}
	if in.Slug, err = requireSlug(r.FormValue("slug")); err != nil {
		return err
	}
	if in.ShortName, err = optionalMinMaxLen("Kürzel", r.FormValue("shortName"), 1, 12); err != nil {
		return err
	}
	if in.Country, err = optionalMinMaxLen("Land", r.FormValue("country"), 2, 40); err != nil {
		return err
	}
	if in.LogoURL, err = optionalURL("logoUrl", r.FormValue("logoUrl")); err != nil {
		return err
	}

	if _, err := s.Store.CreateTeam(r.Context(), in); err != nil {
		return err
	}
	redirect(w, r, "/admin/teams")
	return nil
}

func (s *Server) AdminTeamDelete(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	used, err := s.Store.TeamHasMatches(r.Context(), id)
	if err != nil {
		return err
	}
	if used {
		return validationErr("Team ist noch Matches zugeordnet und kann nicht gelöscht werden")
	}
	if err := s.Store.DeleteTeam(r.Context(), id); err != nil {
		return err
	}
	redirect(w, r, "/admin/teams")
	return nil
}

// parseMatchForm reads the optional match section of the episode form.
// It returns nil (no error) when no teams were picked, meaning "this
// episode is not a match".
func parseMatchForm(r *http.Request) (*store.MatchInput, error) {
	teamA := strings.TrimSpace(r.FormValue("teamAId"))
	teamB := strings.TrimSpace(r.FormValue("teamBId"))
	if teamA == "" && teamB == "" {
		return nil, nil
	}
	if teamA == "" || teamB == "" {
		return nil, validationErr("Für ein Match müssen beide Teams gewählt werden")
	}
	if teamA == teamB {
		return nil, validationErr("Ein Team kann nicht gegen sich selbst spielen")
	}

	in := store.MatchInput{TeamAID: teamA, TeamBID: teamB, GameID: optionalString(r.FormValue("gameId"))}
	var err error
	if in.Stage, err = requireMinMaxLen("Phase", r.FormValue("stage"), 2, 60); err != nil {
		return nil, err
	}
	bestOf, err := requireEnum("Best-of", r.FormValue("bestOf"), "0", "1", "3", "5")
	if err != nil {
		return nil, err
	}
	in.BestOf = int(bestOf[0] - '0')

	if in.ScoreA, err = coerceOptionalInt(r.FormValue("scoreA")); err != nil {
		return nil, err
	}
	if in.ScoreB, err = coerceOptionalInt(r.FormValue("scoreB")); err != nil {
		return nil, err
	}
	if (in.ScoreA == nil) != (in.ScoreB == nil) {
		return nil, validationErr("Ergebnis bitte für beide Teams angeben oder leer lassen")
	}
	// BestOf 0 means "format unknown", so only the lower bound applies then.
	if in.ScoreA != nil && (*in.ScoreA < 0 || *in.ScoreB < 0 || (in.BestOf > 0 && *in.ScoreA+*in.ScoreB > in.BestOf)) {
		return nil, validationErr("Ergebnis passt nicht zum Best-of-Format")
	}
	if in.PlayedOn, err = optionalDate("Spieldatum", r.FormValue("playedOn")); err != nil {
		return nil, err
	}

	for _, line := range strings.Split(r.FormValue("extraVideoUrls"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		u, err := requireURL("Weitere Video-URLs", line)
		if err != nil {
			return nil, err
		}
		in.ExtraVideoURLs = append(in.ExtraVideoURLs, u)
	}
	return &in, nil
}
