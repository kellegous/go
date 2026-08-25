package web

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"

	"github.com/kellegous/golinks"
	"github.com/kellegous/golinks/golinks_connect"
	"github.com/kellegous/golinks/internal"
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

		mux.Handle("/", getDefaultHandler(store, assets))

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
	if _, err := s.store.Put(ctx, link); err != nil {
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

	link, err := s.store.Delete(ctx, prefix)
	if errors.Is(err, store.ErrLinkNotfound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	} else if err != nil {
		return nil, err
	}

	return connect.NewResponse(&golinks.DeleteRes{Link: link}), nil
}

func (s *server) Expand(ctx context.Context, req *connect.Request[golinks.ExpandReq]) (*connect.Response[golinks.ExpandRes], error) {
	link := req.Msg.GetLink()
	if link == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("link is required"))
	}

	uri := req.Msg.GetUri()
	if uri == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("uri is required"))
	}

	l, err := internal.ToLink(link)
	if err != nil {
		return nil, err
	}

	if eu := l.Expand(uri); eu != nil {
		return connect.NewResponse(&golinks.ExpandRes{
			Url:        eu.URL,
			MatchIndex: uint32(eu.MatchIndex),
		}), nil
	}

	return nil, connect.NewError(connect.CodeNotFound, errors.New("uri does not match"))
}
