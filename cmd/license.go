/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// licenseCmd represents the license command
var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Interact with licensing an FME Server",
	Long: `Contains several subcommands for licensing tasks related to FME Server.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("no sub-command specified")
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
}
