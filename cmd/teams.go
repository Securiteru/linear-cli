package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/securiter/linear-cli/api"
)

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List teams",
	RunE: func(cmd *cobra.Command, args []string) error {
		q := `query { teams { nodes { id name key } } }`
		var result struct {
			Teams struct {
				Nodes []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Key  string `json:"key"`
				} `json:"nodes"`
			} `json:"teams"`
		}
		if err := api.Query(q, &result); err != nil {
			return err
		}
		for _, t := range result.Teams.Nodes {
			fmt.Printf("%s  %s  %s\n", t.Key, t.Name, t.ID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(teamsCmd)
}
