package internal

import (
	"regexp"
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

	return &Match{
		Pattern: p,
		URL:     proto.GetUrl(),
	}, nil
}
