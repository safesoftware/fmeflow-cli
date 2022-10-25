package cmd

import (
	"github.com/spf13/cobra"
)

// licenseCmd represents the license command
var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Interact with licensing an FME Server",
	Long: `Contains several subcommands for licensing tasks related to FME Server.
`,
	Args: NoArgs,
}

func init() {
	rootCmd.AddCommand(licenseCmd)
}
