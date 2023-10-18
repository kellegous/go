package cli

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cmdRoot() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "golinks",
		Short: "golinks is a simple go-link/URL shortener",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(cmdServe())
	cmd.AddCommand(cmdDevelop())

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("GOLINKS")

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, err
	}

	return cmd, nil
}

func Execute() {
	root, err := cmdRoot()
	if err != nil {
		log.Panic(err)
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
