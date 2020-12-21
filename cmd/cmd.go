package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

var configFile string

// NewRootCmd creates a new instance of the root command
func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   "Tokenizer",
		Short: "Tokenizer Service",
	}

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")

	cmd.AddCommand(NewVersionCmd(ctx))
	cmd.AddCommand(NewServerStartCmd(ctx))

	return cmd
}
