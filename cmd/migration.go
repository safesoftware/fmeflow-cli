package cmd

import (
	"github.com/spf13/cobra"
)

// migrationCmd represents the migration command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Returns information on migration tasks using the tasks subcommand.",
	Long:  `Returns information on migration tasks using the tasks subcommand.`,
	Args:  NoArgs,
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}
