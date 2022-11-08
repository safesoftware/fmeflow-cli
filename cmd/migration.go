package cmd

import (
	"github.com/spf13/cobra"
)

func newMigrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migration",
		Short: "Returns information on migration tasks using the tasks subcommand.",
		Long:  `Returns information on migration tasks using the tasks subcommand.`,
		Args:  NoArgs,
	}
	cmd.AddCommand(newMigrationTasksCmd())
	return cmd
}
