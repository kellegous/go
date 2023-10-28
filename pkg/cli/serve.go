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

			ctx := context.Background()
			s, err := flags.Store(ctx)
			if err != nil {
				return err
			}
			defer s.Close()

			svr, err := web.NewServer(
				s,
				web.WithAddr(flags.Addr()),
				web.WithHost(flags.Host()))
			if err != nil {
				return err
			}

			return svr.ListenAndServe(ctx)
		},
	}

	flags.Register(cmd.Flags())

	return cmd
}
