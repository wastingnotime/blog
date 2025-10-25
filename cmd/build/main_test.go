package main

import "testing"

func TestPostOutputPath(t *testing.T) {
	cases := map[string]string{
		"/":                  "public/index.html",
		"":                   "public/index.html",
		"/posts/foo/":        "public/posts/foo/index.html",
		"posts/bar":          "public/posts/bar/index.html",
		"/about.html":        "public/about.html",
		"about.html":         "public/about.html",
		"/docs/setup.html":   "public/docs/setup.html",
		"docs/setup.html":    "public/docs/setup.html",
		"/docs/setup/index/": "public/docs/setup/index/index.html",
	}

	for input, want := range cases {
		got := postOutputPath(input)
		if got != want {
			t.Errorf("postOutputPath(%q) = %q; want %q", input, got, want)
		}
	}
}
