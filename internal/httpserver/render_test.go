package httpserver

import (
	"bytes"
	"testing"
	"time"

	"github.com/pure991/streamflix/internal/models"
)

func TestLoadTemplatesParses(t *testing.T) {
	if _, err := loadTemplates(); err != nil {
		t.Fatalf("loadTemplates() failed: %v", err)
	}
}

// TestRenderRightsTemplates exercises the RightsInfo branches added to the
// admin media form and dashboard templates (no rights, unbounded rights,
// fully populated rights, expired rights) so a nil pointer or bad
// {{if}}/{{and}} guard fails the build instead of showing up at runtime.
func TestRenderRightsTemplates(t *testing.T) {
	tmpls, err := loadTemplates()
	if err != nil {
		t.Fatalf("loadTemplates: %v", err)
	}

	future := time.Now().Add(48 * time.Hour)
	past := time.Now().Add(-48 * time.Hour)
	holder := "A24"
	pct := 12.5

	var buf bytes.Buffer
	render := func(name string, data any) {
		t.Helper()
		buf.Reset()
		if err := tmpls[name].ExecuteTemplate(&buf, "layout", data); err != nil {
			t.Errorf("%s render failed: %v", name, err)
		}
	}

	// New media form: no .Media at all.
	render("admin_media_form.html", adminMediaFormData{PageData: PageData{}})

	// Edit media form, no rights tracked yet.
	media := models.Media{ID: "m1", Title: "Test", Slug: "test", Type: models.MediaMovie}
	render("admin_media_form.html", adminMediaFormData{PageData: PageData{}, Media: &media})

	// Edit media form, unbounded rights (ValidUntil nil).
	media.Rights = &models.RightsInfo{ID: "r1", MediaID: "m1", Source: models.RightsOwned, ValidFrom: time.Now()}
	render("admin_media_form.html", adminMediaFormData{PageData: PageData{}, Media: &media})

	// Edit media form, fully populated + expiring rights.
	media.Rights = &models.RightsInfo{
		ID: "r1", MediaID: "m1", Source: models.RightsLicensed,
		RightsHolder: &holder, RevenueSharePct: &pct, ValidFrom: past, ValidUntil: &future,
	}
	render("admin_media_form.html", adminMediaFormData{PageData: PageData{}, Media: &media})

	// Dashboard with no expiring rights.
	render("admin_dashboard.html", adminStatsData{PageData: PageData{}})

	// Dashboard with an upcoming and an already-expired grant.
	render("admin_dashboard.html", adminStatsData{
		PageData: PageData{},
		ExpiringRights: []models.RightsInfo{
			{Source: models.RightsLicensed, ValidUntil: &future, RightsHolder: &holder, Media: &models.Media{ID: "m1", Title: "Test"}},
			{Source: models.RightsRevenueShare, ValidUntil: &past, Media: &models.Media{ID: "m2", Title: "Expired One"}},
		},
	})
}
