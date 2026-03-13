package ping

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/getopsbudget/opsbudget-cli/internal/api"
	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/cobra"
)

// Cmd is the `opsbudget ping` parent command.
var Cmd = &cobra.Command{
	Use:   "ping",
	Short: "Manage uptime monitors",
	Long:  "Add, list, and manage uptime monitors for your sites and APIs.",
}

func init() {
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(rmCmd)
	Cmd.AddCommand(historyCmd)
}

// requireAuth resolves the API key (env var or file) and returns an API client.
// If no key is found, it prints login instructions and returns an error.
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
// Returns true if the error was handled (printed), false otherwise.
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
	case 402:
		fmt.Fprintln(os.Stderr, "💳 Payment method required")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  Add a payment method to start your 14-day free trial:")
		fmt.Fprintln(os.Stderr, "    https://opsbudget.com/billing?ref=cli")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  You won't be charged for 14 days. Plans start at $4/mo.")
	case 403:
		fmt.Fprintln(os.Stderr, "Access denied.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  Contact support: https://opsbudget.com/support")
	default:
		return false
	}

	return true
}

// resolveMonitorID resolves an ID-or-URL argument to a monitor ID.
// If the argument looks like a URL, it fetches all monitors and finds the match.
func resolveMonitorID(client *api.Client, idOrURL string) (string, string, error) {
	if strings.Contains(idOrURL, "://") {
		// Looks like a URL — resolve via list
		monitors, err := client.ListMonitors()
		if err != nil {
			return "", "", err
		}
		for _, m := range monitors {
			if m.URL == idOrURL {
				return m.ID, m.Name, nil
			}
		}
		return "", "", fmt.Errorf("no monitor found for URL: %s", idOrURL)
	}
	return idOrURL, "", nil
}
