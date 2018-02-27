package web

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/HALtheWise/go-links/context"
)

// Regression test
func TestBlankGenerator(t *testing.T) {
	e := needEnv(t)
	defer e.destroy()

	oldnouns := nouns
	oldadjectives := adjectives

	defer func() {
		nouns = oldnouns
		adjectives = oldadjectives
	}()

	nouns = strings.Split("cat dog", " ")
	adjectives = strings.Split("small large", " ")

	randsource = rand.New(rand.NewSource(42))

	desired := []string{"dog", "cat", "smallcat", "largecat", "largedog", "smalldog", "4", "3", "0", "6", "5", "7", "8", "1", "85", "20"}
	var results []string

	for i := range desired {
		uid := rand.Uint64()
		link, err := generateLink(e.ctx, uid)
		if err != nil {
			t.Fatalf("Unable to generate link #%d, #s", i+1, err)
		}
		results = append(results, link)

		err = e.ctx.Put(link, &context.Route{Uid: uid, URL: "https://google.com"})
		if err != nil {
			t.Fatalf("Unable to put route in database: %s", err)
		}
	}

	if !reflect.DeepEqual(desired, results) {
		t.Errorf("Wrong sequence of links generated: Expected \n%#v\n got \n%#v",
			desired, results)
	}
}
