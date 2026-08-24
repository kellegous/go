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
	link := req.Msg.GetLink()
	// TODO(kellegous): validate link
	if err := s.store.Put(ctx, link); err != nil {
		return nil, err
	}

	return connect.NewResponse(&golinks.PutRes{Link: link}), nil
}

func (s *server) Get(ctx context.Context, req *connect.Request[golinks.GetReq]) (*connect.Response[golinks.GetRes], error) {
	link, err := s.store.Get(ctx, req.Msg.GetPrefix())
	if errors.Is(err, store.ErrLinkNotfound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	} else if err != nil {
		return nil, err
	}

	return connect.NewResponse(&golinks.GetRes{Link: link}), nil
}

func (s *server) Delete(ctx context.Context, req *connect.Request[golinks.DeleteReq]) (*connect.Response[golinks.DeleteRes], error) {
	prefix := req.Msg.GetPrefix()

	// TODO(kellegous): This has atomicity issues and the returned link should be returned from the store.Delete() call.
	link, err := s.store.Get(ctx, prefix)
	if errors.Is(err, store.ErrLinkNotfound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	} else if err != nil {
		return nil, err
	}

	if err := s.store.Delete(ctx, prefix); err != nil {
		return nil, err
	}

	return connect.NewResponse(&golinks.DeleteRes{Link: link}), nil
}
