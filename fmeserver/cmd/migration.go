/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
