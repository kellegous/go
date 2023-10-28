package cli

import "github.com/spf13/pflag"

type serveFlags struct {
	withStore
	withHTTP
}

func (f *serveFlags) Register(fs *pflag.FlagSet) {
	f.withStore.Register(fs)
	f.withHTTP.Register(fs)
}
