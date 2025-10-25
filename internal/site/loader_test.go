package site

import "testing"

func TestEnsurePermalink(t *testing.T) {
	tests := map[string]string{
		"":                 "/",
		"/":                "/",
		"posts/foo":        "/posts/foo/",
		"/posts/foo":       "/posts/foo/",
		"/posts/foo/":      "/posts/foo/",
		"/docs/setup.html": "/docs/setup.html",
		"docs/setup.html":  "/docs/setup.html",
	}

	for input, want := range tests {
		if got := ensurePermalink(input); got != want {
			t.Errorf("ensurePermalink(%q) = %q; want %q", input, got, want)
		}
	}
}
