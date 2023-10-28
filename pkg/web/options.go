package web

import "net/url"

type Option func(*Server) error

func WithAssetProxyAt(urlStr string) Option {
	return func(s *Server) error {
		u, err := url.Parse(urlStr)
		if err != nil {
			return err
		}
		s.assetProxyBaseURL = u
		return nil
	}
}

func WithAddr(addr string) Option {
	return func(s *Server) error {
		s.addr = addr
		return nil
	}
}

func WithHost(host string) Option {
	return func(s *Server) error {
		s.host = host
		return nil
	}
}

func WithAdmin(enabled bool) Option {
	return func(s *Server) error {
		s.admin = enabled
		return nil
	}
}
