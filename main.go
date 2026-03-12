package main

import (
	"os"

	"github.com/getopsbudget/opsbudget-cli/cmd"
	"github.com/getopsbudget/opsbudget-cli/internal/api"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	api.Version = version
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
