package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// refreshStatusCmd represents the refreshStatus command
var refreshStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of license refresh",
	Long: `Check the status of a license refresh request.
    
Example:
fmeserver license refresh status`,
	Args: NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/licensing/refresh/status", "GET", nil)
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

		var result RefreshStatus
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Println(result.Status)
				fmt.Println(result.Message)
			} else {
				fmt.Println(string(responseData))
			}

		}
		return nil

	},
}

func init() {
	refreshCmd.AddCommand(refreshStatusCmd)
}
