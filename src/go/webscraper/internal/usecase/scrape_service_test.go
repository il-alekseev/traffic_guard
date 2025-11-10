package usecase

import (
	"net/url"
	"strings"
	"testing"
)

func TestShouldFallbackToRoot(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		raw  string
		want bool
	}{
		{name: "deep path", raw: "https://example.com/news/local/story", want: true},
		{name: "long path", raw: "https://example.com/" + strings.Repeat("a", 60), want: true},
		{name: "long query", raw: "https://example.com/post?id=" + strings.Repeat("b", 60), want: true},
		{name: "root path", raw: "https://example.com/", want: false},
		{name: "single segment", raw: "https://example.com/about", want: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.raw)
			if err != nil {
				t.Fatalf("url.Parse(%q) error = %v", tt.raw, err)
			}
			if got := shouldFallbackToRoot(u); got != tt.want {
				t.Fatalf("shouldFallbackToRoot(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestBuildRootURL(t *testing.T) {
	t.Parallel()

	raw := "https://example.com/news/2024/11/09/story"
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q) error = %v", raw, err)
	}

	got := buildRootURL(u)
	want := "https://example.com/"
	if got != want {
		t.Fatalf("buildRootURL(%q) = %q, want %q", raw, got, want)
	}
}
