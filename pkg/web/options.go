package web

import "net/url"

type Options struct {
	assetProxyBaseURL *url.URL
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
