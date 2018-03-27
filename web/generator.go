package web

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"database/sql"
	"github.com/HALtheWise/o-links/context"
	_ "github.com/lib/pq"
)

// This file is responsible for generating default URLs and validating them.

// hasCollision checks whether the proposed link already exists in the database,
// with a different ID than the one provided.
func hasCollision(ctx *context.Context, link string, uid uint32) (bool, error) {
	route, err := ctx.Get(link)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if route.Uid == uid {
		// The ID's match, so the new route is a legal update to the existing one
		return false, nil
	}

	return true, nil
}

var nouns, adjectives []string

func init() {
	nouns = strings.Split("cat dog phoenix", " ")
	adjectives = strings.Split("fuzzy large", " ")
}

var randsource *rand.Rand

func init() {
	s := rand.NewSource(time.Now().Unix())
	randsource = rand.New(s) // initialize local pseudorandom generator
}

var errCannotGenerate = errors.New("Unable to generate link")

// generateLink() creates, from scratch, a short pronounceable link that does not currently
// collide with anything in the dictionary.
func generateLink(ctx *context.Context, uid uint32) (string, error) {
	const (
		NUM_ATTEMPTS = 10
	)

	// Look for single-noun phrases that work
	for i := 0; i < NUM_ATTEMPTS; i++ {
		link := nouns[randsource.Intn(len(nouns))]
		collides, err := hasCollision(ctx, link, uid)
		if err != nil {
			return "", err
		}

		if !collides {
			// We found something that works!
			return link, nil
		}
	}

	// Look for adjective+noun phrases that work
	for i := 0; i < NUM_ATTEMPTS; i++ {
		link := fmt.Sprintf("%s%s",
			adjectives[randsource.Intn(len(adjectives))], nouns[randsource.Intn(len(nouns))])
		collides, err := hasCollision(ctx, link, uid)
		if err != nil {
			return "", err
		}

		if !collides {
			// We found something that works!
			return link, nil
		}
	}

	// Generate a number, with increasing digit count
	for _, maxval := range []int{9, 99, 9999, 999999, 9999999999} {
		for i := 0; i < NUM_ATTEMPTS; i++ {
			link := fmt.Sprintf("%d", randsource.Intn(maxval))
			collides, err := hasCollision(ctx, link, uid)
			if err != nil {
				return "", err
			}

			if !collides {
				// We found something that works!
				return link, nil
			}
		}
	}

	return "", errCannotGenerate
}
