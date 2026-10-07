package httpserver

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSafeLocalPath(t *testing.T) {
	cases := map[string]string{
		"/series/x":          "/series/x",
		"//evil.example":     "/browse",
		"/\\evil.example":    "/browse",
		"https://evil.test/": "/browse",
		"":                   "/browse",
	}
	for in, want := range cases {
		if got := safeLocalPath(in, "/browse"); got != want {
			t.Errorf("safeLocalPath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseMatchForm(t *testing.T) {
	parse := func(v url.Values) (bool, error) {
		r := httptest.NewRequest("POST", "/", strings.NewReader(v.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		m, err := parseMatchForm(r)
		return m != nil, err
	}
	base := func() url.Values {
		return url.Values{"teamAId": {"a"}, "teamBId": {"b"}, "stage": {"Grand Final"}, "bestOf": {"3"}}
	}

	if isMatch, err := parse(url.Values{}); err != nil || isMatch {
		t.Errorf("no teams must mean 'no match', got match=%v err=%v", isMatch, err)
	}
	if isMatch, err := parse(base()); err != nil || !isMatch {
		t.Errorf("valid match rejected: %v", err)
	}

	bad := map[string]func(url.Values){
		"one team only":       func(v url.Values) { v.Del("teamBId") },
		"same team":           func(v url.Values) { v.Set("teamBId", "a") },
		"bo2":                 func(v url.Values) { v.Set("bestOf", "2") },
		"half a score":        func(v url.Values) { v.Set("scoreA", "2") },
		"score exceeds bo3":   func(v url.Values) { v.Set("scoreA", "3"); v.Set("scoreB", "1") },
		"bad extra video url": func(v url.Values) { v.Set("extraVideoUrls", "not a url") },
		"bad date":            func(v url.Values) { v.Set("playedOn", "23.07.2017") },
	}
	for name, mutate := range bad {
		v := base()
		mutate(v)
		if _, err := parse(v); err == nil {
			t.Errorf("%s: expected validation error", name)
		}
	}

	v := base()
	v.Set("bestOf", "5")
	v.Set("scoreA", "3")
	v.Set("scoreB", "1")
	v.Set("extraVideoUrls", "https://youtu.be/sIQ1Eh11Quk\r\n\r\nhttps://youtu.be/w2MdCEZSBtw\n")
	if _, err := parse(v); err != nil {
		t.Errorf("valid Bo5 with extra parts rejected: %v", err)
	}
}
