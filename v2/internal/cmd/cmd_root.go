package cmd

import (
	"context"
	"os"

	"github.com/kellegous/glue/fn"
	"github.com/kellegous/glue/logging"
	"github.com/kellegous/golinks/internal/web"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const defaultAddr = ":8067"

var defaultStore = Store{
	StoreType: StoreTypeLevelDB,
}

type rootFlags struct {
	Store         Store
	EnableMetrics bool
	Addr          string
}

func (f *rootFlags) Register(fs *pflag.FlagSet) {
	fs.VarP(&f.Store, "store", "s", "The store to use")
	fs.BoolVar(&f.EnableMetrics, "metrics", false, "Enable prometheus metrics")
	fs.StringVar(&f.Addr, "addr", defaultAddr, "The address to listen on")
}

var cmdRoot = func() *cobra.Command {
	flags := rootFlags{
		Store: defaultStore,
	}

	var cmd = &cobra.Command{
		Use:   "golinks",
		Short: `A "go" short-link service`,
		Long:  `A "go" short-link service`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runMain(cmd.Context(), flags); err != nil {
				poop.HitFan(err)
			}
		},
	}

	cmd.AddCommand(cmdVersion)

	flags.Register(cmd.Flags())

	return cmd
}()

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

func runMain(ctx context.Context, flags rootFlags) (err error) {
	lg := logging.MustSetup()
	defer fn.WithAbandon(lg.Sync)

	store, err := flags.Store.Open(ctx)
	if err != nil {
		return poop.Chain(err)
	}
	defer fn.WithCare(store.Close, &err)

	return web.Serve(
		ctx,
		flags.Addr,
		store,
		web.WithEnableMetrics(flags.EnableMetrics),
	)
}
