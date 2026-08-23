package web

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/kellegous/golinks"
	"github.com/kellegous/golinks/golinks_connect"
	"github.com/kellegous/golinks/internal/store"
)

const rpcPrefix = "/rpc"

type server struct {
	store         store.Store
	enableMetrics bool
}

var _ golinks_connect.GoLinksHandler = (*server)(nil)

func NewServer(
	addr string,
	assets http.Handler,
	store store.Store,
	opts ...Option,
) func(context.Context) error {
	s := &server{
		store: store,
	}

	for _, opt := range opts {
		opt(s)
	}

	return func(ctx context.Context) error {
		mux := NewMux()

		mux.Handle("/s/", assets)

		path, handler := golinks_connect.NewGoLinksHandler(s)
		mux.Handle(rpcPrefix+path, handler)

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
}

type Option func(*server)

func WithEnableMetrics(enable bool) Option {
	return func(s *server) {
		s.enableMetrics = enable
	}
}

func (s *server) Put(ctx context.Context, req *connect.Request[golinks.PutReq]) (*connect.Response[golinks.PutRes], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (s *server) Get(ctx context.Context, req *connect.Request[golinks.GetReq]) (*connect.Response[golinks.GetRes], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (s *server) Delete(ctx context.Context, req *connect.Request[golinks.DeleteReq]) (*connect.Response[golinks.DeleteRes], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
