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
