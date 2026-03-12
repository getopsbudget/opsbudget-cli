package cmd

import (
	"fmt"

	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of OpsBudget",
	Long:  "Removes locally stored credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.ClearToken(); err != nil {
			return fmt.Errorf("clearing credentials: %w", err)
		}
		fmt.Println("✓ Logged out successfully")
		return nil
	},
}
