package web

import (
	"context"
	"net/http"

	"github.com/kellegous/golinks/internal/store"
)

type Options struct {
	enableMetrics bool
}

type Option func(*Options)

func WithEnableMetrics(enable bool) Option {
	return func(o *Options) {
		o.enableMetrics = enable
	}
}

func Serve(
	ctx context.Context,
	addr string,
	s store.Store,
	opts ...Option,
) error {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	mux := NewMux()

	svr := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		_ = svr.Shutdown(context.Background())
	}()

	return svr.ListenAndServe()
}
