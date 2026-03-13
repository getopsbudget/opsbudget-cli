package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/getopsbudget/opsbudget-cli/internal/auth"
	"github.com/spf13/pflag"
)

// resetFlags clears the "Changed" state on all flags in the command tree
// so that cobra's required-flag checks work correctly across tests.
func resetFlags() {
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) { f.Changed = false })
	for _, c := range rootCmd.Commands() {
		c.Flags().VisitAll(func(f *pflag.Flag) { f.Changed = false })
		for _, sc := range c.Commands() {
			sc.Flags().VisitAll(func(f *pflag.Flag) { f.Changed = false })
		}
	}
}

// executeCommand runs the root command with the given args and captures output.
func executeCommand(args ...string) (stdout string, stderr string, err error) {
	resetFlags()
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs(args)

	err = rootCmd.Execute()

	return outBuf.String(), errBuf.String(), err
}

func TestUnknownCommand(t *testing.T) {
	_, _, err := executeCommand("notacommand")
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("expected 'unknown command' in error, got: %s", err.Error())
	}
}

func TestUnknownHelpTopic(t *testing.T) {
	_, _, err := executeCommand("help", "notacommand")
	if err == nil {
		t.Fatal("expected error for unknown help topic")
	}
	if !strings.Contains(err.Error(), "unknown help topic") {
		t.Fatalf("expected 'unknown help topic' in error, got: %s", err.Error())
	}
}

func TestValidHelpTopic(t *testing.T) {
	_, _, err := executeCommand("help", "login")
	if err != nil {
		t.Fatalf("expected no error for valid help topic, got: %s", err)
	}
}

func TestAPIKeyCreateMocked(t *testing.T) {
	// Set up mock server.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/api-keys" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token-123" {
			t.Errorf("unexpected auth header: %s", authHeader)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var body struct {
			Name          string `json:"name"`
			ExpiresInDays *int   `json:"expires_in_days"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Name != "ci-key" {
			t.Errorf("unexpected name: %s", body.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp := map[string]any{
			"data": map[string]any{
				"id":         "key_abc123",
				"name":       body.Name,
				"key":        "ob_live_xxxxxxxxxxxxxxxx",
				"prefix":     "ob_live_xx",
				"created_at": "2026-03-13T00:00:00Z",
				"expires_at": nil,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	// Point API client at mock server.
	t.Setenv("OPSBUDGET_API_URL", srv.URL+"/v1")

	// Write temp credentials.
	tmpDir := t.TempDir()
	obDir := tmpDir + "/opsbudget"
	os.MkdirAll(obDir, 0o700)
	cred := `{"token":"test-token-123","saved_at":"2026-01-01T00:00:00Z"}`
	os.WriteFile(obDir+"/credentials.json", []byte(cred), 0o600)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Verify auth setup.
	token, err := auth.LoadToken()
	if err != nil {
		t.Fatalf("auth.LoadToken failed: %s", err)
	}
	if token != "test-token-123" {
		t.Fatalf("expected test-token-123, got %s", token)
	}

	stdout, _, err := executeCommand("api-key", "create", "--name", "ci-key")
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}
	if !strings.Contains(stdout, "ob_live_xxxxxxxxxxxxxxxx") {
		t.Fatalf("expected key in output, got: %s", stdout)
	}
	if !strings.Contains(stdout, "OPSBUDGET_API_KEY") {
		t.Fatalf("expected export snippet in output, got: %s", stdout)
	}
}

func TestAPIKeyCreateRequiresName(t *testing.T) {
	// Provide valid credentials so the test reaches flag validation.
	tmpDir := t.TempDir()
	obDir := tmpDir + "/opsbudget"
	os.MkdirAll(obDir, 0o700)
	cred := `{"token":"test-token","saved_at":"2026-01-01T00:00:00Z"}`
	os.WriteFile(obDir+"/credentials.json", []byte(cred), 0o600)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	_, _, err := executeCommand("api-key", "create")
	if err == nil {
		t.Fatal("expected error when --name is missing")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Fatalf("expected 'name' in error, got: %s", err.Error())
	}
}

func TestUnauthPing(t *testing.T) {
	// Use an empty config dir so no credentials are found.
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	_, _, err := executeCommand("ping", "list")
	if err == nil {
		t.Fatal("expected error for unauthenticated ping")
	}
	if !strings.Contains(err.Error(), "not logged in") {
		t.Fatalf("expected 'not logged in' in error, got: %s", err.Error())
	}
}
