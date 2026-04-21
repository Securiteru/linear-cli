package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linear",
	Short: "Linear GraphQL API CLI",
	Long:  "Fast CLI for Linear issue tracking via GraphQL API. Works with psst for secret injection.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if os.Getenv("LINEAR_API_KEY") == "" {
			return fmt.Errorf("LINEAR_API_KEY not set")
		}
		return nil
	}
}
