package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/securiter/linear-cli/api"
)

var commentBody string

var commentCmd = &cobra.Command{
	Use:   "comment [issue-id] [body]",
	Short: "Add a comment to an issue",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := commentBody
		if len(args) == 1 {
			body = args[0]
		}
		if body == "" {
			return fmt.Errorf("comment body is required (pass as arg or --body)")
		}

		id := args[0]
		if len(args) == 0 {
			return fmt.Errorf("issue-id is required")
		}

		q := fmt.Sprintf(`mutation { commentCreate(input: { issueId: "%s", body: "%s" }) { comment { id body createdAt } } }`, id, escapeGraphQL(body))

		var result struct {
			CommentCreate struct {
				Comment struct {
					ID        string `json:"id"`
					Body      string `json:"body"`
					CreatedAt string `json:"createdAt"`
				} `json:"comment"`
			} `json:"commentCreate"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		fmt.Printf("Comment added to %s at %s\n", id, result.CommentCreate.Comment.CreatedAt)
		return nil
	},
}

var listCommentsCmd = &cobra.Command{
	Use:   "comments [issue-id]",
	Short: "List comments on an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		q := fmt.Sprintf(`query { issue(id: "%s") { comments { nodes { id body user { name } createdAt updatedAt resolvedAt } } } }`, id)

		var result struct {
			Issue *struct {
				Comments struct {
					Nodes []struct {
						ID         string `json:"id"`
						Body       string `json:"body"`
						User       *struct {
							Name string `json:"name"`
						} `json:"user"`
						CreatedAt string `json:"createdAt"`
						UpdatedAt string `json:"updatedAt"`
						ResolvedAt *string `json:"resolvedAt"`
					} `json:"nodes"`
				} `json:"comments"`
			} `json:"issue"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}
		if result.Issue == nil {
			return fmt.Errorf("issue %q not found", id)
		}

		if len(result.Issue.Comments.Nodes) == 0 {
			fmt.Println("No comments.")
			return nil
		}

		for _, c := range result.Issue.Comments.Nodes {
			user := "unknown"
			if c.User != nil {
				user = c.User.Name
			}
			resolved := ""
			if c.ResolvedAt != nil {
				resolved = " [resolved]"
			}
			fmt.Printf("--- %s (%s)%s %s ---\n%s\n\n", user, c.ID, resolved, c.CreatedAt[:10], c.Body)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
	rootCmd.AddCommand(listCommentsCmd)
	commentCmd.Flags().StringVarP(&commentBody, "body", "b", "", "comment body (use for multiline)")
}
