package cmd

import (
	"context"
	"os"

	"github.com/kellegous/glue/fn"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type rootFlags struct {
	Store Store
}

func (f *rootFlags) Register(fs *pflag.FlagSet) {
	fs.VarP(&f.Store, "store", "s", "The store to use")
}

var cmdRoot = func() *cobra.Command {
	flags := rootFlags{
		Store: Store{
			StoreType: StoreTypeLevelDB,
		},
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
	store, err := flags.Store.Open(ctx)
	if err != nil {
		return poop.Chain(err)
	}
	defer fn.WithCare(store.Close, &err)

	return nil
}
