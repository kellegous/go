package store

import (
	"fmt"
	"regexp"
	"testing"
)

func describeRoute(r *Route) string {
	return fmt.Sprintf("{Pattern: %#v, URL: %#v}", r.Pattern.String(), r.URL)
}

func TestExpand(t *testing.T) {
	tests := []struct {
		Route       Route
		URI         string
		ExpectedURL string
		ExpectedOK  bool
	}{
		{
			// matching single capture unnamed
			Route{
				Pattern: regexp.MustCompile("^/a/(.*)$"),
				URL:     "https://a.com/$1",
			},
			"/a/b/c",
			"https://a.com/b/c",
			true,
		},
		{
			// matching single capture named
			Route{
				Pattern: regexp.MustCompile("^/a/(?P<foo>.*)$"),
				URL:     "https://a.com/$foo",
			},
			"/a/b/c",
			"https://a.com/b/c",
			true,
		},
		{
			// matching multiple captures unnamed
			Route{
				Pattern: regexp.MustCompile(`^/(\d+)/(\d+)$`),
				URL:     "https://a.com/$2/$1",
			},
			"/12/13",
			"https://a.com/13/12",
			true,
		},
		{
			// matching multiple captures named
			Route{
				Pattern: regexp.MustCompile(`^/(?P<foo>\d+)/(?P<bar>\d+)$`),
				URL:     "https://a.com/$bar/$foo",
			},
			"/12/13",
			"https://a.com/13/12",
			true,
		},
		{
			// unmatched
			Route{
				Pattern: regexp.MustCompile(`^/a(\d)$`),
				URL:     "https://a.com/$1",
			},
			"/ab",
			"",
			false,
		},
		{
			// literal match
			Route{
				Pattern: regexp.MustCompile(`^literal$`),
				URL:     "https://a.com",
			},
			"literal",
			"https://a.com",
			true,
		},
		{
			// literal match with unbounded captures
			Route{
				Pattern: regexp.MustCompile(`^literal$`),
				URL:     "https://a.com/$1/$2",
			},
			"literal",
			"https://a.com//",
			true,
		},
	}

	for _, test := range tests {
		url, ok := test.Route.Expand(test.URI)
		if url != test.ExpectedURL || ok != test.ExpectedOK {
			t.Fatalf(
				"for Route = %s and uri = %#v, got (%#v, %t) expected (%#v, %t)",
				describeRoute(&test.Route),
				test.URI,
				url,
				ok,
				test.ExpectedURL,
				test.ExpectedOK)
		}
	}
}

func TestPrefix(t *testing.T) {
	tests := []struct {
		Route    Route
		Expected string
	}{
		{Route{Pattern: regexp.MustCompile(`^a/b/c$`)}, "a"},
		{Route{Pattern: regexp.MustCompile(`^a/(\d+)$`)}, "a"},
		{Route{Pattern: regexp.MustCompile(`^\d$`)}, ""},
		{Route{Pattern: regexp.MustCompile(`^a\?b=(\d+)$`)}, "a"},
	}

	for _, test := range tests {
		if p := test.Route.Prefix(); p != test.Expected {
			t.Fatalf("for Route = %s, got %#v expected %#v", describeRoute(&test.Route), p, test.Expected)
		}
	}
}
