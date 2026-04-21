package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/securiter/linear-cli/api"
	"github.com/spf13/cobra"
)

var notifsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "List notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit := notifLimit
		if limit <= 0 {
			limit = 20
		}

		q := fmt.Sprintf(`query { notifications(first: %d) { nodes { id type readAt title actor { name } createdAt } } }`, limit)

		var result struct {
			Notifications struct {
				Nodes []struct {
				ID        string `json:"id"`
				Type      string `json:"type"`
				ReadAt    *string `json:"readAt"`
				Title     string `json:"title"`
					Issue  *struct {
						Identifier string `json:"identifier"`
						Title      string `json:"title"`
					} `json:"issue"`
					Comment *struct {
						Body string `json:"body"`
					} `json:"comment"`
					Actor *struct {
						Name string `json:"name"`
					} `json:"actor"`
					CreatedAt string `json:"createdAt"`
				} `json:"nodes"`
			} `json:"notifications"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}

		if len(result.Notifications.Nodes) == 0 {
			fmt.Println("No notifications.")
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
		fmt.Fprintln(w, "STATUS\tTYPE\tACTOR\tTITLE")
		for _, n := range result.Notifications.Nodes {
			status := "unread"
			if n.ReadAt != nil {
				status = "read"
			}
			actor := "-"
			if n.Actor != nil {
				actor = n.Actor.Name
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", status, n.Type, actor, n.Title)
		}
		w.Flush()
		return nil
	},
}

var notifLimit int

var notifArchiveCmd = &cobra.Command{
	Use:   "notif-archive [notification-id]",
	Short: "Archive a notification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		q := fmt.Sprintf(`mutation { notificationArchive(id: "%s") { success } }`, args[0])

		var result struct {
			NotificationArchive struct {
				Success bool `json:"success"`
			} `json:"notificationArchive"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}
		if result.NotificationArchive.Success {
			fmt.Printf("Archived notification: %s\n", args[0])
		}
		return nil
	},
}

var notifMarkReadCmd = &cobra.Command{
	Use:   "notif-read [notification-id]",
	Short: "Mark notification (and related) as read",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		q := fmt.Sprintf(`mutation { notificationMarkReadAll(id: "%s") { success } }`, args[0])

		var result struct {
			NotificationMarkReadAll struct {
				Success bool `json:"success"`
			} `json:"notificationMarkReadAll"`
		}

		if err := api.Query(q, &result); err != nil {
			return err
		}
		if result.NotificationMarkReadAll.Success {
			fmt.Printf("Marked as read: %s\n", args[0])
		}
		return nil
	},
}

func init() {
	notifsCmd.Flags().IntVarP(&notifLimit, "limit", "n", 20, "max results")

	rootCmd.AddCommand(notifsCmd)
	rootCmd.AddCommand(notifArchiveCmd)
	rootCmd.AddCommand(notifMarkReadCmd)
}
