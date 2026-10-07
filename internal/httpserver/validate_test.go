package httpserver

import "testing"

func TestRequireDate(t *testing.T) {
	if _, err := requireDate("validFrom", "2026-07-16"); err != nil {
		t.Errorf("expected valid date to pass, got %v", err)
	}
	if _, err := requireDate("validFrom", "16.07.2026"); err == nil {
		t.Error("expected invalid date format to fail")
	}
	if _, err := requireDate("validFrom", ""); err == nil {
		t.Error("expected empty date to fail")
	}
}

func TestOptionalDate(t *testing.T) {
	got, err := optionalDate("validUntil", "")
	if err != nil || got != nil {
		t.Errorf("expected nil, nil for empty input, got %v, %v", got, err)
	}

	got, err = optionalDate("validUntil", "2026-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Format("2006-01-02") != "2026-12-31" {
		t.Errorf("unexpected parsed date: %v", got)
	}
}

func TestCoerceOptionalFloat(t *testing.T) {
	got, err := coerceOptionalFloat("revenueSharePct", "")
	if err != nil || got != nil {
		t.Errorf("expected nil, nil for empty input, got %v, %v", got, err)
	}

	got, err = coerceOptionalFloat("revenueSharePct", "12.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || *got != 12.5 {
		t.Errorf("unexpected value: %v", got)
	}

	if _, err := coerceOptionalFloat("revenueSharePct", "not-a-number"); err == nil {
		t.Error("expected error for non-numeric input")
	}
}

func TestOptionalString(t *testing.T) {
	if got := optionalString("   "); got != nil {
		t.Errorf("expected nil for blank input, got %v", *got)
	}
	if got := optionalString(" A24 "); got == nil || *got != "A24" {
		t.Errorf("expected trimmed value 'A24', got %v", got)
	}
}

func TestRequireImageURL(t *testing.T) {
	for _, ok := range []string{"https://example.com/a.jpg", "/static/img/esports/x-poster.svg"} {
		if _, err := requireImageURL("img", ok); err != nil {
			t.Errorf("requireImageURL(%q) rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "img.jpg", "/etc/passwd", "/static/../secret", "javascript:alert(1)"} {
		if _, err := requireImageURL("img", bad); err == nil {
			t.Errorf("requireImageURL(%q) accepted", bad)
		}
	}
}
