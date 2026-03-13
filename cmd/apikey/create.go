package apikey

import (
	"fmt"
	"os"

	"github.com/getopsbudget/opsbudget-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	createName    string
	createExpires int
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Long:  "Create a new API key for programmatic access to the OpsBudget API.\n\nExample:\n  opsbudget api-key create --name \"CI deploy\"\n  opsbudget api-key create --name \"staging\" --expires 90",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := requireAuth()
		if err != nil {
			return err
		}

		req := &api.CreateAPIKeyRequest{
			Name: createName,
		}
		if cmd.Flags().Changed("expires") {
			req.ExpiresInDays = &createExpires
		}

		key, err := client.CreateAPIKey(req)
		if err != nil {
			if handleAPIError(err) {
				return err
			}
			fmt.Fprintf(os.Stderr, "Error creating API key: %s\n", err)
			return err
		}

		w := cmd.OutOrStdout()
		fmt.Fprintln(w, "✓ API key created")
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  Name:    %s\n", key.Name)
		fmt.Fprintf(w, "  Key:     %s\n", key.Key)
		if key.ExpiresAt != nil {
			fmt.Fprintf(w, "  Expires: %s\n", *key.ExpiresAt)
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  ⚠ Copy this key now — it won't be shown again.")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  Export as environment variable:")
		fmt.Fprintf(w, "    export OPSBUDGET_API_KEY=%s\n", key.Key)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&createName, "name", "", "name for the API key (required)")
	createCmd.Flags().IntVar(&createExpires, "expires", 0, "expiration in days (optional)")
	_ = createCmd.MarkFlagRequired("name")
}
