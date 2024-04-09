package cmd

import (
	"github.com/spf13/cobra"
)

// licenseCmd represents the license command
func newLicenseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Interact with licensing an FME Server",
		Long: `Request a license file, refresh the license, check the status of the license, generate a license request file, and get the system code for licensing.
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
