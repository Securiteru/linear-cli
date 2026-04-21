package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/securiter/linear-cli/api"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		q := `query { viewer { id name email } }`
		var result struct {
			Viewer struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"viewer"`
		}
		if err := api.Query(q, &result); err != nil {
			return err
		}
		fmt.Printf("%s <%s>\n", result.Viewer.Name, result.Viewer.Email)
		fmt.Printf("ID: %s\n", result.Viewer.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
