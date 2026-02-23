package commands

import (
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "ga-cli",
	Short: "Google Analytics 4 CLI",
	Long:  "Command-line tool for Google Analytics 4 Data API and Admin API.\nUses service account or Application Default Credentials.",
}

func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	rootCmd.AddCommand(propertiesCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(realtimeCmd)
	rootCmd.AddCommand(metadataCmd)
}
