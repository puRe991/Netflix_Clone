package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

func parseRightsForm(r *http.Request) (store.RightsInfoInput, error) {
	var in store.RightsInfoInput
	var err error

	sourceVal, err := requireEnum("Rechtequelle", r.FormValue("source"),
		string(models.RightsPublicDomain), string(models.RightsCreativeCommons),
		string(models.RightsRevenueShare), string(models.RightsOwned), string(models.RightsLicensed))
	if err != nil {
		return in, err
	}
	in.Source = models.RightsSource(sourceVal)

	in.RightsHolder = optionalString(r.FormValue("rightsHolder"))
	if in.LicenseDocURL, err = optionalURL("licenseDocUrl", r.FormValue("licenseDocUrl")); err != nil {
		return in, err
	}
	if in.RevenueSharePct, err = coerceOptionalFloat("revenueSharePct", r.FormValue("revenueSharePct")); err != nil {
		return in, err
	}
	if in.RevenueSharePct != nil && (*in.RevenueSharePct < 0 || *in.RevenueSharePct > 100) {
		return in, validationErr("revenueSharePct muss zwischen 0 und 100 liegen")
	}
	in.TerritoryLimit = optionalString(r.FormValue("territoryLimit"))
	if in.ValidFrom, err = requireDate("validFrom", r.FormValue("validFrom")); err != nil {
		return in, err
	}
	if in.ValidUntil, err = optionalDate("validUntil", r.FormValue("validUntil")); err != nil {
		return in, err
	}
	in.Notes = optionalString(r.FormValue("notes"))

	return in, nil
}

// AdminRightsUpsert creates or updates the RightsInfo for a media title,
// depending on whether one already exists (1:1 relation via UNIQUE media_id).
func (s *Server) AdminRightsUpsert(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	mediaID := r.PathValue("id")
	in, err := parseRightsForm(r)
	if err != nil {
		return err
	}

	_, err = s.Store.GetRightsInfoByMediaID(r.Context(), mediaID)
	switch err {
	case nil:
		if _, err := s.Store.UpdateRightsInfo(r.Context(), mediaID, in); err != nil {
			return err
		}
	case store.ErrNotFound:
		if _, err := s.Store.CreateRightsInfo(r.Context(), mediaID, in); err != nil {
			return err
		}
	default:
		return err
	}

	redirect(w, r, "/admin/media/"+mediaID+"/edit")
	return nil
}

func (s *Server) AdminRightsDelete(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	mediaID := r.PathValue("id")
	if err := s.Store.DeleteRightsInfo(r.Context(), mediaID); err != nil {
		return err
	}

	redirect(w, r, "/admin/media/"+mediaID+"/edit")
	return nil
}
