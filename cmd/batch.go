package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Securiteru/linear-cli/api"
)

var batchCmd = &cobra.Command{
	Use:   "batch-create",
	Short: "Batch create issues from stdin (JSON lines)",
	Long: `Read JSON lines from stdin, each with title (required), description, and priority.
Team key is required via --team.

Example:
  echo '{"title":"Issue 1"}\n{"title":"Issue 2","priority":2}' | linear batch-create --team ADI`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if batchTeamKey == "" {
			return fmt.Errorf("--team is required")
		}

		teamID, err := resolveTeamID(batchTeamKey)
		if err != nil {
			return err
		}

		inputs := []map[string]interface{}{}
		decoder := json.NewDecoder(cmd.InOrStdin())
		for decoder.More() {
			var item map[string]interface{}
			if err := decoder.Decode(&item); err != nil {
				return fmt.Errorf("parse JSON line: %w", err)
			}
			inputs = append(inputs, item)
		}

		if len(inputs) == 0 {
			return fmt.Errorf("no input (pipe JSON lines)")
		}

		issues := []map[string]interface{}{}
		for _, item := range inputs {
			issue := map[string]interface{}{
				"title":  item["title"],
				"teamId": teamID,
			}
			if desc, ok := item["description"]; ok {
				issue["description"] = desc
			}
			if prio, ok := item["priority"]; ok {
				issue["priority"] = prio
			}
			issues = append(issues, issue)
		}

		issuesJSON, err := json.Marshal(issues)
		if err != nil {
			return err
		}

		q := fmt.Sprintf(`mutation { issueBatchCreate(input: { issues: %s }) { issues { id identifier title } } }`, string(issuesJSON))

		var result struct {
			IssueBatchCreate struct {
				Issues []struct {
					ID         string `json:"id"`
					Identifier string `json:"identifier"`
					Title      string `json:"title"`
				} `json:"issues"`
			} `json:"issueBatchCreate"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		created := result.IssueBatchCreate.Issues

		switch effectiveFormat() {
		case "json":
			return writeJSON(created)
		case "id-only":
			for _, i := range created {
				fmt.Println(i.Identifier)
			}
			return nil
		}
		if optQuiet {
			for _, i := range created {
				fmt.Printf("%s\t%s\n", i.Identifier, i.Title)
			}
			return nil
		}

		fmt.Printf("Created %d issues:\n", len(created))
		for _, i := range created {
			fmt.Printf("  %s - %s\n", i.Identifier, i.Title)
		}
		return nil
	},
}

var batchTeamKey string

func init() {
	batchCmd.Flags().StringVarP(&batchTeamKey, "team", "t", "", "team key (required)")
	rootCmd.AddCommand(batchCmd)
}
