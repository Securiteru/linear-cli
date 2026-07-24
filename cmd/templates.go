package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/Securiteru/linear-cli/api"
	"github.com/spf13/cobra"
)

var templatesTeamFilter string

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List issue templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		q := `query { templates { id name description type team { key } updatedAt } }`

		var result struct {
			Templates []struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Type        string `json:"type"`
				Team        *struct {
					Key string `json:"key"`
				} `json:"team"`
				UpdatedAt string `json:"updatedAt"`
			} `json:"templates"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		nodes := result.Templates
		if cmd.Flags().Changed("team") && templatesTeamFilter != "" {
			filtered := nodes[:0]
			for _, n := range nodes {
				if n.Team != nil && n.Team.Key == templatesTeamFilter {
					filtered = append(filtered, n)
				}
			}
			nodes = filtered
		}

		if len(nodes) == 0 {
			if effectiveFormat() == "json" {
				return writeJSON([]any{})
			}
			fmt.Println("No templates found.")
			return nil
		}

		return outputListItems(toAnySlice(nodes), func(item any) string {
			if n, ok := item.(struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			}); ok {
				return n.Name + "\t" + n.ID
			}
			return ""
		}, []string{"name", "id", "type", "team"}, func() {
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tID\tTYPE\tTEAM\tUPDATED\tDESCRIPTION")
			for _, t := range nodes {
				team := "-"
				if t.Team != nil {
					team = t.Team.Key
				}
				desc := t.Description
				if len(desc) > 60 {
					desc = desc[:60] + "…"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", t.Name, t.ID, t.Type, team, strOr(t.UpdatedAt, "")[:10], desc)
			}
			w.Flush()
		})
	},
}

func strOr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func init() {
	rootCmd.AddCommand(templatesCmd)
	templatesCmd.Flags().StringVarP(&templatesTeamFilter, "team", "t", "", "filter by team key")
}
