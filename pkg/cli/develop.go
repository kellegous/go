package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"

	"github.com/kellegous/golinks/pkg/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startViteServer(
	ctx context.Context,
	root string,
	port int,
) error {
	cmd := exec.CommandContext(
		ctx,
		"npx",
		"vite",
		"--clearScreen=false",
		"--port",
		strconv.Itoa(port))
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func cmdDevelop() *cobra.Command {
	var flags developFlags

	cmd := &cobra.Command{
		Use:          "develop",
		Short:        "run the go-links development server",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			root := flags.Root()
			vitePort := flags.VitePort()

			ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
			defer done()

			if err := startViteServer(ctx, root, vitePort); err != nil {
				return err
			}

			s, err := flags.Store(ctx)
			if err != nil {
				return err
			}
			defer s.Close()

			svr, err := web.NewServer(
				s,
				web.WithAddr(flags.Addr()),
				web.WithHost(flags.Host()),
				web.WithAssetProxyAt(fmt.Sprintf("http://localhost:%d/", vitePort)))
			if err != nil {
				return err
			}

			ch := make(chan error)
			go func() {
				ch <- svr.ListenAndServe(ctx)
			}()

			// start the web
			select {
			case <-ctx.Done():
			case err := <-ch:
				return err
			}
			return nil
		},
	}

	flags.Register(cmd.Flags())

	return cmd
}
