package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/getopsbudget/opsbudget-cli/cmd/apikey"
	"github.com/getopsbudget/opsbudget-cli/cmd/ping"
	"github.com/spf13/cobra"
)

var (
	version  = "dev"
	jsonFlag bool
)

// SetVersion sets the CLI version string (injected via ldflags).
func SetVersion(v string) {
	version = v
}

// JSONOutput returns whether JSON output was requested.
func JSONOutput() bool {
	return jsonFlag
}

var rootCmd = &cobra.Command{
	Use:     "opsbudget",
	Short:   "CLI for Ping by OpsBudget — uptime monitoring",
	Long:    "Manage your uptime monitors from the terminal.\n\nDocs: https://opsbudget.com/docs",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "output as JSON")
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(ping.Cmd)
	rootCmd.AddCommand(apikey.Cmd)

	// Override the default help command so unknown topics produce an error.
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Check whether the topic resolves to a known command.
				target, _, err := rootCmd.Find(args)
				if err != nil || target == nil || target == rootCmd {
					return fmt.Errorf("unknown help topic %q. Run 'opsbudget --help'.", strings.Join(args, " "))
				}
			}
			// Delegate to default help for valid commands.
			defaultHelp(cmd, args)
			return nil
		},
	})
}

// Execute runs the root command.
func Execute() error {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		msg := err.Error()
		// Surface Cobra-level errors that SilenceErrors would otherwise hide.
		if strings.Contains(msg, "unknown command") || strings.Contains(msg, "unknown flag") || strings.Contains(msg, "unknown help topic") {
			fmt.Fprintln(os.Stderr, "Error: "+msg)
		}
	}
	return err
}
