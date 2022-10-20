package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type SystemCode struct {
	SystemCode string `json:"systemCode"`
}

// systemcodeCmd represents the systemcode command
var systemcodeCmd = &cobra.Command{
	Use:   "systemcode",
	Short: "Retrieves system code of the machine running FME Server.",
	Long:  `Retrieves system code of the machine running FME Server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/licensing/systemcode", "GET", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(&request)
		if err != nil {
			return err
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result SystemCode
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Printf(result.SystemCode)
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
	licenseCmd.AddCommand(systemcodeCmd)
}
