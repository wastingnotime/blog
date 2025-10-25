package site

import "testing"

func TestAbsoluteURL(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		in   string
		want string
	}{
		{
			name: "no base url returns relative",
			cfg:  Config{BasePath: ""},
			in:   "/posts/first/",
			want: "/posts/first/",
		},
		{
			name: "base url applied to root",
			cfg:  Config{BasePath: "", BaseURL: "https://example.com"},
			in:   "/",
			want: "https://example.com/",
		},
		{
			name: "base url applied to nested path",
			cfg:  Config{BasePath: "", BaseURL: "https://example.com"},
			in:   "/posts/first/",
			want: "https://example.com/posts/first/",
		},
		{
			name: "base path and base url combined",
			cfg:  Config{BasePath: "/blog", BaseURL: "https://example.com"},
			in:   "/posts/first/",
			want: "https://example.com/blog/posts/first/",
		},
		{
			name: "base path only remains relative",
			cfg:  Config{BasePath: "/blog"},
			in:   "/",
			want: "/blog/",
		},
		{
			name: "html permalinks preserved",
			cfg:  Config{BasePath: "", BaseURL: "https://example.com"},
			in:   "/docs/setup.html",
			want: "https://example.com/docs/setup.html",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cfg.AbsoluteURL(tc.in); got != tc.want {
				t.Fatalf("AbsoluteURL(%q) = %q; want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestNormalizeBaseURL(t *testing.T) {
	tests := map[string]string{
		"":                         "",
		"https://example.com":      "https://example.com",
		"https://example.com/":     "https://example.com",
		" https://example.com/ ":   "https://example.com",
		"https://example.com/foo/": "https://example.com/foo",
	}

	for input, want := range tests {
		if got := normalizeBaseURL(input); got != want {
			t.Errorf("normalizeBaseURL(%q) = %q; want %q", input, got, want)
		}
	}
}
