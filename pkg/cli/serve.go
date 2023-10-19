package cli

import (
	"context"

	"github.com/kellegous/golinks/pkg/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cmdServe() *cobra.Command {
	var flags serveFlags
	cmd := &cobra.Command{
		Use:          "serve",
		Short:        "serve the go-links server",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			be, err := flags.Backend(context.Background())
			if err != nil {
				return err
			}
			defer be.Close()

			return web.ListenAndServe(
				be,
				web.WithAddr(flags.Addr()),
				web.WithHost(flags.Host()))
		},
	}

	flags.Register(cmd.Flags())

	return cmd
}
