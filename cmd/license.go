package cmd

import (
	"github.com/spf13/cobra"
)

// licenseCmd represents the license command
func newLicenseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Interact with licensing an FME Server",
		Long: `Contains several subcommands for licensing tasks related to FME Server.
	`,
		Args: NoArgs,
	}
	cmd.AddCommand(newLicenseStatusCmd())
	cmd.AddCommand(newMachineKeyCmd())
	cmd.AddCommand(newRefreshCmd())
	cmd.AddCommand(newLicenseRequestCmd())
	cmd.AddCommand(newLicenseRequestFileCmd())
	cmd.AddCommand(newSystemCodeCmd())
	return cmd
}
