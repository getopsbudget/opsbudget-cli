package apikey

import (
	"fmt"
	"os"

	"github.com/getopsbudget/opsbudget-cli/internal/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := requireAuth()
		if err != nil {
			return err
		}

		keys, err := client.ListAPIKeys()
		if err != nil {
			if handleAPIError(err) {
				return err
			}
			fmt.Fprintf(os.Stderr, "Error listing API keys: %s\n", err)
			return err
		}

		if len(keys) == 0 {
			fmt.Println("No API keys yet. Create one:")
			fmt.Println("  opsbudget api-key create --name \"my key\"")
			return nil
		}

		if cmd.Root().Flag("json").Changed {
			return output.PrintJSON(keys)
		}

		headers := []string{"ID", "NAME", "PREFIX", "CREATED", "EXPIRES"}
		rows := make([][]string, len(keys))
		for i, k := range keys {
			expires := "never"
			if k.ExpiresAt != nil {
				expires = *k.ExpiresAt
			}
			rows[i] = []string{k.ID, k.Name, k.Prefix, k.CreatedAt, expires}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}
