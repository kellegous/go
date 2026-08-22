package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kellegous/glue/build"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

var cmdVersion = func() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of golinks",
		Long:  "Print the version number of golinks",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runVersion(cmd.Context()); err != nil {
				poop.HitFan(err)
			}
		},
	}
	return cmd
}()

func runVersion(_ context.Context) error {
	b, err := json.MarshalIndent(build.ReadSummary(), "", "  ")
	if err != nil {
		return poop.Chain(err)
	}
	fmt.Printf("%s\n", b)
	return nil
}
