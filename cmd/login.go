package cmd

import (
	"fmt"
	"os"

	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to OpsBudget via browser",
	Long:  "Opens your browser to sign up or log in to OpsBudget.\nYour credentials are saved locally for future CLI use.",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.Login()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Login failed: %s\n\n", err)
			fmt.Fprintln(os.Stderr, "  Try signing up at: https://opsbudget.com/signup?ref=cli")
			fmt.Fprintln(os.Stderr, "  Docs: https://opsbudget.com/docs")
			return err
		}

		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("saving credentials: %w", err)
		}

		fmt.Println("✓ Logged in successfully")
		return nil
	},
}
