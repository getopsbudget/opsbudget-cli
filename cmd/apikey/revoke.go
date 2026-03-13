package apikey

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var revokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Revoke an API key",
	Long:  "Permanently revoke an API key by ID.\n\nExample:\n  opsbudget api-key revoke abc123",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := requireAuth()
		if err != nil {
			return err
		}

		if err := client.RevokeAPIKey(args[0]); err != nil {
			if handleAPIError(err) {
				return err
			}
			fmt.Fprintf(os.Stderr, "Error revoking API key: %s\n", err)
			return err
		}

		fmt.Printf("✓ API key revoked: %s\n", args[0])
		return nil
	},
}
