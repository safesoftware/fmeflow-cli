package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// migrationCmd represents the migration command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Returns information on migration tasks using the tasks subcommand.",
	Long:  `Returns information on migration tasks using the tasks subcommand.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("no sub-command specified")
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}
