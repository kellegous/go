package store

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var ErrNoMatch = errors.New("no match")

type Route struct {
	Pattern *regexp.Regexp `json:"pattern"`
	URL     string         `json:"url"`
	Time    time.Time      `json:"time"`
}

func (r *Route) Expand(uri string) (string, bool) {
	p := r.Pattern
	idx := p.FindStringSubmatchIndex(uri)
	if idx == nil {
		return "", false
	}
	return string(p.ExpandString(nil, r.URL, uri, idx)), true
}

func (r *Route) Prefix() string {
	p, _ := r.Pattern.LiteralPrefix()
	ix := strings.IndexAny(p, "/?")
	if ix == -1 {
		return p
	}
	return p[:ix]
}
