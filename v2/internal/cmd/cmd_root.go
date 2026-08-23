package cmd

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/kellegous/glue/devmode"
	"github.com/kellegous/glue/fn"
	"github.com/kellegous/glue/logging"
	"github.com/kellegous/golinks/internal/ui"
	"github.com/kellegous/golinks/internal/web"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

// TODO(kellegous): This is one of many breaking changes, should I keep 8067?
const defaultAddr = ":4025"

var defaultStore = Store{
	StoreType: StoreTypeLevelDB,
}

type rootFlags struct {
	Store         Store
	EnableMetrics bool
	Addr          string
	DevMode       devmode.Flag
}

func (f *rootFlags) Register(fs *pflag.FlagSet) {
	fs.VarP(&f.Store, "store", "s", "The store to use")
	fs.BoolVar(&f.EnableMetrics, "metrics", false, "Enable prometheus metrics")
	fs.StringVar(&f.Addr, "addr", defaultAddr, "The address to listen on")
	fs.Var(&f.DevMode, "dev-mode", "Enable dev mode")
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

func getAssets(ctx context.Context, devMode *devmode.Flag) (http.Handler, error) {
	if !devMode.IsEnabled() {
		return ui.Assets()
	}

	return devmode.AssetsFromVite(ctx, devMode)
}

func runMain(ctx context.Context, flags rootFlags) (err error) {
	lg := logging.MustSetup()
	defer fn.WithAbandon(lg.Sync)

	store, err := flags.Store.Open(ctx)
	if err != nil {
		return poop.Chain(err)
	}
	defer fn.WithCare(store.Close, &err)

	assets, err := getAssets(ctx, &flags.DevMode)
	if err != nil {
		return poop.Chain(err)
	}

	go func() {
		ctx, done := context.WithTimeout(ctx, 30*time.Second)
		defer done()

		if err := flags.DevMode.ShowBannerWhenReady(
			ctx,
			os.Stdout,
			flags.Addr,
		); err != nil {
			lg.Fatal("failed to show dev mode banner", zap.Error(err))
		}
	}()

	return web.NewServer(
		flags.Addr,
		assets,
		store,
		web.WithEnableMetrics(flags.EnableMetrics),
	)(ctx)
}
