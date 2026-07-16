package models

import (
	"testing"
	"time"
)

func TestRightsInfoExpired(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(24 * time.Hour)

	cases := []struct {
		name string
		r    RightsInfo
		want bool
	}{
		{"unbefristet (nil ValidUntil)", RightsInfo{ValidUntil: nil}, false},
		{"laeuft erst in der Zukunft ab", RightsInfo{ValidUntil: &future}, false},
		{"bereits abgelaufen", RightsInfo{ValidUntil: &past}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.r.Expired(); got != c.want {
				t.Errorf("Expired() = %v, want %v", got, c.want)
			}
		})
	}
}
