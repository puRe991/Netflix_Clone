package httpserver

import "testing"

func TestYouTubeVideoID(t *testing.T) {
	cases := map[string]string{
		"https://www.youtube.com/watch?v=7nWxZin959o":               "7nWxZin959o",
		"https://youtube.com/watch?v=-5laY7bjono&t=120":             "-5laY7bjono",
		"https://m.youtube.com/watch?v=do4FUYUHmUM":                 "do4FUYUHmUM",
		"https://youtu.be/ttbe_A8Ee50":                              "ttbe_A8Ee50",
		"https://www.youtube.com/embed/L9eBTSwDwaM?start=5":         "L9eBTSwDwaM",
		"https://www.youtube-nocookie.com/embed/09N-bQbLrx8":        "09N-bQbLrx8",
		"https://www.youtube.com/live/meijWUaxNYg":                  "meijWUaxNYg",
		"https://commondatastorage.googleapis.com/sample/video.mp4": "",
		"https://www.youtube.com/watch?v=tooShort":                  "",
		"https://evil.example/watch?v=7nWxZin959o":                  "",
		"javascript:alert(1)":                                       "",
		"":                                                          "",
	}
	for in, want := range cases {
		if got := youTubeVideoID(in); got != want {
			t.Errorf("youTubeVideoID(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAllYouTube(t *testing.T) {
	if allYouTube(nil) {
		t.Error("empty list must not count as YouTube")
	}
	if !allYouTube([]string{"https://youtu.be/ttbe_A8Ee50", "https://youtu.be/sIQ1Eh11Quk"}) {
		t.Error("all-YouTube list not detected")
	}
	if allYouTube([]string{"https://youtu.be/ttbe_A8Ee50", "https://example.com/a.mp4"}) {
		t.Error("mixed list must not count as YouTube")
	}
}
