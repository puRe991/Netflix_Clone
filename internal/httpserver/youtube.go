package httpserver

import (
	"net/url"
	"regexp"
	"strings"
)

var youTubeIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// youTubeVideoID extracts the video ID from the common YouTube URL forms
// (watch?v=, youtu.be/, /embed/, /live/, /shorts/, youtube-nocookie.com).
// It returns "" for anything that isn't a YouTube video URL, so callers can
// use it both as a parser and as an "is this a YouTube embed" check.
func youTubeVideoID(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return ""
	}

	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	host = strings.TrimPrefix(host, "m.")

	var id string
	switch host {
	case "youtu.be":
		id = strings.Trim(u.Path, "/")
	case "youtube.com", "youtube-nocookie.com":
		if u.Path == "/watch" {
			id = u.Query().Get("v")
			break
		}
		for _, prefix := range []string{"/embed/", "/live/", "/shorts/"} {
			if strings.HasPrefix(u.Path, prefix) {
				id = strings.Trim(strings.TrimPrefix(u.Path, prefix), "/")
				break
			}
		}
	}

	if !youTubeIDPattern.MatchString(id) {
		return ""
	}
	return id
}

// allYouTube reports whether every URL is a YouTube video. Mixed lists
// (some self-hosted files, some embeds) aren't supported by the player.
func allYouTube(urls []string) bool {
	if len(urls) == 0 {
		return false
	}
	for _, u := range urls {
		if youTubeVideoID(u) == "" {
			return false
		}
	}
	return true
}
