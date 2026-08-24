package internal

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/kellegous/golinks"
	"github.com/kellegous/poop"
)

type Link struct {
	Prefix  string    `json:"prefix"`
	Matches []Match   `json:"matches"`
	Time    time.Time `json:"time"`
}

type Match struct {
	Pattern *regexp.Regexp `json:"pattern"`
	URL     string         `json:"url"`
}

type ExpandedURL struct {
	URL        string `json:"url"`
	MatchIndex int    `json:"match_index"`
}

func (l *Link) Expand(uri string) *ExpandedURL {
	if !strings.HasPrefix(uri, l.Prefix) {
		return nil
	}

	suffix := strings.TrimLeft(uri[len(l.Prefix):], "/")

	for i, match := range l.Matches {
		if expanded, ok := match.Expand(suffix); ok {
			return &ExpandedURL{
				URL:        expanded,
				MatchIndex: i,
			}
		}
	}
	return nil
}

func (m *Match) Expand(uri string) (string, bool) {
	p := m.Pattern
	idx := p.FindStringSubmatchIndex(uri)
	if idx == nil {
		return "", false
	}
	return string(p.ExpandString(nil, m.URL, uri, idx)), true
}

func ToLink(proto *golinks.Link) (*Link, error) {
	if len(proto.GetMatches()) == 0 {
		return nil, poop.New("at least one match is required")
	}

	if proto.GetPrefix() == "" {
		return nil, poop.New("prefix is required")
	}

	matches := make([]Match, 0, len(proto.Matches))
	for _, m := range proto.GetMatches() {
		match, err := toMatch(m)
		if err != nil {
			return nil, poop.Chain(err)
		}
		matches = append(matches, *match)
	}
	return &Link{
		Prefix:  proto.GetPrefix(),
		Matches: matches,
		Time:    proto.CreatedAt.AsTime(),
	}, nil
}

func toMatch(proto *golinks.Match) (*Match, error) {
	p, err := regexp.Compile(proto.GetPattern())
	if err != nil {
		return nil, err
	}

	url := proto.GetUrl()
	if err := validateURL(url); err != nil {
		return nil, err
	}

	return &Match{
		Pattern: p,
		URL:     url,
	}, nil
}

func validateURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return poop.New("invalid URL")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return poop.New("URL must be http or https")
	}

	if strings.Contains(u.Host, "$") {
		return poop.New("URL host cannot contain '$' replacements")
	}

	return nil
}
