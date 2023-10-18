package cli

import "github.com/spf13/pflag"

type serveFlags struct {
	withBackend
	withHTTP
}

func (f *serveFlags) Register(fs *pflag.FlagSet) {
	f.withBackend.Register(fs)
	f.withHTTP.Register(fs)
}
