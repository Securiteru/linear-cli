package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/securiter/linear-cli/api"
)

var getCmd = &cobra.Command{
	Use:   "get [issue-id]",
	Short: "Get issue details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		q := fmt.Sprintf(`query { issue(id: "%s") { id identifier title description state { name } assignee { name } priority labels { nodes { name } } team { key name } url createdAt updatedAt } }`, id)

		var result struct {
			Issue *struct {
				ID          string `json:"id"`
				Identifier  string `json:"identifier"`
				Title       string `json:"title"`
				Description string `json:"description"`
				State       *struct {
					Name string `json:"name"`
				} `json:"state"`
				Assignee *struct {
					Name string `json:"name"`
				} `json:"assignee"`
				Priority int `json:"priority"`
				Labels   *struct {
					Nodes []struct {
						Name string `json:"name"`
					} `json:"nodes"`
				} `json:"labels"`
				Team *struct {
					Key  string `json:"key"`
					Name string `json:"name"`
				} `json:"team"`
				URL       string `json:"url"`
				CreatedAt string `json:"createdAt"`
				UpdatedAt string `json:"updatedAt"`
			} `json:"issue"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}
		if result.Issue == nil {
			return fmt.Errorf("issue %q not found", id)
		}

		issue := result.Issue
		fmt.Printf("%s - %s\n", issue.Identifier, issue.Title)
		fmt.Printf("URL: %s\n", issue.URL)
		if issue.Team != nil {
			fmt.Printf("Team: %s (%s)\n", issue.Team.Name, issue.Team.Key)
		}
		state := "-"
		if issue.State != nil {
			state = issue.State.Name
		}
		fmt.Printf("Status: %s\n", state)
		assignee := "-"
		if issue.Assignee != nil {
			assignee = issue.Assignee.Name
		}
		fmt.Printf("Assignee: %s\n", assignee)
		fmt.Printf("Priority: %s\n", priorityLabel(issue.Priority))
		if issue.Labels != nil && len(issue.Labels.Nodes) > 0 {
			fmt.Print("Labels: ")
			for i, l := range issue.Labels.Nodes {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(l.Name)
			}
			fmt.Println()
		}
		fmt.Printf("Created: %s\n", issue.CreatedAt)
		fmt.Printf("Updated: %s\n", issue.UpdatedAt)
		if issue.Description != "" {
			fmt.Printf("\n%s\n", issue.Description)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
