package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type MachineKey struct {
	MachineKey string `json:"machineKey"`
}

// machinekeyCmd represents the machinekey command
var machinekeyCmd = &cobra.Command{
	Use:   "machinekey",
	Short: "Retrieves machine key of the machine running FME Server.",
	Long:  `Retrieves machine key of the machine running FME Server.`,
	Args:  NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			outputType = "json"
		}
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/licensing/machinekey", "GET", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
			return errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result MachineKey
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Println(result.MachineKey)
			} else {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
			}
		}
		return nil
	},
}

func init() {
	licenseCmd.AddCommand(machinekeyCmd)
}
