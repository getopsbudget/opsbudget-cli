package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/getopsbudget/opsbudget-cli/internal/api"
	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show the currently authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := auth.ResolveAPIKey()
		if err != nil {
			return fmt.Errorf("loading credentials: %w", err)
		}
		if key == "" {
			auth.PrintLoginRequired()
			return fmt.Errorf("not logged in")
		}

		client := api.NewClient(key)
		me, err := client.GetMe()
		if err != nil {
			var apiErr *api.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == 401 {
				fmt.Fprintln(os.Stderr, "Invalid or expired API key.")
				fmt.Fprintln(os.Stderr)
				fmt.Fprintln(os.Stderr, "  Log in again:")
				fmt.Fprintln(os.Stderr, "    opsbudget login")
				return fmt.Errorf("invalid API key")
			}
			return err
		}

		fmt.Printf("Logged in as %s (plan: %s)\n", me.Email, me.Plan)
		return nil
	},
}
