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
		"--debug",
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
		Use:   "develop",
		Short: "run the go-links development server",
		Args:  cobra.NoArgs,
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

			be, err := flags.Backend(ctx)
			if err != nil {
				return err
			}
			defer be.Close()

			ch := make(chan error)
			go func() {
				ch <- web.ListenAndServe(
					be,
					web.WithAssetProxyAt(fmt.Sprintf("http://localhost:%d/", vitePort)))
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
