package apikey

import (
	"errors"
	"fmt"
	"os"

	"github.com/getopsbudget/opsbudget-cli/internal/api"
	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/cobra"
)

// Cmd is the `opsbudget api-key` parent command.
var Cmd = &cobra.Command{
	Use:   "api-key",
	Short: "Manage API keys",
	Long:  "Create, list, and revoke API keys for programmatic access.",
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(revokeCmd)
}

// requireAuth resolves the API key (env var or file) and returns an API client.
func requireAuth() (*api.Client, error) {
	key, err := auth.ResolveAPIKey()
	if err != nil {
		return nil, fmt.Errorf("loading credentials: %w", err)
	}
	if key == "" {
		auth.PrintLoginRequired()
		return nil, fmt.Errorf("not logged in")
	}
	return api.NewClient(key), nil
}

// handleAPIError prints user-friendly messages for common API errors.
func handleAPIError(err error) bool {
	var apiErr *api.APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	switch apiErr.StatusCode {
	case 401:
		fmt.Fprintln(os.Stderr, "⚡ Session expired")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  Log in again:")
		fmt.Fprintln(os.Stderr, "    opsbudget login")
	case 403:
		fmt.Fprintln(os.Stderr, "⚠ Plan inactive or access denied.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  Check your plan: https://opsbudget.com/billing?ref=cli")
	case 409:
		fmt.Fprintln(os.Stderr, "⚠ API key limit reached (max 25).")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  Revoke unused keys first:")
		fmt.Fprintln(os.Stderr, "    opsbudget api-key list")
	case 422:
		fmt.Fprintf(os.Stderr, "Validation error: %s\n", apiErr.Message)
	default:
		return false
	}

	return true
}
