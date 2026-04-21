package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/securiter/linear-cli/api"
)

var (
	searchQuery    string
	teamFilter     string
	statusFilter   string
	assigneeFilter string
	issueLimit     int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List/search issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit := issueLimit
		if limit <= 0 {
			limit = 20
		}

		filterParts := []string{}
		if teamFilter != "" {
			filterParts = append(filterParts, fmt.Sprintf(`team: { key: { eq: "%s" } }`, teamFilter))
		}
		if statusFilter != "" {
			filterParts = append(filterParts, fmt.Sprintf(`state: { name: { eq: "%s" } }`, statusFilter))
		}
		if assigneeFilter != "" {
			filterParts = append(filterParts, fmt.Sprintf(`assignee: { name: { eq: "%s" } }`, assigneeFilter))
		}

		filter := ""
		if len(filterParts) > 0 {
			filter = fmt.Sprintf("filter: { %s }", strings.Join(filterParts, ", "))
		}

		search := ""
		if searchQuery != "" {
			search = fmt.Sprintf(`search: "%s"`, searchQuery)
		}

		combined := strings.TrimSpace(strings.Join([]string{search, filter}, ", "))

		q := fmt.Sprintf(`query { issues(%s first: %d) { nodes { id identifier title state { name } assignee { name } priority labels { nodes { name } } } } }`, combined, limit)

		var result struct {
			Issues struct {
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
					Labels   *struct {
						Nodes []struct {
							Name string `json:"name"`
						} `json:"nodes"`
					} `json:"labels"`
				} `json:"nodes"`
			} `json:"issues"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		nodes := result.Issues.Nodes
		if len(nodes) == 0 {
			if effectiveFormat() == "json" {
				return writeJSON([]any{})
			}
			fmt.Println("No issues found.")
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

func toAnySlice[T any](s []T) []any {
	r := make([]any, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

func priorityLabel(p int) string {
	switch p {
	case 0:
		return "-"
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	default:
		return fmt.Sprintf("%d", p)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&searchQuery, "search", "s", "", "search text")
	listCmd.Flags().StringVarP(&teamFilter, "team", "t", "", "filter by team key (e.g. ADI)")
	listCmd.Flags().StringVarP(&statusFilter, "status", "S", "", "filter by status name")
	listCmd.Flags().StringVarP(&assigneeFilter, "assignee", "a", "", "filter by assignee name")
	listCmd.Flags().IntVarP(&issueLimit, "limit", "n", 20, "max results")
	listCmd.Flags().StringVar(&optFields, "fields", "", "comma-separated fields (e.g. identifier,title,state.name)")
}
