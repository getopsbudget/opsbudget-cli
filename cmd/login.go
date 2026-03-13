package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/getopsbudget/opsbudget-cli/internal/api"
	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/cobra"
)

var loginAPIKeyFlag bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to OpsBudget",
	Long:  "Authenticate the CLI with an API key.\nOpens your browser to sign in and generate an API key, or use --api-key for headless/CI environments.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginAPIKeyFlag {
			fmt.Println("Paste your API key (from https://opsbudget.com/settings):")
		} else {
			fmt.Println("Opening browser to sign in...")
			fmt.Println("After signing in, copy the API key and paste it here.")
			fmt.Println()
			auth.OpenAuthPage()
		}

		key, err := auth.PromptForAPIKey("API key: ")
		if err != nil {
			return fmt.Errorf("reading API key: %w", err)
		}

		key, err = auth.ValidateKeyFormat(key)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		// Verify the key by calling the API.
		client := api.NewClient(key)
		me, err := client.GetMe()
		if err != nil {
			var apiErr *api.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == 401 {
				fmt.Fprintln(os.Stderr, "Invalid API key. Please try again.")
				return fmt.Errorf("invalid API key")
			}
			fmt.Fprintf(os.Stderr, "Error verifying API key: %s\n", err)
			return err
		}

		if err := auth.SaveAPIKey(key); err != nil {
			return fmt.Errorf("saving credentials: %w", err)
		}

		// Zero the key in memory.
		keyBytes := []byte(key)
		for i := range keyBytes {
			keyBytes[i] = 0
		}

		credPath, _ := auth.CredentialsPath()
		fmt.Printf("\nLogged in as %s\n", me.Email)
		fmt.Printf("Your API key is stored in %s\n", credPath)
		return nil
	},
}

func init() {
	loginCmd.Flags().BoolVar(&loginAPIKeyFlag, "api-key", false, "skip browser and paste an API key directly (for headless/CI)")
}
