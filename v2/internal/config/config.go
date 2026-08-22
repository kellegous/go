package config

import (
	"fmt"
	"iter"
	"strings"
)

type Option struct {
	Key string
	Val string
}

func Opt(key, val string) *Option {
	return &Option{Key: key, Val: val}
}

func Parse(opts string) iter.Seq2[*Option, error] {
	return func(yield func(*Option, error) bool) {
		opts := strings.TrimSpace(opts)
		if opts == "" {
			return
		}

		for part := range strings.SplitSeq(opts, ",") {
			part = strings.TrimSpace(part)
			key, val, ok := strings.Cut(part, "=")
			if !ok {
				if !yield(nil, fmt.Errorf("invalid option %q: missing '='", part)) {
					return
				}
				continue
			}
			if key == "" {
				if !yield(nil, fmt.Errorf("invalid option %q: key cannot be empty", part)) {
					return
				}
				continue
			}
			if !yield(&Option{Key: key, Val: val}, nil) {
				return
			}
		}
	}
}

func WithDefaults(
	options iter.Seq2[*Option, error],
	defaults ...*Option,
) iter.Seq2[*Option, error] {
	return func(yield func(*Option, error) bool) {
		for _, def := range defaults {
			if !yield(def, nil) {
				return
			}
		}

		for option, err := range options {
			if err != nil {
				if !yield(option, err) {
					return
				}
			} else if !yield(option, nil) {
				return
			}
		}
	}
}
