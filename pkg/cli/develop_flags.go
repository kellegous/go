package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type developFlags struct {
	withBackend
	withHTTP
}

func (f *developFlags) Root() string {
	return viper.GetString("root")
}

func (f *developFlags) VitePort() int {
	return viper.GetInt("vite.port")
}

func (f *developFlags) Register(fs *pflag.FlagSet) {
	f.withBackend.Register(fs)
	f.withHTTP.Register(fs)

	fs.String(
		"root",
		".",
		"the root directory of the source tree.")

	fs.Int(
		"vite.port",
		3000,
		"the port where the vite server will listen")
}
