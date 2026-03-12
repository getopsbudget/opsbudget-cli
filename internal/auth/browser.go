package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/browser"
)

const (
	callbackPort = 9876
	loginTimeout = 5 * time.Minute
	authURL      = "https://opsbudget.com/cli-auth"
)

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating state: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// Login opens the browser for authentication and waits for the callback token.
func Login() (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	tokenCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "invalid state parameter", http.StatusForbidden)
			return
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html><html><body style="font-family:system-ui;text-align:center;padding:60px">
<h2>✓ Logged in to OpsBudget</h2>
<p>You can close this window and return to your terminal.</p>
</body></html>`)
		tokenCh <- token
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", callbackPort))
	if err != nil {
		return "", fmt.Errorf("starting callback server: %w", err)
	}

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	defer server.Shutdown(context.Background())

	url := fmt.Sprintf("%s?port=%d&state=%s", authURL, callbackPort, state)
	fmt.Printf("Opening browser to log in...\n")
	fmt.Printf("If the browser doesn't open, visit:\n  %s\n\n", url)

	if err := browser.OpenURL(url); err != nil {
		fmt.Printf("Could not open browser automatically.\nPlease open this URL manually:\n  %s\n\n", url)
	}

	fmt.Println("Waiting for login...")

	select {
	case token := <-tokenCh:
		return token, nil
	case err := <-errCh:
		return "", fmt.Errorf("callback server error: %w", err)
	case <-time.After(loginTimeout):
		return "", fmt.Errorf("login timed out after %s — please try again", loginTimeout)
	}
}
