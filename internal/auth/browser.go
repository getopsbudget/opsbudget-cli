package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/browser"
	"golang.org/x/term"
)

const authURL = "https://opsbudget.com/cli-auth"

// OpenAuthPage opens the user's browser to the CLI authentication page.
func OpenAuthPage() {
	if err := browser.OpenURL(authURL); err != nil {
		fmt.Printf("Could not open browser automatically.\nPlease open this URL manually:\n  %s\n\n", authURL)
	}
}

// PromptForAPIKey prints the given prompt and reads a line from stdin
// with terminal echo disabled (masked input). Returns the trimmed input.
func PromptForAPIKey(prompt string) (string, error) {
	fmt.Print(prompt)

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		raw, err := term.ReadPassword(fd)
		fmt.Println() // newline after masked input
		if err != nil {
			return "", fmt.Errorf("reading input: %w", err)
		}
		key := strings.TrimSpace(string(raw))
		// Zero the raw bytes.
		for i := range raw {
			raw[i] = 0
		}
		return key, nil
	}

	// Non-terminal (piped input): read a line from stdin.
	var buf [256]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("reading input: %w", err)
	}
	key := strings.TrimSpace(string(buf[:n]))
	for i := range buf {
		buf[i] = 0
	}
	return key, nil
}
