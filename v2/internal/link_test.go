package internal

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/kellegous/golinks"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestToLink(t *testing.T) {
	createdAt := time.Unix(0x420, 0).UTC()

	type Expected struct {
		Link         *Link
		ErrorMessage string
	}

	tests := []struct {
		Name     string
		Proto    *golinks.Link
		Expected Expected
	}{
		{
			Name: "single match",
			Proto: &golinks.Link{
				Prefix: "gh",
				Matches: []*golinks.Match{
					{Pattern: "issues/(.*)", Url: "https://github.com/example/repo/issues/$1"},
				},
				CreatedAt: timestamppb.New(createdAt),
			},
			Expected: Expected{Link: &Link{
				Prefix: "gh",
				Matches: []Match{{
					Pattern: regexp.MustCompile("issues/(.*)"),
					URL:     "https://github.com/example/repo/issues/$1",
				}},
				Time: createdAt,
			}},
		},
		{
			Name: "multiple matches",
			Proto: &golinks.Link{
				Prefix: "go",
				Matches: []*golinks.Match{
					{Pattern: "docs/(.*)", Url: "https://pkg.go.dev/$1"},
					{Pattern: "issues/(.*)", Url: "https://github.com/kellegous/golinks/issues/$1"},
				},
				CreatedAt: timestamppb.New(createdAt),
			},
			Expected: Expected{Link: &Link{
				Prefix: "go",
				Matches: []Match{
					{Pattern: regexp.MustCompile("docs/(.*)"), URL: "https://pkg.go.dev/$1"},
					{Pattern: regexp.MustCompile("issues/(.*)"), URL: "https://github.com/kellegous/golinks/issues/$1"},
				},
				Time: createdAt,
			}},
		},
		{
			Name:  "missing matches",
			Proto: &golinks.Link{Prefix: "gh"},
			Expected: Expected{
				ErrorMessage: "at least one match is required",
			},
		},
		{
			Name: "missing prefix",
			Proto: &golinks.Link{
				Matches: []*golinks.Match{{Pattern: "issues/(.*)", Url: "https://example.com/$1"}},
			},
			Expected: Expected{
				ErrorMessage: "prefix is required",
			},
		},
		{
			Name: "invalid pattern",
			Proto: &golinks.Link{
				Prefix:  "gh",
				Matches: []*golinks.Match{{Pattern: "["}},
			},
			Expected: Expected{
				ErrorMessage: "error parsing regexp: missing closing ]: `[`",
			},
		},
		{
			Name: "invalid URL",
			Proto: &golinks.Link{
				Prefix:  "gh",
				Matches: []*golinks.Match{{Pattern: "issues/(.*)", Url: "%invalid"}},
			},
			Expected: Expected{
				ErrorMessage: "invalid URL",
			},
		},
		{
			Name: "URL with unsupported scheme",
			Proto: &golinks.Link{
				Prefix:  "gh",
				Matches: []*golinks.Match{{Pattern: "issues/(.*)", Url: "ftp://example.com/issues/$1"}},
			},
			Expected: Expected{
				ErrorMessage: "URL must be http or https",
			},
		},
		{
			Name: "URL with host replacement",
			Proto: &golinks.Link{
				Prefix:  "gh",
				Matches: []*golinks.Match{{Pattern: "issues/(.*)", Url: "https://$1.example.com/issues/$2"}},
			},
			Expected: Expected{
				ErrorMessage: "URL host cannot contain '$' replacements",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := ToLink(tt.Proto)
			if tt.Expected.ErrorMessage != "" {
				if err == nil {
					t.Fatal("ToLink() error = nil, want an error")
				}
				if got := err.Error(); got != tt.Expected.ErrorMessage {
					t.Errorf("ToLink() error = %q, want %q", got, tt.Expected.ErrorMessage)
				}
				if got != nil {
					t.Fatalf("ToLink() = %#v, want nil", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToLink() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.Expected.Link) {
				t.Errorf("ToLink() = %#v, want %#v", got, tt.Expected.Link)
			}
		})
	}
}

func TestMatch_Expand(t *testing.T) {
	type Expected struct {
		URL string
		OK  bool
	}

	tests := []struct {
		Name     string
		Match    *Match
		URI      string
		Expected Expected
	}{
		{
			Name: "match with one replacement",
			Match: &Match{
				Pattern: regexp.MustCompile(`issues/(.*)`),
				URL:     "https://github.com/example/repo/issues/$1",
			},
			URI: "issues/42",
			Expected: Expected{
				URL: "https://github.com/example/repo/issues/42",
				OK:  true,
			},
		},
		{
			Name: "match with multiple replacements",
			Match: &Match{
				Pattern: regexp.MustCompile(`org/(.*)/repo/(.*)`),
				URL:     "https://$1.example.com/$2",
			},
			URI: "org/kellegous/repo/golinks",
			Expected: Expected{
				URL: "https://kellegous.example.com/golinks",
				OK:  true,
			},
		},
		{
			Name: "no match",
			Match: &Match{
				Pattern: regexp.MustCompile(`issues/(.*)`),
				URL:     "https://github.com/example/repo/issues/$1",
			},
			URI:      "pulls/42",
			Expected: Expected{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			gotURL, gotOK := tt.Match.Expand(tt.URI)
			if gotURL != tt.Expected.URL || gotOK != tt.Expected.OK {
				t.Errorf("Match.Expand() = (%q, %v), want (%q, %v)", gotURL, gotOK, tt.Expected.URL, tt.Expected.OK)
			}
		})
	}
}

func TestLink_Expand(t *testing.T) {
	link := &Link{
		Prefix: "gh",
		Matches: []Match{
			{
				Pattern: regexp.MustCompile(`issues/([0-9]+)`),
				URL:     "https://github.com/example/repo/issues/$1",
			},
			{
				Pattern: regexp.MustCompile(`search/(.*)`),
				URL:     "https://github.com/example/repo/search?q=$1",
			},
		},
	}

	tests := []struct {
		Name     string
		URI      string
		Expected *ExpandedURL
	}{
		{
			Name: "first match",
			URI:  "gh/issues/42",
			Expected: &ExpandedURL{
				URL:        "https://github.com/example/repo/issues/42",
				MatchIndex: 0,
			},
		},
		{
			Name: "second match",
			URI:  "gh/search/golinks",
			Expected: &ExpandedURL{
				URL:        "https://github.com/example/repo/search?q=golinks",
				MatchIndex: 1,
			},
		},
		{
			Name: "leading slashes are ignored",
			URI:  "gh///issues/42",
			Expected: &ExpandedURL{
				URL:        "https://github.com/example/repo/issues/42",
				MatchIndex: 0,
			},
		},
		{
			Name: "prefix does not match",
			URI:  "gl/issues/42",
		},
		{
			Name: "no match",
			URI:  "gh/pulls/42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := link.Expand(tt.URI); !reflect.DeepEqual(got, tt.Expected) {
				t.Errorf("Link.Expand(%q) = %#v, want %#v", tt.URI, got, tt.Expected)
			}
		})
	}
}
