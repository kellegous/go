package web

import "net/url"

type Options struct {
	assetProxyBaseURL *url.URL
	addr              string
	host              string
	admin             bool
}

type Option func(*Options) error

func WithAssetProxyAt(urlStr string) Option {
	return func(o *Options) error {
		u, err := url.Parse(urlStr)
		if err != nil {
			return err
		}
		o.assetProxyBaseURL = u
		return nil
	}
}

func WithAddr(addr string) Option {
	return func(o *Options) error {
		o.addr = addr
		return nil
	}
}

func WithHost(host string) Option {
	return func(o *Options) error {
		o.host = host
		return nil
	}
}

func WithAdmin(enabled bool) Option {
	return func(o *Options) error {
		o.admin = enabled
		return nil
	}
}
