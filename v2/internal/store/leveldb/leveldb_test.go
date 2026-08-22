package leveldb

import (
	"os"
	"testing"

	"github.com/kellegous/golinks/internal/store"
	"github.com/kellegous/golinks/internal/store/internal"
)

type testAdapter struct{}

func (ta *testAdapter) NewStore(t *testing.T) (store.Store, func() error) {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}

	s, err := Open(tmp)
	if err != nil {
		t.Fatal(err)
	}

	return s, func() error {
		return s.Close()
	}
}

func TestLevelDB(t *testing.T) {
	internal.Test(t, &testAdapter{})
}
