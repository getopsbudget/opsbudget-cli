package cmd

import (
	"fmt"
	"os"

	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of OpsBudget",
	Long:  "Removes locally stored credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {
		credPath, _ := auth.CredentialsPath()

		if err := auth.ClearCredentials(); err != nil {
			return fmt.Errorf("clearing credentials: %w", err)
		}
		fmt.Printf("Logged out. API key removed from %s\n", credPath)

		if os.Getenv("OPSBUDGET_API_KEY") != "" {
			fmt.Println("Note: OPSBUDGET_API_KEY environment variable is still set.")
		}
		return nil
	},
}
