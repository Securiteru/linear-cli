package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/securiter/linear-cli/api"
	"github.com/spf13/cobra"
)

var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "List webhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		q := `query { webhooks { nodes { id url team { key name } enabled createdAt } } }`

		var result struct {
			Webhooks struct {
				Nodes []struct {
				ID        string `json:"id"`
				URL       string `json:"url"`
				Team      *struct {
					Key  string `json:"key"`
					Name string `json:"name"`
				} `json:"team"`
				Enabled   bool   `json:"enabled"`
				CreatedAt string `json:"createdAt"`
				} `json:"nodes"`
			} `json:"webhooks"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		if len(result.Webhooks.Nodes) == 0 {
			fmt.Println("No webhooks found.")
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tURL\tTEAM\tACTIVE")
		for _, wh := range result.Webhooks.Nodes {
			team := "-"
			if wh.Team != nil {
				team = wh.Team.Key
			}
			active := "no"
			if wh.Enabled {
				active = "yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", wh.ID, wh.URL, team, active)
		}
		w.Flush()
		return nil
	},
}

var (
	whCreateURL    string
	whCreateTeam   string
	whCreateSecret string
)

var whCreateCmd = &cobra.Command{
	Use:   "webhook-create",
	Short: "Create a webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if whCreateURL == "" {
			return fmt.Errorf("--url is required")
		}

		input := fmt.Sprintf(`url: "%s"`, escapeGraphQL(whCreateURL))
		if whCreateTeam != "" {
			teamID, err := resolveTeamID(whCreateTeam)
			if err != nil {
				return err
			}
			input += fmt.Sprintf(`, teamId: "%s"`, teamID)
		}

		q := fmt.Sprintf(`mutation { webhookCreate(input: { %s }) { webhook { id url secret } } }`, input)

		var result struct {
			WebhookCreate struct {
				Webhook struct {
					ID     string `json:"id"`
					URL    string `json:"url"`
					Secret string `json:"secret"`
				} `json:"webhook"`
			} `json:"webhookCreate"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		wh := result.WebhookCreate.Webhook
		fmt.Printf("Created webhook: %s\n", wh.ID)
		fmt.Printf("URL: %s\n", wh.URL)
		if wh.Secret != "" {
			fmt.Printf("Secret: %s\n", wh.Secret)
		}
		return nil
	},
}

var whDeleteCmd = &cobra.Command{
	Use:   "webhook-delete [webhook-id]",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		q := fmt.Sprintf(`mutation { webhookDelete(id: "%s") { success } }`, args[0])

		var result struct {
			WebhookDelete struct {
				Success bool `json:"success"`
			} `json:"webhookDelete"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}
		if result.WebhookDelete.Success {
			fmt.Printf("Deleted webhook: %s\n", args[0])
		}
		return nil
	},
}

func init() {
	whCreateCmd.Flags().StringVarP(&whCreateURL, "url", "u", "", "webhook URL (required)")
	whCreateCmd.Flags().StringVarP(&whCreateTeam, "team", "t", "", "team key")

	rootCmd.AddCommand(webhooksCmd)
	rootCmd.AddCommand(whCreateCmd)
	rootCmd.AddCommand(whDeleteCmd)
}
