package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdRoot = func() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "golinks",
		Short: `A "go" short-link service`,
		Long:  `A "go" short-link service`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}

	cmd.AddCommand(cmdVersion)

	return cmd
}()

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
