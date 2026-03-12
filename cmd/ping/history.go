package ping

import (
	"fmt"
	"os"
	"strconv"

	"github.com/getopsbudget/opsbudget-cli/internal/output"
	"github.com/spf13/cobra"
)

var historyLimit int

var historyCmd = &cobra.Command{
	Use:   "history <id|url>",
	Short: "Show recent check history for a monitor",
	Long:  "Show recent uptime check results for a monitor.\n\nExample:\n  opsbudget ping history https://example.com\n  opsbudget ping history abc123 --limit 50",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := requireAuth()
		if err != nil {
			return err
		}

		id, _, err := resolveMonitorID(client, args[0])
		if err != nil {
			if handleAPIError(err) {
				return err
			}
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return err
		}

		records, err := client.GetMonitorHistory(id, historyLimit)
		if err != nil {
			if handleAPIError(err) {
				return err
			}
			fmt.Fprintf(os.Stderr, "Error fetching history: %s\n", err)
			return err
		}

		if len(records) == 0 {
			fmt.Println("No check history yet.")
			return nil
		}

		if cmd.Root().Flag("json").Changed {
			return output.PrintJSON(records)
		}

		headers := []string{"TIME", "STATUS", "RESPONSE TIME", "STATUS CODE"}
		rows := make([][]string, len(records))
		for i, r := range records {
			rows[i] = []string{
				r.Timestamp,
				output.ColorStatus(r.Status),
				strconv.Itoa(r.ResponseTime) + "ms",
				strconv.Itoa(r.StatusCode),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	historyCmd.Flags().IntVar(&historyLimit, "limit", 20, "number of history records to show")
}
