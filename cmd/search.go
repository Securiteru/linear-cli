package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/securiter/linear-cli/api"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Full-text search issues",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		limit := searchLimit
		if limit <= 0 {
			limit = 20
		}

		q := fmt.Sprintf(`query { searchIssues(term: "%s", first: %d) { nodes { id identifier title state { name } assignee { name } priority } } }`, escapeGraphQL(query), limit)

		var result struct {
			SearchIssues struct {
				Nodes []struct {
					ID         string `json:"id"`
					Identifier string `json:"identifier"`
					Title      string `json:"title"`
					State      *struct {
						Name string `json:"name"`
					} `json:"state"`
					Assignee *struct {
						Name string `json:"name"`
					} `json:"assignee"`
					Priority int `json:"priority"`
				} `json:"nodes"`
			} `json:"searchIssues"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		nodes := result.SearchIssues.Nodes
		if len(nodes) == 0 {
			if effectiveFormat() == "json" {
				return writeJSON([]any{})
			}
			fmt.Println("No results found.")
			return nil
		}

		return outputListItems(toAnySlice(nodes), func(item any) string {
			if n, ok := item.(struct {
				Identifier string `json:"identifier"`
				Title      string `json:"title"`
			}); ok {
				return n.Identifier + "\t" + n.Title
			}
			return ""
		}, []string{"identifier", "title", "state.name", "assignee.name", "priority"}, func() {
			w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tASSIGNEE\tPRIORITY")
			for _, i := range nodes {
				state := "-"
				if i.State != nil {
					state = i.State.Name
				}
				assignee := "-"
				if i.Assignee != nil {
					assignee = i.Assignee.Name
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", i.Identifier, i.Title, state, assignee, priorityLabel(i.Priority))
			}
			w.Flush()
		})
	},
}

var searchLimit int

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 20, "max results")
	searchCmd.Flags().StringVar(&optFields, "fields", "", "comma-separated fields (e.g. identifier,title,state.name)")
	rootCmd.AddCommand(searchCmd)
}
