package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type credentials struct {
	APIKey string `json:"api_key"`
}

// CredentialsPath returns the path to the credentials file.
func CredentialsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("finding config directory: %w", err)
	}
	return filepath.Join(configDir, "opsbudget", "credentials.json"), nil
}

// SaveAPIKey writes the API key to the credentials file.
// The key byte slice is zeroed after writing for security.
func SaveAPIKey(apiKey string) error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	creds := credentials{
		APIKey: apiKey,
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing credentials: %w", err)
	}

	// Zero the marshaled data.
	for i := range data {
		data[i] = 0
	}

	return nil
}

// LoadAPIKey reads the API key from the credentials file.
func LoadAPIKey() (string, error) {
	path, err := CredentialsPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("reading credentials: %w", err)
	}

	var creds credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return "", fmt.Errorf("parsing credentials: %w", err)
	}

	return creds.APIKey, nil
}

// ResolveAPIKey returns the API key by checking, in order:
// 1. OPSBUDGET_API_KEY environment variable
// 2. ~/.config/opsbudget/credentials.json file
func ResolveAPIKey() (string, error) {
	if key := os.Getenv("OPSBUDGET_API_KEY"); key != "" {
		return key, nil
	}
	return LoadAPIKey()
}

// ValidateKeyFormat checks that a key has the expected ob_ prefix and length.
// It trims whitespace before validation and returns the trimmed key.
func ValidateKeyFormat(key string) (string, error) {
	key = strings.TrimSpace(key)
	if !strings.HasPrefix(key, "ob_") {
		return "", fmt.Errorf("invalid API key: must start with \"ob_\"")
	}
	if len(key) != 67 {
		return "", fmt.Errorf("invalid API key: expected 67 characters, got %d", len(key))
	}
	return key, nil
}

// ClearCredentials deletes the credentials file.
func ClearCredentials() error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing credentials: %w", err)
	}
	return nil
}

// PrintLoginRequired prints the standard "login required" message to stderr.
func PrintLoginRequired() {
	fmt.Fprintln(os.Stderr, "Login required.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "  Sign in to get started:")
	fmt.Fprintln(os.Stderr, "    opsbudget login")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "  Or set an API key:")
	fmt.Fprintln(os.Stderr, "    export OPSBUDGET_API_KEY=\"your-key\"")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "  New to OpsBudget? Sign up free:")
	fmt.Fprintln(os.Stderr, "    https://opsbudget.com/signup?ref=cli")
}
